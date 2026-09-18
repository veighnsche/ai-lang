package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// a75: variant elimination. A value match over a variant-typed
// scrutinee takes exactly one on-Case arm per case, binds payloads
// case-specifically, and lowers to an evaluate-once tag switch.
// Probes first: every TestVariantMatch* below fails before the slice.

const variantMatchDecls = `variant Login__State rev 1 (
  case Anonymous()
  case Authenticated(session: Auth__Session)
  case Locked(user_id: str, remaining_seconds: int)
)

type Auth__Session rev 1 (
  user_id: str
)

type M__Out rev 1 (
  message: str
)

`

const variantMatchBody = `match state
    on Login__Anonymous _ => Ok(message = "Sign in")
    on Login__Authenticated a => Ok(message = a.session.user_id)
    on Login__Locked l => Ok(message = "locked: " + l.user_id)`

func variantMatchMod(tests, body string) string {
	return `mod m
  provides [m__label, M__Out, Login__State, Auth__Session]
  uses []
  emits []

` + variantMatchDecls + `fn m__label(state: Login__State) -> M__Out rev 1
  emits []
  tests
` + tests + `
  ` + body + `
`
}

const variantMatchTests = `    anon(Login__Anonymous()) => Ok(message = "Sign in")
    auth(Login__Authenticated(session = Auth__Session(user_id = "u"))) => Ok(message = "u")
    lock(Login__Locked(user_id = "u", remaining_seconds = 3)) => Ok(message = "locked: u")`

// TestVariantMatchClean pins the positive: three cases, three arms,
// payload projection through binders, one test per arm.
func TestVariantMatchClean(t *testing.T) {
	src := variantMatchMod(variantMatchTests, variantMatchBody)
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// TestVariantMatchMissingCase pins exhaustiveness: the Locked arm
// is gone, so the checker names it.
func TestVariantMatchMissingCase(t *testing.T) {
	body := `match state
    on Login__Anonymous _ => Ok(message = "Sign in")
    on Login__Authenticated a => Ok(message = a.session.user_id)`
	src := variantMatchMod(variantMatchTests, body)
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if !hasDiag(diags, "error", "missing Login__Locked") {
		t.Fatalf("expected missing-case error, got %v", diags)
	}
	if !hasCode(diags, "CAN4101") {
		t.Fatalf("expected CAN4101, got %v", diags)
	}
}

// TestVariantMatchWrongUnion pins membership: Pick__A belongs to
// another union, so the arm is stale.
func TestVariantMatchWrongUnion(t *testing.T) {
	body := `match state
    on Login__Anonymous _ => Ok(message = "Sign in")
    on Login__Authenticated a => Ok(message = a.session.user_id)
    on Pick__A _ => Ok(message = "nope")`
	src := `mod m
  provides [m__label, M__Out, Login__State, Auth__Session, Pick__State]
  uses []
  emits []

` + variantMatchDecls + `variant Pick__State rev 1 (
  case A()
)

` + `fn m__label(state: Login__State) -> M__Out rev 1
  emits []
  tests
` + variantMatchTests + `
  ` + body + `
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if !hasDiag(diags, "error", "stale match arm Pick__A") {
		t.Fatalf("expected stale-arm error, got %v", diags)
	}
	if !hasCode(diags, "CAN4102") {
		t.Fatalf("expected CAN4102, got %v", diags)
	}
}

// TestVariantMatchDuplicate pins arm uniqueness: two Anonymous arms.
func TestVariantMatchDuplicate(t *testing.T) {
	body := `match state
    on Login__Anonymous _ => Ok(message = "Sign in")
    on Login__Anonymous _ => Ok(message = "again")
    on Login__Authenticated a => Ok(message = a.session.user_id)
    on Login__Locked l => Ok(message = "locked: " + l.user_id)`
	src := variantMatchMod(variantMatchTests, body)
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if !hasDiag(diags, "error", "duplicate match arm Login__Anonymous") {
		t.Fatalf("expected duplicate-arm error, got %v", diags)
	}
	if !hasCode(diags, "CAN4105") {
		t.Fatalf("expected CAN4105, got %v", diags)
	}
}

// TestVariantMatchWildRejected pins case-arms-only: a catch-all _
// never stands in for a case.
func TestVariantMatchWildRejected(t *testing.T) {
	body := `match state
    on Login__Anonymous _ => Ok(message = "Sign in")
    on Login__Authenticated a => Ok(message = a.session.user_id)
    _ => Ok(message = "rest")`
	src := variantMatchMod(variantMatchTests, body)
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if !hasDiag(diags, "error", "takes case arms only") {
		t.Fatalf("expected case-arms-only error, got %v", diags)
	}
	if !hasCode(diags, "CAN4105") {
		t.Fatalf("expected CAN4105, got %v", diags)
	}
}

// TestVariantMatchMultiSlot pins single-scrutinee: a second slot
// alongside a variant is refused even for _.
func TestVariantMatchMultiSlot(t *testing.T) {
	src := `mod m
  provides [m__label, M__Out, Login__State, Auth__Session]
  uses []
  emits []

` + variantMatchDecls + `fn m__label(state: Login__State, flag: bool) -> M__Out rev 1
  emits []
  tests
    go(Login__Anonymous(), true) => Ok(message = "Sign in")
  match state, flag
    on Login__Anonymous _, _ => Ok(message = "Sign in")
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if !hasDiag(diags, "error", "exactly one scrutinee") {
		t.Fatalf("expected single-scrutinee error, got %v", diags)
	}
	if !hasCode(diags, "CAN4105") {
		t.Fatalf("expected CAN4105, got %v", diags)
	}
}

// TestVariantMatchBinderScope pins case-specific binders: session
// exists on Authenticated, not on Locked.
func TestVariantMatchBinderScope(t *testing.T) {
	body := `match state
    on Login__Anonymous _ => Ok(message = "Sign in")
    on Login__Authenticated a => Ok(message = a.session.user_id)
    on Login__Locked l => Ok(message = l.session)`
	src := variantMatchMod(variantMatchTests, body)
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if !hasDiag(diags, "error", "no field session") {
		t.Fatalf("expected binder-scope error, got %v", diags)
	}
}

// TestVariantMatchUntaken pins CAN4107 per arm: with no test taking
// the Locked arm, coverage names it.
func TestVariantMatchUntaken(t *testing.T) {
	tests := `    anon(state = Login__Anonymous()) => Ok(message = "Sign in")
    auth(state = Login__Authenticated(session = Auth__Session(user_id = "u"))) => Ok(message = "u")`
	src := variantMatchMod(tests, variantMatchBody)
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if !hasDiag(diags, "error", "no test takes on Login__Locked l") {
		t.Fatalf("expected untaken-arm error, got %v", diags)
	}
	if !hasCode(diags, "CAN4107") {
		t.Fatalf("expected CAN4107, got %v", diags)
	}
}

// TestVariantMatchCaseOnBool pins the scrutinee gate: a case pattern
// against a bool scrutinee keeps the legacy rules untouched.
func TestVariantMatchCaseOnBool(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out, Login__State, Auth__Session]
  uses []
  emits []

` + variantMatchDecls + `fn m__go(flag: bool) -> M__Out rev 1
  emits []
  tests
    go(true) => Ok(message = "Sign in")
  match flag
    on Login__Anonymous _ => Ok(message = "Sign in")
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if !hasDiag(diags, "error", "variant pattern on a non-variant match") {
		t.Fatalf("expected non-variant-scrutinee error, got %v", diags)
	}
	if !hasCode(diags, "CAN4106") {
		t.Fatalf("expected CAN4106, got %v", diags)
	}
}

// TestVariantMatchCaseScrutinee pins the elimination boundary: an
// already-eliminated binder is a case, not a union, so matching
// over it is refused instead of re-proving the parent.
func TestVariantMatchCaseScrutinee(t *testing.T) {
	body := `match state
    on Login__Anonymous _ => Ok(message = "Sign in")
    on Login__Authenticated a => match a
      on Login__Authenticated _ => Ok(message = "u")
      on Login__Anonymous _ => Ok(message = "no")
      on Login__Locked _ => Ok(message = "no")
    on Login__Locked l => Ok(message = "locked: " + l.user_id)`
	src := variantMatchMod(variantMatchTests, body)
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	if !hasDiag(diags, "error", "eliminates unions, not cases") {
		t.Fatalf("expected case-scrutinee error, got %v", diags)
	}
	if !hasCode(diags, "CAN4105") {
		t.Fatalf("expected CAN4105, got %v", diags)
	}
}

// TestVariantMatchNested pins nested obligations: an inner match
// over a second variant inside an arm proves separately.
func TestVariantMatchNested(t *testing.T) {
	src := `mod m
  provides [m__pick, M__Out, Login__State, Auth__Session, Pick__State]
  uses []
  emits []

` + variantMatchDecls + `variant Pick__State rev 1 (
  case A()
  case B()
)

fn m__pick(state: Login__State, pick: Pick__State) -> M__Out rev 1
  emits []
  tests
    aa(Login__Anonymous(), Pick__A()) => Ok(message = "A")
    ab(Login__Anonymous(), Pick__B()) => Ok(message = "B")
    au(Login__Authenticated(session = Auth__Session(user_id = "u")), Pick__A()) => Ok(message = "u")
    al(Login__Locked(user_id = "u", remaining_seconds = 1), Pick__A()) => Ok(message = "locked")
  match state
    on Login__Anonymous _ => match pick
      on Pick__A _ => Ok(message = "A")
      on Pick__B _ => Ok(message = "B")
    on Login__Authenticated a => Ok(message = a.session.user_id)
    on Login__Locked l => Ok(message = "locked")
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	if diags := diagnose(dir, "m.can", src); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// TestVariantMatchCrossModule pins cross-module elimination: a
// consumer pinning the parent matches a foreign scrutinee with
// all cases, and the union type rides the existing import.
func TestVariantMatchCrossModule(t *testing.T) {
	prov := `mod prov
  provides [prov__go, Prov__Out, Login__State, Auth__Session]
  uses []
  emits []

variant Login__State rev 1 (
  case Anonymous()
  case Authenticated(session: Auth__Session)
)

type Auth__Session rev 1 (
  user_id: str
)

type Prov__Out rev 1 (
  message: str
)

fn prov__go(state: Login__State) -> Prov__Out rev 1
  emits []
  tests
    anon(Login__Anonymous()) => Ok(message = "in")
    auth(Login__Authenticated(session = Auth__Session(user_id = "u"))) => Ok(message = "u")
  match state
    on Login__Anonymous _ => Ok(message = "in")
    on Login__Authenticated a => Ok(message = a.session.user_id)
`
	cons := `mod cons
  provides [cons__go, Cons__Out]
  uses [Login__State@1, Auth__Session@1]
  emits []

type Cons__Out rev 1 (
  message: str
)

fn cons__go(state: Login__State) -> Cons__Out rev 1
  emits []
  tests
    anon(Login__Anonymous()) => Ok(message = "in")
    auth(Login__Authenticated(session = Auth__Session(user_id = "u"))) => Ok(message = "u")
  match state
    on Login__Anonymous _ => Ok(message = "in")
    on Login__Authenticated a => Ok(message = a.session.user_id)
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "prov.can"), []byte(prov), 0o644); err != nil {
		t.Fatalf("write prov: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "cons.can"), []byte(cons), 0o644); err != nil {
		t.Fatalf("write cons: %v", err)
	}
	if diags := diagnose(dir, "cons.can", cons); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "prov.can"), filepath.Join(dir, "cons.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "cons.ts"))
	if err != nil {
		t.Fatalf("read cons.ts: %v", err)
	}
	ts := string(raw)
	if !strings.Contains(ts, `type Login__State`) || !strings.Contains(ts, `from "./prov"`) {
		t.Fatalf("missing provider union import in emit:\n%s", ts)
	}
	if !strings.Contains(ts, `case "Login__Authenticated": {`) {
		t.Fatalf("missing case arm in emit:\n%s", ts)
	}
}

// TestVariantMatchNoCatalog pins catalog silence: eliminating cases
// creates no errors.json entries.
func TestVariantMatchNoCatalog(t *testing.T) {
	src := variantMatchMod(variantMatchTests, variantMatchBody)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.can"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "errors.json"))
	if err != nil {
		t.Fatalf("read errors.json: %v", err)
	}
	if strings.Contains(string(raw), "Login") {
		t.Fatalf("variant match leaked into errors catalog: %s", raw)
	}
}

// TestVariantMatchEmitSwitch pins TS elimination: one evaluate-once
// switch over the scrutinee temporary, one case per tag, payload
// binders lowering to the narrowed temporary.
func TestVariantMatchEmitSwitch(t *testing.T) {
	src := variantMatchMod(variantMatchTests, variantMatchBody)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.can"), []byte(src), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatalf("read m.ts: %v", err)
	}
	ts := string(raw)
	if !strings.Contains(ts, `.$can_kind`) || !strings.Contains(ts, "switch (") {
		t.Fatalf("missing tag switch in emit:\n%s", ts)
	}
	if !strings.Contains(ts, `case "Login__Authenticated": {`) {
		t.Fatalf("missing case arm in emit:\n%s", ts)
	}
	if !strings.Contains(ts, "const a = $can_m") {
		t.Fatalf("missing payload binder in emit:\n%s", ts)
	}
}
