# B07 — bare variant returns through Ok(value)

Status: implemented; commit recorded in `stdlib-remaining.md`.
B06 supplied generic variant values. B07 lets source functions name a
variant directly as their return type. No grammar change is needed.

## Contract

```can
fn pick__choose<T>(present: bool, value: T) -> Pick__Value<T> rev 1
  emits []
  tests
    some<T=int>(true, 7) => Ok(Pick__Some<int>(7))
    none<T=int>(false, 7) => Ok(Pick__None<int>())
  match present
    true => Ok(Pick__Some<T>(value))
    false => Ok(Pick__None<T>())
```

The success payload has **exactly one field**, `value`, whose nominal type
is the declared variant. Both `Ok(variant)` and `Ok(value = variant)` work.
A nullary case still needs an explicit constructor inside `Ok`.

Callers retain the existing outcome protocol:

```can
match call pick__choose<int>(present, value)
  on Ok r => match r.value
    on Pick__Some<int> s => Ok(s.value)
    on Pick__None<int> _ => Ok(alternative)
```

This excerpt assumes a record return carrying an integer field. Foreign
calls still need the usual per-caller-name exchange rows. Every outcome
arm and nested case arm must be witnessed, per generic function instance.

The `r` binder is the Ok payload, **not** the variant. The checker gives
it a private `ok:<variant>` record shape containing `.value`; it cannot
be named in source or passed around as a first-class outcome. Use
`r.value` to obtain the variant. Case binders retain their existing
case-specific projection rules; reconstruct a case when needed.

Raw variant constructors or references in function result positions are
rejected: return `Ok(...)` or a declared error. Unknown, missing, repeated,
extra, and mistyped Ok fields are rejected in bodies, expectations, and
exchange outcomes, including recovery bodies under error arms. Nested
`Ok` or error constructors cannot masquerade as variant data.

## Forwarding and emission

- `on Ok r => forward r` rebuilds `Ok(value = r.value)` when both source
  and destination returns name the **same nominal variant instance**.
- Different parents, different generic arguments, and variant/record
  conversions are rejected by forwarding. Explicit `Ok(r.value)` remains
  available where a destination record actually declares that field type.
- Existing record-to-record forwarding is unchanged.
- Arm-position `forward call` works for supported local calls; it retains
  the existing expansion into Ok plus each emitted error arm. It does not
  bypass witness requirements or foreign exchange rules.
- Runtime shape is `{ $can_kind: "ok", value: <tagged variant> }`.
  Variant tags stay nested, distinct from the outcome discriminator.
- TypeScript result types come from the declared return type, not example
  literals. Foreign return-only variant references get type imports too.
- The existing evaluator already implements this envelope. No alternate
  runtime path or special case unwrapping was added.

## Deliberate boundaries

JEV (`jev-1.13.0`, Choice) selected:

- `ok_value`: probability 0.99, confidence 0.98 (unwrap binder 0.01,
  raw case results 0.00).
- `calls_only`: probability 0.75, confidence 0.63 (calls plus Fn successes
  0.25, all channels including externs 0.00).

Accordingly, **Fn success types/reference targets and extern returns stay
record-only**. Their admission and invocation/host ABI are separate
follow-ups, pinned by `TestVariantReturnFnAndExternStillRefused`.

Bare scalar return emission is not generalized: the re-probe
`fn probe__identity<T>(value: T) -> T` at `T=int`, with an `Ok(value)`
body and passing row, still fails emit with
`probe__identity$T$int returns unknown type int`. Other existing bare
return restrictions (brands, Bytes, sequences, Fn) remain unchanged.

No first-class `Outcome`, generic error algebra, arbitrary record splat,
nested generic argument support, variant sequences, or async is added.
Universal verification's variant-fragment limits are unchanged. No exact
stdlib option/outcome combinator is claimed by this language slice.

## Evidence and gates

The pre-implementation probe changed the B06 sketch's `option__some<T>`
return from `Option__Box<T>` to `Option__Value<T>`; it failed with
`bare-variant returns are unsupported, return a record`. It now compiles.
The old rejection test now pins the required payload on nullary cases.

`sketches/variant-return` has 27 rows: both generic cases at int/str,
identity returns, local and foreign forwarding, error propagation,
error-to-variant recovery, and consumer case elimination. Tests cover
monomorphic variants, `-> T` instantiated at a monomorphic variant,
malformed payloads, nominal mismatches, unwitnessed arms, and invented
provider successes. The new TS and error catalogue are frozen as goldens;
Node parity includes a bigint above JavaScript number precision.

```sh
go test -count=1 ./...
go run ./tools/modcheck
go run ./tools/gramcheck
go run ./compiler --out /tmp/b07-return \
  sketches/variant-return/use.can sketches/variant-return/pick.can
```

`TestGoldenVariantReturn` compares fresh `pick.ts`, `use.ts`, and
`errors.json` byte-for-byte. `TestVariantReturnNodeParity` executes the
emitted functions, including both nested discriminants and declared errors.
