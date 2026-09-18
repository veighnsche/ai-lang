package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Uses fallback (LSP multi-directory providers): an open file sees
// the same world a group CLI build would. Same-dir siblings load
// as before; uses pins no loaded module satisfies additionally
// resolve to provider files elsewhere under the go.mod-anchored
// tree, transitively and deterministically. Without a go.mod
// ancestor the world stays dir-confined, exactly as before.

// writeFallbackTree creates a rooted file tree: files maps
// slash-separated relative paths (including "go.mod") to bodies,
// and returns the root directory.
func writeFallbackTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, body := range files {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func diagnoseTreeFile(t *testing.T, root, dir, file string) []Diag {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, dir, file))
	if err != nil {
		t.Fatal(err)
	}
	return diagnose(filepath.Join(root, dir), file, string(raw))
}

// A cross-directory function pin resolves through the fallback:
// the lspAuth/lspDB pair from the single-directory tests,
// split across two directories, diagnoses clean.
func TestDiagnoseUsesFallbackAcrossDirs(t *testing.T) {
	root := writeFallbackTree(t, map[string]string{
		"go.mod":    "module fallbacktest\n",
		"a/auth.can": lspAuth,
		"b/db.can":   lspDB,
	})
	if diags := diagnoseTreeFile(t, root, "a", "auth.can"); len(diags) != 0 {
		t.Fatalf("expected clean cross-dir diagnosis, got %v", diags)
	}
}

const fallbackC = `mod c
  provides [c__help, C__T]
  uses []
  emits []

type C__T rev 1 (
  value: int
)

fn c__help(value: int) -> C__T rev 1
  emits []
  tests
    one(1) => Ok(1)
  Ok(value)
`

const fallbackB = `mod b
  provides [b__help, B__U]
  uses [C__T@1]
  emits []

type B__U rev 1 (
  value: int
)

fn b__help(value: int) -> B__U rev 1
  emits []
  tests
    one(1) => Ok(1)
  Ok(value)
`

const fallbackA = `mod a
  provides [a__go, A__O]
  uses [B__U@1]
  emits []

type A__O rev 1 (
  value: int
)

fn a__go(value: int) -> A__O rev 1
  emits []
  tests
    one(1) => Ok(1)
  Ok(value)
`

// Transitive type pins close through two directories: a needs
// B__U from b, b needs C__T from c. Uppercase pins never warn as
// unused, so the carriers need no other reference.
func TestDiagnoseUsesFallbackTransitive(t *testing.T) {
	root := writeFallbackTree(t, map[string]string{
		"go.mod":  "module fallbacktest\n",
		"a/a.can": fallbackA,
		"b/b.can": fallbackB,
		"c/c.can": fallbackC,
	})
	if diags := diagnoseTreeFile(t, root, "a", "a.can"); len(diags) != 0 {
		t.Fatalf("expected clean transitive diagnosis, got %v", diags)
	}
	if diags := diagnoseTreeFile(t, root, "b", "b.can"); len(diags) != 0 {
		t.Fatalf("expected clean middle diagnosis, got %v", diags)
	}
}

// Without a go.mod ancestor the world stays dir-confined: the
// same cross-directory pin keeps today's resolves-nowhere error,
// so temp dirs and go.mod-less projects behave exactly as before.
func TestDiagnoseUsesFallbackNeedsGoMod(t *testing.T) {
	root := writeFallbackTree(t, map[string]string{
		"a/auth.can": lspAuth,
		"b/db.can":   lspDB,
	})
	diags := diagnoseTreeFile(t, root, "a", "auth.can")
	if !hasDiag(diags, "error", "resolves nowhere") {
		t.Fatalf("expected nowhere-uses error without go.mod, got %v", diags)
	}
}

// A same-file pin never loads the open file itself, even with a
// go.mod tree around: self-provides resolve nowhere, as in group
// builds.
func TestDiagnoseUsesFallbackSelfPin(t *testing.T) {
	pinned := `mod m
  provides [m__go, M__S]
  uses [M__S@1]
  emits []

type M__S rev 1 (
  value: int
)

fn m__go(value: int) -> M__S rev 1
  emits []
  tests
    one(1) => Ok(1)
  Ok(value)
`
	root := writeFallbackTree(t, map[string]string{
		"go.mod":  "module fallbacktest\n",
		"a/m.can": pinned,
	})
	diags := diagnoseTreeFile(t, root, "a", "m.can")
	if !hasDiag(diags, "error", "resolves nowhere") {
		t.Fatalf("expected self-pin rejection, got %v", diags)
	}
}

// Revision matching still applies to fallback providers: a pin
// for rev 2 against a rev-1 declaration fails instead of
// silently resolving.
func TestDiagnoseUsesFallbackRevMismatch(t *testing.T) {
	b := `mod b
  provides [b__x, B__O]
  uses []
  emits []

type B__O rev 1 (
  value: int
)

fn b__x(value: int) -> B__O rev 1
  emits []
  tests
    one(1) => Ok(1)
  Ok(value)
`
	a := `mod a
  provides [a__go, A__O]
  uses [b__x@2]
  emits []

type A__O rev 1 (
  value: int
)

fn a__go(value: int) -> A__O rev 1
  emits []
  tests
    one(1) => Ok(1)
  Ok(value)
`
	root := writeFallbackTree(t, map[string]string{
		"go.mod":  "module fallbacktest\n",
		"a/a.can": a,
		"b/b.can": b,
	})
	diags := diagnoseTreeFile(t, root, "a", "a.can")
	if !hasDiag(diags, "error", "rev") {
		t.Fatalf("expected rev-mismatch error, got %v", diags)
	}
}

// An unparseable provider file cannot provide: the pin keeps
// today's error instead of failing the diagnosis.
func TestDiagnoseUsesFallbackSkipsUnparseable(t *testing.T) {
	root := writeFallbackTree(t, map[string]string{
		"go.mod":  "module fallbacktest\n",
		"a/a.can": fallbackA,
		"b/b.can": "mod b\n  provides [",
		"c/c.can": fallbackC,
	})
	diags := diagnoseTreeFile(t, root, "a", "a.can")
	if !hasDiag(diags, "error", "resolves nowhere") {
		t.Fatalf("expected nowhere-uses error past unparseable provider, got %v", diags)
	}
}
