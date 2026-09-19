package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func genericOptionSource(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile("../sketches/generic-option/option.can")
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestGenericVariantStamps(t *testing.T) {
	src := genericOptionSource(t)
	prog, mods, ds := genericProgram(t, map[string]string{"option.can": src})
	if errs := genericErrs(ds); len(errs) != 0 {
		t.Fatal(errs)
	}
	if _, ok := prog.Variants["Option__Value"]; ok {
		t.Fatal("template survived")
	}
	for _, arg := range []string{"int", "str"} {
		parent := mangleInstance("Option__Value", []string{arg})
		v := prog.Variants[parent]
		if v == nil || v.Rev != 1 || len(v.TypeParams) != 0 || len(v.Cases) != 2 {
			t.Fatalf("bad stamp: %+v", v)
		}
		if v.Cases[0].Fields[0][1] != arg || len(v.Cases[1].Fields) != 0 {
			t.Fatalf("bad payloads: %+v", v.Cases)
		}
		for _, short := range []string{"Some", "None"} {
			caseName := mangleInstance("Option__"+short, []string{arg})
			if prog.Cases[caseName] != parent {
				t.Fatalf("%s belongs to %s", caseName, prog.Cases[caseName])
			}
		}
		fn := prog.Fns[mangleInstance("option__value_or", []string{arg})]
		if fn.Params[0][1] != parent || len(fn.Tests) != 2 {
			t.Fatalf("bad fn stamp: %+v", fn)
		}
		for _, arm := range fn.Body.Arms {
			p := arm.Pats[0]
			if len(p.TypeArgs) != 0 || !strings.HasSuffix(p.Name, "$T$"+arg) {
				t.Fatalf("bad pattern: %+v", p)
			}
		}
	}
	// A second expansion is a no-op on already stamped modules.
	if ds := expandGenerics(mods, map[string]string{"option.can": src}); len(ds) != 0 {
		t.Fatal(ds)
	}
}

func TestGenericVariantRejects(t *testing.T) {
	src := genericOptionSource(t)
	tests := []struct{ name, old, new, want string }{
		{"missing_ctor_args", "Option__Some<int>(value = 7)", "Option__Some(value = 7)", "needs explicit type arguments"},
		{"missing_nullary_args", "Option__None<int>()", "Option__None()", "needs explicit type arguments"},
		{"missing_pattern_args", "on Option__Some<T> s", "on Option__Some s", "needs explicit type arguments"},
		{"ctor_arity", "Option__Some<int>(value = 7)", "Option__Some<int,str>(value = 7)", "takes 1 type arguments, got 2"},
		{"pattern_arity", "on Option__Some<T> s", "on Option__Some<T,T> s", "takes 1 type arguments, got 2"},
		{"unknown_pattern_param", "on Option__Some<T> s", "on Option__Some<U> s", "not a type parameter"},
		{"nested_arg", "Option__Some<int>(value = 7)", "Option__Some<Option__Value<int>>(value = 7)", "nested instantiation"},
		{"bare_annotation", "value: Option__Value<T>,", "value: Option__Value,", "unknown type Option__Value"},
		{"wrong_instance_arg", "Option__Some<int>(value = 7)", "Option__Some<str>(value = \"wrong\")", "want Option__Value$T$int"},
		{"wrong_payload", "Option__Some<int>(value = 7)", "Option__Some<int>(value = \"wrong\")", "want int"},
		{"wrong_instance_pattern", "on Option__Some<T> s", "on Option__Some<str> s", "stale match arm"},
		{"missing_case", "    on Option__None<T> _ => Ok(alternative)\n", "", "missing Option__None"},
		{"duplicate_case", "    on Option__None<T> _ => Ok(alternative)", "    on Option__None<T> _ => Ok(alternative)\n    on Option__None<T> _ => Ok(alternative)", "duplicate match arm"},
		{"unwitnessed_instance_arm", "    none_str<T=str>(Option__None<str>(), \"fallback\") => Ok(\"fallback\")\n", "", "no test takes on Option__None"},
		{"unused_parameter", "variant Option__Value<T>", "variant Option__Value<T,U>", "never used in its fields"},
		{"self_recursive", "case Some(value: T)", "case Some(value: T, next: Option__Value<T>)", "fields its own base"},
		{"seq_variant", "value: Option__Value<T>\n", "value: Seq<Option__Value<T>>\n", "variant sequences"},
		{"construct_parent", "Option__Some<int>(value = 7)", "Option__Value<int>(value = 7)", "unknown type Option__Value"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changed := strings.Replace(src, tt.old, tt.new, 1)
			if changed == src {
				t.Fatal("fixture did not change")
			}
			_, _, ds := genericProgram(t, map[string]string{"option.can": changed})
			if !hasDiag(ds, "error", tt.want) {
				t.Fatalf("want %q, got %v", tt.want, ds)
			}
		})
	}
}

func TestGenericVariantParseRejects(t *testing.T) {
	for _, params := range []string{"T,T", "t", "T_"} {
		src := strings.Replace(genericOptionSource(t), "variant Option__Value<T>", "variant Option__Value<"+params+">", 1)
		if _, err := parseModuleText("option.can", src); err == nil {
			t.Fatalf("accepted parameters %s", params)
		}
	}
	for _, pattern := range []string{"Ok<int, str> r", "math.bad<int> e", "Option__Some<> s", "Option__Some<int>"} {
		if _, err := parsePattern(pattern); err == nil {
			t.Fatalf("accepted %s", pattern)
		}
	}
}

func TestGenericVariantCaseCollisions(t *testing.T) {
	for _, extra := range []string{
		"variant Option__Other<U> rev 1 (\n  case Some(other: U)\n)\n",
		"variant Option__Other rev 1 (\n  case Some(other: int)\n)\n",
		"type Option__Some rev 1 (\n  other: int\n)\n",
	} {
		_, _, ds := genericProgram(t, map[string]string{"option.can": genericOptionSource(t) + "\n" + extra})
		if !hasDiag(ds, "error", "case identity collision") {
			t.Fatal(ds)
		}
	}
}

// Generic payload references use the same fixpoint as record fields.
// No direct construction of Wrap__Value is needed to demand its stamp.
func TestGenericVariantPayloadFixpoint(t *testing.T) {
	src := `mod m
  provides [m__read, Wrap__Value, Box__Value, M__Out]
  uses []
  emits []

variant Wrap__Value<T> rev 1 (
  case Packed(box: Box__Value<T>)
)

type Box__Value<T> rev 1 (
  value: T
)

type M__Out rev 1 (
  value: str
)

fn m__read(value: Wrap__Value<str>) -> M__Out rev 1
  emits []
  tests
    one(Wrap__Packed<str>(Box__Value<str>("a"))) => Ok("a")
  match value
    on Wrap__Packed<str> p => Ok(p.box.value)
`
	prog, _, ds := genericProgram(t, map[string]string{"m.can": src})
	if errs := genericErrs(ds); len(errs) != 0 {
		t.Fatal(errs)
	}
	if got := prog.Variants["Wrap__Value$T$str"].Cases[0].Fields[0][1]; got != "Box__Value$T$str" {
		t.Fatal(got)
	}
}

func TestGenericVariantPins(t *testing.T) {
	use, err := os.ReadFile("../sketches/generic-option/use.can")
	if err != nil {
		t.Fatal(err)
	}
	files := map[string]string{"a.can": string(use), "z.can": genericOptionSource(t)}
	_, mods, ds := genericProgram(t, files)
	if errs := genericErrs(ds); len(errs) != 0 {
		t.Fatal(errs)
	}
	if got := strings.Join(mods[0].Hdr["uses"], ","); got != "Option__Value$T$int@1,Option__Value$T$str@1,option__value_or$T$int@1" {
		t.Fatal(got)
	}
	files["a.can"] = strings.Replace(string(use), "Option__Value@1", "Option__Value@2", 1)
	_, _, ds = genericProgram(t, files)
	if len(genericErrs(ds)) == 0 {
		t.Fatal("wrong revision pin accepted")
	}
}

func TestGoldenGenericOption(t *testing.T) {
	out := t.TempDir()
	// Consumer first exercises linked preparation as well as foreign types.
	if err := compile(out, []string{"../sketches/generic-option/use.can", "../sketches/generic-option/option.can"}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"option.ts", "use.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(out, name))
		if err != nil {
			t.Fatal(err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/generic-option", name))
		if err != nil {
			t.Fatal(err)
		}
		checkGoldenFile(t, name, got, want)
	}
}

func TestGenericOptionNodeParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/generic-option/use.can", "../sketches/generic-option/option.can"}); err != nil {
		t.Fatal(err)
	}
	rewriteSpecifier(t, filepath.Join(out, "use.ts"), "./option", "./option.ts")
	driver := `import * as opt from "./option.ts";
import {use__number, use__text} from "./use.ts";
const n = 9007199254740993n;
const someInt = opt.option__some$T$int(n).value;
const noneInt = opt.option__none$T$int(0n).value;
const someStr = opt.option__some$T$str("present").value;
const noneStr = opt.option__none$T$str("").value;
for (const [got,want] of [
 [use__number(someInt).value,n], [use__number(noneInt).value,9n],
 [use__text(someStr).value,"present"], [use__text(noneStr).value,"fallback"],
 [opt.option__value_or$T$str(someStr,"other").value,"present"],
 [opt.option__value_or$T$str(noneStr,"other").value,"other"]]) {
 if (got !== want) throw new Error("generic variant parity failed");
}
if (someInt.$can_kind === someStr.$can_kind || noneInt.$can_kind === noneStr.$can_kind) throw new Error("instance tags collided");
console.log("GENERIC_VARIANT_PARITY_OK");
`
	if err := os.WriteFile(filepath.Join(out, "driver.ts"), []byte(driver), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", "driver.ts")
	cmd.Dir = out
	raw, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(raw), "GENERIC_VARIANT_PARITY_OK") {
		t.Fatalf("%v\n%s", err, raw)
	}
}
