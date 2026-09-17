package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const goodDB = `mod db
  provides [db__get]
  uses []
  emits [db.down]

error db.down()

type Db__U rev 1 (
  id: str
)

fn db__get(id: str) -> Db__U rev 1
  emits [db.down]
  tests
    ok(id = "u") => Ok(id = "u")
=
  match id
    "u" => Ok(id = "u")
    _ => db.down()
`

const goodAuth = `mod auth
  provides [auth__go]
  uses [db__get@1]
  emits [auth.bad]

error auth.bad()

type Auth__S rev 1 (
  id: str
)

fn auth__go(id: str) -> Auth__S rev 1
  emits [auth.bad]
  tests
    ok(id = "u") => Ok(id = "u")
    down(id = "u") => auth.bad
=
  match call db__get(id)
    given
      ok => [Ok(id = "u")]
      down => [db.down()]
    on db.down _ => auth.bad()
    on Ok u => Ok(id = u.id)
`

func writeFixtures(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range files {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestRepoPasses(t *testing.T) {
	if _, err := os.Stat("../../go.mod"); err != nil {
		t.Skip("not in repo checkout")
	}
	_, _, _, errs := check([]string{"../../sketches", "../../std"})
	if len(errs) > 0 {
		t.Fatalf("repo check failed: %v", errs)
	}
}

func TestGoodPair(t *testing.T) {
	dir := writeFixtures(t, map[string]string{"db.ail": goodDB, "auth.ail": goodAuth})
	_, _, _, errs := check([]string{dir})
	if len(errs) > 0 {
		t.Fatalf("expected pass, got %v", errs)
	}
}

func TestUnpinnedUses(t *testing.T) {
	bad := strings.Replace(goodAuth, "db__get@1", "db__get", 1)
	dir := writeFixtures(t, map[string]string{"db.ail": goodDB, "auth.ail": bad})
	_, _, _, errs := check([]string{dir})
	if !contains(errs, "must pin a rev") {
		t.Fatalf("expected pin error, got %v", errs)
	}
}

func TestUnresolvableUses(t *testing.T) {
	bad := strings.Replace(goodAuth, "db__get@1", "ghost@9", 1)
	dir := writeFixtures(t, map[string]string{"db.ail": goodDB, "auth.ail": bad})
	_, _, _, errs := check([]string{dir})
	if !contains(errs, "resolves nowhere") {
		t.Fatalf("expected resolve error, got %v", errs)
	}
}

func TestExternRedeclare(t *testing.T) {
	bad := "extern fn db__get(id: str) -> Db__U\n" + goodAuth
	dir := writeFixtures(t, map[string]string{"db.ail": goodDB, "auth.ail": bad})
	_, _, _, errs := check([]string{dir})
	if !contains(errs, "re-declares") {
		t.Fatalf("expected re-declare error, got %v", errs)
	}
}

func TestBracesBanned(t *testing.T) {
	bad := goodDB + "// { stray brace }\n"
	dir := writeFixtures(t, map[string]string{"db.ail": bad, "auth.ail": goodAuth})
	_, _, _, errs := check([]string{dir})
	if !contains(errs, "curly braces are banned") {
		t.Fatalf("expected braces error, got %v", errs)
	}
}

// R1 (v45): braces inside string literals are data, so a fixture
// carrying "{u}" in a test row passes; a brace in code still fails.
func TestBracesInStringsAllowed(t *testing.T) {
	withBraces := strings.Replace(goodDB,
		"    ok(id = \"u\") => Ok(id = \"u\")\n",
		"    ok(id = \"u\") => Ok(id = \"u\")\n    braced(id = \"{u}\") => db.down()\n", 1)
	dir := writeFixtures(t, map[string]string{"db.ail": withBraces, "auth.ail": goodAuth})
	if _, _, _, errs := check([]string{dir}); len(errs) > 0 {
		t.Fatalf("expected braces in strings to pass, got %v", errs)
	}
	bad := withBraces + "  { stray }\n"
	dir = writeFixtures(t, map[string]string{"db.ail": bad, "auth.ail": goodAuth})
	if _, _, _, errs := check([]string{dir}); !contains(errs, "curly braces are banned") {
		t.Fatalf("expected braces error, got %v", errs)
	}
}

func TestExternSectionGone(t *testing.T) {
	bad := strings.Replace(goodAuth, "  tests", "  externals:\n  tests", 1)
	dir := writeFixtures(t, map[string]string{"db.ail": goodDB, "auth.ail": bad})
	_, _, _, errs := check([]string{dir})
	if !contains(errs, "externals section is gone") {
		t.Fatalf("expected externals error, got %v", errs)
	}
}

func TestGivenMustBeTotal(t *testing.T) {
	bad := strings.Replace(goodAuth, "      down => [db.down()]\n", "", 1)
	dir := writeFixtures(t, map[string]string{"db.ail": goodDB, "auth.ail": bad})
	_, _, _, errs := check([]string{dir})
	if !contains(errs, "given table") {
		t.Fatalf("expected given-table error, got %v", errs)
	}
}

func TestLegacySkipped(t *testing.T) {
	old := goodAuth + "\n// SUPERSEDED-BY: whatever\n"
	dir := writeFixtures(t, map[string]string{"db.ail": goodDB, "auth.ail": old})
	current, skipped, _, errs := check([]string{dir})
	if len(errs) > 0 {
		t.Fatalf("expected pass, got %v", errs)
	}
	if current != 1 || skipped != 1 {
		t.Fatalf("expected 1 current + 1 skipped, got %d + %d", current, skipped)
	}
}

func contains(errs []string, sub string) bool {
	for _, e := range errs {
		if strings.Contains(e, sub) {
			return true
		}
	}
	return false
}

func TestDemoExpects(t *testing.T) {
	partial := strings.Replace(goodAuth, "      down => [db.down()]\n", "", 1)
	marker := "// DEMO-EXPECTS: given table [ok] != tests [down ok]\n"
	dir := writeFixtures(t, map[string]string{"db.ail": goodDB, "auth.ail": marker + partial})
	if _, _, _, errs := check([]string{dir}); len(errs) > 0 {
		t.Fatalf("expected declared demo violation to pass, got %v", errs)
	}
	healed := strings.Replace(partial, "      ok => [Ok(id = \"u\")]\n", "      ok => [Ok(id = \"u\")]\n      down => [db.down()]\n", 1)
	dir = writeFixtures(t, map[string]string{"db.ail": goodDB, "auth.ail": marker + healed})
	if _, _, _, errs := check([]string{dir}); !contains(errs, "demo expects") {
		t.Fatalf("expected healed-demo error, got %v", errs)
	}
	rotted := marker + partial + "// { stray brace }\n"
	dir = writeFixtures(t, map[string]string{"db.ail": goodDB, "auth.ail": rotted})
	if _, _, _, errs := check([]string{dir}); !contains(errs, "curly braces are banned") {
		t.Fatalf("expected extra-violation error, got %v", errs)
	}
}

func TestNewDeclShapes(t *testing.T) {
	db := strings.Replace(goodDB, "type Db__U rev 1 (\n  id: str\n)", "brand Db__Hash is str rev 1\n\ntype Db__U rev 1 (\n  id: str\n  pw_hash: Db__Hash\n)", 1)
	db = strings.Replace(db, "provides [db__get]", "provides [db__get, Db__Hash]", 1)
	auth := strings.Replace(goodAuth, "provides [auth__go]", "provides [auth__go, Auth__Pw, m__use]", 1)
	auth = strings.Replace(auth, "uses [db__get@1]", "uses [db__get@1, Db__Hash@1]", 1)
	auth += "\nbrand Auth__Pw is str rev 1\n\nextern m__use(pw: Auth__Pw) -> Auth__S rev 1\n  emits [auth.bad]\n"
	auth = strings.Replace(auth, "ok(id = \"u\") => Ok(id = \"u\")", "ok(id = \"u\") => Ok(id = d\"1.5\")", 1)
	dir := writeFixtures(t, map[string]string{"db.ail": db, "auth.ail": auth})
	_, _, _, errs := check([]string{dir})
	if len(errs) > 0 {
		t.Fatalf("expected pass, got %v", errs)
	}
}
