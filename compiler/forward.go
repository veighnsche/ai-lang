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
				isInvoke := m.Kind == MatchInvoke
				if (m.Kind != MatchCall && !isInvoke) || len(m.Scruts) != 1 ||
					(!isInvoke && m.Scruts[0].Kind != "call") || (isInvoke && m.Scruts[0].Kind != "ref") {
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
					var rewrite *Small
					var err error
					if isInvoke {
						name := ""
						if len(m.Scruts[0].Ref) == 1 {
							name = m.Scruts[0].Ref[0]
						}
						if m.invokeSig == nil {
							fail(a.Line, "forward cannot resolve invoke target %s", name)
							continue
						}
						rewrite, err = forwardOkRet(prog, fn, name, m.invokeSig.ret, pat.Var)
					} else {
						rewrite, err = forwardOk(prog, fn, m.Scruts[0].Fname, pat.Var)
					}
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

// forwardOk resolves an Ok forward: records must hold exactly the
// destination success fields; variants must have the same nominal type.
// Different wrapper-record names are fine. Reconstruction is explicit,
// never a nominal cast. Lookup covers local, foreign, and extern providers.
func forwardOk(prog *Program, fn *FnDecl, callee, binder string) (*Small, error) {
	calleeRet := ""
	if f, ok := prog.Fns[callee]; ok {
		calleeRet = f.Ret
	} else if e, ok := prog.Externs[callee]; ok {
		calleeRet = e.Ret
	} else {
		return nil, fmt.Errorf("forward cannot resolve callee %s", callee)
	}
	return forwardOkRet(prog, fn, callee, calleeRet, binder)
}

// forwardOkRet rebuilds an Ok payload from a known success type.
// Calls resolve it through the callee; invoke matches read it from the
// callable signature. Variant envelopes retain
// nominal identity; records retain field-shape compatibility.
func forwardOkRet(prog *Program, fn *FnDecl, who, calleeRet, binder string) (*Small, error) {
	if prog.Variants[calleeRet] != nil || prog.Variants[fn.Ret] != nil {
		// Variant success forwarding preserves nominal identity, never
		// coerces between instances or a variant and a wrapper record.
		if calleeRet != fn.Ret {
			return nil, fmt.Errorf("forward needs the same variant return type: %s -> %s, %s -> %s", who, calleeRet, fn.Name, fn.Ret)
		}
		return forwardCtor("Ok", []string{"value"}, binder), nil
	}
	src := recordDecl(prog, calleeRet)
	dst := recordDecl(prog, fn.Ret)
	if src == nil || dst == nil {
		return nil, fmt.Errorf("forward needs record returns: %s -> %s, %s -> %s", who, calleeRet, fn.Name, fn.Ret)
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
		if !sameType(t, f[1]) {
			return nil, fmt.Errorf("forward field %s changes type: %s to %s", f[0], f[1], t)
		}
		fields = append(fields, f[0])
	}
	return forwardCtor("Ok", fields, binder), nil
}
