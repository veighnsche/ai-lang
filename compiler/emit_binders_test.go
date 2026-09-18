package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Error arms bind their payload like Ok arms do; every arm body is
// block-scoped so one binder name can serve two arms; generated
// temporaries live in an unspellable namespace so a source _m1 never
// collides with them.

const binderFixture = `mod audit
  provides [audit__go, audit__wrap, audit__fail, audit__double, Audit__Value]
  uses []
  emits [audit.failed, audit.other]

error audit.failed(code: int)
error audit.other(code: int)

type Audit__Value rev 1 (
  value: int
)

fn audit__fail(which: str) -> Audit__Value rev 1
  emits [audit.failed, audit.other]
  tests
    f("a") => audit.failed(code = 7)
    g("b") => audit.other(code = 8)
    o("z") => Ok(value = 0)
=
  match which
    "a" => audit.failed(code = 7)
    "b" => audit.other(code = 8)
    _ => Ok(value = 0)

fn audit__double(_m1: int) -> Audit__Value rev 1
  emits []
  tests
    d(21) => Ok(value = 42)
=
  Ok(value = _m1 + _m1)

fn audit__go(flag: str) -> Audit__Value rev 1
  emits [audit.failed, audit.other]
  tests
    h("a") => Ok(value = 7)
    k("b") => Ok(value = 8)
    j("z") => Ok(value = 0)
=
  match call audit__fail(flag)
    on audit.failed err => Ok(value = err.code)
    on audit.other err => Ok(value = err.code)
    on Ok r => Ok(value = r.value)

fn audit__wrap(flag: str) -> Audit__Value rev 1
  emits [audit.failed, audit.other]
  tests
    w("a") => Ok(value = 28)
    v("b") => Ok(value = 8)
    u("z") => Ok(value = 0)
=
  match call audit__fail(flag)
    on audit.failed err => match call audit__double(err.code)
      on Ok r => Ok(value = r.value + r.value)
    on audit.other err => Ok(value = err.code)
    on Ok r => Ok(value = r.value)
`

func TestEmitErrorBinders(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"audit.can": binderFixture})
	// diagnose runs every decision table: code = 7 read through the
	// error binder, the shared binder name in two arms, the nested
	// call, and the _m1 source parameter all evaluate here.
	if diags := diagnose(dir, "audit.can", binderFixture); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestEmitErrorBindersShape(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "audit.can"), []byte(binderFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "audit.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(out, "audit.ts"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(got)
	for _, want := range []string{
		`case "audit.failed": {`,
		`case "audit.other": {`,
		`case "ok": {`,
		"const err = $can_m",
		"const $can_m",
		"(_m1 + _m1)",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("emit missing %q\n--- emit ---\n%s", want, src)
		}
	}
	for _, banned := range []string{"const _m", "case \"ok\":\n"} {
		if strings.Contains(src, banned) {
			t.Errorf("emit contains %q (unscoped arm or spellable temp)\n--- emit ---\n%s", banned, src)
		}
	}
}
