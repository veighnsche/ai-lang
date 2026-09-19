package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Blessed-stdlib gates: compiling a std module executes every
// decision table, so a green compile pins the module's behavior
// in the suite. One test per migrated-or-new module below.

// TestStdMapCompiles pins the generic Map<K,V> migration: both
// instantiations check, execute, and emit.
func TestStdMapCompiles(t *testing.T) {
	out := t.TempDir()
	if err := compile(out, []string{"../std/map/map.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
}

// TestStdHostCompiles pins the host shelf: decision tables script
// every extern, so a green compile proves dispatch, not the host.
func TestStdHostCompiles(t *testing.T) {
	out := t.TempDir()
	if err := compile(out, []string{"../std/host/host.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
}

// TestStdHostNodeSmoke executes the committed host implementations
// under node: wall returns positive bigint millis, monotonic never
// goes backwards. Node is required, never skipped (parity
// precedent): an unexecuted host obligation is a gap, not a pass.
func TestStdHostNodeSmoke(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("node missing: host smoke refuses to skip: %v", err)
	}
	dir := t.TempDir()
	raw, err := os.ReadFile("../std/host/host.externs.ts")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "host.externs.ts"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
	harness := `import { host__wall_now, host__mono_now } from "./host.externs.ts";
const w = host__wall_now();
if (w.$can_kind !== "ok" || typeof w.millis !== "bigint" || w.millis <= 0n) {
  console.error("wall clock bad: " + JSON.stringify(w, (_, v) => typeof v === "bigint" ? v.toString() : v));
  process.exit(1);
}
const a = host__mono_now();
await new Promise((r) => setTimeout(r, 2));
const b = host__mono_now();
for (const t of [a, b]) {
  if (t.$can_kind !== "ok" || typeof t.millis !== "bigint" || t.millis < 0n) {
    console.error("mono clock bad");
    process.exit(1);
  }
}
if (b.millis < a.millis) {
  console.error("mono clock went backwards");
  process.exit(1);
}
console.log("HOST_SMOKE_OK");
`
	if err := os.WriteFile(filepath.Join(dir, "smoke.mjs"), []byte(harness), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("node", filepath.Join(dir, "smoke.mjs"))
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("node smoke: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "HOST_SMOKE_OK") {
		t.Fatalf("node smoke missing receipt:\n%s", out)
	}
}
