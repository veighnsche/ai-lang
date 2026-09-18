package main

import (
	"strings"
	"testing"
)

// Decimal semantics live twice: decArith/decCmp/canonDec evaluate
// at compile time (Go, exact via big integers), the $canDec*
// helpers compute at runtime (TypeScript, exact via bigints).
// The two cannot merge across the language boundary, so this
// file pins their contract instead: shared agreement vectors
// plus op-set parity. A new decimal op must land on both sides
// or these tests fail. (TS twins: $canDecAdd/Sub/Mul/Ge/Le/Gt/Lt
// with $canDecSplit/Norm/Mant/FromMant, emit.go ~773-885;
// ==/!= stay native ===/!== on canonical-digit strings.)

// TestDecArithParityVectors pins the Go side of the agreement
// battery: canonical results the TS helpers must reproduce.
func TestDecArithParityVectors(t *testing.T) {
	for _, tc := range []struct{ op, l, r, want string }{
		{"+", "0.1", "0.2", "0.3"},
		{"+", "1.0", "2.0", "3.0"},
		{"-", "1.0", "0.9", "0.1"},
		{"-", "0.0", "0.0", "0.0"},
		{"*", "2.5", "4.0", "10.0"},
		{"*", "0.1", "0.1", "0.01"},
		{"*", "-2.0", "3.0", "-6.0"},
		{"+", "99999999999999999999.9", "0.1", "100000000000000000000.0"},
	} {
		got, err := decArith(tc.op, tc.l, tc.r)
		if err != nil {
			t.Errorf("decArith(%q, %q, %q) errored: %v", tc.op, tc.l, tc.r, err)
			continue
		}
		if got != tc.want {
			t.Errorf("decArith(%q, %q, %q) = %q, want %q", tc.op, tc.l, tc.r, got, tc.want)
		}
	}
	// Division is refused on both sides: decArith errors, emit
	// has no $canDecDiv helper (cannot emit op).
	if _, err := decArith("/", "1.0", "3.0"); err == nil {
		t.Errorf("decArith(/) succeeded: non-terminating division must refuse")
	}
}

// TestDecCmpParityVectors pins comparison semantics, including
// the default (==) and != arms the TS side lowers natively.
func TestDecCmpParityVectors(t *testing.T) {
	lr, _ := decRat("1.5")
	rr, _ := decRat("1.50")
	if lr.Cmp(rr) != 0 {
		t.Fatalf("fixture not equal: 1.5 vs 1.50")
	}
	for _, tc := range []struct {
		op   string
		c    int
		want bool
	}{
		{">=", 1, true}, {">=", 0, true}, {">=", -1, false},
		{"<=", -1, true}, {"<=", 0, true}, {"<=", 1, false},
		{">", 1, true}, {">", 0, false},
		{"<", -1, true}, {"<", 0, false},
		{"!=", 1, true}, {"!=", 0, false},
		{"==", 0, true}, {"==", 1, false},
	} {
		if got := decCmp(tc.op, tc.c); got != tc.want {
			t.Errorf("decCmp(%q, %d) = %v, want %v", tc.op, tc.c, got, tc.want)
		}
	}
}

// TestDecHelperParity guards one-sided extension: every decimal
// op the evaluator models must lower through its TS helper (or
// the documented native equality), and no probe may emit a
// helper the evaluator cannot model.
func TestDecHelperParity(t *testing.T) {
	ts := compileEmit(t, decParityProbe)
	for _, want := range []string{
		"$canDecAdd(", "$canDecSub(", "$canDecMul(",
		"$canDecGe(", "$canDecLe(", "$canDecGt(", "$canDecLt(",
	} {
		if !strings.Contains(ts, want) {
			t.Errorf("probe emit missing %q:\n%s", want, ts)
		}
	}
	// ==/!= on canonical digits are native; they must not grow
	// helpers behind the evaluator's back.
	for _, banned := range []string{"$canDecEq(", "$canDecNe("} {
		if strings.Contains(ts, banned) {
			t.Errorf("probe emit contains unexpected %q:\n%s", banned, ts)
		}
	}
}

const decParityProbe = `mod m
  provides [m__arith, m__ge, m__le, m__gt, m__lt, m__eq, m__ne, M__D, M__B]
  uses []
  emits []

type M__D rev 1 (
  value: dec
)

type M__B rev 1 (
  value: bool
)

fn m__arith(a: dec, b: dec) -> M__D rev 1
  emits []
  tests
    go(d"0.1", d"0.2") => Ok(d"0.1")
  Ok(a + b - b * d"1.0")

fn m__ge(a: dec, b: dec) -> M__B rev 1
  emits []
  tests
    hi(d"0.5", d"0.2") => Ok(true)
    lo(d"0.2", d"0.5") => Ok(false)
  match a >= b
    true => Ok(true)
    false => Ok(false)

fn m__le(a: dec, b: dec) -> M__B rev 1
  emits []
  tests
    lo(d"0.2", d"0.5") => Ok(true)
    hi(d"0.5", d"0.2") => Ok(false)
  match a <= b
    true => Ok(true)
    false => Ok(false)

fn m__gt(a: dec, b: dec) -> M__B rev 1
  emits []
  tests
    hi(d"0.5", d"0.2") => Ok(true)
    lo(d"0.2", d"0.5") => Ok(false)
  match a > b
    true => Ok(true)
    false => Ok(false)

fn m__lt(a: dec, b: dec) -> M__B rev 1
  emits []
  tests
    lo(d"0.2", d"0.5") => Ok(true)
    hi(d"0.5", d"0.2") => Ok(false)
  match a < b
    true => Ok(true)
    false => Ok(false)

fn m__eq(a: dec, b: dec) -> M__B rev 1
  emits []
  tests
    same(d"1.5", d"1.50") => Ok(true)
    diff(d"1.5", d"2.5") => Ok(false)
  match a == b
    true => Ok(true)
    false => Ok(false)

fn m__ne(a: dec, b: dec) -> M__B rev 1
  emits []
  tests
    diff(d"1.5", d"2.5") => Ok(true)
    same(d"1.5", d"1.50") => Ok(false)
  match a != b
    true => Ok(true)
    false => Ok(false)
`
