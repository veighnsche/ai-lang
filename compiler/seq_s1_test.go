package main

import (
	"strings"
	"testing"
)

// a36 S1: typed sequence construction. Finite immutable Seq values
// with explicit element types, including the empty case. Decision
// tables below are the P0 acceptance rows: V-rows must diagnose
// clean (a wrong expectation fails the test run, so clean means the
// evaluator computed exactly the pinned value), T-rows must carry
// the named diagnostic code.

// seqClean asserts a module with Seq fixtures builds with no
// diagnostics: parse, element checking, test execution, and
// exhaustiveness all agree.
func seqClean(t *testing.T, files map[string]string, name string) {
	t.Helper()
	dir := writeLSPDir(t, files)
	body := files[name]
	if diags := diagnose(dir, name, body); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// seqProbeClean scopes the clean assertion for probes appended to real
// blessed files (bytes round-trip/probe tests): blessed scaffolding
// stays compiler-clean (any non-lint diagnostic fails anywhere) but
// its lint is grandfathered, while the appended probe — new code —
// must be fully clean, lint included. probeStart is the probe's
// first 1-based line in the assembled text.
func seqProbeClean(t *testing.T, files map[string]string, name string, probeStart int) {
	t.Helper()
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, name, files[name])
	for _, d := range diags {
		if strings.HasPrefix(d.Code, "CAN341") && d.Line < probeStart {
			continue
		}
		t.Fatalf("unexpected diagnostic (probe starts at line %d): %v", probeStart, diags)
	}
}

// seqCode asserts the rejection row fires precisely: exactly one error
// carries the named code and fragment, and every other error is a
// runtime confirmation (a failed test run or a contradicted script).
// Those confirmations are welcome: they prove the test really
// executed instead of passing vacuously. Anything else is a cascade.
func seqCode(t *testing.T, files map[string]string, name, code, sub string) {
	t.Helper()
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, name, files[name])
	primary := 0
	for _, d := range diags {
		if d.Sev != "error" {
			continue
		}
		switch {
		case d.Code == code && strings.Contains(d.Msg, sub):
			primary++
		case d.Code == CodeTestFailed || d.Code == CodeInconsistentScript:
		default:
			t.Fatalf("unexpected cascade error, got %v", diags)
		}
	}
	if primary != 1 {
		t.Fatalf("expected exactly one %s containing %q, got %v", code, sub, diags)
	}
}

const seqOut = `type M__Out rev 1 (
  vals: Seq<str>
)
`

// V0: empty is a value, not a failure.
const seqV0 = `mod m
  provides [m__go, M__Out]
  uses []
  emits []

` + seqOut + `
fn m__go() -> M__Out rev 1
  emits []
  tests
    go() => Ok(vals = Seq<str>[])
  Ok(vals = Seq<str>[])
`

// V1: singleton empty string is observably distinct from V0
// (length 1, not 0; S2 will measure the difference).
const seqV1 = `mod m
  provides [m__go, M__Out]
  uses []
  emits []

` + seqOut + `
fn m__go() -> M__Out rev 1
  emits []
  tests
    go() => Ok(vals = Seq<str>[""])
  Ok(vals = Seq<str>[""])
`

// V2: order, empties, and repeats preserved exactly.
const seqV2 = `mod m
  provides [m__go, M__Out]
  uses []
  emits []

` + seqOut + `
fn m__go() -> M__Out rev 1
  emits []
  tests
    go() => Ok(vals = Seq<str>["b", "", "a", "b"])
  Ok(vals = Seq<str>["b", "", "a", "b"])
`

// V2b: members containing commas parse through top-level splitting.
const seqV2b = `mod m
  provides [m__go, M__Out]
  uses []
  emits []

` + seqOut + `
fn m__go() -> M__Out rev 1
  emits []
  tests
    go() => Ok(vals = Seq<str>["a,b", "c"])
  Ok(vals = Seq<str>["a,b", "c"])
`

// V3: branded members keep their static identity.
const seqV3 = `mod m
  provides [m__go, M__Safe, M__SeqOut]
  uses []
  emits []

brand M__Safe is str rev 1

type M__SeqOut rev 1 (
  vals: Seq<M__Safe>
)

fn m__go() -> M__SeqOut rev 1
  emits []
  tests
    go() => Ok(vals = Seq<M__Safe>[seal M__Safe("A"), seal M__Safe("")])
  Ok(vals = Seq<M__Safe>[seal M__Safe("A"), seal M__Safe("")])
`

// V4: a second brand family in the same module stays distinct.
const seqV4 = `mod m
  provides [m__go, M__Attr, M__AttrOut]
  uses []
  emits []

brand M__Attr is str rev 1

type M__AttrOut rev 1 (
  vals: Seq<M__Attr>
)

fn m__go() -> M__AttrOut rev 1
  emits []
  tests
    go() => Ok(vals = Seq<M__Attr>[seal M__Attr(""), seal M__Attr("disabled")])
  Ok(vals = Seq<M__Attr>[seal M__Attr(""), seal M__Attr("disabled")])
`

// V5: the sequence rides inside a wider result record.
const seqV5 = `mod m
  provides [m__go, M__Wide]
  uses []
  emits []

type M__Wide rev 1 (
  name: str,
  vals: Seq<str>
)

fn m__go() -> M__Wide rev 1
  emits []
  tests
    go() => Ok(name = "n", vals = Seq<str>["a", ""])
  Ok(name = "n", vals = Seq<str>["a", ""])
`

func TestSeqValuesClean(t *testing.T) {
	for _, c := range []struct {
		name, body string
	}{
		{"V0 empty", seqV0},
		{"V1 singleton empty", seqV1},
		{"V2 order empties repeats", seqV2},
		{"V2b comma member", seqV2b},
		{"V3 branded identity", seqV3},
		{"V4 second brand", seqV4},
		{"V5 record carrier", seqV5},
	} {
		t.Run(c.name, func(t *testing.T) {
			seqClean(t, map[string]string{"m.can": c.body}, "m.can")
		})
	}
}

// The V-rows above prove nothing if expectations never execute: a
// dropped empty and a reordered pair must both fail the test run.
func TestSeqExpectationsExecute(t *testing.T) {
	drop := strings.Replace(seqV2,
		`go() => Ok(vals = Seq<str>["b", "", "a", "b"])`,
		`go() => Ok(vals = Seq<str>["b", "a", "b"])`, 1)
	seqCode(t, map[string]string{"m.can": drop}, "m.can", CodeTestFailed, "go")

	flip := strings.Replace(seqV2,
		`go() => Ok(vals = Seq<str>["b", "", "a", "b"])`,
		`go() => Ok(vals = Seq<str>["a", "", "b", "b"])`, 1)
	seqCode(t, map[string]string{"m.can": flip}, "m.can", CodeTestFailed, "go")
}

// T1: an int member in Seq<str> is CAN6003 in every value position
// with its own checking path.
func TestSeqT1Positions(t *testing.T) {
	body := strings.Replace(seqV0, `Ok(vals = Seq<str>[])`,
		`Ok(vals = Seq<str>[1])`, 1)
	// Calls ride match scrutinees only (v0 subset): the Seq argument
	// checks against the callee parameter at the call site.
	call := `mod m
  provides [m__go, m__id, M__Out]
  uses []
  emits []

` + seqOut + `
fn m__id(xs: Seq<str>) -> M__Out rev 1
  emits []
  tests
    go(Seq<str>[]) => Ok(vals = Seq<str>[])
  Ok(vals = xs)

fn m__go(dummy: str) -> M__Out rev 1
  emits []
  tests
    go(dummy = "d") => Ok(vals = Seq<str>[])
  match call m__id(Seq<str>[1])
    on Ok v => Ok(vals = v.vals)
`
	targ := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

` + seqOut + `
fn m__go(xs: Seq<str>) -> M__Out rev 1
  emits []
  tests
    go(xs = Seq<str>[1]) => Ok(vals = Seq<str>[])
  Ok(vals = xs)
`
	exp := strings.Replace(seqV0,
		`go() => Ok(vals = Seq<str>[])`,
		`go() => Ok(vals = Seq<str>[1])`, 1)
	for _, c := range []struct {
		name, body string
	}{
		{"body literal", body},
		{"call argument", call},
		{"test argument", targ},
		{"expected result", exp},
	} {
		t.Run(c.name, func(t *testing.T) {
			seqCode(t, map[string]string{"m.can": c.body}, "m.can",
				CodeTypeMismatch, "want str")
		})
	}
}

// T1 given position: a scripted Seq argument with a wrong member.
const seqLibEcho = `mod lib
  provides [lib__echo, Lib__Out]
  uses []
  emits []

type Lib__Out rev 1 (
  vals: Seq<str>
)

fn lib__echo(x: str) -> Lib__Out rev 1
  emits []
  tests
    go("a") => Ok(vals = Seq<str>["a", ""])
  Ok(vals = Seq<str>[x, ""])
`

// T1 given position: the provider takes a Seq, and the scripted
// argument smuggles a wrong member. The static element rejection is
// primary; the contradicted script follows as runtime confirmation.
const seqLibTake = `mod lib
  provides [lib__take, Lib__Out]
  uses []
  emits []

type Lib__Out rev 1 (
  vals: Seq<str>
)

fn lib__take(xs: Seq<str>) -> Lib__Out rev 1
  emits []
  tests
    go(Seq<str>[]) => Ok(vals = Seq<str>[])
  Ok(vals = xs)
`

func TestSeqT1Given(t *testing.T) {
	app := `mod app
  provides [app__go, App__Out]
  uses [lib__take@1]
  emits []

type App__Out rev 1 (
  vals: Seq<str>
)

fn app__go(x: str) -> App__Out rev 1
  emits []
  tests
    go("a") => Ok(vals = Seq<str>["a", ""])
  match call lib__take(Seq<str>[x, ""])
    given
      go => [exchange args (xs = Seq<str>[1]) outcome Ok(vals = Seq<str>["a", ""])]
    on Ok v => Ok(vals = v.vals)
`
	seqCode(t, map[string]string{"lib.can": seqLibTake, "app.can": app},
		"app.can", CodeTypeMismatch, "want str")
}

// T2: a branded member in Seq<str> is CAN6003.
func TestSeqT2BrandInStr(t *testing.T) {
	body := `mod m
  provides [m__go, M__B, M__Out]
  uses []
  emits []

brand M__B is str rev 1

` + seqOut + `
fn m__go() -> M__Out rev 1
  emits []
  tests
    go() => Ok(vals = Seq<str>[])
  Ok(vals = Seq<str>[seal M__B("x")])
`
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeTypeMismatch, "want str")
}

// T3: a cross-brand member is CAN6003.
func TestSeqT3CrossBrand(t *testing.T) {
	body := `mod m
  provides [m__go, M__A, M__B, M__SeqOut]
  uses []
  emits []

brand M__A is str rev 1

brand M__B is str rev 1

type M__SeqOut rev 1 (
  vals: Seq<M__B>
)

fn m__go() -> M__SeqOut rev 1
  emits []
  tests
    go() => Ok(vals = Seq<M__B>[])
  Ok(vals = Seq<M__B>[seal M__A("x")])
`
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeTypeMismatch, "want M__B")
}

// T4: emptiness never waives the element check.
func TestSeqT4EmptyTypedMismatch(t *testing.T) {
	body := `mod m
  provides [m__go, M__B, M__Out]
  uses []
  emits []

brand M__B is str rev 1

` + seqOut + `
fn m__go() -> M__Out rev 1
  emits []
  tests
    go() => Ok(vals = Seq<M__B>[])
  Ok(vals = Seq<str>[])
`
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeTypeMismatch, "want Seq<str>")
}

// T5: a foreign seal inside an executable literal is CAN6004.
const seqBrandLib = `mod lib
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
  Ok(echo = seal Lib__B("x"))
`

func TestSeqT5ForeignSeal(t *testing.T) {
	// Brand-only pin: the module mints nothing, it only names the
	// brand inside a literal, which is exactly the forgery.
	app := `mod app
  provides [app__forge, App__Out]
  uses [Lib__B@1]
  emits []

type App__Out rev 1 (
  vals: Seq<Lib__B>
)

fn app__forge() -> App__Out rev 1
  emits []
  tests
    go() => Ok(vals = Seq<Lib__B>[seal Lib__B("x")])
  Ok(vals = Seq<Lib__B>[seal Lib__B("x")])
`
	seqCode(t, map[string]string{"lib.can": seqBrandLib, "app.can": app},
		"app.can", CodeSealForeign, "Lib__B")
}

// T6: the same foreign seal in test data is checked data, not code.
func TestSeqT6ForeignSealInData(t *testing.T) {
	app := `mod app
  provides [app__go, App__Out]
  uses [Lib__B@1]
  emits []

type App__Out rev 1 (
  vals: Seq<Lib__B>
)

fn app__go(xs: Seq<Lib__B>) -> App__Out rev 1
  emits []
  tests
    go(Seq<Lib__B>[seal Lib__B("x")]) => Ok(vals = Seq<Lib__B>[seal Lib__B("x")])
  Ok(vals = xs)
`
	seqClean(t, map[string]string{"lib.can": seqBrandLib, "app.can": app}, "app.can")
}

// T7: a bare [...] in a value position names the typed form.
func TestSeqT7BareList(t *testing.T) {
	body := strings.Replace(seqV0, `Ok(vals = Seq<str>[])`, `Ok(vals = ["a"])`, 1)
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeSeqLiteral, "Seq<T>")
}

// T8: a bare-Seq return is unsupported, like a bare-brand return.
func TestSeqT8BareReturn(t *testing.T) {
	body := `mod m
  provides [m__go]
  uses []
  emits []

fn m__go() -> Seq<str> rev 1
  emits []
  tests
    go() => Ok(vals = Seq<str>[])
  Ok(vals = Seq<str>[])
`
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeTypeMismatch, "bare-Seq")
}

// Unknown element types fail at the annotation, not downstream: the
// body alone names the bad type while the expectation stays clean,
// so the single static error plus the executed mismatch pin both.
func TestSeqUnknownElem(t *testing.T) {
	body := strings.Replace(seqV0, "  Ok(vals = Seq<str>[])", "  Ok(vals = Seq<Nope>[\"a\"])", 1)
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeUnknownType, "Nope")
}

// Malformed literal shapes are parse errors, never binop cascades.
func TestSeqParseShapes(t *testing.T) {
	for _, c := range []struct{ name, frag string }{
		{"empty elem", "Seq<>[]"},
		{"nested", "Seq<Seq<str>>[Seq<str>[\"a\"]]"},
		{"missing bracket", "Seq<str>"},
		{"unbalanced", "Seq<str>[\"a\""},
	} {
		t.Run(c.name, func(t *testing.T) {
			body := strings.Replace(seqV0, `Seq<str>[]`, c.frag, 1)
			dir := writeLSPDir(t, map[string]string{"m.can": body})
			diags := diagnose(dir, "m.can", body)
			found := false
			for _, d := range diags {
				if d.Sev == "error" && d.Code == CodeParse {
					found = true
				}
			}
			if !found {
				t.Fatalf("expected parse error, got %v", diags)
			}
		})
	}
}

// Sequence equality is not a v1 surface, even when both sides agree.
func TestSeqNoEquality(t *testing.T) {
	body := `mod m
  provides [m__go, M__Flag]
  uses []
  emits []

type M__Flag rev 1 (
  flag: bool
)

fn m__go() -> M__Flag rev 1
  emits []
  tests
    go() => Ok(flag = true)
  Ok(flag = (Seq<str>["a"] == Seq<str>["a"]))
`
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeTypeMismatch, "sequence equality")
}

// Linkage: a scripted Seq that drops a field contradicts the provider.
func TestSeqLinkageDrop(t *testing.T) {
	app := `mod app
  provides [app__go, App__Out]
  uses [lib__echo@1]
  emits []

type App__Out rev 1 (
  vals: Seq<str>
)

fn app__go(x: str) -> App__Out rev 1
  emits []
  tests
    go("a") => Ok(vals = Seq<str>["a"])
  match call lib__echo(x)
    given
      go => [exchange args (x = "a") outcome Ok(vals = Seq<str>["a"])]
    on Ok v => Ok(vals = v.vals)
`
	seqCode(t, map[string]string{"lib.can": seqLibEcho, "app.can": app},
		"app.can", CodeInconsistentScript, "contradicts lib__echo")
}

// Linkage: a scripted Seq that reorders fields contradicts too, and
// the truthful script stays clean.
func TestSeqLinkageOrder(t *testing.T) {
	mkapp := func(want string) string {
		return `mod app
  provides [app__go, App__Out]
  uses [lib__echo@1]
  emits []

type App__Out rev 1 (
  vals: Seq<str>
)

fn app__go(x: str) -> App__Out rev 1
  emits []
  tests
    go("a") => Ok(vals = ` + want + `)
  match call lib__echo(x)
    given
      go => [exchange args (x = "a") outcome Ok(vals = ` + want + `)]
    on Ok v => Ok(vals = v.vals)
`
	}
	seqCode(t,
		map[string]string{"lib.can": seqLibEcho, "app.can": mkapp(`Seq<str>["", "a"]`)},
		"app.can", CodeInconsistentScript, "contradicts lib__echo")
	seqClean(t,
		map[string]string{"lib.can": seqLibEcho, "app.can": mkapp(`Seq<str>["a", ""]`)},
		"app.can")
}

// Emit: literals lower to array literals over the erased element type.
func TestSeqEmitArray(t *testing.T) {
	ts := compileEmit(t, seqV2)
	for _, want := range []string{
		`["b", "", "a", "b"]`,
		`string[]`,
	} {
		if !strings.Contains(ts, want) {
			t.Fatalf("emit missing %q:\n%s", want, ts)
		}
	}
}

func TestSeqEmitBrandErases(t *testing.T) {
	ts := compileEmit(t, seqV3)
	for _, want := range []string{
		`["A", ""]`,
		`string[]`,
	} {
		if !strings.Contains(ts, want) {
			t.Fatalf("emit missing %q:\n%s", want, ts)
		}
	}
}
