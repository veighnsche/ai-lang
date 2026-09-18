package main

import (
	"os"
	"testing"
)

// S2 second consumer: linked-pure vectors for
// reserve__request__validate and reserve__request__admit.
// runLinkedPure executes the real bodies across reserve.can,
// quota.can, and scalars.can with no scripts: the shared scalar
// validators, the P1 membership scan, and the public integer
// converter all run. Admit is pure here, so unlike the quota
// pilot it is covered too — rejection-before-confirmation
// executes end to end.

func reserveFiles(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, f := range []string{"reserve.can", "quota.can", "scalars.can"} {
		raw, err := os.ReadFile(map[string]string{
			"reserve.can": "../sketches/reserve/reserve.can",
			"quota.can":   "../std/quota/quota.can",
			"scalars.can": "../std/scalars/scalars.can",
		}[f])
		if err != nil {
			t.Fatal(err)
		}
		out[f] = string(raw)
	}
	return out
}

func TestReserveLinked(t *testing.T) {
	schema := `Reserve__RequestSchema(guest_minimum = 1, guest_maximum = 24, party_lower = 1, party_upper = 8, allowed_slots = Seq<str>["lunch", "dinner"])`
	vectors := []struct {
		name    string
		request string
		schema  string
		expect  string
	}{
		{
			"ok",
			`Reserve__Request(guest = "Ada", party = 3, slot = "dinner")`,
			schema,
			`Ok(Reserve__Request(guest = "Ada", party = 3, slot = "dinner"))`,
		},
		{
			"guest",
			`Reserve__Request(guest = "", party = 3, slot = "dinner")`,
			schema,
			`validation.schema_violation(path = "guest", rule = "str.length_scalars", value = "")`,
		},
		{
			"party_renders_through_real_converter",
			`Reserve__Request(guest = "Ada", party = 9, slot = "dinner")`,
			schema,
			`validation.schema_violation(path = "party", rule = "int.closed_range", value = "9")`,
		},
		{
			"slot",
			`Reserve__Request(guest = "Ada", party = 3, slot = "brunch")`,
			schema,
			`validation.schema_violation(path = "slot", rule = "str.one_of", value = "brunch")`,
		},
		{
			"policy_bounds",
			`Reserve__Request(guest = "Ada", party = 3, slot = "dinner")`,
			`Reserve__RequestSchema(guest_minimum = 5, guest_maximum = 1, party_lower = 1, party_upper = 8, allowed_slots = Seq<str>["lunch", "dinner"])`,
			`validation.schema_violation(path = "", rule = "schema.guest_bounds", value = "minimum=5;maximum=1")`,
		},
		{
			"bigint_exact",
			`Reserve__Request(guest = "Ada", party = -9007199254740993, slot = "dinner")`,
			schema,
			`validation.schema_violation(path = "party", rule = "int.closed_range", value = "-9007199254740993")`,
		},
	}
	for _, order := range [][]string{
		{"reserve.can", "quota.can", "scalars.can"},
		{"scalars.can", "quota.can", "reserve.can"},
	} {
		for _, v := range vectors {
			if err := runLinkedPure(t, reserveFiles(t), order,
				"reserve__request__validate", 1,
				map[string]string{"request": v.request, "schema": v.schema},
				v.expect); err != nil {
				t.Fatalf("linked %s (%v): %v", v.name, order, err)
			}
		}
	}
}

func TestReserveAdmitLinked(t *testing.T) {
	schema := `Reserve__RequestSchema(guest_minimum = 1, guest_maximum = 24, party_lower = 1, party_upper = 8, allowed_slots = Seq<str>["lunch", "dinner"])`
	vectors := []struct {
		name    string
		request string
		schema  string
		expect  string
	}{
		{
			"admitted",
			`Reserve__Request(guest = "Ada", party = 3, slot = "dinner")`,
			schema,
			`Ok(guest = "Ada", party = 3, slot = "dinner")`,
		},
		{
			"rejected_before_confirm",
			`Reserve__Request(guest = "", party = 3, slot = "dinner")`,
			schema,
			`validation.schema_violation(path = "guest", rule = "str.length_scalars", value = "")`,
		},
	}
	for _, order := range [][]string{
		{"reserve.can", "quota.can", "scalars.can"},
		{"scalars.can", "quota.can", "reserve.can"},
	} {
		for _, v := range vectors {
			if err := runLinkedPure(t, reserveFiles(t), order,
				"reserve__request__admit", 1,
				map[string]string{"request": v.request, "schema": v.schema},
				v.expect); err != nil {
				t.Fatalf("linked %s (%v): %v", v.name, order, err)
			}
		}
	}
}

// TestReserveLinkedContradiction pins failure detection: a wrong
// rendered value and a wrong rule token must both fail.
func TestReserveLinkedContradiction(t *testing.T) {
	schema := `Reserve__RequestSchema(guest_minimum = 1, guest_maximum = 24, party_lower = 1, party_upper = 8, allowed_slots = Seq<str>["lunch", "dinner"])`
	request := `Reserve__Request(guest = "Ada", party = 9, slot = "dinner")`
	files := reserveFiles(t)
	order := []string{"reserve.can", "quota.can", "scalars.can"}
	err := runLinkedPure(t, files, order, "reserve__request__validate", 1,
		map[string]string{"request": request, "schema": schema},
		`validation.schema_violation(path = "party", rule = "int.closed_range", value = "10")`)
	if err == nil {
		t.Fatal("expected payload mismatch for wrong render, got pass")
	}
	err = runLinkedPure(t, files, order, "reserve__request__validate", 1,
		map[string]string{"request": request, "schema": schema},
		`validation.schema_violation(path = "party", rule = "str.one_of", value = "9")`)
	if err == nil {
		t.Fatal("expected payload mismatch for wrong rule, got pass")
	}
}
