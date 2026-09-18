package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// One literal encoder (normStr) for values and patterns, and one
// target ordering (UTF-8 byte order via $canStr helpers) matching Go.

const strSemFixture = `mod str
  provides [str__match, str__order, str__len, str__at, str__slice, Str__Value, Int__Value]
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
    bslash("a\nb") => Ok(value = "hit")
    quote("say \"hi\"") => Ok(value = "q")
    plain("xy") => Ok(value = "miss")
=
  match v
    "a\nb" => Ok(value = "hit")
    "say \"hi\"" => Ok(value = "q")
    _ => Ok(value = "miss")

fn str__order(left: str, right: str) -> Int__Value rev 1
  emits []
  tests
    ascii("b", "a") => Ok(value = 1)
    prefix("ab", "a") => Ok(value = 1)
    equal("a", "a") => Ok(value = 0)
    bmp_smp("", "𐀀") => Ok(value = -1)
    smp_bmp("𐀀", "") => Ok(value = 1)
=
  match left == right
    true => Ok(value = 0)
    false => match left >= right
      true => Ok(value = 1)
      false => Ok(value = -1)

fn str__len(v: str) -> Int__Value rev 1
  emits []
  tests
    ascii("abc") => Ok(value = 3)
    mixed("héllo世界") => Ok(value = 7)
    astral("a𝄞b") => Ok(value = 3)
    empty("") => Ok(value = 0)
=
  Ok(value = #v)

fn str__at(v: str, i: int) -> Int__Value rev 1
  emits []
  tests
    ascii("abc", 1) => Ok(value = 98)
    mixed("héllo世界", 5) => Ok(value = 19990)
    astral("a𝄞b", 1) => Ok(value = 119070)
=
  Ok(value = v[i])

fn str__slice(v: str, a: int, b: int) -> Str__Value rev 1
  emits []
  tests
    inner("héllo", 1, 4) => Ok(value = "éll")
    full("abc", 0, 3) => Ok(value = "abc")
    empty("abc", 2, 2) => Ok(value = "")
=
  Ok(value = v[a:b])
`

func TestStrSemanticsEvaluate(t *testing.T) {
	// The bmp_smp/smp_bmp rows hinge on literal U+E000 vs U+10000
	// (invisible in most editors): fail loudly if they ever mangle.
	if !strings.Contains(strSemFixture, "\uE000") || !strings.Contains(strSemFixture, "\U00010000") {
		t.Fatal("fixture lost its U+E000/U+10000 literals")
	}
	dir := writeLSPDir(t, map[string]string{"str.can": strSemFixture})
	// Every expectation above is the source semantics: backslash and
	// quote literals match by raw content, and U+E000 orders below
	// U+10000 in byte order.
	if diags := diagnose(dir, "str.can", strSemFixture); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestStrSemanticsEmit(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "str.can"), []byte(strSemFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "str.can")}); err != nil {
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
		"$canStrGe(left, right)",
		"(BigInt([...v].length))",
		"$canStrAt(v, i)",
		"$canStrSlice(v, a, b)",
		"function $canStrCmp",
		"function $canStrGe",
		"function $canStrAt",
		"function $canStrSlice",
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
	if strings.Contains(src, "function $canStrLe") {
		t.Errorf("emit contains unused $canStrLe helper\n--- emit ---\n%s", src)
	}
}

func TestStrFaultEmit(t *testing.T) {
	// Row 2 (fault contracts): operator sites lower to helpers that
	// throw on invalid domains. The tables exercise valid inputs so
	// the module compiles; the pins assert the loud contract is
	// emitted, not merely present in a frozen golden.
	const faultFixture = `mod fault
  provides [fault__at, fault__slice, Fault__Int, Fault__Str]
  uses []
  emits []

type Fault__Int rev 1 (
  value: int
)

type Fault__Str rev 1 (
  value: str
)

fn fault__at(v: str, i: int) -> Fault__Int rev 1
  emits []
  tests
    go("abc", 1) => Ok(value = 98)
=
  Ok(value = v[i])

fn fault__slice(v: str, a: int, b: int) -> Fault__Str rev 1
  emits []
  tests
    go("abc", 0, 2) => Ok(value = "ab")
=
  Ok(value = v[a:b])
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "fault.can"), []byte(faultFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "fault.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(out, "fault.ts"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(got)
	for _, want := range []string{
		`throw new Error("str index out of range")`,
		`throw new Error("str slice out of range")`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("emit missing loud contract %q\n--- emit ---\n%s", want, src)
		}
	}
}
