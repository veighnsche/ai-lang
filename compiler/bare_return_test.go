package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func bareSources(t *testing.T) map[string]string {
	t.Helper()
	files := map[string]string{}
	for _, n := range []string{"data.can", "bridge.can", "use.can"} {
		b, err := os.ReadFile(filepath.Join("../sketches/bare-returns", n))
		if err != nil {
			t.Fatal(err)
		}
		files[n] = string(b)
	}
	return files
}

func TestBareReturnShapes(t *testing.T) {
	p, _, ds := genericProgram(t, bareSources(t))
	if len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
	e := &emitter{recs: recordShapes(p.Modules), variants: variantShapes(p.Modules), brands: p.Brands, errTypes: errorShapes(p.Modules)}
	for _, name := range []string{"int", "str", "bool", "dec", "token", "bytes", "rows", "choice"} {
		host, err := externUnion(p.Externs["bridge__"+name], p)
		if err != nil {
			t.Fatal(err)
		}
		source, err := fnResultUnion(p.Fns["use__"+name], p)
		if err != nil || source != host {
			t.Fatalf("%s: %s / %s: %v", name, source, host, err)
		}
	}
	for _, name := range []string{"bytes", "token", "items$T$int", "items$T$Data__Row"} {
		fn := p.Fns["use__apply_"+name]
		inv, err := e.sigResultUnion(fn.Body.invokeSig)
		if err != nil {
			t.Fatal(err)
		}
		source, err := fnResultUnion(fn, p)
		if err != nil {
			t.Fatal(err)
		}
		head, err := tsTypeB(fn.Params[0][1], p.Brands, e.recs, e.variants, e.errTypes)
		if err != nil || inv != source || !strings.HasSuffix(head, " => "+source) {
			t.Fatalf("%s / %s / %s: %v", inv, source, head, err)
		}
	}
	// A source factory result is not an admissible callback success.
	ck := newTycker(p, "", "use__factory")
	if _, ok := admissibleRefSig(ck, p, &Small{Kind: "fnref", Fname: "data__factory"}); ok {
		t.Fatal("factory became a higher-order callback")
	}
}

func TestBareReturnRejects(t *testing.T) {
	cases := []struct{ name, file, old, new, want string }{
		{"bytes_raw", "data.can", "  Ok(value)\n", "  value\n", "value return must use Ok"},
		{"bytes_wrong", "data.can", "  Ok(value)\n", "  Ok(Seq<int>[])\n", "want Bytes"},
		{"bytes_nested_outcome", "data.can", "  Ok(value)\n", "  Ok(Ok(value))\n", "needs a value of type Bytes"},
		{"bytes_missing", "data.can", "  Ok(value)\n", "  Ok()\n", "Ok is missing field value"},
		{"bytes_duplicate", "data.can", "  Ok(value)\n", "  Ok(value = value, value = value)\n", "Ok repeats field value"},
		{"bytes_extra", "data.can", "  Ok(value)\n", "  Ok(value = value, other = 1)\n", "Ok has no field other"},
		{"seq_nested_outcome", "data.can", "  Ok(values)\n", "  Ok(Seq<int>[Ok(7)])\n", "Seq literal needs a value of type int"},
		{"seq_element", "data.can", "  Ok(values)\n", "  Ok(Seq<str>[\"x\"])\n", "want Seq<int>"},
		{"brand_unsealed", "data.can", "  Ok(seal Data__Token(value))", "  Ok(value)", "want Data__Token"},
		{"factory_raw", "data.can", "  Ok(fnref data__add(delta = delta))", "  fnref data__add(delta = delta)", "value return must use Ok"},
		{"factory_mismatch", "data.can", "  Ok(fnref data__add(delta = delta))", "  Ok(fnref data__mint())", "want Fn<int, int, []>"},
		{"factory_capture", "data.can", "  Ok(fnref data__add(delta = delta))", "  Ok(fnref data__add(delta = delta + 1))", "computed"},
		{"factory_as_callback", "use.can", "fnref data__bytes()", "fnref data__factory()", "success carriers are data-only"},
		{"extern_as_callback", "use.can", "fnref data__bytes()", "fnref bridge__bytes()", "reference to extern"},
		{"missing_pin", "use.can", "bridge__bytes@1, ", "", "not in uses"},
		{"wrong_pin", "use.can", "bridge__bytes@1", "bridge__bytes@2", "pins rev 2"},
		{"missing_script", "use.can", "    given\n      full => [exchange args () outcome Ok(Bytes(Seq<int>[0, 255]))]\n", "", "has no given table"},
		{"bytes_script_type", "use.can", "outcome Ok(Bytes(Seq<int>[0, 255]))", "outcome Ok(7)", "want Bytes"},
		{"brand_script_type", "use.can", "outcome Ok(seal Data__Token(\"a\"))", "outcome Ok(\"a\")", "want Data__Token"},
		{"seq_script_type", "use.can", "outcome Ok(Seq<Data__Row>[Data__Row(7)])", "outcome Ok(Seq<int>[7])", "want Seq<Data__Row>"},
		{"variant_script_type", "use.can", "outcome Ok(Data__Some<int>(7))", "outcome Ok(7)", "want Data__Choice"},
		{"bool_script_type", "use.can", "outcome Ok(false)", "outcome Ok(0)", "want bool"},
		{"dec_script_type", "use.can", "outcome Ok(d\"0.125\")", "outcome Ok(0)", "want dec"},
		{"int_script_missing", "use.can", "outcome Ok(9007199254740993)", "outcome Ok()", "Ok is missing field value"},
		{"factory_script_type", "use.can", "outcome Ok(fnref data__add(delta = 3))", "outcome Ok(3)", "want Fn<int, int, []>"},
		{"foreign_seal", "use.can", "    on Ok r => forward r\n\nfn use__bytes", "    on Ok r => Ok(seal Data__Token(\"forged\"))\n\nfn use__bytes", "bodies seal only their own module's brands"},
		{"bytes_forward_mismatch", "use.can", "fn use__bytes() -> Bytes", "fn use__bytes() -> Seq<int>", "forward needs the same value return type"},
		{"brand_forward_mismatch", "use.can", "fn use__token() -> Data__Token", "fn use__token() -> str", "forward needs the same scalar return type"},
		{"factory_forward_mismatch", "use.can", "fn use__factory_relay(delta: int) -> Fn<int,int,[]>", "fn use__factory_relay(delta: int) -> Fn<int,bool,[]>", "forward needs the same value return type"},
		{"whole_binder", "use.can", "use__apply_int(f.value, value)", "use__apply_int(f, value)", "select f.value"},
		{"missing_outcome", "use.can", "    on bridge.denied e => forward e\n", "", "missing bridge.denied"},
		{"unwitnessed_outcome", "use.can", "    no(false) => bridge.denied()\n", "", "no test takes on bridge.denied"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			files := bareSources(t)
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

func TestBareReturnNominalBrands(t *testing.T) {
	files := bareSources(t)
	files["data.can"] += "\nbrand Data__Other is str rev 1\n"
	files["data.can"] = strings.Replace(files["data.can"], "  Ok(seal Data__Token(value))", "  Ok(seal Data__Other(value))", 1)
	_, _, ds := genericProgram(t, files)
	if !hasDiag(ds, "error", "got Data__Other, want Data__Token") {
		t.Fatal(ds)
	}
	files = bareSources(t)
	files["use.can"] += "\nbrand Use__Other is str rev 1\n"
	files["use.can"] = strings.Replace(files["use.can"], "fn use__token() -> Data__Token", "fn use__token() -> Use__Other", 1)
	_, _, ds = genericProgram(t, files)
	if !hasDiag(ds, "error", "forward needs the same value return type") {
		t.Fatal(ds)
	}
}

func TestBareReturnContainment(t *testing.T) {
	for _, ret := range []string{"Fn<int,int,[]>", "Bad__Holder", "Seq<Bad__Holder>"} {
		files := bareSources(t)
		files["bridge.can"] += "\ntype Bad__Holder rev 1 (\n  cb: Fn<int,int,[]>\n)\n\nextern bridge__bad() -> " + ret + " rev 1\n  emits []\n"
		p, _, ds := genericProgram(t, files)
		if !hasDiag(ds, "error", "extern signatures are data-only") {
			t.Fatal(ds)
		}
		if _, err := externUnion(p.Externs["bridge__bad"], p); err == nil {
			t.Fatal("emitted callable extern", ret)
		}
	}
	p, _, ds := genericProgram(t, bareSources(t))
	if len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
	ck := newTycker(p, "", "use__factory")
	for _, r := range []string{"Seq<Data__Choice$T$int>", "Seq<Missing>", "Fn<int,int,[]>"} {
		head := "Fn<int," + r + ",[]>"
		if ds := ck.fnHeadDiags(head, "test", 1, "Fn"); len(genericErrs(ds)) == 0 {
			t.Fatalf("admitted %s", head)
		}
		if _, err := tsTypeB(head, p.Brands, recordShapes(p.Modules), variantShapes(p.Modules), errorShapes(p.Modules)); err == nil {
			t.Fatalf("emitted %s", head)
		}
	}
}

func TestBareReturnSequenceRestrictions(t *testing.T) {
	for _, ret := range []string{"Seq<Data__Choice<int>>", "Seq<Seq<int>>", "Seq<Bad__Holder>"} {
		files := bareSources(t)
		files["data.can"] += "\ntype Bad__Holder rev 1 (\n  cb: Fn<int,int,[]>\n)\n\nfn data__bad(value: " + ret + ") -> " + ret + " rev 1\n  emits []\n  tests\n    empty(Seq<int>[]) => Ok(Seq<int>[])\n  Ok(value)\n"
		_, _, ds := genericProgram(t, files)
		if len(genericErrs(ds)) == 0 {
			t.Fatalf("admitted %s", ret)
		}
		if ret == "Seq<Bad__Holder>" && !hasDiag(ds, "error", "sequence elements are data-only") {
			t.Fatal(ds)
		}
	}
}

func TestBareReturnLinkedAndLSP(t *testing.T) {
	files := bareSources(t)
	// Consumer first, so executing its factory callback needs prepared providers.
	files = map[string]string{"a.can": files["use.can"], "b.can": files["bridge.can"], "z.can": files["data.can"]}
	p, _, ds := genericProgram(t, files)
	if len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
	fn := p.Fns["use__factory"]
	for _, row := range fn.Tests {
		ctx, err := freshExecCtx(p, row.Name, true)
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
			t.Fatal(row.Name, err)
		}
	}
	// Externs remain outside linked-pure execution even with data-only returns.
	if err := checkLinkedGraph(p, "use__int"); err == nil {
		t.Fatal("trusted host became linked-pure")
	}
	dir := writeLSPDir(t, files)
	if ds := diagnose(dir, "a.can", files["a.can"]); len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
}

func TestBareReturnFnImports(t *testing.T) {
	files := map[string]string{
		"p.can": `mod p
  provides [P__Out, P__Nested, P__ErrorData, p__run, p__checked]
  uses []
  emits [p.bad]

type P__Nested rev 1 (
  number: int
)

type P__Out rev 1 (
  item: P__Nested
)

type P__ErrorData rev 1 (
  number: int
)

error p.bad(detail: P__ErrorData)

fn p__run(value: int) -> P__Out rev 1
  emits []
  tests
    ok(7) => Ok(P__Nested(7))
  Ok(P__Nested(value))

fn p__checked(value: int) -> int rev 1
  emits [p.bad]
  tests
    ok(7) => Ok(7)
    bad(-1) => p.bad(P__ErrorData(-1))
  match value < 0
    true => p.bad(P__ErrorData(value))
    false => Ok(value)
`,
		"m.can": `mod m
  provides [m__make, m__checked]
  uses [p__run@1, p__checked@1, P__Out@1]
  emits []

fn m__make() -> Fn<int, P__Out, []> rev 1
  emits []
  tests
    ok() => Ok(fnref p__run())
  Ok(fnref p__run())

fn m__checked() -> Fn<int, int, [p.bad]> rev 1
  emits []
  tests
    ok() => Ok(fnref p__checked())
  Ok(fnref p__checked())
`,
	}
	dir := writeLSPDir(t, files)
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can"), filepath.Join(dir, "p.can")}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"type P__Nested", "type P__ErrorData"} {
		if !strings.Contains(string(b), want) {
			t.Fatalf("missing %s:\n%s", want, b)
		}
	}
}

func compileBare(t *testing.T) string {
	t.Helper()
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/bare-returns/use.can", "../sketches/bare-returns/bridge.can", "../sketches/bare-returns/data.can"}); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestGoldenBareReturns(t *testing.T) {
	out := compileBare(t)
	for _, n := range []string{"data.ts", "bridge.ts", "use.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(out, n))
		if err != nil {
			t.Fatal(err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/bare-returns", n))
		if err != nil {
			t.Fatal(err)
		}
		checkGoldenFile(t, n, got, want)
	}
}

func TestBareReturnNodeParity(t *testing.T) {
	out := compileBare(t)
	host, err := os.ReadFile("../sketches/bare-returns/bridge.externs.ts")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(out, "bridge.externs.ts"), host, 0644); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"use.ts", "bridge.ts"} {
		rewriteSpecifier(t, filepath.Join(out, n), "./data", "./data.ts")
		rewriteSpecifier(t, filepath.Join(out, n), "./bridge.externs", "./bridge.externs.ts")
	}
	driver := `import * as u from "./use.ts";
import * as d from "./data.ts";
import * as b from "./bridge.ts";
const show = x => JSON.stringify(x,(_,v)=>typeof v === "bigint" ? v.toString()+"n" : v instanceof Uint8Array ? Array.from(v) : v);
const eq = (a,b) => { if(show(a)!==show(b)) throw new Error(show(a)+" != "+show(b)); };
const ok = value => ({$can_kind:"ok",value});
eq(u.use__int(true),ok(9007199254740993n));
eq(u.use__int(false),{$can_kind:"bridge.denied"});
eq(u.use__str(),ok("hello"));
eq(u.use__bool(),ok(false));
eq(u.use__dec(),ok("0.125"));
eq(u.use__token(),ok("a"));
eq(u.use__bytes(),ok(new Uint8Array([0,255])));
eq(b.bridge__local(),u.use__bytes());
eq(u.use__rows(),ok([{number:7n}]));
eq(u.use__choice(true),ok({$can_kind:"Data__Some$T$int",value:7n}));
eq(u.use__choice(false),ok({$can_kind:"Data__None$T$int"}));
for(const bytes of [new Uint8Array(),new Uint8Array([0,255])]) eq(u.use__apply_bytes(d.data__bytes,bytes),ok(bytes));
for(const text of ["", "a", "λ"]) eq(u.use__apply_token(d.data__mint,text),ok(text));
const bytesId = d.data__reference$T$Bytes().value;
eq(u.use__apply_bytes(bytesId,new Uint8Array([0,255])),ok(new Uint8Array([0,255])));
eq(d.data__reference$T$Data__Token().value("a"),ok("a"));
eq(u.use__apply_items$T$int(d.data__items$T$int,[1n,9007199254740993n]),ok([1n,9007199254740993n]));
eq(u.use__apply_items$T$Data__Row(d.data__items$T$Data__Row,[{number:7n}]),ok([{number:7n}]));
const plus = u.use__factory_relay(3n).value, minus = u.use__factory_relay(-3n).value;
for(const n of [7n,9007199254740993n]) {
 eq(u.use__apply_int(plus,n),ok(n+3n));
 eq(u.use__apply_int(minus,n),ok(n-3n));
 eq(u.use__apply_int(plus,n),ok(n+3n));
 eq(u.use__factory(3n,n),ok(n+3n));
}
console.log("BARE_RETURN_PARITY_OK");
`
	if err := os.WriteFile(filepath.Join(out, "driver.ts"), []byte(driver), 0644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", "driver.ts")
	cmd.Dir = out
	raw, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(raw), "BARE_RETURN_PARITY_OK") {
		t.Fatalf("%v\n%s", err, raw)
	}
}
