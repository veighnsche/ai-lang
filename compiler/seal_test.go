package main

import (
	"testing"
)

// a15: executable code mints only its own module's brands. The
// declaring file owns every executable seal site; tests and given
// rows may name any declared brand because they are checked data.
// A forged brand in a body is CAN6004; the same seal in a test or
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
	dir := writeLSPDir(t, map[string]string{"m.can": sealOwnBody})
	if diags := diagnose(dir, "m.can", sealOwnBody); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestSealLibOwnBrand(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"lib.can": sealLib})
	if diags := diagnose(dir, "lib.can", sealLib); len(diags) != 0 {
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
	dir := writeLSPDir(t, map[string]string{"lib.can": sealLib, "app.can": forge})
	diags := diagnose(dir, "app.can", forge)
	if !hasDiag(diags, "error", "bodies seal only their own module's brands") {
		t.Fatalf("expected CAN6004 foreign-seal error, got %v", diags)
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
	dir := writeLSPDir(t, map[string]string{"lib.can": sealLib, "app.can": uses})
	if diags := diagnose(dir, "app.can", uses); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// a26: explicitly authorized one-way promotion. A destination brand
// names its admitted sources in seals_from; the seal site, source,
// and destination all belong to one module; representation is str on
// both sides by the v0 rule, and evaluation preserves the string by
// erasure. There is no reverse, transitive, or inferred promotion.

const sealPromoteOk = `mod m
  provides [m__node, M__Text, M__Safe, M__SafeResult]
  uses []
  emits []

brand M__Text is str rev 1

brand M__Safe is str rev 1 seals_from [M__Text]

type M__SafeResult rev 1 (
  safe: M__Safe
)

fn m__node(text: M__Text) -> M__SafeResult rev 1
  emits []
  tests
    plain(text = seal M__Text("hi")) => Ok(safe = seal M__Safe("hi"))
    escaped(text = seal M__Text("a&amp;b")) => Ok(safe = seal M__Safe("a&amp;b"))
=
  Ok(safe = seal M__Safe(text))
`

func TestSealPromotionAuthorized(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": sealPromoteOk})
	if diags := diagnose(dir, "m.can", sealPromoteOk); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

const sealPromoteDenied = `mod m
  provides [m__node, M__Text, M__Safe, M__SafeResult]
  uses []
  emits []

brand M__Text is str rev 1

brand M__Safe is str rev 1

type M__SafeResult rev 1 (
  safe: M__Safe
)

fn m__node(text: M__Text) -> M__SafeResult rev 1
  emits []
  tests
    plain(text = seal M__Text("hi")) => Ok(safe = seal M__Safe("hi"))
=
  Ok(safe = seal M__Safe(text))
`

func TestSealPromotionUnauthorized(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": sealPromoteDenied})
	diags := diagnose(dir, "m.can", sealPromoteDenied)
	if !hasDiag(diags, "error", "seal M__Safe cannot promote M__Text") {
		t.Fatalf("expected unauthorized-promotion error, got %v", diags)
	}
}

const sealPromoteReverse = `mod m
  provides [m__back, M__Text, M__Safe, M__BackResult]
  uses []
  emits []

brand M__Text is str rev 1

brand M__Safe is str rev 1 seals_from [M__Text]

type M__BackResult rev 1 (
  text: M__Text
)

fn m__back(frag: M__Safe) -> M__BackResult rev 1
  emits []
  tests
    plain(frag = seal M__Safe("hi")) => Ok(text = seal M__Text("hi"))
=
  Ok(text = seal M__Text(frag))
`

func TestSealPromotionNoReverse(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": sealPromoteReverse})
	diags := diagnose(dir, "m.can", sealPromoteReverse)
	if !hasDiag(diags, "error", "seal M__Text cannot promote M__Safe") {
		t.Fatalf("expected no-reverse-promotion error, got %v", diags)
	}
}

const sealPromoteTransit = `mod m
  provides [m__skip, M__A, M__B, M__C, M__CResult]
  uses []
  emits []

brand M__A is str rev 1

brand M__B is str rev 1 seals_from [M__A]

brand M__C is str rev 1 seals_from [M__B]

type M__CResult rev 1 (
  value: M__C
)

fn m__skip(start: M__A) -> M__CResult rev 1
  emits []
  tests
    plain(start = seal M__A("hi")) => Ok(value = seal M__C("hi"))
=
  Ok(value = seal M__C(start))
`

func TestSealPromotionNoTransitive(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": sealPromoteTransit})
	diags := diagnose(dir, "m.can", sealPromoteTransit)
	if !hasDiag(diags, "error", "seal M__C cannot promote M__A") {
		t.Fatalf("expected no-transitive-promotion error, got %v", diags)
	}
}

const sealPromoteUnknownSrc = `mod m
  provides [M__C]
  uses []
  emits []

brand M__C is str rev 1 seals_from [M__Ghost]
`

func TestSealPromotionUnknownSource(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": sealPromoteUnknownSrc})
	diags := diagnose(dir, "m.can", sealPromoteUnknownSrc)
	if !hasDiag(diags, "error", "seals_from unknown brand M__Ghost") {
		t.Fatalf("expected unknown-source error, got %v", diags)
	}
}

const sealPromoteXFile = `mod app
  provides [app__promote, App__C, App__Out]
  uses [lib__get@1, Lib__B@1]
  emits []

brand App__C is str rev 1 seals_from [Lib__B]

type App__Out rev 1 (
  echo: App__C
)

fn app__promote() -> App__Out rev 1
  emits []
  tests
    go() => Ok(echo = seal App__C("x"))
=
  match call lib__get()
    given
      go => [exchange args () outcome Ok(echo = "x")]
    on Ok got => Ok(echo = seal App__C(got.echo))
`

func TestSealPromotionCrossFile(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"lib.can": sealLib, "app.can": sealPromoteXFile})
	diags := diagnose(dir, "app.can", sealPromoteXFile)
	if !hasDiag(diags, "error", "promotions stay inside one module") {
		t.Fatalf("expected cross-file promotion errors, got %v", diags)
	}
}

const sealBrandReturn = `mod m
  provides [m__bad, M__B]
  uses []
  emits []

brand M__B is str rev 1

fn m__bad(which: str) -> M__B rev 1
  emits []
  tests
    go(which = "x") => Ok(value = seal M__B("x"))
=
  Ok(value = seal M__B("x"))
`

func TestFnBrandReturnRejected(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": sealBrandReturn})
	diags := diagnose(dir, "m.can", sealBrandReturn)
	if !hasDiag(diags, "error", "bare-brand returns are unsupported") {
		t.Fatalf("expected brand-return rejection, got %v", diags)
	}
}
