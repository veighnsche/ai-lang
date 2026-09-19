package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// b02 E2: a shared extern imports from its declaring stem's stub,
// never the caller's. The fallible p__check proves the error-union
// path renders for cross-module externs too; compile() executes
// both decision tables, so the scripted error row runs end to end.
const sharedExternProvider = `mod p
  provides [p__echo, p__check, P__Out]
  uses []
  emits [p.bad]

error p.bad(code: int)

type P__Out rev 1 (
  echo: str
)

extern p__echo(s: str) -> P__Out rev 1
  emits []

extern p__check(s: str) -> P__Out rev 1
  emits [p.bad]
`

const sharedExternCaller = `mod o
  provides [o__go, o__run, O__Out]
  uses [p__echo@1, p__check@1]
  emits []

type O__Out rev 1 (
  id: str
)

fn o__go(id: str) -> O__Out rev 1
  emits []
  tests
    t("u") => Ok("u")
  match call p__echo(id)
    given
      t => [exchange args (s = "u") outcome Ok("u")]
    on Ok v => Ok(id)

fn o__run(id: str) -> O__Out rev 1
  emits []
  tests
    ok("u") => Ok("u")
    bad("u") => Ok("bad")
  match call p__check(id)
    given
      ok => [exchange args (s = "u") outcome Ok("u")]
      bad => [exchange args (s = "u") outcome p.bad(7)]
    on p.bad e => Ok("bad")
    on Ok v => Ok(id)
`

func TestEmitSharedExternDeclaringStem(t *testing.T) {
	dir := t.TempDir()
	for name, body := range map[string]string{"p.can": sharedExternProvider, "o.can": sharedExternCaller} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "p.can"), filepath.Join(dir, "o.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(out, "o.ts"))
	if err != nil {
		t.Fatal(err)
	}
	src := string(got)
	want := `import { p__check, p__echo } from "./p.externs";`
	if !strings.Contains(src, want) {
		t.Errorf("emit missing %q\n--- emit ---\n%s", want, src)
	}
	for _, banned := range []string{"./o.externs", `from "./p";`} {
		if strings.Contains(src, banned) {
			t.Errorf("emit contains %q (caller-stub or module import)\n--- emit ---\n%s", banned, src)
		}
	}
	prov, err := os.ReadFile(filepath.Join(out, "p.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(prov), ".externs") {
		t.Errorf("provider calls no externs but imports a stub\n--- emit ---\n%s", prov)
	}
}
