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
    named(left = 5, right = 8) => Ok(value = -3)
    mixed(5, right = 8) => Ok(value = -3)
    ordered(5, 8) => Ok(value = -3)
=
  Ok(value = left - right)
`

func TestPositionalTestArgsEvaluate(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"demo.can": positionalRows})
	if diags := diagnose(dir, "demo.can", positionalRows); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestPositionalTestArgFaultsAreCAN3205(t *testing.T) {
	cases := []struct {
		name string
		row  string
		want string
	}{
		{"after named", "bad(left = 5, 8) => Ok(value = -3)", "takes positional arg after named arg"},
		{"over arity", "bad(5, 8, 9) => Ok(value = -3)", "takes 3 args for 2 params"},
		{"double claim", "bad(5, left = 6) => Ok(value = -3)", "supplies arg left twice"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := strings.Replace(positionalRows,
				"    ordered(5, 8) => Ok(value = -3)", "    "+tc.row, 1)
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
