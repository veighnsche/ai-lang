package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// a90: arm-level forward call. Probes first: check-time
// elaboration of `=> forward call f(args)` into the full
// match, CAN3013 rejections, the row law on elaborated arms,
// and emit parity with the handwritten shape.

// TestForwardCallValueArm pins the core rewrite: a value arm
// relaying a local call behaves exactly like the handwritten
// match, with every elaborated arm taken by a committed row.
func TestForwardCallValueArm(t *testing.T) {
	src := `mod m
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

fn m__go(x: str) -> M__Out rev 1
  emits [m.bad]
  tests
    ok("a") => Ok("a")
    bad("b") => m.bad("b")
    other("z") => Ok("no")
  match x == "z"
    true => Ok("no")
    false => forward call m__work(x)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("forward call in value arm reported: %v", diags)
	}
}

// TestForwardCallCallArm pins elaboration under a call arm with
// a bound payload argument, including an enclosing binder that
// collides with the elaborated names.
func TestForwardCallCallArm(t *testing.T) {
	src := `mod m
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

fn m__go(x: str) -> M__Out rev 1
  emits [m.bad]
  tests
    ok("a") => Ok("a")
    bad("b") => m.bad("b")
  match call m__work(x)
    on Ok r => forward call m__work(r.value)
    on m.bad e => forward e
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("forward call in call arm reported: %v", diags)
	}
}

// TestForwardCallRecurses pins termination: an elaborated
// self-call under the canonical guard satisfies decreases.
func TestForwardCallRecurses(t *testing.T) {
	src := `mod m
  provides [m__walk, M__W]
  uses []
  emits [m.bad]

error m.bad(value: str)

type M__W rev 1 (
  value: str
)

fn m__walk(s: str, acc: str, n: int) -> M__W rev 1
  decreases n
  emits [m.bad]
  tests
    done("q", "", 0) => Ok("")
    step("ab", "", 2) => Ok("ab")
    bad("aB", "", 2) => m.bad("a")
  match n <= 0
    true => Ok(acc)
    false => match s[0:1] == "B"
      true => m.bad(acc)
      false => forward call m__walk(s[1:#s], acc + s[0:1], n - 1)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("forward self-call under guard reported: %v", diags)
	}
}

// TestForwardCallRejects pins CAN3013: foreign, extern, and
// unknown callees plus malformed sites are refused. A binder
// that merely starts with "call" keeps its CAN3011.
func TestForwardCallRejects(t *testing.T) {
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
	cases := []struct {
		name  string
		files map[string]string
		file  string
	}{
		{
			"foreign",
			map[string]string{"lib.can": lib, "app.can": `mod app
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
  match x == "z"
    true => Ok("no")
    false => forward call lib__work(x)
`},
			"app.can",
		},
		{
			"extern",
			map[string]string{"m.can": `mod m
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
  match x == "z"
    true => Ok("no")
    false => forward call ex__work(x)
`},
			"m.can",
		},
		{
			"unknown",
			map[string]string{"m.can": `mod m
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
  match x == "z"
    true => Ok("no")
    false => forward call nope(x)
`},
			"m.can",
		},
		{
			"bare",
			map[string]string{"m.can": `mod m
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
  match x == "z"
    true => Ok("no")
    false => forward call
`},
			"m.can",
		},
		{
			"noncall",
			map[string]string{"m.can": `mod m
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
  match x == "z"
    true => Ok("no")
    false => forward call 123
`},
			"m.can",
		},
		{
			"unbalanced",
			map[string]string{"m.can": `mod m
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
  match x == "z"
    true => Ok("no")
    false => forward call m__go(
`},
			"m.can",
		},
	}
	for _, c := range cases {
		dir := writeLSPDir(t, c.files)
		if diags := diagnose(dir, c.file, c.files[c.file]); !hasCode(diags, "CAN3013") {
			t.Fatalf("%s reported no CAN3013: %v", c.name, diags)
		}
	}
	// A binder that merely starts with "call" is an ordinary
	// bad operand (CAN3011), never a malformed call site.
	binder := `mod m
  provides [m__go, M__Out]
  uses []
  emits [m.bad]

error m.bad(value: str)

type M__Out rev 1 (
  value: str
)

fn m__go(x: str) -> M__Out rev 1
  emits [m.bad]
  tests
    ok("a") => Ok("a")
  match call m__noop(x)
    on m.bad e => forward calling
    on Ok r => Ok(r.value)

fn m__noop(x: str) -> M__Out rev 1
  emits [m.bad]
  tests
    ok("a") => Ok("a")
  Ok(x)
`
	dir := writeLSPDir(t, map[string]string{"m.can": binder})
	diags := diagnose(dir, "m.can", binder)
	if !hasCode(diags, "CAN3011") {
		t.Fatalf("calling binder reported no CAN3011: %v", diags)
	}
	if hasCode(diags, "CAN3013") {
		t.Fatalf("calling binder wrongly claimed as call site: %v", diags)
	}
}

// TestForwardCallChainPositions pins the v1 boundary: chain
// tail and else keep refusing forward shapes at parse with
// the existing errors — no forward-call expansion runs there.
func TestForwardCallChainPositions(t *testing.T) {
	lib := `mod m
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

fn m__go(x: str) -> M__Out rev 1
  emits [m.bad]
  tests
    ok("a") => Ok("a")
    bad("b") => m.bad("b")
  match chain
    call m__work(x) as w
`
	tail := lib + `    then forward call m__work(w.value)
    else m.bad(x)
`
	dir := writeLSPDir(t, map[string]string{"m.can": tail})
	diags := diagnose(dir, "m.can", tail)
	if !hasError(diags) {
		t.Fatalf("chain tail forward call reported nothing: %v", diags)
	}
	found := false
	for _, d := range diags {
		if strings.Contains(d.Msg, "bare forward") {
			found = true
		}
	}
	if !found {
		t.Fatalf("chain tail refusal must name the bare forward: %v", diags)
	}
	els := lib + `    then Ok(w.value)
    else forward call m__work(x)
`
	dir = writeLSPDir(t, map[string]string{"m.can": els})
	diags = diagnose(dir, "m.can", els)
	if !hasError(diags) {
		t.Fatalf("chain else forward call reported nothing: %v", diags)
	}
	found = false
	for _, d := range diags {
		if strings.Contains(d.Msg, "cannot parse expression") {
			found = true
		}
	}
	if !found {
		t.Fatalf("chain else refusal must fail the expression parse: %v", diags)
	}
}

// TestForwardCallUndeclaredEmits pins the ordering benefit: a
// forwarded kind outside the caller's emits reports CAN4001,
// exactly like spelling the relay out by hand.
func TestForwardCallUndeclaredEmits(t *testing.T) {
	src := `mod m
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

fn m__go(x: str) -> M__Out rev 1
  emits []
  tests
    ok("a") => Ok("a")
  match x == "z"
    true => Ok("no")
    false => forward call m__work(x)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN4001") {
		t.Fatalf("undeclared forwarded kind reported no CAN4001: %v", diags)
	}
}

// TestForwardCallOkShape pins the ordering benefit for Ok: a
// callee success record that drifts from the caller's reports
// CAN3011 through the elaborated arm.
func TestForwardCallOkShape(t *testing.T) {
	src := `mod m
  provides [m__go, m__work, M__Out, M__Work]
  uses []
  emits []

type M__Out rev 1 (
  value: str
)

type M__Work rev 1 (
  other: str
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
  match x == "z"
    true => Ok("no")
    false => forward call m__work(x)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN3011") {
		t.Fatalf("drifted Ok shape reported no CAN3011: %v", diags)
	}
}

// TestForwardCallRowLaw pins CAN4107 on elaborated arms: an
// untaken Ok arm fails, while an untaken local error relay
// carries the certificate (CAN4108), like handwritten arms.
func TestForwardCallRowLaw(t *testing.T) {
	errOnly := `mod m
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

fn m__go(x: str) -> M__Out rev 1
  emits [m.bad]
  tests
    bad("b") => m.bad("b")
  match x == "z"
    true => m.bad(x)
    false => forward call m__work(x)
`
	dir := writeLSPDir(t, map[string]string{"m.can": errOnly})
	if diags := diagnose(dir, "m.can", errOnly); !hasCode(diags, "CAN4107") {
		t.Fatalf("untaken elaborated Ok reported no CAN4107: %v", diags)
	}
	okOnly := `mod m
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

fn m__go(x: str) -> M__Out rev 1
  emits [m.bad]
  tests
    ok("a") => Ok("a")
    other("z") => Ok("no")
  match x == "z"
    true => Ok("no")
    false => forward call m__work(x)
`
	dir = writeLSPDir(t, map[string]string{"m.can": okOnly})
	if diags := diagnose(dir, "m.can", okOnly); hasError(diags) {
		t.Fatalf("certified untaken error relay reported: %v", diags)
	}
}

// TestForwardCallEmitParity pins byte-identical emit: the
// sugared relay and its handwritten expansion compile to the
// same TypeScript.
func TestForwardCallEmitParity(t *testing.T) {
	hand := `mod m
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

fn m__go(x: str) -> M__Out rev 1
  emits [m.bad]
  tests
    ok("a") => Ok("a")
    bad("b") => m.bad("b")
    other("z") => Ok("no")
  match x == "z"
    true => Ok("no")
    false => match call m__work(x)
      on Ok r => forward r
      on m.bad e => forward e
`
	sugared := strings.Replace(hand,
		`    false => match call m__work(x)
      on Ok r => forward r
      on m.bad e => forward e
`, `    false => forward call m__work(x)
`, 1)
	emit := func(name, src string) string {
		t.Helper()
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "m.can"), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
		out := t.TempDir()
		if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
			t.Fatalf("compile %s: %v", name, err)
		}
		raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
		if err != nil {
			t.Fatal(err)
		}
		return string(raw)
	}
	if got, want := emit("sugared", sugared), emit("hand", hand); got != want {
		t.Fatalf("emit differs:\n--- sugared ---\n%s\n--- hand ---\n%s", got, want)
	}
}
