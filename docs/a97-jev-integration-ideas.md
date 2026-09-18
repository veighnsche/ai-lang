**Scope:** these are proposals for the site surfaces you described, not claims about an inspected website. The repository exposes no homepage URL, and search did not establish a deployed site.

**Sizing:** value and effort use S/M/L; effort means a feature prototype on an existing site. **Dogfooding** means `.can` judgment/routing logic with conventional web infrastructure, once native judgments land. All such logic gets complete-observation exchange tests; site-specific model accuracy and calibration need separate evaluation—not a green compiler gate.

## Ideas

1. **Name — Rule-Aware Search**
   **Pitch:** Search “why does my unused branch fail?” and receive the relevant rule, explanation, and runnable example—not merely pages containing “branch.”
   **How it works:** On docs search, retrieve candidates lexically, filter by the selected language revision, then use probabilistic semantic-relevance scoring to rerank a fixed shortlist. Display source excerpts; retain lexical ordering when judgments are uncertain.
   **Why probability:** Queries and documentation use different vocabulary, and several passages may be partially relevant; keyword rules cannot reliably distinguish them.
   **Dogfood angle:** **Yes.** Fan out `score` judgments for relevance and answer completeness; `.can` combines the scores and handles uncertain or unavailable judgments explicitly. Version eligibility remains deterministic.
   **Value/effort + data note:** **Value L / effort M.** Starts with public docs and maintainer-written queries, without click history. Server-side processing; disclose external query processing and offer lexical-only search. Do not retain raw queries by default.

2. **Name — One Good Question**
   **Pitch:** An ambiguous search asks one useful clarification instead of confidently sending the visitor down the wrong path.
   **How it works:** For queries such as “proof failed,” a probabilistic intent classifier assigns weight to termination, effect authority, contract verification, and decision-table evidence. A cost table chooses between showing results and asking a prepared discriminator: “Did compilation report an untaken arm or an unproven contract?”
   **Why probability:** Ambiguity should affect the interaction; a rigid keyword router must choose even when two interpretations are plausible.
   **Dogfood angle:** **Yes.** `choice` supplies the intent distribution; exact-decimal expected-cost calculations select a fixed clarification card. Low confidence never produces an automatic redirect.
   **Value/effort + data note:** **Value M / effort S**, after search exists. Needs an intent taxonomy and labeled ambiguous queries. Keep the query and answer within the session; no identity or navigation history is needed.

3. **Name — Current-Rule Lens**
   **Pitch:** Historical amendment pages show what changed later, with the original wording preserved beside the operative rule.
   **How it works:** A maintainer dashboard retrieves related passages and uses pairwise natural-language inference—supports, contradicts, unrelated—to propose supersession links. Maintainers approve links before the site displays historical/current annotations. Start with the README’s a12 status and the docs index’s older a14 summary.
   **Why probability:** Semantic contradictions can survive valid links and correct filenames; deterministic checks still own explicit version and status metadata.
   **Dogfood angle:** **Yes.** `choice` classifies passage relationships; `score` prioritizes likely visitor confusion. Approval, not model confidence, establishes the displayed relationship.
   **Value/effort + data note:** **Value L / effort S** for a small review queue. Uses public repository text only. Run analysis outside the website build and publish approved static annotations; no visitor tracking.

4. **Name — Find My Example**
   **Pitch:** Describe what you want to build and receive the smallest relevant existing example, already opened in the playground.
   **How it works:** An example-gallery classifier maps a short task description to a fixed catalog: login/retry, bounded retry, private counter, quota validation, or exact division. Strong matches open a prepared card; ambiguous matches show two candidates. These examples already have repository sources.
   **Why probability:** “Retry this without hanging forever” expresses a use case, not the documentation term `decreases`.
   **Dogfood angle:** **Yes.** `choice` selects the example family; `score` distinguishes “direct match” from “useful starting point.” Code is retrieved, not generated.
   **Value/effort + data note:** **Value M / effort S.** Needs curated example descriptions and a small task-label set. Server-side task text only; invite descriptions rather than proprietary code, with no account required.

5. **Name — Receipt Sandbox**
   **Pitch:** Visitors change a confidence threshold and watch the same recorded judgment take a different, fully explained branch.
   **How it works:** A landing-page demo shows a message, its complete Jev observation, routing policy, and selected outcome. Recorded mode executes ordinary scripted playground cases; an explicit live button acquires one new observation through the server. Include tied Choice results, uncertain Noul, and equal-mean/different-shape Scores.
   **Why probability:** It demonstrates why a label alone is insufficient: uncertainty and distribution shape can change application behavior.
   **Dogfood angle:** **Yes—directly.** Use all three primitives and their ordinary exchange scripts. Public teaching receipts are intentional website content, separate from stripped production test tables.
   **Value/effort + data note:** **Value L / effort M.** Needs curated receipts and playground integration, not traffic history. Default to public sample text and recorded mode; live custom text requires explicit submission, server-held credentials, and a request budget.

6. **Name — Diagnostic Companion**
   **Pitch:** Beside a compiler error, the playground offers the most relevant explanation for this particular mistake.
   **How it works:** Deterministically select the official explanation by CAN code, then classify the likely misconception from the diagnostic and a visitor-selected snippet. For CAN4107, distinguish a missing witness from a branch contradicted by its enclosing confidence guard; attach an existing example rather than inventing a fix. 
   **Why probability:** One diagnostic can arise from different misunderstandings; the code identifies the violated rule, not necessarily the lesson the author needs.
   **Dogfood angle:** **Yes.** `choice` selects a misconception category and confidence-routes to a targeted card or the unmodified official explanation. The real compiler remains authoritative.
   **Value/effort + data note:** **Value L / effort M.** Begin with curated broken examples and five common diagnostics. Analyze personal code only after explicit opt-in and snippet preview; avoid source retention.

7. **Name — Hint Ladder**
   **Pitch:** Tutorial hints improve according to which explanations actually help visitors solve the next step.
   **How it works:** Within an exercise/error category, a Thompson-sampling bandit chooses among three approved hint styles: explain the rule, show a smaller example, or ask a guiding question. Reward is subsequent exercise completion, not hint clicks; retain a fixed evaluation holdout. Thompson sampling supplies the exploration/exploitation mechanism. ([arXiv][1])
   **Why probability:** Early success counts are noisy, and always choosing the current leader prevents alternatives from getting evidence.
   **Dogfood angle:** **Partial.** `.can` validates hint eligibility and fallbacks; a separate site service maintains and samples the statistical model. Sampling does not enter the judgment core.
   **Value/effort + data note:** **Value M / effort M.** Requires repeated exercise attempts; defer until traffic supports it. Use session-stable assignments and consented aggregate outcomes, without editor contents or cross-site identifiers.

8. **Name — Mastery Map**
   **Pitch:** Tutorials distinguish “you passed once” from “you probably understand this,” and offer targeted practice.
   **How it works:** Bayesian knowledge tracing maintains per-concept mastery estimates from authored exercise outcomes, accounting for guesses and slips. Concepts include complete exchanges, exhaustive matches, effects, and exact numerics; the tutorial selects practice from a fixed bank. This is a standard knowledge-tracing use case. ([Springer][2])
   **Why probability:** A single success or failure is weak evidence of understanding; a binary completed/not-completed flag hides that uncertainty.
   **Dogfood angle:** **Partial.** Existing `.can` exercises produce deterministic outcomes; the site hosts the knowledge-tracing model. Jev may classify a free-text explanation, but cannot overrule an exercise result.
   **Value/effort + data note:** **Value M / effort M.** Start with explicit provisional priors, then validate against consented attempts. Keep progress locally by default; account synchronization and research use are separate opt-ins.

9. **Name — The Missing Manual**
   **Pitch:** Maintainers see the recurring questions the website fails to answer, grouped by underlying need.
   **How it works:** A dashboard uses Bayesian topic modeling over consented feedback and unsuccessful-search text, then ranks recurring themes with weak coverage in the existing docs. Inspect representative messages before creating a documentation task. Topic mixtures allow one request to concern several concepts. 
   **Why probability:** Different phrasings can describe one missing explanation, while identical words can describe different problems.
   **Dogfood angle:** **Partial.** Topic fitting lives outside can-lang; `score` evaluates whether retrieved documentation addresses each proposed theme. “No good candidate found” is a review suggestion, not proof of absence.
   **Value/effort + data note:** **Value M / effort M.** Needs enough substantive feedback; use manual grouping initially. Explicitly opt in text contributions, redact identifying material, and publish summaries rather than individual queries.

10. **Name — Changelog for My Use Case**
    **Pitch:** “I use typed HTML and exact decimals” produces a focused update digest instead of an amendment-number wall.
    **How it works:** On the changelog, visitors choose topics and a last-read revision. Deterministic metadata selects eligible published changes; semantic relevance scores rank the remaining cards. Preserve direct links, status labels, and the complete chronological view.
    **Why probability:** A change may matter to an application without sharing its feature name—for example, a byte-export change affecting HTML output.
    **Dogfood angle:** **Yes.** `score` evaluates impact on the selected use case; `.can` combines relevance with explicit status rules. Low-confidence relevance appears as “possibly relevant,” not a compatibility warning.
    **Value/effort + data note:** **Value M / effort S.** Uses public amendments and explicit visitor preferences; no behavioral profile needed. Store preferences locally only when requested. Fixed-topic scores can be precomputed outside builds.

11. **Name — Issue Concierge**
    **Pitch:** Before opening a GitHub issue, visitors get a likely destination, related reports, and the missing information that would make their report actionable.
    **How it works:** A site intake form classifies a draft as usage question, documentation problem, site problem, or suspected compiler defect. Retrieval finds existing reports; pairwise relevance scoring suggests duplicates without preventing submission. Prepared checklists request missing reproduction details.
    **Why probability:** Reports describe symptoms in inconsistent language, and superficially similar reports may have different causes.
    **Dogfood angle:** **Yes.** `choice` routes the draft; `score` measures reproduction completeness and candidate similarity. Uncertainty keeps several destinations visible. Nothing is posted or closed automatically.
    **Value/effort + data note:** **Value M / effort M.** Needs public issue content and reviewed example reports. Drafts remain private until the author submits them; show external-processing disclosure and a secret-removal preview.

12. **Name — Journey Breakage Radar**
    **Pitch:** Maintainers learn that a deploy probably broke playground startup or a documentation route before support reports accumulate.
    **How it works:** Apply Bayesian online changepoint detection to volume-conditioned rates of playground startup failures, broken example loads, and documentation 404s. Annotate deployments and alert on sustained posterior evidence of a change; never roll back automatically. ([LIPS][3])
    **Why probability:** A fixed “five failures” alert means different things at ten visits and ten thousand visits; normal variation needs context.
    **Dogfood angle:** **Partial.** Statistical inference runs in the monitoring service; `.can` can implement explicit alert/escalation policy over its outputs. Jev is unnecessary for detecting the change.
    **Value/effort + data note:** **Value M / effort M.** Needs operational event counts and a baseline. Aggregate first-party failures without storing source code, full URLs with query strings, persistent visitor IDs, or raw IP addresses.

13. **Name — First Contribution Match**
    **Pitch:** A visitor selects their skills and available effort, then receives a few suitable website or documentation tasks.
    **How it works:** On the community page, filter a maintainer-curated task catalog by explicit prerequisites, then score fit against the visitor’s stated interests. Offer tasks such as updating an amendment annotation, improving an example explanation, or reproducing a broken site interaction.
    **Why probability:** “I know TypeScript but not language theory” requires softer matching than a `good-first-issue` label.
    **Dogfood angle:** **Yes.** `score` judges task fit and prerequisite clarity; uncertain matches ask a question. The full task list remains accessible, so ranking does not become gatekeeping.
    **Value/effort + data note:** **Value M / effort S.** Needs curated tasks, not contributor histories. Use session-only self-described preferences; do not infer skill from GitHub identity or rank people for maintainers.

14. **Name — Counterexample Arena — MOONSHOT**
    **Pitch:** Visitors challenge their understanding of can-lang and receive a small executable case that exposes the mistaken assumption.
    **How it works:** A playground game maps an explanation to a fixed set of possible misconceptions, then uses Bayesian experimental design to select the challenge expected to distinguish them best. Visitors predict results; the pinned compiler reveals the actual result. Start with claims about coverage, scripts, confidence, and uncertainty.
    **Why probability:** The system is uncertain about the learner’s mental model; choosing the most informative challenge beats presenting the same quiz sequence to everyone.
    **Dogfood angle:** **Yes, with host assistance.** `choice` classifies explanations; `.can` examples and exchange scripts supply executable receipts. The host selects experiments. Out-of-bank claims remain unassessed.
    **Value/effort + data note:** **Value M / effort L.** Needs a carefully authored challenge bank and validated response likelihoods. Default to local session progress; sharing a challenge is explicit and excludes private code.

15. **Name — Synthetic Visitor Lab — MOONSHOT**
    **Pitch:** The website files reproducible draft UX reports against itself, complete with the steps where a simulated visitor got lost.
    **How it works:** A bounded headless-browser runner attempts public tasks such as finding the current `given` rule or launching a retry example. Jev scores a fixed shortlist of visible actions; repeated trials vary starting page and task wording. DOM assertions determine success, while posterior failure estimates prioritize recorded traces for maintainers.
    **Why probability:** Navigation failures can depend on wording and starting context; one deterministic happy-path script misses that variation.
    **Dogfood angle:** **Partial.** `.can` judgments select actions and route uncertainty; the host controls the browser and hard step budget. Synthetic success rates are never presented as real visitor success rates.
    **Value/effort + data note:** **Value M / effort L.** Uses synthetic tasks on staging/public pages, no visitor tracking. Block purchases, posting, credentials, and destructive actions; review drafts before creating public issues.

16. **Name — Calibration Commons — MOONSHOT**
    **Pitch:** Visitors help evaluate the site’s judgment-driven features, and a public dashboard shows where those judgments are dependable or uncertain.
    **How it works:** An annotation page asks concrete questions such as “Does this search result answer this query?” Active learning selects uncertain and diverse cases for review; a separate representative holdout measures performance. Maintainers adjudicate disputed labels before producing reliability plots and Brier scores.
    **Why probability:** It tests whether numerical estimates correspond to observed outcomes, while separating ambiguity from confident mistakes.
    **Dogfood angle:** **Yes for judgment and routing; partial overall.** `noul` predicts relevance, `choice` categorizes disagreement, and `score` grades answer usefulness. Statistical evaluation stays in the site backend; neither crowd votes nor calibration results become proof premises.
    **Value/effort + data note:** **Value M / effort L.** Begin with public, synthetic queries and independent labels. Participation is opt-in; store minimal pseudonymous contribution data, rate-limit submissions, and keep active-learning cases separate from evaluation.

## Top 5 by value/effort

Ranked for an early-stage site: immediate usefulness, little dependence on traffic history, reuse of repository assets, and manageable integration.

| Rank  | Idea                        | One-line justification                                                                                                               |
| ----- | --------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| **1** | **#3 Current-Rule Lens**    | Addresses a concrete documentation problem using public text and a small, reviewable workflow.                                       |
| **2** | **#5 Receipt Sandbox**      | Makes the central product distinction tangible: one observation, explicit uncertainty, deterministic behavior, inspectable evidence. |
| **3** | **#4 Find My Example**      | A small fixed catalog turns visitor intent into an actionable starting point without code generation.                                |
| **4** | **#1 Rule-Aware Search**    | Helps across the entire amendment archive and can start with labeled queries rather than behavioral data.                            |
| **5** | **#6 Diagnostic Companion** | Helps at the moment visitors are blocked, while preserving the compiler’s authority over correctness.                                |

## Start here

Prototype **Current-Rule Lens** first: review a small set of potentially stale passages and publish only approved historical/current links.
Prototype **Receipt Sandbox** alongside it: one mixed-primitive example, curated receipts, visible branch selection, and recorded mode by default.
Together they demonstrate two different benefits: maintaining trustworthy documentation and making probabilistic control flow understandable.
Measure reviewed-link usefulness and whether visitors can correctly predict the sandbox’s low-confidence, tie, and failure behavior.
Defer bandits, topic discovery, and journey models until the site has enough consented data to evaluate them against simpler baselines.

[1]: https://arxiv.org/abs/1707.02038?utm_source=chatgpt.com "A Tutorial on Thompson Sampling"
[2]: https://link.springer.com/article/10.1007/BF01099821 "Knowledge tracing: Modeling the acquisition of procedural knowledge | User Modeling and User-Adapted Interaction | Springer Nature Link"
[3]: https://lips.cs.princeton.edu/bibliography/adams2007changepoint/?utm_source=chatgpt.com "Bayesian Online Changepoint Detection – Laboratory for Intelligent Probabilistic Systems – Princeton University Department of Computer Science"
