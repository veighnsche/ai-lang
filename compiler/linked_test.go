package main

import (
	"strings"
	"testing"
)

// a68: runLinkedPure executes linked pure bodies across
// module boundaries with no given consumed. The verdict
// falsifier: every unit suite is green, the composition
// is wrong, and only linked execution sees it.

// leafMiddleClient returns the three verdict modules. The
// leaf bug ("" for "A") hides behind empty-only rows;
// middle relays through given; client scripts "A".
func leafMiddleClient(broken bool) map[string]string {
	leafTests := `    empty(value = "") => Ok(value = "")`
	leafBody := `  Ok(value = "")`
	if !broken {
		leafTests = `    empty(value = "") => Ok(value = "")
    a(value = "A") => Ok(value = "A")
    other(value = "B") => Ok(value = "")`
		leafBody = `  match value
    on "" => Ok(value = "")
    on "A" => Ok(value = "A")
    on _ => Ok(value = "")`
	}
	leaf := `mod leaf
  provides [leaf__copy]
  uses []
  emits []

fn leaf__copy(value: str) -> Encoding__Text rev 1
  emits []
  tests
` + leafTests + `
=
` + leafBody + `
`
	middle := `mod middle
  provides [middle__copy]
  uses [leaf__copy@1]
  emits []

fn middle__copy(value: str) -> Encoding__Text rev 1
  emits []
  tests
    empty(value = "") => Ok(value = "")
    a(value = "A") => Ok(value = "A")
=
  match call leaf__copy(value = value)
    given
      empty => [exchange args (value = "") outcome Ok(value = "")]
      a => [exchange args (value = "A") outcome Ok(value = "A")]
    on Ok r => Ok(value = r.value)
`
	client := `mod client
  provides [client__go]
  uses [middle__copy@1]
  emits []

fn client__go(value: str) -> Encoding__Text rev 1
  emits []
  tests
    a(value = "A") => Ok(value = "A")
=
  match call middle__copy(value = value)
    given
      a => [exchange args (value = "A") outcome Ok(value = "A")]
    on Ok r => Ok(value = r.value)
`
	return map[string]string{"leaf.can": leaf, "middle.can": middle, "client.can": client}
}

// TestLinkedPureMismatch pins the falsifier: linked
// execution reports actual "" against expected "A".
func TestLinkedPureMismatch(t *testing.T) {
	for _, order := range [][]string{
		{"leaf.can", "middle.can", "client.can"},
		{"client.can", "middle.can", "leaf.can"},
	} {
		files := leafMiddleClient(true)
		err := runLinkedPure(t, files, order, "middle__copy", 1,
			map[string]string{"value": `"A"`}, `Ok(value = "A")`)
		if err == nil {
			t.Fatalf("order %v: expected mismatch failure, got pass", order)
		}
		if !strings.Contains(err.Error(), `""`) || !strings.Contains(err.Error(), `"A"`) {
			t.Fatalf("order %v: failure must name actual and expected, got %v", order, err)
		}
	}
}

// TestLinkedPureControl pins the positive: a correct leaf
// links clean.
func TestLinkedPureControl(t *testing.T) {
	files := leafMiddleClient(false)
	if err := runLinkedPure(t, files,
		[]string{"leaf.can", "middle.can", "client.can"},
		"middle__copy", 1,
		map[string]string{"value": `"A"`}, `Ok(value = "A")`); err != nil {
		t.Fatalf("expected linked pass, got %v", err)
	}
}

// TestLinkedPureRefusals pins the admission walk: externs,
// state ops, effectful fns, and unknown calls refuse
// before execution.
func TestLinkedPureRefusals(t *testing.T) {
	base := `mod m
  provides [m__go]
  uses []
  emits []

`
	cases := map[string]string{
		"extern": `mod m
  provides [m__go, ext__thing]
  uses []
  emits []

` + `extern ext__thing(value: str) -> Encoding__Text rev 1
  emits []

fn m__go(value: str) -> Encoding__Text rev 1
  emits []
  tests
    a(value = "A") => Ok(value = "A")
=
  match call ext__thing(value = value)
    given
      a => [exchange args (value = "A") outcome Ok(value = "A")]
    on Ok r => Ok(value = r.value)
`,
		"unknown": base + `fn m__go(value: str) -> Encoding__Text rev 1
  emits []
  tests
    a(value = "A") => Ok(value = "A")
=
  match call nope__missing(value = value)
    given
      a => [exchange args (value = "A") outcome Ok(value = "A")]
    on Ok r => Ok(value = r.value)
`,
		"state": `mod m
  provides [m__go, M__T]
  uses []
  emits []

state M__C: int = 0

type M__T rev 1 (
  total: int
)

fn m__go(value: str) -> M__T rev 1
  effects [M__C.read]
  emits []
  tests
    a(value = "A") => Ok(total = 0)
=
  match call state__get(M__C)
    on Ok c => Ok(total = c.value)
`,
	}
	for name, src := range cases {
		err := runLinkedPure(t, map[string]string{"m.can": src},
			[]string{"m.can"}, "m__go", 1,
			map[string]string{"value": `"A"`}, `Ok(value = "A")`)
		if err == nil {
			t.Fatalf("%s: expected refusal, got pass", name)
		}
		if name != "unknown" && !strings.Contains(err.Error(), "linked execution refused") {
			t.Fatalf("%s: admission walk must refuse, got %v", name, err)
		}
	}
}
