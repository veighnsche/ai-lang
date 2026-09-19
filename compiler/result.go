package main

// B07: a bare variant return uses the existing Ok protocol with exactly
// one field, value. The variant keeps its own tag inside that field; its
// tag never doubles as the outcome discriminant. Records keep their
// existing field-shaped success payloads.
func variantOkFields(ret string) [][2]string {
	return [][2]string{{"value", ret}}
}

// variantOkType names a checker-only shape for on-Ok binders. The colon
// cannot occur in a source type name: this is not a first-class Outcome
// type or an implicitly declared public record.
func variantOkType(ret string) string {
	return "ok:" + ret
}
