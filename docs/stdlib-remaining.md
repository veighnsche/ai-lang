# stdlib implementation record — overnight goal, 2026-09-19

Goal: implement all planned stdlib from `ASTRA_STDLIB.md` in
PR-sized green commits, JEV for multiple-choice decisions.
Result: everything buildable with the current language is built
(slices 1–6). The remainder is provably blocked on missing
language features, one per row below. No slice was faked around
a blocker: each block below carries its compiler error or design
verdict as evidence.

## Shipped (main, all gates green per commit)

| Slice | Commit | Content |
|-------|--------|---------|
| 1 | `5bad56d` | `std__compare__int/dec/str` close §1.4; 15 rows |
| 2 | `b5780eb` | `std__validate__all` + report records; 15 rows |
| 3 | `b4bfe6e` | `html__url__parse` + `origin_end`; 18 rows |
| 4 | `bb05e2c` | NEW `std/seq`: 7 generic ops; 45 rows |
| 5 | `c3b0075` | `validate__schema__str/int` close §1.5; 14 rows + given rows |
| 6 | `a80e1f5` | NEW `std/map`, `std/set`; 95 rows |

Pre-existing (§1.1–1.3, §1.6, §1.8 text/codecs, §2 elements/render/
assets, quota, schema, ascii) was verified present, not rebuilt.

## Blocked, with evidence

| Roadmap item | Missing feature | Evidence |
|--------------|-----------------|----------|
| §1.7 outcome/option combinators | First-class outcomes + Fn values | b00 staged by JEV (0.97); no `Outcome<T,E>` values exist |
| §1.8 seq map/filter/fold/find/all/any/sort/unique | Function values | Same; callbacks unexpressable |
| §1.8 `Map<K,V>` fully generic | Nested generic Seq elements | `seqElemName` rejects `<>` in elements (`compiler/types.go:164`); shipped str-keyed `Map__Entries<V>` instead |
| §1.8 normalize_nfc, casefold, graphemes | Unicode data kernel | No pinned data, no host path (see host shelf) |
| §1.8 json encode/decode, `schema__migrate` | Functions (+ Schema surface) | Migrate takes `Migration<Old,New>` fn; no Fn values |
| Host shelf (clock/random/hash/secret/log/env) | Shared/callable externs | `calls extern ... from another module: declare your own extern` — externs are module-local by design; a `std/host` module compiles but is uncallable (probed 2026-09-19, artifact deleted per JEV 0.53) |
| §3 HTTP (all) | Async + Resources + Functions | No async surface exists |
| §4 SQL (all) | Async + Resources | Same |
| §5 UI (all) | Functions + Async + Resources | Same |

## Language gaps discovered while slicing (future b-series)

- No `seq[i].field` projection (slice 2; worked around via `_push`).
- Bare `Seq<T>` returns rejected (slice 4; wrapper records required).
- Per-instance arm coverage for generics (slice 4; row cost is real).
- Explicit `<T>` required on recursive generic calls (slice 4).
- `given` rows key on caller test names across module lines (slice 5).
- Downstream sketch goldens embed provider catalogs; regen together (slice 5).
- Row binds must name the declared param (`<V=>`, not `<T=>`, slice 6).

## JEV decision log (all via `jev-1.13.0`, Choice)

a99 forward-looking 0.75; b00 staged 0.97; slice1 compare_trio
0.56; slice3 url_parse 0.40; parse checked_reuse 0.77; slice4
collections_probe (tie 0.49, choice field); slice5 validate_schema
0.65; schema monomorphic_family 1.0; slice6 mapset 0.97, str_keyed
0.70, new_modules 0.56; host bare_externs 0.64; host artifact
delete 0.53.
