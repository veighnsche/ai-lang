package main

import (
	"strings"
	"testing"
)

// a71: machine-actionable payloads on four codes plus
// `canlc explain`. Probes first: Diag has no payload
// fields and explainCode does not exist.

// TestExplainKnown pins per-code docs for the four
// payload codes.
func TestExplainKnown(t *testing.T) {
	for _, code := range []string{CodeForeignRaise, CodeBareErrorKind, CodeGivenDashRetired, CodeMissingArm} {
		text, ok := explainCode(code)
		if !ok {
			t.Fatalf("expected explain entry for %s", code)
		}
		if !strings.Contains(text, code) {
			t.Fatalf("explain for %s must name its code, got %q", code, text)
		}
	}
	if text, _ := explainCode(CodeMissingArm); !strings.Contains(text, "on ") {
		t.Fatalf("CAN4101 explain must show the arm shape, got %q", text)
	}
}

// TestExplainFallback pins the fallback mechanism: a registered
// code without a per-code doc gets family text that says so;
// an unregistered code fails. Uses a temp-registered code so
// the mechanism stays tested at full coverage.
func TestExplainFallback(t *testing.T) {
	allCodes = append(allCodes, "CAN9998")
	defer func() { allCodes = allCodes[:len(allCodes)-1] }()
	text, ok := explainCode("CAN9998")
	if !ok {
		t.Fatal("expected family fallback for temp-registered CAN9998")
	}
	if !strings.Contains(text, "No per-code doc yet") {
		t.Fatalf("fallback must admit its cut, got %q", text)
	}
	if _, ok := explainCode("CAN9999"); ok {
		t.Fatal("unregistered code must fail")
	}
}

// TestExplainComplete pins full coverage: every registered
// code carries a per-code doc with rule, violation, and fix.
// Add the entry with the code, not after.
func TestExplainComplete(t *testing.T) {
	for _, code := range allCodes {
		text, ok := explainCode(code)
		if !ok {
			t.Fatalf("expected explain entry for %s", code)
		}
		for _, want := range []string{code, "rule:", "violation:", "fix:"} {
			if !strings.Contains(text, want) {
				t.Fatalf("explain for %s must contain %q, got %q", code, want, text)
			}
		}
		if strings.Contains(text, "No per-code doc yet") {
			t.Fatalf("explain for %s falls back; write the entry: %q", code, text)
		}
	}
}

// TestPayloadForeignRaise pins CAN4001 fields.
func TestPayloadForeignRaise(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits [m.known]

error m.known()
error m.stray()

type M__Out rev 1 (
  value: int
)

fn m__go(value: int) -> M__Out rev 1
  emits [m.known]
  tests
    go(1) => m.stray()
  m.stray()
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	var found *Diag
	for i, d := range diags {
		if d.Code == CodeForeignRaise {
			found = &diags[i]
		}
	}
	if found == nil {
		t.Fatalf("expected CAN4001, got %v", diags)
	}
	if found.Found != "m.stray" {
		t.Fatalf("Found must name the raised kind, got %+v", found)
	}
	if !strings.Contains(found.Expected, "m.known") {
		t.Fatalf("Expected must list declared emits, got %+v", found)
	}
	if found.Hint == "" {
		t.Fatalf("Hint must be non-empty, got %+v", found)
	}
}

// TestPayloadBareKind pins CAN3204 fields.
func TestPayloadBareKind(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits [m.bad]

error m.bad(id: str, n: int)

type M__Out rev 1 (
  value: int
)

fn m__go(value: int) -> M__Out rev 1
  emits [m.bad]
  tests
    go(1) => m.bad
  Ok(value = value)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	var found *Diag
	for i, d := range diags {
		if d.Code == CodeBareErrorKind {
			found = &diags[i]
		}
	}
	if found == nil {
		t.Fatalf("expected CAN3204, got %v", diags)
	}
	if found.Found != "m.bad" {
		t.Fatalf("Found must name the bare kind, got %+v", found)
	}
	if !strings.Contains(found.Expected, "m.bad(") || !strings.Contains(found.Expected, "id") || !strings.Contains(found.Expected, "n") {
		t.Fatalf("Expected must show the complete construction, got %+v", found)
	}
}

// TestPayloadOmissionBackstop pins the a91 world: a scriptless
// test draws no static diagnostic, and reaching the call without
// a script fails the test at execution (CAN4200).
func TestPayloadOmissionBackstop(t *testing.T) {
	auth := strings.Replace(lspAuth,
		"    down(\"u\") => auth.bad()\n",
		"    down(\"u\") => auth.bad()\n    extra(\"u\") => auth.bad()\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": auth})
	diags := diagnose(dir, "auth.can", auth)
	for _, d := range diags {
		if strings.Contains(d.Msg, "has no script") {
			t.Fatalf("no static diagnostic for omission, got %+v", d)
		}
	}
	var found *Diag
	for i, d := range diags {
		if d.Code == CodeTestFailed && strings.Contains(d.Msg, "extra") {
			found = &diags[i]
		}
	}
	if found == nil {
		t.Fatalf("expected CAN4200 for extra, got %v", diags)
	}
	if !strings.Contains(found.Msg, "no script for call") {
		t.Fatalf("backstop must name the missing script, got %+v", found)
	}
}

// TestPayloadGivenDashRetired pins CAN3111 fields: a `-` row
// is a hard error naming the row and its deletion.
func TestPayloadGivenDashRetired(t *testing.T) {
	auth := strings.Replace(lspAuth,
		"      down => [exchange args (id = \"u\") outcome db.down()]\n",
		"      down => [exchange args (id = \"u\") outcome db.down()]\n      down2 => -\n", 1)
	auth = strings.Replace(auth,
		"    down(\"u\") => auth.bad()\n",
		"    down(\"u\") => auth.bad()\n    down2(\"u\") => auth.bad()\n", 1)
	dir := writeLSPDir(t, map[string]string{"db.can": lspDB, "auth.can": auth})
	diags := diagnose(dir, "auth.can", auth)
	var found *Diag
	for i, d := range diags {
		if d.Code == CodeGivenDashRetired {
			found = &diags[i]
		}
	}
	if found == nil {
		t.Fatalf("expected CAN3111, got %v", diags)
	}
	if found.Found != "down2 => -" {
		t.Fatalf("Found must name the dash row, got %+v", found)
	}
	if !strings.Contains(found.Expected, "omission") {
		t.Fatalf("Expected must point at omission, got %+v", found)
	}
	if found.Hint == "" {
		t.Fatalf("Hint must be non-empty, got %+v", found)
	}
}

// TestPayloadMissingArm pins CAN4101 fields.
func TestPayloadMissingArm(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits [m.bad]

error m.bad()

type M__Out rev 1 (
  value: int
)

fn m__go(value: int) -> M__Out rev 1
  emits [m.bad]
  tests
    go(1) => Ok(value = 1)
  match call m__help(value)
    given
      go => [exchange args (value = 1) outcome Ok(value = 1)]
    on Ok r => Ok(value = r.value)

fn m__help(value: int) -> M__Out rev 1
  emits [m.bad]
  tests
    go(1) => Ok(value = 1)
  Ok(value = value)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	var found *Diag
	for i, d := range diags {
		if d.Code == CodeMissingArm {
			found = &diags[i]
		}
	}
	if found == nil {
		t.Fatalf("expected CAN4101, got %v", diags)
	}
	if found.Expected != "m.bad" {
		t.Fatalf("Expected must name the missing kind, got %+v", found)
	}
	if !strings.Contains(found.Hint, "on m.bad") {
		t.Fatalf("Hint must show the arm shape, got %+v", found)
	}
}

// TestMissingArmFound pins the uniform payload triple: like its
// ForeignRaise/BareKind/Dangling siblings, CAN4101 names what the
// match covers (Found) beside what it wants (Expected).
func TestMissingArmFound(t *testing.T) {
	src := `mod m
  provides [m__go, M__Out]
  uses []
  emits [m.bad]

error m.bad()

type M__Out rev 1 (
  value: int
)

fn m__go(value: int) -> M__Out rev 1
  emits [m.bad]
  tests
    go(1) => Ok(value = 1)
  match call m__help(value)
    given
      go => [exchange args (value = 1) outcome Ok(value = 1)]
    on Ok r => Ok(value = r.value)

fn m__help(value: int) -> M__Out rev 1
  emits [m.bad]
  tests
    go(1) => Ok(value = 1)
  Ok(value = value)
`
	dir := writeLSPDir(t, map[string]string{"m.can": src})
	diags := diagnose(dir, "m.can", src)
	var found *Diag
	for i, d := range diags {
		if d.Code == CodeMissingArm {
			found = &diags[i]
		}
	}
	if found == nil {
		t.Fatalf("expected CAN4101, got %v", diags)
	}
	if found.Found != "Ok" {
		t.Fatalf("Found must name the arms present, got %+v", found)
	}
}
