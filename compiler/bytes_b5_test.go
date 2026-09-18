package main

import (
	"os"
	"testing"
)

// a49 B5: std__utf8__encode is the public stdlib face of the B3
// kernel. The foreign-caller probe below proves the §6 claim that
// ordinary wrappers are normal functions: a separate client module
// calls the wrapper through uses [std__utf8__encode@1], which the
// wrapper's own decision-table rows cannot show. It also pins NUL
// relay (text.can itself stays a text file — no committed raw NUL).

// TestBytesB5ForeignCaller pins cross-module callability against the
// REAL text.can: nul feeds a<NUL>b (REAL NUL, spliced — the only
// .can spelling), computed feeds Aé.
func TestBytesB5ForeignCaller(t *testing.T) {
	raw, err := os.ReadFile("../std/text/text.can")
	if err != nil {
		t.Fatalf("read text.can: %v", err)
	}
	client := `mod client
  provides [client__nul, client__computed]
  uses [std__utf8__encode@1]
  emits []

fn client__nul() -> Bytes__Value rev 1
  emits []
  tests
    nul() => Ok(value = Bytes(Seq<int>[97, 0, 98]))
=
  match call std__utf8__encode(value = "a` + "\x00" + `b")
    given
      nul => [exchange args (value = "a` + "\x00" + `b") outcome Ok(value = Bytes(Seq<int>[97, 0, 98]))]
    on Ok r => Ok(value = r.value)

fn client__computed() -> Bytes__Value rev 1
  emits []
  tests
    computed() => Ok(value = Bytes(Seq<int>[65, 195, 169]))
=
  match call std__utf8__encode(value = "Aé")
    given
      computed => [exchange args (value = "Aé") outcome Ok(value = Bytes(Seq<int>[65, 195, 169]))]
    on Ok r => Ok(value = r.value)
`
	files := map[string]string{"text.can": string(raw), "client.can": client}
	dir := writeLSPDir(t, files)
	if diags := diagnose(dir, "client.can", client); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
