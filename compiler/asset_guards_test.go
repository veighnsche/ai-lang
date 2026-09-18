package main

// Asset fixture-root guards (S5 slice plan): the properties that make
// fixtures prove the check runs, never that an asset is blessed.
// Structural, not behavioral: they fail if the shape rots, whatever
// the rows say.

import (
	"os"
	"testing"
)

func assetGuardModules(t *testing.T) (*Module, *Module) {
	t.Helper()
	htmlRaw, err := os.ReadFile("../std/html/html.ail")
	if err != nil {
		t.Fatal(err)
	}
	schemaRaw, err := os.ReadFile("../std/schema/schema.ail")
	if err != nil {
		t.Fatal(err)
	}
	htmlMod, err := parseModuleText("html.ail", string(htmlRaw))
	if err != nil {
		t.Fatalf("parse html: %v", err)
	}
	schemaMod, err := parseModuleText("schema.ail", string(schemaRaw))
	if err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	return htmlMod, schemaMod
}

// TestAssetNoAmbientAuthority pins the fixture boundary: the schema
// module holds no state cells and declares no effects, so there is
// no ambient registry to poison — every snapshot is an explicit
// value, and fixture snapshots differ from production ones only in
// provenance (which the language cannot see by construction).
func TestAssetNoAmbientAuthority(t *testing.T) {
	_, schemaMod := assetGuardModules(t)
	for _, d := range schemaMod.Decls {
		if _, ok := d.(*StateDecl); ok {
			t.Fatalf("schema declares state: snapshots must stay explicit values")
		}
		if fn, ok := d.(*FnDecl); ok && len(fn.Effects) != 0 {
			t.Fatalf("%s declares effects: the approval path stays pure", fn.Name)
		}
	}
	if len(schemaMod.Hdr["uses"]) != 0 {
		t.Fatalf("schema declares uses: the approval core stays dependency-free")
	}
}

// TestAssetNoScriptedEvidence pins execution over scripting: no
// match node in the schema bodies or the two sinks carries a given
// table, so committed rows execute the real checks — a sealed
// success-shaped value alone never approves, and no row can
// manufacture approval by stubbing a helper.
func TestAssetNoScriptedEvidence(t *testing.T) {
	htmlMod, schemaMod := assetGuardModules(t)
	sinks := map[string]bool{"html__asset__stylesheet": true, "html__asset__script": true}
	for _, m := range []*Module{htmlMod, schemaMod} {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			if m == htmlMod && !sinks[fn.Name] {
				continue
			}
			for _, n := range matchNodes(fn.Body) {
				if n.Given != nil {
					t.Fatalf("%s scripts a callee: approval-path checks must execute", fn.Name)
				}
			}
		}
	}
}

// TestAssetSinkPurity pins the sinks dependency-free: the builders
// call nothing but the certified kernel (no uses entries), so no
// helper can smuggle unapproved strings into the fixed elements.
func TestAssetSinkPurity(t *testing.T) {
	htmlMod, _ := assetGuardModules(t)
	if len(htmlMod.Hdr["uses"]) != 0 {
		t.Fatalf("html declares uses: sinks must stay kernel-only")
	}
}

// TestAssetHostileSet pins the committed hostile set: every planned
// rejection class keeps its row. Deleting a row fails here, not
// silently.
func TestAssetHostileSet(t *testing.T) {
	htmlMod, schemaMod := assetGuardModules(t)
	rows := map[string][]string{}
	for _, m := range []*Module{htmlMod, schemaMod} {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			for _, tst := range fn.Tests {
				rows[fn.Name] = append(rows[fn.Name], tst.Name)
			}
		}
	}
	has := func(fn, row string) bool {
		for _, r := range rows[fn] {
			if r == row {
				return true
			}
		}
		return false
	}
	want := [][2]string{
		{"schema__asset__approve", "approve_unknown_url"},
		{"schema__asset__approve", "approve_wrong_digest"},
		{"schema__asset__approve", "approve_policy_mismatch"},
		{"schema__asset__approve", "approve_revoked"},
		{"schema__asset__approve", "approve_role_mismatch"},
		{"schema__asset__approve", "approve_expired"},
		{"schema__asset__approve", "approve_stale_sequence"},
		{"schema__asset__approve", "approve_no_transitive"},
		{"schema__asset__approve", "approve_conflict"},
		{"schema__asset__approve", "approve_quote_url"},
		{"schema__asset__approve", "approve_control_url"},
		{"schema__asset__approve", "approve_lower_edge"},
		{"schema__asset__approve", "approve_inner_upper"},
		{"schema__asset__recheck", "recheck_tampered"},
		{"schema__asset__recheck", "recheck_mixed_policy"},
		{"html__asset__stylesheet", "stylesheet_role_mismatch"},
		{"html__asset__script", "script_role_mismatch"},
		{"html__asset__script", "script_module_unsupported"},
		{"html__asset__script", "script_worker_unsupported"},
		{"html__asset__script", "script_preload_unsupported"},
	}
	for _, w := range want {
		if !has(w[0], w[1]) {
			t.Fatalf("hostile set lost row %s/%s", w[0], w[1])
		}
	}
}
