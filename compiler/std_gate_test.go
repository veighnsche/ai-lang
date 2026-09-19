package main

import (
	"testing"
)

// Blessed-stdlib gates: compiling a std module executes every
// decision table, so a green compile pins the module's behavior
// in the suite. One test per migrated-or-new module below.

// TestStdMapCompiles pins the generic Map<K,V> migration: both
// instantiations check, execute, and emit.
func TestStdMapCompiles(t *testing.T) {
	out := t.TempDir()
	if err := compile(out, []string{"../std/map/map.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
}
