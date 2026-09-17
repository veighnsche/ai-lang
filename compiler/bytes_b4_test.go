package main

import (
	"os"
	"strings"
	"testing"
)

// a48 B4: html__render__utf8 is the first real Bytes consumer. The
// probe below appends a composition function to a TEMPORARY copy of
// the real html.ail, so the NUL chain exercises real bodies:
// node("A<NUL>&amp;B") -> render -> exactly eight bytes, NUL
// preserved and the entity spelling undecoded.

// TestBytesB4NulChainProbe pins the frozen B4 acceptance row: the
// sealed text A<NUL>&amp;B renders to [65, 0, 38, 97, 109, 112, 59,
// 66]. The probe input carries a REAL NUL byte (Go halves glued
// around "\x00"), the only .ail spelling.
func TestBytesB4NulChainProbe(t *testing.T) {
	raw, err := os.ReadFile("../std/html/html.ail")
	if err != nil {
		t.Fatalf("read html.ail: %v", err)
	}
	probe := "\nfn html__probe__nul_chain() -> Bytes__Value rev 1\n" +
		"  emits []\n" +
		"  tests\n" +
		"    chain() => Ok(value = Bytes(Seq<int>[65, 0, 38, 97, 109, 112, 59, 66]))\n" +
		"=\n" +
		"  match call html__text__node(text = seal Html__Text(\"A\x00&amp;B\"))\n" +
		"    on Ok n => match call html__render__utf8(document = n.safe)\n" +
		"      on Ok b => Ok(value = b.value)\n"
	body := string(raw) + probe
	// The probe fn lives in the temp copy only; provide it there so
	// the copy stays well-formed under the provides rule.
	body = strings.Replace(body, "]\n  uses []", ", html__probe__nul_chain]\n  uses []", 1)
	dir := writeLSPDir(t, map[string]string{"probe.ail": body})
	if diags := diagnose(dir, "probe.ail", body); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
