package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// v10: ints emit as bigint (unbounded, exact) and decs as strings
// carrying canonical digits (exact via the $ailDec helpers below).
// The old number mapping was lossy (0.1+0.2) and is gone.
var tsBase = map[string]string{"str": "string", "int": "bigint", "bool": "boolean", "dec": "string"}

func tsType(t string) (string, error) {
	return tsTypeB(t, nil, nil)
}

// recordShapes indexes declared record fields by type name, first
// wins across modules, matching the checker and the evaluator.
func recordShapes(mods []*Module) map[string][][2]string {
	recs := map[string][][2]string{}
	for _, m := range mods {
		for _, d := range m.Decls {
			if td, ok := d.(*TypeDecl); ok {
				if _, seen := recs[td.Name]; !seen {
					recs[td.Name] = td.Fields
				}
			}
		}
	}
	return recs
}

// tsTypeB maps an ail annotation to TS, erasing brands to their
// underlying type. Branding is proof, not runtime: the emit forgets it.
// Declared record names map to their emitted TS type of the same name.
func tsTypeB(t string, brands map[string]string, recs map[string][][2]string) (string, error) {
	if out, ok := tsBase[t]; ok {
		return out, nil
	}
	if u, ok := brands[t]; ok {
		if out, ok := tsBase[u]; ok {
			return out, nil
		}
	}
	if _, ok := recs[t]; ok {
		return t, nil
	}
	return "", fmt.Errorf("cannot map ail type to TS: %s", t)
}

// tsTag is the outcome discriminator key in emitted TypeScript. It is
// unspellable in ail (identifiers match \w+, so $ never appears in
// source), keeping metadata disjoint from logical payload fields: a
// field named kind stays data, and the tag can never be overwritten by
// one. Payload fields keep their source names verbatim.
const tsTag = "$ail_kind"

// tsField renders one payload field as target storage. Every declared
// field becomes an own data property: ordinarily name: value, except
// __proto__, which an object literal would otherwise install as the
// prototype — the computed-key form defines a real own property.
func tsField(name, value string) string {
	if name == "__proto__" {
		return `["__proto__"]: ` + value
	}
	return name + ": " + value
}

// tsErrMember renders one error kind as a TS union member.
func tsErrMember(ed *ErrorDecl, brands map[string]string, recs map[string][][2]string) (string, error) {
	fs := ""
	for _, f := range ed.Fields {
		t, err := tsTypeB(f[1], brands, recs)
		if err != nil {
			return "", err
		}
		fs += "; " + f[0] + ": " + t
	}
	return fmt.Sprintf("{ %s: \"%s\"%s }", tsTag, ed.Name, fs), nil
}

// externUnion is the TS Result type of a foreign call: ok carrying the
// Ret record plus one member per declared emits kind. The host owns the
// implementation; this is the contract ailc proves against.
func externUnion(ex *ExternDecl, prog *Program) (string, error) {
	var td *TypeDecl
	for _, m := range prog.Modules {
		for _, d := range m.Decls {
			if t, ok := d.(*TypeDecl); ok && t.Name == ex.Ret {
				td = t
			}
		}
	}
	if td == nil {
		return "", fmt.Errorf("extern %s returns unknown type %s", ex.Name, ex.Ret)
	}
	fs := ""
	recs := recordShapes(prog.Modules)
	for _, f := range td.Fields {
		t, err := tsTypeB(f[1], prog.Brands, recs)
		if err != nil {
			return "", err
		}
		fs += "; " + f[0] + ": " + t
	}
	union := fmt.Sprintf("{ %s: \"ok\"%s }", tsTag, fs)
	for _, e := range ex.Emits {
		var ed *ErrorDecl
		for _, m := range prog.Modules {
			for _, d := range m.Decls {
				if er, ok := d.(*ErrorDecl); ok && er.Name == e {
					ed = er
				}
			}
		}
		if ed == nil {
			return "", fmt.Errorf("extern %s emits unknown error %s", ex.Name, e)
		}
		mem, err := tsErrMember(ed, prog.Brands, recs)
		if err != nil {
			return "", err
		}
		union += " | " + mem
	}
	return union, nil
}

// declaredOkShape builds a function's Ok payload shape from its declared
// return record, never from literals in the body or tests. The checker
// already proves every Ok construction against that record (unknown,
// missing, and mistyped fields fail before emission), so the declaration
// is the single source of types: changing evidence without changing the
// signature or body cannot change the emitted type, and computed or
// referenced payloads need no literal witness.
func declaredOkShape(fn *FnDecl, prog *Program) (map[string]string, error) {
	recs := recordShapes(prog.Modules)
	fields, ok := recs[fn.Ret]
	if !ok {
		return nil, fmt.Errorf("%s returns unknown type %s", fn.Name, fn.Ret)
	}
	shape := map[string]string{}
	for _, f := range fields {
		t, err := tsTypeB(f[1], prog.Brands, recs)
		if err != nil {
			return nil, err
		}
		shape[f[0]] = t
	}
	return shape, nil
}

// leafType reports the static type of a literal operand: kinds carry
// their own types, so emit can dispatch exactly even where the checker
// never annotated (direct unit calls). Refs need checker annotations.
func leafType(s *Small) string {
	switch s.Kind {
	case "int":
		return "int"
	case "dec":
		return "dec"
	case "str":
		return "str"
	case "bool":
		return "bool"
	}
	return ""
}

// childType resolves one operand's static type: the checker's T
// annotation first, literals second. Emit runs only after checkSem,
// so unknown here is a loud internal, never a silent default.
func childType(s *Small) string {
	if s == nil {
		return ""
	}
	if s.T != "" {
		return s.T
	}
	return leafType(s)
}

// isScalar reports whether an operand type compares exactly with
// native identity: base types (bigint by value, canonical dec
// strings, strings, booleans) and str-backed brands.
func (e *emitter) isScalar(ot string) bool {
	if _, ok := tsBase[ot]; ok {
		return true
	}
	return e.brands[ot] == "str"
}

// emitEquality lowers ==/!= by resolved operand type. Scalars and
// brands compare natively; records and error payloads compare
// field-by-field over their declared shape, ignoring the outcome
// envelope; cell wrappers compare their .value payloads. Unknown or
// unsupported operand types fail (AIL5005) instead of falling back to
// object identity.
func (e *emitter) emitEquality(op, ot, l, r string) (string, error) {
	ll, rr, shape := l, r, ot
	if strings.HasPrefix(shape, "cell:") {
		// Cell wrappers compare their .value payloads, never the
		// wrapper identity.
		ll += ".value"
		rr += ".value"
		shape = strings.TrimPrefix(shape, "cell:")
	}
	native := e.isScalar(shape)
	var expr string
	if native {
		if op == "==" {
			expr = fmt.Sprintf("(%s === %s)", ll, rr)
		} else {
			expr = fmt.Sprintf("(%s !== %s)", ll, rr)
		}
		return expr, nil
	}
	fields, err := e.equalityFields(shape)
	if err != nil {
		return "", err
	}
	var fs []string
	for _, f := range fields {
		fs = append(fs, strconv.Quote(f))
	}
	e.recEq = true
	expr = fmt.Sprintf("$ailEqRec(%s, %s, [%s])", ll, rr, strings.Join(fs, ", "))
	if op == "!=" {
		expr = "(!" + expr + ")"
	}
	return expr, nil
}

// equalityFields resolves the declared comparison shape for a record
// or error-payload operand: record names to their type fields, err:
// kinds to their error fields. Anything else is not structurally
// comparable (AIL5005).
func (e *emitter) equalityFields(ot string) ([]string, error) {
	if strings.HasPrefix(ot, "err:") {
		fs, ok := e.errFields[strings.TrimPrefix(ot, "err:")]
		if !ok {
			return nil, fmt.Errorf("cannot emit comparison over %s (%s)", ot, CodeBadCompare)
		}
		return fs, nil
	}
	if fs, ok := e.recs[ot]; ok {
		var names []string
		for _, f := range fs {
			names = append(names, f[0])
		}
		return names, nil
	}
	return nil, fmt.Errorf("cannot emit comparison over %s (%s)", ot, CodeBadCompare)
}

// binopOperandType reports the static ail type of a binop's operands
// for dispatch. Same-type operands are enforced by the checker, so one
// type describes both sides. Checker annotations win; literal kinds
// are the fallback. Empty means unknown: the caller fails loud.
func binopOperandType(node *Small) string {
	if node.T != "" {
		return node.T
	}
	if t := leafType(node.L); t != "" {
		return t
	}
	return leafType(node.R)
}

func (e *emitter) emitValue(node *Small) (string, error) {
	switch node.Kind {
	case "str":
		// One literal encoder (normStr): the in-memory string is
		// source-raw (the parser interprets no escapes), so the
		// target literal must re-escape every special char. Values
		// and patterns share this encoder (see stmtMatch).
		return normStr(node.Str), nil
	case "int":
		return node.Num.String() + "n", nil
	case "dec":
		return strconv.Quote(node.Dec), nil
	case "seal":
		// Sealed values emit as their string: brands erase. The
		// checker owns the one-string-operand rule (v25).
		if len(node.Args) != 1 {
			return "", fmt.Errorf("cannot emit seal: want one value")
		}
		return e.emitValue(node.Args[0].V)
	case "bool":
		if node.B {
			return "true", nil
		}
		return "false", nil
	case "ref":
		return strings.Join(node.Ref, "."), nil
	case "binop":
		l, err := e.emitValue(node.L)
		if err != nil {
			return "", err
		}
		r, err := e.emitValue(node.R)
		if err != nil {
			return "", err
		}
		ot := binopOperandType(node)
		if ot == "" {
			return "", fmt.Errorf("cannot emit %s: operand type unknown (run checkSem first)", node.Op)
		}
		// Equality lowers by resolved operand type, preserving the
		// evaluator's structural semantics: scalars (and brands,
		// erased to strings) are exact natively — bigint compares by
		// value, dec strings are canonical — while records, error
		// payloads, and cell wrappers compare field-by-field over
		// their declared shape, ignoring the outcome envelope.
		// Anything else fails instead of falling back to object
		// identity (AIL5005).
		if node.Op == "==" || node.Op == "!=" {
			eq, err := e.emitEquality(node.Op, ot, l, r)
			if err != nil {
				return "", err
			}
			return eq, nil
		}
		if ot == "dec" {
			// Decs are canonical-digit strings: route through the
			// exact $ailDec helpers, never native operators.
			switch node.Op {
			case "+":
				e.decOps["add"] = true
				return fmt.Sprintf("$ailDecAdd(%s, %s)", l, r), nil
			case "-":
				e.decOps["sub"] = true
				return fmt.Sprintf("$ailDecSub(%s, %s)", l, r), nil
			case "*":
				e.decOps["mul"] = true
				return fmt.Sprintf("$ailDecMul(%s, %s)", l, r), nil
			case ">=":
				e.decOps["ge"] = true
				return fmt.Sprintf("$ailDecGe(%s, %s)", l, r), nil
			case "<=":
				e.decOps["le"] = true
				return fmt.Sprintf("$ailDecLe(%s, %s)", l, r), nil
			case ">":
				e.decOps["gt"] = true
				return fmt.Sprintf("$ailDecGt(%s, %s)", l, r), nil
			case "<":
				e.decOps["lt"] = true
				return fmt.Sprintf("$ailDecLt(%s, %s)", l, r), nil
			}
			return "", fmt.Errorf("cannot emit op %s", node.Op)
		}
		// Strings and brands order by UTF-8 bytes through the
		// $ailStr helpers, matching Go's byte-wise evaluator
		// ordering (native target ordering is UTF-16 code units and
		// disagrees past the BMP); string + concatenates natively
		// and ==/!== stay exact in both runtimes.
		if ot == "str" || e.brands[ot] == "str" {
			switch node.Op {
			case ">=":
				e.strOps["ge"] = true
				return fmt.Sprintf("$ailStrGe(%s, %s)", l, r), nil
			case "<=":
				e.strOps["le"] = true
				return fmt.Sprintf("$ailStrLe(%s, %s)", l, r), nil
			case ">":
				e.strOps["gt"] = true
				return fmt.Sprintf("$ailStrGt(%s, %s)", l, r), nil
			case "<":
				e.strOps["lt"] = true
				return fmt.Sprintf("$ailStrLt(%s, %s)", l, r), nil
			}
		}
		// Ints are bigints (native ops exact, except / and %: BigInt
		// truncates toward zero, so Euclidean division rides the
		// $ailDivMod helper, emitted inline only when used).
		if ot == "int" && (node.Op == "/" || node.Op == "%") {
			e.divmod = true
			idx := "0"
			if node.Op == "%" {
				idx = "1"
			}
			return fmt.Sprintf("$ailDivMod(%s, %s)[%s]", l, r, idx), nil
		}
		ops := map[string]string{">=": ">=", "<=": "<=", ">": ">", "<": "<", "+": "+", "-": "-", "*": "*"}
		op, ok := ops[node.Op]
		if !ok {
			return "", fmt.Errorf("cannot emit op %s", node.Op)
		}
		return fmt.Sprintf("(%s %s %s)", l, op, r), nil
	case "strlen":
		v, err := e.emitValue(node.L)
		if err != nil {
			return "", err
		}
		if childType(node.L) != "str" {
			return "", fmt.Errorf("cannot emit #: operand type unknown (run checkSem first)")
		}
		return fmt.Sprintf("(BigInt([...%s].length))", v), nil
	case "stridx":
		b, err := e.emitValue(node.L)
		if err != nil {
			return "", err
		}
		ix, err := e.emitValue(node.R)
		if err != nil {
			return "", err
		}
		if childType(node.L) != "str" {
			return "", fmt.Errorf("cannot emit []: operand type unknown (run checkSem first)")
		}
		e.strOps["at"] = true
		return fmt.Sprintf("$ailStrAt(%s, %s)", b, ix), nil
	case "strslice":
		b, err := e.emitValue(node.L)
		if err != nil {
			return "", err
		}
		lo, err := e.emitValue(node.R)
		if err != nil {
			return "", err
		}
		hi, err := e.emitValue(node.Hi)
		if err != nil {
			return "", err
		}
		if childType(node.L) != "str" {
			return "", fmt.Errorf("cannot emit [:]: operand type unknown (run checkSem first)")
		}
		e.strOps["slice"] = true
		return fmt.Sprintf("$ailStrSlice(%s, %s, %s)", b, lo, hi), nil
	case "ctor":
		var parts []string
		for _, a := range node.Args {
			v, err := e.emitValue(a.V)
			if err != nil {
				return "", err
			}
			parts = append(parts, tsField(a.Name, v))
		}
		inner := strings.Join(parts, ", ")
		if node.Ctor == "Ok" {
			if inner == "" {
				return `{ ` + tsTag + `: "ok" }`, nil
			}
			return `{ ` + tsTag + `: "ok", ` + inner + ` }`, nil
		}
		if strings.Contains(node.Ctor, ".") {
			if inner == "" {
				return `{ ` + tsTag + `: "` + node.Ctor + `" }`, nil
			}
			return `{ ` + tsTag + `: "` + node.Ctor + `", ` + inner + ` }`, nil
		}
		// Declared records in value positions are plain data objects
		// with the checked field set: no outcome tag, since these
		// values never enter the match protocol. Unknown constructors
		// fail loud instead of mistagging as an error kind.
		if _, ok := e.recs[node.Ctor]; !ok {
			return "", fmt.Errorf("cannot emit unknown constructor %s", node.Ctor)
		}
		return `{ ` + inner + ` }`, nil
	case "call":
		// One binding rule (bindSlots): the call is invoked in the
		// resolved parameter order, not source order, so a reordered
		// named call means the same thing here as in evaluation.
		// Store ops are positional-only with a fixed spelling and keep
		// their source order.
		ordered := node.Args
		if !isStoreOp(node.Fname) {
			params, ok := e.params[node.Fname]
			if !ok {
				return "", fmt.Errorf("no params for callee %s", node.Fname)
			}
			slots, err := bindSlots(node.Fname, node.Args, params)
			if err != nil {
				return "", err
			}
			ordered = make([]Arg, len(node.Args))
			for i, a := range node.Args {
				ordered[slots[i]] = a
			}
		}
		var parts []string
		for _, a := range ordered {
			v, err := e.emitValue(a.V)
			if err != nil {
				return "", err
			}
			parts = append(parts, v)
		}
		return fmt.Sprintf("%s(%s)", node.Fname, strings.Join(parts, ", ")), nil
	}
	return "", fmt.Errorf("cannot emit: %s", node.Kind)
}

// decRuntimeShared implements the representation plumbing every dec
// operation needs: split canonical digits, renormalize, and convert
// between canonical strings and scaled BigInt mantissas. It mirrors
// parseDecParts/decArith/canonDec in eval.go exactly: add/sub align to
// the wider scale, mul sums scales, results renormalize (no trailing
// fractional zeros, -0 folds to 0.0).
var decRuntimeShared = []string{
	"function $ailDecSplit(d: string): { neg: boolean; ip: string; fp: string } {",
	"  let neg = false;",
	"  if (d.startsWith(\"-\")) {",
	"    neg = true;",
	"    d = d.slice(1);",
	"  }",
	"  const dot = d.indexOf(\".\");",
	"  return { neg, ip: d.slice(0, dot), fp: d.slice(dot + 1) };",
	"}",
	"function $ailDecNorm(ip: string, fp: string, neg: boolean): string {",
	"  ip = ip.replace(/^0+(?=\\d)/, \"\");",
	"  fp = fp.replace(/0+$/, \"\");",
	"  if (fp === \"\") {",
	"    fp = \"0\";",
	"  }",
	"  if (ip === \"0\" && fp === \"0\") {",
	"    return \"0.0\";",
	"  }",
	"  return (neg ? \"-\" : \"\") + ip + \".\" + fp;",
	"}",
	"function $ailDecMant(p: { neg: boolean; ip: string; fp: string }, scale: number): bigint {",
	"  let f = p.fp;",
	"  while (f.length < scale) {",
	"    f += \"0\";",
	"  }",
	"  const m = BigInt(p.ip + f);",
	"  return p.neg ? -m : m;",
	"}",
	"function $ailDecFromMant(m: bigint, scale: number): string {",
	"  let neg = false;",
	"  if (m < 0n) {",
	"    neg = true;",
	"    m = -m;",
	"  }",
	"  let digits = m.toString();",
	"  while (digits.length < scale + 1) {",
	"    digits = \"0\" + digits;",
	"  }",
	"  return $ailDecNorm(digits.slice(0, digits.length - scale), digits.slice(digits.length - scale), neg);",
	"}",
}

// decRuntimeOps holds one exact operation per ail operator, in fixed
// order for byte-stable output. Only used entries are emitted.
var decRuntimeOps = []struct {
	key  string
	code []string
}{
	{"add", []string{
		"function $ailDecAdd(a: string, b: string): string {",
		"  const A = $ailDecSplit(a);",
		"  const B = $ailDecSplit(b);",
		"  const s = Math.max(A.fp.length, B.fp.length);",
		"  return $ailDecFromMant($ailDecMant(A, s) + $ailDecMant(B, s), s);",
		"}",
	}},
	{"sub", []string{
		"function $ailDecSub(a: string, b: string): string {",
		"  const A = $ailDecSplit(a);",
		"  const B = $ailDecSplit(b);",
		"  const s = Math.max(A.fp.length, B.fp.length);",
		"  return $ailDecFromMant($ailDecMant(A, s) - $ailDecMant(B, s), s);",
		"}",
	}},
	{"mul", []string{
		"function $ailDecMul(a: string, b: string): string {",
		"  const A = $ailDecSplit(a);",
		"  const B = $ailDecSplit(b);",
		"  return $ailDecFromMant($ailDecMant(A, A.fp.length) * $ailDecMant(B, B.fp.length), A.fp.length + B.fp.length);",
		"}",
	}},
	{"ge", []string{
		"function $ailDecGe(a: string, b: string): boolean {",
		"  const A = $ailDecSplit(a);",
		"  const B = $ailDecSplit(b);",
		"  const s = Math.max(A.fp.length, B.fp.length);",
		"  return $ailDecMant(A, s) >= $ailDecMant(B, s);",
		"}",
	}},
	{"le", []string{
		"function $ailDecLe(a: string, b: string): boolean {",
		"  const A = $ailDecSplit(a);",
		"  const B = $ailDecSplit(b);",
		"  const s = Math.max(A.fp.length, B.fp.length);",
		"  return $ailDecMant(A, s) <= $ailDecMant(B, s);",
		"}",
	}},
	{"gt", []string{
		"function $ailDecGt(a: string, b: string): boolean {",
		"  const A = $ailDecSplit(a);",
		"  const B = $ailDecSplit(b);",
		"  const s = Math.max(A.fp.length, B.fp.length);",
		"  return $ailDecMant(A, s) > $ailDecMant(B, s);",
		"}",
	}},
	{"lt", []string{
		"function $ailDecLt(a: string, b: string): boolean {",
		"  const A = $ailDecSplit(a);",
		"  const B = $ailDecSplit(b);",
		"  const s = Math.max(A.fp.length, B.fp.length);",
		"  return $ailDecMant(A, s) < $ailDecMant(B, s);",
		"}",
	}},
	{"parts", []string{
		"function $ailDecParts(d: string): { coefficient: bigint; scale: bigint } {",
		"  const p = $ailDecSplit(d);",
		"  return { coefficient: $ailDecMant(p, p.fp.length), scale: BigInt(p.fp.length) };",
		"}",
	}},
}

// strRuntimeShared compares strings by UTF-8 bytes: Go orders strings
// byte-wise, while native target ordering is UTF-16 code units and
// disagrees past the BMP (U+E000 vs U+10000 sorts opposite). Encoding
// both sides to UTF-8 and comparing bytes reproduces the evaluator
// exactly on every input.
var strRuntimeShared = []string{
	"function $ailStrCmp(a: string, b: string): number {",
	"  const A = new TextEncoder().encode(a);",
	"  const B = new TextEncoder().encode(b);",
	"  const n = Math.min(A.length, B.length);",
	"  for (let i = 0; i < n; i++) {",
	"    if (A[i] !== B[i]) {",
	"      return A[i] < B[i] ? -1 : 1;",
	"    }",
	"  }",
	"  if (A.length === B.length) {",
	"    return 0;",
	"  }",
	"  return A.length < B.length ? -1 : 1;",
	"}",
}

var strRuntimeOps = []struct {
	key  string
	code []string
}{
	{"ge", []string{
		"function $ailStrGe(a: string, b: string): boolean {",
		"  return $ailStrCmp(a, b) >= 0;",
		"}",
	}},
	{"le", []string{
		"function $ailStrLe(a: string, b: string): boolean {",
		"  return $ailStrCmp(a, b) <= 0;",
		"}",
	}},
	{"at", []string{
		"function $ailStrAt(s: string, i: bigint): bigint {",
		"  const cps = [...s];",
		"  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error(\"str index out of range\");",
		"  const k = Number(i);",
		"  if (k >= cps.length) throw new Error(\"str index out of range\");",
		"  return BigInt(cps[k].codePointAt(0));",
		"}",
	}},
	{"slice", []string{
		"function $ailStrSlice(s: string, a: bigint, b: bigint): string {",
		"  const cps = [...s];",
		"  const toIdx = (x: bigint): number => { if (x < 0n || x > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error(\"str slice out of range\"); return Number(x); };",
		"  const lo = toIdx(a), hi = toIdx(b);",
		"  if (lo > hi || hi > cps.length) throw new Error(\"str slice out of range\");",
		"  return cps.slice(lo, hi).join(\"\");",
		"}",
	}},
	{"gt", []string{
		"function $ailStrGt(a: string, b: string): boolean {",
		"  return $ailStrCmp(a, b) > 0;",
		"}",
	}},
	{"lt", []string{
		"function $ailStrLt(a: string, b: string): boolean {",
		"  return $ailStrCmp(a, b) < 0;",
		"}",
	}},
}

// eqHelpers renders the structural equality runtime: $ailEqRec
// compares exactly the declared fields (ignoring the outcome
// envelope, requiring own presence on both sides, matching vEq),
// delegating per-field to $ailEqVal, which recurses into nested
// objects and compares scalars exactly (bigint by value, canonical
// dec strings and strings by content, booleans natively).
var eqHelpers = []string{
	"function $ailEqVal(x: any, y: any): boolean {",
	"  if (typeof x === \"object\" && x !== null && typeof y === \"object\" && y !== null) {",
	"    const kx = Object.keys(x);",
	"    if (kx.length !== Object.keys(y).length) {",
	"      return false;",
	"    }",
	"    for (const k of kx) {",
	"      if (!Object.prototype.hasOwnProperty.call(y, k)) {",
	"        return false;",
	"      }",
	"      if (!$ailEqVal((x as any)[k], (y as any)[k])) {",
	"        return false;",
	"      }",
	"    }",
	"    return true;",
	"  }",
	"  return x === y;",
	"}",
	"function $ailEqRec(a: any, b: any, fields: string[]): boolean {",
	"  for (const f of fields) {",
	"    if (!Object.prototype.hasOwnProperty.call(a, f)) {",
	"      return false;",
	"    }",
	"    if (!Object.prototype.hasOwnProperty.call(b, f)) {",
	"      return false;",
	"    }",
	"    if (!$ailEqVal((a as any)[f], (b as any)[f])) {",
	"      return false;",
	"    }",
	"  }",
	"  return true;",
	"}",
}

// strHelpers renders the byte-order string runtime for exactly the used
// comparisons, shared plumbing first, then ops in fixed order.
func strHelpers(used map[string]bool) []string {
	var out []string
	out = append(out, "// Byte-order string comparison: UTF-8 bytes, matching Go.")
	out = append(out, strRuntimeShared...)
	for _, op := range strRuntimeOps {
		if used[op.key] {
			out = append(out, op.code...)
		}
	}
	return out
}

// decHelpers renders the exact-decimal runtime for exactly the used
// operations, shared plumbing first, then ops in fixed order.
func decHelpers(used map[string]bool) []string {
	var out []string
	out = append(out, "// Exact-decimal runtime (v10): canonical-digit strings, BigInt math.")
	out = append(out, decRuntimeShared...)
	for _, op := range decRuntimeOps {
		if used[op.key] {
			out = append(out, op.code...)
		}
	}
	return out
}

type emitter struct {
	tmp       int
	fnUnions  map[string]string      // callee fn -> its TS Result type
	params    map[string][][2]string // callee fn/extern -> params in order
	brands    map[string]string      // brand -> underlying, for erasure
	cellTypes map[string]string      // cell -> TS type, this module only
	tailUnion string
	decOps    map[string]bool        // exact-decimal helpers used by this module
	strOps    map[string]bool        // byte-order string helpers used by this module
	recEq     bool                   // structural record comparison used by this module
	recs      map[string][][2]string // record name -> declared fields
	errFields map[string][]string    // error kind -> declared field names
	divmod    bool                   // Euclidean division helper used by this module
}

// divModHelper renders the Euclidean integer-division runtime: BigInt
// / and % truncate toward zero, so a negative truncated remainder is
// adjusted into 0 <= r < |b|. The helper text is pinned by
// TestDivEmitHelper and the Go-side contract by TestDivModVectors;
// executing the emitted helper against the contract awaits a node
// gate (tsc verification stays suspended per #11).
var divModHelper = []string{
	"// Euclidean integer division (v17): quotient and remainder with",
	"// 0 <= r < |b| on every sign combination.",
	"function $ailDivMod(a: bigint, b: bigint): [bigint, bigint] {",
	"  let q: bigint = a / b;",
	"  let r: bigint = a - b * q;",
	"  if (r < 0n) {",
	"    q += b > 0n ? -1n : 1n;",
	"    r = a - b * q;",
	"  }",
	"  return [q, r];",
	"}",
}

func (e *emitter) fresh() string {
	e.tmp++
	// Unspellable in ail (identifiers match \w+, so $ never appears in
	// source): generated temporaries can never collide with source
	// bindings such as a parameter named _m1.
	return fmt.Sprintf("$ail_m%d", e.tmp)
}

func (e *emitter) retLines(rhs *Node, env map[string]string) ([]string, error) {
	if rhs.IsMatch {
		var nested []string
		if err := e.stmtMatch(rhs, &nested); err != nil {
			return nil, err
		}
		return nested, nil
	}
	v, err := e.emitValue(rhs.Small)
	if err != nil {
		return nil, err
	}
	return []string{"return " + v + ";"}, nil
}

func indent(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		out = append(out, "  "+ln)
	}
	return out
}

func (e *emitter) stmtMatch(node *Node, out *[]string) error {
	if node.Kind != MatchCall {
		return e.emitValueMatch(node, out)
	}
	scrut := node.Scruts[0]
	if isStoreOp(scrut.Fname) {
		return e.stmtStoreOp(node, scrut, out)
	}
	if isDecParts(scrut.Fname) {
		return e.stmtDecParts(node, scrut, out)
	}
	union, ok := e.fnUnions[scrut.Fname]
	if !ok {
		return fmt.Errorf("no Result type for callee %s", scrut.Fname)
	}
	tmp := e.fresh()
	call, err := e.emitValue(scrut)
	if err != nil {
		return err
	}
	*out = append(*out, fmt.Sprintf("const %s: %s = %s;", tmp, union, call))
	*out = append(*out, fmt.Sprintf("switch (%s."+tsTag+") {", tmp))
	for _, arm := range node.Arms {
		pat := arm.Pats[0]
		// Every arm body is block-scoped: the same binder name in
		// two arms (or a binder shadowing an outer one) must not
		// collide, and success and error binders lower uniformly
		// as one const bound to the matched union value.
		switch {
		case pat.Kind == "variantWild":
			*out = append(*out, fmt.Sprintf("case \"%s\": {", pat.Name))
			lines, err := e.retLines(arm.Rhs, nil)
			if err != nil {
				return err
			}
			*out = append(*out, indent(lines)...)
			*out = append(*out, "}")
		case pat.Kind == "variant" && pat.Name != "Ok":
			*out = append(*out, fmt.Sprintf("case \"%s\": {", pat.Name))
			*out = append(*out, fmt.Sprintf("  const %s = %s;", pat.Var, tmp))
			lines, err := e.retLines(arm.Rhs, map[string]string{pat.Var: pat.Var})
			if err != nil {
				return err
			}
			*out = append(*out, indent(lines)...)
			*out = append(*out, "}")
		case pat.Kind == "variant" && pat.Name == "Ok":
			*out = append(*out, `case "ok": {`)
			*out = append(*out, fmt.Sprintf("  const %s = %s;", pat.Var, tmp))
			lines, err := e.retLines(arm.Rhs, map[string]string{pat.Var: pat.Var})
			if err != nil {
				return err
			}
			*out = append(*out, indent(lines)...)
			*out = append(*out, "}")
		default:
			return fmt.Errorf("call-match arm must be an error kind or Ok")
		}
	}
	*out = append(*out, "}")
	return nil
}

// emitValueMatch emits a value table of any arity 1..N: one arm loop
// for every shape. Scrutinee references preserve the legacy layout
// exactly — the single expression inline (no temporary, golden bytes),
// one $ail_mN temporary per slot past arity 1 — so the only
// arity-dependent choice is output formatting, never match semantics.
// A conditionless arm owns its whole residual: trailing lines mid-chain
// with later arms dead (the historical wild shape, kept for arity-1
// wild arms too), a bare else at the end. A proved-total final arm
// reads the checker's bit, never a recomputed proof; anything else is a
// tested condition. Exhaustiveness already proven; emit assumes it.
func (e *emitter) emitValueMatch(node *Node, out *[]string) error {
	nslot := len(node.Scruts)
	for _, arm := range node.Arms {
		if len(arm.Pats) != nslot {
			return fmt.Errorf("match arm has %d patterns; this match has %d scrutinees", len(arm.Pats), nslot)
		}
		for _, p := range arm.Pats {
			if k := p.Kind; k != "bool" && k != "str" && k != "wild" {
				return fmt.Errorf("variant pattern on a non-call match")
			}
		}
	}
	refs := make([]string, nslot)
	if nslot == 1 {
		sv, err := e.emitValue(node.Scruts[0])
		if err != nil {
			return err
		}
		refs[0] = sv
	} else {
		for i, s := range node.Scruts {
			v, err := e.emitValue(s)
			if err != nil {
				return err
			}
			t := e.fresh()
			*out = append(*out, fmt.Sprintf("const %s = %s;", t, v))
			refs[i] = t
		}
	}
	finalElse := false
	if node.analysis != nil {
		finalElse = node.analysis.emitFinalElse
	}
	known := make([]*bool, nslot)
	for n, arm := range node.Arms {
		last := n == len(node.Arms)-1
		conds, learn, contradiction := e.valueConds(arm, refs, known)
		kw := "if"
		if n > 0 {
			kw = "else if"
		}
		switch {
		case len(conds) == 0 && last && (nslot > 1 || !allWild(arm)):
			// Whole residual at the end: bare else. (Arity-1 wild
			// arms keep the historical trailing-lines shape below.)
			*out = append(*out, "else {")
		case len(conds) == 0:
			// Whole residual mid-chain with later arms dead (the
			// historical wild shape): trailing lines, chain ends.
			lines, err := e.retLines(arm.Rhs, nil)
			if err != nil {
				return err
			}
			*out = append(*out, lines...)
			return nil
		case last && finalElse:
			*out = append(*out, "else {")
		default:
			if contradiction {
				// Static duplicate past the proof: the raw
				// conjunction keeps the chain total; never reached.
				conds, _, _ = e.valueConds(arm, refs, nil)
			}
			*out = append(*out, fmt.Sprintf("%s (%s) {", kw, strings.Join(conds, " && ")))
			if learn >= 0 {
				nb := !arm.Pats[learn].B
				known[learn] = &nb
			}
		}
		lines, err := e.retLines(arm.Rhs, nil)
		if err != nil {
			return err
		}
		*out = append(*out, indent(lines)...)
		*out = append(*out, "}")
	}
	return nil
}

// allWild reports whether every slot pattern is a wildcard.
func allWild(arm Arm) bool {
	for _, p := range arm.Pats {
		if p.Kind != "wild" {
			return false
		}
	}
	return true
}

// valueConds renders one arm's test conjunction over the scrutinee
// references (inline expression or temporaries), dropping bool checks
// implied by failed earlier arms (known). Reports the slot when the
// conjunction is exactly one fresh bool test (its failure fixes that
// slot for later arms), and whether the arm contradicts known facts (a
// static duplicate: raw conditions still render, so the chain stays
// total).
func (e *emitter) valueConds(arm Arm, refs []string, known []*bool) (conds []string, learn int, contradiction bool) {
	learn = -1
	for i, p := range arm.Pats {
		switch p.Kind {
		case "wild":
		case "bool":
			if known != nil && known[i] != nil {
				if *known[i] == p.B {
					continue
				}
				contradiction = true
			}
			c := refs[i]
			if !p.B {
				c = "!(" + refs[i] + ")"
			}
			conds = append(conds, c)
			if known != nil && known[i] == nil && len(conds) == 1 {
				learn = i
			} else {
				learn = -1
			}
		case "str":
			conds = append(conds, fmt.Sprintf("%s === %s", refs[i], normStr(p.Str)))
			learn = -1
		default:
			// Variant slots never survive the proof; render nothing so
			// the arm reads total rather than crashing the emitter.
			learn = -1
		}
	}
	return conds, learn, contradiction
}

// stmtStoreOp emits a cell operation as plain module state access: a
// get reads the cell into an ok+value union, a put assigns then
// proceeds. No imports: shared numeric code is emitted inline (v10).
// Only Ok-variant arms are legal past the exhaustiveness gate.
func (e *emitter) stmtStoreOp(node *Node, scrut *Small, out *[]string) error {
	cell, ok := storeCellName(scrut)
	t, known := e.cellTypes[cell]
	if !ok || !known {
		return fmt.Errorf("unknown cell in %s", scrut.Fname)
	}
	tmp := e.fresh()
	if scrut.Fname == "state__get" {
		*out = append(*out, fmt.Sprintf("const %s: { "+tsTag+": \"ok\", value: %s } = { "+tsTag+": \"ok\", value: %s };", tmp, t, cell))
	} else {
		if len(scrut.Args) != 2 {
			return fmt.Errorf("cannot emit %s: want cell and value", scrut.Fname)
		}
		v, err := e.emitValue(scrut.Args[1].V)
		if err != nil {
			return err
		}
		*out = append(*out, fmt.Sprintf("%s = %s;", cell, v))
		*out = append(*out, fmt.Sprintf("const %s: { "+tsTag+": \"ok\" } = { "+tsTag+": \"ok\" };", tmp))
	}
	*out = append(*out, fmt.Sprintf("switch (%s."+tsTag+") {", tmp))
	for _, arm := range node.Arms {
		pat := arm.Pats[0]
		if pat.Kind != "variant" || pat.Name != "Ok" {
			return fmt.Errorf("store match arm must be Ok")
		}
		*out = append(*out, `case "ok": {`)
		*out = append(*out, fmt.Sprintf("  const %s = %s;", pat.Var, tmp))
		lines, err := e.retLines(arm.Rhs, map[string]string{pat.Var: pat.Var})
		if err != nil {
			return err
		}
		*out = append(*out, indent(lines)...)
		*out = append(*out, "}")
	}
	*out = append(*out, "}")
	return nil
}

// stmtDecParts emits a dec-observation match: the operand evaluates
// once through $ailDecParts into an ok-tagged pair, then the single
// Ok arm binds it like any other payload. Only Ok-variant arms are
// legal past the exhaustiveness gate.
func (e *emitter) stmtDecParts(node *Node, scrut *Small, out *[]string) error {
	if len(scrut.Args) != 1 {
		return fmt.Errorf("cannot emit %s: want one value", scrut.Fname)
	}
	v, err := e.emitValue(scrut.Args[0].V)
	if err != nil {
		return err
	}
	e.decOps["parts"] = true
	tmp := e.fresh()
	*out = append(*out, fmt.Sprintf("const %s: { "+tsTag+": \"ok\", coefficient: bigint, scale: bigint } = { "+tsTag+": \"ok\", ...$ailDecParts(%s) };", tmp, v))
	*out = append(*out, fmt.Sprintf("switch (%s."+tsTag+") {", tmp))
	for _, arm := range node.Arms {
		pat := arm.Pats[0]
		if pat.Kind != "variant" || pat.Name != "Ok" {
			return fmt.Errorf("dec__parts match arm must be Ok")
		}
		*out = append(*out, `case "ok": {`)
		*out = append(*out, fmt.Sprintf("  const %s = %s;", pat.Var, tmp))
		lines, err := e.retLines(arm.Rhs, map[string]string{pat.Var: pat.Var})
		if err != nil {
			return err
		}
		*out = append(*out, indent(lines)...)
		*out = append(*out, "}")
	}
	*out = append(*out, "}")
	return nil
}

func (e *emitter) fn(fn *FnDecl, union string) ([]string, error) {
	var params []string
	for _, p := range fn.Params {
		t, err := tsTypeB(p[1], e.brands, e.recs)
		if err != nil {
			return nil, err
		}
		params = append(params, p[0]+": "+t)
	}
	lines := []string{fmt.Sprintf("export function %s(%s): %s {", fn.Name, strings.Join(params, ", "), union)}
	var body []string
	if fn.Body.IsMatch {
		if err := e.stmtMatch(fn.Body, &body); err != nil {
			return nil, err
		}
	} else {
		v, err := e.emitValue(fn.Body.Small)
		if err != nil {
			return nil, err
		}
		body = append(body, "return "+v+";")
	}
	lines = append(lines, indent(body)...)
	lines = append(lines, "}")
	return lines, nil
}

func emitModule(mod *Module, prog *Program, stemOf, resultOfStem map[string]string, fnUnions map[string]string) (string, error) {
	var L []string
	L = append(L, fmt.Sprintf("// GENERATED from %s by ailc v0.0.0. DO NOT EDIT.", mod.File))
	L = append(L, "// Prod emit: tests + given stripped.")
	// Declared records resolve to their emitted TS type of the same
	// name throughout this module.
	recs := recordShapes(prog.Modules)
	need := map[string]map[string]bool{}
	for _, d := range mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, c := range walkCalls(fn.Body) {
			if prog.Uses[c.Fname] {
				if need[stemOf[c.Fname]] == nil {
					need[stemOf[c.Fname]] = map[string]bool{}
				}
				need[stemOf[c.Fname]][c.Fname] = true
			}
		}
	}
	var stems []string
	for s := range need {
		stems = append(stems, s)
	}
	sort.Strings(stems)
	for _, stem := range stems {
		if stem == mod.Stem {
			continue
		}
		var names []string
		for n := range need[stem] {
			names = append(names, n)
		}
		sort.Strings(names)
		names = append(names, "type "+resultOfStem[stem])
		L = append(L, fmt.Sprintf("import { %s } from \"./%s\";", strings.Join(names, ", "), stem))
	}
	// Foreign imports: every extern this module calls is a typed TS
	// import from its host-provided stub file. Sorted for stability.
	exCalled := map[string]bool{}
	for _, d := range mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, c := range walkCalls(fn.Body) {
			if prog.Externs[c.Fname] != nil {
				exCalled[c.Fname] = true
			}
		}
	}
	var exNames []string
	for n := range exCalled {
		exNames = append(exNames, n)
	}
	sort.Strings(exNames)
	if len(exNames) > 0 {
		L = append(L, fmt.Sprintf("import { %s } from \"./%s.externs\";", strings.Join(exNames, ", "), mod.Stem))
	}
	cap_ := capitalize(mod.Mod) + "Result"
	var members []string
	for _, d := range mod.Decls {
		ed, ok := d.(*ErrorDecl)
		if !ok {
			continue
		}
		mem, err := tsErrMember(ed, prog.Brands, recs)
		if err != nil {
			return "", err
		}
		members = append(members, mem)
	}
	// One ok member per distinct declared Ok shape, sorted for
	// stability. A single-shape module emits exactly the old union; a
	// file that mixes shapes (validators beside their consumers)
	// emits an ok union instead of failing the build. Shapes key on
	// names and types together, so value: bool and value: int stay
	// distinct members instead of collapsing by name.
	shapes := map[string]map[string]string{}
	for _, d := range mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		shape, err := declaredOkShape(fn, prog)
		if err != nil {
			return "", err
		}
		var keys []string
		for k := range shape {
			keys = append(keys, k+":"+shape[k])
		}
		sort.Strings(keys)
		key := strings.Join(keys, ",")
		if _, seen := shapes[key]; !seen {
			shapes[key] = shape
		}
	}
	var skeys []string
	for k := range shapes {
		skeys = append(skeys, k)
	}
	sort.Strings(skeys)
	var okMembers []string
	for _, k := range skeys {
		shape := shapes[k]
		var names []string
		for n := range shape {
			names = append(names, n)
		}
		sort.Strings(names)
		okFs := ""
		for _, n := range names {
			okFs += "; " + n + ": " + shape[n]
		}
		okMembers = append(okMembers, fmt.Sprintf("{ "+tsTag+": \"ok\"%s }", okFs))
	}
	union := strings.Join(okMembers, " | ")
	if union == "" {
		union = "{ " + tsTag + ": \"ok\" }"
	}
	for _, m := range members {
		union += " | " + m
	}
	L = append(L, fmt.Sprintf("export type %s = %s;", cap_, union))
	for _, d := range mod.Decls {
		td, ok := d.(*TypeDecl)
		if !ok {
			continue
		}
		var fs []string
		for _, f := range td.Fields {
			t, err := tsTypeB(f[1], prog.Brands, recs)
			if err != nil {
				return "", err
			}
			fs = append(fs, f[0]+": "+t)
		}
		if len(fs) == 0 {
			L = append(L, fmt.Sprintf("export type %s = {};", td.Name))
		} else {
			L = append(L, fmt.Sprintf("export type %s = { %s };", td.Name, strings.Join(fs, "; ")))
		}
	}
	// Module state: one mutable let per cell, initialized from the
	// decl literal. Tests prove per-scenario behavior from init;
	// prod shares the cell across calls (documented boundary).
	cellTypes := map[string]string{}
	params := map[string][][2]string{}
	for n, f := range prog.Fns {
		params[n] = f.Params
	}
	for n, ex := range prog.Externs {
		params[n] = ex.Params
	}
	em := &emitter{fnUnions: fnUnions, params: params, brands: prog.Brands, cellTypes: cellTypes, decOps: map[string]bool{}, strOps: map[string]bool{}, recs: recs, errFields: prog.Errors, divmod: false}
	for _, d := range mod.Decls {
		sd, ok := d.(*StateDecl)
		if !ok {
			continue
		}
		t, err := tsTypeB(sd.Type, prog.Brands, recs)
		if err != nil {
			return "", err
		}
		cellTypes[sd.Name] = t
		init, err := em.emitValue(sd.Init)
		if err != nil {
			return "", err
		}
		L = append(L, fmt.Sprintf("let %s: %s = %s;", sd.Name, t, init))
	}
	var fnLines []string
	for _, d := range mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		lines, err := em.fn(fn, cap_)
		if err != nil {
			return "", err
		}
		fnLines = append(fnLines, lines...)
	}
	// Exact-decimal runtime: emitted inline only when a dec operation
	// or ordering is used, so files without dec arithmetic gain no code.
	// Helpers carry $ prefixes, which ail naming (domain__verb) cannot
	// spell, so user code can never collide with them.
	if len(em.decOps) > 0 {
		L = append(L, decHelpers(em.decOps)...)
	}
	// Byte-order string runtime: emitted inline only when a string
	// ordering is used, so files without one gain no code.
	if len(em.strOps) > 0 {
		L = append(L, strHelpers(em.strOps)...)
	}
	// Structural equality runtime: emitted inline only when a record,
	// error-payload, or cell comparison is used.
	if em.recEq {
		L = append(L, eqHelpers...)
	}
	if em.divmod {
		L = append(L, divModHelper...)
	}
	L = append(L, fnLines...)
	return strings.Join(L, "\n") + "\n", nil
}
