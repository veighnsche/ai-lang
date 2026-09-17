package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// Values: Kind str,int,bool,dec,rec,ok,err,seq,bytes. Field access
// works on rec/ok/err. Dec holds canonical digits; proofs compare
// exactly via big.Rat, so canonical strings compare equal exactly
// when numeric. Int holds an arbitrary-precision value (v10:
// unbounded, never wraps).
type Value struct {
	Kind    string
	S       string
	N       *big.Int
	B       bool
	D       string
	Dict    map[string]*Value
	ErrKind string
	// Bytes holds the owned octets of a bytes value (v45 S1):
	// exactly the validated literal contents, never shared.
	Bytes []byte
	// Rec names the record constructor for rec values built in value
	// positions; empty for Ok payloads (which render as Ok).
	Rec string
	// Arr holds the ordered members of a seq value (v36 S1);
	// Elem names the checked element type. Brands erase at runtime,
	// so branded members arrive here as their strings: Elem keeps
	// the static identity vEq compares.
	Arr  []*Value
	Elem string
}

// recordDecl finds a record declaration by name, first wins across
// modules, matching the checker and the emitter.
func recordDecl(prog *Program, name string) *TypeDecl {
	if prog == nil {
		return nil
	}
	for _, b := range builtinTypeDecls() {
		if b.Name == name {
			return b
		}
	}
	for _, m := range prog.Modules {
		for _, d := range m.Decls {
			if td, ok := d.(*TypeDecl); ok && td.Name == name {
				return td
			}
		}
	}
	return nil
}

// decRat parses canonical dec digits exactly. Canonical form always
// parses; the error is defense for hand-built nodes.
func decRat(d string) (*big.Rat, error) {
	r, ok := new(big.Rat).SetString(d)
	if !ok {
		return nil, fmt.Errorf("bad dec value %q", d)
	}
	return r, nil
}

// decCmp applies a comparison operator to a big.Rat Cmp result (-1, 0,
// +1): exact decimal ordering for the proof, no float anywhere.
func decCmp(op string, c int) bool {
	switch op {
	case ">=":
		return c >= 0
	case "<=":
		return c <= 0
	case ">":
		return c > 0
	case "<":
		return c < 0
	case "!=":
		return c != 0
	default:
		return c == 0
	}
}

// evArith evaluates exact arithmetic: same-kind operands only (the
// static rule fires first with a code; this is unreachable past the
// gate). Ints are unbounded (v10: no overflow, never wraps — the old
// int64 gate is gone). Dec runs on scaled integers, so 0.1+0.2 is 0.3
// and rendering is total with no rounding rule (add/sub/mul of
// terminating decimals always terminate, which is why division waits
// for its own spec).
func evArith(op string, lv, rv *Value) (*Value, error) {
	// v39 S4: Seq<T> + T appends immutably. The tail array is
	// copied, never shared, so the original value is unchanged;
	// member typing is the checker's (brands erase before values
	// arrive here, so no runtime recheck could be exact anyway).
	if op == "+" && lv.Kind == "seq" {
		arr := append(append([]*Value{}, lv.Arr...), rv)
		return &Value{Kind: "seq", Arr: arr, Elem: lv.Elem}, nil
	}
	// v16: + concatenates strings; every other string computation
	// stays a loud dynamic error past the static gate, as before.
	if lv.Kind == "str" && rv.Kind == "str" {
		if op != "+" {
			return nil, fmt.Errorf("bad %s operands", op)
		}
		return &Value{Kind: "str", S: lv.S + rv.S}, nil
	}
	if lv.Kind != rv.Kind || (lv.Kind != "int" && lv.Kind != "dec") {
		return nil, fmt.Errorf("bad %s operands", op)
	}
	if lv.Kind == "dec" {
		d, err := decArith(op, lv.D, rv.D)
		if err != nil {
			return nil, err
		}
		return &Value{Kind: "dec", D: d}, nil
	}
	// v17: / and % are exact Euclidean integer division through
	// big.Int.DivMod (verified: a == b*q + r with 0 <= r < |b| on
	// every sign combination). A zero divisor is loud, never silent.
	// Non-int operands are loud too: direct evaluator callers bypass
	// the static gate, so this must error, never panic on nil.
	if op == "/" || op == "%" {
		if lv.Kind != "int" || rv.Kind != "int" {
			return nil, fmt.Errorf("bad %s operands", op)
		}
		if rv.N.Sign() == 0 {
			return nil, fmt.Errorf("int division by zero")
		}
		q, m := new(big.Int), new(big.Int)
		q.DivMod(lv.N, rv.N, m)
		if op == "/" {
			return &Value{Kind: "int", N: q}, nil
		}
		return &Value{Kind: "int", N: m}, nil
	}
	r := new(big.Int)
	switch op {
	case "+":
		r.Add(lv.N, rv.N)
	case "-":
		r.Sub(lv.N, rv.N)
	default:
		r.Mul(lv.N, rv.N)
	}
	return &Value{Kind: "int", N: r}, nil
}

// parseDecParts splits canonical digits into mantissa and scale:
// "12.34" becomes (1234, 2). Canonical form always parses; the error
// is defense for hand-built nodes.
func parseDecParts(d string) (*big.Int, int, error) {
	neg := false
	if strings.HasPrefix(d, "-") {
		neg = true
		d = d[1:]
	}
	parts := strings.SplitN(d, ".", 2)
	if len(parts) != 2 {
		return nil, 0, fmt.Errorf("bad dec value %q", d)
	}
	mant := new(big.Int)
	if _, ok := mant.SetString(parts[0]+parts[1], 10); !ok {
		return nil, 0, fmt.Errorf("bad dec value %q", d)
	}
	if neg {
		mant.Neg(mant)
	}
	return mant, len(parts[1]), nil
}

func pow10(n int) *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), nil)
}

// decArith computes on scaled integers and re-normalizes through
// canonDec, so results are canonical by construction.
func decArith(op, l, r string) (string, error) {
	lm, ls, err := parseDecParts(l)
	if err != nil {
		return "", err
	}
	rm, rs, err := parseDecParts(r)
	if err != nil {
		return "", err
	}
	m := new(big.Int)
	s := 0
	switch op {
	case "+":
		s = ls
		if rs > s {
			s = rs
		}
		m.Add(new(big.Int).Mul(lm, pow10(s-ls)), new(big.Int).Mul(rm, pow10(s-rs)))
	case "-":
		s = ls
		if rs > s {
			s = rs
		}
		m.Sub(new(big.Int).Mul(lm, pow10(s-ls)), new(big.Int).Mul(rm, pow10(s-rs)))
	case "*":
		s = ls + rs
		m.Mul(lm, rm)
	default:
		// Unreachable past the static gate (dec / and % are
		// refused in checkSem); loud here so a direct caller can
		// never mistake silence for a quotient.
		return "", fmt.Errorf("bad dec %s operands", op)
	}
	neg := ""
	if m.Sign() < 0 {
		neg = "-"
		m.Abs(m)
	}
	digits := m.String()
	for len(digits) <= s {
		digits = "0" + digits
	}
	raw := neg + digits[:len(digits)-s] + "." + digits[len(digits)-s:]
	return canonDec(raw)
}

type Ctx struct {
	Prog *Program
	Test string
	// Depth counts nested local-helper calls in this test run. The
	// static proofs (cycle ban plus decreases) own termination; the
	// cap is a sandbox resource bound that fails closed and loud.
	Depth int
	// match node -> test -> remaining outcome nodes (nil entry = "-")
	Scripts map[*Node]map[string][]*Small
	// match node -> taken arm indexes, shared across a function's tests
	// for branch coverage. Nil means untracked.
	Cov map[*Node]map[int]bool
	// Store holds one test's private cells, keyed file + name: every
	// test starts from the decl literals, so sequential tests cannot
	// interfere. Helpers share the caller's map, threading state.
	Store map[string]*Value
}

// markTaken records an arm selection for branch coverage: the arm counts
// as taken when its pattern matches, even if its body later errors.
func markTaken(ctx *Ctx, node *Node, i int) {
	if ctx.Cov == nil {
		return
	}
	taken := ctx.Cov[node]
	if taken == nil {
		taken = map[int]bool{}
		ctx.Cov[node] = taken
	}
	taken[i] = true
}

type Program struct {
	Fns     map[string]*FnDecl
	Externs map[string]*ExternDecl
	// FnFile maps every function and extern name to its module file.
	// Same file means a local helper call (executed); anything else
	// means a foreign call (stubbed through given).
	FnFile map[string]string
	// Brands maps brand name to underlying type (v0: always str).
	// The emitter erases brands through this map.
	Brands map[string]string
	// BrandFile maps every brand name to its declaring module file.
	// Same file means an executable seal site (checked); anything
	// else means a foreign brand (refused in bodies, named in data).
	BrandFile map[string]string
	Errors    map[string][]string
	EmitsOf   map[string][]string
	Uses      map[string]bool
	Modules   []*Module
}

func vField(v *Value, field string) (*Value, error) {
	var d map[string]*Value
	switch v.Kind {
	case "rec", "ok", "err":
		d = v.Dict
	default:
		return nil, fmt.Errorf("no field %s on %s", field, v.Kind)
	}
	out, ok := d[field]
	if !ok {
		return nil, fmt.Errorf("missing field %s", field)
	}
	return out, nil
}

func vEq(a, b *Value) (bool, error) {
	if a.Kind != b.Kind {
		return false, nil
	}
	switch a.Kind {
	case "str":
		return a.S == b.S, nil
	case "int":
		return a.N.Cmp(b.N) == 0, nil
	case "dec":
		return a.D == b.D, nil
	case "bool":
		return a.B == b.B, nil
	case "rec", "ok":
		if len(a.Dict) != len(b.Dict) {
			return false, nil
		}
		for k, av := range a.Dict {
			bv, ok := b.Dict[k]
			if !ok {
				return false, nil
			}
			eq, err := vEq(av, bv)
			if err != nil || !eq {
				return eq, err
			}
		}
		return true, nil
	case "err":
		if a.ErrKind != b.ErrKind {
			return false, nil
		}
		return vEq(&Value{Kind: "rec", Dict: a.Dict}, &Value{Kind: "rec", Dict: b.Dict})
	case "seq":
		// Structural sequence comparison for the test evaluator
		// only (v36 S1): same element type, same length, ordered
		// member comparison. This is not a language equality
		// operator; == over sequences stays a compile error.
		if a.Elem != b.Elem || len(a.Arr) != len(b.Arr) {
			return false, nil
		}
		for i := range a.Arr {
			eq, err := vEq(a.Arr[i], b.Arr[i])
			if err != nil || !eq {
				return eq, err
			}
		}
		return true, nil
	case "bytes":
		// Structural byte comparison for the test evaluator only
		// (v45 S1): same length, ordered contents. This is not a
		// language equality operator; == over bytes stays a
		// compile error.
		if len(a.Bytes) != len(b.Bytes) {
			return false, nil
		}
		for i := range a.Bytes {
			if a.Bytes[i] != b.Bytes[i] {
				return false, nil
			}
		}
		return true, nil
	}
	return false, fmt.Errorf("cannot compare %s", a.Kind)
}

func evSmall(node *Small, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	switch node.Kind {
	case "str":
		return &Value{Kind: "str", S: node.Str}, nil
	case "int":
		return &Value{Kind: "int", N: new(big.Int).Set(node.Num)}, nil
	case "dec":
		if _, err := decRat(node.Dec); err != nil {
			return nil, err
		}
		return &Value{Kind: "dec", D: node.Dec}, nil
	case "float":
		return nil, fmt.Errorf("float %s is ungrammatical: write d\"%s\"", node.Str, node.Str)
	case "bool":
		return &Value{Kind: "bool", B: node.B}, nil
	case "seal":
		// Brands erase at runtime: a sealed value is its string.
		if len(node.Args) != 1 {
			return nil, fmt.Errorf("seal takes one string literal")
		}
		v, err := evSmall(node.Args[0].V, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		if v.Kind != "str" {
			return nil, fmt.Errorf("seal takes a string")
		}
		return &Value{Kind: "str", S: v.S}, nil
	case "wild":
		return nil, fmt.Errorf("lone _ is not a value")
	case "binop":
		lv, err := evSmall(node.L, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		rv, err := evSmall(node.R, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		if node.Op == "==" {
			eq, err := vEq(lv, rv)
			return &Value{Kind: "bool", B: eq}, err
		}
		switch node.Op {
		case "+", "-", "*", "/", "%":
			return evArith(node.Op, lv, rv)
		}
		if lv.Kind != rv.Kind || (lv.Kind != "int" && lv.Kind != "str" && lv.Kind != "dec") {
			return nil, fmt.Errorf("bad %s operands", node.Op)
		}
		if lv.Kind == "dec" {
			lr, err := decRat(lv.D)
			if err != nil {
				return nil, err
			}
			rr, err := decRat(rv.D)
			if err != nil {
				return nil, err
			}
			return &Value{Kind: "bool", B: decCmp(node.Op, lr.Cmp(rr))}, nil
		}
		switch node.Op {
		case ">=":
			if lv.Kind == "int" {
				return &Value{Kind: "bool", B: lv.N.Cmp(rv.N) >= 0}, nil
			}
			return &Value{Kind: "bool", B: lv.S >= rv.S}, nil
		case "<=":
			if lv.Kind == "int" {
				return &Value{Kind: "bool", B: lv.N.Cmp(rv.N) <= 0}, nil
			}
			return &Value{Kind: "bool", B: lv.S <= rv.S}, nil
		case ">":
			if lv.Kind == "int" {
				return &Value{Kind: "bool", B: lv.N.Cmp(rv.N) > 0}, nil
			}
			return &Value{Kind: "bool", B: lv.S > rv.S}, nil
		case "<":
			if lv.Kind == "int" {
				return &Value{Kind: "bool", B: lv.N.Cmp(rv.N) < 0}, nil
			}
			return &Value{Kind: "bool", B: lv.S < rv.S}, nil
		default:
			if lv.Kind == "int" {
				return &Value{Kind: "bool", B: lv.N.Cmp(rv.N) != 0}, nil
			}
			return &Value{Kind: "bool", B: lv.S != rv.S}, nil
		}
	case "strlen":
		v, err := evSmall(node.L, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		// v37 S2: sequences count elements, strings count
		// scalars. Members are never inspected, so brands need
		// no case here.
		if v.Kind == "seq" {
			return &Value{Kind: "int", N: big.NewInt(int64(len(v.Arr)))}, nil
		}
		if v.Kind != "str" {
			return nil, fmt.Errorf("bad length operand")
		}
		return &Value{Kind: "int", N: big.NewInt(int64(len([]rune(v.S))))}, nil
	case "stridx":
		b, err := evSmall(node.L, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		ix, err := evSmall(node.R, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		// v38 S3: sequences fetch members, never scalars. Bounds
		// never clamp: an unguarded out-of-range index (or a
		// non-int64 one) fails loud, reachable only without the
		// .ail guards the checked wrapper owns.
		if b.Kind == "seq" {
			if ix.Kind != "int" {
				return nil, fmt.Errorf("bad index operands")
			}
			if !ix.N.IsInt64() {
				return nil, fmt.Errorf("seq index out of range")
			}
			i := ix.N.Int64()
			if i < 0 || i >= int64(len(b.Arr)) {
				return nil, fmt.Errorf("seq index out of range")
			}
			return b.Arr[i], nil
		}
		if b.Kind != "str" || ix.Kind != "int" {
			return nil, fmt.Errorf("bad index operands")
		}
		rs := []rune(b.S)
		if !ix.N.IsInt64() {
			return nil, fmt.Errorf("str index out of range")
		}
		i := ix.N.Int64()
		if i < 0 || i >= int64(len(rs)) {
			return nil, fmt.Errorf("str index out of range")
		}
		return &Value{Kind: "int", N: big.NewInt(int64(rs[i]))}, nil
	case "strslice":
		b, err := evSmall(node.L, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		lo, err := evSmall(node.R, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		hi, err := evSmall(node.Hi, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		if b.Kind != "str" || lo.Kind != "int" || hi.Kind != "int" {
			return nil, fmt.Errorf("bad slice operands")
		}
		rs := []rune(b.S)
		if !lo.N.IsInt64() || !hi.N.IsInt64() {
			return nil, fmt.Errorf("str slice out of range")
		}
		l, h := lo.N.Int64(), hi.N.Int64()
		if l < 0 || h < l || h > int64(len(rs)) {
			return nil, fmt.Errorf("str slice out of range")
		}
		return &Value{Kind: "str", S: string(rs[l:h])}, nil
	case "call":
		return nil, fmt.Errorf("call outside match scrutinee is outside the v0 subset")
	case "exchange":
		return nil, fmt.Errorf("exchange outside a script row is outside the v0 subset")
	case "ctor":
		if node.Ctor == "Bytes" {
			// v45 S1: literal-only construction, validated before
			// creating the value. The checker admits only Seq<int>
			// literals of in-range integer members; re-check here
			// so execution (direct callers included) never invents
			// an unchecked shape. No operation mutates the result.
			if len(node.Args) != 1 || node.Args[0].HasName {
				return nil, fmt.Errorf("Bytes takes one Seq<int> literal")
			}
			seq, err := evSmall(node.Args[0].V, env, ctx, owner)
			if err != nil {
				return nil, err
			}
			if seq.Kind != "seq" || seq.Elem != "int" {
				return nil, fmt.Errorf("Bytes takes one Seq<int> literal")
			}
			out := make([]byte, 0, len(seq.Arr))
			for _, m := range seq.Arr {
				if m.Kind != "int" || m.N == nil || m.N.Sign() < 0 || m.N.Cmp(big.NewInt(256)) >= 0 {
					return nil, fmt.Errorf("Bytes member out of range 0..255")
				}
				out = append(out, byte(m.N.Int64()))
			}
			return &Value{Kind: "bytes", Bytes: out}, nil
		}
		if !isKwargList(node.Args) {
			return nil, fmt.Errorf("positional construction is outside the v0 subset: %s", node.Ctor)
		}
		fields := map[string]*Value{}
		for _, a := range node.Args {
			if _, dup := fields[a.Name]; dup {
				return nil, fmt.Errorf("%s: %s repeats field %s", owner, node.Ctor, a.Name)
			}
			v, err := evSmall(a.V, env, ctx, owner)
			if err != nil {
				return nil, err
			}
			fields[a.Name] = v
		}
		if node.Ctor == "Ok" {
			return &Value{Kind: "ok", Dict: fields}, nil
		}
		if strings.Contains(node.Ctor, ".") {
			if decl, ok := ctx.Prog.Errors[node.Ctor]; ok {
				want := map[string]bool{}
				for _, f := range decl {
					want[f] = true
				}
				for f := range fields {
					if !want[f] {
						return nil, fmt.Errorf("%s: %s has unknown field %s", owner, node.Ctor, f)
					}
				}
				for f := range want {
					if _, ok := fields[f]; !ok {
						return nil, fmt.Errorf("%s: %s missing field %s", owner, node.Ctor, f)
					}
				}
			}
			return &Value{Kind: "err", ErrKind: node.Ctor, Dict: fields}, nil
		}
		// Declared finite monomorphic records in value positions:
		// exact named fields with explicit declared types (checked
		// statically; re-verified here so execution never invents a
		// shape the declaration does not name).
		if decl := recordDecl(ctx.Prog, node.Ctor); decl != nil {
			want := map[string]bool{}
			for _, f := range decl.Fields {
				want[f[0]] = true
			}
			for f := range fields {
				if !want[f] {
					return nil, fmt.Errorf("%s: %s has unknown field %s", owner, node.Ctor, f)
				}
			}
			for f := range want {
				if _, ok := fields[f]; !ok {
					return nil, fmt.Errorf("%s: %s missing field %s", owner, node.Ctor, f)
				}
			}
			return &Value{Kind: "rec", Rec: node.Ctor, Dict: fields}, nil
		}
		return nil, fmt.Errorf("unknown constructor: %s", node.Ctor)
	case "list":
		return nil, fmt.Errorf("list literal outside given is outside the v0 subset")
	case "seqlit":
		// v36 S1: order, empties, and repeats preserved exactly.
		// Members evaluate left to right; the first failure ends
		// the literal, like every other strict position.
		arr := make([]*Value, 0, len(node.Items))
		for _, it := range node.Items {
			v, err := evSmall(it, env, ctx, owner)
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		return &Value{Kind: "seq", Arr: arr, Elem: node.Elem}, nil
	case "ref":
		if len(node.Ref) == 1 {
			v, ok := env[node.Ref[0]]
			if !ok {
				return nil, fmt.Errorf("unbound name: %s", node.Ref[0])
			}
			return v, nil
		}
		v, err := evSmall(&Small{Kind: "ref", Ref: node.Ref[:len(node.Ref)-1]}, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		return vField(v, node.Ref[len(node.Ref)-1])
	}
	return nil, fmt.Errorf("cannot evaluate: %s", node.Kind)
}

// localCallee resolves a same-file helper call: the callee must be a
// function defined in the caller's own file. Anything else (pinned
// foreign ail, externs, unknown names) is not local.
func localCallee(prog *Program, owner, fname string) *FnDecl {
	if prog == nil {
		return nil
	}
	fn, ok := prog.Fns[fname]
	if !ok {
		return nil
	}
	ownerFile, ok := prog.FnFile[owner]
	if !ok {
		return nil
	}
	calleeFile, ok := prog.FnFile[fname]
	if !ok || calleeFile != ownerFile {
		return nil
	}
	return fn
}

// evLocalCall executes a helper body inline: arguments evaluate in
// the caller's environment, bind to the helper params (by name when
// named, by position otherwise), and the body runs under the same
// test, so inner given tables script by caller test name. Coverage
// marks land in the shared map, attributing helper arms to the
// caller test that flowed through them.
func evLocalCall(helper *FnDecl, scrut *Small, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	// Sandbox resource bound, loud on breach: the termination proof
	// is static (cycle ban plus decreases), this only backstops it.
	if ctx.Depth >= 1024 {
		return nil, fmt.Errorf("%s: local call depth exceeded calling %s", owner, helper.Name)
	}
	// One binding rule (bindSlots): argument expressions evaluate in
	// source order here, then land in the resolved parameter slots.
	// Checked programs already passed the same rule, so a bind failure
	// below is a loud internal refusal, never a silent mis-call.
	slots, berr := bindSlots(helper.Name, scrut.Args, helper.Params)
	if berr != nil {
		return nil, fmt.Errorf("%s: %s", owner, berr.Error())
	}
	vals := make([]*Value, len(helper.Params))
	for i, a := range scrut.Args {
		v, err := evSmall(a.V, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		vals[slots[i]] = v
	}
	env2 := map[string]*Value{}
	for i, p := range helper.Params {
		env2[p[0]] = vals[i]
	}
	// No negative-entry check here (v11): every admitted self-call
	// sits under the false arm of its p <= 0 guard, so a negative
	// entry takes the base arm before any recursion. The sandbox depth
	// bound below stays as the resource backstop against hangs.
	ctx.Depth++
	got, err := evNode(helper.Body, env2, ctx, helper.Name)
	ctx.Depth--
	return got, err
}

// storeKey scopes a cell to its file: cells are module-private, so
// two files may hold independent cells under one name.
func storeKey(file, cell string) string {
	return file + "\x00" + cell
}

// freshStore builds one test's private store from the decl literals:
// every test starts from init, so sequential tests cannot interfere.
func freshStore(prog *Program) (map[string]*Value, error) {
	out := map[string]*Value{}
	for _, m := range prog.Modules {
		for _, d := range m.Decls {
			s, ok := d.(*StateDecl)
			if !ok {
				continue
			}
			v, err := evSmall(s.Init, map[string]*Value{}, &Ctx{Prog: prog}, m.Mod)
			if err != nil {
				return nil, err
			}
			out[storeKey(m.ID, s.Name)] = v
		}
	}
	return out, nil
}

// evStoreOp executes a cell operation against the test's store: get
// yields Ok(value), put stores its strictly-evaluated value and
// yields empty Ok. Static checks own shape, cells, and authority;
// here a missing store or cell fails closed.
func evStoreOp(scrut *Small, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	if ctx.Store == nil {
		return nil, fmt.Errorf("%s: no store for call %s", owner, scrut.Fname)
	}
	cell, ok := storeCellName(scrut)
	if !ok {
		return nil, fmt.Errorf("%s: call %s takes the cell first", owner, scrut.Fname)
	}
	file, ok := ctx.Prog.FnFile[owner]
	if !ok {
		return nil, fmt.Errorf("%s: unknown owner file", owner)
	}
	cur, ok := ctx.Store[storeKey(file, cell)]
	if !ok {
		return nil, fmt.Errorf("%s: unknown cell %s", owner, cell)
	}
	if scrut.Fname == "state__get" {
		if len(scrut.Args) != 1 {
			return nil, fmt.Errorf("%s: call %s takes 1 arg", owner, scrut.Fname)
		}
		return &Value{Kind: "ok", Dict: map[string]*Value{"value": cur}}, nil
	}
	if len(scrut.Args) != 2 || scrut.Args[0].HasName || scrut.Args[1].HasName {
		return nil, fmt.Errorf("%s: call %s takes the cell and a value", owner, scrut.Fname)
	}
	v, err := evSmall(scrut.Args[1].V, env, ctx, owner)
	if err != nil {
		return nil, err
	}
	ctx.Store[storeKey(file, cell)] = v
	return &Value{Kind: "ok", Dict: map[string]*Value{}}, nil
}

// evDecPartsOp evaluates the decimal observation kernel: one
// positional dec operand in, Ok(coefficient, scale) out. Canonical
// digits read straight through parseDecParts — literals and
// arithmetic results are canonical by construction, and the sign
// rides on the coefficient. The static gate owns operand misuse, so
// anything else here is loud, never a value.
func evDecPartsOp(scrut *Small, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	if len(scrut.Args) != 1 || scrut.Args[0].HasName {
		return nil, fmt.Errorf("%s: call dec__parts takes one value", owner)
	}
	v, err := evSmall(scrut.Args[0].V, env, ctx, owner)
	if err != nil {
		return nil, err
	}
	if v.Kind != "dec" {
		return nil, fmt.Errorf("bad dec__parts operand")
	}
	mant, scale, err := parseDecParts(v.D)
	if err != nil {
		return nil, err
	}
	return &Value{Kind: "ok", Dict: map[string]*Value{
		"coefficient": {Kind: "int", N: mant},
		"scale":       {Kind: "int", N: big.NewInt(int64(scale))},
	}}, nil
}

// matchSlot tests one value against one pattern, binding a variant
// payload into bind. The per-slot semantics are the historical single
// arm semantics, so arity 1 behaves exactly as the old shared loop.
func matchSlot(v *Value, p Pattern, bind map[string]*Value) bool {
	switch p.Kind {
	case "wild":
		return true
	case "bool":
		return v.Kind == "bool" && v.B == p.B
	case "str":
		return v.Kind == "str" && v.S == p.Str
	case "variantWild":
		return v.Kind == "err" && v.ErrKind == p.Name
	case "variant":
		if p.Name == "Ok" && v.Kind == "ok" {
			bind[p.Var] = &Value{Kind: "rec", Dict: v.Dict}
			return true
		}
		if v.Kind == "err" && v.ErrKind == p.Name {
			bind[p.Var] = &Value{Kind: "rec", Dict: v.Dict}
			return true
		}
		return false
	default:
		return false
	}
}

func evMatch(node *Node, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	if node.Kind == MatchCall {
		return evCallMatch(node, env, ctx, owner)
	}
	return evValueMatch(node, env, ctx, owner)
}

// evValueMatch evaluates a value table of any arity 1..N (docs/v28):
// every scrutinee evaluates exactly once, left to right, then the
// first arm whose every slot matches wins, binding variant payloads
// exactly as the old single loop did. A call scrutinee is a shape
// error owned by the checker; here it fails closed instead of
// evaluating. An empty hit is unreachable past the exhaustiveness
// proof and errors the same way.
func evValueMatch(node *Node, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	vals := make([]*Value, 0, len(node.Scruts))
	for _, s := range node.Scruts {
		if s.Kind == "call" {
			return nil, fmt.Errorf("%s: multi-scrutinee match over call %s is not supported: match a single call per match", owner, s.Fname)
		}
		v, err := evSmall(s, env, ctx, owner)
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	for i, arm := range node.Arms {
		if len(arm.Pats) != len(vals) {
			continue // owned by the arity rule; unmatchable here
		}
		rhsEnv := env
		bound := false
		hit := true
		for j, p := range arm.Pats {
			if p.Kind == "variant" && !bound {
				rhsEnv = copyEnv(env)
				bound = true
			}
			if !matchSlot(vals[j], p, rhsEnv) {
				hit = false
				break
			}
		}
		if hit {
			markTaken(ctx, node, i)
			return evNode(arm.Rhs, rhsEnv, ctx, owner)
		}
	}
	// Legacy rendering for arity 1 is byte-pinned; tuples name themselves.
	if len(vals) == 1 {
		return nil, fmt.Errorf("%s/%s: non-exhaustive match on %s", owner, ctx.Test, vals[0].Kind)
	}
	return nil, fmt.Errorf("%s/%s: non-exhaustive multi-scrutinee match", owner, ctx.Test)
}

// evBytesEncodeOp evaluates UTF-8 encoding (v46 S2, v47 S3): the input
// string as UTF-8 bytes, NUL and BOM preserved. Restricted calls
// refuse loud without a certificate; the refusal must never become
// trusted script evidence (certificates issue before linkage for
// exactly this reason).
func evBytesEncodeOp(scrut *Small, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	if isBytesExport(scrut.Fname) && scrut.ExportBrand == "" {
		return nil, fmt.Errorf("%s: call to %s is not authorized by a valid exports_utf8 grant", owner, scrut.Fname)
	}
	var argv *Small
	if isBytesExport(scrut.Fname) {
		// Restricted calls carry no static signature; their exact
		// shape was proven by the certifier, so re-check it here.
		if len(scrut.Args) != 1 || scrut.Args[0].HasName {
			return nil, fmt.Errorf("%s: call to %s takes the bare granted value", owner, scrut.Fname)
		}
		argv = scrut.Args[0].V
	} else {
		slots, berr := bindSlots(scrut.Fname, scrut.Args, bytesKernels[scrut.Fname].params)
		if berr != nil {
			return nil, fmt.Errorf("%s: %s", owner, berr.Error())
		}
		for i, s := range slots {
			if s == 0 {
				argv = scrut.Args[i].V
			}
		}
	}
	v, err := evSmall(argv, env, ctx, owner)
	if err != nil {
		return nil, err
	}
	if v.Kind != "str" {
		return nil, fmt.Errorf("%s: call to %s takes str", owner, scrut.Fname)
	}
	if !utf8.ValidString(v.S) {
		return nil, fmt.Errorf("%s: host string is not valid UTF-8", owner)
	}
	out := make([]byte, len(v.S))
	copy(out, v.S)
	return &Value{Kind: "ok", Dict: map[string]*Value{"value": {Kind: "bytes", Bytes: out}}}, nil
}

// evBytesDecodeOp evaluates UTF-8 decoding (v50 B6): the input Bytes
// as a string when the whole input is valid UTF-8, else the
// encoding.invalid_utf8 language error carrying the ORIGINAL payload
// unchanged — with a nil Go error. Malformed input is a computed
// comparable result: a Go error here would launder false scripted
// successes into trusted linkage evidence (contradictScriptOk treats
// provider Go errors as "not contradicted").
func evBytesDecodeOp(scrut *Small, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	slots, berr := bindSlots(scrut.Fname, scrut.Args, bytesKernels[scrut.Fname].params)
	if berr != nil {
		return nil, fmt.Errorf("%s: %s", owner, berr.Error())
	}
	var argv *Small
	for i, s := range slots {
		if s == 0 {
			argv = scrut.Args[i].V
		}
	}
	v, err := evSmall(argv, env, ctx, owner)
	if err != nil {
		return nil, err
	}
	if v.Kind != "bytes" {
		return nil, fmt.Errorf("%s: call to %s takes Bytes", owner, scrut.Fname)
	}
	if !utf8.Valid(v.Bytes) {
		out := make([]byte, len(v.Bytes))
		copy(out, v.Bytes)
		return &Value{Kind: "err", ErrKind: encodingInvalidUtf8, Dict: map[string]*Value{"value": {Kind: "bytes", Bytes: out}}}, nil
	}
	return &Value{Kind: "ok", Dict: map[string]*Value{"value": {Kind: "str", S: string(v.Bytes)}}}, nil
}

// evBytesHexEncodeOp evaluates hex encoding (v52 B8): the input
// Bytes as lowercase hex, byte-ordered, empty to "". Bytes are
// never read as text: ASCII-looking octets encode as digits.
func evBytesHexEncodeOp(scrut *Small, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	slots, berr := bindSlots(scrut.Fname, scrut.Args, bytesKernels[scrut.Fname].params)
	if berr != nil {
		return nil, fmt.Errorf("%s: %s", owner, berr.Error())
	}
	var argv *Small
	for i, s := range slots {
		if s == 0 {
			argv = scrut.Args[i].V
		}
	}
	v, err := evSmall(argv, env, ctx, owner)
	if err != nil {
		return nil, err
	}
	if v.Kind != "bytes" {
		return nil, fmt.Errorf("%s: call to %s takes Bytes", owner, scrut.Fname)
	}
	return &Value{Kind: "ok", Dict: map[string]*Value{"value": {Kind: "str", S: hex.EncodeToString(v.Bytes)}}}, nil
}

// evBytesHexDecodeOp evaluates hex decoding (v55 B10): an
// even-length ASCII hex string as Bytes, else the
// encoding.invalid_hex language error carrying the ORIGINAL string
// unchanged — with a nil Go error. The host's partial prefix is
// checked and discarded: it must never escape as success.
func evBytesHexDecodeOp(scrut *Small, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	slots, berr := bindSlots(scrut.Fname, scrut.Args, bytesKernels[scrut.Fname].params)
	if berr != nil {
		return nil, fmt.Errorf("%s: %s", owner, berr.Error())
	}
	var argv *Small
	for i, s := range slots {
		if s == 0 {
			argv = scrut.Args[i].V
		}
	}
	v, err := evSmall(argv, env, ctx, owner)
	if err != nil {
		return nil, err
	}
	if v.Kind != "str" {
		return nil, fmt.Errorf("%s: call to %s takes str", owner, scrut.Fname)
	}
	if !isHexStr(v.S) {
		return &Value{Kind: "err", ErrKind: encodingInvalidHex, Dict: map[string]*Value{"value": {Kind: "str", S: v.S}}}, nil
	}
	out, herr := hex.DecodeString(v.S)
	if herr != nil {
		return &Value{Kind: "err", ErrKind: encodingInvalidHex, Dict: map[string]*Value{"value": {Kind: "str", S: v.S}}}, nil
	}
	return &Value{Kind: "ok", Dict: map[string]*Value{"value": {Kind: "bytes", Bytes: out}}}, nil
}

// isHexStr reports the strict hex grammar: even length, every byte
// an ASCII hex digit. Byte-indexed: AIL strings that reach here are
// validated scalar-by-scalar upstream, and any non-ASCII byte fails
// the digit predicate regardless of position.
func isHexStr(s string) bool {
	if len(s)%2 != 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F' {
			continue
		}
		return false
	}
	return true
}

// evBytesB64EncodeOp evaluates base64 encoding (v58 B12): the input
// Bytes as standard padded base64, byte-ordered, empty to "". Bytes
// are never read as text.
func evBytesB64EncodeOp(scrut *Small, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	slots, berr := bindSlots(scrut.Fname, scrut.Args, bytesKernels[scrut.Fname].params)
	if berr != nil {
		return nil, fmt.Errorf("%s: %s", owner, berr.Error())
	}
	var argv *Small
	for i, s := range slots {
		if s == 0 {
			argv = scrut.Args[i].V
		}
	}
	v, err := evSmall(argv, env, ctx, owner)
	if err != nil {
		return nil, err
	}
	if v.Kind != "bytes" {
		return nil, fmt.Errorf("%s: call to %s takes Bytes", owner, scrut.Fname)
	}
	return &Value{Kind: "ok", Dict: map[string]*Value{"value": {Kind: "str", S: base64.StdEncoding.EncodeToString(v.Bytes)}}}, nil
}

func evCallMatch(node *Node, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	scrut := node.Scruts[0]
	var v *Value
	if scrut.Kind == "call" {
		fname := scrut.Fname
		if isStoreOp(fname) {
			if node.Given != nil {
				return nil, fmt.Errorf("%s: call to %s takes no given table", owner, fname)
			}
			val, err := evStoreOp(scrut, env, ctx, owner)
			if err != nil {
				return nil, err
			}
			v = val
		} else if isDecParts(fname) {
			if node.Given != nil {
				return nil, fmt.Errorf("%s: call to %s takes no given table", owner, fname)
			}
			val, err := evDecPartsOp(scrut, env, ctx, owner)
			if err != nil {
				return nil, err
			}
			v = val
		} else if isBytesKernel(fname) {
			if node.Given != nil {
				return nil, fmt.Errorf("%s: call to %s takes no given table", owner, fname)
			}
			var val *Value
			var err error
			if isBytesDecode(fname) {
				val, err = evBytesDecodeOp(scrut, env, ctx, owner)
			} else if isBytesHexDecode(fname) {
				val, err = evBytesHexDecodeOp(scrut, env, ctx, owner)
			} else if isBytesHexEncode(fname) {
				val, err = evBytesHexEncodeOp(scrut, env, ctx, owner)
			} else if isBytesB64Encode(fname) {
				val, err = evBytesB64EncodeOp(scrut, env, ctx, owner)
			} else {
				val, err = evBytesEncodeOp(scrut, env, ctx, owner)
			}
			if err != nil {
				return nil, err
			}
			v = val
		} else if helper := localCallee(ctx.Prog, owner, fname); helper != nil {
			if node.Given != nil {
				return nil, fmt.Errorf("%s: call to local helper %s takes no given table", owner, fname)
			}
			val, err := evLocalCall(helper, scrut, env, ctx, owner)
			if err != nil {
				return nil, err
			}
			v = val
		} else {
			// Externs are module-local foreign imports: no uses pin, but
			// still scripted through given tables like ail calls.
			if !ctx.Prog.Uses[fname] && ctx.Prog.Externs[fname] == nil {
				return nil, fmt.Errorf("%s: %s not in uses", owner, fname)
			}
			entry, ok := node.Given[ctx.Test]
			if !ok || entry == nil {
				return nil, fmt.Errorf("%s/%s: no script for call %s", owner, ctx.Test, fname)
			}
			perTest, ok := ctx.Scripts[node]
			if !ok {
				perTest = map[string][]*Small{}
				ctx.Scripts[node] = perTest
			}
			script, started := perTest[ctx.Test]
			if !started {
				if entry.Kind == "list" {
					script = append([]*Small{}, entry.Items...)
				} else {
					script = []*Small{entry}
				}
				perTest[ctx.Test] = script
			}
			if len(script) == 0 {
				return nil, fmt.Errorf("%s/%s: call %s script exhausted", owner, ctx.Test, fname)
			}
			item := script[0]
			perTest[ctx.Test] = script[1:]
			if item.Kind != "exchange" {
				return nil, fmt.Errorf("%s/%s: script row must be an exchange with args and outcome", owner, ctx.Test)
			}
			if err := checkExchangeArgs(scrut, item, ctx.Prog, env, ctx, owner); err != nil {
				return nil, err
			}
			if err := checkExchangeArgs(scrut, item, ctx.Prog, env, ctx, owner); err != nil {
				return nil, err
			}
			val, err := evSmall(item.Outcome, env, ctx, owner)
			if err != nil {
				return nil, err
			}
			if val.Kind == "err" {
				allowed := false
				for _, e := range ctx.Prog.EmitsOf[fname] {
					if e == val.ErrKind {
						allowed = true
					}
				}
				if !allowed {
					return nil, fmt.Errorf("%s/%s: stub %s not in %s emits", owner, ctx.Test, val.ErrKind, fname)
				}
			} else if val.Kind != "ok" {
				return nil, fmt.Errorf("%s/%s: stub must be Ok(..) or an error", owner, ctx.Test)
			}
			v = val
		}
	} else {
		// Parser-impossible: MatchCall always carries one call. Fail
		// closed rather than evaluating a non-call as an outcome.
		return nil, fmt.Errorf("%s: call match without a call scrutinee", owner)
	}
	for i, arm := range node.Arms {
		pat := arm.Pats[0]
		switch pat.Kind {
		case "wild":
			markTaken(ctx, node, i)
			return evNode(arm.Rhs, env, ctx, owner)
		case "bool":
			if v.Kind == "bool" && v.B == pat.B {
				markTaken(ctx, node, i)
				return evNode(arm.Rhs, env, ctx, owner)
			}
		case "str":
			if v.Kind == "str" && v.S == pat.Str {
				markTaken(ctx, node, i)
				return evNode(arm.Rhs, env, ctx, owner)
			}
		case "variantWild":
			if v.Kind == "err" && v.ErrKind == pat.Name {
				markTaken(ctx, node, i)
				return evNode(arm.Rhs, env, ctx, owner)
			}
		case "variant":
			if pat.Name == "Ok" && v.Kind == "ok" {
				markTaken(ctx, node, i)
				env2 := copyEnv(env)
				env2[pat.Var] = &Value{Kind: "rec", Dict: v.Dict}
				return evNode(arm.Rhs, env2, ctx, owner)
			}
			if v.Kind == "err" && v.ErrKind == pat.Name {
				markTaken(ctx, node, i)
				env2 := copyEnv(env)
				env2[pat.Var] = &Value{Kind: "rec", Dict: v.Dict}
				return evNode(arm.Rhs, env2, ctx, owner)
			}
		}
	}
	return nil, fmt.Errorf("%s/%s: non-exhaustive match on %s", owner, ctx.Test, v.Kind)
}

// describe renders a value's shape for failure messages: err(db.down), ok.
func describe(v *Value) string {
	if v.Kind == "err" {
		return "err(" + v.ErrKind + ")"
	}
	return v.Kind
}

// normStr escapes exactly \\, \", \n, \r, \t by name and any other C0
// control as \u00XX. Printable Unicode passes through untouched: the
// canonical form is byte-stable, not ASCII-only.
func normStr(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 {
				fmt.Fprintf(&b, `\u%04x`, r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte('"')
	return b.String()
}

// normalizeValue renders a value in canonical form: exactly one textual
// shape per value, so golden tests, diffs, and prompt-pasted values are
// byte-comparable. Records sort keys byte-wise (never declaration order),
// errors name their kind with sorted fields when present.
func normalizeValue(v *Value) string {
	switch v.Kind {
	case "str":
		return normStr(v.S)
	case "int":
		return v.N.String()
	case "dec":
		return "d\"" + v.D + "\""
	case "bool":
		if v.B {
			return "true"
		}
		return "false"
	case "ok", "rec":
		keys := make([]string, 0, len(v.Dict))
		for k := range v.Dict {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, k+" = "+normalizeValue(v.Dict[k]))
		}
		name := "Ok"
		if v.Kind == "rec" && v.Rec != "" {
			name = v.Rec
		}
		return name + "(" + strings.Join(parts, ", ") + ")"
	case "seq":
		parts := make([]string, 0, len(v.Arr))
		for _, m := range v.Arr {
			parts = append(parts, normalizeValue(m))
		}
		return "Seq<" + v.Elem + ">[" + strings.Join(parts, ", ") + "]"
	case "bytes":
		octets := make([]string, 0, len(v.Bytes))
		for _, b := range v.Bytes {
			octets = append(octets, fmt.Sprintf("%d", b))
		}
		return "Bytes(Seq<int>[" + strings.Join(octets, ", ") + "])"
	case "err":
		if len(v.Dict) == 0 {
			return "err(" + v.ErrKind + ")"
		}
		keys := make([]string, 0, len(v.Dict))
		for k := range v.Dict {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, k+" = "+normalizeValue(v.Dict[k]))
		}
		return "err(" + v.ErrKind + "(" + strings.Join(parts, ", ") + "))"
	default:
		return "<unknown>"
	}
}

func copyEnv(env map[string]*Value) map[string]*Value {
	out := make(map[string]*Value, len(env)+1)
	for k, v := range env {
		out[k] = v
	}
	return out
}

func evNode(node *Node, env map[string]*Value, ctx *Ctx, owner string) (*Value, error) {
	if node.IsMatch {
		return evMatch(node, env, ctx, owner)
	}
	return evSmall(node.Small, env, ctx, owner)
}

// runTestValue executes one decision-table row and returns its outcome
// value: argument binding, body evaluation, and stub-consumption check.
// Expectation comparison stays in runTest; normalize prints the outcome
// itself, so it shares this prefix.
func runTestValue(fn *FnDecl, test Test, prog *Program, cov map[*Node]map[int]bool) (*Value, map[string]*Value, error) {
	argVals := map[string]*Small{}
	for _, a := range test.Args {
		argVals[a.Name] = a.V
	}
	env := map[string]*Value{}
	tmpCtx := &Ctx{Prog: prog, Test: test.Name, Scripts: map[*Node]map[string][]*Small{}, Cov: cov}
	for _, p := range fn.Params {
		v, ok := argVals[p[0]]
		if !ok {
			return nil, nil, fmt.Errorf("%s/%s: missing arg %s", fn.Name, test.Name, p[0])
		}
		val, err := evSmall(v, map[string]*Value{}, tmpCtx, fn.Name)
		if err != nil {
			return nil, nil, err
		}
		env[p[0]] = val
	}
	store, err := freshStore(prog)
	if err != nil {
		return nil, nil, err
	}
	ctx := &Ctx{Prog: prog, Test: test.Name, Scripts: map[*Node]map[string][]*Small{}, Cov: cov, Store: store}
	// No negative-entry check here either (v11): the static guard rule
	// admits recursion only under the positive branch, so a negative
	// entry evaluates the base arm and returns a declared outcome.
	got, err := evNode(fn.Body, env, ctx, fn.Name)
	if err != nil {
		return nil, nil, err
	}
	for _, perTest := range ctx.Scripts {
		if len(perTest[test.Name]) > 0 {
			return nil, nil, fmt.Errorf("%s/%s: stub leftovers at a call site", fn.Name, test.Name)
		}
	}
	return got, env, nil
}

// checkExchangeArgs proves "this request received this permitted
// response" (v12): every expected arg must arrive with an equal value
// under the same name, and the call must supply nothing unexpected.
// Positional call args resolve through the callee signature; expected
// args are always named (the parser refuses anything else).
func checkExchangeArgs(scrut, item *Small, prog *Program, env map[string]*Value, ctx *Ctx, owner string) error {
	fname := scrut.Fname
	var params [][2]string
	if fn, ok := prog.Fns[fname]; ok {
		params = fn.Params
	} else if ex, ok := prog.Externs[fname]; ok {
		params = ex.Params
	}
	actual := map[string]*Value{}
	for i, a := range scrut.Args {
		v, err := evSmall(a.V, env, ctx, owner)
		if err != nil {
			return err
		}
		name := a.Name
		if !a.HasName {
			if i >= len(params) {
				return fmt.Errorf("%s/%s: call %s takes %d args", owner, ctx.Test, fname, len(params))
			}
			name = params[i][0]
		}
		if _, dup := actual[name]; dup {
			return fmt.Errorf("%s/%s: call %s supplies arg %s twice", owner, ctx.Test, fname, name)
		}
		actual[name] = v
	}
	seen := map[string]bool{}
	for _, a := range item.Args {
		if seen[a.Name] {
			return fmt.Errorf("%s/%s: exchange repeats arg %s for call %s", owner, ctx.Test, a.Name, fname)
		}
		seen[a.Name] = true
		av, ok := actual[a.Name]
		if !ok {
			return fmt.Errorf("%s/%s: exchange expects arg %s for call %s, but the call supplies none", owner, ctx.Test, a.Name, fname)
		}
		ev, err := evSmall(a.V, env, ctx, owner)
		if err != nil {
			return err
		}
		eq, err := vEq(av, ev)
		if err != nil || !eq {
			return fmt.Errorf("%s/%s: call %s arg %s mismatch: expected %s", owner, ctx.Test, fname, a.Name, normalizeValue(ev))
		}
	}
	for name := range actual {
		if !seen[name] {
			return fmt.Errorf("%s/%s: call %s supplies arg %s with no exchange expectation", owner, ctx.Test, fname, name)
		}
	}
	return nil
}

func runTest(fn *FnDecl, test Test, prog *Program, cov map[*Node]map[int]bool) error {
	got, env, err := runTestValue(fn, test, prog, cov)
	if err != nil {
		return err
	}
	ctx := &Ctx{Prog: prog, Test: test.Name, Scripts: map[*Node]map[string][]*Small{}}
	want := test.Expected
	if want.Kind == "ctor" && want.Ctor == "Ok" {
		if got.Kind != "ok" {
			return fmt.Errorf("%s/%s: expected Ok, got %s", fn.Name, test.Name, describe(got))
		}
		exp, err := evSmall(want, env, ctx, fn.Name)
		if err != nil {
			return err
		}
		eq, err := vEq(got, exp)
		if err != nil || !eq {
			return fmt.Errorf("%s/%s: Ok payload mismatch", fn.Name, test.Name)
		}
		return nil
	}
	if want.Kind == "ctor" && strings.Contains(want.Ctor, ".") {
		// Complete error expectations compare kind and payload: the
		// table asserts the exact error value, not just its kind.
		if got.Kind != "err" || got.ErrKind != want.Ctor {
			return fmt.Errorf("%s/%s: expected %s, got %s", fn.Name, test.Name, want.Ctor, describe(got))
		}
		exp, err := evSmall(want, env, ctx, fn.Name)
		if err != nil {
			return err
		}
		eq, err := vEq(got, exp)
		if err != nil || !eq {
			return fmt.Errorf("%s/%s: error payload mismatch", fn.Name, test.Name)
		}
		return nil
	}
	if want.Kind == "ref" && len(want.Ref) > 1 {
		name := ""
		for i, p := range want.Ref {
			if i > 0 {
				name += "."
			}
			name += p
		}
		if got.Kind != "err" || got.ErrKind != name {
			return fmt.Errorf("%s/%s: expected %s, got %s", fn.Name, test.Name, name, describe(got))
		}
		return nil
	}
	return fmt.Errorf("%s/%s: bad expectation shape", fn.Name, test.Name)
}

var (
	fnNameRe   = regexp.MustCompile(`^[a-z][a-z0-9]*__[a-z][a-z0-9_]*$`)
	typeNameRe = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*__[A-Za-z0-9_]*$`)
)

// checkNaming enforces R3: strict naming grammar, part of every diagnosis.
func checkNaming(m *Module, text string) []Diag {
	var out []Diag
	for _, d := range m.Decls {
		switch d := d.(type) {
		case *FnDecl:
			if d.Name == "state__get" || d.Name == "state__put" {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("name %q is reserved for store operations", d.Name), d.Name, CodeFnNaming))
			} else if !fnNameRe.MatchString(d.Name) {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("function name %q must match domain__verb", d.Name), d.Name, CodeFnNaming))
			}
		case *TypeDecl:
			if !typeNameRe.MatchString(d.Name) {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("type name %q must match Domain__Name", d.Name), d.Name, CodeTypeNaming))
			}
		case *BrandDecl:
			if !typeNameRe.MatchString(d.Name) {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("brand name %q must match Domain__Name", d.Name), d.Name, CodeTypeNaming))
			}
		case *ExternDecl:
			if d.Name == "state__get" || d.Name == "state__put" {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("name %q is reserved for store operations", d.Name), d.Name, CodeFnNaming))
			} else if !fnNameRe.MatchString(d.Name) {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("extern name %q must match domain__verb", d.Name), d.Name, CodeFnNaming))
			}
		case *StateDecl:
			if !typeNameRe.MatchString(d.Name) {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("state name %q must match Domain__Name", d.Name), d.Name, CodeTypeNaming))
			}
		}
	}
	return out
}

// verifyExhaustive proves every match total before anything runs or emits:
// call-matches must cover exactly the callee emits plus ok (missing or
// stale arms are compile errors), bool matches must be exactly true+false,
// value matches must end with _. Proof lives here, not in shipped code.
//
// The walk collects every violation in order; verifyExhaustive reports the
// first (CLI behavior), verifyExhaustiveAll reports all (editor squiggles).
func verifyExhaustive(mods []*Module, prog *Program) error {
	if all := verifyExhaustiveAll(mods, prog); len(all) > 0 {
		return all[0]
	}
	return nil
}

func verifyExhaustiveAll(mods []*Module, prog *Program) []error {
	var out []error
	var walk func(n *Node, owner string)
	walk = func(n *Node, owner string) {
		if n == nil || !n.IsMatch {
			return
		}
		if n.Kind == MatchCall {
			fname := n.Scruts[0].Fname
			if isBytesKernel(fname) {
				// v46 S2: every kernel's contract must exist
				// explicitly; an absent entry is never an empty
				// error set (no dec__parts shortcut).
				if _, ok := prog.EmitsOf[fname]; !ok {
					out = append(out, at(n.Line, fmt.Errorf("%s: bytes kernel %s has no registered contract", owner, fname)))
				}
			}
			want := map[string]bool{"ok": true}
			for _, e := range prog.EmitsOf[fname] {
				want[e] = true
			}
			got := map[string]int{}
			for _, a := range n.Arms {
				switch p := a.Pats[0]; p.Kind {
				case "variantWild":
					got[p.Name] = a.Line
				case "variant":
					k := p.Name
					if k == "Ok" {
						k = "ok"
					}
					got[k] = a.Line
				default:
					out = append(out, at(a.Line, fmt.Errorf("%s: call-match arm must be an error kind or Ok", owner)))
				}
				walk(a.Rhs, owner)
			}
			for _, k := range sortedKeys(want) {
				if _, ok := got[k]; !ok {
					out = append(out, at(n.Line, fmt.Errorf("%s: non-exhaustive match, missing %s", owner, k)))
				}
			}
			for _, k := range sortedKeys(got) {
				if !want[k] {
					out = append(out, at(got[k], fmt.Errorf("%s: stale match arm %s", owner, k)))
				}
			}
			return
		}
		out = append(out, verifyValueMatch(n, owner)...)
		for _, a := range n.Arms {
			walk(a.Rhs, owner)
		}
		return
	}
	for _, m := range mods {
		for _, d := range m.Decls {
			if fn, ok := d.(*FnDecl); ok {
				walk(fn.Body, fn.Name)
			}
		}
	}
	return out
}

// valueAtom is one symbolic inhabitant of a match slot's domain: a bool
// literal, a mentioned string literal, the open string remainder (OTHER),
// or the whole domain of an all-wildcard slot (ANY, rendered as _).
type valueAtom struct {
	render string
	other  bool
}

// armCover is one arm as covered atom indices per slot.
type armCover struct{ slots [][]int }

func containsAtom(xs []int, v int) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

// otherWitness picks a concrete string outside the mentioned literals,
// so an uncovered open remainder renders as a real counterexample like
// (true, "") rather than a diagnostic-only pattern.
func otherWitness(lits []string) string {
	mentioned := map[string]bool{}
	for _, l := range lits {
		mentioned[l] = true
	}
	if !mentioned[""] {
		return ""
	}
	for i := 1; ; i++ {
		s := strings.Repeat("a", i)
		if !mentioned[s] {
			return s
		}
	}
}

// witnessSearch finds one tuple covered by none of arms, or nil when the
// arms cover the whole domain product. Depth-first descent prunes covered
// prefixes slot by slot, so the product is never materialized. Returns
// per-slot atom indices plus their renders.
func witnessSearch(arms []armCover, nslot int, domains [][]valueAtom) (idx []int, render []string) {
	var pi []int
	var pr []string
	var rec func(slot int, live []armCover) bool
	rec = func(slot int, live []armCover) bool {
		if slot == nslot {
			return len(live) == 0
		}
		for ai, a := range domains[slot] {
			var next []armCover
			for _, arm := range live {
				if containsAtom(arm.slots[slot], ai) {
					next = append(next, arm)
				}
			}
			pi = append(pi, ai)
			pr = append(pr, a.render)
			if rec(slot+1, next) {
				return true
			}
			pi = pi[:len(pi)-1]
			pr = pr[:len(pr)-1]
		}
		return false
	}
	if !rec(0, arms) {
		return nil, nil
	}
	return append([]int{}, pi...), append([]string{}, pr...)
}

// valueCoverOf builds the symbolic slot domains and arm covers for a
// multi-scrutinee match: mentioned string literals plus the open
// remainder, both bool literals, or a single whole-domain atom for
// all-wild slots. Shared by the exhaustiveness proof and the emitter's
// residual check. Requires one pattern per slot on every arm (callers
// establish arity first). Bool renders ("true"/"false") and normStr
// renders never collide, and OTHER's concrete witness is chosen outside
// the mentioned literals, so render-keyed lookup is unambiguous.
// Returns the first 1-based slot mixing bool and string literals, or 0.
func valueCoverOf(n *Node, hasBool []bool, strLits [][]string) (domains [][]valueAtom, covers []armCover, atomIndex []map[string]int, mixed int) {
	nslot := len(n.Scruts)
	domains = make([][]valueAtom, nslot)
	for i := 0; i < nslot; i++ {
		switch {
		case len(strLits[i]) > 0 && hasBool[i]:
			return nil, nil, nil, i + 1
		case len(strLits[i]) > 0:
			for _, l := range strLits[i] {
				domains[i] = append(domains[i], valueAtom{render: normStr(l)})
			}
			domains[i] = append(domains[i], valueAtom{render: normStr(otherWitness(strLits[i])), other: true})
		case hasBool[i]:
			domains[i] = []valueAtom{{render: "true"}, {render: "false"}}
		default:
			domains[i] = []valueAtom{{render: "_"}}
		}
	}
	atomIndex = make([]map[string]int, nslot)
	for i, d := range domains {
		atomIndex[i] = map[string]int{}
		for ai := range d {
			atomIndex[i][d[ai].render] = ai
		}
	}
	covers = make([]armCover, 0, len(n.Arms))
	for _, a := range n.Arms {
		c := armCover{slots: make([][]int, nslot)}
		for i, p := range a.Pats {
			switch p.Kind {
			case "wild":
				for ai := range domains[i] {
					c.slots[i] = append(c.slots[i], ai)
				}
			case "bool":
				if p.B {
					c.slots[i] = []int{atomIndex[i]["true"]}
				} else {
					c.slots[i] = []int{atomIndex[i]["false"]}
				}
			case "str":
				c.slots[i] = []int{atomIndex[i][normStr(p.Str)]}
			}
		}
		covers = append(covers, c)
	}
	return domains, covers, atomIndex, 0
}

// valueMatchAnalysis is the proven-table fact emit consumes: whether the
// last arm owns the whole residual, so a bare else is sound. Set only
// for proven tables (nil otherwise); the emitter never recomputes it.
type valueMatchAnalysis struct {
	emitFinalElse bool
}

// residualInLast reports whether every tuple outside the prior arms'
// covers lies inside the last arm's. Bounded: an enormous residual
// stays a tested condition instead. Proof logic lives here, with the
// checker — never in the emitter.
func residualInLast(prior []armCover, last armCover, nslot int, domains [][]valueAtom) bool {
	live := append([]armCover{}, prior...)
	for i := 0; i < 256; i++ {
		idx, _ := witnessSearch(live, nslot, domains)
		if idx == nil {
			return true
		}
		for s, a := range idx {
			if !containsAtom(last.slots[s], a) {
				return false
			}
		}
		lit := armCover{slots: make([][]int, nslot)}
		for s, a := range idx {
			lit.slots[s] = []int{a}
		}
		live = append(live, lit)
	}
	return false
}

// verifyValueMatch proves a value table of any arity 1..N total
// (docs/v28): no call scrutinees, one pattern per slot on every arm,
// bool and string literals never mixed in one slot, and the arms'
// product spaces covering the total space. Arity 1 keeps the historical
// single policy verbatim (AIL4103/AIL4104 byte-pinned; the product
// engine would accept more, and that language change is parked, not
// smuggled in). Diagnostics reuse the single-match codes, tuple-rendered
// past arity 1 and anchored by the same message shapes proofDiag matches
// on. Proven tables record the emit bit; anything else leaves analysis
// nil, which the emitter reads as "test the last arm".
func verifyValueMatch(n *Node, owner string) []error {
	var out []error
	nslot := len(n.Scruts)
	for _, s := range n.Scruts {
		if s.Kind == "call" {
			out = append(out, at(n.Line, fmt.Errorf("%s: multi-scrutinee match over call %s is not supported: match a single call per match", owner, s.Fname)))
		}
	}
	// Per-slot pattern inventory across arms.
	hasBool := make([]bool, nslot)
	strLits := make([][]string, nslot)
	seenLit := make([]map[string]bool, nslot)
	for i := range seenLit {
		seenLit[i] = map[string]bool{}
	}
	shapeOK := true
	for _, a := range n.Arms {
		if len(a.Pats) != nslot {
			out = append(out, at(a.Line, fmt.Errorf("%s: match arm has %d patterns; this match has %d scrutinees", owner, len(a.Pats), nslot)))
			shapeOK = false
			continue
		}
		for i, p := range a.Pats {
			switch p.Kind {
			case "bool":
				hasBool[i] = true
			case "str":
				if !seenLit[i][p.Str] {
					seenLit[i][p.Str] = true
					strLits[i] = append(strLits[i], p.Str)
				}
			case "wild":
			default:
				if nslot == 1 {
					out = append(out, at(a.Line, fmt.Errorf("%s: variant pattern on a non-call match", owner)))
				} else {
					out = append(out, at(a.Line, fmt.Errorf("%s: variant pattern on a non-call match (slot %d)", owner, i+1)))
				}
				shapeOK = false
			}
		}
	}
	if !shapeOK {
		return out
	}
	if nslot == 1 {
		return verifySinglePolicy(n, owner, out, hasBool, strLits)
	}
	domains, covers, _, mixed := valueCoverOf(n, hasBool, strLits)
	if mixed > 0 {
		out = append(out, at(n.Line, fmt.Errorf("%s: bool match must be exactly true+false (slot %d mixes bool and string patterns)", owner, mixed)))
		return out
	}
	// Collect up to three missing witnesses, feeding each back as a
	// literal arm so the next search finds a new cell.
	var wits []string
	otherSlot := 0
	live := covers
	for len(wits) < 3 {
		idx, renders := witnessSearch(live, nslot, domains)
		if idx == nil {
			break
		}
		wits = append(wits, "("+strings.Join(renders, ", ")+")")
		lit := armCover{slots: make([][]int, nslot)}
		for s, a := range idx {
			lit.slots[s] = []int{a}
			if domains[s][a].other && otherSlot == 0 {
				otherSlot = s + 1
			}
		}
		live = append(live, lit)
	}
	if len(wits) == 0 {
		n.analysis = &valueMatchAnalysis{emitFinalElse: residualInLast(covers[:len(covers)-1], covers[len(covers)-1], nslot, domains)}
		return out
	}
	// An uncovered open string remainder dominates: infinitely many
	// cells are missing, so the _-coverage rule fires instead of a
	// finite missing-cell list.
	if otherSlot > 0 {
		out = append(out, at(n.Line, fmt.Errorf("%s: value match without _ is not provably exhaustive (slot %d leaves an open string remainder)", owner, otherSlot)))
		return out
	}
	out = append(out, at(n.Line, fmt.Errorf("%s: non-exhaustive match, missing %s", owner, strings.Join(wits, "; "))))
	return out
}

// verifySinglePolicy enforces the historical arity-1 value rules
// verbatim: bool slots take exactly true+false with no wildcards or
// strings, string tables need _. Anything else keeps its legacy
// message and code. On success records the emit bit like the multi
// path, so the unified emitter reads one fact for every arity.
func verifySinglePolicy(n *Node, owner string, out []error, hasBool []bool, strLits [][]string) []error {
	bools := map[bool]bool{}
	hasStr, hasWild := false, false
	for _, a := range n.Arms {
		switch a.Pats[0].Kind {
		case "bool":
			bools[a.Pats[0].B] = true
		case "str":
			hasStr = true
		case "wild":
			hasWild = true
		}
	}
	if len(bools) > 0 && (len(bools) != 2 || hasStr || hasWild) {
		return append(out, at(n.Line, fmt.Errorf("%s: bool match must be exactly true+false", owner)))
	}
	if hasStr && !hasWild {
		return append(out, at(n.Line, fmt.Errorf("%s: value match without _ is not provably exhaustive", owner)))
	}
	domains, covers, _, mixed := valueCoverOf(n, hasBool, strLits)
	if mixed > 0 {
		// Unreachable through the legacy gates above (any bool
		// presence with a string mix fails the first rule), kept for
		// hand-built ASTs that bypass them.
		return append(out, at(n.Line, fmt.Errorf("%s: bool match must be exactly true+false", owner)))
	}
	n.analysis = &valueMatchAnalysis{emitFinalElse: residualInLast(covers[:len(covers)-1], covers[len(covers)-1], 1, domains)}
	return out
}

// sortedKeys orders string-keyed sets so multi-error reports are stable.
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func walkCalls(node *Node) []*Small {
	var out []*Small
	var walk func(n *Node)
	var walkSmall func(s *Small)
	walkSmall = func(s *Small) {
		if s == nil {
			return
		}
		if s.Kind == "call" {
			out = append(out, s)
		}
		for _, a := range s.Args {
			walkSmall(a.V)
		}
		if s.Kind == "binop" {
			walkSmall(s.L)
			walkSmall(s.R)
		}
		if s.Kind == "strlen" || s.Kind == "stridx" || s.Kind == "strslice" {
			walkSmall(s.L)
			walkSmall(s.R)
			walkSmall(s.Hi)
		}
		for _, it := range s.Items {
			walkSmall(it)
		}
		if s.Kind == "exchange" {
			walkSmall(s.Outcome)
		}
	}
	walk = func(n *Node) {
		if n == nil {
			return
		}
		if n.IsMatch {
			for _, s := range n.Scruts {
				walkSmall(s)
			}
			for _, a := range n.Arms {
				walk(a.Rhs)
			}
			return
		}
		walkSmall(n.Small)
	}
	walk(node)
	return out
}

var pinRe = regexp.MustCompile(`@\d+$`)
