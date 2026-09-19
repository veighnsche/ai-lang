package main

import (
	"strings"
	"testing"
)

// I3 runtime: references create callable values, invocation
// dispatches them against real target bodies, and factory rows
// compare constructions. These tests assert execution, so they
// demand zero error diagnostics (warnings stay welcome).

func wantInvokeClean(t *testing.T, files map[string]string, name string) {
	t.Helper()
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, name, files[name])
	for _, d := range diags {
		if d.Sev == "error" {
			t.Fatalf("expected clean execution, got %v", diags)
		}
	}
}

// Two rows bind different captures against the same target and
// both execute: creation binds per-row values, dispatch runs the
// real body with captures plus input, and the binder threads the
// computed payload. The bind-three/bind-five shape catches a
// live-environment capture implementation.
const invokeExecBody = `mod m
  provides [m__go, m__t, M__O]
  uses []
  emits []

type M__O rev 1 (
  value: int
)

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits []
  tests
    t(3, 7) => Ok(4)
    t2(2, 9) => Ok(7)
  Ok(dividend - divisor)

fn m__go(cb: Fn<int, M__O, []>, n: int) -> M__O rev 1
  emits []
  tests
    g(fnref m__t(divisor = 3), 7) => Ok(4)
    h(fnref m__t(divisor = 2), 9) => Ok(7)
  match invoke cb with n
    on Ok r => Ok(r.value)
`

func TestInvokeExecutes(t *testing.T) {
	wantInvokeClean(t, map[string]string{"m.can": invokeExecBody}, "m.can")
}

// Error dispatch takes the error arm with no relay exemption:
// both arms execute across the run, so coverage holds.
const invokeErrBody = `mod m
  provides [m__go, m__t, M__O]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__O rev 1 (
  value: int
)

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits [m.err]
  tests
    t(3, 7) => Ok(4)
    tz(3, 0) => m.err("zero")
  match dividend
    on 0 => m.err("zero")
    on _ => Ok(dividend - divisor)

fn m__go(cb: Fn<int, M__O, [m.err]>, n: int) -> M__O rev 1
  emits []
  tests
    g(fnref m__t(divisor = 3), 7) => Ok(4)
    gz(fnref m__t(divisor = 3), 0) => Ok(0)
  match invoke cb with n
    on Ok r => Ok(r.value)
    on m.err _ => Ok(0)
`

func TestInvokeErrorArm(t *testing.T) {
	wantInvokeClean(t, map[string]string{"m.can": invokeErrBody}, "m.can")
}

// A factory row compares constructions: the created reference
// equals the expected one by target, signature, slot, and
// captured values, so the row passes. The zero row constructs
// without invoking: bind-zero succeeds, invoke-zero fails.
const invokeFactoryBody = `mod m
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
    gz(0) => Ok(fnref m__t(divisor = 0))
  Ok(fnref m__t(divisor = n))
`

func TestInvokeFactoryRow(t *testing.T) {
	wantInvokeClean(t, map[string]string{"m.can": invokeFactoryBody}, "m.can")
}

// A factory mismatch names the capture but redacts its value:
// payloads never render into diagnostics.
func TestInvokeFactoryMismatchRedacted(t *testing.T) {
	body := strings.Replace(invokeFactoryBody,
		"g(3) => Ok(fnref m__t(divisor = 3))",
		"g(3) => Ok(fnref m__t(divisor = 4))", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": body})
	diags := diagnose(dir, "m.can", body)
	found := false
	for _, d := range diags {
		if d.Sev != "error" {
			continue
		}
		if strings.Contains(d.Msg, "callable capture divisor differs (value redacted)") {
			found = true
		}
		if strings.Contains(d.Msg, "divisor = ") {
			t.Fatalf("mismatch must redact captured payloads, got %v", diags)
		}
	}
	if !found {
		t.Fatalf("expected redacted capture mismatch, got %v", diags)
	}
}

// A foreign invocation runs the real body with linked-pure
// dispatch and no scripts: the target's own foreign call uses
// its given table in its own row, but the callback subtree
// executes the provider body directly.
const invokeForeignProvU = `mod u
  provides [u__inc, U__O]
  uses []
  emits []

type U__O rev 1 (
  value: int
)

fn u__inc(x: int) -> U__O rev 1
  emits []
  tests
    i(1) => Ok(2)
  Ok(x + 1)
`

const invokeForeignProvT = `mod t
  provides [t__dbl, T__O]
  uses [u__inc@1]
  emits []

type T__O rev 1 (
  value: int
)

fn t__dbl(a: int, b: int) -> T__O rev 1
  emits []
  tests
    d(3, 4) => Ok(8)
  match call u__inc(a)
    given
      d => [exchange args (x = 3) outcome Ok(4)]
    on Ok r => Ok(r.value + b)
`

const invokeForeignUse = `mod m
  provides [m__go]
  uses [t__dbl@1]
  emits []

fn m__go(cb: Fn<int, T__O, []>, n: int) -> T__O rev 1
  emits []
  tests
    g(fnref t__dbl(a = 3), 4) => Ok(8)
  match invoke cb with n
    on Ok r => Ok(r.value)
`

func TestInvokeForeignNoScripts(t *testing.T) {
	files := map[string]string{
		"u.can": invokeForeignProvU,
		"t.can": invokeForeignProvT,
		"m.can": invokeForeignUse,
	}
	wantInvokeClean(t, files, "m.can")
}

// A scripted foreign higher-order call compares its callback
// argument by construction: the exchange expectation holds a
// reference, and the call supplies an equal one.
const invokeExchangeProv = `mod p
  provides [p__use, p__id, P__O]
  uses []
  emits []

type P__O rev 1 (
  value: str
)

fn p__id(a: int, b: int) -> P__O rev 1
  emits []
  tests
    i(1, 2) => Ok("q")
  Ok("q")

fn p__use(cb: Fn<int, P__O, []>, n: int) -> P__O rev 1
  emits []
  tests
    u(fnref p__id(a = 1), 1) => Ok("q")
  Ok("q")
`

const invokeExchangeUse = `mod m
  provides [m__go, m__t]
  uses [p__use@1]
  emits []

fn m__t(divisor: int, dividend: int) -> P__O rev 1
  emits []
  tests
    t(3, 7) => Ok("q")
  Ok("q")

fn m__go(n: int) -> P__O rev 1
  emits []
  tests
    g(4) => Ok("q")
  match call p__use(fnref m__t(divisor = 3), n)
    given
      g => [exchange args (cb = fnref m__t(divisor = 3), n = 4) outcome Ok("q")]
    on Ok r => Ok(r.value)
`

func TestInvokeExchangeFnArg(t *testing.T) {
	files := map[string]string{
		"p.can": invokeExchangeProv,
		"m.can": invokeExchangeUse,
	}
	wantInvokeClean(t, files, "m.can")
}

// A generic target becomes a value through explicit
// specialization: the row stamps the reference, and creation,
// dispatch, and comparison all run against the stamp.
const invokeGenericBody = `mod m
  provides [m__go, m__pick, M__O]
  uses []
  emits []

type M__O rev 1 (
  value: int
)

fn m__pick<T>(tag: T, n: int) -> M__O rev 1
  emits []
  tests
    pk<T=str>("s", 7) => Ok(7)
  Ok(n)

fn m__go(cb: Fn<int, M__O, []>, n: int) -> M__O rev 1
  emits []
  tests
    g(fnref m__pick<str>(tag = "s"), 7) => Ok(7)
  match invoke cb with n
    on Ok r => Ok(r.value)
`

func TestInvokeGenericTarget(t *testing.T) {
	wantInvokeClean(t, map[string]string{"m.can": invokeGenericBody}, "m.can")
}

// A pure target with a proved direct self-decrease stays
// admissible: direct recursion is not an invocation cycle, and
// the row below executes three levels through the reference.
const invokeRecursiveBody = `mod m
  provides [m__go, m__sum, M__O]
  uses []
  emits []

type M__O rev 1 (
  value: int
)

fn m__sum(step: int, n: int) -> M__O rev 1
  decreases n
  emits []
  tests
    s0(10, 0) => Ok(0)
    s1(10, 1) => Ok(10)
  match n <= 0
    true => Ok(0)
    false => match call m__sum(step, n - 1)
      on Ok r => Ok(r.value + step)

fn m__go(cb: Fn<int, M__O, []>) -> M__O rev 1
  emits []
  tests
    g(fnref m__sum(step = 10)) => Ok(30)
  match invoke cb with 3
    on Ok r => Ok(r.value)
`

func TestInvokeRecursiveTarget(t *testing.T) {
	wantInvokeClean(t, map[string]string{"m.can": invokeRecursiveBody}, "m.can")
}
