# Decider prompt: how do can-lang compile-time tests treat Jev judgments?

Paste everything below the line to the decider agent. It assumes a checkout
of can-lang at the working root plus network access to docs.typesafe.ai.

------------------------------------------------------------------------

You are the decider for the highest-risk open question in making Jev a
first-class citizen of can-lang: what compile-time decision tables mean when
the branch condition is a probabilistic model judgment. Read, deliberate,
and return a single ratifiable decision. Take your time; this decision
outlives the implementation.

## Settled: do not reopen

1. Jev gets dedicated .can judgment syntax (new expressions with confidence
   handling elaborated by the compiler), not extern-only wrappers.
2. TypeScript emit stays dependency-free: judgment calls emit as plain
   platform fetch against the TypeSafe HTTP API (POST
   api.typesafe.ai/v1/systemone, Bearer key), no SDK import. Same
   zero-dependency invariant as the $canDec inline helpers (R11).
4. Runtime authentication: the TypeSafe Bearer key comes from env vars
   (TYPESAFE_API_KEY per the quickstart convention). Emitted code reads
   it from the environment; it never appears in .can source.
3. can-lang's v0.1 freeze plus tagged-amendment discipline stands
   (REQUIREMENTS.md): no silent rule rewrites, one rule per CAN code,
   every behavior change gated by golden/diagnosis suites.

If you believe a settled item blocks a sound decision, say so only with a
concrete demonstration of the blockage, then still decide under the
settled constraints.

## Context you must read before opining

Repo (paths from root):

- REQUIREMENTS.md: Goal, R5 (errors), R6 (functions), R7 (compile-time
  tests), R8 (call-site given), R10 (tooling/one gate), R11 (TS target)
- docs/a03-branch-coverage.md (test-per-arm law, CAN4107)
- docs/a12-contracts.md (producer-owned emits, exchange rows)
- docs/a81-proof-obligations.md (scripts are never proof facts; proof and
  scripted evidence stay separate)
- docs/a87-acceptance-authority.md (pinned vs proposed rows, CAN6017)
- docs/a14-tsc.md (canlc is the sole verifier)
- sketches/auth-login/auth.can (exchange-table shape in practice)

Live TypeSafe docs (source of truth; append .md to page paths):

- Index: https://docs.typesafe.ai/llms.txt
- Primitives: primitives.md, primitives/choice.md, primitives/noul.md,
  primitives/score.md (Choice/Score carry probabilities plus confidence;
  Noul carries P(yes) and no confidence)
- confidence.md (confidence is derived from the distribution shape;
  three-path high/medium/low pattern; thresholds scale with risk)
- concepts/how-to-build-with-system-one.md, patterns/fan-out.md,
  patterns/confidence-routing.md
- api.md (the HTTP contract the fetch emit targets: single POST endpoint,
  Bearer auth, JSON state/questions in, typed answers out)

Reference: Probably-lang (https://probably-lang.southpolesteve.workers.dev/):
feels with confidence gates, otherwise-maybe, highest-probability match,
chaos sampling, zero-model-call replay, run JSON recordings. Steal what
fits; name concretely what violates can-lang rules and why.

## The question

How do hermetic compile-time decision tables (R7) and call-site scripts
(R8) treat probabilistic Jev judgments? Steelman each candidate, then pick
one or synthesize a hybrid:

A. Scripted distributions in given. Each test scripts label plus
   probabilities in the call-site table; probabilities are just data;
   hermetic; fits R7/R8 as written.
B. Record/replay cassettes. Record live Jev answers once, check the
   cassette in, replay in builds. New artifact; weaker hermeticity.
C. Live calls at build. The compiler calls Jev during evaluation.
   Breaks the hermetic sandbox; builds need credentials and go
   nondeterministic.
D. A hybrid you define. Examples only, not a menu: scripted labels with
   confidence bands instead of exact probabilities; pinned acceptance rows
   verified live in a separate non-blocking gate in the CAN6017 spirit;
   distribution-shape assertions; fan-out questions answered together.

## What good means (criteria, roughly weighted)

1. Hermetic, deterministic, credential-free builds (R7).
2. Proof/evidence separation (a81): scripted outcomes must not become
   proof facts by accident.
3. Test-per-arm law (a03) over probabilistic branches: what coverage means
   when branch selection is probabilistic, including any sampling story.
4. Amendment discipline: deltas expressible as tagged rules with new CAN
   codes; prod-strip (R7) still holds.
5. Agent-authored UX: explicit, verbose, single canonical form; nothing
   the author would have to guess is inferred.
6. Calibration honesty: probabilities are model estimates, not truth;
   builds must not bless thresholds as true; account for drift between
   scripted/recorded and live behavior.
7. Cost, latency, credentials: who pays and when; offline builds.
8. Migration: path from extern-stubbed judgments to first-class syntax.

## Deliverable

1. Decision in one paragraph, your confidence in it, and what evidence
   would change your mind.
2. Rationale per criterion above, each with file/section receipts.
3. Rejected alternatives, each with the concrete failure mode that kills
   it (a specific .can snippet that misbehaves, not a slogan).
4. Ratification-ready amendment sketch: rule deltas against R7/R8 with
   version tags, new CAN codes, one normative .can example (types, fn,
   tests, given, judgment arms including the low-confidence path), the
   expected dependency-free TS emit sketch (fetch plus inline response
   shaping, no imports), errors.json impact.
5. Edge cases adjudicated, each with a ruling: ties; Noul near 0.5; arms
   reachable only below a confidence floor and CAN4107; leftover versus
   unreached scripts for judgments; whether given rows script raw
   probabilities, confidence, or both; whether chaos-style sampling exists
   and where; prod-strip of judgment tables; whether prompt/instruction
   wording is rev-hashed code (reword without bump: error?); multi-question
   fan-out in one model call and how its answers are scripted; endpoint
   and model pinning ("jev-latest"
   floats: allowed, or must the model be pinned like a uses rev per R4?);
   transport failures (non-2xx, network error, malformed answers) mapped
   to which declared emits.
6. Verification plan: which existing golden/diagnosis suites change, which
   new suites are required, modcheck/gramcheck impact.
7. Open questions ranked by blast radius.

## Process

- Read before opining. Cite paths and sections for every load-bearing
  claim. Quote the live docs for every TypeSafe behavior you rely on.
- Do not implement compiler changes. One small throwaway probe is allowed
  only to settle a factual question.
- Return the deliverable as a markdown document plus a 10-line summary.
