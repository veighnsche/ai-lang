package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// b03: projection off an indexed element. The walker, checker,
// evaluator, and emitter cases compose: users stays used, name
// resolves to str, the row executes, and nothing reports.
const seqProjValid = `mod m
  provides [m__first, M__User, M__Out]
  uses []
  emits []

type M__User rev 1 (
  name: str
  age: int
)

type M__Out rev 1 (
  name: str
)

fn m__first(users: Seq<M__User>) -> M__Out rev 1
  emits []
  tests
    t(Seq<M__User>[M__User("a", 3)]) => Ok("a")
  Ok(users[0].name)
`

func TestSeqProjValid(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": seqProjValid})
	if diags := diagnose(dir, "m.can", seqProjValid); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// Chains nest left: os[0].inner.label is two projections over one
// index, resolving through both records to str.
const seqProjChain = `mod m
  provides [m__deep, M__Inner, M__Outer, M__Lab]
  uses []
  emits []

type M__Inner rev 1 (
  label: str
)

type M__Outer rev 1 (
  inner: M__Inner
  n: int
)

type M__Lab rev 1 (
  label: str
)

fn m__deep(os: Seq<M__Outer>) -> M__Lab rev 1
  emits []
  tests
    t(Seq<M__Outer>[M__Outer(M__Inner("x"), 1)]) => Ok("x")
  Ok(os[0].inner.label)
`

func TestSeqProjChain(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": seqProjChain})
	if diags := diagnose(dir, "m.can", seqProjChain); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// Field, index, and projection mix: the stridx base is itself a
// ref path, and the projection reads off the indexed element.
const seqProjMixed = `mod m
  provides [m__cap, M__User, M__Team, M__Out]
  uses []
  emits []

type M__User rev 1 (
  name: str
  age: int
)

type M__Team rev 1 (
  members: Seq<M__User>
)

type M__Out rev 1 (
  name: str
)

fn m__cap(t: M__Team) -> M__Out rev 1
  emits []
  tests
    c(M__Team(Seq<M__User>[M__User("a", 3)])) => Ok("a")
  Ok(t.members[0].name)
`

func TestSeqProjMixed(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": seqProjMixed})
	if diags := diagnose(dir, "m.can", seqProjMixed); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestSeqProjBadBase(t *testing.T) {
	body := strings.Replace(seqProjValid, "users: Seq<M__User>", "users: Seq<str>", 1)
	body = strings.Replace(body, `t(Seq<M__User>[M__User("a", 3)])`, `t(Seq<str>["a"])`, 1)
	dir := writeLSPDir(t, map[string]string{"m.can": body})
	diags := diagnose(dir, "m.can", body)
	if !hasDiag(diags, "error", "cannot project .name of str: base must be a record") {
		t.Fatalf("expected bad-base error, got %v", diags)
	}
}

func TestSeqProjUnknownField(t *testing.T) {
	body := strings.Replace(seqProjValid, "users[0].name", "users[0].nope", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": body})
	diags := diagnose(dir, "m.can", body)
	if !hasDiag(diags, "error", "no field nope on M__User") {
		t.Fatalf("expected unknown-field error, got %v", diags)
	}
}

func TestSeqProjBadSyntax(t *testing.T) {
	body := strings.Replace(seqProjValid, "users[0].name", "users[0].", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": body})
	diags := diagnose(dir, "m.can", body)
	found := false
	for _, d := range diags {
		if strings.Contains(d.Msg, "bad projection") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected bad-projection parse error, got %v", diags)
	}
}

func TestSeqProjEmit(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.can"), []byte(seqProjValid), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "m.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(out, "m.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if want := "$canSeqAt(users, 0n).name"; !strings.Contains(string(got), want) {
		t.Errorf("emit missing %q\n--- emit ---\n%s", want, got)
	}
}
