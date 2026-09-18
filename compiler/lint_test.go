package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// canlc lint goldens: exact report lines for all seven rules.
// Report text is pinned byte-for-byte (a88 golden-test capability);
// behavior first, then refusal shapes, then the proving-fixture
// folder, then idempotence.

const lintRedundant = `mod demo
  provides [demo__go, demo__pos, demo__id, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(left: int, right: int) -> Int__Value rev 1
  emits []
  tests
    named(left = 5, right = 8) => Ok(value = -3)
    mixed(5, right = 8) => Ok(value = -3)
    reordered(right = 8, left = 5) => Ok(value = -3)
  match call demo__id(value = left)
    on Ok r => Ok(value = r.value - right)

fn demo__pos(left: int, right: int) -> Int__Value rev 1
  emits []
  tests
    t(left = 5, right = 8) => Ok(value = -3)
  match call demo__id(left)
    on Ok r => Ok(value = r.value - right)

fn demo__id(value: int) -> Int__Value rev 1
  emits []
  tests
    t(value = 1) => Ok(value = 1)
  Ok(value = value)
`

// TestLintRedundantCtorNames pins a92 Rule 1 on constructions:
// in-order record, error, and Ok fields flag; reordered and
// positional spellings stay; exchange args stay exempt; a
// scripted outcome resolves through the callee return across
// files. Exact count proves no other rule fires.
func TestLintRedundantCtorNames(t *testing.T) {
	demo := `mod demo
  provides [demo__go, demo__box, Demo__Pair, Demo__Box]
  uses [lib__id@1]
  emits [demo.oops]

error demo.oops(code: int, note: str)

type Demo__Pair rev 1 (
  first: str
  second: str
)

type Demo__Box rev 1 (
  pair: Demo__Pair
)

fn demo__go(p: Demo__Pair) -> Demo__Pair rev 1
  emits [demo.oops]
  tests
    fat(Demo__Pair(first = "a", second = "b")) => Ok(first = "a", second = "b")
    slim(Demo__Pair(second = "b", first = "a")) => Ok("a", "b")
  match call lib__id(1)
    given
      fat => [exchange args (x = 1) outcome Ok(value = "ok")]
      slim => [exchange args (x = 1) outcome Ok("ok")]
    on Ok v => Ok(first = "x", second = "y")
    on demo.oops e => demo.oops(code = 0, note = "x")

fn demo__box(p: Demo__Pair) -> Demo__Box rev 1
  emits []
  tests
    nest(Demo__Pair("a", "b")) => Ok(pair = Demo__Pair(first = "a", second = "b"))
  Ok(pair = p)

fn demo__num(x: int) -> int rev 1
  emits []
  tests
    n(1) => Ok(value = 1)
    wild(1) => Ok(answer = 1)
  Ok(value = x)
`
	lib := `mod lib
  provides [lib__id, Lib__Box]
  uses []
  emits []

type Lib__Box rev 1 (
  value: str
)

fn lib__id(x: int) -> Lib__Box rev 1
  emits []
  tests
    t(1) => Ok("ok")
  Ok("ok")
`
	got, skipped := lintFiles(map[string]string{"demo.can": demo, "lib.can": lib})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:20: redundant argument name "first" (field 1 of Demo__Pair is "first"); write positionally`,
		`demo.can:20: redundant argument name "first" (field 1 of Ok is "first"); write positionally`,
		`demo.can:20: redundant argument name "second" (field 2 of Demo__Pair is "second"); write positionally`,
		`demo.can:20: redundant argument name "second" (field 2 of Ok is "second"); write positionally`,
		`demo.can:24: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
		`demo.can:26: redundant argument name "first" (field 1 of Ok is "first"); write positionally`,
		`demo.can:26: redundant argument name "second" (field 2 of Ok is "second"); write positionally`,
		`demo.can:27: redundant argument name "code" (field 1 of demo.oops is "code"); write positionally`,
		`demo.can:27: redundant argument name "note" (field 2 of demo.oops is "note"); write positionally`,
		`demo.can:32: redundant argument name "first" (field 1 of Demo__Pair is "first"); write positionally`,
		`demo.can:32: redundant argument name "pair" (field 1 of Ok is "pair"); write positionally`,
		`demo.can:32: redundant argument name "second" (field 2 of Demo__Pair is "second"); write positionally`,
		`demo.can:33: redundant argument name "pair" (field 1 of Ok is "pair"); write positionally`,
		`demo.can:38: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
		`demo.can:40: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

func TestLintRedundantNames(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintRedundant})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:13: redundant argument name "left" (param 1 of demo__go is "left"); write positionally`,
		`demo.can:13: redundant argument name "right" (param 2 of demo__go is "right"); write positionally`,
		`demo.can:13: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
		`demo.can:14: redundant argument name "right" (param 2 of demo__go is "right"); write positionally`,
		`demo.can:14: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
		`demo.can:15: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
		`demo.can:16: redundant argument name "value" (param 1 of demo__id is "value"); write positionally`,
		`demo.can:17: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
		`demo.can:22: redundant argument name "left" (param 1 of demo__pos is "left"); write positionally`,
		`demo.can:22: redundant argument name "right" (param 2 of demo__pos is "right"); write positionally`,
		`demo.can:22: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
		`demo.can:24: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
		`demo.can:29: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
		`demo.can:29: redundant argument name "value" (param 1 of demo__id is "value"); write positionally`,
		`demo.can:30: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

const lintLadder = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: str) -> Int__Value rev 1
  emits []
  tests
    t(value = "a") => Ok(1)
  match value
    "a" => Ok(1)
    "b" => Ok(1)
    "c" => Ok(2)
    "d" => Ok(1)
    _ => Ok(0)
`

func TestLintMergeableLadder(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintLadder})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	var merges []string
	for _, f := range got {
		if strings.Contains(f.msg, "mergeable") {
			merges = append(merges, f.String())
		}
	}
	want := []string{
		`demo.can:15: mergeable match arms [15 16] (identical outcomes); fold into one or-pattern arm, saves 1 lines`,
	}
	if len(merges) != len(want) || merges[0] != want[0] {
		t.Fatalf("merges = %v, want %v (full report %v)", merges, want, got)
	}
}

// Diagonal bool tables stay tables: or binds tighter than the
// comma, so (true,true)|(false,false) cannot spell as one arm
// (a28 truth-table idiom survives a88).
const lintDiagonal = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(left: int, right: int) -> Int__Value rev 1
  emits []
  tests
    t(left = 1, right = 2) => Ok(0)
  match left < 0, right < 0
    true, true => Ok(1)
    true, false => Ok(-1)
    false, true => Ok(-1)
    false, false => Ok(1)
`

func TestLintDiagonalTablesUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintDiagonal})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "mergeable") {
			t.Fatalf("diagonal table must not merge, got %v (full %v)", f, got)
		}
	}
}

// Variant arms never group: binders cannot be alpha-renamed
// into one, and case-or is not a spelled shape.
const lintVariant = `mod demo
  provides [demo__go, Form__State, Form__Msg]
  uses []
  emits []

variant Form__State rev 1 (
  case Empty()
  case Named(name: str)
)

type Form__Msg rev 1 (
  message: str
)

fn demo__go(state: Form__State) -> Form__Msg rev 1
  emits []
  tests
    t(state = Form__State.Empty()) => Ok("none")
  match state
    on Form__Empty _ => Ok("none")
    on Form__Named n => Ok("none")
`

func TestLintVariantArmsUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintVariant})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "mergeable") {
			t.Fatalf("variant arms must not merge, got %v (full %v)", f, got)
		}
	}
}

// Non-adjacent same-outcome arms stay silent: folding across
// the interleaved "c" arm would steal its overlapping cells.
const lintGap = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: str) -> Int__Value rev 1
  emits []
  tests
    t(value = "a") => Ok(1)
  match value
    "a" => Ok(1)
    "z" => Ok(9)
    "b" => Ok(1)
`

func TestLintGappedPairUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintGap})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "mergeable") {
			t.Fatalf("gapped pair must not merge, got %v (full %v)", f, got)
		}
	}
}

const lintChain = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits [demo.failed]

error demo.failed(value: int)

type Int__Value rev 1 (
  value: int
)

fn demo__step(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(1)
  Ok(value)

fn demo__go(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(1)
  match call demo__step(value)
    on Ok a => match call demo__step(a.value)
      on Ok b => Ok(b.value)
      on demo.failed e => demo.failed(0)
    on demo.failed e => demo.failed(0)
`

func TestLintChainable(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintChain})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:22: chainable ladder: 2 sequential calls sharing one failure outcome; rewrite as match chain`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// Divergent failure payloads share no else: no chain, and the
// constant payloads are not relays either.
const lintChainDiverged = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits [demo.failed]

error demo.failed(value: int)

type Int__Value rev 1 (
  value: int
)

fn demo__step(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(1)
  Ok(value)

fn demo__go(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(1)
  match call demo__step(value)
    on Ok a => match call demo__step(a.value)
      on Ok b => Ok(b.value)
      on demo.failed e => demo.failed(0)
    on demo.failed e => demo.failed(1)
`

func TestLintChainDivergedUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintChainDiverged})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	if len(got) != 0 {
		t.Fatalf("divergent ladder must be clean, got %v", got)
	}
}

const lintTabled = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(a: int, b: int) -> Int__Value rev 1
  emits []
  tests
    t(1, 2) => Ok(3)
  match a < 0
    true => match b < 0
      true => Ok(0)
      false => Ok(1)
    false => match b < 0
      true => Ok(2)
      false => Ok(3)
`

func TestLintTableable(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintTabled})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:14: nested matches share one scrutinee; fold into a multi-scrutinee table`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// Guarded nests stay nested: an impure inner scrutinee must
// keep its guard (CAN4109). Different pure scrutinees fold
// under the ladder rule instead (see lintTabulate).
const lintTabledKept = `mod demo
  provides [demo__guarded, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__guarded(s: str) -> Int__Value rev 1
  emits []
  tests
    t("ab") => Ok(1)
  match #s >= 1
    true => match s[0:1] == "a"
      true => Ok(1)
      false => Ok(0)
    false => Ok(0)
`

func TestLintTabledKeptUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintTabledKept})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "multi-scrutinee") {
			t.Fatalf("kept nests must not fold, got %v (full %v)", f, got)
		}
	}
}

const lintSame = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: int) -> Int__Value rev 1
  emits []
  tests
    t(1) => Ok(0)
  match value < 0
    true => Ok(0)
    false => Ok(0)
`

func TestLintSameOutcome(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintSame})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:14: match always yields Ok(0); drop the match`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// A faulting scrutinee keeps its match even when every arm
// agrees: dropping it would drop the fault.
const lintSameGuarded = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(s: str) -> Int__Value rev 1
  emits []
  tests
    t("ab") => Ok(0)
  match s[0:1] == "a"
    true => Ok(0)
    false => Ok(0)
`

func TestLintSameGuardedUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintSameGuarded})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	// The faulting scrutinee blocks the drop, but the adjacent
	// pair still folds: identical outcomes, scrutinee kept.
	want := []string{
		`demo.can:15: mergeable match arms [15 16] (identical outcomes); fold into one or-pattern arm, saves 1 lines`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

const lintRelay = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits [demo.failed, demo.other]

error demo.failed(value: int)
error demo.other(value: int, note: str)

fn demo__step(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(1)
  Ok(value)

fn demo__go(value: int) -> Int__Value rev 1
  emits [demo.failed, demo.other]
  tests
    t(1) => Ok(1)
  match call demo__step(value)
    on Ok a => Ok(a.value)
    on demo.failed e => demo.failed(e.value)
    on demo.other e => demo.other(e.value, "x")
`

func TestLintForwardable(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintRelay})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:21: handwritten relay of demo.failed; write forward e`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

const lintRanges = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: int) -> Int__Value rev 1
  emits []
  tests
    t(2) => Ok(1)
  match value
    1..3 => Ok(1)
    4..6 => Ok(1)
    8..9 => Ok(2)
    _ => Ok(0)
`

func TestLintRangeMerge(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintRanges})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:15: mergeable ranges 1..3, 4..6; join into 1..6`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// A three-rung equality ladder over one (fault-capable) base:
// the restatement is fault-neutral, so purity is not required.
const lintRestate = `mod demo
  provides [demo__go, Demo__Out]
  uses []
  emits []

type Demo__Out rev 1 (
  value: str
)

fn demo__go(s: str) -> Demo__Out rev 1
  emits []
  tests
    amp("x&y") => Ok("hit")
  match s[0:1] == "&"
    true => Ok("amp")
    false => match s[0:1] == "<"
      true => Ok("lt")
      false => match s[0:1] == ">"
        true => Ok("gt")
        false => Ok("plain")
`

func TestLintRestatable(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintRestate})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:14: equality ladder over one base with 3 rungs; restate as a match on the base, saves 2 lines`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// Ladders that must not restate: a lone rung saves nothing;
// mixed kinds, bool literals, and != never form a str table;
// duplicate literals are dead code, not a table; guards are
// <=-shaped and never rungs.
const lintRestateKept = `mod demo
  provides [demo__single, demo__mixed, demo__bools, demo__neq, demo__dup, demo__guard, Demo__Out]
  uses []
  emits []

type Demo__Out rev 1 (
  value: str
)

fn demo__single(s: str) -> Demo__Out rev 1
  emits []
  tests
    t("ab") => Ok("x")
  match s[0:1] == "a"
    true => Ok("hit")
    false => Ok("miss")

fn demo__mixed(s: str, n: int) -> Demo__Out rev 1
  emits []
  tests
    t("ab", 1) => Ok("x")
  match s[0:1] == "a"
    true => Ok("hit")
    false => match n == 1
      true => Ok("one")
      false => Ok("other")

fn demo__bools(b: bool) -> Demo__Out rev 1
  emits []
  tests
    t(true) => Ok("x")
  match b == true
    true => Ok("hit")
    false => match b == false
      true => Ok("miss")
      false => Ok("dead")

fn demo__neq(s: str) -> Demo__Out rev 1
  emits []
  tests
    t("ab") => Ok("x")
  match s[0:1] != "a"
    true => Ok("hit")
    false => match s[0:1] != "b"
      true => Ok("hit2")
      false => Ok("miss")

fn demo__dup(s: str) -> Demo__Out rev 1
  emits []
  tests
    t("ab") => Ok("x")
  match s[0:1] == "a"
    true => Ok("hit")
    false => match s[0:1] == "a"
      true => Ok("hit2")
      false => Ok("miss")

fn demo__guard(n: int) -> Demo__Out rev 1
  emits []
  tests
    t(1) => Ok("x")
  match n <= 0
    true => Ok("base")
    false => match n == 1
      true => Ok("one")
      false => Ok("other")
`

func TestLintRestateKeptUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintRestateKept})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "restate as a match on the base") {
			t.Fatalf("kept ladders must not restate, got %v (full %v)", f, got)
		}
	}
}

// A pure two-level ladder and a pure diamond over different
// scrutinees: both fold into product tables with don't-care
// slots, all outcomes plain.
const lintTabulate = `mod demo
  provides [demo__ladder, demo__diamond, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__ladder(a: int, b: int) -> Int__Value rev 1
  emits []
  tests
    t(1, 2) => Ok(2)
  match a == 0
    true => Ok(0)
    false => match b == 0
      true => Ok(1)
      false => Ok(2)

fn demo__diamond(a: int, b: int, c: int) -> Int__Value rev 1
  emits []
  tests
    t(1, 2, 3) => Ok(6)
  match a < 0
    true => match b < 0
      true => Ok(0)
      false => Ok(1)
    false => match c < 0
      true => Ok(2)
      false => Ok(3)
`

func TestLintLadderable(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintTabulate})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:14: nested matches over 2 pure scrutinees; fold into a multi-scrutinee table`,
		`demo.can:24: nested matches over 3 pure scrutinees; fold into a multi-scrutinee table`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// Nests that must not tabulate: a guarded (impure) ladder keeps
// its guard; an identical-scrutinee ladder is degenerate;
// tables are already tables; a shared-scrutinee diamond is rule
// 4's, not this rule's; arms that run calls first stay nested
// (can-idioms C1).
const lintTabulateKept = `mod demo
  provides [demo__guarded, demo__same, demo__tabled, demo__shared, demo__help, demo__calls, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__guarded(s: str) -> Int__Value rev 1
  emits []
  tests
    t("ab") => Ok(1)
  match #s >= 1
    true => match s[0:1] == "a"
      true => Ok(1)
      false => Ok(0)
    false => Ok(0)

fn demo__same(a: int) -> Int__Value rev 1
  emits []
  tests
    t(1) => Ok(1)
  match a == 0
    true => Ok(0)
    false => match a == 0
      true => Ok(1)
      false => Ok(2)

fn demo__tabled(a: int, b: int) -> Int__Value rev 1
  emits []
  tests
    t(1, 2) => Ok(3)
  match a == 0, b == 0
    true, true => Ok(0)
    true, false => Ok(1)
    false, true => Ok(2)
    false, false => Ok(3)

fn demo__shared(a: int, b: int) -> Int__Value rev 1
  emits []
  tests
    t(1, 2) => Ok(3)
  match a < 0
    true => match b < 0
      true => Ok(0)
      false => Ok(1)
    false => match b < 0
      true => Ok(2)
      false => Ok(3)

fn demo__help(value: int) -> Int__Value rev 1
  emits []
  tests
    t(1) => Ok(1)
  Ok(value)

fn demo__calls(a: int, b: int) -> Int__Value rev 1
  emits []
  tests
    t(1, 2) => Ok(2)
  match a < 0
    true => match call demo__help(a)
      on Ok r => Ok(r.value)
    false => match b < 0
      true => Ok(1)
      false => Ok(2)
`

func TestLintLadderKeptUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintTabulateKept})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "pure scrutinees") {
			t.Fatalf("kept nests must not tabulate, got %v (full %v)", f, got)
		}
	}
}

// a90 interplay: lint runs on parse, so forward-call sites read
// as single-expression outcomes. A pure ladder with a
// forward-call leaf still tabulates (bodies move verbatim);
// identical forward calls or-fold and drop like any identical
// outcomes. No lint rule needed edits for the shape.
const lintForwardCallInterplay = `mod demo
  provides [demo__ladder, demo__fold, demo__drop, demo__help, Demo__Out]
  uses []
  emits []

type Demo__Out rev 1 (
  value: str
)

fn demo__help(v: str) -> Demo__Out rev 1
  emits []
  tests
    t("a") => Ok("a")
  Ok(v)

fn demo__ladder(a: int, b: int) -> Demo__Out rev 1
  emits []
  tests
    t(1, 2) => Ok("x")
  match a < 0
    true => forward call demo__help("x")
    false => match b < 0
      true => Ok("y")
      false => Ok("z")

fn demo__fold(s: str) -> Demo__Out rev 1
  emits []
  tests
    t("a") => Ok("x")
  match s
    "a" => forward call demo__help("x")
    "b" => forward call demo__help("x")
    _ => Ok("z")

fn demo__drop(v: int) -> Demo__Out rev 1
  emits []
  tests
    t(1) => Ok("x")
  match v < 0
    true => forward call demo__help("x")
    false => forward call demo__help("x")
`

func TestLintForwardCallInterplay(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintForwardCallInterplay})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:20: nested matches over 2 pure scrutinees; fold into a multi-scrutinee table`,
		`demo.can:31: mergeable match arms [31 32] (identical outcomes); fold into one or-pattern arm, saves 1 lines`,
		`demo.can:39: match always yields forward call demo__help("x"); drop the match`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// A two-arm relay and a one-arm infallible relay over
// same-file locals: both rewrite as forward call.
const lintRelayCall = `mod demo
  provides [demo__go, demo__one, demo__work, demo__pure, Demo__Out, Demo__Work]
  uses []
  emits [demo.bad]

error demo.bad(value: str)

type Demo__Out rev 1 (
  value: str
)

type Demo__Work rev 1 (
  value: str
)

fn demo__work(v: str) -> Demo__Work rev 1
  emits [demo.bad]
  tests
    t("a") => Ok("a")
  Ok(v)

fn demo__pure(v: str) -> Demo__Out rev 1
  emits []
  tests
    t("a") => Ok("a")
  Ok(v)

fn demo__go(x: str) -> Demo__Out rev 1
  emits [demo.bad]
  tests
    t("a") => Ok("a")
  match call demo__work(x)
    on Ok r => forward r
    on demo.bad e => forward e

fn demo__one(x: str) -> Demo__Out rev 1
  emits []
  tests
    t("a") => Ok("a")
  match call demo__pure(x)
    on Ok r => forward r
`

func TestLintRelayCallable(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintRelayCall})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:32: relay match forwards every outcome; write forward call demo__work(x)`,
		`demo.can:40: relay match forwards every outcome; write forward call demo__pure(x)`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// Matches that must not rewrite: a rebuilding arm breaks the
// relay; given tables, uses-pinned callees, and externs cannot
// convert (forward call carries no given); an already-relayed
// arm is the fixed shape.
const lintRelayCallKept = `mod demo
  provides [demo__rebuild, demo__given, demo__foreign, demo__extern, demo__fixed, demo__work, ex__work, Demo__Out, Demo__Work]
  uses [other__work@1]
  emits [demo.bad]

error demo.bad(value: str)

extern ex__work(v: str) -> Demo__Work rev 1

type Demo__Out rev 1 (
  value: str
)

type Demo__Work rev 1 (
  value: str
)

fn demo__work(v: str) -> Demo__Work rev 1
  emits [demo.bad]
  tests
    t("a") => Ok("a")
  Ok(v)

fn demo__rebuild(x: str) -> Demo__Out rev 1
  emits [demo.bad]
  tests
    t("a") => Ok("a")
  match call demo__work(x)
    on Ok r => Ok(r.value)
    on demo.bad e => forward e

fn demo__given(x: str) -> Demo__Out rev 1
  emits [demo.bad]
  tests
    t("a") => Ok("a")
  match call other__work(x)
    given
      t => exchange args (x = "a") outcome Ok("a")
    on Ok r => forward r
    on demo.bad e => forward e

fn demo__foreign(x: str) -> Demo__Out rev 1
  emits [demo.bad]
  tests
    t("a") => Ok("a")
  match call other__work(x)
    on Ok r => forward r
    on demo.bad e => forward e

fn demo__extern(x: str) -> Demo__Out rev 1
  emits [demo.bad]
  tests
    t("a") => Ok("a")
  match call ex__work(x)
    on Ok r => forward r
    on demo.bad e => forward e

fn demo__fixed(x: str) -> Demo__Out rev 1
  emits [demo.bad]
  tests
    t("a") => Ok("a")
  match x == "z"
    true => Ok("no")
    false => forward call demo__work(x)
`

func TestLintRelayCallKeptUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintRelayCallKept})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "write forward call") {
			t.Fatalf("kept matches must not rewrite, got %v (full %v)", f, got)
		}
	}
}

// Stale scripts (a91): solo belongs to a fn that never calls the
// scripted fns, so both its rows flag; own and caller keys stay
// silent, as do the unknown key (compiler warning) and the dash
// row (compiler error).
const lintUnreached = `mod demo
  provides [demo__go, demo__chain, demo__help, demo__other, Demo__S, net__fetch]
  uses []
  emits [net.down, demo.bad]

error net.down()
error demo.bad()

type Demo__S rev 1 (
  body: str
)

extern net__fetch() -> Demo__S rev 1
  emits [net.down]

fn demo__go() -> Demo__S rev 1
  emits [demo.bad]
  tests
    ok() => Ok("hi")
    down() => demo.bad()
  match call net__fetch()
    given
      ok => [exchange args () outcome Ok("hi")]
      down => [exchange args () outcome net.down()]
      caller => [exchange args () outcome net.down()]
      solo => [exchange args () outcome net.down()]
      zzz => [exchange args () outcome net.down()]
      dead => -
    on net.down _ => demo.bad()
    on Ok d => Ok(d.body)

fn demo__help() -> Demo__S rev 1
  emits [demo.bad]
  tests
    caller() => Ok("hi")
  match call demo__go()
    on demo.bad e => forward e
    on Ok d => Ok(d.body)

fn demo__other() -> Demo__S rev 1
  emits [demo.bad]
  tests
    solo() => demo.bad()
  demo.bad()

fn demo__chain() -> Demo__S rev 1
  emits [demo.bad]
  tests
    linked() => Ok("hi")
  match chain
    call net__fetch() as d
      given
        linked => exchange args () outcome Ok("hi")
        solo => exchange args () outcome net.down()
    then Ok(d.body)
    else demo.bad()
`

func TestLintUnreachedKeys(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintUnreached})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:26: given key "solo" names a test that cannot reach this call; delete the row`,
		`demo.can:54: given key "solo" names a test that cannot reach this call; delete the row`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
		if got[i].code != CodeLintUnreachedKey {
			t.Fatalf("finding %d code = %q, want CAN3420", i, got[i].code)
		}
	}
}

func TestLintUnreachedKeptUntouched(t *testing.T) {
	clean := strings.Replace(lintUnreached, "      solo => [exchange args () outcome net.down()]\n", "", 1)
	clean = strings.Replace(clean, "        solo => exchange args () outcome net.down()\n", "", 1)
	got, skipped := lintFiles(map[string]string{"demo.can": clean})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "cannot reach this call") {
			t.Fatalf("kept tables must not flag, got %v (full %v)", f, got)
		}
	}
}

// Two sequential calls that cannot fail share no failure outcome:
// the router's divergent plain sides must not establish unity for
// the run (regression: the rejected router leaked failures=true).
const lintChainInfallible = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__step(value: int) -> Int__Value rev 1
  emits []
  tests
    t(1) => Ok(1)
  Ok(value)

fn demo__go(value: int) -> Int__Value rev 1
  emits []
  tests
    t(1) => Ok(1)
  match call demo__step(value)
    on Ok a => match call demo__step(a.value)
      on Ok b => match b.value == 0
        true => Ok(0)
        false => Ok(b.value)
`

func TestLintChainInfallibleUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintChainInfallible})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "chainable") {
			t.Fatalf("infallible ladder must not chain, got %v (full %v)", f, got)
		}
	}
}

// A ladder of relays shares failure TEXT but no failure VALUE:
// each arm rebuilds from its own binder, and chain elaboration
// binds error arms with wildcards, so no shared else can spell
// it. Chainable stays silent; the relay rule reports each arm.
const lintChainRelay = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits [demo.failed]

error demo.failed(value: int)

type Int__Value rev 1 (
  value: int
)

fn demo__step(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(1)
  Ok(value)

fn demo__go(value: int) -> Int__Value rev 1
  emits [demo.failed, demo.other]
  tests
    t(1) => Ok(1)
  match call demo__step(value)
    on Ok a => match call demo__step(a.value)
      on Ok b => Ok(b.value)
      on demo.failed e => demo.failed(e.value)
    on demo.failed e => demo.failed(e.value)
`

// A ladder of forwards shares one spelling but no shared else:
// forward rebuilds from each arm's own binder, which the shared
// else scope cannot see. Chainable stays silent.
const lintChainForward = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits [demo.failed]

error demo.failed(value: int)

type Int__Value rev 1 (
  value: int
)

fn demo__step(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(1)
  Ok(value)

fn demo__go(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(1)
  match call demo__step(value)
    on Ok a => match call demo__step(a.value)
      on Ok b => Ok(b.value)
      on demo.failed e => forward e
    on demo.failed e => forward e
`

func TestLintChainForwardUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintChainForward})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	if len(got) != 0 {
		t.Fatalf("forward ladder must not chain, got %v", got)
	}
}

func TestLintChainRelayUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintChainRelay})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	relays := 0
	for _, f := range got {
		if strings.Contains(f.msg, "chainable") {
			t.Fatalf("relay ladder must not chain, got %v (full %v)", f, got)
		}
		if strings.Contains(f.msg, "handwritten relay") {
			relays++
			if f.code != CodeLintRelay {
				t.Fatalf("relay finding carries %q, want %q", f.code, CodeLintRelay)
			}
		}
	}
	if relays != 2 {
		t.Fatalf("expected 2 relay findings, got %v", got)
	}
}

func TestLintSkipsUnparseable(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"bad.can": "mod demo\n  provides ["})
	if len(got) != 0 {
		t.Fatalf("expected no findings, got %v", got)
	}
	if len(skipped) != 1 || skipped[0] != "bad.can" {
		t.Fatalf("expected one skip, got %v", skipped)
	}
}

// TestLintErrorFixtures is the folder golden for
// sketches/lint-errors: every file must parse, every file must
// yield at least one finding, and the findings are pinned
// exactly — one rule family per file, nothing advisory-only.
// Lint findings are error-severity end to end: the folder run
// exits 1 with every finding on stdout.
func TestLintErrorFixtures(t *testing.T) {
	entries, err := os.ReadDir("../sketches/lint-errors")
	if err != nil {
		t.Fatalf("read fixture dir: %v", err)
	}
	files := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".can") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join("../sketches/lint-errors", e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		files[e.Name()] = string(raw)
	}
	if len(files) != 11 {
		t.Fatalf("want 11 proving fixtures, got %d: %v", len(files), files)
	}
	got, skipped := lintFiles(files)
	if len(skipped) != 0 {
		t.Fatalf("fixtures must parse, skipped: %v", skipped)
	}
	byFile := map[string][]string{}
	for _, f := range got {
		byFile[f.file] = append(byFile[f.file], f.String())
	}
	want := map[string][]string{
		"chain.can": {
			`chain.can:34: chainable ladder: 2 sequential calls sharing one failure outcome; rewrite as match chain`,
		},
		"merge.can": {
			`merge.can:23: mergeable match arms [23 24] (identical outcomes); fold into one or-pattern arm, saves 1 lines`,
		},
		"outcome.can": {
			`outcome.can:20: match always yields Ok(0); drop the match`,
		},
		"ranges.can": {
			`ranges.can:22: mergeable ranges 1..3, 4..6; join into 1..6`,
		},
		"redundant.can": {
			`redundant.can:17: redundant argument name "left" (param 1 of redundant__sub is "left"); write positionally`,
			`redundant.can:17: redundant argument name "right" (param 2 of redundant__sub is "right"); write positionally`,
			`redundant.can:17: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
			`redundant.can:18: redundant argument name "right" (param 2 of redundant__sub is "right"); write positionally`,
			`redundant.can:18: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
			`redundant.can:19: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
			`redundant.can:20: redundant argument name "value" (param 1 of redundant__id is "value"); write positionally`,
			`redundant.can:21: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
			`redundant.can:26: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
			`redundant.can:26: redundant argument name "value" (param 1 of redundant__id is "value"); write positionally`,
			`redundant.can:27: redundant argument name "value" (field 1 of Ok is "value"); write positionally`,
		},
		"relay.can": {
			`relay.can:38: handwritten relay of relay.failed; write forward e`,
		},
		"table.can": {
			`table.can:22: nested matches share one scrutinee; fold into a multi-scrutinee table`,
		},
		"restate.can": {
			`restate.can:21: equality ladder over one base with 2 rungs; restate as a match on the base, saves 1 lines`,
		},
		"ladder.can": {
			`ladder.can:21: nested matches over 2 pure scrutinees; fold into a multi-scrutinee table`,
		},
		"relaycall.can": {
			`relaycall.can:35: relay match forwards every outcome; write forward call relaycall__work(x)`,
		},
		"stale.can": {
			`stale.can:31: given key "solo" names a test that cannot reach this call; delete the row`,
		},
	}
	for name, lines := range want {
		if len(byFile[name]) != len(lines) {
			t.Fatalf("%s findings = %v, want %v", name, byFile[name], lines)
		}
		for i := range lines {
			if byFile[name][i] != lines[i] {
				t.Fatalf("%s finding %d = %q, want %q", name, i, byFile[name][i], lines[i])
			}
		}
	}
	for name := range byFile {
		if _, ok := want[name]; !ok {
			t.Fatalf("unexpected findings for %s: %v", name, byFile[name])
		}
	}
	var stdout, stderr bytes.Buffer
	if code := runLint(&stdout, &stderr, []string{"../sketches/lint-errors"}); code != 1 {
		t.Fatalf("folder lint exit = %d, want 1 (findings are errors): stdout %q stderr %q",
			code, stdout.String(), stderr.String())
	}
	if n := len(strings.Split(strings.TrimSpace(stdout.String()), "\n")); n != len(got) {
		t.Fatalf("folder lint printed %d lines, want %d (one per finding)", n, len(got))
	}
}

// TestLintSpanPositions pins the 2D editor contract: every
// lint-errors fixture finding publishes a [start, end) UTF-16
// span covering exactly the offending token — never the whole
// line. Columns are hand-checked against the fixture sources;
// the covered text is asserted too, so a locator that drifts
// onto the wrong token fails here, not in Cursor.
func TestLintSpanPositions(t *testing.T) {
	entries, err := os.ReadDir("../sketches/lint-errors")
	if err != nil {
		t.Fatalf("read fixture dir: %v", err)
	}
	texts := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".can") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join("../sketches/lint-errors", e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		texts[e.Name()] = string(raw)
	}
	type want struct {
		line, start, end int
		code, cover      string
	}
	wants := map[string][]want{
		"chain.can": {
			{34, 13, 31, CodeLintChain, "chain__step(value)"},
		},
		"merge.can": {
			{23, 4, 7, CodeLintOrFold, `"a"`},
		},
		"outcome.can": {
			{20, 2, 17, CodeLintSameOutcome, "match value < 0"},
		},
		"ranges.can": {
			{22, 4, 8, CodeLintRange, "1..3"},
		},
		"redundant.can": {
			{17, 12, 18, CodeLintRedundant, "left ="},
			{17, 22, 29, CodeLintRedundant, "right ="},
			{17, 39, 46, CodeLintRedundant, "value ="},
			{18, 13, 20, CodeLintRedundant, "right ="},
			{18, 30, 37, CodeLintRedundant, "value ="},
			{19, 41, 48, CodeLintRedundant, "value ="},
			{20, 27, 34, CodeLintRedundant, "value ="},
			{21, 18, 25, CodeLintRedundant, "value ="},
			{26, 25, 32, CodeLintRedundant, "value ="},
			{26, 8, 15, CodeLintRedundant, "value ="},
			{27, 5, 12, CodeLintRedundant, "value ="},
		},
		"relay.can": {
			{38, 25, 46, CodeLintRelay, "relay.failed(e.value)"},
		},
		"table.can": {
			{22, 2, 13, CodeLintTable, "match a < 0"},
		},
		"restate.can": {
			{21, 2, 21, CodeLintRestate, "match s[0:1] == \"&\""},
		},
		"ladder.can": {
			{21, 2, 13, CodeLintLadder, "match a < 0"},
		},
		"relaycall.can": {
			{35, 2, 31, CodeLintRelayCall, "match call relaycall__work(x)"},
		},
		"stale.can": {
			{31, 6, 10, CodeLintUnreachedKey, "solo"},
		},
	}
	for name, ws := range wants {
		lines := strings.Split(texts[name], "\n")
		got := lintDiagsFor(name, texts)
		if len(got) != len(ws) {
			t.Fatalf("%s diags = %d, want %d", name, len(got), len(ws))
		}
		for i, w := range ws {
			d := got[i]
			if d.Line != w.line || d.Start != w.start || d.End != w.end || d.Code != w.code {
				t.Fatalf("%s diag %d = line %d [%d,%d) %s, want line %d [%d,%d) %s",
					name, i, d.Line, d.Start, d.End, d.Code, w.line, w.start, w.end, w.code)
			}
			if cover := lines[d.Line-1][d.Start:d.End]; cover != w.cover {
				t.Fatalf("%s diag %d covers %q, want %q", name, i, cover, w.cover)
			}
		}
	}
	// No finding in the folder may fall back to whole-line:
	// every lint squiggle must be column-precise.
	for name := range texts {
		for _, d := range lintDiagsFor(name, texts) {
			if d.End <= d.Start {
				t.Fatalf("%s:%d %s has no span (whole-line fallback)", name, d.Line, d.Code)
			}
		}
	}
}

// TestLintFindNameEq pins the redundant-name locator: the
// (skip+1)-th `name =` outside strings, with `=>`/`==` and
// dotted-prefix guards.
func TestLintFindNameEq(t *testing.T) {
	cases := []struct {
		code, name string
		skip       int
		ns, ne     int
		ok         bool
	}{
		{"ordered(left = 5, right = 8) => Ok(value = -3)", "left", 0, 8, 14, true},
		{"ordered(left = 5, right = 8) => Ok(value = -3)", "right", 0, 18, 25, true},
		{"one(value = 1) => Ok(value = 1)", "value", 0, 4, 11, true},
		{"one(value = 1) => Ok(value = 1)", "value", 1, 21, 28, true},
		{"relay.failed(value = e.value)", "value", 0, 13, 20, true},
		{"f(x == 1)", "x", 0, 0, 0, false},
		{"m(true => 1)", "true", 0, 0, 0, false},
		{`"left = 5"`, "left", 0, 0, 0, false},
		{"on relay.failed e => x", "failed", 0, 0, 0, false},
		{"f(note = \"a = b\")", "a", 0, 0, 0, false},
	}
	for _, c := range cases {
		ns, ne, ok := lintFindNameEq(c.code, c.name, c.skip)
		if ok != c.ok || ns != c.ns || ne != c.ne {
			t.Fatalf("lintFindNameEq(%q, %q, %d) = [%d,%d) %v, want [%d,%d) %v",
				c.code, c.name, c.skip, ns, ne, ok, c.ns, c.ne, c.ok)
		}
	}
}

// TestLintDiagnoseErrors pins the editor contract: lint findings
// publish through diagnose as error-severity diagnostics with
// registered CAN3410-3416 codes — this is what Cursor draws.
// Clean-world layering holds: a file with real compiler errors
// shows those first, with lint suppressed until they are fixed.
func TestLintDiagnoseErrors(t *testing.T) {
	src := `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: str) -> Int__Value rev 1
  emits []
  tests
    t("a") => Ok(1)
    b("b") => Ok(1)
    c("c") => Ok(2)
    z("z") => Ok(0)
  match value
    "a" => Ok(1)
    "b" => Ok(1)
    "c" => Ok(2)
    _ => Ok(0)
`
	dir := writeLSPDir(t, map[string]string{"demo.can": src})
	diags := diagnose(dir, "demo.can", src)
	var found *Diag
	for i, d := range diags {
		if d.Code == CodeLintOrFold {
			found = &diags[i]
		}
	}
	if found == nil {
		t.Fatalf("expected CAN3411 or-fold diag, got %v", diags)
	}
	if found.Sev != "error" {
		t.Fatalf("lint diag severity = %q, want error", found.Sev)
	}
	if found.Line != 18 {
		t.Fatalf("lint diag line = %d, want 18", found.Line)
	}
}

func TestLintDiagnoseSuppressedOnBroken(t *testing.T) {
	src := `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: str) -> Int__Value rev 1
  emits []
  tests
    t("a") => Ok(1)
  match value
    "a" => Ok(1)
    "b" => Ok(1)
`
	dir := writeLSPDir(t, map[string]string{"demo.can": src})
	diags := diagnose(dir, "demo.can", src)
	for _, d := range diags {
		if strings.HasPrefix(d.Code, "CAN341") {
			t.Fatalf("broken file must show real errors first, got lint %v (full %v)", d, diags)
		}
	}
	if len(diags) == 0 {
		t.Fatalf("expected compiler errors for missing wildcard arm, got none")
	}
}

func TestLintIdempotent(t *testing.T) {
	files := map[string]string{"a.can": lintRedundant, "b.can": lintLadder}
	first, _ := lintFiles(files)
	second, _ := lintFiles(files)
	if len(first) != len(second) {
		t.Fatalf("lengths differ: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("run %d differs at %d", 2, i)
		}
	}
}

// TestLintStdClean pins the stdlib migration (can-idioms C13):
// every blessed std module lints clean — no redundant names,
// no unfolded ladders, no handwritten relays. All files lint
// together so cross-module signatures resolve as in
// `canlc lint std/`.
func TestLintStdClean(t *testing.T) {
	mods := []string{"ascii", "division", "html", "quota", "ratio", "scalars", "schema", "text"}
	files := map[string]string{}
	for _, m := range mods {
		name := m + ".can"
		raw, err := os.ReadFile(filepath.Join("../std", m, name))
		if err != nil {
			t.Fatalf("read std/%s/%s: %v", m, name, err)
		}
		files[name] = string(raw)
	}
	got, skipped := lintFiles(files)
	if len(skipped) != 0 {
		t.Fatalf("std must parse, skipped: %v", skipped)
	}
	if len(got) != 0 {
		var msgs []string
		for _, f := range got {
			msgs = append(msgs, f.String())
		}
		t.Fatalf("std must lint clean, got %d findings:\n%s", len(got), strings.Join(msgs, "\n"))
	}
}

// G2 generics: lint parses pre-expansion templates, so a template
// field naming a parameter (value: T) must read as an opaque sort,
// never as an unknown type or a spurious finding.
func TestLintGenericTypeTolerant(t *testing.T) {
	src := `mod m
  provides [m__wrap, m__pick, M__Box, M__O]
  uses []
  emits []

type M__Box<T> rev 1 (
  value: T
)

type M__O rev 1 (
  value: str
)

fn m__wrap(v: str) -> M__Box<str> rev 1
  emits []
  tests
    one("a") => Ok("a")
  Ok(v)

fn m__pick(b: M__Box<str>) -> M__O rev 1
  emits []
  tests
    hit(M__Box<str>("u")) => Ok("u")
    miss(M__Box<str>("x")) => Ok("x")
  match b.value
    "u" => Ok("u")
    _ => Ok(b.value)
`
	if diags := lintDiagsFor("m.can", map[string]string{"m.can": src}); len(diags) != 0 {
		t.Fatalf("expected no lint findings, got %v", diags)
	}
}
