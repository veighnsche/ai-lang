package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// checkGoldenFile fails the test with a (-want +got) diff when fresh
// compiler output differs from the committed golden. Centralized here
// so every golden gate reports what changed, not just which file.
func checkGoldenFile(t *testing.T, name string, got, want []byte) {
	t.Helper()
	checkGolden(t, name, string(got), string(want))
}

// checkGolden is the string form for goldens built in memory (normalize
// tables, JSON diagnostics) rather than read from files.
func checkGolden(t *testing.T, name, got, want string) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("golden mismatch: %s (-want +got):\n%s", name, diff)
	}
}
