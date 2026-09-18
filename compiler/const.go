package main

// Slice 1: named scalar-literal constants. Declaration checks,
// reference checks, pattern elaboration, and the shared program
// const table. Execution references resolve lazily at each
// consumer (eval, emit, admission, types) through lookupConst, so
// termination stays syntactic: the checker never rewrites
// `n - ONE` into `n - 1`. Patterns elaborate to literals once in
// buildWorld so every pipeline sees one shape.

import (
	"fmt"
	"math/big"
)

// constSorts are the V1 admitted constant sorts.
func constSortOf(t string) bool {
	return t == "int" || t == "str" || t == "dec" || t == "bool"
}

// constLiteral reports whether s is a literal of the declared
// sort: the only V1 initializer shape.
func constLiteral(sort string, s *Small) bool {
	if s == nil {
		return false
	}
	switch sort {
	case "int":
		return s.Kind == "int"
	case "bool":
		return s.Kind == "bool"
	case "str":
		return s.Kind == "str"
	case "dec":
		return s.Kind == "dec"
	}
	return false
}

// checkConstDecls validates constant declarations in one module:
// admitted sort plus a matching literal initializer (AIL6016).
// Cycles are unrepresentable (aliases rejected), duplicates are
// reported at world build (AIL2205).
func checkConstDecls(m *Module, text string) []Diag {
	var out []Diag
	for _, d := range m.Decls {
		c, ok := d.(*ConstDecl)
		if !ok {
			continue
		}
		if !constSortOf(c.Type) {
			out = append(out, spanDiag(text, c.Line, "error",
				fmt.Sprintf("const %s declares unsupported sort %s: V1 admits int, str, dec, bool", c.Name, c.Type),
				c.Name, CodeConstNonliteral))
			continue
		}
		if !constLiteral(c.Type, c.Value) {
			out = append(out, spanDiag(text, c.Line, "error",
				fmt.Sprintf("const %s initializer is not a %s literal: V1 initializers are literals, never computed or aliased", c.Name, c.Type),
				c.Name, CodeConstNonliteral))
		}
	}
	return out
}

// lookupConst resolves a qualified constant name in the built
// program: first declaration wins, deterministically.
func lookupConst(prog *Program, name string) (*ConstDecl, bool) {
	if prog == nil {
		return nil, false
	}
	c, ok := prog.Consts[name]
	return c, ok
}

// constScope is one lexical scope for reference checking.
type constScope map[string]bool

// checkConstRefs validates constant references in execution
// positions (body, tests, given): SCREAMING references must
// resolve (AIL2104), foreign ones need a uses pin (AIL2105).
// Contracts are excluded: admission classifies unknown names
// there, and pins govern execution linkage. Binder tracking
// keeps params, test args, and call binders out of the way.
func checkConstRefs(fn *FnDecl, prog *Program, owner *Module, text string) []Diag {
	var out []Diag
	emit := func(line int, token, code, msg, expected, found, hint string) {
		d := spanDiag(text, line, "error", msg, token, code)
		d.Expected = expected
		d.Found = found
		d.Hint = hint
		out = append(out, d)
	}
	check := func(s *Small, line int, scope constScope) {
		if s.Kind != "ref" || len(s.Ref) == 0 {
			return
		}
		head := s.Ref[0]
		if scope[head] {
			return
		}
		if !constNameRe.MatchString(head) {
			return // legacy territory: unbound lowercase fails at eval
		}
		target, ok := lookupConst(prog, head)
		if !ok {
			emit(line, head, CodeUnknownConst,
				fmt.Sprintf("%s references unknown constant %s", fn.Name, head),
				"a declared constant",
				"reference to "+head,
				"declare the constant, or correct the name")
			return
		}
		if prog.ConstFile[head] != owner.ID && !prog.Uses[head] {
			emit(line, head, CodeConstNotInUses,
				fmt.Sprintf("%s references foreign const %s which is not in uses: add name@rev to uses", fn.Name, head),
				fmt.Sprintf("uses [%s@%d]", head, target.Rev),
				"reference to "+head,
				"pin the provider rev; same-module constants need no pin")
		}
	}
	var walkNode func(n *Node, scope constScope, line int)
	walkNode = func(n *Node, scope constScope, line int) {
		if n == nil {
			return
		}
		if n.Line > 0 {
			line = n.Line
		}
		if n.IsMatch {
			for _, s := range n.Scruts {
				walkSmallTrees(s, func(t *Small) { check(t, line, scope) })
			}
			for _, a := range n.Arms {
				inner := scope
				if len(a.Pats) == 1 && (a.Pats[0].Kind == "variant" || a.Pats[0].Kind == "variantWild") && a.Pats[0].Var != "" {
					inner = copyConstScope(scope)
					inner[a.Pats[0].Var] = true
				}
				walkNode(a.Rhs, inner, a.Line)
			}
			return
		}
		walkSmallTrees(n.Small, func(t *Small) { check(t, line, scope) })
	}
	params := constScope{}
	for _, p := range fn.Params {
		params[p[0]] = true
	}
	walkNode(fn.Body, params, fn.Line)
	// Evidence positions: each row binds its own arg names.
	for _, t := range fn.Tests {
		row := constScope{}
		for _, a := range t.Args {
			row[a.Name] = true
			if a.V != nil {
				walkSmallTrees(a.V, func(s *Small) { check(s, t.Line, row) })
			}
		}
		if t.Expected != nil {
			walkSmallTrees(t.Expected, func(s *Small) { check(s, t.Line, row) })
		}
	}
	for _, m := range matchNodes(fn.Body) {
		for key, s := range m.Given {
			var row constScope
			for _, t := range fn.Tests {
				if t.Name == key {
					row = constScope{}
					for _, a := range t.Args {
						row[a.Name] = true
					}
				}
			}
			if row == nil {
				row = params
			}
			walkSmallTrees(s, func(t *Small) { check(t, m.Line, row) })
		}
	}
	return out
}

func copyConstScope(s constScope) constScope {
	out := constScope{}
	for k, v := range s {
		out[k] = v
	}
	return out
}

// usedConsts collects every constant name referenced anywhere in
// a function, body and evidence alike: a pin used only in a
// test row is still used (no unused-uses warning).
func usedConsts(fn *FnDecl) map[string]bool {
	out := map[string]bool{}
	see := func(s *Small) {
		if s.Kind == "ref" && len(s.Ref) > 0 && constNameRe.MatchString(s.Ref[0]) {
			out[s.Ref[0]] = true
		}
	}
	var walkNode func(n *Node)
	walkNode = func(n *Node) {
		if n == nil {
			return
		}
		if n.IsMatch {
			for _, s := range n.Scruts {
				walkSmallTrees(s, see)
			}
			for _, a := range n.Arms {
				walkNode(a.Rhs)
			}
			for _, s := range n.Given {
				walkSmallTrees(s, see)
			}
			return
		}
		walkSmallTrees(n.Small, see)
	}
	walkNode(fn.Body)
	for _, t := range fn.Tests {
		for _, a := range t.Args {
			if a.V != nil {
				walkSmallTrees(a.V, see)
			}
		}
		if t.Expected != nil {
			walkSmallTrees(t.Expected, see)
		}
	}
	return out
}

// elaborateConstPatterns resolves const-named match patterns in
// one module in place, before checks, runs, or proofs see them.
// Declared bool/str/int constants rewrite to literals (slice 3
// admits integer patterns); undeclared names fall back to variant
// treatment (preserving the old "variant pattern on a non-call
// match" behavior for genuinely unknown names); dec constants
// are AIL6016 with the arm neutralized to wild so the prover
// reports no extra noise. Idempotent: rewritten patterns are
// never Kind "const".
func elaborateConstPatterns(open *Module, prog *Program, text string) []Diag {
	var out []Diag
	pinned := map[string]bool{}
	markUsed := func(name string) {
		if prog.ConstUsed == nil {
			prog.ConstUsed = map[string]map[string]bool{}
		}
		if prog.ConstUsed[open.ID] == nil {
			prog.ConstUsed[open.ID] = map[string]bool{}
		}
		prog.ConstUsed[open.ID][name] = true
	}
	pinCheck := func(fn *FnDecl, line int, name string) {
		// A foreign constant in a pattern needs its rev pin
		// like any other reference (AIL2105, once per name):
		// patterns resolve through the same program table
		// as bodies.
		if prog.ConstFile[name] != open.ID && !prog.Uses[name] && !pinned[name] {
			pinned[name] = true
			out = append(out, spanDiag(text, line, "error",
				fmt.Sprintf("%s references foreign const %s which is not in uses: add name@rev to uses", fn.Name, name),
				name, CodeConstNotInUses))
		}
	}
	for _, d := range open.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, m := range matchNodes(fn.Body) {
			for ai := range m.Arms {
				for pi := range m.Arms[ai].Pats {
					p := &m.Arms[ai].Pats[pi]
					if p.Kind == "range" {
						out = append(out, resolveRangePattern(open, prog, text, fn.Name, m.Arms[ai].Line, p)...)
						continue
					}
					if p.Kind != "const" {
						continue
					}
					name := p.Name
					c, ok := lookupConst(prog, name)
					if !ok {
						*p = Pattern{Kind: "variant", Name: name}
						continue
					}
					// A pattern reference marks its pin used even
					// though elaboration rewrites the pattern away
					// before checkSem collects evidence: record it
					// per module for checkUnusedUses.
					markUsed(name)
					pinCheck(fn, m.Arms[ai].Line, name)
					if c.Value != nil {
						if c.Value.Kind == "bool" {
							*p = Pattern{Kind: "bool", B: c.Value.B}
							continue
						}
						if c.Value.Kind == "str" {
							*p = Pattern{Kind: "str", Str: c.Value.Str}
							continue
						}
						if c.Value.Kind == "int" && c.Value.Num != nil {
							// Slice 3: integer constants are
							// patterns now. The value copies
							// out so the pattern owns its
							// arbitrary-precision bound.
							*p = Pattern{Kind: "int", Num: new(big.Int).Set(c.Value.Num)}
							continue
						}
					}
					*p = Pattern{Kind: "wild"}
					out = append(out, spanDiag(text, m.Arms[ai].Line, "error",
						fmt.Sprintf("const %s is %s: V1 patterns admit bool/str/int constants only", name, c.Type),
						name, CodeConstNonliteral))
				}
			}
		}
	}
	return out
}
