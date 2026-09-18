package main

import (
	"fmt"
)

// Slice 4: or-patterns. One source arm, several alternatives per
// slot; coverage unions over them, credit is sequential (an
// alternative must add space beyond earlier arms and earlier
// alternatives of its own arm, AIL4112).

// patCoverAtoms renders one scalar pattern as its covered atom
// indices in a slot domain: the single-atom cases verbatim from
// the legacy cover loop, so lone and alternative patterns read
// one fact. Missing or unresolvable bounds contribute nothing
// (fail-closed past a world error, which already fails).
func patCoverAtoms(p Pattern, slot int, domains [][]valueAtom, atomIndex []map[string]int) []int {
	switch p.Kind {
	case "wild":
		out := make([]int, 0, len(domains[slot]))
		for ai := range domains[slot] {
			out = append(out, ai)
		}
		return out
	case "bool":
		key := "false"
		if p.B {
			key = "true"
		}
		if a, ok := atomIndex[slot][key]; ok {
			return []int{a}
		}
	case "str":
		if a, ok := atomIndex[slot][normStr(p.Str)]; ok {
			return []int{a}
		}
	case "int":
		if p.Num != nil {
			for k, a := range domains[slot] {
				if a.lo != nil && a.lo.Cmp(p.Num) == 0 {
					return []int{k}
				}
			}
		}
	case "range":
		if p.Num != nil && p.Hi != nil {
			var out []int
			for k, a := range domains[slot] {
				if a.lo != nil && a.lo.Cmp(p.Num) >= 0 && a.lo.Cmp(p.Hi) <= 0 {
					out = append(out, k)
				}
			}
			return out
		}
	case "or":
		seen := map[int]bool{}
		var out []int
		for _, alt := range p.Alts {
			for _, a := range patCoverAtoms(alt, slot, domains, atomIndex) {
				if !seen[a] {
					seen[a] = true
					out = append(out, a)
				}
			}
		}
		return out
	}
	return nil
}

// orAltUseful reports whether alternative k of or-arm i adds
// coverage: some tuple inside its region (its atoms in slot,
// the arm's atoms elsewhere) lies outside earlier arms and
// earlier alternatives of the same arm, modeled as synthetic
// live covers. Pure existence, like armUseful.
func orAltUseful(covers []armCover, i, slot int, altAtoms [][]int, k, nslot int, domains [][]valueAtom) bool {
	live := make([]armCover, 0, i+k)
	live = append(live, covers[:i]...)
	for j := 0; j < k; j++ {
		synth := armCover{slots: make([][]int, nslot)}
		for s := 0; s < nslot; s++ {
			if s == slot {
				synth.slots[s] = altAtoms[j]
			} else {
				synth.slots[s] = covers[i].slots[s]
			}
		}
		live = append(live, synth)
	}
	allowed := make([][]int, nslot)
	for s := 0; s < nslot; s++ {
		if s == slot {
			allowed[s] = altAtoms[k]
		} else {
			allowed[s] = covers[i].slots[s]
		}
	}
	idx, _ := witnessSearchIn(live, nslot, domains, allowed)
	return idx != nil
}

// checkOrAlternatives enforces per-alternative usefulness over
// computed covers: every explicitly written alternative of
// every or-arm must contribute remaining space (AIL4112).
// Arms without alternatives never reach here; callers run this
// only after their legacy gates pass, so domains are sound.
func checkOrAlternatives(n *Node, owner string, covers []armCover, domains [][]valueAtom, atomIndex []map[string]int) []error {
	var out []error
	nslot := len(n.Scruts)
	for i := range n.Arms {
		for s, p := range n.Arms[i].Pats {
			if p.Kind != "or" {
				continue
			}
			altAtoms := make([][]int, len(p.Alts))
			for k, alt := range p.Alts {
				altAtoms[k] = patCoverAtoms(alt, s, domains, atomIndex)
			}
			for k, alt := range p.Alts {
				if !orAltUseful(covers, i, s, altAtoms, k, nslot, domains) {
					out = append(out, at(n.Arms[i].Line, fmt.Errorf("%s: or-alternative %s contributes no remaining space", owner, patAtomRender(alt))))
				}
			}
		}
	}
	return out
}

// hasOrArm reports whether a match carries any or-pattern:
// only such tables pay for cover computation on paths that
// otherwise prove without covers.
func hasOrArm(n *Node) bool {
	for _, a := range n.Arms {
		for _, p := range a.Pats {
			if p.Kind == "or" {
				return true
			}
		}
	}
	return false
}
