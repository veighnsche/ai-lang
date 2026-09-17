package main

import (
	"testing"
)

// v15: executable code mints only its own module's brands. The
// declaring file owns every executable seal site; tests and given
// rows may name any declared brand because they are checked data.
// A forged brand in a body is AIL6004; the same seal in a test or
// script row is clean.

const sealOwnBody = `mod m
  provides [m__mint, M__B, M__Out]
  uses []
  emits []

brand M__B is str rev 1

type M__Out rev 1 (
  echo: M__B
)

fn m__mint(which: str) -> M__Out rev 1
  emits []
  tests
    mint(which = "a") => Ok(echo = seal M__B("a"))
    other(which = "b") => Ok(echo = seal M__B("b"))
=
  match which
    "a" => Ok(echo = seal M__B("a"))
    _ => Ok(echo = seal M__B("b"))
`

const sealLib = `mod lib
  provides [lib__get, Lib__B, Lib__Out]
  uses []
  emits []

brand Lib__B is str rev 1

type Lib__Out rev 1 (
  echo: Lib__B
)

fn lib__get() -> Lib__Out rev 1
  emits []
  tests
    go() => Ok(echo = seal Lib__B("x"))
=
  Ok(echo = seal Lib__B("x"))
`

func TestSealOwnBrandInBody(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.ail": sealOwnBody})
	if diags := diagnose(dir, "m.ail", sealOwnBody); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestSealLibOwnBrand(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"lib.ail": sealLib})
	if diags := diagnose(dir, "lib.ail", sealLib); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestSealForeignBrandInBody(t *testing.T) {
	forge := `mod app
  provides [app__forge, App__Out]
  uses [lib__get@1, Lib__B@1]
  emits []

type App__Out rev 1 (
  echo: Lib__B
)

fn app__forge() -> App__Out rev 1
  emits []
  tests
    go() => Ok(echo = seal Lib__B("x"))
=
  match call lib__get()
    given
      go => [exchange args () outcome Ok(echo = "x")]
    on Ok got => Ok(echo = seal Lib__B("minted"))
`
	dir := writeLSPDir(t, map[string]string{"lib.ail": sealLib, "app.ail": forge})
	diags := diagnose(dir, "app.ail", forge)
	if !hasDiag(diags, "error", "bodies seal only their own module's brands") {
		t.Fatalf("expected AIL6004 foreign-seal error, got %v", diags)
	}
}

func TestSealForeignBrandInData(t *testing.T) {
	// The same foreign seal in test and script rows is checked data,
	// not code: no diagnostic.
	uses := `mod app
  provides [app__pass, App__Out]
  uses [lib__get@1, Lib__B@1]
  emits []

type App__Out rev 1 (
  echo: Lib__B
)

fn app__pass() -> App__Out rev 1
  emits []
  tests
    go() => Ok(echo = seal Lib__B("x"))
=
  match call lib__get()
    given
      go => [exchange args () outcome Ok(echo = seal Lib__B("x"))]
    on Ok got => Ok(echo = got.echo)
`
	dir := writeLSPDir(t, map[string]string{"lib.ail": sealLib, "app.ail": uses})
	if diags := diagnose(dir, "app.ail", uses); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
