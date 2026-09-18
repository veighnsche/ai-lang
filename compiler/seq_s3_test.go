package main

import (
	"strings"
	"testing"
)

// a38 S3: checked access and traversal. Raw xs[i] is guarded-at-use
// access (the s[i] precedent): the compiler checks shapes and types,
// .can guards ensure bounds, and wrappers map guard failures to
// error constructors. There is no catch surface for index failure;
// unguarded access fails loud, never clamped.

// Get wrapper module: guards a la scalar_at plus a brand-preserving
// reader. Every guard arm below is taken by a committed row.
const seqGetMod = `mod t
  provides [t__get, t__get_b, M__B, T__Item, T__BItem]
  uses []
  emits [t.out_of_range]

error t.out_of_range(index: int, length: int)

brand M__B is str rev 1

type T__Item rev 1 (
  item: str
)

type T__BItem rev 1 (
  item: M__B
)

fn t__get(xs: Seq<str>, index: int) -> T__Item rev 1
  emits [t.out_of_range]
  tests
    first(Seq<str>["b", "a"], 0) => Ok(item = "b")
    second(Seq<str>["b", "a"], 1) => Ok(item = "a")
    empty_member(Seq<str>[""], 0) => Ok(item = "")
    negative(Seq<str>["a"], -1) => t.out_of_range(index = -1, length = 1)
    past_end(Seq<str>["a"], 1) => t.out_of_range(index = 1, length = 1)
    empty(Seq<str>[], 0) => t.out_of_range(index = 0, length = 0)
  match index >= 0
    true => match index < #xs
      true => Ok(item = xs[index])
      false => t.out_of_range(index = index, length = #xs)
    false => t.out_of_range(index = index, length = #xs)

fn t__get_b(xs: Seq<M__B>, index: int) -> T__BItem rev 1
  emits [t.out_of_range]
  tests
    one(Seq<M__B>[seal M__B("A")], 0) => Ok(item = seal M__B("A"))
    negative(Seq<M__B>[seal M__B("A")], -1) => t.out_of_range(index = -1, length = 1)
    past_end(Seq<M__B>[seal M__B("A")], 1) => t.out_of_range(index = 1, length = 1)
    empty(Seq<M__B>[], 0) => t.out_of_range(index = 0, length = 0)
  match index >= 0
    true => match index < #xs
      true => Ok(item = xs[index])
      false => t.out_of_range(index = index, length = #xs)
    false => t.out_of_range(index = index, length = #xs)
`

func TestSeqGetWrapper(t *testing.T) {
	seqClean(t, map[string]string{"m.can": seqGetMod}, "m.can")
}

// The brand row proves the static type, not just the runtime
// string: corrupt the expectation to str and the build must reject
// the field before any comparison runs (the erased strings would
// compare equal, so only the checker can catch this).
func TestSeqGetBrandTyped(t *testing.T) {
	bad := strings.Replace(seqGetMod,
		`one(Seq<M__B>[seal M__B("A")], 0) => Ok(item = seal M__B("A"))`,
		`one(Seq<M__B>[seal M__B("A")], 0) => Ok(item = "A")`, 1)
	seqCode(t, map[string]string{"m.can": bad}, "m.can", CodeTypeMismatch, "want M__B")
}

// Traversal module: bounded worker plus the fuel=#xs+1 entry. The
// B arm is untaken by full traversals by invariant
// (position+fuel=L+1 ends in E with fuel=1 left); direct
// zero/negative-budget rows witness it instead.
const seqWalkMod = `mod t
  provides [t__walk_from, t__walk, T__Trace]
  uses []
  emits []

type T__Trace rev 1 (
  count: int,
  trace: str
)

fn t__walk_from(xs: Seq<str>, position: int, fuel: int, count: int, trace: str) -> T__Trace rev 1
  decreases fuel
  emits []
  tests
    zero_budget(Seq<str>["x"], 0, 0, 0, "") => Ok(count = 0, trace = "")
    neg_budget(Seq<str>["x"], 0, -1, 0, "") => Ok(count = 0, trace = "")
    empty(Seq<str>[], 0, 1, 0, "") => Ok(count = 0, trace = "")
    one(Seq<str>["x"], 0, 2, 0, "") => Ok(count = 1, trace = "(x)")
    three(Seq<str>["b", "", "a"], 0, 4, 0, "") => Ok(count = 3, trace = "(b)()(a)")
    short(Seq<str>["a", "b"], 0, 1, 0, "") => Ok(count = 1, trace = "(a)")
  match fuel <= 0
    true => Ok(count = count, trace = trace)
    false => match position < #xs
      true => match call t__walk_from(xs, position + 1, fuel - 1, count + 1, trace + "(" + xs[position] + ")")
        on Ok r => Ok(count = r.count, trace = r.trace)
      false => Ok(count = count, trace = trace)

fn t__walk(xs: Seq<str>) -> T__Trace rev 1
  emits []
  tests
    empty(Seq<str>[]) => Ok(count = 0, trace = "")
    one(Seq<str>[""]) => Ok(count = 1, trace = "()")
    three(Seq<str>["A", "", "B"]) => Ok(count = 3, trace = "(A)()(B)")
  match call t__walk_from(xs, 0, #xs + 1, 0, "")
    on Ok r => Ok(count = r.count, trace = r.trace)
`

func TestSeqTraversal(t *testing.T) {
	seqClean(t, map[string]string{"m.can": seqWalkMod}, "m.can")
}

// B is a real arm, not an aspiration: drop every row that can take
// it (zero/negative budgets and the early-stop short row) and the
// coverage law fires, even though full traversals stay green.
func TestSeqTraversalBudgetArmReal(t *testing.T) {
	lines := strings.Split(seqWalkMod, "\n")
	var keep []string
	for _, l := range lines {
		if strings.Contains(l, "zero_budget") || strings.Contains(l, "neg_budget") ||
			strings.Contains(l, "short(Seq") {
			continue
		}
		keep = append(keep, l)
	}
	bare := strings.Join(keep, "\n")
	seqCode(t, map[string]string{"m.can": bare}, "m.can", CodeArmUntaken, "t__walk_from")
}

// The proof travels with the worker: no decreases, no build.
func TestSeqTraversalNeedsDecreases(t *testing.T) {
	bare := strings.Replace(seqWalkMod, "  decreases fuel\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bare})
	diags := diagnose(dir, "m.can", bare)
	found := false
	for _, d := range diags {
		if d.Sev == "error" && (d.Code == CodeBadDecreases || d.Code == CodeLocalCycle) {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected decreases/cycle refusal, got %v", diags)
	}
}

// Raw access never clamps: unguarded out-of-range fails the run.
func TestSeqRawNoClamp(t *testing.T) {
	body := `mod t
  provides [t__raw, T__Item]
  uses []
  emits []

type T__Item rev 1 (
  item: str
)

fn t__raw(xs: Seq<str>) -> T__Item rev 1
  emits []
  tests
    go(Seq<str>["a", "b"]) => Ok(item = "b")
  Ok(item = xs[5])
`
	seqCode(t, map[string]string{"m.can": body}, "m.can", CodeTestFailed, "seq index out of range")
}

// A non-int index over a sequence keeps the pinned diagnostic.
func TestSeqNonIntIndex(t *testing.T) {
	body := `mod t
  provides [t__raw, T__Item]
  uses []
  emits []

type T__Item rev 1 (
  item: str
)

fn t__raw(xs: Seq<str>) -> T__Item rev 1
  emits []
  tests
    go(Seq<str>["a"]) => Ok(item = "a")
  Ok(item = xs["a"])
`
	seqCode(t, map[string]string{"m.can": body}, "m.can", CodeTypeMismatch, "cannot index with")
}

// Emit lowers sequence indexing through the generic seq helper.
func TestSeqIndexEmit(t *testing.T) {
	ts := compileEmit(t, seqGetMod)
	for _, want := range []string{
		"$canSeqAt",
		"function $canSeqAt<T>(a: T[], i: bigint): T {",
	} {
		if !strings.Contains(ts, want) {
			t.Fatalf("emit missing %q:\n%s", want, ts)
		}
	}
}
