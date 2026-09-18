package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Slice 3: integer range patterns. Probes first: singletons,
// const bounds, cut-point coverage, shadowing (CAN4111), bad
// bounds (CAN4110), mixed slots, witnesses, emit lowering,
// admission stability, and the CAN4107-per-arm obligation.

const rangeLib = `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

`

// TestRangeSingleton pins integer singletons end to end.
func TestRangeSingleton(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    five(5) => Ok(value = 5)
    other(7) => Ok(value = 7)
  match x
    5 => Ok(value = 5)
    _ => Ok(value = 7)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("int singleton reported: %v", diags)
	}
}

// TestRangeConstBounds pins visible integer constants as
// singleton and range bounds.
func TestRangeConstBounds(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out, m__LO, m__HI, m__N]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

const m__LO: int rev 1 = 48

const m__HI: int rev 1 = 57

const m__N: int rev 1 = 58

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    digit(50) => Ok(value = 1)
    colon(58) => Ok(value = 2)
    other(99) => Ok(value = 3)
  match x
    m__LO..m__HI => Ok(value = 1)
    m__N => Ok(value = 2)
    _ => Ok(value = 3)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("const bounds reported: %v", diags)
	}
}

// TestRangeCutPoints pins partial overlap under first-match
// semantics: each arm keeps its exclusive atoms, so no arm is
// shadowed and every arm needs its own row.
func TestRangeCutPoints(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    low(7) => Ok(value = 1)
    high(12) => Ok(value = 2)
    out(99) => Ok(value = 3)
  match x
    1..10 => Ok(value = 1)
    5..15 => Ok(value = 2)
    _ => Ok(value = 3)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("overlapping ranges reported: %v", diags)
	}
}

// TestRangeShadowed pins CAN4111: an arm whose region earlier
// arms fully cover is a static error, not a coverage wait.
func TestRangeShadowed(t *testing.T) {
	for _, arm := range []string{"3..5", "5"} {
		src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    low(2) => Ok(value = 1)
    out(99) => Ok(value = 3)
  match x
    1..10 => Ok(value = 1)
    ` + arm + ` => Ok(value = 2)
    _ => Ok(value = 3)
`
		dir := writeLSPDir(t, map[string]string{"m.can": src})
		if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN4111") {
			t.Fatalf("%q reported no CAN4111: %v", arm, diags)
		}
	}
}

// TestRangeBadBounds pins CAN4110: equal, inverted, and
// non-integer bounds are rejected with the malformed-range code.
func TestRangeBadBounds(t *testing.T) {
	for _, arm := range []string{"5..5", "10..5", "1.5..5", "x..5", "5..", "m__S..5"} {
		src := rangeLib + `const m__S: str rev 1 = "s"

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    one(1) => Ok(value = 1)
  match x
    ` + arm + ` => Ok(value = 1)
    _ => Ok(value = 2)
`
		dir := writeLSPDir(t, map[string]string{"m.can": src})
		if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN4110") {
			t.Fatalf("%q reported no CAN4110: %v", arm, diags)
		}
	}
}

// TestRangeUnknownBound pins CAN2104 for names that resolve
// nowhere: a bound that is not a visible constant is unknown,
// not merely malformed.
func TestRangeUnknownBound(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    one(1) => Ok(value = 1)
  match x
    m__NOPE..5 => Ok(value = 1)
    _ => Ok(value = 2)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN2104") {
		t.Fatalf("unknown bound reported no CAN2104: %v", diags)
	}
}

// TestRangeMixedSlot pins the one-domain rule: integer patterns
// beside bool or string patterns in one slot are rejected.
func TestRangeMixedSlot(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    one(1) => Ok(value = 1)
  match x
    true => Ok(value = 1)
    5 => Ok(value = 5)
    _ => Ok(value = 0)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN4103") {
		t.Fatalf("mixed slot reported no CAN4103: %v", diags)
	}
}

// TestRangeUncovered pins exhaustiveness with concrete int
// witnesses: the open tails render as real counterexamples.
func TestRangeUncovered(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    mid(5) => Ok(value = 5)
  match x
    1..10 => Ok(value = x)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if !hasCode(diags, "CAN4101") {
		t.Fatalf("uncovered tails reported no CAN4101: %v", diags)
	}
	if !hasFound(diags, "(0)") || !hasFound(diags, "(11)") {
		t.Fatalf("expected concrete tail witnesses (0) and (11), got %v", diags)
	}
}

// TestRangeCallArm pins that call arms still take only error
// kinds and Ok: integer patterns there are CAN4105.
func TestRangeCallArm(t *testing.T) {
	src := `mod m
  provides [m__go, m__work, M__Out, M__Work]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

type M__Work rev 1 (
  value: int
)

fn m__work(v: int) -> M__Work rev 1
  emits []
  tests
    w(1) => Ok(value = 1)
  Ok(value = v)

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    one(1) => Ok(value = 1)
  match call m__work(x)
    5 => Ok(value = 5)
    on Ok r => Ok(value = r.value)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN4105") {
		t.Fatalf("int call-arm pattern reported no CAN4105: %v", diags)
	}
}

// TestRangeArmObligation pins one arm one CAN4107: an untaken
// range arm fails execution law even where the proof is clean.
func TestRangeArmObligation(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    low(5) => Ok(value = 1)
    out(99) => Ok(value = 3)
  match x
    1..10 => Ok(value = 1)
    11..20 => Ok(value = 2)
    _ => Ok(value = 3)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN4107") {
		t.Fatalf("untaken range arm reported no CAN4107: %v", diags)
	}
}

// TestRangeForeignConstPin pins the uses rule for bounds: a
// foreign integer constant as a bound needs its rev pin.
func TestRangeForeignConstPin(t *testing.T) {
	lib := `mod lib
  provides [lib__LO, lib__HI]
  uses []
  emits []

const lib__LO: int rev 1 = 48

const lib__HI: int rev 1 = 57
`
	pinned := `mod m
  provides [m__go, M__Out]
  uses [lib__LO@1, lib__HI@1]
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    digit(50) => Ok(value = 1)
    out(99) => Ok(value = 2)
  match x
    lib__LO..lib__HI => Ok(value = 1)
    _ => Ok(value = 2)
`
	dir := writeLSPDir(t, map[string]string{"lib.can": lib, "app.can": pinned})
	if diags := diagnose(dir, "app.can", pinned); hasError(diags) {
		t.Fatalf("pinned foreign bounds reported: %v", diags)
	}
	unpinned := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    digit(50) => Ok(value = 1)
    out(99) => Ok(value = 2)
  match x
    lib__LO..lib__HI => Ok(value = 1)
    _ => Ok(value = 2)
`
	dir2 := writeLSPDir(t, map[string]string{"lib.can": lib, "app.can": unpinned})
	if diags := diagnose(dir2, "app.can", unpinned); !hasCode(diags, "CAN2105") {
		t.Fatalf("unpinned foreign bounds reported no CAN2105: %v", diags)
	}
}

// TestRangeEmit pins the lowering shape: singletons compare
// bigint-exact and ranges lower to closed comparisons.
func TestRangeEmit(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    five(5) => Ok(value = 5)
    digit(50) => Ok(value = 50)
    out(99) => Ok(value = 99)
  match x
    5 => Ok(value = 5)
    48..57 => Ok(value = x)
    _ => Ok(value = 99)
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.can"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	ts := string(raw)
	for _, want := range []string{"=== 5n", ">= 48n", "<= 57n"} {
		if !strings.Contains(ts, want) {
			t.Fatalf("emitted TS lacks %q:\n%s", want, ts)
		}
	}
}

// TestRangeAdmissionContract pins the admission boundary: a
// contracted function with range arms in its body is admitted
// (no shape rejection), while symbolic verification still fails
// closed as inconclusive on the untranslatable body — the
// prover does not move in this slice.
func TestRangeAdmissionContract(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  requires
    x >= 0
  ensures
    on Ok result
      result.value >= 0
  tests
    low(5) => Ok(value = 1)
    out(99) => Ok(value = 2)
  match x
    1..10 => Ok(value = 1)
    _ => Ok(value = 2)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if hasCode(diags, "CAN4302") {
		t.Fatalf("admitted int body reported unsupported: %v", diags)
	}
	// Integer case analysis discharges: range and wildcard arms
	// prove, so the old fail-closed inconclusive is gone. (Point
	// 9: admission no longer outruns the prover here.)
	if hasCode(diags, "CAN4305") {
		t.Fatalf("int body went inconclusive instead of verifying: %v", diags)
	}
	if len(diags) != 0 {
		t.Fatalf("int body did not verify clean: %v", diags)
	}
}
