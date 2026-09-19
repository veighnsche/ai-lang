package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func scalarSources(t *testing.T) map[string]string {
	t.Helper()
	files := map[string]string{}
	for _, name := range []string{"scalar.can", "use.can"} {
		b, err := os.ReadFile(filepath.Join("../sketches/fn-scalar", name))
		if err != nil {
			t.Fatal(err)
		}
		files[name] = string(b)
	}
	return files
}

func TestScalarSuccessSignatures(t *testing.T) {
	prog, _, ds := genericProgram(t, scalarSources(t))
	if errs := genericErrs(ds); len(errs) != 0 {
		t.Fatal(errs)
	}
	ck := newTycker(prog, "", "use__apply$T$int$T$int")
	e := &emitter{recs: recordShapes(prog.Modules), variants: variantShapes(prog.Modules), brands: prog.Brands, errTypes: errorShapes(prog.Modules)}
	for _, typ := range []string{"int", "str", "bool", "dec"} {
		fn := prog.Fns["use__apply$T$"+typ+"$T$"+typ]
		sig := fn.Body.invokeSig
		if sig == nil || sig.in != typ || sig.ret != typ {
			t.Fatalf("%s: %+v", typ, sig)
		}
		ref := fn.Tests[0].Args[0].V
		if got, ok := admissibleRefSig(ck, prog, ref); !ok || !sameType(got, fn.Params[0][1]) {
			t.Fatalf("%s signature: %s %v", typ, got, ok)
		}
		union, err := e.sigResultUnion(sig)
		if err != nil {
			t.Fatal(err)
		}
		source, err := fnResultUnion(prog.Fns["scalar__identity$T$"+typ], prog)
		if err != nil {
			t.Fatal(err)
		}
		head, err := tsTypeB(fn.Params[0][1], prog.Brands, e.recs, e.variants, e.errTypes)
		if err != nil || union != source || !strings.HasSuffix(head, " => "+source) {
			t.Fatalf("%s: %s / %s / %s: %v", typ, union, source, head, err)
		}
		out := fn.Body.Arms[0].Rhs.Small
		if out.Ctor != "Ok" || len(out.Args) != 1 || strings.Join(out.Args[0].V.Ref, ".") != "r.value" {
			t.Fatalf("forward: %+v", out)
		}
	}
	// Test-only foreign generic references must retain their pinned stamps.
	for _, typ := range []string{"int", "str", "bool"} {
		if !prog.Uses["scalar__identity$T$"+typ] {
			t.Fatalf("lost test-only reference pin for %s", typ)
		}
	}
	// An error payload projection is data, not propagation of its error.
	eachRaise(prog.Fns["use__recover"], func(kind string, line int) {
		t.Errorf("recovery falsely raises %s at %d", kind, line)
	})
	var raised []string
	eachRaise(prog.Fns["use__checked"], func(kind string, _ int) { raised = append(raised, kind) })
	if len(raised) != 1 || raised[0] != "scalar.negative" {
		t.Fatalf("lost actual raise: %v", raised)
	}
}

func TestScalarSuccessRejects(t *testing.T) {
	cases := []struct{ name, file, old, new, want string }{
		{"wrong_body", "scalar.can", "  Ok(value + delta)", "  Ok(true)", "got bool, want int"},
		{"wrong_generic_body", "scalar.can", "  Ok(value)\n", "  Ok(7)\n", "got int, want bool"},
		{"missing_value", "scalar.can", "  Ok(value + delta)", "  Ok()", "Ok is missing field value"},
		{"unknown_field", "scalar.can", "  Ok(value + delta)", "  Ok(answer = value)", "Ok has no field answer"},
		{"extra_field", "scalar.can", "  Ok(value + delta)", "  Ok(value = value, extra = delta)", "Ok has no field extra"},
		{"repeat", "scalar.can", "  Ok(value + delta)", "  Ok(value = value, value = delta)", "Ok repeats field value"},
		{"nested_ok", "scalar.can", "  Ok(value + delta)", "  Ok(Ok(value))", "needs a value of scalar int"},
		{"error_as_data", "scalar.can", "    false => Ok(value)", "    false => Ok(scalar.negative(value))", "needs a value of scalar int"},
		{"raw_value", "scalar.can", "  Ok(value + delta)", "  value + delta", "scalar return must use Ok"},
		{"raw_bool", "scalar.can", "  Ok(value > 0)", "  value > 0", "scalar return must use Ok"},
		{"bad_expected", "scalar.can", "plus(3, 7) => Ok(10)", "plus(3, 7) => Ok(false)", "got bool, want int"},
		{"no_int_dec_coercion", "scalar.can", "  Ok(value + value)", "  Ok(2)", "got int, want dec"},
		{"unbound", "scalar.can", "  Ok(value + delta)", "  Ok(missing)", "unbound"},
		{"bad_projection", "scalar.can", "  Ok(value + delta)", "  Ok(value.missing)", "no field missing on value"},
		{"whole_binder", "use.can", "    on Ok r => Ok(r.value)", "    on Ok r => Ok(r)", "select r.value"},
		{"binder_field", "use.can", "    on Ok r => Ok(r.value)", "    on Ok r => Ok(r.missing)", "no field missing on r"},
		{"bad_recovery", "use.can", "Ok(e.value + 3)", "Ok(false)", "got bool, want int"},
		{"bad_script", "use.can", "outcome Ok(7)", "outcome Ok(false)", "got bool, want int"},
		{"script_nested_ok", "use.can", "outcome Ok(7)", "outcome Ok(Ok(7))", "needs a value of scalar int"},
		{"script_lie", "use.can", "outcome Ok(7)", "outcome Ok(8)", "contradicts scalar__checked"},
		{"fn_result_mismatch", "use.can", "Fn<int, int, []>", "Fn<int, bool, []>", "want Fn<int, bool, []>"},
		{"fn_input_mismatch", "use.can", "Fn<int, int, []>", "Fn<str, int, []>", "want Fn<str, int, []>"},
		{"fn_error_mismatch", "use.can", "cb: Fn<int, int, [scalar.negative]>", "cb: Fn<int, int, []>", "want Fn<int, int, []>"},
		{"forward_mismatch", "use.can", "fn use__relay(value: int) -> int", "fn use__relay(value: int) -> bool", "forward needs the same scalar return type"},
		{"missing_outcome", "use.can", "    on scalar.negative e => forward e\n", "", "missing scalar.negative"},
		{"caller_coverage", "use.can", "    no(fnref scalar__checked(), -3) => scalar.negative(-3)\n", "", "no test takes on scalar.negative"},
		{"provider_coverage", "scalar.can", "    no(-3) => scalar.negative(-3)\n", "", "no test takes on true"},
		{"missing_ref_pin", "use.can", "scalar__identity@1, ", "", "not in uses"},
		{"wrong_ref_pin", "use.can", "scalar__identity@1", "scalar__identity@2", "pins rev 2"},
		{"missing_bool_case", "use.can", "      false => Ok(\"not positive\")\n", "", "bool match must be exactly true+false"},
		{"unwitnessed_bool_case", "use.can", "    no(fnref scalar__positive(), 0) => Ok(\"not positive\")\n", "", "no test takes on false"},
		{"computed_capture", "use.can", "Ok(fnref scalar__add(delta = delta))", "Ok(fnref scalar__add(delta = delta + 1))", "computed"},
		{"requires", "scalar.can", "fn scalar__add(delta: int, value: int) -> int rev 1\n  emits []", "fn scalar__add(delta: int, value: int) -> int rev 1\n  emits []\n  requires\n    value >= 0", "has a required precondition"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			files := scalarSources(t)
			if !strings.Contains(files[tt.file], tt.old) {
				t.Fatal("stale fixture")
			}
			files[tt.file] = strings.Replace(files[tt.file], tt.old, tt.new, 1)
			_, _, ds := genericProgram(t, files)
			if !hasDiag(ds, "error", tt.want) {
				t.Fatalf("want %q: %v", tt.want, ds)
			}
		})
	}
}

func TestScalarSuccessPurityAndCycles(t *testing.T) {
	files := scalarSources(t)
	files["scalar.can"] += "\nstate Scalar__Counter: int = 0\n"
	files["scalar.can"] = strings.Replace(files["scalar.can"], "fn scalar__add(delta: int, value: int) -> int rev 1\n", "fn scalar__add(delta: int, value: int) -> int rev 1\n  effects [Scalar__Counter.read]\n", 1)
	_, _, ds := genericProgram(t, files)
	if !hasDiag(ds, "error", "effectful function scalar__add") {
		t.Fatal(ds)
	}
	src := strings.ReplaceAll(fnvaluesInvokeCycle, "-> M__O rev 1", "-> int rev 1")
	src = strings.ReplaceAll(src, "Fn<int, M__O, []>", "Fn<int, int, []>")
	seqCode(t, map[string]string{"m.can": src}, "m.can", CodeLocalCycle, "m__a invokes m__b through a function value; invocation cannot close a cycle")
}

func TestScalarSuccessForwardIsNotStructural(t *testing.T) {
	prog, _, ds := genericProgram(t, scalarSources(t))
	if len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
	// An otherwise identical one-field record is still not the scalar.
	prog.Modules[0].Decls = append(prog.Modules[0].Decls, &TypeDecl{Name: "Box", Fields: [][2]string{{"value", "int"}}})
	for _, pair := range [][2]string{{"int", "Box"}, {"Box", "int"}, {"int", "dec"}, {"str", "dec"}} {
		if _, err := forwardOkRet(prog, &FnDecl{Name: "dst", Ret: pair[1]}, "src", pair[0], "r"); err == nil {
			t.Fatalf("forward admitted %v", pair)
		}
	}
}

func TestScalarSuccessContracts(t *testing.T) {
	// Preserve the declared envelope both at own exits and through verified
	// callee summaries; example rows alone do not prove the contract.
	src := strings.ReplaceAll(verifyAffine, "-> M__Out rev 1", "-> int rev 1")
	prog, texts := admitProg(t, src)
	if ds := VerifyContracts(prog, texts); len(ds) != 0 {
		t.Fatal(ds)
	}
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if err := compile(t.TempDir(), []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatal(err)
	}
	bad := strings.Replace(src, "result.value == 2 * n + 1", "result.value == 2 * n + 2", 1)
	prog, texts = admitProg(t, bad)
	if ds := VerifyContracts(prog, texts); len(ds) == 0 {
		t.Fatal("false scalar contract verified")
	}
}

func TestScalarSuccessBoolContract(t *testing.T) {
	src := `mod m
  provides [m__id, m__go]
  uses []
  emits []

fn m__id(value: bool) -> bool rev 1
  emits []
  ensures
    on Ok r
      r.value == value
  tests
    yes(true) => Ok(true)
    no(false) => Ok(false)
  Ok(value)

fn m__go(value: bool) -> bool rev 1
  emits []
  ensures
    on Ok r
      r.value == value
  tests
    yes(true) => Ok(true)
    no(false) => Ok(false)
  match call m__id(value)
    on Ok r => forward r
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if err := compile(t.TempDir(), []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatal(err)
	}
}

func TestScalarSuccessExternsStillRefused(t *testing.T) {
	for _, typ := range []string{"int", "str", "bool", "dec"} {
		src := "mod m\n  provides [m__host]\n  uses []\n  emits []\n\nextern m__host(value: int) -> " + typ + " rev 1\n  emits []\n"
		_, _, ds := genericProgram(t, map[string]string{"m.can": src})
		if !hasDiag(ds, "error", "externs return a record type") {
			t.Fatal(ds)
		}
	}
}

func TestScalarSuccessLinkedAndLSP(t *testing.T) {
	files := scalarSources(t)
	files = map[string]string{"a.can": files["use.can"], "z.can": files["scalar.can"]}
	prog, _, ds := genericProgram(t, files)
	if len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
	for _, name := range []string{"use__go", "use__relay", "use__recover", "use__report", "use__sugar", "use__classify"} {
		fn := prog.Fns[name]
		for _, row := range fn.Tests {
			ctx, err := freshExecCtx(prog, row.Name, true)
			if err != nil {
				t.Fatal(err)
			}
			env := map[string]*Value{}
			for _, arg := range row.Args {
				v, err := evSmall(arg.V, env, ctx, fn.Name)
				if err != nil {
					t.Fatal(err)
				}
				env[arg.Name] = v
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
				t.Fatalf("%s/%s: %v", name, row.Name, err)
			}
		}
	}
	dir := writeLSPDir(t, files)
	if ds := diagnose(dir, "a.can", files["a.can"]); len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
}

func TestGoldenFnScalar(t *testing.T) {
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/fn-scalar/use.can", "../sketches/fn-scalar/scalar.can"}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"scalar.ts", "use.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(out, name))
		if err != nil {
			t.Fatal(err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/fn-scalar", name))
		if err != nil {
			t.Fatal(err)
		}
		checkGoldenFile(t, name, got, want)
	}
}

func TestScalarSuccessNodeParity(t *testing.T) {
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/fn-scalar/use.can", "../sketches/fn-scalar/scalar.can"}); err != nil {
		t.Fatal(err)
	}
	rewriteSpecifier(t, filepath.Join(out, "use.ts"), "./scalar", "./scalar.ts")
	driver := `import * as s from "./scalar.ts";
import * as u from "./use.ts";
const show = x => JSON.stringify(x,(_,v)=>typeof v === "bigint" ? v.toString()+"n" : v);
const eq = (a,b) => { if(show(a)!==show(b)) throw new Error(show(a)+" != "+show(b)); };
const ok = value => ({$can_kind:"ok",value});
const plus = u.use__bind(3n).cb, minus = u.use__bind(-3n).cb;
for(const n of [0n,7n,-3n,9007199254740993n]) {
 eq(u.use__apply$T$int$T$int(plus,n),ok(n+3n));
 eq(u.use__apply$T$int$T$int(minus,n),ok(n-3n));
 eq(u.use__apply$T$int$T$int(plus,n),ok(n+3n));
 eq(u.use__go(3n,n),ok(n+3n));
 eq(u.use__apply$T$int$T$bool(s.scalar__positive,n),ok(n>0n));
 const result = n < 0n ? {$can_kind:"scalar.negative",value:n} : ok(n);
 eq(u.use__relay(n),result);
 eq(u.use__sugar(n),result);
 eq(u.use__checked(s.scalar__checked,n),result);
 eq(u.use__recover(s.scalar__checked,n),ok(n<0n ? n+3n : n));
 eq(u.use__report(n),{$can_kind:"ok",positive:n>0n});
 eq(u.use__classify(s.scalar__positive,n),ok(n>0n ? "positive" : "not positive"));
}
for(const text of ["","hello","λ 😀"]) eq(u.use__apply$T$str$T$str(s.scalar__identity$T$str,text),ok(text));
for(const b of [false,true]) eq(u.use__apply$T$bool$T$bool(s.scalar__identity$T$bool,b),ok(b));
for(const [d,want] of [["0.125","0.25"],["-2.5","-5.0"]]) {
 eq(u.use__apply$T$dec$T$dec(s.scalar__twice,d),ok(want));
 eq(s.scalar__identity$T$dec(d),ok(d));
}
console.log("SCALAR_SUCCESS_PARITY_OK");
`
	if err := os.WriteFile(filepath.Join(out, "driver.ts"), []byte(driver), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", "driver.ts")
	cmd.Dir = out
	raw, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(raw), "SCALAR_SUCCESS_PARITY_OK") {
		t.Fatalf("%v\n%s", err, raw)
	}
}
