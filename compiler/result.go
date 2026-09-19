package main

// Scalar successes deliberately exclude brands, Bytes, sequences and Fn.
func scalarSuccess(ret string) bool {
	switch ret {
	case "int", "str", "bool", "dec":
		return true
	}
	return false
}

// Value successes use the existing Ok protocol with exactly one field,
// value. A variant keeps its own tag inside that field; its tag never
// doubles as the outcome discriminant. Records keep their field payloads.
func valueOkFields(ret string) [][2]string {
	return [][2]string{{"value", ret}}
}

// successFields is the shared success contract for source returns, Fn
// heads and invocation temporaries. It never infers a shape from examples.
func successFields(ret string, records map[string][][2]string, variant bool) ([][2]string, bool) {
	if variant || scalarSuccess(ret) {
		return valueOkFields(ret), true
	}
	fields, ok := records[ret]
	return fields, ok
}

// successBinderType keeps ordinary calls and invoke on the same ABI:
// records expose their fields; value successes expose exactly .value.
func (c *tycker) successBinderType(ret string) string {
	if c.variants[ret] || scalarSuccess(ret) {
		return valueOkType(ret)
	}
	return ret
}

// valueOkType names a checker-only shape for on-Ok binders. The colon
// cannot occur in a source type name: this is not a first-class Outcome
// type or an implicitly declared public record.
func valueOkType(ret string) string {
	return "ok:" + ret
}
