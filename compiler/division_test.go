package main

import (
	"os"
	"strings"
	"testing"
)

// a17: / and % are exact Euclidean integer division (a == b*q + r
// with 0 <= r < |b| on every sign combination), discharged through
// big.Int.DivMod — no unit-step scans. Dec operands are refused
// statically (CAN6005); a zero divisor is loud at runtime.

func TestDivModVectors(t *testing.T) {
	for _, c := range []struct{ expr, want string }{
		{`7 / 3`, "2"},
		{`7 % 3`, "1"},
		{`-7 / 3`, "-3"},
		{`-7 % 3`, "2"},
		// Negative literals cannot sit right of an operator
		// (pre-existing: unary minus is literal-only, so the +/-
		// level would split there first); parenthesize instead.
		{`7 / (0 - 3)`, "-2"},
		{`7 % (0 - 3)`, "1"},
		{`-7 / (0 - 3)`, "3"},
		{`-7 % (0 - 3)`, "2"},
		{`0 / 5`, "0"},
		{`0 % 5`, "0"},
		{`10 - 8 / 2`, "6"},
		{`100 / 10 / 2`, "5"},
		{`2 + 7 % 3`, "3"},
	} {
		wantInt(t, c.expr, c.want)
	}
}

func TestDivZeroLoud(t *testing.T) {
	for _, expr := range []string{`1 / 0`, `1 % 0`} {
		sm, err := parseSmall(expr)
		if err != nil {
			t.Fatalf("parseSmall(%q): %v", expr, err)
		}
		if _, err := evSmall(sm, map[string]*Value{}, &Ctx{}, "test"); err == nil {
			t.Fatalf("evSmall(%q): expected division-by-zero error, got none", expr)
		}
	}
}

func TestDivMixedLoud(t *testing.T) {
	// Direct evaluator callers bypass the static gate: mixed and
	// non-int operands must error here, never panic on nil. (A
	// string-first expression never parses — quotes swallow
	// operators — so the str case arrives erased through seal.)
	for _, expr := range []string{`seal B("a") / 2`, `d"1.0" / d"2.0"`, `1 / d"2.0"`} {
		sm, err := parseSmall(expr)
		if err != nil {
			t.Fatalf("parseSmall(%q): %v", expr, err)
		}
		if _, err := evSmall(sm, map[string]*Value{}, &Ctx{}, "test"); err == nil {
			t.Fatalf("evSmall(%q): expected error, got none", expr)
		}
	}
}

const typeDivDec = `mod m
  provides [m__div, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  q: dec
)

fn m__div(a: dec, b: dec) -> M__Out rev 1
  emits []
  tests
    t(d"1.0", d"2.0") => Ok(q = d"0.5")
=
  Ok(q = a / b)
`

func TestDecDivisionRefused(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": typeDivDec})
	diags := diagnose(dir, "m.can", typeDivDec)
	if !hasDiag(diags, "error", "has no exact result") {
		t.Fatalf("expected CAN6005 inexact-division error, got %v", diags)
	}
	found := false
	for _, d := range diags {
		if d.Code == "CAN6005" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a coded CAN6005 diagnostic, got %v", diags)
	}
}

func TestDecRemainderRefused(t *testing.T) {
	bad := strings.Replace(typeDivDec, "Ok(q = a / b)", "Ok(q = a % b)", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "has no exact result") {
		t.Fatalf("expected CAN6005 inexact-division error, got %v", diags)
	}
}

func TestDivEmitHelper(t *testing.T) {
	// The Euclidean helper emits inline only when / or % is used.
	dir := writeLSPDir(t, map[string]string{"m.can": typeDivMod})
	out := t.TempDir()
	if err := compile(out, []string{dir + "/m.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(out + "/m.ts")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "function $canDivMod") {
		t.Fatalf("emit missing $canDivMod helper:\n%s", raw)
	}
	if !strings.Contains(string(raw), "$canDivMod(a, 3n)[0]") {
		t.Fatalf("emit missing helper call:\n%s", raw)
	}
}

const typeDivMod = `mod m
  provides [m__div, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  q: int
)

fn m__div(a: int) -> M__Out rev 1
  emits []
  tests
    t(7) => Ok(q = 2)
=
  Ok(q = a / 3)
`
