package main

import (
	"os"
	"strings"
	"testing"
)

// Effects: module-private cells, declared capabilities, per-test
// isolation. Store operations are deterministic, so they take no
// given tables; authority must be proven before anything runs.

const effectsBase = `mod m
  provides [m__bump, M__T]
  uses []
  emits []

state M__C: int = 0

type M__T rev 1 (
  total: int
)

fn m__bump(by: int) -> M__T rev 1
  effects [M__C.read, M__C.write]
  emits []
  tests
    three(3) => Ok(total = 3)
  match call state__get(M__C)
    on Ok c => match call state__put(M__C, c.value + by)
      on Ok _ => Ok(total = c.value + by)
`

func TestEffectsClean(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": effectsBase})
	if diags := diagnose(dir, "m.can", effectsBase); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestEffectsSketchClean(t *testing.T) {
	raw, err := os.ReadFile("../sketches/counter/counter.can")
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	dir := writeLSPDir(t, map[string]string{"counter.can": text})
	if diags := diagnose(dir, "counter.can", text); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestEffectsUndeclaredDirect(t *testing.T) {
	bad := strings.Replace(effectsBase,
		"effects [M__C.read, M__C.write]", "effects [M__C.read]", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "uses M__C.write without declaring it") {
		t.Fatalf("expected undeclared-effect error, got %v", diags)
	}
}

func TestEffectsTransitive(t *testing.T) {
	flow := `mod m
  provides [m__go, m__help, M__T]
  uses []
  emits []

state M__C: int = 0

type M__T rev 1 (
  total: int
)

fn m__help(by: int) -> M__T rev 1
  effects [M__C.read, M__C.write]
  emits []
  tests
    h(1) => Ok(total = 1)
  match call state__get(M__C)
    on Ok c => match call state__put(M__C, c.value + by)
      on Ok _ => Ok(total = c.value + by)

fn m__go(by: int) -> M__T rev 1
  emits []
  tests
    g(1) => Ok(total = 1)
  match call m__help(by)
    on Ok s => Ok(total = s.total)
`
	dir := writeLSPDir(t, map[string]string{"m.can": flow})
	diags := diagnose(dir, "m.can", flow)
	if !hasDiag(diags, "error", "via helper m__help without declaring it") {
		t.Fatalf("expected transitive-effect error, got %v", diags)
	}
	if hasDiag(diags, "error", "fails:") {
		t.Fatalf("gate must block tests before running, got %v", diags)
	}
}

func TestEffectsStale(t *testing.T) {
	stale := strings.Replace(effectsBase,
		"state M__C: int = 0", "state M__C: int = 0\n\nstate M__D: int = 0", 1)
	stale = strings.Replace(stale,
		"effects [M__C.read, M__C.write]", "effects [M__C.read, M__C.write, M__D.read]", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": stale})
	diags := diagnose(dir, "m.can", stale)
	if !hasDiag(diags, "error", "declares effects M__D.read but nothing uses it") {
		t.Fatalf("expected stale-effect error, got %v", diags)
	}
}

func TestEffectsUnknownCap(t *testing.T) {
	bad := strings.Replace(effectsBase,
		"effects [M__C.read, M__C.write]", "effects [Nope.read]", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "declares unknown effect Nope.read") {
		t.Fatalf("expected unknown-effect error, got %v", diags)
	}
}

func TestEffectsUnknownCell(t *testing.T) {
	bad := strings.Replace(effectsBase,
		"match call state__get(M__C)", "match call state__get(M__Nope)", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "names unknown cell M__Nope") {
		t.Fatalf("expected unknown-cell error, got %v", diags)
	}
}

func TestEffectsGivenOnStore(t *testing.T) {
	bad := strings.Replace(effectsBase,
		"  match call state__get(M__C)\n    on Ok c =>",
		"  match call state__get(M__C)\n    given\n      three => [exchange args () outcome Ok(value = 0)]\n    on Ok c =>", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "takes no given table") {
		t.Fatalf("expected given-on-store error, got %v", diags)
	}
}

func TestEffectsPutMismatch(t *testing.T) {
	bad := strings.Replace(effectsBase,
		"state__put(M__C, c.value + by)", `state__put(M__C, "x")`, 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "got str, want int") {
		t.Fatalf("expected put-mismatch error, got %v", diags)
	}
}

func TestEffectsBadInit(t *testing.T) {
	typed := strings.Replace(effectsBase,
		"state M__C: int = 0", `state M__C: int = "x"`, 1)
	dir := writeLSPDir(t, map[string]string{"m.can": typed})
	diags := diagnose(dir, "m.can", typed)
	if !hasDiag(diags, "error", "holds int but init is str") {
		t.Fatalf("expected bad-init error, got %v", diags)
	}
	lit := strings.Replace(effectsBase,
		"state M__C: int = 0", "state M__C: int = by", 1)
	dir = writeLSPDir(t, map[string]string{"m.can": lit})
	diags = diagnose(dir, "m.can", lit)
	if !hasDiag(diags, "error", "init must be a literal") {
		t.Fatalf("expected non-literal-init error, got %v", diags)
	}
}

func TestEffectsBadCellType(t *testing.T) {
	bad := strings.Replace(effectsBase,
		"state M__C: int = 0", "state M__C: Nope = 0", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "cells hold str, int, bool, or dec") {
		t.Fatalf("expected bad-cell-type error, got %v", diags)
	}
}

func TestEffectsStateNaming(t *testing.T) {
	bad := strings.Replace(effectsBase, "state M__C: int = 0", "state m__c: int = 0", 1)
	bad = strings.Replace(bad, "M__C", "m__c", -1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "must match Domain__Name") {
		t.Fatalf("expected state-naming error, got %v", diags)
	}
}

func TestEffectsReserved(t *testing.T) {
	reserved := effectsBase + `
fn state__get() -> M__T rev 1
  emits []
  tests
    q() => Ok(total = 0)
  Ok(total = 0)
`
	reserved = strings.Replace(reserved,
		"provides [m__bump, M__T]", "provides [m__bump, state__get, M__T]", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": reserved})
	diags := diagnose(dir, "m.can", reserved)
	if !hasDiag(diags, "error", "reserved for store operations") {
		t.Fatalf("expected reserved-name error, got %v", diags)
	}
}

const effectsProv = `mod p
  provides [p__do, P__T]
  uses []
  emits []

state P__C: int = 0

type P__T rev 1 (
  total: int
)

fn p__do(by: int) -> P__T rev 1
  effects [P__C.read, P__C.write]
  emits []
  tests
    t(1) => Ok(total = 1)
  match call state__get(P__C)
    on Ok c => match call state__put(P__C, c.value + by)
      on Ok _ => Ok(total = c.value + by)
`

const effectsCon = `mod c
  provides [c__go, C__T]
  uses [p__do@1]
  emits []

type C__T rev 1 (
  total: int
)

fn c__go(by: int) -> C__T rev 1
  emits []
  tests
    g(1) => Ok(total = 1)
  match call p__do(by)
    given
      g => [exchange args (by = 1) outcome Ok(total = 1)]
    on Ok s => Ok(total = s.total)
`

func TestEffectsForeignSuperset(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"p.can": effectsProv, "c.can": effectsCon})
	diags := diagnose(dir, "c.can", effectsCon)
	if !hasDiag(diags, "error", "via uses-pin p__do without declaring it") {
		t.Fatalf("expected foreign-superset error, got %v", diags)
	}
}

func TestEffectsPutPayloadUse(t *testing.T) {
	bad := strings.Replace(effectsBase,
		"on Ok _ => Ok(total = c.value + by)", "on Ok x => Ok(total = x)", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "is the empty Ok of a put") {
		t.Fatalf("expected empty-Ok-use error, got %v", diags)
	}
}

func TestEffectsWrongField(t *testing.T) {
	bad := strings.Replace(effectsBase,
		"Ok(total = c.value + by)", "Ok(total = c.bogus)", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "no field bogus on c") {
		t.Fatalf("expected no-field error, got %v", diags)
	}
}

func TestEffectsNamedArgsRejected(t *testing.T) {
	bad := strings.Replace(effectsBase,
		"match call state__get(M__C)", "match call state__get(cell = M__C)", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "takes positional args") {
		t.Fatalf("expected positional-args error, got %v", diags)
	}
	if hasDiag(diags, "error", "but nothing uses it") {
		t.Fatalf("one mistake must yield one squiggle, got %v", diags)
	}
}

func TestEffectsBoolCell(t *testing.T) {
	flag := `mod m
  provides [m__flag, M__T]
  uses []
  emits []

state M__Flag: bool = false

type M__T rev 1 (
  total: int
)

fn m__flag(by: int) -> M__T rev 1
  effects [M__Flag.read, M__Flag.write]
  emits []
  tests
    on(1) => Ok(total = 1)
  match call state__put(M__Flag, true)
    on Ok _ => match call state__get(M__Flag)
      on Ok c => Ok(total = by)
`
	dir := writeLSPDir(t, map[string]string{"m.can": flag})
	if diags := diagnose(dir, "m.can", flag); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestEffectsSharedCellName(t *testing.T) {
	mk := func(mod, fn, typ, init, want string) string {
		return `mod ` + mod + `
  provides [` + fn + `, ` + typ + `]
  uses []
  emits []

state Shrd__C: int = ` + init + `

type ` + typ + ` rev 1 (
  total: int
)

fn ` + fn + `() -> ` + typ + ` rev 1
  effects [Shrd__C.read]
  emits []
  tests
    g() => Ok(total = ` + want + `)
  match call state__get(Shrd__C)
    on Ok c => Ok(total = c.value)
`
	}
	a := mk("a", "a__go", "A__T", "1", "1")
	b := mk("b", "b__go", "B__T", "2", "2")
	dir := writeLSPDir(t, map[string]string{"a.can": a, "b.can": b})
	if diags := diagnose(dir, "a.can", a); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
	if diags := diagnose(dir, "b.can", b); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestEffectsForeignDeclared(t *testing.T) {
	con := strings.Replace(effectsCon,
		"fn c__go(by: int) -> C__T rev 1\n  emits []",
		"fn c__go(by: int) -> C__T rev 1\n  effects [P__C.read, P__C.write]\n  emits []", 1)
	dir := writeLSPDir(t, map[string]string{"p.can": effectsProv, "c.can": con})
	if diags := diagnose(dir, "c.can", con); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
