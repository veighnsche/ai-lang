package main

import (
	"path/filepath"
	"strings"
	"testing"
)

// R1 (v45): braces inside string literals are data, never
// delimiters, so a module wrapping values in braces compiles and
// its decision table runs.
const braceStringSrc = `mod m
  provides [m__wrap, Str__Value]
  uses []
  emits []

type Str__Value rev 1 (
  value: str
)

fn m__wrap(value: str) -> Str__Value rev 1
  emits []
  tests
    plain(value = "x") => Ok(value = "{x}")
    empty(value = "") => Ok(value = "{}")
=
  Ok(value = "{" + value + "}")
`

func TestBracesInStringsAllowed(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.ail": braceStringSrc})
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.ail")}); err != nil {
		t.Fatalf("braces in strings must compile: %v", err)
	}
}

// Braces in code or comments stay banned (R1).
func TestBracesOutsideStringsBanned(t *testing.T) {
	for name, src := range map[string]string{
		"code":    braceStringSrc + "  { stray }\n",
		"comment": braceStringSrc + "// { stray }\n",
		// A quote inside a comment must not open a string that
		// smuggles a brace past the ban.
		"comment_quote": braceStringSrc + "// say \"hi\" { stray }\n",
	} {
		dir := writeLSPDir(t, map[string]string{"m.ail": src})
		err := compile(t.TempDir(), []string{filepath.Join(dir, "m.ail")})
		if err == nil || !strings.Contains(err.Error(), "curly braces are banned") {
			t.Errorf("%s: expected braces error, got %v", name, err)
		}
	}
}
