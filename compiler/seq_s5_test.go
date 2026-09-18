package main

import (
	"strings"
	"testing"
)

// a42 P1: same-brand str-backed + yields the brand (checker-only;
// eval and emit erase brands already, so no runtime arm exists).
// Fragment assembly rides this rule; nothing else changes.

const seqBrandPlusMod = `mod m
  provides [m__go, M__B, M__C, M__BOut]
  uses []
  emits []

brand M__B is str rev 1

brand M__C is str rev 1

type M__BOut rev 1 (
  v: M__B
)

fn m__go(a: M__B) -> M__BOut rev 1
  emits []
  tests
    go(seal M__B("x")) => Ok(v = seal M__B("xy"))
=
  Ok(v = a + seal M__B("y"))
`

func TestSeqBrandPlus(t *testing.T) {
	seqClean(t, map[string]string{"m.can": seqBrandPlusMod}, "m.can")
}

// The result is the brand, not str: corrupt the field to str and
// the build must refuse, even though both erase identically.
func TestSeqBrandPlusTyped(t *testing.T) {
	// Field and expectation both go str (the erased strings would
	// compare equal, so only the checker can catch this): the one
	// remaining error is the body's brand-vs-str refusal.
	bad := strings.Replace(seqBrandPlusMod, "  v: M__B", "  v: str", 1)
	bad = strings.Replace(bad, `Ok(v = seal M__B("xy"))`, `Ok(v = "xy")`, 1)
	seqCode(t, map[string]string{"m.can": bad}, "m.can", CodeTypeMismatch, "want str")
}

func TestSeqBrandPlusRefusals(t *testing.T) {
	for _, c := range []struct {
		name, expr, sub string
	}{
		// Brand mixed with raw string: no implicit conversion.
		{"brand plus str", `a + "y"`, "no implicit conversions"},
		{"str plus brand", `"x" + a`, "no implicit conversions"},
		// Cross-brand: never unified.
		{"cross brand", `a + seal M__C("y")`, "no implicit conversions"},
		// Brand arithmetic beyond + stays refused.
		{"brand minus", `a - a`, "no implicit conversions"},
	} {
		t.Run(c.name, func(t *testing.T) {
			body := strings.Replace(seqBrandPlusMod,
				`Ok(v = a + seal M__B("y"))`, `Ok(v = `+c.expr+`)`, 1)
			seqCode(t, map[string]string{"m.can": body}, "m.can", CodeTypeMismatch, c.sub)
		})
	}
}

// The emitted concat is native string + over the erased type:
// (a + "y"), with no helper call.
func TestSeqBrandPlusEmit(t *testing.T) {
	ts := compileEmit(t, seqBrandPlusMod)
	if !strings.Contains(ts, `(a + "y")`) {
		t.Fatalf("emit missing native concat:\n%s", ts)
	}
	for _, helper := range []string{"$canStr", "$canSeq", "$canDec", "$canDivMod", "$canEqRec"} {
		if strings.Contains(ts, helper) {
			t.Fatalf("emit should use no helper for brand concat, found %s:\n%s", helper, ts)
		}
	}
}
