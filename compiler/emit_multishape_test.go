package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A file may mix Ok shapes: a validator returns its accepted value
// while its consumer returns a wider record, and both verify and emit
// in one module. The prod union carries one ok member per shape.

const multiShapeSrc = `mod m
  provides [m__check, m__use, M__Out, M__Val]
  uses []
  emits [m.too_big]

error m.too_big(value: int, limit: int)

type M__Out rev 1 (
  got: int
  limit: int
)

type M__Val rev 1 (
  value: int
)

fn m__check(value: int, limit: int) -> M__Val rev 1
  emits [m.too_big]
  tests
    fits(value = 3, limit = 10) => Ok(value = 3)
    spills(value = 11, limit = 10) => m.too_big(value = 11, limit = 10)
=
  match value <= limit
    true => Ok(value = value)
    false => m.too_big(value = value, limit = limit)

fn m__use(value: int, limit: int) -> M__Out rev 1
  emits [m.too_big]
  tests
    ok(value = 3, limit = 10) => Ok(got = 3, limit = 10)
    bad(value = 11, limit = 10) => m.too_big(value = 11, limit = 10)
=
  match call m__check(value, limit)
    on m.too_big _ => m.too_big(value = value, limit = limit)
    on Ok c => Ok(got = c.value, limit = limit)
`

func TestMultiShapeEmit(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": multiShapeSrc})
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	ts := string(raw)
	for _, want := range []string{
		`{ $can_kind: "ok"; value: bigint }`,
		`{ $can_kind: "ok"; got: bigint; limit: bigint }`,
		`{ $can_kind: "m.too_big"; value: bigint; limit: bigint }`,
	} {
		if !strings.Contains(ts, want) {
			t.Errorf("emit missing %s\n%s", want, ts)
		}
	}
}

// Per-function result types: call temporaries and return annotations
// carry the callee's own outcomes, never the module-wide union, so a
// strict checker narrows each handled outcome to its exact payload
// shape (a14: emit narrowing, not suppressions). Every call-switch
// also ends in an unreachable default arm.
func TestPerFnResultUnions(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": multiShapeSrc})
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	ts := string(raw)
	// m__use returns only its own ok shape plus the errors it emits:
	// the value-shape ok member of m__check must not appear.
	wantRet := `export function m__use(value: bigint, limit: bigint): { $can_kind: "ok"; got: bigint; limit: bigint } | { $can_kind: "m.too_big"; value: bigint; limit: bigint } {`
	if !strings.Contains(ts, wantRet) {
		t.Errorf("emit missing per-function return:\n%s", ts)
	}
	// The m__check call temporary carries m__check's union, not MResult.
	wantTmp := `const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "m.too_big"; value: bigint; limit: bigint } = m__check(value, limit);`
	if !strings.Contains(ts, wantTmp) {
		t.Errorf("emit missing per-function call temporary:\n%s", ts)
	}
	if strings.Contains(ts, "const $can_m1: MResult") {
		t.Errorf("call temporary uses module-wide union:\n%s", ts)
	}
	if !strings.Contains(ts, "throw new Error(\"unreachable\");") {
		t.Errorf("call-switch missing unreachable default:\n%s", ts)
	}
}

// A single-shape module keeps exactly one ok member in the union,
// now under the disjoint tag.
func TestSingleShapeEmitUnchanged(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": helperClean})
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	line := ""
	for _, l := range strings.Split(string(raw), "\n") {
		if strings.HasPrefix(l, "export type MResult") {
			line = l
		}
	}
	want := `export type MResult = { $can_kind: "ok"; id: string } | { $can_kind: "m.bad" };`
	if line != want {
		t.Errorf("single-shape union changed:\n got: %s\nwant: %s", line, want)
	}
}
