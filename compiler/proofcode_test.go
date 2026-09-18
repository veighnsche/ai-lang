package main

import (
	"errors"
	"testing"
)

const proofTestText = "mod m\n  provides [m__g]\n  uses []\n  emits []\n"

// Proof-diagnostic codes travel typed on the error from the producer:
// the code must survive line-wrapping and must not depend on message
// wording, so these messages deliberately match no legacy substring.
func TestProofDiagTypedCodes(t *testing.T) {
	cases := []struct {
		code string
		msg  string
	}{
		{CodeProofOther, "m: boom"},
		{CodeMissingArm, "m: non-exhaustive match, missing ok"},
		{CodeStaleArm, "m: stale match arm Other"},
		{CodeBoolArms, "m: bang"},
		{CodeValueNoWild, "m: zap"},
		{CodeBadArmKind, "m: kapow"},
		{CodeVariantOnVal, "m: thud"},
		{CodeUselessArm, "m: wham"},
		{CodeUselessAlt, "m: moof"},
	}
	for _, c := range cases {
		err := at(3, proofErrf(c.code, "%s", c.msg))
		if got := proofDiag(proofTestText, err).Code; got != c.code {
			t.Errorf("typed %s diag code = %s, want %s (msg %q)", c.code, got, c.code, c.msg)
		}
		if got, ok := proofCode(err); !ok || got != c.code {
			t.Errorf("proofCode through LineError = %q, %v; want %q, true", got, ok, c.code)
		}
	}
}

// Untyped errors keep the old default: the message switch is gone, so
// even legacy wording no longer selects a code.
func TestProofDiagUntypedFallsBack(t *testing.T) {
	legacy := []string{
		"m: non-exhaustive match, missing ok",
		"m: stale match arm Other",
		"m: bool match must be exactly true+false",
		"m: arm x is fully covered by earlier arms",
		"m: or-alternative a contributes no remaining space",
		"m: or-pattern alternatives take scalar patterns only",
	}
	for _, msg := range legacy {
		err := at(2, errors.New(msg))
		if got := proofDiag(proofTestText, err).Code; got != CodeProofOther {
			t.Errorf("untyped diag code for %q = %s, want %s", msg, got, CodeProofOther)
		}
	}
}

// The missing-arm hint still reads message data (not identity).
func TestProofDiagMissingHint(t *testing.T) {
	err := at(1, proofErrf(CodeMissingArm, "m: non-exhaustive match, missing ok"))
	d := proofDiag(proofTestText, err)
	if d.Hint == "" || d.Expected == "" {
		t.Fatalf("missing-arm diag lost its hint: %+v", d)
	}
}
