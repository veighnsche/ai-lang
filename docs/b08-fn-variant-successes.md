# B08 — variant successes in Fn

Status: implemented; commit recorded in `stdlib-remaining.md`.
This extends B07's success contract through function values. It changes no
syntax, runtime calling convention, extern ABI, or generic error algebra.

## Surface and semantics

`Fn<A, R, [errors]>` now admits a declared **record or variant** as `R`,
including generic stamps. Input and success identity and concrete error
sets still match exactly; equal-looking variants are not interchangeable.

```can
fn use__apply<T>(cb: Fn<T, Choice__Value<T>, []>, value: T) -> Choice__Value<T> rev 1
  emits []
  tests
    some<T=int>(fnref choice__choose<int>(present = true), 7) => Ok(Choice__Some<int>(7))
    none<T=int>(fnref choice__choose<int>(present = false), 7) => Ok(Choice__None<int>())
  match invoke cb with value
    on Ok r => forward r
```

The complete provider and all its independent witnesses live in
`sketches/fn-variant/choice.can`. The callback above captures a boolean and
leaves its value parameter unbound. Creating the reference still does not
execute it. Factory records may store `Fn<T, Choice__Value<T>, []>`;
existing generic expansion discovers and substitutes those variant types.

Variant success is still exactly `Ok(value = variant)`. Invocation's
`on Ok r` binder exposes the variant at `r.value`, exactly like a named
call. It can be matched, passed as data, returned with `Ok(r.value)`, or
forwarded into the same nominal variant result. The entire Ok binder is
not a first-class outcome value. Declared errors retain their identities
and payloads and can be forwarded or explicitly recovered from.

## Implementation and invariants

- Reference-target admission, `refFnType`, and Fn-head validation admit
  record/variant successes together. Signature discovery used by the
  indirect-cycle checker sees the same newly eligible targets.
- Shared `successFields` supplies record fields or the one-field variant
  envelope to source return checking, source result emission, Fn type
  emission, and invocation temporary emission. No type comes from example
  output literals.
- Shared `successBinderType` keeps call and invoke binders aligned.
- Existing closure thunks, capture snapshots, runtime signature checks,
  provider preparation, and linked-pure execution are unchanged. No
  evaluator fallback or special unwrapping was added.
- Exact input/result/error compatibility, pure reachable graphs, data-only
  captures/input/success, required-precondition refusal, and indirect-cycle
  rejection all remain enforced. Variants carrying Fn remain invalid.
- Invocation still takes no `given` table and grants no target coverage.
  Invocation outcome arms and nested variant cases require real caller
  witnesses. `forward e` does not waive unwitnessed invocation errors.

JEV (`jev-1.13.0`, Choice) selected `record_and_variant`: probability 0.99,
confidence 0.98; admission-only 0.01, generalize-all-data 0.00. The change
therefore covers admission **and** invocation/emission in one green slice.

## Retained boundaries

These describe B08's scope. [B09](b09-scalar-successes.md) subsequently
completes int/str/bool/dec successes across source returns, calls, Fn,
invocation, and emission. Other bare types and externs remain separate.

Primitive scalar, brand, Bytes, sequence, and Fn success types are not
newly admitted. Scalar source-return emission is still a separate gap.
Extern returns remain record-only, and Fn-bearing host signatures remain
forbidden. Error lists stay concrete, exact, and canonically ordered;
there is no generic error-set algebra or first-class Outcome value.

No field-path invocation, captured callbacks, effectful callbacks,
higher-order success carriers, nested user-generic arguments, or universal
callable-contract verification was added. No stdlib combinator is claimed.

### Separate gap found while composing

`forward call use__apply<int>(bound.cb, value)` fails with
`forward call cannot resolve callee use__apply`: the forward-call sugar
stores source text and parses it after generic expansion. The stamped
callee no longer has that source base name. This is pinned by
`TestFnVariantGenericForwardCallStillRefused`, not broadened in B08.

Use the supported explicit form:

```can
on Ok bound => match call use__apply<int>(bound.cb, value)
  on Ok r => forward r
```

The relay linter previously recommended the failing sugar for this valid
form. B08 suppresses that unsafe recommendation for explicit generic
calls; monomorphic relay suggestions are unchanged. The sketch is clean
under LSP diagnosis, which includes lint checks.

## Evidence and gates

The pre-change monomorphic `Fn<int, Probe__Value, []>` probe failed with
`Fn success Probe__Value of param cb is not a record: name the Ok payload
record`; it now compiles through an invoking/forwarding consumer.

`sketches/fn-variant` has **26 rows** covering generic int/str factories
and invocations, captured booleans, checked variant callbacks, success and
error forwarding, recovery, case payload projection, and a factory →
consumer call chain. The provider retains its own case/branch witnesses.

Additional compiler regressions cover monomorphic variants, nominally
incompatible same-shaped parents and generic instances, wrong inputs/error
sets, malformed recovery, missing cases/outcomes/pins, computed captures,
purity/preconditions, Fn containment, indirect cycles, and both provider
and invocation-arm coverage isolation. Go linked execution, LSP diagnosis,
TS goldens, and Node parity are gated; parity holds two captures alive and
checks a bigint beyond JavaScript number precision.

```sh
go test -count=1 ./...
go run ./tools/modcheck
go run ./tools/gramcheck
go run ./compiler --out /tmp/b08-fn-variant \
  sketches/fn-variant/use.can sketches/fn-variant/choice.can
```

`TestGoldenFnVariant` compares fresh `choice.ts`, `use.ts`, and
`errors.json` byte-for-byte. `TestFnVariantNodeParity` executes the emitted
callbacks. Existing record-success Fn and B07 call-return goldens stay
unchanged.
