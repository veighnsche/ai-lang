package main

import (
	"testing"
)

// Issue #1 (linkage): a scripted Ok outcome is a claim about what the
// provider computes. The masquerade pair below is the issue's
// acceptance test: an incorrect scripted result for a verified pure
// helper must fail the build.

const linkLib = `mod lib
  provides [lib__double]
  uses []
  emits []

fn lib__double(x: int) -> int rev 1
  tests
    eight(x = 4) => Ok(value = 8)
=
  match x
    _ => Ok(value = x * 2)
`

const linkAppLie = `mod app
  provides [app__go]
  uses [lib__double@1]
  emits []

fn app__go(x: int) -> int rev 1
  tests
    lie(x = 4) => Ok(value = 999)
=
  match call lib__double(x)
    given
      lie => [exchange args (x = 4) outcome Ok(value = 999)]
    on Ok v => Ok(value = v.value)
`

const linkAppTruth = `mod app
  provides [app__go]
  uses [lib__double@1]
  emits [app.boom]

error app.boom()

fn app__go(x: int) -> int rev 1
  emits [app.boom]
  tests
    truth(x = 4) => Ok(value = 8)
    failure(x = 4) => app.boom()
=
  match call lib__double(x)
    given
      truth => [exchange args (x = 4) outcome Ok(value = 8)]
      failure => [exchange args (x = 4) outcome lib.down()]
    on lib.down _ => app.boom()
    on Ok v => Ok(value = v.value)
`

func TestScriptContradictionFails(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"lib.ail": linkLib, "app.ail": linkAppLie})
	diags := diagnose(dir, "app.ail", linkAppLie)
	found := false
	for _, d := range diags {
		if d.Sev == "error" && d.Code == CodeInconsistentScript {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected AIL3110 contradiction, got %v", diags)
	}
	if !hasDiag(diags, "error", "contradicts lib__double") {
		t.Fatalf("expected contradiction message, got %v", diags)
	}
}

func TestScriptConsistencyClean(t *testing.T) {
	lib := `mod lib
  provides [lib__double]
  uses []
  emits [lib.down]

error lib.down()

fn lib__double(x: int) -> int rev 1
  emits [lib.down]
  tests
    eight(x = 4) => Ok(value = 8)
=
  match x
    _ => Ok(value = x * 2)
`
	dir := writeLSPDir(t, map[string]string{"lib.ail": lib, "app.ail": linkAppTruth})
	if diags := diagnose(dir, "app.ail", linkAppTruth); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// A provider the sandbox cannot model is trusted, never failed: here
// mid__go needs its own extern script to evaluate, so even a wrong Ok
// claim passes the check. The claim is still exercised at runtime by
// the caller's test.
func TestScriptUnmodelableTrusted(t *testing.T) {
	mid := `mod mid
  provides [mid__go, Mid__Out, ext__thing]
  uses []
  emits [mid.boom]

error mid.boom()

type Mid__Out rev 1 (
)

extern ext__thing(x: int) -> Mid__Out rev 1
  emits [mid.boom]

fn mid__go(x: int) -> Mid__Out rev 1
  emits [mid.boom]
  tests
    m1(x = 1) => Ok()
    m2(x = 2) => mid.boom()
=
  match call ext__thing(x)
    given
      m1 => [exchange args (x = 1) outcome Ok()]
      m2 => [exchange args (x = 2) outcome mid.boom()]
    on mid.boom _ => mid.boom()
    on Ok ok => Ok()
`
	app := `mod app
  provides [app__go]
  uses [mid__go@1, Mid__Out@1]
  emits [app.boom]

error app.boom()

fn app__go(x: int) -> Mid__Out rev 1
  emits [app.boom]
  tests
    t1(x = 1) => Ok()
    t2(x = 2) => app.boom()
=
  match call mid__go(x)
    given
      t1 => [exchange args (x = 1) outcome Ok()]
      t2 => [exchange args (x = 2) outcome mid.boom()]
    on mid.boom _ => app.boom()
    on Ok ok => Ok()
`
	dir := writeLSPDir(t, map[string]string{"mid.ail": mid, "app.ail": app})
	if diags := diagnose(dir, "app.ail", app); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
