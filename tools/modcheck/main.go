// Command modcheck is the can-lang module check: every `uses` must resolve
// to another file's `provides` (with a pinned rev), and the retired shapes
// (inline `extern`, `externals` sections, curly braces, partial `given`
// tables) must not reappear. Files marked SUPERSEDED-BY are skipped.
// Demo files in sketches/broken-login/ may carry `// DEMO-EXPECTS: <msg>`
// lines declaring the exact violation they demonstrate; modcheck passes
// them only when they produce exactly those errors and nothing else, so a
// demo can never silently heal or silently rot.
//
// Run from anywhere inside the repo: go run ./tools/modcheck
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/veighnsche/can-lang/internal/scan"
)

// Example programs, each folder named by its contents.
var examples = []string{}

var (
	modRe      = regexp.MustCompile(`(?s)mod\s+(\w+)\s+provides\s*\[(.*?)\]\s+uses\s*\[(.*?)\]`)
	externRe   = regexp.MustCompile(`extern\s+fn\s+(\w+)`)
	pinRe      = regexp.MustCompile(`@\d+$`)
	externsRe  = regexp.MustCompile(`(?m)^\s*externals\s*:`)
	testsRe    = regexp.MustCompile(`(?m)^\s*tests\s*$`)
	eqRe       = regexp.MustCompile(`^=\s*$`)
	caseRe     = regexp.MustCompile(`^\s*(\w+)\(`)
	givenRe    = regexp.MustCompile(`(?m)^\s*given\s*$`)
	givenKeyRe = regexp.MustCompile(`(?m)^\s*(\w+)\s*=>`)
	fnStartRe  = regexp.MustCompile(`(?m)^fn (\w+)`)
	callRe     = regexp.MustCompile(`\bcall\s+(\w+)\s*\(`)
	demoRe     = regexp.MustCompile(`(?m)^\s*//\s*DEMO-EXPECTS:\s*(.+?)\s*$`)
)

func names(raw string) []string {
	var out []string
	for _, p := range strings.Split(raw, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

type owner struct {
	mod, file string
}

// check scans roots (repo layout: roots hold *.can plus one level of
// example dirs) and returns current/skipped counts, scanned dir names,
// and every violation. Pure filesystem in, strings out — see main_test.go.
func check(roots []string) (current, skipped int, scanned []string, errs []string) {
	var files []string
	for _, root := range roots {
		top, _ := filepath.Glob(filepath.Join(root, "*.can"))
		sub, _ := filepath.Glob(filepath.Join(root, "*", "*.can"))
		files = append(files, append(top, sub...)...)
	}
	sort.Strings(files)
	anchor := filepath.Dir(roots[0])
	keyOf := func(path string) string {
		if rel, err := filepath.Rel(anchor, path); err == nil {
			return rel
		}
		return filepath.Base(path)
	}
	// Each directory is its own program: names resolve within it.
	// Slice 1: stdlib shares providers across its directories
	// (ascii consts), so a uses entry also resolves to a
	// provider elsewhere under the same root. Self-provides
	// still resolve nowhere, exactly like the compiler.
	groups := map[string][]string{}
	rootOf := map[string]string{}
	for _, root := range roots {
		top, _ := filepath.Glob(filepath.Join(root, "*.can"))
		sub, _ := filepath.Glob(filepath.Join(root, "*", "*.can"))
		for _, p := range append(top, sub...) {
			rootOf[p] = root
		}
	}
	rootProvides := map[string]map[string][]string{}
	for _, path := range files {
		body, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		m := modRe.FindStringSubmatch(string(body))
		if m == nil {
			continue
		}
		shared := rootProvides[rootOf[path]]
		if shared == nil {
			shared = map[string][]string{}
			rootProvides[rootOf[path]] = shared
		}
		for _, name := range names(m[2]) {
			shared[name] = append(shared[name], keyOf(path))
		}
	}
	var order []string
	for _, path := range files {
		dir := filepath.Dir(path)
		if _, ok := groups[dir]; !ok {
			order = append(order, dir)
		}
		groups[dir] = append(groups[dir], path)
	}
	for _, dir := range order {
		scanned = append(scanned, filepath.Base(dir))
		c, s, es := checkGroup(groups[dir], keyOf, rootProvides[rootOf[groups[dir][0]]])
		current += c
		skipped += s
		errs = append(errs, es...)
	}
	sort.Strings(errs)
	return current, skipped, scanned, errs
}

// checkGroup enforces the module rules within one program directory.
// shared is the root-wide provider index (slice 1): a uses entry
// resolves to another file in its directory or anywhere else under
// the same root, but never to its own file.
func checkGroup(files []string, keyOf func(string) string, shared map[string][]string) (current, skipped int, errs []string) {
	provides := map[string][]owner{}
	uses := map[string][][2]string{}
	externs := map[string][]string{}
	skippedSet := map[string]bool{}
	textOf := map[string]string{}
	buckets := map[string][]string{}
	expected := map[string][]string{}
	add := func(base, format string, args ...any) {
		buckets[base] = append(buckets[base], fmt.Sprintf(base+": "+format, args...))
	}
	for _, path := range files {
		base := keyOf(path)
		text, err := os.ReadFile(path)
		if err != nil {
			add(base, "%v", err)
			continue
		}
		body := string(text)
		textOf[base] = body
		if strings.Contains(body, "SUPERSEDED-BY") {
			skippedSet[base] = true
			skipped++
			continue
		}
		for _, mm := range demoRe.FindAllStringSubmatch(body, -1) {
			expected[base] = append(expected[base], mm[1])
		}
		m := modRe.FindStringSubmatch(body)
		if m == nil {
			add(base, "no parseable mod header")
			continue
		}
		var bare []string
		for _, entry := range names(m[3]) {
			b := pinRe.ReplaceAllString(entry, "")
			bare = append(bare, b)
			if b == entry {
				add(base, "uses %s must pin a rev (name@N)", entry)
			}
		}
		for _, name := range names(m[2]) {
			provides[name] = append(provides[name], owner{m[1], base})
		}
		for _, dep := range bare {
			uses[base] = append(uses[base], [2]string{m[1], dep})
		}
		for _, mm := range externRe.FindAllStringSubmatch(body, -1) {
			externs[base] = append(externs[base], mm[1])
		}
	}
	for name, owners := range provides {
		if len(owners) > 1 {
			errs = append(errs, fmt.Sprintf("%s provided twice: %v", name, owners))
		}
	}
	providedFiles := map[string][]string{}
	for name, owners := range provides {
		for _, o := range owners {
			providedFiles[name] = append(providedFiles[name], o.file)
		}
	}
	for file, deps := range uses {
		for _, dep := range deps {
			others := false
			for _, f := range providedFiles[dep[1]] {
				if f != file {
					others = true
				}
			}
			if !others {
				for _, f := range shared[dep[1]] {
					if f != file {
						others = true
					}
				}
			}
			if !others {
				add(file, "uses %s resolves nowhere", dep[1])
			}
		}
	}
	for file, decls := range externs {
		for _, decl := range decls {
			for _, deps := range uses {
				for _, dep := range deps {
					if dep[1] == decl {
						add(file, "extern fn %s re-declares an can-provided name", decl)
					}
				}
			}
		}
	}
	for _, path := range files {
		base := keyOf(path)
		if skippedSet[base] {
			continue
		}
		body, ok := textOf[base]
		if !ok {
			continue
		}
		if externsRe.MatchString(body) {
			add(base, "externals section is gone, use call-site given")
		}
		if scan.HasBraceOutsideString(body) {
			add(base, "curly braces are banned, use () records")
		}
		testNames := map[string]bool{}
		inTests := false
		for _, line := range strings.Split(body, "\n") {
			switch {
			case testsRe.MatchString(line):
				inTests = true
			case inTests && eqRe.MatchString(line):
				inTests = false
			default:
				if inTests {
					if m := caseRe.FindStringSubmatch(line); m != nil {
						testNames[m[1]] = true
					}
				}
			}
		}
		// Reaching tests per function: a given block must script
		// its own fn's rows plus same-file transitive callers'
		// rows (mirroring check.go reachingTests), not every row
		// in the file. Keys outside the file's rows are still
		// rejected. Single-fn files behave exactly as before.
		spans := fnSpans(body)
		for _, loc := range givenRe.FindAllStringIndex(body, -1) {
			block := body[loc[0]:]
			if i := strings.Index(block, "on "); i >= 0 {
				block = block[:i]
			}
			keys := map[string]bool{}
			for _, m := range givenKeyRe.FindAllStringSubmatch(block, -1) {
				keys[m[1]] = true
			}
			required := testNames
			if fn := enclosingFn(spans, loc[0]); fn != "" {
				required = reachingRows(body, spans, fn)
			}
			if len(required) > 0 && (!subsetOf(required, keys) || !subsetOf(keys, testNames)) {
				add(base, "given table %v != tests %v",
					sortedKeys(keys), sortedKeys(required))
			}
		}
	}
	for base, want := range expected {
		got := map[string]bool{}
		for _, full := range buckets[base] {
			got[strings.TrimPrefix(full, base+": ")] = true
		}
		for _, w := range want {
			if !got[w] {
				errs = append(errs, fmt.Sprintf("%s: demo expects %q but it no longer happens", base, w))
			}
			delete(got, w)
		}
		for extra := range got {
			errs = append(errs, fmt.Sprintf("%s: %s", base, extra))
		}
	}
	for base, fulls := range buckets {
		if _, ok := expected[base]; ok {
			continue
		}
		errs = append(errs, fulls...)
	}
	current = len(files) - skipped
	return current, skipped, errs
}

func equalSets(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

func sortedKeys(m map[string]bool) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// fnSpan is one fn declaration's byte range and name.
type fnSpan struct {
	name       string
	start, end int
}

// fnSpans splits a file body at `fn` lines. A span runs to the
// next fn (or EOF); non-fn decls between fns carry no tests or
// calls, so attribution is harmless.
func fnSpans(body string) []fnSpan {
	locs := fnStartRe.FindAllStringSubmatchIndex(body, -1)
	var out []fnSpan
	for i, l := range locs {
		end := len(body)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		out = append(out, fnSpan{name: body[l[2]:l[3]], start: l[0], end: end})
	}
	return out
}

// enclosingFn names the fn whose span contains off, or "" when
// the offset sits outside every fn (malformed input keeps the
// old file-wide comparison via the caller's fallback).
func enclosingFn(spans []fnSpan, off int) string {
	for _, s := range spans {
		if off >= s.start && off < s.end {
			return s.name
		}
	}
	return ""
}

// spanTests collects decision-table rows within one fn span,
// using the same tests/case scan as the file-wide pass.
func spanTests(span string) map[string]bool {
	out := map[string]bool{}
	inTests := false
	for _, line := range strings.Split(span, "\n") {
		switch {
		case testsRe.MatchString(line):
			inTests = true
		case inTests && eqRe.MatchString(line):
			inTests = false
		default:
			if inTests {
				if m := caseRe.FindStringSubmatch(line); m != nil {
					out[m[1]] = true
				}
			}
		}
	}
	return out
}

// reachingRows mirrors check.go reachingTests textually: the
// fn's own rows plus same-file transitive callers' rows.
// Foreign and intrinsic callees never contribute callers:
// their bodies don't execute under scripts.
func reachingRows(body string, spans []fnSpan, fn string) map[string]bool {
	byName := map[string]fnSpan{}
	for _, s := range spans {
		byName[s.name] = s
	}
	callers := map[string][]string{}
	for _, s := range spans {
		seen := map[string]bool{}
		for _, m := range callRe.FindAllStringSubmatch(body[s.start:s.end], -1) {
			if callee := m[1]; !seen[callee] {
				seen[callee] = true
				if _, ok := byName[callee]; ok {
					callers[callee] = append(callers[callee], s.name)
				}
			}
		}
	}
	out := map[string]bool{}
	if s, ok := byName[fn]; ok {
		for r := range spanTests(body[s.start:s.end]) {
			out[r] = true
		}
	}
	seen := map[string]bool{fn: true}
	queue := []string{fn}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, caller := range callers[u] {
			if seen[caller] {
				continue
			}
			seen[caller] = true
			queue = append(queue, caller)
			for r := range spanTests(body[byName[caller].start:byName[caller].end]) {
				out[r] = true
			}
		}
	}
	return out
}

func subsetOf(a, b map[string]bool) bool {
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

func main() {
	root, err := scan.RepoRoot()
	if err != nil {
		fmt.Println("MODULE CHECK FAILED")
		fmt.Println(" -", err)
		os.Exit(1)
	}
	roots := []string{filepath.Join(root, "sketches"), filepath.Join(root, "std")}
	for _, name := range examples {
		if st, err := os.Stat(filepath.Join(root, "sketches", name)); err == nil && st.IsDir() {
			roots = append(roots, filepath.Join(root, "sketches", name))
		}
	}
	current, skipped, scanned, errs := check(roots)
	if len(errs) > 0 {
		fmt.Println("MODULE CHECK FAILED")
		sort.Strings(errs)
		for _, e := range errs {
			fmt.Println(" -", e)
		}
		os.Exit(1)
	}
	fmt.Printf("modules OK: %d current in %v, %d legacy skipped\n", current, scanned, skipped)
}
