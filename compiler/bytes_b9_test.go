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
// the REAL text.ail.
func TestBytesB9ForeignCaller(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.ail")
	if err != nil {
		t.Fatalf("read text.ail: %v", err)
	}
	client := `mod client
  provides [client__go]
  uses [std__hex__encode@1]
  emits []

fn client__go(value: Bytes) -> Encoding__Text rev 1
  emits []
  tests
    empty(value = Bytes(Seq<int>[])) => Ok(value = "")
    ff(value = Bytes(Seq<int>[255])) => Ok(value = "ff")
    lower(value = Bytes(Seq<int>[171])) => Ok(value = "ab")
    ordered(value = Bytes(Seq<int>[222, 173, 190, 239])) => Ok(value = "deadbeef")
    notext(value = Bytes(Seq<int>[65, 66])) => Ok(value = "4142")
=
  match call std__hex__encode(value = value)
    given
      empty => [exchange args (value = Bytes(Seq<int>[])) outcome Ok(value = "")]
      ff => [exchange args (value = Bytes(Seq<int>[255])) outcome Ok(value = "ff")]
      lower => [exchange args (value = Bytes(Seq<int>[171])) outcome Ok(value = "ab")]
      ordered => [exchange args (value = Bytes(Seq<int>[222, 173, 190, 239])) outcome Ok(value = "deadbeef")]
      notext => [exchange args (value = Bytes(Seq<int>[65, 66])) outcome Ok(value = "4142")]
    on Ok r => Ok(value = r.value)
`
	files := map[string]string{"text.ail": string(raw), "client.ail": client}
	dir := writeLSPDir(t, files)
	if diags := diagnose(dir, "client.ail", client); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
