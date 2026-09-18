## Goal

Close the linter gap behind `html__text__escape_from`: value-match ladders that should be tables pass `canlc lint` silently today. Ship the missing rule(s) as reporter-only findings, migrate the newly-flagged `std/html` workers onto the table shape, and correct the stale "bools cannot spell don't-care slots" rationale in code, docs, and tests.

## Success Criteria

- `canlc lint` flags the `s[0:1] == "lit"` ladders in `html__text__escape_from` and `html__attribute__value_from`.
- `std/` lints clean after migration (`TestLintStdClean` green), with worker behavior unchanged (existing rows green, goldens regenerated).
- The general pure-ladder miss is covered: a pure `match a / … / false => match b …` ladder is flagged, with the don't-care product as the suggested shape.
- No new CI gates, no `--fix`, no rev bumps: reporter-only findings under the existing error-severity lint contract (C13).

## Context And Current Facts

- Rule 4 (`lintTableable`, `compiler/lint.go:1239-1295`) fires only on the diamond: every outer arm nesting into an inner match on the identical, pure scrutinee (`sketches/lint-errors/table.can`). `escape_from` (`std/html/html.can:66-84`) fails all three preconditions: it is a ladder (one arm yields, one nests), each level tests a different scrutinee, and `s[0]`/`s[0:1]` are fault-capable, hence impure per `lintPureSmall` (`compiler/lint.go:448-463`).
- No other rule covers ladders: rule 2 needs identical outcomes with non-match RHS (`compiler/lint.go:425-437` refuses match RHS); rule 3 needs sequential call rungs with binder-free shared failure (`forward e` excluded, `compiler/lint.go:1189`); rule 5 needs all arms identical.
- Proven by probe (scratch `/tmp/esctest`, real tree untouched): both minimized spellings of the worker compile with the existing 11 rows green, and a 1366-case differential run of the emitted TS shows zero mismatches vs the current shape. Pure-ladder and pure different-scrutinee probes compile and lint silent — the miss reproduces with purity removed.
- The "bools cannot spell don't-care slots" rationale (`compiler/lint.go:1236-1238`, can-idioms C9, `table.can` header, `compiler/lint_test.go:343-344`) is not a language fact: `match a == 0, b == 0` with `true, _` arms compiles (blessed `std__int__lcm` does this; CAN4103's no-`_` rule is arity-1 only). Only the impurity half of the refusal is law (CAN4109).
- Forcing functions: `TestLintStdClean` (`compiler/lint_test.go:983`) requires all 8 std modules lint-clean, so any rule that fires on `std/html` must ship with the migration. `TestGoldenStdHtml` pins the emitted `std/html/html.ts`, so migration must re-emit. Fingerprints exclude bodies (`compiler/revision.go:18-20`), so body-only rewrites need no rev bump. Lint is not a CI gate (workflows are `tsc.yml`, `verifier.yml`); the gate is `go test ./...`.
- `std/html/html.can` carries 70 literal NUL bytes; any edit must preserve the count exactly (a88 §4 byte safety).

## Constraints And Non-goals

- Reporter-only. Findings are error-severity per C13; no `--fix` (a88 check-before-fix ethos).
- One rule, one code (a28 decision 5): each new rule takes the next free code (CAN3417, CAN3418; CAN3410-3416 taken, `compiler/code.go:70-76`) and the next C-number (C14, C15).
- `!=` rungs, bool- and int-literal `==` rungs, and multi-arm (non-bool) ladder products are out of v1 scope and stay silent by design.
- No verifier/checker changes: both rules are parse-level, like the existing seven.
- Publication split: slice 1 and slice 2 ship as two separate commits, in order; never collapsed.

## Key Decisions

1. **Two rules, sequenced — `==`-restatement first (C14/CAN3417), ladder→product second (C15/CAN3418).** The motivating case needs the first: slice 2's purity gate excludes `escape_from` (fault-capable scrutinees), while slice 1's fault-neutrality argument covers it. Slice 1 also has no existing-test churn; slice 2 moves one negative test. Rejected: slice-2-only (misses the motivating case) and one combined rule (different soundness stories must not share a code).
2. **Slice 1 trigger: `==`-ladder over a common base with uniform str literals.** A chain of ≥2 value matches, each rung scrutinee a `==` binop (`compiler/parse.go:47-48`, `Op`/`L`/`R`) with the same base on one side and a str literal on the other; true-arm = rung body (anything, verbatim), false-arm = next rung; final false-arm = default. Suggested shape: `match <base>` with literal arms plus `_`. Base identity is structural (`smallEqual` helper — none exists today) rather than textual.
3. **Slice 1 soundness rests on fault-neutrality, not purity.** The table evaluates the base once at the same program point the first rung would; value scrutinees are call-free, so re-evaluation count is unobservable except for faults, and faults occur at the identical point. Guard-safe: guards are `<=`-shaped (`compiler/check.go:1651`), rungs are `==`-shaped — disjoint, and the table sits where the ladder head sat, under the same guards. Coverage preserved: ladder true-arms map to literal arms 1:1, the final false-arm maps to `_` (proven by the unchanged 11 rows on the migrated workers). NUL-policy safe structurally: the NUL rung (`s[0]`, int const) differs in both base and literal kind from the `s[0:1]` str rungs, so it can never join the group — no invisible bytes enter the body.
4. **Slice 1 refusals: single rung (no savings), mixed literal kinds, bool/int literals, duplicate literals, `!=`.** Bool rungs are excluded because the restated arity-1 table would need `_` on a bool slot (CAN4103); int rungs because int-slot `_` is unverified. Duplicates are dead code the checker already rejects (CAN4107); the linter refuses to bless them.
5. **Slice 2 trigger: bool-rung ladder over pure, pairwise-different scrutinees.** Exactly one arm per level nests deeper, the other yields; every rung scrutinee pure per `lintPureSmall`; suggested shape is the product table with don't-care arms (the `std__int__lcm` shape). Purity is required here — unlike slice 1, evaluation moves earlier. Guarded/impure ladders stay silent (CAN4109); after slice 1, `escape_from`'s remaining nest (guard → NUL check → str table) is minimal and silent by design.
6. **Slice 2 corrects the stale rationale at all four sites** (`compiler/lint.go:1236-1238`, can-idioms C9, `table.can` header, `compiler/lint_test.go:343-344`) and moves `demo__diff` from the Untouched test to a positive test. This is a planned semantic change to `TestLintTabledKeptUntouched`, justified by the compiling don't-care probes — not a regression.
7. **Span convention follows rule 4: underline the match head** via `lintHeadSpan`. The base-expression token would be more precise but breaks the one-rule-family visual contract; revisit only if reviewers ask.

## Recommended Approach

Two slices in order. Slice 1 implements C14/CAN3417 (`==`-restatement), adds the 8th fixture, migrates both `std/html` workers to the str-table shape (variant A from the probe), and regenerates goldens. Slice 2 implements C15/CAN3418 (pure-ladder→product), adds the 9th fixture, moves the stale negative test, corrects the four rationale comments, and surveys the repo for new findings. Each slice follows the established rule precedent: `lint.go` check + `code.go` registry entry + inline unit tests (positive and `Untouched` negatives) + one fixture file + span pins + can-idioms C-entry.

## Work Plan

### Slice 1 — C14/CAN3417 `==`-ladder restatement

1. **Check.** Add `smallEqual` (structural `*Small` comparison) and `lintRestatable` to `compiler/lint.go`; wire into `lintFiles`. Trigger per decision 2, refusals per decision 4, head span per decision 7. Add `CodeLintRestate = "CAN3417"` to `compiler/code.go` registry.
2. **Unit tests.** `TestLintRestatable` (positive: 3-rung str ladder over `s[0:1]`, exact finding text) and `TestLintRestateKeptUntouched` (single rung, mixed str/int kinds, bool lits, `!=`, duplicate literals, ladder over a `<=` guard base — all silent).
3. **Fixture.** Add `sketches/lint-errors/restate.can` (compiles clean, lints dirty); bump `TestLintErrorFixtures` count 7→8 with pinned finding lines; add span pins to `TestLintSpanPositions`.
4. **Migrate `std/html`.** Rewrite the inner three rungs of `html__text__escape_from` and `html__attribute__value_from` as `match s[0:1]` tables (guard, NUL rung, tests, and revs untouched). Byte-exact edit; assert the file's NUL count is still 70 before/after (`tr -d -c '\000' | wc -c`). Re-emit `std/html/html.ts` and `errors.json` via `canlc --out`; confirm the TS diff touches only the two workers.
5. **Docs.** Add can-idioms C14 (renumber nothing; append after C13's rule list and fixture paragraph).

### Slice 2 — C15/CAN3418 pure-ladder→product

6. **Check.** Add `lintLadderable` to `compiler/lint.go` (ladder walk, `lintPureSmall` gate, pairwise-different scrutinee texts, bool rungs); wire into `lintFiles`. Add `CodeLintLadder = "CAN3418"` to `compiler/code.go`.
7. **Unit tests.** `TestLintLadderable` (positive: pure 2-rung ladder, exact finding text) and `TestLintLadderKeptUntouched` (impure/guarded ladder over `s[0:1]`, identical-scrutinee ladder, already-tabled multi-match — all silent). Move `demo__diff` from `TestLintTabledKeptUntouched` to a positive folding test; keep `demo__guarded` untouched.
8. **Fixture.** Add `sketches/lint-errors/ladder.can`; bump fixture count 8→9 with pinned lines; add span pins.
9. **Rationale correction.** Rewrite the four stale comments (decision 6) to state the real rule: different pure scrutinees fold with don't-care slots; impure ones stay guarded (CAN4109).
10. **Survey.** Run `canlc lint` over `std/` (must be clean — the slice 1 migration leaves no pure ladder), `sketches/`, and `docs/ALL_EXAMPLES.can` if applicable; migrate trivially-safe gallery hits, report the rest without expanding scope. Add can-idioms C15.

## Validation Plan

- Per slice: `go test ./...` (pins `TestLintStdClean`, fixture goldens, span positions, idempotence), `go run ./tools/modcheck`, `go run ./tools/gramcheck`.
- Slice 1 additionally: `canlc lint std/` exits 0; `tsc -p tsconfig.json` in `tscheck/` (re-emitted `html.ts` must typecheck); NUL-count equality on `std/html/html.can`; re-run the `/tmp/esctest` differential harness shape against the migrated workers (11 rows green at minimum; the 1366-case sweep as the strong check).
- Slice 2 additionally: `canlc lint sketches/lint-errors/` exits 1 with exactly the pinned findings (9 files); repo survey report attached to the change.
- Highest-risk validation: the `std/html` migration review — behavior preservation rests on the existing rows plus the differential sweep, so the TS-diff-must-touch-only-two-workers check is the backstop against accidental semantic drift.

## Risks / Rollback

- **Slice 2 test churn looks like regression.** `TestLintTabledKeptUntouched` loses `demo__diff` by design. Mitigation: the move ships with the probe evidence cited in Context; the `demo__guarded` half stays, proving the rule still refuses impure nests.
- **NUL-byte edit hazard in `std/html/html.can`.** A naive editor pass can normalize the 70 NULs. Mitigation: byte-exact edit path plus the count assertion in step 4; any mismatch fails the slice.
- **Golden churn beyond the two workers.** Re-emit should be local; a wider `html.ts` diff signals emitter drift, not migration. Mitigation: inspect the diff before accepting; rollback is `git checkout -- std/html` plus slice-1 redo.
- **New findings in gallery sketches.** `sketches/broken-login` already lints dirty (pre-existing rule-1 findings), so dirty gallery is precedented. Mitigation: survey-and-report (step 10); migrate only trivially-safe hits.
- **Over-triggering on clever ladders.** A ladder whose rung order encodes priority the table preserves (yielding arm first with don't-cares) is still equivalent; order-dependence is not a refusal. True false positives surface via `TestLintStdClean` and the survey, and narrow the trigger before merge.
- Rollback per slice is a single revert (no checker, CI, or baseline changes in either slice).

## Open Questions

1. Should slice 2 v1 stay bool-rung-only, or also cover multi-arm ladders (exactly one nesting arm with str/int patterns)? Recommendation: bool-only; the product-emission validity argument is simplest there.
2. Gallery policy for slice 2 survey hits: migrate trivially-safe ones in-slice, or report-only and leave all gallery migration to a follow-up? Recommendation: migrate-if-trivial, report the rest.

## Sources

None — every material claim traces to workspace evidence inspected during planning (cited inline by file and line).
