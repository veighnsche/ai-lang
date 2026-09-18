package main

import (
	"testing"
)

// a62: a match-call on a function that resolves nowhere must
// report exactly one diagnostic (CAN3001, the cause). The
// given/proof/execution rungs that vanish with the fix are
// noise: CAN3101, CAN3103, CAN4101, CAN4102, CAN4200, CAN4107.

func diagCodes(diags []Diag) map[string]bool {
	out := map[string]bool{}
	for _, d := range diags {
		out[d.Code] = true
	}
	return out
}

// TestUnknownCallCascadeScripted pins the scripted shape: given
// present, both arms. Only the unknown-callee error may fire.
func TestUnknownCallCascadeScripted(t *testing.T) {
	client := `mod client
  provides [client__go]
  uses []
  emits []

fn client__go(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_base64]
  tests
    go("QQ==") => Ok(value = Bytes(Seq<int>[65]))
    bad("QUI") => encoding.invalid_base64(value = "QUI")
=
  match call nosuch__decode(value = value)
    given
      go => [exchange args (value = "QQ==") outcome Ok(value = Bytes(Seq<int>[65]))]
      bad => [exchange args (value = "QUI") outcome encoding.invalid_base64(value = "QUI")]
    on Ok r => Ok(value = r.value)
    on encoding.invalid_base64 e => encoding.invalid_base64(value = e.value)
`
	dir := writeLSPDir(t, map[string]string{"client.can": client})
	diags := diagnose(dir, "client.can", client)
	codes := diagCodes(diags)
	if !codes[CodeUnknownCall] {
		t.Fatalf("expected CAN3001, got %v", diags)
	}
	for _, c := range []string{CodeNoGiven, CodeStubNotInEmit, CodeMissingArm, CodeStaleArm, CodeTestFailed, CodeArmUntaken} {
		if codes[c] {
			t.Fatalf("cascade rung %s must be suppressed, got %v", c, diags)
		}
	}
}

// TestUnknownCallCascadeBare pins the unscripted shape: no
// given, both arms. Same single-diagnostic contract.
func TestUnknownCallCascadeBare(t *testing.T) {
	client := `mod client
  provides [client__bare]
  uses []
  emits []

fn client__bare(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_base64]
  tests
    go("QQ==") => Ok(value = Bytes(Seq<int>[65]))
=
  match call nosuch__decode(value = value)
    on Ok r => Ok(value = r.value)
    on encoding.invalid_base64 e => encoding.invalid_base64(value = e.value)
`
	dir := writeLSPDir(t, map[string]string{"client.can": client})
	diags := diagnose(dir, "client.can", client)
	codes := diagCodes(diags)
	if !codes[CodeUnknownCall] {
		t.Fatalf("expected CAN3001, got %v", diags)
	}
	for _, c := range []string{CodeNoGiven, CodeStubNotInEmit, CodeMissingArm, CodeStaleArm, CodeTestFailed, CodeArmUntaken} {
		if codes[c] {
			t.Fatalf("cascade rung %s must be suppressed, got %v", c, diags)
		}
	}
}
