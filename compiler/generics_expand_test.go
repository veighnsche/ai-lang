package main

import (
	"sort"
	"strings"
	"testing"
)

// G1 expansion: stamps, routing, headers, and every rejection.
// Fixtures run the full checkProgram pipeline, so passing means
// the stamps check, prove, and execute as ordinary code.
//
// Shape discipline (v0, not generics): `-> T` instances must
// agree on bareness — one body text cannot serve bare and
// record shapes — and callers cannot consume bare returns
// (the `on Ok` channel assumes records; the outcomes design
// owns that gap). Compare-shaped generics (generic params,
// concrete record return) compose fully; select-shaped ones
// prove per-instance rows.

const genericLib = `mod m
  provides [m__sel]
  uses []
  emits []

fn m__sel<T>(condition: bool, when_true: T, when_false: T) -> T rev 1
  emits []
  tests
    pick_str_true<T=str>(true, "a", "b") => Ok("a")
    pick_str_false<T=str>(false, "a", "b") => Ok("b")
    pick_int_true<T=int>(true, 1, 2) => Ok(1)
    pick_int_false<T=int>(false, 1, 2) => Ok(2)
  match condition
    true => Ok(when_true)
    false => Ok(when_false)
`

const genericCmpLib = `mod m
  provides [m__cmp, m__go, M__I]
  uses []
  emits []

type M__I rev 1 (
  value: int
)

fn m__cmp<T>(left: T, right: T) -> M__I rev 1
  emits []
  tests
    lt_str<T=str>("a", "b") => Ok(-1)
    eq_str<T=str>("a", "a") => Ok(0)
    gt_str<T=str>("b", "a") => Ok(1)
    lt_int<T=int>(1, 2) => Ok(-1)
    eq_int<T=int>(2, 2) => Ok(0)
    gt_int<T=int>(2, 1) => Ok(1)
  match left < right, left == right
    true, _ => Ok(-1)
    _, true => Ok(0)
    false, false => Ok(1)

fn m__go(a: str, b: str) -> M__I rev 1
  emits []
  tests
    one("a", "b") => Ok(-1)
  match call m__cmp<str>(a, b)
    on Ok v => Ok(v.value)
`

func genericProgram(t *testing.T, files map[string]string) (*Program, []*Module, []Diag) {
	t.Helper()
	var mods []*Module
	texts := map[string]string{}
	names := []string{}
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		m, err := parseModuleText(name, files[name])
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		mods = append(mods, m)
		texts[m.ID] = files[name]
	}
	prog, collected := checkProgram(mods, texts, nil, nil)
	return prog, mods, collected
}

func genericErrs(collected []Diag) []string {
	var out []string
	for _, d := range collected {
		if d.Sev == "error" {
			out = append(out, d.Code+" "+d.Msg)
		}
	}
	return out
}

func wantGenericDiag(t *testing.T, collected []Diag, code, sub string) {
	t.Helper()
	for _, e := range genericErrs(collected) {
		if strings.HasPrefix(e, code) && strings.Contains(e, sub) {
			return
		}
	}
	t.Fatalf("expected %s containing %q, got %v", code, sub, genericErrs(collected))
}

func fnNames(prog *Program) []string {
	var out []string
	for k := range prog.Fns {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func TestGenericStampsAndRouting(t *testing.T) {
	prog, mods, collected := genericProgram(t, map[string]string{"m.can": genericLib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("expansion must be clean: %v", errs)
	}
	if prog == nil {
		t.Fatal("expected program")
	}
	str, ok := prog.Fns["m__sel$T$str"]
	if !ok {
		t.Fatalf("str stamp missing; fns: %v", fnNames(prog))
	}
	num, ok := prog.Fns["m__sel$T$int"]
	if !ok {
		t.Fatalf("int stamp missing; fns: %v", fnNames(prog))
	}
	if _, ok := prog.Fns["m__sel"]; ok {
		t.Fatal("template must not survive expansion")
	}
	if len(str.Tests) != 2 || str.Tests[0].Name != "pick_str_true" {
		t.Fatalf("str stamp rows = %+v, want 2 pick_str rows", str.Tests)
	}
	if len(num.Tests) != 2 || num.Tests[0].Name != "pick_int_true" {
		t.Fatalf("int stamp rows = %+v, want 2 pick_int rows", num.Tests)
	}
	if len(str.Tests[0].TypeBinds) != 0 {
		t.Fatal("stamps must strip binds")
	}
	if str.Ret != "str" || num.Ret != "int" {
		t.Fatalf("rets = %s/%s, want str/int", str.Ret, num.Ret)
	}
	// Stamps own their trees: no shared subtrees across instances.
	if str.Body == num.Body || str.Body.Scruts[0] == num.Body.Scruts[0] {
		t.Fatal("stamps share AST subtrees")
	}
	// provides names stamps, sorted, in place.
	var provides []string
	for _, m := range mods {
		provides = append(provides, m.Hdr["provides"]...)
	}
	joined := strings.Join(provides, ",")
	if !strings.Contains(joined, "m__sel$T$int,m__sel$T$str") {
		t.Fatalf("provides = %v, want sorted stamps", provides)
	}
}

func TestGenericCallerConcrete(t *testing.T) {
	prog, _, collected := genericProgram(t, map[string]string{"m.can": genericCmpLib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("caller must be clean: %v", errs)
	}
	goFn := prog.Fns["m__go"]
	var callees []string
	for _, c := range walkCalls(goFn.Body) {
		callees = append(callees, c.Fname)
		if len(c.TypeArgs) != 0 {
			t.Fatal("rewritten calls must drop type args")
		}
	}
	if len(callees) != 1 || callees[0] != "m__cmp$T$str" {
		t.Fatalf("caller reaches %v, want [m__cmp$T$str]", callees)
	}
}

func TestGenericPassThrough(t *testing.T) {
	lib := `mod m
  provides [m__inner, m__outer, M__I]
  uses []
  emits []

type M__I rev 1 (
  value: int
)

fn m__inner<T>(first: T, second: T) -> M__I rev 1
  emits []
  tests
    eq_str<T=str>("a", "a") => Ok(1)
    ne_str<T=str>("a", "b") => Ok(0)
  match first == second
    true => Ok(1)
    false => Ok(0)

fn m__outer<T>(first: T, second: T) -> M__I rev 1
  emits []
  tests
    use_str<T=str>("a", "a") => Ok(1)
  match call m__inner<T>(first, second)
    on Ok v => Ok(v.value)
`
	prog, _, collected := genericProgram(t, map[string]string{"m.can": lib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("pass-through must be clean: %v", errs)
	}
	if _, ok := prog.Fns["m__outer$T$str"]; !ok {
		t.Fatalf("outer stamp missing; fns: %v", fnNames(prog))
	}
	inner, ok := prog.Fns["m__inner$T$str"]
	if !ok {
		t.Fatalf("inner stamp missing; fns: %v", fnNames(prog))
	}
	if len(inner.Tests) != 2 {
		t.Fatalf("inner stamp rows = %d, want 2", len(inner.Tests))
	}
}

func TestGenericRowlessInstance(t *testing.T) {
	// An instance nobody rows is missing evidence for one shape:
	// the diagnostic names base and args, not the mangled stamp.
	lib := `mod m
  provides [m__inner, m__outer, M__I]
  uses []
  emits []

type M__I rev 1 (
  value: int
)

fn m__inner<T>(first: T, second: T) -> M__I rev 1
  emits []
  tests
    eq_str<T=str>("a", "a") => Ok(1)
    ne_str<T=str>("a", "b") => Ok(0)
  match first == second
    true => Ok(1)
    false => Ok(0)

fn m__outer<T>(first: T, second: T) -> M__I rev 1
  emits []
  tests
    use_int<T=int>(1, 1) => Ok(1)
  match call m__inner<T>(first, second)
    on Ok v => Ok(v.value)
`
	_, _, collected := genericProgram(t, map[string]string{"m.can": lib})
	wantGenericDiag(t, collected, CodeMissingTests, "generic m__inner at <int> ships no tests")
}

func TestGenericCrossModule(t *testing.T) {
	lib := `mod lib
  provides [lib__cmp, Lib__I]
  uses []
  emits []

type Lib__I rev 1 (
  value: int
)

fn lib__cmp<T>(left: T, right: T) -> Lib__I rev 1
  emits []
  tests
    lt_str<T=str>("a", "b") => Ok(-1)
    eq_str<T=str>("a", "a") => Ok(0)
    gt_str<T=str>("b", "a") => Ok(1)
  match left < right, left == right
    true, _ => Ok(-1)
    _, true => Ok(0)
    false, false => Ok(1)
`
	app := `mod app
  provides [app__go, App__O]
  uses [lib__cmp@1]
  emits []

type App__O rev 1 (
  value: int
)

fn app__go(a: str, b: str) -> App__O rev 1
  emits []
  tests
    one("a", "b") => Ok(-1)
  match call lib__cmp<str>(a, b)
    given
      one => [exchange args (left = "a", right = "b") outcome Ok(-1)]
    on Ok v => Ok(v.value)
`
	prog, mods, collected := genericProgram(t, map[string]string{"app.can": app, "lib.can": lib})
	if errs := genericErrs(collected); len(errs) != 0 {
		t.Fatalf("cross-module must be clean: %v", errs)
	}
	if _, ok := prog.Fns["lib__cmp$T$str"]; !ok {
		t.Fatalf("foreign stamp missing; fns: %v", fnNames(prog))
	}
	for _, m := range mods {
		if m.Mod == "app" {
			joined := strings.Join(m.Hdr["uses"], ",")
			if joined != "lib__cmp$T$str@1" {
				t.Fatalf("app uses = %v, want [lib__cmp$T$str@1]", m.Hdr["uses"])
			}
		}
		if b := prog.GenericBase["lib__cmp$T$str"]; b != "lib__cmp" {
			t.Fatalf("GenericBase = %q, want lib__cmp", b)
		}
	}
}

func TestGenericMissingPinNamesBase(t *testing.T) {
	lib := `mod lib
  provides [lib__cmp, Lib__I]
  uses []
  emits []

type Lib__I rev 1 (
  value: int
)

fn lib__cmp<T>(left: T, right: T) -> Lib__I rev 1
  emits []
  tests
    lt_str<T=str>("a", "b") => Ok(-1)
    eq_str<T=str>("a", "a") => Ok(0)
    gt_str<T=str>("b", "a") => Ok(1)
  match left < right, left == right
    true, _ => Ok(-1)
    _, true => Ok(0)
    false, false => Ok(1)
`
	app := `mod app
  provides [app__go, App__O]
  uses []
  emits []

type App__O rev 1 (
  value: int
)

fn app__go(a: str, b: str) -> App__O rev 1
  emits []
  tests
    one("a", "b") => Ok(-1)
  match call lib__cmp<str>(a, b)
    given
      one => [exchange args (left = "a", right = "b") outcome Ok(-1)]
    on Ok v => Ok(v.value)
`
	_, _, collected := genericProgram(t, map[string]string{"app.can": app, "lib.can": lib})
	wantGenericDiag(t, collected, CodeCallNotInUses, "calls lib__cmp which is not in uses")
}

func TestGenericDeadPinWarns(t *testing.T) {
	lib := `mod lib
  provides [lib__id]
  uses []
  emits []

fn lib__id<T>(value: T) -> T rev 1
  emits []
  tests
    id_str<T=str>("a") => Ok("a")
  match value
    _ => Ok(value)
`
	app := `mod app
  provides [app__go, App__O]
  uses [lib__id@1]
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

func TestGenericRejects(t *testing.T) {
	header := `mod m
  provides [m__f, m__g, M__O, M__V]
  uses []
  emits []

type M__O rev 1 (
  value: str
)

variant M__V rev 1 (
  case A()
  case B(name: str)
)

`
	mkfn := func(name, sig, rows, body string) string {
		return "fn " + name + sig + " rev 1\n  emits []\n  tests\n" + rows + "  " + body + "\n"
	}
	row := func(text string) string { return "    " + text + "\n" }
	valbody := "match value\n    _ => Ok(value)\n"
	cases := []struct {
		name string
		src  string
		want string
	}{
		{"missing args",
			mkfn("m__g<T>", "(value: T) -> T", row(`r<T=str>("a") => Ok("a")`), valbody) +
				mkfn("m__f", "(value: str) -> str", row(`r("a") => Ok("a")`), "match call m__g(value)\n    on Ok v => Ok(v)\n"),
			"needs explicit type arguments"},
		{"call arity",
			mkfn("m__g<T>", "(value: T) -> T", row(`r<T=str>("a") => Ok("a")`), valbody) +
				mkfn("m__f", "(value: str) -> str", row(`r("a") => Ok("a")`), "match call m__g<str,int>(value)\n    on Ok v => Ok(v)\n"),
			"takes 1 type arguments, got 2"},
		{"incomplete row",
			mkfn("m__f<T, U>", "(first: T, second: U) -> T", row(`r<T=str>("a", 1) => Ok("a")`), "match first\n    _ => Ok(first)\n"),
			"must pin every type parameter"},
		{"unknown param row",
			mkfn("m__f<T>", "(value: T) -> T", row(`r<U=str>("a") => Ok("a")`), valbody),
			"unknown type parameter U"},
		{"dup bind",
			mkfn("m__f<T>", "(value: T) -> T", row(`r<T=str,T=int>("a") => Ok("a")`), valbody),
			"pins T twice"},
		{"binds on mono row",
			mkfn("m__f", "(value: str) -> str", row(`r<T=str>("a") => Ok("a")`), valbody),
			"takes no type binds"},
		{"args on mono call",
			mkfn("m__g", "(value: str) -> str", row(`r("a") => Ok("a")`), valbody) +
				mkfn("m__f", "(value: str) -> str", row(`r("a") => Ok("a")`), "match call m__g<str>(value)\n    on Ok v => Ok(v)\n"),
			"type arguments on monomorphic function"},
		{"unknown arg",
			mkfn("m__f<T>", "(value: T) -> T", row(`r<T=Nope>("a") => Ok("a")`), valbody),
			"unknown type Nope"},
		{"nested arg",
			mkfn("m__f<T>", "(value: T) -> T", row(`r<T=Box<str>>("a") => Ok("a")`), valbody),
			"nested instantiation"},
		{"mixed seq param",
			mkfn("m__g<T>", "(value: T) -> T", row(`r<T=str>("a") => Ok("a")`), valbody) +
				mkfn("m__f<U>", "(value: Seq<U>) -> Seq<U>", row(`r<U=str>(Seq<str>["a"]) => Ok(Seq<str>["a"])`), "match call m__g<Seq<U>>(value)\n    on Ok v => Ok(v)\n"),
			"bare parameters or closed types"},
		{"seq variant arg",
			mkfn("m__f<T>", "(value: T) -> T", row(`r<T=Seq<M__V>>("a") => Ok("a")`), valbody),
			"variant sequences are not admitted"},
		{"unused param",
			mkfn("m__f<T>", "(value: str) -> str", row(`r<T=str>("a") => Ok("a")`), valbody),
			"never used in its signature"},
		{"shadow param",
			mkfn("m__f<T>", "(T: str) -> T", row(`r<T=str>("a") => Ok("a")`), "match T\n    _ => Ok(T)\n"),
			"shadows a type parameter"},
		{"dup base",
			mkfn("m__f<T>", "(value: T) -> T", row(`r<T=str>("a") => Ok("a")`), valbody) +
				mkfn("m__f", "(value: str) -> str", row(`r("a") => Ok("a")`), valbody),
			"already declared"},
		{"self different args",
			mkfn("m__f<T>", "(value: T) -> T", row(`r<T=str>("a") => Ok("a")`), "match call m__f<int>(value)\n    on Ok v => Ok(v)\n"),
			"must repeat its own type arguments"},
	}
	for _, c := range cases {
		_, _, collected := genericProgram(t, map[string]string{"m.can": header + c.src})
		wantGenericDiag(t, collected, CodeGenericExpand, c.want)
	}
}

func TestGenericRejectsSpecial(t *testing.T) {
	// Uninstantiated: no rows and no calls. Expansion reports it
	// before the universal rows-required check runs.
	src := `mod m
  provides [m__f]
  uses []
  emits []

fn m__f<T>(value: T) -> T rev 1
  emits []
  tests
  match value
    _ => Ok(value)
`
	_, _, collected := genericProgram(t, map[string]string{"m.can": src})
	wantGenericDiag(t, collected, CodeGenericExpand, "never instantiated")
	// Emits naming a param: error-set parameters are not admitted.
	src2 := `mod m
  provides [m__f, M__O]
  uses []
  emits []

type M__O rev 1 (
  value: str
)

fn m__f<T>(value: T) -> T rev 1
  emits [T]
  tests
    r<T=str>("a") => Ok("a")
  match value
    _ => Ok(value)
`
	_, _, collected2 := genericProgram(t, map[string]string{"m.can": src2})
	wantGenericDiag(t, collected2, CodeGenericExpand, "error-set parameters are not admitted")
}

func TestGenericChainElseGate(t *testing.T) {
	fn := &FnDecl{
		Name: "m__f", Rev: 1, TypeParams: []string{"T"},
		Params: [][2]string{{"value", "T"}}, Ret: "T", Line: 1,
		Body:  &Node{IsMatch: true, Kind: MatchChain, Line: 2, ChainElse: "call m__g<T>(value)", ChainElseLine: 3},
		Tests: []Test{{Name: "r", TypeBinds: [][2]string{{"T", "str"}}, Line: 4}},
	}
	m := &Module{File: "m.can", ID: "m.can", Mod: "m", Hdr: map[string][]string{"provides": {"m__f"}}, Decls: []Decl{fn}}
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
