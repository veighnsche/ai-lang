package main

import (
	"os"
	"strings"
	"testing"
)

// v60 B15: std__base64__decode is the public stdlib face of the
// B14 kernel, and this slice closes the workstream with wrapper-
// level round trips through real bodies. The foreign-caller
// probe pins §6 callability against the REAL text.ail; the
// round-trip probe appends composition fns to a TEMP COPY and
// executes decode(encode(x)) = x per codec plus a cross check,
// so fakerows prove nothing.

// TestBytesB15ForeignCaller pins cross-module callability,
// both arms, the NUL byte, and non-ASCII malformed input.
func TestBytesB15ForeignCaller(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.ail")
	if err != nil {
		t.Fatalf("read text.ail: %v", err)
	}
	client := `mod client
  provides [client__go]
  uses [std__base64__decode@1]
  emits []

fn client__go(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_base64]
  tests
    empty(value = "") => Ok(value = Bytes(Seq<int>[]))
    one(value = "QQ==") => Ok(value = Bytes(Seq<int>[65]))
    nul(value = "AA==") => Ok(value = Bytes(Seq<int>[0]))
    ab(value = "QUI=") => Ok(value = Bytes(Seq<int>[65, 66]))
    ordered(value = "3q2+7w==") => Ok(value = Bytes(Seq<int>[222, 173, 190, 239]))
    odd(value = "QUI") => encoding.invalid_base64(value = "QUI")
    nonascii(value = "é===") => encoding.invalid_base64(value = "é===")
=
  match call std__base64__decode(value = value)
    given
      empty => [exchange args (value = "") outcome Ok(value = Bytes(Seq<int>[]))]
      one => [exchange args (value = "QQ==") outcome Ok(value = Bytes(Seq<int>[65]))]
      nul => [exchange args (value = "AA==") outcome Ok(value = Bytes(Seq<int>[0]))]
      ab => [exchange args (value = "QUI=") outcome Ok(value = Bytes(Seq<int>[65, 66]))]
      ordered => [exchange args (value = "3q2+7w==") outcome Ok(value = Bytes(Seq<int>[222, 173, 190, 239]))]
      odd => [exchange args (value = "QUI") outcome encoding.invalid_base64(value = "QUI")]
      nonascii => [exchange args (value = "é===") outcome encoding.invalid_base64(value = "é===")]
    on Ok r => Ok(value = r.value)
    on encoding.invalid_base64 e => encoding.invalid_base64(value = e.value)
`
	files := map[string]string{"text.ail": string(raw), "client.ail": client}
	dir := writeLSPDir(t, files)
	if diags := diagnose(dir, "client.ail", client); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// TestBytesB15RoundTrip executes wrapper-level round trips
// through real bodies: decode(encode(x)) = x per codec, plus
// hex-of-b64-decoded-bytes as the cross check.
func TestBytesB15RoundTrip(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.ail")
	if err != nil {
		t.Fatalf("read text.ail: %v", err)
	}
	probes := `
fn text__probe__rt_utf8(value: str) -> Encoding__Text rev 1
  emits [encoding.invalid_utf8]
  tests
    rt(value = "Aé") => Ok(value = "Aé")
=
  match call std__utf8__encode(value = value)
    on Ok b => match call std__utf8__decode(value = b.value)
      on Ok t => Ok(value = t.value)
      on encoding.invalid_utf8 e => encoding.invalid_utf8(value = e.value)

fn text__probe__rt_hex(value: str) -> Encoding__Text rev 1
  emits [encoding.invalid_hex]
  tests
    rt(value = "deadbeef") => Ok(value = "deadbeef")
=
  match call std__hex__decode(value = value)
    on Ok b => match call std__hex__encode(value = b.value)
      on Ok t => Ok(value = t.value)
    on encoding.invalid_hex e => encoding.invalid_hex(value = e.value)

fn text__probe__rt_b64(value: str) -> Encoding__Text rev 1
  emits [encoding.invalid_base64]
  tests
    rt(value = "3q2+7w==") => Ok(value = "3q2+7w==")
=
  match call std__base64__decode(value = value)
    on Ok b => match call std__base64__encode(value = b.value)
      on Ok t => Ok(value = t.value)
    on encoding.invalid_base64 e => encoding.invalid_base64(value = e.value)

fn text__probe__cross(value: str) -> Encoding__Text rev 1
  emits [encoding.invalid_base64]
  tests
    cross(value = "3q2+7w==") => Ok(value = "deadbeef")
=
  match call std__base64__decode(value = value)
    on Ok b => match call std__hex__encode(value = b.value)
      on Ok t => Ok(value = t.value)
    on encoding.invalid_base64 e => encoding.invalid_base64(value = e.value)
`
	copy := strings.Replace(string(raw),
		"Split__Result]",
		"Split__Result, text__probe__rt_utf8, text__probe__rt_hex, text__probe__rt_b64, text__probe__cross]",
		1) + probes
	files := map[string]string{"text.ail": copy}
	dir := writeLSPDir(t, files)
	if diags := diagnose(dir, "text.ail", copy); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
	if !strings.Contains(copy, "text__probe__rt_b64") {
		t.Fatal("probe copy lost the b64 round-trip fn")
	}
}
