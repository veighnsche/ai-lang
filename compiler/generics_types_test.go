package main

import (
	"strings"
	"testing"
)

// G2 expansion: generic records stamp one monomorphic copy per
// distinct mention; templates never reach checking, and every
// downstream phase consumes plain records. Fixtures run the full
// checkProgram pipeline, so passing means the stamps check,
// prove, and execute as ordinary code.

func typeDecls(mods []*Module) map[string]*TypeDecl {
	out := map[string]*TypeDecl{}
	for _, m := range mods {
		for _, d := range m.Decls {
			if td, ok := d.(*TypeDecl); ok {
				out[td.Name] = td
			}
		}
	}
	return out
}

func wantFields(t *testing.T, td *TypeDecl, want [][2]string) {
	t.Helper()
	if len(td.Fields) != len(want) {
		t.Fatalf("%s fields = %v, want %v", td.Name, td.Fields, want)
	}
	for i := range want {
		if td.Fields[i] != want[i] {
			t.Fatalf("%s fields = %v, want %v", td.Name, td.Fields, want)
		}
	}
}

const genericTypeLib = `mod m
  provides [m__wrap, m__hold, M__Box]
  uses []
  emits []

type M__Box<T> rev 1 (
  value: T
)

fn m__wrap(v: str) -> M__Box<str> rev 1
  emits []
  tests
    one("a") => Ok("a")
  Ok(v)

fn m__hold(b: M__Box<int>, n: int) -> M__Box<int> rev 1
  emits []
  tests
    one(M__Box<int>(value = 1), 2) => Ok(1)
  match b.value
    _ => Ok(b.value)
`

func TestGenericTypeStampsAndRouting(t *testing.T) {
	prog, mods, collected := genericProgram(t, map[string]string{"m.can": genericTypeLib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("expansion must be clean: %v", errs)
	}
	if prog == nil {
		t.Fatal("expected program")
	}
	tds := typeDecls(mods)
	str, ok := tds["M__Box$T$str"]
	if !ok {
		t.Fatalf("str stamp missing; types: %v", tds)
	}
	num, ok := tds["M__Box$T$int"]
	if !ok {
		t.Fatalf("int stamp missing; types: %v", tds)
	}
	if _, ok := tds["M__Box"]; ok {
		t.Fatal("template must not survive expansion")
	}
	wantFields(t, str, [][2]string{{"value", "str"}})
	wantFields(t, num, [][2]string{{"value", "int"}})
	if str.Rev != 1 || num.Rev != 1 {
		t.Fatal("stamps must keep the template rev")
	}
	// Annotations name stamps after rewriting.
	if prog.Fns["m__wrap"].Ret != "M__Box$T$str" {
		t.Fatalf("wrap ret = %s", prog.Fns["m__wrap"].Ret)
	}
	if p := prog.Fns["m__hold"].Params[0]; p[1] != "M__Box$T$int" {
		t.Fatalf("hold param = %v", p)
	}
	// provides names stamps, sorted, in place.
	var provides []string
	for _, m := range mods {
		provides = append(provides, m.Hdr["provides"]...)
	}
	joined := strings.Join(provides, ",")
	if !strings.Contains(joined, "M__Box$T$int,M__Box$T$str") {
		t.Fatalf("provides = %v, want sorted stamps", provides)
	}
}

func TestGenericTypeConstruction(t *testing.T) {
	lib := `mod m
  provides [m__get, M__Box, M__O]
  uses []
  emits []

type M__Box<T> rev 1 (
  value: T
)

type M__O rev 1 (
  value: str
)

fn m__get(b: M__Box<str>) -> M__O rev 1
  emits []
  tests
    one(M__Box<str>(value = "a")) => Ok("a")
  match b.value
    _ => Ok(b.value)
`
	prog, _, collected := genericProgram(t, map[string]string{"m.can": lib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("construction must be clean: %v", errs)
	}
	get := prog.Fns["m__get"]
	if len(get.Tests) != 1 {
		t.Fatalf("tests = %+v, want 1 row", get.Tests)
	}
	arg := get.Tests[0].Args[0].V
	if arg.Kind != "ctor" || arg.Ctor != "M__Box$T$str" {
		t.Fatalf("row arg = %+v, want M__Box$T$str construction", arg)
	}
	if len(arg.TypeArgs) != 0 {
		t.Fatal("rewritten constructions must drop type args")
	}
}

func TestGenericTypeNestedConstruction(t *testing.T) {
	// Explicit constructions also appear nested: a wrapper field
	// seeds the instance, and bodies and expectations construct it.
	lib := `mod m
  provides [m__nest, M__Box, M__Wrap]
  uses []
  emits []

type M__Box<T> rev 1 (
  value: T
)

type M__Wrap rev 1 (
  inner: M__Box<str>
)

fn m__nest(v: str) -> M__Wrap rev 1
  emits []
  tests
    one("a") => Ok(M__Box<str>(value = "a"))
  Ok(M__Box<str>(value = v))
`
	_, _, collected := genericProgram(t, map[string]string{"m.can": lib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("nested construction must be clean: %v", errs)
	}
}

func TestGenericTypePassThrough(t *testing.T) {
	// One closed mention drives the fixpoint: stamping A<str>
	// reveals B<str>, which stamps in the same phase.
	lib := `mod m
  provides [m__go, M__A, M__B, M__O]
  uses []
  emits []

type M__B<T> rev 1 (
  value: T
)

type M__A<T> rev 1 (
  inner: M__B<T>
)

type M__O rev 1 (
  value: str
)

fn m__go(a: M__A<str>) -> M__O rev 1
  emits []
  tests
    one(M__A<str>(inner = M__B<str>(value = "a"))) => Ok("a")
  match a.inner.value
    _ => Ok(a.inner.value)
`
	_, mods, collected := genericProgram(t, map[string]string{"m.can": lib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("pass-through must be clean: %v", errs)
	}
	tds := typeDecls(mods)
	a, ok := tds["M__A$T$str"]
	if !ok {
		t.Fatalf("A stamp missing; types: %v", tds)
	}
	wantFields(t, a, [][2]string{{"inner", "M__B$T$str"}})
	b, ok := tds["M__B$T$str"]
	if !ok {
		t.Fatalf("B stamp missing; types: %v", tds)
	}
	wantFields(t, b, [][2]string{{"value", "str"}})
}

func TestGenericTypeSeqSelf(t *testing.T) {
	// Seq-wrapped self-reference stays legal, exactly like
	// monomorphic records: values stay finite through the seq.
	lib := `mod m
  provides [m__go, M__Tree, M__O]
  uses []
  emits []

type M__Tree<T> rev 1 (
  value: T
  kids: Seq<M__Tree<T>>
)

type M__O rev 1 (
  value: str
)

fn m__go(tr: M__Tree<str>) -> M__O rev 1
  emits []
  tests
    one(M__Tree<str>(value = "a", kids = Seq<M__Tree<str>>[])) => Ok("a")
  match tr.value
    _ => Ok(tr.value)
`
	_, mods, collected := genericProgram(t, map[string]string{"m.can": lib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("Seq self-reference must be clean: %v", errs)
	}
	tds := typeDecls(mods)
	tr, ok := tds["M__Tree$T$str"]
	if !ok {
		t.Fatalf("Tree stamp missing; types: %v", tds)
	}
	wantFields(t, tr, [][2]string{{"value", "str"}, {"kids", "Seq<M__Tree$T$str>"}})
}

func TestGenericTypeCrossModule(t *testing.T) {
	lib := `mod lib
  provides [lib__wrap, Lib__Box, Lib__Bag]
  uses []
  emits []

type Lib__Box<T> rev 1 (
  value: T
)

type Lib__Bag<T> rev 1 (
  item: T
)

fn lib__wrap(v: str) -> Lib__Box<str> rev 1
  emits []
  tests
    one("a") => Ok("a")
  Ok(v)
`
	app := `mod app
  provides [app__go, App__O]
  uses [Lib__Box@1]
  emits []

type App__O rev 1 (
  value: str
)

fn app__go(b: Lib__Box<str>, g: Lib__Bag<int>) -> App__O rev 1
  emits []
  tests
    one(Lib__Box<str>(value = "a"), Lib__Bag<int>(item = 1)) => Ok("a")
  match b.value
    _ => Ok(b.value)
`
	prog, mods, collected := genericProgram(t, map[string]string{"app.can": app, "lib.can": lib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("cross-module must be clean: %v", errs)
	}
	if prog == nil {
		t.Fatal("expected program")
	}
	tds := typeDecls(mods)
	if _, ok := tds["Lib__Box$T$str"]; !ok {
		t.Fatalf("foreign stamp missing; types: %v", tds)
	}
	if _, ok := tds["Lib__Bag$T$int"]; !ok {
		t.Fatalf("unpinned stamp missing; types: %v", tds)
	}
	for _, m := range mods {
		if m.Mod == "app" {
			joined := strings.Join(m.Hdr["uses"], ",")
			if joined != "Lib__Box$T$str@1" {
				t.Fatalf("app uses = %v, want [Lib__Box$T$str@1]", m.Hdr["uses"])
			}
		}
		if b := prog.GenericBase["Lib__Box$T$str"]; b != "Lib__Box" {
			t.Fatalf("GenericBase = %q, want Lib__Box", b)
		}
	}
}

func TestGenericTypeDeadPinWarns(t *testing.T) {
	lib := `mod lib
  provides [lib__wrap, Lib__Box]
  uses []
  emits []

type Lib__Box<T> rev 1 (
  value: T
)

fn lib__wrap(v: str) -> Lib__Box<str> rev 1
  emits []
  tests
    one("a") => Ok("a")
  Ok(v)
`
	app := `mod app
  provides [app__go, App__O]
  uses [Lib__Box@1]
  emits []

type App__O rev 1 (
  value: str
)

fn app__go(x: str) -> App__O rev 1
  emits []
  tests
    one("a") => Ok("a")
  Ok(x)
`
	_, _, collected := genericProgram(t, map[string]string{"app.can": app, "lib.can": lib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("dead pin must warn, not error: %v", errs)
	}
	found := false
	for _, d := range collected {
		if d.Code == CodeUnusedUses && strings.Contains(d.Msg, "never instantiates it") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected dead-pin warning, got %v", collected)
	}
}

func TestGenericTypeFnInteraction(t *testing.T) {
	// A generic fn mentioning a generic type: the fn stamps
	// first, substituting Box<T> to Box<str>, and the type
	// phase stamps what the fn instances mention.
	lib := `mod m
  provides [m__wrap, M__Box]
  uses []
  emits []

type M__Box<T> rev 1 (
  value: T
)

fn m__wrap<T>(v: T) -> M__Box<T> rev 1
  emits []
  tests
    one<T=str>("a") => Ok("a")
  Ok(v)
`
	prog, mods, collected := genericProgram(t, map[string]string{"m.can": lib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("fn+type interaction must be clean: %v", errs)
	}
	wrap, ok := prog.Fns["m__wrap$T$str"]
	if !ok {
		t.Fatalf("fn stamp missing; fns: %v", fnNames(prog))
	}
	if wrap.Ret != "M__Box$T$str" {
		t.Fatalf("fn stamp ret = %s", wrap.Ret)
	}
	if _, ok := typeDecls(mods)["M__Box$T$str"]; !ok {
		t.Fatal("type stamp missing")
	}
}

func TestGenericTwoTypeTemplatesOneModule(t *testing.T) {
	lib := `mod m
  provides [m__go, M__Box, M__Bag, M__O]
  uses []
  emits []

type M__Box<T> rev 1 (
  value: T
)

type M__Bag<T> rev 1 (
  item: T
)

type M__O rev 1 (
  value: str
)

fn m__go(b: M__Box<str>, g: M__Bag<int>) -> M__O rev 1
  emits []
  tests
    one(M__Box<str>(value = "a"), M__Bag<int>(item = 1)) => Ok("a")
  match b.value
    _ => Ok(b.value)
`
	_, mods, collected := genericProgram(t, map[string]string{"m.can": lib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("two templates must be clean: %v", errs)
	}
	tds := typeDecls(mods)
	for _, stamp := range []string{"M__Box$T$str", "M__Bag$T$int"} {
		if _, ok := tds[stamp]; !ok {
			t.Fatalf("stamp %s missing; types: %v", stamp, tds)
		}
	}
	for _, base := range []string{"M__Box", "M__Bag"} {
		if _, ok := tds[base]; ok {
			t.Fatalf("template %s survived expansion", base)
		}
	}
}

func TestGenericTypeRejects(t *testing.T) {
	header := `mod m
  provides [m__go, M__Box, M__O]
  uses []
  emits []

type M__O rev 1 (
  value: str
)

`
	box := `type M__Box<T> rev 1 (
  value: T
)

`
	goRet := "fn m__go(v: str) -> M__O rev 1\n  emits []\n  tests\n    one(\"a\") => Ok(\"a\")\n  Ok(v)\n"
	cases := []struct {
		name string
		src  string
		code string
		want string
	}{
		{"unmentioned",
			box + goRet,
			CodeGenericExpand, "never instantiated"},
		{"unused param",
			"type M__Box<T> rev 1 (\n  value: str\n)\n\n" + goRet,
			CodeGenericExpand, "never used in its fields"},
		{"dup base",
			box + "type M__Box rev 1 (\n  value: str\n)\n\n" + goRet,
			CodeGenericExpand, "already declared"},
		{"direct self",
			"type M__Loop<T> rev 1 (\n  next: M__Loop<T>\n)\n\n" + goRet,
			CodeGenericExpand, "fields its own base"},
		{"bare ctor",
			box + "fn m__go(b: M__Box<str>) -> M__O rev 1\n  emits []\n  tests\n    one(M__Box(\"a\")) => Ok(\"a\")\n  match b.value\n    _ => Ok(b.value)\n",
			CodeGenericExpand, "needs explicit type arguments"},
		{"args on mono ctor",
			box + "fn m__go(v: str) -> M__O rev 1\n  emits []\n  tests\n    one(\"a\") => Ok(\"a\")\n  Ok(M__O<str>(value = v))\n",
			CodeGenericExpand, "type arguments on monomorphic type M__O"},
		{"arity annotation",
			box + "fn m__go(v: str) -> M__Box<str,int> rev 1\n  emits []\n  tests\n    one(\"a\") => Ok(\"a\")\n  Ok(v)\n",
			CodeGenericExpand, "takes 1 type arguments, got 2"},
		{"arity ctor",
			box + "fn m__go(b: M__Box<str>) -> M__O rev 1\n  emits []\n  tests\n    one(M__Box<str,int>(value = \"a\")) => Ok(\"a\")\n  match b.value\n    _ => Ok(b.value)\n",
			CodeGenericExpand, "takes 1 type arguments, got 2"},
		{"unknown base annotation",
			box + "fn m__go(v: str) -> M__Nope<str> rev 1\n  emits []\n  tests\n    one(\"a\") => Ok(\"a\")\n  Ok(v)\n",
			CodeGenericExpand, "unknown type M__Nope"},
		{"unknown base ctor",
			box + "fn m__go(b: M__Box<str>) -> M__O rev 1\n  emits []\n  tests\n    one(M__Nope<str>(value = \"a\")) => Ok(\"a\")\n  match b.value\n    _ => Ok(b.value)\n",
			CodeGenericExpand, "unknown type M__Nope"},
		{"nested arg",
			box + "fn m__go(v: str) -> M__Box<M__O<str>> rev 1\n  emits []\n  tests\n    one(\"a\") => Ok(\"a\")\n  Ok(v)\n",
			CodeGenericExpand, "nested instantiation"},
		{"args on Ok",
			box + "fn m__go(v: str) -> M__O rev 1\n  emits []\n  tests\n    one(\"a\") => Ok(\"a\")\n  Ok<str>(v)\n",
			CodeGenericExpand, "unknown type Ok"},
		{"bare base annotation",
			box + "fn m__go(x: M__Box, y: M__Box<str>) -> M__O rev 1\n  emits []\n  tests\n    one(\"a\", \"b\") => Ok(\"a\")\n  match y.value\n    _ => Ok(y.value)\n",
			CodeUnknownType, "unknown type"},
		{"param outside template",
			box + "fn m__go(v: str) -> M__Box<Seq<T>> rev 1\n  emits []\n  tests\n    one(\"a\") => Ok(\"a\")\n  Ok(v)\n",
			CodeGenericExpand, "unknown type T"},
		{"seq wrapped unknown",
			box + "fn m__go(v: str) -> Seq<M__Nope<str>> rev 1\n  emits []\n  tests\n    one(\"a\") => Ok(\"a\")\n  Ok(v)\n",
			CodeGenericExpand, "unknown type M__Nope"},
		{"state mention",
			box + "state holding: M__Box<str> = 0\n\n" + goRet,
			CodeGenericExpand, "cells hold str, int, bool, or dec"},
		{"seqlit unknown elem",
			box + "fn m__go(v: Seq<str>) -> M__O rev 1\n  emits []\n  tests\n    one(Seq<M__Nope<str>>[]) => Ok(\"a\")\n  match #v\n    _ => Ok(\"a\")\n",
			CodeGenericExpand, "unknown type M__Nope"},
	}
	for _, c := range cases {
		_, _, collected := genericProgram(t, map[string]string{"m.can": header + c.src})
		wantGenericDiag(t, collected, c.code, c.want)
	}
}

func TestGenericTypeChainElseGate(t *testing.T) {
	// Monomorphic fns are gated too: late-parsed else text never
	// rewrites, so a generic construction there is rejected up
	// front instead of degrading downstream.
	fn := &FnDecl{
		Name: "m__f", Rev: 1,
		Params: [][2]string{{"value", "str"}}, Ret: "M__Box<str>", Line: 1,
		Body:  &Node{IsMatch: true, Kind: MatchChain, Line: 2, ChainElse: "M__Box<str>(value)", ChainElseLine: 3},
		Tests: []Test{{Name: "r", Line: 4}},
	}
	box := &TypeDecl{Name: "M__Box", Rev: 1, TypeParams: []string{"T"}, Fields: [][2]string{{"value", "T"}}, Line: 5}
	m := &Module{File: "m.can", ID: "m.can", Mod: "m", Hdr: map[string][]string{"provides": {"m__f"}}, Decls: []Decl{fn, box}}
	diags := expandGenerics([]*Module{m}, map[string]string{"m.can": "mod m"})
	found := false
	for _, d := range diags {
		if d.Code == CodeGenericExpand && strings.Contains(d.Msg, "chain else") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected chain-else gate, got %v", diags)
	}
}
