package main

import (
	"testing"
)

// a69: requires/ensures parse and store. No verifier, no
// checks, no emit effect. Probes first: the keywords do
// not parse yet.

// TestContractGrammarMax pins the max pilot shape:
// requires, Ok arm with two clauses plus the Boolean
// match-disjunction.
func TestContractGrammarMax(t *testing.T) {
	src := `mod m
  provides [m__max, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__max(left: int, right: int) -> M__Out rev 1
  emits []
  requires
    true
  ensures
    on Ok result
      result.value >= left
      result.value >= right
      match result.value == left
        true => true
        false => result.value == right
  tests
    ordered(left = 1, right = 2) => Ok(value = 2)
    reversed(left = 2, right = 1) => Ok(value = 2)
=
  match left <= right
    on true => Ok(value = right)
    on false => Ok(value = left)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// TestContractGrammarValidator pins the validator shape:
// two error arms with payload predicates, and ensures
// before emits (order independence).
func TestContractGrammarValidator(t *testing.T) {
	src := `mod m
  provides [m__check, M__Out]
  uses []
  emits [m.bad_bounds, m.out_of_range]

error m.bad_bounds(lower: int, upper: int)
error m.out_of_range(value: int, lower: int, upper: int)

type M__Out rev 1 (
  value: int
)

fn m__check(value: int, lower: int, upper: int) -> M__Out rev 1
  ensures
    on Ok result
      result.value == value
    on m.bad_bounds err
      lower >= upper + 1
      err.lower == lower
    on m.out_of_range err
      err.value == value
  emits [m.bad_bounds, m.out_of_range]
  tests
    ok(value = 1, lower = 0, upper = 2) => Ok(value = 1)
    low(value = -1, lower = 0, upper = 2) => m.out_of_range(value = -1, lower = 0, upper = 2)
=
  match value >= lower
    on true => Ok(value = value)
    on false => m.out_of_range(value = value, lower = lower, upper = upper)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// TestContractGrammarRejects pins the malformed shapes:
// unknown outcome, empty requires, empty ensures arm,
// duplicate blocks.
func TestContractGrammarRejects(t *testing.T) {
	base := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(value: int) -> M__Out rev 1
  emits []
`
	cases := map[string]string{
		"unknown outcome": base + `  ensures
    on Ok result
      result.value == value
    on m.nope err
      err.value == value
  tests
    go(value = 1) => Ok(value = 1)
=
  Ok(value = value)
`,
		"empty requires": base + `  requires
  tests
    go(value = 1) => Ok(value = 1)
=
  Ok(value = value)
`,
		"empty arm": base + `  ensures
    on Ok result
  tests
    go(value = 1) => Ok(value = 1)
=
  Ok(value = value)
`,
		"duplicate requires": base + `  requires
    true
  requires
    true
  tests
    go(value = 1) => Ok(value = 1)
=
  Ok(value = value)
`,
		"duplicate ensures": base + `  ensures
    on Ok result
      result.value == value
  ensures
    on Ok result
      result.value == value
  tests
    go(value = 1) => Ok(value = 1)
=
  Ok(value = value)
`,
	}
	for name, src := range cases {
		dir := writeLSPDir(t, map[string]string{"m.ail": src})
		diags := diagnose(dir, "m.ail", src)
		if !hasErrCode(diags, CodeParse) {
			t.Fatalf("%s: expected AIL1000, got %v", name, diags)
		}
	}
}
