# JEV-SURFACE-1 — Three native syntax families for can-lang

**Status:** design alternatives, not a ratified syntax and not implementation authorization.  
**Recommended direction:** **Option B, Question signatures**, with Option C's success/failure layout for multi-question requests.  
**Repository snapshot:** `veighnsche/can-lang`, `main` resolved to `4d447bac858c50d998804aa7fa96f2b66b3b557b`.  
**Review date:** 18 September 2026. No compiler changes, provider calls, or compiler gates were run.

## 0. Reading basis and the source discrepancy

The GitHub read of `docs/a94-jev-judgment-evidence-decision.md` returned **404 at the retrieved main commit**. I therefore read the complete **647-line JEV-1.0 document attached in this conversation**, rather than a repository copy. References to JEV-1.0 below identify sections in that document. The current instruction supplies the overriding transport settlement: **plain platform fetch, no SDK, runtime `TYPESAFE_API_KEY`**. The attached document's SDK requirement is superseded, not silently retained. Differences in an unavailable, subsequently edited a94 copy cannot be verified here. ([JEV1] §§4–6)

I also read the pinned `REQUIREMENTS.md`, `docs/can-idioms.md`, `sketches/auth-login/auth.can`, and `compiler/code.go`. The existing receipts for a03, a77, a81, a87, and a91 establish the amendment interactions. The requested TypeSafe `.md` page fetches failed through the retrieval interface; the corresponding **live HTML pages linked by `llms.txt`** were readable and were used instead. ([R]; [IDIOMS]; [AUTH]; [CODES]; [TS-INDEX])

### The three alternatives at a glance

| Option | What the author writes | Where the question lives | What genuinely differs |
|---|---|---|---|
| **A — Primitive-led expressions** | `noul`, `choice`, or `score` starts the expression; `together` explicitly batches. | Beside the decision arms. | One local unit contains the question, policy, evidence, and continuation. No separate question declaration. |
| **B — Question signatures** | `match noul`, `match choice`, or `match score` references a versioned primitive declaration. | In the declaration section, without transport fields. | Reusable questions have independently named contracts; use sites are short native matches. |
| **C — Judgment scopes** | `judge` obtains a named snapshot; `decide` reads already-classified answers. | Inside one explicit effect scope. | Request acquisition and subsequent decision expressions are distinct. Several decisions can reuse one received snapshot without another request. |

These are **mutually exclusive surface proposals**, not three spellings to ship together. Shared rules below isolate the decisions that should not vary merely to make the alternatives look different.

## 1. Capability inventory: every requested TypeSafe surface has a home

Short quotations anchor the external facts. The syntax and policies in the last column are proposals, not claims about TypeSafe syntax.

| Capability and exact live section | Source anchor | Syntactic home in all three options |
|---|---|---|
| Three typed primitives; named responses; mixed requests — `primitives`, **Define a question**, **What comes back**, **Ask multiple questions together** | “You can mix question types freely.” | Separate `noul`, `choice`, `score` forms; explicit batch question IDs. ([TS-P]) |
| Choice support and selected result — `primitives/choice`, **Request structure**, **Response structure** | “The option with the highest probability.” | Closed `options` maps; selected label, full probabilities, reported confidence remain available. ([TS-C]) |
| Noul question or statement, optional true/false criteria — `primitives/noul`, **Writing a Noul question**, **Noul does not return a separate confidence value** | “Noul does not return a separate confidence value” | `criteria true` / `false`; probability bounds and `yes`, `no`, `uncertain`, never a fabricated confidence field. ([TS-N]) |
| Score ordered levels, fractional result, shape — `primitives/score`, **Levels**, **Response structure**, **Reading a Score** | “Different distributions can produce the same score.” | Ordered `levels`; `.score`, `.probabilities[index]`, `.confidence`; level routing is explicit arithmetic on the mean, not argmax. ([TS-S]) |
| Structured instructions and descriptions — `primitives/advanced`, **Where structure is allowed**, **Structured instructions**, **Structured Choice options**, **Structured Score levels**, **Structured Noul criteria** | “Instructions, Choice options, Score levels, and Noul criteria all accept JSON structure.” | Explicit object/array description literals and metadata-only `unspecified`; examples below use each applicable location. ([TS-A]) |
| Structured Score legends — `primitives/score`, **Structured level descriptions** | The response example has object-valued legend entries. | Decoder compares the legend's complete description trees against the declared ordered levels; not a string-only legend parser. ([TS-S]) |
| Confidence supplied separately from probabilities — `confidence`, **Confidence is derived from the probabilities** | “a statistic computed from the probability distribution” | Preserve reported confidence; derive only the application's routing band. No locally invented confidence formula. ([TS-F]) |
| High/medium/low and action-specific risk policy — `confidence`, **Three paths for using confidence in your code**, **Thresholds scale with risk**; `patterns/confidence-routing`, **Step 2** | “Where you draw those boundaries depends on the stakes.” | Explicit threshold metadata; a branch can additionally inspect confidence for an action-specific bound. All shown thresholds are illustrative policy, not calibrated recommendations. ([TS-F]; [TS-R]) |
| State as text, object, or array — `concepts/state`, **State can be as simple as a string**, **Separate content from questions** | “State can also be a JSON object or array” | Plain `str`, or explicit object/array state construction. No capture of local variables, records, brands, or secrets by implication. ([TS-STATE]) |
| Paths inside questions — `primitives`, **Reference specific fields** | Paths are written inside instructions. | Ordinary literal text such as `` `ticket.message` ``; it is prompt text, not executable field access. ([TS-P]) |
| Independent/speculative questions and ignored results — `primitives`, **Ask speculative questions**, **When one question depends on another**; `patterns/fan-out`, main pattern | “One question’s answer is not hidden context for another.” | One batch, one complete observation, explicit ignored fields. A genuine answer-dependent question uses a later request. No automatic batching. ([TS-P]; [TS-FAN]) |
| Code owns workflow and composition — `concepts/how-to-build-with-system-one`, **Design a System One workflow**, steps 1–8 | “Keep control flow, deterministic rules, and side effects in code.” | Ordinary exhaustive expressions compose observations, compare thresholds, and invoke explicitly handled application calls. ([TS-H]) |
| HTTP transport and response envelopes — `api`, **Evaluation endpoint**, **Request body**, **Response body**, **Errors** | `POST /v1/systemone`; Bearer authentication. | Generated fetch boundary and centrally pinned profile; no request boilerplate in a question or use site. ([TS-API]) |

### Explicit exclusions, not missing designs

**Null instructions are rejected**, although Advanced permits them. A question must contain an explicit, non-null instruction; inferring a question from its identifier would undermine the Goal and TypeSafe's separation of IDs from prompt content. `unspecified` is available for optional criterion/level descriptions, including label-only Choice options. It denotes a foreign metadata entry, not a general can-lang null value. ([R] Goal / Non-goals; [TS-A] §Where structure is allowed; [TS-P] §Define a question)

**Dynamic question counts, runtime-selected option sets, and runtime-generated instruction text are not introduced here.** The requested pages describe taxonomy walks and programmatically generated question sets. A finite tree can be expressed as separately declared questions reached by explicit code, and bounded predeclared questions can be batched. Arbitrary dynamic support would need a separate exhaustive-domain and revision-identity design; it is not smuggled in as a string interpolation feature. ([JEV1] §§4.5, 8.1; [TS-A] §Walking a taxonomy; [TS-P] §When one question depends on another)

**Structured runtime state initially admits unbranded text leaves and nested object/array structure.** Numeric, Boolean, binary, branded, and arbitrary record leaves require a separately specified conversion before entry; this proposal does not silently stringify them. This deliberately does not cover every possible JSON leaf in the documentation's examples. It preserves exact can-lang numerics and disclosure authority rather than inventing an implicit JSON conversion. Metadata trees can contain the closed JSON-compatible scalar literals precisely admitted in §7 because their encoding is fixed at compile time; no dynamic record serialization follows from that permission. ([R] Non-goals / R10–R11; [JEV1] §4.5; [TS-STATE])

**Multimedia state, generation, provider model discovery during builds, floating model aliases, and automatic retries have no syntax here.** Multimedia is not the documented Jev input; generation is outside these three primitives. The other exclusions preserve the settlement. Multiple independently obtained Scores can be combined using existing exact arithmetic and explicit fallible conversion operations; this proposal adds neither floating arithmetic nor implicit normalization. ([JEV1] §§4.3–4.5; [TS-M]; [TS-H] §Combine question outputs in code)

The API reference's narrower string-oriented description types differ from the explicit structured examples and the Advanced page. The designs follow the **more specific structured documentation**, including structured legends. Actual wire conformance remains an offline adapter-fixture shipping requirement, not an assumption that every documented shape has been live-tested here. ([TS-API] §Question types / Score answer; [TS-A]; [TS-S] §Structured level descriptions)


## 2. Shared foundation: what becomes inherent

### 2.1 One transport file, no per-question service plumbing

Every option uses one project-root `jev.config.json` with this complete content:
```json
{
  "format": "can-jev-transport/1",
  "name": "Jev__Runtime",
  "rev": 1,
  "adapter": "typesafe-systemone-fetch/1",
  "model": "jev-1.13.0",
  "timeout_ms": 2000,
  "retries": 0
}
```

The adapter revision fixes the endpoint `https://api.typesafe.ai/v1/systemone`, HTTP method, authentication source, numeric ABI, decoder contract, failure union, and no-redirect policy. There is no configurable URL that an ordinary metadata edit could redirect a Bearer credential toward. Endpoint changes require an adapter-contract amendment. The selected full model ID is listed in the live Models page; moving aliases are refused. ([TS-API] §Evaluation endpoint; [TS-M] §Current models / Aliases)

A module owning judgments pins the config provider once:

```can
uses [Jev__Runtime@1]
```

That is an **ordinary visible dependency pin resolved against a new config-provider kind**, not a default guessed from the working directory. Resolution uses the existing project root; it never searches parent directories for alternate profiles or reads an environment-selected config. `modcheck`, LSP, and the compiler must resolve the same provider. The profile is mandatory only for programs owning judgments. Imported question declarations carry their owning profile identity; a project cannot silently combine incompatible profiles.

A site inherits the model, timeout, retry policy, adapter, endpoint, credential source, failure schema, and generated boundary. It does **not** inherit confidence thresholds: those remain explicit application policy in source. `retries: 0` stays mandatory under JEV-1.0. Moving the field to config does not authorize nonzero retries. A later retry requires another explicitly reached judgment and another exchange. ([JEV1] §§4.4–4.5, 6)

**What is hashed.** Hash the validated config object excluding `rev`, using compact UTF-8 JSON with keys sorted by Unicode scalar order, decoded strings preserved without Unicode normalization, and exact canonical integer literals. Duplicate or unknown keys reject before hashing. The concrete config above has SHA-256:

```text
d642e6c5d8d506cca3ee29ff82d56b4daabcf8f257cb8bd88034c10def2c904a
```

This hash was computed for the document, not read from an existing accepted baseline. The proposed versioned `ai-lock.json` extension records provider name, revision, hash, and format alongside the existing authority mechanism. An actual accepted baseline remains owner-approved history; generating this record never accepts it. ([a77] §Canonicalization / Baseline workflow)

**Runtime does not reread the config file.** The compiler embeds the validated profile and its identity in emitted runtime constants. Only `TYPESAFE_API_KEY` is read from the runtime environment, at the request boundary. Rotating that credential does not change source identity. Model and timeout environment overrides do not exist.

### 2.2 Exact revision consequences

| Change | Required behavior |
|---|---|
| Change model, timeout, retry value, or adapter without changing profile revision | `CAN6013` at the changed profile field; no emit. Existing validation still rejects aliases, nonzero retries, and invalid timeouts. |
| Bump profile revision while a module still pins `Jev__Runtime@1` | `CAN2103` at the stale pin when revision 1 is not supplied. Never silently select revision 2. This is not a promise of unimplemented multi-version loading. |
| Repin profile, but retain the affected judgment contract revision | `CAN6013`: its effective request identity changed. |
| Change instruction words, criteria, level order, policy bounds, state schema, or question membership at the same owning revision | `CAN6013`. Equivalent decoded literals and comments remain inert. |
| Edit runtime input text or a given observation | Not a request-interface revision change. Rerun evidence; apply JEV-1.0's `CAN3511` pinned-premise warning where applicable. |
| Change only JSON formatting or key order | Same canonical profile identity. Score level order and ordered question schemas are **not** reordered. |
| Replace config and accepted baseline together without approval | Candidate-generated history is not accepted authority; preserve the existing trusted-baseline workflow. |

**Ownership differs by option.** A/C give each site an explicit local ID and include the request-and-policy subtree in the enclosing function's revision identity. B gives that subtree its own question declaration and revision. For every option, the proposed request-authority closure of a function declaring `Jev.request` includes the effective judgment contracts it can reach, including through helpers. A changed closure requires a revision bump for the affected request-bearing function; pure consumers of ordinary result data do not acquire that obligation. This is an **explicit extension to a77**, not a claim that ordinary bodies already participate in interface hashing. It intentionally makes a global model change a potentially broad review event. ([a77] §Fingerprint model; [JEV1] §4.5)

Missing providers and stale pins report before dependent identity comparisons; suppress derivative errors until the direct authority fault is fixed. The accepted lock format must be bumped rather than silently comparing unlike fingerprint formats.

**The identity boundary is syntactic, not a detector of intent.** Declared instructions and criteria are request metadata; arbitrary text assembled as runtime state is still input data. State-preparation body edits retain a77's proof/evidence-freshness treatment. The compiler cannot discover an English instruction hidden in a user message and retroactively treat it as declaration metadata. This preserves the limitation already stated by JEV-1.0 rather than claiming stronger prompt authority. ([JEV1] §4.5; [a77] §Fingerprint model)


### 2.3 A failure-only underscore, not a new catch-all

**Draft amendment — `JEV-S1/FAILURE`, R5/R6:**

> In a primitive judgment dispatch, or the outer failure position of a judgment scope, the sole terminal `_ => expression` denotes exactly the closed intrinsic failure set `jev.unavailable` and `jev.invalid_response` from the pinned adapter contract. It is not a residual success pattern. Successful answer-state exhaustiveness is checked independently of this clause. The expression has the enclosing result type and may raise only the enclosing function's declared errors. All ordinary call, Boolean, value, and variant matching rules remain unchanged.

This is a narrow syntactic exception backed by a closed sum. The compiler does **not** create an open-ended `catch`. Removing `low` and hoping `_` handles it remains a missing-success-arm error, `CAN4101`. Adding a new service failure kind changes the adapter/profile contract, so an old underscore cannot quietly absorb an unreviewed new kind.

| Reaches the failure-only `_` | Does not reach it |
|---|---|
| Missing/invalid runtime credential; connection, TLS, body-read, timeout, or redirect refusal; non-success HTTP status, including 401, 422, 429, and 529: `jev.unavailable`. | A valid low-confidence answer, Choice tie, or uncertain Noul: successful states with their own arms. |
| Invalid successful JSON, wrong model/question/type, missing or duplicate support, invalid probabilities/confidence/mean, or mismatched legend: `jev.invalid_response`. | A malformed successful **fixture**: compile-time evidence diagnostic, never injected failure. |
| Either complete failure value explicitly scripted in `given`. | Missing config, stale pin, malformed source, exhausted/mismatched/leftover script, proof refusal, or uncovered source arm: compiler diagnostics. |
| Failure of this particular request's acquisition/validation. | Errors from an application call inside a successful arm, a later nested judgment, or the failure arm itself: ordinary local exhaustive handling. |

HTTP failures are normalized at the boundary; raw server messages and credentials do not enter application values. The status vocabulary is fixed by the adapter. Missing runtime credentials differ from missing build metadata: only the former is a runtime receipt failure. ([TS-API] §Errors; [JEV1] §§4.2–4.3)

All examples use `_ => Ok("failed")` where the enclosing return is a string or route record. An application may instead write `_ => triage.failed()` after declaring that domain error in `emits`. The underscore does not impose a sentinel string or convert failure into success automatically.

**Coverage:** the underscore is **one source arm**, with one CAN4107 obligation. A green row scripting either admitted failure executes it. The generated two-kind expansion does not create two hidden source-arm obligations. The example tables nevertheless test both kinds. Adapter-maintainer suites separately test every failure normalization path. No automatic certificate is granted to a failure handler. ([a03] §Hit semantics / retrospective amendment; [JEV1] §§3.3, 4.1)

### 2.4 Compact evidence, not reduced evidence

All options use the same proposed literal-only observation notation inside the existing exchange grammar:

```can
observed noul(noul = d"0.9")
observed choice(choice = "billing", probabilities = ("billing" = d"0.9", "other" = d"0.1"), confidence = d"0.8")
observed score(score = d"1.3", probabilities = Seq<dec>[d"0", d"0.7", d"0.3"], confidence = d"0.54")
```

These are three **complete evidence literals**, not calls, user-created brands, abbreviated stubs, or an inferred complement. A batch uses an `observed answers` literal containing one complete primitive literal per named question field. An error still uses its complete error construction. At a single-question site the successful logical record is flattened; the fixed wire key is `answer`. At a batch site all explicitly declared question IDs are present. Static request text is bound by the revision fingerprint, not recopied per exchange.

This replaces hand-written observation-record boilerplate with declaration-derived closed shapes. The kind is explicit in the literal and in the question; option support and level order are explicit in the schema. The compiler specializes that known schema just as it already knows a declared call's result binder. It does not infer missing probabilities, confidence, labels, fields, or semantic truth. `observed` is legal only in judgment `given` outcomes; it does not create a second production way to perform a judgment. ([JEV1] §§4.2–4.3; [R8]; [AUTH])

Choice support maps use quoted option keys, so labels need not be can identifiers. An observation binder exposes `.choice`, `.confidence`, and a closed `.probabilities["literal option"]` projection. Score exposes `.score`, `.confidence`, and its complete `Seq<dec>` probabilities; a literal in-range index is available directly. This proposal introduces no general map type, open string lookup, or user-facing generic observation machinery.

Keep JEV-1.0's numeric rules unchanged after the fetch boundary: finite JSON numbers use the specified number-to-canonical-dec image; probability mass allowance is `d"0.000000000001"`; no renormalization; selected Choice label is a maximizer; Score mean consistency uses the declared level indexes; reported confidence is range-checked, not recomputed. A malformed live response becomes a failure; a malformed success fixture is rejected before execution. ([JEV1] §4.3)

Queue identity and consumption are unchanged: one reached request consumes one exchange; missing-and-reached or extra-and-unconsumed fails with CAN4200, including zero visits to a nonempty applicable queue. Whole-batch validation happens before any answer is ignored. Neither model output frequencies nor confidence bands are scripts. ([JEV1] §4.2 / §8.1)

### 2.5 Shared grammar conventions and retained partitions

A description is a literal string, an `object` with explicitly named description fields, an ordered `array` of descriptions, or metadata-only `unspecified`. For example, `object("meaning" = "Billing request", "examples" = array["Refund missing", "Wrong invoice"])` is one complete description tree. Object keys are explicit, unique strings; Score `levels` is always an ordered enumeration. Descriptions are closed compile-time trees. Description-object and state-object keys use a fixed Unicode-scalar key ordering in canonical encoding. Arrays retain their order; explicit question, Choice-support, and Score-level declaration order remains identity-bearing. The emitter uses ordered entries or null-prototype dictionaries and a specified serializer—not ambient host object enumeration—to preserve valid numeric-looking and `__proto__` keys without prototype mutation. A runtime state tree uses the same explicit object/array layout but permits in-scope unbranded text leaves and has the distinct intrinsic type `Jev__State`. The compiler distinguishes those contexts; prompts cannot accidentally capture runtime values.

Noul uses `probability no L yes U` with `0 <= L < 0.5 < U <= 1`. `p <= L` is no, `p >= U` is yes, and the open interval is uncertain. There is **no strictly two-state Noul match** under this settlement. A minimal yes/no *application result* must explicitly decide what uncertainty produces.

Choice/Score use `confidence high H medium M`, with `0 < M < H <= 1`. Confidence bands are computed as in JEV-1.0. Choice ties enter low before threshold comparisons, including ties with reported confidence 1. Score never becomes an argmax primitive. Policy metadata has no defaults. ([JEV1] §4.4)

The state alternative `medium | low _` means either listed successful band, binding/discarding the same observation. `no | uncertain _` similarly makes an explicit application choice about uncertainty. For A/B an observation binder is mandatory, and `_` discards it. For C the declared question name already binds the observation, so arm patterns need no second binder. State alternatives do not include the failure-only underscore. One grouped source arm has one coverage obligation, not one per alternative; boundary rows should still exercise each listed state. Existing shadow/useless-alternative checking extends to these closed patterns.

**How to read the specimens:** the short blocks are complete judgment expressions with `message: str` in scope and all request metadata and result arms shown. String routes return `Ok(str)`; the three minimal numeric grading specimens explicitly state their numeric result and declared error context. They are expression specimens, not claims to be standalone compiler-green modules without tables. Each option then provides a complete module, including all declarations, tests, and full exchanges. All observations are synthetic and make no assertion about what a live model would answer.


## 3. Option A — Primitive-led expressions

### Pitch and owner fit

Make judgments expression heads, alongside ordinary `match`: **`noul` asks, `choice` routes, `score` grades**. A local site ID names the request for evidence and identity; the question and routing policy sit directly beside its arms. `together` is a distinct multi-question expression, never an optimizer transform. This follows all three owner leans: one inherited transport profile, one bounded failure underscore, and visibly different primitive forms. It departs from the declaration-first convention: reusable prompts must either be duplicated or wrapped in an ordinary function. It optimizes local reading, not contract reuse. ([R] Goal / R9; [IDIOMS] §Calls dispatch, never sequence / §Names, surface, revisions)

### A.1 Config and identity

Use the exact `jev.config.json` in §2.1 and `uses [Jev__Runtime@1]` once in the owning module. No site or descriptor carries model, endpoint, timeout, or retries. Each site ID must be unique within its function. The function's revision fingerprint includes its ordered `(site ID, primitive, state schema, question, criteria, policy, transport identity)` entries, excluding tables. Moving an unchanged site physically is inert; changing its ID is a reviewed identity change. Changing prompt words requires the containing function's revision bump. Profile changes follow the complete §2.2 matrix.

### A.2 Noul: minimal yes/no application result

This deliberately maps uncertain evidence to the application's `"no"` result; it does not erase the uncertain state.

```can
noul urgency from message
  asks "Does this message explicitly request help from a person?"
  probability no d"0.2" yes d"0.8"
  yes _ => Ok("yes")
  no | uncertain _ => Ok("no")
  _ => Ok("failed")
```

Noul's probability-only result and optional criteria come from ([TS-N] §Writing a Noul question / §Response / §Noul does not return a separate confidence value). These same receipts apply to A.3–A.4.

### A.3 Noul: structured instructions and true/false criteria

```can
noul human_help from object("ticket" = object("message" = message))
  asks object("question" = "Does `ticket.message` explicitly request a person?", "focus" = "Distinguish a human request from general frustration.")
  criteria
    true = object("meaning" = "The sender asks for a human agent.", "examples" = array["Connect me to a person.", "I want to speak to your staff."])
    false = object("meaning" = "No human contact is requested.", "examples" = array["The page is slow.", "Where is the help article?"])
  probability no d"0.2" yes d"0.8"
  yes _ => Ok("human_queue")
  no _ => Ok("self_service")
  uncertain _ => Ok("clarify_request")
  _ => Ok("failed")
```

The object state is constructed explicitly and sends only the named text. The arbitrary metadata keys `question`, `focus`, `meaning`, and `examples` are data, not reserved API keys. ([TS-A] §Structured instructions / §Structured Noul criteria; [TS-STATE] §State can be as simple as a string; [TS-P] §Reference specific fields)

### A.4 Noul: visibly separate uncertainty

```can
noul urgency from array["Evaluate the next message as the sender's request.", message]
  asks "Does the sender say this must be handled before today's end?"
  probability no d"0.1" yes d"0.9"
  yes _ => Ok("today")
  no _ => Ok("normal_queue")
  uncertain _ => Ok("ask_for_deadline")
  _ => Ok("failed")
```

Exactly 0.5 and both nearby values 0.49/0.51 execute `uncertain`; neither no nor failure is a fallthrough destination. This is the preserved can-lang partition, not TypeSafe confidence. ([JEV1] §4.4)

### A.5 Choice: minimal routing

```can
choice topic from message
  asks "Which queue should handle the message?"
  options
    "billing" = unspecified
    "other" = unspecified
  confidence high d"0.8" medium d"0.5"
  high "billing" _ => Ok("billing")
  high "other" _ => Ok("general")
  medium | low _ => Ok("review")
  _ => Ok("failed")
```

`unspecified` emits an explicit null description for a label-only option; it does not infer a description from a variable name. ([TS-C] §A more complex example; [TS-A] §Where structure is allowed)

### A.6 Choice: full confidence routing and structured option descriptions

```can
choice topic from message
  asks array["Choose the team that owns the sender's primary request.", "Use other when billing is not the primary request."]
  options
    "billing" = object("covers" = "Charges, invoices, and refunds.", "excludes" = "Login problems.", "examples" = array["The invoice is wrong.", "My refund is missing."])
    "other" = object("covers" = "Requests outside billing.", "examples" = array["I cannot sign in.", "How do I change the theme?"])
  confidence high d"0.8" medium d"0.5"
  high "billing" _ => Ok("billing")
  high "other" _ => Ok("general")
  medium _ => Ok("confirm")
  low _ => Ok("review")
  _ => Ok("failed")
```

This exposes array instructions and object-valued Choice criteria without embedding JSON source text. ([TS-A] §Structured instructions / §Structured Choice options; [TS-C] §Structured instructions and criteria; [TS-F] §Three paths for using confidence in your code)

**Tie:** a complete fixture with equal maximal probabilities takes `low`, even when it names a legal maximizer and reports confidence 1. The complete `tie` row in A.7 demonstrates precisely that. There is no first-option tie rule and no `tie` arm to put beside `low` while forgetting that it is a subset of low. ([JEV1] §§4.3–4.4)

### A.7 One request, one ignored answer, complete tests and given

This is a complete proposed module. The batch declares both questions explicitly; all eight selected success/failure scenarios are scripted without abbreviation. `urgency = _` disregards a **received and validated** answer, not an unperformed question.


```can
mod triage
  provides [triage__route, Triage__Route]
  uses [Jev__Runtime@1]
  emits []

type Triage__Route rev 1 (
  destination: str
)

fn triage__route(message: str) -> Triage__Route rev 1
  emits []
  effects [Jev.request]
  tests
    billing("charged twice") => Ok("billing") pinned
    other("charged twice") => Ok("general") pinned
    high_boundary("charged twice") => Ok("billing") pinned
    medium_boundary("charged twice") => Ok("confirm") pinned
    low("charged twice") => Ok("review") pinned
    tie("charged twice") => Ok("review") pinned
    offline("charged twice") => Ok("failed") pinned
    malformed("charged twice") => Ok("failed") pinned
  together triage_request from message
    questions
      topic: choice
        asks "Which queue should handle the message?"
        options
          "billing" = "Charges, invoices, and refunds."
          "other" = "Requests outside billing."
        confidence high d"0.8" medium d"0.5"
      urgency: noul
        asks "Does the sender request help before today's end?"
        probability no d"0.2" yes d"0.8"
    given
      billing => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"1", "other" = d"0"), confidence = d"1"), urgency = observed noul(noul = d"0.9"))]
      other => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "other", probabilities = ("billing" = d"0", "other" = d"1"), confidence = d"1"), urgency = observed noul(noul = d"0.1"))]
      high_boundary => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.9", "other" = d"0.1"), confidence = d"0.8"), urgency = observed noul(noul = d"0.5"))]
      medium_boundary => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.8", "other" = d"0.2"), confidence = d"0.5"), urgency = observed noul(noul = d"0.5"))]
      low => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.6", "other" = d"0.4"), confidence = d"0.49"), urgency = observed noul(noul = d"0.5"))]
      tie => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.5", "other" = d"0.5"), confidence = d"1"), urgency = observed noul(noul = d"0.9"))]
      offline => [exchange args (state = "charged twice") outcome jev.unavailable("network")]
      malformed => [exchange args (state = "charged twice") outcome jev.invalid_response("protocol")]
    on answers (topic = high "billing" _, urgency = _) => Ok("billing")
    on answers (topic = high "other" _, urgency = _) => Ok("general")
    on answers (topic = medium _, urgency = _) => Ok("confirm")
    on answers (topic = low _, urgency = _) => Ok("review")
    _ => Ok("failed")
```


Eight test rows, eight complete exchanges, five source arms. The second question adds **one Noul field to each successful exchange**, not a second table or transport handler. The shared schema supplies only static type information; no numerical observation is filled in. The metadata spelling `together` guarantees one request, with every question seeing the same state. ([TS-P] §Ask multiple questions together / §Ask speculative questions; [TS-FAN]; [JEV1] §8.1)

### A.8 Score: minimal grading

Numeric expression context: return `dec`, declare `error grade.failed()` and include `grade.failed` in `emits`.

```can
error grade.failed()

score completeness from message
  asks "How much diagnostic detail does the report contain?"
  levels ["Names no feature.", "Names the feature but gives no reproduction steps.", "Names the feature and gives reproduction steps."]
  confidence high d"0.8" medium d"0.5"
  high | medium | low estimate => Ok(estimate.score)
  _ => grade.failed()
```

This reports the observed mean in every successful band. It intentionally ignores confidence **for reporting**, not for authorization. Classification still occurs and remains testable. ([TS-S] §Levels / §Response structure)

### A.9 Score: level-routed branch

```can
score severity from message
  asks "How much does the reported problem obstruct the user?"
  levels ["Appearance is wrong but work continues.", "Work is disrupted but a workaround exists.", "Work cannot continue and no workaround is described."]
  confidence high d"0.8" medium d"0.5"
  high estimate =>
    match estimate.score < d"0.5", estimate.score < d"1.5"
      true, _ => Ok("cosmetic")
      false, true => Ok("degraded")
      false, false => Ok("blocking")
  medium _ => Ok("confirm_severity")
  low _ => Ok("inspect_report")
  _ => Ok("failed")
```

The application bins the **reported mean** at explicit 0.5 and 1.5 boundaries. Exactly 0.5 is degraded; exactly 1.5 is blocking. No hidden rounding or highest-probability level is inserted. The inner table follows the existing pure multi-scrutinee idiom. ([TS-S] §Reading a Score; [IDIOMS] §C1)

### A.10 Score: inspect shape, not just the mean

```can
score severity_shape from message
  asks "How much does the reported problem obstruct the user?"
  levels [object("meaning" = "Appearance is wrong but work continues.", "examples" = array["The icon is misaligned."]), object("meaning" = "Work is disrupted but a workaround exists.", "examples" = array["Printing fails but PDF export works."]), object("meaning" = "Work cannot continue.", "examples" = array["The application will not start."])]
  confidence high d"0.8" medium d"0.5"
  high | medium | low estimate =>
    match estimate.probabilities[1] >= d"0.8"
      true => Ok("middle_level_concentrated")
      false => Ok("middle_level_not_concentrated")
  _ => Ok("failed")
```

Complete counterfactual observations distinguishing identical means are:

```can
observed score(score = d"1", probabilities = Seq<dec>[d"0", d"1", d"0"], confidence = d"1")
observed score(score = d"1", probabilities = Seq<dec>[d"0.5", d"0", d"0.5"], confidence = d"0")
```

The first takes the inner true arm; the second takes false. Structured level descriptions retain their order and are checked against the returned legend. ([TS-S] §Reading a Score / §Structured level descriptions)

### A.11 Failure, coverage, and emit

A.7's `_` handles only this request's two intrinsic failure kinds. Neither low nor uncertain is included. `offline` and `malformed` execute that same source arm. `medium_boundary`, `low`, and `tie` exercise their computed success bands; positive probability on `other` in a billing fixture does not cover the other high-label arm. A.9's three numeric arms need witnesses with means in their respective intervals. A.10 needs both center-concentrated and nonconcentrated observations. These are ordinary source-arm obligations, not coverage over every possible probability vector. ([JEV1] §§4.1–4.4; [a03])

The emitted TypeScript reads the single profile's constants, sends `state`, `model`, and the declared questions through the common fetch boundary in §6, validates the whole receipt, and switches on compiler-derived tags. Each inline site becomes a source-mapped generated request schema; no SDK, test table, or model-selection clause is emitted at the use site. A minimal Noul lowering is:

```ts
const receipt = await $jevFetch(message, {
  answer: {
    type: "noul",
    instructions: "Does this message explicitly request help from a person?"
  }
});
switch (receipt.kind) {
  case "jev.unavailable": return { kind: "Ok", value: "failed" };
  case "jev.invalid_response": return { kind: "Ok", value: "failed" };
  case "answers":
    switch ($noulBand(receipt.answers.answer, "0.2", "0.8")) {
      case "yes": return { kind: "Ok", value: "yes" };
      case "no": return { kind: "Ok", value: "no" };
      case "uncertain": return { kind: "Ok", value: "no" };
    }
}
```

`$jevFetch` and `$noulBand` stand for generated, schema-specialized code, not package imports. For scalar returns the emitted payload field is `value`; the complete module's record instead uses `destination`.

**Pros:** smallest conceptual distance from native control flow; complete local context; explicit single-versus-batched acquisition; no independent descriptor ceremony. **Cons:** substantial inline rubrics interrupt business logic; prompt reuse risks duplication; prompt edits and site renames affect the containing function revision; weaker declaration-only retrieval.

**Semantics preserved:** all JEV-1.0 evidence, numeric, routing, failure identity, coverage, and proof boundaries. **Amendments requested:** common `JEV-S1/TRANSPORT`, `/FAILURE`, `/OBS-LITERAL`, `/STRUCTURE`, `/IDENTITY`; option-specific `JEV-S1/A` introduces the primitive expression heads, explicit local site IDs, and `together`. New option codes are `CAN3620` for a duplicate local judgment-site identity and `CAN3621` for an inline primitive outside its whole-expression dispatch position. No general catch-all or sampling exception is requested.


## 4. Option B — Question signatures

### Pitch and owner fit

Give each primitive a native declaration kind, then dispatch it with the familiar **`match noul`**, **`match choice`**, or **`match score`**. A declaration says what is judged and how uncertainty is partitioned; it says nothing about credentials, fetch, retries, or service-error plumbing. This follows all three owner leans while preserving declaration-first reading and explicit revision pins. The deliberate tradeoff is that the smallest use site relies on a nearby or imported question signature. It is not an extern in disguise: the compiler owns the observation schema, routing, evidence consumption, and exhaustive judgment arms. ([R] Goal / R9; [IDIOMS] §Calls dispatch, never sequence / §Names, surface, revisions)

### B.1 Config, metadata, and pinning

Use the exact §2.1 `jev.config.json`. The question-owning module writes `uses [Jev__Runtime@1]` once. Its native question declarations inherit that profile. A consumer imports a question with the existing style, for example `uses [triage__human_help@1]`; it does not repeat the provider's model, timeout, or credentials. If it defines its own questions, it also pins the same project profile explicitly.

`noul name(state: T) rev N`, `choice name(state: T) rev N`, and `score name(state: T) rev N` are declaration kinds with **fixed primitive result schemas**, not functions with inferred return annotations. `T` is explicitly `str` or `Jev__State`. The only request parameter is named `state`. The primitive keyword determines the observation fields; options/levels determine their closed support. These declarations have no executable body, no independent tests, and no per-descriptor `emits` or `effects` rows. Enclosing functions retain explicit `effects [Jev.request]` and their own `emits`. ([JEV1] §§4.2, 4.5)

Prompt wording and policy changes bump the question revision. Its transport identity and any containing request-bearing function's authority closure are versioned as in §2.2. A pinned consumer encountering a changed question revision receives CAN2103 until deliberately repinned. Same-revision prompt/model drift receives CAN6013, not a warning hidden in a lock refresh.

### B.2 Noul: minimal declaration and yes/no application result

```can
noul triage__human_help(state: str) rev 1
  instructions "Does this message explicitly request help from a person?"
  probability no d"0.2" yes d"0.8"

match noul triage__human_help(message)
  on yes _ => Ok("yes")
  on no | uncertain _ => Ok("no")
  _ => Ok("failed")
```

The declaration is reusable; the use site is four lines. Noul remains a probability-valued question, and uncertainty's mapping is explicit. ([TS-N] §Writing a Noul question / §Noul does not return a separate confidence value)

### B.3 Noul: structured true/false criteria

```can
noul triage__structured_help(state: Jev__State) rev 1
  instructions object("question" = "Does `ticket.message` explicitly request a person?", "focus" = "Do not count frustration alone as a human-contact request.")
  criteria
    true = object("meaning" = "The sender requests a human agent.", "examples" = array["Please connect me to a person.", "I need to talk to your staff."])
    false = object("meaning" = "No human contact is requested.", "examples" = array["The screen is confusing.", "Where is the tutorial?"])
  probability no d"0.2" yes d"0.8"

match noul triage__structured_help(object("ticket" = object("message" = message)))
  on yes _ => Ok("human_queue")
  on no _ => Ok("self_service")
  on uncertain _ => Ok("clarify_request")
  _ => Ok("failed")
```

The state type is explicit, and the constructed tree is the entire disclosed input. The literal instruction's backticked path is not evaluated by can-lang. ([TS-A] §Structured instructions / §Structured Noul criteria; [TS-STATE]; [TS-P] §Reference specific fields)

### B.4 Noul: explicit uncertainty

```can
noul triage__deadline(state: Jev__State) rev 1
  instructions "Does the sender say this must be handled before today's end?"
  probability no d"0.1" yes d"0.9"

match noul triage__deadline(array["Evaluate the next message as the sender's request.", message])
  on yes _ => Ok("today")
  on no _ => Ok("normal_queue")
  on uncertain _ => Ok("ask_for_deadline")
  _ => Ok("failed")
```

The explicit uncertain arm is mandatory unless it is explicitly grouped with another successful state. `_` never substitutes for it. At the shown bounds, 0.1 is no, 0.9 is yes, and 0.5 is uncertain. ([JEV1] §4.4)

### B.5 Choice: minimal routing

```can
choice triage__queue(state: str) rev 1
  instructions "Which queue should handle the message?"
  options
    "billing" = unspecified
    "other" = unspecified
  confidence high d"0.8" medium d"0.5"

match choice triage__queue(message)
  on high "billing" _ => Ok("billing")
  on high "other" _ => Ok("general")
  on medium | low _ => Ok("review")
  _ => Ok("failed")
```

Label-only options retain an explicit null description at the wire boundary. Probabilities and reported confidence remain separate data. ([TS-C] §A more complex example / §Response structure; [TS-A] §Where structure is allowed)

### B.6 Choice: high/medium/low with structured metadata

```can
choice triage__queue_detail(state: str) rev 1
  instructions array["Choose the team that owns the primary request.", "Use other when billing is not the primary request."]
  options
    "billing" = object("covers" = "Charges, invoices, and refunds.", "excludes" = "Login problems.", "examples" = array["The invoice is wrong.", "My refund is missing."])
    "other" = object("covers" = "Requests outside billing.", "examples" = array["I cannot sign in.", "How do I change the theme?"])
  confidence high d"0.8" medium d"0.5"

match choice triage__queue_detail(message)
  on high "billing" _ => Ok("billing")
  on high "other" _ => Ok("general")
  on medium _ => Ok("confirm")
  on low _ => Ok("review")
  _ => Ok("failed")
```

The syntax lets the reader inspect the policy without seeing a fetch contract. A higher bound for one action can be an ordinary, exhaustively tested comparison against the bound observation's confidence inside that action's arm; it is not a different provider statistic. ([TS-A] §Structured Choice options; [TS-F] §Thresholds scale with risk; [TS-R] §Step 2)

**Tie:** the `tie` row below has two exact maximizers and reported confidence 1. It must take low, never either high label. A `tie` branch is not added to the language: ties remain part of the existing low domain. ([JEV1] §4.4)

### B.7 Explicit batch declaration; one complete table

A `judgments` declaration is a **schema composition**, not a sequence of function calls. Each member names a primitive declaration with the same sole `state` parameter type. All members receive the same state in one request. Local references resolve locally; foreign question declarations must appear as revision-pinned `uses` entries. Member bindings are named and ordered. No hidden argument mapping or automatic coalescing is permitted.


```can
mod triage
  provides [triage__topic, triage__urgency, triage__overview, triage__route, Triage__Route]
  uses [Jev__Runtime@1]
  emits []

type Triage__Route rev 1 (
  destination: str
)

choice triage__topic(state: str) rev 1
  instructions "Which queue should handle the message?"
  options
    "billing" = "Charges, invoices, and refunds."
    "other" = "Requests outside billing."
  confidence high d"0.8" medium d"0.5"

noul triage__urgency(state: str) rev 1
  instructions "Does the sender request help before today's end?"
  probability no d"0.2" yes d"0.8"

judgments triage__overview(state: str) rev 1
  topic: choice triage__topic
  urgency: noul triage__urgency

fn triage__route(message: str) -> Triage__Route rev 1
  emits []
  effects [Jev.request]
  tests
    billing("charged twice") => Ok("billing") pinned
    other("charged twice") => Ok("general") pinned
    high_boundary("charged twice") => Ok("billing") pinned
    medium_boundary("charged twice") => Ok("confirm") pinned
    low("charged twice") => Ok("review") pinned
    tie("charged twice") => Ok("review") pinned
    offline("charged twice") => Ok("failed") pinned
    malformed("charged twice") => Ok("failed") pinned
  match judgments triage__overview(message)
    given
      billing => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"1", "other" = d"0"), confidence = d"1"), urgency = observed noul(noul = d"0.9"))]
      other => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "other", probabilities = ("billing" = d"0", "other" = d"1"), confidence = d"1"), urgency = observed noul(noul = d"0.1"))]
      high_boundary => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.9", "other" = d"0.1"), confidence = d"0.8"), urgency = observed noul(noul = d"0.5"))]
      medium_boundary => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.8", "other" = d"0.2"), confidence = d"0.5"), urgency = observed noul(noul = d"0.5"))]
      low => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.6", "other" = d"0.4"), confidence = d"0.49"), urgency = observed noul(noul = d"0.5"))]
      tie => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.5", "other" = d"0.5"), confidence = d"1"), urgency = observed noul(noul = d"0.9"))]
      offline => [exchange args (state = "charged twice") outcome jev.unavailable("network")]
      malformed => [exchange args (state = "charged twice") outcome jev.invalid_response("protocol")]
    on answers (topic = high "billing" _, urgency = _) => Ok("billing")
    on answers (topic = high "other" _, urgency = _) => Ok("general")
    on answers (topic = medium _, urgency = _) => Ok("confirm")
    on answers (topic = low _, urgency = _) => Ok("review")
    _ => Ok("failed")
```


The evidence bytes and routes are the same as A.7. Declaration factoring does not demand a second fixture format or a new file of observations. The ignored urgency answer is still complete in every successful exchange; a missing urgency field is CAN3501, not permission to skip that question. ([JEV1] §4.3 / §8.1; [TS-P] §Ask speculative questions)

### B.8 Score: minimal grading

Numeric expression context: return `dec`, declare `error grade.failed()` and include `grade.failed` in `emits`.

```can
error grade.failed()

score grade__completeness(state: str) rev 1
  instructions "How much diagnostic detail does the report contain?"
  levels ["Names no feature.", "Names the feature but gives no reproduction steps.", "Names the feature and gives reproduction steps."]
  confidence high d"0.8" medium d"0.5"

match score grade__completeness(message)
  on high | medium | low estimate => Ok(estimate.score)
  _ => grade.failed()
```

The numeric result is explicitly the received mean, independent of whether it lies exactly on a level. The listed band union intentionally reports it without treating it as a certified grade. ([TS-S] §Levels / §Reading a Score)

### B.9 Score: level-routed branch

```can
score grade__severity(state: str) rev 1
  instructions "How much does the reported problem obstruct the user?"
  levels ["Appearance is wrong but work continues.", "Work is disrupted but a workaround exists.", "Work cannot continue and no workaround is described."]
  confidence high d"0.8" medium d"0.5"

match score grade__severity(message)
  on high estimate =>
    match estimate.score < d"0.5", estimate.score < d"1.5"
      true, _ => Ok("cosmetic")
      false, true => Ok("degraded")
      false, false => Ok("blocking")
  on medium _ => Ok("confirm_severity")
  on low _ => Ok("inspect_report")
  _ => Ok("failed")
```

Levels remain an ordered rubric. The program—not the primitive—chooses the bins; there is no automatic conversion from a distribution to its most likely level. ([TS-S] §Levels / §Reading a Score)

### B.10 Score: structured levels and shape inspection

```can
score grade__severity_shape(state: str) rev 1
  instructions "How much does the reported problem obstruct the user?"
  levels [object("meaning" = "Appearance is wrong but work continues.", "examples" = array["The icon is misaligned."]), object("meaning" = "Work is disrupted but a workaround exists.", "examples" = array["Printing fails but PDF export works."]), object("meaning" = "Work cannot continue.", "examples" = array["The application will not start."])]
  confidence high d"0.8" medium d"0.5"

match score grade__severity_shape(message)
  on high | medium | low estimate =>
    match estimate.probabilities[1] >= d"0.8"
      true => Ok("middle_level_concentrated")
      false => Ok("middle_level_not_concentrated")
  _ => Ok("failed")
```

These complete observations give the same mean but execute opposite inner arms:

```can
observed score(score = d"1", probabilities = Seq<dec>[d"0", d"1", d"0"], confidence = d"1")
observed score(score = d"1", probabilities = Seq<dec>[d"0.5", d"0", d"0.5"], confidence = d"0")
```

Structured legends are checked as description trees, not coerced to strings. ([TS-S] §Structured level descriptions; [TS-A] §Structured Score levels)

### B.11 Failure, coverage, and emit

The wildcard's boundary is the named judgment request, **not all operations in the enclosing match arm**. Reusing a question does not reuse a response: two reached `match choice` sites consume two independent exchange queues and make two runtime requests. Only a `judgments` declaration groups questions. Coverage remains attached to each source use site's arms; a covered use of a question does not cover a different use. B.7 has the same five source-arm witnesses as A.7. ([JEV1] §§4.2, 4.4, 8.1; [a03])

The compiler emits the named question as static request metadata and its policy as deterministic routing code. At the use site it awaits the shared fetch boundary, validates the observation, and switches on a finite derived decision. No external SDK or dynamic config loader appears. A Choice use can lower to:

```ts
const receipt = await $jevFetch(message, {
  answer: {
    type: "choice",
    instructions: "Which queue should handle the message?",
    criteria: { billing: null, other: null }
  }
});
const decision = $triageQueueDecision(receipt, "0.8", "0.5");
switch (decision.kind) {
  case "high.billing": return { kind: "Ok", value: "billing" };
  case "high.other": return { kind: "Ok", value: "general" };
  case "medium": return { kind: "Ok", value: "review" };
  case "low": return { kind: "Ok", value: "review" };
  case "jev.unavailable": return { kind: "Ok", value: "failed" };
  case "jev.invalid_response": return { kind: "Ok", value: "failed" };
}
```

`$triageQueueDecision` is generated for the exact declared support and implements tie-to-low before confidence classification. Its error tags remain available for catalog attribution even though the source uses one handler.

**Pros:** shortest ordinary sites; primitive kind visible at the match; versioned reusable rubrics; preserves declaration-first review; batch composition reuses those exact contracts. **Cons:** more declarations for one-off questions; two locations to read; changes to shared questions have real consumer-pin consequences; the bare batch match is less expressive about acquisition versus later decisions than C.

**Semantics preserved:** all settled JEV-1.0 evidence and routing semantics. **Amendments requested:** the common `JEV-S1` rules plus `JEV-S1/B` for native primitive declarations, `match` scrutinee classes, and non-executable batch composition. New codes: `CAN3630` for a primitive mismatch between syntax and referenced declaration; `CAN3631` for a batch member whose sole state signature does not match the batch. Existing provides/uses, revision, and source-arm diagnostics extend to these declaration kinds. No wrapper-specific transport exception remains in source.


## 5. Option C — Judgment scopes

### Pitch and owner fit

Make the primitive an explicit **acquisition scope**: `judge noul`, `judge choice`, or `judge score` obtains a named, immutable observation; `decide` dispatches its compiler-computed state. The outer underscore handles acquisition failure only. Multiple questions share one scope, and later decisions can read the same answers without another request. This follows the one-config and one-failure-path leans, while giving each primitive a distinct header. It departs from the absolute-minimum-line lean: the separate `decide` line costs space but makes the boundary between asking and acting visible. Unlike A's immediate dispatch and B's reusable declarations, this option gives a successful request a lexical observation scope.

### C.1 Config, identity, and scope rules

Use §2.1's exact `jev.config.json` and the module's single `uses [Jev__Runtime@1]`. A `judge` site inherits the same transport contract as A/B. Each site has a function-local explicit ID. The enclosing function revision owns its question metadata, question order, policy, and effective profile identity; the complete §2.2 revision matrix applies. No anonymous or runtime-generated request identity is admitted.

A single-question scope declares its observation name in the header. A multi-question scope declares each observation name in `questions`. These names are typed by their explicitly declared primitive and support, not by a guessed return type. They are available **only in the successful `decide` expression**, never in the outer failure clause or beyond the scope. `decide` is pure dispatch over an existing observation; it cannot fetch, consume another exchange, or own a `given` table. Calling another `judge` is a distinct effect with its own table.

### C.2 Noul: minimal yes/no application result

```can
judge noul urgency from message
  instructions "Does this message explicitly request help from a person?"
  probability no d"0.2" yes d"0.8"
  decide urgency
    yes => Ok("yes")
    no | uncertain => Ok("no")
  _ => Ok("failed")
```

`urgency.noul` remains available inside the successful continuation. Grouping uncertain with no is an explicit application policy, not a conversion of Noul to a Boolean primitive. ([TS-N] §Writing a Noul question / §Noul does not return a separate confidence value; [JEV1] §4.4)

### C.3 Noul: structured instructions and true/false criteria

```can
judge noul human_help from object("ticket" = object("message" = message))
  instructions object("question" = "Does `ticket.message` explicitly request a person?", "focus" = "Distinguish a human-contact request from general frustration.")
  criteria
    true = object("meaning" = "The sender asks for a human agent.", "examples" = array["Connect me to a person.", "I want to speak to your staff."])
    false = object("meaning" = "No human contact is requested.", "examples" = array["The page is slow.", "Where is the help article?"])
  probability no d"0.2" yes d"0.8"
  decide human_help
    yes => Ok("human_queue")
    no => Ok("self_service")
    uncertain => Ok("clarify_request")
  _ => Ok("failed")
```

Structured content stays metadata; the explicitly constructed text tree is the entire runtime state. The string path in the instruction does not capture a variable implicitly. ([TS-A] §Structured instructions / §Structured Noul criteria; [TS-STATE]; [TS-P] §Reference specific fields)

### C.4 Noul: explicit uncertain handling

```can
judge noul deadline from array["Evaluate the next message as the sender's request.", message]
  instructions "Does the sender say this must be handled before today's end?"
  probability no d"0.1" yes d"0.9"
  decide deadline
    yes => Ok("today")
    no => Ok("normal_queue")
    uncertain => Ok("ask_for_deadline")
  _ => Ok("failed")
```

For this policy, 0.1 executes no, 0.9 executes yes, and 0.5 executes uncertain. A service failure skips `decide` entirely and executes the outer underscore. Missing uncertain is still CAN4101; no error handler repairs it. ([JEV1] §4.4)

### C.5 Choice: minimal routing

```can
judge choice topic from message
  instructions "Which queue should handle the message?"
  options
    "billing" = unspecified
    "other" = unspecified
  confidence high d"0.8" medium d"0.5"
  decide topic
    high "billing" => Ok("billing")
    high "other" => Ok("general")
    medium | low => Ok("review")
  _ => Ok("failed")
```

The option map fixes the labels; null criterion descriptions are explicit. The observation `topic` retains the complete distribution and reported confidence even when only its classified state is matched. ([TS-C] §Request structure / §Response structure; [TS-A] §Where structure is allowed)

### C.6 Choice: high/medium/low and structured descriptions

```can
judge choice topic from message
  instructions array["Choose the team that owns the primary request.", "Use other when billing is not the primary request."]
  options
    "billing" = object("covers" = "Charges, invoices, and refunds.", "excludes" = "Login problems.", "examples" = array["The invoice is wrong.", "My refund is missing."])
    "other" = object("covers" = "Requests outside billing.", "examples" = array["I cannot sign in.", "How do I change the theme?"])
  confidence high d"0.8" medium d"0.5"
  decide topic
    high "billing" => Ok("billing")
    high "other" => Ok("general")
    medium => Ok("confirm")
    low => Ok("review")
  _ => Ok("failed")
```

A tie is a successful low decision before any confidence comparison. C.7's full `tie` exchange gives equal maxima and confidence 1, and still takes low. Optional first-option tie selection and a probability-as-confidence gate are not imported from Probably. ([JEV1] §§4.3–4.4; [TS-A] §Structured Choice options; [TS-F]; [PROBABLY] §§3–4)

### C.7 Multi-question scope, one ignored answer, full evidence

The following complete module receives two answers, deliberately ignores urgency in dispatch, and checks the same eight counterfactual rows as A/B. The inner `_` is a field-discard pattern; the outer `_ =>` is a request-failure clause. Their syntactic positions distinguish them.


```can
mod triage
  provides [triage__route, Triage__Route]
  uses [Jev__Runtime@1]
  emits []

type Triage__Route rev 1 (
  destination: str
)

fn triage__route(message: str) -> Triage__Route rev 1
  emits []
  effects [Jev.request]
  tests
    billing("charged twice") => Ok("billing") pinned
    other("charged twice") => Ok("general") pinned
    high_boundary("charged twice") => Ok("billing") pinned
    medium_boundary("charged twice") => Ok("confirm") pinned
    low("charged twice") => Ok("review") pinned
    tie("charged twice") => Ok("review") pinned
    offline("charged twice") => Ok("failed") pinned
    malformed("charged twice") => Ok("failed") pinned
  judge triage_snapshot from message
    questions
      topic: choice
        instructions "Which queue should handle the message?"
        options
          "billing" = "Charges, invoices, and refunds."
          "other" = "Requests outside billing."
        confidence high d"0.8" medium d"0.5"
      urgency: noul
        instructions "Does the sender request help before today's end?"
        probability no d"0.2" yes d"0.8"
    given
      billing => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"1", "other" = d"0"), confidence = d"1"), urgency = observed noul(noul = d"0.9"))]
      other => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "other", probabilities = ("billing" = d"0", "other" = d"1"), confidence = d"1"), urgency = observed noul(noul = d"0.1"))]
      high_boundary => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.9", "other" = d"0.1"), confidence = d"0.8"), urgency = observed noul(noul = d"0.5"))]
      medium_boundary => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.8", "other" = d"0.2"), confidence = d"0.5"), urgency = observed noul(noul = d"0.5"))]
      low => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.6", "other" = d"0.4"), confidence = d"0.49"), urgency = observed noul(noul = d"0.5"))]
      tie => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.5", "other" = d"0.5"), confidence = d"1"), urgency = observed noul(noul = d"0.9"))]
      offline => [exchange args (state = "charged twice") outcome jev.unavailable("network")]
      malformed => [exchange args (state = "charged twice") outcome jev.invalid_response("protocol")]
    decide topic, urgency
      high "billing", _ => Ok("billing")
      high "other", _ => Ok("general")
      medium, _ => Ok("confirm")
      low, _ => Ok("review")
    _ => Ok("failed")
```


The four inner rows exhaust every successful answer product, including all three urgency states. They do not need nine separate source arms merely because urgency can have three states. A malformed ignored answer still invalidates the whole response. One complete batch consumes one exchange, not one exchange per `decide`. ([TS-P] §Ask multiple questions together / §Ask speculative questions; [JEV1] §8.1)

**The scope-specific benefit:** this complete alternative expression reuses the received urgency result only when the topic is high billing. It still asks both questions in its one request; the inner decision is not a conditional second model call.

```can
judge triage_snapshot from message
  questions
    topic: choice
      instructions "Which queue should handle the message?"
      options
        "billing" = "Charges, invoices, and refunds."
        "other" = "Requests outside billing."
      confidence high d"0.8" medium d"0.5"
    urgency: noul
      instructions "Does the sender request help before today's end?"
      probability no d"0.2" yes d"0.8"
  decide topic
    high "billing" =>
      decide urgency
        yes => Ok("billing_today")
        no => Ok("billing_normal")
        uncertain => Ok("billing_clarify_deadline")
    high "other" => Ok("general")
    medium => Ok("confirm")
    low => Ok("review")
  _ => Ok("failed")
```

Each inner urgency arm needs its own reached witness. A no/uncertain urgency fixture paired with high other does not cover it. If deciding urgency required an instruction derived from the topic's answer, this batch would be wrong: that question requires a later explicit request. This distinction is about invocation and disclosure, not just indentation. ([TS-P] §When one question depends on another; [JEV1] §8.1)

### C.8 Score: minimal grading

Numeric expression context: return `dec`, declare `error grade.failed()` and include `grade.failed` in `emits`.

```can
error grade.failed()

judge score completeness from message
  instructions "How much diagnostic detail does the report contain?"
  levels ["Names no feature.", "Names the feature but gives no reproduction steps.", "Names the feature and gives reproduction steps."]
  confidence high d"0.8" medium d"0.5"
  decide completeness
    high | medium | low => Ok(completeness.score)
  _ => grade.failed()
```

The result is the reported weighted mean; the successful grouping explicitly declines to filter it by confidence. There is no omitted result on uncertainty or failure. ([TS-S] §Levels / §Reading a Score)

### C.9 Score: level-routed branch

```can
judge score severity from message
  instructions "How much does the reported problem obstruct the user?"
  levels ["Appearance is wrong but work continues.", "Work is disrupted but a workaround exists.", "Work cannot continue and no workaround is described."]
  confidence high d"0.8" medium d"0.5"
  decide severity
    high =>
      match severity.score < d"0.5", severity.score < d"1.5"
        true, _ => Ok("cosmetic")
        false, true => Ok("degraded")
        false, false => Ok("blocking")
    medium => Ok("confirm_severity")
    low => Ok("inspect_report")
  _ => Ok("failed")
```

The bins belong to ordinary application logic. The scope introduces neither rounding nor sampling, and the rubric's position indexes remain zero-based. ([TS-S] §Levels / §Reading a Score)

### C.10 Score: structured levels and distribution shape

```can
judge score severity from message
  instructions "How much does the reported problem obstruct the user?"
  levels [object("meaning" = "Appearance is wrong but work continues.", "examples" = array["The icon is misaligned."]), object("meaning" = "Work is disrupted but a workaround exists.", "examples" = array["Printing fails but PDF export works."]), object("meaning" = "Work cannot continue.", "examples" = array["The application will not start."])]
  confidence high d"0.8" medium d"0.5"
  decide severity
    high | medium | low =>
      match severity.probabilities[1] >= d"0.8"
        true => Ok("middle_level_concentrated")
        false => Ok("middle_level_not_concentrated")
  _ => Ok("failed")
```

These complete observations share mean 1 but select opposite shape branches:

```can
observed score(score = d"1", probabilities = Seq<dec>[d"0", d"1", d"0"], confidence = d"1")
observed score(score = d"1", probabilities = Seq<dec>[d"0.5", d"0", d"0.5"], confidence = d"0")
```

The explicit level descriptions are preserved in request metadata and checked against the response legend. ([TS-S] §Structured level descriptions / §Reading a Score; [TS-A] §Structured Score levels)

### C.11 Failure, coverage, and emit

The outer underscore runs only if the acquisition fails. It is **not** a `try` around the entire lexical scope. A declared application error produced in a `decide` arm remains that arm's expression result; a later effect must be handled by its own ordinary match or judgment. A programmer defect in generated deterministic code is not silently reclassified as service unavailability.

C.7 has four successful source decision arms and one failure arm. The same fixture-to-arm mapping as A/B applies. In the nested specimen after C.7, the urgency arms are separate source obligations reached only under high billing. Re-reading an observation creates no new exchange obligation and grants no coverage to an untaken branch. ([a03]; [JEV1] §§4.1–4.2)

The TS lowering awaits one shared fetch request, then introduces immutable local bindings for the validated answers. The pure `decide` expressions become finite switches over the already-computed bands. No request is inserted into a nested `decide`, and the catch boundary never wraps the continuation. For C.7, a small shape sketch is:

```ts
const receipt = await $jevFetch(message, $triageQuestions);
switch (receipt.kind) {
  case "jev.unavailable": return { kind: "Ok", destination: "failed" };
  case "jev.invalid_response": return { kind: "Ok", destination: "failed" };
  case "answers": {
    const topic = receipt.answers.topic;
    const urgency = receipt.answers.urgency;
    const topicState = $topicDecision(topic, "0.8", "0.5");
    const urgencyState = $noulBand(urgency, "0.2", "0.8");
    switch (topicState.kind) {
      case "high.billing": return { kind: "Ok", destination: "billing" };
      case "high.other": return { kind: "Ok", destination: "general" };
      case "medium": return { kind: "Ok", destination: "confirm" };
      case "low": return { kind: "Ok", destination: "review" };
    }
  }
}
```

`urgencyState` is intentionally unused by this source decision; whole-response decoding has still checked it. An implementation may omit a dead, total band computation, but never skip that answer's validation or question acquisition. `$triageQuestions` is the literal schema from C.7, not an opaque runtime service or fixture file. All helpers are generated and dependency-free; §6 specifies the actual HTTP boundary.

**Pros:** one visible failure boundary per request; explicit distinction between asking and deciding; natural reuse of batch observations; no hidden requests in follow-on decisions. **Cons:** deeper indentation and an extra keyword; local rubrics still duplicate; observation scope needs a new checker rule; a single binary question is less compact than B; users must understand that outer failure does not catch continuation failures.

**Semantics preserved:** complete observations, one exchange per request, deterministic partitions, whole-batch validation, coverage, and proof separation. **Amendments requested:** common `JEV-S1` rules plus `JEV-S1/C` for acquisition scopes and pure decisions. Reuse the shared `CAN3620` site-identity rule; add `CAN3640` for a `decide` referencing an observation outside its declaring scope's successful continuation. No ambient exceptions, state mutation, sampling, or implicit batching are introduced.


## 6. Shared fetch emit: inherent transport, not hidden semantics

All three alternatives target **the same request and response contract**. A reached request sends a JSON body containing `state`, `model`, and the explicit question map to `POST https://api.typesafe.ai/v1/systemone`, authenticated with the runtime Bearer key. The emitter translates `options` to Choice `criteria`, preserves ordered Score `criteria`, and emits a Noul `criteria` object only when the source declares one. Description trees remain JSON trees, including metadata nulls. ([TS-API] §Evaluation endpoint / §Request body / §Question types; [TS-A] §Where structure is allowed)

The generated runtime has **no imports or runtime package dependencies**. It uses platform fetch, AbortController, timers, and the server host's `process.env.TYPESAFE_API_KEY`. This is a server-host contract, not permission to embed a key in browser output. No SDK retries, global SDK defaults, dynamic config loader, token logging, or fallback provider exists. The owner pays for production requests; compilation spends no model budget.

The following is a boundary sketch. `WireQuestions`, `JevState`, `ValidatedReceipt`, `$encodeRequest`, and `$decodeResponse` denote **compiler-generated types/functions specialized to the declared questions**, whose specified behavior appears below. They are not packages or fabricated TypeSafe API functions. The decoder specification is complete at the policy level; its generated implementation is outside this syntax task and has not been tested here.

```ts
declare const process:
  | { env: Readonly<Record<string, string | undefined>> }
  | undefined;

const $jevProfile = Object.freeze({
  model: "jev-1.13.0",
  timeout_ms: 2000,
  retries: 0,
  adapter: "typesafe-systemone-fetch/1"
});

type JevFailure =
  | { kind: "jev.unavailable"; reason: string }
  | { kind: "jev.invalid_response"; reason: string };

declare function $encodeRequest(
  state: JevState,
  model: string,
  questions: WireQuestions
): string;

declare function $decodeResponse(
  text: string,
  model: string,
  questions: WireQuestions
): ValidatedReceipt | JevFailure;

async function $jevFetch(
  state: JevState,
  questions: WireQuestions
): Promise<ValidatedReceipt | JevFailure> {
  const key = typeof process === "undefined"
    ? undefined
    : process.env.TYPESAFE_API_KEY;
  if (!key) {
    return { kind: "jev.unavailable", reason: "credentials" };
  }

  const body = $encodeRequest(state, $jevProfile.model, questions);
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), $jevProfile.timeout_ms);
  let text: string;
  try {
    const response = await fetch("https://api.typesafe.ai/v1/systemone", {
      method: "POST",
      redirect: "error",
      headers: {
        "Authorization": `Bearer ${key}`,
        "Content-Type": "application/json"
      },
      body,
      signal: controller.signal
    });
    if (!response.ok) {
      return { kind: "jev.unavailable", reason: "http_status" };
    }
    text = await response.text();
  } catch {
    return {
      kind: "jev.unavailable",
      reason: controller.signal.aborted ? "timeout" : "network"
    };
  } finally {
    clearTimeout(timer);
    controller.abort();
  }
  return $decodeResponse(text, $jevProfile.model, questions);
}
```

**The catch is deliberately narrow.** Serialization and deterministic validation are not inside a catch that can turn arbitrary compiler/runtime defects into an ordinary “failed” result. `$decodeResponse` returns `jev.invalid_response` for specified parse/schema violations; application continuation runs only after this boundary has returned. A malformed response is a normal boundary outcome, but an implementation defect is not permission to fabricate one.

`$decodeResponse` must reject duplicate response keys before ordinary object materialization loses that evidence; a plain `response.json()` cannot by itself establish the “exactly once” schema rule. It checks the model ID, complete question set, answer kinds, complete support, selected maxima, numeric ABI, and ordered Score legend—including structured legend values. It produces the same canonical observations that the compiler evaluates in `given`. Extra usage/envelope metadata is not fabricated into business evidence. The exact-decimal helpers and strict JSON decoder are generated implementation code with their own fixtures, not new dependencies. ([JEV1] §4.3; [TS-API] §Response body; [TS-S] §Structured level descriptions)

Every invocation makes **at most one fetch attempt**. Neither a timeout nor a disconnected response proves that the remote service did not already evaluate the request. Source-level retries remain explicit. The configured timer bounds the intended network wait; it is not a mathematical proof that the remote server terminates or that a blocked host event loop meets a wall-clock deadline.

### Build artifacts and coverage attribution

Production retains prompts, schemas, bounds, effective profile constants, deterministic dispatch, and declared domain behavior. It strips tests, observations, test names, fixture-only strings, evidence-bearing source maps, and packaged evidence reports. Source schemas are necessary runtime code; test observations are not. No cassette reader ships. ([R7]; [JEV1] §§4.1, 8)

Keep `errors.json`'s current array and fields: `kind`, `fields`, `raised_by`, `handled_by`, `hit_by_tests`. Record the intrinsic adapter as the producer of its two failure kinds. Attribute the source underscore as handling each admitted kind, while preserving **one source coverage-arm identity**. `hit_by_tests` remains script/expectation references, not measured execution or escaping-error proof. Programs without judgments keep their existing catalogs and emit byte-identical. Build-side catalogs containing evidence are not production assets. ([CATALOG] §catalogEntry / §buildCatalog; [JEV1] §6)

Async lowering follows the explicitly declared `Jev.request` capability through callers. Non-judgment functions do not become async merely because a project has a config file. The compiler remains responsible for behavioral checking; strict TypeScript checking validates output shape and must not be presented as model or coverage verification. ([R6] / R11; [a14] §Reinstated)

## 7. Amendment ledger and diagnostic allocation

**`JEV-S1` is the proposed amendment family for this surface task.** It is subordinate to the settled JEV-1.0 evidence rules. Reserve only the chosen option's optional rules at ratification; do not ship all three grammars. The reviewed `compiler/code.go` has no allocated CAN36xx entries. The CAN35xx numbers from the attached settlement remain reserved for their original rules; recheck both blocks against the actual integration revision before allocating. ([CODES]; [JEV1] §4.6)

### Common rule deltas

| Tag | Ratification-ready delta | Diagnostic / gate |
|---|---|---|
| **JEV-S1/TRANSPORT** — R2/R11 | A judgment-owning module pins one config provider. Runtime transport policy is the compiled value of that provider; source questions/sites cannot override it. The only environment input is the runtime credential. | `CAN3601`: config key-set mismatch (missing, duplicate, or unknown field); dedicated missing-model/timeout rules take precedence for their fields. `CAN3602`: transport-policy field appears in a question or site. `CAN3603`: adapter ABI is not supported by this compiler. Required provider/pin faults use CAN2101–2103; model, retry, timeout violations keep CAN3510/CAN3509/CAN3512. |
| **JEV-S1/HTTP** — R11 | Replace the attached settlement's SDK emit with platform fetch and a generated boundary implementing §6. Preserve its logical observation ABI and declared intrinsic failures. | Wire/body/auth/error/no-retry and Go/TS parity fixtures; no fabricated diagnostic for a runtime outage. CAN3508 still names numeric representability, not an SDK dependency. |
| **JEV-S1/FAILURE** — R5/R6 | Adopt §2.3's failure-only underscore over exactly the pinned two-kind failure set. Success exhaustiveness is independent. The clause does not catch failures of its continuation. | `CAN3604`: judgment failure clause is duplicated, nonterminal, or combined with individual intrinsic failure arms—the single failure-clause shape rule. Missing success/failure coverage remains CAN4101; wrong ordinary wildcard positions retain existing match diagnostics. |
| **JEV-S1/COVERAGE** — R7/R10 | Elaboration preserves source-arm identity. One failure clause or explicit alternative arm is one obligation; no generated expansion fabricates hits or creates additional author arms. | Existing CAN4107; exact source-hit golden tests, including a failure-only witness and an omitted low path. |
| **JEV-S1/OBS-LITERAL** — R7/R8 | Admit only the complete `observed` literals in §2.4 at judgment exchange outcomes. They elaborate to the same closed observations; no numbers, bands, or fields are inferred. | `CAN3608`: observation literal outside a judgment exchange outcome. Completeness and kind retain CAN3501; closed-value requirement retains CAN3513; numeric validity retains CAN3502–3508. |
| **JEV-S1/STRUCTURE** — R1/R2/R6 | Admit explicit closed description trees and text-only structured state. `unspecified` exists only in metadata and maps to foreign null; it is not a can value. `Jev__State` is an intrinsic type fixed by the adapter ABI, not an unresolved user type. | `CAN3605`: description is not a closed admitted metadata tree. `CAN3606`: instruction root is null or absent. `CAN3607`: state construction contains a leaf outside the admitted unbranded-text domain. |
| **JEV-S1/OBS-PROJECTION** — type rules | Permit only declared Choice support-key projections and in-range literal Score probability indexes as total observation projections. No dynamic map lookup is added. | `CAN3609`: observation probability projection does not name a statically admitted support member. |
| **JEV-S1/IDENTITY** — R4/a77 | Hash config, request metadata, and request-bearing function authority closure as in §2.2. A profile repin cannot silently preserve an unchanged effective request revision. Version the accepted-baseline format. | Existing CAN6013 / CAN2103; changed inputs/evidence retain CAN3511 where pinned. No new code that duplicates revision drift. |
| **JEV-S1/SURFACE** — Goal/R6 | Admit explicit state alternatives with identical observation bindings. Keep Noul uncertainty and Choice/Score low states in the successful domain. A single-question request uses fixed wire key `answer`; a batch has two or more explicitly named questions. | Existing missing/useless arm and type rules; `CAN3610` for a batch declaring fewer than two questions. The chosen grammar's singleton form is canonical. |

Metadata admits strings, booleans, exact integers within the JSON safe-integer range, representable `dec` literals, nested description trees, and `unspecified`. Their wire encoding is explicit; unsupported scalar literals are CAN3605. Runtime structured state remains deliberately narrower. Empty-string instruction policy is not inferred: an empty string is explicit text and is permitted, although it is a poor rubric; the compiler is not a prompt-quality judge. A missing Noul criterion block is the documented optional case; when present, both `true` and `false` entries are required, and either description may explicitly be `unspecified`. Choice option keys are unique; Score has 2–10 ordered entries and Choice has 2–255 options under the retained JEV-1.0 profile. Schema-shape faults keep CAN3501. ([TS-A] §Where structure is allowed; [TS-N] §Writing a Noul question; [JEV1] §5)

### Option-specific deltas

| Option tag | New rule | Code |
|---|---|---|
| **JEV-S1/A** | A function's explicit local judgment-site IDs are unique. | CAN3620. Shared with C when relevant. |
| **JEV-S1/A** | An inline primitive is a whole dispatch expression, not a nested value-producing call or bare arithmetic operand. | CAN3621. |
| **JEV-S1/B** | A `match noul`/`choice`/`score` head agrees with the referenced primitive declaration. | CAN3630. |
| **JEV-S1/B** | Every batch member's sole `state` parameter type agrees with its batch signature. | CAN3631. |
| **JEV-S1/C** | `decide` references only observations available in the declaring judgment's successful continuation. | CAN3640. |

Every new code needs a unique registry entry, ordered `allCodes`, an `explain` entry, an exact offending-token span, structured Expected/Found/Hint data, and shared CLI/LSP behavior. Existing code meanings do not change merely to make a new syntax easier. General syntax errors still use the existing parser diagnostic; ordinary effect-authority violations remain CAN3107/CAN3108. A judgment under an unsupported universal-contract proof remains unsupported/unverified, never proved from fixtures. ([CODES]; [R10]; [a81] §Obligations / Findings)

### Concrete negative controls

With B.5's `triage__queue` declaration, this is **not** exhaustive:

```can
match choice triage__queue(message)
  on high "billing" _ => Ok("billing")
  on high "other" _ => Ok("general")
  on medium _ => Ok("confirm")
  _ => Ok("failed")
```

`low` is missing: CAN4101, despite the underscore. The same mutation of A or C must fail at its missing-success boundary.

This remains an invalid successful fixture, not a convenient way to execute `_`:

```can
low => [exchange args (state = "charged twice") outcome observed choice(choice = "billing", probabilities = ("billing" = d"0.6"), confidence = d"0.49")]
```

The declared `other` probability is absent: CAN3501 before row execution. To test service failure, supply `jev.invalid_response("protocol")` explicitly. Likewise, a high billing arm containing a branch `confidence < d"0.5"` under high 0.8 cannot gain a witness by scripting that low confidence: the policy routes it away first, leaving the inner arm uncovered under CAN4107. ([JEV1] §§4.1, 8.2)

## 8. Recommendation and a constrained best-of design

**Recommend B as the foundation.** Its primitive name is immediately visible, use sites are compact, and versioned question signatures restore the declaration-first review surface instead of moving that information into business-control nesting. It removes actual transport repetition without hiding routing policy or asking an agent to infer which primitive a function-shaped wrapper contains. A is stronger for one-off local readability; C is stronger for distinguishing batch acquisition from reuse. Neither outweighs B's independently reviewable request contracts for this language.

| Criterion | A | B | C |
|---|---|---|---|
| Short single-question sites | Small, but repeat metadata | **Strongest: metadata declared once** | Extra `decide` nesting |
| Read formal declarations first | Request lives in body | **Strongest** | Request lives in body |
| Reuse and revision ownership | Function-owned sites | **Independent question contracts** | Function-owned scopes |
| One failure path without catch-all | Bounded terminal clause | Bounded terminal clause | Explicit acquisition boundary |
| Reusing several received answers | Bound answer records and nested matches | Bound answer records and nested matches | **Most direct scope model** |
| Evidence cost | Same complete observation bytes | Same complete observation bytes | Same complete observation bytes |
| Main implementation burden | Inline declaration analysis | New declaration and scrutinee kinds | Observation scope and pure `decide` |

### Best-of: B singletons, B declarations plus C scope for batches

The combination worth pursuing is **not all three syntaxes together**. Keep B's three primitive declaration kinds and singleton `match` forms. Keep B's named `judgments` schema composition for two or more questions. For a batch, use C's one-acquisition scope and pure `decide` layout instead of also admitting `match judgments`. Thus request cardinality selects one canonical form: singleton primitive match versus explicitly multi-question scope. No anonymous A-style inline alternative and no one-question batch is admitted.

The complete batch example below combines B's declarations with C's success/failure boundary. Its evidence is unchanged—not abbreviated, imported, or preclassified. This is a proposed hybrid, not a fourth semantic model.


```can
mod triage
  provides [triage__topic, triage__urgency, triage__overview, triage__route, Triage__Route]
  uses [Jev__Runtime@1]
  emits []

type Triage__Route rev 1 (
  destination: str
)

choice triage__topic(state: str) rev 1
  instructions "Which queue should handle the message?"
  options
    "billing" = "Charges, invoices, and refunds."
    "other" = "Requests outside billing."
  confidence high d"0.8" medium d"0.5"

noul triage__urgency(state: str) rev 1
  instructions "Does the sender request help before today's end?"
  probability no d"0.2" yes d"0.8"

judgments triage__overview(state: str) rev 1
  topic: choice triage__topic
  urgency: noul triage__urgency

fn triage__route(message: str) -> Triage__Route rev 1
  emits []
  effects [Jev.request]
  tests
    billing("charged twice") => Ok("billing") pinned
    other("charged twice") => Ok("general") pinned
    high_boundary("charged twice") => Ok("billing") pinned
    medium_boundary("charged twice") => Ok("confirm") pinned
    low("charged twice") => Ok("review") pinned
    tie("charged twice") => Ok("review") pinned
    offline("charged twice") => Ok("failed") pinned
    malformed("charged twice") => Ok("failed") pinned
  judge triage__overview(message)
    given
      billing => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"1", "other" = d"0"), confidence = d"1"), urgency = observed noul(noul = d"0.9"))]
      other => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "other", probabilities = ("billing" = d"0", "other" = d"1"), confidence = d"1"), urgency = observed noul(noul = d"0.1"))]
      high_boundary => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.9", "other" = d"0.1"), confidence = d"0.8"), urgency = observed noul(noul = d"0.5"))]
      medium_boundary => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.8", "other" = d"0.2"), confidence = d"0.5"), urgency = observed noul(noul = d"0.5"))]
      low => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.6", "other" = d"0.4"), confidence = d"0.49"), urgency = observed noul(noul = d"0.5"))]
      tie => [exchange args (state = "charged twice") outcome observed answers(topic = observed choice(choice = "billing", probabilities = ("billing" = d"0.5", "other" = d"0.5"), confidence = d"1"), urgency = observed noul(noul = d"0.9"))]
      offline => [exchange args (state = "charged twice") outcome jev.unavailable("network")]
      malformed => [exchange args (state = "charged twice") outcome jev.invalid_response("protocol")]
    decide topic, urgency
      high "billing", _ => Ok("billing")
      high "other", _ => Ok("general")
      medium, _ => Ok("confirm")
      low, _ => Ok("review")
    _ => Ok("failed")
```


Under this hybrid, the named descriptor owns request metadata; repeated uses remain distinct source AST sites with distinct queues, just like existing `match call` sites. The baseline uses structural site identity rather than line numbers. The scope imports exactly the descriptor's declared answer names into its successful continuation. No arbitrary dynamic binding or question lookup occurs. Its rules are `JEV-S1/B` plus `JEV-S1/C` scope semantics and the common batch-minimum rule; A's inline-site grammar is not admitted.

**Do not compromise on the underscore distinction.** The attractive but unsound shortcut is to let it mean “anything not matched above.” That would make deleting a low-confidence or newly added label arm silently redirect success to failure. The bounded failure clause is what makes the short syntax compatible with the language's contract, not the character `_` by itself.

### What is retained from Probably's feel

Take short judgment heads, confidence-aware control flow, an explicit uncertainty path, and replay that makes zero model calls. Do not take optional fallthrough without an expression result, winning-probability gates presented as TypeSafe confidence, first-option tie resolution, or chaos sampling. Here replay means executing the same ordinary full-observation `given` rows; importing a recording is an explicit conversion to those rows, not another runtime/build mode. ([PROBABLY] §§3–4, §6; [JEV1] §§7–8)

## 9. Verification obligations before a surface is called usable

These are required future checks, **not reported green runs**. Keep the language's existing single behavioral gate; do not create a privileged Jev-only mode that ordinary builds or LSP checks omit.

| Existing suite/tool | Change required for the chosen option |
|---|---|
| `compiler/emit_golden_test.go` | Preserve non-judgment AuthLogin, RetryLoop, Counter, and std goldens. Add singleton Noul/Choice/Score and complete batch goldens, with dependency-free helpers and catalogs. |
| `TestGoldenJSONDiags`, diagnosis and LSP suites | New one-fault fixtures for each selected rule; precise spans; independent failure/success exhaustiveness. Preserve ordinary bool/call wildcard rejections. |
| Registry and explain tests | `TestCodesUnique`, `TestAllCodesSequenced`, structured explanation and registry coverage for each new rule; no number reuse. |
| Revision and accepted-baseline suites | Config same-rev mutation, stale pin, deliberate repin, unchanged request rev after repin, transitive authority closure, prompt rewording, comment-only control, incompatible lock format, candidate-history rejection. |
| Acceptance suites | CAN6017 unchanged. CAN3511 catches pinned observation/input/policy premise drift without turning it into proof or a blocking acceptance decision. |
| Contract verification suites | Fixture confidence cannot prove natural-language truth; unsupported judgment effects cannot disappear from proof traversal. |
| Catalog tests including `TestCatalogNoPhantomHits` | One source underscore, two handled intrinsic kinds; no probability-mass hits; reference-based `hit_by_tests` unchanged. |
| `tools/modcheck` | Resolve config-provider identity and selected declaration kinds, state-type references, effects, and per-site given keys. No filesystem-default model, no service calls, no retired `-` rows. |
| `tools/gramcheck` and editor grammar | Add only the chosen surface's scopes, metadata trees, failure clause, observation literals, and state alternatives. Keep indentation, raw-string and exact-decimal rules. |
| Strict TS gate | Check the real emitted fetch/helper code under the project's pinned checker and explicit server platform declarations. No SDK stub or SDK dependency is needed. canlc still owns behavioral verification. |

**New conformance corpus:** exact request JSON for all structured metadata/state locations; label-only null descriptions; structured Score legends; valid and malformed full responses; duplicate keys; numeric-looking, `constructor`, and `__proto__` field names; model mismatches; numeric endpoints; ties; Noul boundary/equality cases; equal-mean different-shape Score observations; entire-batch rejection for a malformed ignored answer. Fake fetch fixtures must capture method, body, headers, redirect policy, and attempt count without a live service.

**Coverage corpus:** assert actual source-arm IDs for each of the eight complete example rows. Removing the only failed-request witness produces CAN4107 on `_`; removing low success handling produces CAN4101 even with `_`; leaving a covered low handler but never scripting low produces CAN4107. A nested decision needs a reaching witness, not merely probability mass. Explicitly grouped states share a source obligation, but primitive classifier tests cover each state and all boundaries independently.

**Evidence lifecycle:** missing table, omitted/unreached, omitted/reached, empty-on-hit, repeated site, identical state with differing answers, wrong state, partial queue, zero-consumed nonempty queue, shared helper, unrelated root queue, unknown key, and retired dash. Full-observation literals must never become legal production assertions or solver assumptions.

**Hermeticity and stripping:** compile with credentials absent and network denied; vary runtime environment settings and obtain identical artifacts. Put unique fixture-only sentinels in rows and confirm their absence from production TS, bundles, maps, and packaged reports. A repeated `decide` in C/hybrid must issue zero additional fetch calls; an explicit second judgment must issue exactly its own attempt and consume its own exchange.

Required commands remain `go test -count=1 ./...`, `go run ./tools/modcheck`, `go run ./tools/gramcheck`, and the repository's pinned strict TypeScript gate. No claim of passing these commands is made here. Migration removes the rejected per-descriptor transport fields, old SDK import/manifest, and duplicated per-error source arms; it does not preserve obsolete wrapper syntax as a second spelling. Legacy label-only extern stubs require explicit reviewed full observations, never invented complements or confidence. ([R10]; [JEV1] §§3.8, 9)

## 10. Open questions ranked by blast radius

These are remaining ratification or implementation choices, not gaps in the individual examples' routing semantics.

| Rank | Question / needed evidence | Binding position in these designs |
|---|---|---|
| **1 — Revision blast radius** | Does the owner accept that a global model/config change can require revisions across the reachable `Jev.request` contract closure? A representative multi-module application is needed to measure the review burden. | Yes in these proposals: explicit transitive identity protects pinned consumers. Do not replace it with runtime overrides or silently weaken prompt/model hashing. A narrower closure requires a separate concrete counterexample-based amendment. |
| **2 — Failure abstraction** | Is losing per-site transport-error discrimination acceptable for the first-class surface? A real application that needs different behavior for protocol failure versus outage would test this lean. | One terminal failure clause handles the pinned two-kind set. Keep both kinds in compiler/runtime diagnostics and catalogs. Do not add an open catch-all or quietly restore per-kind site boilerplate. Any future detailed handler needs one separately ratified canonical form. |
| **3 — Canonical syntax selection** | Choose A, B, C, or the constrained B/C hybrid after reviewing complete source specimens, not only hero snippets. | Recommend B foundation with C batch layout. Never ship all alternatives as interchangeable syntax. |
| **4 — State/disclosure scope** | Should dynamic numeric/Boolean state leaves, approved binary data, or branded values enter provider state? Missing: explicit numeric representation and disclosure-authority rules. | Only unbranded text trees here. Broader structured metadata is closed, not implicit serialization of arbitrary program values. |
| **5 — Wire/host conformance** | Which minimum server runtime and generated strict parser meet the platform-fetch, timeout, duplicate-key, and structured-legend requirements? Missing: committed host/decoder conformance vectors. | No SDK and no client credential exposure. Failing a boundary fixture blocks shipping, not the syntax settlement; actual API credentials are unnecessary for those tests. |
| **6 — Batching and checker scale** | What maximum question product is needed by real applications, and what deterministic resource budget can the exhaustive checker guarantee? | Exact checking or an explicit diagnostic/refusal; no sampling or acceptance after an incomplete analysis. No automatic request batching. |
| **7 — Operational quality** | Who owns independent labels, calibration evaluation, drift review, and recording retention? | Outside compiler acceptance. A full-observation row tests conditional code behavior; it never blesses the model, confidence, or chosen risk threshold. |

## 11. Ten-line summary

1. Three alternatives are presented: inline primitive expressions, reusable question signatures, and acquisition/decision scopes.
2. Recommend question signatures as the foundation, with a scope layout only for explicitly multi-question batches.
3. Put model, timeout, retries, and adapter identity in one versioned `jev.config.json`; use one visible provider pin.
4. Compile that config into runtime constants; read only `TYPESAFE_API_KEY` from the runtime environment.
5. Make `_ =>` a closed request-failure clause, never a catch-all for low confidence, uncertainty, or missing labels.
6. Keep Noul probability-only, Choice ties low, and Score a weighted mean with its complete distribution available.
7. Use complete `observed` literals inside existing request-matching exchanges; no bands, complements, or confidence are invented.
8. Cover actual source arms, including failure and uncertainty; one batch consumes one exchange even when an answer is ignored.
9. Preserve prompt/config revision identity, offline builds, proof separation, production stripping, and dependency-free fetch emit.
10. No compiler changes or gates were run; the missing repository a94 was replaced only as a reading source by the attached JEV-1.0 document.

## Source receipts

Repository references below use the retrieved commit, not a floating main link. JEV-1.0 refers to the earlier document supplied in this conversation, not to the unavailable a94 path. The `.md` and HTML TypeSafe pages are alternate representations of the requested documentation; successful reads were of the HTML pages.

[JEV1]: jev-judgment-evidence-decision.md "JEV-1.0 — attached conversation document, §§1–10"
[R]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/REQUIREMENTS.md "Requirements — Goal, R1–R11"
[R7]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/REQUIREMENTS.md "Requirements — R7"
[R8]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/REQUIREMENTS.md "Requirements — R8"
[R10]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/REQUIREMENTS.md "Requirements — R10"
[R6]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/REQUIREMENTS.md "Requirements — R6"
[IDIOMS]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/can-idioms.md "can-idioms — Tier 1 and Tier 2"
[AUTH]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/sketches/auth-login/auth.can "auth-login normative source"
[CODES]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/compiler/code.go "Diagnostic registry"
[CATALOG]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/compiler/catalog.go "Error catalog schema and reference semantics"
[a03]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a03-branch-coverage.md "Branch coverage — Hit semantics and retrospective amendment"
[a77]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a77-revision-identity.md "Revision identity — Fingerprint model, Canonicalization, Baseline workflow"
[a81]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a81-proof-obligations.md "Proof obligations — Obligations and Findings"
[a87]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a87-acceptance-authority.md "Acceptance authority — Semantics and Interplay"
[a91]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a91-given-omission-amendment.md "Partial given tables and unknown keys"
[a14]: https://github.com/veighnsche/can-lang/blob/4d447bac858c50d998804aa7fa96f2b66b3b557b/docs/a14-tsc.md "TypeScript gate — Reinstated (a70)"
[TS-INDEX]: https://docs.typesafe.ai/llms.txt "TypeSafe documentation index"
[TS-P]: https://docs.typesafe.ai/primitives "Primitives — definitions, state references, batching and dependent questions"
[TS-C]: https://docs.typesafe.ai/primitives/choice "Choice — request/response and structured criteria"
[TS-N]: https://docs.typesafe.ai/primitives/noul "Noul — instructions, true/false criteria, probability-only response"
[TS-S]: https://docs.typesafe.ai/primitives/score "Score — ordered levels, mean, distribution, structured level descriptions"
[TS-A]: https://docs.typesafe.ai/primitives/advanced "Advanced primitives — structured instructions and criteria"
[TS-F]: https://docs.typesafe.ai/confidence "Confidence — distribution statistic, three paths, risk-scaled thresholds"
[TS-STATE]: https://docs.typesafe.ai/concepts/state "State — strings, objects, arrays and separation from questions"
[TS-R]: https://docs.typesafe.ai/patterns/confidence-routing "Confidence routing — Step 2 action-specific thresholds"
[TS-H]: https://docs.typesafe.ai/concepts/how-to-build-with-system-one "How to build with System One — workflow design"
[TS-API]: https://docs.typesafe.ai/api "HTTP API — endpoint, body, question and response types, errors"
[TS-M]: https://docs.typesafe.ai/models "Models — full IDs and moving aliases"
[TS-FAN]: https://docs.typesafe.ai/patterns/fan-out "Fan-out — shared-state multi-question request"
[PROBABLY]: https://probably-lang.southpolesteve.workers.dev/ "Probably — feels, match, chaos and replay"
