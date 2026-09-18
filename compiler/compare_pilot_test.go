package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// G1 pilot (a93): the generic std__compare stamps one checked
// copy per ordered type (bool stays monomorphic: >= refuses
// bool operands, so no single body serves all four). Linked
// vectors execute the real stamped bodies; the parity driver
// executes the stamped TS emit under node. Stamp names ($)
// never appear in .can source, only here as test addresses,
// the same precedent as the expansion tests.
func comparePilotFiles(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../std/scalars/scalars.can")
	if err != nil {
		t.Fatal(err)
	}
	return map[string]string{"scalars.can": string(raw)}
}

func TestComparePilotLinked(t *testing.T) {
	vectors := []struct {
		name   string
		root   string
		args   map[string]string
		expect string
	}{
		{"int gt", "std__compare$T$int", map[string]string{"left": "5", "right": "4"}, "Ok(1)"},
		{"int eq", "std__compare$T$int", map[string]string{"left": "4", "right": "4"}, "Ok(0)"},
		{"int lt", "std__compare$T$int", map[string]string{"left": "3", "right": "4"}, "Ok(-1)"},
		{"dec gt", "std__compare$T$dec", map[string]string{"left": `d"1.6"`, "right": `d"1.5"`}, "Ok(1)"},
		{"dec eq", "std__compare$T$dec", map[string]string{"left": `d"1.5"`, "right": `d"1.5"`}, "Ok(0)"},
		{"dec lt", "std__compare$T$dec", map[string]string{"left": `d"1.4"`, "right": `d"1.5"`}, "Ok(-1)"},
		{"str gt", "std__compare$T$str", map[string]string{"left": `"b"`, "right": `"a"`}, "Ok(1)"},
		{"str eq", "std__compare$T$str", map[string]string{"left": `"a"`, "right": `"a"`}, "Ok(0)"},
		{"str lt", "std__compare$T$str", map[string]string{"left": `"a"`, "right": `"b"`}, "Ok(-1)"},
		{"bool tf", "std__compare__bool", map[string]string{"left": "true", "right": "false"}, "Ok(1)"},
		{"bool ft", "std__compare__bool", map[string]string{"left": "false", "right": "true"}, "Ok(-1)"},
		{"bool tt", "std__compare__bool", map[string]string{"left": "true", "right": "true"}, "Ok(0)"},
	}
	for _, v := range vectors {
		if err := runLinkedPure(t, comparePilotFiles(t), []string{"scalars.can"}, v.root, 1, v.args, v.expect); err != nil {
			t.Fatalf("linked %s: %v", v.name, err)
		}
	}
}

func TestComparePilotParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("node missing: runtime parity refuses to skip: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{"../std/scalars/scalars.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	driver := `import { std__compare$T$int, std__compare$T$dec, std__compare$T$str, std__compare__bool } from "./scalars.ts";
const show = (x) => JSON.stringify(x, (k, v) => typeof v === "bigint" ? { $bigint: v.toString() } : v);
const eq = (got, want, name) => {
  if (show(got) !== show(want)) throw new Error(name + ": got " + show(got));
};
const O = (v) => ({ $can_kind: "ok", value: v });
eq(std__compare$T$int(5n, 4n), O(1n), "int gt");
eq(std__compare$T$int(4n, 4n), O(0n), "int eq");
eq(std__compare$T$int(3n, 4n), O(-1n), "int lt");
eq(std__compare$T$dec("1.6", "1.5"), O(1n), "dec gt");
eq(std__compare$T$dec("1.5", "1.5"), O(0n), "dec eq");
eq(std__compare$T$dec("1.4", "1.5"), O(-1n), "dec lt");
eq(std__compare$T$str("b", "a"), O(1n), "str gt");
eq(std__compare$T$str("a", "a"), O(0n), "str eq");
eq(std__compare$T$str("a", "b"), O(-1n), "str lt");
eq(std__compare__bool(true, false), O(1n), "bool tf");
eq(std__compare__bool(false, true), O(-1n), "bool ft");
eq(std__compare__bool(true, true), O(0n), "bool tt");
eq(std__compare__bool(false, false), O(0n), "bool ff");
console.log("COMPARE_PARITY_OK");
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
	if !strings.Contains(string(raw), "COMPARE_PARITY_OK") {
		t.Fatalf("parity driver silent:\n%s", raw)
	}
}
