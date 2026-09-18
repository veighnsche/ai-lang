package main

import (
	"strings"
	"testing"
)

// R3's single-separator reading was wrong: the grammar admits
// hierarchical names (std__int__abs) because the trailing verb class
// swallows further separators, while still rejecting separator abuse.
// Pin both sides: brief-style names verify clean, and malformed names
// keep their naming diagnostic.

const multiUnderscoreSrc = `mod m
  provides [std__int__abs, M__Val]
  uses []
  emits []

type M__Val rev 1 (
  value: int
)

fn std__int__abs(value: int) -> M__Val rev 1
  emits []
  tests
    neg(-3) => Ok(value = 3)
    pos(3) => Ok(value = 3)
  match value >= 0
    true => Ok(value = value)
    false => Ok(value = 0 - value)
`

func TestMultiUnderscoreNameAccepted(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": multiUnderscoreSrc})
	if diags := diagnose(dir, "m.can", multiUnderscoreSrc); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestMalformedNamesRejected(t *testing.T) {
	for _, name := range []string{"go", "_foo__bar", "foo__", "foo____bar"} {
		bad := strings.ReplaceAll(multiUnderscoreSrc, "std__int__abs", name)
		dir := writeLSPDir(t, map[string]string{"m.can": bad})
		if diags := diagnose(dir, "m.can", bad); !hasDiag(diags, "error", "must match domain__verb") {
			t.Fatalf("expected naming error for %q, got %v", name, diags)
		}
	}
}
