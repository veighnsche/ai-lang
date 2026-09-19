package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func fnVariantSources(t *testing.T) map[string]string {
	t.Helper()
	files := map[string]string{}
	for _, name := range []string{"choice.can", "use.can"} {
		b, err := os.ReadFile(filepath.Join("../sketches/fn-variant", name))
		if err != nil {
			t.Fatal(err)
		}
		files[name] = string(b)
	}
	return files
}

func TestFnVariantSignatures(t *testing.T) {
	prog, mods, ds := genericProgram(t, fnVariantSources(t))
	if errs := genericErrs(ds); len(errs) != 0 {
		t.Fatal(errs)
	}
	ck := newTycker(prog, "", "use__apply$T$int")
	for _, arg := range []string{"int", "str"} {
		fn := prog.Fns["use__apply$T$"+arg]
		sig := fn.Body.invokeSig
		want := "Choice__Value$T$" + arg
		if sig == nil || sig.ret != want || sig.in != arg {
			t.Fatalf("sig = %+v", sig)
		}
		ref := fn.Tests[0].Args[0].V
		got, ok := admissibleRefSig(ck, prog, ref)
		if !ok || !sameType(got, fn.Params[0][1]) {
			t.Fatalf("reference signature = %q %v", got, ok)
		}
		// Forwarding remains explicit reconstruction of the success payload.
		out := fn.Body.Arms[0].Rhs.Small
		if out.Ctor != "Ok" || len(out.Args) != 1 || strings.Join(out.Args[0].V.Ref, ".") != "r.value" {
			t.Fatalf("forward = %+v", out)
		}
		bound := typeDecls(mods)["Use__Bound$T$"+arg]
		if !sameType(bound.Fields[0][1], fn.Params[0][1]) {
			t.Fatalf("bound field = %v", bound.Fields)
		}
		// Fn type, invoke temporary, and source-function success must agree.
		e := &emitter{recs: recordShapes(prog.Modules), variants: variantShapes(prog.Modules), brands: prog.Brands, errTypes: errorShapes(prog.Modules)}
		union, err := e.sigResultUnion(sig)
		if err != nil {
			t.Fatal(err)
		}
		source, err := fnResultUnion(prog.Fns["choice__choose$T$"+arg], prog)
		if err != nil {
			t.Fatal(err)
		}
		head, err := tsTypeB(fn.Params[0][1], prog.Brands, e.recs, e.variants, e.errTypes)
		if err != nil {
			t.Fatal(err)
		}
		if union != source || !strings.HasSuffix(head, " => "+source) {
			t.Fatalf("%s / %s / %s", union, source, head)
		}
	}
}

func TestFnVariantRejects(t *testing.T) {
	tests := []struct{ name, file, old, new, want string }{
		{"wrong_instance", "use.can", "cb: Fn<T, Choice__Value<T>, []>\n", "cb: Fn<T, Choice__Value<str>, []>\n", "want Fn<int"},
		{"scalar_mismatch", "use.can", "cb: Fn<T, Choice__Value<T>, []>\n", "cb: Fn<T, int, []>\n", "want Fn<int,int,[]>"},
		{"wrong_input", "use.can", "cb: Fn<int, Choice__Value<int>, [choice.denied]>", "cb: Fn<str, Choice__Value<int>, [choice.denied]>", "want Fn<str"},
		{"wrong_error_set", "use.can", "cb: Fn<int, Choice__Value<int>, [choice.denied]>", "cb: Fn<int, Choice__Value<int>, []>", "want Fn<int"},
		{"whole_binder", "use.can", "on Ok r => match r.value", "on Ok r => match r", "select r.value"},
		{"binder_as_value", "use.can", "on Ok r => Ok(r.value)", "on Ok r => Ok(r)", "select r.value"},
		{"bad_recovery", "use.can", "on choice.denied e => Ok(Choice__None<int>())", "on choice.denied e => Ok(7)", "want Choice__Value"},
		{"forward_wrong_instance", "use.can", "fn use__checked(cb: Fn<int, Choice__Value<int>, [choice.denied]>, value: int) -> Choice__Value<int>", "fn use__checked(cb: Fn<int, Choice__Value<int>, [choice.denied]>, value: int) -> Choice__Value<str>", "forward needs the same variant return type"},
		{"missing_case", "use.can", "      on Choice__None<int> _ => Ok(0)\n", "", "missing Choice__None"},
		{"missing_outcome", "use.can", "    on choice.denied e => forward e\n", "", "missing choice.denied"},
		{"computed_capture", "use.can", "Ok(fnref choice__choose<T>(present = present))", "Ok(fnref choice__choose<T>(present = not present))", "computed"},
		{"missing_pin", "use.can", "choice__choose@1, ", "", "uses"},
		{"wrong_pin", "use.can", "choice__choose@1", "choice__choose@2", "rev"},
		{"provider_coverage", "choice.can", "    none_int<T=int>(false, 7) => Ok(Choice__None<int>())\n", "", "no test takes on false"},
		{"caller_coverage", "use.can", "    denied(fnref choice__checked(), -1) => choice.denied(-1)\n", "", "no test takes on choice.denied"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := fnVariantSources(t)
			changed := strings.Replace(files[tt.file], tt.old, tt.new, 1)
			if changed == files[tt.file] {
				t.Fatal("fixture unchanged")
			}
			files[tt.file] = changed
			_, _, ds := genericProgram(t, files)
			if !hasDiag(ds, "error", tt.want) {
				t.Fatalf("want %q: %v", tt.want, ds)
			}
		})
	}
}

func TestFnVariantMonomorphic(t *testing.T) {
	src := strings.Replace(monoVariantReturn, "provides [", "provides [m__apply, ", 1) + `
fn m__apply(cb: Fn<bool, M__State, []>, full: bool) -> M__State rev 1
  emits []
  tests
    full(fnref m__choose(), true) => Ok(M__Full(7))
    empty(fnref m__choose(), false) => Ok(M__Empty())
  match invoke cb with full
    on Ok r => forward r
`
	_, _, ds := genericProgram(t, map[string]string{"m.can": src})
	if errs := genericErrs(ds); len(errs) != 0 {
		t.Fatal(errs)
	}
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if ds := diagnose(dir, "m.can", src); len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
}

func TestFnVariantSuccessIsNominal(t *testing.T) {
	src := strings.Replace(monoVariantReturn, "provides [", "provides [Other__State, m__wrong, ", 1) + `
variant Other__State rev 1 (
  case Full(number: int)
  case Empty()
)

fn m__wrong(cb: Fn<bool, Other__State, []>) -> M__State rev 1
  emits []
  tests
    full(fnref m__choose()) => Ok(M__Full(7))
  Ok(M__Full(7))
`
	_, _, ds := genericProgram(t, map[string]string{"m.can": src})
	if !hasDiag(ds, "error", "want Fn<bool, Other__State, []>") {
		t.Fatal(ds)
	}
}

func TestFnVariantPurityAndContainment(t *testing.T) {
	for _, mode := range []string{"effects", "requires", "variant_payload"} {
		t.Run(mode, func(t *testing.T) {
			files := fnVariantSources(t)
			want := ""
			switch mode {
			case "effects":
				files["choice.can"] += "\nstate Choice__Counter: int = 0\n"
				files["choice.can"] = strings.Replace(files["choice.can"], "fn choice__checked(value: int) -> Choice__Value<int> rev 1\n", "fn choice__checked(value: int) -> Choice__Value<int> rev 1\n  effects [Choice__Counter.read]\n", 1)
				want = "effectful function choice__checked"
			case "requires":
				files["choice.can"] = strings.Replace(files["choice.can"], "  emits [choice.denied]\n  tests", "  emits [choice.denied]\n  requires\n    value >= 0\n  tests", 1)
				want = "has a required precondition"
			case "variant_payload":
				files["use.can"] += "\ntype Bad__Holder rev 1 (\n  cb: Fn<int, Choice__Value<int>, []>\n)\n\nvariant Bad__Value rev 1 (\n  case Hidden(holder: Bad__Holder)\n)\n"
				want = "field holder of Bad__Hidden contains a function value"
			}
			_, _, ds := genericProgram(t, files)
			if !hasDiag(ds, "error", want) {
				t.Fatalf("want %s: %v", want, ds)
			}
		})
	}
}

func TestFnVariantInvokeCycle(t *testing.T) {
	src := strings.Replace(fnvaluesInvokeCycle, "type M__O rev 1 (\n  value: int\n)", "variant M__O rev 1 (\n  case Value(value: int)\n)", 1)
	src = strings.ReplaceAll(src, "Ok(3)", "Ok(M__Value(3))")
	src = strings.ReplaceAll(src, "Ok(1)", "Ok(M__Value(1))")
	src = strings.ReplaceAll(src, "on Ok r => Ok(r.value + b)", "on Ok r => forward r")
	src = strings.ReplaceAll(src, "on Ok r => Ok(r.value)", "on Ok r => forward r")
	seqCode(t, map[string]string{"m.can": src}, "m.can", CodeLocalCycle, "m__a invokes m__b through a function value; invocation cannot close a cycle")
}

func TestFnUnsupportedSuccessesStillRefused(t *testing.T) {
	files := fnVariantSources(t)
	prog, _, ds := genericProgram(t, files)
	if errs := genericErrs(ds); len(errs) != 0 {
		t.Fatal(errs)
	}
	ck := newTycker(prog, "", "use__apply$T$int")
	for _, r := range []string{"Fn<int,Choice__Value$T$int,[]>"} {
		head := "Fn<int," + r + ",[]>"
		if ds := ck.fnHeadDiags(head, "test", 1, "Fn"); !hasDiag(ds, "error", "contains a function value") {
			t.Fatalf("admitted %s: %v", head, ds)
		}
		if _, err := tsTypeB(head, prog.Brands, recordShapes(prog.Modules), variantShapes(prog.Modules), errorShapes(prog.Modules)); err == nil {
			t.Fatalf("emitted %s", head)
		}
	}
}

// This separate pre-existing sugar gap was discovered while composing
// the B08 sketch. Explicit match call + forward works; late-parsed
// forward-call text does not see the generic target's stamped name.
func TestFnVariantGenericForwardCallStillRefused(t *testing.T) {
	files := fnVariantSources(t)
	files["use.can"] = strings.Replace(files["use.can"], "on Ok bound => match call use__apply<int>(bound.cb, value)\n      on Ok r => forward r", "on Ok bound => forward call use__apply<int>(bound.cb, value)", 1)
	_, _, ds := genericProgram(t, files)
	if !hasDiag(ds, "error", "forward call cannot resolve callee use__apply") {
		t.Fatal(ds)
	}
}

func TestFnVariantLinkedAndLSP(t *testing.T) {
	files := fnVariantSources(t)
	// Consumer-first ordering matches the CLI golden gate.
	files = map[string]string{"a.can": files["use.can"], "z.can": files["choice.can"]}
	prog, _, ds := genericProgram(t, files)
	if errs := genericErrs(ds); len(errs) != 0 {
		t.Fatal(errs)
	}
	fn := prog.Fns["use__go"]
	for _, row := range fn.Tests {
		ctx, err := freshExecCtx(prog, row.Name, true)
		if err != nil {
			t.Fatal(err)
		}
		env := map[string]*Value{}
		for _, a := range row.Args {
			v, err := evSmall(a.V, env, ctx, fn.Name)
			if err != nil {
				t.Fatal(err)
			}
			env[a.Name] = v
		}
		got, err := evNode(fn.Body, env, ctx, fn.Name)
		if err != nil {
			t.Fatal(err)
		}
		want, err := evSmall(row.Expected, env, ctx, fn.Name)
		if err != nil {
			t.Fatal(err)
		}
		eq, err := vEq(got, want)
		if err != nil || !eq {
			t.Fatalf("linked %s: %v", row.Name, err)
		}
	}
	dir := writeLSPDir(t, files)
	if ds := diagnose(dir, "a.can", files["a.can"]); len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
}

func TestGoldenFnVariant(t *testing.T) {
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/fn-variant/use.can", "../sketches/fn-variant/choice.can"}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"choice.ts", "use.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(out, name))
		if err != nil {
			t.Fatal(err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/fn-variant", name))
		if err != nil {
			t.Fatal(err)
		}
		checkGoldenFile(t, name, got, want)
	}
}

func TestFnVariantNodeParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/fn-variant/use.can", "../sketches/fn-variant/choice.can"}); err != nil {
		t.Fatal(err)
	}
	rewriteSpecifier(t, filepath.Join(out, "use.ts"), "./choice", "./choice.ts")
	driver := `import * as c from "./choice.ts";
import * as u from "./use.ts";
const show = x => JSON.stringify(x,(_,v)=>typeof v === "bigint" ? v.toString()+"n" : v);
const eq = (a,b) => { if(show(a)!==show(b)) throw new Error(show(a)+" != "+show(b)); };
const ok = value => ({$can_kind:"ok",value});
const some = value => ({$can_kind:"Choice__Some$T$int",value});
const none = {$can_kind:"Choice__None$T$int"};
const yes = u.use__bind$T$int(true).cb;
const no = u.use__bind$T$int(false).cb;
for(const n of [7n,9007199254740993n]) {
 eq(u.use__apply$T$int(yes,n),ok(some(n)));
 eq(u.use__apply$T$int(no,n),ok(none));
 eq(u.use__apply$T$int(yes,n),ok(some(n))); // older capture survives
 eq(u.use__go(true,n),ok(some(n)));
 eq(u.use__go(false,n),ok(none));
}
for(const present of [true,false]) {
 const cb = u.use__bind$T$str(present).cb;
 const value = present ? {$can_kind:"Choice__Some$T$str",value:"a"} : {$can_kind:"Choice__None$T$str"};
 eq(u.use__apply$T$str(cb,"a"),ok(value));
}
for(const n of [7n,0n,-1n,9007199254740993n]) {
 const want = n < 0n ? {$can_kind:"choice.denied",value:n} : ok(n === 0n ? none : some(n));
 eq(u.use__checked(c.choice__checked,n),want);
 eq(u.use__recover(c.choice__checked,n),n < 0n ? ok(none) : want);
 eq(u.use__number(c.choice__checked,n),ok(n < 0n ? -1n : n));
}
console.log("FN_VARIANT_PARITY_OK");
`
	if err := os.WriteFile(filepath.Join(out, "driver.ts"), []byte(driver), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", "driver.ts")
	cmd.Dir = out
	raw, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(raw), "FN_VARIANT_PARITY_OK") {
		t.Fatalf("%v\n%s", err, raw)
	}
}
