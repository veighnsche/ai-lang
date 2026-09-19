# B09 — scalar successes

Status: implemented; commit recorded in `stdlib-remaining.md`.

## Contract

Source functions and `Fn<A, R, [errors]>` now support **int, str, bool, and
dec** successes alongside the existing record/variant successes. This
includes generic instantiations. Success remains an outcome, not a raw
scalar:

```can
fn scalar__positive(value: int) -> bool rev 1
  emits []
  tests
    yes(7) => Ok(true)
    no(0) => Ok(false)
  Ok(value > 0)
```

`Ok(x)` and `Ok(value = x)` mean the same thing. There must be exactly one
field named `value`, with the declared scalar type. Missing, repeated,
extra, wrong-typed, and nested-outcome payloads are rejected in bodies,
expected results, and exchange scripts. This deliberately closes the old
unchecked named-scalar-Ok path pinned by `TestPositionalScalarOk`.

Named calls and invocations bind the **Ok envelope**. Read `r.value`, not
`r`; match that value or pass it as data. Raw scalar bodies are rejected,
including raw boolean bodies. Boolean expressions and contract predicates
remain ordinary boolean expressions; no implicit success wrapping was added.

```can
match invoke cb with value
  on Ok r => forward r
```

Forwarding into a scalar result requires the **identical scalar type**.
There is no int/dec coercion, and a one-field record is not interchangeable
with the scalar it contains. Record forwarding keeps its existing
field-shape compatibility; variant forwarding keeps nominal identity.

JEV (`jev-1.13.0`, Choice) selected `primitive_envelope`, probability 1.00,
confidence 1.00; source-only, all-data/extern, and raw-scalar alternatives
all 0.00.

## Implementation

- `successFields` describes the declared success shape across constructor
  checking, source return emission, Fn admission/types, and invocation
  temporaries. Scalar successes use `valueOkFields`; no shape is inferred
  from test literals.
- `successBinderType` gives scalar/variant Ok binders private checker-only
  shapes. Those shapes cannot be named in source or passed as outcome data.
- Scalar Ok fields use the normal constructor checker, replacing the old
  `checkScalarOk` pass-through. Recovery arms are checked against the
  enclosing function's scalar return too.
- Contract admission and callee-summary proof binding use the same envelope
  for int/bool successes. True contracts verify, false ones fail even when
  the example rows agree with the implementation. str/dec remain outside
  the existing proof fragment; no new solver theory or callback-contract
  assumption was added.
- No evaluator representation or closure runtime changed. Emission retains
  bigint for int and the existing exact decimal-string representation for
  dec. Reference captures, arity, purity, required-precondition restrictions,
  exact error sets, containment, indirect-cycle detection, and independent
  provider/caller coverage remain enforced.

### Two existing defects exposed by the sketch

1. Generic header rewriting counted executable calls/references but dropped
   pins for references appearing **only in tests**. A valid
   `fnref scalar__identity<int>()` then failed with `references scalar__identity
   which is not in uses`. References require authority in test data too;
   expansion now retains their used stamps. Missing/wrong pins still fail.
2. `eachRaise` counted every reference rooted at an error binder as a raise,
   including `e.value` used to compute a successful recovery. Only a reference
   to the whole error is a raise; payload projections are data. Explicit
   error constructors and identity forwarding still count. This also removes
   a false `html.html__render__document` entry from the `html.invalid_identifier`
   catalogue's `raised_by` list. HTML source and emitted TS are unchanged.

## Evidence and boundaries

The pre-change four-type generic identity probe passed all four rows, then
failed emission with `probe__identity$T$bool returns unknown type bool`.
It now compiles and emits for all four types. The earlier B07 int-only
emission blocker is closed too.

`sketches/fn-scalar` has **40 rows**: generic scalar identity/application,
int-to-bool callbacks, captured callback factories, decimal arithmetic,
foreign exchange rows, local/generic call forwarding, arm-position
`forward call`, checked invocation, recovery from an error payload,
record projection, and matching a boolean invocation result. Goldens,
linked execution, LSP, and Node parity are tested. Node holds independent
captures alive and exercises bigints, false booleans, empty/Unicode strings,
and decimal fractions.

Bare brand, Bytes, sequence, and Fn returns remain unsupported. Externs
still return records (`TestScalarSuccessExternsStillRefused`); unsupported
Fn successes are pinned by `TestFnUnsupportedSuccessesStillRefused`.
Generic `forward call` text remains the separate B08-recorded gap; explicit
`match call f<T>(...)` plus `forward r` works. No first-class Outcome,
generic error algebra, record Ok splat, nested user-generic arguments, or
stdlib API migration is claimed. Re-probe exact optional APIs before
claiming a stdlib slice.

Gates:

```sh
go test -count=1 ./...
go run ./tools/modcheck
go run ./tools/gramcheck
go run ./compiler --out /tmp/b09-fn-scalar \
  sketches/fn-scalar/use.can sketches/fn-scalar/scalar.can
```

Fresh `scalar.ts`, `use.ts`, and `errors.json` must match committed goldens.
The HTML catalogue is regenerated from HTML + ASCII + scalars; its TS is
byte-identical. The JSON stack's 785 rows and goldens remain unchanged.
