// ailc lsp: minimal Language Server over stdio.
//
// Speaks just enough JSON-RPC to drive editor squiggles: initialize,
// textDocument/didOpen, textDocument/didChange, shutdown/exit. Every
// keystroke re-runs the full diagnosis (parse, naming, uses resolution,
// exhaustiveness proof, signature tests) and publishes diagnostics.
//
// Run: ailc lsp [--baseline BASE.json]   (editors connect stdout/stdin
// with Content-Length framing). With --baseline, every diagnosis also
// runs revision-identity enforcement (a79): the same AIL6013 findings
// the CLI reports, as editor squiggles.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Diag is one squiggle: 1-based Line, Sev "error"/"warning", human
// message, an optional 0-based UTF-16 [Start, End) token span, a stable
// Code from compiler/code.go, and the File (basename) it belongs to. When
// End <= Start the editor falls back to the whole line, so checks that
// cannot name a token stay line-precise without pretending otherwise.
type Diag struct {
	File       string
	Line       int
	Sev        string
	Msg        string
	Start, End int
	Code       string
	// Expected, Found, and Hint carry machine-actionable
	// payloads (a71, a61 item 1): what the rule wanted, what
	// it saw, and the suggested fix shape. Empty means the
	// code carries no payload; JSON omits empty fields so
	// payload-free lines stay byte-identical.
	Expected string
	Found    string
	Hint     string
}

func sevCode(sev string) int {
	if sev == "warning" {
		return 2
	}
	return 1
}

// diagLine finds the message's own line info, defaulting to the fallback.
func diagLine(err error, fallback int) int {
	var le *LineError
	if errors.As(err, &le) && le.Line > 0 {
		return le.Line
	}
	return fallback
}

// locateLine returns the 1-based line first containing sub, or fallback.
func locateLine(text, sub string, fallback int) int {
	for n, line := range strings.Split(text, "\n") {
		if strings.Contains(line, sub) {
			return n + 1
		}
	}
	return fallback
}

// u16len counts UTF-16 code units: editor columns are UTF-16, Go strings
// are bytes, and the two only agree on ASCII.
func u16len(s string) int {
	n := 0
	for _, r := range s {
		if r > 0xFFFF {
			n += 2
		} else {
			n++
		}
	}
	return n
}

// u16col converts a byte offset inside line to a UTF-16 column.
func u16col(line string, byteOff int) int {
	if byteOff < 0 {
		byteOff = 0
	}
	if byteOff > len(line) {
		byteOff = len(line)
	}
	return u16len(line[:byteOff])
}

// tokenSpan returns the 0-based UTF-16 [start, end) span of token on the
// given 1-based line. It prefers the applied form (name + "(" hits the
// call, not the mention) then the declared form (name + ":" hits the
// parameter, not a use), then the bare token. ok=false means the token
// is not on that line and the caller must fall back to whole-line.
func tokenSpan(text string, line int, token string) (start, end int, ok bool) {
	lines := strings.Split(text, "\n")
	if line < 1 || line > len(lines) || token == "" {
		return 0, 0, false
	}
	ln := lines[line-1]
	for _, cand := range []string{token + "(", token + ":", token} {
		if i := strings.Index(ln, cand); i >= 0 {
			if cand != token {
				return u16col(ln, i), u16col(ln, i+len(token)), true
			}
			return u16col(ln, i), u16col(ln, i+len(cand)), true
		}
	}
	return 0, 0, false
}

// spanDiag builds a Diag whose squiggle covers token on line, falling
// back to the whole line when the token is not there.
func spanDiag(text string, line int, sev, msg, token, code string) Diag {
	if s, e, ok := tokenSpan(text, line, token); ok {
		return Diag{Line: line, Sev: sev, Msg: msg, Start: s, End: e, Code: code}
	}
	return Diag{Line: line, Sev: sev, Msg: msg, Code: code}
}

// locateLineFrom is locateLine starting at 1-based fromLine: given-table
// entries repeat across retry-nested matches, so each match searches
// downward from itself instead of stealing the first match's row.
func locateLineFrom(text, sub string, fromLine, fallback int) int {
	lines := strings.Split(text, "\n")
	if fromLine < 1 {
		fromLine = 1
	}
	for n := fromLine - 1; n < len(lines); n++ {
		if strings.Contains(lines[n], sub) {
			return n + 1
		}
	}
	return locateLine(text, sub, fallback)
}

// diagnose runs every check on the open file (sibling .ail files in dir
// provide the uses/provides world) and returns sorted diagnostics.
// Without a baseline no identity findings report: the editor stays
// quiet exactly as before.
func diagnose(dir, name, text string) []Diag {
	return diagnoseWith(dir, name, text, nil)
}

// diagnoseWith threads an accepted revision baseline through the
// same pipeline: when the world otherwise checks clean, identity
// drift against the baseline appends AIL6013 findings, mirroring
// the CLI's firstError gate so broken programs never gain drift
// noise on top of their real errors.
func diagnoseWith(dir, name, text string, base *RevisionBaseline) []Diag {
	var out []Diag
	var open *Module
	entries, _ := os.ReadDir(dir)
	files := []string{}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".ail") && e.Name() != name {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)
	// all is ordered deterministically: the open file first (its
	// declarations win shadowing), then siblings by filename, so repeated
	// diagnoses of the same directory agree with each other.
	all := []*Module{}
	open, err := parseModuleText(name, text)
	if err != nil {
		return append(out, Diag{File: name, Line: diagLine(err, 1), Sev: "error", Msg: stripLinePrefix(err), Code: CodeParse})
	}
	all = append(all, open)
	// An unparseable sibling only matters when it could be the missing
	// provider: the warning fires solely alongside "resolves nowhere",
	// so one broken demo never yellows its whole gallery.
	var sibWarn []Diag
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			continue
		}
		m, err := parseModuleText(f, string(data))
		if err != nil {
			sibWarn = append(sibWarn, Diag{File: name, Line: 1, Sev: "warning",
				Msg:  fmt.Sprintf("sibling %s does not parse, uses-checks may over-report: %v", f, err),
				Code: CodeSiblingParse})
			continue
		}
		all = append(all, m)
	}
	texts := map[string]string{name: text}
	for _, f := range files {
		if data, err := os.ReadFile(filepath.Join(dir, f)); err == nil {
			texts[f] = string(data)
		}
	}
	out = append(out, checkStatic(open, text)...)
	prog, world := buildWorld(open, all, texts)
	out = append(out, world...)
	for _, d := range world {
		if d.Sev == "error" && strings.Contains(d.Msg, "resolves nowhere") {
			out = append(out, sibWarn...)
			break
		}
	}
	if len(world) > 0 {
		sortDiags(out)
		return withFile(out, name)
	}
	// a46 S2 barrier: same whole-program certification as the CLI, so
	// "no squiggles" and "compiles" cannot diverge on authority.
	out = append(out, certifyExports(all, prog, texts)...)
	// S2 slice plan barrier: same whole-program rule for asset bridges.
	out = append(out, certifyAssetBridge(all, prog, texts)...)

	for _, err := range verifyExhaustiveAll([]*Module{open}, prog) {
		out = append(out, proofDiag(text, err))
	}
	// World-level termination refusal (a11): cross-file cycles block
	// execution through the same gate, per the R10 world-error rule
	// (report per-line, suppress only execution-dependent checks).
	global := checkGlobalCycles(all, texts, prog)
	out = append(out, global...)
	recCycles := checkRecordCycles(all, texts)
	out = append(out, recCycles...)
	out = append(out, checkSem(open, text, prog, nil, hasErrors(global) || hasErrors(recCycles))...)
	if base != nil && !hasErrors(out) {
		out = append(out, CheckRevisionIdentity(prog, texts, base)...)
	}
	// a87: pinned weakening reports under the same clean-world rule,
	// after identity, so drift noise never stacks atop real errors.
	if base != nil && !hasErrors(out) {
		out = append(out, CheckPinnedRows(prog, texts, base)...)
	}
	// a82: verifier activation in the editor. Like identity
	// enforcement, proof findings append only when the world
	// otherwise checks clean, so broken programs never gain
	// proof noise atop real errors.
	if !hasErrors(out) {
		out = append(out, VerifyContracts(prog, texts)...)
	}
	sortDiags(out)
	return withFile(out, name)
}

// withFile stamps every unattributed diagnostic with its owner file.
// buildWorld sets File at creation (it knows every module); everything
// else describes the open document.
func withFile(out []Diag, name string) []Diag {
	for i := range out {
		if out[i].File == "" {
			out[i].File = name
		}
	}
	return out
}

// checkStatic runs the checks that need no world: naming, header
// integrity, and decision-table shapes. Shared by the editor and the CLI
// so "no squiggles" and "compiles" cannot diverge.
func checkStatic(open *Module, text string) []Diag {
	var out []Diag
	out = append(out, checkNaming(open, text)...)
	out = append(out, checkModIntegrity(open, text)...)
	for _, d := range open.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		if len(fn.Tests) == 0 {
			out = append(out, spanDiag(text, fn.Line, "error",
				fmt.Sprintf("%s ships no tests: every function needs its decision table", fn.Name), fn.Name, CodeMissingTests))
		}
		out = append(out, checkTestShapes(fn, text)...)
	}
	return withFile(out, open.File)
}

// checkSem runs the world-dependent checks: calls, given, emits, and
// unused items. onPass fires per passing test (the CLI prints PASS; the
// editor passes nil). Only call on a clean world. extBlocked carries a
// world-level refusal (a11: cross-file cycles) into the same
// prove-first gate as the per-module termination proofs.
func checkSem(open *Module, text string, prog *Program, onPass func(fn, test string), extBlocked bool) []Diag {
	var out []Diag
	// Slice 1: resolve const-named patterns to literals before
	// any other per-function check, test run, or proof sees
	// them. Idempotent: rewritten literals are not revisited.
	out = append(out, elaborateConstPatterns(open, prog, text)...)
	called := map[string]bool{}
	localExtern := map[string]bool{}
	for _, d := range open.Decls {
		if ex, ok := d.(*ExternDecl); ok {
			localExtern[ex.Name] = true
		}
	}
	for _, d := range open.Decls {
		switch d := d.(type) {
		case *FnDecl:
			fn := d
			for k := range calledFns(fn) {
				called[k] = true
			}
			out = append(out, checkCalls(fn, prog, localExtern, text)...)
			out = append(out, checkEagerScrutinee(fn, text)...)
			out = append(out, checkDecreases(fn, prog, text)...)
			out = append(out, checkEffects(fn, prog, text)...)
			out = append(out, checkGiven(fn, prog, text)...)
			out = append(out, checkScriptConsistency(fn, prog, text)...)
			out = append(out, checkEmits(fn, prog, text)...)
			out = append(out, checkTypes(fn, prog, text)...)
			out = append(out, checkUnusedParams(fn, text)...)
			out = append(out, checkConstRefs(fn, prog, open, text)...)
			for k := range usedConsts(fn) {
				called[k] = true
			}
		case *ExternDecl:
			out = append(out, checkExternSig(d, prog, text)...)
		case *BrandDecl:
			out = append(out, checkBrandDecl(d, prog, text)...)
		case *StateDecl:
			out = append(out, checkStateDecl(d, text)...)
		case *TypeDecl:
			out = append(out, checkDeclFields(d.Name, d.Fields, d.Line, prog, text)...)
		case *VariantDecl:
			for _, c := range d.Cases {
				out = append(out, checkDeclFields(qualifyCase(d.Name, c.Short), c.Fields, c.Line, prog, text)...)
			}
		case *ErrorDecl:
			out = append(out, checkDeclFields(d.Name, d.Fields, d.Line, prog, text)...)
		}
	}
	out = append(out, checkConstDecls(open, prog, text)...)
	for k := range prog.ConstUsed[open.ID] {
		called[k] = true
	}
	out = append(out, checkUnusedUses(open, text, called)...)
	out = append(out, checkLocalCycles(open, prog, text)...)
	// A local call cycle would hang test execution: prove acyclic
	// before running anything. Open termination proofs (bad, stale,
	// or unproven decreases, unguarded recursion, cross-file cycles),
	// ill-formed script evidence (outcome-only rows, malformed
	// outcomes, Ok claims the provider body contradicts), and unproven
	// authority (effects) block the same way:
	// the gate is prove-first, run-after.
	blocked := extBlocked
	for _, d := range out {
		if d.Sev == "error" && (d.Code == CodeLocalCycle ||
			d.Code == CodeBadDecreases || d.Code == CodeStaleDecreases ||
			d.Code == CodeNoDecrease || d.Code == CodeNoGuard ||
			d.Code == CodeNoExchange || d.Code == CodeBadStub || d.Code == CodeInconsistentScript ||
			d.Code == CodeUndeclaredEffect || d.Code == CodeStaleEffect) {
			blocked = true
		}
	}
	// Coverage is assessed over green tests module-wide: each test
	// runs with a private map, and only passing runs merge into the
	// shared map, so a helper arm counts caller flow-through without
	// letting a failing test fake coverage. Functions without tests
	// are already flagged elsewhere.
	shared := map[*Node]map[int]bool{}
	merge := func(src map[*Node]map[int]bool) {
		for n, arms := range src {
			dst := shared[n]
			if dst == nil {
				dst = map[int]bool{}
				shared[n] = dst
			}
			for i := range arms {
				dst[i] = true
			}
		}
	}
	failed := map[string]bool{}
	if !blocked {
		for _, d := range open.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			for _, t := range fn.Tests {
				tc := map[*Node]map[int]bool{}
				if err := runTest(fn, t, prog, tc); err != nil {
					var uce *UnknownCallError
					if errors.As(err, &uce) && calleeUnknown(prog, uce.Fname) {
						// a62: the row can only fail on the
						// unknown call checkCalls already
						// reported; suppress the AIL4200 but
						// mark the fn failed so coverage
						// stays silent too.
						failed[fn.Name] = true
						continue
					}
					failed[fn.Name] = true
					out = append(out, spanDiag(text, t.Line, "error",
						fmt.Sprintf("test %s fails: %s", t.Name, stripLinePrefix(err)), t.Name, CodeTestFailed))
				} else {
					merge(tc)
					if onPass != nil {
						onPass(fn.Name, t.Name)
					}
				}
			}
		}
	}
	for _, d := range open.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		if !failed[fn.Name] && !blocked && len(fn.Tests) > 0 {
			out = append(out, checkCoverage(fn, prog, text, shared)...)
		}
	}
	return withFile(out, open.File)
}

// checkCoverage enforces the test-per-arm law: every match arm must
// execute at least once across the function's decision-table run.
// Untaken arms are dead code or missing tests, both compile errors.
// The one exception is a checked identity relay: a bound error arm of
// a local call whose body is exactly the same-kind reconstruction
// with every payload field unchanged. Such an arm is a total,
// transparent re-raise — no behavior remains to witness — so it
// carries a structural certificate instead of an execution one. A
// relay-shaped arm that fails the check (wrong kind, dropped field,
// changed value) is an invalid certificate, not an uncovered arm.
func checkCoverage(fn *FnDecl, prog *Program, text string, cov map[*Node]map[int]bool) []Diag {
	var out []Diag
	for _, n := range matchNodes(fn.Body) {
		for i, a := range n.Arms {
			if cov[n][i] {
				continue
			}
			// First slot describes the arm.
			desc, tok := patDesc(a.Pats[0])
			if ok, reason := relayStatus(prog, fn.Name, n, n.Arms, i); ok {
				if reason == "" {
					continue
				}
				out = append(out, spanDiag(text, a.Line, "error",
					fmt.Sprintf("invalid identity relay %s in %s: %s", desc, fn.Name, reason), tok, CodeInvalidRelay))
				continue
			}
			out = append(out, spanDiag(text, a.Line, "error",
				fmt.Sprintf("no test takes %s in %s", desc, fn.Name), tok, CodeArmUntaken))
		}
	}
	return out
}

// relayStatus checks an untaken arm for the identity-relay shape: a
// bound error arm of a local call whose body rebuilds an error.
// It returns ok=false for anything else (the arm stays under the
// execution law). For a relay shape it returns ok=true with reason=""
// when the certificate verifies (same kind, complete fields, unchanged
// bound values, no other content), or ok=true with a reason naming the
// defect when the certificate is invalid. A shadowed arm is never
// certified: an earlier arm matching the same error kind already
// consumes every value this one could take, so structural evidence
// cannot substitute for the execution the shadowing removed.
func relayStatus(prog *Program, owner string, n *Node, arms []Arm, idx int) (bool, string) {
	if n.Kind != MatchCall {
		return false, ""
	}
	if localCallee(prog, owner, n.Scruts[0].Fname) == nil {
		return false, ""
	}
	a := arms[idx]
	pat := a.Pats[0]
	if pat.Kind != "variant" || pat.Name == "Ok" {
		return false, ""
	}
	for _, prev := range arms[:idx] {
		pp := prev.Pats[0]
		if (pp.Kind == "variant" || pp.Kind == "variantWild") && pp.Name == pat.Name {
			return false, ""
		}
	}
	fields, known := prog.Errors[pat.Name]
	if !known {
		return false, ""
	}
	rhs := a.Rhs
	if rhs == nil || rhs.IsMatch || rhs.Small == nil || rhs.Small.Kind != "ctor" {
		return false, ""
	}
	s := rhs.Small
	if !strings.Contains(s.Ctor, ".") {
		return false, ""
	}
	if s.Ctor != pat.Name {
		return true, fmt.Sprintf("reconstructs %s instead of %s", s.Ctor, pat.Name)
	}
	want := map[string]bool{}
	for _, f := range fields {
		want[f] = true
	}
	seen := map[string]bool{}
	for _, arg := range s.Args {
		if !want[arg.Name] {
			return true, fmt.Sprintf("rebuilds unexpected field %s", arg.Name)
		}
		seen[arg.Name] = true
		v := arg.V
		if v == nil || v.Kind != "ref" || len(v.Ref) != 2 || v.Ref[0] != pat.Var || v.Ref[1] != arg.Name {
			return true, fmt.Sprintf("field %s is not %s.%s", arg.Name, pat.Var, arg.Name)
		}
	}
	for _, f := range fields {
		if !seen[f] {
			return true, fmt.Sprintf("drops field %s", f)
		}
	}
	return true, ""
}

// patDesc renders a pattern the way its arm row reads it, plus the token
// the squiggle should cover on that row.
func patDesc(p Pattern) (desc, tok string) {
	switch p.Kind {
	case "variantWild":
		return "on " + p.Name + " _", p.Name
	case "variant":
		return "on " + p.Name + " " + p.Var, p.Name
	case "bool":
		if p.B {
			return "on true", "true"
		}
		return "on false", "false"
	case "str":
		// a66: an interpreted pattern's decoded value is not
		// searchable in source; locate its verbatim spelling.
		if strings.HasPrefix(p.Raw, `e"`) {
			return "on " + p.Raw, p.Raw
		}
		return "on " + strconv.Quote(p.Str), p.Str
	default:
		return "_", "_"
	}
}

// sortDiags orders squiggles file-top to file-bottom, errors before
// warnings on the same line, for a stable editor presentation.
func sortDiags(out []Diag) {
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Line != out[j].Line {
			return out[i].Line < out[j].Line
		}
		if out[i].Sev != out[j].Sev {
			return out[i].Sev < out[j].Sev
		}
		return out[i].Msg < out[j].Msg
	})
}

func stripLinePrefix(err error) string {
	var le *LineError
	if errors.As(err, &le) {
		return le.Err.Error()
	}
	return err.Error()
}

// proofDiag anchors an exhaustiveness-proof failure to a token. A stale
// arm names its kind on its own row, so the kind gets the squiggle; every
// other proof failure (missing arm, bool shape, bad pattern) lands on the
// match or arm keyword of its row. The shape matching mirrors the exact
// message shapes built in eval.go's verifyExhaustiveAll, which is also
// what assigns each failure its stable code.
func proofDiag(text string, err error) Diag {
	line := diagLine(err, 1)
	msg := stripLinePrefix(err)
	if i := strings.Index(msg, "stale match arm "); i >= 0 {
		// a63: the kind is the first field after the marker;
		// a nesting hint may follow it (see verifyExhaustiveAll).
		kind := strings.TrimSpace(msg[i+len("stale match arm "):])
		if j := strings.IndexAny(kind, " ;"); j >= 0 {
			kind = kind[:j]
		}
		return spanDiag(text, line, "error", msg, kind, CodeStaleArm)
	}
	code := CodeProofOther
	var want, hint string
	switch {
	case strings.Contains(msg, "non-exhaustive match, missing "):
		code = CodeMissingArm
		if i := strings.Index(msg, "missing "); i >= 0 {
			want = strings.TrimSpace(msg[i+len("missing "):])
			hint = fmt.Sprintf("add an `on %s ...` arm covering the missing outcome", want)
		}
	case strings.Contains(msg, "bool match must be exactly "):
		code = CodeBoolArms
	case strings.Contains(msg, "value match without _ "):
		code = CodeValueNoWild
	case strings.Contains(msg, "call-match arm must be "):
		code = CodeBadArmKind
	case strings.Contains(msg, "patterns; this match has "):
		code = CodeBadArmKind
	case strings.Contains(msg, "variant pattern on a non-call "):
		code = CodeVariantOnVal
	case strings.Contains(msg, "mixes integer patterns with"):
		code = CodeBoolArms
	case strings.Contains(msg, "is fully covered by earlier arms"):
		code = CodeUselessArm
	case strings.Contains(msg, "contributes no remaining space"):
		code = CodeUselessAlt
	case strings.Contains(msg, "alternatives take scalar patterns only"):
		code = CodeUselessAlt
	}
	kw := "match"
	if lines := strings.Split(text, "\n"); line >= 1 && line <= len(lines) {
		if strings.HasPrefix(strings.TrimSpace(lines[line-1]), "on ") {
			kw = "on"
		}
	}
	d := spanDiag(text, line, "error", msg, kw, code)
	d.Expected = want
	d.Hint = hint
	return d
}

// ------------------------------------------------------------- protocol ---
type rpcMsg struct {
	ID     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params json.RawMessage  `json:"params"`
}

type docID struct {
	URI string `json:"uri"`
}

func pathFromURI(uri string) string {
	return strings.TrimPrefix(uri, "file://")
}

func writeFrame(w *bufio.Writer, v any) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Content-Length: %d\r\n\r\n", len(body))
	w.Write(body)
	return w.Flush()
}

func publishDiagnostics(w *bufio.Writer, uri, text string, diags []Diag) error {
	lines := strings.Split(text, "\n")
	items := []any{}
	for _, d := range diags {
		line := d.Line - 1
		if line < 0 {
			line = 0
		}
		if line >= len(lines) {
			line = len(lines) - 1
		}
		width := u16len(lines[line])
		start, end := 0, width
		if d.End > d.Start {
			start, end = d.Start, d.End
			if start < 0 {
				start = 0
			}
			if end > width {
				end = width
			}
			if end <= start {
				start, end = 0, width
			}
		}
		item := map[string]any{
			"range": map[string]any{
				"start": map[string]any{"line": line, "character": start},
				"end":   map[string]any{"line": line, "character": end},
			},
			"severity": sevCode(d.Sev),
			"source":   "ailc",
			"message":  d.Msg,
		}
		// Slice 0: the stable code and the a71 payloads ride
		// the wire (code string, data object); empty fields
		// stay omitted so payload-free lines are byte-identical.
		if d.Code != "" {
			item["code"] = d.Code
		}
		data := map[string]any{}
		if d.Expected != "" {
			data["expected"] = d.Expected
		}
		if d.Found != "" {
			data["found"] = d.Found
		}
		if d.Hint != "" {
			data["hint"] = d.Hint
		}
		if len(data) > 0 {
			item["data"] = data
		}
		items = append(items, item)
	}
	if items == nil {
		items = []any{}
	}
	return writeFrame(w, map[string]any{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params":  map[string]any{"uri": uri, "diagnostics": items},
	})
}

// parseLSPArgs takes the flags after `lsp`: only --baseline PATH.
// Bare means unenforced; anything else is a usage error.
func parseLSPArgs(argv []string) (string, error) {
	var baseline string
	for i := 0; i < len(argv); {
		if argv[i] == "--baseline" && i+1 < len(argv) {
			baseline = argv[i+1]
			i += 2
		} else {
			return "", fmt.Errorf("usage: ailc lsp [--baseline BASE.json]")
		}
	}
	return baseline, nil
}

func runLSP(argv []string) int {
	baselinePath, err := parseLSPArgs(argv)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	// The baseline loads once at startup: per-keystroke reloads
	// would re-read the file on every edit, and a mid-session
	// baseline change takes effect on editor restart. A missing
	// or unreadable baseline warns and runs unenforced: a bad
	// flag must not brick editing, and the CLI stays the
	// authority (it fails hard on the same input).
	var base *RevisionBaseline
	if baselinePath != "" {
		base, err = LoadBaseline(baselinePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ailc lsp: cannot load baseline, revision enforcement off: %v\n", err)
			base = nil
		}
	}
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	docs := map[string]string{}
	respond := func(id *json.RawMessage, result any) {
		writeFrame(out, map[string]any{"jsonrpc": "2.0", "id": id, "result": result})
	}
	for {
		var length int
		for {
			line, err := in.ReadString('\n')
			if err != nil {
				return 0
			}
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			if strings.HasPrefix(strings.ToLower(line), "content-length:") {
				length, _ = strconv.Atoi(strings.TrimSpace(line[len("content-length:"):]))
			}
		}
		body := make([]byte, length)
		if _, err := io.ReadFull(in, body); err != nil {
			return 0
		}
		var msg rpcMsg
		if err := json.Unmarshal(body, &msg); err != nil {
			continue
		}
		switch msg.Method {
		case "initialize":
			respond(msg.ID, map[string]any{"capabilities": map[string]any{"textDocumentSync": 1}})
		case "initialized", "$/cancelRequest", "textDocument/didClose":
			if msg.Method == "textDocument/didClose" {
				var p struct {
					TextDocument docID `json:"textDocument"`
				}
				if json.Unmarshal(msg.Params, &p) == nil {
					delete(docs, p.TextDocument.URI)
				}
			}
		case "textDocument/didOpen":
			var p struct {
				TextDocument struct {
					URI  string `json:"uri"`
					Text string `json:"text"`
				} `json:"textDocument"`
			}
			if json.Unmarshal(msg.Params, &p) != nil {
				continue
			}
			docs[p.TextDocument.URI] = p.TextDocument.Text
			path := pathFromURI(p.TextDocument.URI)
			diags := diagnoseWith(filepath.Dir(path), filepath.Base(path), p.TextDocument.Text, base)
			publishDiagnostics(out, p.TextDocument.URI, p.TextDocument.Text, diags)
		case "textDocument/didChange":
			var p struct {
				TextDocument docID `json:"textDocument"`
				Changes      []struct {
					Text string `json:"text"`
				} `json:"contentChanges"`
			}
			if json.Unmarshal(msg.Params, &p) != nil || len(p.Changes) == 0 {
				continue
			}
			text := p.Changes[len(p.Changes)-1].Text
			docs[p.TextDocument.URI] = text
			path := pathFromURI(p.TextDocument.URI)
			diags := diagnoseWith(filepath.Dir(path), filepath.Base(path), text, base)
			publishDiagnostics(out, p.TextDocument.URI, text, diags)
		case "shutdown":
			respond(msg.ID, nil)
		case "exit":
			return 0
		}
	}
}
