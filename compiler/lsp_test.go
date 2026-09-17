package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const lspDB = `mod db
  provides [db__get, Db__U]
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
    other(id = "x") => db.down()
=
  match id
    "u" => Ok(id = "u")
    _ => db.down()
`

const lspAuth = `mod auth
  provides [auth__go, Auth__S]
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
    down(id = "u") => auth.bad()
=
  match call db__get(id)
    given
      ok => [exchange args (id = "u") outcome Ok(id = "u")]
      down => [exchange args (id = "u") outcome db.down()]
    on db.down _ => auth.bad()
    on Ok u => Ok(id = u.id)
`

func writeLSPDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func hasDiag(diags []Diag, sev, sub string) bool {
	for _, d := range diags {
		if d.Sev == sev && strings.Contains(d.Msg, sub) {
			return true
		}
	}
	return false
}

func TestDiagnoseClean(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": lspAuth})
	if diags := diagnose(dir, "auth.ail", lspAuth); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestDiagnoseMissingTests(t *testing.T) {
	noTests := strings.Replace(lspAuth, "  tests\n    ok(id = \"u\") => Ok(id = \"u\")\n    down(id = \"u\") => auth.bad()\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": noTests})
	diags := diagnose(dir, "auth.ail", noTests)
	if !hasDiag(diags, "warning", "ships no tests") {
		t.Fatalf("expected missing-tests warning, got %v", diags)
	}
}

func TestDiagnoseBadName(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(", "fn go(", 1)
	bad = strings.Replace(bad, "provides [auth__go]", "provides [go]", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "must match domain__verb") {
		t.Fatalf("expected naming error, got %v", diags)
	}
}

func TestDiagnoseMissingRev(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(id: str) -> Auth__S rev 1", "fn auth__go(id: str) -> Auth__S", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "missing rev") {
		t.Fatalf("expected missing-rev error, got %v", diags)
	}
}

// Complete error expectations compare kind and payload (v12): a
// wrong payload fails the test even when the kind matches, and a
// bare kind is refused statically.
const expectErrFixture = `mod m
  provides [m__go, M__S]
  uses []
  emits [m.bad]

error m.bad(id: str)

type M__S rev 1 (
  id: str
)

fn m__go(id: str) -> M__S rev 1
  emits [m.bad]
  tests
    hit(id = "u") => Ok(id = "u")
    miss(id = "x") => m.bad(id = "x")
=
  match id
    "u" => Ok(id = id)
    _ => m.bad(id = id)
`

func TestExpectErrorPayloadMismatch(t *testing.T) {
	bad := strings.Replace(expectErrFixture, `miss(id = "x") => m.bad(id = "x")`, `miss(id = "x") => m.bad(id = "y")`, 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": bad})
	diags := diagnose(dir, "m.ail", bad)
	if !hasDiag(diags, "error", "error payload mismatch") {
		t.Fatalf("expected payload-mismatch failure, got %v", diags)
	}
}

func TestExpectBareErrorKindRefused(t *testing.T) {
	bad := strings.Replace(expectErrFixture, `miss(id = "x") => m.bad(id = "x")`, `miss(id = "x") => m.bad`, 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": bad})
	diags := diagnose(dir, "m.ail", bad)
	if !hasDiag(diags, "error", "expects bare error kind m.bad") {
		t.Fatalf("expected bare-kind refusal, got %v", diags)
	}
	var coded bool
	for _, d := range diags {
		if d.Code == CodeBareErrorKind {
			coded = true
		}
	}
	if !coded {
		t.Fatalf("expected AIL3204 code, got %v", diags)
	}
}

func TestDiagnoseFailingTest(t *testing.T) {
	bad := strings.Replace(lspAuth, "down(id = \"u\") => auth.bad()", "down(id = \"u\") => Ok(id = \"u\")", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "test down fails") {
		t.Fatalf("expected failing-test error, got %v", diags)
	}
}

func TestDiagnoseMissingArm(t *testing.T) {
	bad := strings.Replace(lspAuth, "    on db.down _ => auth.bad()\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "non-exhaustive") {
		t.Fatalf("expected exhaustiveness error, got %v", diags)
	}
}

func TestDiagnoseUnknownCallee(t *testing.T) {
	bad := strings.Replace(lspAuth, "match call db__get(id)", "match call db__nope(id)", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "calls unknown function db__nope") {
		t.Fatalf("expected unknown-callee error, got %v", diags)
	}
}

func TestDiagnoseCallNotInUses(t *testing.T) {
	bad := strings.Replace(lspAuth, "uses [db__get@1]", "uses []", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "not in uses") {
		t.Fatalf("expected not-in-uses error, got %v", diags)
	}
}

func TestDiagnoseMissingGiven(t *testing.T) {
	bad := strings.Replace(lspAuth, "    given\n      ok => [exchange args (id = \"u\") outcome Ok(id = \"u\")]\n      down => [exchange args (id = \"u\") outcome db.down()]\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "has no given table") {
		t.Fatalf("expected missing-given error, got %v", diags)
	}
}

func TestDiagnoseDanglingTest(t *testing.T) {
	bad := strings.Replace(lspAuth, "    down(id = \"u\") => auth.bad()\n", "    down(id = \"u\") => auth.bad()\n    extra(id = \"u\") => auth.bad\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "test extra has no script") {
		t.Fatalf("expected dangling-test error, got %v", diags)
	}
}

func TestDiagnoseDeadScript(t *testing.T) {
	bad := strings.Replace(lspAuth, "      down => [exchange args (id = \"u\") outcome db.down()]\n", "      down => [exchange args (id = \"u\") outcome db.down()]\n      zzz => [exchange args (id = \"u\") outcome db.down()]\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "warning", "script zzz never runs") {
		t.Fatalf("expected dead-script warning, got %v", diags)
	}
}

func TestDiagnoseBadStub(t *testing.T) {
	bad := strings.Replace(lspAuth, "down => [exchange args (id = \"u\") outcome db.down()]", "down => [exchange args (id = \"u\") outcome db.bogus()]", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "stub db.bogus not in db__get emits") {
		t.Fatalf("expected bad-stub error, got %v", diags)
	}
}

// Exchanges prove "this request received this permitted response"
// (v12): a wrong request value, a wrong request name, or a row
// without an exchange all fail loudly.
func TestExchangeArgMismatch(t *testing.T) {
	bad := strings.Replace(lspAuth, `ok => [exchange args (id = "u") outcome Ok(id = "u")]`, `ok => [exchange args (id = "x") outcome Ok(id = "u")]`, 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "arg id mismatch") {
		t.Fatalf("expected arg-mismatch failure, got %v", diags)
	}
}

func TestExchangeArgNameMissing(t *testing.T) {
	bad := strings.Replace(lspAuth, `ok => [exchange args (id = "u") outcome Ok(id = "u")]`, `ok => [exchange args (idd = "u") outcome Ok(id = "u")]`, 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "but the call supplies none") {
		t.Fatalf("expected missing-arg failure, got %v", diags)
	}
}

func TestExchangeRowRequired(t *testing.T) {
	bad := strings.Replace(lspAuth, `ok => [exchange args (id = "u") outcome Ok(id = "u")]`, `ok => [Ok(id = "u")]`, 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "script row must be an exchange") {
		t.Fatalf("expected exchange-required error, got %v", diags)
	}
	var coded bool
	for _, d := range diags {
		if d.Code == CodeNoExchange {
			coded = true
		}
	}
	if !coded {
		t.Fatalf("expected AIL3109 code, got %v", diags)
	}
}

func TestDiagnoseForeignRaise(t *testing.T) {
	bad := strings.Replace(lspAuth, "on db.down _ => auth.bad()", "on db.down _ => db.down()", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "not in its emits") {
		t.Fatalf("expected foreign-raise error, got %v", diags)
	}
}

// Unrealized emits entries are a conservative upper bound, not an
// error (v12) — but every entry must name a declared error.
func TestDiagnoseUpperBoundEmits(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": lspAuth})
	if diags := diagnose(dir, "auth.ail", lspAuth); len(diags) != 0 {
		t.Fatalf("expected no diagnostics for a realized subset, got %v", diags)
	}
}

func TestDiagnoseUnknownEmitsEntry(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(id: str) -> Auth__S rev 1\n  emits [auth.bad]", "fn auth__go(id: str) -> Auth__S rev 1\n  emits [auth.bad, auth.stale]", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "declares unknown error kind auth.stale in emits") {
		t.Fatalf("expected unknown-emits-entry error, got %v", diags)
	}
}

func TestDiagnoseUnknownErrorKind(t *testing.T) {
	bad := strings.Replace(lspAuth, "on db.down _ => auth.bad()", "on db.down _ => auth.bogus()", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "unknown error kind auth.bogus") {
		t.Fatalf("expected unknown-kind error, got %v", diags)
	}
}

func TestDiagnoseProvidesMismatch(t *testing.T) {
	undeclared := strings.Replace(lspAuth, "provides [auth__go, Auth__S]", "provides [auth__go, Auth__S, auth__ghost]", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": undeclared})
	if diags := diagnose(dir, "auth.ail", undeclared); !hasDiag(diags, "error", "never defines it") {
		t.Fatalf("expected provides-ghost error, got %v", diags)
	}
	undefined := strings.Replace(lspAuth, "provides [auth__go, Auth__S]", "provides [auth__go]", 1)
	dir = writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": undefined})
	if diags := diagnose(dir, "auth.ail", undefined); !hasDiag(diags, "error", "missing from provides") {
		t.Fatalf("expected unprovided-def error, got %v", diags)
	}
}

func TestDiagnoseUnusedUses(t *testing.T) {
	db := strings.Replace(lspDB, "  provides [db__get, Db__U]\n", "  provides [db__get, db__ping, Db__U]\n", 1)
	db += `
fn db__ping() -> Db__U rev 1
  emits []
  tests
    ok() => Ok(id = "u")
=
  Ok(id = "u")
`
	auth := strings.Replace(lspAuth, "uses [db__get@1]", "uses [db__get@1, db__ping@1]", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": db, "auth.ail": auth})
	diags := diagnose(dir, "auth.ail", auth)
	if !hasDiag(diags, "warning", "never calls it") {
		t.Fatalf("expected unused-uses warning, got %v", diags)
	}
}

func TestDiagnoseDupTests(t *testing.T) {
	bad := strings.Replace(lspAuth, "    down(id = \"u\") => auth.bad()\n", "    down(id = \"u\") => auth.bad()\n    ok(id = \"u\") => Ok(id = \"u\")\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "duplicate test ok") {
		t.Fatalf("expected dup-test error, got %v", diags)
	}
}

func TestDiagnoseBadTestArgs(t *testing.T) {
	bad := strings.Replace(lspAuth, "ok(id = \"u\") => Ok(id = \"u\")", "ok(bogus = \"u\") => Ok(id = \"u\")", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "unknown arg bogus") {
		t.Fatalf("expected unknown-arg error, got %v", diags)
	}
	if !hasDiag(diags, "error", "missing arg id") {
		t.Fatalf("expected missing-arg error, got %v", diags)
	}
}

func TestDiagnoseMultiWorld(t *testing.T) {
	bad := strings.Replace(lspAuth, "uses [db__get@1]", "uses [db__get, db__nope@1]", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "must pin a rev") {
		t.Fatalf("expected unpinned-uses error, got %v", diags)
	}
	if !hasDiag(diags, "error", "resolves nowhere") {
		t.Fatalf("expected nowhere-uses error, got %v", diags)
	}
}

func TestDiagnoseCallOutsideScrutinee(t *testing.T) {
	bad := strings.Replace(lspAuth, "    on Ok u => Ok(id = u.id)\n", "    on Ok u => call db__get(id)\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "outside a match scrutinee") {
		t.Fatalf("expected outside-scrutinee error, got %v", diags)
	}
}

func TestDiagnoseSharedErrorAcrossSiblings(t *testing.T) {
	other := strings.Replace(lspAuth, "fn auth__go(", "fn auth__second(", 1)
	other = strings.Replace(other, "provides [auth__go, Auth__S]", "provides [auth__second, Auth__S]", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": lspAuth, "other.ail": other})
	diags := diagnose(dir, "auth.ail", lspAuth)
	if hasDiag(diags, "error", "double definition") {
		t.Fatalf("sibling sharing error auth.bad must stay silent, got %v", diags)
	}
}

func TestDiagnoseSiblingFnDup(t *testing.T) {
	dup := strings.Replace(lspDB, "fn db__get(", "fn auth__go(", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": dup, "auth.ail": lspAuth})
	diags := diagnose(dir, "auth.ail", lspAuth)
	if !hasDiag(diags, "error", "sibling db.ail also defines auth__go") {
		t.Fatalf("expected sibling-dup error, got %v", diags)
	}
}

func findDiag(diags []Diag, sub string) *Diag {
	for i, d := range diags {
		if strings.Contains(d.Msg, sub) {
			return &diags[i]
		}
	}
	return nil
}

// expectSpan is the independent oracle: plain substring search on the
// fixture line, never tokenSpan. Fixtures are ASCII so byte index equals
// the UTF-16 column the editor receives.
func expectSpan(t *testing.T, body, token string, line int) (int, int) {
	t.Helper()
	lines := strings.Split(body, "\n")
	if line < 1 || line > len(lines) {
		t.Fatalf("line %d out of range", line)
	}
	i := strings.Index(lines[line-1], token)
	if i < 0 {
		t.Fatalf("token %q not on line %d: %q", token, line, lines[line-1])
	}
	return i, i + len(token)
}

func expectLine(t *testing.T, body, sub string) int {
	t.Helper()
	for n, line := range strings.Split(body, "\n") {
		if strings.Contains(line, sub) {
			return n + 1
		}
	}
	t.Fatalf("no line contains %q", sub)
	return 0
}

func checkSpan(t *testing.T, body string, diags []Diag, sub, token string, line int) {
	t.Helper()
	d := findDiag(diags, sub)
	if d == nil {
		t.Fatalf("no diagnostic containing %q in %v", sub, diags)
	}
	wantS, wantE := expectSpan(t, body, token, line)
	if d.Line != line || d.Start != wantS || d.End != wantE {
		t.Fatalf("%q: got L%d [%d,%d), want L%d [%d,%d)", sub, d.Line, d.Start, d.End, line, wantS, wantE)
	}
}

func TestSpanNaming(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(", "fn go(", 1)
	bad = strings.Replace(bad, "provides [auth__go, Auth__S]", "provides [go, Auth__S]", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "must match domain__verb", "go", expectLine(t, bad, "fn go("))
}

func TestSpanUnknownCall(t *testing.T) {
	bad := strings.Replace(lspAuth, "match call db__get(id)", "match call db__nope(id)", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "calls unknown function", "db__nope", expectLine(t, bad, "match call db__nope"))
}

func TestSpanDanglingTest(t *testing.T) {
	bad := strings.Replace(lspAuth, "    down(id = \"u\") => auth.bad()\n", "    down(id = \"u\") => auth.bad()\n    extra(id = \"u\") => auth.bad\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "has no script at the call", "extra", expectLine(t, bad, "extra(id = "))
}

func TestSpanDeadScript(t *testing.T) {
	bad := strings.Replace(lspAuth, "      down => [exchange args (id = \"u\") outcome db.down()]\n", "      down => [exchange args (id = \"u\") outcome db.down()]\n      zzz => [exchange args (id = \"u\") outcome db.down()]\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "script zzz never runs", "zzz", expectLine(t, bad, "zzz =>"))
}

func TestSpanUnknownEmitsEntry(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(id: str) -> Auth__S rev 1\n  emits [auth.bad]", "fn auth__go(id: str) -> Auth__S rev 1\n  emits [auth.bad, auth.stale]", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "declares unknown error kind auth.stale in emits", "auth.stale", expectLine(t, bad, "emits [auth.bad, auth.stale]"))
}

func TestSpanForeignRaise(t *testing.T) {
	bad := strings.Replace(lspAuth, "on db.down _ => auth.bad()", "on db.down _ => db.down()", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	line := expectLine(t, bad, "=> db.down()")
	d := findDiag(diags, "not in its emits")
	if d == nil {
		t.Fatalf("no foreign-raise diagnostic in %v", diags)
	}
	ln := strings.Split(bad, "\n")[line-1]
	i := strings.Index(ln, "db.down(")
	if d.Line != line || d.Start != i || d.End != i+len("db.down") {
		t.Fatalf("got L%d [%d,%d), want L%d [%d,%d)", d.Line, d.Start, d.End, line, i, i+len("db.down"))
	}
}

func TestSpanUsesEntry(t *testing.T) {
	bad := strings.Replace(lspAuth, "uses [db__get@1]", "uses [db__get, db__nope@1]", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "must pin a rev", "db__get", expectLine(t, bad, "uses [db__get,"))
}

func TestSpanMissingArm(t *testing.T) {
	bad := strings.Replace(lspAuth, "    on db.down _ => auth.bad()\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "non-exhaustive match, missing", "match", expectLine(t, bad, "match call db__get"))
}

func TestSpanStaleArm(t *testing.T) {
	bad := strings.Replace(lspAuth, "    on db.down _ => auth.bad()\n", "    on db.down _ => auth.bad()\n    on db.bogus _ => auth.bad()\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "stale match arm", "db.bogus", expectLine(t, bad, "on db.bogus _"))
}

func TestSpanFailingTest(t *testing.T) {
	bad := strings.Replace(lspAuth, "down(id = \"u\") => auth.bad()", "down(id = \"u\") => Ok(id = \"u\")", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "test down fails", "down", expectLine(t, bad, "down(id = \"u\") => Ok"))
}

func TestSpanUnusedUses(t *testing.T) {
	db := strings.Replace(lspDB, "  provides [db__get, Db__U]\n", "  provides [db__get, db__ping, Db__U]\n", 1)
	db += `
fn db__ping() -> Db__U rev 1
  emits []
  tests
    ok() => Ok(id = "u")
=
  Ok(id = "u")
`
	auth := strings.Replace(lspAuth, "uses [db__get@1]", "uses [db__get@1, db__ping@1]", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": db, "auth.ail": auth})
	diags := diagnose(dir, "auth.ail", auth)
	checkSpan(t, auth, diags, "never calls it", "db__ping@1", expectLine(t, auth, "uses [db__get@1, db__ping@1]"))
}

func TestSpanFallbackWholeLine(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(id: str) -> Auth__S rev 1", "fn auth__go(id: str) -> Auth__S", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	d := findDiag(diags, "missing rev")
	if d == nil {
		t.Fatalf("no missing-rev diagnostic in %v", diags)
	}
	if d.Start != 0 || d.End != 0 {
		t.Fatalf("parser errors must stay whole-line, got [%d,%d)", d.Start, d.End)
	}
}

func TestU16Len(t *testing.T) {
	if u16len("abc") != 3 || u16len("aé") != 2 || u16len("a\U0001D11E") != 3 {
		t.Fatalf("bad UTF-16 counts: %d %d %d", u16len("abc"), u16len("aé"), u16len("a\U0001D11E"))
	}
}

func TestDiagnoseUntakenArm(t *testing.T) {
	bad := strings.Replace(lspAuth, "    down(id = \"u\") => auth.bad()\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "no test takes", "db.down", expectLine(t, bad, "on db.down _"))
}

func TestDiagnoseShadowedArm(t *testing.T) {
	bad := strings.Replace(lspAuth, "    on db.down _ => auth.bad()\n", "    on db.down _ => auth.bad()\n    on db.down _ => auth.bad()\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	first := expectLine(t, bad, "on db.down _")
	d := findDiag(diags, "no test takes")
	if d == nil {
		t.Fatalf("expected shadowed-arm error, got %v", diags)
	}
	if d.Line != first+1 {
		t.Fatalf("shadowed arm should be the second copy at L%d, got L%d", first+1, d.Line)
	}
}

func TestCoverageGreenOnly(t *testing.T) {
	bad := strings.Replace(lspAuth, "    down(id = \"u\") => auth.bad()\n", "    down(id = \"u\") => auth.bad()\n    extra(id = \"u\") => auth.bad\n", 1)
	bad = strings.Replace(bad, "    on db.down _ => auth.bad()\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	if !hasDiag(diags, "error", "no script for call") {
		t.Fatalf("expected the failing test, got %v", diags)
	}
	if hasDiag(diags, "error", "no test takes") {
		t.Fatalf("coverage must wait for green tests, got %v", diags)
	}
}

// Type discipline (v0.4): float ban, dec exactness, brands, seal, extern.

func TestDiagnoseFloatLiteral(t *testing.T) {
	bad := strings.Replace(lspAuth, "ok(id = \"u\") => Ok(id = \"u\")", "ok(id = \"u\") => Ok(id = 0.5)", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "has no spelling", "0.5", expectLine(t, bad, "Ok(id = 0.5)"))
}

const typeShop = `mod shop
  provides [shop__price, Shop__Item]
  uses []
  emits []

type Shop__Item rev 1 (
  price: dec
)

fn shop__price() -> Shop__Item rev 1
  emits []
  tests
    half() => Ok(price = d"1.5")
=
  Ok(price = d"1.50")
`

func TestDiagnoseDecExact(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"shop.ail": typeShop})
	if diags := diagnose(dir, "shop.ail", typeShop); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

const typeLeak = `mod m
  provides [m__leak, M__B, M__Out]
  uses []
  emits []

brand M__B is str rev 1

type M__Out rev 1 (
  echo: str
)

fn m__leak(pw: M__B) -> M__Out rev 1
  emits []
  tests
    t(pw = seal M__B("s")) => Ok(echo = "s")
=
  Ok(echo = pw)
`

func TestDiagnoseBrandSink(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.ail": typeLeak})
	diags := diagnose(dir, "m.ail", typeLeak)
	checkSpan(t, typeLeak, diags, "got M__B, want str", "echo", expectLine(t, typeLeak, "Ok(echo = pw)"))
}

const typeCmp = `mod m
  provides [m__cmp, M__A, M__B, M__Out]
  uses []
  emits []

brand M__A is str rev 1
brand M__B is str rev 1

type M__Out rev 1 (
  echo: str
)

fn m__cmp(a: M__A, b: M__B) -> M__Out rev 1
  emits []
  tests
    t(a = seal M__A("x"), b = seal M__B("y")) => Ok(echo = "n")
=
  match a == b
    true => Ok(echo = "y")
    false => Ok(echo = "n")
`

func TestDiagnoseCrossBrandCompare(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.ail": typeCmp})
	diags := diagnose(dir, "m.ail", typeCmp)
	checkSpan(t, typeCmp, diags, "cannot compare M__A with M__B", "==", expectLine(t, typeCmp, "match a == b"))
}

const typeSeal = `mod m
  provides [m__seal, M__B, M__Out]
  uses []
  emits []

brand M__B is str rev 1

type M__Out rev 1 (
  echo: str
)

fn m__seal(pw: str) -> M__Out rev 1
  emits []
  tests
    t(pw = "s") => Ok(echo = "s")
=
  Ok(echo = seal M__B(pw))
`

func TestDiagnoseSealLiteral(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.ail": typeSeal})
	diags := diagnose(dir, "m.ail", typeSeal)
	checkSpan(t, typeSeal, diags, "takes a string literal", "M__B", expectLine(t, typeSeal, "seal M__B(pw)"))
}

func TestDiagnoseUnknownType(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(id: str)", "fn auth__go(id: Nope)", 1)
	dir := writeLSPDir(t, map[string]string{"db.ail": lspDB, "auth.ail": bad})
	diags := diagnose(dir, "auth.ail", bad)
	checkSpan(t, bad, diags, "unknown type Nope in param id", "id", expectLine(t, bad, "fn auth__go(id: Nope)"))
}

const typeExtern = `mod m
  provides [m__check, M__Pw, M__Verdict, M__Out, m__use]
  uses []
  emits [m.nope]

error m.nope()

brand M__Pw is str rev 1

type M__Verdict rev 1 (
)

type M__Out rev 1 (
  ok: bool
)

extern m__use(pw: M__Pw) -> M__Verdict rev 1
  emits [m.nope]

fn m__check(pw: M__Pw) -> M__Out rev 1
  emits [m.nope]
  tests
    yes(pw = seal M__Pw("s")) => Ok(ok = true)
    no(pw = seal M__Pw("s")) => m.nope()
=
  match call m__use(pw)
    given
      yes => [exchange args (pw = seal M__Pw("s")) outcome Ok()]
      no => [exchange args (pw = seal M__Pw("s")) outcome m.nope()]
    on m.nope _ => m.nope()
    on Ok v => Ok(ok = true)
`

func TestDiagnoseExternClean(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.ail": typeExtern})
	if diags := diagnose(dir, "m.ail", typeExtern); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestDiagnoseExternBoolRet(t *testing.T) {
	bad := strings.Replace(typeExtern, "extern m__use(pw: M__Pw) -> M__Verdict rev 1", "extern m__use(pw: M__Pw) -> bool rev 1", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": bad})
	diags := diagnose(dir, "m.ail", bad)
	if !hasDiag(diags, "error", "extern m__use returns bool") {
		t.Fatalf("expected extern record-return error, got %v", diags)
	}
}

const typeForeign = `mod o
  provides [o__go, O__Out]
  uses []
  emits []

type O__Out rev 1 (
  id: str
)

fn o__go(id: str) -> O__Out rev 1
  emits []
  tests
    t(id = "u") => Ok(id = "u")
=
  match call m__use(id)
    given
      t => [exchange args (pw = "u") outcome Ok()]
    on Ok v => Ok(id = id)
`

func TestDiagnoseForeignExtern(t *testing.T) {
	files := map[string]string{"m.ail": typeExtern, "o.ail": typeForeign}
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, "o.ail", typeForeign)
	if !hasDiag(diags, "error", "declare your own extern") {
		t.Fatalf("expected foreign-extern error, got %v", diags)
	}
}

// Arithmetic (v0.6): precedence, left associativity, exact dec, mismatch codes.

const typeArith = `mod m
  provides [m__calc, m__prec, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  n: int
)

fn m__calc(a: int, b: int, c: int) -> M__Out rev 1
  emits []
  tests
    t(a = 10, b = 3, c = 2) => Ok(n = 5)
=
  Ok(n = a - b - c)

fn m__prec(a: int, b: int, c: int) -> M__Out rev 1
  emits []
  tests
    t(a = 2, b = 3, c = 4) => Ok(n = 14)
=
  Ok(n = a + b * c)
`

func TestDiagnoseArithClean(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.ail": typeArith})
	if diags := diagnose(dir, "m.ail", typeArith); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

const typeExact = `mod m
  provides [m__exact, M__Bit]
  uses []
  emits []

type M__Bit rev 1 (
  yes: bool
)

fn m__exact(a: dec, b: dec, c: dec) -> M__Bit rev 1
  emits []
  tests
    third(a = d"0.1", b = d"0.2", c = d"0.3") => Ok(yes = true)
    fourth(a = d"0.1", b = d"0.2", c = d"0.4") => Ok(yes = false)
=
  match a + b == c
    true => Ok(yes = true)
    false => Ok(yes = false)
`

func TestDiagnoseArithDecExact(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.ail": typeExact})
	if diags := diagnose(dir, "m.ail", typeExact); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestDiagnoseArithMixed(t *testing.T) {
	bad := strings.Replace(typeArith, "Ok(n = a - b - c)", `Ok(n = a + d"1.5")`, 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": bad})
	diags := diagnose(dir, "m.ail", bad)
	checkSpan(t, bad, diags, "cannot add int with dec", "+", expectLine(t, bad, `Ok(n = a + d"1.5")`))
}

func TestDiagnoseArithStr(t *testing.T) {
	// v16: + concatenates strings, so the refused string operation
	// is now -. The contract under test is unchanged: strings do
	// no arithmetic besides explicit construction.
	bad := strings.Replace(typeArith, "(a: int, b: int, c: int) -> M__Out rev 1\n  emits []\n  tests\n    t(a = 10, b = 3, c = 2) => Ok(n = 5)\n=\n  Ok(n = a - b - c)",
		"(a: str, b: str, c: int) -> M__Out rev 1\n  emits []\n  tests\n    t(a = \"x\", b = \"y\", c = 2) => Ok(n = 5)\n=\n  Ok(n = a - b)", 1)
	dir := writeLSPDir(t, map[string]string{"m.ail": bad})
	diags := diagnose(dir, "m.ail", bad)
	checkSpan(t, bad, diags, "cannot subtract str with str", "-", expectLine(t, bad, "Ok(n = a - b)"))
}

const typeConcat = `mod m
  provides [m__cat, M__Cat]
  uses []
  emits []

type M__Cat rev 1 (
  s: str
)

fn m__cat(left: str, right: str) -> M__Cat rev 1
  emits []
  tests
    basic(left = "x", right = "y") => Ok(s = "xy")
    empty(left = "", right = "y") => Ok(s = "y")
=
  Ok(s = left + right)
`

func TestDiagnoseStrConcatClean(t *testing.T) {
	// v16: str + str is explicit construction, not a mismatch.
	dir := writeLSPDir(t, map[string]string{"m.ail": typeConcat})
	if diags := diagnose(dir, "m.ail", typeConcat); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

const typeIntDiv = `mod m
  provides [m__div, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  q: int
)

fn m__div(a: int, b: int) -> M__Out rev 1
  emits []
  tests
    t(a = 7, b = 3) => Ok(q = 2)
=
  Ok(q = a / b)
`

func TestDiagnoseDivisionDeferred(t *testing.T) {
	// v17 discharges the v06 deferral for integers: / is exact
	// Euclidean division with a loud zero divisor. What stays
	// deferred is decimal division (AIL6005, pinned separately).
	dir := writeLSPDir(t, map[string]string{"m.ail": typeIntDiv})
	if diags := diagnose(dir, "m.ail", typeIntDiv); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
