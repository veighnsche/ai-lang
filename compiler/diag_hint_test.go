package main

import (
	"os"
	"strings"
	"testing"
)

// a63: a stale arm nested inside a match that is missing the
// same outcome almost always means the arm was attached to the
// wrong match. The stale diagnostic must point at the
// enclosing match; genuinely independent matches must not get
// the hint.

// TestStaleArmNestingHint pins the hint on the B15 shape: the
// outer failable call's error arm misattached to the inner
// total call.
func TestStaleArmNestingHint(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.can")
	if err != nil {
		t.Fatalf("read text.can: %v", err)
	}
	client := `mod client
  provides [client__go]
  uses [std__hex__decode@1, std__hex__encode@1]
  emits []

fn client__go(value: str) -> Encoding__Text rev 1
  emits [encoding.invalid_hex]
  tests
    rt("deadbeef") => Ok("deadbeef")
  match call std__hex__decode(value = value)
    given
      rt => [exchange args (value = "deadbeef") outcome Ok(Bytes(Seq<int>[222, 173, 190, 239]))]
    on Ok b => match call std__hex__encode(value = b.value)
      given
        rt => [exchange args (value = Bytes(Seq<int>[222, 173, 190, 239])) outcome Ok("deadbeef")]
      on Ok t => Ok(t.value)
      on encoding.invalid_hex e => encoding.invalid_hex(value = e.value)
`
	files := map[string]string{"text.can": string(raw), "client.can": client}
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, "client.can", client)
	var stale *Diag
	for i, d := range diags {
		if d.Code == CodeStaleArm {
			stale = &diags[i]
		}
	}
	if stale == nil {
		t.Fatalf("expected CAN4102, got %v", diags)
	}
	if !strings.Contains(stale.Msg, "enclosing match at line") {
		t.Fatalf("stale arm must hint at the enclosing match, got %q", stale.Msg)
	}
	if !hasErrCode(diags, CodeMissingArm) {
		t.Fatalf("expected companion CAN4101, got %v", diags)
	}
}

// TestStaleArmNoHintAcrossMatches pins precision: a genuinely
// stale arm with no enclosing failable match gets no hint.
func TestStaleArmNoHintAcrossMatches(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.can")
	if err != nil {
		t.Fatalf("read text.can: %v", err)
	}
	client := `mod client
  provides [client__missing, client__total]
  uses [std__hex__decode@1, std__hex__encode@1]
  emits []

fn client__missing(value: str) -> Encoding__Text rev 1
  emits [encoding.invalid_hex]
  tests
    rt("deadbeef") => Ok("deadbeef")
  match call std__hex__decode(value = value)
    given
      rt => [exchange args (value = "deadbeef") outcome Ok(Bytes(Seq<int>[222, 173, 190, 239]))]
    on Ok b => Ok("deadbeef")

fn client__total(value: Bytes) -> Encoding__Text rev 1
  emits [encoding.invalid_hex]
  tests
    ff(Bytes(Seq<int>[255])) => Ok("ff")
  match call std__hex__encode(value = value)
    given
      ff => [exchange args (value = Bytes(Seq<int>[255])) outcome Ok("ff")]
    on Ok t => Ok(t.value)
    on encoding.invalid_hex e => encoding.invalid_hex(value = e.value)
`
	files := map[string]string{"text.can": string(raw), "client.can": client}
	dir := writeLSPDir(t, files)
	diags := diagnose(dir, "client.can", client)
	var stale *Diag
	for i, d := range diags {
		if d.Code == CodeStaleArm {
			stale = &diags[i]
		}
	}
	if stale == nil {
		t.Fatalf("expected CAN4102, got %v", diags)
	}
	if strings.Contains(stale.Msg, "enclosing match at line") {
		t.Fatalf("independent stale arm must not hint, got %q", stale.Msg)
	}
}
