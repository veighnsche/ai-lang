# B10 — remaining bare returns and data-only extern successes

Status: implemented; commit recorded in `stdlib-remaining.md`.

## Scope and convention

B10 completes the existing `Ok(value)` convention for the remaining
supported non-record source returns and applies it to synchronous externs.
Records still flatten their declared fields; no existing ABI is replaced.
Extern declarations stay monomorphic; their return types may name concrete
generic data instances, as the variant example does.

| Return type | Source function | Fn invocation success | Extern success |
|---|---|---|---|
| Record, scalar, variant | yes | data-only | data-only |
| Declared brand | yes | yes | yes |
| Bytes | yes | yes | yes |
| Seq of an already-supported data element | yes | yes | yes |
| Fn | yes, via ordinary source calls | **no** | **no** |

Source factories already returned callables inside records. They can now
return them without that record declaration:

```can
fn data__factory(delta: int) -> Fn<int, int, []> rev 1
  emits []
  tests
    plus(3) => Ok(fnref data__add(delta = 3))
    minus(-3) => Ok(fnref data__add(delta = -3))
  Ok(fnref data__add(delta = delta))
```

A caller receives the callable at `on Ok f` → `f.value`, then passes it to
an ordinary consumer with a Fn parameter. Invocation heads remain bare
names: B10 does not add `invoke f.value`. Generic factories also work;
`data__reference<T>() -> Fn<T,T,[]>` in the sketch returns identity callbacks
for Bytes and a nominal brand.

Every non-record success must use `Ok(x)` or `Ok(value = x)`. Its payload
must match the declared type exactly. Raw values/binders, missing/extra/
repeated fields, wrong types, and nested outcome payloads are rejected.
Sequence literals also reject untyped outcome constructors as elements;
a `Seq<int>[Ok(7)]` must not bypass the payload check via its annotation.

Forwarding reconstructs `Ok(value = r.value)`. For brands, Bytes, Seq, and
Fn it requires the same type (Fn spacing is insignificant), just as scalar
and variant forwarding already did. A brand is not its underlying scalar,
two brands are not interchangeable, a Seq<int> is not Bytes, and a
single-field record is not a value-return envelope.

## Authority boundaries, unchanged

- **Externs are data-only**, transitively through records/sequences. Direct
  Fn results, Fn-bearing records, and Fn-bearing inputs remain prohibited.
  Externs are not referenceable and do not become linked-pure or verified
  merely because their success shape is now expressible.
- **Invocation remains first-order**: Fn input/success/captures are data-only.
  A source factory returning Fn cannot itself be used as a callback. Fn
  returning a function-bearing record remains prohibited too.
- Reference creation still requires a pure reachable graph, no required
  preconditions, exactly one unbound parameter, and data-only captures.
  No runtime closure or capture representation changed.
- **Brand sealing remains module-owned**. A caller may pass/forward a brand
  received from a provider but may not mint that provider's brand in its
  body. Test/script seals remain checked data. Host implementations remain
  trusted at their declared boundary, as with brand fields in records.
- Existing sequence-element rules remain: no variant sequences, nested
  sequences, or callable-bearing sequence elements. Success annotations
  are checked even when no literal exposes the bad element type.
- Foreign uses/revision pins, caller-keyed given tables, complete outcome
  handling, identity error forwarding, and per-arm witnesses are unchanged.
  A missing invocation/error witness is not supplied by provider examples.
- Specialized compiler-owned bridge certificates retain their existing
  shapes. No bridge authority, solver theory, async, or resource design was
  generalized; B05 is untouched.

## Implementation

`valueSuccess` recognizes supported one-value shapes, with brands/variants
requiring declaration lookup. Shared `successFields` drives source and
extern emission, Fn types, invocation temporary unions, and Ok checking.
`successBinderType` now installs private envelope shapes lazily, covering
instantiated Seq and Fn spellings without pre-enumerating them.

`externUnion` uses this same shape instead of searching only record
declarations. Its record field ordering is preserved. Data-only checks
also protect direct emitter entry points. `typeContainsFn` shares the
existing recursive containment traversal between checker and emitter.
Foreign type imports include return-only Seq and Fn dependencies (including
inlined record/error payload fields) and extern-return variants; the sketch's consumer imports its row and variant
types from `data`, but host functions from the declaring `bridge.externs`.

No evaluator changes were needed. Brands still erase to their established
host base representation, Bytes to Uint8Array, Seq to arrays, and Fn to
existing closure thunks. Source and host implementations must still return
an outcome envelope; no runtime auto-unwrapping or host-result validation
was added.

## Decision and evidence

JEV (`jev-1.13.0`, Choice) narrowly selected `all_safe_returns`: probability
0.46, confidence 0.27. Alternatives: externs-first 0.41, data-only-first
0.13, unrestricted/higher-order-host 0.00. This was a low-confidence scope
choice, not a safety verdict. The implementation keeps all stated authority
and containment restrictions and tests the channels independently.

The pre-change probe passed its Bytes row but failed at the extern header:
`extern probe__host returns int: externs return a record type, Ok/errors
script the outcome`. It now compiles. Existing brand/Bytes/Seq bare-return
rejection tests became admission tests; the bare-Fn tests now reject wrong
payloads/unknown heads rather than the type category. Extern callable
rejection is pinned to data-only containment, not record-only syntax.

`sketches/bare-returns` has **40 rows** over three modules. It includes
local/shared extern calls for all four scalar types, a brand, Bytes, a
record sequence, and a generic variant; pure Fn consumers for the newly
admitted data types; generic identities/factories; and forwarding a source
factory's callable through a foreign call and into a local consumer.

`compiler/bare_return_test.go` covers ABI agreement, malformed returns and
scripts, nominal brands/seal authority, callback containment, sequence
restrictions, pins, witnesses, linked execution, LSP, and Node parity.
The committed `bridge.externs.ts` supplies real synchronous host functions
for that parity run. Scripted Can rows prove the interface behavior they
exercise, **not** arbitrary host correctness.

Gates: `go test -count=1 ./...`, `go run ./tools/modcheck`, and
`go run ./tools/gramcheck`; compile the three sketch modules consumer-first
and compare `data.ts`, `bridge.ts`, `use.ts`, and `errors.json` byte-for-byte.
The emitted sketch and host implementation also pass strict TypeScript.
Existing stdlib/host and B06–B09 goldens remain unchanged. No stdlib API
migration, first-class Outcome, generic error algebra, or record Ok splat
is claimed.
