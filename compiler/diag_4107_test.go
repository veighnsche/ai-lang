package main

import (
	"os"
	"testing"
)

// Chatbot verdict C: CAN4107 is KEPT. Scripts are not
// coverage: each caller's own handlers need passing-test
// evidence. This file commits the verdict's two-function
// negative fixture as the regression.

// TestArmTakenPerCaller pins the two CAN4107 locations:
// client__success corrupts the rejected payload and
// client__failure discards decoded bytes, each exercised
// only on its correct path.
func TestArmTakenPerCaller(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.can")
	if err != nil {
		t.Fatalf("read text.can: %v", err)
	}
	client := `mod client
  provides [client__success, client__failure]
  uses [std__hex__decode@1]
  emits [encoding.invalid_hex]

fn client__success(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_hex]
  tests
    good("41") => Ok(Bytes(Seq<int>[65]))
  match call std__hex__decode(value = value)
    given
      good => exchange args (value = "41") outcome Ok(Bytes(Seq<int>[65]))
      bad => exchange args (value = "zz") outcome encoding.invalid_hex(value = "zz")
    on Ok r => Ok(r.value)
    on encoding.invalid_hex e => encoding.invalid_hex(value = "lost")

fn client__failure(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_hex]
  tests
    bad("zz") => encoding.invalid_hex(value = "zz")
  match call std__hex__decode(value = value)
    given
      good => exchange args (value = "41") outcome Ok(Bytes(Seq<int>[65]))
      bad => exchange args (value = "zz") outcome encoding.invalid_hex(value = "zz")
    on Ok r => Ok(Bytes(Seq<int>[]))
    on encoding.invalid_hex e => encoding.invalid_hex(value = e.value)
`
	files := map[string]string{"text.can": string(raw), "client.can": client}
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, "client.can", client)
	var taken []Diag
	for _, d := range diags {
		if d.Code == CodeArmUntaken {
			taken = append(taken, d)
		}
		if d.Sev == "error" && d.Code != CodeArmUntaken {
			t.Fatalf("unexpected error diag %v in %v", d, diags)
		}
	}
	if len(taken) != 2 {
		t.Fatalf("expected exactly two CAN4107, got %v", diags)
	}
}

// TestArmTakenMissingRowsExposeBugs pins the second half:
// adding the missing rows with correct expectations fails
// behaviorally (CAN4200) against the corrupt handlers.
func TestArmTakenMissingRowsExposeBugs(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.can")
	if err != nil {
		t.Fatalf("read text.can: %v", err)
	}
	client := `mod client
  provides [client__success, client__failure]
  uses [std__hex__decode@1]
  emits [encoding.invalid_hex]

fn client__success(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_hex]
  tests
    good("41") => Ok(Bytes(Seq<int>[65]))
    bad("zz") => encoding.invalid_hex(value = "zz")
  match call std__hex__decode(value = value)
    given
      good => exchange args (value = "41") outcome Ok(Bytes(Seq<int>[65]))
      bad => exchange args (value = "zz") outcome encoding.invalid_hex(value = "zz")
    on Ok r => Ok(r.value)
    on encoding.invalid_hex e => encoding.invalid_hex(value = "lost")

fn client__failure(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_hex]
  tests
    bad("zz") => encoding.invalid_hex(value = "zz")
    good("41") => Ok(Bytes(Seq<int>[65]))
  match call std__hex__decode(value = value)
    given
      good => exchange args (value = "41") outcome Ok(Bytes(Seq<int>[65]))
      bad => exchange args (value = "zz") outcome encoding.invalid_hex(value = "zz")
    on Ok r => Ok(Bytes(Seq<int>[]))
    on encoding.invalid_hex e => encoding.invalid_hex(value = e.value)
`
	files := map[string]string{"text.can": string(raw), "client.can": client}
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, "client.can", client)
	var failed []Diag
	for _, d := range diags {
		if d.Code == CodeTestFailed {
			failed = append(failed, d)
		}
	}
	if len(failed) != 2 {
		t.Fatalf("expected exactly two CAN4200, got %v", diags)
	}
	if hasErrCode(diags, CodeArmUntaken) {
		t.Fatalf("no arm may stay untaken, got %v", diags)
	}
}
