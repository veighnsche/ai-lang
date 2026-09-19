## Goal

Produce the async language design for can-lang: surface syntax with look-and-feel samples, static discipline, test-scripting model, host boundary, and gate plan — as a b-series design doc — without any async/await form and without implementing anything yet.

## Success Criteria

- A b-series design doc (`docs/b05-async-design.md` if the number is still free, else the next free b number) exists in the established Q&A-verdict format with checked falsification notes.
- The doc resolves all five design areas below with explicit verdicts; every `ASTRA_AUDIT.md` §4 proof obligation (join exactly once, exhaustive outcome product, partitioned authority, explicit liveness bound) is addressed or explicitly deferred with a reason.
- async/await and first-class pending values are explicitly rejected with grounding in user constraint and workspace proof discipline, not taste.
- No compiler, sketch, stdlib, grammar, or tool file is modified by this work.

## Context And Current Facts

- The user asked to start the async design now (timing: design ahead of implementation) and rules out JS/TS-style async/await entirely (user constraint).
- `docs/a/ASTRA_AUDIT.md` §4 prescribes the direction: structured async with exhaustive joins and explicit authority partitioning, illustrated as a `match join` with per-branch effects and product arms over branch outcomes. It requires: every spawned task joined exactly once; exhaustive combined outcomes; disjoint branch effects unless a protocol proves safety; terminating can-tasks first; host async only behind an explicit trusted completion or timeout contract; liveness assumptions never smuggled into proofs.
- `docs/a/a05-expressiveness.md` ("Explicitly later") defers async ("distant. Nothing in v0 wants it yet") and rules that no build starts without its own spec doc in the series — so the design doc is the mandatory first artifact, and the audit's sketch (a rating illustration, not a spec) must not be mistaken for the design.
- Effects discipline (`docs/a/a09-effects.md`): effects only via `call` plus exhaustive outcome handling, declared where readers see them; store ops are a deterministic third callee kind needing no `given`. Async tasks become a new concurrency surface inside this discipline, not beside it.
- Test evidence law (a03/a81/a91, as exercised through b00): coverage is observed execution; scripts are counterfactual observations, never proof facts; omission claims non-reach. Async completion in tests must therefore be scripted observation, never evidence of liveness.
- Host boundary (`docs/b02-shared-externs-design.md`): shared externs with `uses` pins and declaring-stem emit; host implementations themselves still unwritten. Host-async completion/timeout contracts will extend this boundary.
- b00 boundary (`docs/b00-function-values-design.md` §9/Unresolved): no effectful or host callbacks, no Fn-bearing invocation inputs. Async tasks must be joined values of a join construct, never callbacks — the design must not reopen the callback surface.
- Demand context (`docs/stdlib-remaining.md`, `docs/a/a13-stdlib.md`, `docs/ASTRA_STDLIB.md`): HTTP/SQL/UI need Async (+ Resources); HTML and strings/brands ship before async. This design unblocks planning those layers; it does not reorder implementation.

## Constraints And Non-goals

- No await points, no `async` function modifier, no function coloring in any form (user constraint).
- No first-class pending/deferred values: nothing inhabits an incomplete computation as a passable, storable value.
- No bare spawn/detach: every task is joined exactly once in the surface (audit obligation).
- No implementation: no compiler edits, prototypes, goldens, or grammar changes. A later implementation plan follows approval of the design doc.
- Resources (scopes/lifetimes/cancellation tokens as a surface) are a separate design; this plan covers async only and records the seam. Assumption, stated for the approval checkpoint.
- No restatement of the audit sketch as the design: every sketch element must be re-derived or explicitly adopted with grounding.

## Key Decisions

1. **Design now, build later.** User timing plus the a05 spec-doc rule: the doc is the deliverable; implementation waits for the HTTP/SQL/UI work that needs it.
2. **Concurrency surface is a match variant, not a modifier or a value.** A join is written and dispatched like the existing outcome matches (per the audit's `match join` direction), so branch outcomes compose through the existing exhaustiveness and coverage machinery instead of inventing a parallel one.
3. **No coloring, no suspension points.** Await points are invisible to decision tables: a proof cannot enumerate what it cannot see, and a rejection channel outside declared outcomes violates the outcome discipline. Rejected on proof-discipline grounds in addition to the user's constraint.
4. **Tasks are joined, never called back.** b00 explicitly excluded effectful/host callbacks; async tasks complete into join arms, and no design element may reintroduce callback-shaped completion.
5. **Completion is scripted observation.** Test rows script task completion like `given`/`exchange` rows script foreign outcomes; a passing row means "given this completion, this path," never "the task terminates."
6. **Terminating can-tasks first; host completion by explicit contract.** Matches the audit's staging: the core design proves what it can (join structure, outcome products, authority partition) and puts liveness/timeout trust in named contracts at the shared-extern boundary, never in types.

## Recommended Approach

Write the design doc in five slices, each ending in verdicts with checked falsification against the cited sources above:

- **D1 — Surface and look-and-feel.** Join-match syntax, branch declarations with effect/authority annotations, product-arm patterns, and the zero/one-branch edge shapes. Samples for: two-branch join, single-branch timeout wrapper, nested joins, and a join inside an ordinary outcome arm. Rejections: await points, deferred values, bare spawn.
- **D2 — Static discipline.** Spawn/join-once rule and its checker formulation; outcome-product construction and exhaustiveness checking (including the audit's nine-cases-from-two-operations blowup and what abstraction, if any, is justified); authority/effect partitioning rule with the disjointness check; termination stance for can-tasks.
- **D3 — Test scripting.** Scripted completion rows: shape, placement (`given`-style vs new block), partial-script omission semantics, and what counts as coverage for join arms. Determinism argument for joined execution order vs `order [...]`.
- **D4 — Host boundary and lowering sketch.** Completion/timeout contract shape over shared externs; what the TypeScript emit produces for a join (existing emit target per the tsc gate, `docs/a/a14-tsc.md`); why no host promise type leaks into Can semantics; state-interference rules for branches touching cells.
- **D5 — Gate plan and acceptance.** b00-§8-style gate table (parse/types, static admission, value semantics, rows/coverage, golden emission, linked, parity, modcheck, gramcheck, LSP, lint) scoped to what async adds, plus the open-questions-and-risks section.

## Work Plan

1. D1 surface draft with samples; verify each sample parses mentally against the existing match grammar and name every grammar touch it would eventually need (no edits).
2. D2 static rules; falsify each against a05/a08-termination/a09/a11-recursion and the audit obligations; record what the checker can vs cannot see.
3. D3 scripting model; falsify against a03/a81/a91 (scripts-are-observations, no proof smuggling, omission semantics).
4. D4 host contracts + lowering sketch; falsify against b02 (pin/stem discipline) and a14 (canlc owns behavior, tsc owns shape).
5. D5 gate table + risks; assemble the b-series doc; self-review pass for verdict completeness (every §4 obligation addressed or explicitly deferred).

Dependencies are linear (D2 needs D1's surface; D3 needs D2's outcome product; D4 needs D1+D2; D5 needs all). No code surfaces are touched.

## Validation Plan

- Each slice ends with verdicts plus the falsification source that was checked (doc path + section); a verdict without a checked source is rework, not done.
- Audit-trace check: map every `ASTRA_AUDIT.md` §4 obligation sentence to a verdict or an explicit deferral; no obligation may be silently dropped.
- User-constraint check: the words "await", "async fn", "colored", "suspend" appear in the doc only inside rejections.
- b00-boundary check: no task-completion element is expressible as an effectful/host callback; the Fn-bearing input/success ban still holds for every new position.
- Manual review gate: present the assembled doc for approval; implementation planning is a separate follow-up and must not start here.
- Highest-risk validation step: the D2 outcome-product blowup — if no justified abstraction contains the case explosion, the doc must say so plainly (per the audit's "why humans hate it") rather than hand-wave.

## Risks / Rollback

- Risk: the design drifts into implementation (prototypes, grammar edits). Mitigation: the plan forbids tracked-file edits outside the doc; rollback is deleting the doc draft.
- Risk: re-deriving the audit sketch without re-examination. Mitigation: D1 must adopt-or-replace each sketch element with fresh grounding.
- Risk: scope bleed into Resources (cancellation/deadlines as values). Mitigation: recorded seam; deadlines appear only as host-contract parameters unless a verdict with grounding says otherwise.
- No compatibility impact: design-only work changes no shipped behavior.

## Open Questions

1. Should cancellation/deadline surfaces be designed jointly with async, or stay split as this plan assumes (async now, Resources later)? Recommendation: stay split; joint design roughly doubles D2/D4. Awaiting the approval checkpoint, not blocking drafting.
