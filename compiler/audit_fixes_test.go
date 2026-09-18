package main

import (
	"errors"
	"math/big"
	"strings"
	"testing"
)

// Subagent-audit fixes: one focused regression test per confirmed
// finding. Each test fails on the pre-fix behavior (or panics) and
// passes after.

// A mixed call binds each argument by its own form (bindSlots), not
// by the whole call's first argument: m__pick(n, c = 3, b = 2)
// passes b = 2 against a requires b == 3, so the call-site
// precondition must NOT verify. Judging the call from Args[0]
// alone misbound c's value into b and proved it.
func TestAuditMixedCallPreconditionCaught(t *testing.T) {
	src := `mod m
  provides [m__pick, m__caller, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__pick(a: int, b: int, c: int) -> M__Out rev 1
  emits []
  requires
    b == 3
  ensures
    on Ok result
      result.value == a + b + c
  tests
    v(1, 3, 5) => Ok(9)
  Ok(a + b + c)

fn m__caller(n: int) -> M__Out rev 1
  emits []
  requires
    n == 1
  ensures
    on Ok result
      result.value == 6
  tests
    go(1) => Ok(6)
  match call m__pick(n, c = 3, b = 2)
    on Ok r => Ok(r.value)
`
	prog, texts := admitProg(t, src)
	diags := VerifyContracts(prog, texts)
	if !hasFound(diags, "requires of m__pick unproven") {
		t.Fatalf("violated mixed-call precondition verified: %v", diags)
	}
}

// The other direction: a satisfied mixed call verifies clean. The
// old Args[0] rule bound c's value into b here and spuriously
// refuted the precondition.
func TestAuditMixedCallSatisfiedVerifies(t *testing.T) {
	src := `mod m
  provides [m__pick, m__caller, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__pick(a: int, b: int, c: int) -> M__Out rev 1
  emits []
  requires
    b == 3
  ensures
    on Ok result
      result.value == a + b + c
  tests
    v(1, 3, 5) => Ok(9)
  Ok(a + b + c)

fn m__caller(n: int) -> M__Out rev 1
  emits []
  requires
    n == 1
  ensures
    on Ok result
      result.value == 9
  tests
    go(1) => Ok(9)
  match call m__pick(n, c = 5, b = 3)
    on Ok r => Ok(r.value)
`
	prog, texts := admitProg(t, src)
	if diags := VerifyContracts(prog, texts); len(diags) != 0 {
		t.Fatalf("satisfied mixed call did not verify: %v", diags)
	}
}

// Raw string patterns validate like raw literals: a lone quote
// panicked the slice, an unterminated literal silently lost its
// last byte.
func TestAuditRawPatternValidation(t *testing.T) {
	for _, bad := range []string{`"`, `"abc`, `"a" extra`} {
		if _, err := parsePattern(bad); err == nil {
			t.Fatalf("parsePattern(%q): expected error, got none", bad)
		}
	}
	p, err := parsePattern(`"abc"`)
	if err != nil {
		t.Fatalf("parsePattern: %v", err)
	}
	if p.Kind != "str" || p.Str != "abc" {
		t.Fatalf("parsePattern(%q) = %+v", `"abc"`, p)
	}
	p, err = parsePattern(`""`)
	if err != nil {
		t.Fatalf("parsePattern empty: %v", err)
	}
	if p.Kind != "str" || p.Str != "" {
		t.Fatalf("parsePattern empty = %+v", p)
	}
}

// Sealed values canon from their operand: seal B("x") and seal
// B("y") must hash differently, or the drift is silent.
func TestAuditSealCanonDistinguishes(t *testing.T) {
	x, err := parseSmall(`seal B("x")`)
	if err != nil {
		t.Fatalf("parseSmall seal: %v", err)
	}
	y, err := parseSmall(`seal B("y")`)
	if err != nil {
		t.Fatalf("parseSmall seal: %v", err)
	}
	if canonSmall(x) == canonSmall(y) {
		t.Fatalf("seal canon collides: %q", canonSmall(x))
	}
	again, err := parseSmall(`seal B("x")`)
	if err != nil {
		t.Fatalf("parseSmall seal: %v", err)
	}
	if canonSmall(x) != canonSmall(again) {
		t.Fatalf("seal canon unstable: %q vs %q", canonSmall(x), canonSmall(again))
	}
}

// not/neg/forward canon with their operands: distinct operands
// must hash distinctly, never collapse to one unknown-kind.
func TestAuditUnaryCanonDistinguishes(t *testing.T) {
	nx, err := parseSmall("not a")
	if err != nil {
		t.Fatalf("parseSmall not: %v", err)
	}
	ny, err := parseSmall("not b")
	if err != nil {
		t.Fatalf("parseSmall not: %v", err)
	}
	if canonSmall(nx) == canonSmall(ny) {
		t.Fatalf("not canon collides: %q", canonSmall(nx))
	}
	negx := &Small{Kind: "neg", L: &Small{Kind: "ref", Ref: []string{"a"}}}
	negy := &Small{Kind: "neg", L: &Small{Kind: "ref", Ref: []string{"b"}}}
	if canonSmall(negx) == canonSmall(negy) {
		t.Fatalf("neg canon collides: %q", canonSmall(negx))
	}
	fx := &Small{Kind: "forward", Str: "e"}
	fy := &Small{Kind: "forward", Str: "f"}
	if canonSmall(fx) == canonSmall(fy) {
		t.Fatalf("forward canon collides: %q", canonSmall(fx))
	}
	for _, c := range []string{canonSmall(nx), canonSmall(negx), canonSmall(fx)} {
		if strings.Contains(c, "unknown-kind") {
			t.Fatalf("canon fell through to unknown-kind: %q", c)
		}
	}
}

// Repeated parameter names fail at the declaration: the binding
// map keeps one slot while positional calls fill both.
func TestAuditDuplicateParamsRejected(t *testing.T) {
	fn := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(a: int, a: int) -> M__Out rev 1
  emits []
  tests
    v(1, 2) => Ok(1)
  Ok(a)
`
	if _, err := parseModuleText("m.can", fn); err == nil || !strings.Contains(err.Error(), "duplicate param a") {
		t.Fatalf("fn dup params: got %v", err)
	}
	ext := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

extern e__f(a: int, a: int) -> int rev 1
  emits []

fn m__go(a: int) -> M__Out rev 1
  emits []
  tests
    v(1) => Ok(1)
  Ok(a)
`
	if _, err := parseModuleText("m.can", ext); err == nil || !strings.Contains(err.Error(), "duplicate param a") {
		t.Fatalf("extern dup params: got %v", err)
	}
}

// != is the negation of == over the same structural equality:
// bools and records compare, exactly as == accepts.
func TestAuditNotEqualStructural(t *testing.T) {
	for _, c := range []struct {
		expr string
		want bool
	}{
		{`true != false`, true},
		{`true != true`, false},
		{`1 != 2`, true},
		{`2 != 2`, false},
	} {
		sm, err := parseSmall(c.expr)
		if err != nil {
			t.Fatalf("parseSmall(%q): %v", c.expr, err)
		}
		v, err := evSmall(sm, map[string]*Value{}, &Ctx{}, "test")
		if err != nil {
			t.Fatalf("evSmall(%q): %v", c.expr, err)
		}
		if v.Kind != "bool" || v.B != c.want {
			t.Fatalf("evSmall(%q) = %+v, want bool %v", c.expr, v, c.want)
		}
	}
	rec := func(n int64) *Value {
		return &Value{Kind: "rec", Dict: map[string]*Value{"v": {Kind: "int", N: big.NewInt(n)}}}
	}
	for _, c := range []struct {
		name string
		l, r *Value
		want bool
	}{
		{"differ", rec(1), rec(2), true},
		{"equal", rec(1), rec(1), false},
	} {
		sm := &Small{Kind: "binop", Op: "!=",
			L: &Small{Kind: "ref", Ref: []string{"l"}},
			R: &Small{Kind: "ref", Ref: []string{"r"}}}
		env := map[string]*Value{"l": c.l, "r": c.r}
		v, err := evSmall(sm, env, &Ctx{}, "test")
		if err != nil {
			t.Fatalf("record %s: %v", c.name, err)
		}
		if v.Kind != "bool" || v.B != c.want {
			t.Fatalf("record %s = %+v, want %v", c.name, v, c.want)
		}
	}
}

// Ordering needs an ordered domain: records and bools compare
// for equality only.
func TestAuditOrderingRefused(t *testing.T) {
	rec := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(a: M__Out, b: M__Out) -> M__Out rev 1
  emits []
  tests
    v(M__Out(1), M__Out(2)) => Ok(1)
  match a < b
    true => Ok(1)
    false => Ok(2)
`
	dir := writeLSPDir(t, map[string]string{"m.can": rec})
	diags := diagnose(dir, "m.can", rec)
	if !hasDiag(diags, "error", "cannot order M__Out with M__Out") {
		t.Fatalf("record ordering accepted: %v", diags)
	}
	boolean := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(a: bool, b: bool) -> M__Out rev 1
  emits []
  tests
    v(true, false) => Ok(1)
  match a > b
    true => Ok(1)
    false => Ok(2)
`
	dir = writeLSPDir(t, map[string]string{"m.can": boolean})
	diags = diagnose(dir, "m.can", boolean)
	if !hasDiag(diags, "error", "cannot order bool with bool") {
		t.Fatalf("bool ordering accepted: %v", diags)
	}
}

// The emitter fails ordering over unordered operand types like
// == does, instead of emitting a native object comparison.
func TestAuditEmitOrderingGate(t *testing.T) {
	lit := func(n int64) *Small { return &Small{Kind: "int", Num: big.NewInt(n)} }
	e := &emitter{}
	for _, ot := range []string{"M__Out", "bool"} {
		node := &Small{Kind: "binop", Op: "<", T: ot, L: lit(1), R: lit(2)}
		if _, err := e.emitValue(node); err == nil || !strings.Contains(err.Error(), CodeBadCompare) {
			t.Fatalf("emit <%s over %s: got %v", node.Op, ot, err)
		}
	}
	ok, err := e.emitValue(&Small{Kind: "binop", Op: "<", T: "int", L: lit(1), R: lit(2)})
	if err != nil || !strings.Contains(ok, "<") {
		t.Fatalf("emit int ordering: %q, %v", ok, err)
	}
}

// A trapping scrutinee is impure: / and % with a non-literal
// divisor match firstTrappingOp, so same-outcome never advises
// dropping a match whose scrutinee can fault.
func TestAuditLintPurityTrapping(t *testing.T) {
	lit := func(v int64) *Small { return &Small{Kind: "int", Num: big.NewInt(v)} }
	ref := &Small{Kind: "ref", Ref: []string{"d"}}
	cases := []struct {
		name string
		in   *Small
		want bool
	}{
		{"divRef", &Small{Kind: "binop", Op: "/", L: lit(7), R: ref}, false},
		{"modRef", &Small{Kind: "binop", Op: "%", L: lit(7), R: ref}, false},
		{"modLit", &Small{Kind: "binop", Op: "%", L: lit(7), R: lit(2)}, true},
		{"divZeroLit", &Small{Kind: "binop", Op: "/", L: lit(7), R: lit(0)}, false},
		{"add", &Small{Kind: "binop", Op: "+", L: lit(7), R: ref}, true},
	}
	for _, tc := range cases {
		if got := lintPureSmall(tc.in); got != tc.want {
			t.Fatalf("%s: pure = %v, want %v", tc.name, got, tc.want)
		}
	}
	trapping := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(n: int, d: int) -> M__Out rev 1
  emits []
  tests
    v(7, 3) => Ok(0)
    w(6, 3) => Ok(0)
  match n % d
    1 => Ok(0)
    _ => Ok(0)
`
	dir := writeLSPDir(t, map[string]string{"m.can": trapping})
	diags := diagnose(dir, "m.can", trapping)
	if hasCode(diags, CodeLintSameOutcome) {
		t.Fatalf("same-outcome fired over a trapping scrutinee: %v", diags)
	}
	pure := strings.Replace(trapping, "match n % d", "match n % 2", 1)
	pure = strings.Replace(pure, "v(7, 3)", "v(7, 2)", 1)
	pure = strings.Replace(pure, "w(6, 3)", "w(6, 2)", 1)
	dir = writeLSPDir(t, map[string]string{"m.can": pure})
	diags = diagnose(dir, "m.can", pure)
	if !hasCode(diags, CodeLintSameOutcome) {
		t.Fatalf("same-outcome silent over a pure scrutinee: %v", diags)
	}
}

// A uses entry with a malformed pin is a pin-format error, never
// a resolution error.
func TestAuditMalformedPinFormat(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses [db__get@x]
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(n: int) -> M__Out rev 1
  emits []
  tests
    v(1) => Ok(1)
  Ok(n)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if !hasCode(diags, CodeUsesPin) {
		t.Fatalf("malformed pin reported no %s: %v", CodeUsesPin, diags)
	}
	if hasCode(diags, CodeUsesResolve) {
		t.Fatalf("malformed pin misreported as %s: %v", CodeUsesResolve, diags)
	}
}

// A match with no arms reports the match line, not the indent width.
func TestAuditNoArmsLine(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(n: int) -> M__Out rev 1
  emits []
  tests
    v(1) => Ok(1)
  match n
`
	_, err := parseModuleText("m.can", src)
	if err == nil || !strings.Contains(err.Error(), "match with no arms") {
		t.Fatalf("expected no-arms error, got %v", err)
	}
	var le *LineError
	if !errors.As(err, &le) {
		t.Fatalf("expected a line error, got %T", err)
	}
	if want := expectLine(t, src, "match n"); le.Line != want {
		t.Fatalf("no-arms line = %d, want %d", le.Line, want)
	}
}

// Hand-built shapes fail closed, never panic: an int node without
// its big value, a call arg without its expression.
func TestAuditNilGuards(t *testing.T) {
	mid := &Small{Kind: "binop", Op: "/",
		L: &Small{Kind: "binop", Op: "+",
			L: &Small{Kind: "ref", Ref: []string{"lo"}},
			R: &Small{Kind: "ref", Ref: []string{"hi"}}},
		R: &Small{Kind: "int"}}
	if isMid(mid, "lo", "hi") {
		t.Fatal("isMid accepted a divisorless midpoint")
	}
	if name, ok := storeCellName(&Small{Args: []Arg{{V: nil}}}); ok || name != "" {
		t.Fatalf("storeCellName accepted a valueless arg: %q", name)
	}
}
