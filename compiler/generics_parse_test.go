package main

import (
	"strings"
	"testing"
)

// G1 parse: a generic fn declares type parameters, rows pin
// complete instantiations, and calls instantiate explicitly.
// Expansion (not parsing) validates names against the decl.
const genericParseGood = `mod m
  provides [m__sel, M__O]
  uses []
  emits []

type M__O rev 1 (
  value: str
)

fn m__sel<T>(condition: bool, when_true: T, when_false: T) -> T rev 1
  emits []
  tests
    pick<T=str>(true, "a", "b") => Ok("a")
  match condition
    true => Ok(when_true)
    false => Ok(when_false)

fn m__go(flag: bool) -> M__O rev 1
  emits []
  tests
    one(true) => Ok("a")
  match call m__sel<str>(flag, "a", "b")
    on Ok v => Ok(v.value)
`

func TestParseGenericDecl(t *testing.T) {
	m, err := parseModuleText("m.can", genericParseGood)
	if err != nil {
		t.Fatalf("generic decl must parse: %v", err)
	}
	var fn *FnDecl
	for _, d := range m.Decls {
		if f, ok := d.(*FnDecl); ok && f.Name == "m__sel" {
			fn = f
		}
	}
	if fn == nil {
		t.Fatal("m__sel missing after parse")
	}
	if len(fn.TypeParams) != 1 || fn.TypeParams[0] != "T" {
		t.Fatalf("type params = %v, want [T]", fn.TypeParams)
	}
	if len(fn.Tests) != 1 || len(fn.Tests[0].TypeBinds) != 1 {
		t.Fatalf("row binds missing: %+v", fn.Tests)
	}
	b := fn.Tests[0].TypeBinds[0]
	if b[0] != "T" || b[1] != "str" {
		t.Fatalf("bind = %v, want [T str]", b)
	}
}

func TestParseGenericCall(t *testing.T) {
	m, err := parseModuleText("m.can", genericParseGood)
	if err != nil {
		t.Fatalf("generic call must parse: %v", err)
	}
	var found []string
	for _, d := range m.Decls {
		fn, ok := d.(*FnDecl)
		if !ok || fn.Name != "m__go" {
			continue
		}
		for _, c := range walkCalls(fn.Body) {
			if c.Fname == "m__sel" {
				found = c.TypeArgs
			}
		}
	}
	if len(found) != 1 || found[0] != "str" {
		t.Fatalf("call type args = %v, want [str]", found)
	}
}

func TestParseGenericMultiParam(t *testing.T) {
	src := strings.Replace(genericParseGood,
		"fn m__sel<T>(condition: bool, when_true: T, when_false: T) -> T rev 1",
		"fn m__sel<T, U>(first: T, second: U) -> T rev 1", 1)
	src = strings.Replace(src,
		`pick<T=str>(true, "a", "b") => Ok("a")`,
		`pick<T=str,U=int>("a", 1) => Ok("a")`, 1)
	m, err := parseModuleText("m.can", src)
	if err != nil {
		t.Fatalf("multi-param decl must parse: %v", err)
	}
	for _, d := range m.Decls {
		fn, ok := d.(*FnDecl)
		if !ok || fn.Name != "m__sel" {
			continue
		}
		if len(fn.TypeParams) != 2 || fn.TypeParams[1] != "U" {
			t.Fatalf("type params = %v, want [T U]", fn.TypeParams)
		}
		if len(fn.Tests[0].TypeBinds) != 2 {
			t.Fatalf("binds = %v, want 2 pairs", fn.Tests[0].TypeBinds)
		}
	}
}

func TestParseGenericSeqArg(t *testing.T) {
	src := strings.Replace(genericParseGood,
		"match call m__sel<str>(flag, \"a\", \"b\")",
		"match call m__sel<Seq<str>>(flag, \"a\", \"b\")", 1)
	m, err := parseModuleText("m.can", src)
	if err != nil {
		t.Fatalf("Seq type arg must parse: %v", err)
	}
	for _, d := range m.Decls {
		fn, ok := d.(*FnDecl)
		if !ok || fn.Name != "m__go" {
			continue
		}
		for _, c := range walkCalls(fn.Body) {
			if c.Fname == "m__sel" && (len(c.TypeArgs) != 1 || c.TypeArgs[0] != "Seq<str>") {
				t.Fatalf("call type args = %v, want [Seq<str>]", c.TypeArgs)
			}
		}
	}
}

func TestParseGenericScrutineeSplit(t *testing.T) {
	// Multi-param generic calls survive scrutinee splitting;
	// comparisons and Seq heads split exactly as before.
	lists := []struct {
		in   string
		want []string
	}{
		{`call f<str,int>(x), y`, []string{`call f<str,int>(x)`, `y`}},
		{`call f<Seq<str>>(x)`, []string{`call f<Seq<str>>(x)`}},
		{`a < b, c`, []string{`a < b`, `c`}},
		{`a<b, c`, []string{`a<b`, `c`}},
		{`call f(a<b), y`, []string{`call f(a<b)`, `y`}},
	}
	for _, l := range lists {
		got, err := splitMatchList(l.in)
		if err != nil {
			t.Fatalf("split %q: %v", l.in, err)
		}
		if len(got) != len(l.want) {
			t.Fatalf("split %q = %v, want %v", l.in, got, l.want)
		}
		for i := range got {
			if got[i] != l.want[i] {
				t.Fatalf("split %q = %v, want %v", l.in, got, l.want)
			}
		}
	}
}

// G2 parse: a generic record declares type parameters, and
// constructions instantiate explicitly. Expansion (not parsing)
// validates bases, arity, and argument shape.
const genericTypeParseGood = `mod m
  provides [m__wrap, M__Box]
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
`

func TestParseGenericTypeDecl(t *testing.T) {
	m, err := parseModuleText("m.can", genericTypeParseGood)
	if err != nil {
		t.Fatalf("generic type decl must parse: %v", err)
	}
	var td *TypeDecl
	for _, d := range m.Decls {
		if tt, ok := d.(*TypeDecl); ok && tt.Name == "M__Box" {
			td = tt
		}
	}
	if td == nil {
		t.Fatal("M__Box missing after parse")
	}
	if len(td.TypeParams) != 1 || td.TypeParams[0] != "T" {
		t.Fatalf("type params = %v, want [T]", td.TypeParams)
	}
	if len(td.Fields) != 1 || td.Fields[0][0] != "value" || td.Fields[0][1] != "T" {
		t.Fatalf("fields = %v, want [[value T]]", td.Fields)
	}
	var fn *FnDecl
	for _, d := range m.Decls {
		if f, ok := d.(*FnDecl); ok && f.Name == "m__wrap" {
			fn = f
		}
	}
	if fn == nil || fn.Ret != "M__Box<str>" {
		t.Fatalf("ret = %q, want M__Box<str>", fn.Ret)
	}
}

func TestParseGenericCtor(t *testing.T) {
	s, err := parseSmall(`M__Box<str>("a")`)
	if err != nil {
		t.Fatalf("generic construction must parse: %v", err)
	}
	if s.Kind != "ctor" || s.Ctor != "M__Box" {
		t.Fatalf("ctor = %+v, want M__Box construction", s)
	}
	if len(s.TypeArgs) != 1 || s.TypeArgs[0] != "str" {
		t.Fatalf("ctor type args = %v, want [str]", s.TypeArgs)
	}
	// Head-first and head-last expressions split around the
	// head, never inside it.
	for _, expr := range []string{`M__Box<str>(x) == y`, `x == M__Box<str>(v)`, `M__Box<str>(x) or M__Box<str>(y)`} {
		if _, err := parseSmall(expr); err != nil {
			t.Fatalf("expr %s must parse: %v", expr, err)
		}
	}
	// Chained comparisons keep their old readings: only
	// `w<x>(` is a head.
	cmp, err := parseSmall(`a<b>c`)
	if err != nil {
		t.Fatalf("chained comparison must parse: %v", err)
	}
	if cmp.Kind != "binop" || cmp.Op != "<" {
		t.Fatalf("a<b>c = %+v, want outer <", cmp)
	}
	// A bare mention in value position fails precisely.
	if _, err := parseSmall(`M__Box<str>`); err == nil || !strings.Contains(err.Error(), "is not a value") {
		t.Fatalf("bare mention must fail precisely, got %v", err)
	}
}

func TestParseGenericHeadSplit(t *testing.T) {
	lists := []struct {
		in   string
		want []string
	}{
		{`M__Box<str,int>(x), y`, []string{`M__Box<str,int>(x)`, `y`}},
		{`M__Box<str>(x), M__Box<int>(y)`, []string{`M__Box<str>(x)`, `M__Box<int>(y)`}},
		{`a<b, c`, []string{`a<b`, `c`}},
		{`f(a<b), y`, []string{`f(a<b)`, `y`}},
	}
	for _, l := range lists {
		got, err := splitMatchList(l.in)
		if err != nil {
			t.Fatalf("split %q: %v", l.in, err)
		}
		if len(got) != len(l.want) {
			t.Fatalf("split %q = %v, want %v", l.in, got, l.want)
		}
		for i := range got {
			if got[i] != l.want[i] {
				t.Fatalf("split %q = %v, want %v", l.in, got, l.want)
			}
		}
	}
}

func TestParseGenericSeqElem(t *testing.T) {
	s, err := parseSmall(`Seq<M__Box<str>>[M__Box<str>("a")]`)
	if err != nil {
		t.Fatalf("Seq of instance must parse: %v", err)
	}
	if s.Kind != "seqlit" || s.Elem != "M__Box<str>" {
		t.Fatalf("seqlit = %+v, want elem M__Box<str>", s)
	}
	if _, err := parseSmall(`Seq<Seq<str>>["a"]`); err == nil || !strings.Contains(err.Error(), "no nesting") {
		t.Fatalf("nested Seq must keep its parse error, got %v", err)
	}
}

func TestParseGenericRejects(t *testing.T) {
	base := genericParseGood
	cases := []struct {
		name string
		old  string
		new  string
		want string
	}{
		{"lowercase param", "fn m__sel<T>(", "fn m__sel<t>(", "bad type parameter"},
		{"duplicate param", "fn m__sel<T>(", "fn m__sel<T, T>(", "duplicate type parameter"},
		{"empty args", "m__sel<str>(flag", "m__sel<>(flag", "bad call syntax"},
		{"dangling comma", "m__sel<str>(flag", "m__sel<str,>(flag", "empty slot"},
		{"bind without equals", "pick<T=str>(", "pick<T>(", "bad type bind"},
		{"bind empty arg", "pick<T=str>(", "pick<T=>(", "bad type bind"},
	}
	for _, c := range cases {
		src := strings.Replace(base, c.old, c.new, 1)
		if src == base {
			t.Fatalf("%s: fixture anchor missing", c.name)
		}
		_, err := parseModuleText("m.can", src)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Fatalf("%s: expected %q, got %v", c.name, c.want, err)
		}
	}
}
