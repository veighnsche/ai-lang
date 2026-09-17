package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// v74: variant value admission. Cases construct by qualified name,
// carry the parent type outward, evaluate to a tagged carrier that
// is never an error, and emit as TS discriminated unions. Probes
// first: every TestVariantValue* below fails before the slice.

// Shared fixture: a two-case variant, a record payload type, and a
// wrapper record carrying a variant-typed field. BODY/EXPECT swap
// per probe through variantValueMod.

const variantValueDecls = `variant Login__State rev 1 (
  case Anonymous()
  case Authenticated(session: Auth__Session)
)

type Auth__Session rev 1 (
  user_id: str
)

type M__Box rev 1 (
  state: Login__State
)

type M__Out rev 1 (
  box: M__Box
)

`

func variantValueMod(body, expect string) string {
	return `mod m
  provides [m__wrap, M__Out, M__Box, Login__State, Auth__Session]
  uses []
  emits []

` + variantValueDecls + `fn m__wrap(user: str) -> M__Out rev 1
  emits []
  tests
    go(user = "u") => ` + expect + `
=
  ` + body + `
`
}

const variantValueBody = `Ok(box = M__Box(state = Login__Authenticated(session = Auth__Session(user_id = user))))`
const variantValueExpect = `Ok(box = M__Box(state = Login__Authenticated(session = Auth__Session(user_id = "u"))))`

// TestVariantValueClean pins the positive: exact-field qualified
// construction inside an Ok payload, through a variant-typed
// record field, diagnoses clean with emits [] (no error
// accounting sees the case).
func TestVariantValueClean(t *testing.T) {
	src := variantValueMod(variantValueBody, variantValueExpect)
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// TestVariantValueUnknownField pins exact fields: a case has no
// field bogus.
func TestVariantValueUnknownField(t *testing.T) {
	body := `Ok(box = M__Box(state = Login__Authenticated(bogus = Auth__Session(user_id = user))))`
	src := variantValueMod(body, variantValueExpect)
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if !hasDiag(diags, "error", "Login__Authenticated has no field bogus") {
		t.Fatalf("expected unknown-field error, got %v", diags)
	}
}

// TestVariantValueMissingField pins exact fields: the session
// payload is required.
func TestVariantValueMissingField(t *testing.T) {
	body := `Ok(box = M__Box(state = Login__Authenticated()))`
	src := variantValueMod(body, variantValueExpect)
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if !hasDiag(diags, "error", "Login__Authenticated is missing field session") {
		t.Fatalf("expected missing-field error, got %v", diags)
	}
}

// TestVariantValueFieldType pins field typing through the case:
// session wants Auth__Session, not str.
func TestVariantValueFieldType(t *testing.T) {
	body := `Ok(box = M__Box(state = Login__Authenticated(session = "x")))`
	src := variantValueMod(body, variantValueExpect)
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if !hasDiag(diags, "error", "got str, want Auth__Session") {
		t.Fatalf("expected field-type mismatch, got %v", diags)
	}
}

// TestVariantValueNominal pins nominal typing: a same-shape record
// is not its case, and a record is not the parent. The
// variant-typed field rejects Auth__Session with the parent named.
func TestVariantValueNominal(t *testing.T) {
	body := `Ok(box = M__Box(state = Auth__Session(user_id = user)))`
	src := variantValueMod(body, variantValueExpect)
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if !hasDiag(diags, "error", "got Auth__Session, want Login__State") {
		t.Fatalf("expected nominal mismatch, got %v", diags)
	}
}

// TestVariantBareReturnRejected pins the wrapper discipline:
// functions return records, never a bare variant.
func TestVariantBareReturnRejected(t *testing.T) {
	src := `mod m
  provides [m__id, Login__State]
  uses []
  emits []

variant Login__State rev 1 (
  case Anonymous()
)

fn m__id(user: str) -> Login__State rev 1
  emits []
  tests
    go(user = "u") => Ok()
=
  Ok()
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if !hasDiag(diags, "error", "bare-variant returns are unsupported, return a record") {
		t.Fatalf("expected bare-variant return rejection, got %v", diags)
	}
}

// TestVariantBareBodyRejected pins the wrapper discipline at the
// body: a bare case is its parent type, not the return record.
func TestVariantBareBodyRejected(t *testing.T) {
	src := variantValueMod(`Login__Anonymous()`, variantValueExpect)
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if !hasDiag(diags, "error", "got Login__State, want M__Out") {
		t.Fatalf("expected bare-body mismatch, got %v", diags)
	}
}

// TestVariantBareStubRejected locks the stub rule: given tables
// take Ok or errors, never a bare case.
func TestVariantBareStubRejected(t *testing.T) {
	src := `mod m
  provides [m__go, ext__help, M__Out, M__Mid, Login__State]
  uses []
  emits []

variant Login__State rev 1 (
  case Anonymous()
)

type M__Mid rev 1 (
  value: int
)

type M__Out rev 1 (
  value: int
)

extern ext__help(value: int) -> M__Mid rev 1
  emits []

fn m__go(value: int) -> M__Out rev 1
  emits []
  tests
    g(value = 1) => Ok(value = 1)
=
  match call ext__help(value)
    given
      g => [exchange args (value = 1) outcome Login__Anonymous()]
    on Ok v => Ok(value = v.value)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if !hasDiag(diags, "error", "neither Ok nor an error") {
		t.Fatalf("expected bad-stub rejection, got %v", diags)
	}
}

// TestVariantBareExpectRejected locks the expectation rule: tests
// expect Ok(...) or an error kind, never a bare case.
func TestVariantBareExpectRejected(t *testing.T) {
	src := variantValueMod(variantValueBody, `Login__Authenticated(session = Auth__Session(user_id = "u"))`)
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if !hasDiag(diags, "error", "write Ok(...) or an error kind") {
		t.Fatalf("expected bare-expectation rejection, got %v", diags)
	}
}

// TestVariantEqualityRefused pins the == boundary: cases compare
// by matching (v75), never by ==. Like Bytes and Seq, the refusal
// lands at check time so no variant operand sails through to a
// loud emit failure.
func TestVariantEqualityRefused(t *testing.T) {
	src := `mod m
  provides [m__eq, M__Out, Login__State]
  uses []
  emits []

variant Login__State rev 1 (
  case Anonymous()
  case Authenticated(session: str)
)

type M__Out rev 1 (
  value: int
)

fn m__eq(flag: int) -> M__Out rev 1
  emits []
  tests
    go(flag = 1) => Ok(value = 1)
=
  match Login__Anonymous() == Login__Anonymous()
    true => Ok(value = 1)
    false => Ok(value = 0)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if !hasDiag(diags, "error", "variant equality is not in v1") {
		t.Fatalf("expected variant-equality refusal, got %v", diags)
	}
}

// TestVariantParentRecordCollision pins the identity prerequisite:
// a variant parent colliding with a record type name is rejected,
// so emit never maps one TS name to two shapes.
func TestVariantParentRecordCollision(t *testing.T) {
	src := `mod m
  provides [m__noop, M__Out, Login__State]
  uses []
  emits []

variant Login__State rev 1 (
  case Anonymous()
)

type Login__State rev 1 (
  user_id: str
)

type M__Out rev 1 (
  value: int
)

fn m__noop(value: int) -> M__Out rev 1
  emits []
  tests
    go(value = 1) => Ok(value = 1)
=
  Ok(value = value)
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if !hasDiag(diags, "error", "collides with record type Login__State") {
		t.Fatalf("expected parent/record collision, got %v", diags)
	}
}

// TestVariantTagIdentity pins tag-first equality: two cases with
// the same payload shape are different values. Static checking
// passes; the runner reports the payload mismatch.
func TestVariantTagIdentity(t *testing.T) {
	src := `mod m
  provides [m__pick, M__Out, Pick__State]
  uses []
  emits []

variant Pick__State rev 1 (
  case A(x: int)
  case B(x: int)
)

type M__Out rev 1 (
  pick: Pick__State
)

fn m__pick(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(pick = Pick__B(x = 1))
=
  Ok(pick = Pick__A(x = x))
`
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	diags := diagnose(dir, "m.ail", src)
	if !hasDiag(diags, "error", "Ok payload mismatch") {
		t.Fatalf("expected tag-identity mismatch, got %v", diags)
	}
}

// TestVariantCrossModule pins cross-module admission: a consumer
// pinning the parent constructs the foreign case, carries the
// foreign union through its own record field, and emits a typed
// import from the provider stem (records ride the same fix: the
// pre-existing unbound cross-module type reference failed strict
// tsc for both).
func TestVariantCrossModule(t *testing.T) {
	prov := `mod prov
  provides [prov__make, Prov__Data, Prov__Session, Login__State]
  uses []
  emits []

variant Login__State rev 1 (
  case Anonymous()
  case Authenticated(session: Prov__Session)
)

type Prov__Session rev 1 (
  user_id: str
)

type Prov__Data rev 1 (
  id: str
)

fn prov__make(id: str) -> Prov__Data rev 1
  emits []
  tests
    m(id = "u") => Ok(id = "u")
=
  Ok(id = id)
`
	cons := `mod cons
  provides [cons__go, Cons__Out]
  uses [prov__make@1, Prov__Data@1, Login__State@1]
  emits []

type Cons__Out rev 1 (
  held: Prov__Data
  state: Login__State
)

fn cons__go(id: str) -> Cons__Out rev 1
  emits []
  tests
    g(id = "u") => Ok(held = Prov__Data(id = "u"), state = Login__Anonymous())
=
  match call prov__make(id)
    given
      g => [exchange args (id = "u") outcome Ok(id = "u")]
    on Ok v => Ok(held = Prov__Data(id = v.id), state = Login__Anonymous())
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "prov.ail"), []byte(prov), 0o644); err != nil {
		t.Fatalf("write prov: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "cons.ail"), []byte(cons), 0o644); err != nil {
		t.Fatalf("write cons: %v", err)
	}
	if diags := diagnose(dir, "cons.ail", cons); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "prov.ail"), filepath.Join(dir, "cons.ail")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "cons.ts"))
	if err != nil {
		t.Fatalf("read cons.ts: %v", err)
	}
	ts := string(raw)
	if !strings.Contains(ts, `import { prov__make, type Login__State, type Prov__Data, type ProvResult } from "./prov";`) {
		t.Fatalf("missing provider type imports in emit:\n%s", ts)
	}
	if !strings.Contains(ts, `{ $ail_kind: "Login__Anonymous" }`) {
		t.Fatalf("missing nullary case construction in emit:\n%s", ts)
	}
}

// TestVariantNoCatalog pins catalog silence: constructing cases
// creates no errors.json entries.
func TestVariantNoCatalog(t *testing.T) {
	src := variantValueMod(variantValueBody, variantValueExpect)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.ail"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.ail")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "errors.json"))
	if err != nil {
		t.Fatalf("read errors.json: %v", err)
	}
	if strings.Contains(string(raw), "Login") {
		t.Fatalf("variant leaked into errors catalog: %s", raw)
	}
}

// TestVariantEmitUnion pins TS emission: one exported union type
// per variant, qualified tags as discriminants, payload fields
// mapped through their record types.
func TestVariantEmitUnion(t *testing.T) {
	src := variantValueMod(variantValueBody, variantValueExpect)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.ail"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.ail")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatalf("read m.ts: %v", err)
	}
	ts := string(raw)
	if !strings.Contains(ts, `export type Login__State = { $ail_kind: "Login__Anonymous" } | { $ail_kind: "Login__Authenticated"; session: Auth__Session }`) {
		t.Fatalf("missing variant union in emit:\n%s", ts)
	}
	if !strings.Contains(ts, `{ $ail_kind: "Login__Authenticated", session:`) {
		t.Fatalf("missing tagged construction in emit:\n%s", ts)
	}
}
