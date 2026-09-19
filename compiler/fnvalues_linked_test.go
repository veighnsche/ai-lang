package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Linked b00 vectors: a real source factory builds the callback,
// the carrier record crosses into a real consumer, and invocation
// executes the original checked body. All four sign combinations
// plus zero dividend and zero divisor run through use__go; the
// bind root returns a callback the receipt comparator checks;
// one vector passes a callback inline as a linked argument; and
// a 1023-deep driver proves nested invocation inherits the
// caller's depth accounting instead of resetting it.
func fnCallbackLinkedFiles(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, f := range []string{"ops.can", "use.can"} {
		raw, err := os.ReadFile("../sketches/fn-callback/" + f)
		if err != nil {
			t.Fatal(err)
		}
		out[f] = string(raw)
	}
	return out
}

func TestFnCallbackLinked(t *testing.T) {
	files := fnCallbackLinkedFiles(t)
	order := []string{"ops.can", "use.can"}
	vectors := []struct {
		name   string
		root   string
		args   map[string]string
		expect string
	}{
		{"pos pos", "use__go", map[string]string{"d": "3", "n": "7"}, "Ok(2, 1)"},
		{"pos neg", "use__go", map[string]string{"d": "3", "n": "-7"}, "Ok(-3, 2)"},
		{"neg pos", "use__go", map[string]string{"d": "-3", "n": "7"}, "Ok(-2, 1)"},
		{"neg neg", "use__go", map[string]string{"d": "-3", "n": "-7"}, "Ok(3, 2)"},
		{"zero dividend", "use__go", map[string]string{"d": "3", "n": "0"}, "Ok(0, 0)"},
		{"zero divisor", "use__go", map[string]string{"d": "0", "n": "7"}, "Ok(0, 0)"},
		{"factory carry", "use__bind", map[string]string{"d": "3"},
			"Ok(fnref ops__divmod(divisor = 3))"},
		{"inline callback arg", "use__run",
			map[string]string{"cb": "fnref ops__divmod(divisor = 3)", "n": "7"}, "Ok(2, 1)"},
		{"decimal callback", "use__rundec",
			map[string]string{"cb": `fnref ops__bump(step = d"0.5")`, "x": `d"1.0"`},
			`Ok(d"1.5")`},
	}
	for _, v := range vectors {
		if err := runLinkedPure(t, files, order, v.root, 1, v.args, v.expect); err != nil {
			t.Fatalf("linked %s: %v", v.name, err)
		}
	}
}

// A driver 1023 calls deep invokes a callback whose target calls
// one more helper. Inherited accounting trips the 1024 backstop
// inside the target; a reset would let the helper run clean.
const fnDepthSrc = `mod m
  provides [m__deep, m__t, m__leaf, M__O]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__O rev 1 (
  value: str
)

fn m__leaf(x: int) -> M__O rev 1
  emits []
  tests
    l(1) => Ok("q")
  Ok("q")

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits [m.err]
  tests
    t(3, 7) => Ok("q")
    tz(3, 0) => m.err("zero")
  match dividend
    0 => m.err("zero")
    _ => match call m__leaf(dividend)
      on Ok _ => Ok("q")

fn m__deep(cb: Fn<int, M__O, [m.err]>, n: int, seed: int) -> M__O rev 1
  decreases n
  emits []
  tests
    d0(fnref m__t(divisor = 3), 0, 5) => Ok("q")
    dneg(fnref m__t(divisor = 3), 0, 0) => Ok("e")
    d1(fnref m__t(divisor = 3), 1, 5) => Ok("q")
    d1e(fnref m__t(divisor = 3), 1, 0) => Ok("q")
  match n <= 0
    true => match invoke cb with seed
      on Ok _ => Ok("q")
      on m.err _ => Ok("e")
    false => match call m__deep(cb, n - 1, seed)
      on Ok _ => Ok("q")
`

func TestFnLinkedInheritsDepth(t *testing.T) {
	files := map[string]string{"m.can": fnDepthSrc}
	args := map[string]string{"cb": "fnref m__t(divisor = 3)", "n": "1023", "seed": "5"}
	err := runLinkedPure(t, files, []string{"m.can"}, "m__deep", 1, args, `Ok("q")`)
	if err == nil || !strings.Contains(err.Error(), "local call depth exceeded") {
		t.Fatalf("deep invoke must trip the inherited backstop, got %v", err)
	}
}

// TestFnLinkedConsumerFirst pins the CLI mirror of
// prepareProviders: checkSem checks and RUNS each module in
// paths order, but an invoked callback executes its provider
// body linked — including multi-arg constructors that only bind
// after the provider's static phase. Consumer-first order must
// still compile: the CLI prepares every module's statics before
// any test row runs. (S9 found this via dec__encode_apply
// invoking into scalars' two-field std__dec__parts; the minimal
// trigger is a depth-two provider chain under the invoke.)
func TestFnLinkedConsumerFirst(t *testing.T) {
	lib := `mod probelib
  provides [lib__pair2, lib__d2s, Pair__Value, Str__Value]
  uses []
  emits []

type Str__Value rev 1 (
  value: str
)

type Pair__Value rev 1 (
  a: int
  b: int
)

fn lib__pair2(x: dec) -> Pair__Value rev 1
  emits []
  tests
    one(d"12.34") => Ok(1234, 2)
  match call dec__parts(x)
    on Ok p => Ok(p.coefficient, p.scale)

fn lib__d2s(value: dec) -> Str__Value rev 1
  emits []
  tests
    s2s_simple(d"12.34") => Ok("12.34")
  match call lib__pair2(value)
    on Ok p => Ok("12.34")
`
	use := `mod probeuse
  provides [use__apply, use__target, Int__Value]
  uses [lib__d2s@1]
  emits []

type Int__Value rev 1 (
  value: int
)

fn use__target(x: dec) -> Int__Value rev 1
  emits []
  tests
    one(d"12.34") => Ok(5)
  match call lib__d2s(x)
    given
      one => [exchange args (value = d"12.34") outcome Ok("12.34")]
    on Ok r => Ok(5)

fn use__apply(x: dec, f: Fn<dec, Int__Value, []>) -> Int__Value rev 1
  emits []
  tests
    one(d"12.34", fnref use__target()) => Ok(5)
  match invoke f with x
    on Ok r => Ok(r.value)
`
	dir := writeLSPDir(t, map[string]string{"use.can": use, "lib.can": lib})
	out := t.TempDir()
	// Consumer first: the order that executed raw bodies.
	if err := compile(out, []string{filepath.Join(dir, "use.can"), filepath.Join(dir, "lib.can")}); err != nil {
		t.Fatalf("consumer-first compile: %v", err)
	}
}
