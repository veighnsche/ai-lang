package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// b04: a Seq over a multi-argument instantiation — the Map<K,V>
// entries shape. The literal head's comma must not split the
// argument list; expansion stamps the nested mention; the value
// reads back through indexed projection.
const seqNestedPair = `mod m
  provides [m__go, M__Pair, M__Map, M__Out]
  uses []
  emits []

type M__Pair<K, V> rev 1 (
  key: K
  val: V
)

type M__Map<K, V> rev 1 (
  entries: Seq<M__Pair<K, V>>
)

type M__Out rev 1 (
  key: str
)

fn m__go(mp: M__Map<str, int>) -> M__Out rev 1
  emits []
  tests
    t(M__Map<str, int>(Seq<M__Pair<str, int>>[M__Pair<str, int>("a", 1)])) => Ok("a")
  Ok(mp.entries[0].key)
`

func TestSeqNestedPair(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": seqNestedPair})
	if diags := diagnose(dir, "m.can", seqNestedPair); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestSeqNestedPairEmit(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "m.can"), []byte(seqNestedPair), 0o644); err != nil {
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
	src := string(got)
	for _, want := range []string{"M__Pair$", "entries: ", "[]", ").key"} {
		if !strings.Contains(src, want) {
			t.Errorf("emit missing %q\n--- emit ---\n%s", want, src)
		}
	}
}
