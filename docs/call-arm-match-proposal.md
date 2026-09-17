# Proposal: call scrutinees in multi-scrutinee matches (slice 0)

Status: proposed, awaiting approval. No compiler edits follow from this
document until it is approved; the work plan in section 9 is structural
commitment on approval.

## 1. Goal

Extend `match x, y` (docs/v28-multi-scrutinee-match.md) to the corpus v1
cannot serve. Migration frequency data: division's only nested pair, all
html nesting (`std/html/html.ail` escape ladders), and the whole
`docs/ALL_EXAMPLES.ail` effect corpus (auth-login, retry-loop, counter)
are `match call` chains or call-bearing nests. v1 serves pure value
tables; the deepest remaining nests stay nested.

Success: product matches over call scrutinees and per-slot variant
patterns compile, prove exhaustive before emit, keep one test obligation
per source arm, and convert the section-8 targets with no change in
meaning to untouched modules' `errors.json`.

## 2. Background and current split

- `MatchValue`: one or more scrutinees, bool/str/wild patterns only.
  Variant patterns rejected (`AIL4106`, `emit.go` arm gate); call
  scrutinees rejected (checker shape rule mirrored by
  `verifyValueMatch`, `eval.go`; `evValueMatch` fails closed).
- `MatchCall`: exactly one call scrutinee, Ok/error patterns, `given`
  support (`check.go:checkGiven`). Parser assigns `MatchCall` iff one
  scrutinee and it is a call (`parse.go:parseMatchArms`).
- `parseScrutList` already parses comma-separated scrutinee lists and
  `parseSmall` already parses `call f(x)`, so `match call db__get(id),
  attempt` parses today and dies only in the checker. No grammar change
  is proposed; slice 1 verifies the splitter is paren-aware for the
  mixed shape (it already handles calls-with-args singly and value
  lists plurally — only the combination is untested).

## 3. Semantic decisions

- **D1 (abort-on-raise).** Scrutinees evaluate left to right; a scrutinee
  that raises aborts the whole match and later scrutinees never
  evaluate. The error propagates; no arm is tried. This matches the
  nested operational reading (`match call f / arm => match call g`
  never runs `g` when `f` raises).
- **D2 (exactly-once LTR retained).** Call order is observable and will
  be pinned by an abort-ordering test (highest-risk validation).
- **D3 (per-slot variants; raises never match).** A raising scrutinee is
  not a value and never reaches arm testing (follows from D1). Value
  slots keep bool/str/wild; call slots take Ok-binds, error kinds, `_`.
- **D4 (`given` iff exactly one call slot).** A `given` table scripts one
  foreign call; with two or more call slots the scripting target is
  ambiguous, so that stays an error (existing `CodeGivenOnLocal`
  message, extended). One call slot reuses the existing per-callee
  logic unchanged. This refines the pre-plan note ("disallow given past
  arity 1"): auth-login and retry-loop each carry one foreign call per
  match, so the refined rule keeps them convertible while the blunt one
  would not.
- **D5 (witness budget 3, slot-major enumeration).** Product coverage
  enumerates the cartesian slot-major for deterministic messages.
- **D6 (recursion arms allowed).** A call is a call; termination is
  governed by the existing decreases machinery per function regardless
  of match shape. Slice-4 proof: retry-loop converts and its decreases
  proof still passes unchanged.
- **D7 (new `MatchKind`, legacy path untouched).** Parser keeps emitting
  `MatchCall` for the single-call shape, which keeps every legacy
  golden byte-identical. Any match with a call scrutinee past arity 1,
  or a mixed call/value slot list, becomes the new kind and takes the
  new check/eval/emit branches. No legacy churn by construction.
- **D8 (`_` allowed on call slots).** Uniform with v28 value slots.

## 4. Syntax

No new syntax. `match call f(x), y, z == 0` — each slot a value
expression or a `call` expression, comma-separated as today. Mixed
call/value slot lists are the expected common shape (call outcome
crossed with a local condition).

## 5. Patterns per slot

Value slots: bool/str/wild (unchanged). Call slots: `on Ok w` binds,
`on mod.kind payload` error patterns, `_` (D8). Slot-count mismatch
and kind errors reuse the v28 tuple-rendered messages.

## 6. Coverage and proof

Per-slot domains: bool slots `{true, false}`; call slots `{Ok} ∪
emits-of(callee)`, where emits stays the v12 upper bound (unrealized
entries admitted, so tables must still cover them — same rule as
single-call today, crossed per slot). Product coverage with D5
enumeration; missing cells render slot-major up to the budget of 3.
Proof-before-emit and one-obligation-per-arm hold as today.

## 7. Codes

Proposed: zero new codes. Shape violations reuse the v28
tuple-rendered messages; multi-call `given` reuses `CodeGivenOnLocal`
with an extended message. Slice 2 confirms; any new rule discovered
there returns here for a code decision before implementation.

## 8. Emit composition

One temporary per call scrutinee (per-function scope, just shipped),
kind-tests per call slot composed into the existing conjunction chain
(`if ($c1ok && $b2)`), proved `else` preserved, totality as today.
Detail belongs to slice 3; the principle is composition, not a second
chain shape.

## 9. Work plan (structural commitment on approval)

1. **Eval** (`compiler/eval.go`): LTR call-scrutinee evaluation with
   D1 abort; per-slot matching reusing `matchSlot`; `markTaken` per
   winning arm. Tests: abort ordering via two-effect stubs; per-slot
   payload binding.
2. **Check/coverage** (`compiler/check.go`, prover): per-slot variant
   allowance; product coverage with D5; D4 `given` rule; D7 kind
   assignment. Tests: missing-witness message golden for a 2-call
   table; `given`-with-two-calls error.
3. **Emit** (`compiler/emit.go`): section-8 composition.
   Tests: chain assertions for a 2-call table, `multimatch_test.go`
   style.
4. **Dogfood + goldens**: convert html ladders, auth-login,
   retry-loop, counter; division pair stays nested (trapping
   scrutinee — now `AIL4109`-linted). Regen all goldens with
   normalized-diff review (temp-numbers stripped); full gates
   (`go test ./...`, `modcheck`, `gramcheck`).

## 10. Validation

Per slice: full gates. Epic gate: `errors.json` byte-identical except
intentionally converted tables; converted tables' new witnesses
reviewed by hand. Highest-risk step: slice-1 abort-ordering test —
it pins observable effect semantics.

## 11. Risks

- Effect ordering becomes load-bearing in goldens (desired: loud on
  later eval changes).
- Legacy single-call output must stay byte-identical (D7 structural
  guarantee; suite enforces).
- Normalized golden review is the control for mechanical churn.

## 12. Non-goals

Guards, OR-patterns (reaffirmed rare by migration data);
lazy/short-circuit scrutinees (eager LTR is doctrine);
`given` with two or more call slots (D4).

## 13. Open questions

None. D4 refines the approved pre-plan note as marked; D6 carries a
slice-4 proof obligation instead of a question.
