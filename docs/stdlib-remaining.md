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

Pre-existing (§1.1–1.3, §1.6, §1.8 text/codecs, §2 elements/render/
assets, quota, schema, ascii) was verified present, not rebuilt.

## Blocked, with evidence

| Roadmap item | Missing feature | Evidence |
|--------------|-----------------|----------|
| §1.7 outcome/option combinators | First-class outcomes + generic error algebra | VERDICT (slice 10): still blocked. Catalogue: "Do not pretend these are available as generic libraries today." Every signature takes/returns `Outcome<T,…>`/`Option<T>`; probe: `unknown type Outcome`. b00 excludes "generic error-set algebra" and "generic outcome values" (`docs/b00-function-values-design.md`); generic variants do not exist (`bad variant decl`), so no `Option<T>` either. No monomorphic fallback matches the catalogue shape. |
| §1.8 seq map/filter/fold/find/all/any | SHIPPED (slice 7): generic workers invoke total callbacks; find reports `sequence.not_found()` (generic variants do not exist — `bad variant decl` — so no `Option<T>`); predicates return per-module `Bool__Value` (std/set precedent) | 62 rows; 107 pass on compile; `TestStdSeqCompiles` gates. |
| §1.8 seq sort/unique | SHIPPED (slice 8): `Seq__Order` value Asc/Desc over per-instance built-in order (insertion sort, stable by construction); unique keeps first occurrences via per-instance `==` | 40 rows; 147 pass on compile; lexicographic orders deferred. |
| §1.8 `Map<K,V>` fully generic | SHIPPED (slice M1): `Map<K,V>` over `Seq<Map__Pair<K,V>>` with catalogue names; str-keyed `Map__Entries<V>` migrated away (no downstream users); errors payloadless (payloads cannot be generic) | 45 rows across `<str,int>` + `<int,str>`; goldens regen via documented flow; `TestStdMapCompiles` gates. |
| §1.8 normalize_nfc, casefold, graphemes | Unicode data kernel | No pinned data, no host path (see host shelf) |
| §1.8 json encode/decode, `schema__migrate` | PARTIAL (slice 9): value layer shipped — `Json__Value` AST, fuel-bounded render machine, escape, 8 scalar codecs, monomorphic `Json__*Schema` family with schema-carried Fn dispatch; bytes-level parse + text drivers remain (slice 11). `schema__migrate`: VERDICT (slice 10) — generic migrate is inexpressible: no `Ok` splat (`Ok field y: got M__B, want int`) and `forward` is refused outside call-outcome arms, so no body returns an arbitrary `New` through invoke dispatch (b00's missing generic error algebra independently blocks `! E`). | 125 rows; 548 pass on compile; `TestStdJsonCompiles` gates. |
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
- No `Outcome<T,E>` / `Option<T>` types; no generic error-set algebra (slice 10; §1.7 blocked).

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
