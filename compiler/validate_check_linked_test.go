package main

import (
	"os"
	"testing"
)

// S2a (a91): linked-pure vectors for the closed check ops.
// runLinkedPure executes the real bodies across quota.can and
// scalars.can with no scripts: producers, the P1 membership
// scan, and the public integer converter all run under each
// check. Both module orders, like the pilot vectors.

func checkFiles(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, f := range []string{"quota.can", "scalars.can"} {
		raw, err := os.ReadFile("../std/" + map[string]string{
			"quota.can":   "quota/quota.can",
			"scalars.can": "scalars/scalars.can",
		}[f])
		if err != nil {
			t.Fatal(err)
		}
		out[f] = string(raw)
	}
	return out
}

func TestLengthCheckLinked(t *testing.T) {
	vectors := []struct {
		name       string
		value      string
		constraint string
		expect     string
	}{
		{
			"ok",
			`"Ada"`,
			`Validate__Length(path = "guest", tag = "guest", minimum = 1, maximum = 24)`,
			`Ok("Ada")`,
		},
		{
			"short",
			`""`,
			`Validate__Length(path = "guest", tag = "guest", minimum = 1, maximum = 24)`,
			`validation.schema_violation(path = "guest", rule = "str.length_scalars", value = "")`,
		},
		{
			"min_negative",
			`"Ada"`,
			`Validate__Length(path = "guest", tag = "guest", minimum = -1, maximum = 24)`,
			`validation.schema_violation(path = "", rule = "schema.guest_minimum", value = "-1")`,
		},
		{
			"bounds",
			`"Ada"`,
			`Validate__Length(path = "guest", tag = "guest", minimum = 5, maximum = 1)`,
			`validation.schema_violation(path = "", rule = "schema.guest_bounds", value = "minimum=5;maximum=1")`,
		},
		{
			"tag_passthrough",
			`"Ada"`,
			`Validate__Length(path = "label", tag = "label", minimum = 5, maximum = 1)`,
			`validation.schema_violation(path = "", rule = "schema.label_bounds", value = "minimum=5;maximum=1")`,
		},
	}
	for _, order := range [][]string{
		{"quota.can", "scalars.can"},
		{"scalars.can", "quota.can"},
	} {
		for _, v := range vectors {
			if err := runLinkedPure(t, checkFiles(t), order,
				"std__validate__length_check", 1,
				map[string]string{"value": v.value, "constraint": v.constraint},
				v.expect); err != nil {
				t.Fatalf("linked %s (%v): %v", v.name, order, err)
			}
		}
	}
}

func TestRangeCheckLinked(t *testing.T) {
	vectors := []struct {
		name       string
		value      string
		constraint string
		expect     string
	}{
		{
			"ok",
			`3`,
			`Validate__Range(path = "party", tag = "party", lower = 1, upper = 8)`,
			`Ok(3)`,
		},
		{
			"above_renders_through_real_converter",
			`9`,
			`Validate__Range(path = "party", tag = "party", lower = 1, upper = 8)`,
			`validation.schema_violation(path = "party", rule = "int.closed_range", value = "9")`,
		},
		{
			"bounds",
			`3`,
			`Validate__Range(path = "party", tag = "party", lower = 8, upper = 1)`,
			`validation.schema_violation(path = "", rule = "schema.party_bounds", value = "lower=8;upper=1")`,
		},
		{
			"bigint_exact",
			`-9007199254740993`,
			`Validate__Range(path = "party", tag = "party", lower = 1, upper = 8)`,
			`validation.schema_violation(path = "party", rule = "int.closed_range", value = "-9007199254740993")`,
		},
	}
	for _, order := range [][]string{
		{"quota.can", "scalars.can"},
		{"scalars.can", "quota.can"},
	} {
		for _, v := range vectors {
			if err := runLinkedPure(t, checkFiles(t), order,
				"std__validate__range_check", 1,
				map[string]string{"value": v.value, "constraint": v.constraint},
				v.expect); err != nil {
				t.Fatalf("linked %s (%v): %v", v.name, order, err)
			}
		}
	}
}

func TestMembershipCheckLinked(t *testing.T) {
	vectors := []struct {
		name       string
		value      string
		constraint string
		expect     string
	}{
		{
			"ok",
			`"dinner"`,
			`Validate__Membership(path = "slot", allowed = Seq<str>["lunch", "dinner"])`,
			`Ok("dinner")`,
		},
		{
			"bad",
			`"brunch"`,
			`Validate__Membership(path = "slot", allowed = Seq<str>["lunch", "dinner"])`,
			`validation.schema_violation(path = "slot", rule = "str.one_of", value = "brunch")`,
		},
		{
			"path_passthrough",
			`"admin"`,
			`Validate__Membership(path = "mode", allowed = Seq<str>["interactive", "batch"])`,
			`validation.schema_violation(path = "mode", rule = "str.one_of", value = "admin")`,
		},
	}
	for _, order := range [][]string{
		{"quota.can", "scalars.can"},
		{"scalars.can", "quota.can"},
	} {
		for _, v := range vectors {
			if err := runLinkedPure(t, checkFiles(t), order,
				"std__validate__membership_check", 1,
				map[string]string{"value": v.value, "constraint": v.constraint},
				v.expect); err != nil {
				t.Fatalf("linked %s (%v): %v", v.name, order, err)
			}
		}
	}
}

// TestCheckLinkedContradiction pins failure detection: a wrong
// tag infix and a wrong rendered bound must both fail.
func TestCheckLinkedContradiction(t *testing.T) {
	files := checkFiles(t)
	order := []string{"quota.can", "scalars.can"}
	err := runLinkedPure(t, files, order, "std__validate__length_check", 1,
		map[string]string{
			"value":      `"Ada"`,
			"constraint": `Validate__Length(path = "guest", tag = "guest", minimum = 5, maximum = 1)`,
		},
		`validation.schema_violation(path = "", rule = "schema.label_bounds", value = "minimum=5;maximum=1")`)
	if err == nil {
		t.Fatal("expected payload mismatch for wrong tag, got pass")
	}
	err = runLinkedPure(t, files, order, "std__validate__range_check", 1,
		map[string]string{
			"value":      `3`,
			"constraint": `Validate__Range(path = "party", tag = "party", lower = 8, upper = 1)`,
		},
		`validation.schema_violation(path = "", rule = "schema.party_bounds", value = "lower=8;upper=2")`)
	if err == nil {
		t.Fatal("expected payload mismatch for wrong bound render, got pass")
	}
}
