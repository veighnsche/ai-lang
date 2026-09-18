package main

import (
	"os"
	"testing"
)

// S2b third consumer: linked-pure vectors for
// notify__request__validate and notify__request__admit.
// runLinkedPure executes the real bodies across notify.can,
// quota.can, and scalars.can with no scripts: the S2a check
// ops run under the consumer's binding chain, with the
// shared producers and converter beneath them. Both module
// orders.

func notifyFiles(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	for _, f := range []string{"notify.can", "quota.can", "scalars.can"} {
		raw, err := os.ReadFile(map[string]string{
			"notify.can":  "../sketches/notify/notify.can",
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

func TestNotifyLinked(t *testing.T) {
	schema := `Notify__RequestSchema(handle_minimum = 3, handle_maximum = 16, retries_lower = 0, retries_upper = 5, allowed_channels = Seq<str>["email", "sms"])`
	vectors := []struct {
		name    string
		request string
		schema  string
		expect  string
	}{
		{
			"ok",
			`Notify__Request(handle = "sam", retries = 2, channel = "email")`,
			schema,
			`Ok(Notify__Request(handle = "sam", retries = 2, channel = "email"))`,
		},
		{
			"handle",
			`Notify__Request(handle = "ab", retries = 2, channel = "email")`,
			schema,
			`validation.schema_violation(path = "handle", rule = "str.length_scalars", value = "ab")`,
		},
		{
			"retries_renders_through_real_check",
			`Notify__Request(handle = "sam", retries = 6, channel = "email")`,
			schema,
			`validation.schema_violation(path = "retries", rule = "int.closed_range", value = "6")`,
		},
		{
			"channel",
			`Notify__Request(handle = "sam", retries = 2, channel = "push")`,
			schema,
			`validation.schema_violation(path = "channel", rule = "str.one_of", value = "push")`,
		},
		{
			"policy_bounds",
			`Notify__Request(handle = "sam", retries = 2, channel = "email")`,
			`Notify__RequestSchema(handle_minimum = 9, handle_maximum = 4, retries_lower = 0, retries_upper = 5, allowed_channels = Seq<str>["email", "sms"])`,
			`validation.schema_violation(path = "", rule = "schema.handle_bounds", value = "minimum=9;maximum=4")`,
		},
		{
			"bigint_exact",
			`Notify__Request(handle = "sam", retries = -9007199254740993, channel = "email")`,
			schema,
			`validation.schema_violation(path = "retries", rule = "int.closed_range", value = "-9007199254740993")`,
		},
	}
	for _, order := range [][]string{
		{"notify.can", "quota.can", "scalars.can"},
		{"scalars.can", "quota.can", "notify.can"},
	} {
		for _, v := range vectors {
			if err := runLinkedPure(t, notifyFiles(t), order,
				"notify__request__validate", 1,
				map[string]string{"request": v.request, "schema": v.schema},
				v.expect); err != nil {
				t.Fatalf("linked %s (%v): %v", v.name, order, err)
			}
		}
	}
}

func TestNotifyAdmitLinked(t *testing.T) {
	schema := `Notify__RequestSchema(handle_minimum = 3, handle_maximum = 16, retries_lower = 0, retries_upper = 5, allowed_channels = Seq<str>["email", "sms"])`
	vectors := []struct {
		name    string
		request string
		schema  string
		expect  string
	}{
		{
			"admitted",
			`Notify__Request(handle = "sam", retries = 2, channel = "email")`,
			schema,
			`Ok(handle = "sam", retries = 2, channel = "email")`,
		},
		{
			"rejected_before_ticket",
			`Notify__Request(handle = "ab", retries = 2, channel = "email")`,
			schema,
			`validation.schema_violation(path = "handle", rule = "str.length_scalars", value = "ab")`,
		},
	}
	for _, order := range [][]string{
		{"notify.can", "quota.can", "scalars.can"},
		{"scalars.can", "quota.can", "notify.can"},
	} {
		for _, v := range vectors {
			if err := runLinkedPure(t, notifyFiles(t), order,
				"notify__request__admit", 1,
				map[string]string{"request": v.request, "schema": v.schema},
				v.expect); err != nil {
				t.Fatalf("linked %s (%v): %v", v.name, order, err)
			}
		}
	}
}

// TestNotifyLinkedContradiction pins failure detection: a wrong
// tag infix and a wrong rendered bound must both fail through
// the check-op layer.
func TestNotifyLinkedContradiction(t *testing.T) {
	schema := `Notify__RequestSchema(handle_minimum = 9, handle_maximum = 4, retries_lower = 0, retries_upper = 5, allowed_channels = Seq<str>["email", "sms"])`
	request := `Notify__Request(handle = "sam", retries = 2, channel = "email")`
	files := notifyFiles(t)
	order := []string{"notify.can", "quota.can", "scalars.can"}
	err := runLinkedPure(t, files, order, "notify__request__validate", 1,
		map[string]string{"request": request, "schema": schema},
		`validation.schema_violation(path = "", rule = "schema.guest_bounds", value = "minimum=9;maximum=4")`)
	if err == nil {
		t.Fatal("expected payload mismatch for wrong tag, got pass")
	}
	err = runLinkedPure(t, files, order, "notify__request__validate", 1,
		map[string]string{"request": request, "schema": schema},
		`validation.schema_violation(path = "", rule = "schema.handle_bounds", value = "minimum=9;maximum=5")`)
	if err == nil {
		t.Fatal("expected payload mismatch for wrong bound render, got pass")
	}
}
