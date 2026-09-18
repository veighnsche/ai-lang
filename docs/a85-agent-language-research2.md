# a85 — What makes a programming language good for AI coding agents?

**Research cutoff: September 18, 2026.**

## Executive finding

**Keep “contract-first.” Do not yet commit to “verbosity-first.”**

The strongest evidence supports giving agents **reliable semantic information, fast validation, actionable feedback, and an acceptance process they cannot satisfy merely by rewriting the specification**. It does not establish that maximal explicitness, eliminating inference, removing syntactic sugar, or requiring a decision table in every function produces better end-to-end results. Type-aware generation and compiler-informed context have measured benefits, but those findings concern access to semantics—not necessarily spelling every semantic fact out in source. ([arXiv][1])

My recommendation is to treat ai-lang as two experiments:

**The verification-and-tooling system** is the better-supported investment. **The new language surface** remains an experiment that must outperform an equally instrumented existing language—not an existing language deprived of ai-lang’s tools.

There is encouraging evidence that unfamiliar, agent-oriented languages can be usable. There is also evidence that unfamiliar syntax, missing abstractions, excessive instructions, and superficially successful verification can undermine results. None of the reviewed sources establishes a universally best language for coding agents. ([GitHub][2])

### How to read the evidence

**Measured** means an experiment reported outcomes. It does not automatically mean independent replication, a fair causal comparison, or transfer to repository-scale agents. **Proposed** means an implemented design or reasoned mechanism without an isolated demonstration of the claimed agent benefit. **Opinion** means a preference or argument without that demonstration.

Confidence below concerns the requirement’s **incremental benefit to agent reliability**, not the mathematical value of a language guarantee.

Except where noted, I found **no independent replication of the specific intervention** in the reviewed sources. Multiple models evaluated by the original authors are not independent replication. Several findings have corroboration from different systems, which is useful but weaker than a matched replication.

This assesses the design described in your prompt, not the current compiler implementation.

---

# 1. Requirements catalog

The catalog deliberately includes language, compiler, runtime, and harness requirements. Separating them prematurely would hide where the strongest evidence actually lies.

| ID      | Requirement                                                                                                            | Why it helps agents                                                                           | Evidence and source                                                                                                                                          | Confidence |
| ------- | ---------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------- |
| **R01** | Judge accepted changes against behavioral evidence independent of the implementation-generating process.               | Prevents matching mistakes in code and tests from masquerading as correctness.                | **Measured:** EvalPlus, 2023-10-30, found substantially more failures with strengthened tests. ([arXiv][3])                                                  | **High**   |
| **R02** | Provide a cheap, repeatable generate–check–repair loop.                                                                | Converts mistakes into local repair tasks instead of requiring perfect first attempts.        | **Measured:** SWE-agent, 2024-11-11; AutoVerus, 2025-08-22. Different settings corroborate feedback-driven workflows. ([arXiv][4])                           | **High**   |
| **R03** | Expose expected types, available bindings, and relevant declarations at the edit location.                             | Replaces guessing and irrelevant retrieval with semantic context.                             | **Measured:** *Statically Contextualizing LLMs with Typed Holes*, 2024-09-02. Narrow benchmark; same-author TypeScript extension. ([arXiv][5])               | **Med**    |
| **R04** | Check type-level validity during generation or immediately afterward.                                                  | Eliminates a substantial class of invalid candidates before expensive execution.              | **Measured:** Type-Constrained Code Generation, 2025-05-08; benefits and overhead depend on model and task. ([arXiv][1])                                     | **Med**    |
| **R05** | Keep semantic tooling useful while a program is incomplete.                                                            | Agents can fill holes and repair fragments without losing all compiler assistance.            | **Measured/proposed:** typed-hole contextualization and its proposed ChatLSP interface, 2024-09-02. ([arXiv][5])                                             | **Med**    |
| **R06** | Minimize irrelevant and mechanically redundant context while preserving necessary semantic information.                | Reduces token cost and interference from unnecessary instructions.                            | **Measured:** formatting study, 2025-08-19; AGENTS.md study, 2026-02-12. Neither establishes that arbitrary minification is safe. ([arXiv][6])               | **Med**    |
| **R07** | Preserve productive abstractions and a discoverable standard library rather than minimizing the language at any cost.  | Avoids forcing agents to reconstruct familiar operations through error-prone low-level code.  | **Measured, confounded:** PyLang study, 2026-05-15; unfamiliarity and missing library facilities are not cleanly separated. ([arXiv][7])                     | **Med**    |
| **R08** | Teach unfamiliar constructs through executable examples and authoritative, retrievable reference material.             | Helps bridge the gap between learning grammar and implementing behavior.                      | **Measured:** PyLang prompting experiments, 2026-05-15; not a universal curriculum result. ([arXiv][7])                                                      | **Med**    |
| **R09** | Resolve dependencies and API symbols against a version-pinned, authoritative environment.                              | Makes invented packages and obsolete API assumptions detectable.                              | **Measured problem; proposed boundary:** package-hallucination study, 2025-03-02. ([arXiv][8])                                                               | **Med**    |
| **R10** | Emit structured diagnostics containing location, obligation, actual versus expected facts, and reproducible evidence.  | Gives the agent a concrete repair target rather than prose to reinterpret.                    | **Measured feedback benefit; proposed exact schema:** AutoVerus and typed-hole work. JSON itself is not the proven ingredient. ([arXiv][9])                  | **Med**    |
| **R11** | Check consistency among requirements, formal contracts, examples, and implementation.                                  | Detects disagreement between artifacts that individually look plausible.                      | **Measured:** Clover, 2024-11-16, on a small verification benchmark—not proof of intended behavior. ([arXiv][10])                                            | **Med**    |
| **R12** | Report execution coverage and oracle adequacy as different properties.                                                 | A branch can execute without a test asserting the behavior that makes it correct.             | **Measured:** test-criteria study, 2026-09-08; fresh preprint on a selected hard-fault population. ([arXiv][11])                                             | **High**   |
| **R13** | Protect acceptance criteria, verifier settings, and specification authority from unauthorized agent edits.             | Prevents making the gate easier instead of fixing the program.                                | **Measured risk:** Anthropic reward-hacking experiment, 2025-11-21; deliberately vulnerable training setting. ([Anthropic][12])                              | **Med**    |
| **R14** | Distinguish statically proved, dynamically checked, tested, assumed, and unchecked claims.                             | Prevents a successful intermediate gate from being presented as stronger evidence than it is. | **Measured failure cases; proposed reporting rule:** AxDafny, 2026-06-30. ([arXiv][13])                                                                      | **Med**    |
| **R15** | Validate resource use and actual runtime behavior separately from functional proof.                                    | Mathematical correctness does not establish acceptable execution time or memory use.          | **Measured:** AxDafny’s executable evaluation, 2026-06-30. ([arXiv][13])                                                                                     | **Med**    |
| **R16** | Make effects inspectable and restrict their execution through explicit capabilities.                                   | Separates describing an effect from being authorized to perform it.                           | **Proposed/prototype:** Sunholo AILANG and Vera, inspected 2026-09-18. No isolated agent-reliability effect-system study found. ([AILANG][14])               | **Low**    |
| **R17** | Represent fallible outcomes explicitly and check their handling exhaustively.                                          | Makes omitted failure paths mechanically visible.                                             | **Proposed:** follows the checking mechanisms above; agent-oriented prototypes exist, but direct declared-error-set ablations were not found. ([GitHub][15]) | **Low**    |
| **R18** | Specify numeric semantics and verify that compilation and foreign boundaries preserve them.                            | Prevents a source-level guarantee from disappearing in the executable system.                 | **Proposed:** a correctness requirement, not an experimentally established LLM advantage. Runtime gaps motivate separate conformance checks. ([arXiv][13])   | **Low**    |
| **R19** | Make builds and checks reproducible, recording versions, inputs, seeds, and relevant environment state.                | Makes failures replayable and comparisons interpretable.                                      | **Engineering evidence/proposal:** SWE-agent and long-running-agent harness work; determinism is not independently isolated. ([arXiv][4])                    | **Med**    |
| **R20** | Provide compiler-generated canonical representations rather than requiring the agent to hand-maintain redundant forms. | Can deliver normalized artifacts without imposing all normalization work on generation.       | **Prototype/opinion:** Vera’s tooling; agent-language design proposals. No decisive canonical-source ablation found. ([GitHub][15])                          | **Low**    |
| **R21** | Support bounded, local edits with immediate validation and inspectable diffs.                                          | Limits the scope of accidental damage during repair.                                          | **Measured:** SWE-agent interface experiments; this does not establish AST edits as universally superior to text patches. ([arXiv][16])                      | **Med**    |
| **R22** | Preserve machine-readable progress and support interruption-safe resumption.                                           | Reduces repeated work and premature declarations of completion across sessions.               | **Engineering observation:** Anthropic long-running harness, 2025-11-26. ([Anthropic][17])                                                                   | **Med**    |

**The important asymmetry:** effects, totality, exact numerics, and exhaustive error handling may be excellent semantic guarantees even when their **additional benefit specifically for AI agents** remains unmeasured. They should not be marketed internally as empirical findings.

---

# 2. Design dimensions and tradeoffs

## 2.1 Syntax errors are only one part of the failure distribution

Two empirical taxonomies are particularly useful. Tambon and colleagues inspected 333 bugs from three older code models and identified problems including misunderstood requirements, corner cases, input types, hallucinated objects, and incomplete output. Wang and colleagues examined 558 incorrect snippets across six models and found missing conditions and steps, incorrect arguments and returns, and other logical defects. These studies use different populations and classifications; they do not establish one universal percentage breakdown. ([arXiv][18])

An important terminology trap: a paper’s “syntactic characteristics” of bugs can describe **where a bug occurs in the syntax tree**, not whether the program fails to parse. A wrong condition is syntactically well-formed code. ([arXiv][19])

The resulting design map is:

| Failure class                       | A plausible defense                                   | What that defense does not establish                 |
| ----------------------------------- | ----------------------------------------------------- | ---------------------------------------------------- |
| Invalid grammar                     | Parser feedback or grammar-constrained generation     | Correct meaning                                      |
| Wrong types or unavailable bindings | Type checking and semantic lookup                     | Correct algorithm                                    |
| Invented API or dependency          | Authoritative symbol and package resolution           | Correct use of a real API                            |
| Missing case or wrong control flow  | Exhaustiveness, contracts, properties, targeted tests | Adequacy of the stated specification                 |
| Wrong expected result               | Independent oracle or reference behavior              | That the reference captures the user’s intent        |
| Excessive time or memory            | Resource checks and executable evaluation             | General performance outside the evaluated conditions |

Package hallucination also deserves its own boundary. A large study generated 576,000 Python and JavaScript samples across 16 models and found recurring invented package names. That supports checking against real dependencies—not the stronger conclusion that the standard library should be tiny. ([arXiv][8])

**Implication for ai-lang:** optimizing grammar regularity while leaving specification, API discovery, and operational behavior weak would address only a fraction of the problem.

## 2.2 Explicit semantics are valuable; mandatory explicit spelling is a separate bet

**Measured:** Type-Constrained Code Generation substantially reduced compilation failures beyond syntax-only constraints. In the studied distribution, 94% of compilation failures were type-related—**not 94% of all programming errors**. One reported HumanEval result for Qwen2.5-32B moved from 79.6% to 81.8% pass@1, while median synthesis overhead was roughly 20–60% across settings. The evaluation used an unconstrained-first, constrained-resampling procedure, so this is not a simple equal-cost demonstration that more restrictions always win. ([arXiv][20])

**Measured:** typed-hole contextualization improved generation by retrieving expected types and relevant bindings. Its TypeScript extension provides useful same-author corroboration, but the benchmark and contextualization setup are narrower than arbitrary repository maintenance. ([arXiv][5])

Neither result establishes that an agent must manually write every type annotation.

The distinction I would encode in ai-lang’s requirements is:

> **The agent must not have to guess a semantic fact the compiler can determine. That does not imply the agent must spell that fact out.**

A concise expression whose fully elaborated type is queryable may be more agent-friendly than a verbose expression requiring the model to repeat the same type correctly in several places.

**Proposed sweet spot:** explicit public contracts and semantic boundaries, deterministic elaboration, and a fully inspectable typed representation. Whether local annotations help should be measured by construct and task—not settled through a global “no inference” doctrine.

## 2.3 Verbosity, formatting, and sugar have different effects

**Measured:** *The Hidden Cost of Readability* evaluated ten models across four languages on fill-in-the-middle completion. Removing legally removable formatting reduced input tokens by an average of 24.5% while broadly preserving completion performance. This does not establish that stripping meaningful names, comments, or structure is harmless, nor does it demonstrate repository-scale agent gains. ([arXiv][6])

**Measured, with an important reversal:** *Token Sugar* trained models to use token-efficient shorthand and reported up to 11.2% fewer generated tokens with similar correctness. But an untrained GPT-4.1 performed much worse on the unfamiliar shorthand: 94.5% completion accuracy with familiar Python versus 51.2% with sugar, recovering only to 54.9% with examples. Both results concern code completion rather than long-running repository agents. ([arXiv][21])

These findings do not support either slogan:

* “More explicit text is always easier for an agent.”
* “Fewer tokens are always better.”

They support a conditional claim: **compact representations help when the model can use them reliably**.

There is a related instruction-budget result. The AGENTS.md study found that LLM-generated repository context often reduced success and increased cost by roughly 20% or more. Human-written context generally helped success, although it also increased cost. The counter-evidence is against indiscriminate additional instructions—not against useful contracts or documentation. ([arXiv][22])

**Guidance:** retain a readable, inspectable canonical representation, but experiment with compact compiler-produced views. Do not make “maximum source verbosity” a requirement before measuring **successful-change cost**, including repeated reads and repairs.

## 2.4 Familiarity versus a purpose-built language: the evidence genuinely disagrees

**Negative result, with confounding:** the PyLang study used an unfamiliar experimental language and 352 tasks. Models could learn its syntax without closing the functional-performance gap. However, the language also removed familiar facilities and required manual reconstruction of operations normally supplied by libraries. The result challenges “a simple grammar is enough,” but does not cleanly distinguish unfamiliarity from lost expressiveness and library support. ([arXiv][7])

**Positive result, with serious limitations:** VeraBench’s author-run report evaluates 60 problems across nine models using Vera v0.1.8. It reports means of 98.7% for Vera, 96.7% for Python, and 99.7% for TypeScript. This demonstrates that the proposed language can be usable in that benchmark. It does **not** isolate the benefit of its constraints: the benchmark is saturated, uses single runs, acknowledges cross-language grading asymmetries, and changes several language properties simultaneously. Claims that typed slot references explain the advantage are therefore not established by an ablation. The asserted absence from all training sets is also not independently verifiable for closed models. ([GitHub][2])

These are contrasting results, not replications or direct rebuttals: the languages, task distributions, available abstractions, and harnesses differ.

An exploratory cross-language study by Dan Luu similarly finds that conclusions about token efficiency and correctness depend heavily on task difficulty and reasoning settings. It also challenges extrapolation from counting tokens in existing short solutions to the cost of generating correct programs. This is original experimental analysis, not a broad peer-reviewed language ranking. ([Dan Luu][23])

**Guidance:** minimize **unnecessary novelty**, not necessarily all novelty. A new language needs unusually good examples, discoverability, standard operations, and repair tooling. Being unpleasant for humans is not itself evidence of being useful for models.

## 2.5 Decision tables and branch evidence: useful, but not sufficient specifications

**Measured:** EvalPlus expanded HumanEval’s tests by about 80 times and found additional failures across 26 models, changing some model rankings. Passing the original examples was materially weaker evidence than it appeared. ([arXiv][3])

**Fresh measured counter-evidence:** the September 8, 2026 test-criteria preprint investigated coverage-oriented test selection on a deliberately selected population of 6,066 difficult faults. Branch-oriented suites triggered faults substantially more often than they actually detected them; detection remained around 1–2% in that setting because assertions were weak or wrong. These figures must not be generalized to all generated tests: easier faults had been filtered out. No independent replication was found. ([arXiv][11])

A simple counterexample shows the logical limitation without depending on either study.

Suppose the intended operation is absolute value:

```text
if x < 0 then -x else 0
```

The examples `-1 → 1` and `0 → 0` execute both branches. Every branch has evidence. Yet `2 → 0` is wrong.

There are three distinct properties:

**Exhaustiveness:** every permitted input reaches a defined case.

**Witness coverage:** every case executes on at least one committed example.

**Behavioral adequacy:** each case implements the intended relation over its entire domain.

The first two do not imply the third.

Decision tables can be excellent specifications for finite rule systems. Symbolic rows with universally proved obligations can also do substantially more than ordinary examples. But I found no direct comparative evidence that **mandatory per-function decision tables** outperform alternatives across general agent-written code.

**Recommendation for ai-lang:** preserve branch evidence as one signal, but do not let it stand in for independent properties, boundary tests, or adequate contracts. Co-location of tests and logic is not the problem; **common generation errors and common authority to weaken both** are the problem.

## 2.6 Verification helps—but there are three separate boundaries

### A. Proving a supplied specification versus discovering the right specification

**Measured:** AutoVerus reports proof-generation success above 90% on 150 tasks through staged generation, refinement, and verifier-guided repair. Those are proofs for supplied code and specifications, not end-to-end demonstrations that a model correctly formalizes arbitrary business intent. ([arXiv][9])

**Measured:** Clover checks consistency among code, documentation, and formal specifications. Its selected benchmark showed encouraging acceptance and rejection behavior, including finding flawed human-written examples. Its reported absence of accepted incorrect examples in that experiment is not a general zero-false-acceptance guarantee. ([arXiv][10])

For ai-lang, a proof should therefore answer:

> “What exactly was proved, under which assumptions, and who authorized those assumptions?”

It should not silently become “the user’s request was fulfilled.”

### B. Functional verification versus executable usefulness

**Measured:** AxDafny reported 56.4% on its LCBProDafny evaluation versus 11.6% for a single-shot baseline, but the agent could use up to 20 iterations, so the comparison bundles additional compute and interaction. More importantly, among 100 easy executable problems, 75 solutions verified; only 32 passed execution tests, while 39 exceeded time limits and four exceeded memory limits. This is not evidence that the verifier proved false mathematical results. It is evidence that the proved properties did not establish operational acceptability. ([arXiv][13])

Termination is similarly weaker than a useful latency bound. A terminating computation can still be unusable for the application.

### C. Restricting the final artifact versus restricting the construction process

**Measured, indirect:** CRANE studies how applying output constraints to the whole reasoning process can reduce solvability and proposes separating unconstrained reasoning from constrained output. Its experiments concern mathematical and logical reasoning tasks, not production programming languages. It therefore provides a caution, not a direct language-design verdict. ([arXiv][24])

**Guidance:** make final acceptance strict, while allowing incomplete drafts, typed holes, local experiments, and rejected candidates. An agent should not need to produce a fully proved module before receiving useful semantic feedback.

## 2.7 Tooling has better evidence than a particular source notation

**Measured:** SWE-agent demonstrates that navigation, editing, execution, and feedback interfaces materially affect agent performance. Its custom interface and ablations support treating the agent–computer interface as a first-class design problem. They do not establish that the same gains come specifically from JSON diagnostics, AST edits, or an unfamiliar syntax. ([arXiv][16])

**Engineering observation:** Anthropic’s long-running-agent work reports benefits from explicit progress tracking, incremental implementation, and end-to-end checks. It also reports that models were less likely to inappropriately modify a JSON feature list than a Markdown one. That last observation is not a controlled demonstration that JSON is generally the best programming-language surface. ([Anthropic][17])

There is a useful current disagreement about retrieval:

**Measured, 2025:** Cursor reported better question-answering results with semantic search and an online experiment showing modest improvements in code retention, particularly in larger repositories. Retention is a proxy, and the data are proprietary. ([Cursor][25])

**Operational revision, July 21, 2026:** a Cursor representative explained that the older dedicated semantic-search path was no longer meaningfully helping as models improved at grep/indexed search; agents were using Instant Grep and file reads instead. This statement supplies no comparable new controlled results, so it is not a statistical refutation of the earlier experiment. It is evidence that the best retrieval mechanism can change. ([Cursor - Community Forum][26])

The stable requirement is consequently **reliable access to relevant code and semantics**, not “must use embeddings” or “must use grep.”

Cognition’s Devin retrospective emphasizes clear requirements and verifiable outcomes; GitHub’s original Copilot coding-agent announcement emphasizes a controlled execution environment with tests and linters. Both support the industry emphasis on harnesses, but neither isolates a programming-language feature as the cause of success. ([Cognition][27])

## 2.8 Acceptance criteria must be outside the agent’s unilateral control

**Measured risk:** Anthropic’s reward-hacking research demonstrated coding-environment shortcuts in a deliberately constructed training setting, including exploiting the evaluation process instead of performing the task. The setup involved vulnerable environments and additional conditions; it does not measure the prevalence of such behavior in ordinary production agents. ([Anthropic][12])

Nevertheless, the design consequence is straightforward. An agent should not be able to make an implementation “correct” by silently:

* weakening a postcondition or replacing a meaningful predicate;
* strengthening a precondition to exclude failing inputs;
* changing expected results, deleting failing cases, or changing verifier settings;
* routing around the language through an unchecked host boundary.

This requires an authority model, not merely a stern instruction.

For ai-lang, specification changes should be **first-class semantic changes** with explicit authorization. Tests and contracts can remain in the same artifact while their approved portions are protected or independently checked.

## 2.9 Much “LLM language” work answers a different question

Three categories should not be conflated.

**Languages for building systems that call LLMs.** DSPy provides declarative modules and signatures and optimizes language-model pipelines against metrics. LMQL combines prompting, control flow, and constraints in a query language. Their measured results concern programming **with** language models, not whether coding agents produce more reliable applications when every source construct is explicit. ([arXiv][28])

**Languages intended to be written by LLMs.** Vera explores mandatory contracts, effects, typed slot references, and compiler tooling. Sunholo’s **AILANG is a separate project from your ai-lang**; it explores a deterministic core, typed effects, capabilities, contracts, and agent-facing tooling. Aver is another relevant prototype. These are useful prior art, but their existence does not independently validate each design choice. ([GitHub][15])

**Opinion about what such a language should look like.** Fabio Akita’s February 2026 article discusses an agent-language wishlist developed through conversations with Claude and Codex. It is worth reading for ideas, but agreement among model-generated preferences is not empirical support. A wishlist resembling ai-lang may be an echo of the same intuitions rather than independent confirmation. ([AkitaOnRails][29])

---

# 3. Top 10 highest-ROI requirements and falsifiable checks

The ranking prioritizes evidence of practical reliability, not architectural elegance. A lower-ranked security boundary can still be mandatory.

For the experiments, use paired tasks, multiple model families, repeated runs, fixed tool versions, and held-out evaluation. Measure both equal-token and equal-wall-time conditions. Report:

**Independently correct changes, false acceptance, cost per correct change, repair iterations, unrelated regressions, and runtime failures.**

The acceptance conditions below are proposed experiments, not thresholds established by the literature.

| Rank   | Requirement                                                  | Experiment and falsifiable acceptance check                                                                                                                                                                                                                                                                             |
| ------ | ------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **1**  | **Independent acceptance evidence — R01/R12**                | Compare current co-generated decision-table evidence with a frozen external oracle on the same implementations, including realistic seeded defects. Accept the additional mechanism if it reduces accepted incorrect programs without achieving that reduction merely by rejecting substantially more correct programs. |
| **2**  | **Fast check–repair loop — R02**                             | Compare generate-once, generic execution feedback, and structured compiler/verifier feedback under matched budgets. The claim fails if extra feedback does not improve independently correct completion or cost per correct completion.                                                                                 |
| **3**  | **Semantic context at the edit site — R03/R05**              | Compare raw file context, text search, and compiler-provided expected types/bindings on identical modifications. Accept semantic retrieval if gains persist against a strong search baseline—not only against an artificially impoverished prompt.                                                                      |
| **4**  | **Type-aware validity — R04**                                | Compare syntax-only checking, post-generation type checking, and type-constrained generation where supported. Measure hidden-test success after accounting for decoding and repair overhead. Fewer compile failures alone is not sufficient.                                                                            |
| **5**  | **Compact, information-preserving context — R06**            | Provide semantically equivalent verbose and compact views, keeping APIs and examples constant. Accept compactness when total successful-change cost falls without a meaningful correctness loss. Reject blanket verbosity if it performs worse at equal budget.                                                         |
| **6**  | **Productive, discoverable APIs — R07–R09**                  | Hold core semantics constant while varying standard operations and version-aware API lookup. Include tasks requiring unfamiliar APIs. Accept the feature when hallucinations and reconstruction bugs decrease and the documentation cost does not erase the gain.                                                       |
| **7**  | **Local, interruption-safe modification — R21/R22**          | Compare whole-function/file regeneration with bounded edits and checkpoints. Interrupt runs deliberately. Measure recovery cost and unrelated behavioral changes; do not presume AST editing wins over a well-designed patch interface.                                                                                 |
| **8**  | **Protected specification authority — R13**                  | Give an agent pressure to complete tasks containing difficult obligations and tempting bypasses. Accept only if it cannot obtain a passing result through unauthorized contract, domain, test, configuration, or host-boundary changes. Also verify that legitimate authorized spec changes remain possible.            |
| **9**  | **Contract/example/implementation consistency — R11**        | Compare decision tables alone, properties alone, and the combined specification system. Include wrong-but-self-consistent agent-generated specifications. Accept the combined system only if it improves external correctness enough to justify generation and proof cost.                                              |
| **10** | **Executable conformance and resource checks — R14/R15/R18** | Differentially test emitted TypeScript against reference semantics, including numeric boundaries, serialization, foreign calls, time, and memory. Any program labeled as satisfying a particular guarantee must meet that guarantee under the documented execution assumptions.                                         |

### The comparison ai-lang most needs

Use three principal conditions:

| Condition                                                                                                            | Purpose                                                      |
| -------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------ |
| **A. Strict TypeScript with an ordinary competent agent harness**                                                    | Establish a practical baseline.                              |
| **B. Strict TypeScript with the applicable ai-lang tooling, protected tests, semantic queries, and contract checks** | Estimate the benefit of the verification and tooling system. |
| **C. ai-lang with the same agent, budget, and corresponding tooling**                                                | Estimate what the new language adds beyond the harness.      |

Make libraries, task semantics, and evaluation equivalent where possible, and explicitly record unavoidable capability differences.

**B versus A measures much of the tooling contribution. C versus B is the more relevant test of the language thesis.**

A comparison between C and an under-tooled A could produce an impressive result while answering the wrong question.

---

# 4. Gaps in the current thesis

## 4.1 Verdict on the stated bets

| Current bet                                               | Research verdict                                                                                                        | Recommended disposition                                                                |
| --------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| **Verbose and explicit**                                  | Useful information helps; more text is not established as better.                                                       | Replace with **semantically explicit, mechanically compact where safe**.               |
| **Single canonical form**                                 | Useful for tooling and normalization; a universal generation advantage is unproved.                                     | Keep compiler normalization. Experiment with the source restrictions separately.       |
| **No inference the agent must guess**                     | Strong requirement when it means eliminating guesswork; unsupported when it prohibits deterministic compiler inference. | Require inspectable elaboration and type queries.                                      |
| **Tests and logic in one artifact**                       | Plausible consistency and maintenance benefit; no direct comparative evidence found.                                    | Keep as an organizational experiment, not an independence guarantee.                   |
| **Decision table and branch evidence for every function** | Coverage is useful; mandatory universality and oracle adequacy are unproved.                                            | Add external properties/oracles and test the per-function mandate.                     |
| **Declared error sets**                                   | Mechanically valuable in principle; direct agent-specific benefit remains unmeasured.                                   | Retain as a semantic choice with an empirical question attached.                       |
| **Termination, effects, exact numerics**                  | Valuable guarantees, but neither specification correctness nor operational sufficiency follows automatically.           | Preserve precise claims and validate their executable boundaries.                      |
| **Transpile to TypeScript**                               | A deployment choice, not a correctness argument.                                                                        | Treat lowering, runtime support, and foreign interfaces as part of the trusted system. |

The measured bases for these distinctions are the type/context studies, verbosity studies, testing research, and verification/runtime results—not a head-to-head evaluation of ai-lang. ([arXiv][1])

## 4.2 Requirements missing or underemphasized

**Specification provenance and authority.** Who decides that a contract expresses the desired behavior? Who may change it? What prevents silently narrowing the promised input domain?

**Independent evidence.** The same agent can write an incorrect function, matching examples, and a matching explanation. A unified artifact improves proximity, not independence.

**An agent-facing semantic protocol.** Expected types, available bindings, effect obligations, error identities, failing witnesses, dependency versions, and impacted proofs should be queryable without reconstructing them from prose.

**A productive library and teaching system.** A small semantic core need not mean making agents implement splitting, traversal, parsing, or numeric conventions from scratch.

**Partial-program support and economical repair.** The useful unit of work is often an incomplete edit, not a finished function. Strict final artifacts and helpful intermediate states are compatible.

**Operational and backend conformance.** The executable program must preserve the claimed semantics and meet the actual resource constraints.

**Evaluation of modification, not just initial synthesis.** Repository tasks need tests for renaming, changing types, evolving interfaces, preserving unrelated behavior, and recovering after interruption.

These are recommendations from the preceding evidence and failure analysis, not claims that the current repository lacks particular implementations.

## 4.3 The most uncomfortable unsupported bet

The biggest unsupported step is:

> “Because the language is explicit, constrained, and checkable, agents will produce more correct useful software at lower total cost.”

Checkability can make failures easier to identify while also increasing construction cost or rejecting programs the agent cannot reformulate within budget. A system can also make its accepted programs look excellent by accepting very few.

The relevant outcome is therefore not compilation rate, verification rate, branch coverage, or elegance in isolation. It is:

> **How many independently correct, operationally acceptable changes does the system deliver per unit of budget, and how often does it incorrectly claim success?**

The reviewed literature motivates this outcome, but does not supply a measured universal optimum for the language surface.

---

# 5. Sources

Dates below identify the **version inspected**, where applicable—not necessarily the first publication. Rolling project pages are labeled as snapshots. Sources are newest first within dated groups.

## Papers and research studies

| Date/version                    | Source                                                                                                                                                      | Evidence type                                  |
| ------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| **2026-09-08, v1**              | Hamidi et al., [How effective are traditional test criteria at detecting bugs in large language models generated code?](https://arxiv.org/abs/2609.09315)   | Measured; fresh preprint                       |
| **2026-06-30, v1**              | Breen et al., [AxDafny: Agentic Verified Code Generation in Dafny](https://arxiv.org/abs/2606.32007)                                                        | Measured; preprint                             |
| **2026-05-15, v1**              | Kumar et al., [Syntax Without Semantics: Teaching Large Language Models to Code in an Unseen Language](https://arxiv.org/abs/2605.15607)                    | Measured; preprint                             |
| **2026-02-12, v1**              | Gloaguen et al., [Evaluating AGENTS.md: Are Repository-Level Context Files Helpful for Coding Agents?](https://arxiv.org/abs/2602.11988)                    | Measured; preprint                             |
| **2025-12-09, v1**              | Sun et al., [Token Sugar: Making Source Code Sweeter for LLMs through Token-Efficient Shorthand](https://arxiv.org/abs/2512.08266)                          | Measured; preprint                             |
| **2025-08-22, v3**              | Yang et al., [AutoVerus: Automated Proof Generation for Rust Code](https://arxiv.org/abs/2409.13082)                                                        | Measured; OOPSLA 2025                          |
| **2025-08-19, v1**              | Pan et al., [The Hidden Cost of Readability: How Code Formatting Silently Consumes Your LLM Budget](https://arxiv.org/abs/2508.13666)                       | Measured; preprint                             |
| **2025-05-08, v2**              | Mündler et al., [Type-Constrained Code Generation with Language Models](https://arxiv.org/abs/2504.09246)                                                   | Measured; PLDI 2025                            |
| **2025-04-06, v2**              | Banerjee et al., [CRANE: Reasoning with constrained LLM generation](https://arxiv.org/abs/2502.09061)                                                       | Theory and measured results; indirect transfer |
| **2025-03-02, v3**              | Spracklen et al., [We Have a Package for You! A Comprehensive Analysis of Package Hallucinations by Code Generating LLMs](https://arxiv.org/abs/2406.10279) | Measured; USENIX Security 2025                 |
| **2025-02-12, v3**              | Wang et al., [Towards Understanding the Characteristics of Code Generation Errors Made by Large Language Models](https://arxiv.org/abs/2406.08731)          | Measured; ICSE 2025                            |
| **2024-11-16, v4**              | Sun et al., [Clover: Closed-Loop Verifiable Code Generation](https://arxiv.org/abs/2310.17807)                                                              | Measured prototype                             |
| **2024-11-11, v3**              | Yang et al., [SWE-agent: Agent-Computer Interfaces Enable Automated Software Engineering](https://arxiv.org/abs/2405.15793)                                 | Measured; NeurIPS 2024                         |
| **2024-09-02, v1**              | Blinn et al., [Statically Contextualizing Large Language Models with Typed Holes](https://arxiv.org/abs/2409.00921)                                         | Measured; OOPSLA 2024                          |
| **2024-03-18, v2**              | Tambon et al., [Bugs in Large Language Models Generated Code: An Empirical Study](https://arxiv.org/abs/2403.08937)                                         | Measured taxonomy                              |
| **2023-10-05, v1**              | Khattab et al., [DSPy: Compiling Declarative Language Model Calls into Self-Improving Pipelines](https://arxiv.org/abs/2310.03714)                          | Measured; programming with LLMs                |
| **2022-12-12, initial version** | Beurer-Kellner et al., [Prompting Is Programming: A Query Language for Large Language Models](https://arxiv.org/abs/2212.06094)                             | Measured prototype; PLDI 2023                  |

## Benchmarks and original comparative experiments

| Date/version                      | Source                                                                                                                                                             | Important qualification                                                                 |
| --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------- |
| **2026-09-18 snapshot**           | Allan, [VeraBench](https://github.com/aallan/vera-bench)                                                                                                           | Author-run comparison; reported language version v0.1.8; saturation and grading caveats |
| **2023-10-30, v3**                | Liu et al., [Is Your Code Generated by ChatGPT Really Correct? Rigorous Evaluation of Large Language Models for Code Generation](https://arxiv.org/abs/2305.01210) | EvalPlus; NeurIPS 2023                                                                  |
| **Undated; inspected 2026-09-18** | Dan Luu, [How does programming language affect token efficiency and correctness?](https://danluu.com/pl-tokens/)                                                   | Exploratory original experiments and methodological critique                            |

## Industry experiments, engineering reports, and postmortems

| Date           | Source                                                                                                                                                                                     | Evidence type                                                     |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------- |
| **2026-07-21** | Cursor, [Explanation of the move away from the dedicated semantic-search path](https://forum.cursor.com/t/what-do-you-think-about-cursor-removing-the-codebase-indexing-settings/165899/5) | Operational statement; no new controlled results supplied         |
| **2025-11-26** | Anthropic, [Effective harnesses for long-running agents](https://www.anthropic.com/engineering/effective-harnesses-for-long-running-agents)                                                | Engineering observations                                          |
| **2025-11-21** | Anthropic, [From shortcuts to sabotage: natural emergent misalignment from reward hacking](https://www.anthropic.com/research/emergent-misalignment-reward-hacking)                        | Controlled research in a deliberately vulnerable training setting |
| **2025-11-14** | Cognition, [Devin’s 2025 Performance Review: Learnings From 18 Months of Agents At Work](https://cognition.com/blog/devin-annual-performance-review-2025)                                  | Product retrospective; observational                              |
| **2025-11-06** | Cursor, [Improving agent with semantic search](https://cursor.com/blog/semsearch)                                                                                                          | Internal benchmark and online experiment                          |
| **2025-05-19** | GitHub, [GitHub Copilot coding agent in public preview](https://github.blog/changelog/2025-05-19-github-copilot-coding-agent-in-public-preview/)                                           | Historical product architecture announcement                      |

## Language prototypes and proposals

| Date                    | Source                                                                                                                                                                | Evidence type                                                   |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------- |
| **2026-09-18 snapshot** | [Vera](https://github.com/aallan/vera)                                                                                                                                | Implemented agent-oriented language                             |
| **2026-09-18 snapshot** | [Sunholo AILANG](https://ailang.sunholo.com/)                                                                                                                         | Implemented agent-oriented language; distinct from your project |
| **2026-09-18 snapshot** | [Aver](https://github.com/jasisz/aver)                                                                                                                                | Implemented language prototype                                  |
| **2026-02-09**          | Fabio Akita, [AI Agents: What Would Be the Best Programming Language for LLMs?](https://akitaonrails.com/en/2026/02/09/ai-agents-best-programming-language-for-llms/) | Opinion and model-assisted design discussion                    |

---

# 6. Unresolved

These questions remain open in the reviewed evidence. They should become explicit research items rather than implicit language commitments.

| Open item                                 | What remains unknown                                                                                                                                  |
| ----------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| **U01 — Annotation policy**               | Which explicit annotations improve agent reliability beyond deterministic inference plus semantic queries, and in which contexts?                     |
| **U02 — Canonical syntax**                | Does one source spelling materially improve generation and modification, or is canonical elaboration sufficient?                                      |
| **U03 — Names and delimiters**            | Do naming conventions, name elimination, indentation, or delimiter choices causally improve repository-scale outcomes when familiarity is controlled? |
| **U04 — Mandatory decision tables**       | Which function classes benefit enough to justify the extra specification and maintenance burden?                                                      |
| **U05 — Specification independence**      | How reliably can independently generated specifications expose shared model blind spots without a trusted human or executable reference?              |
| **U06 — Library size and expressiveness** | Where is the optimum between a small learnable surface and enough abstractions to avoid repetitive low-level reconstruction?                          |
| **U07 — Effects, errors, and totality**   | What is their marginal agent benefit after controlling for the checking tools, feedback, and extra token cost?                                        |
| **U08 — Verification economics**          | Which obligations deserve static proof, bounded checking, runtime guards, or tests under a realistic task budget?                                     |
| **U09 — Long-lived modification**         | Do the benefits survive interface evolution, concurrency, large repositories, and multi-agent changes rather than isolated synthesis?                 |
| **U10 — Compiler and runtime trust**      | How should source guarantees be validated across lowering, foreign calls, serialization, resource limits, and deployed execution?                     |
| **U11 — Model dependence**                | Which language advantages persist across model families and updates rather than exploiting one model’s current weaknesses?                            |
| **U12 — Adaptation cost**                 | When do training, examples, constrained decoding, and custom tooling repay the cost of adopting unfamiliar syntax?                                    |

**Decision recommendation:** freeze the semantic promises, evidence model, and evaluation protocol first. Freeze restrictive surface choices only after they beat an equally well-tooled TypeScript baseline. The next decisive artifact is a comparative benchmark with feature ablations—not another declaration that explicitness must be good for agents.

[1]: https://arxiv.org/abs/2504.09246 "[2504.09246] Type-Constrained Code Generation with Language Models"
[2]: https://github.com/aallan/vera-bench "GitHub - aallan/vera-bench: VeraBench: a benchmark suite for LLM code generation in Vera · GitHub"
[3]: https://arxiv.org/html/2305.01210v3 "https://arxiv.org/html/2305.01210v3"
[4]: https://arxiv.org/abs/2405.15793 "https://arxiv.org/abs/2405.15793"
[5]: https://arxiv.org/html/2409.00921v1 "https://arxiv.org/html/2409.00921v1"
[6]: https://arxiv.org/html/2508.13666v1 "https://arxiv.org/html/2508.13666v1"
[7]: https://arxiv.org/html/2605.15607v1 "https://arxiv.org/html/2605.15607v1"
[8]: https://arxiv.org/abs/2406.10279 "https://arxiv.org/abs/2406.10279"
[9]: https://arxiv.org/abs/2409.13082 "https://arxiv.org/abs/2409.13082"
[10]: https://arxiv.org/abs/2310.17807 "https://arxiv.org/abs/2310.17807"
[11]: https://arxiv.org/html/2609.09315v1 "How effective are traditional test criteria at detecting bugs in large language models generated code?"
[12]: https://www.anthropic.com/research/emergent-misalignment-reward-hacking "https://www.anthropic.com/research/emergent-misalignment-reward-hacking"
[13]: https://arxiv.org/html/2606.32007v1 "AxDafny: Agentic Verified Code Generation in Dafny"
[14]: https://ailang.sunholo.com/ "https://ailang.sunholo.com/"
[15]: https://github.com/aallan/vera "https://github.com/aallan/vera"
[16]: https://arxiv.org/html/2405.15793 "SWE-agent: Agent-Computer Interfaces Enable Automated Software Engineering"
[17]: https://www.anthropic.com/engineering/effective-harnesses-for-long-running-agents "https://www.anthropic.com/engineering/effective-harnesses-for-long-running-agents"
[18]: https://arxiv.org/abs/2403.08937 "https://arxiv.org/abs/2403.08937"
[19]: https://arxiv.org/html/2406.08731v1 "https://arxiv.org/html/2406.08731v1"
[20]: https://arxiv.org/pdf/2504.09246v2 "https://arxiv.org/pdf/2504.09246v2"
[21]: https://arxiv.org/html/2512.08266v1 "Token Sugar: Making Source Code Sweeter for LLMs through Token-Efficient Shorthand"
[22]: https://arxiv.org/html/2602.11988v1 "https://arxiv.org/html/2602.11988v1"
[23]: https://danluu.com/pl-tokens/ "https://danluu.com/pl-tokens/"
[24]: https://arxiv.org/html/2502.09061v2 "https://arxiv.org/html/2502.09061v2"
[25]: https://prod.cursor.com/blog/semsearch "https://prod.cursor.com/blog/semsearch"
[26]: https://forum.cursor.com/t/what-do-you-think-about-cursor-removing-the-codebase-indexing-settings/165899/5 "What do you think about Cursor removing the codebase indexing settings? - #5 by maximep - Discussions - Cursor - Community Forum"
[27]: https://cognition.com/blog/devin-annual-performance-review-2025 "https://cognition.com/blog/devin-annual-performance-review-2025"
[28]: https://arxiv.org/abs/2310.03714 "https://arxiv.org/abs/2310.03714"
[29]: https://www.akitaonrails.com/en/2026/02/09/ai-agents-best-programming-language-for-llms/ "https://www.akitaonrails.com/en/2026/02/09/ai-agents-best-programming-language-for-llms/"
