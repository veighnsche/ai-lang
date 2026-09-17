package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var jsonDB = strings.Replace(lspDB, "provides [db__get, Db__U]", "provides [db__get, db__ping, Db__U]", 1) + `
fn db__ping() -> Db__U rev 1
  emits []
  tests
    ok() => Ok(id = "u")
=
  Ok(id = "u")
`

const jsonAuth = `mod auth
  provides [auth__go, Auth__S]
  uses [db__get@1, db__ping@1]
  emits [auth.bad]

error auth.bad()

type Auth__S rev 1 (
  id: str
)

fn auth__go(id: str) -> Auth__S rev 1
  emits [auth.bad, auth.stale]
  tests
    ok(id = "u") => Ok(id = "u")
    down(id = "u") => auth.bad()
    extra(id = "u") => auth.bad()
=
  match call db__get(id)
    given
      ok => [exchange args (id = "u") outcome Ok(id = "u")]
      down => [exchange args (id = "u") outcome db.down()]
      zzz => [exchange args (id = "u") outcome db.down()]
    on db.down _ => auth.bad()
    on Ok u => Ok(id = u.id)
`

func TestGoldenJSONDiags(t *testing.T) {
	dir := t.TempDir()
	for name, body := range map[string]string{"db.ail": jsonDB, "auth.ail": jsonAuth} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mods, texts, collected, err := parsePaths([]string{filepath.Join(dir, "db.ail"), filepath.Join(dir, "auth.ail")})
	if err != nil {
		t.Fatal(err)
	}
	_, collected = checkProgram(mods, texts, collected, nil)
	var buf bytes.Buffer
	reportDiags(&buf, collected)
	got := buf.String()
	if got != jsonGolden {
		t.Fatalf("JSON golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, jsonGolden)
	}
}

const jsonGolden = `{"code":"AIL3401","sev":"warning","file":"auth.ail","line":3,"start":19,"end":29,"msg":"uses db__ping@1 but auth never calls it"}
{"code":"AIL4002","sev":"error","file":"auth.ail","line":13,"start":19,"end":29,"msg":"auth__go declares unknown error kind auth.stale in emits"}
{"code":"AIL4200","sev":"error","file":"auth.ail","line":17,"start":4,"end":9,"msg":"test extra fails: auth__go/extra: no script for call db__get"}
{"code":"AIL3105","sev":"error","file":"auth.ail","line":17,"start":4,"end":9,"msg":"test extra has no script at the call to db__get (line 19)","expected":"extra =\u003e [exchange args (...) outcome ...]","found":"extra","hint":"add a script row for this test under the given table at line 19"}
{"code":"AIL3104","sev":"warning","file":"auth.ail","line":23,"start":6,"end":9,"msg":"script zzz never runs: no test named zzz in auth__go"}
`

func TestCodesUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, c := range allCodes {
		if c == "" || seen[c] {
			t.Fatalf("duplicate or empty code %q", c)
		}
		seen[c] = true
		if !strings.HasPrefix(c, "AIL") || len(c) != 7 {
			t.Fatalf("code %q breaks the AILnnnn shape", c)
		}
	}
}

func TestAllDiagsCoded(t *testing.T) {
	registered := map[string]bool{}
	for _, c := range allCodes {
		registered[c] = true
	}
	check := func(body string, diags []Diag) {
		t.Helper()
		for _, d := range diags {
			if d.Code == "" || !registered[d.Code] {
				t.Fatalf("uncoded diagnostic %v in %q", d, body[:20])
			}
		}
	}
	if _, err := os.Stat("../../go.mod"); err == nil {
		for _, dir := range []string{"../../sketches/auth-login", "../../sketches/broken-login"} {
			entries, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				data, err := os.ReadFile(filepath.Join(dir, e.Name()))
				if err != nil {
					t.Fatal(err)
				}
				check(e.Name(), diagnose(dir, e.Name(), string(data)))
			}
		}
	}
	mutations := []string{
		strings.Replace(lspAuth, "match call db__get(id)", "match call db__nope(id)", 1),
		strings.Replace(lspAuth, "uses [db__get@1]", "uses []", 1),
		strings.Replace(lspAuth, "on db.down _ => auth.bad()", "on db.down _ => db.down()", 1),
		strings.Replace(lspAuth, "on db.down _ => auth.bad()", "on db.down _ => auth.bogus()", 1),
		strings.Replace(lspAuth, "down => [exchange args (id = \"u\") outcome db.down()]", "down => [exchange args (id = \"u\") outcome db.bogus()]", 1),
		strings.Replace(lspAuth, "ok(id = \"u\") => Ok(id = \"u\")", "ok(bogus = \"u\") => Ok(id = \"u\")", 1),
	}
	for i, bad := range mutations {
		dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
		check(string(rune('a'+i)), diagnose(dir, "auth.ail", bad))
	}
}
