package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// a82: verifier activation. Contracted functions must verify at
// compile time (CLI) and surface in the editor (LSP); passing
// compilations report the verification scope. Probes first: these
// observe the unwired behavior (contracts ignored) before the
// wiring lands.

// activateMutant is max+1 with rows following the mutant body:
// execution stays green, so only the proof can reject it.
var activateMutant = strings.Replace(admitMax,
	"on true => Ok(value = right)", "on true => Ok(value = right + 1)", 1)

var activateMutantRows = strings.Replace(activateMutant,
	"ordered(left = 1, right = 2) => Ok(value = 2)", "ordered(left = 1, right = 2) => Ok(value = 3)", 1)

func writeActivateDir(t *testing.T, files map[string]string) (string, []string) {
	t.Helper()
	dir := t.TempDir()
	var paths []string
	for name, text := range files {
		fp := filepath.Join(dir, name)
		if err := os.WriteFile(fp, []byte(text), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		paths = append(paths, fp)
	}
	return dir, paths
}

// TestActivateBlocksUnproven pins the CLI gate: a contracted
// program whose proof fails does not compile, and the failure
// names the unproven obligation.
func TestActivateBlocksUnproven(t *testing.T) {
	dir, paths := writeActivateDir(t, map[string]string{"m.can": activateMutantRows})
	out := filepath.Join(dir, "out")
	err := compile(out, paths)
	if err == nil {
		t.Fatalf("unproven contract compiled clean")
	}
	// Prose failures carry the message; the stable code rides
	// the diagnose squiggle pinned below.
	if !strings.Contains(err.Error(), "unproven") {
		t.Fatalf("failure names no unproven obligation: %v", err)
	}
}

// TestActivateAcceptsVerified pins the open gate: a verifying
// contracted program still compiles.
func TestActivateAcceptsVerified(t *testing.T) {
	dir, paths := writeActivateDir(t, map[string]string{"m.can": admitMax})
	out := filepath.Join(dir, "out")
	if err := compile(out, paths); err != nil {
		t.Fatalf("verified contract failed to compile: %v", err)
	}
}

// TestActivateDiagnoseSquiggles pins the editor gate: drift from
// the contract surfaces CAN4304 on the open document.
func TestActivateDiagnoseSquiggles(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"m.can": activateMutantRows})
	diags := diagnose(dir, "m.can", activateMutantRows)
	if !hasCode(diags, CodeContractUnproven) {
		t.Fatalf("editor reported no unproven finding: %v", diags)
	}
}

// TestActivateSkipsBroken pins the error gate: programs that fail
// ordinary checks gain no proof noise on top of real errors.
func TestActivateSkipsBroken(t *testing.T) {
	broken := strings.Replace(admitMax,
		"on false => Ok(value = left)", "on false => Ok(value = nosuch)", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": broken})
	diags := diagnose(dir, "m.can", broken)
	if len(diags) == 0 {
		t.Fatalf("broken program reported nothing")
	}
	for _, d := range diags {
		if strings.HasPrefix(d.Code, "CAN43") {
			t.Fatalf("broken program gained proof noise: %v", d)
		}
	}
}

// TestActivateReport pins evidence reporting: a passing program
// names its verified declarations and marks the rest.
func TestActivateReport(t *testing.T) {
	prog, _ := admitProg(t, admitMax)
	verified, uncontracted := contractIdentities(prog)
	if len(verified) != 1 || verified[0] != "m__max@1" {
		t.Fatalf("verified identities: %v", verified)
	}
	if len(uncontracted) != 0 {
		t.Fatalf("uncontracted identities: %v", uncontracted)
	}
	mixed := strings.Replace(admitMax, "provides [m__max, M__Out]", "provides [m__max, m__plain, M__Out]", 1) + `
fn m__plain(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 1) => Ok(value = 1)
=
  Ok(value = x)
`
	prog, _ = admitProg(t, mixed)
	verified, uncontracted = contractIdentities(prog)
	if len(verified) != 1 || len(uncontracted) != 1 || uncontracted[0] != "m__plain@1" {
		t.Fatalf("mixed identities: verified=%v uncontracted=%v", verified, uncontracted)
	}
}
