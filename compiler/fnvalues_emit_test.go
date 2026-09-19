package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// I4 emit: references lower to self-applying closure thunks,
// invocations dispatch through the shared arm chains, and the
// sketch below freezes the whole surface (thunk shape, capture
// temporaries, parameter order, imported targets, Fn-bearing
// records, and the error catalogue) byte-for-byte.
func TestGoldenFnCallback(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../sketches/fn-callback/ops.can",
		"../sketches/fn-callback/use.can",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"ops.ts", "use.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/fn-callback", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		checkGoldenFile(t, f, got, want)
	}
}

// The language refuses function equality: directly on callables
// and through bearing records. Factory receipt comparison stays
// in the test evaluator; == never reaches it.
const fnEqDirect = `mod m
  provides [m__go, m__t, M__O]
  uses []
  emits []

type M__O rev 1 (
  value: str
)

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits []
  tests
    t(3, 7) => Ok("q")
  Ok("q")

fn m__go(cb: Fn<int, M__O, []>, n: int) -> M__O rev 1
  emits []
  tests
    g(fnref m__t(divisor = 3), 4) => Ok("q")
  match cb == cb
    on true => Ok("q")
    on false => Ok("e")
`

const fnEqBearing = `mod m
  provides [m__go, m__t, M__O, M__H]
  uses []
  emits []

type M__O rev 1 (
  value: str
)

type M__H rev 1 (
  cb: Fn<int, M__O, []>
)

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits []
  tests
    t(3, 7) => Ok("q")
  Ok("q")

fn m__go(h: M__H, k: M__H) -> M__O rev 1
  emits []
  tests
    g(M__H(fnref m__t(divisor = 3)), M__H(fnref m__t(divisor = 3))) => Ok("q")
  match h == k
    on true => Ok("q")
    on false => Ok("e")
`

func TestFnEqualityRefused(t *testing.T) {
	seqCode(t, map[string]string{"m.can": fnEqDirect}, "m.can",
		CodeTypeMismatch, "function values do not compare")
	seqCode(t, map[string]string{"m.can": fnEqBearing}, "m.can",
		CodeTypeMismatch, "function values do not compare")
}

// The emitter refuses callable comparison too (CAN5005), so a
// broken program can never lower to object identity.
func TestEmitEqualityRejectsFn(t *testing.T) {
	em := &emitter{
		recs: map[string][][2]string{
			"M__H": {{"cb", "Fn<int, M__O, []>"}},
		},
	}
	if _, err := em.emitEquality("==", "Fn<int, M__O, []>", "a", "b"); err == nil ||
		!strings.Contains(err.Error(), "function values do not compare") {
		t.Fatalf("direct callable comparison must fail closed, got %v", err)
	}
	if _, err := em.emitEquality("==", "M__H", "a", "b"); err == nil ||
		!strings.Contains(err.Error(), "record contains a function value") {
		t.Fatalf("bearing record comparison must fail closed, got %v", err)
	}
}

// Every Go-side callable refusal has an emitter counterpart:
// an unchecked program (unit-built, or a future pipeline slip)
// fails the build instead of emitting a guessed closure.
func TestEmitCallableRefusals(t *testing.T) {
	em := &emitter{}
	var out []string
	// No static certificate: the Go side reports CAN6003 for a
	// non-callable or unbound target; emit refuses the same node.
	invoke := &Node{Kind: MatchInvoke,
		Scruts: []*Small{{Kind: "ref", Ref: []string{"cb"}}},
		InvokeArg: &Small{Kind: "int"}}
	if err := em.stmtInvoke(invoke, &out); err == nil ||
		!strings.Contains(err.Error(), "without a checked signature") {
		t.Fatalf("uncertified invoke must fail closed, got %v", err)
	}
	// Unknown reference target: check reports CAN6020-family or
	// unbound; emit refuses before inventing parameters.
	if _, err := em.emitFnref(&Small{Kind: "fnref", Fname: "m__nope"}); err == nil ||
		!strings.Contains(err.Error(), "no params for reference target") {
		t.Fatalf("unknown fnref target must fail closed, got %v", err)
	}
	// Bare-Fn return: check reports CAN6003 demanding a record;
	// the Ok-shape lookup finds no record and fails the build.
	fn := &FnDecl{Name: "m__mk", Ret: "Fn<int, M__O, []>"}
	if _, err := declaredOkShape(fn, &Program{}); err == nil ||
		!strings.Contains(err.Error(), "returns unknown type") {
		t.Fatalf("bare-Fn return must fail closed, got %v", err)
	}
}

// Go/TypeScript parity over the callback sketch: the same
// vectors execute on both sides. Two closures stay alive
// simultaneously with independent captures, the older one
// invokes again, a large integer crosses host-number exactness,
// the decimal callback keeps canonical strings, and the zero
// divisor takes the error arm. Node absence fails the gate.
func TestFnCallbackParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("node missing: runtime parity refuses to skip: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/fn-callback/ops.can", "../sketches/fn-callback/use.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	// Same harness shim as the S1a parity tests: node ESM needs
	// explicit extensions the bundler-style emit omits.
	rewriteSpecifier(t, filepath.Join(out, "use.ts"), "./ops", "./ops.ts")
	driver := `import { use__bind, use__run, use__go, use__binddec, use__rundec } from "./use.ts";
const show = (x) => JSON.stringify(x, (k, v) => typeof v === "bigint" ? { $bigint: v.toString() } : v);
const eq = (got, want, name) => {
  if (show(got) !== show(want)) throw new Error(name + ": got " + show(got));
};
const Q = (q, r) => ({ $can_kind: "ok", quotient: q, remainder: r });
const c3 = use__bind(3n);
const c5 = use__bind(5n);
eq(use__run(c3.cb, 7n), Q(2n, 1n), "run 3");
eq(use__run(c5.cb, 7n), Q(1n, 2n), "run 5");
eq(use__run(c3.cb, 8n), Q(2n, 2n), "run 3 again");
eq(use__run(use__bind(3n).cb, 9007199254740993n), Q(3002399751580331n, 0n), "beyond exactness");
eq(use__run(use__bind(0n).cb, 7n), Q(0n, 0n), "zero divisor");
const d = use__binddec("0.5");
eq(use__rundec(d.cb, "1.0"), { $can_kind: "ok", value: "1.5" }, "dec bump");
eq(use__go(3n, 7n), Q(2n, 1n), "linked");
console.log("FN_PARITY_OK");
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
	if !strings.Contains(string(raw), "FN_PARITY_OK") {
		t.Fatalf("parity driver silent:\n%s", raw)
	}
}
