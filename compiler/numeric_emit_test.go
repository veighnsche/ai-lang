package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEmitNumerics pins v10's exact target: ints emit as bigint,
// decs as canonical-digit strings through exact helpers. The fixture
// exercises arithmetic and ordering over refs (typed dispatch via
// checker annotations) and literals in bodies.
const numericFixture = `mod num
  provides [num__add, num__big, num__lit, num__next, Num__Out]
  uses []
  emits []

type Num__Out rev 1 (
  total: dec
  count: int
)

fn num__add(a: dec, b: dec) -> Num__Out rev 1
  emits []
  tests
    t1(a = d"0.1", b = d"0.2") => Ok(total = d"0.3", count = 0)
=
  Ok(total = a + b, count = 0)

fn num__big(a: dec, b: dec) -> Num__Out rev 1
  emits []
  tests
    t1(a = d"0.5", b = d"0.2") => Ok(total = d"0.5", count = 1)
    t2(a = d"0.1", b = d"0.9") => Ok(total = d"0.9", count = 0)
=
  match a >= b
    true => Ok(total = a, count = 1)
    false => Ok(total = b, count = 0)

fn num__lit() -> Num__Out rev 1
  emits []
  tests
    t1() => Ok(total = d"3.14", count = 7)
=
  Ok(total = d"3.14", count = 7)

fn num__next(n: int) -> Num__Out rev 1
  emits []
  tests
    t1(n = 41) => Ok(total = d"1.0", count = 42)
=
  Ok(total = d"1.0", count = n + 1)
`

func compileFixture(t *testing.T, name, src string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(src), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, name)}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(out, "num.ts"))
	if err != nil {
		t.Fatalf("read emit: %v", err)
	}
	return string(got)
}

func TestEmitNumerics(t *testing.T) {
	got := compileFixture(t, "num.ail", numericFixture)
	// compile() already ran every decision table, so the expectations
	// above passed under exact proof semantics.
	for _, want := range []string{
		"total: string",
		"count: bigint",
		"a: string",
		"n: bigint",
		"$ailDecAdd(a, b)",
		"$ailDecGe(a, b)",
		`"3.14"`,
		"(n + 1n)",
		"function $ailDecSplit",
		"function $ailDecAdd",
		"function $ailDecGe",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("emit missing %q\n--- emit ---\n%s", want, got)
		}
	}
	for _, banned := range []string{"$ailDecSub", "$ailDecMul", "$ailDecLe"} {
		if strings.Contains(got, banned) {
			t.Errorf("emit contains unused helper %q\n--- emit ---\n%s", banned, got)
		}
	}
	// The only number-typed values in exact emit are the helpers'
	// own scale plumbing; every ail value is bigint or string.
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, ": number") && !strings.Contains(line, "scale: number") {
			t.Errorf("emit maps an ail value to lossy number: %q", line)
		}
	}
}
