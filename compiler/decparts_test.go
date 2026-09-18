package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// a21: dec__parts is the blessed decimal observation kernel — one
// positional dec operand in, Ok(coefficient, scale) out, total. These
// tests pin the kernel contract directly, its static misuse shapes,
// and its TS helper wiring. Decision-table coverage of the blessed
// std__dec__parts wrapper lives in std/scalars/scalars.can and runs
// at every compile.

func decPartsScrut(t *testing.T, arg string) *Small {
	t.Helper()
	v, err := parseSmall(arg)
	if err != nil {
		t.Fatalf("parseSmall(%q): %v", arg, err)
	}
	return &Small{Kind: "call", Fname: "dec__parts", Args: []Arg{{V: v}}}
}

func TestDecPartsKernel(t *testing.T) {
	for _, c := range []struct {
		arg         string
		coefficient string
		scale       string
	}{
		{`d"7.0"`, "70", "1"},
		{`d"0.25"`, "25", "2"},
		{`d"-0.5"`, "-5", "1"},
		{`d"12.340"`, "1234", "2"},
		{`d"0.00"`, "0", "1"},
		{`d"-0.0"`, "0", "1"},
		{`d"100.0"`, "1000", "1"},
		{`d"123456789.123456789"`, "123456789123456789", "9"},
	} {
		got, err := evDecPartsOp(decPartsScrut(t, c.arg), map[string]*Value{}, &Ctx{}, "test")
		if err != nil {
			t.Fatalf("evDecPartsOp(%s): %v", c.arg, err)
		}
		if got.Kind != "ok" {
			t.Fatalf("evDecPartsOp(%s) kind = %s, want ok", c.arg, got.Kind)
		}
		if got.Dict["coefficient"].N.String() != c.coefficient || got.Dict["scale"].N.String() != c.scale {
			t.Fatalf("evDecPartsOp(%s) = (%s, %s), want (%s, %s)", c.arg,
				got.Dict["coefficient"].N, got.Dict["scale"].N, c.coefficient, c.scale)
		}
	}
}

func TestDecPartsKernelStoredDigits(t *testing.T) {
	// The kernel observes stored digits; canonicality is the
	// language invariant (literals canonicalize at parse, arithmetic
	// re-normalizes), not a kernel re-check. The TS helper mirrors
	// this raw reading, so parity holds by construction.
	v, err := evDecPartsOp(&Small{Kind: "call", Fname: "dec__parts",
		Args: []Arg{{V: &Small{Kind: "dec", Dec: "12.340"}}}}, map[string]*Value{}, &Ctx{}, "test")
	if err != nil {
		t.Fatalf("evDecPartsOp(hand-built): %v", err)
	}
	if v.Dict["coefficient"].N.String() != "12340" || v.Dict["scale"].N.String() != "3" {
		t.Fatalf("hand-built 12.340 = (%s, %s), want (12340, 3)",
			v.Dict["coefficient"].N, v.Dict["scale"].N)
	}
}

func TestDecPartsKernelLoud(t *testing.T) {
	// Direct evaluator callers bypass the static gate: non-dec
	// operands and bad arity must error here, never panic on nil.
	bad, err := parseSmall("7")
	if err != nil {
		t.Fatalf("parseSmall: %v", err)
	}
	str, err := parseSmall(`"ab"`)
	if err != nil {
		t.Fatalf("parseSmall: %v", err)
	}
	for _, s := range []*Small{
		{Kind: "call", Fname: "dec__parts", Args: []Arg{{V: bad}}},
		{Kind: "call", Fname: "dec__parts", Args: []Arg{{V: str}}},
		{Kind: "call", Fname: "dec__parts", Args: nil},
		{Kind: "call", Fname: "dec__parts", Args: []Arg{{V: bad}, {V: bad}}},
		{Kind: "call", Fname: "dec__parts", Args: []Arg{{V: bad, HasName: true}}},
	} {
		if _, err := evDecPartsOp(s, map[string]*Value{}, &Ctx{}, "test"); err == nil {
			t.Fatalf("evDecPartsOp(%v): expected error, got none", s)
		}
	}
}

const typeDecPartsMisuse = `mod m
  provides [m__str, m__named, m__two, m__given, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  n: int
)

fn m__str(v: str) -> M__Out rev 1
  emits []
  tests
    t(v = "ab") => Ok(n = 0)
=
  match call dec__parts(v)
    on Ok p => Ok(n = p.coefficient)

fn m__named(v: dec) -> M__Out rev 1
  emits []
  tests
    t(v = d"1.0") => Ok(n = 0)
=
  match call dec__parts(value = v)
    on Ok p => Ok(n = p.coefficient)

fn m__two(v: dec) -> M__Out rev 1
  emits []
  tests
    t(v = d"1.0") => Ok(n = 0)
=
  match call dec__parts(v, v)
    on Ok p => Ok(n = p.coefficient)

fn m__given(v: dec) -> M__Out rev 1
  emits []
  tests
    t(v = d"1.0") => Ok(n = 0)
=
  match call dec__parts(v)
    given
      t => [exchange args (value = d"1.0") outcome Ok(coefficient = 1, scale = 0)]
    on Ok p => Ok(n = p.coefficient)
`

func TestDiagnoseDecPartsMisuse(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": typeDecPartsMisuse})
	diags := diagnose(dir, "m.can", typeDecPartsMisuse)
	for _, want := range []string{
		"call dec__parts value: got str, want dec",
		"call dec__parts takes positional args",
		"call dec__parts takes 1 arg",
		"call to dec__parts takes no given table",
	} {
		if !hasDiag(diags, "error", want) {
			t.Fatalf("expected %q, got %v", want, diags)
		}
	}
}

const typeDecPartsTotal = `mod m
  provides [m__stale, m__missing, M__Out]
  uses []
  emits [m.nope]

error m.nope()

type M__Out rev 1 (
  n: int
)

fn m__stale(v: dec) -> M__Out rev 1
  emits []
  tests
    t(v = d"1.0") => Ok(n = 1)
=
  match call dec__parts(v)
    on Ok p => Ok(n = p.coefficient)
    on m.nope _ => Ok(n = 0)

fn m__missing(v: dec) -> M__Out rev 1
  emits []
  tests
    t(v = d"1.0") => Ok(n = 0)
=
  match call dec__parts(v)
    on m.nope _ => Ok(n = 0)
`

func TestDiagnoseDecPartsTotal(t *testing.T) {
	// The kernel emits nothing, so want = {ok}: an error arm is
	// stale, and an error-only match misses ok. Totality at the gate.
	dir := writeLSPDir(t, map[string]string{"m.can": typeDecPartsTotal})
	diags := diagnose(dir, "m.can", typeDecPartsTotal)
	for _, want := range []string{
		"stale match arm m.nope",
		"non-exhaustive match, missing ok",
	} {
		if !hasDiag(diags, "error", want) {
			t.Fatalf("expected %q, got %v", want, diags)
		}
	}
}

func TestDecPartsEmitHelper(t *testing.T) {
	// The $canDecParts helper emits only when the kernel is used.
	used := `mod m
  provides [m__parts, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  coefficient: int
  scale: int
)

fn m__parts(v: dec) -> M__Out rev 1
  emits []
  tests
    t(v = d"1.0") => Ok(coefficient = 10, scale = 1)
=
  match call dec__parts(v)
    on Ok p => Ok(coefficient = p.coefficient, scale = p.scale)
`
	dir := writeLSPDir(t, map[string]string{"m.can": used})
	out := t.TempDir()
	if err := compile(out, []string{dir + "/m.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "function $canDecParts") {
		t.Fatalf("emit missing $canDecParts helper:\n%s", raw)
	}
	if !strings.Contains(string(raw), "$canDecParts(v)") {
		t.Fatalf("emit missing helper call:\n%s", raw)
	}

	plain := used
	plain = strings.Replace(plain, "match call dec__parts(v)\n    on Ok p => Ok(coefficient = p.coefficient, scale = p.scale)",
		"Ok(coefficient = 1, scale = 0)", 1)
	plain = strings.Replace(plain, "t(v = d\"1.0\") => Ok(coefficient = 10, scale = 1)",
		"t(v = d\"1.0\") => Ok(coefficient = 1, scale = 0)", 1)
	dir2 := writeLSPDir(t, map[string]string{"m.can": plain})
	out2 := t.TempDir()
	if err := compile(out2, []string{dir2 + "/m.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw2, err := os.ReadFile(filepath.Join(out2, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw2), "$canDecParts") {
		t.Fatalf("emit leaks $canDecParts when unused:\n%s", raw2)
	}
}
