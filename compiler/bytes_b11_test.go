package main

import (
	"os"
	"testing"
)

// a57 B11: std__hex__decode is the public stdlib face of the B10
// kernel. The foreign-caller probe scripts TRUE outcomes for both
// sides, including the NUL/non-ASCII malformed inputs the
// committed rows omit (text.ail stays a text file).

// TestBytesB11ForeignCaller pins cross-module callability against
// the REAL text.ail. Each client fn covers both arms (AIL4107).
func TestBytesB11ForeignCaller(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.ail")
	if err != nil {
		t.Fatalf("read text.ail: %v", err)
	}
	client := `mod client
  provides [client__a, client__b]
  uses [std__hex__decode@1]
  emits [encoding.invalid_hex]

fn client__a(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_hex]
  tests
    ok(value = "41") => Ok(value = Bytes(Seq<int>[65]))
    bad(value = "41zz42") => encoding.invalid_hex(value = "41zz42")
    nul(value = "a` + "\x00" + `b") => encoding.invalid_hex(value = "a` + "\x00" + `b")
=
  match call std__hex__decode(value = value)
    given
      ok => [exchange args (value = "41") outcome Ok(value = Bytes(Seq<int>[65]))]
      bad => [exchange args (value = "41zz42") outcome encoding.invalid_hex(value = "41zz42")]
      nul => [exchange args (value = "a` + "\x00" + `b") outcome encoding.invalid_hex(value = "a` + "\x00" + `b")]
    on Ok r => Ok(value = r.value)
    on encoding.invalid_hex e => encoding.invalid_hex(value = e.value)

fn client__b(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_hex]
  tests
    ok(value = "ff") => Ok(value = Bytes(Seq<int>[255]))
    cjk(value = "0中") => encoding.invalid_hex(value = "0中")
    bad(value = "f") => encoding.invalid_hex(value = "f")
=
  match call std__hex__decode(value = value)
    given
      ok => [exchange args (value = "ff") outcome Ok(value = Bytes(Seq<int>[255]))]
      cjk => [exchange args (value = "0中") outcome encoding.invalid_hex(value = "0中")]
      bad => [exchange args (value = "f") outcome encoding.invalid_hex(value = "f")]
    on Ok r => Ok(value = r.value)
    on encoding.invalid_hex e => encoding.invalid_hex(value = e.value)
`
	files := map[string]string{"text.ail": string(raw), "client.ail": client}
	dir := writeLSPDir(t, files)
	if diags := diagnose(dir, "client.ail", client); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
