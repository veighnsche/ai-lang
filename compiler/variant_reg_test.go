package main

import (
	"strings"
	"testing"
)

// a73: variant declaration registry. Frontend foundation
// only: no construction, proof, or emit. Probes first:
// `variant` is an unknown top-level decl.

const variantNoop = `
type M__Out rev 1 (
  value: int
)

fn m__noop(value: int) -> M__Out rev 1
  emits []
  tests
    go(value = 1) => Ok(value = 1)
=
  Ok(value = value)
`

func variantMod(provides string, decls string) string {
	return `mod m
  provides [` + provides + `]
  uses []
  emits []

` + decls + variantNoop
}

const variantLoginDecls = `variant Login__State rev 1 (
  case Anonymous()
  case Authenticated(session: Auth__Session)
  case Locked(user_id: str, remaining_seconds: int)
  case Nested(event: Login__Event)
)

variant Login__Event rev 1 (
  case Ticked()
)

type Auth__Session rev 1 (
  user_id: str
)

`

// TestVariantRegistryClean pins the positive: a valid
// three-case declaration (nullary, record payload,
// scalar payloads) diagnoses clean.
func TestVariantRegistryClean(t *testing.T) {
	src := variantMod("m__noop, M__Out, Login__State, Login__Event, Auth__Session", variantLoginDecls)
	dir := writeLSPDir(t, map[string]string{"m.ail": src})
	if diags := diagnose(dir, "m.ail", src); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

// TestVariantRegistryPin pins parent-revision resolution:
// correct rev clean, wrong rev fails CodeUsesRev.
func TestVariantRegistryPin(t *testing.T) {
	base := variantMod("m__noop, M__Out, Login__State, Login__Event, Auth__Session", variantLoginDecls)
	consumer := func(pin string) string {
		return `mod c
  provides [c__noop, C__Out]
  uses [` + pin + `]
  emits []

type C__Out rev 1 (
  value: int
)

fn c__noop(value: int) -> C__Out rev 1
  emits []
  tests
    go(value = 1) => Ok(value = 1)
=
  Ok(value = value)
`
	}
	dir := writeLSPDir(t, map[string]string{"m.ail": base, "c.ail": consumer("Login__State@1")})
	if diags := diagnose(dir, "c.ail", consumer("Login__State@1")); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
	dir2 := writeLSPDir(t, map[string]string{"m.ail": base, "c.ail": consumer("Login__State@2")})
	if diags := diagnose(dir2, "c.ail", consumer("Login__State@2")); !hasErrCode(diags, CodeUsesRev) {
		t.Fatalf("expected AIL2103, got %v", diags)
	}
}

// TestVariantRegistryRejects pins each negative rule with
// its diagnostic message.
func TestVariantRegistryRejects(t *testing.T) {
	type neg struct {
		provides string
		decls    string
		want     string
	}
	cases := map[string]neg{
		"empty variant": {
			"m__noop, M__Out, Login__State",
			"variant Login__State rev 1 (\n)\n",
			"at least one case",
		},
		"duplicate case": {
			"m__noop, M__Out, Login__State",
			"variant Login__State rev 1 (\n  case Anonymous()\n  case Anonymous()\n)\n",
			"duplicate case",
		},
		"cross-variant collision": {
			"m__noop, M__Out, Login__State, Login__Event",
			"variant Login__State rev 1 (\n  case Anonymous()\n)\n\nvariant Login__Event rev 1 (\n  case Anonymous()\n)\n",
			"collision",
		},
		"record collision": {
			"m__noop, M__Out, Login__Anonymous, Login__State",
			"type Login__Anonymous rev 1 (\n  user_id: str\n)\n\nvariant Login__State rev 1 (\n  case Anonymous()\n)\n",
			"collision",
		},
		"unknown payload": {
			"m__noop, M__Out, Login__State",
			"variant Login__State rev 1 (\n  case Authenticated(session: Nope__Missing)\n)\n",
			"unknown type",
		},
		"seq variant": {
			"m__noop, M__Out, Login__State",
			"variant Login__State rev 1 (\n  case Many(items: Seq<Login__State>)\n)\n",
			"deferred",
		},
		"direct cycle": {
			"m__noop, M__Out, Login__State",
			"variant Login__State rev 1 (\n  case Looping(next: Login__State)\n)\n",
			"cycle",
		},
		"record cycle": {
			"m__noop, M__Out, Login__State, Login__Link",
			"variant Login__State rev 1 (\n  case Wrapped(link: Login__Link)\n)\n\ntype Login__Link rev 1 (\n  next: Login__State\n)\n",
			"cycle",
		},
		"duplicate variant": {
			"m__noop, M__Out, Login__State",
			"variant Login__State rev 1 (\n  case Anonymous()\n)\n\nvariant Login__State rev 1 (\n  case Other()\n)\n",
			"double definition",
		},
		"qualified row": {
			"m__noop, M__Out, Login__State",
			"variant Login__State rev 1 (\n  case Login__Anonymous()\n)\n",
			"short case name",
		},
	}
	for name, tc := range cases {
		src := variantMod(tc.provides, tc.decls)
		dir := writeLSPDir(t, map[string]string{"m.ail": src})
		diags := diagnose(dir, "m.ail", src)
		hit := false
		for _, d := range diags {
			if d.Sev != "error" {
				continue
			}
			if strings.Contains(strings.ToLower(d.Msg), strings.ToLower(tc.want)) {
				hit = true
			}
		}
		if !hit {
			t.Fatalf("%s: expected error containing %q, got %v", name, tc.want, diags)
		}
	}
}
