package main

import "strings"

// invokeSig is the checked monomorphic callable signature behind
// one `match invoke cb with n`: the input type A, the success
// record R, and the declared error kinds in written order. It is
// the static half of the invocation certificate: resolution
// stashes it on the node, and check, proof, forward elaboration,
// run, and emit all consume the stash, failing closed when it is
// absent.
type invokeSig struct {
	// head is the parameter's raw Fn annotation, kept so the
	// checker can tell a shadowed name (binder type in env)
	// from the callable it resolves to.
	head string
	in   string
	ret  string
	errs []string
}

// resolveInvokeSites stashes the callable signature on every
// MatchInvoke in one module. The target is a bare name, so the
// only possible callable denotation is an Fn-typed parameter of
// the enclosing function: arm binders type as records and error
// payloads, never as Fn, so a shadowed name is simply not
// callable and the checker owns that error. Resolution is pure
// (no diagnostics) and runs once in buildWorld, after expansion
// and before elaboration, check, and proof in both pipelines —
// the LSP runs proof before checkSem, so neither consumer may
// set or clear the stash; both only read it.
func resolveInvokeSites(m *Module) {
	for _, d := range m.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		params := map[string]string{}
		for _, p := range fn.Params {
			params[p[0]] = p[1]
		}
		for _, n := range matchNodes(fn.Body) {
			if n.Kind != MatchInvoke || len(n.Scruts) != 1 {
				continue
			}
			s := n.Scruts[0]
			if s.Kind != "ref" || len(s.Ref) != 1 {
				continue
			}
			t, ok := params[s.Ref[0]]
			if !ok {
				continue
			}
			a, r, e, ok := fnTypeShape(t)
			if !ok {
				continue
			}
			var errs []string
			if raw := strings.TrimSpace(e[1 : len(e)-1]); raw != "" {
				for _, k := range strings.Split(raw, ",") {
					errs = append(errs, strings.TrimSpace(k))
				}
			}
			n.invokeSig = &invokeSig{head: t, in: a, ret: r, errs: errs}
		}
	}
}
