package main

// Slice 2: forward arms. `on KIND v => forward v` elaborates at
// check time into the ordinary complete constructor
// (`KIND(f = v.f, ...)`), so admission, the prover, emit, and
// CAN4107 see exactly the handwritten shape. Only an exact
// arm-binder operand elaborates; every other shape is CAN3011,
// through the ordinary check-diagnostic mechanism (no new one).

import (
	"fmt"
	"strings"
)

// elaborateForwards rewrites forward arms in one module in place.
// It runs once in buildWorld, after the const table (and every
// declaration it reads) is complete and before any proof, test
// run, or emit sees the match. Unlike const patterns it must run
// exactly once: unresolvable forwards stay forward-shaped, so a
// second run would duplicate their diagnostics.
func elaborateForwards(open *Module, prog *Program, text string) []Diag {
	var out []Diag
	fail := func(line int, format string, args ...any) {
		out = append(out, spanDiag(text, line, "error", fmt.Sprintf(format, args...), "forward", CodeBadForward))
	}
	for _, d := range open.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, m := range matchNodes(fn.Body) {
			for ai := range m.Arms {
				a := &m.Arms[ai]
				if a.Rhs == nil || a.Rhs.IsMatch || a.Rhs.Small == nil || a.Rhs.Small.Kind != "forward" {
					continue
				}
				operand := strings.TrimSpace(a.Rhs.Small.Str)
				if m.Kind == MatchInvoke {
					continue // invoke arms are call-shaped; forward validity lands with invocation
				}
				if m.Kind != MatchCall || len(m.Scruts) != 1 || m.Scruts[0].Kind != "call" {
					fail(a.Line, "forward is a call-outcome arm shape: value matches cannot forward")
					continue
				}
				if len(a.Pats) != 1 {
					fail(a.Line, "forward needs a single-pattern call arm")
					continue
				}
				pat := a.Pats[0]
				if pat.Kind != "variant" || pat.Var == "" || pat.Var == "_" {
					// `_` never binds: it declares the payload
					// unwanted, so forwarding from it contradicts
					// the declaration even where `_.value`
					// would technically resolve.
					fail(a.Line, "forward needs a bound payload binder: on KIND binder => forward binder")
					continue
				}
				if operand != pat.Var {
					fail(a.Line, "forward operand must be this arm's binder %s", pat.Var)
					continue
				}
				if pat.Name == "Ok" {
					rewrite, err := forwardOk(prog, fn, m.Scruts[0].Fname, pat.Var)
					if err != nil {
						fail(a.Line, "%s", err.Error())
						continue
					}
					a.Rhs.Small = rewrite
					continue
				}
				fields, known := prog.Errors[pat.Name]
				if !known {
					fail(a.Line, "forward of unknown error kind %s", pat.Name)
					continue
				}
				a.Rhs.Small = forwardCtor(pat.Name, fields, pat.Var)
			}
		}
	}
	return out
}

// forwardCtor builds the complete reconstruction the arm would
// spell by hand: KIND(f = v.f, ...) over the declared fields.
func forwardCtor(kind string, fields []string, binder string) *Small {
	args := make([]Arg, 0, len(fields))
	for _, f := range fields {
		args = append(args, Arg{Name: f, HasName: true,
			V: &Small{Kind: "ref", Ref: []string{binder, f}}})
	}
	return &Small{Kind: "ctor", Ctor: kind, Args: args}
}

// forwardOk resolves an Ok forward: the callee's success record
// must hold exactly the enclosing function's success fields with
// resolved types. Different wrapper-record names are fine; the
// operation stays explicit field reconstruction, never a nominal
// cast. Callee lookup covers local, foreign, and extern providers.
func forwardOk(prog *Program, fn *FnDecl, callee, binder string) (*Small, error) {
	calleeRet := ""
	if f, ok := prog.Fns[callee]; ok {
		calleeRet = f.Ret
	} else if e, ok := prog.Externs[callee]; ok {
		calleeRet = e.Ret
	} else {
		return nil, fmt.Errorf("forward cannot resolve callee %s", callee)
	}
	src := recordDecl(prog, calleeRet)
	dst := recordDecl(prog, fn.Ret)
	if src == nil || dst == nil {
		return nil, fmt.Errorf("forward needs record returns: %s -> %s, %s -> %s", callee, calleeRet, fn.Name, fn.Ret)
	}
	if len(src.Fields) != len(dst.Fields) {
		return nil, fmt.Errorf("forward needs exactly the destination success fields: %s holds %d, %s holds %d", calleeRet, len(src.Fields), fn.Ret, len(dst.Fields))
	}
	want := map[string]string{}
	for _, f := range dst.Fields {
		want[f[0]] = f[1]
	}
	fields := make([]string, 0, len(src.Fields))
	for _, f := range src.Fields {
		t, ok := want[f[0]]
		if !ok {
			return nil, fmt.Errorf("forward needs exactly the destination success fields: %s has no %s", fn.Ret, f[0])
		}
		if t != f[1] {
			return nil, fmt.Errorf("forward field %s changes type: %s to %s", f[0], f[1], t)
		}
		fields = append(fields, f[0])
	}
	return forwardCtor("Ok", fields, binder), nil
}
