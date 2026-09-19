package main

import (
	"strings"
	"testing"
)

// a45 S1 (B1): Bytes value admission and literal construction. Finite
// immutable byte values denoted only by Bytes(Seq<int>[...]) with
// integer-literal members in 0..255. Decision tables below are the
// scope-doc acceptance rows: V-rows must diagnose clean (a wrong
// expectation fails the run, so clean means the evaluator computed
// exactly the pinned value), T-rows must carry the named code. Reuses
// the Seq S1 harnesses seqClean/seqCode; compileEmit pins TS lowering.

const bytesV0 = `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  vals: Bytes
)

fn m__go() -> M__Out rev 1
  emits []
  tests
    go() => Ok(Bytes(Seq<int>[]))
  Ok(Bytes(Seq<int>[]))
`

// V0: empty is a value, not a failure.
func TestBytesV0Empty(t *testing.T) {
	seqClean(t, map[string]string{"m.can": bytesV0}, "m.can")
}

// V1: nonempty ordered contents with repeats are observably exact.
func TestBytesV1Ordered(t *testing.T) {
	body := strings.ReplaceAll(bytesV0, `Bytes(Seq<int>[])`, `Bytes(Seq<int>[0, 127, 128, 255, 255, 0])`)
	seqClean(t, map[string]string{"m.can": body}, "m.can")
}

// V2: construction is admitted in every value position: params, test
// args, refs, record fields, call args, expectations.
func TestBytesV2Positions(t *testing.T) {
	body := `mod m
  provides [m__id, m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  vals: Bytes
)

fn m__id(x: Bytes) -> M__Out rev 1
  emits []
  tests
    id(Bytes(Seq<int>[1])) => Ok(Bytes(Seq<int>[1]))
  Ok(x)

fn m__go() -> M__Out rev 1
  emits []
  tests
    go() => Ok(Bytes(Seq<int>[2]))
  match call m__id(Bytes(Seq<int>[2]))
    on Ok r => Ok(r.vals)
`
	seqClean(t, map[string]string{"m.can": body}, "m.can")
}

// V3: Seq<Bytes> follows the uniform element rule: construction,
// append, and checked access compose; member mismatches refuse.
func TestBytesV3SeqBytes(t *testing.T) {
	body := `mod m
  provides [m__mk, m__go, B__Out]
  uses []
  emits []

type B__Out rev 1 (
  vals: Seq<Bytes>
)

fn m__mk() -> B__Out rev 1
  emits []
  tests
    mk() => Ok(Seq<Bytes>[Bytes(Seq<int>[0]), Bytes(Seq<int>[255])])
  Ok(Seq<Bytes>[Bytes(Seq<int>[0]), Bytes(Seq<int>[255])])

fn m__go() -> B__Out rev 1
  emits []
  tests
    go() => Ok(Seq<Bytes>[Bytes(Seq<int>[0]), Bytes(Seq<int>[1])])
  Ok(Seq<Bytes>[Bytes(Seq<int>[0])] + Bytes(Seq<int>[1]))
`
	seqClean(t, map[string]string{"m.can": body}, "m.can")
}

func TestBytesV3SeqBytesMismatch(t *testing.T) {
	for _, c := range []struct{ name, body, sub string }{
		{"wrong member", `mod m
  provides [m__go, B__Out]
  uses []
  emits []

type B__Out rev 1 (
  vals: Seq<Bytes>
)

fn m__go() -> B__Out rev 1
  emits []
  tests
    go() => Ok(Seq<Bytes>[Bytes(Seq<int>[0])])
  Ok(Seq<Bytes>[Bytes(Seq<int>[0])] + "s")
`, "Bytes"},
		{"elem mismatch", `mod m
  provides [m__go, B__Out]
  uses []
  emits []

type B__Out rev 1 (
  vals: Seq<Bytes>
)

fn m__go() -> B__Out rev 1
  emits []
  tests
    go() => Ok(Seq<Bytes>[Bytes(Seq<int>[0])])
  Ok(Seq<int>[0])
`, "Seq<Bytes>"},
	} {
		t.Run(c.name, func(t *testing.T) {
			seqCode(t, map[string]string{"m.can": c.body}, "m.can",
				CodeTypeMismatch, c.sub)
		})
	}
}

// V4: structural == is preserved for admitted containers of Bytes:
// records, nested records, and error payloads compare by contents.
func TestBytesV4StructuralEq(t *testing.T) {
	rec := `mod m
  provides [m__go, M__Pair]
  uses []
  emits []

type M__Pair rev 1 (
  a: Bytes,
  b: Bytes
)

fn m__go() -> M__Pair rev 1
  emits []
  tests
    go() => Ok(Bytes(Seq<int>[0, 255]), Bytes(Seq<int>[1]))
  Ok(Bytes(Seq<int>[0, 255]), Bytes(Seq<int>[1]))
`
	seqClean(t, map[string]string{"m.can": rec}, "m.can")
	nested := `mod m
  provides [m__go, M__Inner, M__Outer]
  uses []
  emits []

type M__Inner rev 1 (
  vals: Bytes
)

type M__Outer rev 1 (
  inner: M__Inner,
  tag: str
)

fn m__go() -> M__Outer rev 1
  emits []
  tests
    go() => Ok(M__Inner(Bytes(Seq<int>[7])), "t")
  Ok(M__Inner(Bytes(Seq<int>[7])), "t")
`
	seqClean(t, map[string]string{"m.can": nested}, "m.can")
	payload := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  vals: Bytes
)

error m.boom(value: Bytes)

fn m__go(flag: bool) -> M__Out rev 1
  emits [m.boom]
  tests
    hit(true) => m.boom(Bytes(Seq<int>[9]))
    miss(false) => Ok(Bytes(Seq<int>[]))
  match flag
    true => m.boom(Bytes(Seq<int>[9]))
    false => Ok(Bytes(Seq<int>[]))
`
	seqClean(t, map[string]string{"m.can": payload}, "m.can")
}

// V5: linkage and expectation evidence. A provider computing different
// bytes than the script claims must contradict (CAN3110), never trust;
// a directly wrong expectation fails the run (CAN4200).
const bytesLib = `mod lib
  provides [lib__byte, L__Out]
  uses []
  emits []

type L__Out rev 1 (
  vals: Bytes
)

fn lib__byte(flag: bool) -> L__Out rev 1
  emits []
  tests
    go(true) => Ok(Bytes(Seq<int>[0, 255]))
  Ok(Bytes(Seq<int>[0, 255]))
`

const bytesAppLie = `mod app
  provides [app__go, A__Out]
  uses [lib__byte@1]
  emits []

type A__Out rev 1 (
  vals: Bytes
)

fn app__go(flag: bool) -> A__Out rev 1
  emits []
  tests
    lie(true) => Ok(Bytes(Seq<int>[0]))
  match call lib__byte(flag)
    given
      lie => [exchange args (flag = true) outcome Ok(Bytes(Seq<int>[0]))]
    on Ok v => Ok(v.vals)
`

func TestBytesV5LinkageContradiction(t *testing.T) {
	dir := writeLSPDir(t, map[string]string{"lib.can": bytesLib, "app.can": bytesAppLie})
	diags := diagnose(dir, "app.can", bytesAppLie)
	found := false
	for _, d := range diags {
		if d.Sev == "error" && d.Code == CodeInconsistentScript {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected CAN3110 contradiction, got %v", diags)
	}
	if !hasDiag(diags, "error", "contradicts lib__byte") {
		t.Fatalf("expected contradiction message, got %v", diags)
	}
}

func TestBytesV5WrongExpectation(t *testing.T) {
	body := strings.Replace(bytesV0,
		"  Ok(Bytes(Seq<int>[]))",
		"  Ok(Bytes(Seq<int>[0, 255]))", 1)
	dir := writeLSPDir(t, map[string]string{"m.can": body})
	diags := diagnose(dir, "m.can", body)
	found := false
	for _, d := range diags {
		if d.Sev == "error" && d.Code == CodeTestFailed {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected CAN4200 test failure, got %v", diags)
	}
}

// V6: emitted TypeScript pins. Byte literals lower to numbers (never
// bigints); the equality runtime stays absent when unused and arrives
// exactly once when records of Bytes compare.
func TestBytesV6EmitPins(t *testing.T) {
	body := strings.ReplaceAll(bytesV0, `Bytes(Seq<int>[])`, `Bytes(Seq<int>[0, 127, 128, 255, 255, 0])`)
	ts := compileEmit(t, body)
	if !strings.Contains(ts, `Uint8Array.from([0, 127, 128, 255, 255, 0])`) {
		t.Fatalf("emit missing numeric byte lowering:\n%s", ts)
	}
	if strings.Contains(ts, "255n") {
		t.Fatalf("emit leaked bigint into bytes:\n%s", ts)
	}
	if strings.Contains(ts, "$canEqVal") {
		t.Fatalf("emit inserted unused equality runtime:\n%s", ts)
	}
}

func TestBytesV6EmitNestedEqOnce(t *testing.T) {
	body := `mod m
  provides [m__eq1, m__eq2, M__Pair, M__Flag]
  uses []
  emits []

type M__Pair rev 1 (
  a: Bytes,
  b: Bytes
)

type M__Flag rev 1 (
  flag: bool
)

fn m__eq1(p: M__Pair, q: M__Pair) -> M__Flag rev 1
  emits []
  tests
    eq(M__Pair(a = Bytes(Seq<int>[0]), b = Bytes(Seq<int>[1])), M__Pair(a = Bytes(Seq<int>[0]), b = Bytes(Seq<int>[1]))) => Ok(true)
  Ok((p == q))

fn m__eq2(p: M__Pair, q: M__Pair) -> M__Flag rev 1
  emits []
  tests
    ne(M__Pair(a = Bytes(Seq<int>[0]), b = Bytes(Seq<int>[1])), M__Pair(a = Bytes(Seq<int>[9]), b = Bytes(Seq<int>[1]))) => Ok(false)
  Ok((p == q))
`
	ts := compileEmit(t, body)
	if got := strings.Count(ts, "function $canEqBytes"); got != 1 {
		t.Fatalf("expected one $canEqBytes helper, got %d:\n%s", got, ts)
	}
	if !strings.Contains(ts, "$canEqVal") {
		t.Fatalf("emit missing structural equality runtime:\n%s", ts)
	}
}

// T0: members outside 0..255 are static range failures with index and
// exact value. Only the body offends; the expectation stays clean.
func TestBytesT0ElementRange(t *testing.T) {
	for _, c := range []struct{ name, frag, sub string }{
		{"negative", `Bytes(Seq<int>[-1])`, "-1"},
		{"over", `Bytes(Seq<int>[256])`, "256"},
		{"huge", `Bytes(Seq<int>[123456789012345678901234567890])`, "123456789012345678901234567890"},
	} {
		t.Run(c.name, func(t *testing.T) {
			body := strings.Replace(bytesV0,
				"  Ok(Bytes(Seq<int>[]))",
				"  Ok("+c.frag+")", 1)
			seqCode(t, map[string]string{"m.can": body}, "m.can",
				CodeBytesElementRange, c.sub)
		})
	}
}

// T1: runtime members and runtime sequences are shape failures, not
// conversions. There is no general Seq<int> -> Bytes spelling in v1.
func TestBytesT1LiteralShape(t *testing.T) {
	runtimeMember := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  vals: Bytes
)

fn m__go(x: int) -> M__Out rev 1
  emits []
  tests
    go(x = 3) => Ok(Bytes(Seq<int>[]))
  Ok(Bytes(Seq<int>[x]))
`
	seqCode(t, map[string]string{"m.can": runtimeMember}, "m.can",
		CodeBytesLiteral, "not an integer literal")
	runtimeSeq := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  vals: Bytes
)

fn m__go(xs: Seq<int>) -> M__Out rev 1
  emits []
  tests
    go(xs = Seq<int>[1]) => Ok(Bytes(Seq<int>[]))
  Ok(Bytes(xs))
`
	seqCode(t, map[string]string{"m.can": runtimeSeq}, "m.can",
		CodeBytesLiteral, "Seq<int> literal")
}

// T2: the diagnostic-position table. Bare lists, annotations, malformed
// heads, and unbound names keep their own codes; Bytes adds nothing.
func TestBytesT2Positions(t *testing.T) {
	bareList := strings.Replace(bytesV0,
		"  Ok(Bytes(Seq<int>[]))",
		"  Ok(Bytes([0, 1]))", 1)
	seqCode(t, map[string]string{"m.can": bareList}, "m.can",
		CodeSeqLiteral, "Seq<T>")
	bareAnnot := `mod m
  provides [m__go, M__Out]
  uses []
  emits []

type M__Out rev 1 (
  vals: Bytes
)

fn m__go(x: Seq) -> M__Out rev 1
  emits []
  tests
    go(0) => Ok(Bytes(Seq<int>[]))
  Ok(Bytes(Seq<int>[]))
`
	seqCode(t, map[string]string{"m.can": bareAnnot}, "m.can",
		CodeUnknownType, "Seq")
	unbound := strings.Replace(bytesV0,
		"  Ok(Bytes(Seq<int>[]))",
		"  Ok(Bytes(Seq))", 1)
	seqCode(t, map[string]string{"m.can": unbound}, "m.can",
		CodeTypeMismatch, "unbound name Seq")
}

func TestBytesT2ParseShapes(t *testing.T) {
	for _, c := range []struct{ name, frag string }{
		{"empty elem", "Seq<>[]"},
		{"nested", "Seq<Seq<int>>[Seq<int>[0]]"},
		{"missing bracket", "Seq<int>"},
	} {
		t.Run(c.name, func(t *testing.T) {
			body := strings.Replace(bytesV0, `Bytes(Seq<int>[])`, `Bytes(`+c.frag+`)`, -1)
			dir := writeLSPDir(t, map[string]string{"m.can": body})
			diags := diagnose(dir, "m.can", body)
			found := false
			for _, d := range diags {
				if d.Sev == "error" && d.Code == CodeParse {
					found = true
				}
			}
			if !found {
				t.Fatalf("expected parse error, got %v", diags)
			}
		})
	}
}

// T3: direct byte operators are deferred. Structural record comparison
// (V4/V6) is the test surface, not these spellings.
func TestBytesT3NoDirectComparison(t *testing.T) {
	for _, c := range []struct{ name, op string }{
		{"eq", "=="},
		{"ne", "!="},
		{"order", ">"},
	} {
		t.Run(c.name, func(t *testing.T) {
			body := strings.Replace(`mod m
  provides [m__go, M__Flag]
  uses []
  emits []

type M__Flag rev 1 (
  flag: bool
)

fn m__go() -> M__Flag rev 1
  emits []
  tests
    go() => Ok(true)
  Ok((Bytes(Seq<int>[0]) OP Bytes(Seq<int>[0])))
`, "OP", c.op, 1)
			seqCode(t, map[string]string{"m.can": body}, "m.can",
				CodeTypeMismatch, "Bytes comparison")
		})
	}
}

// B10 admits bare-Bytes successes; Bytes state stays out.
func TestBytesT4BareReturn(t *testing.T) {
	body := `mod m
  provides [m__go]
  uses []
  emits []

fn m__go() -> Bytes rev 1
  emits []
  tests
    go() => Ok(Bytes(Seq<int>[]))
  Ok(Bytes(Seq<int>[]))
`
	_, _, ds := genericProgram(t, map[string]string{"m.can": body})
	if len(genericErrs(ds)) != 0 {
		t.Fatal(ds)
	}
}

func TestBytesT4NoBytesState(t *testing.T) {
	body := strings.Replace(bytesV0, `type M__Out rev 1 (`,
		"state M__C: Bytes = Bytes(Seq<int>[])\n\ntype M__Out rev 1 (", 1)
	seqCode(t, map[string]string{"m.can": body}, "m.can",
		CodeUnknownType, "cells hold")
}

// T5: the Bytes name is reserved. Source records and brands cannot
// shadow the primitive.
func TestBytesT5NoShadow(t *testing.T) {
	rec := strings.Replace(bytesV0, `type M__Out rev 1 (`, `type Bytes rev 1 (
  x: int
)

type M__Out rev 1 (`, 1)
	rec = strings.Replace(rec, "provides [m__go, M__Out]", "provides [m__go, M__Out, Bytes]", 1)
	// The shadow name also breaks the Domain__Name rule (CAN2002);
	// both fire, each owning its rule, so pin presence rather than
	// the single-primary shape.
	dir := writeLSPDir(t, map[string]string{"m.can": rec})
	diags := diagnose(dir, "m.can", rec)
	if !hasDiag(diags, "error", "shadows the Bytes primitive") {
		t.Fatalf("expected CAN6012 shadow rejection, got %v", diags)
	}
	brand := strings.Replace(bytesV0, `type M__Out rev 1 (`, `brand Bytes is str rev 1

type M__Out rev 1 (`, 1)
	brand = strings.Replace(brand, "provides [m__go, M__Out]", "provides [m__go, M__Out, Bytes]", 1)
	// Companion CAN2002 as above: presence, not single-primary.
	bdir := writeLSPDir(t, map[string]string{"m.can": brand})
	bdiags := diagnose(bdir, "m.can", brand)
	if !hasDiag(bdiags, "error", "shadows the Bytes primitive") {
		t.Fatalf("expected CAN6012 shadow rejection, got %v", bdiags)
	}
}
