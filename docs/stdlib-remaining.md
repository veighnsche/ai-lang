# stdlib implementation record — overnight goal, 2026-09-19

Goal: implement all planned stdlib from `ASTRA_STDLIB.md` in
PR-sized green commits, JEV for multiple-choice decisions.
Result: everything buildable with the current language is built
(slices 1–6). The remainder is provably blocked on missing
language features, one per row below. No slice was faked around
a blocker: each block below carries its compiler error or design
verdict as evidence.

Continuation (post-b00, same contract): b00 function values
landed, unblocking the §1.8 callback rows. Slices 7+ below.
B06 subsequently lands explicit generic variants; B07 admits bare variant
returns through the existing Ok(value) protocol; B08 carries that success
contract through Fn references and invocation; B09 completes primitive
scalar successes across source returns, calls, and Fn. These language
slices do not claim the catalogue's option/outcome combinators.

## Shipped (main, all gates green per commit)

| Slice | Commit | Content |
|-------|--------|---------|
| 1 | `5bad56d` | `std__compare__int/dec/str` close §1.4; 15 rows |
| 2 | `b5780eb` | `std__validate__all` + report records; 15 rows |
| 3 | `b4bfe6e` | `html__url__parse` + `origin_end`; 18 rows |
| 4 | `bb05e2c` | NEW `std/seq`: 7 generic ops; 45 rows |
| 5 | `c3b0075` | `validate__schema__str/int` close §1.5; 14 rows + given rows |
| 6 | `a80e1f5` | NEW `std/map`, `std/set`; 95 rows |
| 7 | `44b188b` | NEW `std/seq` map/filter/fold/all/any/find; 62 rows |
| 8 | `02495f4` | NEW `std/seq` sort/unique; 40 rows |
| 9 | `43f0dfc` | NEW `std/json` value layer: AST + render frame machine + escape + 8 scalar codecs + monomorphic schema family with Fn dispatch; 125 rows |
| 11a | `7092467` | `std/json`: render goes total (variant tags, fuel-exhaust `Ok(acc)`, budget error deleted) + parse leaves (ws/head/literal/unescape/9-state numcheck/contains/pop/attach); 234 rows |
| 11b | `f4845bf` | `std/json`: byte-level parse — `parse_value` entry + `parse_step` 12-state machine (value states push, continuations replace, `StrKey` inherits fields/keys, `Tail` rejects trailers); 100 rows (15 value + 85 step), 334 file total |
| 12 | `e7f7e5f` | `std/json`: 8 monomorphic text drivers (`int/str/bool/dec` × `encode_text`/`decode_text`) closing the schema round-trip; 28 rows, 362 file total |
| B06 | `e7ac40c` | Language unblocker: explicit generic variants, case construction/matching, instance-specific tags, cross-module pins, payload dependency discovery; `sketches/generic-option` has 12 rows plus byte-identical TS/catalogue goldens and Node parity. No stdlib combinator slice claimed. |
| B07 | `cb457c1` | Language unblocker: bare variant returns via checked `Ok(value)`, ordinary call binders expose `.value`, same-nominal forwarding, declared TS success shapes/imports; `sketches/variant-return` has 27 rows, TS/catalogue goldens, and Node parity. Fn successes and extern returns were still record-only at B07. |
| B08 | `147260a` | Language unblocker: Fn successes admit records or variants, including generic stamps; reference admission/signature discovery, invocation binders/forwarding, and TS types agree on the B07 envelope. `sketches/fn-variant` has 26 rows, TS/catalogue goldens, linked/LSP/Node checks, and purity/cycle/coverage regressions. Externs and scalar successes are unchanged. |
| B09 | TBD | Language unblocker: int/str/bool/dec successes through strict `Ok(value)` checking, generic source returns, calls, Fn references/invocation, exact-type forwarding, and TS emission. `sketches/fn-scalar` has 40 rows, goldens, linked/LSP/Node parity, and int/bool contract proof tests. Fixes test-only generic-reference pins and false error-projection raises; HTML catalogue corrected, emitted HTML unchanged. |

Pre-existing (§1.1–1.3, §1.6, §1.8 text/codecs, §2 elements/render/
assets, quota, schema, ascii) was verified present, not rebuilt.

## Blocked, with evidence

| Roadmap item | Missing feature | Evidence |
|--------------|-----------------|----------|
| §1.7 outcome/option combinators | First-class outcomes + generic error algebra; generic record-success reconstruction | PARTIAL UNBLOCK (B06–B09): user-declared optional variants support construction, matching, source returns, calls/forwarding, and Fn successes/invocation. B09 also closes primitive scalar source/callback successes, including the generic identity emission probe. No stdlib combinators claimed. Outcome probe remains `unknown type Outcome`; b00 excludes generic error-set algebra and generic outcome values, and the no-Ok-splat limitation below remains for arbitrary record successes. `TestFnUnsupportedSuccessesStillRefused` pins other bare Fn restrictions; `TestScalarSuccessExternsStillRefused` and `TestVariantReturnExternStillRefused` pin record-only extern returns. Re-probe exact optional APIs individually before a stdlib slice; no monomorphic fallback is being substituted. |
| §1.8 seq map/filter/fold/find/all/any | SHIPPED (slice 7): generic workers invoke total callbacks; find reports `sequence.not_found()` (B06–B07 remove the original generic-variant/return blockers; optional-return API migration has not been attempted); predicates return per-module `Bool__Value` (std/set precedent) | 62 rows; 107 pass on compile; `TestStdSeqCompiles` gates. `find` is unchanged. |
| §1.8 seq sort/unique | SHIPPED (slice 8): `Seq__Order` value Asc/Desc over per-instance built-in order (insertion sort, stable by construction); unique keeps first occurrences via per-instance `==` | 40 rows; 147 pass on compile; lexicographic orders deferred. |
| §1.8 `Map<K,V>` fully generic | SHIPPED (slice M1): `Map<K,V>` over `Seq<Map__Pair<K,V>>` with catalogue names; str-keyed `Map__Entries<V>` migrated away (no downstream users); errors payloadless (payloads cannot be generic) | 45 rows across `<str,int>` + `<int,str>`; goldens regen via documented flow; `TestStdMapCompiles` gates. |
| §1.8 normalize_nfc, casefold, graphemes | Unicode data kernel | No pinned data, no host path (see host shelf) |
| §1.8 json encode/decode, `schema__migrate` | PARTIAL (slices 9+11): value layer + byte-level parse shipped — `Json__Value` AST, fuel-bounded render machine, escape, 8 scalar codecs, monomorphic `Json__*Schema` family with schema-carried Fn dispatch, `parse_value`/`parse_step` 12-state machine, 8 monomorphic text drivers (slice 12). `schema__migrate`: VERDICT (slice 10) — generic migrate is inexpressible: no `Ok` splat (`Ok field y: got M__B, want int`) and `forward` is refused outside call-outcome arms, so no body returns an arbitrary `New` through invoke dispatch (b00's missing generic error algebra independently blocks `! E`). | 362 rows; 785 pass on compile; `TestStdJsonCompiles` gates. |
| Host shelf (clock/random/hash/secret/log/env) | SHIPPED (slices H1–H5): `std/host` carries all 7 catalogue fns over pinned externs with real node-backed `host.externs.ts` impls; `TestStdHostNodeSmoke` executes every impl (incl. sha256 known vector); `sketches/host-clock` consumes both the wrapper and the shared extern directly | Millis instants (JEV 0.97); sealed profile/secret/env brands (JEV 0.98/1.0); denied + sub-millis documented v1 limits. |
| §3 HTTP (all) | Async + Resources + Functions | No async surface exists |
| §4 SQL (all) | Async + Resources | Same |
| §5 UI (all) | Functions + Async + Resources | Same |

## Language gaps discovered while slicing (future b-series)

- ~~No `seq[i].field` projection (slice 2; worked around via `_push`)~~ — LANDED as b03 (`proj` node; chains nest, `$canSeqAt(m, i).f` emit).
- ~~`substNode` skipped `InvokeArg` (slice 7; generic invoke args stamped verbatim)~~ — FIXED in-slice (one line; `TestExpandInvokeArgSubstituted`).
- Bare `Seq<T>` returns rejected (slice 4; wrapper records required).
- Per-instance arm coverage for generics (slice 4; row cost is real).
- Explicit `<T>` required on recursive generic calls (slice 4).
- `given` rows key on caller test names across module lines (slice 5).
- Downstream sketch goldens embed provider catalogs; regen together (slice 5).
- Row binds must name the declared param (`<V=>`, not `<T=>`, slice 6).
- ~~CLI ran each module's tests before later modules' statics, so consumer-first invoke executed raw provider bodies (slice 9; `Ok takes 2 args for 1 fields`)~~ — FIXED in-slice (CLI mirror of `prepareProviders`; `TestFnLinkedConsumerFirst`).
- Only direct self-recursion admitted (slice 9; mutual value/fields/array recursion refused) — worked around via single-fn frame machine + fuel.
- Variant sequences not admitted; sequence concatenation not in v1; slice operator is str-only (slice 9) — worked around via tag-dispatched record frames, append-only back stack, copy-by-index `pop`.
- `invoke` heads must be bare names (slice 9; field paths do not parse) — worked around via apply wrappers taking the Fn as a param.
- `forward call` is arm-position-only, never a bare body (slice 9).
- No `Ok` splat: `Ok(r)` binds the whole record to the first field (slice 10; generic migrate blocked).
- No `Outcome<T,E>` / `Option<T>` builtins; no generic error-set algebra (slice 10). B06 now admits user-declared generic variants such as `Option__Value<T>`; exact §1.7 APIs are not yet claimed.
- ~~Generic variants fail with `bad variant decl`~~ — B06 re-probed that exact error before implementation; declarations, explicit case constructors/patterns, and record-carried optional results now work. See `docs/b06-generic-variants.md`.
- ~~Bare variant returns remain unsupported after B06 (`bare-variant returns are unsupported, return a record`)~~ — B07 re-probed and removed this restriction for source functions and ordinary calls/forwarding. Success is exactly `Ok(value)` and callers select the variant with `r.value`; no first-class outcome value was added. See `docs/b07-variant-returns.md`.
- ~~Fn successes/reference targets remain record-only~~ — B08 admits declared record or variant successes end-to-end; B09 adds primitive scalar successes. Exact type/error-set matching remains. Extern returns remain record-only. Nested instantiation as a type argument and variant sequences remain deferred.
- ~~Generic scalar return emission is incomplete~~ — B09 re-probed the four-type identity (four passing rows, then `probe__identity$T$bool returns unknown type bool`) and closed emission plus Fn/invocation/forwarding for int/str/bool/dec. The earlier B07 int-only blocker is closed. Other bare successes remain refused (`TestFnUnsupportedSuccessesStillRefused`). See `docs/b09-scalar-successes.md`.
- ~~Test-only foreign generic references lose uses pins~~ — B09 fixes header rewriting to include reference sites in tests; no authority rule is relaxed. Missing/wrong pins still fail.
- ~~Error payload reads count as raises~~ — B09 fixes `eachRaise`: `e.value` is data, not an error outcome. Identity relays still count. HTML's catalogue loses one falsely attributed raise; emitted code is unchanged.
- Generic `forward call` sugar is late-parsed source text, not rewritten by expansion: `forward call use__apply<int>(...)` fails with `forward call cannot resolve callee use__apply` (`TestFnVariantGenericForwardCallStillRefused`). Explicit `match call ...<int>` plus `forward r` works. B08 fixes the relay linter's unsafe suggestion of that sugar; the language gap itself remains.
- Downstream-unwitnessable error arms are an API bug: render's budget/mismatch errors removed in slice 11a (fuel-exhaust now `Ok(acc)` per `int_to_str_from`).
- Variant matches take case arms only, never `_` (slice 11a; attach carries all 12 PTag arms).
- `and`/`or` stay eager, no short-circuit (slice 11a; guards nest instead).
- `forward call` expands to `Ok` + one arm per emits kind, and arm coverage witnesses the expansion (slice 11b; `no test takes on Ok r` on a cannot-succeed arm — the empty-entry arm returns its error directly).

## JEV decision log (all via `jev-1.13.0`, Choice)

a99 forward-looking 0.75; b00 staged 0.97; slice1 compare_trio
0.56; slice3 url_parse 0.40; parse checked_reuse 0.77; slice4
collections_probe (tie 0.49, choice field); slice5 validate_schema
0.65; schema monomorphic_family 1.0; slice6 mapset 0.97, str_keyed
0.70, new_modules 0.56; host bare_externs 0.64; host artifact
delete 0.53; b02 extern admission uses_pin 0.99, host_resolution
declaring_stem 0.84; H1 instant_repr millis_int 0.97; H2
hash_profile brand 0.98, secret_repr brand 1.0; slice8
sort_surface order_value 0.97; slice9 schema_shape
monomorphic_family 0.74.

Next-unblocker recommendation (`jev-1.13.0`, Choice): `generic_values`
(probability 0.64, confidence 0.53; alternatives async/resources 0.21,
Unicode kernel 0.15, clarify priority 0.00). Recommend staged generic
variants/Option first, then first-class outcomes and generic error algebra.
This was a sequencing recommendation; the user subsequently approved focusing
on generic variants. B06 re-probed and removed that blocker only.

B06 case-pattern surface (`jev-1.13.0`, Choice): `explicit`
(probability 1.00, confidence 0.99; contextual inference 0.00,
parse-only deferral 0.00). Require `on Option__Some<T> s` alongside
`Option__Some<T>(value)`, extending the existing explicit monomorphizer.
Outcomes, error algebra, async, and bare variant returns stayed out of B06 scope.

B07 return convention (`jev-1.13.0`, Choice): `ok_value`
(probability 0.99, confidence 0.98; unwrap binder 0.01, raw case 0.00).
Keep `Ok(variant)` and `on Ok r` with the variant at `r.value`.
B07 scope (`jev-1.13.0`, Choice): `calls_only`
(probability 0.75, confidence 0.63; calls + Fn successes 0.25, all channels
including externs 0.00). Source returns, local/foreign calls, and same-nominal
forwarding landed together; Fn successes and extern returns stayed record-only
at B07.

B08 Fn-success scope (`jev-1.13.0`, Choice): `record_and_variant`
(probability 0.99, confidence 0.98; admission-only 0.01, all-data 0.00).
Admit declared record/variant successes through references, signature discovery,
invocation binders, forwarding, and emission together. Preserve B07's
`Ok(value)` / `r.value` convention and all purity, nominal typing, concrete
error-set, containment, cycle, and witness rules. Scalar/other bare successes
and extern returns stayed separate at B08.

B09 scalar-success scope (`jev-1.13.0`, Choice): `primitive_envelope`
(probability 1.00, confidence 1.00; source-only, all-data/extern, raw-scalar
alternatives 0.00). Complete int/str/bool/dec through strict `Ok(value)`
checking, source/call emission, Fn, invocation, and exact-type forwarding;
keep `r.value` binders. Close the legacy unchecked named-field path, rather
than preserve an unsound exception. Brands/Bytes/sequences/Fn bare successes
and extern ABI widening remain out of scope.
