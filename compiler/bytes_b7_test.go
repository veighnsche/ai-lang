package main

import (
	"os"
	"testing"
)

// a51 B7: std__utf8__decode is the public stdlib face of the B6
// kernel, and the first stdlib function over the fallible path.
// The foreign-caller probe scripts TRUE outcomes for both sides:
// success, failure with identical payload, and the NUL/BOM
// inputs the committed rows omit (text.can stays a text file).

// TestBytesB7ForeignCaller pins cross-module callability against
// the REAL text.can. Each client fn covers both arms (CAN4107):
// ok+bad share one fn with nul, bom shares one with a second
// bad input. NUL/BOM spliced (only .can spellings).
func TestBytesB7ForeignCaller(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.can")
	if err != nil {
		t.Fatalf("read text.can: %v", err)
	}
	client := `mod client
  provides [client__a, client__b]
  uses [std__utf8__decode@1]
  emits [encoding.invalid_utf8]

fn client__a(value: Bytes) -> Encoding__Text rev 1
  emits [encoding.invalid_utf8]
  tests
    ok(Bytes(Seq<int>[65])) => Ok("A")
    bad(Bytes(Seq<int>[65, 255, 66])) => encoding.invalid_utf8(value = Bytes(Seq<int>[65, 255, 66]))
    nul(Bytes(Seq<int>[97, 0, 98])) => Ok("a` + "\x00" + `b")
  match call std__utf8__decode(value)
    given
      ok => [exchange args (value = Bytes(Seq<int>[65])) outcome Ok("A")]
      bad => [exchange args (value = Bytes(Seq<int>[65, 255, 66])) outcome encoding.invalid_utf8(value = Bytes(Seq<int>[65, 255, 66]))]
      nul => [exchange args (value = Bytes(Seq<int>[97, 0, 98])) outcome Ok("a` + "\x00" + `b")]
    on Ok r => Ok(r.value)
    on encoding.invalid_utf8 e => forward e

fn client__b(value: Bytes) -> Encoding__Text rev 1
  emits [encoding.invalid_utf8]
  tests
    bom(Bytes(Seq<int>[239, 187, 191, 65])) => Ok("` + "\uFEFFA" + `")
    bad(Bytes(Seq<int>[255])) => encoding.invalid_utf8(value = Bytes(Seq<int>[255]))
  match call std__utf8__decode(value)
    given
      bom => [exchange args (value = Bytes(Seq<int>[239, 187, 191, 65])) outcome Ok("` + "\uFEFFA" + `")]
      bad => [exchange args (value = Bytes(Seq<int>[255])) outcome encoding.invalid_utf8(value = Bytes(Seq<int>[255]))]
    on Ok r => Ok(r.value)
    on encoding.invalid_utf8 e => forward e
`
	files := map[string]string{"text.can": string(raw), "client.can": client}
	dir := writeLSPDir(t, files)
	if diags := diagnose(dir, "client.can", client); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
