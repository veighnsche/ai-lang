package main

import (
	"sort"
	"strings"
)

// The error catalog is the build-side inventory of every error kind in
// the program: who raises it, who handles it, which tests hit it. Agents
// grep qualified names; the catalog hands them the index instead of
// making them build it. Only meaningful on clean builds: unknown kinds
// fail compilation before a catalog is written.
type catalogArm struct {
	Fn  string `json:"fn"`
	Arm string `json:"arm"`
}

type catalogEntry struct {
	Kind       string       `json:"kind"`
	Fields     []string     `json:"fields"`
	RaisedBy   []string     `json:"raised_by"`
	HandledBy  []catalogArm `json:"handled_by"`
	HitByTests []string     `json:"hit_by_tests"`
}

type handledArm struct {
	kind, fn, arm string
}

// stubKindsOf extracts the error kinds one test scripts at one call node:
// dotted constructors (Ok is not an error and is skipped), through lists.
func stubKindsOf(sm *Small, kinds map[string]bool) {
	if sm == nil {
		return
	}
	if sm.Kind == "list" {
		for _, it := range sm.Items {
			stubKindsOf(it, kinds)
		}
		return
	}
	if sm.Kind == "exchange" {
		stubKindsOf(sm.Outcome, kinds)
		return
	}
	if sm.Kind == "strlen" || sm.Kind == "stridx" || sm.Kind == "strslice" || sm.Kind == "proj" {
		stubKindsOf(sm.L, kinds)
		stubKindsOf(sm.R, kinds)
		stubKindsOf(sm.Hi, kinds)
		return
	}
	if sm.Kind == "ctor" && sm.Ctor != "Ok" && strings.Contains(sm.Ctor, ".") {
		kinds[sm.Ctor] = true
	}
}

// buildCatalog inventories every error kind across all modules. raised_by
// reuses eachRaise, the same walk checkEmits diagnoses with; handled_by
// records every call-match arm patterning a dotted kind with its source
// row; hit_by_tests records stub references plus expectation references
// as mod.fn/test. Arrays are sorted, kinds are sorted: byte-identical
// reruns or it is a bug. hit_by_tests records references, not
// executions: a certified-but-unexecuted relay arm contributes nothing
// (no fabricated hits), and a declared-but-unrealized error is listed
// with empty lists (conservative upper bound) — see
// TestCatalogNoPhantomHits. Consumers must not read hits as proof a
// kind was observed escaping.
func buildCatalog(mods []*Module, prog *Program, texts map[string]string) []catalogEntry {
	raisedBy := map[string]map[string]bool{}
	var handled []handledArm
	handledSeen := map[string]bool{}
	hitBy := map[string]map[string]bool{}
	mark := func(dst map[string]map[string]bool, kind, val string) {
		if dst[kind] == nil {
			dst[kind] = map[string]bool{}
		}
		dst[kind][val] = true
	}
	for _, m := range mods {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			who := m.Mod + "." + fn.Name
			eachRaise(fn, func(kind string, _ int) { mark(raisedBy, kind, who) })
			rows := strings.Split(texts[m.ID], "\n")
			calls := []*Node{}
			for _, n := range matchNodes(fn.Body) {
				if n.Kind != MatchCall {
					continue
				}
				calls = append(calls, n)
				for _, a := range n.Arms {
					p := a.Pats[0]
					if p.Kind != "variant" && p.Kind != "variantWild" {
						continue
					}
					if !strings.Contains(p.Name, ".") {
						continue
					}
					row := ""
					if a.Line >= 1 && a.Line <= len(rows) {
						row = strings.TrimSpace(rows[a.Line-1])
					}
					// Dedup by row text: two arms with identical source
					// (the retry pair) are one handling site each in
					// spirit, and listing the same row twice reads as a
					// catalog bug, not extra information.
					key := who + "\x00" + p.Name + "\x00" + row
					if !handledSeen[key] {
						handledSeen[key] = true
						handled = append(handled, handledArm{p.Name, who, row})
					}
				}
			}
			// Test references live outside the match loop: value-match
			// bodies (db__get_user matches on id, not on a call) have
			// expectations too, and nesting this under call-matches
			// silently drops them.
			for _, t := range fn.Tests {
				kinds := map[string]bool{}
				for _, n := range calls {
					stubKindsOf(n.Given[t.Name], kinds)
				}
				if w := t.Expected; w.Kind == "ref" && len(w.Ref) > 1 {
					kinds[strings.Join(w.Ref, ".")] = true
				}
				if w := t.Expected; w.Kind == "ctor" && strings.Contains(w.Ctor, ".") {
					kinds[w.Ctor] = true
				}
				for k := range kinds {
					mark(hitBy, k, who+"/"+t.Name)
				}
			}
		}
	}
	// An extern's emits row is the foreign code's contract: the catalog
	// lists the extern as a raiser, or foreign-raised kinds would show
	// no source at all.
	for _, m := range mods {
		for _, d := range m.Decls {
			ex, ok := d.(*ExternDecl)
			if !ok {
				continue
			}
			for _, e := range ex.Emits {
				mark(raisedBy, e, m.Mod+"."+ex.Name)
			}
		}
	}
	// A kernel's declared failures attribute to the compiler
	// intrinsic, not to an .can source location: availability is not
	// an execution claim. Total kernels emit nothing and contribute
	// no rows here.
	for name, k := range bytesKernels {
		for _, e := range k.emits {
			mark(raisedBy, e, "kernel."+name)
		}
	}
	byKind := map[string][]catalogArm{}
	for _, h := range handled {
		byKind[h.kind] = append(byKind[h.kind], catalogArm{Fn: h.fn, Arm: h.arm})
	}
	for k := range byKind {
		arms := byKind[k]
		sort.Slice(arms, func(i, j int) bool {
			if arms[i].Fn != arms[j].Fn {
				return arms[i].Fn < arms[j].Fn
			}
			return arms[i].Arm < arms[j].Arm
		})
		byKind[k] = arms
	}
	kinds := map[string]bool{}
	for k := range prog.Errors {
		kinds[k] = true
	}
	for k := range raisedBy {
		kinds[k] = true
	}
	for k := range byKind {
		kinds[k] = true
	}
	for k := range hitBy {
		kinds[k] = true
	}
	sorted := make([]string, 0, len(kinds))
	for k := range kinds {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	out := make([]catalogEntry, 0, len(sorted))
	for _, k := range sorted {
		fields := prog.Errors[k]
		if fields == nil {
			fields = []string{}
		}
		arms := byKind[k]
		if arms == nil {
			arms = []catalogArm{}
		}
		out = append(out, catalogEntry{
			Kind:       k,
			Fields:     fields,
			RaisedBy:   sortedSet(raisedBy[k]),
			HandledBy:  arms,
			HitByTests: sortedSet(hitBy[k]),
		})
	}
	return out
}

func sortedSet(m map[string]bool) []string {
	out := []string{}
	for k := range m {
		if k != "" {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}
