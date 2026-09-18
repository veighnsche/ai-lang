package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// P1 (a89): generated-runtime agreement for
// std__validate__str_one_of. Compiles the real blessed
// std/quota module and executes membership under node:
// hits, misses, empty policy, duplicates, case
// sensitivity, and unicode must agree with the
// interpreter's decision tables, including the complete
// validation.not_allowed payload.
func TestValidateStrOneOfRuntimeParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("node missing: runtime parity refuses to skip: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{"../std/quota/quota.can", "../std/scalars/scalars.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	// Harness shim for the S1a packaging gap (a89 amendment 1):
	// emit writes extensionless "./scalars", which node ESM cannot
	// resolve. Rewrite the specifier in the temp copy only; the
	// golden pins the emitted bytes verbatim.
	rewriteSpecifier(t, filepath.Join(out, "quota.ts"), "./scalars", "./scalars.ts")
	driver := `import { std__validate__str_one_of } from "./quota.ts";
const eq = (got, want, name) => {
  if (JSON.stringify(got) !== JSON.stringify(want)) throw new Error(name + ": got " + JSON.stringify(got));
};
const ok = (v) => ({ $can_kind: "ok", value: v });
const deny = (v) => ({ $can_kind: "validation.not_allowed", value: v });
eq(std__validate__str_one_of("a", []), deny("a"), "empty policy denies");
eq(std__validate__str_one_of("a", ["a"]), ok("a"), "singleton hit");
eq(std__validate__str_one_of("b", ["a"]), deny("b"), "singleton miss");
eq(std__validate__str_one_of("a", ["a", "b", "c"]), ok("a"), "first");
eq(std__validate__str_one_of("b", ["a", "b", "c"]), ok("b"), "middle");
eq(std__validate__str_one_of("c", ["a", "b", "c"]), ok("c"), "last");
eq(std__validate__str_one_of("z", ["a", "b", "c"]), deny("z"), "absent");
eq(std__validate__str_one_of("a", ["a", "a"]), ok("a"), "dup hit");
eq(std__validate__str_one_of("Batch", ["batch"]), deny("Batch"), "case sensitive");
eq(std__validate__str_one_of("世界", ["héllo", "世界"]), ok("世界"), "unicode hit");
eq(std__validate__str_one_of("", ["", "a"]), ok(""), "empty value hit");
eq(std__validate__str_one_of("", ["a"]), deny(""), "empty value miss");
console.log("ONEOF_PARITY_OK");
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
	if !strings.Contains(string(raw), "ONEOF_PARITY_OK") {
		t.Fatalf("parity driver silent:\n%s", raw)
	}
}

// rewriteSpecifier swaps one module specifier in a temp-dir emit
// so node ESM can resolve it. Test-harness only: committed
// goldens keep the emitted bytes, and the cross-dir import story
// (paths and extensions) belongs to a future emit slice.
func rewriteSpecifier(t *testing.T, path, old, new string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"`+old+`"`) {
		t.Fatalf("specifier %q not found in %s", old, path)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(string(raw), `"`+old+`"`, `"`+new+`"`, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
}
