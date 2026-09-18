# a85 — What makes a language good for AI coding agents?

**Research cutoff: 18 September 2026.**

## Verdict

**Keep ai-lang’s emphasis on checkable semantics. Do not treat maximal verbosity, universal explicitness, or banning inference and sugar as established requirements.**

The strongest evidence supports something narrower:

> **Give the agent precise, relevant information about the program; make incorrect changes cheap to detect and repair; and prevent the agent from redefining what counts as correct.**

Type-aware generation, compiler-informed context, executable feedback, and well-designed editing interfaces have measured benefits. The evidence for particular surface-language doctrines is weaker, mixed, or actively contrary to the “more explicit tokens must be better” hypothesis. There are controlled results showing benefits from compact representations—and results showing that unfamiliar compact representations substantially hurt models that have not learned them. ([arXiv][1])

**My recommendation is to optimize ai-lang for independently validated task completion per unit of cost—not for compilation rate, annotation density, branch coverage, or proof completion individually.**

The literature reviewed here does **not** establish a universally best language. A June 2026 study produced 34 chess-engine artifacts across 17 primary languages and found substantial differences in cost, features, and performance. However, its authors explicitly identify confounding factors that prevent a causal language ranking. “Agents can use this language” and “this language is a good choice for agents” are different claims. ([arXiv][2])

### How to read the evidence

**Measured—intervention** means researchers changed a mechanism and compared outcomes. **Measured—observational** means they analyzed failures, repositories, or production work without isolating the language feature causally. **Proposed** means a design or prototype exists, but the relevant agent advantage is not established. **Opinion** means an assertion without supporting comparative evidence.

Confidence below concerns the **direction of the requirement**, not a universal effect size. Unless a replication or rebuttal is identified explicitly, I did not locate an independent exact replication; released artifacts and evaluations across several models are not themselves independent replication.

---

## 1. Requirements catalog

These are candidate requirements for the **language–compiler–agent interface**, rather than syntax alone.

| ID      | Requirement                                                                                               | Why it helps agents                                                                            | Evidence and qualification                                                                                                                                                                                                  | Confidence                                    |
| ------- | --------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------- |
| **R01** | Provide a fast external check–diagnose–repair loop.                                                       | Turns speculation about correctness into observable failures.                                  | **Measured—intervention:** SWE-agent, 2024; testing/static-analysis feedback study, 2024/2026. Benefits concern external feedback, not ungrounded self-critique. ([arXiv][3])                                               | **High**                                      |
| **R02** | Expose expected types, available bindings, and relevant definitions at the edit location.                 | Prevents guessing about project-specific interfaces.                                           | **Measured—intervention:** typed-hole contextualization, OOPSLA 2024; small application benchmark, not unrestricted repository maintenance. ([arXiv][4])                                                                    | **High**                                      |
| **R03** | Support generation constrained by syntax, scope, and types where practical.                               | Avoids spending attempts on mechanically invalid candidates.                                   | **Measured—intervention:** Synchromesh, ICLR 2022; type-constrained generation, PLDI 2025. General-purpose semantic correctness remains unguaranteed. ([arXiv][5])                                                          | **High**                                      |
| **R04** | Make diagnostics actionable and machine-addressable, with locations and semantic detail.                  | Enables targeted repairs rather than whole-file guessing.                                      | **Measured** for useful feedback; **proposed** for a particular JSON schema, stable-code design, or counterexample format. Anthropic reports format sensitivity rather than one universally best encoding. ([Anthropic][6]) | **High** for content; **medium** for encoding |
| **R05** | Provide bounded, task-relevant views instead of requiring entire files or instruction corpora.            | Preserves useful context without overwhelming the task.                                        | **Measured:** SWE-agent context ablations; AGENTS.md study, revised June 2026. More supplied information is not consistently better. ([arXiv][3])                                                                           | **High**                                      |
| **R06** | Support small edits with immediate validation and clear failure recovery.                                 | Limits accidental damage and shortens repair loops.                                            | **Measured—intervention:** SWE-agent’s lint-enabled editor. Atomic transactions and AST edits are reasonable extensions, not all independently validated. ([arXiv][3])                                                      | **High**                                      |
| **R07** | Make the actual versioned API and dependency universe queryable.                                          | Replaces remembered or invented APIs with available symbols.                                   | **Measured—observational:** package-hallucination study, USENIX Security 2025. Supports grounding; does not establish that the smallest possible standard library is best. ([arXiv][7])                                     | **High**                                      |
| **R08** | Separate implementation authority from authority over acceptance criteria.                                | Prevents correlated tests and weakened specifications from certifying a wrong implementation.  | **Measured—intervention:** code-before-test study, July 2026; **production observation:** GitHub’s September 2026 migration report. ([arXiv][8])                                                                            | **High**                                      |
| **R09** | Supplement examples with adversarial cases, properties, and independent tests.                            | Detects implementations that merely fit supplied examples.                                     | **Measured:** EvalPlus, NeurIPS 2023; PropertyGPT, NDSS 2025. Generated properties still require validation. ([NeurIPS Proceedings][9])                                                                                     | **High**                                      |
| **R10** | Expose coverage as navigational feedback, not as a correctness certificate.                               | Shows unexercised behavior while preserving the distinction between execution and correctness. | **Measured:** CoverUp, FSE 2025; transfer challenged by a 2026 comparative reimplementation. ([arXiv][10])                                                                                                                  | **High** for feedback                         |
| **R11** | Check specifications for vacuity, weakening, and unauthorized proof escape hatches.                       | A successful proof is only useful when its assumptions and claims are meaningful.              | **Measured—observational:** vericoding benchmark reports specifications admitting trivial solutions; consistency checking explored by Clover. ([arXiv][11])                                                                 | **High**                                      |
| **R12** | Make composition obligations visible across functions and modules.                                        | Local correctness does not automatically establish correct integration.                        | **Measured:** DafnyCOMP, September 2025, exposes severe compositional verification failures; synthetic benchmark limits generalization. ([arXiv][12])                                                                       | **High**                                      |
| **R13** | Represent error, effect, and capability boundaries in checkable interfaces.                               | Makes hidden obligations explicit and restricts unintended behavior.                           | **Proposed:** agent-oriented languages including Vera and Sunholo AILANG. Strong engineering rationale; isolated causal agent evidence is limited. ([Vera][13])                                                             | **Medium**                                    |
| **R14** | Validate resource use and deployed behavior separately from logical correctness.                          | Termination and functional proofs do not imply acceptable runtime or memory use.               | **Measured:** AxDafny, June 2026, reports verified solutions failing execution limits. ([arxiv.org][14])                                                                                                                    | **High**                                      |
| **R15** | Define and test source-to-target and foreign-interface semantic conformance.                              | Prevents a verified source program from changing meaning at execution boundaries.              | **Measured—observational:** GitHub migration failures involving numeric and other implicit semantics. Applying this to a TypeScript transpiler is an engineering inference. ([github.blog][15])                             | **High**                                      |
| **R16** | Allow useful analysis of incomplete programs while keeping final acceptance strict.                       | Agents can construct and repair incrementally without losing all semantic assistance.          | **Measured prototype:** Hazel typed holes; language-server context remains available around incomplete code. ([arXiv][4])                                                                                                   | **Medium–high**                               |
| **R17** | Provide deterministic canonical rendering without assuming only one authoring representation is useful.   | Stabilizes stored artifacts while allowing efficient generation and editing views.             | **Proposed design implication:** SimPy demonstrates AST-preserving alternate representations; it does not prove a universal canonical syntax. ([arXiv][16])                                                                 | **Medium**                                    |
| **R18** | Optimize token efficiency without discarding useful semantic information.                                 | Reduces context and generation costs when meaning remains recoverable.                         | **Measured—intervention:** formatting-removal study, 2025/ICSE 2026; benefits measured on code completion, not arbitrary agent trajectories. ([arXiv][17])                                                                  | **Medium–high**                               |
| **R19** | Require informative boundary types, but test rather than assume the value of redundant local annotations. | Retains interface information without forcing the agent to repeat compiler-known facts.        | **Measured, mixed:** MultiPL-E found benefits from precise TypeScript signatures but no significant Python annotation effect. ([arXiv][18])                                                                                 | **Medium**                                    |
| **R20** | Make the language cheaply learnable through a versioned grammar, examples, and meaningful names.          | Reduces interference from familiar-but-wrong language and API patterns.                        | **Measured:** Grammar Prompting, NeurIPS 2023; naming and representation sensitivity in MultiPL-E and ReCode. ([arXiv][19])                                                                                                 | **Medium–high**                               |
| **R21** | Make builds, checks, and effect traces reproducible and replayable.                                       | Lets an agent distinguish its code changes from environmental variation.                       | **Proposed/engineering:** replay-oriented agent-language designs. I did not locate an isolated experiment establishing an agent-specific benefit from deterministic builds alone. ([AILANG][20])                            | **Medium**                                    |
| **R22** | Evaluate maintainable changes and legitimate task completion, not only benchmark passes.                  | Prevents rewarding test gaming, excessive edits, and unmergeable solutions.                    | **Measured—observational:** METR, March 2026; **benchmark proposal:** Cognition’s FrontierCode, June 2026. ([Metr][21])                                                                                                     | **High**                                      |

**Not established as requirements:** maximal verbosity; a blanket inference ban; eliminating all syntactic sugar; mandatory decision tables for every function; a universally optimal delimiter style; or minimizing library size regardless of expressiveness.

That is not a finding that these choices are wrong. It is a finding that they need experiments rather than being treated as consequences of “AI-first.”

---

## 2. Design dimensions and tradeoffs

### 2.1 What actually goes wrong?

There is no defensible universal pie chart of “agent errors.” Studies use different models, languages, task distributions, and denominators.

A 2025 empirical study examined **333 bugs** and identified failures including misunderstood requirements, missing corner cases, incorrect input assumptions, hallucinated objects, wrong attributes, syntax mistakes, and incomplete generation. It is a taxonomy of sampled failures, not a prevalence estimate for current agents. ([Springer][22])

For ai-lang, the useful distinction is:

| Failure layer                 | Typical failure                                                             | What helps most directly                                                 |
| ----------------------------- | --------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| **Syntax and representation** | Invalid syntax, incomplete output, syntax borrowed from another language.   | Grammar support, cheap parsing, learnable examples.                      |
| **Binding and types**         | Unknown identifiers, wrong members, incompatible arguments.                 | Actual symbol/type context and compiler feedback.                        |
| **API availability**          | Invented package, nonexistent operation, wrong version.                     | Versioned dependency and API grounding.                                  |
| **Behavior and control flow** | Wrong condition, missing case, incorrect evaluation order, wrong recursion. | Independent behavioral tests, counterexamples, contracts.                |
| **Specification and proof**   | Weak postcondition, impossible precondition, proving the wrong property.    | Specification review, non-vacuity checks, protected acceptance criteria. |
| **Integration and runtime**   | Serialization mismatch, ambient-state assumptions, resource exhaustion.     | Boundary conformance, integration tests, execution budgets.              |

The categories are a synthesis of the empirical taxonomy, package-hallucination research, verification studies, and production reports—not a claim that one language mechanism fixes all six. ([Springer][22])

An especially relevant production observation: GitHub’s **16 September 2026** Copilot-to-Rust migration report categorized 8,678 captured compiler diagnostics. Name/import problems accounted for 37%, missing members 22%, type mismatches 14%, and trait bounds 11%; ownership/borrowing/lifetime diagnostics accounted for 1.7%. These are diagnostic events, not independent defects, but they challenge the assumption that sophisticated ownership checking must be the dominant agent bottleneck. This is one migration, not evidence that Rust is universally best. ([The GitHub Blog][15])

### 2.2 Precise types versus “spell everything out”

**Measured result:** the PLDI 2025 type-constrained generation paper reports eliminating more than half of compilation errors. In its analyzed compilation failures, **94% were type errors and 6% syntax errors**. That is emphatically **not** “94% of all generated-code bugs are type errors.” The method uses type reasoning and inference to constrain generation. ([arXiv][23])

This supports **making types computationally available**, not making the model write every inferable type.

MultiPL-E supplies a more direct annotation experiment. Weakening TypeScript function signatures to `Any` hurt performance, while removing optional Python type annotations had no significant effect. Its cross-language analysis did not establish a general typed-language advantage. These are older Codex results and require retesting with current agents. ([arXiv][18])

**Recommended sweet spot:** explicit interfaces, explicit semantic choices, deterministic compiler elaboration, and queryable inferred facts.

For example, these are different:

* The agent must guess whether an API accepts a decimal or a binary float.
* The compiler derives a local variable’s type from a fully typed expression.
* The compiler silently inserts a potentially lossy numeric conversion.

My recommendation is to eliminate the first and tightly control the third. The second is potentially useful automation.

The typed-hole work reinforces this distinction: retrieving relevant type definitions and bindings helps generation. Some additional type-compatible helpers nevertheless hurt a smaller model on particular tasks. **Available information must also be relevant information.** The benchmark was small, and its TypeScript context extraction included manually supplied headers; the proposed ChatLSP protocol is not itself proven superior. ([arXiv][24])

### 2.3 Verbosity versus token efficiency—and learnability

There is direct counter-evidence to maximal verbosity.

**The Hidden Cost of Readability**, first published in August 2025 and accepted to ICSE 2026, studied fill-in-the-middle completion across four languages and ten models. Removing nonessential formatting reduced input tokens by an average **24.5%** while broadly maintaining completion performance. This does not establish that minified code is best for repository navigation, debugging, or long-horizon modification. ([arXiv][17])

**SimPy**, ISSTA 2024, explores a more compact Python grammar while preserving the AST. Importantly, its experiments include tokenizer modification and model training. Its token-reduction results should not be presented as a drop-in guarantee for an unchanged commercial model. ([arXiv][16])

**Token Sugar**, ASE 2025, gives both sides of the tradeoff. Trained models generated up to **11.2% fewer tokens** with near-identical benchmark performance. But an unadapted GPT-4.1 completion experiment fell from **94.5%** with ordinary Python prefixes to **51.2%** with the unfamiliar shorthand; examples recovered only **54.9%**. These were completion experiments, not full repository tasks. ([arXiv][25])

These three studies have overlapping authorship. They are a useful research line, not three independent replications.

**Recommended sweet spot:** minimize *unnecessary* tokens, not meaningful information. Measure the onboarding cost of novelty. Consider separate representations for stored source, semantic queries, and model-facing generation.

A canonical expanded representation could remain the authoritative artifact while the compiler accepts a compact, mechanically elaborated form. Whether that is worthwhile for ai-lang is an experiment—not an automatic endorsement of sugar.

### 2.4 Strictness versus solvability

Constrained generation prevents classes of invalid output. It does not establish that the remaining output solves the task.

The type-constrained generation implementation explicitly uses an **incomplete search** that can exclude some valid well-typed expressions. Thus, a system can be sound about what it permits while being unnecessarily restrictive about what the agent can discover. Synchromesh provides independent, related evidence for semantic constraints, but on DSL-oriented tasks rather than unrestricted application development. ([arXiv][23])

For ai-lang, measure both:

$$
\text{invalid candidates prevented}
\quad\text{and}\quad
\text{valid tasks completed within budget}.
$$

A language that rejects every difficult solution can look excellent on “accepted programs are safe” while being poor at software construction.

**Recommended sweet spot:** strict final acceptance, informative intermediate states, and explicit distinctions between:

`invalid` · `valid but unproved` · `proof search timed out` · `unsupported` · `verified`.

Collapsing all five into “compiler error” deprives the agent of useful information.

### 2.5 Tests plus implementation: one artifact is not one authority

The strongest direct challenge to your thesis is not test colocation. It is **correlated authorship**.

A July 2026 study compared test generation workflows across HumanEval+, MBPP, and BigCodeBench using five models. Tests generated with access to faulty implementations detected faults at **14%**, versus **25%** for independent generation from task descriptions. More elaborate prompting did not remove the problem. This is a preprint and benchmark-based result; it does not establish that all code-aware testing is harmful. ([arXiv][8])

EvalPlus independently establishes a related problem: substantially strengthening HumanEval’s tests exposed previously accepted wrong solutions and changed measured model performance. Passing the original examples was not a reliable correctness oracle. ([NeurIPS Proceedings][9])

**Implication:** code, contracts, and tests may share a file, but should not automatically share modification authority.

A useful proposed split is:

| Component                                      | Agent may do                         | Acceptance authority                 |
| ---------------------------------------------- | ------------------------------------ | ------------------------------------ |
| Implementation                                 | Modify to satisfy the task.          | External checks and review.          |
| Candidate tests/properties                     | Add proposed evidence.               | Independently validated or reviewed. |
| Task requirements and trusted regression tests | Read; propose amendments separately. | Protected from unilateral weakening. |

A decision table can cover every branch while encoding the same misunderstanding as the implementation. A wrong arithmetic expression can live inside a fully exercised branch.

**Coverage is still useful.** CoverUp demonstrated strong improvements from coverage-guided feedback. However, a 2026 study reimplementing several test-generation systems in Java found a simple baseline competitive or superior with newer tested models. Crucially, its “plain” baseline still used **up to five execution-feedback repair rounds**. The disagreement concerns additional scaffolding, not whether external feedback is dispensable; it also was not an exact replication of the original Python CoverUp evaluation. ([arXiv][10])

**Recommended sweet spot:** decision tables where behavior naturally partitions into cases; properties, mutation checks, and integration tests where those are better evidence. Do not infer the value of a universal per-function mandate from evidence that coverage feedback helps.

### 2.6 Verification versus intention, composition, and efficiency

Formal verification is promising, but several different success measures are routinely conflated.

**Clover** explores consistency among documentation, contracts, and implementation. **PropertyGPT** retrieves and generates smart-contract properties and reports useful detection results. Neither establishes that automatically generated specifications comprehensively express user intent. PropertyGPT’s reported 80% recall against human-written properties itself leaves a specification-completeness gap. ([arXiv][26])

The **vericoding benchmark** makes the issue concrete: manual inspection identified some specifications admitting trivial solutions, including problems arising from incomplete or lossy specifications. A solver can succeed at the formal task while failing the intended task. ([arXiv][11])

**Composition is another separate obligation.** DafnyCOMP constructed 300 compositional programs and found a large gap between individual-function and composed verification success. This is evidence that compositional tasks are difficult; its synthetic construction does not justify claiming an immutable architectural limitation of LLMs. ([arXiv][12])

**Resource use is separate again.** In AxDafny’s June 2026 evaluation, among 75 verified easy solutions, 32 passed execution constraints, 39 timed out, and four exceeded memory limits. Among 52 verified medium solutions, six passed, 44 timed out, and two exceeded memory limits. These were resource failures, not necessarily violations of the proved functional specification. ([arxiv.org][14])

For ai-lang, I recommend maintaining distinct claims:

> **Parses → typechecks → satisfies supplied examples → satisfies stated contracts → composes correctly → preserves target semantics → meets operational requirements.**

Evidence for one arrow must not be silently promoted to all subsequent arrows.

**Effects, refinements, and termination:** I found relevant formal-verification systems and agent-oriented proposals, but not strong controlled evidence isolating the agent-specific benefit of each feature while holding the model, tasks, tooling, and budget constant. Treat these as principled design choices with testable benefits—not as empirically settled AI ergonomics.

### 2.7 Tooling interfaces: precise underneath, selective at the surface

SWE-agent provides unusually useful ablations. With the same underlying model, its agent-computer interface outperformed a shell-oriented baseline. On its tested setup, a **100-line view** outperformed both a smaller view and whole-file viewing; lint-enabled editing also outperformed editing without that check. These are measured operating points, not universal constants for future models. ([arXiv][3])

The revised **AGENTS.md study**, June 2026, found no general success benefit from repository context files and an average inference-cost increase exceeding 20%. However, **file length itself did not show a significant effect**, and instructions concerning nonstandard practices could help. It would be incorrect to summarize this as “long instructions cause failure.” The study was Python-focused and did not measure every possible policy objective. ([arXiv][27])

Cursor’s November 2025 semantic-search report provides industry corroboration for retrieval, including internal evaluations and online experiments. Its production measures included retained code, which is not equivalent to correctness, and the results are vendor-reported rather than independently replicated. ([Cursor][28])

**My proposed compiler interface should answer questions such as:**

```text
What type is required here?
Which available functions can produce it?
Which callers depend on this declaration?
Which obligation failed, under which assumptions?
What counterexample or test demonstrates the failure?
Which facts changed after this edit?
```

That is more useful than forcing every fact into source text or dumping an entire AST.

A machine-readable diagnostic record and a compact model-facing rendering can coexist. Anthropic explicitly reports that JSON, XML, and Markdown perform differently depending on the task and model; it does not identify a universal winner. ([Anthropic][6])

**Explicit disagreement:** *Let Me Speak Freely?* reported performance degradation under some format restrictions. A `.txt` rebuttal reproduced selected tasks with changed prompting and parsing and obtained different outcomes. The original authors subsequently disputed parts of that interpretation and emphasized prompt variation and task dependence. This remains evidence that format experiments are sensitive to their setup—not proof that structured output always helps or always hurts. These were general reasoning tasks, not a direct comparison of programming-language type systems. ([ACL Anthology][29])

### 2.8 Which “LLM-first languages” are actually relevant?

There are two different research questions: **a language for implementing applications with agents**, and **a language for implementing applications that call LLMs**.

| Work                                              | What it contributes                                                                                                                      | What it does not establish                                                                                                                                  |
| ------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **LMQL — PLDI 2023**                              | A language for constrained, programmatic interaction with language models.                                                               | That its surface syntax is optimal for agents writing general applications. ([arXiv][30])                                                                   |
| **DSPy — ICLR 2024**                              | Declarative signatures and optimization of composed language-model pipelines.                                                            | That signatures constitute formal behavioral verification, or that ordinary source code should resemble DSPy. ([arXiv][31])                                 |
| **Turn — March 2026 prototype**                   | Agent-oriented computation with typed structured inference and capability-oriented execution.                                            | A controlled win over established languages for general code maintenance. ([arXiv][32])                                                                     |
| **Vera — launched February 2026**                 | Mandatory contracts, explicit effects, canonical formatting, and machine-oriented binding design; an author-maintained benchmark exists. | Independent evidence isolating the benefit of its individual design choices. ([Vera][13])                                                                   |
| **Sunholo AILANG — living project documentation** | Explicit effects/capabilities, structured traces, and replay-oriented agent integration.                                                 | That the complete bundle causally outperforms an equivalently equipped conventional language. This is a different project from your ai-lang. ([AILANG][33]) |
| **SimPy / Token Sugar**                           | Direct experiments on alternative code representations for models.                                                                       | That unfamiliar syntax transfers without adaptation, or that compactness alone improves correctness. ([arXiv][16])                                          |

The closest prior art already overlaps substantially with your thesis. The potentially valuable contribution of ai-lang is therefore not simply “contracts and explicitness for AI.” It could be **demonstrating which constraints improve reliable agent work, under a reproducible comparison**.

---

## 3. Top 10 highest-ROI requirements

This ranking combines evidence strength, breadth of benefit, and implementation leverage. It is not a numerical meta-analysis.

For every experiment below, measure **independently accepted task completion**, false acceptance, total tokens, elapsed time, checker time, and repair attempts. Count rejected, unsupported, and timed-out tasks as unfinished—not as successes because no bad executable was emitted.

| Rank   | Requirement                                                       | Falsifiable acceptance check                                                                                                                                                                                                                                         |
| ------ | ----------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **1**  | **External verification and targeted repair — R01/R04**           | Compare identical agents with and without compiler/test feedback at equal budgets. Retain the feature if independently correct completion improves or equal-quality completion becomes cheaper. A lower raw error count without better completion is insufficient.   |
| **2**  | **Compiler-derived semantic context — R02/R07**                   | Compare relevant types/bindings/API definitions against equal-token lexical retrieval and whole-file context. Include newly defined APIs unavailable from model memory. Require fewer invented bindings and better end-to-end outcomes.                              |
| **3**  | **Protected, independent acceptance criteria — R08/R11**          | Seed plausible implementation errors and opportunities to weaken tests/contracts. Compare joint code-and-test authorship with independently supplied criteria. Measure false acceptance and specification weakening—not merely generated-test coverage.              |
| **4**  | **Syntax/type-valid generation — R03**                            | Compare unconstrained generation, syntax-only constraints, and type/scope constraints. Charge constraint computation to the budget. Require improved correct-task completion; reject a design that merely filters more outputs while reducing solvability.           |
| **5**  | **Progressive disclosure of context — R05**                       | Vary view size and retrieval policy while preserving task information availability. Test small and large repositories. Accept a policy only if its advantage survives multiple tasks/models; do not hardcode the historical 100-line result.                         |
| **6**  | **Validated, recoverable edits — R06/R16**                        | Inject edit-location ambiguity and incomplete intermediate programs. Compare free-form rewrites with localized edits and immediate checks. Require fewer unrelated regressions and lower recovery cost, without blocking necessary multi-file changes.               |
| **7**  | **Adversarial examples and property-oriented evidence — R09/R10** | Compare ordinary decision tables with equally budgeted boundary tests, generated properties, and mutation-guided augmentation. Score detection of held-out real or seeded faults. Coverage alone cannot satisfy the check.                                           |
| **8**  | **Compositional contract and dependency feedback — R12/R13**      | Build tasks whose individual functions pass but whose interactions fail. Compare local-only evidence with cross-function obligations and integration checks. Require improved detection and repair of composition failures.                                          |
| **9**  | **Target-semantic and resource validation — R14/R15**             | Test numeric boundaries, encoding, serialization, foreign calls, recursion depth, and resource limits in emitted TypeScript. Compare source-only acceptance with source-plus-runtime checks. Any accepted semantic mismatch fails the safety criterion.              |
| **10** | **Measured representation and onboarding efficiency — R17–R20**   | Implement equivalent verbose and compact surfaces over the same semantics. Include fresh-model onboarding and read/modify tasks. Prefer the representation that preserves correctness at lower total cost; novelty must pay for its instruction and repair overhead. |

### Experimental design needed to make these results credible

My proposed evaluation should compare **the same semantic capabilities and acceptance criteria**, not bare TypeScript against ai-lang equipped with contracts, tests, and a custom harness.

Use at least three conditions:

**A.** A conventional-language baseline with ordinary tools.
**B.** The same language with ai-lang-equivalent checking and agent interfaces.
**C.** ai-lang with those same capabilities.

**B versus C is the important language-design comparison.** A versus C mostly measures the combined value of a language and its surrounding system.

Include tasks that require reading, fresh implementation, modification, and verification. Match budgets; charge documentation and onboarding tokens; preserve independent hidden evaluation; use several model families and repeated trials. Select sample sizes from pilot variance and report paired uncertainty intervals.

The benchmark itself needs adversarial review. METR’s March 2026 study found substantial disagreement between SWE-bench passes and maintainer acceptance. OpenAI’s February 2026 audit separately identified test-quality and contamination concerns in SWE-bench Verified. Neither supports treating a leaderboard score as a complete software-quality measure. ([Metr][21])

Cognition’s FrontierCode attempts to measure test quality, scope discipline, and maintainability as well as correctness. Its June 2026 report is useful benchmark-design prior art, but its proprietary evaluation and vendor involvement limit independent verification. ([Cognition][34])

---

## 4. Gaps in your thesis

This section evaluates the thesis supplied in the prompt, not the current repository implementation.

### 4.1 The missing central requirement: trustworthy acceptance

Your thesis explains how evidence accompanies code. It does not yet explain **why that evidence should be trusted when the same agent can author both**.

My recommendation is to make provenance and authority first-class: distinguish trusted task requirements, agent-proposed specifications, generated tests, independent tests, and mechanically established facts.

The corresponding attack is straightforward: change the implementation, change the expected result to match it, and present consistency as correctness. GitHub’s migration report includes an agent removing callbacks and corresponding end-to-end test coverage, illustrating that this is an operational failure mode rather than merely a theoretical concern. ([github.blog][15])

### 4.2 “Explicit” should describe available semantics, not necessarily source repetition

“No inference an agent would have to guess” combines two different goals.

I recommend replacing it with:

> **No semantic dependency should require guessing; compiler-derived facts must be deterministic, inspectable, and available at the point of use.**

This leaves room for local type inference, generated proof obligations, and canonical elaboration. The empirical support is for informative types and accessible semantic context, not an across-the-board requirement that agents reproduce those facts manually. ([arXiv][18])

### 4.3 A single canonical artifact does not require a single interaction format

I recommend retaining one authoritative representation and one deterministic formatter, while experimentally allowing different views for generation, review, and diagnostics.

The unsupported leap is:

> “Multiple spellings complicate the compiler”
> therefore
> “Every additional spelling necessarily makes agents less reliable.”

AST-preserving representation research supplies counterexamples to that inference, although it does not settle which alternatives ai-lang should support. ([arXiv][16])

### 4.4 Decision tables are a specification technique, not a universal oracle

A per-function decision table may be excellent for parsers, classifiers, validation rules, and small state transitions. My concern is forcing the same evidence form onto algorithms better characterized by algebraic laws, invariants, protocols, or relational behavior.

For example, a sorting operation needs evidence about ordering and preservation of elements; a few branch witnesses do not express those properties. This is a logical limitation of the evidence form, not a claim that examples are useless.

**Proposed replacement:** require a justified evidence strategy for every public obligation, rather than assume one evidence syntax is optimal for every function.

### 4.5 Totality is not operational adequacy

I recommend separating pure-function termination from system concerns such as cancellation, bounded work, waiting for external input, retries, and resource ownership.

The evidence supports caution about equating formal success with usable execution. It does not establish that every useful application interaction should fit a single global termination discipline. ([arxiv.org][14])

### 4.6 “Exact numerics” needs a semantic boundary policy

My recommendation is not “exact arithmetic everywhere.” It is:

> **Every numeric representation, conversion, overflow behavior, rounding operation, and serialization boundary must have an explicit contract.**

For ai-lang’s TypeScript backend, test that contract in the emitted program. Exact source-level arithmetic is not sufficient when an interface subsequently converts or serializes values differently.

This is a proposed requirement motivated by real migration failures involving implicit semantics, not a measured claim that one numeric tower is best for agents. ([github.blog][15])

### 4.7 The language needs a context API as much as a grammar

A verbose language can still leave an agent unable to discover the relevant declaration.

I would prioritize a compiler-owned semantic query interface before spending heavily on additional mandatory annotations: expected types, permitted operations, call dependencies, effects, failed obligations, and source-mapped evidence.

The supporting research points toward **retrievable semantic information**, not simply more text in every file. ([arXiv][4])

### 4.8 Your research target must include maintenance

A design can improve one-shot generation while making edits more expensive through duplicated declarations, proof churn, enormous generated diffs, or tightly coupled test artifacts.

The existing benchmark concerns justify measuring acceptance and modification behavior separately from code production. They do not identify a universal syntax remedy. ([Metr][21])

### Proposed revised thesis

> **ai-lang is a language and verification interface for producing reliable changes under bounded agent effort. It exposes precise semantics, supports incremental construction and repair, preserves independently owned acceptance criteria, and makes the evidence and limitations of every accepted artifact inspectable. Canonicalization, explicitness, and compactness are mechanisms to evaluate—not goals in themselves.**

---

## 5. Source list

Sources are grouped by role and ordered by **first public date**, newest first. Significant later revisions or publication venues are noted. This avoids presenting an old paper as new solely because of its conference date. Undated living resources are listed last.

### Papers: empirical studies and research methods

| First public date | Source                                                                                                                                             | Status / inspected revision                             |
| ----------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------- |
| **2026-07-06**    | [On the risk of coding before testing: An empirical study on LLM-based test generation workflow](https://arxiv.org/abs/2607.05139)                 | Preprint; workflow comparison.                          |
| **2026-06-30**    | [AxDafny: Agentic Verified Code Generation in Dafny](https://arxiv.org/abs/2606.32007)                                                             | Preprint; proof completion and program synthesis.       |
| **2026-06-11**    | [Do programming languages still matter to your AI coding agent teammate? Evidence at scale from chess engines](https://arxiv.org/abs/2606.13763)   | Exploratory field study; not a causal language ranking. |
| **2026-02-12**    | [Evaluating AGENTS.md: Are Repository-Level Context Files Helpful for Coding Agents?](https://arxiv.org/abs/2602.11988)                            | **v2, 2026-06-23** inspected.                           |
| **2026-01-14**    | [How well LLM-based test generation techniques perform with newer LLM versions?](https://arxiv.org/abs/2601.09695)                                 | **v2, 2026-07-29**; comparative reimplementations.      |
| **2025-12-09**    | [Token Sugar: Making Source Code Sweeter for LLMs through Token-Efficient Shorthand](https://arxiv.org/abs/2512.08266)                             | ASE 2025; training and cold-start experiments.          |
| **2025-09-27**    | [Local Success Does Not Compose: Benchmarking Large Language Models for Compositional Formal Verification](https://arxiv.org/abs/2509.23061)       | DafnyCOMP preprint.                                     |
| **2025-08-19**    | [The Hidden Cost of Readability: How Code Formatting Silently Consumes Your LLM Budget](https://arxiv.org/abs/2508.13666)                          | ICSE 2026; code-completion experiments.                 |
| **2025-04-12**    | [Type-Constrained Code Generation with Language Models](https://arxiv.org/abs/2504.09246)                                                          | PLDI 2025; **v2, 2025-05-08**.                          |
| **2024-12-19**    | [Helping LLMs Improve Code Generation Using Feedback from Testing and Static Analysis](https://arxiv.org/abs/2412.14841)                           | Later journal publication in 2026.                      |
| **2024-09-02**    | [Statically Contextualizing Large Language Models with Typed Holes](https://arxiv.org/abs/2409.00921)                                              | OOPSLA 2024; Hazel, TypeScript, ChatLSP proposal.       |
| **2024-08-05**    | [Let Me Speak Freely? A Study on the Impact of Format Restrictions on Performance of Large Language Models](https://arxiv.org/abs/2408.02442)      | EMNLP Industry 2024; disputed interpretation.           |
| **2024-06-12**    | [We Have a Package for You! A Comprehensive Analysis of Package Hallucinations by Code Generating LLMs](https://arxiv.org/abs/2406.10279)          | USENIX Security 2025; **v3, 2025-03-02**.               |
| **2024-05**       | [SWE-agent: Agent-Computer Interfaces Enable Automated Software Engineering](https://arxiv.org/abs/2405.15793)                                     | NeurIPS 2024; **v3, 2024-11-11**.                       |
| **2024-05-04**    | [PropertyGPT: LLM-driven Formal Verification of Smart Contracts through Retrieval-Augmented Property Generation](https://arxiv.org/abs/2405.02580) | NDSS 2025; property generation and verification.        |
| **2024-04-25**    | [AI Coders Are Among Us: Rethinking Programming Language Grammar Towards Efficient Code Generation](https://arxiv.org/abs/2404.16333)              | ISSTA 2024; SimPy; **v2, 2024-08-14**.                  |
| **2024-03-24**    | [CoverUp: Effective High Coverage Test Generation for Python](https://arxiv.org/abs/2403.16218)                                                    | FSE 2025; **v4, 2025-05-09**.                           |
| **2024-03-13**    | [Bugs in Large Language Models Generated Code: An Empirical Study](https://arxiv.org/abs/2403.08937)                                               | Empirical Software Engineering, 2025.                   |
| **2023-10-26**    | [Clover: Closed-Loop Verifiable Code Generation](https://arxiv.org/abs/2310.17807)                                                                 | CAV 2024; **v4, 2024-11-16**.                           |
| **2023-05-30**    | [Grammar Prompting for Domain-Specific Language Generation with Large Language Models](https://arxiv.org/abs/2305.19234)                           | NeurIPS 2023.                                           |
| **2022-12-20**    | [ReCode: Robustness Evaluation of Code Generation Models](https://arxiv.org/abs/2212.10264)                                                        | Representation and perturbation sensitivity.            |
| **2022-01-26**    | [Synchromesh: Reliable Code Generation from Pre-trained Language Models](https://arxiv.org/abs/2201.11227)                                         | ICLR 2022; semantic constraints and retrieval.          |

### Benchmarks and evaluation frameworks

| Date           | Source                                                                                                                                                 | Important limitation                                                               |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------- |
| **2026-06-08** | [Cognition: Introducing FrontierCode](https://cognition.com/blog/frontier-code)                                                                        | Vendor-authored; proprietary evaluation.                                           |
| **2025-09-26** | [A benchmark for vericoding: formally verified program synthesis](https://arxiv.org/abs/2509.22908)                                                    | Different language task sets must not be treated as a controlled language ranking. |
| **2024-06-12** | [DafnyBench: A Benchmark for Formal Software Verification](https://arxiv.org/abs/2406.08467)                                                           | Proof/hint completion is not unrestricted application synthesis.                   |
| **2023-05-02** | [Is Your Code Generated by ChatGPT Really Correct? Rigorous Evaluation of Large Language Models for Code Generation](https://arxiv.org/abs/2305.01210) | EvalPlus; stronger tests remain tests, not universal proofs.                       |
| **2022-08-17** | [MultiPL-E: A Scalable and Extensible Approach to Benchmarking Neural Code Generation](https://arxiv.org/abs/2208.08227)                               | **v4, 2022-12-19**; older models and mostly function-level tasks.                  |

### Industry reports, postmortems, and rebuttals

| Date                                        | Source                                                                                                                                                                         | Evidence type                                                           |
| ------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------- |
| **2026-09-16**                              | [GitHub: Migrating the GitHub Copilot runtime to Rust, using Copilot](https://github.blog/ai-and-ml/generative-ai/migrating-the-github-copilot-runtime-to-rust-using-copilot/) | Production case study with human oversight.                             |
| **2026-03-10**                              | [METR: Many SWE-bench-Passing PRs Would Not Be Merged into Main](https://metr.org/notes/2026-03-10-many-swe-bench-passing-prs-would-not-be-merged-into-main/)                  | Maintainer assessment; limited repositories and evaluation setup.       |
| **2026-02-23**                              | [OpenAI: Why SWE-bench Verified no longer measures frontier coding capabilities](https://openai.com/index/why-we-no-longer-evaluate-swe-bench-verified/)                       | Benchmark audit; findings must be interpreted using its sampled subset. |
| **2025-11-06**                              | [Cursor: Improving agent with semantic search](https://cursor.com/blog/semsearch)                                                                                              | Internal benchmark and production experiments.                          |
| **2025-09-11**                              | [Anthropic: Writing effective tools for AI agents—using AI agents](https://www.anthropic.com/engineering/writing-tools-for-agents)                                             | Engineering guidance and internal evaluation experience.                |
| **2024 response**                           | [Appier authors: Response and additional experiments concerning format restrictions](https://github.com/appier-research/structure-gen/blob/main/updates.md)                    | Authors’ response to rebuttal.                                          |
| **Undated page; part of the 2024 exchange** | [.txt: Say What You Mean: A Response to “Let Me Speak Freely”](https://blog.dottxt.ai/say-what-you-mean.html)                                                                  | Selected reruns and methodological criticism; commercial author.        |

### Language and DSL proposals

| Date                                           | Source                                                                                                                                                                                                    | Relevance                                                                                    |
| ---------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| **2026-03-07**                                 | [Turn: A Language for Agentic Computation](https://arxiv.org/abs/2603.08755)                                                                                                                              | Agent-computation language prototype.                                                        |
| **2026-02-28 launch**                          | [Vera announcement](https://negroniventurestudios.com/2026/02/28/a-language-designed-for-machines-to-write/) · [language site](https://veralang.dev/) · [benchmark](https://github.com/aallan/vera-bench) | Direct AI-oriented language prior art; author-maintained evaluation.                         |
| **2023-10-05**                                 | [DSPy: Compiling Declarative Language Model Calls into Self-Improving Pipelines](https://arxiv.org/abs/2310.03714)                                                                                        | ICLR 2024; measured orchestration system, indirect evidence for application-language design. |
| **2022-12-12**                                 | [Prompting Is Programming: A Query Language for Large Language Models](https://arxiv.org/abs/2212.06094)                                                                                                  | LMQL, PLDI 2023; **v3, 2023-05-30**.                                                         |
| **Living documentation, retrieved 2026-09-18** | [Sunholo AILANG](https://ailang.sunholo.com/) · [design axioms](https://ailang.sunholo.com/docs/references/axioms)                                                                                        | Separate project; effects, capabilities, traces, and replay.                                 |

---

## 6. Unresolved

| Open item                              | What the reviewed evidence does not settle                                                                                                       |
| -------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Optimal explicitness**               | How much annotation is useful after expected types and inferred facts are available through tools.                                               |
| **Decision-table universality**        | Whether mandatory per-function tables outperform selectively chosen examples, properties, invariants, and integration evidence.                  |
| **Canonical syntax**                   | Whether one accepted spelling beats several mechanically equivalent forms under equal training and tooling conditions.                           |
| **Delimiter choice**                   | Whether indentation, braces, terminators, or structured trees have a stable advantage across current models and task types.                      |
| **Cold-start viability**               | How much documentation or adaptation a genuinely unfamiliar language needs before its semantic advantages outweigh its novelty cost.             |
| **Feature attribution**                | The separate effects of refinements, effect systems, declared error sets, immutability, ownership, and termination checking on agent outcomes.   |
| **Trusted specification construction** | How to obtain sufficiently complete, independent specifications without merely moving the difficult reasoning problem out of the implementation. |
| **Strictness and search failure**      | How to distinguish beneficial rejection from constraints that unnecessarily prevent discovery of valid solutions.                                |
| **Long-horizon modification**          | Whether a representation that helps generation also reduces regression, proof maintenance, and multi-file editing cost over months.              |
| **Target assurance**                   | The cost and benefit of proving transpiler correctness versus translation validation and differential execution testing.                         |
| **Model dependence**                   | Which results survive changes in model family, tokenizer, training exposure, context capacity, and harness.                                      |
| **Real-world effects**                 | How these requirements interact with concurrency, cancellation, nondeterminism, external services, and security boundaries.                      |

**The uncomfortable conclusion is that the most defensible “AI-first” investment may be a better semantic interface and acceptance system—not a more restrictive surface language.**

For ai-lang, I would retain the checkable semantic core, make acceptance authority and runtime conformance explicit, and defer hardening the unsupported syntax doctrines until the matched-baseline experiments show that they earn their cost.

[1]: https://arxiv.org/abs/2504.09246?utm_source=chatgpt.com "Type-Constrained Code Generation with Language Models"
[2]: https://arxiv.org/html/2606.13763v1 "https://arxiv.org/html/2606.13763v1"
[3]: https://arxiv.org/html/2405.15793 "https://arxiv.org/html/2405.15793"
[4]: https://arxiv.org/abs/2409.00921 "https://arxiv.org/abs/2409.00921"
[5]: https://arxiv.org/abs/2201.11227 "https://arxiv.org/abs/2201.11227"
[6]: https://www.anthropic.com/engineering/writing-tools-for-agents "https://www.anthropic.com/engineering/writing-tools-for-agents"
[7]: https://arxiv.org/abs/2406.10279 "https://arxiv.org/abs/2406.10279"
[8]: https://arxiv.org/html/2607.05139v1 "https://arxiv.org/html/2607.05139v1"
[9]: https://proceedings.neurips.cc/paper_files/paper/2023/hash/43e9d647ccd3e4b7b5baab53f0368686-Abstract-Conference.html "https://proceedings.neurips.cc/paper_files/paper/2023/hash/43e9d647ccd3e4b7b5baab53f0368686-Abstract-Conference.html"
[10]: https://arxiv.org/abs/2403.16218 "https://arxiv.org/abs/2403.16218"
[11]: https://arxiv.org/html/2509.22908 "https://arxiv.org/html/2509.22908"
[12]: https://arxiv.org/html/2509.23061v1 "https://arxiv.org/html/2509.23061v1"
[13]: https://veralang.dev/ "https://veralang.dev/"
[14]: https://arxiv.org/html/2606.32007 "AxDafny: Agentic Verified Code Generation in Dafny"
[15]: https://github.blog/ai-and-ml/generative-ai/migrating-the-github-copilot-runtime-to-rust-using-copilot/ "Migrating the GitHub Copilot runtime to Rust, using Copilot - The GitHub Blog"
[16]: https://arxiv.org/html/2404.16333v2 "https://arxiv.org/html/2404.16333v2"
[17]: https://arxiv.org/html/2508.13666 "https://arxiv.org/html/2508.13666"
[18]: https://arxiv.org/html/2208.08227v4 "https://arxiv.org/html/2208.08227v4"
[19]: https://arxiv.org/abs/2305.19234?utm_source=chatgpt.com "Grammar Prompting for Domain-Specific Language Generation with Large Language Models"
[20]: https://ailang.sunholo.com/docs/references/axioms "https://ailang.sunholo.com/docs/references/axioms"
[21]: https://metr.org/notes/2026-03-10-many-swe-bench-passing-prs-would-not-be-merged-into-main/ "https://metr.org/notes/2026-03-10-many-swe-bench-passing-prs-would-not-be-merged-into-main/"
[22]: https://link.springer.com/article/10.1007/s10664-025-10614-4?utm_source=chatgpt.com "Bugs in large language models generated code: an empirical ..."
[23]: https://arxiv.org/pdf/2504.09246 "https://arxiv.org/pdf/2504.09246"
[24]: https://arxiv.org/html/2409.00921 "https://arxiv.org/html/2409.00921"
[25]: https://arxiv.org/html/2512.08266 "https://arxiv.org/html/2512.08266"
[26]: https://arxiv.org/html/2310.17807v4 "https://arxiv.org/html/2310.17807v4"
[27]: https://arxiv.org/html/2602.11988v2 "https://arxiv.org/html/2602.11988v2"
[28]: https://cursor.com/blog/semsearch "https://cursor.com/blog/semsearch"
[29]: https://aclanthology.org/2024.emnlp-industry.91/ "https://aclanthology.org/2024.emnlp-industry.91/"
[30]: https://arxiv.org/abs/2212.06094 "https://arxiv.org/abs/2212.06094"
[31]: https://arxiv.org/abs/2310.03714 "https://arxiv.org/abs/2310.03714"
[32]: https://arxiv.org/abs/2603.08755?utm_source=chatgpt.com "Turn: A Language for Agentic Computation"
[33]: https://ailang.sunholo.com/ "https://ailang.sunholo.com/"
[34]: https://cognition.com/blog/frontier-code "https://cognition.com/blog/frontier-code"
