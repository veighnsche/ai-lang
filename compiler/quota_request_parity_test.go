package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// S1a (a89): generated-runtime agreement for the quota-request
// pilot. Compiles the real blessed quota + scalars modules and
// executes validate/admit under node: success preservation,
// every violation shape, the frozen payload mapping, fail-fast
// priority, exact beyond-safe-integer rendering, and the
// admission gate (rejected requests never consume quota).
func TestQuotaRequestRuntimeParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("node missing: runtime parity refuses to skip: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{"../std/quota/quota.can", "../std/scalars/scalars.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	// Same harness shim as the P1 parity test (a89 amendment 1).
	rewriteSpecifier(t, filepath.Join(out, "quota.ts"), "./scalars", "./scalars.ts")
	driver := `import { quota__request__validate, quota__request__admit } from "./quota.ts";
const show = (x) => JSON.stringify(x, (k, v) => typeof v === "bigint" ? { $bigint: v.toString() } : v);
const eq = (got, want, name) => {
  if (show(got) !== show(want)) throw new Error(name + ": got " + show(got));
};
const V = (path, rule, value) => ({ $can_kind: "validation.schema_violation", path, rule, value });
const S = { label_minimum: 1n, label_maximum: 12n, amount_lower: 1n, amount_upper: 10n, allowed_modes: ["interactive", "batch"] };
const R = (label, amount, mode) => ({ label, amount, mode });
// validate: success preserves every field.
eq(quota__request__validate(R("Ada", 5n, "batch"), S), { $can_kind: "ok", value: R("Ada", 5n, "batch") }, "validate ok");
eq(quota__request__validate(R(" ab ", 5n, "batch"), S), { $can_kind: "ok", value: R(" ab ", 5n, "batch") }, "whitespace preserved");
// validate: each violation shape and the frozen mapping.
eq(quota__request__validate(R("", 5n, "batch"), S), V("label", "str.length_scalars", ""), "label short");
eq(quota__request__validate(R("Ada", 11n, "batch"), S), V("amount", "int.closed_range", "11"), "amount above");
eq(quota__request__validate(R("Ada", 5n, "admin"), S), V("mode", "str.one_of", "admin"), "mode bad");
eq(quota__request__validate(R("Ada", 5n, "batch"), { ...S, label_minimum: 5n, label_maximum: 1n }), V("", "schema.label_bounds", "minimum=5;maximum=1"), "label bounds");
eq(quota__request__validate(R("Ada", 5n, "batch"), { ...S, label_minimum: -1n }), V("", "schema.label_minimum", "-1"), "label min negative");
eq(quota__request__validate(R("Ada", 5n, "batch"), { ...S, allowed_modes: [] }), V("mode", "str.one_of", "batch"), "modes empty");
// validate: fail-fast priority — earlier data failure wins.
eq(quota__request__validate(R("", 5n, "batch"), { ...S, amount_lower: 10n, amount_upper: 0n }), V("label", "str.length_scalars", ""), "label beats policy");
eq(quota__request__validate(R("Ada", 99n, "admin"), S), V("amount", "int.closed_range", "99"), "amount beats mode");
// validate: exact rendering past the host safe-integer range.
eq(quota__request__validate(R("Ada", -9007199254740993n, "batch"), S), V("amount", "int.closed_range", "-9007199254740993"), "bigint exact");
// admit: rejected requests never consume; successes accumulate.
eq(quota__request__admit(R("", 3n, "batch"), S, 10n), V("label", "str.length_scalars", ""), "admit rejected");
eq(quota__request__admit(R("Ada", 3n, "batch"), S, 10n), { $can_kind: "ok", used: 3n, remaining: 7n }, "admit first");
eq(quota__request__admit(R("Ada", 2n, "batch"), S, 10n), { $can_kind: "ok", used: 5n, remaining: 5n }, "admit second");
eq(quota__request__admit(R("Ada", 10n, "batch"), S, 5n), { $can_kind: "validation.out_of_range", value: 15n, lower: 0n, upper: 5n }, "admit over");
eq(quota__request__admit(R("Ada", 5n, "batch"), S, -1n), { $can_kind: "validation.invalid_bounds", lower: 0n, upper: -1n }, "admit bad quota");
eq(quota__request__admit(R("Ada", -2n, "batch"), { ...S, amount_lower: -5n }, 10n), { $can_kind: "validation.negative_value", value: -2n }, "admit negative");
console.log("REQUEST_PARITY_OK");
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
	if !strings.Contains(string(raw), "REQUEST_PARITY_OK") {
		t.Fatalf("parity driver silent:\n%s", raw)
	}
}
