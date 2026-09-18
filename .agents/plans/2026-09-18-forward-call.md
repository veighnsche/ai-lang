## Goal

Implement the final recommendation from the relay-toll discussion: per-arm `forward call` — `<pat> => forward call f(args)` as a whole arm RHS, elaborated at check time into the full `match call` with per-outcome `forward` arms. All-forward relay matches (starting with the `html__text__escape_from` / `html__attribute__value_from` workers that motivated this) collapse from 3 lines to 1 per arm, with no new trusted machinery.

## Success Criteria

- `=> forward call f(args)` compiles in value and call arms; elaborates to exactly today's handwritten shape; rows, diagnostics, and coverage behave identically to the expanded form.
- Both escape workers migrated (inner blocks one line per arm); `std/` lints clean; goldens regenerated; differential old-vs-new sweep zero-mismatch.
- Prover, coverage, eval, and emit untouched (they see the elaborated shape); one new code (CAN3013); full gates green (`go test ./...`, `modcheck`, `gramcheck`, `tsc`).

## Context And Current Facts

- The toll is per syntactic call site: every call dispatches `Ok` plus every emitted kind (CAN3003/CAN4109). Helper-abstraction re-pays it per use site; the four escape calls must stay four (distinct argument tuples). 97 `=> forward` arms exist across `std/` + sketches, so the feature has consumers beyond the motivator.
- Rejected alternatives, with grounds established in discussion: the fused single-site sketch (no scrutinee binding; positional pattern↔piece zip fails silently with all-`str` pieces; template substitution is macro machinery; inverts the deliberate or-pattern shared-outcome rule C7/CAN4111); `let`-binding (names the callee once but yields 11–13 lines vs 9 — multiplicity is conserved — with bigger machinery; demand-driven follow-up, `a28` holds its slot); user-defined metaprogramming (collides with CAN4107 coverage, golden diagnostics/attribution, CAN3005–3009 syntactic termination, `given` site identity, `a77` revision fingerprints, effects transitivity, and build determinism).
- Elaboration precedents: `elaborateForwards` (`compiler/forward.go`, runs in `buildWorld`, `compiler/check.go:397`) and `elabChain` (`compiler/chain.go`), ordered const→forwards→chains, each once, before proofs, runs, and emit. `markTaken` keys arms by `(node, index)` (`compiler/eval.go:298`), so elaborated arms sharing one source line stay distinct obligations.
- Parse hook already exists: `forward call f(args)` parses today through the arm-RHS path as `Small{Kind: "forward", Str: "call f(args)"}`; `elaborateForwards` currently rejects it as a bad operand (CAN3011). No `parse.go` change needed for match arms. Chain tail/`else` refuse bare `forward` at parse (`compiler/chain.go:148,166`).
- One-code family precedent: CAN3011 covers all `forward` misuses (not one code per refusal). Effects subset messages (`compiler/check.go:1906-1933`, CAN3107 "needs X via Y without declaring it: add it") shape the new error text. Test precedent: `forward_test.go` / `chain_test.go` (`writeLSPDir` + `diagnose` + `hasCode` tables, elaboration positives, emit parity).
- `a89` is owned by a concurrent session (verified free: no `docs/a90*`); its file is untouched outsider state and stays out of this work.

## Constraints And Non-goals

- Slice 0 (design note + approval) precedes ANY compiler change, including preparatory refactors — the call-arm deferral precedent binds here.
- v1 scope: same-file local callees only (self + siblings, resolved via `prog.Fns` + `FnFile`); foreign/uses-pinned/extern callees refused (they need a `given` design that does not exist); whole-arm-RHS only; chain tail/`else` keep refusing (existing parse errors cover).
- Out of v1, explicitly: `let`-binding, call fusion, metaprogramming of any other kind, foreign/`given` support, chain-tail support, and a lint rule suggesting `forward call` (natural follow-up once the shape ships).
- Elaboration only: no changes to prover, coverage, eval, emit, or lint trusted paths. No new CI gates, no rev bumps (bodies never enter fingerprints, `compiler/revision.go:18-20`).
- Publication split: slice 0 ships the note only (no code); slice 1 ships implement + probes + migration as one atomic commit (feature and goldens revert together).

## Key Decisions

1. **Per-arm `forward call` is the final recommendation** — not fusion, not `let`, not macros. Smallest mechanism (desugar of a fixed shape), shortest outcome on the motivator (19-line body → 9), generalizes to differing callees and multi-arg differences for free, keeps every pair adjacent per the language's explicitness bias.
2. **New pass `expandForwardCall` in `buildWorld`, BEFORE `elaborateForwards`** (order const→forwardcall→forwards→chains). Required by construction, not optional: its output contains `forward r/e` arms that the existing pass must then expand — which also yields the `Ok`-shape check (`forwardOk`), relay certificates (CAN4108), and normal arm coverage for free. New file `compiler/forwardcall.go`, mirroring `forward.go`/`chain.go`.
3. **Allowed in value AND call arms** (mechanism is position-independent; the auth-login dependent shape benefits too). Refused in chain tail/`else` (existing parse errors cover; message polish only if trivial). The value-arm asymmetry with plain `forward` (CAN3011 refuses it there) is principled: `forward` relays an outcome, `forward call` creates the dispatch.
4. **One new code CAN3013 `CodeBadForwardCall`** for all static refusals (non-local callee, undeclared callee emits, malformed site) — CAN3011-family precedent wins over strict one-code-per-refusal. Verified free; CAN3417/3418 are taken by the lint slices.
5. **Emits inclusion via `prog.EmitsOf`**: callee emits ⊆ caller `emits`, else CAN3013 with a CAN3107-shaped message ("`f` forwards `g`'s outcome but does not declare …: add it to emits").
6. **Migration breadth: the two motivating workers only** (`escape_from`, `value_from`) — the `a86` approve-only precedent. The other ~90 `forward` arms stay valid code; a lint suggestion rule is the follow-up vehicle, not this slice.
7. **No lint changes.** Existing rules treat the shape sanely without edits: C14/C15 need match nodes (`forward call` sites are Smalls until elaboration, and lint runs on parse); or-fold/same-outcome text identity stays sound for identical call texts; the chain rule correctly ignores non-rungs. Pinned by probes, not by code.

## Recommended Approach

Two slices. Slice 0 writes `docs/a90-forward-call.md` from these decisions (semantics, checker obligations, the fusion/`let`/macro rejection record, v1 scope, reopening conditions) and stops for approval with zero code touched. Slice 1 implements as one unit — registry + explain, elaborator + wire-up, probes, migration + goldens, docs — following the `a86` one-slice shape, with emit parity (elaborated vs handwritten byte-identical TS) as the strong check.

## Work Plan

### Slice 0 — design note (no code)

0. Draft `docs/a90-forward-call.md`: surface syntax, elaboration rule, checker obligations (emits inclusion CAN3013, locals-only, whole-RHS, chain-position refusals), coverage/termination arguments (free via ordering, `(node, index)` identity), rejected alternatives with grounds (fusion, `let`, macros per Context), v1 non-goals with reopening conditions (foreign/`given`, lint suggestion, chain tail). Approval gate; no compiler edits, no refactors.

### Slice 1 — implement, probe, migrate

1. **Registry.** `CodeBadForwardCall = "CAN3013"` const + `allCodes` sequence slot (call family) + `explain.go` entry (rule/violation/fix). Pinned by `TestExplainComplete`, `TestCodesUnique`, `TestAllCodesSequenced`.
2. **Elaborator.** New `compiler/forwardcall.go`: branch on `forward` Smalls whose `Str` has the `call ` prefix; resolve callee to a same-file `FnDecl` (`prog.Fns` + `FnFile`, else CAN3013); check callee emits ⊆ caller emits (`prog.EmitsOf`, else CAN3013); build the `match call` node (`Ok` + per-emit forward arms) at the site line, replacing the arm RHS in place. Wire into `buildWorld` before `elaborateForwards` (`compiler/check.go:~393`). Must run exactly once (same doubling rationale as the existing passes).
3. **Probes.** New `compiler/forwardcall_test.go`, mirroring `forward_test.go`: elaboration positives (value arm, call arm with bound payload arg, self-recursion under a `<=` guard); rejection table (foreign/uses-pinned/extern callee, undeclared emits, malformed site, chain tail/`else`, non-arm position); `Ok`-shape mismatch surfacing via ordering; row law (untaken elaborated arm fails CAN4107); emit parity (handwritten vs `forward call` spellings compile to byte-identical TS); no-op on plain `forward v`.
4. **Migration.** `escape_from` + `value_from` inner blocks → one `forward call` line per arm (byte-exact `html.can` edit, NUL count asserted 70 before/after); regen `html.ts` + `errors.json` via `canlc --out` (with ascii+scalars siblings); TS diff confined to the 2 functions, `errors.json` predicted identical; differential old-vs-new sweep reusing the `/tmp/esctest/diff_std.mjs` inputs (escape + value cases).
5. **Docs.** `REQUIREMENTS.md` amendment per the note (the R6 match/call-arm grammar line gains the `forward call` RHS shape); `docs/can-idioms.md` Tier 1 amendment ("Calls dispatch, never sequence" gains the second legal call shape — law text, no C-number since no lint rule ships in v1).
6. **Gates.** `go test ./...`, `go run ./tools/modcheck`, `go run ./tools/gramcheck`, `tsc -p tsconfig.json` in `tscheck/`; `canlc lint std/` exits 0; `sketches/lint-errors/` exits 1 with exactly the 9 pinned findings (no new-rule interference).

## Validation Plan

- Focused probes first (unit 3 green before migration), then the full suite; repo gates at slice end per unit 6.
- Strong checks: emit-parity test (byte-identical TS for both spellings over identical rows); migration differential sweep zero-mismatch; TS-diff confinement; NUL-count equality.
- Highest-risk validation: elaboration-order interactions (output consumed by the forwards pass; chain-`else` reparse never seeing the shape; LSP prove-before-check path) — covered by the chain-position probes plus the existing chain/lsp suites staying green with zero edits.

## Risks / Rollback

- **Order bugs** (output missed by a later pass, double diagnostics on re-run): the exactly-once wiring mirrors the three existing passes; probes pin end-to-end expansion through `diagnose`.
- **Span quality on shared lines**: elaborated arms share the site line; `markTaken` identity is structural, and LSP anchors fall back honestly (whole-line) if a token span ever misses.
- **Scope creep** (foreign/`given`, lint suggestion, `let`): explicitly out with reopening conditions in the note; any request for them during review becomes a follow-up, not slice-1 growth.
- **`a89` concurrent session**: only shared surface is `docs/` numbering (`a90` verified free); no code overlap. Its file stays untouched and out of any commit.
- Rollback is a single revert: feature, probes, migration, and goldens ship atomically.

## Open Questions

1. Confirm locals-only v1 (recommend yes; foreign needs a `given` design that does not exist).
2. Confirm `forward call` allowed in call-match arms as well as value arms (recommend yes — mechanism is position-independent, and payload-dependent relays benefit).
3. Confirm migration breadth: motivating-2-only (recommend; `a86` precedent) vs all mechanically-convertible sites.

## Sources

None — every material claim traces to workspace evidence inspected during planning (cited inline by file and line).
