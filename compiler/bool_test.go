package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Slice 5: eager boolean and/or/not. Probes first: truth
// tables, precedence (or < and < not < comparisons),
// bool-only operands, no calls in operands, symbol-alias
// rejection, eager fault loudness, helper lowering, and a
// generated-TS runtime parity gate.

const boolLib = `mod m
  provides [m__go, m__or, m__not, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: bool
)

`

// TestBoolBasic pins truth tables end to end.
func TestBoolBasic(t *testing.T) {
	src := boolLib + `fn m__go(left: bool, right: bool) -> M__Out rev 1
  emits []
  tests
    tt(left = true, right = true) => Ok(value = true)
    tf(left = true, right = false) => Ok(value = false)
    ft(left = false, right = true) => Ok(value = false)
    ff(left = false, right = false) => Ok(value = false)
=
  Ok(value = left and right)
` + `fn m__or(left: bool, right: bool) -> M__Out rev 1
  emits []
  tests
    tt(left = true, right = true) => Ok(value = true)
    tf(left = true, right = false) => Ok(value = true)
    ft(left = false, right = true) => Ok(value = true)
    ff(left = false, right = false) => Ok(value = false)
=
  Ok(value = left or right)
` + `fn m__not(x: bool) -> M__Out rev 1
  emits []
  tests
    t(x = true) => Ok(value = false)
    f(x = false) => Ok(value = true)
=
  Ok(value = not x)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("truth tables reported: %v", diags)
	}
}

// TestBoolPrecedence pins or < and < not < comparisons:
// `not a == b and c` reads `(not (a == b)) and c`, and
// `a or b and c` reads `a or (b and c)`.
func TestBoolPrecedence(t *testing.T) {
	src := `mod m
  provides [m__go, m__or2, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: bool
)

fn m__go(a: bool, b: bool, c: bool) -> M__Out rev 1
  emits []
  tests
    row(a = false, b = false, c = false) => Ok(value = false)
    row2(a = true, b = false, c = true) => Ok(value = true)
=
  Ok(value = not a == b and c)
` + `fn m__or2(a: bool, b: bool, c: bool) -> M__Out rev 1
  emits []
  tests
    row(a = true, b = true, c = false) => Ok(value = true)
=
  Ok(value = a or b and c)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); hasError(diags) {
		t.Fatalf("precedence reported: %v", diags)
	}
}

// TestBoolOperandTypes pins bool-only operands: integers
// are CAN6003, never truthy.
func TestBoolOperandTypes(t *testing.T) {
	for _, body := range []string{
		"Ok(value = 1 and right)",
		"Ok(value = left and 1)",
		"Ok(value = not 1)",
		"Ok(value = left or 1)",
	} {
		src := boolLib + `fn m__go(left: bool, right: bool) -> M__Out rev 1
  emits []
  tests
    go(left = true, right = true) => Ok(value = true)
=
  ` + body + "\n"
		dir := writeLSPDir(t, map[string]string{"m.can": src})
		if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN6003") {
			t.Fatalf("%q reported no CAN6003: %v", body, diags)
		}
	}
}

// TestBoolCallOperand pins no calls inside operands: a bare
// call stays CAN3003 outside a match scrutinee.
func TestBoolCallOperand(t *testing.T) {
	src := `mod m
  provides [m__go, m__flag, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: bool
)

fn m__flag() -> M__Out rev 1
  emits []
  tests
    go() => Ok(value = true)
=
  Ok(value = true)
` + `fn m__go(left: bool) -> M__Out rev 1
  emits []
  tests
    go(left = true) => Ok(value = true)
=
  Ok(value = left and call m__flag())
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN3003") {
		t.Fatalf("call operand reported no CAN3003: %v", diags)
	}
}

// TestBoolSymbolsRejected guards the no-alias rule: &&, ||,
// and ! stay parse errors.
func TestBoolSymbolsRejected(t *testing.T) {
	for _, body := range []string{
		"Ok(value = left && right)",
		"Ok(value = left || right)",
		"Ok(value = !left)",
	} {
		src := boolLib + `fn m__go(left: bool, right: bool) -> M__Out rev 1
  emits []
  tests
    go(left = true, right = true) => Ok(value = true)
=
  ` + body + "\n"
		dir := writeLSPDir(t, map[string]string{"m.can": src})
		if diags := diagnose(dir, "m.can", src); !hasError(diags) {
			t.Fatalf("%q reported nothing", body)
		}
	}
}

// TestBoolEagerFault pins loud faults: the well-typed fault
// cases fail the build (CAN4200), never settling into a
// boolean or a typed outcome.
func TestBoolEagerFault(t *testing.T) {
	for _, body := range []string{
		"Ok(value = false and ((1 / 0) == 0))",
		"Ok(value = true or ((1 / 0) == 0))",
	} {
		src := boolLib + `fn m__go() -> M__Out rev 1
  emits []
  tests
    go() => Ok(value = false)
=
  ` + body + "\n"
		dir := writeLSPDir(t, map[string]string{"m.can": src})
		if diags := diagnose(dir, "m.can", src); !hasCode(diags, "CAN4200") {
			t.Fatalf("%q reported no CAN4200: %v", body, diags)
		}
	}
}

// TestBoolEmitHelpers pins strict helper lowering: and/or
// emit $canBoolAnd/$canBoolOr calls (never bare &&), and a
// module using no boolean operator emits neither helper.
func TestBoolEmitHelpers(t *testing.T) {
	src := boolLib + `fn m__go(left: bool, right: bool) -> M__Out rev 1
  emits []
  tests
    go(left = true, right = true) => Ok(value = true)
=
  Ok(value = left and right)
` + `fn m__or(left: bool, right: bool) -> M__Out rev 1
  emits []
  tests
    go(left = true, right = true) => Ok(value = true)
=
  Ok(value = left or right)
` + `fn m__not(x: bool) -> M__Out rev 1
  emits []
  tests
    go(x = true) => Ok(value = false)
=
  Ok(value = not x)
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.can"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	ts := string(raw)
	for _, want := range []string{"$canBoolAnd(", "$canBoolOr(", "!("} {
		if !strings.Contains(ts, want) {
			t.Fatalf("emit missing %q", want)
		}
	}
	plain := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: bool
)

fn m__go(left: bool) -> M__Out rev 1
  emits []
  tests
    go(left = true) => Ok(value = true)
=
  Ok(value = left)
`
	dir2 := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir2, "m.can"), []byte(plain), 0o644); err != nil {
		t.Fatal(err)
	}
	out2 := t.TempDir()
	if err := compile(out2, []string{filepath.Join(dir2, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw2, err := os.ReadFile(filepath.Join(out2, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"$canBoolAnd", "$canBoolOr"} {
		if strings.Contains(string(raw2), want) {
			t.Fatalf("unused helper emitted: %q", want)
		}
	}
}

// TestBoolRuntimeParity executes the generated TS under node:
// truth tables must agree with the interpreter, and the
// eager discriminator (false left, faulting right) must
// throw instead of settling to false.
func TestBoolRuntimeParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("node missing: runtime parity refuses to skip: %v", err)
	}
	src := `mod m
  provides [m__and, m__or, m__not, m__eager, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: bool
)

fn m__and(left: bool, right: bool) -> M__Out rev 1
  emits []
  tests
    go(left = true, right = true) => Ok(value = true)
=
  Ok(value = left and right)

fn m__or(left: bool, right: bool) -> M__Out rev 1
  emits []
  tests
    go(left = true, right = true) => Ok(value = true)
=
  Ok(value = left or right)

fn m__not(x: bool) -> M__Out rev 1
  emits []
  tests
    go(x = true) => Ok(value = false)
=
  Ok(value = not x)

fn m__eager(x: int) -> M__Out rev 1
  emits []
  tests
    big(x = 200) => Ok(value = false)
    small(x = 50) => Ok(value = false)
=
  Ok(value = (x > 100) and ((10 / x) > 1))
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.can"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	driver := `import { m__and, m__or, m__not, m__eager } from "./m.ts";
const eq = (got, want, name) => {
  if (JSON.stringify(got) !== JSON.stringify(want)) throw new Error(name + ": got " + JSON.stringify(got));
};
eq(m__and(true, true), { $can_kind: "ok", value: true }, "and tt");
eq(m__and(true, false), { $can_kind: "ok", value: false }, "and tf");
eq(m__and(false, true), { $can_kind: "ok", value: false }, "and ft");
eq(m__or(false, false), { $can_kind: "ok", value: false }, "or ff");
eq(m__or(false, true), { $can_kind: "ok", value: true }, "or ft");
eq(m__not(true), { $can_kind: "ok", value: false }, "not t");
eq(m__not(false), { $can_kind: "ok", value: true }, "not f");
eq(m__eager(200n), { $can_kind: "ok", value: false }, "eager big");
let threw = false;
try { m__eager(0n); } catch { threw = true; }
if (!threw) throw new Error("eager zero did not throw");
console.log("BOOL_PARITY_OK");
`
	if err := os.WriteFile(filepath.Join(out, "driver.ts"), []byte(driver), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", "driver.ts")
	cmd.Dir = out
	raw, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("node parity failed: %v\n%s", err, raw)
	}
	if !strings.Contains(string(raw), "BOOL_PARITY_OK") {
		t.Fatalf("parity driver silent:\n%s", raw)
	}
}
