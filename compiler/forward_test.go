package main

import (
	"os"
	"strings"
	"testing"
)

// Slice 2: forward arms. Probes first: check-time elaboration of
// `forward v` into the complete constructor, CAN3011 rejections,
// the CAN4107 certificate boundary, and the manual-trap guard.

const forwardLib = `mod m
  provides [m__go, m__work, M__Out, M__Work]
  uses []
  emits [m.bad]

error m.bad(value: str)

type M__Out rev 1 (
  value: str
)

type M__Work rev 1 (
  value: str
)

fn m__work(v: str) -> M__Work rev 1
  emits [m.bad]
  tests
    w("x") => Ok("x")
    u("b") => m.bad("b")
  match v == "b"
    true => m.bad(v)
    false => Ok(v)
`

const forwardGoHead = `
fn m__go(x: str) -> M__Out rev 1
  emits [m.bad]
  tests
    ok("a") => Ok("a")
    bad("b") => m.bad("b")
  match call m__work(x)
`

// TestForwardErrorElaborates pins the core rewrite: an error arm
// forwarding its binder behaves exactly like the handwritten
// reconstruction.
func TestForwardErrorElaborates(t *testing.T) {
	src := forwardLib + forwardGoHead + `    on m.bad e => forward e
    on Ok r => Ok(r.value)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("forwarded error arm reported: %v", diags)
	}
}

// TestForwardOkElaborates pins the Ok rewrite: same field set on
// both records reconstructs explicitly.
func TestForwardOkElaborates(t *testing.T) {
	src := forwardLib + forwardGoHead + `    on m.bad e => m.bad(e.value)
    on Ok r => forward r
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	// The forwarded Ok arm stays silent; the handwritten m.bad relay
	// alongside it is a lint error (CAN3415), and nothing else reports.
	relay := 0
	for _, d := range diags {
		if d.Code == CodeLintRelay {
			relay++
			continue
		}
		if d.Sev == "error" {
			t.Fatalf("unexpected error beside the relay finding: %v", diags)
		}
	}
	if relay != 1 {
		t.Fatalf("expected exactly one CAN3415 relay finding, got %v", diags)
	}
}

// TestForwardRejects pins CAN3011: only `forward <arm-binder>`
// elaborates. Underscores, outer binders, projections,
// bare forward, constructor operands, wild patterns, and value
// matches are all rejected.
func TestForwardRejects(t *testing.T) {
	badArm := []string{
		"on m.bad e => forward _",
		"on m.bad e => forward x",
		"on m.bad e => forward e.value",
		"on m.bad e => forward",
		"on m.bad e => forward m.bad(e.value)",
		"on m.bad _ => forward _",
	}
	for _, arm := range badArm {
		src := forwardLib + forwardGoHead + "    " + arm + "\n    on Ok r => Ok(r.value)\n"
		dir := writeLSPDir(t, map[string]string{"m.can": src})
		if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN3011") {
			t.Fatalf("%q reported no CAN3011: %v", arm, diags)
		}
	}
	valueMatch := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: str
)

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    one(1) => Ok("yes")
    two(2) => Ok("no")
  match x <= 1
    true => forward x
    false => Ok("no")
`
	dir := writeLSPDir(t, map[string]string{"m.can": valueMatch})
	if diags := diagnose(dir, "m.can", valueMatch); !hasCode(diags, "CAN3011") {
		t.Fatalf("value-match forward reported no CAN3011: %v", diags)
	}
}

// TestForwardOkMismatch pins CAN3011 for shape drift: the source
// payload must hold exactly the destination success fields.
func TestForwardOkMismatch(t *testing.T) {
	src := `mod m
  provides [m__go, m__work, M__Out, M__Work]
  uses []
  emits []

type M__Out rev 1 (
  other: str
)

type M__Work rev 1 (
  value: str
)

fn m__work(v: str) -> M__Work rev 1
  emits []
  tests
    w("x") => Ok("x")
  Ok(v)

fn m__go(x: str) -> M__Out rev 1
  emits []
  tests
    ok("a") => Ok("a")
  match call m__work(x)
    given
      ok => [Ok("a")]
    on Ok r => forward r
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN3011") {
		t.Fatalf("mismatched Ok forward reported no CAN3011: %v", diags)
	}
}

// TestForwardUnknownCallee pins fail-closed elaboration: a forward
// on a call nobody declares is CAN3011, never a downstream crash
// on a missing shape.
func TestForwardUnknownCallee(t *testing.T) {
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
    ok("a") => Ok("a")
  match call nope__missing(x)
    on Ok r => forward r
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN3011") {
		t.Fatalf("unknown-callee forward reported no CAN3011: %v", diags)
	}
}

// TestForwardExternOk pins Ok forwarding across an extern
// boundary: the callee's declared return record resolves the
// same way as a local one.
func TestForwardExternOk(t *testing.T) {
	src := `mod m
  provides [m__go, ex__work, M__Out, M__Work]
  uses []
  emits []

extern ex__work(v: str) -> M__Work rev 1

type M__Out rev 1 (
  value: str
)

type M__Work rev 1 (
  value: str
)

fn m__go(x: str) -> M__Out rev 1
  emits []
  tests
    ok("a") => Ok("a")
  match call ex__work(x)
    given
      ok => [exchange args (v = "a") outcome Ok("a")]
    on Ok r => forward r
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("extern Ok forward reported: %v", diags)
	}
}

// TestForwardSameEmitsCheck pins that elaborated arms face the
// handwritten checks: a forwarded kind outside the caller's
// emits is CAN4001, exactly like spelling the constructor out.
func TestForwardSameEmitsCheck(t *testing.T) {
	src := forwardLib + `
fn m__go(x: str) -> M__Out rev 1
  emits []
  tests
    ok("a") => Ok("a")
    bad("b") => m.bad("b")
  match call m__work(x)
    on m.bad e => forward e
    on Ok r => Ok(r.value)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN4001") {
		t.Fatalf("forward outside emits reported no CAN4001: %v", diags)
	}
}

// TestForwardCertificate pins the unchanged CAN4107 boundary: an
// untaken forwarded error arm of a local call carries the relay
// certificate instead of an execution failure.
func TestForwardCertificate(t *testing.T) {
	src := forwardLib + `
fn m__go(x: str) -> M__Out rev 1
  emits [m.bad]
  tests
    ok("a") => Ok("a")
  match call m__work(x)
    on m.bad e => forward e
    on Ok r => Ok(r.value)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("certified forward relay reported: %v", diags)
	}
}

// TestForwardOkUntaken pins the execution law for Ok: a forwarded
// Ok arm no test takes is CAN4107, never certified.
func TestForwardOkUntaken(t *testing.T) {
	src := forwardLib + `
fn m__go(x: str) -> M__Out rev 1
  emits [m.bad]
  tests
    bad("b") => m.bad("b")
  match call m__work(x)
    on m.bad e => m.bad(e.value)
    on Ok r => forward r
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN4107") {
		t.Fatalf("untaken Ok forward reported no CAN4107: %v", diags)
	}
}

// TestForwardForeignNoCertificate pins the locality boundary: a
// forwarded error arm of a foreign call acquires no certificate,
// so untaken it is CAN4107.
func TestForwardForeignNoCertificate(t *testing.T) {
	lib := `mod lib
  provides [lib__work, L__Work]
  uses []
  emits [lib.bad]

error lib.bad(value: str)

type L__Work rev 1 (
  value: str
)

fn lib__work(v: str) -> L__Work rev 1
  emits [lib.bad]
  tests
    w("x") => Ok("x")
  Ok(v)
`
	app := `mod app
  provides [app__go, A__Out]
  uses [lib__work@1]
  emits [lib.bad]

type A__Out rev 1 (
  value: str
)

fn app__go(x: str) -> A__Out rev 1
  emits [lib.bad]
  tests
    ok("a") => Ok("a")
  match call lib__work(x)
    given
      ok => [exchange args (v = "a") outcome Ok("a")]
    on lib.bad e => forward e
    on Ok r => Ok(r.value)
`
	dir := writeLSPDir(t, map[string]string{"lib.can": lib, "app.can": app})
	if diags := diagnose(dir, "app.can", app); !hasCode(diags, "CAN4107") {
		t.Fatalf("untaken foreign forward reported no CAN4107: %v", diags)
	}
}

// TestForwardTrapStaysManual guards the migration trap:
// html__attribute__id selects the enclosing input, not e.value,
// so its nul_byte arm must remain handwritten, never forwarded.
func TestForwardTrapStaysManual(t *testing.T) {
	raw, err := os.ReadFile("../std/html/html.can")
	if err != nil {
		t.Fatal(err)
	}
	// a92: the manual spelling is positional now; the guard
	// stays the same, the arm reconstructs (never forward e).
	if !strings.Contains(string(raw), "on html.nul_byte e => html.nul_byte(value)") {
		t.Fatalf("trap arm left its manual shape")
	}
}

// TestForwardTrapDistinctValues pins why the trap stays manual: a
// reconstruction selecting the enclosing value passes exactly
// because the two values differ — forwarding e would change the
// verdict, and the rows prove it.
func TestForwardTrapDistinctValues(t *testing.T) {
	src := `mod m
  provides [m__go, m__check, M__Out, M__Mid]
  uses []
  emits [m.dirty]

error m.dirty(value: str)

type M__Out rev 1 (
  value: str
)

type M__Mid rev 1 (
  clean: str
)

fn m__check(v: str, tag: str) -> M__Mid rev 1
  emits [m.dirty]
  tests
    clean("a", "t") => Ok("a")
    dirty("a", "bad") => m.dirty("CALLEE")
  match tag == "bad"
    true => m.dirty("CALLEE")
    false => Ok(v)

fn m__go(value: str, tag: str) -> M__Out rev 1
  emits [m.dirty]
  tests
    hit("a-b", "bad") => m.dirty("a-b")
    pass("a", "t") => Ok("a")
  match call m__check(value, tag)
    on m.dirty e => m.dirty(value)
    on Ok c => Ok(c.clean)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("enclosing-value reconstruction reported: %v", diags)
	}
}
