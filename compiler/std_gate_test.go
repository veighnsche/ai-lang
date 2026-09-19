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

// TestStdSeqCompiles pins the higher-order seq extension: generic
// consumers invoke total callbacks through map/filter/fold/all/
// any/find, so a green compile proves generic invoke end to end.
func TestStdSeqCompiles(t *testing.T) {
	out := t.TempDir()
	if err := compile(out, []string{"../std/seq/seq.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
}

// TestStdHostNodeSmoke executes the committed host implementations
// under node: wall returns positive bigint millis, monotonic never
// goes backwards, random returns exact-length bytes and rejects
// negative counts, sha256 matches the abc known vector and md5 is
// refused, secrets compare equal/unequal/length-mismatch, env reads
// a set var and reports a missing one absent, and log emits one
// exact JSON line on stderr. Node is required, never skipped
// (parity precedent): an unexecuted host obligation is a gap, not
// a pass.
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
	harness := `import { host__wall_now, host__mono_now, host__rand_bytes, host__hash_digest, host__secret_equal, host__env_read, host__log_write } from "./host.externs.ts";
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
const r = host__rand_bytes(16n);
if (r.$can_kind !== "ok" || !(r.bytes instanceof Uint8Array) || r.bytes.length !== 16) {
  console.error("random bytes bad");
  process.exit(1);
}
const neg = host__rand_bytes(-1n);
if (neg.$can_kind !== "random.invalid_count") {
  console.error("random negative count not rejected");
  process.exit(1);
}
const abc = new Uint8Array([97, 98, 99]);
const h = host__hash_digest(abc, "sha256");
const hex = [...h.bytes].map((x) => x.toString(16).padStart(2, "0")).join("");
if (h.$can_kind !== "ok" || hex !== "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad") {
  console.error("sha256 abc vector mismatch: " + hex);
  process.exit(1);
}
const md5 = host__hash_digest(abc, "md5");
if (md5.$can_kind !== "hash.unsupported_profile") {
  console.error("md5 not rejected");
  process.exit(1);
}
const eq = host__secret_equal("s3cr3t", "s3cr3t");
const ne = host__secret_equal("s3cr3t", "other!");
const len = host__secret_equal("short", "much longer value");
if (eq.$can_kind !== "ok" || eq.equal !== true) {
  console.error("secret same not equal");
  process.exit(1);
}
if (ne.$can_kind !== "ok" || ne.equal !== false) {
  console.error("secret diff not unequal");
  process.exit(1);
}
if (len.$can_kind !== "ok" || len.equal !== false) {
  console.error("secret length mismatch not unequal");
  process.exit(1);
}
process.env.NODE_SMOKE_HOST_TEST = "smoke-ok";
const found = host__env_read("NODE_SMOKE_HOST_TEST");
const missing = host__env_read("NODE_SMOKE_HOST_ABSENT_XYZ");
if (found.$can_kind !== "ok" || found.value !== "smoke-ok") {
  console.error("env read bad");
  process.exit(1);
}
if (missing.$can_kind !== "environment.absent") {
  console.error("env absent bad");
  process.exit(1);
}
const lines = [];
const origError = console.error;
console.error = (s) => { lines.push(String(s)); };
const logged = host__log_write({ level: 1n, message: "hello-smoke" });
console.error = origError;
if (logged.$can_kind !== "ok") {
  console.error("log write bad");
  process.exit(1);
}
if (lines.length !== 1 || lines[0] !== '{"level":"1","message":"hello-smoke"}') {
  console.error("log line bad: " + JSON.stringify(lines));
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
