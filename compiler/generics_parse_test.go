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
