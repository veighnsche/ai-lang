package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// One literal encoder (normStr) for values and patterns, and one
// target ordering (UTF-8 byte order via $ailStr helpers) matching Go.

const strSemFixture = `mod str
  provides [str__match, str__order, Str__Value, Int__Value]
  uses []
  emits []

type Str__Value rev 1 (
  value: str
)

type Int__Value rev 1 (
  value: int
)

fn str__match(v: str) -> Str__Value rev 1
  emits []
  tests
    bslash(v = "a\nb") => Ok(value = "hit")
    quote(v = "say \"hi\"") => Ok(value = "q")
    plain(v = "xy") => Ok(value = "miss")
=
  match v
    "a\nb" => Ok(value = "hit")
    "say \"hi\"" => Ok(value = "q")
    _ => Ok(value = "miss")

fn str__order(left: str, right: str) -> Int__Value rev 1
  emits []
  tests
    ascii(left = "b", right = "a") => Ok(value = 1)
    prefix(left = "ab", right = "a") => Ok(value = 1)
    equal(left = "a", right = "a") => Ok(value = 0)
    bmp_smp(left = "", right = "𐀀") => Ok(value = -1)
    smp_bmp(left = "𐀀", right = "") => Ok(value = 1)
=
  match left == right
    true => Ok(value = 0)
    false => match left >= right
      true => Ok(value = 1)
      false => Ok(value = -1)
`

func TestStrSemanticsEvaluate(t *testing.T) {
	// The bmp_smp/smp_bmp rows hinge on literal U+E000 vs U+10000
	// (invisible in most editors): fail loudly if they ever mangle.
	if !strings.Contains(strSemFixture, "\uE000") || !strings.Contains(strSemFixture, "\U00010000") {
		t.Fatal("fixture lost its U+E000/U+10000 literals")
	}
	dir := writeLSPDir(t, map[string]string{"str.ail": strSemFixture})
	// Every expectation above is the source semantics: backslash and
	// quote literals match by raw content, and U+E000 orders below
	// U+10000 in byte order.
	if diags := diagnose(dir, "str.ail", strSemFixture); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestStrSemanticsEmit(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "str.ail"), []byte(strSemFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "str.ail")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(out, "str.ts"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(got)
	// Patterns use the same encoder as values: the raw backslash and
	// the raw quotes round-trip instead of becoming TS escapes.
	for _, want := range []string{
		`(v === "a\\nb")`,
		`(v === "say \\\"hi\\\"")`,
		"$ailStrGe(left, right)",
		"function $ailStrCmp",
		"function $ailStrGe",
		"TextEncoder",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("emit missing %q\n--- emit ---\n%s", want, src)
		}
	}
	// No native string ordering remains; the old unescaped pattern
	// form is gone.
	for _, banned := range []string{"(left >= right)", `(v === "a\nb")`} {
		if strings.Contains(src, banned) {
			t.Errorf("emit contains %q\n--- emit ---\n%s", banned, src)
		}
	}
	if strings.Contains(src, "function $ailStrLe") {
		t.Errorf("emit contains unused $ailStrLe helper\n--- emit ---\n%s", src)
	}
}
