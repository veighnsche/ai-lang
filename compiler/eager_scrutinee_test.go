package main

// Scrutinee-trap lint (AIL4109): a multi-scrutinee match evaluates every
// scrutinee eagerly, so a trapping operation (string index/slice,
// division/remainder, calls) warns; pure tables, literal-nonzero
// divisors, and arity-1 matches stay silent.

import (
	"math/big"
	"strings"
	"testing"
)

const eagerText = `mod m
  provides [m__f, m__g, m__h, m__p, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  v: bool
)

fn m__f(n: int) -> M__Out rev 1
  emits []
  tests
    z(n = 0) => Ok(v = false)
    t(n = 1) => Ok(v = true)
    f(n = 2) => Ok(v = false)
=
  match n == 0, 10 % (n + 1) == 0
    true, _ => Ok(v = false)
    false, true => Ok(v = true)
    false, false => Ok(v = false)

fn m__g(a: bool, b: bool) -> M__Out rev 1
  emits []
  tests
    tt(a = true, b = true) => Ok(v = true)
    tf(a = true, b = false) => Ok(v = false)
    ft(a = false, b = true) => Ok(v = false)
    ff(a = false, b = false) => Ok(v = true)
=
  match a, b
    true, true => Ok(v = true)
    true, false => Ok(v = false)
    false, true => Ok(v = false)
    false, false => Ok(v = true)

fn m__h(n: int) -> M__Out rev 1
  emits []
  tests
    t(n = 2) => Ok(v = true)
    f(n = 3) => Ok(v = false)
=
  match 10 % n == 0
    true => Ok(v = true)
    false => Ok(v = false)

fn m__p(n: int, b: bool) -> M__Out rev 1
  emits []
  tests
    et(n = 2, b = true) => Ok(v = true)
    ef(n = 2, b = false) => Ok(v = false)
    ot(n = 3, b = true) => Ok(v = false)
    of(n = 3, b = false) => Ok(v = true)
=
  match n % 2 == 0, b
    true, true => Ok(v = true)
    true, false => Ok(v = false)
    false, true => Ok(v = false)
    false, false => Ok(v = true)
`

func eagerDiags(t *testing.T) []Diag {
	t.Helper()
	dir := writeLSPDir(t, map[string]string{"m.ail": eagerText})
	return diagnose(dir, "m.ail", eagerText)
}

func TestEagerScrutineeWarns(t *testing.T) {
	var warns []Diag
	for _, d := range eagerDiags(t) {
		if d.Sev == "error" {
			t.Fatalf("unexpected error diagnostic: %v", d)
		}
		if d.Code == CodeEagerScrutinee {
			warns = append(warns, d)
		}
	}
	if len(warns) != 1 {
		t.Fatalf("expected exactly one AIL4109 warning, got %v", warns)
	}
	w := warns[0]
	if w.Sev != "warning" {
		t.Fatalf("AIL4109 must be warning-grade, got %q", w.Sev)
	}
	if !strings.Contains(w.Msg, "scrutinee 2") || !strings.Contains(w.Msg, "remainder") || !strings.Contains(w.Msg, "eagerly") {
		t.Fatalf("warning names slot, cause, and hazard, got %q", w.Msg)
	}
}

func TestEagerScrutineeSilent(t *testing.T) {
	// Pure tables (m__g), arity-1 trapping scrutinees (m__h), and
	// literal-nonzero divisors (m__p) produce no AIL4109. The warns test
	// above already pins the total count at one; here the pipeline must
	// additionally be error-free on those functions, which it checks by
	// failing on any error diagnostic.
	for _, d := range eagerDiags(t) {
		if d.Sev == "error" {
			t.Fatalf("unexpected error diagnostic: %v", d)
		}
	}
}

func TestFirstTrappingOp(t *testing.T) {
	intLit := func(v int64) *Small { return &Small{Kind: "int", Num: big.NewInt(v)} }
	cases := []struct {
		name string
		in   *Small
		want string
	}{
		{"index", &Small{Kind: "stridx"}, "string index"},
		{"slice", &Small{Kind: "strslice"}, "string slice"},
		{"divRef", &Small{Kind: "binop", Op: "/", R: &Small{Kind: "ref"}}, "division"},
		{"modRef", &Small{Kind: "binop", Op: "%", R: &Small{Kind: "ref"}}, "remainder"},
		{"modLit", &Small{Kind: "binop", Op: "%", R: intLit(2)}, ""},
		{"modZeroLit", &Small{Kind: "binop", Op: "%", R: intLit(0)}, "remainder"},
		{"pure", &Small{Kind: "binop", Op: "=="}, ""},
		{"call", &Small{Kind: "call", Fname: "db__get"}, "call to db__get"},
		{"nested", &Small{Kind: "binop", Op: "==", L: &Small{Kind: "stridx"}}, "string index"},
	}
	for _, tc := range cases {
		if got := firstTrappingOp(tc.in); got != tc.want {
			t.Fatalf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}
