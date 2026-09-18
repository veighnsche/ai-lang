package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// S2 second consumer: generated-runtime agreement for the
// reservation desk. Compiles reserve + quota + scalars and
// executes validate/admit under node: success preservation,
// every violation shape, the frozen payload mapping, fail-fast
// priority, exact beyond-safe-integer rendering, and the
// admission gate (rejected requests never confirm).
func TestReserveRuntimeParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("node missing: runtime parity refuses to skip: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/reserve/reserve.can", "../std/quota/quota.can", "../std/scalars/scalars.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	// Same harness shim as the S1a parity test (a89 amendment 1):
	// extensionless same-dir specifiers become explicit .ts paths
	// in the temp copy only; goldens pin emitted bytes verbatim.
	rewriteSpecifier(t, filepath.Join(out, "reserve.ts"), "./quota", "./quota.ts")
	// No ./scalars rewrite since migration: reserve never calls
	// the converter directly, so its emit imports quota alone.
	rewriteSpecifier(t, filepath.Join(out, "quota.ts"), "./scalars", "./scalars.ts")
	driver := `import { reserve__request__validate, reserve__request__admit } from "./reserve.ts";
const show = (x) => JSON.stringify(x, (k, v) => typeof v === "bigint" ? { $bigint: v.toString() } : v);
const eq = (got, want, name) => {
  if (show(got) !== show(want)) throw new Error(name + ": got " + show(got));
};
const V = (path, rule, value) => ({ $can_kind: "validation.schema_violation", path, rule, value });
const S = { guest_minimum: 1n, guest_maximum: 24n, party_lower: 1n, party_upper: 8n, allowed_slots: ["lunch", "dinner"] };
const R = (guest, party, slot) => ({ guest, party, slot });
// validate: success preserves every field.
eq(reserve__request__validate(R("Ada", 3n, "dinner"), S), { $can_kind: "ok", value: R("Ada", 3n, "dinner") }, "validate ok");
eq(reserve__request__validate(R(" ab ", 3n, "dinner"), S), { $can_kind: "ok", value: R(" ab ", 3n, "dinner") }, "whitespace preserved");
// validate: each violation shape and the frozen mapping.
eq(reserve__request__validate(R("", 3n, "dinner"), S), V("guest", "str.length_scalars", ""), "guest short");
eq(reserve__request__validate(R("Ada", 9n, "dinner"), S), V("party", "int.closed_range", "9"), "party above");
eq(reserve__request__validate(R("Ada", 3n, "brunch"), S), V("slot", "str.one_of", "brunch"), "slot bad");
eq(reserve__request__validate(R("Ada", 3n, "dinner"), { ...S, guest_minimum: 5n, guest_maximum: 1n }), V("", "schema.guest_bounds", "minimum=5;maximum=1"), "guest bounds");
eq(reserve__request__validate(R("Ada", 3n, "dinner"), { ...S, guest_minimum: -1n }), V("", "schema.guest_minimum", "-1"), "guest min negative");
eq(reserve__request__validate(R("Ada", 3n, "dinner"), { ...S, allowed_slots: [] }), V("slot", "str.one_of", "dinner"), "slots empty");
// validate: fail-fast priority — earlier data failure wins.
eq(reserve__request__validate(R("", 3n, "dinner"), { ...S, party_lower: 8n, party_upper: 1n }), V("guest", "str.length_scalars", ""), "guest beats policy");
eq(reserve__request__validate(R("Ada", 99n, "brunch"), S), V("party", "int.closed_range", "99"), "party beats slot");
// validate: exact rendering past the host safe-integer range.
eq(reserve__request__validate(R("Ada", -9007199254740993n, "dinner"), S), V("party", "int.closed_range", "-9007199254740993"), "bigint exact");
// admit: rejected requests never confirm.
eq(reserve__request__admit(R("", 3n, "dinner"), S), V("guest", "str.length_scalars", ""), "admit rejected");
eq(reserve__request__admit(R("Ada", 3n, "dinner"), S), { $can_kind: "ok", guest: "Ada", party: 3n, slot: "dinner" }, "admit ok");
console.log("RESERVE_PARITY_OK");
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
	if !strings.Contains(string(raw), "RESERVE_PARITY_OK") {
		t.Fatalf("parity driver silent:\n%s", raw)
	}
}
