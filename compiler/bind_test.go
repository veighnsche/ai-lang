package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// One binding rule: positional args bind by source index, named args by
// parameter name, and every parameter receives exactly one argument.

func TestBindSlots(t *testing.T) {
	params := [][2]string{{"left", "int"}, {"right", "int"}}
	intLit := func() *Small { return &Small{Kind: "int"} }

	cases := []struct {
		name string
		args []Arg
		want []int
		fail string
	}{
		{"positional", []Arg{{V: intLit()}, {V: intLit()}}, []int{0, 1}, ""},
		{"named in order", []Arg{{Name: "left", HasName: true, V: intLit()}, {Name: "right", HasName: true, V: intLit()}}, []int{0, 1}, ""},
		{"named reordered", []Arg{{Name: "right", HasName: true, V: intLit()}, {Name: "left", HasName: true, V: intLit()}}, []int{1, 0}, ""},
		{"mixed consistent", []Arg{{V: intLit()}, {Name: "right", HasName: true, V: intLit()}}, []int{0, 1}, ""},
		{"unknown name", []Arg{{Name: "up", HasName: true, V: intLit()}, {Name: "left", HasName: true, V: intLit()}}, nil, "has no param up"},
		{"duplicate named", []Arg{{Name: "left", HasName: true, V: intLit()}, {Name: "left", HasName: true, V: intLit()}}, nil, "supplies arg left twice"},
		{"mixed double assignment", []Arg{{V: intLit()}, {Name: "left", HasName: true, V: intLit()}}, nil, "supplies arg left twice"},
		{"positional overflow", []Arg{{V: intLit()}, {V: intLit()}, {V: intLit()}}, nil, "takes 3 args for 2 params"},
		{"missing arg", []Arg{{Name: "left", HasName: true, V: intLit()}}, nil, "is missing arg right"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := bindSlots("f", tc.args, params)
			if tc.fail == "" {
				if err != nil {
					t.Fatalf("bindSlots: %v", err)
				}
				if len(got) != len(tc.want) {
					t.Fatalf("slots = %v, want %v", got, tc.want)
				}
				for i := range got {
					if got[i] != tc.want[i] {
						t.Fatalf("slots = %v, want %v", got, tc.want)
					}
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.fail) {
				t.Fatalf("bindSlots err = %v, want %q", err, tc.fail)
			}
		})
	}
}

// A reordered same-type named call is the case TypeScript cannot catch:
// both orders type-check, only the binding vector decides 7 vs -7.
const bindReorder = `mod audit
  provides [audit__go, audit__subtract, Audit__Value]
  uses []
  emits []

type Audit__Value rev 1 (
  value: int
)

fn audit__subtract(left: int, right: int) -> Audit__Value rev 1
  emits []
  tests
    s(9, 2) => Ok(value = 7)
  Ok(value = left - right)

fn audit__go() -> Audit__Value rev 1
  emits []
  tests
    g() => Ok(value = 7)
  match call audit__subtract(right = 2, left = 9)
    on Ok r => Ok(value = r.value)
`

func TestBindReorderedCallEvaluates(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"audit.can": bindReorder})
	// diagnose runs every decision table: g() => 7 passes only if the
	// evaluator binds right = 2, left = 9 (-7 would fail the table).
	if diags := diagnose(dir, "audit.can", bindReorder); len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}

func TestBindReorderedCallEmitsInParamOrder(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "audit.can"), []byte(bindReorder), 0o644); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	if err := compile(out, []string{filepath.Join(dir, "audit.can")}); err != nil {
		t.Fatalf("compile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(out, "audit.ts"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "audit__subtract(9n, 2n)") {
		t.Errorf("emit did not lower the reordered call to parameter order:\n%s", got)
	}
	if strings.Contains(string(got), "audit__subtract(2n, 9n)") {
		t.Errorf("emit kept source order for the reordered call:\n%s", got)
	}
}

func TestBindBadVectorIsCAN3010(t *testing.T) {
	dup := strings.Replace(bindReorder,
		"match call audit__subtract(right = 2, left = 9)",
		"match call audit__subtract(left = 9, left = 9)", 1)
	dir := writeLSPDir(t, map[string]string{"audit.can": dup})
	diags := diagnose(dir, "audit.can", dup)
	if !hasDiag(diags, "error", "supplies arg left twice") {
		t.Fatalf("expected duplicate-binding error, got %v", diags)
	}
	if !hasCode(diags, "CAN3010") {
		t.Fatalf("expected CAN3010, got %v", diags)
	}

	missing := strings.Replace(bindReorder,
		"match call audit__subtract(right = 2, left = 9)",
		"match call audit__subtract(right = 2)", 1)
	dir = writeLSPDir(t, map[string]string{"audit.can": missing})
	diags = diagnose(dir, "audit.can", missing)
	if !hasDiag(diags, "error", "is missing arg left") {
		t.Fatalf("expected missing-binding error, got %v", diags)
	}
	if !hasCode(diags, "CAN3010") {
		t.Fatalf("expected CAN3010, got %v", diags)
	}
}

func hasCode(diags []Diag, code string) bool {
	for _, d := range diags {
		if d.Code == code {
			return true
		}
	}
	return false
}
