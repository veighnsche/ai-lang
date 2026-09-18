package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// S2a (a91): generated-runtime agreement for the closed check
// ops. Compiles the real quota + scalars modules and executes
// all three checks under node: success passthrough, every
// violation shape, tag/path parameterization, malformed
// policy before data, and exact beyond-safe-integer renders.
// Separate driver from the S1a parity test so the pilot
// evidence stays frozen.
func TestCheckRuntimeParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("node missing: runtime parity refuses to skip: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{"../std/quota/quota.can", "../std/scalars/scalars.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	rewriteSpecifier(t, filepath.Join(out, "quota.ts"), "./scalars", "./scalars.ts")
	driver := `import { std__validate__length_check, std__validate__range_check, std__validate__membership_check } from "./quota.ts";
const show = (x) => JSON.stringify(x, (k, v) => typeof v === "bigint" ? { $bigint: v.toString() } : v);
const eq = (got, want, name) => {
  if (show(got) !== show(want)) throw new Error(name + ": got " + show(got));
};
const V = (path, rule, value) => ({ $can_kind: "validation.schema_violation", path, rule, value });
const L = (path, tag, minimum, maximum) => ({ path, tag, minimum, maximum });
const R = (path, tag, lower, upper) => ({ path, tag, lower, upper });
const M = (path, allowed) => ({ path, allowed });
// length: passthrough, data violation, malformed policy, tag pins.
eq(std__validate__length_check("Ada", L("guest", "guest", 1n, 24n)), { $can_kind: "ok", value: "Ada" }, "length ok");
eq(std__validate__length_check("", L("guest", "guest", 1n, 24n)), V("guest", "str.length_scalars", ""), "length short");
eq(std__validate__length_check("Ada", L("guest", "guest", -1n, 24n)), V("", "schema.guest_minimum", "-1"), "length min negative");
eq(std__validate__length_check("", L("guest", "guest", -1n, 24n)), V("", "schema.guest_minimum", "-1"), "length min beats data");
eq(std__validate__length_check("Ada", L("guest", "guest", 5n, 1n)), V("", "schema.guest_bounds", "minimum=5;maximum=1"), "length bounds");
eq(std__validate__length_check("", L("label", "label", 1n, 12n)), V("label", "str.length_scalars", ""), "length other path");
eq(std__validate__length_check("Ada", L("label", "label", 5n, 1n)), V("", "schema.label_bounds", "minimum=5;maximum=1"), "length other tag");
eq(std__validate__length_check("Ada", L("guest", "guest", 9007199254740993n, 1n)), V("", "schema.guest_bounds", "minimum=9007199254740993;maximum=1"), "length bigint bound");
// range: passthrough, data violation, malformed policy, tag pins.
eq(std__validate__range_check(3n, R("party", "party", 1n, 8n)), { $can_kind: "ok", value: 3n }, "range ok");
eq(std__validate__range_check(0n, R("party", "party", 1n, 8n)), V("party", "int.closed_range", "0"), "range below");
eq(std__validate__range_check(9n, R("party", "party", 1n, 8n)), V("party", "int.closed_range", "9"), "range above");
eq(std__validate__range_check(3n, R("party", "party", 8n, 1n)), V("", "schema.party_bounds", "lower=8;upper=1"), "range bounds");
eq(std__validate__range_check(11n, R("amount", "amount", 1n, 10n)), V("amount", "int.closed_range", "11"), "range other path");
eq(std__validate__range_check(5n, R("amount", "amount", 10n, 0n)), V("", "schema.amount_bounds", "lower=10;upper=0"), "range other tag");
eq(std__validate__range_check(-9007199254740993n, R("party", "party", 1n, 8n)), V("party", "int.closed_range", "-9007199254740993"), "range bigint");
// membership: passthrough, rejection, empty deny-all, path pins.
eq(std__validate__membership_check("dinner", M("slot", ["lunch", "dinner"])), { $can_kind: "ok", value: "dinner" }, "member ok");
eq(std__validate__membership_check("brunch", M("slot", ["lunch", "dinner"])), V("slot", "str.one_of", "brunch"), "member bad");
eq(std__validate__membership_check("dinner", M("slot", [])), V("slot", "str.one_of", "dinner"), "member empty allowed");
eq(std__validate__membership_check("admin", M("mode", ["interactive", "batch"])), V("mode", "str.one_of", "admin"), "member other path");
eq(std__validate__membership_check("Dinner", M("slot", ["lunch", "dinner"])), V("slot", "str.one_of", "Dinner"), "member case");
console.log("CHECK_PARITY_OK");
`
	if err := os.WriteFile(filepath.Join(out, "driver.ts"), []byte(driver), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", "driver.ts")
	cmd.Dir = out
	raw, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("node parity failed: %v\n%s", err, raw)
	}
	if !strings.Contains(string(raw), "CHECK_PARITY_OK") {
		t.Fatalf("parity driver silent:\n%s", raw)
	}
}
