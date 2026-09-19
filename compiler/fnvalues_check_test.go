package main

import (
	"testing"
)

// B00 stage 1a check: fnref creation parses but is refused with
// CAN6018 until invocation lands. The Fn annotation itself is
// accepted by shape so the report stays singular.
const fnvaluesCheckBody = `mod m
  provides [m__go, m__t, M__O]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__O rev 1 (
  value: str
)

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits [m.err]
  tests
    t(3, 7) => Ok("q")
  Ok("q")

fn m__go(n: int) -> M__O rev 1
  emits []
  tests
    g(3) => Ok("q")
  Ok(fnref m__t(divisor = 3))
`

func TestCheckFnrefDeferred(t *testing.T) {
	seqCode(t, map[string]string{"m.can": fnvaluesCheckBody}, "m.can",
		CodeFnValueDeferred, "deferred")
}

// A row-held reference reports the same code: the Fn annotation
// admits the argument by shape (no unknown-type cascade) while
// creation itself is refused. The row then fails evaluation
// (CAN4200), which seqCode admits as a cascade.
const fnvaluesCheckRow = `mod m
  provides [m__go, m__t, M__O]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__O rev 1 (
  value: str
)

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits [m.err]
  tests
    t(3, 7) => Ok("q")
  Ok("q")

fn m__go(cb: Fn<int, M__O, [m.err]>) -> M__O rev 1
  emits []
  tests
    g(fnref m__t(divisor = 3)) => Ok("q")
  Ok("q")
`

func TestCheckFnrefRowDeferred(t *testing.T) {
	seqCode(t, map[string]string{"m.can": fnvaluesCheckRow}, "m.can",
		CodeFnValueDeferred, "deferred")
}

// Invocation is refused with the same stage code: the row-held
// reference reports once, the match invoke node reports once, and
// nothing else fires (no value-table proof, no exhaustiveness, no
// given validation). The row then fails evaluation (CAN4200),
// which is admitted below like seqCode admits it.
const fnvaluesCheckInvoke = `mod m
  provides [m__go, m__t, M__O]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__O rev 1 (
  value: str
)

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits [m.err]
  tests
    t(3, 7) => Ok("q")
  Ok("q")

fn m__go(cb: Fn<int, M__O, [m.err]>, n: int) -> M__O rev 1
  emits [m.err]
  tests
    g(fnref m__t(divisor = 3), 4) => Ok("q")
  match invoke cb with n
    on Ok _ => Ok("q")
    on m.err _ => Ok("e")
`

func TestCheckInvokeDeferred(t *testing.T) {
	files := map[string]string{"m.can": fnvaluesCheckInvoke}
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, "m.can", files["m.can"])
	deferred := 0
	proof := 0
	for _, d := range diags {
		if d.Sev != "error" {
			continue
		}
		switch {
		case d.Code == CodeFnValueDeferred:
			deferred++
		case d.Code == CodeProofOther:
			proof++
		case d.Code == CodeTestFailed || d.Code == CodeInconsistentScript:
		default:
			t.Fatalf("unexpected cascade error, got %v", diags)
		}
	}
	if deferred != 2 {
		t.Fatalf("expected exactly two deferrals (row + node), got %v", diags)
	}
	if proof != 1 {
		t.Fatalf("expected exactly one proof fail-closed, got %v", diags)
	}
}
