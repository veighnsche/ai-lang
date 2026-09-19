package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Declared finite monomorphic records work end to end: constructed in
// value positions, passed as parameters, nested in fields, read
// through projections, and compared inside outcomes — with constructor
// identity, exact fields, and field types enforced before tests.

const recordFixture = `mod audit
  provides [audit__main, audit__use, audit__wrap, audit__box, audit__eq, Audit__Attrs, Audit__Page, Audit__Box, Audit__Out]
  uses []
  emits []

type Audit__Attrs rev 1 (
  id: str
)

type Audit__Page rev 1 (
  attrs: Audit__Attrs
  title: str
)

type Audit__Box rev 1 (
  attrs: Audit__Attrs
)

type Audit__Out rev 1 (
  value: str
)

fn audit__use(attrs: Audit__Attrs) -> Audit__Out rev 1
  emits []
  tests
    u(Audit__Attrs("main")) => Ok("main")
  Ok(attrs.id)

fn audit__wrap(page: Audit__Page) -> Audit__Out rev 1
  emits []
  tests
    w(Audit__Page(Audit__Attrs("deep"), "t")) => Ok("deep")
  Ok(page.attrs.id)

fn audit__main() -> Audit__Out rev 1
  emits []
  tests
    m() => Ok("main")
  match call audit__use(Audit__Attrs("main"))
    on Ok r => Ok(r.value)

fn audit__box(flag: str) -> Audit__Box rev 1
  emits []
  tests
    x("a") => Ok(Audit__Attrs("a"))
    y("b") => Ok(Audit__Attrs("b"))
  Ok(Audit__Attrs(flag))

fn audit__eq(flag: str) -> Audit__Out rev 1
  emits []
  tests
    e("a") => Ok("same")
    n("b") => Ok("diff")
  match call audit__box(flag)
    on Ok left => match call audit__box("a")
      on Ok right => match left == right
        true => Ok("same")
        false => Ok("diff")
`

func TestRecordProductsEvaluate(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"audit.can": recordFixture})
	// Construction, passing, nested reads, boxed returns, and outcome
	// comparison agree in the evaluator.
	if diags := diagnose(dir, "audit.can", recordFixture); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestRecordProductsEmit(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "audit.can"), []byte(recordFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "audit.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(out, "audit.ts"))
	if err != nil {
		t.Fatal(err)
	}
	ts := string(got)
	for _, want := range []string{
		`export type Audit__Attrs = { id: string };`,
		`export type Audit__Page = { attrs: Audit__Attrs; title: string };`,
		`export function audit__use(attrs: Audit__Attrs)`,
		`{ id: "main" }`,
		`$canEqRec(left, right, ["attrs"])`,
		`{ $can_kind: "ok"; attrs: Audit__Attrs }`,
	} {
		if !strings.Contains(ts, want) {
			t.Errorf("emit missing %q\n--- emit ---\n%s", want, ts)
		}
	}
}

func TestRecordConstructionRefusals(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string
		code string
	}{
		{"wrong constructor", "Ok(value = Audit__Nope(id = \"x\"))", "unknown record Audit__Nope", "CAN6002"},
		{"missing field", "Audit__Attrs()", "is missing field id", "CAN6003"},
		{"duplicate field", "Audit__Attrs(id = \"a\", id = \"b\")", "repeats field id", "CAN6003"},
		{"wrong field type", "Audit__Attrs(id = 1)", "got int, want str", "CAN6003"},
		{"unknown field", "Audit__Attrs(id = \"a\", bogus = 1)", "has no field bogus", "CAN6003"},
		{"positional over arity", "Audit__Attrs(\"a\", \"b\")", "takes 2 args for 1 fields", "CAN6003"},
		{"positional double claim", "Audit__Attrs(\"a\", id = \"b\")", "supplies field id twice", "CAN6003"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			use := map[string]string{
				"wrong constructor":       "Audit__Nope(id = \"x\")",
				"missing field":           "Audit__Attrs()",
				"duplicate field":         "Audit__Attrs(id = \"a\", id = \"b\")",
				"wrong field type":        "Audit__Attrs(id = 1)",
				"unknown field":           "Audit__Attrs(id = \"a\", bogus = 1)",
				"positional over arity":   "Audit__Attrs(\"a\", \"b\")",
				"positional double claim": "Audit__Attrs(\"a\", id = \"b\")",
			}[tc.name]
			mod := `mod audit
  provides [audit__probe, Audit__Attrs]
  uses []
  emits []

type Audit__Attrs rev 1 (
  id: str
)

fn audit__probe() -> Audit__Attrs rev 1
  emits []
  tests
    g() => Audit__Attrs(id = "v")
  USE
`
			mod = strings.Replace(mod, "USE", use, 1)
			dir := writeLSPDir(t, map[string]string{"audit.can": mod})
			diags := diagnose(dir, "audit.can", mod)
			if !hasDiag(diags, "error", tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, diags)
			}
			if !hasCode(diags, tc.code) {
				t.Fatalf("expected %s, got %v", tc.code, diags)
			}
		})
	}
}

// TestPositionalConstruction pins a92: records, errors, and Ok
// accept positional args binding fields front-to-back. Clean
// diagnose proves check, execution, and comparison all agree
// through resolution.
func TestPositionalConstruction(t *testing.T) {
	mod := `mod audit
  provides [audit__probe, Audit__Pair]
  uses []
  emits [audit.oops]

error audit.oops(code: int, note: str)

type Audit__Pair rev 1 (
  first: str
  second: str
)

fn audit__probe(which: int) -> Audit__Pair rev 1
  emits [audit.oops]
  tests
    pair(0) => Ok("a", "b")
    oops(9) => audit.oops(3, "bad")
  match which
    0 => Ok("a", "b")
    _ => audit.oops(3, "bad")
`
	dir := writeLSPDir(t, map[string]string{"audit.can": mod})
	if diags := diagnose(dir, "audit.can", mod); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// TestPositionalMixedBindsByIndex pins the a92 mixing rule (the
// calls rule: an unnamed arg claims its own list index, so
// trailing positionals bind past a named head). The head name
// sits in its own slot, so the lint notes exactly that and
// nothing else; no static fault and no test failure prove the
// binding ran.
func TestPositionalMixedBindsByIndex(t *testing.T) {
	mod := `mod audit
  provides [audit__probe, Audit__Triple]
  uses []
  emits []

type Audit__Triple rev 1 (
  first: str
  second: str
  third: str
)

fn audit__probe(which: int) -> Audit__Triple rev 1
  emits []
  tests
    mixed(1) => Ok(first = "a", "b", "c")
  match which
    _ => Ok("a", "b", "c")
`
	dir := writeLSPDir(t, map[string]string{"audit.can": mod})
	diags := diagnose(dir, "audit.can", mod)
	if len(diags) != 1 || diags[0].Code != CodeLintRedundant {
		t.Fatalf("expected exactly the redundant-name lint, got %v", diags)
	}
	if !strings.Contains(diags[0].Msg, `"first"`) {
		t.Fatalf("lint must name the head arg, got %v", diags)
	}
}

// TestPositionalBindsByListIndex pins what the mixing rule is
// NOT: positional-ordinal. Under ordinal binding the second
// positional would take the second unclaimed field; under the
// calls rule it claims its own list index, colliding with the
// named head — a double-claim fault.
func TestPositionalBindsByListIndex(t *testing.T) {
	mod := `mod audit
  provides [audit__probe, Audit__Triple]
  uses []
  emits []

type Audit__Triple rev 1 (
  first: str
  second: str
  third: str
)

fn audit__probe(which: int) -> Audit__Triple rev 1
  emits []
  tests
    mixed(1) => Ok("a", "b", "c")
  match which
    _ => Ok("a", third = "c", "b")
`
	dir := writeLSPDir(t, map[string]string{"audit.can": mod})
	diags := diagnose(dir, "audit.can", mod)
	if !hasDiag(diags, "error", "supplies field third twice") {
		t.Fatalf("expected list-index double-claim, got %v", diags)
	}
}

// TestPositionalScalarOk pins a92 scalar Ok: exactly the
// conventional `value`. A lone positional binds it and runs;
// B09 also checks named fields statically: a wild name is a
// shape error, not just a failed example. Extra positionals fault.
func TestPositionalScalarOk(t *testing.T) {
	mod := `mod m
  provides [m__go]
  uses []
  emits []

fn m__go(x: int) -> int rev 1
  emits []
  tests
    pos(4) => Ok(8)
  match x
    _ => Ok(8)
`
	dir := writeLSPDir(t, map[string]string{"m.can": mod})
	if diags := diagnose(dir, "m.can", mod); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
	wild := `mod m
  provides [m__go]
  uses []
  emits []

fn m__go(x: int) -> int rev 1
  emits []
  tests
    wild(4) => Ok(answer = 8)
  match x
    _ => Ok(8)
`
	dir = writeLSPDir(t, map[string]string{"m.can": wild})
	diags := diagnose(dir, "m.can", wild)
	if !hasDiag(diags, "error", "Ok has no field answer") || !hasDiag(diags, "error", "Ok is missing field value") {
		t.Fatalf("wild scalar name must fail statically, got %v", diags)
	}
	faults := `mod m
  provides [m__go]
  uses []
  emits []

fn m__go(x: int) -> int rev 1
  emits []
  tests
    over(4) => Ok(8, 9)
    dbl(4) => Ok(9, value = 8)
  match x
    _ => Ok(8)
`
	dir = writeLSPDir(t, map[string]string{"m.can": faults})
	diags = diagnose(dir, "m.can", faults)
	if !hasDiag(diags, "error", "Ok takes 2 args for 1 field") {
		t.Fatalf("expected scalar over-arity fault, got %v", diags)
	}
	if !hasDiag(diags, "error", "Ok supplies field value twice") {
		t.Fatalf("expected scalar double-claim fault, got %v", diags)
	}
}

// TestPositionalOkUnderErrorArm pins the a92 error-arm rule: an
// Ok arm top is return-positioned even where threading goes
// want-free, so positional args bind there too. (std/text
// regressed here: Ok under not_found arms died at execution.)
func TestPositionalOkUnderErrorArm(t *testing.T) {
	mod := `mod m
  provides [m__go, M__Out, net__fetch]
  uses []
  emits [net.down]

error net.down()

type M__Out rev 1 (
  body: str
)

extern net__fetch() -> M__Out rev 1
  emits [net.down]

fn m__go() -> M__Out rev 1
  emits []
  tests
    ok() => Ok("hi")
    down() => Ok("down")
  match call net__fetch()
    given
      ok => [exchange args () outcome Ok("hi")]
      down => [exchange args () outcome net.down()]
    on net.down _ => Ok("down")
    on Ok d => Ok(d.body)
`
	dir := writeLSPDir(t, map[string]string{"m.can": mod})
	if diags := diagnose(dir, "m.can", mod); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestRecordBrandField(t *testing.T) {
	mod := `mod audit
  provides [audit__probe, Audit__Token, Audit__Other, Audit__Session]
  uses []
  emits []

brand Audit__Token is str rev 1
brand Audit__Other is str rev 1

type Audit__Session rev 1 (
  token: Audit__Token
)

fn audit__probe() -> Audit__Session rev 1
  emits []
  tests
    g() => Audit__Session(token = seal Audit__Token("s"))
  Audit__Session(token = seal Audit__Other("x"))
`
	dir := writeLSPDir(t, map[string]string{"audit.can": mod})
	diags := diagnose(dir, "audit.can", mod)
	if !hasDiag(diags, "error", "got Audit__Other, want Audit__Token") {
		t.Fatalf("expected brand mismatch, got %v", diags)
	}
}

func TestRecordCycleIsCAN6006(t *testing.T) {
	cases := []struct {
		name string
		decl string
		want string
	}{
		{"mutual", "type Audit__A rev 1 (\n  b: Audit__B\n)\n\ntype Audit__B rev 1 (\n  a: Audit__A\n)", "Audit__A -> Audit__B -> Audit__A"},
		{"self", "type Audit__A rev 1 (\n  me: Audit__A\n)", "Audit__A -> Audit__A"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mod := "mod audit\n  provides [Audit__A, Audit__B]\n  uses []\n  emits []\n\n" + tc.decl + "\n"
			dir := writeLSPDir(t, map[string]string{"audit.can": mod})
			diags := diagnose(dir, "audit.can", mod)
			if !hasDiag(diags, "error", "record type cycle") {
				t.Fatalf("expected cycle error, got %v", diags)
			}
			if !hasDiag(diags, "error", tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, diags)
			}
			if !hasCode(diags, "CAN6006") {
				t.Fatalf("expected CAN6006, got %v", diags)
			}
		})
	}
}
