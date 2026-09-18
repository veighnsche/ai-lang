package main

import (
	"os"
	"testing"
)

// a58 B13: std__base64__encode is the public stdlib face of the
// B12 kernel. The foreign-caller probe proves the §6 claim that
// ordinary wrappers are normal functions, plus padding, order,
// and no-text-interpretation vectors.

// TestBytesB13ForeignCaller pins cross-module callability against
// the REAL text.can.
func TestBytesB13ForeignCaller(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.can")
	if err != nil {
		t.Fatalf("read text.can: %v", err)
	}
	client := `mod client
  provides [client__go]
  uses [std__base64__encode@1]
  emits []

fn client__go(value: Bytes) -> Encoding__Text rev 1
  emits []
  tests
    empty(Bytes(Seq<int>[])) => Ok(value = "")
    one(Bytes(Seq<int>[65])) => Ok(value = "QQ==")
    two(Bytes(Seq<int>[65, 66])) => Ok(value = "QUI=")
    ordered(Bytes(Seq<int>[222, 173, 190, 239])) => Ok(value = "3q2+7w==")
    notext(Bytes(Seq<int>[0, 65])) => Ok(value = "AEE=")
=
  match call std__base64__encode(value)
    given
      empty => [exchange args (value = Bytes(Seq<int>[])) outcome Ok(value = "")]
      one => [exchange args (value = Bytes(Seq<int>[65])) outcome Ok(value = "QQ==")]
      two => [exchange args (value = Bytes(Seq<int>[65, 66])) outcome Ok(value = "QUI=")]
      ordered => [exchange args (value = Bytes(Seq<int>[222, 173, 190, 239])) outcome Ok(value = "3q2+7w==")]
      notext => [exchange args (value = Bytes(Seq<int>[0, 65])) outcome Ok(value = "AEE=")]
    on Ok r => Ok(value = r.value)
`
	files := map[string]string{"text.can": string(raw), "client.can": client}
	dir := writeLSPDir(t, files)
	if diags := diagnose(dir, "client.can", client); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
