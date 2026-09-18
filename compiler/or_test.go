package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Slice 4: or-patterns. Probes first: basic union arms, string
// pipes, range/`|` precedence, per-alternative usefulness
// (AIL4112), mixed slots, const alternatives, emit lowering,
// and the one-obligation-per-source-arm law.

// TestOrBasic pins `|` alternatives selecting end to end.
func TestOrBasic(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    one(x = 1) => Ok(value = 1)
    two(x = 2) => Ok(value = 2)
    three(x = 3) => Ok(value = 3)
    other(x = 9) => Ok(value = 9)
=
  match x
    1 | 2 | 3 => Ok(value = x)
    _ => Ok(value = 9)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); hasError(diags) {
		t.Fatalf("or alternatives reported: %v", diags)
	}
}

// TestOrStringPipe pins that `|` inside a string literal is
// not an alternative separator.
func TestOrStringPipe(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: str
)

fn m__go(x: str) -> M__Out rev 1
  emits []
  tests
    piped(x = "a|b") => Ok(value = "a|b")
    plain(x = "c") => Ok(value = "c")
    other(x = "z") => Ok(value = "z")
=
  match x
    "a|b" | "c" => Ok(value = x)
    _ => Ok(value = "z")
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); hasError(diags) {
		t.Fatalf("string pipe reported: %v", diags)
	}
}

// TestOrRangePrecedence pins ranges binding tighter than `|`:
// `1..10 | 5..15` is two range alternatives, and the second
// keeps its exclusive region (11..15), so no arm is shadowed.
func TestOrRangePrecedence(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    low(x = 7) => Ok(value = 1)
    high(x = 12) => Ok(value = 1)
    out(x = 99) => Ok(value = 2)
=
  match x
    1..10 | 5..15 => Ok(value = 1)
    _ => Ok(value = 2)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); hasError(diags) {
		t.Fatalf("range alternatives reported: %v", diags)
	}
}

// TestOrShadowedAlt pins AIL4112: an alternative whose
// effective region is empty (covered by earlier arms and
// earlier alternatives) is rejected statically.
func TestOrShadowedAlt(t *testing.T) {
	for _, arm := range []string{"1..10 | 5..8", "1 | 1", "1 | 2 | 2"} {
		src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 1)
    out(x = 99) => Ok(value = 2)
=
  match x
    ` + arm + ` => Ok(value = 1)
    _ => Ok(value = 2)
`
		dir := writeLSPDir(t, map[string]string{"m.ail": src})
		if diags := diagnose(dir, "m.ail", src); !hasCode(diags, "AIL4112") {
			t.Fatalf("%q reported no AIL4112: %v", arm, diags)
		}
	}
}

// TestOrWildcardRejected pins AIL4112 for `_` inside
// alternative lists: the wildcard stays its own arm.
func TestOrWildcardRejected(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 1)
=
  match x
    _ | 1 => Ok(value = 1)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); !hasCode(diags, "AIL4112") {
		t.Fatalf("wildcard alternative reported no AIL4112: %v", diags)
	}
}

// TestOrMixedSlot pins one scalar domain per slot: int
// beside string alternatives reuses the mixed-slot rule.
func TestOrMixedSlot(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 1)
=
  match x
    1 | "a" => Ok(value = 1)
    _ => Ok(value = 2)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); !hasError(diags) {
		t.Fatalf("mixed alternatives reported nothing")
	}
}

// TestOrConstAlts pins visible integer constants as
// alternatives, resolved like their literals.
func TestOrConstAlts(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out, m__A, m__B]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

const m__A: int rev 1 = 47

const m__B: int rev 1 = 63

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    slash(x = 47) => Ok(value = 1)
    query(x = 63) => Ok(value = 1)
    other(x = 9) => Ok(value = 2)
=
  match x
    m__A | m__B => Ok(value = 1)
    _ => Ok(value = 2)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); hasError(diags) {
		t.Fatalf("const alternatives reported: %v", diags)
	}
}

// TestOrBoolAlts pins boolean alternatives satisfying the
// both-arms rule: `true | false` in one source arm is the
// complete bool case, no `_` needed.
func TestOrBoolAlts(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(x: bool) -> M__Out rev 1
  emits []
  tests
    t(x = true) => Ok(value = 1)
    f(x = false) => Ok(value = 1)
=
  match x
    true | false => Ok(value = 1)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); hasError(diags) {
		t.Fatalf("bool alternatives reported: %v", diags)
	}
}

// TestOrEmit pins the lowering shape: alternative
// conditions join with ||, ranges keep closed comparisons.
func TestOrEmit(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    one(x = 1) => Ok(value = 1)
    teen(x = 12) => Ok(value = 12)
    out(x = 99) => Ok(value = 99)
=
  match x
    1 | 10..15 => Ok(value = x)
    _ => Ok(value = 99)
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.ail"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.ail")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	ts := string(raw)
	for _, want := range []string{"(x === 1n || (x >= 10n && x <= 15n))"} {
		found := false
		for i := 0; i+len(want) <= len(ts); i++ {
			if ts[i:i+len(want)] == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("emit missing %q", want)
		}
	}
}

// TestOrVariantRejected pins no variant/error/Ok alternatives:
// they keep the historical value-match refusal.
func TestOrVariantRejected(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 1)
=
  match x
    Ok | 1 => Ok(value = 1)
    _ => Ok(value = 2)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); !hasError(diags) {
		t.Fatalf("variant alternative reported nothing")
	}
}

// TestOrMultiSlot pins alternatives inside one slot of a
// wider table: alternation never spans slots.
func TestOrMultiSlot(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(x: int, y: str) -> M__Out rev 1
  emits []
  tests
    hit(x = 1, y = "a") => Ok(value = 1)
    miss(x = 9, y = "z") => Ok(value = 2)
=
  match x, y
    1 | 2, "a" | "b" => Ok(value = 1)
    _, _ => Ok(value = 2)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); hasError(diags) {
		t.Fatalf("multi-slot alternatives reported: %v", diags)
	}
}

// TestOrCallMatchPipe pins `|` staying out of call-match `on`
// branches: alternatives are value patterns only.
func TestOrCallMatchPipe(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 1)
=
  match call m__id(x = x)
    on Ok r | Ok q => Ok(value = 1)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); !hasError(diags) {
		t.Fatalf("call-match pipe reported nothing")
	}
}

// TestOrOneObligation pins one AIL4107 obligation per source
// arm: an or-arm taken by any alternative satisfies the law,
// and an untaken or-arm is still exactly one obligation.
func TestOrOneObligation(t *testing.T) {
	src := rangeLib + `fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    one(x = 1) => Ok(value = 1)
    other(x = 9) => Ok(value = 2)
=
  match x
    1 | 2 => Ok(value = 1)
    _ => Ok(value = 2)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); hasError(diags) {
		t.Fatalf("taken or-arm reported: %v", diags)
	}
}
