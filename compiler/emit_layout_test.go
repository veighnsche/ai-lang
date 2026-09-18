package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Outcome metadata and logical payload fields are disjoint: the tag
// lives under the unspellable $can_kind key, and every declared field
// — including kind, __proto__, and constructor — becomes an own data
// property with its source name unchanged.

const layoutFixture = `mod audit
  provides [audit__go, audit__proto, audit__ctor, Audit__Value, Audit__Proto, Audit__Ctor]
  uses []
  emits [audit.bad]

error audit.bad(kind: str, code: int)

type Audit__Value rev 1 (
  kind: str
)

type Audit__Proto rev 1 (
  __proto__: str
)

type Audit__Ctor rev 1 (
  constructor: str
)

fn audit__go(flag: str) -> Audit__Value rev 1
  emits [audit.bad]
  tests
    k("x") => Ok(kind = "request")
    p("y") => Ok(kind = "other")
    e("z") => audit.bad(kind = "broken", code = 7)
=
  match flag
    "x" => Ok(kind = "request")
    "y" => Ok(kind = "other")
    _ => audit.bad(kind = "broken", code = 7)

fn audit__proto() -> Audit__Proto rev 1
  emits []
  tests
    w() => Ok(__proto__ = "kept")
=
  Ok(__proto__ = "kept")

fn audit__ctor() -> Audit__Ctor rev 1
  emits []
  tests
    c() => Ok(constructor = "mine")
=
  Ok(constructor = "mine")
`

func TestPayloadLayoutEvaluates(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"audit.can": layoutFixture})
	// The evaluator keeps tag and payload disjoint by construction;
	// these tables pin the data side: kind as a value, the error tag
	// beside its kind field, and the exact code.
	if diags := diagnose(dir, "audit.can", layoutFixture); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestPayloadLayoutEmit(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "audit.can"), []byte(layoutFixture), 0o644); err != nil {
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
		`$can_kind: "ok"`,
		`$can_kind: "audit.bad"`,
		`kind: "request"`,
		`{ $can_kind: "audit.bad"; kind: string; code: bigint }`,
		`["__proto__"]: "kept"`,
		`constructor: "mine"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("emit missing %q\n--- emit ---\n%s", want, src)
		}
	}
	// The flat envelope is gone: no bare kind tag in values, types,
	// or discriminants.
	for _, banned := range []string{
		`{ kind: "ok"`,
		`{ kind: "audit.bad"`,
		`.kind) {`,
	} {
		if strings.Contains(src, banned) {
			t.Errorf("emit contains flat envelope %q\n--- emit ---\n%s", banned, src)
		}
	}
}
