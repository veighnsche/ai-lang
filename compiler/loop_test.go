package main

import (
	"os"
	"strings"
	"testing"
)

// Termination + iteration: a direct self-call is a loop only with a
// decreases line proving an int param steps down by one under the
// positive branch on every site (a11). Without the proof the shape is
// a cycle error, refused before anything runs; recursion across files
// is refused program-wide by the same rule.

const loopPoll = `mod m
  provides [m__poll, M__S]
  uses []
  emits []

type M__S rev 1 (
  n: int
)

fn m__poll(n: int) -> M__S rev 1
  decreases n
  emits []
  tests
    now(n = 0) => Ok(n = 0)
    later(n = 2) => Ok(n = 0)
=
  match n <= 0
    true => Ok(n = 0)
    false => match call m__poll(n - 1)
      on Ok s => Ok(n = s.n)
`

func TestLoopClean(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.ail": loopPoll})
	if diags := diagnose(dir, "m.ail", loopPoll); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestLoopNamedDecreaseClean(t *testing.T) {
	named := strings.Replace(loopPoll,
		"match call m__poll(n - 1)", "match call m__poll(n = n - 1)", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": named})
	if diags := diagnose(dir, "m.ail", named); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestLoopBadName(t *testing.T) {
	bad := strings.Replace(loopPoll, "decreases n", "decreases nosuch", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": bad})
	diags := diagnose(dir, "m.ail", bad)
	if !hasDiag(diags, "error", "decreases nosuch: no such param") {
		t.Fatalf("expected bad-decreases error, got %v", diags)
	}
}

func TestLoopNonInt(t *testing.T) {
	bad := strings.Replace(loopPoll,
		"fn m__poll(n: int) -> M__S rev 1\n  decreases n",
		"fn m__poll(n: str) -> M__S rev 1\n  decreases n", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": bad})
	diags := diagnose(dir, "m.ail", bad)
	if !hasDiag(diags, "error", "must be an int param, got str") {
		t.Fatalf("expected non-int decreases error, got %v", diags)
	}
}

func TestLoopStale(t *testing.T) {
	stale := strings.Replace(loopPoll,
		"    false => match call m__poll(n - 1)\n      on Ok s => Ok(n = s.n)",
		"    false => Ok(n = 0)", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": stale})
	diags := diagnose(dir, "m.ail", stale)
	if !hasDiag(diags, "error", "decreases n but never calls itself") {
		t.Fatalf("expected stale-decreases error, got %v", diags)
	}
}

func TestLoopNoDecreaseShapes(t *testing.T) {
	cases := map[string]string{
		"unchanged": "match call m__poll(n)",
		"up":        "match call m__poll(n + 1)",
		"zero":      "match call m__poll(n - 0)",
		"computed":  "match call m__poll(n - n)",
		"literal":   "match call m__poll(1)",
	}
	for name, site := range cases {
		bad := strings.Replace(loopPoll, "match call m__poll(n - 1)", site, 1)
		dir := writeLSPDir(t, map[string]string{"m.ail": bad})
		diags := diagnose(dir, "m.ail", bad)
		if !hasDiag(diags, "error", "without decreasing n by one: pass n - 1") {
			t.Fatalf("%s: expected no-decrease error, got %v", name, diags)
		}
	}
}

// Larger steps terminate but break the one-spelling rule: only the
// unit step is admitted (a11).
func TestLoopLargerStepRefused(t *testing.T) {
	for _, site := range []string{"match call m__poll(n - 2)", "match call m__poll(n = n - 3)"} {
		bad := strings.Replace(loopPoll, "match call m__poll(n - 1)", site, 1)
		dir := writeLSPDir(t, map[string]string{"m.ail": bad})
		diags := diagnose(dir, "m.ail", bad)
		if !hasDiag(diags, "error", "without decreasing n by one: pass n - 1") {
			t.Fatalf("%s: expected unit-step error, got %v", site, diags)
		}
	}
}

// Unguarded self-calls are refused even with a unit step: a site in
// the true arm diverges (each entry steps further negative), and a
// site with no bound guard at all proves nothing (a11, AIL3009).
func TestLoopUnguardedRefused(t *testing.T) {
	trueArm := strings.Replace(loopPoll,
		"    true => Ok(n = 0)\n    false => match call m__poll(n - 1)\n      on Ok s => Ok(n = s.n)",
		"    true => match call m__poll(n - 1)\n      on Ok s => Ok(n = s.n)\n    false => Ok(n = 0)", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": trueArm})
	diags := diagnose(dir, "m.ail", trueArm)
	if !hasDiag(diags, "error", "outside the positive branch") {
		t.Fatalf("expected unguarded-recursion error, got %v", diags)
	}
	bare := strings.Replace(loopPoll,
		`  match n <= 0
    true => Ok(n = 0)
    false => match call m__poll(n - 1)
      on Ok s => Ok(n = s.n)`,
		`  match call m__poll(n - 1)
    on Ok s => Ok(n = s.n)`, 1)
	dir = writeLSPDir(t, map[string]string{"m.ail": bare})
	diags = diagnose(dir, "m.ail", bare)
	if !hasDiag(diags, "error", "outside the positive branch") {
		t.Fatalf("expected unguarded-recursion error, got %v", diags)
	}
}

// Strict and < 1 guards cannot prove termination: the rule admits
// only the p <= 0 spelling, so a self-call under any other guard is
// refused even with a unit step (issue #40 review).
func TestLoopStrictGuardRefused(t *testing.T) {
	for _, guard := range []string{"match n > 0", "match n < 1"} {
		bad := strings.Replace(loopPoll, "match n <= 0", guard, 1)
		dir := writeLSPDir(t, map[string]string{"m.ail": bad})
		diags := diagnose(dir, "m.ail", bad)
		if !hasDiag(diags, "error", "outside the positive branch") {
			t.Fatalf("%s: expected guard-shape error, got %v", guard, diags)
		}
	}
}

// Blessed schemas (a19, issue #41): euclid and narrowing admit
// efficient recursion the unit step cannot spell. Each pairs one
// canonical step with one canonical guard; anything else is refused
// with the same code families as the unit loop.
const loopEuclid = `mod m
  provides [m__gcd, M__S]
  uses []
  emits []

type M__S rev 1 (
  n: int
)

fn m__gcd(a: int, b: int) -> M__S rev 1
  decreases a, b by euclid
  emits []
  tests
    basic(a = 12, b = 8) => Ok(n = 4)
    coprime(a = 8, b = 9) => Ok(n = 1)
    zero_b(a = 5, b = 0) => Ok(n = 5)
=
  match b <= 0
    true => Ok(n = a)
    false => match call m__gcd(b, a % b)
      on Ok r => Ok(n = r.n)
`

func TestLoopEuclidClean(t *testing.T) {
	for _, line := range []string{"decreases a, b by euclid", "decreases a,b by euclid", "decreases a , b by euclid"} {
		variant := strings.Replace(loopEuclid, "decreases a, b by euclid", line, 1)
		dir := writeLSPDir(t, map[string]string{"m.ail": variant})
		if diags := diagnose(dir, "m.ail", variant); len(diags) != 0 {
			t.Fatalf("%s: expected no diagnostics, got %v", line, diags)
		}
	}
}

func TestLoopEuclidBadStep(t *testing.T) {
	for _, site := range []string{"match call m__gcd(a - 1, b)", "match call m__gcd(b, a - b)"} {
		bad := strings.Replace(loopEuclid, "match call m__gcd(b, a % b)", site, 1)
		dir := writeLSPDir(t, map[string]string{"m.ail": bad})
		diags := diagnose(dir, "m.ail", bad)
		if !hasDiag(diags, "error", "without a euclid step") {
			t.Fatalf("%s: expected euclid-step error, got %v", site, diags)
		}
	}
}

func TestLoopEuclidUnguarded(t *testing.T) {
	trueArm := strings.Replace(loopEuclid,
		"    true => Ok(n = a)\n    false => match call m__gcd(b, a % b)\n      on Ok r => Ok(n = r.n)",
		"    true => match call m__gcd(b, a % b)\n      on Ok r => Ok(n = r.n)\n    false => Ok(n = a)", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": trueArm})
	diags := diagnose(dir, "m.ail", trueArm)
	if !hasDiag(diags, "error", "must sit under the false arm of b <= 0") {
		t.Fatalf("expected euclid-guard error, got %v", diags)
	}
}

const loopNarrow = `mod m
  provides [m__sqrt, M__S]
  uses []
  emits []

type M__S rev 1 (
  n: int
)

fn m__sqrt(value: int, lo: int, hi: int) -> M__S rev 1
  decreases lo, hi by narrowing
  emits []
  tests
    ten(value = 10, lo = 0, hi = 11) => Ok(n = 3)
    exact(value = 9, lo = 3, hi = 4) => Ok(n = 3)
    nine(value = 9, lo = 0, hi = 10) => Ok(n = 3)
=
  match (hi - lo) <= 1
    true => Ok(n = lo)
    false => match ((lo + hi) / 2) * ((lo + hi) / 2) <= value
      true => match call m__sqrt(value, (lo + hi) / 2, hi)
        on Ok r => Ok(n = r.n)
      false => match call m__sqrt(value, lo, (lo + hi) / 2)
        on Ok r => Ok(n = r.n)
`

func TestLoopNarrowingClean(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.ail": loopNarrow})
	if diags := diagnose(dir, "m.ail", loopNarrow); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestLoopNarrowingBadStep(t *testing.T) {
	bad := strings.Replace(loopNarrow,
		"match call m__sqrt(value, (lo + hi) / 2, hi)",
		"match call m__sqrt(value, lo, hi - 1)", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": bad})
	diags := diagnose(dir, "m.ail", bad)
	if !hasDiag(diags, "error", "without a narrowing step") {
		t.Fatalf("expected narrowing-step error, got %v", diags)
	}
}

func TestLoopNarrowingBadGuard(t *testing.T) {
	bad := strings.Replace(loopNarrow, "match (hi - lo) <= 1", "match hi <= lo", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": bad})
	diags := diagnose(dir, "m.ail", bad)
	if !hasDiag(diags, "error", "must sit under the false arm of (hi - lo) <= 1") {
		t.Fatalf("expected narrowing-guard error, got %v", diags)
	}
}

func TestLoopBadSchemaLine(t *testing.T) {
	for _, line := range []string{"decreases a, b by half", "decreases a by euclid", "decreases a, b, c by euclid"} {
		bad := strings.Replace(loopEuclid, "decreases a, b by euclid", line, 1)
		dir := writeLSPDir(t, map[string]string{"m.ail": bad})
		diags := diagnose(dir, "m.ail", bad)
		if !hasDiag(diags, "error", "bad decreases line") {
			t.Fatalf("%s: expected bad-line error, got %v", line, diags)
		}
	}
}

func TestLoopSchemaNonInt(t *testing.T) {
	bad := strings.Replace(loopEuclid, "fn m__gcd(a: int, b: int)", "fn m__gcd(a: int, b: str)", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": bad})
	diags := diagnose(dir, "m.ail", bad)
	if !hasDiag(diags, "error", "must be an int param") {
		t.Fatalf("expected int-param error, got %v", diags)
	}
}

func TestLoopNamedNoDecrease(t *testing.T) {
	bad := strings.Replace(loopPoll,
		"match call m__poll(n - 1)", "match call m__poll(n = n)", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": bad})
	diags := diagnose(dir, "m.ail", bad)
	if !hasDiag(diags, "error", "without decreasing n by one: pass n - 1") {
		t.Fatalf("expected no-decrease error, got %v", diags)
	}
}

func TestLoopMutualStaysBanned(t *testing.T) {
	pair := `mod m
  provides [m__a, m__b, M__S]
  uses []
  emits []

type M__S rev 1 (
  n: int
)

fn m__a(n: int) -> M__S rev 1
  decreases n
  emits []
  tests
    a0(n = 0) => Ok(n = 0)
=
  match n <= 0
    true => Ok(n = 0)
    false => match call m__b(n - 1)
      on Ok s => Ok(n = s.n)

fn m__b(n: int) -> M__S rev 1
  decreases n
  emits []
  tests
    b0(n = 0) => Ok(n = 0)
=
  match n <= 0
    true => Ok(n = 0)
    false => match call m__a(n - 1)
      on Ok s => Ok(n = s.n)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": pair})
	diags := diagnose(dir, "m.ail", pair)
	if !hasDiag(diags, "error", "local call cycle") {
		t.Fatalf("expected mutual-cycle error, got %v", diags)
	}
}

// Negative entries take the base arm and return a declared outcome
// (a11): the guard rule admits recursion only under the positive
// branch, so fault-bounded evaluation is no longer the theorem.
func TestLoopNegativeEntryTakesBase(t *testing.T) {
	neg := strings.Replace(loopPoll,
		"    later(n = 2) => Ok(n = 0)",
		"    later(n = 2) => Ok(n = 0)\n    neg(n = -1) => Ok(n = 0)", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": neg})
	if diags := diagnose(dir, "m.ail", neg); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestLoopDupDecreases(t *testing.T) {
	dup := strings.Replace(loopPoll,
		"  decreases n", "  decreases n\n  decreases n", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": dup})
	diags := diagnose(dir, "m.ail", dup)
	if !hasDiag(diags, "error", "duplicate decreases line") {
		t.Fatalf("expected duplicate-decreases parse error, got %v", diags)
	}
}

// Cross-file recursion is refused program-wide (a11): the sandbox
// stubs foreign calls, so without a whole-program check a mutual
// cycle across two files passes every test, then links into an
// unproved recursive cycle in production. The gate blocks test
// execution entirely (no test failures, only the cycle error).
func TestLoopCrossFileCycleRefused(t *testing.T) {
	alpha := `mod alpha
  provides [alpha__run, Alpha__Value]
  uses [beta__run@1, Beta__Value@1]
  emits []

type Alpha__Value rev 1 (
  value: int
)

fn alpha__run(n: int) -> Alpha__Value rev 1
  emits []
  tests
    sample(n = 0) => Ok(value = 0)
=
  match call beta__run(n)
    given
      sample => [exchange args (n = 0) outcome Ok(value = 0)]
    on Ok result => Ok(value = result.value)
`
	beta := `mod beta
  provides [beta__run, Beta__Value]
  uses [alpha__run@1, Alpha__Value@1]
  emits []

type Beta__Value rev 1 (
  value: int
)

fn beta__run(n: int) -> Beta__Value rev 1
  emits []
  tests
    sample(n = 0) => Ok(value = 0)
=
  match call alpha__run(n)
    given
      sample => [exchange args (n = 0) outcome Ok(value = 0)]
    on Ok result => Ok(value = result.value)
`
	files := map[string]string{"alpha.ail": alpha, "beta.ail": beta}
	for _, open := range []string{"alpha.ail", "beta.ail"} {
		dir := writeLSPDir(t, files)
		diags := diagnose(dir, open, files[open])
		if !hasDiag(diags, "error", "call cycle alpha__run -> beta__run -> alpha__run") {
			t.Fatalf("%s: expected cross-file cycle error, got %v", open, diags)
		}
		if hasDiag(diags, "error", "fails:") {
			t.Fatalf("%s: gate must block tests before running, got %v", open, diags)
		}
		var attributed bool
		for _, d := range diags {
			if d.Code == CodeLocalCycle && d.File == "beta.ail" {
				attributed = true
			}
		}
		if !attributed {
			t.Fatalf("%s: cycle must attribute to the closing call site in beta.ail, got %v", open, diags)
		}
	}
}

// The shipped sketch minus its proof is refused before anything runs:
// minus the decreases line the same file is a cycle error, so the
// gate blocks test execution entirely (no test failures, only the
// cycle error). The refusal is proof-gated admission: this finite
// fuel shape terminates either way, and the error says exactly that
// the proof is missing — not that the build would hang.
func TestLoopSketchMinusDecreasesRefused(t *testing.T) {
	raw, err := os.ReadFile("../sketches/retry-loop/retry.ail")
	if err != nil {
		t.Fatal(err)
	}
	stripped := strings.Replace(string(raw), "  decreases fuel\n", "", 1)
	if stripped == string(raw) {
		t.Fatal("sketch has no decreases line to strip")
	}
	dir := writeLSPDir(t, map[string]string{"retry.ail": stripped})
	diags := diagnose(dir, "retry.ail", stripped)
	if !hasDiag(diags, "error", "local call cycle retry__fetch -> retry__fetch") {
		t.Fatalf("expected cycle error, got %v", diags)
	}
	if hasDiag(diags, "error", "fails:") {
		t.Fatalf("gate must block tests before running, got %v", diags)
	}
}
