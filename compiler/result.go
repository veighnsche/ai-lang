package main

// B07: a bare variant return uses the existing Ok protocol with exactly
// one field, value. The variant keeps its own tag inside that field; its
// tag never doubles as the outcome discriminant. Records keep their
// existing field-shaped success payloads.
func variantOkFields(ret string) [][2]string {
	return [][2]string{{"value", ret}}
}

// successFields is the shared record/variant success contract used by
// source returns, Fn heads, and invocation temporaries. Other bare
// successes are deliberately not admitted here.
func successFields(ret string, records map[string][][2]string, variant bool) ([][2]string, bool) {
	if variant {
		return variantOkFields(ret), true
	}
	fields, ok := records[ret]
	return fields, ok
}

// successBinderType keeps ordinary calls and invoke on the same ABI:
// records expose their fields; variant successes expose exactly .value.
func (c *tycker) successBinderType(ret string) string {
	if c.variants[ret] {
		return variantOkType(ret)
	}
	return ret
}

// variantOkType names a checker-only shape for on-Ok binders. The colon
// cannot occur in a source type name: this is not a first-class Outcome
// type or an implicitly declared public record.
func variantOkType(ret string) string {
	return "ok:" + ret
}
