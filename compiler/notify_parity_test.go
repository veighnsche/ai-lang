package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// S2b third consumer: generated-runtime agreement for the
// notification desk. Compiles notify + quota + scalars and
// executes validate/admit under node: success preservation,
// every violation shape, the frozen payload mapping, fail-fast
// priority, exact beyond-safe-integer rendering, and the
// admission gate (rejected requests never ticket).
func TestNotifyRuntimeParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("node missing: runtime parity refuses to skip: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{"../sketches/notify/notify.can", "../std/quota/quota.can", "../std/scalars/scalars.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	// Same harness shim as the S1a parity test (a89 amendment 1).
	rewriteSpecifier(t, filepath.Join(out, "notify.ts"), "./quota", "./quota.ts")
	// No ./scalars rewrite: notify never calls the converter
	// directly, so its emit imports quota alone.
	rewriteSpecifier(t, filepath.Join(out, "quota.ts"), "./scalars", "./scalars.ts")
	driver := `import { notify__request__validate, notify__request__admit } from "./notify.ts";
const show = (x) => JSON.stringify(x, (k, v) => typeof v === "bigint" ? { $bigint: v.toString() } : v);
const eq = (got, want, name) => {
  if (show(got) !== show(want)) throw new Error(name + ": got " + show(got));
};
const V = (path, rule, value) => ({ $can_kind: "validation.schema_violation", path, rule, value });
const S = { handle_minimum: 3n, handle_maximum: 16n, retries_lower: 0n, retries_upper: 5n, allowed_channels: ["email", "sms"] };
const R = (handle, retries, channel) => ({ handle, retries, channel });
// validate: success preserves every field.
eq(notify__request__validate(R("sam", 2n, "email"), S), { $can_kind: "ok", value: R("sam", 2n, "email") }, "validate ok");
eq(notify__request__validate(R(" ab ", 2n, "email"), S), { $can_kind: "ok", value: R(" ab ", 2n, "email") }, "whitespace preserved");
// validate: each violation shape and the frozen mapping.
eq(notify__request__validate(R("ab", 2n, "email"), S), V("handle", "str.length_scalars", "ab"), "handle short");
eq(notify__request__validate(R("sam", 6n, "email"), S), V("retries", "int.closed_range", "6"), "retries above");
eq(notify__request__validate(R("sam", 2n, "push"), S), V("channel", "str.one_of", "push"), "channel bad");
eq(notify__request__validate(R("sam", 2n, "email"), { ...S, handle_minimum: 9n, handle_maximum: 4n }), V("", "schema.handle_bounds", "minimum=9;maximum=4"), "handle bounds");
eq(notify__request__validate(R("sam", 2n, "email"), { ...S, handle_minimum: -1n }), V("", "schema.handle_minimum", "-1"), "handle min negative");
eq(notify__request__validate(R("sam", 2n, "email"), { ...S, allowed_channels: [] }), V("channel", "str.one_of", "email"), "channels empty");
// validate: fail-fast priority — earlier data failure wins.
eq(notify__request__validate(R("ab", 2n, "email"), { ...S, retries_lower: 4n, retries_upper: 1n }), V("handle", "str.length_scalars", "ab"), "handle beats policy");
eq(notify__request__validate(R("sam", 99n, "push"), S), V("retries", "int.closed_range", "99"), "retries beats channel");
// validate: exact rendering past the host safe-integer range.
eq(notify__request__validate(R("sam", -9007199254740993n, "email"), S), V("retries", "int.closed_range", "-9007199254740993"), "bigint exact");
// admit: rejected requests never ticket.
eq(notify__request__admit(R("ab", 2n, "email"), S), V("handle", "str.length_scalars", "ab"), "admit rejected");
eq(notify__request__admit(R("sam", 2n, "email"), S), { $can_kind: "ok", handle: "sam", retries: 2n, channel: "email" }, "admit ok");
console.log("NOTIFY_PARITY_OK");
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
	if !strings.Contains(string(raw), "NOTIFY_PARITY_OK") {
		t.Fatalf("parity driver silent:\n%s", raw)
	}
}
