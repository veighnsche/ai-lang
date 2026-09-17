// types.go: the first real type checker.
//
// Until now annotations were documentation the evaluator ignored.
// checkTypes enforces them: every value position has a declared type,
// brands are nominal, and nothing converts implicitly. A branded value
// flows only into a position annotated with the same brand, so the sink
// rule (no secrets into str fields, no cross-brand comparison) falls out
// of exact matching with no separate taint pass. Branding is proof, not
// runtime: the evaluator and emitter already treat sealed values as
// plain strings.
//
// Runs in checkSem after checkEmits, before test execution, shared by
// editor and CLI. Three entry points: checkTypes per function,
// checkExternSig per extern, checkBrandDecl per brand.
package main

import (
	"fmt"
	"strings"
)

// tycker carries per-function checking state: the world's declared
// record and error shapes plus the value names in scope.
type tycker struct {
	prog   *Program
	text   string
	fn     string
	recs   map[string][][2]string
	errs   map[string][][2]string
	brands map[string]bool
	// brandFiles maps brand name to declaring file (first wins).
	brandFiles map[string]string
	// exec is true inside function bodies (executable positions)
	// and false in tests and given rows (checked data positions).
	exec bool
	// cells maps visible cell names to base types (this file only).
	cells map[string]string
	out   []Diag
}

func newTycker(prog *Program, text, fn string) *tycker {
	c := &tycker{prog: prog, text: text, fn: fn,
		recs:       map[string][][2]string{},
		errs:       map[string][][2]string{},
		brands:     map[string]bool{},
		brandFiles: map[string]string{},
		cells:      map[string]string{}}
	for _, m := range prog.Modules {
		for _, d := range m.Decls {
			switch d := d.(type) {
			case *TypeDecl:
				if _, ok := c.recs[d.Name]; !ok {
					c.recs[d.Name] = d.Fields
				}
			case *ErrorDecl:
				if _, ok := c.errs[d.Name]; !ok {
					c.errs[d.Name] = d.Fields
				}
			case *BrandDecl:
				c.brands[d.Name] = true
				if _, ok := c.brandFiles[d.Name]; !ok {
					if f, ok := prog.BrandFile[d.Name]; ok {
						c.brandFiles[d.Name] = f
					}
				}
			case *StateDecl:
				// Cells resolve in the checking function's own file;
				// only base types enter (the decl rule owns the rest).
				if prog.FnFile[fn] != m.ID {
					continue
				}
				switch d.Type {
				case "str", "int", "bool", "dec":
					if _, ok := c.cells[d.Name]; !ok {
						c.cells[d.Name] = d.Type
					}
				}
			}
		}
	}
	return c
}

// knownType reports whether a name is a legal annotation: a base type,
// a declared record, or a declared brand. Error kinds are not values.
func (c *tycker) knownType(t string) bool {
	switch t {
	case "str", "int", "bool", "dec":
		return true
	}
	if _, ok := c.recs[t]; ok {
		return true
	}
	return c.brands[t]
}

type calleeSig struct {
	params [][2]string
	ret    string
}

func (c *tycker) callee(name string) *calleeSig {
	if f, ok := c.prog.Fns[name]; ok {
		return &calleeSig{f.Params, f.Ret}
	}
	if e, ok := c.prog.Externs[name]; ok {
		return &calleeSig{e.Params, e.Ret}
	}
	return nil
}

func (c *tycker) mismatch(line int, where, got, want, token string) {
	c.out = append(c.out, spanDiag(c.text, line, "error",
		fmt.Sprintf("%s: got %s, want %s", where, got, want), token, CodeTypeMismatch))
}

// resolveRef types a variable or field path. ok=false means uncheckable;
// the caller (value, which visits each node once) owns the diagnostic.
func (c *tycker) resolveRef(ref []string, env map[string]string) (string, bool) {
	if len(ref) == 0 {
		return "", false
	}
	t, ok := env[ref[0]]
	if !ok || t == "" {
		return "", false
	}
	if strings.HasPrefix(t, "cell:") {
		// Synthetic get payload: exactly .value of the cell type.
		if len(ref) == 2 && ref[1] == "value" {
			return strings.TrimPrefix(t, "cell:"), true
		}
		return "", false
	}
	if t == "parts" {
		// Synthetic observation payload: exactly .coefficient and
		// .scale, both int. Lowercase, so no declared type collides.
		if len(ref) == 2 && (ref[1] == "coefficient" || ref[1] == "scale") {
			return "int", true
		}
		return "", false
	}
	for _, f := range ref[1:] {
		var fields [][2]string
		if strings.HasPrefix(t, "err:") {
			fields = c.errs[strings.TrimPrefix(t, "err:")]
		} else {
			fields = c.recs[t]
			if fields == nil && !c.brands[t] && t != "str" && t != "int" && t != "bool" && t != "dec" {
				return "", false
			}
		}
		if fields == nil {
			return "", false
		}
		found := false
		for _, fd := range fields {
			if fd[0] == f {
				t = fd[1]
				found = true
				break
			}
		}
		if !found {
			return "", false
		}
	}
	return t, true
}

// typeOf is the pure half of checking: the static type of a Small, or
// ok=false when the position is unchecked (calls resolve through their
// own rule, stubs through the callee contract, wildcards are patterns).
func (c *tycker) typeOf(s *Small, env map[string]string) (string, bool) {
	switch s.Kind {
	case "str":
		return "str", true
	case "int":
		return "int", true
	case "bool":
		return "bool", true
	case "dec":
		return "dec", true
	case "seal":
		if !c.brands[s.Seal] {
			return "", false
		}
		return s.Seal, true
	case "ref":
		t, ok := c.resolveRef(s.Ref, env)
		if !ok {
			return "", false
		}
		if strings.HasPrefix(t, "err:") {
			return t, true
		}
		// A variable of undeclared type is unchecked here: the
		// param/field declaration owns the AIL6002, and comparing
		// through it would cascade one typo into many.
		if !c.knownType(t) {
			return "", false
		}
		return t, true
	case "ctor":
		// A named record constructor carries its record type outward
		// so outer positions (Ok fields, call arguments, test
		// expectations, equality) check the constructor identity, not
		// just the inner fields. Ok and error ctors resolve through
		// their own rules; unknown records stay silent here because
		// checkCtor owns the unknown-record diagnostic.
		if s.Ctor != "Ok" && !strings.Contains(s.Ctor, ".") {
			if _, ok := c.recs[s.Ctor]; ok {
				return s.Ctor, true
			}
		}
		return "", false
	case "strlen":
		return "int", true
	case "stridx":
		return "int", true
	case "strslice":
		return "str", true
	case "binop":
		if !isArith(s.Op) {
			return "bool", true
		}
		l, lok := c.typeOf(s.L, env)
		r, rok := c.typeOf(s.R, env)
		// v16: + concatenates strings (construction needs no
		// indexing); - and * stay numeric-only, and brands and
		// bools compute nothing even when both sides agree.
		if s.Op == "+" && lok && rok && l == "str" && r == "str" {
			return "str", true
		}
		// v17: decimal division and remainder have no exact result;
		// keep them untyped so parents stay silent and the value
		// rule reports the one refusal.
		if (s.Op == "/" || s.Op == "%") && lok && rok && l == "dec" && r == "dec" {
			return "", false
		}
		if !lok || !rok || l != r || (l != "int" && l != "dec") {
			return "", false
		}
		return l, true
	}
	return "", false
}

// isArith reports the computing operators: comparisons ask, these do.
// v17 adds / and % (exact Euclidean integer division); dec operands
// for either are refused in the value rule, not here.
func isArith(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/" || op == "%"
}

// arithVerb names the operator class for mismatch messages, so agents
// do not file arithmetic mistakes under comparison.
func arithVerb(op string) string {
	switch op {
	case "+":
		return "add"
	case "-":
		return "subtract"
	case "/":
		return "divide"
	case "%":
		return "modulo"
	default:
		return "multiply"
	}
}

// tokenOf picks the squiggle token for a value: the applied name for
// calls, the leaf for paths, the operator for comparisons.
func tokenOf(s *Small) string {
	switch s.Kind {
	case "call":
		return s.Fname
	case "ref":
		if len(s.Ref) > 0 {
			return s.Ref[len(s.Ref)-1]
		}
	case "binop":
		return s.Op
	case "strlen":
		return "#"
	case "stridx":
		return "[]"
	case "strslice":
		return "[:]"
	case "seal":
		return s.Seal
	case "ctor":
		return s.Ctor
	case "exchange":
		return "exchange"
	}
	return ""
}

// value checks one Small against its wanted type and recurses into
// children with their own wants. want="" means the position is
// want-free (arm roots producing errors, scrutinees): children are
// still checked structurally, so one broken row never hides the rest.
func (c *tycker) value(s *Small, want string, line int, env map[string]string, where string) {
	if s == nil {
		return
	}
	if s.Kind == "float" {
		msg := fmt.Sprintf("float %s has no spelling: write d\"%s\"", s.Str, s.Str)
		if strings.ContainsAny(s.Str, "eE") {
			msg = fmt.Sprintf("float %s has no spelling: write the value as d\"12.34\"", s.Str)
		}
		c.out = append(c.out, spanDiag(c.text, line, "error", msg, s.Str, CodeFloatLiteral))
		return
	}
	if s.Kind == "seal" {
		if !c.brands[s.Seal] {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("unknown brand %s in seal", s.Seal), s.Seal, CodeUnknownType))
			return
		}
		// v15: executable code mints only its own module's brands.
		// The declaring file owns every executable seal site (grep
		// seal is the audit); tests and given rows may name any
		// declared brand because they are checked data, not code.
		if c.exec {
			if owner, ok := c.brandFiles[s.Seal]; ok && owner != c.prog.FnFile[c.fn] {
				c.out = append(c.out, spanDiag(c.text, line, "error",
					fmt.Sprintf("seal %s in %s mints a brand declared in %s: bodies seal only their own module's brands", s.Seal, c.fn, owner), s.Seal, CodeSealForeign))
				return
			}
		}
		// v10: brands erase to strings at runtime; emit reads T.
		s.T = s.Seal
		if len(s.Args) != 1 || s.Args[0].V.Kind != "str" {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("seal %s takes a string literal", s.Seal), s.Seal, CodeTypeMismatch))
			return
		}
		if want != "" && want != s.Seal {
			c.mismatch(line, where, s.Seal, want, s.Seal)
		}
		return
	}
	if s.Kind == "ref" {
		t, ok := env[s.Ref[0]]
		if !ok {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("unbound name %s in %s", s.Ref[0], c.fn), s.Ref[0], CodeTypeMismatch))
			return
		}
		if t == "" {
			// Bound to an unknown contract (unknown callee or
			// return): the owner reports it; checking here would
			// cascade one broken line into many.
			return
		}
		if t == "empty-ok" {
			// Bound to a put's empty Ok: it carries no fields and
			// no value, so any use is a mistake at this line.
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("%s is the empty Ok of a put: it carries no fields", s.Ref[0]), s.Ref[0], CodeTypeMismatch))
			return
		}
		got, ok := c.resolveRef(s.Ref, env)
		if !ok {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("no field %s on %s", strings.Join(s.Ref[1:], "."), s.Ref[0]), tokenOf(s), CodeTypeMismatch))
			return
		}
		if want != "" && !strings.HasPrefix(got, "err:") && got != want {
			c.mismatch(line, where, got, want, tokenOf(s))
		}
		// v10: record the resolved type for emit's typed dispatch.
		s.T = got
		return
	}
	if want != "" {
		if got, ok := c.typeOf(s, env); ok && got != want {
			c.mismatch(line, where, got, want, tokenOf(s))
		}
	}
	switch s.Kind {
	case "binop":
		where := "comparison"
		if isArith(s.Op) {
			where = arithVerb(s.Op)
		}
		c.value(s.L, "", line, env, where)
		c.value(s.R, "", line, env, where)
		l, lok := c.typeOf(s.L, env)
		r, rok := c.typeOf(s.R, env)
		if !lok || !rok {
			return
		}
		// v10: record the operand type for emit's typed dispatch.
		// Same-type operands are enforced below; emit runs only on
		// success, so the annotation always agrees there.
		if l == r {
			s.T = l
		}
		if isArith(s.Op) {
			// Arithmetic yields the operand type, but only int
			// and dec compute, plus str under + (v16: explicit
			// construction). Same-brand seals and bools do not,
			// even when both sides agree; neither do - and * on
			// strings.
			if l == r && l == "str" && s.Op == "+" {
				return
			}
			// v17: 1/3 does not terminate, so decimal division and
			// remainder are refused per operation. Integers divide
			// exactly (Euclidean); use them.
			if l == r && l == "dec" && (s.Op == "/" || s.Op == "%") {
				c.out = append(c.out, spanDiag(c.text, line, "error",
					fmt.Sprintf("dec %s has no exact result: divide integers, not decimals", s.Op), s.Op, CodeInexactDivision))
				return
			}
			if l != r || (l != "int" && l != "dec") {
				c.out = append(c.out, spanDiag(c.text, line, "error",
					fmt.Sprintf("cannot %s %s with %s: no implicit conversions", arithVerb(s.Op), l, r), s.Op, CodeTypeMismatch))
			}
			return
		}
		if l != r {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("cannot compare %s with %s: no implicit conversions", l, r), s.Op, CodeTypeMismatch))
		}
	case "strlen":
		c.value(s.L, "", line, env, "length")
		if t, ok := c.typeOf(s.L, env); ok && t != "str" {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("cannot count scalars of %s: length needs str", t), "#", CodeTypeMismatch))
			return
		}
		s.T = "int"
	case "stridx":
		c.value(s.L, "", line, env, "index base")
		c.value(s.R, "", line, env, "index")
		if t, ok := c.typeOf(s.L, env); ok && t != "str" {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("cannot index into %s: base must be str", t), "[]", CodeTypeMismatch))
			return
		}
		if t, ok := c.typeOf(s.R, env); ok && t != "int" {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("cannot index with %s: index must be int", t), "[]", CodeTypeMismatch))
			return
		}
		s.T = "int"
	case "strslice":
		c.value(s.L, "", line, env, "slice base")
		c.value(s.R, "", line, env, "slice start")
		c.value(s.Hi, "", line, env, "slice end")
		if t, ok := c.typeOf(s.L, env); ok && t != "str" {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("cannot slice %s: base must be str", t), "[:]", CodeTypeMismatch))
			return
		}
		for _, b := range []*Small{s.R, s.Hi} {
			if t, ok := c.typeOf(b, env); ok && t != "int" {
				c.out = append(c.out, spanDiag(c.text, line, "error",
					fmt.Sprintf("cannot slice with %s: bounds must be int", t), "[:]", CodeTypeMismatch))
				return
			}
		}
		s.T = "str"
	case "call":
		if isStoreOp(s.Fname) {
			c.checkStoreOp(s, line, env)
			return
		}
		if isDecParts(s.Fname) {
			c.checkDecParts(s, line, env)
			return
		}
		sig := c.callee(s.Fname)
		if sig == nil {
			return
		}
		slots, berr := bindSlots(s.Fname, s.Args, sig.params)
		if berr != nil {
			c.out = append(c.out, spanDiag(c.text, line, "error", berr.Error(), s.Fname, CodeBadBinding))
			// Still check the argument expressions themselves so one
			// bad vector never hides nested errors inside the args.
			for _, a := range s.Args {
				c.value(a.V, "", line, env, "call "+s.Fname+" arg")
			}
			return
		}
		for i, a := range s.Args {
			p := sig.params[slots[i]]
			want, label := p[1], "call "+s.Fname+" arg "+p[0]
			c.value(a.V, "", line, env, label)
			if c.knownType(want) {
				if got, ok := c.typeOf(a.V, env); ok && got != want {
					c.mismatch(line, label, got, want, tokenOf(a.V))
				}
			}
		}
	case "ctor":
		c.checkCtor(s, want, line, env, where)
	case "list":
		for _, it := range s.Items {
			c.value(it, "", line, env, where)
		}
	default:
		for _, a := range s.Args {
			c.value(a.V, "", line, env, where)
		}
	}
}

// checkStoreOp validates a store call against its cell: positional
// args only (one spelling), the cell first, and a put value of
// exactly the cell type. Unknown cells belong to checkEffects and
// stay silent here, like unknown callees.
func (c *tycker) checkStoreOp(s *Small, line int, env map[string]string) {
	for _, a := range s.Args {
		if a.HasName {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("call %s takes positional args", s.Fname), s.Fname, CodeTypeMismatch))
			return
		}
	}
	cell, ok := storeCellName(s)
	if !ok {
		c.out = append(c.out, spanDiag(c.text, line, "error",
			fmt.Sprintf("call %s takes the cell first", s.Fname), s.Fname, CodeTypeMismatch))
		return
	}
	want, known := c.cells[cell]
	if !known {
		return
	}
	count := 1
	if s.Fname == "state__put" {
		count = 2
	}
	if len(s.Args) != count {
		c.out = append(c.out, spanDiag(c.text, line, "error",
			fmt.Sprintf("call %s takes %d args", s.Fname, count), s.Fname, CodeTypeMismatch))
		return
	}
	if s.Fname == "state__put" {
		label := fmt.Sprintf("call %s value", s.Fname)
		c.value(s.Args[1].V, "", line, env, label)
		if got, ok := c.typeOf(s.Args[1].V, env); ok && got != want {
			c.mismatch(line, label, got, want, tokenOf(s.Args[1].V))
		}
	}
}

// checkDecParts validates the decimal observation kernel: positional
// args only (one spelling), exactly one arg, and a dec operand.
// Anything else is AIL6003, the operand-rule family.
func (c *tycker) checkDecParts(s *Small, line int, env map[string]string) {
	for _, a := range s.Args {
		if a.HasName {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("call %s takes positional args", s.Fname), s.Fname, CodeTypeMismatch))
			return
		}
	}
	if len(s.Args) != 1 {
		c.out = append(c.out, spanDiag(c.text, line, "error",
			fmt.Sprintf("call %s takes 1 arg", s.Fname), s.Fname, CodeTypeMismatch))
		return
	}
	label := fmt.Sprintf("call %s value", s.Fname)
	c.value(s.Args[0].V, "", line, env, label)
	if got, ok := c.typeOf(s.Args[0].V, env); ok && got != "dec" {
		c.mismatch(line, label, got, "dec", tokenOf(s.Args[0].V))
	}
}

// nodeStoreArms types a store-op match: a get payload binds its var
// to the cell type through the synthetic value field; a put yields
// empty Ok, so its var binds nothing checkable. The cell type flows
// from the declaration, never from inference.
func (c *tycker) nodeStoreArms(n *Node, env map[string]string, want string) {
	cell, ok := storeCellName(n.Scrut)
	t, known := "", false
	if ok {
		t, known = c.cells[cell]
	}
	for _, a := range n.Arms {
		env2 := map[string]string{}
		for k, v := range env {
			env2[k] = v
		}
		armWant := ""
		if a.Pat.Kind == "variant" && a.Pat.Name == "Ok" && a.Pat.Var != "" {
			if n.Scrut.Fname == "state__get" && known {
				env2[a.Pat.Var] = "cell:" + t
			} else if n.Scrut.Fname == "state__put" {
				env2[a.Pat.Var] = "empty-ok"
			} else {
				env2[a.Pat.Var] = ""
			}
			armWant = want
		}
		if a.Pat.Kind == "variantWild" && a.Pat.Name == "Ok" {
			armWant = want
		}
		c.node(a.Rhs, env2, armWant)
	}
}

// nodePartsArms types a dec-observation match: the Ok payload binds
// its var to the synthetic parts shape, whose only fields are the
// int coefficient and scale. The shape is fixed by the kernel, never
// by inference, so no declared record is consulted.
func (c *tycker) nodePartsArms(n *Node, env map[string]string, want string) {
	for _, a := range n.Arms {
		env2 := map[string]string{}
		for k, v := range env {
			env2[k] = v
		}
		armWant := ""
		if a.Pat.Kind == "variant" && a.Pat.Name == "Ok" && a.Pat.Var != "" {
			env2[a.Pat.Var] = "parts"
			armWant = want
		}
		if a.Pat.Kind == "variantWild" && a.Pat.Name == "Ok" {
			armWant = want
		}
		c.node(a.Rhs, env2, armWant)
	}
}

// checkCtor validates a construction against its contract: Ok against
// the wanted record, dotted errors against their ErrorDecl, named
// records against their TypeDecl. want="" (error-producing positions)
// still checks error fields; unknown contracts belong to other codes,
// so undeclared kinds are skipped, never double-reported.
func (c *tycker) checkCtor(s *Small, want string, line int, env map[string]string, where string) {
	name := s.Ctor
	var fields [][2]string
	label := where
	if name == "Ok" {
		if want == "" || strings.HasPrefix(want, "err:") {
			for _, a := range s.Args {
				c.value(a.V, "", line, env, where)
			}
			return
		}
		rec, ok := c.recs[want]
		if !ok {
			return
		}
		fields, label = rec, "Ok"
	} else if strings.Contains(name, ".") {
		ed, ok := c.errs[name]
		if !ok {
			for _, a := range s.Args {
				c.value(a.V, "", line, env, where)
			}
			return
		}
		fields, label = ed, name
	} else {
		rec, ok := c.recs[name]
		if !ok {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("unknown record %s in construction", name), name, CodeUnknownType))
			return
		}
		fields, label = rec, name
		if want != "" && want != name {
			c.mismatch(line, where, name, want, name)
		}
	}
	byName := map[string]string{}
	for _, f := range fields {
		byName[f[0]] = f[1]
	}
	seenArg := map[string]bool{}
	for _, a := range s.Args {
		if seenArg[a.Name] {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("%s repeats field %s", label, a.Name), a.Name, CodeTypeMismatch))
		}
		seenArg[a.Name] = true
		ft, ok := byName[a.Name]
		if !ok {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("%s has no field %s", label, a.Name), a.Name, CodeTypeMismatch))
			c.value(a.V, "", line, env, where)
			continue
		}
		flabel := fmt.Sprintf("%s field %s", label, a.Name)
		c.value(a.V, "", line, env, flabel)
		if c.knownType(ft) {
			if got, ok := c.typeOf(a.V, env); ok && got != ft {
				c.mismatch(line, flabel, got, ft, a.Name)
			}
		}
	}
	for _, f := range fields {
		found := false
		for _, a := range s.Args {
			if a.Name == f[0] {
				found = true
				break
			}
		}
		if !found {
			c.out = append(c.out, spanDiag(c.text, line, "error",
				fmt.Sprintf("%s is missing field %s", label, f[0]), f[0], CodeTypeMismatch))
		}
	}
}

// node walks a body threading the wanted type to final-value positions:
// Ok arms and value matches produce the function's return, error arms
// produce errors (checked against their decl, want-free). Bindings copy
// the map down; shadowing stays undiagnosed (open question).
func (c *tycker) node(n *Node, env map[string]string, want string) {
	if n == nil {
		return
	}
	if !n.IsMatch {
		c.value(n.Small, want, n.Line, env, "returns")
		return
	}
	c.value(n.Scrut, "", n.Line, env, "match scrutinee")
	if n.Scrut != nil && n.Scrut.Kind == "call" {
		if isStoreOp(n.Scrut.Fname) {
			c.nodeStoreArms(n, env, want)
			return
		}
		if isDecParts(n.Scrut.Fname) {
			c.nodePartsArms(n, env, want)
			return
		}
		sig := c.callee(n.Scrut.Fname)
		for _, a := range n.Arms {
			env2 := map[string]string{}
			for k, v := range env {
				env2[k] = v
			}
			armWant := ""
			if a.Pat.Kind == "variant" && a.Pat.Var != "" {
				if a.Pat.Name == "Ok" {
					if sig != nil && c.knownType(sig.ret) {
						env2[a.Pat.Var] = sig.ret
					} else {
						env2[a.Pat.Var] = ""
					}
					armWant = want
				} else if _, ok := c.errs[a.Pat.Name]; ok {
					env2[a.Pat.Var] = "err:" + a.Pat.Name
				} else {
					env2[a.Pat.Var] = ""
				}
			}
			if a.Pat.Kind == "variantWild" && a.Pat.Name == "Ok" {
				armWant = want
			}
			c.node(a.Rhs, env2, armWant)
		}
		return
	}
	for _, a := range n.Arms {
		c.node(a.Rhs, env, want)
	}
}

// checkStubs validates given-table outcomes against the callee
// contract: Ok shapes against its Ret, error fields against their
// decl. Kinds were already checked; undeclared kinds are skipped here.
// Stubs evaluate in the test's argument environment, so param names
// resolve here exactly as they do at runtime.
func (c *tycker) checkStubs(fn *FnDecl, text string, env map[string]string) {
	for _, m := range matchNodes(fn.Body) {
		if m.Scrut == nil || m.Scrut.Kind != "call" || m.Given == nil {
			continue
		}
		sig := c.callee(m.Scrut.Fname)
		if sig == nil {
			continue
		}
		for key, sm := range m.Given {
			if sm == nil {
				continue
			}
			line := locateLineFrom(text, key+" =>", m.Line, m.Line)
			where := fmt.Sprintf("script %s for %s", key, m.Scrut.Fname)
			var items []*Small
			if sm.Kind == "list" {
				items = sm.Items
			} else {
				items = []*Small{sm}
			}
			for _, it := range items {
				if it.Kind != "exchange" {
					continue // check.go owns the row shape
				}
				if it.Outcome.Kind == "ctor" && it.Outcome.Ctor == "Ok" {
					c.checkCtor(it.Outcome, sig.ret, line, env, where)
				} else {
					c.value(it.Outcome, "", line, env, where)
				}
				// Expected args type against callee params, like test
				// args type against function params. Names themselves
				// are proven dynamically at each hit.
				for _, a := range it.Args {
					c.value(a.V, "", line, env, "exchange arg "+a.Name+" for "+m.Scrut.Fname)
					want := ""
					for _, p := range sig.params {
						if p[0] == a.Name {
							want = p[1]
						}
					}
					if want != "" && c.knownType(want) {
						if got, ok := c.typeOf(a.V, env); ok && got != want {
							c.mismatch(line, "exchange arg "+a.Name+" for "+m.Scrut.Fname, got, want, tokenOf(a.V))
						}
					}
				}
			}
		}
	}
}

// checkTypes enforces annotations on one function: params and ret name
// known types, tests match the signature, stubs match the callee, and
// the body checks internally with the return threaded to Ok positions.
func checkTypes(fn *FnDecl, prog *Program, text string) []Diag {
	c := newTycker(prog, text, fn.Name)
	for _, p := range fn.Params {
		if !c.knownType(p[1]) {
			c.out = append(c.out, spanDiag(text, fn.Line, "error",
				fmt.Sprintf("unknown type %s in param %s", p[1], p[0]), p[0], CodeUnknownType))
		}
	}
	if !c.knownType(fn.Ret) {
		c.out = append(c.out, spanDiag(text, fn.Line, "error",
			fmt.Sprintf("unknown type %s in returns", fn.Ret), fn.Ret, CodeUnknownType))
	}
	env := map[string]string{}
	for _, p := range fn.Params {
		env[p[0]] = p[1]
	}
	for _, t := range fn.Tests {
		for _, a := range t.Args {
			want := ""
			for _, p := range fn.Params {
				if p[0] == a.Name {
					want = p[1]
				}
			}
			label := fmt.Sprintf("test %s arg %s", t.Name, a.Name)
			c.value(a.V, "", t.Line, env, label)
			if want != "" && c.knownType(want) {
				if got, ok := c.typeOf(a.V, env); ok && got != want {
					c.mismatch(t.Line, label, got, want, tokenOf(a.V))
				}
			}
		}
		exp := t.Expected
		label := fmt.Sprintf("test %s expects", t.Name)
		switch {
		case exp.Kind == "ctor" && exp.Ctor == "Ok":
			if c.knownType(fn.Ret) {
				c.checkCtor(exp, fn.Ret, t.Line, env, label)
			} else {
				c.value(exp, "", t.Line, env, label)
			}
		case exp.Kind == "ctor" && strings.Contains(exp.Ctor, "."):
			// Complete error expectations carry the payload and
			// check against the error decl, like constructions.
			if _, ok := c.errs[exp.Ctor]; !ok {
				c.out = append(c.out, spanDiag(text, t.Line, "error",
					fmt.Sprintf("unknown error kind %s in %s", exp.Ctor, label), exp.Ctor, CodeUnknownType))
			} else {
				c.checkCtor(exp, "", t.Line, env, label)
			}
		case exp.Kind == "ref" && len(exp.Ref) > 1:
			// Bare error kinds prove nothing about the payload (v12):
			// expectations must construct the complete error value.
			c.out = append(c.out, spanDiag(text, t.Line, "error",
				fmt.Sprintf("test %s expects bare error kind %s: write the complete error value", t.Name, strings.Join(exp.Ref, ".")), exp.Ref[len(exp.Ref)-1], CodeBareErrorKind))
		case exp.Kind == "float":
			c.value(exp, "", t.Line, env, label)
		default:
			c.out = append(c.out, spanDiag(text, t.Line, "error",
				fmt.Sprintf("%s %s: write Ok(...) or an error kind", label, exp.Kind), t.Name, CodeTypeMismatch))
		}
	}
	c.checkStubs(fn, text, env)
	ret := ""
	if c.knownType(fn.Ret) {
		ret = fn.Ret
	}
	// Bodies are executable positions: the seal rule applies from
	// here on. Everything above checked data (tests, scripts).
	c.exec = true
	c.node(fn.Body, env, ret)
	return c.out
}

// checkExternSig validates a foreign import's contract: params and ret
// name known types, and ret is a record (foreign outcomes script as
// Ok/errors only; bool-returning calls are unscriptable in v0).
func checkExternSig(ex *ExternDecl, prog *Program, text string) []Diag {
	c := newTycker(prog, text, ex.Name)
	for _, p := range ex.Params {
		if !c.knownType(p[1]) {
			c.out = append(c.out, spanDiag(text, ex.Line, "error",
				fmt.Sprintf("unknown type %s in param %s", p[1], p[0]), p[0], CodeUnknownType))
		}
	}
	if !c.knownType(ex.Ret) {
		c.out = append(c.out, spanDiag(text, ex.Line, "error",
			fmt.Sprintf("unknown type %s in returns", ex.Ret), ex.Ret, CodeUnknownType))
	} else if _, ok := c.recs[ex.Ret]; !ok {
		// Comma-ok, not nil: an empty record (Verdict) is declared
		// and legal; only non-records are rejected.
		c.out = append(c.out, spanDiag(text, ex.Line, "error",
			fmt.Sprintf("extern %s returns %s: externs return a record type, Ok/errors script the outcome", ex.Name, ex.Ret), ex.Ret, CodeTypeMismatch))
	}
	// v12: extern manifests are upper bounds like function emits —
	// every entry must name a declared error.
	for _, e := range ex.Emits {
		if _, ok := c.errs[e]; !ok {
			c.out = append(c.out, spanDiag(text, ex.Line, "error",
				fmt.Sprintf("extern %s declares unknown error kind %s in emits", ex.Name, e), e, CodeUnknownKind))
		}
	}
	return c.out
}

// checkDeclFields validates record and error field annotations: every
// field names a known type. Uses of an undeclared field type stay
// silent (typeOf suppresses them), so the declaration owns the error.
func checkDeclFields(name string, fields [][2]string, line int, prog *Program, text string) []Diag {
	c := newTycker(prog, text, name)
	var out []Diag
	for _, f := range fields {
		if !c.knownType(f[1]) {
			out = append(out, spanDiag(text, line, "error",
				fmt.Sprintf("unknown type %s in field %s", f[1], f[0]), f[0], CodeUnknownType))
		}
	}
	return out
}

// checkBrandDecl enforces the v0 brand boundary: string-backed only.
// int-backed brands wait for a second underlying type with something
// to prove about it.
func checkBrandDecl(b *BrandDecl, text string) []Diag {
	if b.Under != "str" {
		return []Diag{spanDiag(text, b.Line, "error",
			fmt.Sprintf("brand %s wraps %s: v0 brands wrap str only", b.Name, b.Under), b.Under, CodeTypeMismatch)}
	}
	return nil
}
