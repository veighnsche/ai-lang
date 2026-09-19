package main

import (
	"strings"
	"testing"
)

// B00 stage 1a parse: Fn<A,R,[errs]> survives field splitting as
// one annotation, and fnref creates a distinct node kind. Static
// validation (known types, capture binding, invocation) arrives
// in later slices; parsing owns shape only.
const fnvaluesParseGood = `mod m
  provides [m__go, m__t, M__O]
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

fn m__go(cb: Fn<int, M__O, [m.err]>, n: int) -> M__O rev 1
  emits []
  tests
    g(3) => Ok("q")
  match n < 10
    true => Ok("q")
    false => Ok("q")
`

func TestParseFnFieldType(t *testing.T) {
	m, err := parseModuleText("m.can", fnvaluesParseGood)
	if err != nil {
		t.Fatalf("Fn annotation must parse: %v", err)
	}
	var fn *FnDecl
	for _, d := range m.Decls {
		if f, ok := d.(*FnDecl); ok && f.Name == "m__go" {
			fn = f
		}
	}
	if fn == nil {
		t.Fatal("m__go missing after parse")
	}
	if len(fn.Params) != 2 {
		t.Fatalf("params = %v, want 2", fn.Params)
	}
	if fn.Params[0][0] != "cb" || fn.Params[0][1] != "Fn<int, M__O, [m.err]>" {
		t.Fatalf("cb param = %v, want Fn<int, M__O, [m.err]>", fn.Params[0])
	}
	if fn.Params[1][0] != "n" || fn.Params[1][1] != "int" {
		t.Fatalf("n param = %v, want int", fn.Params[1])
	}
}

func TestParseFnEmptyErrors(t *testing.T) {
	src := strings.Replace(fnvaluesParseGood,
		"cb: Fn<int, M__O, [m.err]>", "cb: Fn<int, M__O, []>", 1)
	m, err := parseModuleText("m.can", src)
	if err != nil {
		t.Fatalf("empty error list must parse: %v", err)
	}
	for _, d := range m.Decls {
		if f, ok := d.(*FnDecl); ok && f.Name == "m__go" {
			if f.Params[0][1] != "Fn<int, M__O, []>" {
				t.Fatalf("cb param = %v", f.Params[0])
			}
			return
		}
	}
	t.Fatal("m__go missing after parse")
}

func TestParseFnref(t *testing.T) {
	src := strings.Replace(fnvaluesParseGood,
		"  match n < 10\n    true => Ok(\"q\")\n    false => Ok(\"q\")\n",
		"  Ok(fnref m__t(divisor = 3))\n", 1)
	m, err := parseModuleText("m.can", src)
	if err != nil {
		t.Fatalf("fnref must parse: %v", err)
	}
	var body *Small
	for _, d := range m.Decls {
		if f, ok := d.(*FnDecl); ok && f.Name == "m__go" {
			body = f.Body.Small
		}
	}
	if body == nil || body.Kind != "ctor" {
		t.Fatalf("body = %+v, want Ok ctor", body)
	}
	if len(body.Args) != 1 {
		t.Fatalf("Ok args = %+v, want 1", body.Args)
	}
	ref := body.Args[0].V
	if ref.Kind != "fnref" {
		t.Fatalf("kind = %q, want fnref", ref.Kind)
	}
	if ref.Fname != "m__t" {
		t.Fatalf("target = %q, want m__t", ref.Fname)
	}
	if len(ref.Args) != 1 || !ref.Args[0].HasName || ref.Args[0].Name != "divisor" {
		t.Fatalf("captures = %+v, want named divisor", ref.Args)
	}
}

func TestParseFnrefGeneric(t *testing.T) {
	src := strings.Replace(fnvaluesParseGood,
		"  match n < 10\n    true => Ok(\"q\")\n    false => Ok(\"q\")\n",
		"  Ok(fnref m__t<int>(divisor = 3))\n", 1)
	m, err := parseModuleText("m.can", src)
	if err != nil {
		t.Fatalf("generic fnref must parse: %v", err)
	}
	for _, d := range m.Decls {
		if f, ok := d.(*FnDecl); ok && f.Name == "m__go" {
			ref := f.Body.Small.Args[0].V
			if ref.Kind != "fnref" || len(ref.TypeArgs) != 1 || ref.TypeArgs[0] != "int" {
				t.Fatalf("ref = %+v, want fnref with [int]", ref)
			}
			return
		}
	}
	t.Fatal("m__go missing after parse")
}

func TestParseFnrefRejectsPositional(t *testing.T) {
	src := strings.Replace(fnvaluesParseGood,
		"  match n < 10\n    true => Ok(\"q\")\n    false => Ok(\"q\")\n",
		"  Ok(fnref m__t(3))\n", 1)
	if _, err := parseModuleText("m.can", src); err == nil {
		t.Fatal("positional capture must be rejected at parse")
	}
}

const fnvaluesInvokeGood = `mod m
  provides [m__go, m__t, M__O]
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

fn m__go(cb: Fn<int, M__O, [m.err]>, n: int) -> M__O rev 1
  emits [m.err]
  tests
    g(3, 4) => Ok("q")
  match invoke cb with n
    on Ok v => Ok(v.value)
    on m.err _ => Ok("e")
`

func TestParseMatchInvoke(t *testing.T) {
	m, err := parseModuleText("m.can", fnvaluesInvokeGood)
	if err != nil {
		t.Fatalf("match invoke must parse: %v", err)
	}
	var body *Node
	for _, d := range m.Decls {
		if f, ok := d.(*FnDecl); ok && f.Name == "m__go" {
			body = f.Body
		}
	}
	if body == nil || !body.IsMatch || body.Kind != MatchInvoke {
		t.Fatalf("body kind = %+v, want MatchInvoke", body)
	}
	if len(body.Scruts) != 1 || body.Scruts[0].Kind != "ref" {
		t.Fatalf("scruts = %+v, want one ref", body.Scruts)
	}
	if body.InvokeArg == nil {
		t.Fatal("InvokeArg missing")
	}
	if len(body.Arms) != 2 {
		t.Fatalf("arms = %d, want 2", len(body.Arms))
	}
	if body.Arms[0].Pats[0].Name != "Ok" || body.Arms[1].Pats[0].Name != "m.err" {
		t.Fatalf("arm patterns parsed as values, want call-style: %+v", body.Arms)
	}
}

func TestParseMatchInvokeRejects(t *testing.T) {
	for _, head := range []string{
		"  match invoke cb\n",
		"  match invoke cb with 6, 7\n",
		"  match invoke call m__t(1) with 2\n",
		"  match invoke 6 with 7\n",
	} {
		src := fnvaluesInvokeGood
		start := strings.Index(src, "  match invoke cb with n\n")
		src = src[:start] + head + "    on Ok v => Ok(v.value)\n    on m.err _ => Ok(\"e\")\n"
		if _, err := parseModuleText("m.can", src); err == nil {
			t.Fatalf("invoke head %q must be rejected at parse", head)
		}
	}
}

// Invocation lives only in match heads: value expressions,
// argument lists, and multi-scrutinee slots never reach the
// invoke path, so each fails ordinary expression parsing.
func TestParseInvokePositionsRefused(t *testing.T) {
	const body = "  match invoke cb with n\n    on Ok v => Ok(v.value)\n    on m.err _ => Ok(\"e\")\n"
	for _, b := range []string{
		"  match n\n    _ => Ok(invoke cb with n)\n",
		"  match n\n    _ => Ok(call m__t(invoke cb with n, 1))\n",
		"  match n, invoke cb with n\n    _, _ => Ok(\"q\")\n",
	} {
		src := strings.Replace(fnvaluesInvokeGood, body, b, 1)
		if _, err := parseModuleText("m.can", src); err == nil {
			t.Fatalf("invoke in %q must be rejected at parse", b)
		}
	}
}

// A bare `match invoke` is not an invoke head: it falls through to
// an ordinary value match, so a variable named invoke keeps working.
func TestParseMatchInvokeBareFallsThrough(t *testing.T) {
	src := fnvaluesInvokeGood
	start := strings.Index(src, "  match invoke cb with n\n")
	arms := "    true => Ok(\"q\")\n    false => Ok(\"e\")\n"
	src = src[:start] + "  match invoke\n" + arms
	m, err := parseModuleText("m.can", src)
	if err != nil {
		t.Fatalf("bare match invoke must fall through: %v", err)
	}
	for _, d := range m.Decls {
		if f, ok := d.(*FnDecl); ok && f.Name == "m__go" {
			if !f.Body.IsMatch || f.Body.Kind != MatchValue {
				t.Fatalf("kind = %v, want MatchValue", f.Body.Kind)
			}
			return
		}
	}
	t.Fatal("m__go missing after parse")
}

func TestParseMatchInvokeGiven(t *testing.T) {
	src := strings.Replace(fnvaluesInvokeGood,
		"  match invoke cb with n\n",
		"  match invoke cb with n\n    given\n      g => [exchange args (n = 4) outcome Ok(\"q\")]\n", 1)
	m, err := parseModuleText("m.can", src)
	if err != nil {
		t.Fatalf("given on invoke must parse: %v", err)
	}
	for _, d := range m.Decls {
		if f, ok := d.(*FnDecl); ok && f.Name == "m__go" {
			if f.Body.Given == nil {
				t.Fatal("Given missing on invoke node")
			}
			return
		}
	}
	t.Fatal("m__go missing after parse")
}

func TestParseFnrefRejectsMalformed(t *testing.T) {
	for _, bad := range []string{
		"  Ok(fnref m__t(divisor = ))\n",
		"  Ok(fnref (divisor = 3))\n",
		"  Ok(fnref m__t(divisor = 3)\n",
	} {
		src := strings.Replace(fnvaluesParseGood,
			"  match n < 10\n    true => Ok(\"q\")\n    false => Ok(\"q\")\n", bad, 1)
		if _, err := parseModuleText("m.can", src); err == nil {
			t.Fatalf("malformed fnref %q must be rejected", bad)
		}
	}
}
