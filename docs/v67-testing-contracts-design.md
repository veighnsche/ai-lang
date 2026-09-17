# v67: testing + contracts whole-story design

Status: proposal (pre-decision). No rule, no code, no golden.
This doc settles the testing story (including the v68
`runLinkedPure` shape) and proposes the contracts-epic
design (F# borrow #2, whose first step is a design
proposal). Refinements stay parked. Implementation slices
follow only on greenlight, probe-first, one per version.

## Part 1 — Testing story (settled)

### 1.1 Hermetic units: kept

Decision tables + `given` keep their semantics (verdict B).
A caller tests its response to any permitted outcome,
including errors the current provider never produces,
without depending on incidental provider behavior.
Failure injection is the point, not a limitation.

### 1.2 Arm coverage: kept

AIL4107 stays exactly as is (verdict C, pinned by
`compiler/diag_4107_test.go`): taken-ness is execution of
the particular arm under a passing test, or the existing
local identity-relay certificate. Script presence is not
coverage. The relay exception stays local-only; foreign
calls never qualify.

### 1.3 Linked-pure integration: the one addition (v68)

Adopt verdict B's `runLinkedPure` contract unchanged:

- Explicitly selected by committed Go tests only. No new
  `.ail` call syntax; `checkGiven` default untouched.
- Inputs: explicit module set, root function + revision,
  signature-checked args. Module presence never implies
  execution (verdict's core objection to presence-based
  admission).
- Precondition: ordinary whole-program checks, tables,
  revisions, export certification green first.
- Reachable graph: checked pure AIL or deterministic
  kernels only, all branches inspected. Reachable
  externs, state ops, effects, unresolved calls, or
  cycles refuse. Cross-module cycles stay rejected;
  checked direct self-recursion allowed.
- Execution: real bodies across modules, no `given`
  consumed, no script fallback, no mixed mode.
- Identity preserved (`Module.ID`, ownership,
  certificates); fresh context per vector.
- Failure fails; nothing is "trusted". Traces stay
  separate and never satisfy AIL4107.
- Falsifier (verdict): `runLinkedPure(middle__copy@1,
  "A")` returns `""` against `"A"`, order-reversed,
  plus a correct-leaf positive control.

### 1.4 Diagnostics: done

v62 (unknown-callee single diagnostic) and v63
(stale-arm nesting hint) shipped. No further
diagnostic work is queued by this story.

### 1.5 What testing does not give

Execution is not verification. A scripted response
consistent with a contract is still not execution or
verification of the provider body (borrow-doc §2
boundary). §1.3 closes the execution gap; Part 2
closes the universality gap. Neither subsumes the
other, and real execution alone cannot exercise a
hypothetical provider failure admitted by an
upper-bound `emits`.

## Part 2 — Contracts epic (proposal)

### 2.1 Outcome-indexed contracts

Borrow-doc §2 shape, unchanged: `requires` preconditions
plus `ensures` postconditions indexed by the actual
outcome (`on Ok result ...`, `on <err> err ...`). Every
declared error outcome gets its clause set: when the
error is permissible and what its payload means. The
`max` and validator sketches in the borrow doc are the
reference examples.

### 2.2 First-cut verification logic

Booleans, record projections, equality, linear integer
arithmetic (borrow-doc proof cost). Unsupported
propositions and inconclusive verification fail closed:
never runtime assertions, never assumed facts.
Checking a declared condition from known facts is
verification; it is not permission to invent types,
preconditions, or effects.

### 2.3 Contracts participate in revision review

An agent must not "fix" a failing proof by silently
narrowing the public input domain (borrow-doc §2).
`requires`/`ensures` edits are contract changes and go
through the same rev discipline as signatures and
`emits`. Private kernels may `require` (callers prove);
public functions validate untrusted input into declared
errors.

### 2.4 Tables stay mandatory

Concrete admitted test inputs are still required for
executable functions (borrow-doc acceptance test):
contracts quantify universally, tables witness
concretely. AIL3110 (scripted-Ok contradiction) stays
as the linkage check between the two.

### 2.5 The three proofs, stated once

- AIL3110: a scripted success must agree with what the
  provider body computes (example-level honesty).
- AIL4107: every handler must execute under a passing
  test (caller-level evidence).
- Contracts: the body must satisfy its outcome-indexed
  predicates for every admitted input (universality).
  A bug outside all example inputs must fail here
  while tables stay green (borrow-doc acceptance
  test); an impossible precondition masquerading as a
  proof must fail for lack of admitted inputs.

### 2.6 Parked, explicitly

Standalone refinements, ghosts, lexicographic
decreases: parked (unchanged). First follow-up on top
of this epic is cell state contracts
(`quota__consume` unchanged-on-error), already queued
behind #2. Refinements stay parked per 2026-09-17
decision even though §2.1 founds them.

## Part 3 — Slice order (needs greenlight)

1. v68: `runLinkedPure` (§1.3), probe-first with the
   verdict falsifier + control.
2. v69+: contracts, sliced after this design lands:
   `requires`/`ensures` grammar + AST, verifier core
   (§2.2 logic), revision-review hook (§2.3), stdlib
   pilot (borrow-doc `max`/validator shapes as rows,
   not new semantics). Detailed per-slice docs at
   that time; no multi-slice singletons.
3. Closed tagged unions (borrow #3) and everything
   downstream are unaffected by this doc.

Each slice: failing probes first, existing `.ail`
untouched except additive rows, committed
docs/generated, `go test -count=1 ./...` +
`modcheck` + `gramcheck` green.

## Rollback

Delete this file. Nothing else references it yet.
