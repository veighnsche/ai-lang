package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Slice 1: named scalar-literal constants. Probes first: parse,
// naming, duplicates, literal-only initializers, uses/provides,
// lazy resolution in eval/emit, syntactic termination.

// hasError reports whether any diagnostic is an error: probes
// asserting clean builds ignore warnings (unused params still
// warn when a body ignores its row, const or not).
func hasError(diags []Diag) bool {
	for _, d := range diags {
		if d.Sev == "error" {
			return true
		}
	}
	return false
}

const constDeclSrc = `mod m
  provides [m__go, M__Out, m__COLON]
  uses []
  emits []

const m__COLON: int rev 1 = 58
`

// TestConstParses pins the declaration shape: parsePaths yields a
// ConstDecl carrying name, type, revision, and literal value.
func TestConstParses(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "m.ail")
	if err := os.WriteFile(fp, []byte(constDeclSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	mods, _, _, err := parsePaths([]string{fp})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var found *ConstDecl
	for _, d := range mods[0].Decls {
		if c, ok := d.(*ConstDecl); ok {
			found = c
		}
	}
	if found == nil {
		t.Fatalf("no ConstDecl parsed")
	}
	if found.Name != "m__COLON" || found.Type != "int" || found.Rev != 1 {
		t.Fatalf("bad decl: %+v", found)
	}
	if found.Value == nil || found.Value.Kind != "int" || found.Value.Num.String() != "58" {
		t.Fatalf("bad value: %+v", found.Value)
	}
}

// TestConstNaming pins AIL2003: const names are domain__SCREAMING.
func TestConstNaming(t *testing.T) {
	src := `mod m
  provides [m__go, Foo]
  uses []
  emits []

const Foo: int rev 1 = 1
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); !hasCode(diags, "AIL2003") {
		t.Fatalf("bad const name reported no AIL2003: %v", diags)
	}
}

// TestConstNamingMulti pins the decided stdlib convention: the
// domain is the first segment, so std__ascii__COLON is accepted
// (no AIL2003) while the SCREAMING tail stays mandatory.
func TestConstNamingMulti(t *testing.T) {
	src := `mod m
  provides [m__go, std__ascii__COLON]
  uses []
  emits []

const std__ascii__COLON: int rev 1 = 58
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); hasCode(diags, "AIL2003") {
		t.Fatalf("multi-segment const name reported AIL2003: %v", diags)
	}
}

// TestConstDuplicate pins AIL2205 for same-module duplicates.
func TestConstDuplicate(t *testing.T) {
	src := `mod m
  provides [m__go, m__A]
  uses []
  emits []

const m__A: int rev 1 = 1

const m__A: int rev 1 = 2
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); !hasCode(diags, "AIL2205") {
		t.Fatalf("duplicate const reported no AIL2205: %v", diags)
	}
}

// TestConstLiteralOnly pins AIL6016: initializers are literals,
// never computed, never aliases.
func TestConstLiteralOnly(t *testing.T) {
	for name, init := range map[string]string{"computed": "1 + 2", "alias": "m__A"} {
		src := `mod m
  provides [m__go, m__A, m__B]
  uses []
  emits []

const m__A: int rev 1 = 1

const m__B: int rev 1 = ` + init + `
`
		dir := writeLSPDir(t, map[string]string{"m.ail": src})
		if diags := diagnose(dir, "m.ail", src); !hasCode(diags, "AIL6016") {
			t.Fatalf("%s initializer reported no AIL6016: %v", name, diags)
		}
	}
}

// TestConstUseBeforeDecl pins order independence: a use above the
// declaration resolves by qualified name.
func TestConstUseBeforeDecl(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out, m__N]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 58)
=
  Ok(value = m__N)

const m__N: int rev 1 = 58
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); hasError(diags) {
		t.Fatalf("use-before-decl reported: %v", diags)
	}
}

// TestConstExecutes pins lazy resolution through execution: the
// test passes only if the const value flows into eval, and the
// emitted program inlines the literal.
func TestConstExecutes(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out, m__N]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

const m__N: int rev 1 = 58

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 58)
=
  Ok(value = m__N)
`
	dir := t.TempDir()
	fp := filepath.Join(dir, "m.ail")
	if err := os.WriteFile(fp, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out")
	if err := compile(out, []string{fp}); err != nil {
		t.Fatalf("const program failed to compile: %v", err)
	}
	ts, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(ts), "58n") {
		t.Fatalf("emitted program does not inline the literal: %s", ts)
	}
}

// TestConstExecutesValue pins the value, not just success: a wrong
// expectation fails, proving the const (not a default) flowed in.
func TestConstExecutesValue(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out, m__N]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

const m__N: int rev 1 = 58

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 59)
=
  Ok(value = m__N)
`
	dir := t.TempDir()
	fp := filepath.Join(dir, "m.ail")
	if err := os.WriteFile(fp, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := compile(filepath.Join(dir, "out"), []string{fp}); err == nil {
		t.Fatalf("wrong expectation compiled clean: const value did not flow into eval")
	}
}

// TestConstUnknown pins AIL2104 for references that resolve nowhere.
func TestConstUnknown(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 1)
=
  Ok(value = m__NOPE)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); !hasCode(diags, "AIL2104") {
		t.Fatalf("unknown const reported no AIL2104: %v", diags)
	}
}

// TestConstMissingUses pins AIL2105: a foreign const without a pin.
func TestConstMissingUses(t *testing.T) {
	lib := `mod lib
  provides [lib__K]
  uses []
  emits []

const lib__K: int rev 1 = 7
`
	app := `mod app
  provides [app__go, A__Out]
  uses []
  emits []

type A__Out rev 1 (
  value: int
)

fn app__go(x: int) -> A__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 7)
=
  Ok(value = lib__K)
`
	dir := writeLSPDir(t, map[string]string{"lib.ail": lib, "app.ail": app})
	if diags := diagnose(dir, "app.ail", app); !hasCode(diags, "AIL2105") {
		t.Fatalf("unpinnned foreign const reported no AIL2105: %v", diags)
	}
}

// TestConstUsesPin pins the open gate: a pinned foreign const
// resolves across files.
func TestConstUsesPin(t *testing.T) {
	lib := `mod lib
  provides [lib__K]
  uses []
  emits []

const lib__K: int rev 1 = 7
`
	app := `mod app
  provides [app__go, A__Out]
  uses [lib__K@1]
  emits []

type A__Out rev 1 (
  value: int
)

fn app__go(x: int) -> A__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 7)
=
  Ok(value = lib__K)
`
	dir := writeLSPDir(t, map[string]string{"lib.ail": lib, "app.ail": app})
	if diags := diagnose(dir, "app.ail", app); hasError(diags) {
		t.Fatalf("pinned foreign const reported: %v", diags)
	}
}

// TestConstProvidesMiss pins provides coverage for the new decl.
func TestConstProvidesMiss(t *testing.T) {
	src := `mod m
  provides [m__go]
  uses []
  emits []

const m__K: int rev 1 = 7
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); !hasCode(diags, CodeProvidesMiss) {
		t.Fatalf("unprovided const reported no AIL2301: %v", diags)
	}
}

// TestConstPatternBoolStr pins const patterns in value matches: a
// declared bool/str constant behaves exactly like its literal.
func TestConstPatternBoolStr(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out, m__FLAG, m__SEP]
  uses []
  emits []

type M__Out rev 1 (
  value: str
)

const m__FLAG: bool rev 1 = true

const m__SEP: str rev 1 = ", "

fn m__go(x: int, s: str) -> M__Out rev 1
  emits []
  tests
    one(x = 1, s = ", ") => Ok(value = "yes")
    two(x = 2, s = ", ") => Ok(value = "sep")
    three(x = 2, s = ";") => Ok(value = "no")
=
  match x <= 1
    m__FLAG => Ok(value = "yes")
    false => match s
      m__SEP => Ok(value = "sep")
      _ => Ok(value = "no")
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); hasError(diags) {
		t.Fatalf("const patterns reported: %v", diags)
	}
}

// TestConstPatternInt pins slice-3 admission: integer constants
// are patterns now (the range slice landed), behaving exactly
// like their literal.
func TestConstPatternInt(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out, m__N]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

const m__N: int rev 1 = 58

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    colon(x = 58) => Ok(value = 2)
    other(x = 1) => Ok(value = 3)
=
  match x
    m__N => Ok(value = 2)
    _ => Ok(value = 3)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); hasError(diags) {
		t.Fatalf("int const pattern reported: %v", diags)
	}
}

// TestConstEvidenceMarksUsed pins the unused-import interplay: a
// constant referenced only in evidence still marks its pin used.
func TestConstEvidenceMarksUsed(t *testing.T) {
	lib := `mod lib
  provides [lib__K]
  uses []
  emits []

const lib__K: int rev 1 = 7
`
	app := `mod app
  provides [app__go, A__Out]
  uses [lib__K@1]
  emits []

type A__Out rev 1 (
  value: int
)

fn app__go(x: int) -> A__Out rev 1
  emits []
  tests
    go(x = lib__K) => Ok(value = 7)
=
  Ok(value = x)
`
	dir := writeLSPDir(t, map[string]string{"lib.ail": lib, "app.ail": app})
	if diags := diagnose(dir, "app.ail", app); hasCode(diags, "AIL3401") {
		t.Fatalf("evidence-only const pin reported unused: %v", diags)
	}
}

// TestConstPatternForeignPin pins AIL2105 for patterns: a foreign
// constant in pattern position needs its rev pin like a body
// reference.
func TestConstPatternForeignPin(t *testing.T) {
	lib := `mod lib
  provides [lib__FLAG]
  uses []
  emits []

const lib__FLAG: bool rev 1 = true
`
	app := `mod app
  provides [app__go, A__Out]
  uses []
  emits []

type A__Out rev 1 (
  value: str
)

fn app__go(x: int) -> A__Out rev 1
  emits []
  tests
    one(x = 1) => Ok(value = "yes")
    two(x = 2) => Ok(value = "no")
=
  match x <= 1
    lib__FLAG => Ok(value = "yes")
    false => Ok(value = "no")
`
	dir := writeLSPDir(t, map[string]string{"lib.ail": lib, "app.ail": app})
	if diags := diagnose(dir, "app.ail", app); !hasCode(diags, "AIL2105") {
		t.Fatalf("unpinned foreign const pattern reported no AIL2105: %v", diags)
	}
}

// TestConstPatternMarksUsed pins the pattern evidence interplay: a
// pinned foreign const referenced only in a pattern still marks
// its pin used (no AIL3401).
func TestConstPatternMarksUsed(t *testing.T) {
	lib := `mod lib
  provides [lib__FLAG]
  uses []
  emits []

const lib__FLAG: bool rev 1 = true
`
	app := `mod app
  provides [app__go, A__Out]
  uses [lib__FLAG@1]
  emits []

type A__Out rev 1 (
  value: str
)

fn app__go(x: int) -> A__Out rev 1
  emits []
  tests
    one(x = 1) => Ok(value = "yes")
    two(x = 2) => Ok(value = "no")
=
  match x <= 1
    lib__FLAG => Ok(value = "yes")
    false => Ok(value = "no")
`
	dir := writeLSPDir(t, map[string]string{"lib.ail": lib, "app.ail": app})
	if diags := diagnose(dir, "app.ail", app); hasCode(diags, "AIL3401") {
		t.Fatalf("pattern-only const pin reported unused: %v", diags)
	}
}

// TestConstIdentityDrift pins canonConst (R4 clarification): same
// name and rev with a changed value is baseline drift ("value
// changed") — a changed value cannot hide behind its spelling.
func TestConstIdentityDrift(t *testing.T) {
	base := `mod m
  provides [m__N]
  uses []
  emits []

const m__N: int rev 1 = 58
`
	files := map[string]string{"m.ail": base}
	progB, _ := revisionProg(t, files, []string{"m.ail"})
	bsln := revisionBaseline(t, progB, "review-base:B")
	changed := `mod m
  provides [m__N]
  uses []
  emits []

const m__N: int rev 1 = 59
`
	progC, textsC := revisionProg(t, map[string]string{"m.ail": changed}, []string{"m.ail"})
	diags := CheckRevisionIdentity(progC, textsC, bsln)
	if !hasDiag(diags, "error", "differs from accepted baseline") {
		t.Fatalf("expected value drift, got %v", diags)
	}
	if !hasFound(diags, "value changed") {
		t.Fatalf("expected value-changed detail, got %v", diags)
	}
}

// TestConstIdentityClean pins the quiet path: an unchanged const
// compares clean against its own baseline.
func TestConstIdentityClean(t *testing.T) {
	src := `mod m
  provides [m__N]
  uses []
  emits []

const m__N: int rev 1 = 58
`
	files := map[string]string{"m.ail": src}
	prog, texts := revisionProg(t, files, []string{"m.ail"})
	bsln := revisionBaseline(t, prog, "review-base:B")
	if diags := CheckRevisionIdentity(prog, texts, bsln); len(diags) != 0 {
		t.Fatalf("expected clean identity, got %v", diags)
	}
}

// TestConstStaysSyntactic pins the termination rule: `n - m__ONE`
// (const m__ONE = 1) is not another spelling of the certified
// `n - 1` step, so the guarded recursion must fail.
func TestConstStaysSyntactic(t *testing.T) {
	src := `mod m
  provides [m__poll, M__S, m__ONE]
  uses []
  emits []

type M__S rev 1 (
  n: int
)

const m__ONE: int rev 1 = 1

fn m__poll(n: int) -> M__S rev 1
  decreases n
  emits []
  tests
    now(n = 0) => Ok(n = 0)
    later(n = 2) => Ok(n = 0)
=
  match n <= 0
    true => Ok(n = 0)
    false => match call m__poll(n - m__ONE)
      on Ok s => Ok(n = s.n)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if len(diags) == 0 {
		t.Fatalf("const-spelled step certified the recursion")
	}
}
