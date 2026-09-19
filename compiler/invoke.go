package main

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// invokeSig is the checked monomorphic callable signature behind
// one `match invoke cb with n`: the input type A, the success
// type R (data-only), and the declared error kinds in written order. It is
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
// admitted invocation head is still an Fn-typed parameter of
// the enclosing function. B11 typed Ok binders can themselves carry Fn,
// but do not extend this resolver: pass them to an ordinary consumer.
// The checker also rejects binders shadowing parameters with a different type. Resolution is pure
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

// sigKey renders a monomorphic callable signature canonically so
// both sides of an invocation edge compare structurally: target
// annotations and written heads may space their error lists
// differently, but the denoted sets are identical.
func sigKey(a, r string, errs []string) string {
	cp := slices.Clone(errs)
	slices.Sort(cp)
	return "Fn<" + a + ", " + r + ", [" + strings.Join(cp, ", ") + "]>"
}

// admissibleRefSig reports the denoted signature key when s is an
// admissible address-taken reference: the target resolves with
// exactly one parameter unbound and a supported success type (refFnType),
// the denoted input and success are data-only, the reachable
// graph is linked-pure, and no reachable function needs a
// precondition. It mirrors checkFnrefTarget without diagnostics;
// capture type-correctness stays with the checker (a mistyped
// capture fails the build regardless of any cycle edge).
func admissibleRefSig(ck *tycker, prog *Program, s *Small) (string, bool) {
	t, ok := ck.refFnType(s)
	if !ok {
		return "", false
	}
	a, r, _, ok := fnTypeShape(t)
	if !ok || ck.typeHasFn(a) || ck.typeHasFn(r) {
		return "", false
	}
	if checkLinkedGraph(prog, s.Fname) != nil {
		return "", false
	}
	if firstRequiresHolder(prog, s.Fname) != "" {
		return "", false
	}
	return t, true
}

// walkFnrefs collects every reference creation in executable
// positions: arm bodies and plain bodies, call scrutinee
// arguments, and invoke arguments. Creation never executes, so
// references stay out of walkCalls (and the linked-execution
// graph) while import and cycle collection still see the address
// dependency.
func walkFnrefs(body *Node) []*Small {
	var out []*Small
	bodySmalls(body, func(s *Small, line int) {
		if s.Kind == "fnref" {
			out = append(out, s)
		}
	})
	for _, m := range matchNodes(body) {
		for _, s := range m.Scruts {
			walkSmallTrees(s, func(x *Small) {
				if x.Kind == "fnref" {
					out = append(out, x)
				}
			})
		}
	}
	return out
}

// eachFnref visits every reference creation in executable and
// fixture positions: bodies (walkFnrefs) plus test rows (args
// and expectations). Given tables are mocks, not addresses.
func eachFnref(fn *FnDecl, f func(*Small)) {
	see := func(s *Small) {
		walkSmallTrees(s, func(x *Small) {
			if x.Kind == "fnref" {
				f(x)
			}
		})
	}
	for _, s := range walkFnrefs(fn.Body) {
		f(s)
	}
	for _, t := range fn.Tests {
		for _, a := range t.Args {
			see(a.V)
		}
		see(t.Expected)
	}
}

// refTargetsBySig collects every admissible address-taken target
// stamp in the program keyed by denoted signature, from bodies
// and fixtures alike.
func refTargetsBySig(mods []*Module, prog *Program) map[string]map[string]bool {
	ck := newTycker(prog, "", "")
	out := map[string]map[string]bool{}
	for _, m := range mods {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			eachFnref(fn, func(s *Small) {
				key, ok := admissibleRefSig(ck, prog, s)
				if !ok {
					return
				}
				if out[key] == nil {
					out[key] = map[string]bool{}
				}
				out[key][s.Fname] = true
			})
		}
	}
	return out
}

// checkInvokeCycles rejects cycles closed by invocation (b00 I2).
// For each monomorphic Fn signature the program takes addresses
// of, every invocation of that signature contributes potential
// edges from its owner to every admissible target; direct-call
// edges are preserved. Any cycle containing an indirect edge is
// refused — dynamic dispatch admits no decreases proof, not even
// for a self-call (unreachable by construction today, since an
// invoker's Fn parameter can never bind admissibly, but refused
// all the same). Direct-only cycles belong to the local and
// global checks, so this reports only cycles with an indirect
// edge, once, at the first invoke site in cycle order.
func checkInvokeCycles(mods []*Module, texts map[string]string, prog *Program) []Diag {
	var out []Diag
	type edge struct {
		to       string
		token    string
		line     int
		file     string
		indirect bool
	}
	edges := map[string][]edge{}
	add := func(from, to, token string, line int, file string, indirect bool) {
		for _, e := range edges[from] {
			if e.to == to && e.indirect == indirect {
				return
			}
		}
		edges[from] = append(edges[from], edge{to, token, line, file, indirect})
	}
	targets := refTargetsBySig(mods, prog)
	ordered := append([]*Module{}, mods...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].File < ordered[j].File })
	for _, m := range ordered {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			for _, n := range matchNodes(fn.Body) {
				if n.Kind != MatchCall {
					continue
				}
				ns := n.Scruts[0]
				if _, isEx := prog.Externs[ns.Fname]; isEx {
					continue
				}
				if _, ok := prog.Fns[ns.Fname]; !ok {
					continue
				}
				add(fn.Name, ns.Fname, ns.Fname, n.Line, m.ID, false)
			}
			bodySmalls(fn.Body, func(s *Small, line int) {
				if s.Kind != "call" {
					return
				}
				if _, isEx := prog.Externs[s.Fname]; isEx {
					return
				}
				if _, ok := prog.Fns[s.Fname]; !ok {
					return
				}
				add(fn.Name, s.Fname, s.Fname, line, m.ID, false)
			})
			for _, n := range matchNodes(fn.Body) {
				if n.Kind != MatchInvoke || n.invokeSig == nil || len(n.Scruts) != 1 {
					continue
				}
				s := n.Scruts[0]
				if s.Kind != "ref" || len(s.Ref) != 1 {
					continue
				}
				sig := n.invokeSig
				var tos []string
				for to := range targets[sigKey(sig.in, sig.ret, sig.errs)] {
					tos = append(tos, to)
				}
				sort.Strings(tos)
				for _, to := range tos {
					add(fn.Name, to, s.Ref[0], n.Line, m.ID, true)
				}
			}
		}
	}
	invokeSite := map[string]edge{}
	for from, es := range edges {
		for _, e := range es {
			if e.indirect {
				key := from + "\x00" + e.to
				if _, ok := invokeSite[key]; !ok {
					invokeSite[key] = e
				}
			}
		}
	}
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := map[string]int{}
	var stack []string
	var visit func(u string)
	visit = func(u string) {
		color[u] = gray
		stack = append(stack, u)
		for _, e := range edges[u] {
			switch color[e.to] {
			case gray:
				start := 0
				for i, n := range stack {
					if n == e.to {
						start = i
						break
					}
				}
				path := append(append([]string{}, stack[start:]...), e.to)
				// Every directed cycle holds a back edge, and
				// this one's edges are the stack links plus
				// the closer: report only when one is
				// indirect, at the first invoke site in
				// cycle order.
				links := [][2]string{}
				for i := start; i+1 < len(stack); i++ {
					links = append(links, [2]string{stack[i], stack[i+1]})
				}
				links = append(links, [2]string{u, e.to})
				for _, l := range links {
					if s, ok := invokeSite[l[0]+"\x00"+l[1]]; ok {
						d := spanDiag(texts[s.file], s.line, "error",
							fmt.Sprintf("call cycle %s: %s invokes %s through a function value; invocation cannot close a cycle", strings.Join(path, " -> "), l[0], l[1]), s.token, CodeLocalCycle)
						d.File = qualifiedID(mods, s.file)
						out = append(out, d)
						break
					}
				}
			case white:
				visit(e.to)
			}
		}
		stack = stack[:len(stack)-1]
		color[u] = black
	}
	var seeds []string
	for _, m := range ordered {
		for _, d := range m.Decls {
			if fn, ok := d.(*FnDecl); ok {
				seeds = append(seeds, fn.Name)
			}
		}
	}
	for _, s := range seeds {
		if color[s] == white {
			visit(s)
		}
	}
	return out
}
