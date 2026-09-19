# B11 — whole generic success values, including records

Status: implemented; commit recorded in `stdlib-remaining.md`.

## Surface

```can
fn success__id<T>(value: T) -> T rev 1
  emits []
  tests
    pair<T=Success__Pair>(Success__Pair(7, true)) => Ok(7, true)
    integer<T=int>(7) => Ok(7)
  Ok<T>(value)
```

`Ok<T>(value)` (or `Ok<T>(value = expression)`) takes **one whole value
of exactly T**, and T must match the enclosing success boundary. This applies
to source exits, expected successes, and scripted successes. It is not a
first-class outcome or a record constructor usable as data.

```can
match invoke transform with input
  on Ok<U> result => Ok<U>(result)
```

`on Ok<T> result` checks T against the call/invoke's declared success type and
binds **the T itself**. It works for records, supported value successes, and
ordinary-call Fn factories. A record field literally named `value` stays an
ordinary field; a one-field record is not interchangeable with its scalar.
An empty record binds an empty record, not a missing value.

For example, mapping a callback's record result into an optional value is now:

```can
match invoke transform with s.value
  on Ok<U> mapped => Ok<Success__Option<U>>(Success__Some<U>(mapped))
```

Types may be parameters or supported concrete type expressions, including
`Success__Box<T>`, `Seq<T>`, and `Fn<T, T, []>`. These are annotations, not
requests to instantiate a generic data type named `Ok`. Existing limitations
on nested user-generic **arguments to functions/cases** are unchanged.

## Compatibility and boundaries

- **No ABI change.** Record successes still flatten declared fields;
  other supported successes still use `{ $can_kind: "ok", value: ... }`.
- Existing `Ok(...)` and `on Ok r` keep their meanings. `Ok(record)` is
  **not** silently reinterpreted as a splat. Existing TS/catalogue goldens stay
  unchanged.
- Typed construction checks nominal identity, not just field compatibility.
  It rejects missing/extra arguments, wrong annotations, raw outcomes as data,
  wrong brands, and int/dec coercion. `forward` keeps its existing exact-value
  and record-field compatibility rules; typed scalar binders need no `.value`.
- Callback input/success/captures and extern signatures remain data-only.
  Source factories can return Fn, but cannot themselves be callbacks. There
  is no new host callback, sealing, bridge-certificate, effect, or proof authority.
- Invocation resolution still admits Fn-typed **parameters**, not arbitrary
  new binders. Pass a factory's typed Fn binder to an ordinary consumer;
  `on Ok<Fn<...>> cb` does not by itself enable `invoke cb`. Field-path invoke
  is also unchanged.
- Synthetic store/dec-parts payloads have no declared success type and reject
  typed Ok patterns. Public kernel record successes retain their exact named
  types (e.g. `Bytes__Value`, not `Bytes`).
- Error identity relays, complete dispatch, per-instance arm witnesses,
  invocation witnesses, and foreign uses/revision pins remain mandatory.
- Contracts still use their existing outcome-envelope binders. Int/bool and
  supported records are proved at typed exits and through typed call binders;
  this does not add a str/dec/Fn proof theory.

## Implementation

Parser/expansion retain typed Ok annotations while stamping their component
types, including annotations in patterns and scripts. The checker validates
whole-value construction and exact binder annotations independently of the
legacy field-based construction path. No outcome gets a data type.

The evaluator and emitter convert at the outcome boundary. A record operand
is evaluated once, then its declared fields form the success. Typed record
binders project a plain record without the outcome discriminant. TS uses a
single-evaluation IIFE with a structural parameter type, avoiding unnecessary
nominal record imports. No object spread or implicit host validation is added.
Proof lowering uses equivalent pure record projections; contracts continue
to see the unchanged envelope shape.

**Generic given-row routing fix:** mixed int/record specializations exposed
scripts from other test specializations being checked against the wrong
signature. Given entries now follow the caller test rows during stamping,
including chain-step tables. Only keys naming tests assigned to another
stamp are removed; unknown keys remain, missing scripts still fail, and no
witness is synthesized.

## Evidence and scope

Pre-change optional `value_or<T> -> T`, instantiated with a two-field record,
failed with `Ok field left: got Probe__Pair, want int`. Replacing the two exits
with `Ok<T>(...)` closes that exact probe without a wrapper return type.

`sketches/success-values/` has **54 rows over two modules**, including:

- generic identity over records (multi-field, `value` field, empty), int/bool/
  str/dec, Bytes, Seq<int>, and a brand;
- generic-record construction, local/foreign calls, typed expected/scripted
  results, typed ordinary-call Fn factories, and linked callback execution;
- the §1.7 optional **API shapes**: `require`, `value_or`, and `map`, including
  int-to-record mapping, absence without invocation, and exact error forwarding;
- a real synchronous host returning a record or an error.

This is a language/prototype slice, **not a canonical `std/option` release**.
The optional shapes are now expressible with user-declared variants; naming,
packaging, and `seq.find` migration remain a separate stdlib slice. First-class
`Outcome<T,E>` and generic error-set algebra still block outcome combinators
and the catalogue's fully generic `schema__migrate`. A fixed-error or total
callback can now return arbitrary supported record `New`; `! E` is not added.

`compiler/success_value_test.go` covers static negatives, routing, linked/LSP
parity, int/bool/record proofs and false contracts, kernel binding, byte-identical
goldens, and Node parity (including actual host execution and absence of nested
outcome tags). Higher-order/host/sequence boundary regressions remain green.

## Reproduction

```sh
go run ./compiler --out /tmp/b11-success \
  sketches/success-values/use.can sketches/success-values/success.can
# Compare success.ts, use.ts, errors.json with the committed sketch.
go test -count=1 ./...
go run ./tools/modcheck
go run ./tools/gramcheck
./tscheck/node_modules/.bin/tsc --noEmit --strict --target ES2020 \
  --module ESNext --moduleResolution bundler sketches/success-values/*.ts
```

## JEV decisions

All are `jev-1.13.0`, Choice; distributions are scope/surface judgments, not
proofs of soundness.

- Next focus: `generic_success`, probability/confidence **0.98/0.98**;
  outcome algebra 0.01, async 0.01, higher-order source/host callbacks 0.00.
  User approved this focus.
- Surface: `typed_ok`, probability **1.00**, confidence **0.99**;
  dedicated intrinsics, splat-plus-extractor, and breaking uniform ABI 0.00.
- Related routing fix: `route_givens`, probability **0.97**, confidence
  **0.95**; defer 0.03. Route scripts with their test specializations without
  relaxing missing-script or witness rules.
