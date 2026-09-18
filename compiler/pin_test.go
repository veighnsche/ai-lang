package main

import (
	"strings"
	"testing"
)

// Revision pins resolve name + revision: the provider's declared rev
// must equal the pinned rev, for functions, records, and brands.

const pinProvider = `mod alpha
  provides [alpha__read, Alpha__Data, Alpha__Tag]
  uses []
  emits []

brand Alpha__Tag is str rev 2

type Alpha__Data rev 1 (
  id: str
)

fn alpha__read(id: str) -> Alpha__Data rev 1
  emits []
  tests
    r(id = "u") => Ok(id = "u")
=
  Ok(id = id)
`

const pinConsumerHead = `mod beta
  provides [beta__go, Beta__Out]
  uses [USES]
  emits []

type Beta__Out rev 1 (
  id: str
)

fn beta__go(id: str) -> Beta__Out rev 1
  emits []
  tests
    g(id = "u") => Ok(id = "u")
=
  match call alpha__read(id)
    given
      g => [exchange args (id = "u") outcome Ok(id = "u")]
    on Ok v => Ok(id = v.id)
`

func pinConsumer(uses string) string {
	return strings.Replace(pinConsumerHead, "USES", uses, 1)
}

func TestPinExactResolves(t *testing.T) {
	consumer := pinConsumer("alpha__read@1, Alpha__Data@1, Alpha__Tag@2")
	dir := writeLSPDir(t, map[string]string{"alpha.can": pinProvider, "beta.can": consumer})
	// Clean, and the scripted foreign call executes: the pin resolved.
	if diags := diagnose(dir, "beta.can", consumer); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestPinMismatchIsCAN2103(t *testing.T) {
	cases := []struct {
		name string
		uses string
		want string
	}{
		{"fn future rev", "alpha__read@999, Alpha__Data@1, Alpha__Tag@2", "pins rev 999, but alpha__read declares rev 1"},
		{"fn zero rev", "alpha__read@0, Alpha__Data@1, Alpha__Tag@2", "pins rev 0, but alpha__read declares rev 1"},
		{"type mismatch", "alpha__read@1, Alpha__Data@2, Alpha__Tag@2", "pins rev 2, but Alpha__Data declares rev 1"},
		{"brand mismatch", "alpha__read@1, Alpha__Data@1, Alpha__Tag@1", "pins rev 1, but Alpha__Tag declares rev 2"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			consumer := pinConsumer(tc.uses)
			dir := writeLSPDir(t, map[string]string{"alpha.can": pinProvider, "beta.can": consumer})
			diags := diagnose(dir, "beta.can", consumer)
			if !hasDiag(diags, "error", tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, diags)
			}
			if !hasCode(diags, "CAN2103") {
				t.Fatalf("expected CAN2103, got %v", diags)
			}
		})
	}
}

func TestPinStaleAfterProviderBump(t *testing.T) {
	bumped := strings.Replace(pinProvider, "fn alpha__read(id: str) -> Alpha__Data rev 1", "fn alpha__read(id: str) -> Alpha__Data rev 2", 1)
	consumer := pinConsumer("alpha__read@1, Alpha__Data@1, Alpha__Tag@2")
	dir := writeLSPDir(t, map[string]string{"alpha.can": bumped, "beta.can": consumer})
	diags := diagnose(dir, "beta.can", consumer)
	if !hasDiag(diags, "error", "pins rev 1, but alpha__read declares rev 2") {
		t.Fatalf("expected stale-pin error, got %v", diags)
	}
	if !hasCode(diags, "CAN2103") {
		t.Fatalf("expected CAN2103, got %v", diags)
	}
}

func TestParsePinRev(t *testing.T) {
	cases := []struct {
		in  string
		rev int
		ok  bool
	}{
		{"alpha__read@1", 1, true},
		{"alpha__read@0", 0, true},
		{"Db__User@4", 4, true},
		{"alpha__read", 0, false},
		{"alpha__read@x", 0, false},
		{"alpha__read@", 0, false},
		{"alpha__read@99999999999999999999999", -1, true},
	}
	for _, tc := range cases {
		rev, ok := parsePinRev(tc.in)
		if rev != tc.rev || ok != tc.ok {
			t.Errorf("parsePinRev(%q) = (%d, %v), want (%d, %v)", tc.in, rev, ok, tc.rev, tc.ok)
		}
	}
}
