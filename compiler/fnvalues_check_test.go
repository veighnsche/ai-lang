package main

import (
	"strings"
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

// Calls in invoke-argument position are outside-call violations:
// the argument executes off any scrutinee, so the v0 discipline
// applies exactly as elsewhere.
const fnvaluesCheckWithCall = `mod m
  provides [m__go, m__h, m__t, M__O]
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

fn m__h(x: int) -> M__O rev 1
  emits []
  tests
    h(1) => Ok("q")
  Ok("q")

fn m__go(cb: Fn<int, M__O, [m.err]>, n: int) -> M__O rev 1
  emits [m.err]
  tests
    g(fnref m__t(divisor = 3), 4) => Ok("q")
  match invoke cb with call m__h(n)
    on Ok _ => Ok("q")
    on m.err _ => Ok("e")
`

// wantFnHeadDiag asserts one error carries the code and fragment.
// The fnref row arg always adds its own CAN6018 deferral (plus a
// possible runtime confirmation), so only presence is asserted,
// never singularity.
func wantFnHeadDiag(t *testing.T, body, code, sub string) {
	t.Helper()
	dir := writeLSPDir(t, map[string]string{"m.can": body})
	diags := diagnose(dir, "m.can", body)
	for _, d := range diags {
		if d.Sev == "error" && d.Code == code && strings.Contains(d.Msg, sub) {
			return
		}
	}
	t.Fatalf("expected %s containing %q, got %v", code, sub, diags)
}

// B00 stage 1d check: deep head validation. The input must name a
// known type (CAN6002), the success must name a record (CAN6019),
// and kinds must name declared errors (CAN4002), distinctly and in
// canonical order (CAN6019).
const fnHeadBase = `mod m
  provides [m__go, m__t, M__O]
  uses []
  emits [m.err, m.odd]

error m.err(detail: str)

error m.odd(detail: str)

type M__O rev 1 (
  value: str
)

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits [m.err]
  tests
    t(3, 7) => Ok("q")
  Ok("q")

fn m__go(cb: HEAD) -> M__O rev 1
  emits []
  tests
    g(fnref m__t(divisor = 3)) => Ok("q")
  Ok("q")
`

func TestCheckFnHeadParams(t *testing.T) {
	cases := []struct {
		name string
		head string
		code string
		sub  string
	}{
		{"unknown input", `Fn<Nope, M__O, [m.err]>`, CodeUnknownType, "unknown type Nope in Fn input of param cb"},
		{"int success", `Fn<int, int, [m.err]>`, CodeFnHeadInvalid, "Fn success int of param cb is not a record"},
		{"kind success", `Fn<int, m.err, [m.err]>`, CodeFnHeadInvalid, "Fn success m.err of param cb is not a record"},
		{"unknown kind", `Fn<int, M__O, [m.nope]>`, CodeUnknownKind, "unknown error kind m.nope in Fn list of param cb"},
		{"duplicate kind", `Fn<int, M__O, [m.err, m.err]>`, CodeFnHeadInvalid, "duplicate error kind m.err in Fn list of param cb"},
		{"unordered kinds", `Fn<int, M__O, [m.odd, m.err]>`, CodeFnHeadInvalid, "not canonically ordered: want [m.err, m.odd]"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body := strings.Replace(fnHeadBase, "HEAD", c.head, 1)
			wantFnHeadDiag(t, body, c.code, c.sub)
		})
	}
}

// Bare-Fn returns are unsupported like bare-Seq returns: results
// name a record. An invalid head in return position still reports
// its own head diagnostic.
func TestCheckFnHeadReturns(t *testing.T) {
	bare := strings.Replace(fnHeadBase, "cb: HEAD", "n: int", 1)
	bare = strings.Replace(bare, "fn m__go(n: int) -> M__O rev 1",
		"fn m__go(n: int) -> Fn<int, M__O, [m.err]> rev 1", 1)
	bare = strings.Replace(bare, "g(fnref m__t(divisor = 3))", "g(3)", 1)
	wantFnHeadDiag(t, bare, CodeTypeMismatch, "bare-Fn returns are unsupported, return a record")

	bad := strings.Replace(bare,
		"-> Fn<int, M__O, [m.err]> rev 1", "-> Fn<int, int, [m.err]> rev 1", 1)
	wantFnHeadDiag(t, bad, CodeFnHeadInvalid, "Fn success int of returns is not a record")
}

// Record fields carry callables, so field heads are deep-checked
// exactly like params.
func TestCheckFnHeadField(t *testing.T) {
	body := strings.Replace(fnHeadBase, "provides [m__go, m__t, M__O]",
		"provides [m__go, m__t, M__O, M__Box]", 1)
	body = strings.Replace(body, "type M__O rev 1 (",
		"type M__Box rev 1 (\n  cb: Fn<int, M__O, [m.nope]>\n)\n\ntype M__O rev 1 (", 1)
	body = strings.Replace(body, "HEAD", `Fn<int, M__O, [m.err]>`, 1)
	wantFnHeadDiag(t, body, CodeUnknownKind, "unknown error kind m.nope in Fn list of field cb")
}

// Extern signatures are deep-checked too; an extern returning a
// bare Fn trips the existing record-return rule.
func TestCheckFnHeadExtern(t *testing.T) {
	body := `mod m
  provides [m__go, m__t, M__O, M__Verdict, m__use]
  uses []
  emits [m.err, m.odd]

error m.err(detail: str)

error m.odd(detail: str)

type M__O rev 1 (
  value: str
)

type M__Verdict rev 1 (
)

extern m__use(cb: Fn<int, M__O, [m.nope]>) -> M__Verdict rev 1
  emits [m.err]

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
	wantFnHeadDiag(t, body, CodeUnknownKind, "unknown error kind m.nope in Fn list of param cb")

	ret := strings.Replace(body,
		"extern m__use(cb: Fn<int, M__O, [m.nope]>) -> M__Verdict rev 1",
		"extern m__use(n: int) -> Fn<int, M__O, [m.err]> rev 1", 1)
	wantFnHeadDiag(t, ret, CodeTypeMismatch, "externs return a record type")
}

// A valid multi-kind head with an empty-list field reports no head
// diagnostic: only the row's CAN6018 deferral (and its runtime
// confirmation) may fire.
func TestCheckFnHeadValid(t *testing.T) {
	body := strings.Replace(fnHeadBase, "provides [m__go, m__t, M__O]",
		"provides [m__go, m__t, M__O, M__Box]", 1)
	body = strings.Replace(body, "type M__O rev 1 (",
		"type M__Box rev 1 (\n  cb: Fn<int, M__O, []>\n)\n\ntype M__O rev 1 (", 1)
	body = strings.Replace(body, "HEAD", `Fn<int, M__O, [m.err, m.odd]>`, 1)
	dir := writeLSPDir(t, map[string]string{"m.can": body})
	diags := diagnose(dir, "m.can", body)
	for _, d := range diags {
		if d.Sev != "error" {
			continue
		}
		switch d.Code {
		case CodeFnHeadInvalid, CodeUnknownType, CodeUnknownKind:
			t.Fatalf("unexpected head diagnostic, got %v", diags)
		}
	}
}

func TestCheckInvokeWithCall(t *testing.T) {
	files := map[string]string{"m.can": fnvaluesCheckWithCall}
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, "m.can", files["m.can"])
	deferred, outside, proof := 0, 0, 0
	for _, d := range diags {
		if d.Sev != "error" {
			continue
		}
		switch {
		case d.Code == CodeFnValueDeferred:
			deferred++
		case d.Code == CodeCallOutside:
			outside++
		case d.Code == CodeProofOther:
			proof++
		case d.Code == CodeTestFailed || d.Code == CodeInconsistentScript:
		default:
			t.Fatalf("unexpected cascade error, got %v", diags)
		}
	}
	if deferred != 2 || outside != 1 || proof != 1 {
		t.Fatalf("want 2 deferrals + outside call + proof guard, got %v", diags)
	}
}
