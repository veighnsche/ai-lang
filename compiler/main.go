// Package main is ailc v0.0.0, the ai-lang transpiler.
//
// Golden outputs in sketches/auth-login are gated by TestGoldenAuthLogin:
// any parse, test-evaluation, or emit change that alters output fails
// `go test ./...`.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func failf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

func main() {
	os.Exit(run(os.Args[1:]))
}

// version is stamped at build time via:
//
//	go build -ldflags "-X main.version=<v>" ./compiler
//
// Unstamped builds (e.g. plain `go install ...@latest`) report "dev".
var version = "dev"

func run(argv []string) int {
	if len(argv) > 0 && (argv[0] == "--version" || argv[0] == "-version" || argv[0] == "version") {
		fmt.Printf("ailc %s\n", version)
		return 0
	}
	if len(argv) > 0 && argv[0] == "lsp" {
		return runLSP()
	}
	if len(argv) > 0 && argv[0] == "normalize" {
		if err := runNormalize(os.Stdout, argv[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "ailc FAILED: %v\n", err)
			return 1
		}
		return 0
	}
	var out, format string
	var args []string
	for i := 0; i < len(argv); {
		if argv[i] == "--out" && i+1 < len(argv) {
			out = argv[i+1]
			i += 2
		} else if argv[i] == "--format" && i+1 < len(argv) {
			format = argv[i+1]
			i += 2
		} else {
			args = append(args, argv[i])
			i++
		}
	}
	if out == "" || len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: ailc [--format json] --out OUT_DIR file.ail [...]")
		return 2
	}
	if format != "" && format != "json" {
		fmt.Fprintf(os.Stderr, "ailc: unknown --format %q (want json)\n", format)
		return 2
	}
	if err := compileEx(out, args, format == "json"); err != nil {
		fmt.Fprintf(os.Stderr, "ailc FAILED: %v\n", err)
		return 1
	}
	return 0
}

func compile(out string, paths []string) error {
	return compileEx(out, paths, false)
}

// compileEx runs the same full suite the editor runs (checkStatic,
// world, checkSem, proof, tests), so "no squiggles" and "compiles" mean
// the same thing. Prose mode preserves the old fail-fast behavior by
// returning the first error; json mode prints every diagnostic as one
// JSON object per line on stdout and stays silent on success.
func compileEx(out string, paths []string, jsonOut bool) error {
	mods, texts, collected, err := parsePaths(paths)
	if err != nil {
		return err
	}
	pass := func(mod, fn, test string) {}
	if !jsonOut {
		pass = func(mod, fn, test string) {
			fmt.Printf("PASS %s.%s/%s\n", mod, fn, test)
		}
	}
	prog, collected := checkProgram(mods, texts, collected, pass)
	if err := firstError(collected); err != nil {
		return failDiags(collected, jsonOut)
	}
	for _, m := range mods {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			seen := map[string]bool{}
			for _, c := range walkCalls(fn.Body) {
				if prog.Uses[c.Fname] && !seen[c.Fname] {
					seen[c.Fname] = true
					fn.UsesHere = append(fn.UsesHere, c.Fname)
				}
			}
		}
	}
	total := 0
	for _, m := range mods {
		for _, d := range m.Decls {
			if fn, ok := d.(*FnDecl); ok {
				total += len(fn.Tests)
			}
		}
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	stemOf := map[string]string{}
	resultOfStem := map[string]string{}
	fnUnions := map[string]string{}
	for _, m := range mods {
		resultOfStem[m.Stem] = capitalize(m.Mod) + "Result"
		for _, d := range m.Decls {
			switch d := d.(type) {
			case *FnDecl:
				stemOf[d.Name] = m.Stem
				// Per-function unions (see fnResultUnion): call
				// temporaries and return annotations carry the
				// callee's own outcomes, never the module-wide
				// union, so strict checkers narrow exactly.
				u, err := fnResultUnion(d, prog)
				if err != nil {
					return err
				}
				fnUnions[d.Name] = u
			case *ExternDecl:
				u, err := externUnion(d, prog)
				if err != nil {
					return err
				}
				fnUnions[d.Name] = u
			case *TypeDecl:
				stemOf[d.Name] = m.Stem
			}
		}
	}
	for _, m := range mods {
		text, err := emitModule(m, prog, stemOf, resultOfStem, fnUnions)
		if err != nil {
			return err
		}
		name := m.Stem + ".ts"
		if err := os.WriteFile(out+"/"+name, []byte(text), 0o644); err != nil {
			return err
		}
		if !jsonOut {
			fmt.Printf("EMIT %s\n", name)
		}
	}
	catalog, err := json.MarshalIndent(buildCatalog(mods, prog, texts), "", "  ")
	if err != nil {
		return err
	}
	catalog = append(catalog, '\n')
	if err := os.WriteFile(out+"/errors.json", catalog, 0o644); err != nil {
		return err
	}
	if !jsonOut {
		fmt.Printf("ailc: %d tests passed, %d modules emitted to %s\n", total, len(mods), out)
	}
	return nil
}

// runNormalize implements `ailc normalize file.ail [...]`: gate on the
// full suite, then print every decision-table outcome in canonical form,
// one `mod.fn/test => value` line per test, sorted. Expectation mismatches
// still print (the outcome is the artifact); only execution failures —
// missing scripts, unbound names, leftover stubs — fail the command.
// Exit 0 on print: normalize observes, compile gates.
func runNormalize(w io.Writer, paths []string) error {
	if len(paths) == 0 {
		return fmt.Errorf("usage: ailc normalize file.ail [...]")
	}
	mods, _, collected, err := parsePaths(paths)
	if err != nil {
		return err
	}
	prog, collected := checkProgram(mods, nil, collected, nil)
	if ferr := firstError(collected); ferr != nil {
		return ferr
	}
	var lines []string
	for _, m := range mods {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			for _, t := range fn.Tests {
				got, _, err := runTestValue(fn, t, prog, nil)
				if err != nil {
					return err
				}
				lines = append(lines, fmt.Sprintf("%s.%s/%s => %s", m.Mod, fn.Name, t.Name, normalizeValue(got)))
			}
		}
	}
	sort.Strings(lines)
	for _, l := range lines {
		fmt.Fprintln(w, l)
	}
	return nil
}

// parsePaths reads and parses every path, collecting AIL1000 diagnostics
// for files that do not parse instead of failing fast, so one broken file
// never hides the rest. Raw IO errors still fail immediately.
//
// Identity is the cleaned input path: two inputs with different
// identities are different modules even when their basenames match.
// Output stems stay bare while unique, then disambiguate by directory;
// the same identity twice is an AIL5007 collision, rejected before
// evaluation or writing.
func parsePaths(paths []string) (mods []*Module, texts map[string]string, collected []Diag, err error) {
	texts = map[string]string{}
	seen := map[string]bool{}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, nil, nil, err
		}
		m, perr := parseModuleText(p, string(data))
		if perr != nil {
			base := filepath.Base(filepath.Clean(p))
			collected = append(collected, Diag{File: base, Line: diagLine(perr, 1), Sev: "error", Msg: stripLinePrefix(perr), Code: CodeParse})
			continue
		}
		if seen[m.ID] {
			collected = append(collected, Diag{File: m.File, Line: 1, Sev: "error",
				Msg: fmt.Sprintf("module %s inputs twice: one canonical identity per file", m.ID), Code: CodeModuleCollision})
			continue
		}
		seen[m.ID] = true
		mods = append(mods, m)
		texts[m.ID] = string(data)
	}
	assignStems(mods)
	return mods, texts, collected, nil
}

// assignStems gives every module an injective output stem. The bare
// stem (filename without extension) wins while unique, so single-file
// and distinct-name inputs emit exactly as before; sharers after the
// first (in sorted identity order) take the sanitized identity path,
// with numeric suffixes breaking residual ties. Deterministic in the
// input set, never silently merging two owners into one artifact.
func assignStems(mods []*Module) {
	count := map[string]int{}
	for _, m := range mods {
		count[m.Stem]++
	}
	used := map[string]bool{}
	ordered := append([]*Module{}, mods...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ID < ordered[j].ID })
	for _, m := range ordered {
		if count[m.Stem] == 1 && !used[m.Stem] {
			used[m.Stem] = true
			continue
		}
		candidate := sanitizeStem(strings.TrimSuffix(m.ID, ".ail"))
		for n := 2; used[candidate]; n++ {
			candidate = fmt.Sprintf("%s_%d", sanitizeStem(strings.TrimSuffix(m.ID, ".ail")), n)
		}
		m.Stem = candidate
		used[candidate] = true
	}
}

// sanitizeStem maps an identity path to stem characters: every run of
// non-letters-and-digits (separators included) becomes one underscore,
// with leading/trailing underscores trimmed. Underscore itself maps to
// itself, so "a/b" and "a_b" can still collide — assignStems breaks
// that tie with a numeric suffix.
func sanitizeStem(id string) string {
	var b strings.Builder
	prev := '_'
	for _, r := range id {
		var c rune
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			c = r
		default:
			c = '_'
		}
		if c == '_' && prev == '_' {
			continue
		}
		b.WriteRune(c)
		prev = c
	}
	return strings.Trim(b.String(), "_")
}

// checkProgram runs the full shared suite over parsed modules: static
// checks, world build, semantic checks, exhaustiveness proof, and test
// runs. pass fires per passing test. Execution-dependent phases are
// skipped on a dirty world, mirroring the editor.
func checkProgram(mods []*Module, texts map[string]string, collected []Diag, pass func(mod, fn, test string)) (*Program, []Diag) {
	for _, m := range mods {
		collected = append(collected, checkStatic(m, texts[m.ID])...)
	}
	if len(mods) == 0 {
		return nil, collected
	}
	prog, world := buildWorld(mods[0], mods, texts)
	collected = append(collected, world...)
	if hasErrors(world) {
		return prog, collected
	}

	// World-level termination refusal (v11): cross-file cycles are
	// reported per-line and suppress only execution-dependent checks,
	// per the R10 world-error rule.
	global := checkGlobalCycles(mods, texts, prog)
	collected = append(collected, global...)
	// Finite products only (first cut): record-type cycles fail before
	// tests or output through the same execution gate.
	recCycles := checkRecordCycles(mods, texts)
	collected = append(collected, recCycles...)
	gblocked := hasErrors(global) || hasErrors(recCycles)
	for _, m := range mods {
		text := texts[m.ID]
		var hook func(fn, test string)
		if pass != nil {
			mod := m.Mod
			hook = func(fn, test string) { pass(mod, fn, test) }
		}
		collected = append(collected, checkSem(m, text, prog, hook, gblocked)...)
	}
	for _, m := range mods {
		for _, err := range verifyExhaustiveAll([]*Module{m}, prog) {
			d := proofDiag(texts[m.ID], err)
			d.File = m.File
			collected = append(collected, d)
		}
	}
	return prog, collected
}

// hasErrors reports whether any diagnostic is an error.
func hasErrors(diags []Diag) bool {
	return firstError(diags) != nil
}

// firstError returns the first error diagnostic as a file:line: message,
// dropping a redundant filename the message already carries.
func firstError(diags []Diag) error {
	for _, d := range diags {
		if d.Sev == "error" {
			msg := d.Msg
			if d.File != "" {
				msg = strings.TrimPrefix(msg, d.File+": ")
				return fmt.Errorf("%s:%d: %s", d.File, d.Line, msg)
			}
			return fmt.Errorf("%s", msg)
		}
	}
	return nil
}

// failDiags prints every diagnostic as JSON lines in json mode, then
// returns the first error. Prose mode just returns it; run() prints the
// familiar single line.
func failDiags(collected []Diag, jsonOut bool) error {
	if jsonOut {
		reportDiags(os.Stdout, collected)
	}
	return firstError(collected)
}

// jsonDiag is the machine rendering of a Diag: same facts as the editor
// squiggle plus the stable code. Field order is fixed for golden tests.
type jsonDiag struct {
	Code  string `json:"code"`
	Sev   string `json:"sev"`
	File  string `json:"file"`
	Line  int    `json:"line"`
	Start int    `json:"start"`
	End   int    `json:"end"`
	Msg   string `json:"msg"`
}

// reportDiags writes diagnostics sorted file-top to file-bottom, one JSON
// object per line. stdout stays pure JSON: callers must suppress every
// other print in json mode.
func reportDiags(w io.Writer, diags []Diag) {
	cp := append([]Diag(nil), diags...)
	sortDiags(cp)
	for _, d := range cp {
		body, err := json.Marshal(jsonDiag{
			Code: d.Code, Sev: d.Sev, File: d.File, Line: d.Line,
			Start: d.Start, End: d.End, Msg: d.Msg,
		})
		if err != nil {
			continue
		}
		fmt.Fprintln(w, string(body))
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	c := s[0]
	if c >= 'a' && c <= 'z' {
		c -= 'a' - 'A'
	}
	return string(c) + s[1:]
}
