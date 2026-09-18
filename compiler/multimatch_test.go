package main

// Multi-scrutinee value matches (docs/a28): parse gates, product-space
// exhaustiveness, evaluation, and TS emit — plus the native-match
// refactor (MatchKind families, one 1..N value path per phase).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func multiModule(t *testing.T, text string) (*Module, *Program) {
	t.Helper()
	m, err := parseModuleText("m.can", text)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	prog, _ := buildWorld(m, []*Module{m}, map[string]string{"m.can": text})
	return m, prog
}

func multiExhaustive(t *testing.T, text string) []error {
	t.Helper()
	m, prog := multiModule(t, text)
	return verifyExhaustiveAll([]*Module{m}, prog)
}

func errMsgs(errs []error) []string {
	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		msgs = append(msgs, e.Error())
	}
	return msgs
}

func wantErr(t *testing.T, errs []error, sub string) {
	t.Helper()
	for _, m := range errMsgs(errs) {
		if strings.Contains(m, sub) {
			return
		}
	}
	t.Fatalf("expected error containing %q, got %v", sub, errMsgs(errs))
}

func wantClean(t *testing.T, errs []error) {
	t.Helper()
	if len(errs) > 0 {
		t.Fatalf("expected no errors, got %v", errMsgs(errs))
	}
}

const multiHead = `mod m
  provides [m__f]
  uses []
  emits []

fn m__f(x: bool, y: bool) -> bool rev 1
  emits []
  tests
    t(x = true, y = true) => true
=
`

func TestMultiFullTableClean(t *testing.T) {
	text := multiHead + `  match x, y
    true, true => true
    true, false => false
    false, true => false
    false, false => true
`
	wantClean(t, multiExhaustive(t, text))
}

func TestMultiWildcardTableClean(t *testing.T) {
	// (true,_) covers two cells; the residual arms complete the table.
	text := multiHead + `  match x, y
    true, _ => true
    false, true => true
    false, false => false
`
	wantClean(t, multiExhaustive(t, text))
}

func TestMultiMissingTuple(t *testing.T) {
	text := multiHead + `  match x, y
    true, true => true
    true, false => false
    false, true => false
`
	errs := multiExhaustive(t, text)
	wantErr(t, errs, "non-exhaustive match, missing (false, false)")
}

func TestMultiStrNeedsWild(t *testing.T) {
	// The open string remainder in slot 1 is uncovered: the _-coverage
	// rule fires, not a finite missing-cell list.
	errs := multiExhaustive(t, multiStrText())
	wantErr(t, errs, "value match without _ is not provably exhaustive (slot 1")
}

func multiStrText() string {
	return `mod m
  provides [m__g]
  uses []
  emits []

fn m__g(s: str, b: bool) -> bool rev 1
  emits []
  tests
    t(s = "a", b = true) => true
=
  match s, b
    "a", true => true
    "a", false => false
`
}

func TestMultiArity(t *testing.T) {
	text := multiHead + `  match x, y
    true => true
    false, _ => false
`
	errs := multiExhaustive(t, text)
	wantErr(t, errs, "match arm has 1 patterns; this match has 2 scrutinees")
	for _, e := range errs {
		if strings.Contains(e.Error(), "match arm has 1 patterns") {
			if got := proofDiag(text, e).Code; got != CodeBadArmKind {
				t.Fatalf("arity error code = %s, want %s", got, CodeBadArmKind)
			}
		}
	}
}

func TestMultiMissingCode(t *testing.T) {
	text := multiHead + `  match x, y
    true, _ => true
    false, false => false
`
	errs := multiExhaustive(t, text)
	wantErr(t, errs, "non-exhaustive match, missing (false, true)")
	for _, e := range errs {
		if strings.Contains(e.Error(), "non-exhaustive match, missing ") {
			if got := proofDiag(text, e).Code; got != CodeMissingArm {
				t.Fatalf("missing-tuple code = %s, want %s", got, CodeMissingArm)
			}
		}
	}
}

func TestMultiCallRejected(t *testing.T) {
	text := multiHead + `  match call m__f(x, y), y
    true, true => true
    true, false => false
    false, true => false
    false, false => true
`
	errs := multiExhaustive(t, text)
	wantErr(t, errs, "multi-scrutinee match over call m__f is not supported")
}

func TestMultiGivenRejected(t *testing.T) {
	text := multiHead + `  match x, y
    given
      t => true
    true, _ => true
    false, _ => false
`
	m, prog := multiModule(t, text)
	var fn *FnDecl
	for _, d := range m.Decls {
		if f, ok := d.(*FnDecl); ok {
			fn = f
		}
	}
	found := false
	for _, d := range checkGiven(fn, prog, text) {
		if strings.Contains(d.Msg, "multi-scrutinee match takes no given table") {
			found = true
			if d.Code != CodeGivenOnLocal {
				t.Fatalf("given-on-multi code = %s, want %s", d.Code, CodeGivenOnLocal)
			}
		}
	}
	if !found {
		t.Fatal("expected given-on-multi rejection")
	}
}

func TestMultiMixedSlot(t *testing.T) {
	text := multiHead + `  match x, y
    true, _ => true
    "a", _ => false
`
	errs := multiExhaustive(t, text)
	wantErr(t, errs, "bool match must be exactly true+false (slot 1 mixes bool and string patterns)")
}

func TestMultiEmptySlotParse(t *testing.T) {
	text := multiHead + `  match x, y
    true,, true => true
    false, _ => false
`
	if _, err := parseModuleText("m.can", text); err == nil {
		t.Fatal("expected parse error for empty pattern slot")
	} else if !strings.Contains(err.Error(), "empty slot in match list") {
		t.Fatalf("wrong parse error: %v", err)
	}
}

func TestMultiTrailingCommaParse(t *testing.T) {
	text := multiHead + `  match x,
    true => true
    false => false
`
	if _, err := parseModuleText("m.can", text); err == nil {
		t.Fatal("expected parse error for trailing comma in match list")
	} else if !strings.Contains(err.Error(), "empty slot in match list") {
		t.Fatalf("wrong parse error: %v", err)
	}
}

const multiEvalText = `mod m
  provides [m__multi, m__nested, m__words, M__Out, M__Word]
  uses []
  emits []

type M__Out rev 1 (
  v: bool
)

type M__Word rev 1 (
  w: str
)

fn m__multi(a: bool, b: bool) -> M__Out rev 1
  emits []
  tests
    tt(a = true, b = true) => Ok(v = true)
    tf(a = true, b = false) => Ok(v = false)
    ft(a = false, b = true) => Ok(v = false)
    ff(a = false, b = false) => Ok(v = true)
=
  match a, b
    true, true => Ok(v = true)
    true, false => Ok(v = false)
    false, true => Ok(v = false)
    false, false => Ok(v = true)

fn m__nested(a: bool, b: bool) -> M__Out rev 1
  emits []
  tests
    tt(a = true, b = true) => Ok(v = true)
    tf(a = true, b = false) => Ok(v = false)
    ft(a = false, b = true) => Ok(v = false)
    ff(a = false, b = false) => Ok(v = true)
=
  match a
    true => match b
      true => Ok(v = true)
      false => Ok(v = false)
    false => match b
      true => Ok(v = false)
      false => Ok(v = true)

fn m__words(s: str, b: bool) -> M__Word rev 1
  emits []
  tests
    admin(s = "admin", b = false) => Ok(w = "root")
    user(s = "bob", b = true) => Ok(w = "user")
    guest(s = "bob", b = false) => Ok(w = "guest")
=
  match s, b
    "admin", _ => Ok(w = "root")
    _, true => Ok(w = "user")
    _, false => Ok(w = "guest")
`

func multiFn(t *testing.T, mods []*Module, name string) *FnDecl {
	t.Helper()
	for _, m := range mods {
		for _, d := range m.Decls {
			if fn, ok := d.(*FnDecl); ok && fn.Name == name {
				return fn
			}
		}
	}
	t.Fatalf("fn %s not found", name)
	return nil
}

func TestMultiEvalRunsTables(t *testing.T) {
	// Every decision-table row runs green: a wrong arm pick would fail
	// expectation comparison, so green tables prove selection.
	m, prog := multiModule(t, multiEvalText)
	for _, name := range []string{"m__multi", "m__nested", "m__words"} {
		fn := multiFn(t, []*Module{m}, name)
		for _, tc := range fn.Tests {
			if err := runTest(fn, tc, prog, map[*Node]map[int]bool{}); err != nil {
				t.Fatalf("%s/%s: %v", name, tc.Name, err)
			}
		}
	}
}

func TestMultiNestedEquivalence(t *testing.T) {
	// The tuple table and its nested-single spelling agree on the full
	// input product: multi-match is sugar over nesting, not new meaning.
	m, prog := multiModule(t, multiEvalText)
	multi := multiFn(t, []*Module{m}, "m__multi")
	nested := multiFn(t, []*Module{m}, "m__nested")
	for _, ab := range [][2]bool{{true, true}, {true, false}, {false, true}, {false, false}} {
		args := []Arg{
			{Name: "a", V: &Small{Kind: "bool", B: ab[0]}},
			{Name: "b", V: &Small{Kind: "bool", B: ab[1]}},
		}
		syn := Test{Name: "syn", Args: args, Expected: &Small{Kind: "bool", B: true}}
		gotM, _, err := runTestValue(multi, syn, prog, map[*Node]map[int]bool{})
		if err != nil {
			t.Fatalf("multi %v: %v", ab, err)
		}
		gotN, _, err := runTestValue(nested, syn, prog, map[*Node]map[int]bool{})
		if err != nil {
			t.Fatalf("nested %v: %v", ab, err)
		}
		if normalizeValue(gotM) != normalizeValue(gotN) {
			t.Fatalf("divergence on %v: multi=%s nested=%s", ab, normalizeValue(gotM), normalizeValue(gotN))
		}
	}
}

func TestMultiDiagnoseClean(t *testing.T) {
	// The full editor pipeline (checks, proof, test runs, test-per-arm
	// coverage) reports nothing on exhaustive, fully-hit tables.
	dir := writeLSPDir(t, map[string]string{"m.can": multiEvalText})
	for _, d := range diagnose(dir, "m.can", multiEvalText) {
		if d.Sev == "error" {
			t.Fatalf("unexpected diagnostic: %v", d)
		}
	}
}

const multiEmitText = `mod m
  provides [m__w, M__Word]
  uses []
  emits []

type M__Word rev 1 (
  w: str
)

fn m__w(s: str, b: bool) -> M__Word rev 1
  emits []
  tests
    admin(s = "admin", b = false) => Ok(w = "root")
    user(s = "bob", b = true) => Ok(w = "user")
    guest(s = "bob", b = false) => Ok(w = "guest")
=
  match s, b
    "admin", _ => Ok(w = "root")
    _, true => Ok(w = "user")
    _, false => Ok(w = "guest")
`

const multiEmitBoolText = `mod m
  provides [m__q, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  v: bool
)

fn m__q(a: bool, b: bool) -> M__Out rev 1
  emits []
  tests
    tt(a = true, b = true) => Ok(v = true)
    tf(a = true, b = false) => Ok(v = false)
    ft(a = false, b = true) => Ok(v = false)
    ff(a = false, b = false) => Ok(v = true)
=
  match a, b
    true, true => Ok(v = true)
    true, false => Ok(v = false)
    false, true => Ok(v = false)
    false, false => Ok(v = true)
`

func compileEmit(t *testing.T, text string) string {
	t.Helper()
	dir := writeLSPDir(t, map[string]string{"m.can": text})
	out := t.TempDir()
	if err := compile(out, []string{dir + "/m.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func TestMultiEmitChain(t *testing.T) {
	// One temporary per scrutinee, evaluated once; the implied check
	// drops out of the second arm; the proved-total tail is a bare else.
	ts := compileEmit(t, multiEmitText)
	for _, want := range []string{
		"const $can_m1=s;",
		"const $can_m2=b;",
		`if($can_m1==="admin"){`,
		"else if($can_m2){",
		"else{",
	} {
		if !strings.Contains(ts, want) {
			t.Fatalf("emit missing %q:\n%s", want, ts)
		}
	}
	if strings.Contains(ts, "!($can_m1)") {
		t.Fatalf("emit retests the implied scrutinee:\n%s", ts)
	}
}

func TestMultiEmitBoolTotal(t *testing.T) {
	// Full conjunctions where nothing is implied, and the last cell
	// still lowers to a bare else by the residual proof.
	ts := compileEmit(t, multiEmitBoolText)
	for _, want := range []string{
		"const $can_m1=a;",
		"const $can_m2=b;",
		"if($can_m1&&$can_m2){",
		"else{",
	} {
		if !strings.Contains(ts, want) {
			t.Fatalf("emit missing %q:\n%s", want, ts)
		}
	}
}

func TestSingleCommaStrStillSingle(t *testing.T) {
	// A comma inside a string is not a separator: one scrutinee, one
	// pattern, and the legacy single-match rules apply verbatim.
	text := `mod m
  provides [m__h]
  uses []
  emits []

fn m__h(s: str) -> bool rev 1
  emits []
  tests
    t(s = "a,b") => true
=
  match s
    "a,b" => true
    _ => false
`
	wantClean(t, multiExhaustive(t, text))
}
