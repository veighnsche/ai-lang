package main

// a90: arm-level forward call. `=> forward call f(args)` elaborates
// at check time into the ordinary full match (`Ok` plus one arm
// per callee emits kind, every arm forwarding), so admission, the
// prover, coverage, and emit see exactly the handwritten shape.
// Only same-file local callees elaborate; every other shape is
// CAN3013, through the ordinary check-diagnostic mechanism.

import (
	"fmt"
	"strings"
)

// expandForwardCalls rewrites forward-call arms in one module in
// place. It runs once in buildWorld, after the const table (and
// every declaration it reads) is complete and before
// elaborateForwards: its output arms are plain `forward v`
// shapes the existing pass must then expand. Like the other
// passes it must run exactly once: refused sites stay
// forward-shaped, so a second run would duplicate their
// diagnostics.
func expandForwardCalls(open *Module, prog *Program, text string) []Diag {
	var out []Diag
	fail := func(line int, format string, args ...any) {
		out = append(out, spanDiag(text, line, "error", fmt.Sprintf(format, args...), "forward", CodeBadForwardCall))
	}
	for _, d := range open.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		here := prog.FnFile[fn.Name]
		for _, m := range matchNodes(fn.Body) {
			for ai := range m.Arms {
				a := &m.Arms[ai]
				if a.Rhs == nil || a.Rhs.IsMatch || a.Rhs.Small == nil || a.Rhs.Small.Kind != "forward" {
					continue
				}
				rest := strings.TrimSpace(a.Rhs.Small.Str)
				if rest != "call" && !strings.HasPrefix(rest, "call ") {
					continue // plain forward v: elaborateForwards owns it
				}
				call, err := parseSmall(rest)
				if err != nil || call == nil || call.Kind != "call" {
					fail(a.Line, "forward call needs a call shaped as a match scrutinee: forward call f(args)")
					continue
				}
				if _, isExtern := prog.Externs[call.Fname]; isExtern {
					fail(a.Line, "forward call takes same-file locals only: %s is an extern (extern calls need given tables, so they keep their match)", call.Fname)
					continue
				}
				file, known := prog.FnFile[call.Fname]
				if !known {
					fail(a.Line, "forward call cannot resolve callee %s", call.Fname)
					continue
				}
				if file != here {
					fail(a.Line, "forward call takes same-file locals only: %s is not defined in this file (foreign calls need given tables, so they keep their match)", call.Fname)
					continue
				}
				a.Rhs = forwardCallMatch(call, prog.EmitsOf[call.Fname], a.Line)
			}
		}
	}
	return out
}

// forwardCallMatch builds the relay match a site would spell by
// hand: Ok plus one arm per callee emits kind in declared
// order, every arm forwarding its binder. Binders use the
// handwritten r/e convention and shadow exactly like
// handwritten nested binders (pinned by the call-arm probe);
// every arm shares the site line, and arm identity stays
// (node, index), so coverage and spans keep working.
func forwardCallMatch(call *Small, emits []string, line int) *Node {
	fwd := func(binder string) *Node {
		return &Node{Small: &Small{Kind: "forward", Str: binder}, Line: line}
	}
	arms := []Arm{
		{Pats: []Pattern{{Kind: "variant", Name: "Ok", Var: "r"}}, Rhs: fwd("r"), Line: line},
	}
	for _, e := range emits {
		arms = append(arms, Arm{
			Pats: []Pattern{{Kind: "variant", Name: e, Var: "e"}},
			Rhs:  fwd("e"), Line: line,
		})
	}
	return &Node{IsMatch: true, Kind: MatchCall, Scruts: []*Small{call}, Arms: arms, Line: line}
}
