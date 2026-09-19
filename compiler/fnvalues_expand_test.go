package main

import (
	"sort"
	"strings"
	"testing"
)

func typeNames(tds map[string]*TypeDecl) []string {
	var out []string
	for k := range tds {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// B00 expansion: fnref type arguments demand stamps like calls,
// Fn heads descend for instance collection, and every walker
// covers the new positions. Full-pipeline fixtures: assertions
// target stamps and slot strings, plus the expansion refusals
// (bare generics, monomorphic arguments, chain-else).

const fnvaluesExpandLib = `mod m
  provides [m__id, m__go, M__Box, M__O]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__Box<T> rev 1 (
  item: T
)

type M__O rev 1 (
  value: str
)

fn m__id<T>(x: T) -> M__Box<T> rev 1
  emits []
  tests
    i<T=int>(1) => Ok(1)
  Ok(x)

fn m__go(n: int) -> M__Box<int> rev 1
  emits []
  tests
    g(1) => Ok(1)
  Ok(fnref m__id<int>(x = n))
`

func TestExpandFnrefDemand(t *testing.T) {
	prog, mods, _ := genericProgram(t, map[string]string{"m.can": fnvaluesExpandLib})
	if prog == nil {
		t.Fatal("expected program")
	}
	if _, ok := prog.Fns["m__id$T$int"]; !ok {
		t.Fatalf("fnref demand missing; fns: %v", fnNames(prog))
	}
	// The reference site rewrites to the stamp like a call site.
	var fn *FnDecl
	for _, d := range mods[0].Decls {
		if f, ok := d.(*FnDecl); ok && f.Name == "m__go" {
			fn = f
		}
	}
	if fn == nil {
		t.Fatal("m__go missing")
	}
	ref := fn.Body.Small.Args[0].V
	if ref.Kind != "fnref" || ref.Fname != "m__id$T$int" || len(ref.TypeArgs) != 0 {
		t.Fatalf("ref = %+v, want rewritten fnref to m__id$T$int", ref)
	}
}

func TestExpandFnrefBareGeneric(t *testing.T) {
	src := `mod m
  provides [m__id, m__go, M__Box]
  uses []
  emits []

type M__Box<T> rev 1 (
  item: T
)

fn m__id<T>(x: T) -> M__Box<T> rev 1
  emits []
  tests
    i<T=int>(1) => Ok(1)
  Ok(x)

fn m__go(n: int) -> M__Box<int> rev 1
  emits []
  tests
    g(1) => Ok(1)
  Ok(fnref m__id(x = n))
`
	_, _, collected := genericProgram(t, map[string]string{"m.can": src})
	wantGenericDiag(t, collected, CodeGenericExpand, "needs explicit type arguments")
}

func TestExpandFnrefChainElseGate(t *testing.T) {
	src := `mod m
  provides [m__id, m__go, M__Box]
  uses []
  emits []

type M__Box<T> rev 1 (
  item: T
)

fn m__id<T>(x: T) -> M__Box<T> rev 1
  emits []
  tests
    i<T=int>(1) => Ok(1)
  Ok(x)

fn m__go(n: int) -> M__Box<int> rev 1
  emits []
  tests
    g(1) => Ok(1)
  match chain
    call m__id<int>(1) as a
    then a
    else Ok(fnref m__id<int>(x = n))
`
	_, _, collected := genericProgram(t, map[string]string{"m.can": src})
	wantGenericDiag(t, collected, CodeGenericExpand, "chain else")
}

// Two generic targets referenced in reverse declaration order
// both stamp: demand collection is order-independent, and the
// call site and the reference site each rewrite to their stamp.
func TestExpandTwoGenericsReverseOrder(t *testing.T) {
	src := `mod m
  provides [m__first, m__second, m__go, M__Box]
  uses []
  emits []

type M__Box<T> rev 1 (
  item: T
)

fn m__first<T>(x: T) -> M__Box<T> rev 1
  emits []
  tests
    f<T=int>(1) => Ok(1)
  Ok(x)

fn m__second<T>(x: T) -> M__Box<T> rev 1
  emits []
  tests
    s<T=int>(1) => Ok(1)
  Ok(x)

fn m__go(cb: Fn<int, M__Box<int>, []>, n: int) -> M__Box<int> rev 1
  emits []
  tests
    g(fnref m__second<int>(x = 1), 2) => Ok(2)
  match call m__first<int>(n)
    on Ok r => match invoke cb with r.item
      on Ok o => Ok(o.item)
`
	prog, mods, _ := genericProgram(t, map[string]string{"m.can": src})
	if prog == nil {
		t.Fatal("expected program")
	}
	for _, stamp := range []string{"m__first$T$int", "m__second$T$int"} {
		if _, ok := prog.Fns[stamp]; !ok {
			t.Fatalf("stamp %s missing; fns: %v", stamp, fnNames(prog))
		}
	}
	var fn *FnDecl
	for _, d := range mods[0].Decls {
		if f, ok := d.(*FnDecl); ok && f.Name == "m__go" {
			fn = f
		}
	}
	if fn == nil {
		t.Fatal("m__go missing")
	}
	call := fn.Body.Scruts[0]
	if call.Kind != "call" || call.Fname != "m__first$T$int" {
		t.Fatalf("call site = %+v, want rewritten call to m__first$T$int", call)
	}
	ref := fn.Tests[0].Args[0].V
	if ref.Kind != "fnref" || ref.Fname != "m__second$T$int" || len(ref.TypeArgs) != 0 {
		t.Fatalf("row ref = %+v, want rewritten fnref to m__second$T$int", ref)
	}
}

func TestExpandFnrefMonomorphicArgs(t *testing.T) {
	src := `mod m
  provides [m__t, m__go, M__O]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__O rev 1 (
  value: str
)

fn m__t(divisor: int, dividend: int) -> M__O rev 1
  emits [m.err]
  tests
    t(3, 7) => Ok("q")
  Ok("q")

fn m__go(n: int) -> M__O rev 1
  emits []
  tests
    g(1) => Ok("q")
  Ok(fnref m__t<int>(divisor = n))
`
	_, _, collected := genericProgram(t, map[string]string{"m.can": src})
	wantGenericDiag(t, collected, CodeGenericExpand, "type arguments on monomorphic function")
}

const fnvaluesExpandSlots = `mod m
  provides [m__hold, m__go, M__Box, M__Wrap, M__O]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__Box<T> rev 1 (
  item: T
)

type M__Wrap<T> rev 1 (
  cb: Fn<T, M__Box<T>, []>
)

type M__O rev 1 (
  value: str
)

fn m__hold(w: M__Wrap<int>) -> M__O rev 1
  emits []
  tests
    h(1) => Ok("q")
  Ok("q")

fn m__go(cb: Fn<int, M__Box<str>, []>) -> M__O rev 1
  emits []
  tests
    g(1) => Ok("q")
  Ok("q")
`

func TestExpandFnSlotDescent(t *testing.T) {
	prog, mods, collected := genericProgram(t, map[string]string{"m.can": fnvaluesExpandSlots})
	if prog == nil {
		t.Fatal("expected program")
	}
	for _, e := range genericErrs(collected) {
		if strings.Contains(e, "never used in its fields") {
			t.Fatalf("Fn field must count as parameter use, got %v", genericErrs(collected))
		}
	}
	tds := typeDecls(mods)
	if _, ok := tds["M__Box$T$str"]; !ok {
		t.Fatalf("closed Fn success demand missing; types: %v", typeNames(tds))
	}
	if _, ok := tds["M__Box$T$int"]; !ok {
		t.Fatalf("relative Fn success demand missing; types: %v", typeNames(tds))
	}
	wrap, ok := tds["M__Wrap$T$int"]
	if !ok {
		t.Fatalf("wrap stamp missing; types: %v", typeNames(tds))
	}
	if wrap.Fields[0][1] != "Fn<int,M__Box$T$int,[]>" {
		t.Fatalf("wrap field = %q, want substituted + rewritten Fn head", wrap.Fields[0][1])
	}
}

const fnvaluesExpandInvoke = `mod m
  provides [m__go, M__O]
  uses []
  emits [m.err]

error m.err(detail: str)

type M__O rev 1 (
  value: str
)

fn m__go<T>(cb: Fn<T, M__O, []>, n: T) -> M__O rev 1
  emits [m.err]
  tests
    i<T=int>(1, 2) => Ok("q")
    s<T=str>("a", "b") => Ok("q")
  match invoke cb with n
    on Ok _ => Ok("q")
    on m.err _ => Ok("e")
`

func TestExpandInvokeArgCloned(t *testing.T) {
	prog, _, _ := genericProgram(t, map[string]string{"m.can": fnvaluesExpandInvoke})
	if prog == nil {
		t.Fatal("expected program")
	}
	intFn, ok := prog.Fns["m__go$T$int"]
	if !ok {
		t.Fatalf("int stamp missing; fns: %v", fnNames(prog))
	}
	strFn, ok := prog.Fns["m__go$T$str"]
	if !ok {
		t.Fatalf("str stamp missing; fns: %v", fnNames(prog))
	}
	if intFn.Body.InvokeArg == nil || strFn.Body.InvokeArg == nil {
		t.Fatal("InvokeArg missing after stamping")
	}
	if intFn.Body.InvokeArg == strFn.Body.InvokeArg {
		t.Fatal("stamps alias the template InvokeArg")
	}
}

// Invoke arguments substitute like scrutinees: a generic
// consumer invoking with a constructed generic argument stamps
// closed, rewritten construction sites per instance.
func TestExpandInvokeArgSubstituted(t *testing.T) {
	src := `mod m
  provides [m__go, m__t, m__u, M__O, M__Step]
  uses []
  emits []

type M__O rev 1 (
  value: str
)

type M__Step<A, T> rev 1 (
  acc: A
  elem: T
)

fn m__t(step: M__Step<int, int>) -> M__O rev 1
  emits []
  tests
    t(M__Step<int, int>(1, 2)) => Ok("q")
  Ok("q")

fn m__u(step: M__Step<int, str>) -> M__O rev 1
  emits []
  tests
    u(M__Step<int, str>(1, "b")) => Ok("q")
  Ok("q")

fn m__go<T>(cb: Fn<M__Step<int, T>, M__O, []>, n: T) -> M__O rev 1
  emits []
  tests
    i<T=int>(fnref m__t(), 2) => Ok("q")
    s<T=str>(fnref m__u(), "b") => Ok("q")
  match invoke cb with M__Step<int, T>(1, n)
    on Ok _ => Ok("q")
`
	prog, _, _ := genericProgram(t, map[string]string{"m.can": src})
	if prog == nil {
		t.Fatal("expected program")
	}
	for stamp, want := range map[string]string{
		"m__go$T$int": "M__Step$T$int$T$int",
		"m__go$T$str": "M__Step$T$int$T$str",
	} {
		fn, ok := prog.Fns[stamp]
		if !ok {
			t.Fatalf("stamp %s missing; fns: %v", stamp, fnNames(prog))
		}
		arg := fn.Body.InvokeArg
		if arg == nil || arg.Kind != "ctor" || arg.Ctor != want || len(arg.TypeArgs) != 0 {
			t.Fatalf("%s InvokeArg = %+v, want rewritten %s", stamp, arg, want)
		}
	}
}
