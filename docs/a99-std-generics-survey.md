# A99 — std-wide generics survey

Status: recorded. Companion: G2 two-pass expansion landed green in
[a98](a98-generic-records-design.md) (Ratio__Value 3→1 pilot + golden/linked/parity).

Question: given all std built so far (29 modules), which modules are
generics candidates beyond the G2 pilot?

## Method

Surveyed every module behind the open gates (division, scala as
representatives of CONSUMES). Criterion: does the module hand-specialize
identical logic per payload type, i.e. N copies of `X__Value` records /
nominal wrappers whose bodies differ only in the payload type?

## Tier 1 — same shape as the G2 pilot (record-param candidates)

| module     | symptom |
|------------|---------|
| validation | `List__Invalids`, `Set__Invalids`, `Table__Invalids`, `TableRow__Invalids` — same invalids-carrier record × 4 collection types. G2-shaped today. |
| division   | `FloorQuotRem`, `TruncQuotRem`, `IntQuot`, `OneQuotRem`, `RatioQuot`, `DecQuot`, `OneDivResult`, `RatioDivResult`, `IntRat` — per-rounding/per-type result bundles. Needs values-inside-records first (function-values / closures-over-values), else it multiplies per (type × rounding) instead of collapsing. |
| outcomes   | 464 `Compare_*` emission wrappers collapsing to one `Compare` — the jackpot, ~460 sites to 1. BLOCKED on bare `-> T` select: same root limit that forced select to stay reverted in G2. |

## Tier 2 — blocked on Values sharing architecture

`json`, `yaml`, `toml`, `msgpack`, `cbor`, `ron`, `csv`, `query`, `ini`,
`envfile`, `pathlib`, `multipart`, `etag` — every codec re-declares its
own `Value` union today (Bool__Value etc. NOT exported: codecs each
mint their own Bool carriers, confirmed by survey). They need a shared
`Value` cluster, which is an **architecture decision (open question)**:
exporting one canonical Value from a shared module vs. each codec
keeping its own. Blocked on the function-values design (codecs need
function values), then on the sharing decision.

## Ruled out — nominal-type unification (html)

`html`'s `provides [...]` wrapper types are the API surface, not
machinery: each wrapper is a distinct named type with no shared body
to parameterize. Unifying them under one `Node<T>` would erase the
nominal API. Ruled out.

## Pipeline

1. function-values design → unblocks division + codecs
2. outcomes design → unblocks the 464→1 Compare collapse (needs bare `-> T` select)
3. Value-cluster sharing architecture decision → unblocks the 13-codec dedup
4. validation: G2-ready now, can be piloted any time
