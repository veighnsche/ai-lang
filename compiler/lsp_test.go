package main

import (
	"encoding/json"
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
    ok("u") => Ok(id = "u")
    other("x") => db.down()
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
    ok("u") => Ok(id = "u")
    down("u") => auth.bad()
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
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": lspAuth})
	if diags := diagnose(dir, "auth.can", lspAuth); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestDiagnoseMissingTests(t *testing.T) {
	noTests := strings.Replace(lspAuth, "  tests\n    ok(\"u\") => Ok(id = \"u\")\n    down(\"u\") => auth.bad()\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": noTests})
	diags := diagnose(dir, "auth.can", noTests)
	if !hasDiag(diags, "error", "ships no tests") {
		t.Fatalf("expected missing-tests error, got %v", diags)
	}
	if !hasErrors(diags) {
		t.Fatalf("expected missing-tests to fail the shared build gate, got %v", diags)
	}
}

func TestDiagnoseBadName(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(", "fn go(", 1)
	bad = strings.Replace(bad, "provides [auth__go]", "provides [go]", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "must match domain__verb") {
		t.Fatalf("expected naming error, got %v", diags)
	}
}

func TestDiagnoseMissingRev(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(id: str) -> Auth__S rev 1", "fn auth__go(id: str) -> Auth__S", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "missing rev") {
		t.Fatalf("expected missing-rev error, got %v", diags)
	}
}

// Complete error expectations compare kind and payload (a12): a
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
    hit("u") => Ok(id = "u")
    miss("x") => m.bad(id = "x")
  match id
    "u" => Ok(id = id)
    _ => m.bad(id = id)
`

func TestExpectErrorPayloadMismatch(t *testing.T) {
	bad := strings.Replace(expectErrFixture, `miss("x") => m.bad(id = "x")`, `miss("x") => m.bad(id = "y")`, 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	if !hasDiag(diags, "error", "error payload mismatch") {
		t.Fatalf("expected payload-mismatch failure, got %v", diags)
	}
}

func TestExpectBareErrorKindRefused(t *testing.T) {
	bad := strings.Replace(expectErrFixture, `miss("x") => m.bad(id = "x")`, `miss("x") => m.bad`, 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
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
		t.Fatalf("expected CAN3204 code, got %v", diags)
	}
}

func TestDiagnoseFailingTest(t *testing.T) {
	bad := strings.Replace(lspAuth, "down(\"u\") => auth.bad()", "down(\"u\") => Ok(id = \"u\")", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "test down fails") {
		t.Fatalf("expected failing-test error, got %v", diags)
	}
}

func TestDiagnoseMissingArm(t *testing.T) {
	bad := strings.Replace(lspAuth, "    on db.down _ => auth.bad()\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "non-exhaustive") {
		t.Fatalf("expected exhaustiveness error, got %v", diags)
	}
}

func TestDiagnoseUnknownCallee(t *testing.T) {
	bad := strings.Replace(lspAuth, "match call db__get(id)", "match call db__nope(id)", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "calls unknown function db__nope") {
		t.Fatalf("expected unknown-callee error, got %v", diags)
	}
}

func TestDiagnoseCallNotInUses(t *testing.T) {
	bad := strings.Replace(lspAuth, "uses [db__get@1]", "uses []", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "not in uses") {
		t.Fatalf("expected not-in-uses error, got %v", diags)
	}
}

func TestDiagnoseMissingGiven(t *testing.T) {
	bad := strings.Replace(lspAuth, "    given\n      ok => [exchange args (id = \"u\") outcome Ok(id = \"u\")]\n      down => [exchange args (id = \"u\") outcome db.down()]\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "has no given table") {
		t.Fatalf("expected missing-given error, got %v", diags)
	}
}

func TestDiagnoseOmissionBackstop(t *testing.T) {
	bad := strings.Replace(lspAuth, "    down(\"u\") => auth.bad()\n", "    down(\"u\") => auth.bad()\n    extra(\"u\") => auth.bad()\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if hasDiag(diags, "error", "has no script") {
		t.Fatalf("omission draws no static diagnostic, got %v", diags)
	}
	if !hasDiag(diags, "error", "test extra fails") {
		t.Fatalf("expected execution backstop for extra, got %v", diags)
	}
}

func TestDiagnoseDeadScript(t *testing.T) {
	bad := strings.Replace(lspAuth, "      down => [exchange args (id = \"u\") outcome db.down()]\n", "      down => [exchange args (id = \"u\") outcome db.down()]\n      zzz => [exchange args (id = \"u\") outcome db.down()]\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "warning", "script zzz never runs") {
		t.Fatalf("expected dead-script warning, got %v", diags)
	}
}

func TestDiagnoseBadStub(t *testing.T) {
	bad := strings.Replace(lspAuth, "down => [exchange args (id = \"u\") outcome db.down()]", "down => [exchange args (id = \"u\") outcome db.bogus()]", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "stub db.bogus not in db__get emits") {
		t.Fatalf("expected bad-stub error, got %v", diags)
	}
}

// Exchanges prove "this request received this permitted response"
// (a12): a wrong request value, a wrong request name, or a row
// without an exchange all fail loudly.
func TestExchangeArgMismatch(t *testing.T) {
	bad := strings.Replace(lspAuth, `ok => [exchange args (id = "u") outcome Ok(id = "u")]`, `ok => [exchange args (id = "x") outcome db.down()]`, 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "arg id mismatch") {
		t.Fatalf("expected arg-mismatch failure, got %v", diags)
	}
}

func TestExchangeArgNameMissing(t *testing.T) {
	bad := strings.Replace(lspAuth, `ok => [exchange args (id = "u") outcome Ok(id = "u")]`, `ok => [exchange args (idd = "u") outcome Ok(id = "u")]`, 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "but the call supplies none") {
		t.Fatalf("expected missing-arg failure, got %v", diags)
	}
}

func TestExchangeRowRequired(t *testing.T) {
	bad := strings.Replace(lspAuth, `ok => [exchange args (id = "u") outcome Ok(id = "u")]`, `ok => [Ok(id = "u")]`, 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
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
		t.Fatalf("expected CAN3109 code, got %v", diags)
	}
}

func TestDiagnoseForeignRaise(t *testing.T) {
	bad := strings.Replace(lspAuth, "on db.down _ => auth.bad()", "on db.down _ => db.down()", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "not in its emits") {
		t.Fatalf("expected foreign-raise error, got %v", diags)
	}
}

// Unrealized emits entries are a conservative upper bound, not an
// error (a12) — but every entry must name a declared error.
func TestDiagnoseUpperBoundEmits(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": lspAuth})
	if diags := diagnose(dir, "auth.can", lspAuth); len(diags) != 0 {
		t.Fatalf("expected no diagnostics for a realized subset, got %v", diags)
	}
}

func TestDiagnoseUnknownEmitsEntry(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(id: str) -> Auth__S rev 1\n  emits [auth.bad]", "fn auth__go(id: str) -> Auth__S rev 1\n  emits [auth.bad, auth.stale]", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "declares unknown error kind auth.stale in emits") {
		t.Fatalf("expected unknown-emits-entry error, got %v", diags)
	}
}

func TestDiagnoseUnknownErrorKind(t *testing.T) {
	bad := strings.Replace(lspAuth, "on db.down _ => auth.bad()", "on db.down _ => auth.bogus()", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "unknown error kind auth.bogus") {
		t.Fatalf("expected unknown-kind error, got %v", diags)
	}
}

func TestDiagnoseProvidesMismatch(t *testing.T) {
	undeclared := strings.Replace(lspAuth, "provides [auth__go, Auth__S]", "provides [auth__go, Auth__S, auth__ghost]", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": undeclared})
	if diags := diagnose(dir, "auth.can", undeclared); !hasDiag(diags, "error", "never defines it") {
		t.Fatalf("expected provides-ghost error, got %v", diags)
	}
	undefined := strings.Replace(lspAuth, "provides [auth__go, Auth__S]", "provides [auth__go]", 1)
	dir = writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": undefined})
	if diags := diagnose(dir, "auth.can", undefined); !hasDiag(diags, "error", "missing from provides") {
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
  Ok(id = "u")
`
	auth := strings.Replace(lspAuth, "uses [db__get@1]", "uses [db__get@1, db__ping@1]", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": db, "auth.can": auth})
	diags := diagnose(dir, "auth.can", auth)
	if !hasDiag(diags, "warning", "never calls it") {
		t.Fatalf("expected unused-uses warning, got %v", diags)
	}
}

func TestDiagnoseDupTests(t *testing.T) {
	bad := strings.Replace(lspAuth, "    down(\"u\") => auth.bad()\n", "    down(\"u\") => auth.bad()\n    ok(\"u\") => Ok(id = \"u\")\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "duplicate test ok") {
		t.Fatalf("expected dup-test error, got %v", diags)
	}
}

func TestDiagnoseBadTestArgs(t *testing.T) {
	bad := strings.Replace(lspAuth, "ok(\"u\") => Ok(id = \"u\")", "ok(bogus = \"u\") => Ok(id = \"u\")", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "unknown arg bogus") {
		t.Fatalf("expected unknown-arg error, got %v", diags)
	}
	if !hasDiag(diags, "error", "missing arg id") {
		t.Fatalf("expected missing-arg error, got %v", diags)
	}
}

func TestDiagnoseMultiWorld(t *testing.T) {
	bad := strings.Replace(lspAuth, "uses [db__get@1]", "uses [db__get, db__nope@1]", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "must pin a rev") {
		t.Fatalf("expected unpinned-uses error, got %v", diags)
	}
	if !hasDiag(diags, "error", "resolves nowhere") {
		t.Fatalf("expected nowhere-uses error, got %v", diags)
	}
}

func TestDiagnoseCallOutsideScrutinee(t *testing.T) {
	bad := strings.Replace(lspAuth, "    on Ok u => Ok(id = u.id)\n", "    on Ok u => call db__get(id)\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "outside a match scrutinee") {
		t.Fatalf("expected outside-scrutinee error, got %v", diags)
	}
}

func TestDiagnoseSharedErrorAcrossSiblings(t *testing.T) {
	other := strings.Replace(lspAuth, "fn auth__go(", "fn auth__second(", 1)
	other = strings.Replace(other, "provides [auth__go, Auth__S]", "provides [auth__second, Auth__S]", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": lspAuth, "other.can": other})
	diags := diagnose(dir, "auth.can", lspAuth)
	if hasDiag(diags, "error", "double definition") {
		t.Fatalf("sibling sharing error auth.bad must stay silent, got %v", diags)
	}
}

func TestDiagnoseSiblingFnDup(t *testing.T) {
	dup := strings.Replace(lspDB, "fn db__get(", "fn auth__go(", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": dup, "auth.can": lspAuth})
	diags := diagnose(dir, "auth.can", lspAuth)
	if !hasDiag(diags, "error", "sibling db.can also defines auth__go") {
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
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	checkSpan(t, bad, diags, "must match domain__verb", "go", expectLine(t, bad, "fn go("))
}

func TestSpanUnknownCall(t *testing.T) {
	bad := strings.Replace(lspAuth, "match call db__get(id)", "match call db__nope(id)", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	checkSpan(t, bad, diags, "calls unknown function", "db__nope", expectLine(t, bad, "match call db__nope"))
}

func TestSpanGivenDashRetired(t *testing.T) {
	bad := strings.Replace(lspAuth, "    down(\"u\") => auth.bad()\n", "    down(\"u\") => auth.bad()\n    down2(\"u\") => auth.bad()\n", 1)
	bad = strings.Replace(bad, "      down => [exchange args (id = \"u\") outcome db.down()]\n", "      down => [exchange args (id = \"u\") outcome db.down()]\n      down2 => -\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	checkSpan(t, bad, diags, "retired `-` spelling", "down2", expectLine(t, bad, "down2 =>"))
}

func TestSpanDeadScript(t *testing.T) {
	bad := strings.Replace(lspAuth, "      down => [exchange args (id = \"u\") outcome db.down()]\n", "      down => [exchange args (id = \"u\") outcome db.down()]\n      zzz => [exchange args (id = \"u\") outcome db.down()]\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	checkSpan(t, bad, diags, "script zzz never runs", "zzz", expectLine(t, bad, "zzz =>"))
}

func TestSpanUnknownEmitsEntry(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(id: str) -> Auth__S rev 1\n  emits [auth.bad]", "fn auth__go(id: str) -> Auth__S rev 1\n  emits [auth.bad, auth.stale]", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	checkSpan(t, bad, diags, "declares unknown error kind auth.stale in emits", "auth.stale", expectLine(t, bad, "emits [auth.bad, auth.stale]"))
}

func TestSpanForeignRaise(t *testing.T) {
	bad := strings.Replace(lspAuth, "on db.down _ => auth.bad()", "on db.down _ => db.down()", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
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
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	checkSpan(t, bad, diags, "must pin a rev", "db__get", expectLine(t, bad, "uses [db__get,"))
}

func TestSpanMissingArm(t *testing.T) {
	bad := strings.Replace(lspAuth, "    on db.down _ => auth.bad()\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	checkSpan(t, bad, diags, "non-exhaustive match, missing", "match", expectLine(t, bad, "match call db__get"))
}

func TestSpanStaleArm(t *testing.T) {
	bad := strings.Replace(lspAuth, "    on db.down _ => auth.bad()\n", "    on db.down _ => auth.bad()\n    on db.bogus _ => auth.bad()\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	checkSpan(t, bad, diags, "stale match arm", "db.bogus", expectLine(t, bad, "on db.bogus _"))
}

func TestSpanFailingTest(t *testing.T) {
	bad := strings.Replace(lspAuth, "down(\"u\") => auth.bad()", "down(\"u\") => Ok(id = \"u\")", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	checkSpan(t, bad, diags, "test down fails", "down", expectLine(t, bad, "down(\"u\") => Ok"))
}

func TestSpanUnusedUses(t *testing.T) {
	db := strings.Replace(lspDB, "  provides [db__get, Db__U]\n", "  provides [db__get, db__ping, Db__U]\n", 1)
	db += `
fn db__ping() -> Db__U rev 1
  emits []
  tests
    ok() => Ok(id = "u")
  Ok(id = "u")
`
	auth := strings.Replace(lspAuth, "uses [db__get@1]", "uses [db__get@1, db__ping@1]", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": db, "auth.can": auth})
	diags := diagnose(dir, "auth.can", auth)
	checkSpan(t, auth, diags, "never calls it", "db__ping@1", expectLine(t, auth, "uses [db__get@1, db__ping@1]"))
}

func TestSpanFallbackWholeLine(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(id: str) -> Auth__S rev 1", "fn auth__go(id: str) -> Auth__S", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
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
	bad := strings.Replace(lspAuth, "    down(\"u\") => auth.bad()\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	checkSpan(t, bad, diags, "no test takes", "db.down", expectLine(t, bad, "on db.down _"))
}

func TestDiagnoseShadowedArm(t *testing.T) {
	bad := strings.Replace(lspAuth, "    on db.down _ => auth.bad()\n", "    on db.down _ => auth.bad()\n    on db.down _ => auth.bad()\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	first := expectLine(t, bad, "on db.down _")
	d := findDiag(diags, "no test takes")
	if d == nil {
		t.Fatalf("expected shadowed-arm error, got %v", diags)
	}
	if d.Line != first+1 {
		t.Fatalf("shadowed arm should be the second copy at L%d, got L%d", first+1, d.Line)
	}
}

// relayMath exercises the identity-relay certificate: no test takes the
// relay arm (pos recurses 3->2->1->0 through the Ok arm only, neg and
// zero never enter the call match), so silence proves the certificate.
const relayMath = `mod math
  provides [math__sum_to, Int__Value]
  uses []
  emits [math.negative_input]

error math.negative_input(value: int)

type Int__Value rev 1 (
  value: int
)

fn math__sum_to(n: int) -> Int__Value rev 1
  decreases n
  emits [math.negative_input]
  tests
    neg(-3) => math.negative_input(value = -3)
    zero(0) => Ok(value = 0)
    pos(3) => Ok(value = 6)
  match n <= 0
    true => match n == 0
      true => Ok(value = 0)
      false => math.negative_input(value = n)
    false => match call math__sum_to(n - 1)
      on math.negative_input err => math.negative_input(value = err.value)
      on Ok r => Ok(value = n + r.value)
`

func TestDiagnoseIdentityRelayCertified(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"math.can": relayMath})
	diags := diagnose(dir, "math.can", relayMath)
	// The relay certificate still holds (no untaken-arm or invalid-
	// relay errors), but the handwritten spelling is a lint error
	// (CAN3415): certified is not idiomatic.
	relay := 0
	for _, d := range diags {
		if d.Code == CodeLintRelay {
			relay++
			continue
		}
		if d.Sev == "error" {
			t.Fatalf("unexpected error beside the relay finding: %v", diags)
		}
	}
	if relay == 0 {
		t.Fatalf("expected a CAN3415 relay finding, got %v", diags)
	}
}

func TestDiagnoseInvalidRelayValue(t *testing.T) {
	bad := strings.Replace(relayMath,
		"on math.negative_input err => math.negative_input(value = err.value)",
		"on math.negative_input err => math.negative_input(value = 0)", 1)
	dir := writeLSPDir(t, map[string]string{"math.can": bad})
	diags := diagnose(dir, "math.can", bad)
	if !hasCode(diags, "CAN4108") {
		t.Fatalf("expected CAN4108 for a changed relay value, got %v", diags)
	}
	if findDiag(diags, "field value is not err.value") == nil {
		t.Fatalf("expected the value reason, got %v", diags)
	}
}

func TestDiagnoseInvalidRelayKind(t *testing.T) {
	bad := strings.Replace(relayMath, "emits [math.negative_input]",
		"emits [math.negative_input, math.other]", 1)
	bad = strings.Replace(bad, "error math.negative_input(value: int)",
		"error math.negative_input(value: int)\n\nerror math.other(value: int)", 1)
	bad = strings.Replace(bad,
		"on math.negative_input err => math.negative_input(value = err.value)",
		"on math.negative_input err => math.other(value = err.value)", 1)
	dir := writeLSPDir(t, map[string]string{"math.can": bad})
	diags := diagnose(dir, "math.can", bad)
	if !hasCode(diags, "CAN4108") {
		t.Fatalf("expected CAN4108 for a kind-swapping relay, got %v", diags)
	}
	if findDiag(diags, "reconstructs math.other instead of math.negative_input") == nil {
		t.Fatalf("expected the kind reason, got %v", diags)
	}
}

const relayMathTwo = `mod math2
  provides [math2__go, Int__Value]
  uses []
  emits [math2.fail]

error math2.fail(a: int, b: int)

type Int__Value rev 1 (
  value: int
)

fn math2__go(n: int) -> Int__Value rev 1
  decreases n
  emits [math2.fail]
  tests
    fail(-1) => math2.fail(a = 1, b = 2)
    zero(0) => Ok(value = 0)
    pos(2) => Ok(value = 3)
  match n <= 0
    true => match n == 0
      true => Ok(value = 0)
      false => math2.fail(a = 1, b = 2)
    false => match call math2__go(n - 1)
      on math2.fail e => math2.fail(a = e.a, b = e.b)
      on Ok r => Ok(value = n + r.value)
`

func TestDiagnoseInvalidRelayDrops(t *testing.T) {
	bad := strings.Replace(relayMathTwo,
		"on math2.fail e => math2.fail(a = e.a, b = e.b)",
		"on math2.fail e => math2.fail(a = e.a)", 1)
	dir := writeLSPDir(t, map[string]string{"math2.can": bad})
	diags := diagnose(dir, "math2.can", bad)
	if !hasCode(diags, "CAN4108") {
		t.Fatalf("expected CAN4108 for a field-dropping relay, got %v", diags)
	}
	if findDiag(diags, "drops field b") == nil {
		t.Fatalf("expected the drops reason, got %v", diags)
	}
}

func TestDiagnoseInvalidRelayUnexpected(t *testing.T) {
	bad := strings.Replace(relayMathTwo,
		"on math2.fail e => math2.fail(a = e.a, b = e.b)",
		"on math2.fail e => math2.fail(a = e.a, b = e.b, z = e.a)", 1)
	dir := writeLSPDir(t, map[string]string{"math2.can": bad})
	diags := diagnose(dir, "math2.can", bad)
	if !hasCode(diags, "CAN4108") {
		t.Fatalf("expected CAN4108 for an extra-field relay, got %v", diags)
	}
	if findDiag(diags, "rebuilds unexpected field z") == nil {
		t.Fatalf("expected the unexpected-field reason, got %v", diags)
	}
}

// A valid identity relay over a local call is exempt from execution
// coverage by structural certificate: relayMathTwo's relay arm never
// executes (fail raises before any call; pos unwinds through Ok), yet
// the fixture must diagnose clean — same shape as the a44 skip-level
// relay arm, which compiles without a dedicated row.
func TestDiagnoseValidRelayExempt(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"math2.can": relayMathTwo})
	diags := diagnose(dir, "math2.can", relayMathTwo)
	// Certified against untaken-arm and invalid-relay errors, but
	// the handwritten spelling is a lint error (CAN3415).
	relay := 0
	for _, d := range diags {
		if d.Code == CodeLintRelay {
			relay++
			continue
		}
		if d.Sev == "error" {
			t.Fatalf("unexpected error beside the relay finding: %v", diags)
		}
	}
	if relay == 0 {
		t.Fatalf("expected a CAN3415 relay finding, got %v", diags)
	}
}

// Row 3 (D1): a shadowed arm is never certified, even relay-shaped.
// The second copy is dead — structural evidence cannot substitute
// for the execution the shadowing removed — so it stays under the
// execution law like every other duplicate (TestDiagnoseShadowedArm).
const relayDup = `mod dup
  provides [dup__go, dup__help, Dup__Out, Dup__Val]
  uses []
  emits [dup.bad]

error dup.bad(value: int)

type Dup__Val rev 1 (
  value: int
)

type Dup__Out rev 1 (
  value: int
)

fn dup__help(x: int) -> Dup__Val rev 1
  emits [dup.bad]
  tests
    go(1) => Ok(value = 1)
  Ok(value = x)

fn dup__go(x: int) -> Dup__Out rev 1
  emits [dup.bad]
  tests
    go(1) => Ok(value = 1)
  match call dup__help(x)
    on Ok r => Ok(value = r.value)
    on dup.bad e => dup.bad(value = e.value)
    on dup.bad e2 => dup.bad(value = e2.value)
`

func TestDiagnoseShadowedRelay(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"dup.can": relayDup})
	diags := diagnose(dir, "dup.can", relayDup)
	if !hasCode(diags, "CAN4107") {
		t.Fatalf("expected CAN4107 for the shadowed relay arm, got %v", diags)
	}
	if findDiag(diags, "no test takes on dup.bad e2") == nil {
		t.Fatalf("expected the duplicate arm named, got %v", diags)
	}
	if hasCode(diags, "CAN4108") {
		t.Fatalf("shadowed arm must not reach the relay check, got %v", diags)
	}
}

// Row 3 (D1): certificates are owner-local. A relay-shaped arm over
// a foreign call must not certify — the untaken arm stays under the
// execution law even when the reconstruction is field-perfect.
const relayForeignLib = `mod lib
  provides [lib__get, Lib__Out]
  uses []
  emits [lib.bad]

error lib.bad(value: int)

type Lib__Out rev 1 (
  value: int
)

fn lib__get() -> Lib__Out rev 1
  emits [lib.bad]
  tests
    go() => Ok(value = 1)
  Ok(value = 1)
`

const relayForeignApp = `mod app
  provides [app__go, App__Out]
  uses [lib__get@1]
  emits [lib.bad]

type App__Out rev 1 (
  value: int
)

fn app__go() -> App__Out rev 1
  emits [lib.bad]
  tests
    go() => Ok(value = 1)
  match call lib__get()
    given
      go => [exchange args () outcome Ok(value = 1)]
    on Ok r => Ok(value = r.value)
    on lib.bad e => lib.bad(value = e.value)
`

func TestDiagnoseForeignRelayUntaken(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"lib.can": relayForeignLib, "app.can": relayForeignApp})
	diags := diagnose(dir, "app.can", relayForeignApp)
	if !hasCode(diags, "CAN4107") {
		t.Fatalf("expected CAN4107 for the unscripted foreign arm, got %v", diags)
	}
	if hasCode(diags, "CAN4108") {
		t.Fatalf("foreign arm must not reach the relay check, got %v", diags)
	}
}

// Row 3 (D1): the catalog records references, not executions. A
// certified-but-unexecuted relay arm manufactures no hit, and a
// declared-but-unrealized error compiles with empty lists — the
// conservative upper bound, pinned against regression either way.
func TestCatalogNoPhantomHits(t *testing.T) {
	const catSrc = `mod cat
  provides [cat__go, cat__help, Cat__Out, Cat__Val]
  uses []
  emits [cat.used, cat.free]

error cat.used(value: int)
error cat.free(value: int)

type Cat__Val rev 1 (
  value: int
)

type Cat__Out rev 1 (
  value: int
)

fn cat__help(x: int) -> Cat__Val rev 1
  emits [cat.used]
  tests
    go(1) => Ok(value = 1)
  Ok(value = x)

fn cat__go(x: int) -> Cat__Out rev 1
  emits [cat.used, cat.free]
  tests
    go(1) => Ok(value = 1)
  match call cat__help(x)
    on Ok r => Ok(value = r.value)
    on cat.used e => cat.used(value = e.value)
`
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "cat.can")
	if err := os.WriteFile(srcPath, []byte(catSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compileEx(out, []string{srcPath}, true); err != nil {
		t.Fatalf("unrealized entries must compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "errors.json"))
	if err != nil {
		t.Fatal(err)
	}
	var entries []catalogEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		t.Fatal(err)
	}
	byKind := map[string]catalogEntry{}
	for _, e := range entries {
		byKind[e.Kind] = e
	}
	used, ok := byKind["cat.used"]
	if !ok {
		t.Fatalf("catalog omits cat.used, got %v", entries)
	}
	if len(used.HitByTests) != 0 {
		t.Fatalf("certified arm must manufacture no hits, got %v", used.HitByTests)
	}
	if len(used.HandledBy) != 1 {
		t.Fatalf("relay arm is a handling site, got %v", used.HandledBy)
	}
	free, ok := byKind["cat.free"]
	if !ok {
		t.Fatalf("catalog omits unrealized cat.free, got %v", entries)
	}
	if len(free.RaisedBy) != 0 || len(free.HitByTests) != 0 {
		t.Fatalf("unrealized entry must stay empty, got %+v", free)
	}
}

func TestCoverageGreenOnly(t *testing.T) {
	bad := strings.Replace(lspAuth, "    down(\"u\") => auth.bad()\n", "    down(\"u\") => auth.bad()\n    extra(\"u\") => auth.bad\n", 1)
	bad = strings.Replace(bad, "    on db.down _ => auth.bad()\n", "", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
	if !hasDiag(diags, "error", "no script for call") {
		t.Fatalf("expected the failing test, got %v", diags)
	}
	if hasDiag(diags, "error", "no test takes") {
		t.Fatalf("coverage must wait for green tests, got %v", diags)
	}
}

// Type discipline (v0.4): float ban, dec exactness, brands, seal, extern.

func TestDiagnoseFloatLiteral(t *testing.T) {
	bad := strings.Replace(lspAuth, "ok(\"u\") => Ok(id = \"u\")", "ok(id = \"u\") => Ok(id = 0.5)", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
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
  Ok(price = d"1.50")
`

func TestDiagnoseDecExact(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"shop.can": typeShop})
	if diags := diagnose(dir, "shop.can", typeShop); len(diags) != 0 {
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
    t(seal M__B("s")) => Ok(echo = "s")
  Ok(echo = pw)
`

func TestDiagnoseBrandSink(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": typeLeak})
	diags := diagnose(dir, "m.can", typeLeak)
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
    t(seal M__A("x"), seal M__B("y")) => Ok(echo = "n")
  match a == b
    true => Ok(echo = "y")
    false => Ok(echo = "n")
`

func TestDiagnoseCrossBrandCompare(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": typeCmp})
	diags := diagnose(dir, "m.can", typeCmp)
	checkSpan(t, typeCmp, diags, "cannot compare M__A with M__B", "==", expectLine(t, typeCmp, "match a == b"))
}

const typeStrictMixed = `mod m
  provides [m__cmp, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  echo: str
)

fn m__cmp(a: int, b: dec) -> M__Out rev 1
  emits []
  tests
    t(1, d"1.0") => Ok(echo = "n")
  match a > b
    true => Ok(echo = "y")
    false => Ok(echo = "n")
`

func TestDiagnoseStrictMixedCompare(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": typeStrictMixed})
	diags := diagnose(dir, "m.can", typeStrictMixed)
	checkSpan(t, typeStrictMixed, diags, "cannot compare int with dec", ">", expectLine(t, typeStrictMixed, "match a > b"))
}

const typeTextOps = `mod m
  provides [m__len, m__at, m__at_str, m__slice, m__slice_str, M__Out, M__Str]
  uses []
  emits []

type M__Out rev 1 (
  n: int
)

type M__Str rev 1 (
  s: str
)

fn m__len(v: int) -> M__Out rev 1
  emits []
  tests
    t(1) => Ok(n = 0)
  Ok(n = #v)

fn m__at(v: int, i: int) -> M__Out rev 1
  emits []
  tests
    t(1, 0) => Ok(n = 0)
  Ok(n = v[i])

fn m__at_str(v: str, i: str) -> M__Out rev 1
  emits []
  tests
    t("ab", "x") => Ok(n = 0)
  Ok(n = v[i])

fn m__slice(v: int, a: int, b: int) -> M__Str rev 1
  emits []
  tests
    t(1, 0, 1) => Ok(s = "x")
  Ok(s = v[a:b])

fn m__slice_str(v: str, a: str, b: int) -> M__Str rev 1
  emits []
  tests
    t("ab", "x", 1) => Ok(s = "x")
  Ok(s = v[a:b])
`

func TestDiagnoseTextOpsMismatch(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": typeTextOps})
	diags := diagnose(dir, "m.can", typeTextOps)
	for _, want := range []string{
		"cannot count scalars of int",
		"cannot index into int",
		"cannot index with str",
		"cannot slice int",
		"cannot slice with str",
	} {
		if !hasDiag(diags, "error", want) {
			t.Fatalf("expected %q, got %v", want, diags)
		}
	}
}

const typeSeal = `mod m
  provides [m__seal, M__B, M__Out]
  uses []
  emits []

brand M__B is str rev 1

type M__Out rev 1 (
  echo: str
)

fn m__seal(pw: str, n: int) -> M__Out rev 1
  emits []
  tests
    t("s", 0) => Ok(echo = "s")
  Ok(echo = seal M__B(pw))
`

// a25 supersedes the literal-only seal rule: seals take string
// literals or string-typed refs and fields, so decision-tabled
// constructors can mint computed brands. What stays rejected is
// anything not statically a string — int literals and int refs.
func TestDiagnoseSealString(t *testing.T) {
	badLit := strings.Replace(typeSeal, "Ok(echo = seal M__B(pw))", "Ok(echo = seal M__B(7))", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": badLit})
	diags := diagnose(dir, "m.can", badLit)
	checkSpan(t, badLit, diags, "seal M__B takes a string", "M__B", expectLine(t, badLit, "seal M__B(7)"))

	badRef := strings.Replace(typeSeal, "Ok(echo = seal M__B(pw))", "Ok(echo = seal M__B(n))", 1)
	dir2 := writeLSPDir(t, map[string]string{"m.can": badRef})
	diags2 := diagnose(dir2, "m.can", badRef)
	checkSpan(t, badRef, diags2, "seal M__B takes a string", "M__B", expectLine(t, badRef, "seal M__B(n)"))
}

func TestDiagnoseUnknownType(t *testing.T) {
	bad := strings.Replace(lspAuth, "fn auth__go(id: str)", "fn auth__go(id: Nope)", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": bad})
	diags := diagnose(dir, "auth.can", bad)
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
    yes(seal M__Pw("s")) => Ok(ok = true)
    no(seal M__Pw("s")) => m.nope()
  match call m__use(pw)
    given
      yes => [exchange args (pw = seal M__Pw("s")) outcome Ok()]
      no => [exchange args (pw = seal M__Pw("s")) outcome m.nope()]
    on m.nope _ => m.nope()
    on Ok v => Ok(ok = true)
`

func TestDiagnoseExternClean(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": typeExtern})
	if diags := diagnose(dir, "m.can", typeExtern); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestDiagnoseExternBoolRet(t *testing.T) {
	bad := strings.Replace(typeExtern, "extern m__use(pw: M__Pw) -> M__Verdict rev 1", "extern m__use(pw: M__Pw) -> bool rev 1", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
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
    t("u") => Ok(id = "u")
  match call m__use(id)
    given
      t => [exchange args (pw = "u") outcome Ok()]
    on Ok v => Ok(id = id)
`

func TestDiagnoseForeignExtern(t *testing.T) {
	files := map[string]string{"m.can": typeExtern, "o.can": typeForeign}
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, "o.can", typeForeign)
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
    t(10, 3, 2) => Ok(n = 5)
  Ok(n = a - b - c)

fn m__prec(a: int, b: int, c: int) -> M__Out rev 1
  emits []
  tests
    t(2, 3, 4) => Ok(n = 14)
  Ok(n = a + b * c)
`

func TestDiagnoseArithClean(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": typeArith})
	if diags := diagnose(dir, "m.can", typeArith); len(diags) != 0 {
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
    third(d"0.1", d"0.2", d"0.3") => Ok(yes = true)
    fourth(d"0.1", d"0.2", d"0.4") => Ok(yes = false)
  match a + b == c
    true => Ok(yes = true)
    false => Ok(yes = false)
`

func TestDiagnoseArithDecExact(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": typeExact})
	if diags := diagnose(dir, "m.can", typeExact); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestDiagnoseArithMixed(t *testing.T) {
	bad := strings.Replace(typeArith, "Ok(n = a - b - c)", `Ok(n = a + d"1.5")`, 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
	checkSpan(t, bad, diags, "cannot add int with dec", "+", expectLine(t, bad, `Ok(n = a + d"1.5")`))
}

func TestDiagnoseArithStr(t *testing.T) {
	// a16: + concatenates strings, so the refused string operation
	// is now -. The contract under test is unchanged: strings do
	// no arithmetic besides explicit construction.
	bad := strings.Replace(typeArith, "(a: int, b: int, c: int) -> M__Out rev 1\n  emits []\n  tests\n    t(10, 3, 2) => Ok(n = 5)\n  Ok(n = a - b - c)",
		"(a: str, b: str, c: int) -> M__Out rev 1\n  emits []\n  tests\n    t(\"x\", \"y\", 2) => Ok(n = 5)\n  Ok(n = a - b)", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": bad})
	diags := diagnose(dir, "m.can", bad)
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
    basic("x", "y") => Ok(s = "xy")
    empty("", "y") => Ok(s = "y")
  Ok(s = left + right)
`

func TestDiagnoseStrConcatClean(t *testing.T) {
	// a16: str + str is explicit construction, not a mismatch.
	dir := writeLSPDir(t, map[string]string{"m.can": typeConcat})
	if diags := diagnose(dir, "m.can", typeConcat); len(diags) != 0 {
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
    t(7, 3) => Ok(q = 2)
  Ok(q = a / b)
`

func TestDiagnoseDivisionDeferred(t *testing.T) {
	// a17 discharges the a06 deferral for integers: / is exact
	// Euclidean division with a loud zero divisor. What stays
	// deferred is decimal division (CAN6005, pinned separately).
	dir := writeLSPDir(t, map[string]string{"m.can": typeIntDiv})
	if diags := diagnose(dir, "m.can", typeIntDiv); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
