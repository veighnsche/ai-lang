package main

import (
	"strings"
	"testing"
)

// a37 S2: # counts sequence elements. Each row is a tested function
// returning a record int field; a miscount fails the test run, so
// clean means the evaluator computed exactly the pinned count.

const seqLenOut = `type M__Len rev 1 (
  n: int
)
`

// seqLenMod builds a module whose body computes a length and whose
// expectation pins the literal count. The expectation is always a
// literal: repeating the expression there would make a miscount
// agree with itself and pass vacuously.
func seqLenMod(lenExpr, want string) string {
	return `mod m
  provides [m__go, M__Len]
  uses []
  emits []

` + seqLenOut + `
fn m__go() -> M__Len rev 1
  emits []
  tests
    go() => Ok(n = ` + want + `)
  Ok(n = ` + lenExpr + `)
`
}

func TestSeqLengthRows(t *testing.T) {
	rows := map[string][2]string{
		"zero":                {`#Seq<str>[]`, "0"},
		"singleton empty":     {`#Seq<str>[""]`, "1"},
		"two empties":         {`#Seq<str>["", ""]`, "2"},
		"multiscalar members": {`#Seq<str>["abc", "𝌆"]`, "2"},
		"repeats":             {`#Seq<str>["b", "", "a", "b"]`, "4"},
	}
	for name, row := range rows {
		t.Run(name, func(t *testing.T) {
			seqClean(t, map[string]string{"m.can": seqLenMod(row[0], row[1])}, "m.can")
		})
	}
}

// The rows above prove nothing if a miscount passes: off-by-one in
// both directions must fail the run.
func TestSeqLengthMiscountsFail(t *testing.T) {
	for _, want := range []string{"3", "5"} {
		bad := seqLenMod(`#Seq<str>["b", "", "a", "b"]`, want)
		seqCode(t, map[string]string{"m.can": bad}, "m.can", CodeTestFailed, "go")
	}
}

// Branded sequences count members without inspecting them.
func TestSeqLengthBranded(t *testing.T) {
	body := `mod m
  provides [m__go, M__B, M__Len]
  uses []
  emits []

brand M__B is str rev 1

` + seqLenOut + `
fn m__go() -> M__Len rev 1
  emits []
  tests
    go() => Ok(n = #Seq<M__B>[seal M__B("x"), seal M__B("")])
  Ok(n = #Seq<M__B>[seal M__B("x"), seal M__B("")])
`
	seqClean(t, map[string]string{"m.can": body}, "m.can")
}

// Length over a parameter, the S3 traversal-bound shape.
func TestSeqLengthParam(t *testing.T) {
	body := `mod m
  provides [m__go, M__Len]
  uses []
  emits []

` + seqLenOut + `
fn m__go(xs: Seq<str>) -> M__Len rev 1
  emits []
  tests
    empty(Seq<str>[]) => Ok(n = 0)
    two(Seq<str>["a", ""]) => Ok(n = 2)
  Ok(n = #xs)
`
	seqClean(t, map[string]string{"m.can": body}, "m.can")
}

// Non-sequence operands keep the pinned scalar diagnostic verbatim.
func TestSeqLengthNonSeq(t *testing.T) {
	body := seqLenMod(`#5`, "5")
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeTypeMismatch, "cannot count scalars of int")
}

// Emit lowers sequence length through the shared spread shape.
func TestSeqLengthEmit(t *testing.T) {
	ts := compileEmit(t, seqLenMod(`#Seq<str>["b", "a"]`, "2"))
	for _, want := range []string{".length", "bigint"} {
		if !strings.Contains(ts, want) {
			t.Fatalf("emit missing %q:\n%s", want, ts)
		}
	}
}
