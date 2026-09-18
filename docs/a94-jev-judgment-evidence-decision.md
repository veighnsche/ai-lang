# JEV-1.0 — Hermetic judgment evidence and deterministic routing

**Status:** ratification candidate; one decision, not implementation authorization.  
**Repository reviewed:** `veighnsche/can-lang`, commit `4d447bac858c50d998804aa7fa96f2b66b3b557b`.  
**External references reviewed:** live TypeSafe and Probably documentation, 18 September 2026.  
**Execution status:** no compiler changes, model calls, or compiler gates were run. The proposed syntax below is new syntax, not a claim that the reviewed compiler accepts it.

## 1. Decision

**Ratify a hybrid whose compiler semantics are full-observation scripting: each reached Jev request consumes an ordinary, request-matching `given` exchange containing the complete logical response; canlc then validates its protocol shape and executes the same deterministic confidence/selection policy as production.** Choice scripts include the reported label, complete probabilities, and reported confidence; Score scripts include the reported score, complete probabilities, and reported confidence; Noul scripts contain only P(yes). Scripts are counterfactual environment observations, never model-truth assertions, calibration certificates, or proof premises. All builds and editor checks remain offline and credential-free. Recordings may be explicitly converted into these source rows, but are not another compiler input format. Live evaluation is a separately authorized operational report, never a substitute for the compiler gate. Sampling is absent from JEV-1.0. **My confidence in this architectural decision is 0.93**, a subjective assessment, not a measured probability. A published, versioned confidence algorithm with executable conformance vectors could justify deriving confidence instead; substantial real-world evidence that source exchanges cannot preserve necessary request/response identity could justify an external fixture format. Neither would justify live builds or treating model estimates as proof facts.

**Meaning of a passing row:** “For this exact request, if the environment supplies this admissible observation, this implementation takes this path and produces this result.” It does **not** mean “Jev will return this observation,” “the selected label is correct,” or “this threshold is safe.” This preserves the boundary established by [R7/R8][R], [a12 §Theorem][A12], and [a81 §Obligations][A81].

## 2. Source findings that constrain the decision

### 2.1 Read the amendments, not just the historical headings

| Source receipt | Operative finding |
|---|---|
| `REQUIREMENTS.md`, Goal, R5–R8, R10–R11 | Explicit declarations; hermetic decision tables; request/response exchanges; producer-owned error bounds; one behavioral verifier; TypeScript output. [R] |
| `docs/a03-branch-coverage.md`, §Hit semantics and retrospective amendment | Coverage is observed execution or the specifically authorized structural certificate. A certificate is never a fabricated hit. The historical stale-emits requirement is superseded. [A03] |
| `docs/a12-contracts.md`, §Rule / §Theorem | Scripts bind arguments and permitted responses. A caller cannot manufacture producer honesty. Unrealized declared errors are allowed. [A12] |
| `docs/a81-proof-obligations.md`, §Obligations / §Deliberate decisions | Scripts never become proof facts. A scripted execution and a universal proof can deliberately describe different reachable-outcome assumptions. [A81] |
| `docs/a87-acceptance-authority.md`, §Semantics / §Interplay | `pinned` changes acceptance authority, not execution or proof. CAN6017 warns after a clean gate, does not reject, and currently does not protect `given` scripts. [A87] |
| `docs/a14-tsc.md`, §Reinstated (a70) and checker-upgrade note | `tsc` has been reinstated for output type/contract **shape**. canlc alone owns behavioral proof, termination, exhaustiveness, tables, and coverage. Do not revive the earlier “no tsc gate” state. [A14] |
| `docs/a77-revision-identity.md`, §Fingerprint model / §Enforcement | Interface identity, body-proof freshness, and test evidence are distinct. Ordinary bodies and `given` are not interface fingerprints. Prompt identity therefore needs an explicit amendment. [A77] |
| `docs/a91-given-omission-amendment.md`, §S0 rulings / follow-up | Tables are partial. Omission claims non-reach; reaching without a script fails. CAN3105 is retired; `-` is CAN3111; CAN3420 is an advisory static-unreachability lint. [A91] |
| `sketches/auth-login/auth.can`, `auth__login` / `auth__verify` | The normative example follows the existing one-line exchange shape and same-file helper/reaching-test organization, rather than inventing a separate test map. [AUTH] |

The requested `.md` TypeSafe endpoints were unavailable through the retrieval interface. I read the corresponding live HTML pages identified by the official [documentation index][TS-I], including the JavaScript API references. No local skill snapshot substitutes for those live pages.

### 2.2 TypeSafe behavior ledger

Quoted fragments below are short excerpts; the right column states their consequences. Proposed can-lang policy is identified as such and is not attributed to TypeSafe.

| Live page and short quotation | Behavior relied on / consequence |
|---|---|
| [Primitives][TS-P]: “Question IDs are for your code. They are not sent to the model.” | Request questions have named response slots. IDs still belong to can-lang schema identity even though they are not prompt text. |
| [Choice, response structure][TS-C]: “The option with the highest probability.” “The sum of all values is 1.” | Preserve the reported selection and full support. Validate that the selected label maximizes the returned distribution. **Routing ties to low is our rule**, not a claimed SDK tie guarantee. |
| [Noul][TS-N]: “Noul does not return a separate confidence value” | Script P(yes), not an invented confidence field. Our uncertain interval is a probability policy. |
| [Score][TS-S]: “It’s each level number multiplied by its probability, added up” | Score is an ordered-level weighted mean, not the most likely level. Preserve the full distribution so different shapes with the same mean remain distinguishable. |
| [Confidence][TS-F]: “a statistic computed from the probability distribution”; “Where you draw those boundaries depends on the stakes.” | TypeSafe reports confidence for Choice/Score, but this page supplies no numerical algorithm to reproduce it. Preserve it as reported data; thresholds remain application policy. |
| [How to build with System One][TS-H]: “Keep control flow, deterministic rules, and side effects in code.” | The compiler implements the policy, not the model. Model observations cannot certify authorization, business predicates, or side effects. |
| [Fan-out][TS-B]: “All questions are evaluated in parallel” | One explicit batch is one request containing all independent questions. Ignoring an answer after receipt is allowed. Parallel evaluation is **not** evidence of statistically independent errors. |
| [Confidence routing][TS-R]: “some actions are riskier than others and thus demand a higher confidence threshold.” | Risk policy must be reviewable source. The example thresholds in this decision are routing-test values, not calibrated recommendations. |
| [JavaScript SDK][TS-J]: “Answer types are inferred from your questions.” | Emit the actual `@typesafe-ai/sdk` client and typed question constructors; do not invent an SDK API for confidence routing. |
| [TypeSafeClient][TS-JC]: “The request cannot connect or times out after retries.” | SDK failures must be translated at the foreign boundary into declared can-lang errors. No SDK exception escapes into `.can` control flow. |
| [Client configuration][TS-JF]: “Explicit values take precedence over environment variables, then SDK defaults.” | Explicit model, endpoint, timeout, retry, and logging policy prevent environmental defaults from silently changing those behaviors. Credentials remain runtime-only. |
| [RetryPolicy][TS-JY]: “Maximum retries after the initial attempt; `0` disables retries.” | Emit `maxRetries: 0`. A source-level retry is another explicit judgment invocation, not an invisible SDK attempt. |
| [RequestOptions][TS-JO]: “Timeout per attempt in milliseconds; there is no total retry budget.” | Declare a finite timeout and disable automatic retries. This bounds the configured wait policy, not the number of inferences a remote server might internally perform. |
| [Models][TS-M]: “An alias moves when a new release ships”; “The response’s `model` field reports the versioned ID that answered” | Require a full versioned model ID and validate the response ID. The current listed ID used below is `jev-1.13.0`, not `jev-1.13` or `jev-latest`. A version string is not itself a proof of immutable service behavior. |

## 3. Rationale against the eight criteria

### 3.1 Hermetic, deterministic, credential-free builds

The compiler evaluates only source inputs, committed source exchanges, and deterministic adapter/routing rules. Missing scripts, unavailable recordings, absent credentials, and an offline machine must never trigger a model call or a fallback result. Installing locked tool dependencies is distinct from executing decision tables; the latter needs no SDK client, network, clock, or randomness. This applies to CLI builds, normalization that evaluates rows, and LSP checks. **A live provider is not a build dependency.** Receipt: [R7/R10][R], [a12 §Theorem][A12].

### 3.2 Proof/evidence separation

The proof domain contains every admitted observation, independently of the samples in `given`. An author can test a conservatively admitted observation without claiming that the provider actually produces it. Confidence values and label strings do not establish propositions about the input. A successful routing test therefore cannot prove “this message is billing-related,” even at confidence 1.

When a universal contract is requested, the first implementation must explicitly reject judgment/effect constructs that the contracts verifier cannot soundly model, using its unsupported/unverified-dependency machinery; it must not silently omit their obligations. Ordinary shape, effect, termination, exhaustiveness, coverage, and decision-table checks remain mandatory. Later proof support can verify the deterministic routing relation against symbolic admissible observations, not against fixture frequencies. No new probabilistic theorem prover is required for this decision. Receipt: [a81 §Obligations / §Findings][A81], [a77 §Fingerprint model][A77].

### 3.3 Test-per-arm coverage

A source arm is covered only when a green decision-table execution selects it, or when it qualifies for the existing narrowly authorized structural certificate. Positive probability mass on a label is not a hit. A thousand recorded high-confidence responses do not cover low confidence. Conversely, a synthetic low-confidence response really can test low-path business logic without proving that the current model will take that path.

Compiler-generated validation and dispatch scaffolding has compiler conformance tests. It must not create hidden user coverage obligations or fabricate coverage for source arms. Track source-arm identities through elaboration. Receipt: [a03 §Hit semantics / retrospective amendment][A03], [R10][R].

### 3.4 Amendment and production-strip discipline

Add the rules under the explicit tag **JEV-1.0**, including the necessary R2/R4/R6/R10/R11 companions; do not pretend an R8-only patch explains new effects, declarations, or async emit. New violations receive unique registered CAN codes. Old unrelated goldens remain unchanged. All tables and fixture payloads remain build-side; the emitted program retains only prompts, schemas, policies, adapter logic, and business behavior. Receipt: [R freeze / R7 / R10 / R11][R], [compiler/code.go][CODES].

### 3.5 Agent-authored UX

There is one source form for judgment evidence: `exchange args (...) outcome ...`. There is one request form for both one and several questions: a named judgment descriptor with an explicit answer record. No confidence-band stubs, inferred probabilities, implicit normalization, implicit model alias, automatic batching, fallback-to-live switch, or alternate cassette syntax.

The verbosity has a purpose: an agent can review the request contract, its policy, and the counterfactual observation separately. Repeated test states with different observations are legal and often essential. Receipt: [Goal / R8][R], [auth-login exchange sites][AUTH].

### 3.6 Calibration honesty and drift

The compiler checks the implemented threshold relation, not the appropriateness of the threshold. Retaining both probabilities and reported confidence makes this boundary visible and supports future distribution-shape checks without reconstructing missing evidence.

Live quality evaluation must use separately reviewed expected labels or outcomes. **Do not treat pinned routing rows as model ground truth:** the same state may deliberately be scripted as billing, other, uncertain, and unavailable. A report should distinguish observed error rate, coverage/abstention rate, confidence calibration, model/request version, sample population, and uncertainty. Lack of a live run is “not measured,” never green. Receipt: [a87 §Interplay][A87]; external basis: ledger [TS-F], [TS-R], [TS-M].

### 3.7 Cost, latency, credentials

Builds and editor checks spend zero model-call budget. Production requests are paid by the deployment account; an explicitly invoked live-evaluation job is paid by its designated account and budget. Recording, evaluation, and production must never source credentials from committed `.can` or fixture data. The adapter disables browser credential exposure and body logging; recordings require a separate data-handling review.

Batching reduces the number of requested round trips; it is not a promise of free extra questions. A configured timeout is an operational limit, not a mathematical termination proof for the remote service. Receipt: [R7/R8][R]; SDK configuration/retry/timeout and fan-out behavior: [TS-JF], [TS-JY], [TS-JO], [TS-B], [TS-C].

### 3.8 Migration

Preserve the old extern contract until a judgment replacement has an explicit descriptor and full observation fixtures. A legacy stub returning only `Ok(label = "billing")` cannot be mechanically upgraded by inventing probabilities or confidence. Its routing intention must be rewritten as an explicitly synthetic observation or supplied by a reviewed recording; these are different provenance claims.

Keep domain outputs and errors stable where possible. Declare the new request capability in callers, update explicit pins where interface authority changes, remove replaced extern imports, and freeze old-versus-new behavior over the reviewed fixture corpus. Do not retain two spellings for the same first-class operation after that module migrates. Receipt: [R2/R6/R8][R], [a12 §Consequences][A12], [a77 §Enforcement][A77].

## 4. Normative amendment text

### 4.1 JEV-1.0/R7 — what a judgment test means

**MUST:** A judgment in a compile-time test obtains its response exclusively from the selected call-site exchange queue. No compiler, editor, or normalization mode may call the model, consult a recording file, read runtime credentials, or generate a random draw to complete a row.

**MUST:** Judgment exchange outcomes are closed literal values (including already-supported typed constant expansions), not computations that manufacture probabilities or confidence. Validate each scripted successful observation against the descriptor and the adapter rules below before executing rows. Malformed successful evidence is a diagnostic, not an injected service failure. To test failure handling, script the declared error value explicitly.

**MUST:** Execute validation/selection/confidence policy before choosing a source arm. Scripts cannot directly supply a selected confidence arm.

**MUST:** Assess coverage over the executed source arms of green tables, retaining the existing structural-certificate distinction. Probabilities, replay counts, recorded provenance, and pinning confer no extra coverage credit.

**MUST NOT:** Use scripts, pinned status, reported confidence, or live-evaluation results as proof premises. A high-confidence label is a model estimate, never a proof of the natural-language predicate.

**MUST:** Strip tests, exchanges, fixture observations, and their identifying strings from shipped production artifacts. Build reports containing evidence remain outside the production package.

### 4.2 JEV-1.0/R8 — exchange identity and consumption

A judgment site is a new effectful site class, including when its descriptor is in the same file. It is **not** a deterministic same-file helper and must not inherit the prohibition on `given` for such helpers.

Each exchange binds all dynamic descriptor arguments by name and equality, and supplies either the complete successful answer record or one producer-declared error construction. Static request identity comes from the referenced descriptor and its revision fingerprint; authors do not duplicate the prompt in every exchange. At runtime, the compiler serializes exactly that descriptor and those evaluated arguments.

Tables remain partial. The following outcomes are normative:

| Situation for the current root test and site | Ruling |
|---|---|
| No queue and no visit | Legal non-reach; no evidence for that site. |
| Visit with no queue or an exhausted queue | CAN4200; never a live call. |
| Wrong argument name/value or unexpected argument | Existing argument/exchange checks and CAN4200 for execution-time mismatch. |
| Queue has N exchanges; site visited N times | Consume in order, one exchange per invocation. |
| Queue has N > 0 exchanges; site visited fewer times, including zero | CAN4200: leftover evidence. Initialize the audit before the first possible visit, so zero-consumption is not invisible. |
| Valid test name with no static path to site | Retain CAN3420 as an advisory. It never turns unused evidence into coverage or exempts applicable leftover checks. |
| Unknown test key | CAN3104 error in CLI, LSP, and modcheck under the explicit severity amendment below; do not reintroduce CAN3105. |
| `key => -` | CAN3111, as under a91. |

The queue audit covers the root function and its structurally reachable same-file helper bodies, stopping at foreign scripted-call boundaries. Only that root test's entries participate. It must not charge one function for unrelated functions' queues. This explicitly specifies zero-consumption for the new judgment class; it is not a claim that every older evaluator path already implements that audit.

Equal requests may have unequal observations, at the same site or different sites. Do not memoize by arguments or apply deterministic-provider body checking to a judgment. Calling a real same-file helper still executes that helper; a judgment inside it consumes the reaching root test's queue. Existing same-file and foreign-call distinctions otherwise stay intact. Basis: [R8][R], [a12 §Theorem][A12], [a91][A91].

### 4.3 JEV-1.0/OBS — the observation contract

“Raw” below means the complete **logical model answer**, not an untyped HTTP body. SDK envelope fields, such as usage, are not business-logic evidence. The adapter checks the envelope's requested model identity separately.

| Primitive | Required scripted successful fields | Dispatch |
|---|---|---|
| Choice | `choice: str`, every declared option's probability, `confidence: dec` | Validate support and maximizer; ties route low; otherwise classify reported confidence. |
| Score | `score: dec`, one probability per ordered level, `confidence: dec` | Validate weighted-mean consistency; classify reported confidence. Never replace the score by argmax. |
| Noul | `noul: dec` | Classify P(yes) using explicit no/yes bounds; no confidence field or synthesized confidence. |

All response questions and probability support entries must appear exactly once, with no extras. Question type tags and Score legend must match the descriptor at the wire boundary; they are declaration-derived schema, not independently editable fixture fields. A malformed or missing answer invalidates the **whole** response; there is no partial-success branch in JEV-1.0.

**Numeric ABI, chosen here—not claimed to be a TypeSafe formula:**

- Source scalars are exact `dec` literals or their closed typed-constant expansions. At SDK ingress, each finite JavaScript number is converted through its shortest round-trip decimal rendering, expanding exponent notation into canonical can-lang decimal form. Signed zero becomes zero.
- A fixture scalar must be in that same number-to-decimal image. An over-precise decimal that the SDK boundary cannot preserve is rejected, not silently rounded. This restriction is specific to observation fixtures, not general `dec` arithmetic.
- Probabilities and reported confidence lie in [0,1]. Let ε = `d"0.000000000001"`. A probability vector must have `abs(sum(p)-1) <= ε`. **Do not renormalize, clamp, or fill omitted entries.** ε is a versioned serialization-consistency allowance, not a calibration claim and not a threshold-comparison tolerance.
- Choice's reported label must be one of the exact maximizers of the canonical values. Any multiple maximum is a tie. No epsilon-based winner selection.
- For Score with K declared ordered levels, validate `0 <= score <= K-1` and `abs(score - sum(i*p_i)) <= K*ε`; preserve the reported score. Level order belongs to descriptor identity.
- Confidence receives range validation, **not** a fabricated `confidence = f(probabilities)` check. The reviewed documentation supplies no such executable formula. Accepted fixtures are therefore a conservative protocol model, not a certification of provider-realizable confidence/distribution pairs.

The compiler evaluator and emitted adapter must share conformance vectors for this ABI. A malformed success fixture fails before row execution; a malformed live success becomes `jev.invalid_response(reason = "protocol")`. This separates test-author mistakes from production boundary failures. Existing exact-decimal representation is documented in [R10/R11][R] and [compiler/eval.go, decimal comparison][EVAL].

### 4.4 JEV-1.0/ROUTE — deterministic policies

Choice and Score declare exact constants `medium = M` and `high = H`, constrained by `0 < M < H <= 1`. Classification is **high when c >= H; medium when M <= c < H; low when c < M**. For Choice only, a probability tie takes low before these comparisons, even if the reported confidence is high. High Choice states retain their selected label. Medium/low carry the complete observation but do not automatically act on its label.

Noul declares `no = L` and `yes = U`, with `0 <= L < 0.5 < U <= 1`. Classify **no when p <= L; yes when p >= U; uncertain otherwise**. Exactly 0.5 is uncertain. This is intentionally not the Choice/Score confidence statistic.

Policies have no defaults. Source must handle all declared service-error outcomes and cover the successful answer-state space. Ignoring one question with a field wildcard is permitted; omitting the low/uncertain combinations from an otherwise discriminating match is not. Arms select top-to-bottom after all answers in that one response have been obtained and validated. Coverage is per source arm, not per product cell.

No probabilistic sampling, seeded or otherwise, exists in these expressions. Ordinary comparisons against an observation's fields remain ordinary program logic and have their ordinary test/coverage obligations. The compiler does not infer that high means “safe to authorize.”

### 4.5 Companion declaration, effect, identity, and acceptance rules

**R2/R6:** Introduce a versioned `judgment` descriptor and `match judge` expression. Descriptors are provided/used and revision-pinned like other declarations, but have no executable body or independent function tests. Their sites are exercised by enclosing functions' tables. `Jev__Service@1` is a compiler-distributed, locked intrinsic provider manifest—not an unresolved ordinary module. It owns `jev.unavailable(reason: str)` and `jev.invalid_response(reason: str)`. The manifest and descriptor schema fix the producer's error bound; a consumer cannot shorten it to evade handling.

The new `Jev.request` effect capability must be explicitly declared on each judgment-using function and its transitive callers. Ordinary functions stay synchronous and pure where they were before. JEV-1.0 initially accepts **text state**, literal instructions/criteria, finite declared question sets, a full versioned model ID, a finite positive timeout, and `retries 0`. More state representations are a separate surface, not implicit serialization of arbitrary records or branded secrets.

**R4/a77:** A judgment's request schema, ordered questions/options/levels, decoded instruction and criterion strings, routing thresholds, service/adapter ABI, model ID, timeout, and retry policy are interface-identity fields. **Rewording those strings without a descriptor revision bump is CAN6013.** Comments or equivalent literal encodings are not rewording. Changing runtime input data is not a revision change. Ordinary state-preparation code retains a77's body/freshness rules; this decision does not pretend the compiler can detect an English instruction hidden inside arbitrary input data. Evidence changes rerun tests, not API revisions. Baseline format expansion is versioned and requires the existing explicit acceptance workflow. [A77]

**R10, explicit severity convergence:** a91 records an existing discrepancy: unknown script keys are rejected by modcheck but CAN3104 is a compiler warning. JEV-1.0 promotes CAN3104 to error in the common CLI/LSP path, for ordinary and judgment sites alike. Its rule and number do not change. This is an expressly tagged diagnostic behavior change, with the affected dead-script diagnosis golden updated; unrelated golden output stays fixed. It prevents the new evidence surface from inheriting a split gate. [A91]

**R7/a87:** Preserve CAN6017 exactly for expectation weakening/removal/demotion. Add a distinct, non-blocking **pinned judgment premise-drift warning**. For a pinned root test, the accepted baseline additionally fingerprints its test inputs, the statically reachable judgment descriptors, and applicable call-site exchange queues through same-file helpers—including exchanges that determine whether the judgment is reached. Identify sites structurally, not by source line number. Changing or removing this evidence bundle warns even if the pinned expected output is unchanged. Generated baselines remain unaccepted; warnings occur only after the ordinary clean gate. This closes the fixture-premise hole without silently upgrading A-light to rejection or treating pinning as truth. [A87]

### 4.6 Diagnostic allocation

The following numbers are proposed in the unused CAN35xx block of the reviewed registry. Ratification reserves them; implementation must confirm no intervening allocation. Each code denotes one rule. Existing diagnostics are reused where their rule genuinely applies. [CODES]

| Code | Rule / severity |
|---|---|
| CAN3501 | Judgment declaration/answer does not match its complete declared protocol schema. Error. |
| CAN3502 | Probability outside [0,1]. Error. |
| CAN3503 | Probability mass outside the versioned serialization allowance. Error. |
| CAN3504 | Choice selection is not a maximizer of its supplied distribution. Error. |
| CAN3505 | Reported confidence outside [0,1]. Error. |
| CAN3506 | Score violates its declared weighted-mean/range consistency relation. Error. |
| CAN3507 | Invalid routing partition bounds. Error. |
| CAN3508 | Observation scalar is not preserved by the specified SDK numeric boundary. Error. |
| CAN3509 | Judgment requests automatic retries; JEV-1.0 requires zero. Error. |
| CAN3510 | Model selection is not an explicit full versioned ID. Error. |
| CAN3511 | Pinned judgment premise bundle changed or disappeared. Warning after clean gate; never blocks emit. |
| CAN3512 | Missing or nonpositive/nonrepresentable finite request timeout. Error. |
| CAN3513 | Judgment exchange outcome is not a closed literal value. Error. |

Schema checking owns missing/extra fields, wrong question kind, forbidden Noul confidence, and mismatched support/legend. Scalar type mismatches remain CAN6003; malformed exchange syntax remains CAN3109. CAN3101, CAN3111, CAN4101, CAN4107, CAN4200, CAN6013, CAN6017, and effect-authority diagnostics retain their established meanings. No diagnostic is called “uncalibrated”: canlc has no evidence to make that judgment. Every new rule needs `canlc explain`, structured Expected/Found/Hint data, precise spans, and CLI/LSP parity. [R10][R], [CODES]

## 5. Normative `.can` example

This is the proposed canonical grammar, not pseudocode for an extern wrapper. The descriptor's return record declares the **successful observation schema**. `match judge` validates and classifies it before matching `on answers` patterns; there is no raw `on Ok` escape from that dispatch. Each answer field is named explicitly, including in a one-question request. Ordinary record types remain data types, not truth-bearing brands.

All observations below are **synthetic**. In particular, the two point-mass answers for the identical state test opposite routing possibilities; they do not both claim to describe the true subject of that message. The pinned rows pin application behavior conditional on those observations.

```can
mod triage
  provides [triage__questions, triage__route, Triage__Probabilities, Triage__Topic, Triage__Answers, Triage__Route]
  uses [Jev__Service@1]
  emits [jev.unavailable, jev.invalid_response, triage.unavailable]

error triage.unavailable(reason: str)

type Triage__Probabilities rev 1 (
  billing: dec
  other: dec
)

type Triage__Topic rev 1 (
  choice: str
  probabilities: Triage__Probabilities
  confidence: dec
)

type Triage__Answers rev 1 (
  topic: Triage__Topic
)

type Triage__Route rev 1 (
  route: str
)

judgment triage__questions(state: str) -> Triage__Answers rev 1
  service Jev__Service
  model "jev-1.13.0"
  timeout_ms 2000
  retries 0
  emits [jev.unavailable, jev.invalid_response]
  effects [Jev.request]
  questions
    topic: Triage__Topic choice
      instructions "Which team should handle this message?"
      criteria
        billing = "Charges, invoices, and payment problems"
        other = "A subject outside billing"
      confidence high d"0.8" medium d"0.5"

fn triage__route(message: str) -> Triage__Route rev 1
  emits [triage.unavailable]
  effects [Jev.request]
  tests
    billing(message = "charged twice") => Ok(route = "billing") pinned
    other(message = "charged twice") => Ok(route = "general") pinned
    high_boundary(message = "charged twice") => Ok(route = "billing") pinned
    medium_boundary(message = "charged twice") => Ok(route = "confirm") pinned
    low(message = "charged twice") => Ok(route = "review") pinned
    tie(message = "charged twice") => Ok(route = "review") pinned
    offline(message = "charged twice") => triage.unavailable(reason = "model_unavailable") pinned
    malformed(message = "charged twice") => triage.unavailable(reason = "invalid_model_response") pinned
  match judge triage__questions(state = message)
    given
      billing => [exchange args (state = "charged twice") outcome Ok(topic = Triage__Topic(choice = "billing", probabilities = Triage__Probabilities(billing = d"1", other = d"0"), confidence = d"1"))]
      other => [exchange args (state = "charged twice") outcome Ok(topic = Triage__Topic(choice = "other", probabilities = Triage__Probabilities(billing = d"0", other = d"1"), confidence = d"1"))]
      high_boundary => [exchange args (state = "charged twice") outcome Ok(topic = Triage__Topic(choice = "billing", probabilities = Triage__Probabilities(billing = d"0.9", other = d"0.1"), confidence = d"0.8"))]
      medium_boundary => [exchange args (state = "charged twice") outcome Ok(topic = Triage__Topic(choice = "billing", probabilities = Triage__Probabilities(billing = d"0.8", other = d"0.2"), confidence = d"0.5"))]
      low => [exchange args (state = "charged twice") outcome Ok(topic = Triage__Topic(choice = "billing", probabilities = Triage__Probabilities(billing = d"0.6", other = d"0.4"), confidence = d"0.49"))]
      tie => [exchange args (state = "charged twice") outcome Ok(topic = Triage__Topic(choice = "billing", probabilities = Triage__Probabilities(billing = d"0.5", other = d"0.5"), confidence = d"0"))]
      offline => [exchange args (state = "charged twice") outcome jev.unavailable(reason = "sdk_failure")]
      malformed => [exchange args (state = "charged twice") outcome jev.invalid_response(reason = "protocol")]
    on jev.unavailable _ => triage.unavailable(reason = "model_unavailable")
    on jev.invalid_response _ => triage.unavailable(reason = "invalid_model_response")
    on answers (topic = high "billing" _) => Ok(route = "billing")
    on answers (topic = high "other" _) => Ok(route = "general")
    on answers (topic = medium _) => Ok(route = "confirm")
    on answers (topic = low _) => Ok(route = "review")
```

**Grammar commitments:** after `on answers`, the parenthesized field patterns name every question exactly once, in declaration order. Choice accepts `high "label" binder`, `medium binder`, and `low binder`. Score accepts `high binder`, `medium binder`, and `low binder`. Noul accepts `yes binder`, `no binder`, and `uncertain binder`. A binder is a name or `_`. A whole field may be `_`, explicitly ignoring that answer. Error outcomes remain separate exhaustive `on error.kind` arms; no wildcard swallows service errors. The answer-product match is over **one completed request**, not an eager sequence of separate model calls.

Score declarations use an ordered `levels` list, a `score` question kind, the same explicit confidence partition, and an observation record with `score: dec`, `probabilities: Seq<dec>`, and `confidence: dec`. Sequence position is the numeric level. Noul declarations use a `noul` question kind, explicit `yes`/`no` criterion descriptions, `probability no L yes U`, and an observation record containing `noul: dec` only. These are primitive-specific forms, not alternative spellings for the same operation. The initial service profile admits 2–255 Choice options and 2–10 Score levels; timeout is an integer from 1 through 2147483647 milliseconds. The model token must match the versioned Jev form `jev-MAJOR.MINOR.PATCH`, with numeric components; floating aliases are rejected. Provider availability is checked only at runtime, never by querying a model list during compilation. Primitive structures: [TS-C], [TS-S], [TS-N].

A wildcard is explicit disregard, **not** a compiler endorsement that acting without regard to uncertainty is safe. Likewise a low arm can contain an application error or review path; the language cannot infer that an arbitrary business action is safe merely from the arm's name.

## 6. Expected TypeScript emit sketch

This sketch shows the actual SDK boundary and the emitted business control flow. `$canJevDecodeTriage` and `$canJevClassifyTopic` are **compiler-owned generated helpers**, not claimed SDK functions. Their bodies must implement §4.3–4.4 and be frozen by parity tests; their declarations below describe the sketch's boundary, not a completed implementation.

```ts
import { choice, TypeSafeClient } from "@typesafe-ai/sdk";

type TriageTopic = {
  choice: "billing" | "other";
  probabilities: { billing: string; other: string };
  confidence: string; // canonical exact can dec, not JS arithmetic
};

type JudgmentReply =
  | { kind: "answers"; topic: TriageTopic }
  | { kind: "jev.unavailable"; reason: string }
  | { kind: "jev.invalid_response"; reason: string };

type TopicDecision =
  | { kind: "high.billing" }
  | { kind: "high.other" }
  | { kind: "medium" }
  | { kind: "low" };

type RouteResult =
  | { kind: "Ok"; route: string }
  | { kind: "triage.unavailable"; reason: string };

declare function $canJevDecodeTriage(
  raw: unknown,
  expectedModel: string,
): JudgmentReply;

declare function $canJevClassifyTopic(
  observation: TriageTopic,
  high: string,
  medium: string,
): TopicDecision;

async function $triageRequest(state: string): Promise<JudgmentReply> {
  let raw: unknown;
  try {
    // Deliberately not at module initialization: construction can fail.
    const client = new TypeSafeClient({
      baseURL: "https://api.typesafe.ai",
      defaultModel: "jev-1.13.0",
      dangerouslyAllowBrowser: false,
      logLevel: "warn",
      retry: { maxRetries: 0 },
      timeout: 2000,
    }); // API key is obtained only at runtime, outside source fixtures.
    raw = await client.systemOne(
      {
        model: "jev-1.13.0",
        state,
        questions: {
          topic: choice("Which team should handle this message?", {
            billing: "Charges, invoices, and payment problems",
            other: "A subject outside billing",
          }),
        },
      },
      { timeout: 2000, retry: { maxRetries: 0 } },
    );
  } catch {
    // Catch only the foreign boundary, never arbitrary business logic.
    // Do not expose SDK messages, request bodies, or credentials.
    return { kind: "jev.unavailable", reason: "sdk_failure" };
  }
  return $canJevDecodeTriage(raw, "jev-1.13.0");
}

export async function triage__route(message: string): Promise<RouteResult> {
  const response = await $triageRequest(message);
  switch (response.kind) {
    case "jev.unavailable":
      return { kind: "triage.unavailable", reason: "model_unavailable" };
    case "jev.invalid_response":
      return { kind: "triage.unavailable", reason: "invalid_model_response" };
    case "answers": {
      const decision = $canJevClassifyTopic(response.topic, "0.8", "0.5");
      switch (decision.kind) {
        case "high.billing": return { kind: "Ok", route: "billing" };
        case "high.other": return { kind: "Ok", route: "general" };
        case "medium": return { kind: "Ok", route: "confirm" };
        case "low": return { kind: "Ok", route: "review" };
      }
    }
  }
}
```

No production path reads `given`, test names, or a cassette. The decoder checks the returned model ID, question membership and kind, all support entries, and numeric protocol. Invalid data becomes the declared protocol error, not a thrown assertion. Adapter protocol failures and SDK failures remain distinguishable. SDK constructor failures are inside the foreign boundary too. The sketch intentionally uses no `default`/`never` exhaustiveness scaffold. Async lowering is propagated through the explicitly declared request capability, not guessed from imported TypeScript. Source receipts: [R6/R11][R], [a14 §Reinstated][A14]; SDK signatures and options: [TS-J], [TS-JC], [TS-JF], [TS-JO], [TS-JY].

**Scope of the attempt guarantee:** one reached `judge` performs at most one SDK request attempt. A network failure does not establish whether the server performed an inference; no exactly-once remote-execution claim is made. Any subsequent attempt requires another explicit invocation governed by the existing termination and script-consumption rules.

### `errors.json` impact

Keep the current array schema: `kind`, `fields`, `raised_by`, `handled_by`, `hit_by_tests`. For judgment programs, include declared `jev.unavailable`, `jev.invalid_response`, and application errors, with the service manifest identified as producer and actual source error handlers attributed. The example adds `triage.unavailable`; it does not turn billing, high, low, or uncertainty into error kinds.

**Do not reinterpret `hit_by_tests`.** The current catalog explicitly collects script/expectation references, not observed executions or escaping errors. Extend its walkers to the new sites while preserving that meaning and `TestCatalogNoPhantomHits`. Measured arm execution belongs to coverage reporting, not a new meaning for that field. Programs without judgments keep byte-identical catalogs. Full evidence-bearing catalogs remain build reports outside the production package, including source maps or archives that would otherwise leak stripped tables. Receipt: [compiler/catalog.go, catalogEntry / buildCatalog][CATALOG].

## 7. Alternatives: strongest case and concrete counterexample

The fragments below use §5's declarations unless explicitly marked as alternative syntax. They demonstrate policy mistakes, not claims that the current compiler already parses the proposed judgment grammar.

### A — distributions in `given`

**Strongest case:** this is already the right hermetic abstraction: explicit external data at the exact request site. Accept it as the foundation.

**What is rejected:** an incomplete variant that scripts label/probabilities but guesses confidence, or treats the selected probability as confidence:

```can
// Invalid successful fixture: confidence has been omitted.
billing => [exchange args (state = "charged twice") outcome Ok(topic = Triage__Topic(choice = "billing", probabilities = Triage__Probabilities(billing = d"0.63", other = d"0.37")))]
```

Inferring confidence `0.63` can select high at a `0.6` threshold where a separately reported confidence `0.53` would not. This distinction is real: the live Choice example reports a selected probability of 0.63 with confidence 0.53. It does not license using either value for the other. Our example's two-option fixture is illustrative, not a reproduction of that complete provider response. [TS-C]

**Disposition:** full-observation A survives and is adopted. It needs explicit validation, dispatch, and proof-separation rules; “probabilities are just data” alone does not specify them.

### B — record/replay cassettes

**Strongest case:** a checked-in immutable cassette with complete request identity, a versioned decoder, and no live fallback **is hermetic**. Calling all replay “weaker hermeticity” would be incorrect. It can preserve useful regression examples and provenance.

**Failure of weak replay identity:** consider this rejected alternative spelling:

```can
match judge triage__questions(state = "a different customer's message")
  given
    billing => replay "charged-twice-run.json" at 0
```

A site/ordinal replay that supplies the recorded billing answer without checking the new state recreates a12's wrong-request bug. A complete-request cassette fixes that bug; acknowledge the fix rather than claiming replay is inherently unsound. [A12]

**Failure of passive recordings as the complete test strategy:**

```can
// The sole recorded response was high billing.
on answers (topic = high "billing" _) => Ok(route = "billing")
on answers (topic = high "other" _) => Ok(route = "general")
on answers (topic = medium _) => Ok(route = "confirm")
on answers (topic = low _) => Ok(route = "review")
```

The last three arms still lack witnesses and incur CAN4107. Running the model repeatedly until all appear does not give a deterministic, bounded way to construct the test suite. Synthetic counterfactual fixtures remain necessary. [A03]

**Disposition:** recordings are allowed as external evidence sources, not as a second canonical compiler input. A strong cassette importer that preserves identity and emits ordinary full exchanges conforms to this decision. A fully checked external fixture format is not logically unsound; it loses here on duplicated artifact/authority paths and the Goal's colocation requirement, not on a fabricated impossibility proof. [R]

### C — live calls during builds

**Strongest case:** live results can expose integration failure and operational drift that fixtures cannot. That is valuable operational evidence.

**Killing counterexample:**

```can
// Alternative live-build policy: no fixture for this reaching test.
fn triage__live_example(message: str) -> Triage__Route rev 1
  emits [triage.unavailable]
  effects [Jev.request]
  tests
    stable(message = "charged twice") => Ok(route = "billing") pinned
  match judge triage__questions(state = message)
    // Under C, the compiler asks the service here.
```

Suppose identical source receives billing with confidence 0.81 on one build and 0.79 on another. The declared 0.8 threshold produces different routes. Even a temporarily stable provider leaves the gate dependent on credentials, availability, quota, and network behavior. Missing arms in this abbreviated fragment are beside the point: completing them cannot repair the external build dependency. Live-per-keystroke behavior would also make ordinary editor validation spend money and disclose inputs. This directly contradicts R7, so reject it rather than adding a “best effort” fallback. [R]

### D — bands, guessed observations, or hybrids that blur authority

**Strongest case for bands:** authors often care about routing categories, not arbitrary decimal details. A symbolic band could quantify over many distributions.

**Killing counterexample for preclassified scripts:**

```can
// Rejected alternative: scripts bypass the policy being tested.
billing => [exchange args (state = "charged twice") outcome High(choice = "billing")]
```

Changing the policy from high 0.8 to high 0.99 leaves this preselected branch unchanged. The table no longer tests the threshold elaboration. A sound symbolic alternative would have to define a set of observations and verify that **every** member yields the expected behavior, including ties and distribution-shape predicates. That is a different quantified-testing feature, not a harmless shorter spelling.

**Killing counterexample for live-checking pinned rows as labels:** §5 pins both `billing` and `other` for the same message under different fixtures. A live answer cannot simultaneously satisfy both as label-truth expectations. Treating either as ground truth silently changes what the row means. Keep live model-quality acceptance separate from deterministic routing acceptance. [A87]

**Disposition:** adopt only the hybrid with full source observations plus optional, independently authorized recording and live-evaluation reports.

### What to take from Probably—and what not to take

Probably documents highest-probability matching, a winning-probability gate, first-option tie behavior, sampling, and replay of responses and random draws. Its “otherwise maybe” path can be omitted, leaving neither ordinary branch to run. Its replay promises “without making model calls.” Those are its documented semantics, not TypeSafe's separate confidence statistic. [PROB]

Take explicit uncertainty handling, provenance-rich recordings, and zero-model-call replay. Reject optional fallthrough with no result under can-lang's exhaustive expression rule; reject first-option tie selection as our safety policy; reject chaos in the judgment core; and reject importing a winning-probability gate under the name of TypeSafe confidence. Replay random draws would be necessary for a future explicit sampling effect, but replay alone never proves statistical coverage. [R6/R7][R], [A03]

## 8. Edge-case rulings

| Case | Ruling |
|---|---|
| **Choice tie** | Any shared maximum routes low, before confidence gating. The reported label must still be a maximizer. Neither declaration order nor the SDK's tie preference chooses an actionable high branch. A tie between two options above other options is also a tie. |
| **Noul near 0.5** | It is probability uncertainty, not missing confidence. With no 0.2 / yes 0.8, 0.49, 0.5, and 0.51 are all uncertain. Exact bounds 0.2 and 0.8 choose no and yes respectively. |
| **Score ties / multimodality** | Score is the returned weighted mean. A distribution split between levels 0 and 2 can have the same mean as mass at level 1. Preserve shape and reported confidence; do not invent a winning-level branch or equate equal means with equal evidence. |
| **An arm reachable only below a floor** | Test it through the low path. A high-path nested condition requiring `confidence < M`, where high requires `confidence >= H > M`, cannot acquire a witness by assigning its label positive probability. CAN4107 remains; remove or relocate the impossible arm. No new “model says unreachable” certificate. |
| **Missing versus leftover scripts** | Missing and unvisited is legal. Missing and reached fails. A scripted nonempty queue with zero visits is leftover, as is a partially consumed queue. Applicable queues must be registered before execution; see §4.2. |
| **Probabilities, confidence, or both?** | Both for Choice/Score, plus their reported result. Only P(yes) for Noul. Confidence bands are computed, never supplied. No inferred complement for a missing Choice option, even in binary Choice. |
| **Distribution-shape assertions** | Explicit deterministic logic may inspect the full observation and require, for example, that a competitor's mass stay below a declared bound. Such logic needs its own rows/arms. The compiler supplies no extra formula or calibration inference. |
| **Chaos sampling** | Not admitted in JEV-1.0 at build or runtime. A later randomness effect would need its own request/draw scripts, seeded replay semantics, and ordinary arm witnesses. This decision does not reserve sampling as an implicit mode of `judge`. |
| **Production stripping** | Strip source tables and evidence from TS, bundles, source maps, embedded manifests, and shipped reports. Keep runtime instructions, criteria, schemas, thresholds, and error handling. A full evidence catalog is not a production artifact. |
| **Prompt reworded without bump** | Descriptor instruction/criteria text changes its interface fingerprint: CAN6013. Literal decoding is canonical; comments are ignored. Model ID and policy changes are also descriptor changes. User input changes do not require revisions. State-preparation code is not magically reclassified as prompt metadata; §4.5 states the precise boundary. |
| **Provider drift under the same version ID** | Offline checks cannot detect changed model behavior. Preserve model/SDK/adapter/request identity in authorized operational reports and re-evaluate against independently labeled cases. A same-ID behavior change cannot be silently “repaired” by refreshing fixtures. |
| **Invalid live numbers or missing answers** | Reject the entire batch as `jev.invalid_response`. No NaN, infinity, clamping, fabricated zeros, absent-field defaults, or partial-success recovery. Invalid successful source fixtures are diagnostics instead. |
| **Unavailable credentials / timeout / SDK failure** | Production receives `jev.unavailable`; tests script it. Builds never read credentials. Catch SDK construction and request failures only at the generated boundary. |
| **Repeated identical requests** | Separate occurrences consume separate exchanges; observations may differ. Argument equality is not a promise of a deterministic provider, and no result cache is inserted. |
| **A probability of 1** | It is still a reported estimate. It does not mint a brand, prove a predicate, authorize a payment, or eliminate a required error arm. |

### 8.1 Fan-out, concretely

**Admit explicit batching, forbid automatic batching.** A descriptor may declare several questions, all evaluated against its one supplied state. Add each question to the answer-record schema. One `match judge` consumes one exchange whose `Ok` record contains every answer. This is a **joint scripted observation**; there is no multiplication of marginal probabilities or assumption that errors are independent.

For example, after adding a fully declared `urgency` question, the shape is:

```can
// Fragment: TopicObservation and UrgencyObservation stand for complete
// typed constructors, not permitted abbreviated fixture syntax.
row => [exchange args (state = "charged twice") outcome Ok(topic = TopicObservation, urgency = UrgencyObservation)]
```

Real source must replace those illustrative names with complete constructors like §5. An answer-product arm can then deliberately ignore an already-received result:

```can
on answers (topic = high "billing" topic, urgency = _) => Ok(route = "billing")
```

Consuming the response consumes the exchange **once**. Ignoring `urgency` afterward is not a leftover exchange. Omitting `urgency` from the successful fixture, by contrast, is a schema error. If it is present but malformed, the whole live response is rejected even though that branch would ignore it. That fail-closed rule avoids branch-dependent protocol validation.

All required arguments and questions must exist before this one request. A question whose instructions, criteria, or state depend on an earlier answer requires a later, explicitly reached `judge` site with a separate exchange queue. Rewriting a conditional second request into unconditional fan-out changes disclosure, cost, and invocation behavior; the compiler must not perform that transformation. Operational parallelism does not weaken can-lang's eager, explicit effect rules. Source basis: [R6/R8][R]; TypeSafe batching basis: [TS-P], [TS-B], [TS-H].

### 8.2 The impossible-arm example

This complete arm body fragment uses the declared partition from §5:

```can
on answers (topic = high "billing" observed) =>
  match observed.confidence < d"0.5"
    true => Ok(route = "review")
    false => Ok(route = "billing")
```

No admissible observation reaches the inner true arm: the enclosing high state requires confidence at least 0.8, and the tie case has already gone low. A fixture with label billing and confidence 0.49 takes the outer low arm instead; a `High(...)` stub cannot override the policy. The inner true arm is uncovered, CAN4107. This is a logical inconsistency in the implemented routing, not a statistical rarity to waive. [A03]

## 9. Verification and delivery plan

The following are **required acceptance checks**, not reports of green runs. Keep each independently shippable compiler surface isolated; no slice may advertise first-class judgments as usable until its full proof/evaluation/emit path is gated.

### 9.1 Existing suites and artifacts

| Existing surface | Required treatment |
|---|---|
| `compiler/emit_golden_test.go` | Preserve `TestGoldenAuthLogin`, `TestGoldenRetryLoop`, `TestGoldenCounter`, and other non-judgment outputs as byte-identical controls. Add judgment goldens rather than silently updating those programs. Their established compile-plus-emit pattern is the precedent. [GOLDEN] |
| `TestGoldenJSONDiags`, diagnosis/LSP tests | Add one-fault judgment fixtures for each new rule and freeze code, severity, Expected/Found/Hint, and exact token/UTF-16 spans. Assert editor and CLI agreement. Update the existing CAN3104 dead-script severity golden expressly; preserve unrelated diagnosis output. [CODES], [R10][R] |
| `TestCodesUnique`, `TestAllCodesSequenced`, `compiler/explain.go` coverage | Register the new block in order; never reuse retired numbers; provide minimal counterexample and repair for each code. [CODES] |
| `compiler/contract_verify_test.go` | Add judgment and scripted-evidence negative controls: a label/confidence fixture must not prove an input predicate; unsupported effects cannot disappear from a proof walk; unverified dependencies must remain unverified. Retain the existing proved-summary/scripted-outage separation controls. [A81] |
| Revision identity / accepted-baseline suites | Same-revision prompt, criterion, model, gate, timeout, schema, and numeric-ABI changes reject. Comments and equivalent decoded literals remain inert. Fixtures change evidence, not API identity. Explicitly test incompatible baseline versions and refusal to use candidate-generated baselines as authority. [A77] |
| `compiler/acceptance_test.go` | CAN6017 behavior remains unchanged. Add independent CAN3511 cases: mutate a fixture, input, request descriptor, or reach-controlling exchange while retaining the expected output; remove the last judgment site; proposed rows remain untracked; generated baselines do not accept themselves. [A87] |
| Catalog / `TestCatalogNoPhantomHits` | Add intrinsic producer attribution and source-handler traversal for judgments. Preserve reference-based `hit_by_tests`; do not count probability support as hits, certificates as executions, or labels as errors. [CATALOG] |
| `.github/workflows/tsc.yml`, `tscheck/` | Type-check the real generated helpers against a lockfile-pinned TypeSafe SDK with the repository's pinned strict checker. Do not substitute hand-written fake SDK signatures. Preserve the separate canlc behavioral gate. [A14] |

### 9.2 New offline suites

**Protocol and numeric parity.** Compare the Go evaluator and emitted JS decoder/selector on shared committed vectors: all scalar boundaries; canonical number conversion; rejected excess precision; mass exactly at and beyond ε; no normalization; selected-label mismatch; ties including ties at high reported confidence; Score weighted means and equal-mean/different-shape observations; all Noul endpoints; absent or extra fields; wrong question kinds; wrong response model; nonfinite wire values. Independently compute expected outcomes for boundary vectors so a shared bug is not mistaken for parity.

**Dispatch and coverage.** Exercise each source high label, medium, low, no/yes/uncertain, service error, and distribution-shape branch. Include shadowed arms, the impossible-arm fragment above, and a non-winning positive-probability label. Assert exact source-arm hits, not merely successful outputs. Test batch product matches and explicit ignored answers. Exhaustiveness analysis must fail closed on resource exhaustion; it must never sample product cells or silently treat an incomplete analysis as total.

**Exchange lifecycle.** Test missing table, omitted/unvisited, omitted/reached, empty-on-hit, wrong request, repeated sites, identical requests with differing answers, partly consumed lists, zero-consumed lists, shared helpers, unrelated root-test queues, unknown keys, retired `-`, and non-reached second questions. The same fixture corpus must run through CLI and editor semantics.

**SDK boundary without the service.** Use the SDK's documented injectable transport in compiler-maintainer tests. Capture the exact serialized request; test constructor failure, request failure, timeout behavior, one-attempt behavior, reported-model mismatch, malformed response, and successful fan-out. Assert that no real transport is reachable. A wire-validation failure test is not implemented by merely stubbing the final `.can` protocol error; test both layers separately. The SDK supports custom fetch specifically for transport/tests. [TS-JF]

**Proof and authority adversaries.** Mutate an agent-authored model answer from wrong to right while preserving a pinned expected route. CAN3511 must expose the premise change; no solver fact may change because of it. Test a forged “confidence equals truth” claim and a mocked outcome impossible under an unrelated verified provider summary. Evidence/proof separation must remain explicit. [A81], [A87]

**Production and hermeticity.** Build with network denied and credentials absent; vary locale, time, environment model aliases, and random seed and compare deterministic artifacts. Put unique sentinels in test names, fixture inputs, and synthetic probabilities, then verify their absence from production TS, bundles, maps, and packaged sidecars. Ensure production imports no fixture loader. Non-judgment code must not acquire async wrappers or SDK imports.

### 9.3 modcheck, gramcheck, and the one gate

`modcheck` must recognize judgment declarations in `provides`/`uses`, the intrinsic service manifest and its revision, declared answer types, and explicit effect closure. It must distinguish same-file judgment sites from same-file ordinary helpers. Check site/test identity without reintroducing total tables or `-`; CAN3104 is an error in the shared compiler path under the explicit severity amendment. Never contact a provider or model catalog. Existing module resolution rules, including provider ownership, remain authoritative. [R2/R8][R], [A91]

`gramcheck` and the editor grammar must add scopes/samples for `judgment`, `judge`, descriptor fields, primitive kinds, confidence/probability policy, `on answers`, and typed answer-product patterns. Preserve literal handling, decimal syntax, indentation, and the braces ban. The retired `-` token can retain its historical coloring without becoming legal syntax again. [R1/R10][R], [A91]

Required ordinary commands remain `go test -count=1 ./...`, `go run ./tools/modcheck`, `go run ./tools/gramcheck`, and the repository's pinned strict TypeScript check. Wire the new diagnostics, proof refusals, coverage, and table runs into the existing common canlc checking path—not a special “Jev test” command that ordinary builds omit. Live evaluation is expressly outside these commands. [R10][R], [A14]

### 9.4 Phased order

1. **Ratify text and identities:** register JEV-1.0, reserve codes, define intrinsic manifest and baseline-format additions, commit grammar and diagnosis specifications. No runtime capability is advertised yet.
2. **One-question Choice, end to end:** descriptor checking, observations, exact request exchanges, compiler dispatch, source coverage, error catalog, SDK adapter, production stripping, strict emit and offline parity. Land it only with negative controls.
3. **Revision and acceptance enforcement:** integrate descriptor fingerprints and CAN3511 into the accepted-baseline path before applications are allowed to treat migrated judgment contracts as reviewed. Generation stays distinct from acceptance.
4. **Noul and Score separately:** each primitive gets its own shape, policy, boundaries, and parity fixtures; neither is coerced into Choice.
5. **Explicit multi-question batching:** extend the same response/queue semantics to product patterns, independent of external recording tooling. No optimizer-driven batching.
6. **Migration and optional operational tooling:** convert one real extern-stubbed consumer, remove its superseded imports/manifests, document changed authority, then add recording-to-source import and independent live-evaluation reports only as separately authorized tools.

During staged implementation, unsupported surfaces fail explicitly; there is no label-only interim mode or hidden live fallback. Schema/refinement typing, richer state formats, and sampling must not hitchhike on the first Choice slice.

## 10. Remaining questions, ranked by blast radius

These do not leave the central decision undecided. They identify evidence or mechanisms required before stronger claims or broader deployments.

| Rank | Question / exact missing evidence | Binding disposition meanwhile |
|---|---|---|
| **1 — epistemic authority** | Who owns independently labeled live acceptance data, and what deployment consequences follow from degraded quality? Missing: the application's approved evaluation corpus, acceptance owner, and release policy. | No live result is a compiler proof or a CAN6017 event. No unattended fixture/threshold/model refresh. The language retains A-light, not invented A-strong authorization. |
| **2 — service identity and confidence** | What numerical confidence algorithm and immutability guarantees attach to a versioned service ID? Missing: a provider specification of the algorithm, versioning guarantees, and conformance vectors. | Script reported confidence, preserve probabilities, check returned model identity, and monitor behavior separately. Never claim that canlc verifies confidence/distribution consistency beyond the specified checks. |
| **3 — adapter release conformance** | Does the selected SDK release preserve all requested wire shapes and the numeric ABI over actual responses? Missing: an exact SDK lockfile plus committed wire, error, decimal-conversion, and emitted-runtime conformance fixtures. | This is a **shipping gate** for the adapter. Revise a failing ABI assumption by a tagged amendment; do not silently relax checks or label an untested sketch working code. |
| **4 — batching scale** | What deterministic exhaustiveness implementation handles the intended number and size of question products? Missing: a representative application's largest descriptors and checker performance tests. | Exact analysis or explicit failure; never approximate coverage, statistical independence, or proof. Start with bounded tested implementation surfaces, not hidden acceptance. |
| **5 — richer input and disclosure** | Which record/sequence/brand values can become provider state without violating representation or disclosure boundaries? Missing: a separately ratified state-serialization and brand-authority design. | JEV-1.0 uses explicit text state. No generic stringification, automatic unsealing, or guessed context attachment. |
| **6 — recording lifecycle** | What provenance, redaction, retention, and import format should recordings use? Missing: a reviewed operational format and data-handling policy. | Optional tooling only. Compiler input remains ordinary source exchanges; recording provenance never auto-pins acceptance. |

**Ratification action:** accept JEV-1.0 as the conditional-observation semantics and its associated amendments. Do not interpret ratification as green compiler gates, completed SDK integration, approved production thresholds, or a calibration certificate.

---

## Ten-line summary

1. Ratify full-observation `given` scripts plus deterministic compiler-owned routing; no model calls in builds.
2. Choice and Score carry complete probabilities and reported confidence; Noul carries only P(yes).
3. Passing rows certify conditional program behavior, never model correctness or threshold calibration.
4. Coverage requires a selected source arm or the existing authorized certificate; probability mass is not a hit.
5. Choice ties route low; Noul has an explicit uncertain interval; Score remains a weighted mean.
6. One explicit batch consumes one complete exchange; ignored answers are not leftover requests.
7. Strong replay can be hermetic, but recordings enter builds only after conversion to ordinary source evidence.
8. Prompts and policies are revision-bearing; pinned fixture-premise changes get a separate non-blocking warning.
9. Emit the TypeSafe JavaScript SDK with explicit model, timeout, zero retries, typed errors, and no fixture bytes.
10. Keep live evaluation independently authorized and non-blocking; no compiler changes or gates were run here.

## Source references

Repository links below are pinned to the reviewed commit. Section names in the text identify the receipt within each file.

[R]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/REQUIREMENTS.md "REQUIREMENTS.md — Goal; R2; R4–R8; R10–R11"
[A03]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a03-branch-coverage.md "a03 — Rule; Hit semantics; retrospective amendment"
[A12]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a12-contracts.md "a12 — Rule; Theorem"
[A81]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a81-proof-obligations.md "a81 — Obligations; Findings; Deliberate decisions"
[A87]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a87-acceptance-authority.md "a87 — Semantics; Interplay"
[A14]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a14-tsc.md "a14 — Reinstated (a70); checker-upgrade note"
[A77]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a77-revision-identity.md "a77 — Fingerprint model; Enforcement; Baseline workflow"
[A91]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a91-given-omission-amendment.md "a91 — S0 rulings; stale-script lint"
[AUTH]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/sketches/auth-login/auth.can "auth.can — auth__login; auth__verify"
[CODES]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/compiler/code.go "compiler/code.go — diagnostic registry"
[CATALOG]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/compiler/catalog.go "compiler/catalog.go — catalogEntry; buildCatalog"
[EVAL]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/compiler/eval.go "compiler/eval.go — canonical decimals and decCmp"
[GOLDEN]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/compiler/emit_golden_test.go "compiler/emit_golden_test.go — existing golden suites"
[TS-I]: https://docs.typesafe.ai/llms.txt "TypeSafe documentation index"
[TS-P]: https://docs.typesafe.ai/primitives "TypeSafe Primitives"
[TS-C]: https://docs.typesafe.ai/primitives/choice "TypeSafe Choice — response structure and example"
[TS-N]: https://docs.typesafe.ai/primitives/noul "TypeSafe Noul"
[TS-S]: https://docs.typesafe.ai/primitives/score "TypeSafe Score"
[TS-F]: https://docs.typesafe.ai/confidence "TypeSafe Confidence"
[TS-H]: https://docs.typesafe.ai/concepts/how-to-build-with-system-one "TypeSafe How to build with System One"
[TS-B]: https://docs.typesafe.ai/patterns/fan-out "TypeSafe Fan-out"
[TS-R]: https://docs.typesafe.ai/patterns/confidence-routing "TypeSafe Confidence routing"
[TS-J]: https://docs.typesafe.ai/sdk/javascript "TypeSafe JavaScript SDK"
[TS-JC]: https://docs.typesafe.ai/sdk/javascript/api/classes/TypeSafeClient "TypeSafeClient API"
[TS-JF]: https://docs.typesafe.ai/sdk/javascript/api/interfaces/TypeSafeClientConfig "TypeSafeClientConfig API"
[TS-JO]: https://docs.typesafe.ai/sdk/javascript/api/interfaces/RequestOptions "RequestOptions API"
[TS-JY]: https://docs.typesafe.ai/sdk/javascript/api/interfaces/RetryPolicy "RetryPolicy API"
[TS-M]: https://docs.typesafe.ai/models "TypeSafe model IDs and aliases"
[PROB]: https://probably-lang.southpolesteve.workers.dev/ "Probably language reference"
