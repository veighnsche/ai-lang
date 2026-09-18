package main

import (
	"os"
	"strings"
	"testing"
)

// a48 B4: html__render__utf8 is the first real Bytes consumer. The
// probe below appends a composition function to a TEMPORARY copy of
// the real html.can, so the NUL chain exercises real bodies:
// node("A<NUL>&amp;B") -> render -> exactly eight bytes, NUL
// preserved and the entity spelling undecoded.

// TestBytesB4NulChainProbe pins the frozen B4 acceptance row: the
// sealed text A<NUL>&amp;B renders to [65, 0, 38, 97, 109, 112, 59,
// 66]. The probe input carries a REAL NUL byte (Go halves glued
// around "\x00"), the only .can spelling.
func TestBytesB4NulChainProbe(t *testing.T) {
	raw, err := os.ReadFile("../std/html/html.can")
	if err != nil {
		t.Fatalf("read html.can: %v", err)
	}
	// Slice 1: html.can pins ascii consts, so the temp copy needs
	// the provider beside it.
	asciiRaw, err := os.ReadFile("../std/ascii/ascii.can")
	if err != nil {
		t.Fatalf("read ascii.can: %v", err)
	}
	// Slice 5: ascii pins Bool__Value, so the temp copy needs
	// the scalars provider beside it.
	scalarsRaw, err := os.ReadFile("../std/scalars/scalars.can")
	if err != nil {
		t.Fatalf("read scalars.can: %v", err)
	}
	probe := "\nfn html__probe__nul_chain() -> Bytes__Value rev 1\n" +
		"  emits []\n" +
		"  tests\n" +
		"    chain() => Ok(value = Bytes(Seq<int>[65, 0, 38, 97, 109, 112, 59, 66]))\n" +
		"=\n" +
		"  match call html__text__node(seal Html__Text(\"A\x00&amp;B\"))\n" +
		"    on Ok n => match call html__render__utf8(n.safe)\n" +
		"      on Ok b => Ok(value = b.value)\n"
	body := string(raw) + probe
	// The probe fn lives in the temp copy only; provide it there so
	// the copy stays well-formed under the provides rule.
	body = strings.Replace(body, ", html__asset__script]", ", html__asset__script, html__probe__nul_chain]", 1)
	// The temp copy carries real html.can scaffolding: blessed lint is
	// grandfathered, but the appended probe must be fully clean and
	// the scaffolding compiler-clean.
	seqProbeClean(t, map[string]string{"probe.can": body, "ascii.can": string(asciiRaw), "scalars.can": string(scalarsRaw)},
		"probe.can", strings.Count(string(raw), "\n")+1)
}
