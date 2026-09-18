package main

import (
	"os"
	"testing"
)

// S1a (a89): linked-pure vectors for quota__request__validate.
// runLinkedPure executes the real bodies across quota.can and
// scalars.can with no scripts: the scalar validators, the P1
// membership scan, and the public integer converter all run.
// Effectful quota__request__admit is excluded by design (a68
// refuses effectful graphs); its table coverage stands alone.

func quotaRequestFiles(t *testing.T) map[string]string {
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

func TestQuotaRequestLinked(t *testing.T) {
	schema := `Quota__RequestSchema(label_minimum = 1, label_maximum = 12, amount_lower = 1, amount_upper = 10, allowed_modes = Seq<str>["interactive", "batch"])`
	vectors := []struct {
		name    string
		request string
		schema  string
		expect  string
	}{
		{
			"ok",
			`Quota__Request(label = "Ada", amount = 5, mode = "batch")`,
			schema,
			`Ok(value = Quota__Request(label = "Ada", amount = 5, mode = "batch"))`,
		},
		{
			"label",
			`Quota__Request(label = "", amount = 5, mode = "batch")`,
			schema,
			`validation.schema_violation(path = "label", rule = "str.length_scalars", value = "")`,
		},
		{
			"amount_renders_through_real_converter",
			`Quota__Request(label = "Ada", amount = 11, mode = "batch")`,
			schema,
			`validation.schema_violation(path = "amount", rule = "int.closed_range", value = "11")`,
		},
		{
			"mode",
			`Quota__Request(label = "Ada", amount = 5, mode = "admin")`,
			schema,
			`validation.schema_violation(path = "mode", rule = "str.one_of", value = "admin")`,
		},
		{
			"policy_bounds",
			`Quota__Request(label = "Ada", amount = 5, mode = "batch")`,
			`Quota__RequestSchema(label_minimum = 5, label_maximum = 1, amount_lower = 1, amount_upper = 10, allowed_modes = Seq<str>["interactive", "batch"])`,
			`validation.schema_violation(path = "", rule = "schema.label_bounds", value = "minimum=5;maximum=1")`,
		},
		{
			"bigint_exact",
			`Quota__Request(label = "Ada", amount = -9007199254740993, mode = "batch")`,
			schema,
			`validation.schema_violation(path = "amount", rule = "int.closed_range", value = "-9007199254740993")`,
		},
	}
	for _, order := range [][]string{
		{"quota.can", "scalars.can"},
		{"scalars.can", "quota.can"},
	} {
		for _, v := range vectors {
			if err := runLinkedPure(t, quotaRequestFiles(t), order,
				"quota__request__validate", 1,
				map[string]string{"request": v.request, "schema": v.schema},
				v.expect); err != nil {
				t.Fatalf("linked %s (%v): %v", v.name, order, err)
			}
		}
	}
}

// S1b (a89): linked vectors for the envelope witness. The real
// S1a body runs under the wrapper with no scripts: prefixing,
// rule/value preservation, and whole-envelope success all
// execute.
func TestQuotaEnvelopeLinked(t *testing.T) {
	schema := `Quota__RequestSchema(label_minimum = 1, label_maximum = 12, amount_lower = 1, amount_upper = 10, allowed_modes = Seq<str>["interactive", "batch"])`
	vectors := []struct {
		name     string
		envelope string
		schema   string
		expect   string
	}{
		{
			"ok",
			`Quota__Envelope(request = Quota__Request(label = "Ada", amount = 5, mode = "batch"))`,
			schema,
			`Ok(value = Quota__Envelope(request = Quota__Request(label = "Ada", amount = 5, mode = "batch")))`,
		},
		{
			"root",
			`Quota__Envelope(request = Quota__Request(label = "Ada", amount = 5, mode = "batch"))`,
			`Quota__RequestSchema(label_minimum = 5, label_maximum = 1, amount_lower = 1, amount_upper = 10, allowed_modes = Seq<str>["interactive", "batch"])`,
			`validation.schema_violation(path = "request", rule = "schema.label_bounds", value = "minimum=5;maximum=1")`,
		},
		{
			"field",
			`Quota__Envelope(request = Quota__Request(label = "", amount = 5, mode = "batch"))`,
			schema,
			`validation.schema_violation(path = "request.label", rule = "str.length_scalars", value = "")`,
		},
		{
			"rendered",
			`Quota__Envelope(request = Quota__Request(label = "Ada", amount = 11, mode = "batch"))`,
			schema,
			`validation.schema_violation(path = "request.amount", rule = "int.closed_range", value = "11")`,
		},
	}
	for _, order := range [][]string{
		{"quota.can", "scalars.can"},
		{"scalars.can", "quota.can"},
	} {
		for _, v := range vectors {
			if err := runLinkedPure(t, quotaRequestFiles(t), order,
				"quota__envelope__validate", 1,
				map[string]string{"envelope": v.envelope, "schema": v.schema},
				v.expect); err != nil {
				t.Fatalf("linked %s (%v): %v", v.name, order, err)
			}
		}
	}
}

// TestQuotaRequestLinkedContradiction pins failure detection:
// a wrong rendered value and a wrong rule token must both fail.
func TestQuotaRequestLinkedContradiction(t *testing.T) {
	schema := `Quota__RequestSchema(label_minimum = 1, label_maximum = 12, amount_lower = 1, amount_upper = 10, allowed_modes = Seq<str>["interactive", "batch"])`
	request := `Quota__Request(label = "Ada", amount = 11, mode = "batch")`
	files := quotaRequestFiles(t)
	order := []string{"quota.can", "scalars.can"}
	err := runLinkedPure(t, files, order, "quota__request__validate", 1,
		map[string]string{"request": request, "schema": schema},
		`validation.schema_violation(path = "amount", rule = "int.closed_range", value = "12")`)
	if err == nil {
		t.Fatal("expected payload mismatch for wrong render, got pass")
	}
	err = runLinkedPure(t, files, order, "quota__request__validate", 1,
		map[string]string{"request": request, "schema": schema},
		`validation.schema_violation(path = "amount", rule = "str.one_of", value = "11")`)
	if err == nil {
		t.Fatal("expected payload mismatch for wrong rule, got pass")
	}
}
