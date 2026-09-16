package main

import "testing"

func TestVersionFlag(t *testing.T) {
	for _, argv := range [][]string{{"--version"}, {"-version"}, {"version"}} {
		if code := run(argv); code != 0 {
			t.Fatalf("run(%q) = %d, want 0", argv, code)
		}
	}
}
