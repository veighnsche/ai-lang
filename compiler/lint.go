package main

// Command canlc lint is the strictness linter: two checks the
// compiler deliberately leaves advisory. `canlc lint file.can
// [...]` (files or directories; directories walk *.can) prints
// one line per finding, `path:line: message`, sorted by file
// then line, and exits 1 when findings exist, 0 when clean.
// Usage errors exit 2. Files that do not parse are outside
// every check's scope: they report `lint skipped` on stderr and
// never fail the run.
//
// Rule 1 — redundant argument names. A named argument sitting
// in its own parameter slot adds nothing over the positional
// spelling, so it is reported: with params (left, right),
// `f(right = 8)` is fine (reordered) but `f(left = 5)` must be
// `f(5)`. Applies everywhere a positional spelling exists —
// call args and test rows. Exempt by construction: Ok/error
// constructions, seal payloads, and `given` exchange args,
// which have no positional form (the parser refuses anything
// else there, a12).
//
// Rule 2 — mergeable match arms (docs/a88-arm-minimization.md,
// draft). Value matches only, one match at a time: arms with
// identical outcomes (modulo whitespace) whose patterns differ
// in exactly one slot of discrete atoms (bool, str, int
// singleton, const) fold into one or-pattern arm. Wildcard,
// variant, range, and or-containing arms never group; multi-slot
// differences cannot spell as one arm (or binds tighter than
// the comma) and stay tables. Call matches, chains, guards, and
// given tables are out of scope. Read-only: sources are never
// written, so NUL/control bytes are preserved trivially.
// Report lines carry file and line only — the AST has no
// columns, and fabricated columns would rot (deviation from
// a88's file:line:col noted here instead of hidden).

import (
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// lintFinding is one report line. No columns: see the header.
// code carries the CAN3410-3416 rule identity for editor
// diagnostics; the CLI prints file:line: message only.
type lintFinding struct {
	file string
	line int
	msg  string
	code string
}

func (f lintFinding) String() string {
	return fmt.Sprintf("%s:%d: %s", f.file, f.line, f.msg)
}

// lintModule is one parsed lint input with its source text.
type lintModule struct {
	name string
	text string
	mod  *Module
}

// lintFiles runs both checks over named in-memory sources and
// returns sorted findings plus the names it had to skip.
func lintFiles(files map[string]string) (findings []lintFinding, skipped []string) {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	var mods []lintModule
	for _, name := range names {
		m, err := parseModuleText(name, files[name])
		if err != nil {
			skipped = append(skipped, name)
			continue
		}
		mods = append(mods, lintModule{name, files[name], m})
	}
	for _, p := range mods {
		findings = append(findings, lintRedundantNames(p, mods)...)
		findings = append(findings, lintMergeableArms(p)...)
		findings = append(findings, lintChainable(p)...)
		findings = append(findings, lintTableable(p)...)
		findings = append(findings, lintSameOutcome(p)...)
		findings = append(findings, lintForwardable(p, mods)...)
		findings = append(findings, lintRangeMerge(p)...)
	}
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].file != findings[j].file {
			return findings[i].file < findings[j].file
		}
		if findings[i].line != findings[j].line {
			return findings[i].line < findings[j].line
		}
		return findings[i].msg < findings[j].msg
	})
	sort.Strings(skipped)
	return findings, skipped
}

// runLint implements `canlc lint path [...]`.
func runLint(out, errOut io.Writer, argv []string) int {
	var paths []string
	for _, a := range argv {
		fi, err := os.Stat(a)
		if err != nil {
			fmt.Fprintf(errOut, "canlc lint: %v\n", err)
			return 2
		}
		if fi.IsDir() {
			var inner []string
			_ = filepath.Walk(a, func(p string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() && strings.HasSuffix(p, ".can") {
					inner = append(inner, p)
				}
				return nil
			})
			paths = append(paths, inner...)
			continue
		}
		paths = append(paths, a)
	}
	if len(paths) == 0 {
		fmt.Fprintln(errOut, "usage: canlc lint file.can [...]")
		return 2
	}
	files := map[string]string{}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			fmt.Fprintf(errOut, "canlc lint: %v\n", err)
			return 2
		}
		files[p] = string(data)
	}
	findings, skipped := lintFiles(files)
	for _, s := range skipped {
		fmt.Fprintf(errOut, "lint skipped (no parse): %s\n", s)
	}
	for _, f := range findings {
		fmt.Fprintln(out, f.String())
	}
	if len(findings) > 0 {
		return 1
	}
	return 0
}

// lintSigParams resolves a callee's parameter list for the
// redundancy check: same-file function, own-module extern,
// compiler kernel, then a unique cross-file provider. Anything
// ambiguous or unknown resolves nowhere — the linter reports
// redundancy only where it knows the parameter list.
func lintSigParams(mods []lintModule, cur, fname string) ([][2]string, bool) {
	for _, m := range mods {
		if m.name != cur {
			continue
		}
		for _, d := range m.mod.Decls {
			if f, ok := d.(*FnDecl); ok && f.Name == fname {
				return f.Params, true
			}
		}
		for _, d := range m.mod.Decls {
			if e, ok := d.(*ExternDecl); ok && e.Name == fname {
				return e.Params, true
			}
		}
	}
	if k, ok := bytesKernels[fname]; ok {
		return k.params, true
	}
	var found [][2]string
	n := 0
	for _, m := range mods {
		if m.name == cur {
			continue
		}
		for _, d := range m.mod.Decls {
			if f, ok := d.(*FnDecl); ok && f.Name == fname {
				found = f.Params
				n++
			}
		}
	}
	if n == 1 {
		return found, true
	}
	return nil, false
}

// lintRedundantNames implements rule 1 over one module: test
// rows against their own params, calls against the resolved
// callee signature.
func lintRedundantNames(lm lintModule, mods []lintModule) []lintFinding {
	var out []lintFinding
	report := func(line, idx int, name, owner string) {
		out = append(out, lintFinding{lm.name, line, fmt.Sprintf(
			"redundant argument name %q (param %d of %s is %q); write positionally",
			name, idx+1, owner, name), CodeLintRedundant})
	}
	for _, d := range lm.mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, t := range fn.Tests {
			for i, a := range t.Args {
				if !a.HasName {
					continue
				}
				if i < len(fn.Params) && fn.Params[i][0] == a.Name {
					report(t.Line, i, a.Name, fn.Name)
				}
			}
		}
		checkCall := func(s *Small, line int) {
			params, ok := lintSigParams(mods, lm.name, s.Fname)
			if !ok {
				return
			}
			for i, a := range s.Args {
				if !a.HasName {
					continue
				}
				if i < len(params) && params[i][0] == a.Name {
					report(line, i, a.Name, s.Fname)
				}
			}
		}
		for _, m := range matchNodes(fn.Body) {
			for _, sc := range m.Scruts {
				walkSmallTrees(sc, func(c *Small) {
					if c.Kind == "call" {
						checkCall(c, m.Line)
					}
				})
			}
		}
		bodySmalls(fn.Body, func(s *Small, line int) {
			if s.Kind == "call" {
				checkCall(s, line)
			}
		})
	}
	return out
}

// splitArmRHS cuts a match-arm row at the first => outside
// string literals and comments. .can strings open with " and
// honor backslash escapes; // outside a string starts a
// comment whose => is not a separator.
func splitArmRHS(line string) (string, bool) {
	inStr := false
	for i := 0; i < len(line); {
		ch := line[i]
		if inStr {
			if ch == '\\' && i+1 < len(line) {
				i += 2
				continue
			}
			if ch == '"' {
				inStr = false
			}
			i++
			continue
		}
		if ch == '"' {
			inStr = true
			i++
			continue
		}
		if ch == '/' && i+1 < len(line) && line[i+1] == '/' {
			return "", false
		}
		if ch == '=' && i+1 < len(line) && line[i+1] == '>' {
			return line[i+2:], true
		}
		i++
	}
	return "", false
}

// lintAtomKind reports whether a slot pattern is a discrete
// atom eligible for or-folding: exact-equality leaves with no
// binders and no substructure. Ranges need synthesis judgment
// (a88 non-goal), ors are already folded, wildcards stay
// separate (CAN4112), variants bind.
func lintAtomKind(p Pattern) bool {
	if p.Var != "" {
		return false
	}
	switch p.Kind {
	case "bool", "str", "int", "const":
		return true
	}
	return false
}

// patEqual compares slot patterns structurally: kind, payload,
// and (for or-free atoms, the only ones compared) no children.
func patEqual(a, b Pattern) bool {
	if a.Kind != b.Kind || a.B != b.B || a.Str != b.Str ||
		a.Name != b.Name || a.Var != b.Var || a.Raw != b.Raw ||
		a.LoS != b.LoS || a.HiS != b.HiS {
		return false
	}
	if (a.Num == nil) != (b.Num == nil) || (a.Hi == nil) != (b.Hi == nil) {
		return false
	}
	if a.Num != nil && a.Num.Cmp(b.Num) != 0 {
		return false
	}
	if a.Hi != nil && a.Hi.Cmp(b.Hi) != 0 {
		return false
	}
	if len(a.Alts) != len(b.Alts) {
		return false
	}
	for i := range a.Alts {
		if !patEqual(a.Alts[i], b.Alts[i]) {
			return false
		}
	}
	return true
}

// lintMergeableGroup reports whether same-outcome arms fold
// into one or-pattern arm: identical patterns in every slot
// but one, where every pattern is a discrete atom of one kind.
// Discrete-distinct atoms are pairwise disjoint, so the union
// is exact, first-match order is preserved by outcome, and
// every alternative adds space (CAN4112-clean on clean code).
func lintMergeableGroup(arms []lintArm, g []int) bool {
	if len(g) < 2 {
		return false
	}
	nslot := len(arms[g[0]].pats)
	differ := -1
	for s := 0; s < nslot; s++ {
		same := true
		for _, gi := range g[1:] {
			if !patEqual(arms[g[0]].pats[s], arms[gi].pats[s]) {
				same = false
				break
			}
		}
		if !same {
			if differ != -1 {
				return false
			}
			differ = s
		}
	}
	if differ == -1 {
		return false
	}
	kind := arms[g[0]].pats[differ].Kind
	for _, gi := range g {
		p := arms[gi].pats[differ]
		if p.Kind != kind || !lintAtomKind(p) {
			return false
		}
	}
	for x := 0; x < len(g); x++ {
		for y := x + 1; y < len(g); y++ {
			if patEqual(arms[g[x]].pats[differ], arms[g[y]].pats[differ]) {
				return false
			}
		}
	}
	return true
}

// lintArm is a groupable match arm: line, slot patterns, and
// the whitespace-normalized outcome text.
type lintArm struct {
	line int
	idx  int
	pats []Pattern
	key  string
}

// lintArmKey normalizes one arm's outcome text for identity
// comparison. Nested-match right-hand sides are refused: their
// first line collides across different bodies (schema scan
// true/false arms share a scrutinee), so first-line text cannot
// key them. Only single-expression outcomes group.
func lintArmKey(lines []string, a Arm) (string, bool) {
	if a.Rhs == nil || a.Rhs.IsMatch || a.Rhs.Small == nil {
		return "", false
	}
	if a.Line < 1 || a.Line > len(lines) {
		return "", false
	}
	rhs, ok := splitArmRHS(lines[a.Line-1])
	if !ok {
		return "", false
	}
	return strings.Join(strings.Fields(rhs), " "), true
}

// lintPureSmall reports whether a scrutinee can fault at
// runtime: no calls (kernels included), no index or slice
// reads. Everything else — refs, literals, arithmetic, ctors,
// strlen — is total. Unknown kinds default closed (impure): a
// new faulting operator must opt into purity explicitly.
func lintPureSmall(s *Small) bool {
	pure := true
	walkSmallTrees(s, func(c *Small) {
		switch c.Kind {
		case "str", "int", "bool", "dec", "ref", "binop", "not", "neg",
			"strlen", "ctor", "list", "seqlit", "seal", "forward":
		default:
			pure = false
		}
	})
	return pure
}

// lintCodePart cuts // comments outside string literals.
func lintCodePart(line string) string {
	inStr := false
	for i := 0; i < len(line); {
		ch := line[i]
		if inStr {
			if ch == '\\' && i+1 < len(line) {
				i += 2
				continue
			}
			if ch == '"' {
				inStr = false
			}
			i++
			continue
		}
		if ch == '"' {
			inStr = true
			i++
			continue
		}
		if ch == '/' && i+1 < len(line) && line[i+1] == '/' {
			return line[:i]
		}
		i++
	}
	return line
}

// lintInnerKey renders a nested match's scrutinee list for
// identity comparison. The nested match opens on its arm's own
// row (`true => match b < 0`), so the arm line — not the inner
// node, which carries no line — is the source of truth. Same
// file and scope: identical text is identical evaluation.
func lintInnerKey(lines []string, a Arm) (string, bool) {
	if a.Line < 1 || a.Line > len(lines) {
		return "", false
	}
	rhs, ok := splitArmRHS(lines[a.Line-1])
	if !ok {
		return "", false
	}
	fields := strings.Fields(lintCodePart(rhs))
	if len(fields) < 2 || fields[0] != "match" {
		return "", false
	}
	return strings.Join(fields[1:], " "), true
}

// lintDeadMatch reports whether a match decides nothing: every
// arm yields the identical plain expression over pure
// scrutinees. Same-outcome reports it; or-folding skips it
// (dropping subsumes folding).
func lintDeadMatch(lines []string, m *Node) bool {
	for _, s := range m.Scruts {
		if !lintPureSmall(s) {
			return false
		}
	}
	_, same := lintSameKey(lines, m)
	return same
}

// lintKernelErrors mirrors the compiler-owned codec errors,
// each carrying the single `value` field (compiler/emit.go
// helpers, used as such in std/text/text.can).
var lintKernelErrors = map[string][]string{
	"encoding.invalid_utf8":   {"value"},
	"encoding.invalid_hex":    {"value"},
	"encoding.invalid_base64": {"value"},
}

// lintMergeableArms implements rule 2 over one module.
func lintMergeableArms(lm lintModule) []lintFinding {
	lines := strings.Split(lm.text, "\n")
	var out []lintFinding
	for _, d := range lm.mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, m := range matchNodes(fn.Body) {
			if m.Kind != MatchValue || m.Given != nil || len(m.Arms) < 2 {
				continue
			}
			if lintDeadMatch(lines, m) {
				continue
			}
			var arms []lintArm
			for ai, a := range m.Arms {
				groupable := true
				for _, p := range a.Pats {
					if !lintAtomKind(p) {
						groupable = false
						break
					}
				}
				if !groupable {
					continue
				}
				key, ok := lintArmKey(lines, a)
				if !ok {
					continue
				}
				arms = append(arms, lintArm{a.Line, ai, a.Pats, key})
			}
			// Groups must be adjacent source blocks: an
			// interleaved arm with overlapping cells would
			// change dispatch under the fold, so same-outcome
			// arms separated by anything stay silent.
			start := 0
			flush := func(end int) {
				g := make([]int, end-start)
				for i := range g {
					g[i] = start + i
				}
				if !lintMergeableGroup(arms, g) {
					return
				}
				ls := make([]int, len(g))
				for i, gi := range g {
					ls[i] = arms[gi].line
				}
				out = append(out, lintFinding{lm.name, ls[0], fmt.Sprintf(
					"mergeable match arms %v (identical outcomes); fold into one or-pattern arm, saves %d lines",
					ls, len(g)-1), CodeLintOrFold})
			}
			for i := 1; i <= len(arms); i++ {
				if i < len(arms) && arms[i].key == arms[start].key &&
					arms[i].idx == arms[i-1].idx+1 {
					continue
				}
				flush(i)
				start = i
			}
		}
	}
	return out
}

// lintOkBinder returns the bound payload binder of a call
// match's Ok arm, or false. Chain steps thread through it, so
// a missing or blank binder ends any run.
func lintOkBinder(n *Node) (string, bool) {
	var binder string
	found := false
	for i := range n.Arms {
		a := &n.Arms[i]
		if len(a.Pats) == 1 && a.Pats[0].Kind == "variant" && a.Pats[0].Name == "Ok" {
			if found {
				return "", false
			}
			found = true
			if a.Pats[0].Var == "" || a.Pats[0].Var == "_" {
				return "", false
			}
			binder = a.Pats[0].Var
		}
	}
	if !found {
		return "", false
	}
	return binder, true
}

// lintChainable implements the chain-combinator rule (a86): a
// linear run of two or more `match call` rungs whose Ok arms
// continue directly — or through one pure bool router whose
// other side yields the shared failure — with every error arm
// yielding the identical plain expression folds into one
// `match chain` with that expression as the shared else. Runs
// need at least one error arm: without a shared failure there
// is nothing to share. Maximal runs only; sub-runs contained
// in a reported run stay silent.
func lintChainable(lm lintModule) []lintFinding {
	lines := strings.Split(lm.text, "\n")
	var out []lintFinding
	var runs [][]*Node
	for _, d := range lm.mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, start := range matchNodes(fn.Body) {
			if start.Kind != MatchCall {
				continue
			}
			if r, ok := lintChainRun(start, lines); ok && len(r) >= 2 {
				runs = append(runs, r)
			}
		}
	}
	contained := func(outer, inner []*Node) bool {
		if len(inner) >= len(outer) {
			return false
		}
		set := map[*Node]bool{}
		for _, n := range outer {
			set[n] = true
		}
		for _, n := range inner {
			if !set[n] {
				return false
			}
		}
		return true
	}
	for i, r := range runs {
		sub := false
		for j, o := range runs {
			if i != j && contained(o, r) {
				sub = true
				break
			}
		}
		if sub {
			continue
		}
		out = append(out, lintFinding{lm.name, r[0].Line, fmt.Sprintf(
			"chainable ladder: %d sequential calls sharing one failure outcome; rewrite as match chain",
			len(r)), CodeLintChain})
	}
	return out
}

// lintChainRun extends one maximal chainable run from start.
// Unity is the shared failure text: every error arm in the run
// — plus any bool-router failure side — must yield it plainly.
// Router failure sides may establish unity (a rung that only
// ever fails through its guard still shares the outcome).
// Unity must also be binder-free: chain elaboration binds error
// arms with wildcards (compiler/chain.go, elabChain), so an
// expression referencing an arm binder cannot become the shared
// else. Binder-shaped unity (relays) is rule 6 territory.
func lintChainRun(start *Node, lines []string) ([]*Node, bool) {
	unity, hasUnity, failures := "", false, false
	adopt := func(texts []string) bool {
		for _, t := range texts {
			if !hasUnity {
				unity, hasUnity = t, true
				continue
			}
			if t != unity {
				return false
			}
		}
		return true
	}
	rungErrors := func(n *Node) ([]string, []*Small, bool) {
		var texts []string
		var smalls []*Small
		for _, a := range n.Arms {
			if len(a.Pats) == 1 && a.Pats[0].Kind == "variant" && a.Pats[0].Name == "Ok" {
				continue
			}
			key, ok := lintArmKey(lines, a)
			if !ok {
				return nil, nil, false
			}
			texts = append(texts, key)
			smalls = append(smalls, a.Rhs.Small)
		}
		return texts, smalls, true
	}
	okNext := func(n *Node) (*Node, []*Small, bool) {
		var body *Node
		count := 0
		for i := range n.Arms {
			a := &n.Arms[i]
			if len(a.Pats) == 1 && a.Pats[0].Kind == "variant" && a.Pats[0].Name == "Ok" {
				count++
				if a.Rhs == nil {
					return nil, nil, false
				}
				body = a.Rhs
			}
		}
		if count != 1 || body == nil {
			return nil, nil, false
		}
		if _, ok := lintOkBinder(n); !ok {
			return nil, nil, false
		}
		if body.IsMatch && body.Kind == MatchCall {
			return body, nil, true
		}
		// Bool router: pure single scrutinee, true/false arms,
		// one side continuing, the other yielding unity.
		if !body.IsMatch || body.Kind != MatchValue || body.Given != nil {
			return nil, nil, false
		}
		if len(body.Scruts) != 1 || !lintPureSmall(body.Scruts[0]) {
			return nil, nil, false
		}
		if len(body.Arms) != 2 {
			return nil, nil, false
		}
		var cont *Node
		var sides []string
		var sideSmalls []*Small
		for i := range body.Arms {
			ra := &body.Arms[i]
			if len(ra.Pats) != 1 || ra.Pats[0].Kind != "bool" {
				return nil, nil, false
			}
			if ra.Rhs != nil && ra.Rhs.IsMatch && ra.Rhs.Kind == MatchCall {
				if cont != nil {
					return nil, nil, false
				}
				cont = ra.Rhs
				continue
			}
			key, ok := lintArmKey(lines, *ra)
			if !ok {
				return nil, nil, false
			}
			sides = append(sides, key)
			sideSmalls = append(sideSmalls, ra.Rhs.Small)
		}
		if cont == nil {
			return nil, nil, false
		}
		// Commit to the shared unity only now: a rejected router
		// must not establish unity or failures for the run (two
		// infallible calls share no failure outcome).
		if !adopt(sides) {
			return nil, nil, false
		}
		if len(sides) > 0 {
			failures = true
		}
		return cont, sideSmalls, true
	}
	var run []*Node
	var errSmalls []*Small
	visited := map[*Node]bool{}
	cur := start
	for cur != nil && cur.Kind == MatchCall && !visited[cur] {
		visited[cur] = true
		texts, smalls, ok := rungErrors(cur)
		if !ok || !adopt(texts) {
			break
		}
		if len(texts) > 0 {
			failures = true
		}
		errSmalls = append(errSmalls, smalls...)
		run = append(run, cur)
		nxt, sides, ok := okNext(cur)
		if !ok {
			break
		}
		errSmalls = append(errSmalls, sides...)
		cur = nxt
	}
	if len(run) < 2 || !failures {
		return nil, false
	}
	if lintChainBound(run, errSmalls) {
		return nil, false
	}
	return run, true
}

// lintChainBound reports whether the run's shared failure text
// leans on an arm binder: every rung's Ok binder and every error
// arm's binder is collected, and any head reference to one inside
// a unity expression disqualifies the run. forward counts as
// bound unconditionally: it rebuilds from its arm binder by
// definition. Chain elaboration re-parses one else text for
// wildcard-bound error arms at every level, so binder-bound text
// (relays, forwards) cannot become that else — however
// textually identical the arms read.
func lintChainBound(run []*Node, smalls []*Small) bool {
	bound := map[string]bool{}
	for _, n := range run {
		if b, ok := lintOkBinder(n); ok {
			bound[b] = true
		}
		for i := range n.Arms {
			a := &n.Arms[i]
			if len(a.Pats) != 1 {
				continue
			}
			if v := a.Pats[0].Var; v != "" && v != "_" {
				bound[v] = true
			}
		}
	}
	for _, s := range smalls {
		hit := false
		walkSmallTrees(s, func(c *Small) {
			if hit {
				return
			}
			// forward rebuilds from its arm binder, which the
			// shared else scope cannot see (wildcard binders).
			if c.Kind == "forward" {
				hit = true
				return
			}
			if c.Kind != "ref" || len(c.Ref) == 0 {
				return
			}
			if bound[c.Ref[0]] {
				hit = true
			}
		})
		if hit {
			return true
		}
	}
	return false
}

// lintTableable implements the multi-scrutinee rule: nested
// matches whose every outer arm body is directly an inner
// value match over the identical pure scrutinee fold into one
// multi-scrutinee table (a28 shape). Scrutinees must be pure —
// eager table evaluation must not fault where guarded nesting
// would not (CAN4109) — and identical: different inner
// scrutinees would need don't-care slots bools cannot spell,
// so those nests stay nested.
func lintTableable(lm lintModule) []lintFinding {
	lines := strings.Split(lm.text, "\n")
	var out []lintFinding
	for _, d := range lm.mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, m := range matchNodes(fn.Body) {
			if m.Kind != MatchValue || m.Given != nil || len(m.Arms) < 2 {
				continue
			}
			var key string
			foldable := true
			for _, a := range m.Arms {
				inner := a.Rhs
				if inner == nil || !inner.IsMatch || inner.Kind != MatchValue ||
					inner.Given != nil || len(inner.Arms) < 2 {
					foldable = false
					break
				}
				for _, s := range inner.Scruts {
					if !lintPureSmall(s) {
						foldable = false
						break
					}
				}
				if !foldable {
					break
				}
				k, ok := lintInnerKey(lines, a)
				if !ok {
					foldable = false
					break
				}
				if key == "" {
					key = k
				} else if k != key {
					foldable = false
					break
				}
			}
			if !foldable || key == "" {
				continue
			}
			out = append(out, lintFinding{lm.name, m.Line,
				"nested matches share one scrutinee; fold into a multi-scrutinee table", CodeLintTable})
		}
	}
	return out
}

// lintSameKey reports whether every arm of m yields the
// identical plain expression, returning its normalized text.
func lintSameKey(lines []string, m *Node) (string, bool) {
	var key string
	for i, a := range m.Arms {
		k, ok := lintArmKey(lines, a)
		if !ok {
			return "", false
		}
		if i == 0 {
			key = k
		} else if k != key {
			return "", false
		}
	}
	return key, true
}

// lintSameOutcome drops matches that decide nothing: every arm
// yields the identical plain expression over pure scrutinees,
// so the match always yields it. Faulting scrutinees stay —
// dropping the match would drop the fault.
func lintSameOutcome(lm lintModule) []lintFinding {
	lines := strings.Split(lm.text, "\n")
	var out []lintFinding
	for _, d := range lm.mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, m := range matchNodes(fn.Body) {
			if m.Kind != MatchValue || m.Given != nil || len(m.Arms) < 2 {
				continue
			}
			same := true
			for _, s := range m.Scruts {
				if !lintPureSmall(s) {
					same = false
					break
				}
			}
			var key string
			if same {
				key, same = lintSameKey(lines, m)
			}
			if !same {
				continue
			}
			rhs := key
			if len(rhs) > 72 {
				rhs = rhs[:69] + "..."
			}
			out = append(out, lintFinding{lm.name, m.Line, fmt.Sprintf(
				"match always yields %s; drop the match", rhs), CodeLintSameOutcome})
		}
	}
	return out
}

// lintForwardable implements the relay spelling rule: an error
// arm rebuilding its own kind field-for-field from its binder
// is exactly what `forward` elaborates into (compiler/forward.go),
// so the handwritten form must be `on KIND binder => forward binder`.
// Mirrors relayStatus without the coverage half: locality and
// firstness certify tests, never spelling.
func lintForwardable(lm lintModule, mods []lintModule) []lintFinding {
	fields := map[string][]string{}
	for _, m := range mods {
		for _, d := range m.mod.Decls {
			if e, ok := d.(*ErrorDecl); ok {
				var names []string
				for _, f := range e.Fields {
					names = append(names, f[0])
				}
				fields[e.Name] = names
			}
		}
	}
	for k, v := range lintKernelErrors {
		fields[k] = v
	}
	var out []lintFinding
	for _, d := range lm.mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, m := range matchNodes(fn.Body) {
			if m.Kind != MatchCall {
				continue
			}
			for _, a := range m.Arms {
				if len(a.Pats) != 1 {
					continue
				}
				pat := a.Pats[0]
				if pat.Kind != "variant" || pat.Name == "Ok" ||
					pat.Var == "" || pat.Var == "_" {
					continue
				}
				if a.Rhs == nil || a.Rhs.IsMatch || a.Rhs.Small == nil ||
					a.Rhs.Small.Kind != "ctor" || a.Rhs.Small.Ctor != pat.Name {
					continue
				}
				want, known := fields[pat.Name]
				if !known {
					continue
				}
				s := a.Rhs.Small
				if len(s.Args) != len(want) {
					continue
				}
				seen := map[string]bool{}
				relay := true
				for _, arg := range s.Args {
					known := false
					for _, f := range want {
						if f == arg.Name {
							known = true
							break
						}
					}
					if !arg.HasName || !known || seen[arg.Name] {
						relay = false
						break
					}
					seen[arg.Name] = true
					v := arg.V
					if v == nil || v.Kind != "ref" || len(v.Ref) != 2 ||
						v.Ref[0] != pat.Var || v.Ref[1] != arg.Name {
						relay = false
						break
					}
				}
				if relay {
					for _, f := range want {
						if !seen[f] {
							relay = false
							break
						}
					}
				}
				if relay {
					out = append(out, lintFinding{lm.name, a.Line, fmt.Sprintf(
						"handwritten relay of %s; write forward %s", pat.Name, pat.Var), CodeLintRelay})
				}
			}
		}
	}
	return out
}

// lintRangeMerge joins contiguous int-literal ranges with one
// outcome: `1..3` beside `4..6` yielding X is `1..6` yielding
// X. Adjacent source arms only (dispatch proof as in or-fold);
// const bounds stay out (values resolve at check, lint reads
// parse); contained ranges are the checker's shadow business.
func lintRangeMerge(lm lintModule) []lintFinding {
	lines := strings.Split(lm.text, "\n")
	var out []lintFinding
	one := big.NewInt(1)
	for _, d := range lm.mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, m := range matchNodes(fn.Body) {
			if m.Kind != MatchValue || m.Given != nil {
				continue
			}
			type cand struct {
				line int
				idx  int
				lo   *big.Int
				hi   *big.Int
				key  string
			}
			var cands []cand
			for ai, a := range m.Arms {
				if len(a.Pats) != 1 {
					continue
				}
				p := a.Pats[0]
				if p.Kind != "range" || p.Num == nil || p.Hi == nil {
					continue
				}
				key, ok := lintArmKey(lines, a)
				if !ok {
					continue
				}
				cands = append(cands, cand{a.Line, ai, new(big.Int).Set(p.Num), new(big.Int).Set(p.Hi), key})
			}
			// Chains extend forward over adjacent source arms:
			// next range starts exactly where the running max
			// ends, same outcome, strictly extending. Anything
			// else (gap, overlap, interleaving arm) flushes.
			flush := func(chain []cand) {
				if len(chain) < 2 {
					return
				}
				parts := make([]string, len(chain))
				for i, c := range chain {
					parts[i] = c.lo.String() + ".." + c.hi.String()
				}
				out = append(out, lintFinding{lm.name, chain[0].line, fmt.Sprintf(
					"mergeable ranges %s; join into %s..%s",
					strings.Join(parts, ", "),
					chain[0].lo.String(), chain[len(chain)-1].hi.String()), CodeLintRange})
			}
			var chain []cand
			maxHi := new(big.Int)
			prevIdx := -2
			for _, c := range cands {
				want := new(big.Int).Add(maxHi, one)
				if len(chain) > 0 && c.key == chain[0].key && c.idx == prevIdx+1 &&
					c.lo.Cmp(want) == 0 && c.hi.Cmp(maxHi) > 0 {
					chain = append(chain, c)
					maxHi.Set(c.hi)
					prevIdx = c.idx
					continue
				}
				flush(chain)
				chain = []cand{c}
				maxHi.Set(c.hi)
				prevIdx = c.idx
			}
			flush(chain)
		}
	}
	return out
}

// lintDiagsFor converts one file's lint findings into editor
// diagnostics. Severity is always error: the linter has no
// advisory level (docs/can-idioms.md C6-C12). Only the named
// file's findings convert: sibling lint belongs to its own
// document, and publishDiagnostics draws every line against
// the open URI.
func lintDiagsFor(name string, texts map[string]string) []Diag {
	findings, _ := lintFiles(texts)
	var out []Diag
	for _, f := range findings {
		if f.file != name {
			continue
		}
		out = append(out, Diag{File: name, Line: f.line, Sev: "error", Msg: f.msg, Code: f.code})
	}
	return out
}
