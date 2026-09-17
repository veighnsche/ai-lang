package main

import (
	"strings"
	"testing"
)

// a80: contract admission for the verifier first cut. A contracted
// function is admitted only when its complete contract and body lie
// inside the declared proof fragment (mathematical ints, Booleans,
// finite acyclic records of those sorts); anything else fails closed
// before any solver runs. Probes first.

const admitMax = `mod m
  provides [m__max, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__max(left: int, right: int) -> M__Out rev 1
  emits []
  requires
    true
  ensures
    on Ok result
      result.value >= left
      result.value >= right
      match result.value == left
        true => true
        false => result.value == right
  tests
    ordered(left = 1, right = 2) => Ok(value = 2)
    reversed(left = 2, right = 1) => Ok(value = 2)
=
  match left <= right
    on true => Ok(value = right)
    on false => Ok(value = left)
`

const admitValidator = `mod m
  provides [m__check, M__Out]
  uses []
  emits [m.bad_bounds, m.out_of_range]

error m.bad_bounds(lower: int, upper: int)
error m.out_of_range(value: int, lower: int, upper: int)

type M__Out rev 1 (
  value: int
)

fn m__check(value: int, lower: int, upper: int) -> M__Out rev 1
  ensures
    on Ok result
      result.value == value
    on m.bad_bounds err
      lower >= upper + 1
      err.lower == lower
    on m.out_of_range err
      err.value == value
  emits [m.bad_bounds, m.out_of_range]
  tests
    ok(value = 1, lower = 0, upper = 2) => Ok(value = 1)
    low(value = -1, lower = 0, upper = 2) => m.out_of_range(value = -1, lower = 0, upper = 2)
=
  match value >= lower
    on true => Ok(value = value)
    on false => m.out_of_range(value = value, lower = lower, upper = upper)
`

func admitProg(t *testing.T, src string) (*Program, map[string]string) {
	t.Helper()
	return revisionProg(t, map[string]string{"m.ail": src}, []string{"m.ail"})
}

// TestAdmitPilots pins the open gate: both public pilots lie inside
// the fragment and report no admission finding.
func TestAdmitPilots(t *testing.T) {
	for name, src := range map[string]string{"max": admitMax, "validator": admitValidator} {
		prog, texts := admitProg(t, src)
		if diags := CheckContractAdmission(prog, texts); len(diags) != 0 {
			t.Fatalf("%s: expected admission, got %v", name, diags)
		}
	}
}

// TestAdmitEnsuresFalse pins verdict 3C: a declared but unrealized
// error with `ensures false` is admitted, not rejected.
func TestAdmitEnsuresFalse(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits [m.never]

error m.never()

type M__Out rev 1 (
  value: int
)

fn m__go(x: int) -> M__Out rev 1
  emits [m.never]
  requires
    true
  ensures
    on Ok result
      result.value == x
    on m.never err
      false
  tests
    go(x = 1) => Ok(value = 1)
=
  Ok(value = x)
`
	prog, texts := admitProg(t, src)
	if diags := CheckContractAdmission(prog, texts); len(diags) != 0 {
		t.Fatalf("expected admission, got %v", diags)
	}
}

// TestAdmitRequiresNonBool pins malformed: a requires row must be a
// Boolean predicate.
func TestAdmitRequiresNonBool(t *testing.T) {
	prog, texts := admitProg(t, strings.Replace(admitMax, "  requires\n    true\n", "  requires\n    1\n", 1))
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractMalformed) {
		t.Fatalf("non-Boolean requires reported no malformed finding")
	}
}

// TestAdmitUnknownBinder pins malformed: predicates may only name
// the arm binder and the function params.
func TestAdmitUnknownBinder(t *testing.T) {
	prog, texts := admitProg(t, strings.Replace(admitMax, "result.value >= left", "result2.value >= left", 1))
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractMalformed) {
		t.Fatalf("unknown binder reported no malformed finding")
	}
}

// TestAdmitMissingOutcome pins annotation completeness: every
// declared emits kind needs exactly one ensures arm.
func TestAdmitMissingOutcome(t *testing.T) {
	lines := strings.Split(admitValidator, "\n")
	var kept []string
	skip := false
	for _, l := range lines {
		if strings.TrimSpace(l) == "on m.out_of_range err" {
			skip = true
			continue
		}
		if skip && (strings.HasPrefix(l, "    on ") || strings.HasPrefix(l, "  emits")) {
			skip = false
		}
		if !skip {
			kept = append(kept, l)
		}
	}
	prog, texts := admitProg(t, strings.Join(kept, "\n"))
	diags := CheckContractAdmission(prog, texts)
	if !hasCode(diags, CodeContractMalformed) {
		t.Fatalf("missing outcome arm reported no malformed finding")
	}
	if !hasFound(diags, "m.out_of_range") {
		t.Fatalf("missing-outcome finding names no kind: %v", diags)
	}
}

// TestAdmitDuplicateOutcome pins malformed: one outcome, one arm.
func TestAdmitDuplicateOutcome(t *testing.T) {
	dup := strings.Replace(admitMax, "      match result.value == left",
		"      result.value == left\n    on Ok again\n      result.value == left\n      match result.value == left", 1)
	prog, texts := admitProg(t, dup)
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractMalformed) {
		t.Fatalf("duplicate outcome arm reported no malformed finding")
	}
}

// TestAdmitMatchArmSort pins malformed Boolean cases: every arm of
// an ensures match must return a Boolean predicate.
func TestAdmitMatchArmSort(t *testing.T) {
	bad := strings.Replace(admitMax, "        false => result.value == right",
		"        false => Ok(value = 1)", 1)
	prog, texts := admitProg(t, bad)
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractMalformed) {
		t.Fatalf("non-Boolean match arm reported no malformed finding")
	}
}

// TestAdmitNonlinear pins the verdict counterexample: an
// unsupported clause rejects the complete contract, never a
// partial proof of the remaining clauses.
func TestAdmitNonlinear(t *testing.T) {
	src := `mod m
  provides [m__sq, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__sq(x: int) -> M__Out rev 1
  emits []
  requires
    true
  ensures
    on Ok result
      result.value >= 0
      result.value == x * x
  tests
    zero(x = 0) => Ok(value = 0)
=
  Ok(value = 0)
`
	prog, texts := admitProg(t, src)
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractUnsupported) {
		t.Fatalf("nonlinear clause reported no unsupported finding")
	}
}

// TestAdmitDivision pins unsupported: not even division by a
// constant enters the first cut.
func TestAdmitDivision(t *testing.T) {
	bad := strings.Replace(admitMax, "      result.value >= right",
		"      result.value >= right\n      result.value == left / 2", 1)
	prog, texts := admitProg(t, bad)
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractUnsupported) {
		t.Fatalf("division reported no unsupported finding")
	}
}

// TestAdmitStringSort pins unsupported sorts: str predicates have
// no verification theory in this cut.
func TestAdmitStringSort(t *testing.T) {
	src := `mod m
  provides [m__tag, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  name: str
)

fn m__tag(x: int) -> M__Out rev 1
  emits []
  requires
    true
  ensures
    on Ok result
      result.name == "x"
  tests
    go(x = 1) => Ok(name = "x")
=
  Ok(name = "x")
`
	prog, texts := admitProg(t, src)
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractUnsupported) {
		t.Fatalf("string predicate reported no unsupported finding")
	}
}

// TestAdmitBrandSort pins unsupported sorts: even same-brand
// inputs are outside the fragment.
func TestAdmitBrandSort(t *testing.T) {
	src := `mod m
  provides [m__go, M__B, M__Out]
  uses []
  emits []

brand M__B is str rev 1

type M__Out rev 1 (
  value: int
)

fn m__go(tok: M__B) -> M__Out rev 1
  emits []
  requires
    true
  ensures
    on Ok result
      result.value == 1
  tests
    go(tok = seal M__B("a")) => Ok(value = 1)
=
  Ok(value = 1)
`
	prog, texts := admitProg(t, src)
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractUnsupported) {
		t.Fatalf("brand-typed input reported no unsupported finding")
	}
}

// TestAdmitDecSort pins unsupported sorts: dec has no theory.
func TestAdmitDecSort(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(x: dec) -> M__Out rev 1
  emits []
  requires
    x == d"1.5"
  ensures
    on Ok result
      result.value == 1
  tests
    go(x = d"1.5") => Ok(value = 1)
=
  Ok(value = 1)
`
	prog, texts := admitProg(t, src)
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractUnsupported) {
		t.Fatalf("dec predicate reported no unsupported finding")
	}
}

// TestAdmitExternCall pins the trust hole: a reachable extern call
// rejects the contracted function, declaration or not.
func TestAdmitExternCall(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out, ext__thing]
  uses []
  emits []

type M__Out rev 1 (
  n: int
)

extern ext__thing(x: int) -> M__Out rev 1
  emits []

fn m__go(x: int) -> M__Out rev 1
  emits []
  requires
    true
  ensures
    on Ok result
      result.n == x
  tests
    go(x = 1) => Ok(n = 1)
=
  match call ext__thing(x)
    given
      go => [exchange args (x = 1) outcome Ok(n = 1)]
    on Ok ok => Ok(n = ok.n)
`
	prog, texts := admitProg(t, src)
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractUnsupported) {
		t.Fatalf("extern call reported no unsupported finding")
	}
}

// TestAdmitRecursion pins the induction gap: a terminating,
// provable loop is still rejected without an induction rule.
func TestAdmitRecursion(t *testing.T) {
	src := `mod m
  provides [m__poll, M__S]
  uses []
  emits []

type M__S rev 1 (
  n: int
)

fn m__poll(n: int) -> M__S rev 1
  decreases n
  emits []
  requires
    n >= 0
  ensures
    on Ok s
      s.n == 0
  tests
    now(n = 0) => Ok(n = 0)
    later(n = 2) => Ok(n = 0)
=
  match n <= 0
    true => Ok(n = 0)
    false => match call m__poll(n - 1)
      on Ok s => Ok(n = s.n)
`
	prog, texts := admitProg(t, src)
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractUnsupported) {
		t.Fatalf("recursive proof reported no unsupported finding")
	}
}

// TestAdmitContractedCallee pins unavailable summaries: calling a
// contracted function proves nothing until that callee is verified
// in the proving run, so the caller is not admitted now.
func TestAdmitContractedCallee(t *testing.T) {
	src := `mod m
  provides [m__one, m__two, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__one() -> M__Out rev 1
  emits []
  requires
    true
  ensures
    on Ok result
      result.value == 1
  tests
    go() => Ok(value = 1)
=
  Ok(value = 1)

fn m__two() -> M__Out rev 1
  emits []
  requires
    true
  ensures
    on Ok result
      result.value == 1
  tests
    go() => Ok(value = 1)
=
  match call m__one()
    on Ok r => Ok(value = r.value)
`
	prog, texts := admitProg(t, src)
	diags := CheckContractAdmission(prog, texts)
	if !hasCode(diags, CodeContractUnverifiedDep) {
		t.Fatalf("contracted callee reported no unavailable finding")
	}
	for _, d := range diags {
		if d.Code == CodeContractUnverifiedDep && strings.Contains(d.Msg+d.Expected+d.Found, "m__one") &&
			!strings.Contains(d.Msg+d.Expected+d.Found, "m__two") {
			t.Fatalf("finding attributes the callee instead of the caller: %v", d)
		}
	}
}

// TestAdmitUncontractedCallee pins unstated summaries: a verified
// caller may not rely on a callee that carries no contract.
func TestAdmitUncontractedCallee(t *testing.T) {
	src := `mod m
  provides [m__one, m__two, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__one() -> M__Out rev 1
  emits []
  tests
    go() => Ok(value = 1)
=
  Ok(value = 1)

fn m__two() -> M__Out rev 1
  emits []
  requires
    true
  ensures
    on Ok result
      result.value == 1
  tests
    go() => Ok(value = 1)
=
  match call m__one()
    on Ok r => Ok(value = r.value)
`
	prog, texts := admitProg(t, src)
	if diags := CheckContractAdmission(prog, texts); !hasCode(diags, CodeContractUnverifiedDep) {
		t.Fatalf("uncontracted callee reported no unavailable finding")
	}
}

// TestAdmitSkipsUncontracted pins scope: functions without
// contracts keep ordinary status and report nothing.
func TestAdmitSkipsUncontracted(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  value: int
)

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 2) => Ok(value = 4)
=
  Ok(value = x * x)
`
	prog, texts := admitProg(t, src)
	if diags := CheckContractAdmission(prog, texts); len(diags) != 0 {
		t.Fatalf("uncontracted function reported findings: %v", diags)
	}
}
