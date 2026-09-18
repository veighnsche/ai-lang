package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Equality lowers by resolved operand type: scalars compare natively
// (exact in both runtimes), records compare field-by-field over their
// declared shape. Distinct equal records compare equal; identity never
// decides.

const eqFixture = `mod audit
  provides [audit__cmp, audit__make, Audit__Pair]
  uses []
  emits []

type Audit__Pair rev 1 (
  first: int
  second: str
)

fn audit__make(first: int, second: str) -> Audit__Pair rev 1
  emits []
  tests
    m(first = 1, second = "a") => Ok(first = 1, second = "a")
=
  Ok(first = first, second = second)

fn audit__cmp(flag: str) -> Audit__Pair rev 1
  emits []
  tests
    s(flag = "a") => Ok(first = 1, second = "a")
    d(flag = "b") => Ok(first = 0, second = "neq")
=
  match call audit__make(first = 1, second = "a")
    on Ok left => match call audit__make(first = 1, second = flag)
      on Ok right => match left == right
        true => Ok(first = left.first, second = left.second)
        false => Ok(first = 0, second = "neq")
`

func TestRecordEqualityEvaluates(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"audit.can": eqFixture})
	// The evaluator compares structurally: distinct equal records
	// match (same), one changed field does not (diff).
	if diags := diagnose(dir, "audit.can", eqFixture); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestRecordEqualityEmit(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "audit.can"), []byte(eqFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "audit.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(out, "audit.ts"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(got)
	for _, want := range []string{
		`$canEqRec(left,right,["first","second"])`,
		"function $canEqRec",
		"function $canEqVal",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("emit missing %q\n--- emit ---\n%s", want, src)
		}
	}
	// Same-object identity would read (left === right): the lowering
	// must not depend on allocation.
	if strings.Contains(src, "(left===right)") {
		t.Errorf("emit compares records by identity:\n%s", src)
	}
}
