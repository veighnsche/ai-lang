package main

// a80: contract admission for the verifier first cut. A contracted
// function is admitted only when its complete contract and body lie
// inside the declared proof fragment: mathematical integers,
// Booleans, and finite acyclic records of those sorts, with total
// bodies free of externs, kernels, state, and recursion. Anything
// else fails closed here, before any solver runs, and one
// unsupported clause rejects the complete contract: supported
// clauses are never proved while the rest is silently omitted.
//
// Pure classifier, not wired into any pipeline yet: activation
// ("verifier checking is active") waits for the proving slices.
// Uncontracted functions keep ordinary status and report nothing.

import (
	"fmt"
	"strings"
)

// admitSort is a proof-fragment sort: int, bool, a nominal record
// with its field shapes, or outside (with the reason named).
type admitSort struct {
	kind   string // "int", "bool", "rec", "outside"
	rec    string // nominal record name when kind == "rec"
	why    string // human reason when kind == "outside"
	fields map[string]admitSort // record field shapes (error payloads synthesize these)
}

// admitTerm is the outcome of classifying one Small node:
// ok with a sort, unknown (malformed: names nothing resolvable),
// mismatch (ill-sorted combination), outside (names an unsupported
// sort/operator), or call (needs the call rules, not the term rules).
type admitTerm struct {
	sort admitSort
	flag string // "", "unknown", "mismatch", "outside", "call"
	what string // offending detail for messages
}

// admission indexes the declarations one CheckContractAdmission run
// needs: functions by qualified name, externs, nominal sorts.
type admission struct {
	prog     *Program
	texts    map[string]string
	fns      map[string]*FnDecl
	extern   map[string]bool
	types    map[string]*TypeDecl
	errors   map[string]*ErrorDecl
	brands   map[string]bool
	variants map[string]bool
	out      []Diag
	// verified names callees a proving run has established in
	// this run (verdict position 4): their summaries are usable
	// assumptions. Empty during standalone admission, so every
	// callee summary is unavailable there.
	verified map[string]bool
}

func indexAdmission(prog *Program, texts map[string]string) *admission {
	a := &admission{
		prog:     prog,
		texts:    texts,
		fns:      map[string]*FnDecl{},
		extern:   map[string]bool{},
		types:    map[string]*TypeDecl{},
		errors:   map[string]*ErrorDecl{},
		brands:   map[string]bool{},
		variants: map[string]bool{},
		verified: map[string]bool{},
	}
	for _, m := range prog.Modules {
		for _, d := range m.Decls {
			switch t := d.(type) {
			case *FnDecl:
				a.fns[t.Name] = t
			case *ExternDecl:
				a.extern[t.Name] = true
			case *TypeDecl:
				a.types[t.Name] = t
			case *ErrorDecl:
				a.errors[t.Name] = t
			case *BrandDecl:
				a.brands[t.Name] = true
			case *VariantDecl:
				a.variants[t.Name] = true
			}
		}
	}
	return a
}

func isContracted(fn *FnDecl) bool {
	return len(fn.Requires) > 0 || len(fn.Ensures) > 0
}

// CheckContractAdmission classifies every contracted function in the
// program: malformed contracts (AIL4301), fragment-external sorts,
// operators, and body dependencies (AIL4302), and reliance on
// callee summaries no proving run has established (AIL4303).
func CheckContractAdmission(prog *Program, texts map[string]string) []Diag {
	a := indexAdmission(prog, texts)
	for _, m := range prog.Modules {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok || !isContracted(fn) {
				continue
			}
			a.checkFn(m, fn)
		}
	}
	a.checkRecursion()
	return a.out
}

func (a *admission) textOf(m *Module) string {
	if t, ok := a.texts[m.ID]; ok {
		return t
	}
	return ""
}

// emit records one finding with a71 payloads: what the rule wanted,
// what it saw, and the fix direction.
func (a *admission) emit(text string, line int, token, code, msg, expected, found, hint string) {
	d := spanDiag(text, line, "error", msg, token, code)
	d.Expected = expected
	d.Found = found
	d.Hint = hint
	a.out = append(a.out, d)
}

func (a *admission) ownerOf(fn *FnDecl) *Module {
	for _, m := range a.prog.Modules {
		for _, d := range m.Decls {
			if f, ok := d.(*FnDecl); ok && f == fn {
				return m
			}
		}
	}
	return nil
}

// sortOfType maps a declared type name into the fragment.
func (a *admission) sortOfType(t string) admitSort {
	switch t {
	case "int":
		return admitSort{kind: "int"}
	case "bool":
		return admitSort{kind: "bool"}
	case "str", "dec", "float", "Bytes":
		return admitSort{kind: "outside", why: t}
	}
	if strings.HasPrefix(t, "Seq<") {
		return admitSort{kind: "outside", why: t}
	}
	if a.brands[t] {
		return admitSort{kind: "outside", why: "brand " + t}
	}
	if a.variants[t] {
		return admitSort{kind: "outside", why: "variant " + t}
	}
	if td, ok := a.types[t]; ok {
		s := admitSort{kind: "rec", rec: t, fields: map[string]admitSort{}}
		for _, f := range td.Fields {
			s.fields[f[0]] = a.sortOfType(f[1])
		}
		return s
	}
	return admitSort{kind: "outside", why: "unknown type " + t}
}

// errorShape synthesizes a record shape from an error declaration's
// payload fields so err.field projections resolve.
func (a *admission) errorShape(name string) (admitSort, bool) {
	ed, ok := a.errors[name]
	if !ok {
		return admitSort{}, false
	}
	s := admitSort{kind: "rec", rec: name, fields: map[string]admitSort{}}
	for _, f := range ed.Fields {
		s.fields[f[0]] = a.sortOfType(f[1])
	}
	return s, true
}

// linearInt reports whether s is a linear integer term: literals,
// projections, addition/subtraction of linear terms, and
// multiplication with a literal coefficient on one side.
func (a *admission) linearInt(s *Small, scope map[string]admitSort) bool {
	t := a.term(s, scope)
	if t.flag != "" || t.sort.kind != "int" {
		return false
	}
	switch s.Kind {
	case "int", "ref":
		return true
	case "binop":
		switch s.Op {
		case "+", "-":
			return a.linearInt(s.L, scope) && a.linearInt(s.R, scope)
		case "*":
			if isIntLit(s.L) {
				return a.linearInt(s.R, scope)
			}
			if isIntLit(s.R) {
				return a.linearInt(s.L, scope)
			}
			return false
		}
	}
	return false
}

func isIntLit(s *Small) bool {
	return s != nil && s.Kind == "int"
}

// term classifies one Small node without deciding the diagnostic:
// callers map unknown to malformed, outside to unsupported, and
// mismatch to malformed.
func (a *admission) term(s *Small, scope map[string]admitSort) admitTerm {
	if s == nil {
		return admitTerm{flag: "unknown", what: "missing expression"}
	}
	switch s.Kind {
	case "int":
		return admitTerm{sort: admitSort{kind: "int"}}
	case "bool":
		return admitTerm{sort: admitSort{kind: "bool"}}
	case "str", "dec", "float":
		return admitTerm{flag: "outside", what: s.Kind + " literal"}
	case "wild":
		return admitTerm{flag: "unknown", what: "wildcard"}
	case "call":
		return admitTerm{flag: "call", what: s.Fname}
	case "ctor", "seal", "exchange", "strlen", "stridx", "strslice", "seqlit", "list":
		return admitTerm{flag: "outside", what: s.Kind + " operation"}
	case "ref":
		return a.projection(s, scope)
	case "binop":
		return a.binop(s, scope)
	}
	return admitTerm{flag: "outside", what: "unknown expression kind " + s.Kind}
}

// projection resolves a reference path through record shapes: every
// segment must name an actual field, and every step must stay
// inside admitted sorts.
func (a *admission) projection(s *Small, scope map[string]admitSort) admitTerm {
	if len(s.Ref) == 0 {
		return admitTerm{flag: "unknown", what: "empty reference"}
	}
	cur, ok := scope[s.Ref[0]]
	if !ok {
		return admitTerm{flag: "unknown", what: "unknown name " + s.Ref[0]}
	}
	if cur.kind == "outside" {
		return admitTerm{flag: "outside", what: cur.why}
	}
	path := s.Ref[0]
	for _, seg := range s.Ref[1:] {
		path += "." + seg
		if cur.kind != "rec" {
			return admitTerm{flag: "unknown", what: path + " projects into non-record sort"}
		}
		next, ok := cur.fields[seg]
		if !ok {
			return admitTerm{flag: "unknown", what: "unknown field " + path}
		}
		if next.kind == "outside" {
			return admitTerm{flag: "outside", what: next.why + " at " + path}
		}
		cur = next
	}
	return admitTerm{sort: cur}
}

func (a *admission) binop(s *Small, scope map[string]admitSort) admitTerm {
	l := a.term(s.L, scope)
	r := a.term(s.R, scope)
	for _, t := range []admitTerm{l, r} {
		if t.flag == "unknown" || t.flag == "mismatch" {
			return admitTerm{flag: "unknown", what: t.what}
		}
		if t.flag == "outside" {
			return t
		}
		if t.flag == "call" {
			return t
		}
	}
	ints := admitSort{kind: "int"}
	switch s.Op {
	case "+", "-":
		if l.sort.kind == "int" && r.sort.kind == "int" {
			return admitTerm{sort: ints}
		}
		return admitTerm{flag: "mismatch", what: describeSmall(s)}
	case "*":
		if l.sort.kind == "int" && r.sort.kind == "int" {
			if isIntLit(s.L) || isIntLit(s.R) {
				return admitTerm{sort: ints}
			}
			// Both sides integer but neither a literal
			// coefficient (x * x, (x+1) * (y+2)): nonlinear.
			return admitTerm{flag: "outside", what: "variable-by-variable multiplication"}
		}
		return admitTerm{flag: "mismatch", what: describeSmall(s)}
	case "/", "%":
		if l.sort.kind == "int" && r.sort.kind == "int" {
			return admitTerm{flag: "outside", what: s.Op + " has no integer-division theory in this cut"}
		}
		return admitTerm{flag: "mismatch", what: describeSmall(s)}
	case "==", "!=":
		if l.sort.kind == "int" && r.sort.kind == "int" {
			return admitTerm{sort: admitSort{kind: "bool"}}
		}
		if l.sort.kind == "bool" && r.sort.kind == "bool" {
			return admitTerm{sort: admitSort{kind: "bool"}}
		}
		if l.sort.kind == "rec" && r.sort.kind == "rec" {
			if l.sort.rec != r.sort.rec {
				return admitTerm{flag: "mismatch", what: describeSmall(s)}
			}
			if bad := admittedRecordFields(l.sort); bad != "" {
				return admitTerm{flag: "outside", what: bad}
			}
			return admitTerm{sort: admitSort{kind: "bool"}}
		}
		return admitTerm{flag: "mismatch", what: describeSmall(s)}
	case "<", "<=", ">", ">=":
		if l.sort.kind == "int" && r.sort.kind == "int" {
			return admitTerm{sort: admitSort{kind: "bool"}}
		}
		return admitTerm{flag: "mismatch", what: describeSmall(s)}
	}
	return admitTerm{flag: "outside", what: "unknown operator " + s.Op}
}

// admittedRecordFields names the first fragment-external field of a
// record shape, recursing through nested records.
func admittedRecordFields(s admitSort) string {
	if s.kind != "rec" {
		return ""
	}
	for name, f := range s.fields {
		if f.kind == "outside" {
			return f.why + " in field " + name
		}
		if f.kind == "rec" {
			if bad := admittedRecordFields(f); bad != "" {
				return bad
			}
		}
	}
	return ""
}

// describeSmall renders one Small node for Found payloads.
func describeSmall(s *Small) string {
	if s == nil {
		return "missing expression"
	}
	switch s.Kind {
	case "ref":
		return strings.Join(s.Ref, ".")
	case "binop":
		return describeSmall(s.L) + " " + s.Op + " " + describeSmall(s.R)
	case "call":
		return "call " + s.Fname
	case "int":
		if s.Num != nil {
			return s.Num.String()
		}
		return "int literal"
	case "bool":
		if s.B {
			return "true"
		}
		return "false"
	case "str":
		return "string literal"
	default:
		return s.Kind + " expression"
	}
}

// checkFn runs every admission check for one contracted function:
// annotation completeness, requires predicates, ensures arms, and
// the executable body.
func (a *admission) checkFn(m *Module, fn *FnDecl) {
	text := a.textOf(m)
	name, line := declNameLine(fn)
	scope := map[string]admitSort{}
	for _, p := range fn.Params {
		st := a.sortOfType(p[1])
		if st.kind == "outside" {
			a.emit(text, line, p[0], CodeContractUnsupported,
				fmt.Sprintf("%s takes %s of unsupported proof sort (%s): the first cut quantifies over integers, Booleans, and records of those sorts only", name, p[0], st.why),
				"params of sort int, bool, or records of int/bool",
				fmt.Sprintf("%s: %s", p[0], p[1]),
				"express the input with supported sorts, or defer this contract-bearing interface")
			continue
		}
		scope[p[0]] = st
	}
	a.checkCompleteness(m, fn, name, line, text)
	for _, r := range fn.Requires {
		a.checkPred(m, fn, name, line, text, scope, "requires", r)
	}
	// An empty admitted domain proves vacuously: fail closed.
	for _, r := range fn.Requires {
		if r != nil && r.Kind == "bool" && !r.B {
			a.emit(text, line, name, CodeContractUnsupported,
				fmt.Sprintf("%s requires false: no input is admitted, so verification would prove vacuously", name),
				"a requires clause admitting at least one input",
				"requires false",
				"admit the intended inputs, or defer this contract-bearing interface")
			break
		}
	}
	for _, arm := range fn.Ensures {
		a.checkArm(m, fn, name, text, scope, arm)
	}
	if len(fn.Effects) > 0 {
		a.emit(text, line, name, CodeContractUnsupported,
			fmt.Sprintf("%s declares effects %s: state is outside the verified closure", name, strings.Join(fn.Effects, ", ")),
			"a total body with no cell capabilities",
			strings.Join(fn.Effects, ", "),
			"verify a total, effect-free formulation, or defer this contract-bearing interface")
	}
	a.checkBody(m, fn, name, text, fn.Body, line)
}

// checkCompleteness enforces verdict 3A: exactly one ensures arm
// for Ok and every declared emits kind.
func (a *admission) checkCompleteness(m *Module, fn *FnDecl, name string, line int, text string) {
	needed := append([]string{"Ok"}, fn.Emits...)
	seen := map[string]int{}
	for _, arm := range fn.Ensures {
		seen[arm.Outcome]++
		if seen[arm.Outcome] == 2 {
			a.emit(text, arm.Line, arm.Outcome, CodeContractMalformed,
				fmt.Sprintf("%s has two ensures arms for %s: one outcome, one arm", name, arm.Outcome),
				"exactly one ensures arm per outcome",
				"duplicate on "+arm.Outcome,
				"merge the predicates into a single arm")
		}
	}
	for _, want := range needed {
		if seen[want] == 0 {
			a.emit(text, line, name, CodeContractMalformed,
				fmt.Sprintf("%s declares %s but has no ensures arm for it: a missing outcome is never true", name, want),
				"exactly one ensures arm for Ok and every emits kind",
				"missing on "+want,
				"write the arm, or remove the outcome from the interface through reviewed revision changes")
		}
	}
}

// checkPred enforces one Boolean predicate over the scope: unknown
// names and ill-sorted combinations are malformed, fragment-external
// sorts and operators are unsupported, calls are not predicates.
func (a *admission) checkPred(m *Module, fn *FnDecl, name string, line int, text string, scope map[string]admitSort, where string, s *Small) {
	t := a.term(s, scope)
	found := describeSmall(s)
	switch t.flag {
	case "unknown":
		a.emit(text, line, name, CodeContractMalformed,
			fmt.Sprintf("%s %s names the unresolvable %s: predicates resolve over the arm binder and the function params through actual fields", name, where, t.what),
			"a Boolean predicate over known binders and fields",
			found,
			"correct the name, or project through a declared field")
	case "mismatch":
		a.emit(text, line, name, CodeContractMalformed,
			fmt.Sprintf("%s %s is ill-sorted: %s", name, where, found),
			"sort-correct predicates: integer arithmetic and comparisons, Boolean terms",
			found,
			"correct the operand sorts")
	case "outside":
		// A predicate position demands a Boolean expression: a
		// never-Boolean shape (constructor, seal, literal of the
		// wrong sort) is a malformed case, while a Boolean-shaped
		// comparison or reference reaching outside the fragment
		// is unsupported.
		if s.Kind != "binop" && s.Kind != "ref" {
			a.emit(text, line, name, CodeContractMalformed,
				fmt.Sprintf("%s %s is not a Boolean predicate: %s", name, where, found),
				"a Boolean predicate over the binder and params",
				found,
				"write a Boolean comparison, input, projection, or case expression")
			break
		}
		a.emit(text, line, name, CodeContractUnsupported,
			fmt.Sprintf("%s %s uses %s, outside the first-cut fragment", name, where, t.what),
			"integer, Boolean, and record-of-those terms",
			found,
			"use supported logic, or defer this contract-bearing interface")
	case "call":
		a.emit(text, line, name, CodeContractUnsupported,
			fmt.Sprintf("%s %s calls %s: calls are not predicates", name, where, t.what),
			"call-free predicates",
			found,
			"state the property without calling")
	default:
		if t.sort.kind != "bool" {
			a.emit(text, line, name, CodeContractMalformed,
				fmt.Sprintf("%s %s is not Boolean: %s", name, where, found),
				"a Boolean predicate",
				found,
				"write a Boolean predicate over the binder and params")
		}
	}
}

// checkArm checks one ensures arm: the binder scope plus every
// predicate row and Boolean match block (conjoined).
func (a *admission) checkArm(m *Module, fn *FnDecl, name, text string, scope map[string]admitSort, arm ContractArm) {
	ascope := map[string]admitSort{}
	for k, v := range scope {
		ascope[k] = v
	}
	if arm.Outcome == "Ok" {
		ascope[arm.Bind] = a.sortOfType(fn.Ret)
	} else if shape, ok := a.errorShape(arm.Outcome); ok {
		ascope[arm.Bind] = shape
	} else {
		// The parser rejects unknown outcomes; this is unreachable
		// on clean programs, kept so admission never trusts it.
		a.emit(text, arm.Line, arm.Outcome, CodeContractMalformed,
			fmt.Sprintf("%s has an ensures arm for unknown outcome %s", name, arm.Outcome),
			"arms for Ok or a declared emits kind",
			"on "+arm.Outcome,
			"correct the outcome name")
		ascope[arm.Bind] = admitSort{kind: "outside", why: "unknown outcome " + arm.Outcome}
	}
	for _, p := range arm.Preds {
		a.checkPred(m, fn, name, arm.Line, text, ascope, "ensures on "+arm.Outcome, p)
	}
	for _, mt := range arm.Matches {
		a.checkContractMatch(m, fn, name, text, ascope, arm, mt)
	}
}

// checkContractMatch enforces verdict position 2: the scrutinee has
// Boolean sort, exactly one true arm and one false arm, each arm a
// Boolean expression or a nested supported match. These are logical
// expressions, not handler tables: no calls, no given, no omitted
// branch.
func (a *admission) checkContractMatch(m *Module, fn *FnDecl, name, text string, scope map[string]admitSort, arm ContractArm, node *Node) {
	where := "ensures on " + arm.Outcome + " match"
	line := arm.Line
	if node != nil && node.Line > 0 {
		line = node.Line
	}
	if node == nil || !node.IsMatch {
		a.emit(text, line, arm.Outcome, CodeContractMalformed,
			fmt.Sprintf("%s %s is not a match block", name, where),
			"a Boolean match with true and false arms",
			"missing match",
			"write the Boolean case expression")
		return
	}
	if node.Kind == MatchCall {
		a.emit(text, line, arm.Outcome, CodeContractUnsupported,
			fmt.Sprintf("%s %s matches a call outcome: calls are not predicates", name, where),
			"a Boolean scrutinee with true and false arms",
			"match call ...",
			"state the property without calling")
		return
	}
	if len(node.Scruts) != 1 {
		a.emit(text, line, arm.Outcome, CodeContractMalformed,
			fmt.Sprintf("%s %s matches %d scrutinees: one Boolean scrutinee", name, where, len(node.Scruts)),
			"exactly one scrutinee",
			fmt.Sprintf("%d scrutinees", len(node.Scruts)),
			"match a single Boolean expression")
		return
	}
	st := a.term(node.Scruts[0], scope)
	switch {
	case st.flag == "call":
		a.emit(text, line, arm.Outcome, CodeContractUnsupported,
			fmt.Sprintf("%s %s scrutinizes a call: calls are not predicates", name, where),
			"a call-free Boolean scrutinee",
			describeSmall(node.Scruts[0]),
			"state the property without calling")
		return
	case st.flag == "outside":
		a.emit(text, line, arm.Outcome, CodeContractUnsupported,
			fmt.Sprintf("%s %s scrutinizes %s, outside the fragment", name, where, st.what),
			"a Boolean scrutinee",
			describeSmall(node.Scruts[0]),
			"use supported logic")
		return
	case st.flag != "" || st.sort.kind != "bool":
		a.emit(text, line, arm.Outcome, CodeContractMalformed,
			fmt.Sprintf("%s %s scrutinee is not Boolean: %s", name, where, describeSmall(node.Scruts[0])),
			"a Boolean scrutinee",
			describeSmall(node.Scruts[0]),
			"match a Boolean comparison, input, or projection")
		return
	}
	if len(node.Given) > 0 {
		a.emit(text, line, arm.Outcome, CodeContractMalformed,
			fmt.Sprintf("%s %s carries a given table: specification matches take no scripts", name, where),
			"bare true and false arms",
			"given table present",
			"remove the given table")
	}
	seen := map[bool]bool{}
	for _, armNode := range node.Arms {
		if len(armNode.Pats) != 1 || armNode.Pats[0].Kind != "bool" {
			a.emit(text, armNode.Line, arm.Outcome, CodeContractMalformed,
				fmt.Sprintf("%s %s has a non-Boolean arm pattern: exactly one true arm and one false arm", name, where),
				"true and false arms only",
				"non-Boolean pattern",
				"write exactly true => and false => arms")
			continue
		}
		b := armNode.Pats[0].B
		if seen[b] {
			a.emit(text, armNode.Line, arm.Outcome, CodeContractMalformed,
				fmt.Sprintf("%s %s repeats its %v arm: one true arm, one false arm", name, where, b),
				"exactly one true arm and one false arm",
				"duplicate arm",
				"merge the duplicate arm")
			continue
		}
		seen[b] = true
		a.checkMatchArmBool(m, fn, name, text, scope, arm, where, armNode)
	}
	for _, b := range []bool{true, false} {
		if !seen[b] {
			a.emit(text, line, arm.Outcome, CodeContractMalformed,
				fmt.Sprintf("%s %s omits its %v arm: structural completeness, not witnesses", name, where, b),
				"exactly one true arm and one false arm",
				fmt.Sprintf("missing %v arm", b),
				fmt.Sprintf("add the %v => arm", b))
		}
	}
}

// checkMatchArmBool enforces Boolean arm results: a Boolean term or
// a nested supported Boolean match.
func (a *admission) checkMatchArmBool(m *Module, fn *FnDecl, name, text string, scope map[string]admitSort, arm ContractArm, where string, armNode Arm) {
	rhs := armNode.Rhs
	line := armNode.Line
	if rhs == nil {
		a.emit(text, line, arm.Outcome, CodeContractMalformed,
			fmt.Sprintf("%s %s has an empty arm: every arm returns a Boolean predicate", name, where),
			"Boolean arm results",
			"empty arm",
			"write the arm predicate")
		return
	}
	if rhs.IsMatch {
		a.checkContractMatch(m, fn, name, text, scope, arm, rhs)
		return
	}
	if rhs.Small == nil {
		a.emit(text, line, arm.Outcome, CodeContractMalformed,
			fmt.Sprintf("%s %s arm returns no predicate", name, where),
			"Boolean arm results",
			"missing predicate",
			"write the arm predicate")
		return
	}
	a.checkPred(m, fn, name, line, text, scope, where+" arm", rhs.Small)
}

// checkBody walks the executable body: every call resolves through
// the dependency rules, every value stays inside admitted sorts.
// Given tables hold scripts, not proof facts, and are ignored.
func (a *admission) checkBody(m *Module, fn *FnDecl, name, text string, node *Node, line int) {
	if node == nil {
		return
	}
	if node.Line > 0 {
		line = node.Line
	}
	if node.IsMatch {
		if node.Kind == MatchCall {
			if len(node.Scruts) == 1 && node.Scruts[0] != nil && node.Scruts[0].Kind == "call" {
				a.callRules(m, fn, name, text, node.Scruts[0].Fname, line)
			} else {
				a.emit(text, line, name, CodeContractMalformed,
					fmt.Sprintf("%s body has a call match with no call scrutinee", name),
					"match call f(args) with on arms",
					"missing call",
					"write the call scrutinee")
			}
			for _, armNode := range node.Arms {
				a.checkBody(m, fn, name, text, armNode.Rhs, armNode.Line)
			}
			return
		}
		if len(node.Scruts) != 1 {
			a.emit(text, line, name, CodeContractUnsupported,
				fmt.Sprintf("%s body matches %d scrutinees: the first cut proves single Boolean cases only", name, len(node.Scruts)),
				"a single Boolean scrutinee",
				fmt.Sprintf("%d scrutinees", len(node.Scruts)),
				"match one Boolean expression")
			return
		}
		if st := a.term(node.Scruts[0], a.bodyScope(fn)); st.flag != "" || st.sort.kind != "bool" {
			a.emit(text, line, name, CodeContractUnsupported,
				fmt.Sprintf("%s body matches %s: value cases prove over Booleans only", name, describeSmall(node.Scruts[0])),
				"a Boolean scrutinee",
				describeSmall(node.Scruts[0]),
				"match a Boolean expression, or defer this contract-bearing interface")
			return
		}
		seen := map[bool]bool{}
		for _, armNode := range node.Arms {
			if len(armNode.Pats) == 1 && armNode.Pats[0].Kind == "bool" {
				seen[armNode.Pats[0].B] = true
			}
			a.checkBody(m, fn, name, text, armNode.Rhs, armNode.Line)
		}
		if !seen[true] || !seen[false] {
			a.emit(text, line, name, CodeContractUnsupported,
				fmt.Sprintf("%s body has a Boolean match without both arms: the prover needs the complete case", name),
				"true and false arms",
				"incomplete Boolean case",
				"write both arms")
		}
		return
	}
	if node.Small == nil {
		return
	}
	a.checkBodySmall(m, fn, name, text, node.Small, line)
}

// bodyScope resolves body references to params only: match binders
// introduce fresh symbolic values the admission walk does not track
// by name, so any other head is unresolvable here. Binders flow
// through structurally (their constructor payloads are checked),
// never by name.
func (a *admission) bodyScope(fn *FnDecl) map[string]admitSort {
	scope := map[string]admitSort{}
	for _, p := range fn.Params {
		scope[p[0]] = a.sortOfType(p[1])
	}
	return scope
}

// checkBodySmall enforces admitted terms in value position: calls
// go through the dependency rules, constructors check their
// payloads, everything else must be an admitted term.
func (a *admission) checkBodySmall(m *Module, fn *FnDecl, name, text string, s *Small, line int) {
	switch s.Kind {
	case "call":
		a.callRules(m, fn, name, text, s.Fname, line)
		return
	case "ctor":
		for _, arg := range s.Args {
			if arg.V == nil {
				continue
			}
			if arg.V.Kind == "call" {
				a.callRules(m, fn, name, text, arg.V.Fname, line)
				continue
			}
			if arg.V.Kind == "ctor" {
				a.checkBodySmall(m, fn, name, text, arg.V, line)
				continue
			}
			t := a.term(arg.V, a.bodyScope(fn))
			switch t.flag {
			case "unknown":
				// Constructor payloads reference match binders
				// (on Ok r => Ok(value = r.value)), which the
				// admission walk does not track by name: the
				// structural proof (a later slice) follows them.
				// Only flag payloads that resolve to something
				// inadmissible, never bare binder paths.
			case "mismatch":
				a.emit(text, line, name, CodeContractMalformed,
					fmt.Sprintf("%s body constructs an ill-sorted payload: %s", name, describeSmall(arg.V)),
					"sort-correct payloads",
					describeSmall(arg.V),
					"correct the operand sorts")
			case "outside":
				a.emit(text, line, name, CodeContractUnsupported,
					fmt.Sprintf("%s body constructs %s, outside the fragment", name, t.what),
					"integer, Boolean, and record payloads",
					describeSmall(arg.V),
					"use supported payloads, or defer this contract-bearing interface")
			case "call":
				a.callRules(m, fn, name, text, t.what, line)
			}
		}
		return
	}
	t := a.term(s, a.bodyScope(fn))
	switch t.flag {
	case "unknown":
		// Binder paths (see above): structural proof follows them.
	case "mismatch":
		a.emit(text, line, name, CodeContractMalformed,
			fmt.Sprintf("%s body is ill-sorted: %s", name, describeSmall(s)),
			"sort-correct terms",
			describeSmall(s),
			"correct the operand sorts")
	case "outside":
		a.emit(text, line, name, CodeContractUnsupported,
			fmt.Sprintf("%s body uses %s, outside the fragment", name, t.what),
			"integer, Boolean, and record terms",
			describeSmall(s),
			"use supported logic, or defer this contract-bearing interface")
	case "call":
		a.callRules(m, fn, name, text, t.what, line)
	}
}

// callRules enforces verdict position 4: an extern or kernel call
// has no proof rules (unsupported); a function call has no usable
// summary until a proving run verifies it (unavailable), whether or
// not it carries a contract. Self-calls stay silent here: the
// recursion pass reports them as one cycle finding.
func (a *admission) callRules(m *Module, fn *FnDecl, name, text, callee string, line int) {
	if callee == fn.Name {
		return
	}
	if a.extern[callee] {
		a.emit(text, line, name, CodeContractUnsupported,
			fmt.Sprintf("%s calls extern %s: no proof-import mechanism exists, so an extern declaration is never a theorem", name, callee),
			"calls to verified source callees only",
			"call "+callee,
			"verify a source formulation, or defer this contract-bearing interface")
		return
	}
	if target, ok := a.fns[callee]; ok {
		if isContracted(target) {
			if a.verified[callee] {
				return
			}
			a.emit(text, line, name, CodeContractUnverifiedDep,
				fmt.Sprintf("%s calls contracted %s: its summary is unavailable until a proving run verifies it", name, callee),
				"callees verified in this run",
				"call "+callee,
				"verify the source callee and its closure; do not treat its declaration as proof")
		} else {
			a.emit(text, line, name, CodeContractUnverifiedDep,
				fmt.Sprintf("%s calls %s, which carries no contract: a verified caller may not rely on unstated summaries", name, callee),
				"callees with verified contracts",
				"call "+callee,
				"give the callee a verified contract, inline the reasoning into supported logic, or defer this interface")
		}
		return
	}
	if isBytesKernel(callee) || isStoreOp(callee) || isDecParts(callee) {
		a.emit(text, line, name, CodeContractUnsupported,
			fmt.Sprintf("%s calls kernel %s: tested behavior is not a proof rule, and this cut models no kernel theory", name, callee),
			"calls to verified source callees only",
			"call "+callee,
			"verify a source formulation, or defer this contract-bearing interface")
		return
	}
	// Unknown on a clean program is unreachable (resolution runs
	// first); staying silent beats inventing a rule for it.
}

// checkRecursion rejects direct and mutual recursion in the
// contracted closure: verifying a postcondition across a self-call
// needs an induction rule this cut does not have. Language
// termination rules are unchanged; this is a proof-domain
// restriction, not a program rejection.
func (a *admission) checkRecursion() []string {
	edges := map[string]map[string]bool{}
	var contracted []string
	for name, fn := range a.fns {
		if !isContracted(fn) {
			continue
		}
		contracted = append(contracted, name)
		edges[name] = map[string]bool{}
		collectCalls(fn.Body, func(callee string) {
			if target, ok := a.fns[callee]; ok && isContracted(target) {
				edges[name][callee] = true
			}
		})
	}
	var involved []string
	for _, name := range contracted {
		if reaches(edges, name, name) {
			involved = append(involved, name)
			fn := a.fns[name]
			m := a.ownerOf(fn)
			var text string
			if m != nil {
				text = a.textOf(m)
			}
			short, line := declNameLine(fn)
			a.emit(text, line, short, CodeContractUnsupported,
				fmt.Sprintf("%s is recursive in the contracted closure: the first cut has no induction rule", short),
				"an acyclic contracted call graph",
				"recursive proof over "+short,
				"verify a non-recursive formulation, or defer this contract-bearing interface")
		}
	}
	return involved
}

// collectCalls gathers every called name in a body tree.
func collectCalls(node *Node, f func(string)) {
	if node == nil {
		return
	}
	if node.IsMatch {
		for _, s := range node.Scruts {
			collectCallSmalls(s, f)
		}
		for _, armNode := range node.Arms {
			collectCalls(armNode.Rhs, f)
		}
		return
	}
	collectCallSmalls(node.Small, f)
}

func collectCallSmalls(s *Small, f func(string)) {
	if s == nil {
		return
	}
	if s.Kind == "call" {
		f(s.Fname)
	}
	collectCallSmalls(s.L, f)
	collectCallSmalls(s.R, f)
	collectCallSmalls(s.Hi, f)
	collectCallSmalls(s.Outcome, f)
	for _, arg := range s.Args {
		if arg.V != nil {
			collectCallSmalls(arg.V, f)
		}
	}
	for _, it := range s.Items {
		collectCallSmalls(it, f)
	}
}

// reaches reports whether start reaches target through the call
// edges, consuming at least one edge (so a self-loop counts).
func reaches(edges map[string]map[string]bool, start, target string) bool {
	seen := map[string]bool{}
	stack := []string{start}
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for next := range edges[n] {
			if next == target {
				return true
			}
			if !seen[next] {
				seen[next] = true
				stack = append(stack, next)
			}
		}
	}
	return false
}
