package main

import (
	"testing"
)

// v66 (chatbot verdict A): interpreted string literals
// e"..." with exactly six escapes (\" \\ \n \r \t \0),
// decoded once left to right. Ordinary "..." stays raw.
// Probes first: the e-forms do not parse yet.

// TestEscapedLiterals pins value bytes through the utf8
// kernel and control survival through the base64 kernel.
// Numeric expectations guard against a broken decoder
// corrupting input and expectation together.
func TestEscapedLiterals(t *testing.T) {
	probe := `mod probe
  provides [probe__utf8, probe__b64]
  uses []
  emits []

fn probe__utf8(value: str) -> Bytes__Value rev 1
  emits []
  tests
    raw_slash_n(value = "\n") => Ok(value = Bytes(Seq<int>[92, 110]))
    escaped_lf(value = e"\n") => Ok(value = Bytes(Seq<int>[10]))
    decoded_once(value = e"\\n") => Ok(value = Bytes(Seq<int>[92, 110]))
    escaped_nul(value = e"\0") => Ok(value = Bytes(Seq<int>[0]))
    not_octal(value = e"\01") => Ok(value = Bytes(Seq<int>[0, 49]))
=
  match call bytes__utf8__encode(value)
    on Ok r => Ok(value = r.value)

fn probe__b64(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_base64]
  tests
    trailing_lf(value = e"QQ==\n") => encoding.invalid_base64(value = e"QQ==\n")
    crlf_tail(value = e"QUJD\r\n\r\n") => encoding.invalid_base64(value = e"QUJD\r\n\r\n")
=
  match call bytes__base64__decode(value)
    on Ok r => Ok(value = r.value)
    on encoding.invalid_base64 e => encoding.invalid_base64(value = e.value)
`
	dir := writeLSPDir(t, map[string]string{"probe.ail": probe})
	if diags := diagnose(dir, "probe.ail", probe); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// TestEscapedLiteralRejectsHostEscape pins the closed set:
// a host-style \x escape inside e"..." is AIL1000, so it can
// never silently rewrite a codec vector into valid base64.
func TestEscapedLiteralRejectsHostEscape(t *testing.T) {
	probe := `mod probe
  provides [probe__go]
  uses []
  emits []

fn probe__go(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_base64]
  tests
    hexescape(value = e"\x51Q==") => encoding.invalid_base64(value = e"\x51Q==")
=
  match call bytes__base64__decode(value)
    on Ok r => Ok(value = r.value)
    on encoding.invalid_base64 e => encoding.invalid_base64(value = e.value)
`
	dir := writeLSPDir(t, map[string]string{"probe.ail": probe})
	diags := diagnose(dir, "probe.ail", probe)
	if !hasErrCode(diags, CodeParse) {
		t.Fatalf("expected AIL1000 for unknown escape, got %v", diags)
	}
}
