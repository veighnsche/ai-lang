package main

import (
	"strings"
	"testing"
)

// a81: proof obligations + z3 boundary. Admitted functions prove
// their per-exit obligations (Requires ∧ Path ⟹ Ensures) and
// call-site preconditions through verified callee summaries,
// with z3 deciding QF_LIA queries. Probes first.

// TestVerifyMax pins the pilot: max verifies end to end.
func TestVerifyMax(t *testing.T) {
	prog, texts := admitProg(t, admitMax)
	if diags := VerifyContracts(prog, texts); len(diags) != 0 {
		t.Fatalf("max did not verify: %v", diags)
	}
}

// TestVerifyValidator pins the pilot: the range validator verifies.
func TestVerifyValidator(t *testing.T) {
	prog, texts := admitProg(t, admitValidator)
	if diags := VerifyContracts(prog, texts); len(diags) != 0 {
		t.Fatalf("validator did not verify: %v", diags)
	}
}

const verifyAffine = `mod m
  provides [m__apply, m__affine, M__Out, M__Req, M__Pol]
  uses []
  emits []

type M__Pol rev 1 (
  minimum: int
)

type M__Req rev 1 (
  enabled: bool
  value: int
  policy: M__Pol
)

type M__Out rev 1 (
  value: int
)

fn m__affine(n: int) -> M__Out rev 1
  emits []
  requires
    n >= 0
  ensures
    on Ok result
      result.value == 2 * n + 1
  tests
    z(n = 0) => Ok(value = 1)
    two(n = 2) => Ok(value = 5)
=
  Ok(value = 2 * n + 1)

fn m__apply(req: M__Req) -> M__Out rev 1
  emits []
  requires
    req.policy.minimum >= 0
  ensures
    on Ok result
      match req.enabled
        true => match req.value >= req.policy.minimum
          true => result.value == 2 * req.value + 1
          false => result.value == req.policy.minimum
        false => result.value == req.policy.minimum
  tests
    hit(req = M__Req(enabled = true, value = 3, policy = M__Pol(minimum = 2))) => Ok(value = 7)
    miss(req = M__Req(enabled = true, value = 1, policy = M__Pol(minimum = 2))) => Ok(value = 2)
    off(req = M__Req(enabled = false, value = 3, policy = M__Pol(minimum = 2))) => Ok(value = 2)
=
  match req.enabled
    true => match req.value >= req.policy.minimum
      true => match call m__affine(req.value)
        on Ok r => Ok(value = r.value)
      false => Ok(value = req.policy.minimum)
    false => Ok(value = req.policy.minimum)
`

// TestVerifyComposition pins the third pilot: nested projections,
// a Boolean input, nontrivial entry assumptions, path-sensitive
// call admission, constant multiplication, and a verified result
// summary, all in one program.
func TestVerifyComposition(t *testing.T) {
	prog, texts := admitProg(t, verifyAffine)
	if diags := VerifyContracts(prog, texts); len(diags) != 0 {
		t.Fatalf("composition did not verify: %v", diags)
	}
}

// TestVerifyTheoryControl pins backend selection: 2*x==1 is
// unsatisfiable over integers, so the constant body verifies. A
// real-arithmetic backend would offer x=1/2 and wrongly reject.
func TestVerifyTheoryControl(t *testing.T) {
	src := `mod m
  provides [m__int, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__int(x: int) -> M__Out rev 1
  emits []
  requires
    true
  ensures
    on Ok result
      match 2 * x == 1
        true => false
        false => true
  tests
    z(x = 0) => Ok(value = 0)
=
  Ok(value = 0)
`
	prog, texts := admitProg(t, src)
	if diags := VerifyContracts(prog, texts); len(diags) != 0 {
		t.Fatalf("integer control did not verify: %v", diags)
	}
}

// TestVerifyRecordEquality pins whole-record equality: same
// nominal type flattens into field equalities.
func TestVerifyRecordEquality(t *testing.T) {
	src := `mod m
  provides [m__pass, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__pass(a: M__Out, b: M__Out) -> M__Out rev 1
  emits []
  requires
    a == b
  ensures
    on Ok result
      result == a
  tests
    same(a = M__Out(value = 1), b = M__Out(value = 1)) => Ok(value = 1)
=
  Ok(value = a.value)
`
	prog, texts := admitProg(t, src)
	if diags := VerifyContracts(prog, texts); len(diags) != 0 {
		t.Fatalf("record equality did not verify: %v", diags)
	}
}

const verifyOutageSvc = `mod svc
  provides [svc__get, Svc__Out]
  uses []
  emits [service.unavailable]

error service.unavailable(code: int)

type Svc__Out rev 1 (
  value: int
)

fn svc__get(x: int) -> Svc__Out rev 1
  emits [service.unavailable]
  requires
    true
  ensures
    on Ok result
      result.value == x
    on service.unavailable err
      false
  tests
    go(x = 2) => Ok(value = 2)
=
  Ok(value = x)
`

const verifyOutageApp = `mod app
  provides [app__use]
  uses [svc__get@1, Svc__Out@1]
  emits [app.down]

error app.down(code: int)

fn app__use(x: int) -> Svc__Out rev 1
  emits [app.down, service.unavailable]
  requires
    true
  ensures
    on Ok result
      result.value == x
    on app.down err
      false
    on service.unavailable err
      false
  tests
    normal(x = 2) => Ok(value = 2)
    outage(x = 2) => app.down(code = 2)
=
  match call svc__get(x)
    given
      normal => [exchange args (x = 2) outcome Ok(value = 2)]
      outage => [exchange args (x = 2) outcome service.unavailable(code = 2)]
    on Ok r => Ok(value = r.value)
    on service.unavailable e => app.down(code = e.code)
`

// TestVerifyOutageSeparation pins verdict position 6: the provider
// proof covers actual executions, the injected outage row covers
// handler execution, and neither disturbs the other. Both functions
// verify while the outage test still passes ordinary checks.
func TestVerifyOutageSeparation(t *testing.T) {
	files := map[string]string{"svc.can": verifyOutageSvc, "app.can": verifyOutageApp}
	prog, texts := revisionProg(t, files, []string{"svc.can", "app.can"})
	if diags := VerifyContracts(prog, texts); len(diags) != 0 {
		t.Fatalf("outage separation did not verify: %v", diags)
	}
}

// TestVerifyMaxPlusOne pins refutation with a countermodel: the
// overlarge result violates membership, and the finding names the
// inputs that witness it.
func TestVerifyMaxPlusOne(t *testing.T) {
	// Rows follow the mutant body (execution stays green); the
	// proof is what rejects the overlarge result.
	bad := strings.Replace(admitMax,
		"on true => Ok(value = right)", "on true => Ok(value = right + 1)", 1)
	bad = strings.Replace(bad,
		"ordered(left = 1, right = 2) => Ok(value = 2)", "ordered(left = 1, right = 2) => Ok(value = 3)", 1)
	prog, texts := admitProg(t, bad)
	diags := VerifyContracts(prog, texts)
	if !hasCode(diags, CodeContractUnproven) {
		t.Fatalf("overlarge max reported no unproven finding: %v", diags)
	}
	if !hasFound(diags, "counterexample") || !hasFound(diags, "left") || !hasFound(diags, "right") {
		t.Fatalf("unproven finding carries no input countermodel: %v", diags)
	}
}

// TestVerifyHelperBody pins per-function blame: the mutated helper
// fails its own postcondition, and the caller degrades to an
// unavailable summary (not a false proof).
func TestVerifyHelperBody(t *testing.T) {
	// Rows follow the mutant body; the helper postcondition is
	// what fails.
	bad := strings.Replace(verifyAffine, "Ok(value = 2 * n + 1)", "Ok(value = 2 * n + 2)", 1)
	bad = strings.Replace(bad, "z(n = 0) => Ok(value = 1)", "z(n = 0) => Ok(value = 2)", 1)
	bad = strings.Replace(bad, "two(n = 2) => Ok(value = 5)", "two(n = 2) => Ok(value = 6)", 1)
	bad = strings.Replace(bad, "hit(req = M__Req(enabled = true, value = 3, policy = M__Pol(minimum = 2))) => Ok(value = 7)",
		"hit(req = M__Req(enabled = true, value = 3, policy = M__Pol(minimum = 2))) => Ok(value = 8)", 1)
	prog, texts := admitProg(t, bad)
	diags := VerifyContracts(prog, texts)
	if !hasCode(diags, CodeContractUnproven) {
		t.Fatalf("mutated helper reported no unproven finding: %v", diags)
	}
	if !hasFound(diags, "m__affine") {
		t.Fatalf("unproven finding names no helper: %v", diags)
	}
	if !hasCode(diags, CodeContractUnverifiedDep) {
		t.Fatalf("caller of failed helper reported no unavailable finding: %v", diags)
	}
}

// TestVerifyWeakSummary pins modular blame: the helper verifies a
// weaker contract, so the caller proof fails against the available
// assumptions, reported as a countermodel modulo those summaries
// rather than a demonstrated execution.
func TestVerifyWeakSummary(t *testing.T) {
	weak := strings.Replace(verifyAffine, "result.value == 2 * n + 1", "result.value >= 0", 1)
	prog, texts := admitProg(t, weak)
	diags := VerifyContracts(prog, texts)
	if !hasCode(diags, CodeContractUnproven) {
		t.Fatalf("weak summary reported no unproven finding: %v", diags)
	}
	if !hasFound(diags, "m__apply") {
		t.Fatalf("unproven finding names no caller: %v", diags)
	}
	if !hasFound(diags, "modulo") {
		t.Fatalf("weak-summary failure not marked modulo assumptions: %v", diags)
	}
}

const verifyStep = `mod m
  provides [m__step, m__call, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__step(n: int) -> M__Out rev 1
  emits []
  requires
    n >= 0
  ensures
    on Ok result
      result.value == n + 1
  tests
    z(n = 0) => Ok(value = 1)
=
  Ok(value = n + 1)

fn m__call(x: int) -> M__Out rev 1
  emits []
  requires
    true
  ensures
    on Ok result
      result.value >= 1
  tests
    neg(x = -2) => Ok(value = 3)
    pos(x = 2) => Ok(value = 3)
=
  match x <= 0
    true => match call m__step(ARG)
      on Ok r => Ok(value = r.value)
    false => match call m__step(x)
      on Ok r => Ok(value = r.value)
`

// TestVerifyStepCall pins path-sensitive admission both ways: the
// correct argument verifies, and the x-mutation fails the call
// precondition with an x countermodel.
func TestVerifyStepCall(t *testing.T) {
	good := strings.Replace(verifyStep, "m__step(ARG)", "m__step(0 - x)", 1)
	prog, texts := admitProg(t, good)
	if diags := VerifyContracts(prog, texts); len(diags) != 0 {
		t.Fatalf("path-sensitive call did not verify: %v", diags)
	}
	// The row follows the mutant body (x = -2 yields -1 through
	// the bad call); the call precondition is what fails.
	bad := strings.Replace(verifyStep, "m__step(ARG)", "m__step(x)", 1)
	bad = strings.Replace(bad, "neg(x = -2) => Ok(value = 3)", "neg(x = -2) => Ok(value = -1)", 1)
	prog, texts = admitProg(t, bad)
	diags := VerifyContracts(prog, texts)
	if !hasCode(diags, CodeContractUnproven) {
		t.Fatalf("bad call argument reported no unproven finding: %v", diags)
	}
	if !hasFound(diags, "m__step") || !hasFound(diags, "x") {
		t.Fatalf("precondition failure names no callee or input: %v", diags)
	}
}

// TestVerifyBadRow pins inadmissible test inputs: the row violates
// requires while passing execution, so the program is clean but
// the row is rejected.
func TestVerifyBadRow(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(x: int) -> M__Out rev 1
  emits []
  requires
    x >= 0
  ensures
    on Ok result
      result.value == x
  tests
    neg(x = -1) => Ok(value = -1)
=
  Ok(value = x)
`
	prog, texts := admitProg(t, src)
	diags := VerifyContracts(prog, texts)
	if !hasCode(diags, CodeContractInadmissibleTest) {
		t.Fatalf("out-of-requires row reported no inadmissible finding: %v", diags)
	}
	if !hasFound(diags, "neg") {
		t.Fatalf("inadmissible finding names no row: %v", diags)
	}
}

// TestVerifyMissingSolver pins fail-closed tooling: with no solver
// binary, contracted functions are inconclusive, never accepted.
func TestVerifyMissingSolver(t *testing.T) {
	t.Setenv("CANLC_Z3", "/nonexistent/z3-bogus")
	prog, texts := admitProg(t, admitMax)
	diags := VerifyContracts(prog, texts)
	if len(diags) == 0 {
		t.Fatalf("missing solver verified nothing yet reported nothing")
	}
	if !hasCode(diags, CodeContractInconclusive) {
		t.Fatalf("missing solver reported no inconclusive finding: %v", diags)
	}
	if hasCode(diags, CodeContractUnproven) {
		t.Fatalf("missing solver reported unproven instead of inconclusive: %v", diags)
	}
}

// TestVerifyIntArms pins integer case analysis: range and
// wildcard arms encode as value constraints, so an admitted
// integer match discharges instead of going inconclusive.
func TestVerifyIntArms(t *testing.T) {
	prog, texts := admitProg(t, verifyClamp)
	if diags := VerifyContracts(prog, texts); len(diags) != 0 {
		t.Fatalf("clamp did not verify: %v", diags)
	}
}

const verifyClamp = `mod m
  provides [m__clamp, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__clamp(x: int) -> M__Out rev 1
  emits []
  requires
    x >= 0
  ensures
    on Ok result
      result.value >= 0
      result.value <= 10
  tests
    lo(x = 3) => Ok(value = 3)
    hi(x = 15) => Ok(value = 10)
    edge(x = 0) => Ok(value = 0)
=
  match x
    0..10 => Ok(value = x)
    _ => Ok(value = 10)
`

// TestParseSMTResult pins response classification: only exact
// sat/unsat/unknown decide; anything else is a tooling failure.
func TestParseSMTResult(t *testing.T) {
	for out, want := range map[string]string{
		"sat\n":     "sat",
		"unsat\n":   "unsat",
		"unknown\n": "unknown",
	} {
		if got, err := parseSMTResult(out); err != nil || got != want {
			t.Fatalf("parse %q: got %q err %v", out, got, err)
		}
	}
	for _, out := range []string{"", "sat unknown\n", "error\n", "(\n"} {
		if _, err := parseSMTResult(out); err == nil {
			t.Fatalf("parse %q: expected tooling error", out)
		}
	}
}

// TestParseSMTModel pins countermodel extraction: input variables
// read back with mathematical-integer values.
func TestParseSMTModel(t *testing.T) {
	got := parseSMTModel("((|left| 0)\n (|right| (- 3))\n (|flag| false))", []string{"left", "right", "flag"})
	if got["left"] != "0" || got["right"] != "-3" || got["flag"] != "false" {
		t.Fatalf("bad model parse: %v", got)
	}
}
