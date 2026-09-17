package main

import (
	"testing"
)

// Arithmetic (v0.6, unbounded v10) unit proofs: exactness and
// normalization at the evaluator level. Decision-table coverage lives
// in lsp_test.go; this file pins the value semantics underneath it.
func wantInt(t *testing.T, expr, want string) {
	t.Helper()
	if got := evArithExpr(t, expr); got.Kind != "int" || got.N.String() != want {
		t.Fatalf("%s = %v, want int %s", expr, got, want)
	}
}
func evArithExpr(t *testing.T, expr string) *Value {
	t.Helper()
	sm, err := parseSmall(expr)
	if err != nil {
		t.Fatalf("parseSmall(%q): %v", expr, err)
	}
	v, err := evSmall(sm, map[string]*Value{}, &Ctx{}, "test")
	if err != nil {
		t.Fatalf("evSmall(%q): %v", expr, err)
	}
	return v
}

func TestArithDecExact(t *testing.T) {
	cases := []struct{ expr, want string }{
		{`d"0.1" + d"0.2"`, "0.3"},
		{`d"1.50" + d"0.5"`, "2.0"},
		{`d"2.0" - d"0.75"`, "1.25"},
		{`d"1.5" * d"2.0"`, "3.0"},
		{`d"0.1" * d"0.1"`, "0.01"},
		{`d"1.0" - d"1.0"`, "0.0"},
		{`d"0.5" - d"1.0"`, "-0.5"},
	}
	for _, c := range cases {
		if got := evArithExpr(t, c.expr); got.Kind != "dec" || got.D != c.want {
			t.Fatalf("%s = %v, want dec %s", c.expr, got, c.want)
		}
	}
}

func TestArithInt(t *testing.T) {
	wantInt(t, "2 + 3 * 4", "14")
	wantInt(t, "10 - 3 - 2", "5")
}

func wantBool(t *testing.T, expr string, want bool) {
	t.Helper()
	got := evArithExpr(t, expr)
	if got.Kind != "bool" || got.B != want {
		t.Fatalf("%s = %v, want bool %v", expr, got, want)
	}
}

// TestStrictComparison pins < and > (issue #40) at the evaluator
// level: ints by big-int order, decs by exact rational order.
// Strings compare by byte order like >= and <=; str literals cannot
// sit beside an operator (quotes swallow it, as ever), so string
// strictness is pinned via refs in TestStrictComparisonStr below.
func TestStrictComparison(t *testing.T) {
	for _, c := range []struct {
		expr string
		want bool
	}{
		{`3 > 2`, true},
		{`2 > 3`, false},
		{`2 > 2`, false},
		{`2 < 3`, true},
		{`3 < 2`, false},
		{`2 < 2`, false},
		{`-1 < 0`, true},
		{`d"1.5" > d"1.0"`, true},
		{`d"1.0" > d"1.5"`, false},
		{`d"1.5" > d"1.5"`, false},
		{`d"1.0" < d"1.5"`, true},
		{`d"1.5" < d"1.0"`, false},
		{`d"-0.5" < d"0.0"`, true},
		{`d"0.10" < d"0.2"`, true},
	} {
		wantBool(t, c.expr, c.want)
	}
}

func TestStrictComparisonStr(t *testing.T) {
	env := map[string]*Value{
		"x": {Kind: "str", S: "a"},
		"y": {Kind: "str", S: "b"},
	}
	for _, c := range []struct {
		expr string
		want bool
	}{
		{`x < y`, true},
		{`y > x`, true},
		{`x > x`, false},
		{`x < x`, false},
	} {
		sm, err := parseSmall(c.expr)
		if err != nil {
			t.Fatalf("parseSmall(%q): %v", c.expr, err)
		}
		got, err := evSmall(sm, env, &Ctx{}, "test")
		if err != nil {
			t.Fatalf("evSmall(%q): %v", c.expr, err)
		}
		if got.Kind != "bool" || got.B != c.want {
			t.Fatalf("%s = %v, want bool %v", c.expr, got, c.want)
		}
	}
}

// TestArithUnbounded pins v10: ints never overflow, never wrap. The
// old int64 gate is gone; chains that once failed loud now compute
// exactly, matching the bigint target.
func TestArithUnbounded(t *testing.T) {
	wantInt(t, "3037000500 * 3037000500", "9223372037000250000")
	wantInt(t, "1000000000000000000 - 1", "999999999999999999")
	wantInt(t, "123456789012345678901234567890 + 1", "123456789012345678901234567891")
	wantInt(t, "99999999999999999999999999 * 99999999999999999999999999",
		"9999999999999999999999999800000000000000000000000001")
}

func TestArithBadOperands(t *testing.T) {
	// Note: `"a" + "b"` is one string literal (quotes swallow
	// operators, as ever); str operands arrive via seals and refs.
	// v16: str + str concatenates (seals erase to strings at
	// runtime), so the loud string case is now -.
	for _, expr := range []string{`1 + d"1.5"`, `seal B("x") - seal B("y")`, `true * 2`} {
		sm, err := parseSmall(expr)
		if err != nil {
			t.Fatalf("parseSmall(%q): %v", expr, err)
		}
		if _, err := evSmall(sm, map[string]*Value{}, &Ctx{}, "test"); err == nil {
			t.Fatalf("evSmall(%q): expected error, got none", expr)
		}
	}
}
