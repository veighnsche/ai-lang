package main

import (
	"strings"
	"testing"
)

// B00 invocation slice I1: well-formed reference creation checks
// clean. The factory returns a record holding the callable
// (bare-Fn returns stay refused); the row then fails evaluation
// (CAN4200) until the I3 runtime constructs values, which the
// assertions below admit explicitly.
const fnvaluesCheckBody = `mod m
  provides [m__go, m__t, M__O, M__H]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__O rev 1 (
  value: str
)

type M__H rev 1 (
  cb: Fn<int, M__O, [m.err]>
)

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits [m.err]
  tests
    t(3, 7) => Ok("q")
  Ok("q")

fn m__go(n: int) -> M__H rev 1
  emits []
  tests
    g(3) => Ok(fnref m__t(divisor = 3))
  M__H(fnref m__t(divisor = n))
`

// wantFnvaluesClean asserts static acceptance: every error is the
// admitted runtime confirmation (the row cannot evaluate until
// the I3 runtime lands), never a static rejection.
func wantFnvaluesClean(t *testing.T, files map[string]string, name string) {
	t.Helper()
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, name, files[name])
	for _, d := range diags {
		if d.Sev != "error" {
			continue
		}
		if d.Code == CodeTestFailed || d.Code == CodeInconsistentScript {
			continue
		}
		t.Fatalf("expected static acceptance, got %v", diags)
	}
}

func TestCheckFnrefFactory(t *testing.T) {
	wantFnvaluesClean(t, map[string]string{"m.can": fnvaluesCheckBody}, "m.can")
}

// A row-held reference checks clean too: the denoted Fn type
// matches the callback parameter, and only evaluation waits on
// the runtime.
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

func TestCheckFnrefRow(t *testing.T) {
	wantFnvaluesClean(t, map[string]string{"m.can": fnvaluesCheckRow}, "m.can")
}

// Invocation stays deferred in I1: the row-held reference now
// checks clean (only its CAN4200 evaluation failure remains,
// admitted below), the match invoke node reports once, and
// nothing else fires (no value-table proof, no exhaustiveness,
// no given validation).
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
	if deferred != 1 {
		t.Fatalf("expected exactly one deferral (the invoke node), got %v", diags)
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
// The fnref row arg always adds its own CAN4200 runtime
// confirmation (references check but do not evaluate yet), so
// only presence is asserted, never singularity.
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
// diagnostic: only the row's runtime confirmation may fire.
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

// I1 negatives: each creation rule reports its own code with a
// singular report (seqCode admits only the runtime confirmation).
const fnvaluesNegBase = `mod m
  provides [m__go, m__t, m__h, M__O]
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

fn m__go(cb: Fn<int, M__O, [m.err]>) -> M__O rev 1
  emits []
  tests
    ROW
  Ok("q")
`

func TestCheckFnrefNegatives(t *testing.T) {
	cases := []struct {
		name string
		row  string
		code string
		sub  string
	}{
		{"empty leaves two unbound", `g(fnref m__t()) => Ok("q")`,
			CodeFnResidualArity, "leaves 2 parameters unbound"},
		{"full binding is a call", `g(fnref m__t(divisor = 3, dividend = 4)) => Ok("q")`,
			CodeFnResidualArity, "binds every parameter"},
		{"call capture is computed", `g(fnref m__t(divisor = call m__h(3))) => Ok("q")`,
			CodeFnComputedCapture, "is computed (call)"},
		{"captures follow declaration order", `g(fnref m__t(dividend = 4, divisor = 3)) => Ok("q")`,
			CodeBadBinding, "out of order"},
		{"unknown target", `g(fnref m__nope(divisor = 3)) => Ok("q")`,
			CodeUnknownCall, "references unknown function"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body := strings.Replace(fnvaluesNegBase, "ROW", c.row, 1)
			seqCode(t, map[string]string{"m.can": body}, "m.can", c.code, c.sub)
		})
	}
}

// A capture of callback type is refused (CAN6022) alongside the
// slot mismatch: captures are data-only.
func TestCheckFnrefCaptureContainsFn(t *testing.T) {
	body := `mod m
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
  emits []
  tests
    g(fnref m__t(divisor = cb), 4) => Ok("q")
  Ok("q")
`
	wantFnHeadDiag(t, body, CodeFnContainment, "contains a function value: captures are data-only")
}

// Externs are not referenceable: only source functions take
// addresses.
func TestCheckFnrefExternRefused(t *testing.T) {
	body := `mod m
  provides [m__go, m__t, m__use, M__O, M__Verdict]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__O rev 1 (
  value: str
)

type M__Verdict rev 1 (
)

extern m__use(n: int) -> M__Verdict rev 1
  emits []

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits [m.err]
  tests
    t(3, 7) => Ok("q")
  Ok("q")

fn m__go(cb: Fn<int, M__O, [m.err]>) -> M__O rev 1
  emits []
  tests
    g(fnref m__use(n = 3)) => Ok("q")
  Ok("q")
`
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeFnTargetRefused, "reference to extern m__use refused")
}

// An effectful target is refused through the whole-graph linked
// criterion, even though the target itself checks clean. The
// zero-capture reference against the unary target also proves the
// residual rule needs no captures when one parameter is already
// the whole input.
func TestCheckFnrefEffectfulRefused(t *testing.T) {
	body := `mod m
  provides [m__go, m__bump, M__T]
  uses []
  emits []

state M__C: int = 0

type M__T rev 1 (
  total: int
)

fn m__bump(by: int) -> M__T rev 1
  effects [M__C.read, M__C.write]
  emits []
  tests
    three(3) => Ok(3)
  match call state__get(M__C)
    on Ok c => match call state__put(M__C, c.value + by)
      on Ok _ => Ok(c.value + by)

fn m__go(cb: Fn<int, M__T, []>) -> M__T rev 1
  emits []
  tests
    g(fnref m__bump()) => Ok(3)
  Ok(0)
`
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeFnTargetRefused, "reference to m__bump refused: effectful function")
}

// A target with a required precondition is refused: captures
// never discharge preconditions.
func TestCheckFnrefRequiresRefused(t *testing.T) {
	body := `mod m
  provides [m__go, m__pick, M__O]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__O rev 1 (
  value: int
)

fn m__pick(a: int, b: int) -> M__O rev 1
  emits []
  requires
    a == 3
  tests
    v(3, 5) => Ok(8)
  Ok(a + b)

fn m__go(cb: Fn<int, M__O, []>) -> M__O rev 1
  emits []
  tests
    g(fnref m__pick(a = 3)) => Ok(8)
  Ok(0)
`
	wantFnHeadDiag(t, body, CodeFnTargetRefused, "has a required precondition")
}

// Q3a: a foreign reference needs the creator's own uses pin,
// exactly like a foreign call.
const fnvaluesForeignProv = `mod t
  provides [t__dbl, T__O]
  uses []
  emits []

type T__O rev 1 (
  value: int
)

fn t__dbl(a: int, b: int) -> T__O rev 1
  emits []
  tests
    d(3, 4) => Ok(7)
  Ok(a + b)
`

const fnvaluesForeignBase = `mod m
  provides [m__go]
  uses [PINS]
  emits []

fn m__go(cb: Fn<int, T__O, []>) -> T__O rev 1
  emits []
  tests
    g(fnref t__dbl(a = 3)) => Ok(7)
  Ok(0)
`

func TestCheckFnrefForeignNeedsPin(t *testing.T) {
	body := strings.Replace(fnvaluesForeignBase, "PINS", "", 1)
	files := map[string]string{"t.can": fnvaluesForeignProv, "m.can": body}
	seqCode(t, files, "m.can", CodeCallNotInUses, "references t__dbl which is not in uses")
}

// The pinned foreign reference checks clean, and the address
// counts as use: no unused-pin warning fires even though no call
// ever reaches the target.
func TestCheckFnrefForeignPinClean(t *testing.T) {
	body := strings.Replace(fnvaluesForeignBase, "PINS", "t__dbl@1", 1)
	files := map[string]string{"t.can": fnvaluesForeignProv, "m.can": body}
	wantFnvaluesClean(t, files, "m.can")
	dir := writeLSPDir(t, files)
	for _, d := range diagnose(dir, "m.can", body) {
		if d.Code == CodeUnusedUses {
			t.Fatalf("row-held address should count as use, got %v", d)
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
	if deferred != 1 || outside != 1 || proof != 1 {
		t.Fatalf("want 1 deferral + outside call + proof guard, got %v", diags)
	}
}
