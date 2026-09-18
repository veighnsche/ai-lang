package main

import (
	"strings"
	"testing"
)

// a39 S4: immutable append. Fixtures return BOTH sequences so a
// mutation cannot hide behind a correct final result; the aliasing
// row shares one base value between two appends.

// appendPair builds a module returning before/after around one append.
func appendPair(before, item, wantBefore, wantAfter string) string {
	return `mod t
  provides [t__app, T__Pair]
  uses []
  emits []

type T__Pair rev 1 (
  before: Seq<str>,
  after: Seq<str>
)

fn t__app() -> T__Pair rev 1
  emits []
  tests
    go() => Ok(before = ` + wantBefore + `, after = ` + wantAfter + `)
  Ok(before = ` + before + `, after = ` + before + ` + ` + item + `)
`
}

func TestSeqAppendRows(t *testing.T) {
	for _, c := range []struct {
		name, before, item, wantBefore, wantAfter string
	}{
		{"empty base", `Seq<str>[]`, `"a"`, `Seq<str>[]`, `Seq<str>["a"]`},
		{"empty member kept", `Seq<str>["a"]`, `""`, `Seq<str>["a"]`, `Seq<str>["a", ""]`},
		{"order repeats kept", `Seq<str>["b", "", "a"]`, `"b"`, `Seq<str>["b", "", "a"]`, `Seq<str>["b", "", "a", "b"]`},
	} {
		t.Run(c.name, func(t *testing.T) {
			body := appendPair(c.before, c.item, c.wantBefore, c.wantAfter)
			seqClean(t, map[string]string{"m.can": body}, "m.can")
		})
	}
}

// Two appends to one base value share nothing: if the evaluator
// mutated in place, base, first, and second would all disagree.
func TestSeqAppendAliasing(t *testing.T) {
	body := `mod t
  provides [t__fork, T__Fork]
  uses []
  emits []

type T__Fork rev 1 (
  base: Seq<str>,
  first: Seq<str>,
  second: Seq<str>
)

fn t__fork(base: Seq<str>) -> T__Fork rev 1
  emits []
  tests
    go(Seq<str>["a"]) => Ok(base = Seq<str>["a"], first = Seq<str>["a", "b"], second = Seq<str>["a", "c"])
  Ok(base = base, first = base + "b", second = base + "c")
`
	seqClean(t, map[string]string{"m.can": body}, "m.can")

	// And the mutation it guards against really fails: swap one
	// expectation and the run must disagree.
	bad := strings.Replace(body, `second = Seq<str>["a", "c"]`, `second = Seq<str>["a", "b", "c"]`, 1)
	seqCode(t, map[string]string{"m.can": bad}, "m.can", CodeTestFailed, "go")
}

// Branded appends keep the brand; a raw string member does not admit.
func TestSeqAppendBranded(t *testing.T) {
	body := `mod t
  provides [t__app, M__B, T__BPair]
  uses []
  emits []

brand M__B is str rev 1

type T__BPair rev 1 (
  before: Seq<M__B>,
  after: Seq<M__B>
)

fn t__app() -> T__BPair rev 1
  emits []
  tests
    go() => Ok(before = Seq<M__B>[seal M__B("A")], after = Seq<M__B>[seal M__B("A"), seal M__B("B")])
  Ok(before = Seq<M__B>[seal M__B("A")], after = Seq<M__B>[seal M__B("A")] + seal M__B("B"))
`
	seqClean(t, map[string]string{"m.can": body}, "m.can")

	bad := strings.Replace(body,
		`after = Seq<M__B>[seal M__B("A")] + seal M__B("B")`,
		`after = Seq<M__B>[seal M__B("A")] + "B"`, 1)
	seqCode(t, map[string]string{"m.can": bad}, "m.can", CodeTypeMismatch, "want M__B")
}

func TestSeqAppendRejects(t *testing.T) {
	for _, c := range []struct {
		name, expr, code, sub string
	}{
		{"int member", `Seq<str>["a"] + 1`, CodeTypeMismatch, "want str"},
		{"concat cut", `Seq<str>["a"] + Seq<str>["b"]`, CodeTypeMismatch, "concatenation is not in v1"},
		{"member left", `1 + Seq<str>["a"]`, CodeTypeMismatch, "no implicit conversions"},
	} {
		t.Run(c.name, func(t *testing.T) {
			// The bad expression stands alone in the after
			// field; only its diagnostic may fire.
			solo := `mod t
  provides [t__app, T__Pair]
  uses []
  emits []

type T__Pair rev 1 (
  before: Seq<str>,
  after: Seq<str>
)

fn t__app() -> T__Pair rev 1
  emits []
  tests
    go() => Ok(before = Seq<str>[], after = Seq<str>[])
  Ok(before = Seq<str>[], after = ` + c.expr + `)
`
			seqCode(t, map[string]string{"m.can": solo}, "m.can", c.code, c.sub)
		})
	}
}

// Emit lowers append to spread over the erased element type.
func TestSeqAppendEmit(t *testing.T) {
	body := appendPair(`Seq<str>["a"]`, `"b"`, `Seq<str>["a"]`, `Seq<str>["a", "b"]`)
	ts := compileEmit(t, body)
	for _, want := range []string{`[...`, `"b"`, "string[]"} {
		if !strings.Contains(ts, want) {
			t.Fatalf("emit missing %q:\n%s", want, ts)
		}
	}
}
