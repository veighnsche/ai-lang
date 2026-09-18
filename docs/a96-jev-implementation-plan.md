## Decision recap

1. Implement versioned B declarations for `noul`, `choice`, and `score`, with ordinary `provides`/`uses` ownership and revision pins.
2. Admit only B singleton dispatch and C acquisition/`decide` scopes for B-declared batches containing at least two questions.
3. Preserve complete observations, deterministic routing, partial `given` queues, source-arm coverage, and proof/fixture separation.
4. Compile one pinned transport profile into dependency-free TypeScript using platform fetch; read only `TYPESAFE_API_KEY` at runtime.
5. Activate the surface only through the ordinary compiler/editor gate after identity, evidence, scope, transport, and emit conformance are integrated.

## 1. Assumptions and HEAD-drift notes

**Reviewed HEAD:** `71cd229155fc6087b2e938360e33d5de1a3b442c`. A second branch read returned the same commit. The plan below targets that snapshot, not a95’s original `4d447bac858c` baseline. **No compiler changes or gates were run.**

### Verified source changes and implications

| Finding at HEAD                                                                                                                                            | Implementation consequence                                                                                                                                                                                                 |
| ---------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Both a94 and a95 now exist.** Their repository blobs are byte-identical to the earlier attached documents, established by comparing Git blob hashes.     | Use repository a94 as the JEV-1.0 source. a95’s “a94 unavailable” receipt is historical, not a present limitation. Its `[JEV1]` section references remain applicable.                                                      |
| a94 retains the SDK, per-descriptor transport, old descriptor syntax, and individual intrinsic-error arms. a95 retains all three candidate grammars.       | Append a tagged selection/supersession record. Implement the selected surface directly; do not implement the rejected surface as an intermediate compatibility layer.                                                      |
| **Neither CAN35xx nor CAN36xx is registered.** `CAN3014` now belongs to generic expansion.                                                                 | Budget the a94 evidence foundation as new implementation work, not an existing service beneath a grammar change. Reserve only the selected rules.                                                                          |
| `expand.go` now stamps generic functions before checking.                                                                                                  | New nodes, observations, source identities, and request effects must survive expansion and every associated AST walk. Do not design against the older nongeneric pipeline.                                                 |
| `emit.go` uses **`$can_kind`**, not a payload-visible `kind`, and already handles `__proto__` specially.                                                   | Extend the actual emitted ABI. a95’s illustrative transport/result sketches are not authority to replace existing tag or field-encoding conventions.                                                                       |
| a81’s “unwired” status is superseded by **a82 activation**.                                                                                                | Judgment effects must be handled in the active contract admission/proof pipeline from first integration.                                                                                                                   |
| `RevisionFormat` is currently `1`; ordinary CLI compilation can omit a baseline. a79 also documents missing-baseline LSP fallback to unenforced operation. | Version the baseline extension explicitly. Distinguish “compiled” from “revision enforcement established”; fix missing-*selected*-baseline CLI/LSP disagreement rather than implying automatic enforcement already exists. |
| `CAN3104` remains a compiler warning, including in `TestGoldenJSONDiags`.                                                                                  | Implement a94’s explicitly tagged promotion to error, including resulting diagnostic-order/suppression changes. Do not silently inherit the split from modcheck.                                                           |
| Both modcheck and the TextMate `given` region depend on finding `on`.                                                                                      | Batch `decide` syntax requires boundary handling, not merely adding keywords.                                                                                                                                              |
| Current goldens include `Notify`, revised `Reserve`, and changed stdlib outputs.                                                                           | Preserve the **HEAD inventory**, not only the old auth/retry/counter examples.                                                                                                                                             |

Receipts: a94/a95; current registry and expansion; emitter; verifier activation; revision/baseline implementation; diagnosis golden; module scanner and editor grammar; current emit goldens.

### Two integration resolutions to record before implementation

**CAN3620 needs a tagged clarification.** a95 §7 assigns it to duplicate *explicit* local site IDs, while the selected §8 hybrid removes those IDs and uses structural identities. This plan retains **one site-identity uniqueness rule**, but applies it to the resolved/persisted structural-site inventory. Repeated calls to the same descriptor are legal distinct sites; repeated execution of one site consumes its queue. A duplicate identity in an accepted-baseline inventory is a diagnosable input fault; a collision manufactured internally is a compiler defect. Do not add explicit site-ID syntax merely to make the old wording applicable. This clarification must precede registration.

**Exact-analysis exhaustion needs an honest diagnostic.** This plan proposes **CAN3641 / `CodeJevAnalysisLimit`**, under `JEV-S1/EXACT-ANALYSIS`, for exhaustion of the deterministic judgment exhaustiveness-analysis budget. It must not masquerade as CAN4101 “missing arm,” CAN4107 “untaken arm,” or a Z3 finding. This is an additional implementation-rule allocation, not another judgment grammar.

The TypeSafe API pages could not be freshly retrieved in this pass. The transport plan therefore implements the repository’s settled a95 contract; it does not claim new live-service conformance evidence.

## 2. Phased implementation plan

### Shared gate commands

The following commands are **future acceptance commands**, not reported results:

```sh
go test -count=1 ./...
go run ./tools/modcheck
go run ./tools/gramcheck
(cd tscheck && ./node_modules/.bin/tsc -p tsconfig.json)
```

Call this set **G** below. Locked dependency provisioning is separate:

```sh
(cd tscheck && npm ci --no-audit --no-fund)
```

The current strict-TS job uses Node 24 and the repository pins TypeScript `7.0.2`. Preserve those checked-in toolchain inputs unless a separately documented change is necessary. Ensure the job running Go parity tests also has the required Node runtime; currently the verifier workflow explicitly installs Z3 and Go, but not Node. No test may silently skip a required Jev conformance check because its host is unavailable.

**Dependency order**

```text
P0 → P1 → P2 → P3
                ├→ P4 identity / acceptance / proof
                └→ P5 batch scopes

P6 transport work can start after P2;
its integration depends on P3 and P5.

P4 + P5 + P6 → P7 ordinary-path activation
             → P8 migration
             → P9 final qualification
```

Tests and editor support are developed alongside their owning phases. P7 is the integration checkpoint, not permission to postpone all tooling work until the end. Before activation, unfinished components remain internal, with the public compiler refusing unsupported entry syntax. There is no feature flag or privileged Jev compilation mode.

### P0 — Freeze the selected contract and integration ledger

**Goal.** Turn the winning combination into one implementable grammar/identity contract, with historical alternatives clearly nonnormative.

**Files.** Amend `REQUIREMENTS.md`, `docs/a94-jev-judgment-evidence-decision.md`, `docs/a95-jev-native-syntax-options.md`, `docs/can-idioms.md`, and the documentation index. Add a selected-surface implementation record, proposed as `docs/jev-native-implementation-plan.md`.

Record the ten shared amendment tags from a95 §7, the selected B/C rules, the CAN3620 clarification, and the proposed CAN3641 rule. Record that `match judgments`, A inline forms, explicit A/C local IDs, old `judgment`/`match judge`, and singleton acquisition scopes are not legal alternatives.

Define the common internal contracts before parallel work: resolved question schema; profile identity; complete observation representation; structural source-site identity; source-arm identity; acquisition scope; and request-authority closure. These are compiler concepts, not new user-facing generic types.

**Diagnostics.** Reserve the a94 CAN3501–CAN3513 rules, common CAN3601–CAN3610, selected CAN3620/CAN3630/CAN3631/CAN3640, and proposed CAN3641. **Do not allocate CAN3621**, which belongs to excluded A syntax.

**Fixtures planned.** Freeze the selected a95 §8 source specimen as the principal batch acceptance fixture. Specify corresponding Noul, Choice, Score, and structured-state singleton fixtures. Inventory all committed nonjudgment TS/catalog goldens at HEAD.

**Done criteria.** No implementation ticket has to guess the selected grammar, identity boundary, source coverage unit, or error abstraction. Every reserved code has one written rule and an owning phase.

**Gates/risk.** Documentation review in this phase; G remains the merge baseline for later implementation. Main risk: allowing a specification inconsistency to become an accidental permanent grammar.

### P1 — Transport profile, provider resolution, and diagnostic infrastructure

**Goal.** Establish `jev.config.json` as a real pinned provider without creating ambient defaults.

**Files.** Add a small shared profile reader/canonicalizer under `internal/jevconfig/`. Integrate through `compiler/check.go`, `compiler/main.go`, `compiler/lsp.go`, and `tools/modcheck/main.go`. Extend `compiler/code.go`, `compiler/explain.go`, and their tests.

Admit exactly the a95 profile:

* Format `can-jev-transport/1`.
* Name `Jev__Runtime`, positive revision, adapter `typesafe-systemone-fetch/1`.
* Full versioned model ID; no aliases or discovery.
* Finite permitted timeout and `retries: 0`.
* No configurable endpoint, credential value, environment-selected profile, or per-descriptor override.

Resolve the provider from the existing project root. Imported declarations retain their owning profile identity; incompatible profiles cannot silently merge. Canonical hashing excludes `rev`, preserves decoded strings without Unicode normalization, and rejects duplicate/unknown keys before materializing an ordinary map. Preserve a95’s diagnostic precedence for missing/invalid model and timeout fields.

**Diagnostics.** Implement CAN3601–CAN3603 and retained CAN3509, CAN3510, CAN3512. Reuse provider/pin/shadowing diagnostics where their existing rules apply.

**Tests.** New `compiler/jev_config_test.go` and shared-package tests:

* Canonical hash golden; formatting/key-order neutrality.
* Duplicate keys, escaped duplicate keys, missing/extra fields, alias, invalid timeout, nonzero retries.
* Missing provider, stale pin, incompatible ownership, user shadowing.
* Different working directories and unrelated parent config files.
* Profile-free ordinary programs remain unchanged.

**Done criteria.** Compiler, editor, and modcheck select the same provider and canonical identity without reading credentials or consulting the service.

**Gates/risk.** G. Avoid loading Jev intrinsic errors into every program indiscriminately: that would pollute nonjudgment catalogs, whose current builder inventories all registered program errors.

### P2 — Shared grammar, metadata/state trees, observations, and numeric validation

**Goal.** Build the reusable representation needed by both singleton and batch forms.

**Files.** Extend AST/parser definitions in `compiler/parse.go`; add focused files such as `compiler/jev_parse.go`, `compiler/jev_schema.go`, and `compiler/jev_observation.go`. Integrate `types.go`, `bind.go`, and the cloning/traversal paths in `expand.go`.

Use separate typed contexts for **closed metadata** and **runtime text-state trees**. Metadata admits the exact §7 literal domain, including safe integer/representable decimal literals and metadata-only `unspecified`; runtime state remains unbranded text/object/array only. Instruction roots must be present and non-null. Optional Noul criteria are either absent or contain both true and false descriptions. Question/support/level order remains identity-bearing.

Implement complete `observed noul`, `observed choice`, `observed score`, and `observed answers` literals exclusively at judgment exchange outcomes. No confidence-band fixtures, inferred complements, fabricated fields, or general observation constructor becomes available in production expressions.

Implement the shared pure observation validator and deterministic classifiers. Preserve a94’s numeric image, mass allowance, maximizer requirement, Score consistency check, and exact routing comparisons; do not invent a confidence formula.

**Diagnostics.** CAN3501–CAN3508, CAN3513, CAN3605–CAN3609. Existing type mismatches remain CAN6003.

**Tests.** New `jev_parse_test.go`, `jev_schema_test.go`, and `jev_observation_test.go`, backed by `compiler/testdata/jev/observations/`:

* Every metadata location, nested structure, `unspecified`, null/missing instruction rejection.
* Numeric-looking, Unicode, `constructor`, and `__proto__` keys.
* Forbidden runtime numeric/Boolean/binary/branded leaves.
* Missing/extra observation fields, wrong primitive, incomplete batch, forbidden Noul confidence.
* Mass/score tolerance boundaries, unrepresentable literals, selected-label mismatch.
* Literal-only valid Choice support projections and Score indexes; dynamic/out-of-support projections reject.

**Done criteria.** One validated observation representation feeds evaluation and later emit conformance. No parser or expansion path loses new nodes or their source locations.

**Gates/risk.** G. Highest local risk is context confusion: permitting a runtime value in metadata or treating `unspecified` as a general can value.

### P3 — B declarations and singleton execution

**Goal.** Implement each primitive’s declaration and singleton dispatch through the existing evidence/coverage machinery.

**Files.** Extend `check.go`, `types.go`, `eval.go`, `lsp.go`, and `lint.go`; add `jev_singleton.go`, `jev_eval.go`, and `jev_coverage.go` as focused integration units.

Deliver three subordinate slices: **Noul → Choice → Score**. Each includes declaration resolution, singleton head checking, binder typing, exhaustive successful-state checking, the terminal failure clause, and actual table execution—not parser-only support.

A same-file question declaration is still a foreign judgment boundary. It does not inherit the “same-file helper executes without `given`” rule. Ordinary local helpers continue to execute normally. Register applicable queues before execution so an entirely unvisited nonempty queue is detected; keep root tests and foreign scripted boundaries distinct.

Treat `Jev.request` as an explicit request capability, not a fictitious state cell. Validate transitive authority and stale/undeclared capabilities. Update lint walkers so they cannot erase an acquisition, fold it into ordinary pure matching, or automatically batch requests.

**Diagnostics.** CAN3630 for primitive/declaration mismatch; CAN3604 for malformed failure-clause shape. Extend existing CAN4101, CAN4107, CAN4111/CAN4112 and effect diagnostics. Promote CAN3104 to error under the explicit a94 amendment.

**Tests.** New `jev_singleton_test.go`, `jev_evidence_test.go`, and `jev_coverage_test.go`; extend `effects_test.go`, `helper_test.go`, `diag_4107_test.go`, and `json_golden_test.go`.

Cover every Noul state; every high Choice label, medium, low, and high-confidence tie; Score confidence bands and equal-mean/different-distribution cases. Exercise missing/reached scripts, exhaustion, wrong arguments, repeated differing observations, partial leftovers, zero-consumption leftovers, helper reach, unknown keys, and retired `-`.

For coverage, one grouped successful arm or failure clause owns **one source obligation**. Either intrinsic failure can execute `_`; testing both kinds is required in compiler conformance, not a second author coverage obligation.

**Done criteria.** Table execution chooses source arms only after validation/classification. A malformed successful fixture diagnoses before execution; it never reaches `_`.

**Gates/risk.** G. Preserve the existing executed/certified/uncovered distinction; do not extend the relay certificate to judgment uncertainty or failure clauses.

### P4 — Revision closure, accepted evidence, and active proof integration

**Goal.** Make request identity enforceable before the surface is presented as reviewed or usable.

**Files.** Extend `compiler/revision.go`, `revision_identity_test.go`, `acceptance_test.go`, `lsp_baseline_test.go`, `verify_admission.go`, `verify_prove.go`, and the contract suites. Add `jev_revision_test.go`, `jev_acceptance_test.go`, and `jev_proof_test.go`.

Version the baseline format explicitly. Add profile, question, batch, and **transitive reachable request-contract closure** identity. Keep that separate from ordinary body/proof freshness and fixture evidence. Include the selected contracts reachable through local helpers and cross-module calls—not all declarations in the project and not pure consumers of returned data.

Enforce the §2.2 precedence: resolve provider/pin faults first; then compare identities. Same-revision metadata drift is CAN6013; a stale provider pin is CAN2103; repinning a changed profile cannot silently preserve an affected request contract’s revision.

Use the existing general baseline-selection mechanism for CLI and LSP. The accepted project workflow supplies the same baseline to both. Missing *selected* authority must not become a clean unenforced editor session. Unbaselined bootstrap/generation must remain visibly unverified for revision acceptance; do not invent a Jev-only verifier command. Current baseline behavior makes this an actual integration obligation.

Implement CAN3511 over pinned inputs, relevant exchanges, reach-controlling exchanges, effective descriptors, and structural site identities through helpers. Preserve CAN6017 independently. Generated baselines never accept themselves.

Add explicit judgment cases to active proof admission/traversal. Unsupported effects produce unsupported/unverified findings; neither `given`, confidence, pinned status, nor a disappearing AST node can create a proof premise.

**Diagnostics.** CAN3511 and clarified CAN3620; existing CAN6013/CAN2103, CAN6017, CAN4302/CAN4303.

**Tests and measurement.** Full metadata mutation matrix; baseline format/authority rejection; comments/literal-equivalence controls; fixture-only changes; removal of the final judgment; generic expansion identity controls. Add the multi-module measurement fixture specified in §4 below.

**Done criteria.** The identity matrix is executable, the measured affected set equals the intended request closure, and proof results cannot change merely because a scripted observation changes.

**Gates/risk.** G. **Q1 is the principal release risk:** excessive revision propagation must be measured, but never silently “fixed” by weakening the settled closure.

### P5 — B batch composition and C acquisition/`decide` scopes

**Goal.** Add batches without introducing a second singleton spelling or repeated requests during decision reuse.

**Files.** Extend the shared parser/schema/evaluator and add `compiler/jev_scope.go` plus `jev_scope_test.go`.

Resolve `judgments` members as B primitive declarations with matching sole-state signatures. Reject zero/one-question batches. `judge descriptor(...)` evaluates state once, consumes one exchange, validates the entire response, then establishes the descriptor’s observation names in its successful continuation.

`decide` is pure selection over those observations. Nested decisions reuse them without another fetch or queue consumption. An arm may invoke a later explicitly handled operation; its failures are not caught by the outer acquisition’s `_`.

Check successful product coverage independently from the failure clause. Ignoring an answer does not excuse missing or malformed evidence. Dependent questions remain later acquisitions, never compiler-created fan-out.

**Diagnostics.** CAN3610, CAN3631, CAN3640, and proposed CAN3641; reuse successful-arm and failure-clause diagnostics.

**Tests.** New `jev_batch_test.go`, `jev_scope_test.go`, and `jev_analysis_test.go`:

* Mixed primitive batches, imported members, state mismatch.
* Scope leakage, unavailable observations, nested acquisition ownership.
* Ignored valid answer, ignored malformed answer, incomplete joint fixture.
* Repeated `decide` versus a second `judge`, with exact request/queue counts.
* Repeated descriptor uses as distinct structural sites.
* Success product gaps, alternatives, shadowing, and deterministic analysis-limit refusal.

Use an explicit, versioned analysis work budget. The initial engineering assumption is **100,000 analysis work units per decision**, with the work unit defined and counted by the checker—not wall-clock timeout or sampling. P5/P9 measurements may change that provisional limit before qualification; the checker must never accept after incomplete analysis.

**Done criteria.** The a95 §8 fixture has four successful routing arms and one source failure arm; its two failure rows hit that same arm. Repeated decision reuse performs zero additional acquisitions.

**Gates/risk.** G. Main risks are scope escape, wrongly charging ignored answers as leftovers, and hiding incomplete product analysis.

### P6 — Strict fetch boundary, exact-decimal parity, and async emit

**Goal.** Implement production behavior with the settled HTTP contract and no SDK/runtime package dependency.

**Files.** Add `compiler/jev_emit.go`, `jev_http_emit.go`, and focused generated-decoder templates/helpers. Extend `emit.go` and applicable `emit_*` tests. Add `jev_wire_test.go`, `jev_parity_test.go`, and committed wire vectors.

Separate the implementation into reviewable sub-slices: ordered request encoding; duplicate-rejecting decoder; narrow fetch boundary; async business-flow lowering.

The emitted boundary must:

* Use the fixed adapter endpoint and request schema; no configurable credential destination.
* Read only `TYPESAFE_API_KEY` at request time; never reread config or load fixtures.
* Make at most one fetch attempt, reject redirects, apply AbortController timeout through body reading, and clean up the timer.
* Convert only specified boundary failures to the retained two intrinsic error kinds.
* Read raw response text and reject duplicate **decoded** keys before lossy object materialization.
* Check model, question membership/kind, support, selected Choice result, complete structured Score legend, and all numeric obligations—even for ignored answers.
* Perform exact decision arithmetic after the specified numeric ingress conversion; never normalize or substitute floating calculations for exact comparisons.

Do not wrap application continuation code in the transport catch. A compiler defect is not a successful execution of the user’s failure policy.

Use the actual `$can_kind` result ABI. Metadata/config providers must not accidentally become runtime function imports. Ordinary can-module imports remain ordinary imports; “dependency-free” means no SDK or external runtime-helper package.

Lower async functions and call sites only after validated `Jev.request` authority closure. Preserve synchronous output byte-for-byte for unaffected functions and programs.

**Diagnostics.** No new compile-time outage diagnostic. Runtime failures use the two retained intrinsic kinds; malformed source evidence retains its existing CAN rule.

**Tests/goldens.** Fake-fetch and fake-timer fixtures capture exact body, headers, endpoint, redirect policy, attempts, abort, and cleanup. Test absent credentials, timeout, non-2xx status, failed body read, malformed JSON, escaped duplicate keys, prototype-sensitive keys, wrong model/support/legend, and continuation failure isolation. No real key is required; a nonsecret sentinel tests header construction.

Shared Go/JS vectors must include independently specified expected values, not only agreement between implementations. Add the four Jev TS/catalog goldens and cross-module async parity.

**Done criteria.** The emitted program and compiler evaluator agree on the committed corpus; strict TS checks the actual generated helpers. No SDK or hand-written fake SDK surface is involved.

**Gates/risk.** G, with Node available to the Go-driven parity tests. The provisional host is server-side Node 24, already used by the repository’s TS job; official Node documentation supplies the required platform fetch and AbortController APIs.  ([Node.js][1])

### P7 — Ordinary-path activation, diagnostics, catalogs, and editor completion

**Goal.** Make the feature accepted through the same path in CLI, LSP, normalization, module checking, and emit.

**Files.** `compiler/main.go`, `lsp.go`, `check.go`, `catalog.go`, `normalize_golden_test.go`, `json_golden_test.go`, `diag_hint_test.go`, `lsp_wire_test.go`; `tools/modcheck/main.go`/tests; `tools/gramcheck/main.go`/tests; `editors/vscode/syntaxes/can.tmGrammar.json` and client integration where needed.

Activate public grammar only after P4–P6 are complete. Every new node must have explicit treatment in world construction, expansion, effects, revision identity, proof admission, evaluation, coverage, cataloging, normalization, linting, and emit. Unknown forms must fail closed rather than disappear in a default visitor case.

Replace modcheck’s `on`-substring boundary assumption with indentation-aware evidence-region handling. Update TextMate regions for singleton and batch syntax without swallowing `decide` or failure clauses into `given`. Add scopes and precedence tests for all selected metadata, primitive, pattern, observation, and scope forms; retain generics/string precedence.

Publish exact source locations, including config-file diagnostics and related dependency locations. Profile edits must invalidate affected open judgments. CLI and LSP must agree on missing selected baseline authority instead of retaining a stderr-only editor fallback.

For `errors.json`, preserve the array schema and reference meaning of `hit_by_tests`. Attribute one source `_` handler to **both** intrinsic kinds without manufacturing two coverage obligations or probability-mass hits. Load Jev catalog entries only into applicable worlds. Preserve `TestCatalogNoPhantomHits` in `compiler/lsp_test.go`.

**Tests.** `TestJevDiagnosticMatrix` must cover every new code, a legal neighboring case, precise UTF-16 spans, severity, Expected/Found/Hint, and CLI/LSP equality. Extend all-code ordering/uniqueness and explain tests. Add profile-edit invalidation, ordinary-wildcard negative controls, normalization, and catalog goldens.

**Done criteria.** No accepted source form requires a separate Jev command, editor-only exception, or skipped verification phase.

**Gates/risk.** G. Main risk is a split gate: especially JSON config errors, stale cached profiles, and syntax recognized by one tool but not another.

### P8 — Migration and removal of rejected artifacts

**Goal.** Leave one documented and executable surface.

**Files.** Add normative programs under `sketches/jev/`, root profile where appropriate, generated TS/catalog goldens, and README coverage. Update `docs/can-idioms.md`, a94/a95 selection notices, and compiler/editor documentation.

Inventory actual legacy matches before removing anything. Remove executable per-descriptor transport fields, old service/provider manifests, SDK imports/dependencies, old observation-record boilerplate, individual intrinsic handlers at native sites, and stale generated outputs **where they exist**. The inspected HEAD does not justify pretending that an old SDK integration has already shipped and been removed.

Keep historical decision documents as history, explicitly superseded. Keep rejection fixtures as rejection fixtures. Do not delete unrelated `auth.externs.ts` or `retry.externs.ts`, and do not turn retained historical prose into an accepted parser alternative.

Legacy label-only stubs require explicitly synthetic full observations or reviewed recorded evidence; never invent probabilities while calling the migration mechanical. Repin and obtain accepted baselines through the existing approval workflow.

**Diagnostics/tests.** Add negative fixtures for old `judgment`/`match judge`, `match judgments`, A inline grammar, singleton scopes, transport overrides, and mixed per-error/native failure handling. Add `TestJevLegacySurfaceRejected` and a scoped production-dependency inventory test.

**Done criteria.** No accepted fixture, generated artifact, runtime import, or current guide offers the rejected surface. Nonjudgment goldens remain byte-identical.

**Gates/risk.** G. Main risk is leaving migration scaffolding in production—or deleting legitimate non-Jev extern support while cleaning it up.

### P9 — Cross-boundary qualification and release evidence

**Goal.** Close the accumulated verification matrix through ordinary entry points.

**Files/suites.** Add `jev_hermetic_test.go`, `jev_strip_test.go`, and `jev_end_to_end_test.go`. Complete `compiler/testdata/jev/` corpora and record Q1/Q6 measurements in `docs/jev-revision-blast-radius.md` and the implementation record.

Run the provisioned toolchain with network access denied and real credentials absent. Compare deterministic artifacts across environment changes. Execute emitted runtime fixtures only through fake transport/timers. Ensure CI does not turn a missing Node/Z3 prerequisite into skipped mandatory evidence.

Audit production TS, generated declarations, source maps where produced, and packaging sidecars for fixture/test/pinned-premise sentinels. Runtime prompts, policies, schemas, and profile constants remain; evidence-bearing catalogs and accepted-baseline data are build/review artifacts, not production payload.

Complete the source mutation tests: deleted low arm, deleted low witness, both failure witnesses removed, malformed ignored answer, changed pinned premise with unchanged expected result, config repin without dependent revision, and contracted judgment incorrectly treated as pure.

**Done criteria.** Every checklist item below has a committed test and a recorded future execution result. Qualification is not complete merely because the four happy-path goldens compile.

**Gates/risk.** G plus the ordinary baselined CLI/LSP acceptance fixtures invoked by those suites. No live evaluation is part of this phase.

## 3. a95 §9 verification traceability

**New suite names below are proposed; existing suite/file names identify inspected integration surfaces.**

| a95 §9 obligation                            | Owning phases            | Concrete verification                                                                                                                                                                                 |
| -------------------------------------------- | ------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Emit goldens and unchanged ordinary output   | P0, P6, P8, P9           | Extend `emit_golden_test.go` with `TestGoldenJevNoul`, `TestGoldenJevChoice`, `TestGoldenJevScore`, `TestGoldenJevBatch`; preserve current AuthLogin/RetryLoop/Counter/Reserve/Notify/stdlib outputs. |
| Diagnosis and LSP suites                     | P1–P5, completed P7      | `TestJevDiagnosticMatrix`; `json_golden_test.go`, `diag_hint_test.go`, `lsp_test.go`, `lsp_wire_test.go`; independent success/failure completeness and ordinary-wildcard controls.                    |
| Registry and explain                         | P0–P7                    | `TestCodesUnique`, `TestAllCodesSequenced`, `TestAllDiagsCoded`, `explain_test.go`; every registered code has a violating and legal control.                                                          |
| Revision and accepted baselines              | P1, P4, P7               | `jev_revision_test.go`, `revision_identity_test.go`, `lsp_baseline_test.go`; full §2.2 matrix, format upgrade, authority rejection, structural identities, transitive request closure.                |
| Acceptance authority                         | P4                       | `acceptance_test.go`, `jev_acceptance_test.go`; CAN6017 unchanged; CAN3511 detects input/queue/request/reach changes without expected-output churn.                                                   |
| Contract verification                        | P3–P5, completed P7      | `contract_admission_test.go`, `contract_verify_test.go`, `contract_activate_test.go`, `jev_proof_test.go`; unsupported effects remain visible, fixture changes cannot prove truth.                    |
| Catalog, including no phantom hits           | P3, P5, P7               | `catalog.go`, existing `TestCatalogNoPhantomHits`, new `jev_catalog_test.go`; one source handler, two handled kinds, unchanged reference semantics.                                                   |
| modcheck                                     | P1–P5, completed P7      | `tools/modcheck/main_test.go`; profile/declaration resolution, proper evidence-region boundaries, state signatures, no default profile, no `-`.                                                       |
| gramcheck and editor grammar                 | P2, P3, P5, completed P7 | `tools/gramcheck/main_test.go`, `can.tmGrammar.json`; selected syntax only, nested region termination, string/generic precedence.                                                                     |
| Strict TS gate                               | P6–P9                    | Actual generated TS under `tscheck/tsconfig.json`; existing pinned checker; no SDK stubs, suppressions, or runtime package dependency.                                                                |
| New HTTP/metadata/numeric conformance corpus | P1, P2, P6               | `testdata/jev/observations/` and `wire/`; all structure positions, null rules, ordered keys, duplicates, model/support/legend, numeric ABI, exact comparisons.                                        |
| New source-coverage corpus                   | P3, P5, P9               | `jev_coverage_test.go`; five source obligations for the eight-row §8 specimen; source alternatives, ties, uncertainty, missing arms versus missing witnesses.                                         |
| New evidence-lifecycle corpus                | P3, P5, P9               | `jev_evidence_test.go`; omitted/unvisited, omitted/reached, exhaustion, wrong args, differing repeats, partial/zero leftovers, helper/root isolation, ignored answers.                                |
| New hermeticity and stripping corpus         | P6, P9                   | `jev_hermetic_test.go`, `jev_strip_test.go`; network denied, no real credentials, fake runtime boundary, deterministic artifacts, zero fixture bytes in production.                                   |
| New scope/request-count corpus               | P5, P6, P9               | `jev_scope_test.go`, `jev_end_to_end_test.go`; nested `decide` adds no request, later `judge` owns its exchange and failures.                                                                         |

The existing JSON diagnostic suite contains the registry-order and uniqueness tests, so those are extensions to a real suite rather than a proposed replacement registry mechanism.

## 4. a95 §10 assumptions and overturning evidence

### Q1 — Revision blast radius

**Assumption:** retain the full reachable `Jev.request` contract closure. A global profile change is intentionally a potentially broad revision event. Do not narrow it merely to reduce edits.

**Required representative measurement in P4:** create a seven-module fixture with three primitive declarations, one mixed batch, two singleton entry functions, one batch router, two workflow functions, one top-level entry, and a pure result-consuming module. Include a diamond dependency, a same-file helper, an unused question, and repeated descriptor sites.

For model, timeout, one prompt, one threshold, batch membership, a pure body, and a fixture-only mutation, record:

| Measurement                                       | Required distinction                                          |
| ------------------------------------------------- | ------------------------------------------------------------- |
| Declarations whose effective fingerprints change  | Actual reachable closure versus unrelated declarations        |
| Revisions and `uses` pins requiring edits         | Necessary changes versus collateral churn                     |
| Pinned-premise warnings                           | Evidence churn versus interface churn                         |
| Diagnostic sequence                               | Direct provider/pin fault versus derivative identity failures |
| Review artifact size and changed source locations | Mechanical edits versus policy decisions                      |
| Generic-instantiation effects                     | Original declaration identity versus generated instances      |

This is a **future measurement**, not a benchmark claimed here. Its report is a P4/P9 acceptance artifact.

**Overturning evidence:** a concrete application demonstrates that the closure forces unrelated contracts to change, or that a narrower definition preserves the same pinned authority with materially less burden. The remedy is a separately tagged identity amendment, not an implementation shortcut.

### Q2 — Failure abstraction

**Assumption:** one terminal clause remains sufficient for first-class judgment control flow. Keep the two intrinsic kinds internally and in catalogs; do not expose a second detailed-handler grammar.

**Overturning evidence:** a real application requires meaningfully different recovery for protocol failure and unavailability, with an exhaustive alternative that preserves the bounded failure abstraction. That evidence can commission a later design; it does not reintroduce per-error arms in this implementation.

### Q4 — State and disclosure

**Assumption:** dynamic state leaves are unbranded text only. Closed metadata has its separately specified scalar domain. No inferred conversion or unsealing enters this work.

**Overturning evidence:** a concrete consumer cannot express necessary input through explicit existing text preparation, and supplies an exact representation/disclosure proposal. Broader state remains a separate design.

### Q5 — Host and wire conformance

**Assumption:** qualify server-side Node 24 first, using a generated duplicate-rejecting decoder and no SDK. Record the exact runtime/toolchain used for qualification. Do not advertise browser or other host support from TypeScript type-checking alone.

**Overturning evidence:** committed raw-wire or host fixtures show an admitted structured legend, numeric value, timeout/body-read behavior, or required platform operation cannot be represented correctly. Correct the adapter implementation—or amend an incompatible ABI assumption explicitly—before shipping. Live credentials are not needed to demonstrate these boundary failures.

### Q6 — Batch/checker scale

**Assumption:** exact finite analysis with an explicit deterministic budget is adequate for the initial applications. No partial analysis is acceptance.

**Required evidence:** P5/P9 include the seven-module fixture’s largest decision, wide Choice products, multiple ignored fields, alternatives, and nested decisions. Record analysis work and resource usage, not just successful exit status.

**Overturning evidence:** realistic required batches consistently exceed the budget. Optimize exact analysis or revise the declared implementation limit; never sample cells, infer independence, or auto-split the request.

**Q3 is closed by the owner’s decision. Q7 operational calibration, labels, recording retention, and live evaluation remain outside this implementation plan.** These dispositions follow a95 §10 without reopening the selected grammar.

## 5. Risks, in blast-radius order

| Rank  | Risk                                                                                                | Mitigation and release condition                                                                                                                             |
| ----- | --------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **1** | **Q1: request-authority changes either propagate too little or cause uncontrolled revision churn.** | P4’s measured multi-module closure; direct/dependent diagnostic precedence; explicit owner acceptance of the measured burden.                                |
| **2** | New AST nodes disappear from expansion, effects, proof, identity, or emit walks.                    | Explicit visitor cases; generic/local/cross-module negative controls; no public activation before P7.                                                        |
| **3** | Runtime decoding/routing differs from compile-time evidence.                                        | Shared vectors with independently specified answers, strict duplicate-key parsing, exact-decimal parity, full ignored-answer validation.                     |
| **4** | Failure `_` becomes an accidental catch-all or catches continuation errors.                         | Independent success exhaustiveness; narrow acquisition boundary; deleted-low-arm and faulting-continuation tests.                                            |
| **5** | Baseline/editor behavior creates a falsely clean review.                                            | Same selected authority for CLI/LSP; visible unenforced status; missing selected baseline rejects acceptance; config-file spans and dependency invalidation. |
| **6** | Scope reuse makes extra requests or misaccounts scripts.                                            | One acquisition/one exchange invariant; explicit request counters; nested `decide`, repeated-site, helper, and zero-leftover tests.                          |
| **7** | New intrinsics or async lowering change unrelated output.                                           | HEAD golden inventory; conditional Jev-world/catalog inclusion; synchronous controls; retain current `$can_kind` ABI.                                        |
| **8** | Migration leaves a second grammar or removes legitimate extern code.                                | Scoped accepted-source/import inventory; legacy rejection tests; historical docs preserved but marked superseded.                                            |
| **9** | Tests accidentally require the service, credentials, or optional host tools.                        | Denied-network qualification, nonsecret boundary sentinels, fake fetch/timers, mandatory CI prerequisites, no skipped conformance checks.                    |

## 6. Sizing and parallelism

Sizes describe relative implementation/review effort, **not elapsed-time estimates**.

| Phase | Size    | Reason                                                                                                                |
| ----- | ------- | --------------------------------------------------------------------------------------------------------------------- |
| P0    | **S**   | Selection ledger and two bounded integration clarifications; no runtime mechanism.                                    |
| P1    | **M**   | New provider kind, strict config parsing/canonicalization, and three-tool resolution agreement.                       |
| P2    | **L**   | New structured literal contexts, typed observation shapes, numeric validation, and expansion-safe AST representation. |
| P3    | **L**   | Three primitive surfaces plus effect, queue, exhaustive-state, coverage, and lint integration.                        |
| P4    | **L**   | Transitive identity, accepted-baseline format, premise authority, active proof refusal, and Q1 measurement.           |
| P5    | **L**   | New acquisition scope, observation availability, nested decisions, product analysis, and request accounting.          |
| P6    | **L**   | Security-sensitive parser/HTTP boundary, exact numeric parity, and cross-module async lowering.                       |
| P7    | **M–L** | Many existing entry points; complexity is consistency rather than new semantics.                                      |
| P8    | **S–M** | Primarily normative fixtures/docs at inspected HEAD; increases with actual legacy consumers found.                    |
| P9    | **M**   | Cross-boundary qualification; most component corpora must already exist.                                              |

**Can parallelize:** after P2 fixes the resolved schema and observation ABI, transport/decoder work can proceed alongside singleton checking. P4 identity/proof work and P5 scope work can proceed in parallel after P3 establishes source sites. Editor coloring and diagnostic fixtures can be prepared with each grammar slice.

**Strictly ordered:** profile/schema identity before accepted baselines; validated observation semantics before trusting parity; source-site/source-arm identity before scope expansion; identity, proof, scope, and emit completion before activation; migration before final production inventory/stripping acceptance.

Do not distribute simultaneous uncoordinated edits to `parse.go`, `check.go`, `lsp.go`, or `revision.go`. Assign an integration owner for shared AST and gate changes; parallel work should consume the agreed structures rather than create competing ones.

## 7. Acceptance checklist — “surface usable”

* [ ] **Canonical syntax:** all three B declarations and singleton forms work; only 2+ question B descriptors use `judge`/`decide`. A forms, `match judgments`, singleton scopes, and old a94 syntax reject.
* [ ] **Complete metadata/state surface:** every admitted instruction/criterion/level structure and `unspecified` position is tested; runtime state remains text-only; no implicit capture, conversion, or unsealing exists.
* [ ] **Pinned transport:** config/provider resolution is identical across tools; aliases, overrides, invalid timeout, nonzero retries, and incompatible profiles reject; runtime reads only the credential.
* [ ] **Identity and acceptance:** the full §2.2 mutation matrix passes against accepted baselines; Q1 measurements are committed; CAN3511 and CAN6017 remain distinct; generation is never acceptance.
* [ ] **Evidence lifecycle:** full observations are mandatory; wrong/missing/exhausted/unconsumed exchanges fail correctly, including zero visits; identical requests may receive different scripted observations.
* [ ] **Coverage and scope:** actual source arms receive hits; failure expansion does not invent obligations; uncertainty/ties/low paths remain explicit; ignored malformed answers fail; `decide` reuse performs no request.
* [ ] **Proof separation:** contracted unsupported judgment effects reject through the active verifier; changing fixtures, confidence, or pinned status cannot establish a proof fact.
* [ ] **Wire and emit:** exact request, duplicate-key, model/support/legend, numeric, timeout, redirect, one-attempt, and continuation-boundary vectors pass on the qualified host; emitted code uses no SDK/runtime package.
* [ ] **One gate and diagnostics:** every selected code has registry/explain/structured-payload/span coverage and CLI/LSP parity; modcheck, normalization, lint, grammar, config invalidation, and catalogs agree.
* [ ] **No collateral output changes:** every nonjudgment HEAD golden remains byte-identical; legitimate extern support remains; pure functions do not become async or acquire Jev catalog entries.
* [ ] **No migration residue or evidence shipment:** accepted examples have one surface; production TS/maps/packages contain no tests, scripts, fixture sentinels, or pinned-premise data.
* [ ] **Recorded qualification:** G and the committed offline end-to-end suites have actual implementation-stage results, with no skipped mandatory checks and no live model dependency. A green result establishes conditional program behavior—not model calibration or threshold safety.

[1]: https://nodejs.org/download/release/v24.4.0/docs/api/globals.html?utm_source=chatgpt.com "Global objects | Node.js v24.4.0 Documentation"
