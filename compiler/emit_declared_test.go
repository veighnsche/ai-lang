package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Result types come from declared return records, never from test
// literals: changing evidence without changing the signature or body
// cannot change the emitted function's type, and computed or referenced
// payloads need no literal witness.

const declaredHead = `mod audit
  provides [PROVIDES]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

type Bool__Value rev 1 (
  value: bool
)

type Dec__Value rev 1 (
  value: dec
)

type Str__Value rev 1 (
  value: str
)

`

const declaredFns = `fn audit__int(value: int) -> Int__Value rev 1
  emits []
  tests
    i(value = 1) => Ok(value = value)
=
  Ok(value = value)

fn audit__bool(flag: bool) -> Bool__Value rev 1
  emits []
  tests
    b(flag = true) => Ok(value = flag)
=
  Ok(value = flag)

fn audit__dec(amount: dec) -> Dec__Value rev 1
  emits []
  tests
    d(amount = d"1.5") => Ok(value = amount)
=
  Ok(value = amount)

fn audit__str(name: str) -> Str__Value rev 1
  emits []
  tests
    s(name = "x") => Ok(value = name)
=
  Ok(value = name)
`

func declaredSrc() string {
	return strings.Replace(declaredHead, "PROVIDES",
		"audit__int, audit__bool, audit__dec, audit__str, Int__Value, Bool__Value, Dec__Value, Str__Value", 1) + declaredFns
}

func TestDeclaredResultTypes(t *testing.T) {
	src := declaredSrc()
	dir := writeLSPDir(t, map[string]string{"audit.can": src})
	// No Ok literal anywhere: every payload is a reference, yet every
	// function carries its declared result through checking and eval.
	if diags := diagnose(dir, "audit.can", src); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestDeclaredResultEmit(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "audit.can"), []byte(declaredSrc()), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "audit.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(out, "audit.ts"))
	if err != nil {
		t.Fatal(err)
	}
	ts := string(got)
	// One file, four declared results: bool, int, and str-shaped (dec
	// erases to string, sharing the identical member) each survive.
	for _, want := range []string{
		`{ $can_kind: "ok"; value: bigint }`,
		`{ $can_kind: "ok"; value: boolean }`,
		`{ $can_kind: "ok"; value: string }`,
		"export function audit__int(value: bigint)",
		"export function audit__bool(flag: boolean)",
	} {
		if !strings.Contains(ts, want) {
			t.Errorf("emit missing %q\n--- emit ---\n%s", want, ts)
		}
	}
}

func TestEvidenceCannotRetype(t *testing.T) {
	body := `fn audit__two(value: int) -> Int__Value rev 1
  emits []
  tests
    two(value = 2) => EXPECTED
=
  Ok(value = value)
`
	head := `mod audit
  provides [audit__two, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

`
	emit := func(t *testing.T, expected string) string {
		t.Helper()
		src := head + strings.Replace(body, "EXPECTED", expected, 1)
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "audit.can"), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
		out := t.TempDir()
		if err := compile(out, []string{filepath.Join(dir, "audit.can")}); err != nil {
			t.Fatalf("compile: %v", err)
		}
		got, err := os.ReadFile(filepath.Join(out, "audit.ts"))
		if err != nil {
			t.Fatal(err)
		}
		return string(got)
	}
	// Same signature, same body, different evidence: the literal-free
	// variant has no witness at all, yet both compile ...
	a := emit(t, "Ok(value = value)")
	b := emit(t, "Ok(value = 2)")
	// ... to the same production code and signature.
	if prod(a) != prod(b) {
		t.Errorf("evidence changed production code:\n--- without ---\n%s\n--- with ---\n%s", prod(a), prod(b))
	}
}

// prod extracts the production contract of audit__two: its signature
// and the module result type. Test tables never appear in prod emit.
func prod(ts string) string {
	var keep []string
	for _, ln := range strings.Split(ts, "\n") {
		if strings.HasPrefix(ln, "export function audit__two") ||
			strings.HasPrefix(ln, "export type AuditResult") {
			keep = append(keep, ln)
		}
	}
	return strings.Join(keep, "\n")
}
