package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func variantReturnSources(t *testing.T) map[string]string {
	t.Helper()
	files := map[string]string{}
	for _, name := range []string{"pick.can", "use.can"} {
		b, err := os.ReadFile(filepath.Join("../sketches/variant-return", name))
		if err != nil {
			t.Fatal(err)
		}
		files[name] = string(b)
	}
	return files
}

func TestVariantReturnTypesAndForward(t *testing.T) {
	files := variantReturnSources(t)
	prog, _, ds := genericProgram(t, files)
	if errs := genericErrs(ds); len(errs) != 0 {
		t.Fatal(errs)
	}
	for _, name := range []string{"pick__checked", "pick__relay", "use__relay", "use__recover"} {
		if got := prog.Fns[name].Ret; got != "Pick__Value$T$int" {
			t.Fatalf("%s return = %s", name, got)
		}
	}
	// forward rebuilds the envelope from .value, never spreads the
	// variant tag over the Ok discriminant or returns the binder itself.
	f := prog.Fns["use__relay"]
	rhs := f.Body.Arms[0].Rhs.Small
	if rhs.Ctor != "Ok" || len(rhs.Args) != 1 || rhs.Args[0].Name != "value" || strings.Join(rhs.Args[0].V.Ref, ".") != "r.value" {
		t.Fatalf("forward = %+v", rhs)
	}
	for _, arg := range []string{"int", "str"} {
		fn := prog.Fns["pick__choose$T$"+arg]
		shape, err := declaredOkShape(fn, prog)
		if err != nil || len(shape) != 1 || shape["value"] != "Pick__Value$T$"+arg {
			t.Fatalf("shape = %v, %v", shape, err)
		}
	}
	// The call binder is a private envelope; its .value is the nominal
	// variant, whose case-specific payload is typed after elimination.
	label := prog.Fns["use__label"]
	got := label.Body.Arms[0].Rhs.Scruts[0].T
	if got != "Pick__Value$T$int" {
		t.Fatal(got)
	}
}

func TestVariantReturnRejects(t *testing.T) {
	tests := []struct{ name, file, old, new, want string }{
		{"missing_value", "pick.can", "true => Ok(Pick__Some<T>(value))", "true => Ok()", "Ok is missing field value"},
		{"extra_positional", "pick.can", "true => Ok(Pick__Some<T>(value))", "true => Ok(Pick__Some<T>(value), Pick__None<T>())", "Ok takes 2 args for 1 fields"},
		{"wrong_field", "pick.can", "true => Ok(Pick__Some<T>(value))", "true => Ok(other = Pick__Some<T>(value))", "Ok has no field other"},
		{"duplicate_field", "pick.can", "true => Ok(Pick__Some<T>(value))", "true => Ok(value = Pick__Some<T>(value), value = Pick__None<T>())", "Ok repeats field value"},
		{"wrong_payload_type", "pick.can", "true => Ok(Pick__Some<T>(value))", "true => Ok(7)", "want Pick__Value"},
		{"wrong_instance", "pick.can", "true => Ok(Pick__Some<T>(value))", "true => Ok(Pick__None<str>())", "want Pick__Value$T$int"},
		{"nested_ok", "pick.can", "true => Ok(Pick__Some<T>(value))", "true => Ok(Ok(Pick__Some<T>(value)))", "needs a value of variant"},
		{"error_as_data", "pick.can", "true => Ok(Pick__Some<T>(value))", "true => Ok(pick.denied(1))", "needs a value of variant"},
		{"raw_case_body", "pick.can", "true => Ok(Pick__Some<T>(value))", "true => Pick__Some<T>(value)", "variant return must use Ok"},
		{"raw_reference_body", "pick.can", "  Ok(value)\n", "  value\n", "variant return must use Ok"},
		{"bad_expectation", "pick.can", "some_int<T=int>(true, 7) => Ok(Pick__Some<int>(7))", "some_int<T=int>(true, 7) => Ok(7)", "want Pick__Value"},
		{"raw_expectation", "pick.can", "some_int<T=int>(true, 7) => Ok(Pick__Some<int>(7))", "some_int<T=int>(true, 7) => Pick__Some<int>(7)", "write Ok"},
		{"bad_script", "use.can", "outcome Ok(Pick__Some<int>(7))", "outcome Ok(7)", "want Pick__Value"},
		{"nested_ok_script", "use.can", "outcome Ok(Pick__Some<int>(7))", "outcome Ok(Ok(Pick__Some<int>(7)))", "needs a value of variant"},
		{"wrong_recovery", "use.can", "on pick.denied e => Ok(Pick__None<int>())", "on pick.denied e => Ok(7)", "want Pick__Value"},
		{"unwrapped_binder", "use.can", "on Ok r => match r.value", "on Ok r => match r", "select r.value"},
		{"binder_as_variant", "use.can", "on Ok r => Ok(r.value)", "on Ok r => Ok(r)", "select r.value"},
		{"missing_case", "use.can", "      on Pick__None<int> _ => Ok(\"absent\")\n", "", "missing Pick__None"},
		{"unwitnessed_case", "use.can", "    none(0) => Ok(\"absent\")\n", "", "no test takes on Pick__None"},
		{"forward_wrong_instance", "use.can", "fn use__relay(value: int) -> Pick__Value<int>", "fn use__relay(value: int) -> Pick__Value<str>", "forward needs the same variant return type"},
		{"forward_variant_to_record", "use.can", "fn use__relay(value: int) -> Pick__Value<int>", "fn use__relay(value: int) -> Use__Text", "forward needs the same variant return type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := variantReturnSources(t)
			changed := strings.Replace(files[tt.file], tt.old, tt.new, 1)
			if changed == files[tt.file] {
				t.Fatal("fixture unchanged")
			}
			files[tt.file] = changed
			_, _, ds := genericProgram(t, files)
			if !hasDiag(ds, "error", tt.want) {
				t.Fatalf("want %q, got %v", tt.want, ds)
			}
		})
	}
}

// A non-generic variant must work too. The caller eliminates .value,
// projects each case's own fields, and constructs a fresh outcome.
const monoVariantReturn = `mod m
  provides [M__State, m__choose, m__copy]
  uses []
  emits []

variant M__State rev 1 (
  case Full(number: int)
  case Empty()
)

fn m__choose(full: bool) -> M__State rev 1
  emits []
  tests
    full(true) => Ok(M__Full(7))
    empty(false) => Ok(M__Empty())
  match full
    true => Ok(M__Full(7))
    false => Ok(M__Empty())

fn m__copy(full: bool) -> M__State rev 1
  emits []
  tests
    full(true) => Ok(M__Full(7))
    empty(false) => Ok(M__Empty())
  match call m__choose(full)
    on Ok r => match r.value
      on M__Full f => Ok(M__Full(f.number))
      on M__Empty _ => Ok(M__Empty())
`

func TestVariantReturnMonomorphic(t *testing.T) {
	_, _, ds := genericProgram(t, map[string]string{"m.can": monoVariantReturn})
	if errs := genericErrs(ds); len(errs) != 0 {
		t.Fatal(errs)
	}
	dir := writeLSPDir(t, map[string]string{"m.can": monoVariantReturn})
	if ds := diagnose(dir, "m.can", monoVariantReturn); len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
}

func TestVariantReturnBareTypeParameter(t *testing.T) {
	src := strings.Replace(monoVariantReturn, "provides [", "provides [m__id, ", 1)
	src += `
fn m__id<T>(value: T) -> T rev 1
  emits []
  tests
    full<T=M__State>(M__Full(7)) => Ok(M__Full(7))
    empty<T=M__State>(M__Empty()) => Ok(M__Empty())
  Ok(value)
`
	prog, _, ds := genericProgram(t, map[string]string{"m.can": src})
	if errs := genericErrs(ds); len(errs) != 0 {
		t.Fatal(errs)
	}
	if shape, err := declaredOkShape(prog.Fns["m__id$T$M__State"], prog); err != nil || shape["value"] != "M__State" {
		t.Fatalf("%v, %v", shape, err)
	}
}

func TestVariantReturnNominalParent(t *testing.T) {
	src := strings.Replace(monoVariantReturn, "provides [", "provides [Other__State, ", 1)
	src += "\nvariant Other__State rev 1 (\n  case Empty()\n)\n"
	src = strings.Replace(src, "false => Ok(M__Empty())", "false => Ok(Other__Empty())", 1)
	_, _, ds := genericProgram(t, map[string]string{"m.can": src})
	if !hasDiag(ds, "error", "got Other__State, want M__State") {
		t.Fatal(ds)
	}
}

func TestVariantReturnFnAndExternStillRefused(t *testing.T) {
	src := monoVariantReturn + `
extern host__choose(full: bool) -> M__State rev 1
  emits []

fn m__invoke(cb: Fn<bool, M__State, []>, full: bool) -> M__State rev 1
  emits []
  tests
    full(fnref m__choose(), true) => Ok(M__Full(7))
  match invoke cb with full
    on Ok r => Ok(r.value)
`
	src = strings.Replace(src, "provides [", "provides [host__choose, m__invoke, ", 1)
	_, _, ds := genericProgram(t, map[string]string{"m.can": src})
	for _, want := range []string{"externs return a record type", "Fn success M__State", "success M__State is not a record"} {
		if !hasDiag(ds, "error", want) {
			t.Fatalf("want %s: %v", want, ds)
		}
	}
}

func TestVariantReturnScriptCannotInventSuccess(t *testing.T) {
	files := variantReturnSources(t)
	files["use.can"] = strings.Replace(files["use.can"], "outcome Ok(Pick__Some<int>(7))", "outcome Ok(Pick__Some<int>(99))", 1)
	_, _, ds := genericProgram(t, files)
	if len(genericErrs(ds)) == 0 {
		t.Fatal("invented provider success accepted")
	}
}

func TestGoldenVariantReturn(t *testing.T) {
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/variant-return/use.can", "../sketches/variant-return/pick.can"}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"pick.ts", "use.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(out, name))
		if err != nil {
			t.Fatal(err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/variant-return", name))
		if err != nil {
			t.Fatal(err)
		}
		checkGoldenFile(t, name, got, want)
	}
}

func TestVariantReturnNodeParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/variant-return/use.can", "../sketches/variant-return/pick.can"}); err != nil {
		t.Fatal(err)
	}
	rewriteSpecifier(t, filepath.Join(out, "use.ts"), "./pick", "./pick.ts")
	driver := `import * as p from "./pick.ts";
import * as u from "./use.ts";
const show = x => JSON.stringify(x, (_,v) => typeof v === "bigint" ? v.toString()+"n" : v);
const eq = (a,b) => { if(show(a)!==show(b)) throw new Error(show(a)+" != "+show(b)); };
const ok = value => ({$can_kind:"ok",value});
const some = value => ({$can_kind:"Pick__Some$T$int",value});
const none = {$can_kind:"Pick__None$T$int"};
for(const n of [7n,9007199254740993n,0n,-1n]) {
 const want = n < 0n ? {$can_kind:"pick.denied",value:n} : ok(n === 0n ? none : some(n));
 eq(p.pick__checked(n),want);
 eq(p.pick__relay(n),want);
 eq(u.use__relay(n),want);
 eq(u.use__recover(n),n < 0n ? ok(none) : want);
 eq(u.use__label(n),ok(n < 0n ? "denied" : n === 0n ? "absent" : "present"));
}
for(const [present,value] of [[true,7n],[false,7n]]) {
 const want = ok(present ? some(value) : none);
 eq(p.pick__choose$T$int(present,value),want);
 eq(p.pick__relay_generic$T$int(present,value),want);
 eq(p.pick__identity$T$int(want.value),want);
}
for(const present of [true,false]) {
 const value = present ? {$can_kind:"Pick__Some$T$str",value:"a"} : {$can_kind:"Pick__None$T$str"};
 eq(p.pick__choose$T$str(present,"a"),ok(value));
 eq(p.pick__relay_generic$T$str(present,"a"),ok(value));
 eq(p.pick__identity$T$str(value),ok(value));
}
console.log("VARIANT_RETURN_PARITY_OK");
`
	if err := os.WriteFile(filepath.Join(out, "driver.ts"), []byte(driver), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", "driver.ts")
	cmd.Dir = out
	raw, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(raw), "VARIANT_RETURN_PARITY_OK") {
		t.Fatalf("%v\n%s", err, raw)
	}
}
