package main

// Slice 3: integer range patterns. Bound resolution runs in
// buildWorld (consts need the finished table); the cut-point
// proof lives here beside the legacy atom machinery it extends.
// Singletons are degenerate spans [v, v]; ranges are inclusive
// [lo, hi]. Cut points are span lowers and upper+1 bounds at
// arbitrary precision, so every arm spans whole atoms and
// membership is constant per atom by construction.

import (
	"fmt"
	"math/big"
	"strings"
)

// resolveRangePattern finishes one parsed range in place: symbolic
// bounds (const names, anything the parser deferred) resolve to
// arbitrary-precision integers, foreign bound constants take
// CAN2105 with used-marking like every other pattern reference,
// and the strict lower-under-upper gap is enforced (CAN4110).
// Literal-resolved ranges still pass through the gap check.
// Unknown const names are CAN2104, never malformed bounds.
// Idempotent: resolved ranges carry Num/Hi with empty LoS/HiS.
func resolveRangePattern(open *Module, prog *Program, text, fnName string, line int, p *Pattern) []Diag {
	var out []Diag
	bad := func(format string, args ...any) {
		out = append(out, spanDiag(text, line, "error", fmt.Sprintf(format, args...), rangeToken(p), CodeBadRange))
	}
	need := func(name string) {
		markConstUsed(prog, open.ID, name)
		if prog.ConstFile[name] != open.ID && !prog.Uses[name] {
			out = append(out, spanDiag(text, line, "error",
				foreignConstMsg(fnName, name),
				name, CodeConstNotInUses))
		}
	}
	if p.Num != nil && p.Hi != nil && p.LoS == "" && p.HiS == "" {
		if p.Num.Cmp(p.Hi) >= 0 {
			bad("invalid integer range bounds %s..%s: lower must be under upper (equal endpoints use the singleton %s)",
				p.Num.String(), p.Hi.String(), p.Num.String())
		}
		return out
	}
	lo, lok := resolveRangeBound(prog, p.LoS)
	hi, hik := resolveRangeBound(prog, p.HiS)
	if !lok {
		out = append(out, rangeBoundDiag(open, prog, text, fnName, line, p.LoS, need)...)
	}
	if !hik {
		out = append(out, rangeBoundDiag(open, prog, text, fnName, line, p.HiS, need)...)
	}
	if !lok || !hik {
		return out
	}
	if lo.Cmp(hi) >= 0 {
		bad("invalid integer range bounds %s..%s: lower must be under upper (equal endpoints use the singleton %s)",
			lo.String(), hi.String(), lo.String())
		return out
	}
	// Resolved bound constants mark their pins used and take
	// CAN2105 when foreign: bounds are references, and a bound
	// used nowhere else must still pin its provider.
	for _, side := range []string{p.LoS, p.HiS} {
		if constNameRe.MatchString(side) {
			need(side)
		}
	}
	p.Num, p.Hi = lo, hi
	p.LoS, p.HiS = "", ""
	return out
}

// rangeToken renders the source spelling for the squiggle
// locator: resolved ranges print numerically, symbolic ones
// verbatim.
func rangeToken(p *Pattern) string {
	if p.LoS != "" || p.HiS != "" {
		return p.LoS + ".." + p.HiS
	}
	if p.Num != nil && p.Hi != nil {
		return p.Num.String() + ".." + p.Hi.String()
	}
	return ".."
}

// resolveRangeBound maps one bound side to its integer: decimal
// literals parse directly, const-shaped names resolve through
// the program table. ok=false covers unknown, malformed, and
// wrong-sort bounds alike; rangeBoundDiag classifies.
func resolveRangeBound(prog *Program, s string) (*big.Int, bool) {
	if isDigits(s) {
		n, _ := new(big.Int).SetString(s, 10)
		return n, true
	}
	if constNameRe.MatchString(s) {
		c, ok := lookupConst(prog, s)
		if !ok || c.Type != "int" || c.Value == nil || c.Value.Kind != "int" || c.Value.Num == nil {
			return nil, false
		}
		return new(big.Int).Set(c.Value.Num), true
	}
	return nil, false
}

// rangeBoundDiag classifies one unresolvable bound side: unknown
// const names are CAN2104 (with used-marking, so no unused-pin
// noise joins it); everything else is CAN4110.
func rangeBoundDiag(open *Module, prog *Program, text, fnName string, line int, side string, need func(string)) []Diag {
	if constNameRe.MatchString(side) {
		if _, ok := lookupConst(prog, side); !ok {
			need(side)
			return []Diag{spanDiag(text, line, "error",
				fmt.Sprintf("%s references unknown constant %s", fnName, side),
				side, CodeUnknownConst)}
		}
		c, _ := lookupConst(prog, side)
		need(side)
		return []Diag{spanDiag(text, line, "error",
			fmt.Sprintf("invalid integer range bounds: %s is %s, bounds are integer literals or visible integer constants", side, c.Type),
			side, CodeBadRange)}
	}
	return []Diag{spanDiag(text, line, "error",
		fmt.Sprintf("invalid integer range bounds: %q is not an integer literal or integer constant", side),
		side, CodeBadRange)}
}

// intSlotSpans collects the inclusive spans of one slot's
// integer patterns: singletons as degenerate [v, v] spans.
// Or-alternatives contribute each alternative's span.
// Unresolved (nil-bound) patterns contribute nothing: they only
// reach the proof past a world error, which already fails.
func intSlotSpans(n *Node, slot int) [][2]*big.Int {
	var out [][2]*big.Int
	var spans func(p Pattern)
	spans = func(p Pattern) {
		switch p.Kind {
		case "int":
			if p.Num != nil {
				out = append(out, [2]*big.Int{p.Num, p.Num})
			}
		case "range":
			if p.Num != nil && p.Hi != nil {
				out = append(out, [2]*big.Int{p.Num, p.Hi})
			}
		case "or":
			for _, alt := range p.Alts {
				spans(alt)
			}
		}
	}
	for _, a := range n.Arms {
		if slot >= len(a.Pats) {
			continue
		}
		spans(a.Pats[slot])
	}
	return out
}

// intCuts returns the sorted unique cut points of a slot's spans:
// every span lower and every span upper+1, at full precision.
func intCuts(spans [][2]*big.Int) []*big.Int {
	seen := map[string]bool{}
	var cuts []*big.Int
	add := func(v *big.Int) {
		key := v.String()
		if seen[key] {
			return
		}
		seen[key] = true
		cuts = append(cuts, new(big.Int).Set(v))
	}
	one := big.NewInt(1)
	for _, s := range spans {
		add(s[0])
		add(new(big.Int).Add(s[1], one))
	}
	for i := 1; i < len(cuts); i++ {
		for j := i; j > 0 && cuts[j].Cmp(cuts[j-1]) < 0; j-- {
			cuts[j], cuts[j-1] = cuts[j-1], cuts[j]
		}
	}
	return cuts
}

// intAtoms builds the slot domain from cut points: the open tail
// below the first cut, one closed-open atom per cut gap, and the
// open tail past the last cut. Renders are concrete members: the
// gap floor, or one step outside for tails. Empty cuts (only past
// a world error) fall back to a whole-domain atom, fail-closed.
func intAtoms(cuts []*big.Int) []valueAtom {
	if len(cuts) == 0 {
		return []valueAtom{{render: "_"}}
	}
	one := big.NewInt(1)
	atoms := []valueAtom{{render: new(big.Int).Sub(cuts[0], one).String(), lo: nil, hi: cuts[0]}}
	for i := 0; i < len(cuts); i++ {
		var hiCut *big.Int
		if i+1 < len(cuts) {
			hiCut = cuts[i+1]
		}
		atoms = append(atoms, valueAtom{render: cuts[i].String(), lo: cuts[i], hi: hiCut})
	}
	return atoms
}

// patRender describes one arm for the useless-arm diagnostic:
// single slots render the pattern, tuples render parenthesized.
func patRender(n *Node, ai int) string {
	if len(n.Arms[ai].Pats) == 1 {
		return patAtomRender(n.Arms[ai].Pats[0])
	}
	parts := make([]string, 0, len(n.Arms[ai].Pats))
	for _, p := range n.Arms[ai].Pats {
		parts = append(parts, patAtomRender(p))
	}
	return "(" + strings.Join(parts, ", ") + ")"
}

// patAtomRender renders one pattern by value: integers print
// numerically (resolved const bounds included), strings stay
// quoted, wild stays wild.
func patAtomRender(p Pattern) string {
	switch p.Kind {
	case "int":
		if p.Num != nil {
			return p.Num.String()
		}
		return "int"
	case "range":
		if p.Num != nil && p.Hi != nil {
			return p.Num.String() + ".." + p.Hi.String()
		}
		if p.LoS != "" || p.HiS != "" {
			return p.LoS + ".." + p.HiS
		}
		return "range"
	case "or":
		parts := make([]string, 0, len(p.Alts))
		for _, alt := range p.Alts {
			parts = append(parts, patAtomRender(alt))
		}
		return strings.Join(parts, " | ")
	case "bool":
		if p.B {
			return "true"
		}
		return "false"
	case "str":
		return normStr(p.Str)
	case "wild":
		return "_"
	case "variant", "variantWild":
		return p.Name
	default:
		return p.Kind
	}
}

// armHasInt reports whether an arm carries an integer pattern in
// any slot (inside or-alternatives counts): only such arms face
// the static usefulness rule, so legacy slots keep legacy
// behavior exactly.
func armHasInt(a Arm) bool {
	var has func(p Pattern) bool
	has = func(p Pattern) bool {
		if p.Kind == "int" || p.Kind == "range" {
			return true
		}
		if p.Kind == "or" {
			for _, alt := range p.Alts {
				if has(alt) {
					return true
				}
			}
		}
		return false
	}
	for _, p := range a.Pats {
		if has(p) {
			return true
		}
	}
	return false
}

// armUseful reports whether arm i adds coverage beyond earlier
// arms: some atom tuple inside its cover lies outside every
// earlier cover. Pure existence — the reporting path derives
// witnesses itself, never from here.
func armUseful(covers []armCover, i, nslot int, domains [][]valueAtom) bool {
	allowed := make([][]int, nslot)
	for s := 0; s < nslot; s++ {
		allowed[s] = covers[i].slots[s]
	}
	idx, _ := witnessSearchIn(covers[:i], nslot, domains, allowed)
	return idx != nil
}

// witnessSearchIn is witnessSearch restricted to per-slot atom
// subsets: the first tuple inside allowed covered by none of
// live, or nil/nil when every allowed tuple is already covered.
func witnessSearchIn(live []armCover, nslot int, domains [][]valueAtom, allowed [][]int) (idx []int, render []string) {
	var pi []int
	var pr []string
	var rec func(slot int, rest []armCover) bool
	rec = func(slot int, rest []armCover) bool {
		if slot == nslot {
			return len(rest) == 0
		}
		for _, ai := range allowed[slot] {
			var next []armCover
			for _, arm := range rest {
				if containsAtom(arm.slots[slot], ai) {
					next = append(next, arm)
				}
			}
			pi = append(pi, ai)
			pr = append(pr, domains[slot][ai].render)
			if rec(slot+1, next) {
				return true
			}
			pi = pi[:len(pi)-1]
			pr = pr[:len(pr)-1]
		}
		return false
	}
	if !rec(0, live) {
		return nil, nil
	}
	return append([]int{}, pi...), append([]string{}, pr...)
}

// verifyIntMatch proves a value table containing integer
// patterns, any arity. Inventory arrives from the shared loop;
// domains reuse the legacy atom machinery with int slots cut at
// span bounds; messages mirror the legacy multi-slot shapes so
// mixed tables read uniformly. Usefulness is static and scoped:
// an integer arm earlier arms fully cover is CAN4111 even before
// any test runs, while legacy arms keep the execution-law path.
func verifyIntMatch(n *Node, owner string, out []error, hasBool []bool, strLits [][]string, intSpans [][][2]*big.Int) []error {
	nslot := len(n.Scruts)
	for i := 0; i < nslot; i++ {
		if len(intSpans[i]) > 0 && (hasBool[i] || len(strLits[i]) > 0) {
			out = append(out, at(n.Line, proofErrf(CodeBoolArms, "%s: slot %d mixes integer patterns with bool/string patterns: one slot takes one scalar domain", owner, i+1)))
			return out
		}
	}
	domains, covers, atomIndex, mixed := valueCoverOf(n, hasBool, strLits, intSpans)
	if mixed > 0 {
		out = append(out, at(n.Line, proofErrf(CodeBoolArms, "%s: bool match must be exactly true+false (slot %d mixes bool and string patterns)", owner, mixed)))
		return out
	}
	for i := range n.Arms {
		if !armHasInt(n.Arms[i]) {
			continue
		}
		if !armUseful(covers, i, nslot, domains) {
			out = append(out, at(n.Arms[i].Line, proofErrf(CodeUselessArm, "%s: arm %s is fully covered by earlier arms", owner, patRender(n, i))))
		}
	}
	// Slice 4: per-alternative usefulness over the same covers.
	out = append(out, checkOrAlternatives(n, owner, covers, domains, atomIndex)...)
	// Collect up to three missing witnesses, feeding each back as
	// a literal arm so the next search finds a new cell: the
	// legacy loop verbatim, over extended domains.
	var wits []string
	otherSlot := 0
	live := covers
	for len(wits) < 3 {
		idx, renders := witnessSearch(live, nslot, domains)
		if idx == nil {
			break
		}
		wits = append(wits, "("+strings.Join(renders, ", ")+")")
		lit := armCover{slots: make([][]int, nslot)}
		for s, a := range idx {
			lit.slots[s] = []int{a}
			if domains[s][a].other && otherSlot == 0 {
				otherSlot = s + 1
			}
		}
		live = append(live, lit)
	}
	if len(wits) == 0 {
		n.analysis = &valueMatchAnalysis{emitFinalElse: residualInLast(covers[:len(covers)-1], covers[len(covers)-1], nslot, domains)}
		return out
	}
	// An uncovered open string remainder dominates: infinitely
	// many cells are missing, so the _-coverage rule fires instead
	// of a finite missing-cell list. Integer tails always render
	// concrete witnesses, never remainder.
	if otherSlot > 0 {
		out = append(out, at(n.Line, proofErrf(CodeValueNoWild, "%s: value match without _ is not provably exhaustive (slot %d leaves an open string remainder)", owner, otherSlot)))
		return out
	}
	out = append(out, at(n.Line, proofErrf(CodeMissingArm, "%s: non-exhaustive match, missing %s", owner, strings.Join(wits, "; "))))
	return out
}
