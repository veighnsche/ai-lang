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
    u(attrs = Audit__Attrs(id = "main")) => Ok(value = "main")
=
  Ok(value = attrs.id)

fn audit__wrap(page: Audit__Page) -> Audit__Out rev 1
  emits []
  tests
    w(page = Audit__Page(attrs = Audit__Attrs(id = "deep"), title = "t")) => Ok(value = "deep")
=
  Ok(value = page.attrs.id)

fn audit__main() -> Audit__Out rev 1
  emits []
  tests
    m() => Ok(value = "main")
=
  match call audit__use(attrs = Audit__Attrs(id = "main"))
    on Ok r => Ok(value = r.value)

fn audit__box(flag: str) -> Audit__Box rev 1
  emits []
  tests
    x(flag = "a") => Ok(attrs = Audit__Attrs(id = "a"))
    y(flag = "b") => Ok(attrs = Audit__Attrs(id = "b"))
=
  Ok(attrs = Audit__Attrs(id = flag))

fn audit__eq(flag: str) -> Audit__Out rev 1
  emits []
  tests
    e(flag = "a") => Ok(value = "same")
    n(flag = "b") => Ok(value = "diff")
=
  match call audit__box(flag = flag)
    on Ok left => match call audit__box(flag = "a")
      on Ok right => match left == right
        true => Ok(value = "same")
        false => Ok(value = "diff")
`

func TestRecordProductsEvaluate(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"audit.ail": recordFixture})
	// Construction, passing, nested reads, boxed returns, and outcome
	// comparison agree in the evaluator.
	if diags := diagnose(dir, "audit.ail", recordFixture); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestRecordProductsEmit(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "audit.ail"), []byte(recordFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "audit.ail")}); err != nil {
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
		`$ailEqRec(left, right, ["attrs"])`,
		`{ $ail_kind: "ok"; attrs: Audit__Attrs }`,
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
		{"wrong constructor", "Ok(value = Audit__Nope(id = \"x\"))", "unknown record Audit__Nope", "AIL6002"},
		{"missing field", "Audit__Attrs()", "is missing field id", "AIL6003"},
		{"duplicate field", "Audit__Attrs(id = \"a\", id = \"b\")", "repeats field id", "AIL6003"},
		{"wrong field type", "Audit__Attrs(id = 1)", "got int, want str", "AIL6003"},
		{"unknown field", "Audit__Attrs(id = \"a\", bogus = 1)", "has no field bogus", "AIL6003"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			use := map[string]string{
				"wrong constructor": "Audit__Nope(id = \"x\")",
				"missing field":     "Audit__Attrs()",
				"duplicate field":   "Audit__Attrs(id = \"a\", id = \"b\")",
				"wrong field type":  "Audit__Attrs(id = 1)",
				"unknown field":     "Audit__Attrs(id = \"a\", bogus = 1)",
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
=
  USE
`
			mod = strings.Replace(mod, "USE", use, 1)
			dir := writeLSPDir(t, map[string]string{"audit.ail": mod})
			diags := diagnose(dir, "audit.ail", mod)
			if !hasDiag(diags, "error", tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, diags)
			}
			if !hasCode(diags, tc.code) {
				t.Fatalf("expected %s, got %v", tc.code, diags)
			}
		})
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
=
  Audit__Session(token = seal Audit__Other("x"))
`
	dir := writeLSPDir(t, map[string]string{"audit.ail": mod})
	diags := diagnose(dir, "audit.ail", mod)
	if !hasDiag(diags, "error", "got Audit__Other, want Audit__Token") {
		t.Fatalf("expected brand mismatch, got %v", diags)
	}
}

func TestRecordCycleIsAIL6006(t *testing.T) {
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
			dir := writeLSPDir(t, map[string]string{"audit.ail": mod})
			diags := diagnose(dir, "audit.ail", mod)
			if !hasDiag(diags, "error", "record type cycle") {
				t.Fatalf("expected cycle error, got %v", diags)
			}
			if !hasDiag(diags, "error", tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, diags)
			}
			if !hasCode(diags, "AIL6006") {
				t.Fatalf("expected AIL6006, got %v", diags)
			}
		})
	}
}
