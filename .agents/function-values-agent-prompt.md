# Agent prompt — function-values design for Can-lang

Paste everything below the line to the external agent.

---

You are designing a language feature for **Can-lang**, a small typed
language with a Go compiler. Repo (clone it, read it, ground your
design in the actual code):

**https://github.com/veighnsche/can-lang**

## What Can-lang is

- Source language: `.can` modules under `std/` (29 stdlib modules:
  codecs like json/yaml/toml/msgpack/cbor/csv, plus division, scala,
  validation, outcomes, html, etc.).
- Compiler: Go, under `compiler/` (`parse.go`, `expand.go`, `eval.go`,
  ...). Verified with `go test -count=1 ./...` plus gates:
  `go run ./tools/modcheck`, `go run ./tools/gramcheck`, LSP and lint
  probes. The suite is currently green — your design must keep every
  gate passing.
- Design docs live in `docs/aNN-*.md`. Read these before designing:
  - `docs/a93-generics-design.md` — generics design (Seq precedent)
  - `docs/a98-generic-records-design.md` — G2 generic records:
    two-pass monomorphic expansion with identity splice, per-instance
    type stamps, Q1a ratio pilot, Q2a reject-unmentioned verdicts
  - `docs/a99-std-generics-survey.md` — std-wide generics survey; this
    is the requirements source for your design (Tier 1 + Tier 2)
  - `docs/ASTRA_STDLIB.md` — stdlib roadmap (`std__outcome__and_then`
    and friends show the `Fn<T,U,E>` shape the stdlib expects)

## The problem (from A99)

Two clusters of stdlib code cannot be deduplicated today because
Can-lang has no function values:

1. **division** — per-rounding/per-type result bundles (`FloorQuotRem`,
   `TruncQuotRem`, `IntQuot`, `RatioQuot`, `DecQuot`, ...). Needs
   values-inside-records / closures-over-values: without it, adding
   generics here multiplies per (type × rounding) instead of collapsing.
2. **13 codecs** (json, yaml, toml, msgpack, cbor, ron, csv, query, ini,
   envfile, pathlib, multipart, etag) — each re-declares its own `Value`
   union today (confirmed: `Bool__Value` etc. are NOT shared). They need
   function values first, then a Value-cluster sharing architecture
   decision (your design should frame that decision, not necessarily
   settle it).

Out of scope: unifying html's nominal wrapper types (A99 ruled this
out — the wrappers ARE the API), and the outcomes 464→1 Compare
collapse (that is the follow-up design; it additionally needs bare
`-> T` select, currently blocked — note the interface your design
leaves for it).

## Known constraints (do not regress these)

- G2 expansion is two-pass monomorphic with per-instance type stamps;
  extend that machinery, don't fork it.
- Fixed during G2, keep fixed: type-stamp `Domain__Name` exemption via
  `GenericBase`, precise chain-else gate, grammar without lookbehind,
  no LSP/lint name collisions.
- `select` stays reverted on the bare-return emit limit — your design
  must not depend on bare `-> T` select.

## Deliverable

Write `docs/b00-function-values-design.md` in the repo, following the
style of a98/a99 (status header, numbered design questions with
verdicts, e.g. Q1a-style). It must contain:

1. Syntax + typing rules for function values (and closures if you
   propose them), grounded in `compiler/parse.go` / `expand.go` as
   they exist today — cite the functions/areas you would extend.
2. Expansion strategy: how function values interact with two-pass
   monomorphic expansion and per-instance stamps.
3. Worked pilot: pick ONE division bundle family and show the
   before/after (N declarations → 1), with the exact `.can` code.
4. The Value-cluster sharing decision, framed as explicit options with
   a recommendation.
5. Gate plan: what new tests/probes each gate needs (unit, golden,
   linked, parity, modcheck, gramcheck, LSP, lint) to keep the suite
   green.
6. Open questions + risks section, including the interface left for the
   follow-up outcomes design.

Do NOT implement. Design doc only. If the repo contradicts anything in
this prompt, the repo wins — say so in the doc.

## Evidence contract (how your doc will be judged)

- Every claim about current compiler behavior must cite the inspected
  body as `file:line` (e.g. `compiler/expand.go:123`). Search/grep
  output only locates candidates — open the function before citing it.
- Every design decision must name the concrete falsification you
  checked: what existing test, gate, or std module would break if you
  were wrong, and why it survives. Unchecked decisions go in Open
  questions, not in the design.
- End the doc with an explicit `Unresolved` list: anything you could
  not verify from the repo, stated as exact gaps, never silently
  dropped or reworded into the design.
- State omitted scope: which std modules, gates, or edge cases you did
  not examine. A bounded sample is never presented as exhaustive.
