package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func successSources(t *testing.T) map[string]string {
	t.Helper()
	files := map[string]string{}
	for _, name := range []string{"success.can", "use.can"} {
		b, err := os.ReadFile(filepath.Join("../sketches/success-values", name))
		if err != nil {
			t.Fatal(err)
		}
		files[name] = string(b)
	}
	return files
}

func TestWholeSuccessRejects(t *testing.T) {
	for _, tc := range []struct{ name, file, old, new, want string }{
		{"legacy_record_not_splat", "success.can", "  Ok<T>(value)", "  Ok(value)", "got Success__Pair, want int"},
		{"wrong_nominal_record", "success.can", "  Ok<T>(value)", "  Ok<T>(Success__Value(7))", "got Success__Value, want Success__Pair"},
		{"wrong_annotation", "success.can", "  Ok<T>(value)", "  Ok<int>(value)", "requires a matching success boundary"},
		{"unknown_type", "success.can", "  Ok<T>(value)", "  Ok<Missing>(value)", "known success type"},
		{"missing_value", "success.can", "  Ok<T>(value)", "  Ok<T>()", "exactly one whole value"},
		{"extra_value", "success.can", "  Ok<T>(value)", "  Ok<T>(value, value)", "exactly one whole value"},
		{"unknown_field", "success.can", "  Ok<T>(value)", "  Ok<T>(left = value)", "exactly one whole value"},
		{"nested_outcome", "success.can", "  Ok<T>(value)", "  Ok<T>(Ok(value))", "not an outcome"},
		{"nested_typed_outcome", "success.can", "  Ok<T>(value)", "  Ok<T>(Ok<T>(value))", "matching success boundary"},
		{"error_as_value", "success.can", "  Ok<T>(value)", "  Ok<T>(option.absent())", "not an outcome"},
		{"pattern_type", "use.can", "on Ok<T> result => Ok<T>(result)", "on Ok<int> result => Ok<T>(result)", "exact success type Success__Pair"},
		{"pattern_wrapper_not_equivalent", "use.can", "on Ok<T> result => Ok<T>(result)", "on Ok<Success__Value> result => Ok<T>(result)", "exact success type int"},
		{"scalar_no_implicit_projection", "use.can", "on Ok<T> result => Ok<T>(result)", "on Ok<T> result => Ok<T>(result.value)", "no field value on result"},
		{"typed_outcome_as_call_arg", "use.can", "use__apply<T>(fnref success__id<T>(), value)", "use__apply<T>(fnref success__id<T>(), Ok<T>(value))", "matching success boundary"},
		{"missing_given", "use.can", "      integer => [exchange args (value = 7) outcome Ok<Success__Box<int>>(Success__Box<int>(7))]\n", "", "no script for call"},
		{"wrong_script", "use.can", "outcome Ok<Success__Box<int>>(Success__Box<int>(7))", "outcome Ok<Success__Box<int>>(Success__Box<int>(8))", "contradicts"},
		{"missing_invoke_witness", "success.can", "    int_some<T=int, U=int>(Success__Some<int>(7), fnref success__id<int>()) => Ok(Success__Some<int>(7))\n", "", "no test takes"},
		{"foreign_seal", "use.can", "Ok<Success__Pair>(pair)", "Ok<Success__Pair>(Success__Pair(7, seal Success__Token(\"x\")))", "bodies seal only"},
		{"missing_pin", "use.can", "success__id@1, ", "", "not in uses"},
		{"binder_invocation_not_added", "use.can", "match call use__apply<T>(cb, value)", "match invoke cb with value", "want a function value"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			files := successSources(t)
			if !strings.Contains(files[tc.file], tc.old) {
				t.Fatal("missing mutation")
			}
			files[tc.file] = strings.Replace(files[tc.file], tc.old, tc.new, 1)
			_, _, ds := genericProgram(t, files)
			var messages []string
			for _, d := range ds {
				messages = append(messages, d.Msg)
			}
			if !strings.Contains(strings.Join(messages, "\n"), tc.want) {
				t.Fatalf("want %q: %s", tc.want, strings.Join(messages, "\n"))
			}
		})
	}
}

func TestWholeSuccessParse(t *testing.T) {
	for _, src := range []string{"Ok<int>()", "Ok<Fn<int, int, []>>(cb)", "Ok<Seq<int>>(items)", "Ok<M__Box<T>>(box)"} {
		s, err := parseSmall(src)
		if err != nil || s.Ctor != "Ok" || len(s.TypeArgs) != 1 {
			t.Fatalf("%s: %+v %v", src, s, err)
		}
	}
	for _, src := range []string{"Ok<>(x)", "Ok<int, bool>(x)"} {
		if _, err := parseSmall(src); err == nil {
			t.Fatal(src)
		}
	}
	for _, src := range []string{"Ok<Fn<int, int, []>> cb", "Ok<M__Box<T>> box"} {
		if _, err := parsePattern(src); err != nil {
			t.Fatal(err)
		}
	}
}

func TestWholeSuccessNominalRecords(t *testing.T) {
	src := strings.Replace(wholeSuccessProof, "fn m__pair(value:", "type M__Other rev 1 (\n  left: int\n  right: bool\n)\nfn m__pair(value:", 1)
	for _, tc := range []struct{ old, new, want string }{
		{"Ok<M__Pair>(value)", "Ok<M__Pair>(M__Other(value.left, value.right))", "got M__Other, want M__Pair"},
		{"on Ok<M__Pair> result", "on Ok<M__Other> result", "exact success type M__Pair"},
	} {
		_, _, ds := genericProgram(t, map[string]string{"m.can": strings.Replace(src, tc.old, tc.new, 1)})
		wantGenericDiag(t, ds, CodeTypeMismatch, tc.want)
	}
}

func TestWholeSuccessGivenRouting(t *testing.T) {
	fn := &FnDecl{Tests: []Test{{Name: "keep"}}, Body: &Node{Given: map[string]*Small{"keep": {}, "other": {}, "typo": {}}, ChainSteps: []ChainStep{{Given: map[string]*Small{"keep": {}, "other": {}, "typo": {}}}}}}
	routeGivenRows(fn, []Test{{Name: "keep"}, {Name: "other"}})
	for _, g := range []map[string]*Small{fn.Body.Given, fn.Body.ChainSteps[0].Given} {
		if len(g) != 2 || g["keep"] == nil || g["typo"] == nil {
			t.Fatal(g)
		}
	}
}

func TestWholeSuccessLinkedAndLSP(t *testing.T) {
	files := successSources(t)
	files = map[string]string{"a.can": files["use.can"], "z.can": files["success.can"]}
	prog, _, ds := genericProgram(t, files)
	if len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
	for name, fn := range prog.Fns {
		if name == "use__host_result" {
			continue
		} // Hosts are scripted, never linked-pure.
		for _, row := range fn.Tests {
			ctx, err := freshExecCtx(prog, row.Name, true)
			if err != nil {
				t.Fatal(err)
			}
			env := map[string]*Value{}
			for _, arg := range row.Args {
				v, err := evSmall(arg.V, env, ctx, name)
				if err != nil {
					t.Fatal(err)
				}
				env[arg.Name] = v
			}
			got, err := evNode(fn.Body, env, ctx, name)
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			want, err := evSmall(row.Expected, env, ctx, name)
			if err != nil {
				t.Fatal(err)
			}
			eq, _, err := expectEq(want, got)
			if err != nil || !eq {
				t.Fatalf("%s/%s: %v %v != %v", name, row.Name, err, got, want)
			}
		}
	}
	dir := writeLSPDir(t, files)
	if ds := diagnose(dir, "a.can", files["a.can"]); len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
}

const wholeSuccessProof = `mod m
  provides [M__Pair, m__pair, m__go, m__bool, m__flag]
  uses []
  emits []
type M__Pair rev 1 (
  left: int
  right: bool
)
fn m__pair(value: M__Pair) -> M__Pair rev 1
  emits []
  ensures
    on Ok r
      r.left == value.left
      r.right == value.right
  tests
    yes(M__Pair(7, true)) => Ok(7, true)
  Ok<M__Pair>(value)
fn m__go(value: M__Pair) -> M__Pair rev 1
  emits []
  ensures
    on Ok r
      r.left == value.left
      r.right == value.right
  tests
    yes(M__Pair(7, true)) => Ok(7, true)
  match call m__pair(value)
    on Ok<M__Pair> result => Ok<M__Pair>(result)
fn m__bool(value: bool) -> bool rev 1
  emits []
  ensures
    on Ok r
      r.value == value
  tests
    yes(true) => Ok(true)
    no(false) => Ok(false)
  Ok<bool>(value)
fn m__flag(value: bool) -> bool rev 1
  emits []
  ensures
    on Ok r
      r.value == value
  tests
    yes(true) => Ok(true)
    no(false) => Ok(false)
  match call m__bool(value)
    on Ok<bool> result => Ok<bool>(result)
`

func TestWholeSuccessContracts(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": wholeSuccessProof})
	if err := compile(t.TempDir(), []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatal(err)
	}
	for _, src := range []string{
		strings.Replace(wholeSuccessProof, "r.left == value.left", "r.left == value.left + 1", 1),
		strings.ReplaceAll(wholeSuccessProof, "r.value == value", "r.value == not value"),
	} {
		prog, texts := admitProg(t, src)
		if ds := VerifyContracts(prog, texts); len(ds) == 0 {
			t.Fatal("false contract verified")
		}
	}
	if ds := diagnose(dir, "m.can", wholeSuccessProof); len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
}

func TestWholeSuccessIntContracts(t *testing.T) {
	src := strings.ReplaceAll(verifyAffine, "-> M__Out rev 1", "-> int rev 1")
	src = strings.ReplaceAll(src, "Ok(", "Ok<int>(")
	src = strings.Replace(src, "on Ok r => Ok<int>(r.value)", "on Ok<int> r => Ok<int>(r)", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if err := compile(t.TempDir(), []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatal(err)
	}
	bad := strings.Replace(src, "result.value == 2 * n + 1", "result.value == 2 * n + 2", 1)
	prog, texts := admitProg(t, bad)
	if ds := VerifyContracts(prog, texts); len(ds) == 0 {
		t.Fatal("false int contract verified")
	}
}

func TestWholeSuccessKernel(t *testing.T) {
	src := strings.Replace(bytesEncodeFull(), "on Ok r", "on Ok<Bytes__Value> r", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "const r = { value:") {
		t.Fatal("specialized encoder failed to project record binder")
	}
	bad := strings.Replace(src, "Ok<Bytes__Value> r", "Ok<Bytes> r", 1)
	_, _, ds := genericProgram(t, map[string]string{"m.can": bad})
	wantGenericDiag(t, ds, CodeTypeMismatch, "exact success type Bytes__Value")
}

func compileSuccess(t *testing.T) string {
	t.Helper()
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/success-values/use.can", "../sketches/success-values/success.can"}); err != nil {
		t.Fatal(err)
	}
	return out
}
func TestGoldenSuccessValues(t *testing.T) {
	out := compileSuccess(t)
	for _, name := range []string{"success.ts", "use.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(out, name))
		if err != nil {
			t.Fatal(err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/success-values", name))
		if err != nil {
			t.Fatal(err)
		}
		checkGoldenFile(t, name, got, want)
	}
}
func TestWholeSuccessNodeParity(t *testing.T) {
	out := compileSuccess(t)
	host, err := os.ReadFile("../sketches/success-values/use.externs.ts")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(out, "use.externs.ts"), host, 0644); err != nil {
		t.Fatal(err)
	}
	rewriteSpecifier(t, filepath.Join(out, "use.ts"), "./success", "./success.ts")
	rewriteSpecifier(t, filepath.Join(out, "use.ts"), "./use.externs", "./use.externs.ts")
	driver := `import * as s from "./success.ts";
import * as u from "./use.ts";
const show=x=>JSON.stringify(x,(_,v)=>typeof v==="bigint"?v.toString()+"n":v instanceof Uint8Array?Array.from(v):v);
const eq=(a,b)=>{if(show(a)!==show(b))throw new Error(show(a)+" != "+show(b));};
const ok=value=>({$can_kind:"ok",value});
const some=value=>({$can_kind:"Success__Some$T$Success__Pair",value});
const none={$can_kind:"Success__None$T$Success__Pair"};
let reads=0;
eq(s.success__id$T$Success__Pair({get left(){reads++;return 7n;},get right(){reads++;return true;}}),{$can_kind:"ok",left:7n,right:true});
eq(reads,2);
for(const left of [0n,-3n,9007199254740993n]) for(const right of [false,true]) {
 const pair={left,right}, result={$can_kind:"ok",left,right};
 eq(s.success__id$T$Success__Pair(pair),result);
 eq(u.use__apply$T$Success__Pair(s.success__id$T$Success__Pair,pair),result);
 eq(u.use__relay$T$Success__Pair(pair),result);
 eq(u.use__factory$T$Success__Pair().value(pair),result);
 eq(u.use__factory_apply$T$Success__Pair(pair),result);
 eq(u.use__box$T$Success__Pair(pair),{$can_kind:"ok",item:pair});
 eq(s.success__map$T$Success__Pair$T$Success__Pair(some(pair),s.success__id$T$Success__Pair),ok(some(pair)));
 eq(s.success__value_or$T$Success__Pair(none,pair),result);
 eq(u.use__recover(some(pair)),result);
 eq(s.success__mapped({$can_kind:"Success__Some$T$int",value:left}),ok(some({left,right:left>0n})));
 eq(u.use__relay$T$int(left),ok(left));
 eq(u.use__factory_apply$T$int(left),ok(left));
}
eq(s.success__map$T$Success__Pair$T$Success__Pair(none,()=>{throw new Error("must not invoke for None")}),ok(none));
eq(s.success__mapped({$can_kind:"Success__None$T$int"}),ok(none));
eq(u.use__recover(none),{$can_kind:"ok",left:0n,right:false});
eq(u.use__apply$T$Success__Empty(s.success__id$T$Success__Empty,{}),{$can_kind:"ok"});
eq(u.use__apply$T$Success__Value(s.success__id$T$Success__Value,{value:9n}),ok(9n));
eq(u.use__apply$T$Bytes(s.success__id$T$Bytes,new Uint8Array([0,255])),ok(new Uint8Array([0,255])));
eq(u.use__apply$T$Seq$L$int$G$(s.success__id$T$Seq$L$int$G$,[1n,2n]),ok([1n,2n]));
eq(u.use__host_result(true),{$can_kind:"ok",left:7n,right:true});
eq(u.use__host_result(false),{$can_kind:"option.absent"});
console.log("WHOLE_SUCCESS_PARITY_OK");`
	if err := os.WriteFile(filepath.Join(out, "driver.ts"), []byte(driver), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", "driver.ts")
	cmd.Dir = out
	raw, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(raw), "WHOLE_SUCCESS_PARITY_OK") {
		t.Fatalf("%v\n%s", err, raw)
	}
}
