package main

import (
	"strings"
	"testing"
)

// Positional test args resolve front-to-back through the parameter
// list; named args keep their names (CAN3205 rule). diagnose runs
// every decision table, so clean diags prove both acceptance and
// evaluation order: ordered(5, 8) => -3 passes only if 5 binds left.

const positionalRows = `mod demo
  provides [demo__sub, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__sub(left: int, right: int) -> Int__Value rev 1
  emits []
  tests
    named(left = 5, right = 8) => Ok(-3)
    mixed(5, right = 8) => Ok(-3)
    ordered(5, 8) => Ok(-3)
  Ok(left - right)
`

func TestPositionalTestArgsEvaluate(t *testing.T) {
	// Every mixed shape evaluates (tables run: order proven by
	// execution), and the in-slot names are legal (no CAN3205)
	// but unidiomatic: exactly the three redundant findings
	// report — named left, named right, mixed right — nothing else.
	dir := writeLSPDir(t, map[string]string{"demo.can": positionalRows})
	diags := diagnose(dir, "demo.can", positionalRows)
	redundant := 0
	for _, d := range diags {
		if d.Code == CodeLintRedundant {
			redundant++
			continue
		}
		if d.Sev == "error" {
			t.Fatalf("unexpected error beside the redundant findings: %v", diags)
		}
	}
	if redundant != 3 {
		t.Fatalf("expected exactly three CAN3410 findings, got %v", diags)
	}
}

func TestPositionalTestArgFaultsAreCAN3205(t *testing.T) {
	cases := []struct {
		name string
		row  string
		want string
	}{
		{"after named", "bad(left = 5, 8) => Ok(-3)", "takes positional arg after named arg"},
		{"over arity", "bad(5, 8, 9) => Ok(-3)", "takes 3 args for 2 params"},
		{"double claim", "bad(5, left = 6) => Ok(-3)", "supplies arg left twice"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := strings.Replace(positionalRows,
				"    ordered(5, 8) => Ok(-3)", "    "+tc.row, 1)
			dir := writeLSPDir(t, map[string]string{"demo.can": src})
			diags := diagnose(dir, "demo.can", src)
			if !hasDiag(diags, "error", tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, diags)
			}
			if !hasCode(diags, "CAN3205") {
				t.Fatalf("expected CAN3205, got %v", diags)
			}
		})
	}
}
