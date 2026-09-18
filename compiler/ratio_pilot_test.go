package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// G2 pilot (a97): Ratio__Int/Bool/Dec unify into Ratio__Value<T>;
// one shape stamps int, bool, and dec, and rows stay positional
// Ok(...) with .value projections. Linked vectors execute through
// the stamped records; the parity driver executes the TS emit
// under node. Stamp names ($) never appear in .can source, only
// in the emit the golden pins.
func ratioPilotFiles(t *testing.T) map[string]string {
	t.Helper()
	raw, err := os.ReadFile("../std/ratio/ratio.can")
	if err != nil {
		t.Fatal(err)
	}
	return map[string]string{"ratio.can": string(raw)}
}

func TestRatioPilotLinked(t *testing.T) {
	vectors := []struct {
		name   string
		root   string
		args   map[string]string
		expect string
	}{
		{"abs neg", "std__ratio__abs", map[string]string{"value": "-3"}, "Ok(3)"},
		{"abs zero", "std__ratio__abs", map[string]string{"value": "0"}, "Ok(0)"},
		{"gcd basic", "std__ratio__gcd_from", map[string]string{"a": "12", "b": "8"}, "Ok(4)"},
		{"gcd coprime", "std__ratio__gcd_from", map[string]string{"a": "8", "b": "9"}, "Ok(1)"},
		{"seen hit", "std__ratio__seen_from", map[string]string{"seen": "Seq<int>[2, 1]", "target": "1", "pos": "0", "fuel": "3"}, "Ok(true)"},
		{"seen miss", "std__ratio__seen_from", map[string]string{"seen": "Seq<int>[2, 3]", "target": "1", "pos": "0", "fuel": "3"}, "Ok(false)"},
		{"digit five", "std__ratio__digit_from", map[string]string{"remaining": "5", "acc": `d"0.0"`, "fuel": "6"}, `Ok(d"5.0")`},
		{"fold eighth", "std__ratio__fold_from", map[string]string{"digits": "Seq<int>[1, 2, 5]", "pos": "0", "place": `d"0.1"`, "acc": `d"0.0"`, "fuel": "4"}, `Ok(d"0.125")`},
		{"exact half", "std__ratio__to_dec_exact", map[string]string{"value": "Ratio__Fraction(1, 2)"}, `Ok(d"0.5")`},
		{"exact whole", "std__ratio__to_dec_exact", map[string]string{"value": "Ratio__Fraction(7, 2)"}, `Ok(d"3.5")`},
	}
	for _, v := range vectors {
		if err := runLinkedPure(t, ratioPilotFiles(t), []string{"ratio.can"}, v.root, 1, v.args, v.expect); err != nil {
			t.Fatalf("linked %s: %v", v.name, err)
		}
	}
}

func TestRatioPilotParity(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Fatalf("node missing: runtime parity refuses to skip: %v", err)
	}
	out := t.TempDir()
	if err := compile(out, []string{"../std/ratio/ratio.can"}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	driver := `import { std__ratio__abs, std__ratio__gcd_from, std__ratio__seen_from, std__ratio__digit_from, std__ratio__fold_from, std__ratio__to_dec_exact } from "./ratio.ts";
const show = (x) => JSON.stringify(x, (k, v) => typeof v === "bigint" ? { $bigint: v.toString() } : v);
const eq = (got, want, name) => {
  if (show(got) !== show(want)) throw new Error(name + ": got " + show(got));
};
const O = (v) => ({ $can_kind: "ok", value: v });
eq(std__ratio__abs(-3n), O(3n), "abs neg");
eq(std__ratio__abs(0n), O(0n), "abs zero");
eq(std__ratio__gcd_from(12n, 8n), O(4n), "gcd basic");
eq(std__ratio__gcd_from(8n, 9n), O(1n), "gcd coprime");
eq(std__ratio__seen_from([2n, 1n], 1n, 0n, 3n), O(true), "seen hit");
eq(std__ratio__seen_from([2n, 3n], 1n, 0n, 3n), O(false), "seen miss");
eq(std__ratio__digit_from(5n, "0.0", 6n), O("5.0"), "digit five");
eq(std__ratio__fold_from([1n, 2n, 5n], 0n, "0.1", "0.0", 4n), O("0.125"), "fold eighth");
eq(std__ratio__to_dec_exact({ numerator: 1n, denominator: 2n }), O("0.5"), "exact half");
eq(std__ratio__to_dec_exact({ numerator: 7n, denominator: 2n }), O("3.5"), "exact whole");
console.log("RATIO_PARITY_OK");
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
	if !strings.Contains(string(raw), "RATIO_PARITY_OK") {
		t.Fatalf("parity driver silent:\n%s", raw)
	}
}
