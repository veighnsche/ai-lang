package main

import (
	"strings"
	"testing"
)

// a66 (chatbot verdict A): interpreted string literals
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
    raw_slash_n("\n") => Ok(value = Bytes(Seq<int>[92, 110]))
    escaped_lf(e"\n") => Ok(value = Bytes(Seq<int>[10]))
    decoded_once(e"\\n") => Ok(value = Bytes(Seq<int>[92, 110]))
    escaped_nul(e"\0") => Ok(value = Bytes(Seq<int>[0]))
    not_octal(e"\01") => Ok(value = Bytes(Seq<int>[0, 49]))
=
  match call bytes__utf8__encode(value)
    on Ok r => Ok(value = r.value)

fn probe__b64(value: str) -> Bytes__Value rev 1
  emits [encoding.invalid_base64]
  tests
    ok_abc(e"QUJD") => Ok(value = Bytes(Seq<int>[65, 66, 67]))
    trailing_lf(e"QQ==\n") => encoding.invalid_base64(value = e"QQ==\n")
    crlf_tail(e"QUJD\r\n\r\n") => encoding.invalid_base64(value = e"QUJD\r\n\r\n")
=
  match call bytes__base64__decode(value)
    on Ok r => Ok(value = r.value)
    on encoding.invalid_base64 e => forward e
`
	dir := writeLSPDir(t, map[string]string{"probe.can": probe})
	if diags := diagnose(dir, "probe.can", probe); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// TestEscapedPatterns pins value/pattern parity: an e-pattern
// matches its decoded value at runtime. The e-arm is taken by
// its row, so a clean run proves both directions agree.
func TestEscapedPatterns(t *testing.T) {
	probe := `mod probe
  provides [probe__pat, Str__Value]
  uses []
  emits []

type Str__Value rev 1 (
  value: str
)

fn probe__pat(value: str) -> Str__Value rev 1
  emits []
  tests
    lf(e"a\nb") => Ok(value = "nl")
    other("x") => Ok(value = "other")
=
  match value
    on e"a\nb" => Ok(value = "nl")
    on _ => Ok(value = "other")
`
	dir := writeLSPDir(t, map[string]string{"probe.can": probe})
	if diags := diagnose(dir, "probe.can", probe); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// TestEscapedPatternDiagLocation pins the a66 squiggle: an
// untaken e-arm reports its source spelling, never a decoded
// control character.
func TestEscapedPatternDiagLocation(t *testing.T) {
	probe := `mod probe
  provides [probe__pat, Str__Value]
  uses []
  emits []

type Str__Value rev 1 (
  value: str
)

fn probe__pat(value: str) -> Str__Value rev 1
  emits []
  tests
    other("x") => Ok(value = "other")
=
  match value
    on e"a\nb" => Ok(value = "nl")
    on _ => Ok(value = "other")
`
	dir := writeLSPDir(t, map[string]string{"probe.can": probe})
	diags := diagnose(dir, "probe.can", probe)
	if !hasErrCode(diags, CodeArmUntaken) {
		t.Fatalf("expected CAN4107 for the untaken e-arm, got %v", diags)
	}
	for _, d := range diags {
		if d.Code != CodeArmUntaken {
			continue
		}
		if !strings.Contains(d.Msg, `on e"a\nb"`) {
			t.Fatalf("arm diag must name the source spelling, got %q", d.Msg)
		}
		if strings.Contains(d.Msg, "\n") {
			t.Fatalf("arm diag must not contain a decoded newline, got %q", d.Msg)
		}
	}
}

// TestEscapedLiteralRejectsHostEscape pins the closed set:
// a host-style \x escape inside e"..." is CAN1000, so it can
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
	dir := writeLSPDir(t, map[string]string{"probe.can": probe})
	diags := diagnose(dir, "probe.can", probe)
	if !hasErrCode(diags, CodeParse) {
		t.Fatalf("expected CAN1000 for unknown escape, got %v", diags)
	}
}
