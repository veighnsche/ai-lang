package main

import (
	"strings"
	"testing"
)

// Bodies need no `=` separator: the first non-metadata row is the
// body by construction (isMetaHead mirrors the metadata switch).
const eqBodyFree = `mod m
  provides [m__go, M__O]
  uses []
  emits []

type M__O rev 1 (
  value: int
)

fn m__go(value: int) -> M__O rev 1
  emits []
  tests
    one(1) => Ok(value = 1)
  Ok(value = value)
`

func TestParseBodyWithoutEquals(t *testing.T) {
	m, err := parseModuleText("m.can", eqBodyFree)
	if err != nil {
		t.Fatalf("equals-free body must parse: %v", err)
	}
	var body bool
	for _, d := range m.Decls {
		if fn, ok := d.(*FnDecl); ok && fn.Name == "m__go" && fn.Body != nil {
			body = true
		}
	}
	if !body {
		t.Fatalf("m__go body missing after parse")
	}
}

// The old lone `=` fails loudly with a migration pointer,
// never silently.
func TestParseLoneEqualsRejected(t *testing.T) {
	old := strings.Replace(eqBodyFree, "    one(1) => Ok(value = 1)\n", "    one(1) => Ok(value = 1)\n=\n", 1)
	_, err := parseModuleText("m.can", old)
	if err == nil || !strings.Contains(err.Error(), "lone = separator removed") {
		t.Fatalf("expected lone-= rejection, got %v", err)
	}
}

// A header with no body still fails, naming the function.
func TestParseMissingBody(t *testing.T) {
	nobody := strings.Replace(eqBodyFree, "  Ok(value = value)\n", "", 1)
	_, err := parseModuleText("m.can", nobody)
	if err == nil || !strings.Contains(err.Error(), "missing body") {
		t.Fatalf("expected missing-body error, got %v", err)
	}
}
