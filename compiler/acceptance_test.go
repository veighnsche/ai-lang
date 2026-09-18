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
    held(7) => Ok(7, "hi") pinned
    loose(8) => Ok(8, "hi")
    wordy(9) => Ok(9, "pinned")
    pinned(10) => Ok(10, "hi")
  match n == 9
    true => Ok(n, "pinned")
    false => Ok(n, "hi")
`

func acceptanceTests(t *testing.T, src string) []Test {
	t.Helper()
	m, err := parseModuleText("acc.can", src)
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
	dir := writeLSPDir(t, map[string]string{"acc.can": acceptanceSrc})
	if diags := diagnose(dir, "acc.can", acceptanceSrc); len(diags) != 0 {
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
    keep(5) => Ok(5, "hi") pinned
    flex(6) => Ok(6, "hi")
  match n == 5
    true => Ok(5, "hi")
    false => Ok(n, "hi")
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
	files := map[string]string{"acc.can": src}
	return revisionProg(t, files, []string{"acc.can"})
}

func TestPinnedWeakeningWarns(t *testing.T) {
	prog, _ := acceptanceReprog(t, acceptanceBehav)
	base := acceptanceBase(t, prog)
	// The attack: move implementation and expectation together so the
	// suite stays green while acceptance drifts.
	moved := strings.Replace(acceptanceBehav, "keep(5) => Ok(5,", "keep(5) => Ok(6,", 1)
	moved = strings.Replace(moved, "true => Ok(5,", "true => Ok(6,", 1)
	prog2, texts2 := acceptanceReprog(t, moved)
	diags := CheckPinnedRows(prog2, texts2, base)
	if !hasCode(diags, "CAN6017") {
		t.Fatalf("expected CAN6017 weakening warning, got %v", diags)
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
	moved := strings.Replace(acceptanceBehav, "false => Ok(n,", "false => Ok(n + 1,", 1)
	moved = strings.Replace(moved, "flex(6) => Ok(6,", "flex(6) => Ok(7,", 1)
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
    keep(5) => Ok(5, "hi") pinned
    flex(6) => Ok(6, "hi")
  Ok(n, "hi")
`

func TestPinnedRemovedWarns(t *testing.T) {
	prog, _ := acceptanceReprog(t, acceptanceFlat)
	base := acceptanceBase(t, prog)
	lines := strings.Split(acceptanceFlat, "\n")
	var kept []string
	for _, l := range lines {
		if strings.Contains(l, "keep(5)") {
			continue
		}
		kept = append(kept, l)
	}
	prog2, texts2 := acceptanceReprog(t, strings.Join(kept, "\n"))
	diags := CheckPinnedRows(prog2, texts2, base)
	if !hasCode(diags, "CAN6017") || !hasFound(diags, "removed") {
		t.Fatalf("expected CAN6017 removal warning, got %v", diags)
	}
}

func TestPinnedDemoteWarns(t *testing.T) {
	prog, _ := acceptanceReprog(t, acceptanceBehav)
	base := acceptanceBase(t, prog)
	demoted := strings.Replace(acceptanceBehav, `) pinned`, `)`, 1)
	prog2, texts2 := acceptanceReprog(t, demoted)
	diags := CheckPinnedRows(prog2, texts2, base)
	if !hasCode(diags, "CAN6017") || !hasFound(diags, "demoted") {
		t.Fatalf("expected CAN6017 demotion warning, got %v", diags)
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
