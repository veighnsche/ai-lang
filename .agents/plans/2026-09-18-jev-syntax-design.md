# Design prompt: first-class Jev judgment syntax for can-lang

Paste everything below the line to the designer chatbot. It assumes web
access to GitHub and docs.typesafe.ai (a local can-lang checkout works too;
local paths are given beside blob URLs).

------------------------------------------------------------------------

You are designing the surface syntax that makes Jev judgments feel like
they were always part of can-lang. This is a syntax and look-and-feel
task: produce multiple distinct options, not one verdict.

## The problem

JEV-1.0 (docs/a94-jev-judgment-evidence-decision.md) settled the semantics
of judgment evidence, and that settlement stands. But its section 5 syntax
is rejected: it reads like a fetch wrapper. Per-descriptor service, model,
timeout, and retry fields plus per-site jev.unavailable and
jev.invalid_response arms put transport plumbing at every call site.
Judgments must feel like native control flow. The network reality stays,
but it must be handled inherently, not re-spelled everywhere.

## What stays settled

Preserve these, or challenge one only with a concrete demonstration that
it blocks a sound design (then still deliver options under the settled
constraints):

1. JEV-1.0 evidence semantics: complete-observation given scripts,
   deterministic compiler-owned routing, offline credential-free builds,
   test-per-arm coverage (CAN4107), scripts never proof premises, no
   sampling in the judgment core.
2. Dependency-free TypeScript emit via plain platform fetch against the
   TypeSafe HTTP API; no SDK import. The Bearer key comes from the
   TYPESAFE_API_KEY environment variable at runtime, never from source.
3. can-lang amendment discipline (REQUIREMENTS.md): tagged rules, one rule
   per CAN code, behavior changes gated by golden/diagnosis suites.

## Read first, then design

can-lang on GitHub (veighnsche/can-lang, branch main):

- REQUIREMENTS.md (Goal, R1-R11, especially R5-R8 exhaustiveness and
  no-catch-all): https://github.com/veighnsche/can-lang/blob/main/REQUIREMENTS.md
- JEV-1.0, the semantics you preserve and the syntax you replace:
  https://github.com/veighnsche/can-lang/blob/main/docs/a94-jev-judgment-evidence-decision.md
- Normative .can style: sketches/auth-login/auth.can and docs/can-idioms.md
  (local: sketches/auth-login/auth.can, docs/can-idioms.md)
- Diagnostic registry shape: compiler/code.go

Live TypeSafe docs (source of truth; append .md to page paths). Squeeze
every capability below into a syntactic home; nothing in these pages may
lack an expression in your designs without an explicit reason:

- Index: https://docs.typesafe.ai/llms.txt
- primitives.md, primitives/choice.md, primitives/noul.md,
  primitives/score.md, primitives/advanced.md (structured instructions
  and criteria, Noul true/false criteria, Choice criteria maps, Score
  ordered levels, multi-question fan-out in one request)
- confidence.md (reported confidence on Choice/Score, none on Noul;
  thresholds as application policy), concepts/state.md,
  patterns/confidence-routing.md, concepts/how-to-build-with-system-one.md
- api.md (the HTTP contract the fetch emit targets)

Probably-lang as feel inspiration (not semantics to copy blindly):
https://probably-lang.southpolesteve.workers.dev/
Study feels with confidence gates, otherwise-maybe, highest-probability
match, and zero-model-call replay. Steal the lightness; respect that
can-lang requires exhaustive expressions and refuses optional fallthrough.

## The owner's leans

Design toward these. Alternatives are welcome if argued, but each option
must say where it follows or departs from them:

1. Model choice, timeout, and retries live in a separate config file, not
   per site and not per descriptor. Transport identity is configured once.
2. A failed network call collapses to a wildcard arm, shaped like
   _ => "failed": one small failure path, not exhaustive per-error arms
   at every judgment.
3. Noul, Choice, and Score each get a clear, distinct, minimal look. An
   agent reader should tell which primitive a site uses at a glance, and
   the simplest case of each should be a few lines.

## Tensions each option must resolve

a. A wildcard failure arm against R5/R6: no catch-all, exhaustive match.
   Reconcile this or propose a principled, narrowly scoped exception with
   draft amendment text. "Wildcards are fine now" is not a resolution.
b. Config-file transport policy against R4 revision identity (JEV-1.0
   section 4.5: model, policy, and prompt wording are rev-hashed; change
   without bump is CAN6013). How is the config pinned, hashed, and
   versioned? What breaks, and with which diagnostic, when it changes
   under a pinned consumer?
c. CAN4107 coverage over the new arm shapes, including the wildcard arm
   and the low-confidence / uncertain paths. Say exactly what executes
   what.
d. tests and given shape: keep JEV-1.0 full-observation evidence and show
   it does not balloon under the new syntax. One fenced example per
   option is mandatory.

## Deliverable

At least three genuinely distinct syntax options. For EACH option:

- Name and one-paragraph pitch.
- The config/metadata shape (file, keys, pinning) and what a use site
  inherits without restating.
- Noul examples: minimal yes/no; with true/false criteria; with explicit
  uncertain handling.
- Choice examples: minimal routing; confidence-routed high/medium/low;
  tie behavior; multi-question fan-out with one ignored answer.
- Score examples: minimal grading; level-routed branch; one example that
  inspects distribution shape beyond the mean.
- One example with its tests and given rows, proving the evidence shape
  stays tolerable.
- The failure path: the wildcard arm in context, and which failures it
  catches versus which stay explicit.
- An emit note: what the fetch-based, dependency-free TypeScript looks
  like in one paragraph plus a small sketch.
- Pros, cons, and which JEV-1.0 semantics it preserves versus asks to
  amend (each amendment tagged with new CAN codes).

Then:

- Your recommendation with reasoning, and a combined best-of sketch if
  one emerges across options.
- Open questions ranked by blast radius.
- A 10-line summary.

## Process

- Read before designing. Cite the exact doc section behind every TypeSafe
  capability your syntax exposes.
- Every example must be complete enough to read aloud: no "etc.", no
  illustrative abbreviations in place of real rows.
- Do not implement compiler changes. Return the deliverable as a markdown
  document plus the 10-line summary.
