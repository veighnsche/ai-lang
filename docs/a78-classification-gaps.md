# a78: classification gaps (verified investigation + proposal)

Status: investigation closed (all four gaps verified against
source and compiler probes, 2026-09-18); proposals below are
pre-decision. No rule, no code, no golden. Each item
graduates separately.

## Motivating case

`html__url__scheme_token` (std/html/html.ail): 60 lines,
~20 arms, 14 test rows, 10 magic ASCII bounds (`58 43 45
46 47 57 64 90 96 122`) — to express a 5-way character
class (letter / digit / `+ - .` / `:` / reject) with a
first-char-must-be-letter rule. The function is correct
and green. The complaint is not readability (agents do
not mind nesting) but obligation inflation: every ladder
rung is an arm, every arm is a test-per-arm obligation
(A4107), and every bound is an off-by-one slot the
checker then blesses. Verbosity the checker blesses is
more dangerous than verbosity it rejects (Goal
amendment: a weak table is satisfied by a wrong
implementation).

## Gap 1: int/range patterns (verified; deferred by decision)

- `parsePattern` (compiler/parse.go:1601) accepts exactly
  `_`, `true`/`false`, `"str"`/`e"str"`, variant/case
  patterns. Anything else is `bad match pattern`.
- Probe: arm `5 => 5` over an int scrutinee fails:
  `ailc FAILED: bad match pattern: 5`.
- Prior decision: a28 ("keep OR-patterns/guards/
  multi-call out of v1"). Deferred, not overlooked —
  revisiting needs a trigger, and the motivating case
  is now that trigger's best exhibit.
- Forcing mechanism, demonstrated by stdlib itself: int
  discrimination must go through comparison-then-bool-
  match ladders. No alternative exists.

## Gap 2: boolean combination (verified; never proposed)

- Operator tables (compiler/parse.go:498,872,889,900)
  are `== >= <= != > < + - * / %` plus `#`. No
  `&&`/`||`/`!`; no `and`/`or`/`not` keywords anywhere
  in parse/check/eval.
- Probe: `match a and b` fails: `cannot parse
  expression: a and b`.
- Partial mitigation only: a28 multi-match combines
  discrimination slots (`bool`/`str`/`_` per slot, no
  `given`) — the two-predicate case, no negation, no
  disjunction, no composition of arbitrary exprs. (a28's
  "boolean combination" phrase means tuple coverage,
  not operators.)
- Consequence: predicates cannot compose, which is why
  classifiers inline as ladders instead of composing
  as calls.

## Gap 3: named constants (verified; never proposed)

- Six decl kinds exist (compiler/parse.go:1304–1445):
  Error, Type, Brand, State, Extern, Fn. No Const.
- Probe: `const M__C: int = 58` fails: `unknown
  top-level decl`.
- Every workaround fails: magic literals (current —
  a direct violation of R5's "agents grep qualified
  names, never magic numbers"); zero-arg fns (`match
  call` + arms + test rows per constant); `state`
  cells (mutable, need effects, and cannot be provided
  per R2 — unusable as shared constants). Magic
  numbers are forced, not lazy.
- Cheapest gap to close; most embarrassing to lack.

## Gap 4: forward arms (verified; already proposed)

- `relay` exists only as a coverage-checker concept
  (certified identity relay, A4108), not as writable
  syntax. Cost receipt: 7x `on Ok r => Ok(has =
  r.has)` plus 7x full 4-arg self-call in the
  motivating case.
- Proposal already filed: a61 item 2. Not re-proposed
  here; listed for completeness and ranking.

## Bonus finding: grammar dead weight (verified)

Of the TextMate grammar's keyword/operator vocabulary,
`when then and reach expect mock case spec externals`
and `|` appear nowhere in the parser — the single hit
for the whole set is `"case"` as a `parseFields` label
argument (compiler/parse.go:1349), not a keyword. The
grammar promises a materially larger language than
exists. Cleanup is editorial (grammar + gramcheck
samples), not a language change; filed here so it is
not lost.

## What already exists (proposals must compose)

- The stdlib predicate pattern is established over
  int-denoted scalars (`std__int__is_even/odd/prime/
  multiple`, `std__str__is_whitespace`). But no
  char-class predicates (`alpha`/`digit`/`alnum` all
  absent) and nothing to combine them with (Gap 2).
- a28 multi-match (2+ slots, product-space proof).
- R2 (cells never provided) and R5 (grep names) bound
  the const design: whatever a constant is, it must be
  providible and greppable, or it is not the fix.

## Proposal direction (ranked by cost-of-absence)

1. Named constants first. Shape open (decl kind vs
   evaluable form), but the acceptance test is fixed:
   the ten magic bounds in the motivating case become
   greppable names with zero call/test overhead, and
   R4's future content-hash rule must treat them as
   inert (a typo fix in a name is not a breaking
   change — same canonicalization note as comments).
2. Forward arms second. Already specified (a61 item
   2); this doc adds the second consumer.
3. Range patterns third. Revisit the a28 deferral
   with the trigger in hand: `65..90`-style arms (or
   the smallest equivalent) collapse the ladder class.
   Biggest design surface of the four (proof rules,
   emit, coverage interplay) — design proposal first,
   like any epic.
4. Boolean operators fourth. Narrowed by multi-match
   but still needed for predicate composition. Pairs
   with a stdlib row of char-class predicates; neither
   alone fixes classification.

## Non-goals

- Guards/or-patterns/multi-call beyond the range-arm
  question (still out; a28 stands until item 3
  graduates).
- Any const design that is secretly a fn (call +
  arms + tests per constant) or secretly a cell
  (mutable, unprovidible). Both fail the acceptance
  test above.
- Readability refactors of the motivating case under
  current rules. It is correct; rewriting it without
  new expressive power is churn.

## Ordering and gates

1. Item 1 (const): design note with rev/hash
   semantics, then grammar + check + eval + emit +
   goldens; the motivating case's bounds are the
   migration proof.
2. Item 2 (forward): per a61 ordering (after
   diagnostic payloads).
3. Item 3 (ranges): design proposal first; needs the
   a28 authors' (or verdict-style) review since it
   reopens a settled deferral.
4. Item 4 (bool ops): with the char-predicate stdlib
   row; either may lead.

Each item graduates separately with a verdict-style
review. Item 1 may proceed alone. Evidence for all
claims above: three compiler probes
(`bad match pattern: 5`; `cannot parse expression:
a and b`; `unknown top-level decl`), rerunnable from
/tmp/gapprobe fixtures (scratch, not committed).
