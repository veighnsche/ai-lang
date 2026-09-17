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
  provides [num__add, num__big, num__lit, num__next, num__strict, Num__Out]
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

fn num__strict(a: dec, n: int, s: str) -> Num__Out rev 1
  emits []
  tests
    t1(a = d"0.5", n = 3, s = "b") => Ok(total = d"0.5", count = 2)
    t2(a = d"0.1", n = 3, s = "b") => Ok(total = d"0.1", count = -1)
    t3(a = d"0.5", n = 7, s = "b") => Ok(total = d"0.5", count = 0)
    t4(a = d"0.5", n = 3, s = "z") => Ok(total = d"0.5", count = 1)
    t5(a = d"0.2", n = 3, s = "b") => Ok(total = d"0.2", count = 0)
=
  match a > d"0.2"
    true => match n < 5
      true => match s < "m"
        true => Ok(total = a, count = 2)
        false => Ok(total = a, count = 1)
      false => Ok(total = a, count = 0)
    false => match a < d"0.2"
      true => Ok(total = a, count = -1)
      false => Ok(total = a, count = 0)
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
		"$ailDecGt(a,",
		"$ailDecLt(a,",
		"$ailStrLt(s,",
		"(n < 5n)",
		"function $ailDecGt",
		"function $ailDecLt",
		"function $ailStrLt",
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
	// own plumbing (dec scale handling, the 3-way str comparator);
	// every ail value is bigint or string.
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, ": number") &&
			!strings.Contains(line, "scale: number") &&
			!strings.Contains(line, "function $ailStrCmp(") {
			t.Errorf("emit maps an ail value to lossy number: %q", line)
		}
	}
}
