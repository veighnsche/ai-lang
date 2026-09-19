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
  Ok(fnref m__t(divisor = n))
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

// The proving invocation: the row-held reference checks against
// the callback parameter, the argument matches the Fn input, and
// the arms cover Ok plus the declared error set. Only the row's
// runtime confirmation remains until I3 evaluates invocations.
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
    tz(3, 0) => m.err("zero")
  match dividend
    on 0 => m.err("zero")
    on _ => Ok("q")

fn m__go(cb: Fn<int, M__O, [m.err]>, n: int) -> M__O rev 1
  emits [m.err]
  tests
    g(fnref m__t(divisor = 3), 4) => Ok("q")
    gz(fnref m__t(divisor = 3), 0) => Ok("e")
  match invoke cb with n
    on Ok _ => Ok("q")
    on m.err _ => Ok("e")
`

// I2: the invocation proves. The reference and argument check
// against the Fn signature, the arms cover Ok plus the declared
// error set, and only the row's runtime confirmation remains
// (evaluation lands in I3).
func TestCheckInvokeProves(t *testing.T) {
	wantFnvaluesClean(t, map[string]string{"m.can": fnvaluesCheckInvoke}, "m.can")
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
// Only presence is asserted, never singularity: rows and arm
// bodies may add their own failures around the head under test.
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
		{"arithmetic capture is computed", `g(fnref m__t(divisor = 1 + 2)) => Ok("q")`,
			CodeFnComputedCapture, "is computed (binop)"},
		{"captures follow declaration order", `g(fnref m__t(dividend = 4, divisor = 3)) => Ok("q")`,
			CodeBadBinding, "out of order"},
		{"capture value shares unbound-name rule", `g(fnref m__t(divisor = nope)) => Ok("q")`,
			CodeTypeMismatch, "unbound name nope in m__go"},
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
	seqCode(t, map[string]string{"m.can": fnvaluesCheckWithCall}, "m.can",
		CodeCallOutside, "outside a match scrutinee")
}

// I2 negatives: each invocation rule reports its own reused code
// with a singular report (seqCode admits only the runtime
// confirmation). The proof stays silent on unresolvable targets
// (a62): the checker owns those, so no CAN4100 joins them.
const fnvaluesInvokeBase = `mod m
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
    tz(3, 0) => m.err("zero")
  match dividend
    on 0 => m.err("zero")
    on _ => Ok("q")

fn m__go(cb: Fn<int, M__O, [m.err]>, n: int) -> M__O rev 1
  emits [m.err]
  tests
    g(fnref m__t(divisor = 3), 4) => Ok("q")
    gz(fnref m__t(divisor = 3), 0) => Ok("e")
MATCH
`

func TestCheckInvokeNegatives(t *testing.T) {
	cases := []struct {
		name  string
		match string
		code  string
		sub   string
	}{
		{"non-callable target", `  match invoke n with n
    on Ok _ => Ok("q")
    on m.err _ => Ok("e")`,
			CodeTypeMismatch, "invoke target n has type int: want a function value"},
		{"unknown target", `  match invoke nope with n
    on Ok _ => Ok("q")
    on m.err _ => Ok("e")`,
			CodeTypeMismatch, "unbound name nope in m__go"},
		{"argument mismatch", `  match invoke cb with "s"
    on Ok _ => Ok("q")
    on m.err _ => Ok("e")`,
			CodeTypeMismatch, "invoke cb argument: got str, want int"},
		{"missing arm", `  match invoke cb with n
    on Ok _ => Ok("q")`,
			CodeMissingArm, "non-exhaustive match, missing m.err"},
		{"given scripts nothing", `  match invoke cb with n
    given
      g => [exchange args (n = 4) outcome Ok("q")]
    on Ok _ => Ok("q")
    on m.err _ => Ok("e")`,
			CodeGivenOnInvoke, "takes no given table"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body := strings.Replace(fnvaluesInvokeBase, "MATCH", c.match, 1)
			seqCode(t, map[string]string{"m.can": body}, "m.can", c.code, c.sub)
		})
	}
}

// A stale arm is presence-asserted like every other stale arm:
// the untakeable arm also trips coverage, so the report is never
// singular.
func TestCheckInvokeStaleArm(t *testing.T) {
	body := strings.Replace(fnvaluesInvokeBase, "MATCH", `  match invoke cb with n
    on Ok _ => Ok("q")
    on m.err _ => Ok("e")
    on m.odd _ => Ok("o")`, 1)
	dir := writeLSPDir(t, map[string]string{"m.can": body})
	diags := diagnose(dir, "m.can", body)
	if !hasErrCode(diags, CodeStaleArm) || !hasDiag(diags, "error", "stale match arm m.odd") {
		t.Fatalf("expected stale m.odd arm, got %v", diags)
	}
}

// A binder shadowing the callback parameter is not callable: the
// checker reports the shadowed type while the proof still reads
// the parameter's signature (exhaustive here, so silent).
func TestCheckInvokeShadowedTarget(t *testing.T) {
	body := strings.Replace(fnvaluesInvokeBase, "MATCH", `  match invoke cb with n
    on Ok _ => Ok("q")
    on m.err cb => match invoke cb with n
      on Ok _ => Ok("q")
      on m.err _ => Ok("e")`, 1)
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeTypeMismatch, "invoke target cb has type m.err: want a function value")
}

// Ok forwards elaborate in invoke arms exactly like call arms:
// the success record comes off the resolved signature.
func TestCheckInvokeForwardOk(t *testing.T) {
	body := strings.Replace(fnvaluesInvokeBase, "MATCH", `  match invoke cb with n
    on Ok r => forward r
    on m.err _ => Ok("e")`, 1)
	wantFnvaluesClean(t, map[string]string{"m.can": body}, "m.can")
}

// A mistargeted forward in an invoke arm is CAN3011, like
// everywhere else.
func TestCheckInvokeForwardRejects(t *testing.T) {
	body := strings.Replace(fnvaluesInvokeBase, "MATCH", `  match invoke cb with n
    on Ok r => forward q
    on m.err _ => Ok("e")`, 1)
	dir := writeLSPDir(t, map[string]string{"m.can": body})
	if diags := diagnose(dir, "m.can", body); !hasCode(diags, "CAN3011") {
		t.Fatalf("mistargeted invoke forward reported no CAN3011: %v", diags)
	}
}

// I2b: an invocation-closed cycle is refused. A invokes its
// callback, the rows take addresses of B for that signature, B
// calls C, and C calls A back: every static check passes, but
// the loop through the function value cannot terminate, so
// CAN3005 reports at the invoke site and execution stays gated.
const fnvaluesInvokeCycle = `mod m
  provides [m__a, m__b, m__c, M__O]
  uses []
  emits []

type M__O rev 1 (
  value: int
)

fn m__b(a: int, b: int) -> M__O rev 1
  emits []
  tests
    b(1, 2) => Ok(3)
  match call m__c(a)
    on Ok r => Ok(r.value + b)

fn m__c(x: int) -> M__O rev 1
  emits []
  tests
    c(1) => Ok(1)
  match call m__a(fnref m__b(a = 1), x)
    on Ok r => Ok(r.value)

fn m__a(cb: Fn<int, M__O, []>, n: int) -> M__O rev 1
  emits []
  tests
    a(fnref m__b(a = 1), 2) => Ok(3)
  match invoke cb with n
    on Ok r => Ok(r.value)
`

func TestCheckInvokeCycle(t *testing.T) {
	seqCode(t, map[string]string{"m.can": fnvaluesInvokeCycle}, "m.can",
		CodeLocalCycle, "m__a invokes m__b through a function value; invocation cannot close a cycle")
}

// One signature written tight and one written loose still
// match: Fn heads compare structurally, never byte-identical.
func TestCheckFnrefTightSpacing(t *testing.T) {
	body := strings.Replace(fnvaluesCheckBody, "cb: Fn<int, M__O, [m.err]>", "cb: Fn<int,M__O,[m.err]>", 1)
	wantFnvaluesClean(t, map[string]string{"m.can": body}, "m.can")
}

// Source-level field access on a callable fails through the
// ordinary projection rule: a function value has no fields.
func TestCheckInvokeFieldRefused(t *testing.T) {
	body := strings.Replace(fnvaluesInvokeBase, "MATCH", `  match invoke cb with n
    on Ok r => Ok(cb.value)
    on m.err _ => Ok("e")`, 1)
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeTypeMismatch, "no field value on cb")
}

// Only admissible addresses close cycles: a target with a
// required precondition contributes no edge, so the same shape
// reports the target refusal without any cycle.
func TestCheckInvokeCycleRequiresTargetClean(t *testing.T) {
	body := strings.Replace(fnvaluesInvokeCycle,
		"fn m__b(a: int, b: int) -> M__O rev 1\n  emits []",
		"fn m__b(a: int, b: int) -> M__O rev 1\n  emits []\n  requires\n    a == 1", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": body})
	diags := diagnose(dir, "m.can", body)
	if !hasCode(diags, CodeFnTargetRefused) {
		t.Fatalf("requires target reported no CAN6023: %v", diags)
	}
	for _, d := range diags {
		if d.Sev == "error" && d.Code == CodeLocalCycle {
			t.Fatalf("inadmissible address must not close a cycle, got %v", diags)
		}
	}
}
