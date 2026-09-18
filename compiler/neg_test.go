package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Slice 6: prefix unary minus for exact int and dec. Probes
// first: tight binding (a * -b, a - -b), literal stability
// (-3, d"-0.5" keep bytes; spaced negations normalize),
// emit shapes ((-value), $canDecSub("0.0", value)), no unary
// plus, and bool/str refusals.

// TestNegBasic pins integer and decimal negation end to end.
func TestNegBasic(t *testing.T) {
	src := `mod m
  provides [m__go, m__dec, M__Out, M__Dec]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

type M__Dec rev 1 (
  value: dec
)

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    pos(3) => Ok(-3)
    neg(-3) => Ok(3)
    zero(0) => Ok(0)
  Ok(-x)

fn m__dec(x: dec) -> M__Dec rev 1
  emits []
  tests
    pos(d"1.5") => Ok(d"-1.5")
    neg(d"-1.5") => Ok(d"1.5")
    zero(d"0.0") => Ok(d"0.0")
  Ok(-x)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("negation reported: %v", diags)
	}
}

// TestNegTight parses the required tight shapes: a * -b and
// a - -b read with the minus on b, evaluating exactly.
func TestNegTight(t *testing.T) {
	src := `mod m
  provides [m__go, m__sub, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(a: int, b: int) -> M__Out rev 1
  emits []
  tests
    row(3, 2) => Ok(-6)
  Ok(a * -b)

fn m__sub(a: int, b: int) -> M__Out rev 1
  emits []
  tests
    row(3, 2) => Ok(5)
  Ok(a - -b)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("tight negation reported: %v", diags)
	}
}

// TestNegEmitShape pins the lowering: integers emit
// (-value), decimals reuse $canDecSub("0.0", value) and never
// JavaScript -x.
func TestNegEmitShape(t *testing.T) {
	src := `mod m
  provides [m__go, m__dec, M__Out, M__Dec]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

type M__Dec rev 1 (
  value: dec
)

fn m__go(a: int, b: int) -> M__Out rev 1
  emits []
  tests
    row(3, 2) => Ok(-6)
  Ok(-a * b)

fn m__dec(x: dec) -> M__Dec rev 1
  emits []
  tests
    pos(d"1.5") => Ok(d"-1.5")
  Ok(-x)
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.can"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	ts := string(raw)
	for _, want := range []string{"((-a) * b)", `$canDecSub("0.0", x)`} {
		if !strings.Contains(ts, want) {
			t.Fatalf("emit missing %q", want)
		}
	}
	// Decimals never lower to JavaScript unary minus on a string.
	if strings.Contains(ts, `-("`) {
		t.Fatalf("decimal negation lowered to JavaScript -x")
	}
}

// TestNegLiteralsStable pins literal behavior: -3 and
// d"-0.5" keep their bytes, and spaced negations (- 3)
// normalize back to the literal forms.
func TestNegLiteralsStable(t *testing.T) {
	src := `mod m
  provides [m__go, m__dec, M__Out, M__Dec]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

type M__Dec rev 1 (
  value: dec
)

fn m__go() -> M__Out rev 1
  emits []
  tests
    go() => Ok(-3)
  Ok(- 3)

fn m__dec() -> M__Dec rev 1
  emits []
  tests
    go() => Ok(d"-0.5")
  Ok(-d"0.5")
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.can"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	ts := string(raw)
	for _, want := range []string{"-3n", `"-0.5"`} {
		if !strings.Contains(ts, want) {
			t.Fatalf("emit missing literal %q", want)
		}
	}
	for _, want := range []string{"(-3n)", `$canDecSub("0.0", "0.5")`} {
		if strings.Contains(ts, want) {
			t.Fatalf("literal did not normalize: %q", want)
		}
	}
}

// TestNegNoPlus guards the asymmetry: unary plus stays a
// parse error.
func TestNegNoPlus(t *testing.T) {
	for _, body := range []string{
		"Ok(+b)",
		"Ok(a + +b)",
	} {
		src := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(a: int, b: int) -> M__Out rev 1
  emits []
  tests
    row(3, 2) => Ok(5)
  ` + body + "\n"
		dir := writeLSPDir(t, map[string]string{"m.can": src})
		if diags := diagnose(dir, "m.can", src); !hasError(diags) {
			t.Fatalf("%q reported nothing", body)
		}
	}
}

// TestNegTypes pins int/dec-only negation: bool and string
// operands are CAN6003.
func TestNegTypes(t *testing.T) {
	src := `mod m
  provides [m__go, m__str, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: bool
)

fn m__go(x: bool) -> M__Out rev 1
  emits []
  tests
    go(true) => Ok(true)
  Ok(-x)

fn m__str(y: str) -> M__Out rev 1
  emits []
  tests
    go("a") => Ok(true)
  Ok(-y)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN6003") {
		t.Fatalf("mistyped negation reported no CAN6003: %v", diags)
	}
}
