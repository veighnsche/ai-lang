package main

import (
	"strings"
	"testing"
)

// Helpers are same-file calls that execute instead of stubbing: the
// call site carries no given table and no uses pin, and inner tables
// script every reaching test by name.

const helperClean = `mod m
  provides [m__go, m__help, M__S]
  uses []
  emits [m.bad]

error m.bad()

type M__S rev 1 (
  id: str
)

fn m__help(id: str) -> M__S rev 1
  emits [m.bad]
  tests
    h_ok(id = "a") => Ok(id = "a")
    h_bad(id = "b") => m.bad()
=
  match id
    "a" => Ok(id = "a")
    _ => m.bad()

fn m__go(id: str) -> M__S rev 1
  emits [m.bad]
  tests
    g_ok(id = "a") => Ok(id = "a")
    g_bad(id = "b") => m.bad()
=
  match call m__help(id)
    on m.bad _ => m.bad()
    on Ok s => Ok(id = s.id)
`

func TestHelperClean(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": helperClean})
	if diags := diagnose(dir, "m.can", helperClean); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestHelperGivenIsError(t *testing.T) {
	bad := strings.Replace(helperClean,
		"  match call m__help(id)\n    on m.bad _ => m.bad()",
		"  match call m__help(id)\n    given\n      g_ok => [exchange args (id = \"a\") outcome Ok(id = \"a\")]\n      g_bad => [exchange args (id = \"a\") outcome m.bad()]\n    on m.bad _ => m.bad()", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "takes no given table") {
		t.Fatalf("expected given-on-local error, got %v", diags)
	}
}

func TestHelperSelfCycle(t *testing.T) {
	loop := `mod m
  provides [m__loop, M__S]
  uses []
  emits [m.bad]

error m.bad()

type M__S rev 1 (
  id: str
)

fn m__loop(id: str) -> M__S rev 1
  emits [m.bad]
  tests
    l1(id = "a") => m.bad()
=
  match call m__loop(id)
    on m.bad _ => m.bad()
    on Ok s => Ok(id = s.id)
`
	dir := writeLSPDir(t, map[string]string{"m.can": loop})
	diags := diagnose(dir, "m.can", loop)
	if !hasDiag(diags, "error", "local call cycle m__loop -> m__loop") {
		t.Fatalf("expected self-cycle error, got %v", diags)
	}
}

func TestHelperMutualCycle(t *testing.T) {
	pair := `mod m
  provides [m__a, m__b, M__S]
  uses []
  emits [m.bad]

error m.bad()

type M__S rev 1 (
  id: str
)

fn m__a(id: str) -> M__S rev 1
  emits [m.bad]
  tests
    a1(id = "a") => m.bad()
=
  match call m__b(id)
    on m.bad _ => m.bad()
    on Ok s => Ok(id = s.id)

fn m__b(id: str) -> M__S rev 1
  emits [m.bad]
  tests
    b1(id = "a") => m.bad()
=
  match call m__a(id)
    on m.bad _ => m.bad()
    on Ok s => Ok(id = s.id)
`
	dir := writeLSPDir(t, map[string]string{"m.can": pair})
	diags := diagnose(dir, "m.can", pair)
	if !hasDiag(diags, "error", "local call cycle m__a -> m__b -> m__a") {
		t.Fatalf("expected mutual-cycle error, got %v", diags)
	}
}

func TestHelperNeedsNoPin(t *testing.T) {
	pinned := strings.Replace(helperClean, "uses []", "uses [m__help@1]", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": pinned})
	diags := diagnose(dir, "m.can", pinned)
	if !hasDiag(diags, "error", "resolves nowhere") {
		t.Fatalf("expected same-file uses rejection, got %v", diags)
	}
}

// A helper arm taken only through a caller test counts as covered:
// coverage is module-wide over green tests.
func TestHelperCoverageThroughCaller(t *testing.T) {
	flow := `mod m
  provides [m__go, m__help, M__S]
  uses []
  emits [m.bad]

error m.bad()

type M__S rev 1 (
  id: str
)

fn m__help(id: str) -> M__S rev 1
  emits [m.bad]
  tests
    h_ok(id = "a") => Ok(id = "a")
=
  match id
    "a" => Ok(id = "a")
    _ => m.bad()

fn m__go(id: str) -> M__S rev 1
  emits [m.bad]
  tests
    g_ok(id = "a") => Ok(id = "a")
    g_bad(id = "b") => m.bad()
=
  match call m__help(id)
    on m.bad _ => m.bad()
    on Ok s => Ok(id = s.id)
`
	dir := writeLSPDir(t, map[string]string{"m.can": flow})
	if diags := diagnose(dir, "m.can", flow); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// An inner given table must script every reaching test: the caller's
// tests flow through the helper under their own names.
func TestHelperInnerGivenNeedsCallerTests(t *testing.T) {
	needsDB := `mod db
  provides [db__get, Db__U]
  uses []
  emits [db.down]

error db.down()

type Db__U rev 1 (
  id: str
)

fn db__get(id: str) -> Db__U rev 1
  emits [db.down]
  tests
    ok(id = "u") => Ok(id = "u")
    other(id = "x") => db.down()
=
  match id
    "u" => Ok(id = "u")
    _ => db.down()
`
	needsGo := `mod m
  provides [m__go, m__help, M__S]
  uses [db__get@1]
  emits [m.bad]

error m.bad()

type M__S rev 1 (
  id: str
)

fn m__help(id: str) -> M__S rev 1
  emits [m.bad]
  tests
    h_ok(id = "a") => Ok(id = "a")
=
  match call db__get(id)
    given
      h_ok => [exchange args (id = "a") outcome Ok(id = "a")]
    on db.down _ => m.bad()
    on Ok u => Ok(id = u.id)

fn m__go(id: str) -> M__S rev 1
  emits [m.bad]
  tests
    g_ok(id = "a") => Ok(id = "a")
=
  match call m__help(id)
    on m.bad _ => m.bad()
    on Ok s => Ok(id = s.id)
`
	dir := writeLSPDir(t, map[string]string{"db.can": needsDB, "m.can": needsGo})
	diags := diagnose(dir, "m.can", needsGo)
	if !hasDiag(diags, "error", "test g_ok has no script at the call to db__get") {
		t.Fatalf("expected dangling caller-test error, got %v", diags)
	}
}
