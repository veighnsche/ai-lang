package main

import (
	"os"
	"testing"
)

// a53 B9: std__hex__encode is the public stdlib face of the B8
// kernel. The foreign-caller probe proves the §6 claim that
// ordinary wrappers are normal functions, plus lowercase,
// order, and no-text-interpretation vectors.

// TestBytesB9ForeignCaller pins cross-module callability against
// the REAL text.can.
func TestBytesB9ForeignCaller(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.can")
	if err != nil {
		t.Fatalf("read text.can: %v", err)
	}
	client := `mod client
  provides [client__go]
  uses [std__hex__encode@1]
  emits []

fn client__go(value: Bytes) -> Encoding__Text rev 1
  emits []
  tests
    empty(Bytes(Seq<int>[])) => Ok("")
    ff(Bytes(Seq<int>[255])) => Ok("ff")
    lower(Bytes(Seq<int>[171])) => Ok("ab")
    ordered(Bytes(Seq<int>[222, 173, 190, 239])) => Ok("deadbeef")
    notext(Bytes(Seq<int>[65, 66])) => Ok("4142")
  match call std__hex__encode(value)
    given
      empty => [exchange args (value = Bytes(Seq<int>[])) outcome Ok("")]
      ff => [exchange args (value = Bytes(Seq<int>[255])) outcome Ok("ff")]
      lower => [exchange args (value = Bytes(Seq<int>[171])) outcome Ok("ab")]
      ordered => [exchange args (value = Bytes(Seq<int>[222, 173, 190, 239])) outcome Ok("deadbeef")]
      notext => [exchange args (value = Bytes(Seq<int>[65, 66])) outcome Ok("4142")]
    on Ok r => Ok(r.value)
`
	files := map[string]string{"text.can": string(raw), "client.can": client}
	dir := writeLSPDir(t, files)
	if diags := diagnose(dir, "client.can", client); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
