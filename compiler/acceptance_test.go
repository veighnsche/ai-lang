package main

import (
	"strings"
	"testing"
)

// A-light acceptance authority (a87): a trailing `pinned` word marks a
// decision-table row as trusted acceptance. Slice 1 covers syntax only:
// the marker parses, sets the flag, and changes no behavior.

const acceptanceSrc = `mod acc
  provides [acc__f, Acc__V]
  uses []
  emits []

type Acc__V rev 1 (
  value: int
  note: str
)

fn acc__f(n: int) -> Acc__V rev 1
  emits []
  tests
    held(n = 7) => Ok(value = 7, note = "hi") pinned
    loose(n = 8) => Ok(value = 8, note = "hi")
    wordy(n = 9) => Ok(value = 9, note = "pinned")
    pinned(n = 10) => Ok(value = 10, note = "hi")
=
  match n == 9
    true => Ok(value = n, note = "pinned")
    false => Ok(value = n, note = "hi")
`

func acceptanceTests(t *testing.T, src string) []Test {
	t.Helper()
	m, err := parseModuleText("acc.ail", src)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	for _, d := range m.Decls {
		if fn, ok := d.(*FnDecl); ok && fn.Name == "acc__f" {
			return fn.Tests
		}
	}
	t.Fatalf("acc__f not found")
	return nil
}

func TestPinnedRowSetsFlag(t *testing.T) {
	tests := acceptanceTests(t, acceptanceSrc)
	if len(tests) != 4 {
		t.Fatalf("expected 4 tests, got %d", len(tests))
	}
	if !tests[0].Pinned {
		t.Errorf("held row with marker: Pinned = false, want true")
	}
	for _, tc := range tests[1:] {
		if tc.Pinned {
			t.Errorf("row %q without marker: Pinned = true, want false", tc.Name)
		}
	}
}

func TestPinnedRowDiagnosesClean(t *testing.T) {
	// The marker is behavior-free in slice 1: a pinned file gates
	// exactly like its unpinned twin.
	dir := writeLSPDir(t, map[string]string{"acc.ail": acceptanceSrc})
	if diags := diagnose(dir, "acc.ail", acceptanceSrc); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// The behavior fixture is branch-free so rows can be deleted without
// tripping coverage: every check below targets acceptance drift only.

const acceptanceBehav = `mod acc
  provides [acc__w, Acc__V]
  uses []
  emits []

type Acc__V rev 1 (
  value: int
  note: str
)

fn acc__w(n: int) -> Acc__V rev 1
  emits []
  tests
    keep(n = 5) => Ok(value = 5, note = "hi") pinned
    flex(n = 6) => Ok(value = 6, note = "hi")
=
  match n == 5
    true => Ok(value = 5, note = "hi")
    false => Ok(value = n, note = "hi")
`

func acceptanceBase(t *testing.T, prog *Program) *RevisionBaseline {
	t.Helper()
	return &RevisionBaseline{
		Format: RevisionFormat, Origin: "test", Accepted: true,
		Scope:   []string{"acc"},
		Entries: FingerprintProgram(prog),
		Pinned:  PinnedRows(prog),
	}
}

func acceptanceReprog(t *testing.T, src string) (*Program, map[string]string) {
	t.Helper()
	files := map[string]string{"acc.ail": src}
	return revisionProg(t, files, []string{"acc.ail"})
}

func TestPinnedWeakeningWarns(t *testing.T) {
	prog, _ := acceptanceReprog(t, acceptanceBehav)
	base := acceptanceBase(t, prog)
	// The attack: move implementation and expectation together so the
	// suite stays green while acceptance drifts.
	moved := strings.Replace(acceptanceBehav, "keep(n = 5) => Ok(value = 5,", "keep(n = 5) => Ok(value = 6,", 1)
	moved = strings.Replace(moved, "true => Ok(value = 5,", "true => Ok(value = 6,", 1)
	prog2, texts2 := acceptanceReprog(t, moved)
	diags := CheckPinnedRows(prog2, texts2, base)
	if !hasCode(diags, "AIL6017") {
		t.Fatalf("expected AIL6017 weakening warning, got %v", diags)
	}
	if !hasFound(diags, "keep") {
		t.Fatalf("expected row name in diagnostic, got %v", diags)
	}
	for _, d := range diags {
		if d.Sev != "warning" {
			t.Fatalf("weakening must warn, never error: got %v", diags)
		}
	}
	if err := firstError(diags); err != nil {
		t.Fatalf("warnings must not gate: firstError = %v", err)
	}
}

func TestPinnedUnpinnedChurnSilent(t *testing.T) {
	prog, _ := acceptanceReprog(t, acceptanceBehav)
	base := acceptanceBase(t, prog)
	moved := strings.Replace(acceptanceBehav, "false => Ok(value = n,", "false => Ok(value = n + 1,", 1)
	moved = strings.Replace(moved, "flex(n = 6) => Ok(value = 6,", "flex(n = 6) => Ok(value = 7,", 1)
	prog2, texts2 := acceptanceReprog(t, moved)
	if diags := CheckPinnedRows(prog2, texts2, base); len(diags) != 0 {
		t.Fatalf("unpinned churn must stay silent, got %v", diags)
	}
}

// Removal needs a branch-free body: deleting a row from a branched
// body would trip coverage, which is a different law's complaint.
const acceptanceFlat = `mod acc
  provides [acc__w, Acc__V]
  uses []
  emits []

type Acc__V rev 1 (
  value: int
  note: str
)

fn acc__w(n: int) -> Acc__V rev 1
  emits []
  tests
    keep(n = 5) => Ok(value = 5, note = "hi") pinned
    flex(n = 6) => Ok(value = 6, note = "hi")
=
  Ok(value = n, note = "hi")
`

func TestPinnedRemovedWarns(t *testing.T) {
	prog, _ := acceptanceReprog(t, acceptanceFlat)
	base := acceptanceBase(t, prog)
	lines := strings.Split(acceptanceFlat, "\n")
	var kept []string
	for _, l := range lines {
		if strings.Contains(l, "keep(n = 5)") {
			continue
		}
		kept = append(kept, l)
	}
	prog2, texts2 := acceptanceReprog(t, strings.Join(kept, "\n"))
	diags := CheckPinnedRows(prog2, texts2, base)
	if !hasCode(diags, "AIL6017") || !hasFound(diags, "removed") {
		t.Fatalf("expected AIL6017 removal warning, got %v", diags)
	}
}

func TestPinnedDemoteWarns(t *testing.T) {
	prog, _ := acceptanceReprog(t, acceptanceBehav)
	base := acceptanceBase(t, prog)
	demoted := strings.Replace(acceptanceBehav, `) pinned`, `)`, 1)
	prog2, texts2 := acceptanceReprog(t, demoted)
	diags := CheckPinnedRows(prog2, texts2, base)
	if !hasCode(diags, "AIL6017") || !hasFound(diags, "demoted") {
		t.Fatalf("expected AIL6017 demotion warning, got %v", diags)
	}
}

func TestRegenWithoutMarkerSilent(t *testing.T) {
	// Negative control: regen output records nothing for unmarked rows,
	// so acceptance can never sneak in through generation.
	plain := strings.ReplaceAll(acceptanceSrc, `) pinned`, `)`)
	prog, _ := acceptanceReprog(t, plain)
	if pins := PinnedRows(prog); len(pins) != 0 {
		t.Fatalf("unmarked regen must record no pins, got %v", pins)
	}
}
