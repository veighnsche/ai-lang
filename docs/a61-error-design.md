# a61: Rust-grade errors for agents (proposal)

Status: proposal (pre-decision). No rule, no code, no golden.
Decision and sequencing are open; the ordering below is the
proposal. Nothing here redesigns R5: every item is additive.

Goal: error handling as good as Rust's, optimized for AI
coding agents instead of humans. The language already beats
Rust on three dimensions (no escape hatches, declared `emits`
sets, machine-readable catalog); this proposal closes the
three gaps where Rust is ahead, in agent-ROI order.

## What stays (not negotiable in this proposal)

- Errors as values, exhaustive match, no catch-all. There is
  no `unwrap`, no panic path for typed errors, and no
  `anyhow`-style untyped escape hatch. An agent cannot write
  the lazy path; that is the point.
- `emits [...]` stays the declared per-function error set
  (checked-error precision Rust never got). `CAN4001`
  stands.
- Deny-by-default diagnostics. A warning an agent can ignore
  is noise: every diagnostic is either an error or deleted.
  No Clippy-style configurable lint levels.
- `?` is NOT borrowed. `?` is early return, and hidden
  control flow is a non-goal. Item 2 takes its terseness
  through an explicit form instead.

## Item 1: structured diagnostic fixes + explain (first)

Problem: canlc diagnostics point but never suggest. There is
no fix/suggestion path in the LSP or CLI; rustc ships
"did you mean", "consider adding", and rustfix. Agents act
on structured edits far more reliably than on prose, so
this is the highest-ROI error investment and needs no
language change.

Proposal:

- High-frequency codes carry machine-actionable payloads:
  `expected` / `found`, the candidate token where one
  exists, and a suggested edit expressed as (span,
  replacement text). First codes: `CAN4001` (raise
  outside emits: suggest the declared kind or the emits
  entry), `CAN3204` (bare kind: suggest the complete
  construction), `CAN3105` (dangling test: suggest the
  missing script row shape), `CAN4101` (missing arm:
  list the missing patterns — rustc already does this
  for non-exhaustive matches).
- Payloads ride the existing JSON diagnostic lines, so
  CLI and editor share them like every other squiggle.
- `canlc explain CANxxxx` prints the per-code doc: the
  rule, a minimal violation, the legal fix. Codes stay
  frozen; explain text is documentation, versioned with
  the compiler.
- Goldens freeze payload shapes the way `TestGoldenJSONDiags`
  freezes codes today.

## Item 2: explicit forward arms (second)

Problem: every error mapping is a full arm, and relay arms
(`on auth.login_failed _ => auth.login_failed(user_id =
user.id)`) are boilerplate the agent writes and the
reviewer re-verifies field by field. Rust's `?` solves the
terseness but violates expression-position semantics.

Proposal (name TBD, one canonical form):

```
on auth.login_failed e => forward e
```

Rules, all checkable:

- The arm still exists. Coverage law and test-per-arm
  (`CAN4107`) apply unchanged: `forward` is a body shape,
  not an arm exemption.
- The compiler checks the forwarded kind is in the
  function's `emits` and the payload is complete — the
  same checks a hand-written relay passes, minus the
  handwriting.
- Binder form follows existing arm shapes (`e` binds the
  payload record; `_` stays legal only where nothing is
  forwarded). No new binding semantics.
- Needs grammar + eval + emit + goldens; TS emit is the
  same union construction a manual relay emits today.

## Item 3: cause chains + arm traces (third)

Problem: mapping drops provenance. When `auth__verify`
maps `auth.mismatch` to `auth.login_failed`, the original
value dies at `_`. Rust has `source()` chains and
`.context()`; canlc has neither. For an agent debugging a
failing test, "which hops did this error pass through" is
the missing datum.

Proposal, two parts:

- `cause` payload convention: a mapped error may carry the
  consumed error value as a declared field, constructed
  explicitly like any other payload (complete-construction
  rule unchanged — `CAN3204` still rejects bare kinds).
  Convention first, compiler support only if a consumer
  needs it.
- Structured `CAN4200`: today the most-read diagnostic in
  the language renders as one flat string (`test X fails:
  ...`). It should carry a field-level diff (which field
  differed, expected vs got — tag-alone comparison already
  fails, so the diff must reach the payload) plus the
  arm-firing trace for the failed test (which exchanges
  scripted, which arms fired, in order). Trace format is
  shared with any future journey analysis (a56).

## Item 4: structured loud-fault payloads (fourth)

Problem: the typed/loud split (fault-contracts.md) is
principled, but loud faults surface as message strings
(`str index out of range`). An agent debugs faster from
`{op: index, arg: 5, len: 3}`.

Proposal: evaluator faults carry structured payloads
(operation, arguments, domain bound) in test output and
JSON diagnostics; plain message strings survive only at
the TS emit boundary, where the host contract owns them.
No language change; evaluator + diagnostic rendering only.

## Item 5: catalog mapping edges (last)

Problem: `handled_by[].arm` in errors.json is a source row
(`on auth.mismatch _ => ...`), not data. Every consumer
that wants the A→B mapping graph re-parses prose.

Proposal: the catalog gains structured mapping edges
(from-kind, to-kind, fn, test rows exercising the edge)
plus per-kind metadata (one doc line; retryable or not).
This is what makes "every path from `db.down` to a
user-visible error" answerable — the journey analysis a56
will eventually want. Catalog stays byte-identical and
sorted; new keys only, no shape break beyond addition.

## Non-goals

- Any untyped error escape (`anyhow`, `Box<dyn Error>`,
  stringly failures). Agents abuse escape hatches first.
- Panics or exceptions as control flow. Loud faults stay
  loud and uncatchable; typed errors stay values.
- Natural-language step semantics or prose-driven fixes.
  Fixes are spans and replacements, never sentences the
  agent interprets.
- Replacing `tests` + `given`. Error tooling layers above
  the kernel; the kernel does not shrink.

## Ordering and gates

1. Item 1 (diagnostic payloads + explain). No grammar
   change; goldens + `go test ./...` + LSP parity.
2. Item 2 (forward arms). Grammar, eval, emit, goldens;
   auth-login relay arms are the migration proof.
3. Item 3 (cause + CAN4200 structure). Evaluator +
   diagnostics; failing-test fixtures pin the diff.
4. Item 4 (fault payloads). Evaluator only.
5. Item 5 (catalog edges). Additive keys; regen goldens.

Each item graduates separately: a verdict-style review
(counterexamples welcome) before any code, in the order
above. Item 1 may proceed alone; items 2–5 wait their
turn.
