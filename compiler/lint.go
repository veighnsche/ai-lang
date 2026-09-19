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
// call args, test rows, and constructions (records, cases,
// errors, provable Ok: a92). Exempt: seal payloads and `given`
// exchange args, which have no positional form (the parser
// refuses anything else there, a12); Ok without a provable
// return record, whose slot names are uncheckable.
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
// columns (deviation from a88's file:line:col noted here
// instead of hidden). Editor squiggles still land in 2D: each
// finding carries a [start, end) UTF-16 span located by
// searching the verbatim source line for the exact offending
// token the check already names — never synthesized from the
// AST. A span that cannot be located stays whole-line (end <=
// start), so the editor falls back honestly instead of
// underlining the wrong columns.

import (
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// lintFinding is one report line. code carries the
// CAN3410-3416 rule identity for editor diagnostics; the CLI
// prints file:line: message only. start/end is the 0-based
// UTF-16 [start, end) span of the offending token on line,
// located verbatim in source by the reporting check; end <=
// start means unlocated and the editor falls back to the
// whole line.
type lintFinding struct {
	file       string
	line       int
	msg        string
	code       string
	start, end int
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
		findings = append(findings, lintRestatable(p)...)
		findings = append(findings, lintLadderable(p)...)
		findings = append(findings, lintRelayCallable(p)...)
		findings = append(findings, lintUnreachedKeys(p)...)
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
	lines := strings.Split(lm.text, "\n")
	var out []lintFinding
	report := func(line int, args []Arg, i int, owner, kind string, rhs bool) {
		name := args[i].Name
		// a92: a duplicated name is a fault for the checker,
		// not a rewrite for the linter: removing one of two
		// identical names never fixes the list, so neither
		// occurrence is reported.
		dups := 0
		for _, b := range args {
			if b.HasName && b.Name == name {
				dups++
			}
		}
		if dups > 1 {
			return
		}
		f := lintFinding{file: lm.name, line: line, msg: fmt.Sprintf(
			"redundant argument name %q (%s %d of %s is %q); write positionally",
			name, kind, i+1, owner, name), code: CodeLintRedundant}
		// Underline the name through its equals sign. same
		// counts earlier same-named args, so duplicates and
		// right-hand sides repeating the name never win.
		same := 0
		for _, b := range args[:i] {
			if b.HasName && b.Name == name {
				same++
			}
		}
		if line >= 1 && line <= len(lines) {
			code := lintCodePart(lines[line-1])
			search, off := code, 0
			if rhs {
				// Expectations, arm right-hand sides, and
				// outcomes live past their row arrow: start
				// there so an identical name before the
				// arrow never steals the span. No arrow
				// (plain bodies) falls back to the whole line.
				if sub, ok := splitArmRHS(code); ok {
					search, off = sub, len(code)-len(sub)
				}
			}
			if ns, ne, ok := lintFindNameEq(search, name, same); ok {
				f.start, f.end = u16col(code, off+ns), u16col(code, off+ne)
			}
		}
		out = append(out, f)
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
					report(t.Line, t.Args, i, fn.Name, "param", false)
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
					report(line, s.Args, i, s.Fname, "param", false)
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
		// a92: constructions. A top Ok resolves through the
		// fn return record; records, cases, and errors
		// through their decls. Nested values recurse with no
		// Ok shape: only outermost Ok positions prove, so
		// anything deeper keeps its names unflagged, never
		// misflagged.
		var checkCtorVal func(s *Small, line int, okRet string, rhs bool)
		checkCtorVal = func(s *Small, line int, okRet string, rhs bool) {
			if s == nil {
				return
			}
			if s.Kind == "ctor" && s.Ctor == "Ok" {
				if okFields, ok := lintOkFields(mods, lm.name, okRet); ok {
					for i, a := range s.Args {
						if !a.HasName {
							continue
						}
						if i < len(okFields) && okFields[i][0] == a.Name {
							report(line, s.Args, i, s.Ctor, "field", rhs)
						}
					}
				}
				for _, a := range s.Args {
					checkCtorVal(a.V, line, "", rhs)
				}
				return
			}
			if s.Kind == "ctor" {
				if fields, ok := lintCtorFields(mods, lm.name, s.Ctor); ok {
					for i, a := range s.Args {
						if !a.HasName {
							continue
						}
						if i < len(fields) && fields[i][0] == a.Name {
							report(line, s.Args, i, s.Ctor, "field", rhs)
						}
					}
				}
				for _, a := range s.Args {
					checkCtorVal(a.V, line, "", rhs)
				}
				return
			}
			// Calls, binops, lists, exchanges: names here are
			// call args (owned by the call pass), exchange
			// args (grammar-mandated, exempt), or operators
			// — never constructions. Recurse for nested
			// values with no Ok shape.
			for _, a := range s.Args {
				checkCtorVal(a.V, line, "", rhs)
			}
			checkCtorVal(s.L, line, "", rhs)
			checkCtorVal(s.R, line, "", rhs)
			checkCtorVal(s.Hi, line, "", rhs)
			for _, it := range s.Items {
				checkCtorVal(it, line, "", rhs)
			}
			checkCtorVal(s.Outcome, line, "", rhs)
		}
		for _, t := range fn.Tests {
			for _, a := range t.Args {
				checkCtorVal(a.V, t.Line, "", false)
			}
			checkCtorVal(t.Expected, t.Line, fn.Ret, true)
		}
		lintTops(fn.Body, func(s *Small, line int) {
			checkCtorVal(s, line, fn.Ret, true)
		})
		for _, m := range matchNodes(fn.Body) {
			for _, sc := range m.Scruts {
				checkCtorVal(sc, m.Line, "", false)
			}
			if m.InvokeArg != nil {
				checkCtorVal(m.InvokeArg, m.Line, "", false)
			}
			if m.Kind != MatchCall || m.Given == nil {
				continue
			}
			ms := m.Scruts[0]
			outRet := ""
			if ret, ok := lintSigRet(mods, lm.name, ms.Fname); ok {
				outRet = ret
			}
			for key, sm := range m.Given {
				if sm == nil {
					continue
				}
				line := locateLineFrom(lm.text, key+" =>", m.Line, m.Line)
				items := []*Small{sm}
				if sm.Kind == "list" {
					items = sm.Items
				}
				for _, it := range items {
					if it.Kind != "exchange" {
						continue
					}
					checkCtorVal(it.Outcome, line, outRet, true)
					for _, a := range it.Args {
						checkCtorVal(a.V, line, "", false)
					}
				}
			}
		}
	}
	return out
}

// lintTops visits arm right-hand sides and plain expression
// bodies at their top Small (bodySmalls without the deep
// walk): the only body positions whose Ok proves.
func lintTops(n *Node, f func(s *Small, line int)) {
	var walk func(x *Node)
	walk = func(x *Node) {
		if x == nil {
			return
		}
		if x.IsMatch {
			for _, a := range x.Arms {
				walk(a.Rhs)
			}
			return
		}
		f(x.Small, x.Line)
	}
	walk(n)
}

// lintSigRet resolves a callee's return for the construction
// check: same discipline as lintSigParams, Ret instead of
// Params.
func lintSigRet(mods []lintModule, cur, fname string) (string, bool) {
	for _, m := range mods {
		if m.name != cur {
			continue
		}
		for _, d := range m.mod.Decls {
			if f, ok := d.(*FnDecl); ok && f.Name == fname {
				return f.Ret, true
			}
		}
		for _, d := range m.mod.Decls {
			if e, ok := d.(*ExternDecl); ok && e.Name == fname {
				return e.Ret, true
			}
		}
	}
	if k, ok := bytesKernels[fname]; ok {
		return k.ret, true
	}
	var found string
	n := 0
	for _, m := range mods {
		if m.name == cur {
			continue
		}
		for _, d := range m.mod.Decls {
			if f, ok := d.(*FnDecl); ok && f.Name == fname {
				found = f.Ret
				n++
			}
		}
	}
	if n == 1 {
		return found, true
	}
	return "", false
}

// lintCtorFields resolves a constructor's field list: records
// through their TypeDecl, variant cases through the qualified
// case, dotted errors through their ErrorDecl. Own module
// first, then a unique cross-module match — the lintSigParams
// discipline. Ok and Bytes never resolve here: Ok wants its
// return record (caller-supplied), Bytes takes no names.
func lintCtorFields(mods []lintModule, cur, ctor string) ([][2]string, bool) {
	if ctor == "Ok" || ctor == "Bytes" {
		return nil, false
	}
	for _, m := range mods {
		if m.name != cur {
			continue
		}
		if f, ok := lintDeclFields(m.mod, ctor); ok {
			return f, true
		}
	}
	var found [][2]string
	n := 0
	for _, m := range mods {
		if m.name == cur {
			continue
		}
		if f, ok := lintDeclFields(m.mod, ctor); ok {
			found = f
			n++
		}
	}
	if n == 1 {
		return found, true
	}
	return nil, false
}

// lintOkFields resolves the provable Ok shape: the return
// record's fields, or the scalar singleton `value` (a92).
// Anything else (unknown, unresolvable) stays silent.
func lintOkFields(mods []lintModule, cur, ret string) ([][2]string, bool) {
	if ret == "" {
		return nil, false
	}
	if fields, ok := lintCtorFields(mods, cur, ret); ok {
		return fields, true
	}
	if ret == "int" || ret == "str" || ret == "bool" || ret == "dec" {
		return [][2]string{{"value", ret}}, true
	}
	return nil, false
}

func lintDeclFields(m *Module, ctor string) ([][2]string, bool) {
	for _, d := range m.Decls {
		switch d := d.(type) {
		case *TypeDecl:
			if d.Name == ctor {
				return d.Fields, true
			}
		case *VariantDecl:
			for _, c := range d.Cases {
				if qualifyCase(d.Name, c.Short) == ctor {
					return c.Fields, true
				}
			}
		case *ErrorDecl:
			if d.Name == ctor {
				return d.Fields, true
			}
		}
	}
	return nil, false
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
// The differing slot comes back for the squiggle: the fold
// starts at that atom on the first arm's row.
func lintMergeableGroup(arms []lintArm, g []int) (bool, int) {
	if len(g) < 2 {
		return false, -1
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
				return false, -1
			}
			differ = s
		}
	}
	if differ == -1 {
		return false, -1
	}
	kind := arms[g[0]].pats[differ].Kind
	for _, gi := range g {
		p := arms[gi].pats[differ]
		if p.Kind != kind || !lintAtomKind(p) {
			return false, -1
		}
	}
	for x := 0; x < len(g); x++ {
		for y := x + 1; y < len(g); y++ {
			if patEqual(arms[g[x]].pats[differ], arms[g[y]].pats[differ]) {
				return false, -1
			}
		}
	}
	return true, differ
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
// reads, no trapping division or remainder. Everything else —
// refs, literals, total arithmetic, ctors, strlen — is total.
// Unknown kinds default closed (impure): a new faulting
// operator must opt into purity explicitly. The / and %
// exemption mirrors firstTrappingOp exactly (a literal nonzero
// divisor cannot fail), so same-outcome findings never advise
// dropping a match whose scrutinee can fault.
func lintPureSmall(s *Small) bool {
	pure := true
	walkSmallTrees(s, func(c *Small) {
		if c.Kind == "binop" && (c.Op == "/" || c.Op == "%") && !isNonzeroIntLit(c.R) {
			pure = false
			return
		}
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

// lintIsWord reports whether b can sit inside an identifier-ish
// token: letters, digits, underscore, and the dot of dotted
// names (so `failed` never matches inside `relay.failed`).
func lintIsWord(b byte) bool {
	if b == '_' || b == '.' {
		return true
	}
	if b >= '0' && b <= '9' {
		return true
	}
	if b >= 'a' && b <= 'z' {
		return true
	}
	if b >= 'A' && b <= 'Z' {
		return true
	}
	return false
}

// lintFindBounded locates needle in hay with identifier
// boundaries on both sides: `1` never matches inside `15`,
// `match` never inside `rematch`. Returns the byte offset or
// -1. hay must already be comment-stripped; strings are NOT
// skipped (string literals are legitimate targets).
func lintFindBounded(hay, needle string) int {
	if needle == "" {
		return -1
	}
	for off := 0; off+len(needle) <= len(hay); {
		i := strings.Index(hay[off:], needle)
		if i < 0 {
			return -1
		}
		i += off
		if i > 0 && lintIsWord(hay[i-1]) {
			off = i + 1
			continue
		}
		if j := i + len(needle); j < len(hay) && lintIsWord(hay[j]) {
			off = i + 1
			continue
		}
		return i
	}
	return -1
}

// lintFindBare locates word outside string literals with
// identifier boundaries: the `match` keyword search must not
// hit a "match" string in the scrutinee. from limits the
// search to code[from:]. Returns the byte offset or -1.
func lintFindBare(code, word string, from int) int {
	if word == "" || from < 0 {
		return -1
	}
	if from > len(code) {
		return -1
	}
	inStr := false
	for i := from; i < len(code); {
		ch := code[i]
		if inStr {
			if ch == '\\' && i+1 < len(code) {
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
		if strings.HasPrefix(code[i:], word) &&
			(i == 0 || !lintIsWord(code[i-1])) {
			if j := i + len(word); j >= len(code) || !lintIsWord(code[j]) {
				return i
			}
		}
		i++
	}
	return -1
}

// lintFindNameEq locates the (skip+1)-th `name =` on comment-
// stripped code, outside string literals: the redundant
// argument name through its equals sign. The `=` must not open
// `=>` or `==`, and name must stand alone. skip counts
// earlier same-named args, so `f(value = 1) => Ok(value = 1)`
// underlines the argument, never the right-hand side.
// Returns byte offsets [ns, ne).
func lintFindNameEq(code, name string, skip int) (ns, ne int, ok bool) {
	if name == "" || skip < 0 {
		return 0, 0, false
	}
	inStr := false
	for i := 0; i < len(code); {
		ch := code[i]
		if inStr {
			if ch == '\\' && i+1 < len(code) {
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
		if strings.HasPrefix(code[i:], name) &&
			(i == 0 || !lintIsWord(code[i-1])) {
			j := i + len(name)
			for j < len(code) && (code[j] == ' ' || code[j] == '\t') {
				j++
			}
			if j < len(code) && code[j] == '=' &&
				(j+1 >= len(code) || (code[j+1] != '=' && code[j+1] != '>')) {
				if skip == 0 {
					return i, j + 1, true
				}
				skip--
				i = j + 1
				continue
			}
		}
		i++
	}
	return 0, 0, false
}

// lintMatchParen returns the byte offset just past the paren
// closing the `(` at open, skipping string literals. -1 when
// the line never closes it.
func lintMatchParen(code string, open int) int {
	depth := 0
	inStr := false
	for i := open; i < len(code); i++ {
		ch := code[i]
		if inStr {
			if ch == '\\' && i+1 < len(code) {
				i++
				continue
			}
			if ch == '"' {
				inStr = false
			}
			continue
		}
		switch ch {
		case '"':
			inStr = true
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return -1
}

// lintFindCall locates the applied form `fname(` outside
// strings and returns the byte span through its matching
// close paren: the full call text. from limits the search to
// code[from:] (relay searches past `=>` so the pattern side,
// which names the same constructor, never wins).
func lintFindCall(code, fname string, from int) (ns, ne int, ok bool) {
	if fname == "" || from < 0 {
		return 0, 0, false
	}
	if from > len(code) {
		return 0, 0, false
	}
	inStr := false
	for i := from; i < len(code); {
		ch := code[i]
		if inStr {
			if ch == '\\' && i+1 < len(code) {
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
		if strings.HasPrefix(code[i:], fname) &&
			(i == 0 || !lintIsWord(code[i-1])) {
			if j := i + len(fname); j < len(code) && code[j] == '(' {
				if end := lintMatchParen(code, j); end > 0 {
					return i, end, true
				}
				return 0, 0, false
			}
		}
		i++
	}
	return 0, 0, false
}

// lintArrowOff returns the byte offset where an arm row's `=>`
// separator starts, or -1 when the row has none outside
// strings and comments.
func lintArrowOff(line string) int {
	rhs, ok := splitArmRHS(line)
	if !ok {
		return -1
	}
	return len(line) - len(rhs) - len("=>")
}

// lintPatSideSpan returns the UTF-16 span to underline on one
// arm row: the verbatim atom text when it sits on the pattern
// side (before `=>`, so right-hand sides repeating the
// literal never win), else the whole trimmed pattern side.
// ok=false keeps the whole-line fallback.
func lintPatSideSpan(raw, atom string) (s, e int, ok bool) {
	code := lintCodePart(raw)
	arrow := lintArrowOff(raw)
	if arrow < 0 || arrow > len(code) {
		return 0, 0, false
	}
	lhs := code[:arrow]
	if atom != "" {
		if i := lintFindBounded(lhs, atom); i >= 0 {
			return u16col(code, i), u16col(code, i+len(atom)), true
		}
	}
	s0 := 0
	for s0 < len(lhs) && (lhs[s0] == ' ' || lhs[s0] == '\t') {
		s0++
	}
	e0 := len(lhs)
	for e0 > s0 && (lhs[e0-1] == ' ' || lhs[e0-1] == '\t') {
		e0--
	}
	if e0 <= s0 {
		return 0, 0, false
	}
	return u16col(code, s0), u16col(code, e0), true
}

// lintHeadSpan underlines a match head on its own row: the
// `match` keyword through the end of code. String-aware, so a
// "match" string in the scrutinee never wins. ok=false keeps
// the whole-line fallback.
func lintHeadSpan(lines []string, ln int) (s, e int, ok bool) {
	if ln < 1 || ln > len(lines) {
		return 0, 0, false
	}
	code := lintCodePart(lines[ln-1])
	i := lintFindBare(code, "match", 0)
	if i < 0 {
		return 0, 0, false
	}
	e0 := len(code)
	for e0 > i && (code[e0-1] == ' ' || code[e0-1] == '\t') {
		e0--
	}
	if e0 <= i {
		return 0, 0, false
	}
	return u16col(code, i), u16col(code, e0), true
}

// lintPatText renders a discrete atom exactly as it reads in
// source, for verbatim span location: quoted spelling for
// strings (Raw), numerals for ints, the word for bools, the
// name for consts. ok=false asks the caller to underline the
// whole pattern side instead of pretending.
func lintPatText(p Pattern) (string, bool) {
	switch p.Kind {
	case "str":
		if p.Raw == "" {
			return "", false
		}
		return p.Raw, true
	case "int":
		if p.Num == nil {
			return "", false
		}
		return p.Num.String(), true
	case "bool":
		if p.B {
			return "true", true
		}
		return "false", true
	case "const":
		if p.Name == "" {
			return "", false
		}
		return p.Name, true
	}
	return "", false
}

// lintCallFname names the call a match-call rung scrutinizes:
// exactly one scrutinee, and it is applied.
func lintCallFname(n *Node) (string, bool) {
	if n == nil || len(n.Scruts) != 1 {
		return "", false
	}
	s := n.Scruts[0]
	if s == nil || s.Kind != "call" || s.Fname == "" {
		return "", false
	}
	return s.Fname, true
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
				merge, differ := lintMergeableGroup(arms, g)
				if !merge {
					return
				}
				ls := make([]int, len(g))
				for i, gi := range g {
					ls[i] = arms[gi].line
				}
				f := lintFinding{file: lm.name, line: ls[0], msg: fmt.Sprintf(
					"mergeable match arms %v (identical outcomes); fold into one or-pattern arm, saves %d lines",
					ls, len(g)-1), code: CodeLintOrFold}
				// Underline the differing atom on the first
				// arm's row (the fold starts there), else its
				// whole pattern side.
				atom, _ := lintPatText(arms[g[0]].pats[differ])
				if ls[0] >= 1 && ls[0] <= len(lines) {
					if s, e, ok := lintPatSideSpan(lines[ls[0]-1], atom); ok {
						f.start, f.end = s, e
					}
				}
				out = append(out, f)
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
		f := lintFinding{file: lm.name, line: r[0].Line, msg: fmt.Sprintf(
			"chainable ladder: %d sequential calls sharing one failure outcome; rewrite as match chain",
			len(r)), code: CodeLintChain}
		// Underline the first rung's full call text (the
		// ladder's anchor), else the callee name, else the
		// whole line.
		if fname, ok := lintCallFname(r[0]); ok {
			if ln := r[0].Line; ln >= 1 && ln <= len(lines) {
				code := lintCodePart(lines[ln-1])
				if ns, ne, ok := lintFindCall(code, fname, 0); ok {
					f.start, f.end = u16col(code, ns), u16col(code, ne)
				} else if s, e, ok := tokenSpan(lm.text, ln, fname); ok {
					f.start, f.end = s, e
				}
			}
		}
		out = append(out, f)
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
// would not (CAN4109) — and identical: different pure
// scrutinees fold under the ladder rule instead (don't-care
// slots, the lcm shape); only identical-scrutinee diamonds
// are this rule's.
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
			f := lintFinding{file: lm.name, line: m.Line, msg: "nested matches share one scrutinee; fold into a multi-scrutinee table", code: CodeLintTable}
			// Underline the outer match head: the decision
			// point the table replaces.
			if s, e, ok := lintHeadSpan(lines, m.Line); ok {
				f.start, f.end = s, e
			}
			out = append(out, f)
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
			f := lintFinding{file: lm.name, line: m.Line, msg: fmt.Sprintf(
				"match always yields %s; drop the match", rhs), code: CodeLintSameOutcome}
			// Underline the match head: the text to delete.
			if s, e, ok := lintHeadSpan(lines, m.Line); ok {
				f.start, f.end = s, e
			}
			out = append(out, f)
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
	lines := strings.Split(lm.text, "\n")
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
				for i, arg := range s.Args {
					name := arg.Name
					if !arg.HasName {
						// a92: positional relays spell without
						// names; the slot is the list index.
						if i >= len(want) {
							relay = false
							break
						}
						name = want[i]
					} else {
						known := false
						for _, f := range want {
							if f == name {
								known = true
								break
							}
						}
						if !known {
							relay = false
							break
						}
					}
					if seen[name] {
						relay = false
						break
					}
					seen[name] = true
					v := arg.V
					if v == nil || v.Kind != "ref" || len(v.Ref) != 2 ||
						v.Ref[0] != pat.Var || v.Ref[1] != name {
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
					f := lintFinding{file: lm.name, line: a.Line, msg: fmt.Sprintf(
						"handwritten relay of %s; write forward %s", pat.Name, pat.Var), code: CodeLintRelay}
					// Underline the handwritten right-hand
					// side: the ctor text `forward` replaces.
					// Past `=>`, so the pattern side naming
					// the same constructor never wins.
					if ln := a.Line; ln >= 1 && ln <= len(lines) {
						raw := lines[ln-1]
						code := lintCodePart(raw)
						arrow := lintArrowOff(raw)
						if arrow >= 0 && arrow <= len(code) {
							if i := lintFindBare(code, s.Ctor, arrow+len("=>")); i >= 0 {
								j := i + len(s.Ctor)
								if j < len(code) && code[j] == '(' {
									if end := lintMatchParen(code, j); end > 0 {
										f.start, f.end = u16col(code, i), u16col(code, end)
									}
								} else {
									f.start, f.end = u16col(code, i), u16col(code, j)
								}
							}
						}
					}
					out = append(out, f)
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
				f := lintFinding{file: lm.name, line: chain[0].line, msg: fmt.Sprintf(
					"mergeable ranges %s; join into %s..%s",
					strings.Join(parts, ", "),
					chain[0].lo.String(), chain[len(chain)-1].hi.String()), code: CodeLintRange}
				// Underline the first range literal (the join
				// starts there), else its whole pattern side.
				want := chain[0].lo.String() + ".." + chain[0].hi.String()
				if ln := chain[0].line; ln >= 1 && ln <= len(lines) {
					if s, e, ok := lintPatSideSpan(lines[ln-1], want); ok {
						f.start, f.end = s, e
					}
				}
				out = append(out, f)
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

// lintSmallEqual reports whether two Smalls are structurally
// identical: same shape and leaves throughout. There are no
// paren nodes and lint runs pre-check (types unset), so the
// comparison is over parse structure only. Update with the
// Small struct: a field added there and missed here silently
// compares unequal.
func lintSmallEqual(a, b *Small) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Kind != b.Kind || a.Op != b.Op || a.Str != b.Str ||
		a.B != b.B || a.Dec != b.Dec || a.T != b.T ||
		a.Seal != b.Seal || a.Elem != b.Elem || a.Fname != b.Fname ||
		a.Ctor != b.Ctor || a.ExportBrand != b.ExportBrand {
		return false
	}
	if (a.Num == nil) != (b.Num == nil) {
		return false
	}
	if a.Num != nil && a.Num.Cmp(b.Num) != 0 {
		return false
	}
	if len(a.Ref) != len(b.Ref) {
		return false
	}
	for i := range a.Ref {
		if a.Ref[i] != b.Ref[i] {
			return false
		}
	}
	if !lintSmallEqual(a.L, b.L) || !lintSmallEqual(a.R, b.R) ||
		!lintSmallEqual(a.Hi, b.Hi) || !lintSmallEqual(a.Outcome, b.Outcome) {
		return false
	}
	if len(a.Args) != len(b.Args) {
		return false
	}
	for i := range a.Args {
		if a.Args[i].Name != b.Args[i].Name || a.Args[i].HasName != b.Args[i].HasName ||
			!lintSmallEqual(a.Args[i].V, b.Args[i].V) {
			return false
		}
	}
	if len(a.Items) != len(b.Items) {
		return false
	}
	for i := range a.Items {
		if !lintSmallEqual(a.Items[i], b.Items[i]) {
			return false
		}
	}
	return true
}

// lintRung splits one ==-rung: a single-scrutinee value match
// testing base == "lit" with true/false arms in either order.
// The literal may sit on either side; both-literal comparisons
// are constant conditions (dead arms elsewhere) and never rungs.
func lintRung(n *Node) (base *Small, lit string, falseArm *Arm, ok bool) {
	if n == nil || !n.IsMatch || n.Kind != MatchValue || n.Given != nil {
		return nil, "", nil, false
	}
	if len(n.Scruts) != 1 || len(n.Arms) != 2 {
		return nil, "", nil, false
	}
	s := n.Scruts[0]
	if s == nil || s.Kind != "binop" || s.Op != "==" {
		return nil, "", nil, false
	}
	base = s.L
	if s.L != nil && s.L.Kind == "str" {
		base, lit = s.R, s.L.Str
	} else if s.R != nil && s.R.Kind == "str" {
		lit = s.R.Str
	} else {
		return nil, "", nil, false
	}
	if base == nil || base.Kind == "str" {
		return nil, "", nil, false
	}
	var f *Arm
	seenTrue := false
	for i := range n.Arms {
		a := &n.Arms[i]
		if len(a.Pats) != 1 || a.Pats[0].Kind != "bool" {
			return nil, "", nil, false
		}
		if a.Pats[0].B {
			if seenTrue {
				return nil, "", nil, false
			}
			seenTrue = true
			continue
		}
		if f != nil {
			return nil, "", nil, false
		}
		f = a
	}
	if !seenTrue || f == nil {
		return nil, "", nil, false
	}
	return base, lit, f, true
}

// lintRestatable implements the equality-ladder rule: a chain of
// two or more ==-rungs over one base with distinct string
// literals restates as a single match on the base with one arm
// per literal plus _. The restatement is fault-neutral — the
// table evaluates the base once at the same program point the
// first rung would — so purity is not required, and bodies move
// verbatim (no outcome comparison). Only false-nesting ladders
// qualify: an inverted ladder over distinct literals always
// carries a dead arm, which is the checker's business.
func lintRestatable(lm lintModule) []lintFinding {
	lines := strings.Split(lm.text, "\n")
	var out []lintFinding
	for _, d := range lm.mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		nodes := matchNodes(fn.Body)
		chains := map[*Node][]*Node{}
		for _, m := range nodes {
			base, lit, f, ok := lintRung(m)
			if !ok {
				continue
			}
			chain := []*Node{m}
			seen := map[string]bool{lit: true}
			next := f.Rhs
			for next != nil {
				nbase, nlit, nf, ok := lintRung(next)
				if !ok || !lintSmallEqual(nbase, base) || seen[nlit] {
					break
				}
				chain = append(chain, next)
				seen[nlit] = true
				next = nf.Rhs
			}
			chains[m] = chain
		}
		// Maximal ladders only: a rung continuing another rung's
		// chain never heads its own report.
		continued := map[*Node]bool{}
		for _, c := range chains {
			for _, n := range c[1:] {
				continued[n] = true
			}
		}
		for _, m := range nodes {
			c := chains[m]
			if len(c) < 2 || continued[m] {
				continue
			}
			f := lintFinding{file: lm.name, line: m.Line, msg: fmt.Sprintf(
				"equality ladder over one base with %d rungs; restate as a match on the base, saves %d lines",
				len(c), len(c)-1), code: CodeLintRestate}
			// Underline the ladder head: the decision point
			// the table replaces.
			if s, e, ok := lintHeadSpan(lines, m.Line); ok {
				f.start, f.end = s, e
			}
			out = append(out, f)
		}
	}
	return out
}

// lintPlainArm reports whether an arm yields a single expression:
// no nested match of either family. The ladder rule folds pure
// decisions only — arms that run calls first stay nested
// (can-idioms C1) — so only plain bodies move into tables.
func lintPlainArm(a *Arm) bool {
	return a.Rhs != nil && !a.Rhs.IsMatch && a.Rhs.Small != nil
}

// lintBoolArms splits a two-arm bool match in either arm order.
func lintBoolArms(m *Node) (tr, fl *Arm, ok bool) {
	if len(m.Arms) != 2 {
		return nil, nil, false
	}
	for i := range m.Arms {
		a := &m.Arms[i]
		if len(a.Pats) != 1 || a.Pats[0].Kind != "bool" {
			return nil, nil, false
		}
		if a.Pats[0].B {
			if tr != nil {
				return nil, nil, false
			}
			tr = a
		} else {
			if fl != nil {
				return nil, nil, false
			}
			fl = a
		}
	}
	return tr, fl, true
}

// lintValueNest returns the directly-nested value match of an arm
// body, or nil. Call matches and chains never nest into tables.
func lintValueNest(a *Arm) *Node {
	if a.Rhs != nil && a.Rhs.IsMatch && a.Rhs.Kind == MatchValue && a.Rhs.Given == nil {
		return a.Rhs
	}
	return nil
}

// lintBoolLevel splits one ladder level: a single-scrutinee bool
// value match over a pure scrutinee. Tables evaluate eagerly,
// so a faulting scrutinee must keep its guard (CAN4109).
func lintBoolLevel(m *Node) (scrut *Small, tr, fl *Arm, ok bool) {
	if m == nil || !m.IsMatch || m.Kind != MatchValue || m.Given != nil {
		return nil, nil, nil, false
	}
	if len(m.Scruts) != 1 || m.Scruts[0] == nil || !lintPureSmall(m.Scruts[0]) {
		return nil, nil, nil, false
	}
	tr, fl, ok = lintBoolArms(m)
	if !ok {
		return nil, nil, nil, false
	}
	return m.Scruts[0], tr, fl, true
}

// lintLadderChain walks a pure ladder from its head: every level
// nests exactly one arm deeper while the other yields plainly,
// down to a tail level whose arms both yield. All scrutinees
// must be pairwise different — a repeated scrutinee is
// statically decided below its first test (degenerate, the
// checker's business). Returns the scrutinee count and the
// spine below the head.
func lintLadderChain(m *Node) (int, []*Node, bool) {
	scrut, tr, fl, ok := lintBoolLevel(m)
	if !ok {
		return 0, nil, false
	}
	scruts := []*Small{scrut}
	var spine []*Node
	curTr, curFl := tr, fl
	for {
		trNest, flNest := lintValueNest(curTr), lintValueNest(curFl)
		if trNest != nil && flNest != nil {
			return 0, nil, false // diamond, not a ladder
		}
		var next *Node
		switch {
		case trNest != nil:
			if !lintPlainArm(curFl) {
				return 0, nil, false
			}
			next = trNest
		case flNest != nil:
			if !lintPlainArm(curTr) {
				return 0, nil, false
			}
			next = flNest
		default:
			if !lintPlainArm(curTr) || !lintPlainArm(curFl) {
				return 0, nil, false
			}
			if len(scruts) < 2 {
				return 0, nil, false // lone level, nothing to fold
			}
			return len(scruts), spine, true
		}
		nscrut, ntr, nfl, ok := lintBoolLevel(next)
		if !ok {
			return 0, nil, false
		}
		for _, s := range scruts {
			if lintSmallEqual(s, nscrut) {
				return 0, nil, false
			}
		}
		scruts = append(scruts, nscrut)
		spine = append(spine, next)
		curTr, curFl = ntr, nfl
	}
}

// lintDiamondDiff splits a diamond over different scrutinees:
// every arm nests directly into a leaf bool level (both arms
// plain) over a pure scrutinee, and no two scrutinees in the
// shape are equal. Identical-scrutinee diamonds belong to rule
// 4. Returns the scrutinee count and the inner levels.
func lintDiamondDiff(m *Node) (int, []*Node, bool) {
	scrut, tr, fl, ok := lintBoolLevel(m)
	if !ok {
		return 0, nil, false
	}
	ti, fi := lintValueNest(tr), lintValueNest(fl)
	if ti == nil || fi == nil {
		return 0, nil, false
	}
	scruts := []*Small{scrut}
	var spine []*Node
	for _, inner := range []*Node{ti, fi} {
		nscrut, ntr, nfl, ok := lintBoolLevel(inner)
		if !ok || !lintPlainArm(ntr) || !lintPlainArm(nfl) {
			return 0, nil, false
		}
		for _, s := range scruts {
			if lintSmallEqual(s, nscrut) {
				return 0, nil, false
			}
		}
		scruts = append(scruts, nscrut)
		spine = append(spine, inner)
	}
	return len(scruts), spine, true
}

// lintLadderable implements the pure-nest rule: ladders and
// diamonds over pure, pairwise-different scrutinees with plain
// outcomes fold into one multi-scrutinee table with don't-care
// slots (the std__int__lcm shape). Disjoint slot-1 values make
// arm order irrelevant, so yielding arms keep their priority in
// any position. Maximal shapes only: a level continuing another
// shape's spine never heads its own report.
func lintLadderable(lm lintModule) []lintFinding {
	lines := strings.Split(lm.text, "\n")
	var out []lintFinding
	for _, d := range lm.mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		nodes := matchNodes(fn.Body)
		type shape struct {
			count int
			spine []*Node
		}
		shapes := map[*Node]shape{}
		for _, m := range nodes {
			if n, spine, ok := lintLadderChain(m); ok {
				shapes[m] = shape{n, spine}
				continue
			}
			if n, spine, ok := lintDiamondDiff(m); ok {
				shapes[m] = shape{n, spine}
			}
		}
		continued := map[*Node]bool{}
		for _, s := range shapes {
			for _, n := range s.spine {
				continued[n] = true
			}
		}
		for _, m := range nodes {
			s, ok := shapes[m]
			if !ok || continued[m] {
				continue
			}
			f := lintFinding{file: lm.name, line: m.Line, msg: fmt.Sprintf(
				"nested matches over %d pure scrutinees; fold into a multi-scrutinee table",
				s.count), code: CodeLintLadder}
			// Underline the nest head: the decision point the
			// table replaces.
			if st, en, ok := lintHeadSpan(lines, m.Line); ok {
				f.start, f.end = st, en
			}
			out = append(out, f)
		}
	}
	return out
}

// lintMatchScrut extracts the scrutinee text of a match head row:
// the source after "match", comments stripped, whitespace
// normalized (same normalization as lintInnerKey).
func lintMatchScrut(lines []string, m *Node) (string, bool) {
	if m.Line < 1 || m.Line > len(lines) {
		return "", false
	}
	row := lines[m.Line-1]
	i := strings.Index(row, "match ")
	if i < 0 {
		return "", false
	}
	fields := strings.Fields(lintCodePart(row[i+len("match "):]))
	if len(fields) == 0 {
		return "", false
	}
	return strings.Join(fields, " "), true
}

// lintRelayCallable implements the relay-call rule: a call match
// with no given table, over a same-file local callee, whose
// every arm forwards its own binder rewrites as one forward
// call relay. Foreign, uses-pinned, and extern callees keep
// their match (given tables cannot attach to the relay shape),
// as does any arm that rebuilds instead of forwarding. Unknown
// callees stay silent: broken code gets check errors, not lint
// advice.
func lintRelayCallable(lm lintModule) []lintFinding {
	lines := strings.Split(lm.text, "\n")
	locals := map[string]bool{}
	externs := map[string]bool{}
	for _, d := range lm.mod.Decls {
		switch t := d.(type) {
		case *FnDecl:
			locals[t.Name] = true
		case *ExternDecl:
			externs[t.Name] = true
		}
	}
	pinned := map[string]bool{}
	for _, u := range lm.mod.Hdr["uses"] {
		name := u
		if i := strings.IndexByte(name, '@'); i >= 0 {
			name = name[:i]
		}
		pinned[strings.TrimSpace(name)] = true
	}
	var out []lintFinding
	for _, d := range lm.mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, m := range matchNodes(fn.Body) {
			if !m.IsMatch || m.Kind != MatchCall || m.Given != nil {
				continue
			}
			if len(m.Scruts) != 1 || m.Scruts[0] == nil || m.Scruts[0].Kind != "call" {
				continue
			}
			fname := m.Scruts[0].Fname
			if !locals[fname] || pinned[fname] || externs[fname] {
				continue
			}
			if len(m.Arms) == 0 {
				continue
			}
			all := true
			for i := range m.Arms {
				a := &m.Arms[i]
				if len(a.Pats) != 1 || a.Pats[0].Kind != "variant" {
					all = false
					break
				}
				binder := a.Pats[0].Var
				if binder == "" || binder == "_" || a.Rhs == nil || a.Rhs.IsMatch ||
					a.Rhs.Small == nil || a.Rhs.Small.Kind != "forward" ||
					a.Rhs.Small.Str != binder {
					all = false
					break
				}
			}
			if !all {
				continue
			}
			call, ok := lintMatchScrut(lines, m)
			if !ok {
				continue
			}
			if len(call) > 72 {
				call = call[:72]
			}
			f := lintFinding{file: lm.name, line: m.Line, msg: "relay match forwards every outcome; write forward " + call, code: CodeLintRelayCall}
			// Underline the match head: the site the relay
			// replaces.
			if s, e, ok := lintHeadSpan(lines, m.Line); ok {
				f.start, f.end = s, e
			}
			out = append(out, f)
		}
	}
	return out
}

// lintUnreachedKeys implements the stale-script rule (a91): a
// given key naming a real test with no static path to the call
// site never selects — scripts run under the executing test's
// name — so the row is dead weight omission was meant to shed.
// Unknown names stay the compiler's warning (CAN3104) and
// retired `-` rows its error (CAN3111); this rule only flags
// keys that name a real test elsewhere. Sound but partial: no
// static path means truly unreachable (every call is static),
// while a static path may still never execute. Chain-step
// tables are covered like call-match tables.
func lintUnreachedKeys(lm lintModule) []lintFinding {
	lines := strings.Split(lm.text, "\n")
	fns := map[string]*FnDecl{}
	testsOf := map[string][]string{}
	var order []string
	for _, d := range lm.mod.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		fns[fn.Name] = fn
		order = append(order, fn.Name)
		for _, t := range fn.Tests {
			testsOf[fn.Name] = append(testsOf[fn.Name], t.Name)
		}
	}
	fileTests := map[string]bool{}
	for _, names := range testsOf {
		for _, n := range names {
			fileTests[n] = true
		}
	}
	callers := map[string][]string{}
	for _, name := range order {
		seen := map[string]bool{}
		for _, callee := range lintCalledNames(fns[name].Body) {
			if _, ok := fns[callee]; !ok || seen[callee] {
				continue
			}
			seen[callee] = true
			callers[callee] = append(callers[callee], name)
		}
	}
	var out []lintFinding
	for _, name := range order {
		selectable := map[string]bool{}
		for _, t := range testsOf[name] {
			selectable[t] = true
		}
		seen := map[string]bool{name: true}
		queue := []string{name}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, caller := range callers[u] {
				if seen[caller] {
					continue
				}
				seen[caller] = true
				queue = append(queue, caller)
				for _, t := range testsOf[caller] {
					selectable[t] = true
				}
			}
		}
		for _, site := range lintGivenSites(fns[name].Body) {
			for key, sm := range site.table {
				if sm == nil || selectable[key] || !fileTests[key] {
					continue
				}
				ln := locateLineFrom(lm.text, key+" =>", site.line, site.line)
				f := lintFinding{file: lm.name, line: ln, msg: fmt.Sprintf(
					"given key %q names a test that cannot reach this call; delete the row", key),
					code: CodeLintUnreachedKey}
				if s, e, ok := lintGivenKeySpan(lines, ln, key); ok {
					f.start, f.end = s, e
				}
				out = append(out, f)
			}
		}
	}
	return out
}

// lintGivenSite is one given table with the source line anchoring
// it: the match line for call tables, the step line for chains.
type lintGivenSite struct {
	table map[string]*Small
	line  int
}

// lintGivenSites collects every given table in a body: call-match
// tables plus chain-step tables, innermost included.
func lintGivenSites(n *Node) []lintGivenSite {
	var out []lintGivenSite
	var walk func(x *Node)
	walk = func(x *Node) {
		if x == nil {
			return
		}
		if x.IsMatch {
			if x.Given != nil {
				out = append(out, lintGivenSite{x.Given, x.Line})
			}
			for _, s := range x.ChainSteps {
				if s.Given != nil {
					out = append(out, lintGivenSite{s.Given, s.Line})
				}
			}
			walk(x.ChainTail)
			for _, a := range x.Arms {
				walk(a.Rhs)
			}
		}
	}
	walk(n)
	return out
}

// lintCalledNames collects every called name in a body, descending
// into chain steps, guards, and tails where walkCalls stops: a
// missed caller edge would become a false finding.
func lintCalledNames(n *Node) []string {
	seen := map[string]bool{}
	var out []string
	add := func(fname string) {
		if !seen[fname] {
			seen[fname] = true
			out = append(out, fname)
		}
	}
	for _, c := range walkCalls(n) {
		add(c.Fname)
	}
	var walkChain func(x *Node)
	walkChain = func(x *Node) {
		if x == nil {
			return
		}
		for _, s := range x.ChainSteps {
			for _, c := range walkCalls(&Node{Small: s.Call}) {
				add(c.Fname)
			}
			for _, c := range walkCalls(&Node{Small: s.Guard}) {
				add(c.Fname)
			}
		}
		walkChain(x.ChainTail)
		for _, a := range x.Arms {
			walkChain(a.Rhs)
		}
	}
	walkChain(n)
	return out
}

// lintGivenKeySpan underlines the key token on its `key => ...`
// row: leading whitespace, the exact key, then `=>`. Anything
// else stays whole-line.
func lintGivenKeySpan(lines []string, ln int, key string) (s, e int, ok bool) {
	if ln < 1 || ln > len(lines) {
		return 0, 0, false
	}
	code := lintCodePart(lines[ln-1])
	i := 0
	for i < len(code) && (code[i] == ' ' || code[i] == '\t') {
		i++
	}
	if !strings.HasPrefix(code[i:], key) {
		return 0, 0, false
	}
	j := i + len(key)
	if j < len(code) && (code[j] == '_' ||
		(code[j] >= '0' && code[j] <= '9') ||
		(code[j] >= 'a' && code[j] <= 'z') ||
		(code[j] >= 'A' && code[j] <= 'Z')) {
		return 0, 0, false
	}
	k := j
	for k < len(code) && (code[k] == ' ' || code[k] == '\t') {
		k++
	}
	if !strings.HasPrefix(code[k:], "=>") {
		return 0, 0, false
	}
	return u16col(code, i), u16col(code, j), true
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
		out = append(out, Diag{File: name, Line: f.line, Sev: "error", Msg: f.msg, Code: f.code,
			Start: f.start, End: f.end})
	}
	return out
}
