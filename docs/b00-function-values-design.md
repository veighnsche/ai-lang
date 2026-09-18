# B00 — Function values and explicit captured values

Status: **proposed bounded design; migration premise contradicted by the inspected snapshot**.
Design only: no compiler, standard-library source, generated output, or tests changed.
Date: 2026-09-19.
Evidence revision: `veighnsche/can-lang@de2bafd3a7a4c858e3a8cf3fee53f650f77876c6`.
Required inputs: a93, a98, a99, and the Functions / outcome-combinator portions of `ASTRA_STDLIB.md`.
Verification status: **source-body inspection and static counterexample review only. No compiler gates were run.** The prompt's green-suite statement is not an independently established result here.

## 0. Evidence boundary and corrected premises

The supplied `veighnsche/ai-lang` repository resolves to `veighnsche/can-lang`. The revision above was selected once and used for the source reads in this document. Direct `git clone` failed because the shell could not resolve GitHub. Repository contents were subsequently read through the connected GitHub reader, not a local checkout. This document is delivered at the requested repository-relative path with a single-file patch; neither a successful clone nor a remote commit is claimed.

References below are **source-file line ranges at that revision**, not search-result line numbers. Function bodies and the cited test bodies were opened. A “checked falsification” means the candidate rule was traced against that inspected source; it does not mean a proposed feature was compiled or a test executed. Verdicts are recommendations for this design, not a claim of maintainer approval.

### Q0a — Which requirements survive source inspection?

**Verdict: the repository wins; retain the higher-order capability request, but reject the claimed migration inventory as established evidence.**

| Prompt / A99 premise | What the inspected snapshot establishes | Consequence |
| --- | --- | --- |
| 29 modules under `std/` | The pinned `std/` directory response lists eight directories: `ascii`, `division`, `html`, `quota`, `ratio`, `scalars`, `schema`, `text`. | Do not claim a 29-module source audit. This counts directories, not all `.can` modules. The directory inventory is GitHub tree evidence, not a compiler inference. |
| Division has `FloorQuotRem`, `TruncQuotRem`, `IntQuot`, `RatioQuot`, `DecQuot`, etc. | The complete `std/division/division.can` instead declares `Int__DivMod`, `Int__Value`, and `Bool__Value`, and five integer operations (`std/division/division.can:9–83`). | No rounding-family N→1 migration can honestly be demonstrated from that file. Section 6 supplies a real, explicitly different pilot. |
| Thirteen codec directories each contain a distinct `Value` union | None of those thirteen directory names appears in the pinned `std/` listing. A direct read of `std/json/json.can` returned not found. | A99's inventory is a requirements claim, not verified implementation evidence. Section 7 frames the architecture but does not invent union definitions. This says nothing exhaustive about sketches elsewhere. |
| An outcomes module has 464 wrappers to collapse | `outcomes` is also absent from that `std/` directory listing. | Leave the requested interface for a follow-up; do not certify the count or a migration. |
| a98 says self-mention is always rejected | The current template checker rejects direct self-fields, while the committed Seq-self test permits `Seq<M__Tree<T>>` (`compiler/expand.go:1460–1490`; `compiler/generics_types_test.go:212–248`). | Preserve the code/test distinction, not the broader prose in a98. |

**Checked falsification.** The complete division source is the counterexample to “function values are necessary for its existing result-bundle deduplication.” Its two one-field wrappers differ only in their field type (`std/division/division.can:21–27`); G2's existing construction and function/type interaction tests already exercise precisely that specialization mechanism (`compiler/generics_types_test.go:101–143, 342–381`). The design survives by crediting G2, not function values, for that collapse. The missing rounding and codec evidence is retained in **Unresolved**, not silently replaced with invented implementations.

### Q0b — What does B00 actually authorize?

**Verdict: specify one small function-value core; do not authorize a std-wide migration.**

The core is an explicitly typed, unary, immutable callable value, created from a named Can function by binding named value parameters. It can be stored in records, passed to ordinary functions, and invoked with explicit outcome dispatch. It is pure with respect to external observations and state, but it may return a declared error. There are no anonymous source bodies, inferred captures, dynamically supplied host callbacks, or changes to bare-return `select`.

The catalogue spells `and_then` with `Fn<T,U,E2>` plus an explicit combined error set (`docs/ASTRA_STDLIB.md:190–207`). That supports a function-value interface. It does not prove that the alleged division and codec migrations require it.

**Checked falsification.** `declaredOkShape` obtains success fields from `recs[fn.Ret]` and errors if that record is absent (`compiler/emit.go:285–308`). Treating the catalogue's `U` as a new bare-return ABI would immediately contradict that body. B00 keeps a named record success carrier. No existing `select` rejection needs to change.

## 1. Syntax and type contract

### Q1a — One callable type, with concrete errors and no hidden capabilities

**Verdict: introduce the built-in type head `Fn<A, R, [e1, e2]>`.**

`A` is the one invocation input type. `R` is the **declared success record**, not an `Outcome` value and not an implicit scalar box. The third slot is a literal list of declared error kinds. `[]` means no declared errors. The callable's capability contract is fixed to **no state and no foreign observations** in B00. Thus there is no second `PureFn` spelling: `Fn<A,R,[]>` is the error-free case of the same pure callable mechanism.

The list must contain known, distinct error kinds in canonical lexical order. Compare error sets exactly when checking callable types; do not infer a set, widen one implicitly, or introduce an untyped error. The target's `emits` list remains an upper bound, not a claim that every listed error actually occurs: the inspected `checkEmits` body checks actual raises against declarations and checks declared kinds exist, but does not demand that each declaration be raised (`compiler/check.go:2089–2120`). Canonical ordering of this new type syntax does not impose a new ordering rule on existing `emits` declarations.

Assignment requires the same input type, same nominal success-record identity, same concrete error set, and the same empty capability contract. There is no variance or brand conversion. A primitive or nominal value may be `A`; a primitive cannot be `R`. `Fn<int, Int__DivMod, [math.zero_divisor]>` is a valid proposed type. `Fn<int, int, []>`, bare `Fn`, and an error-set variable such as `Fn<T,R,E>` are not new admitted types. A bare identifier that is an existing in-scope type parameter is still interpreted under the existing type-parameter rules; recognition of this built-in is specifically the applied `Fn<...>` head.

After expansion, `R` may be a G2 record stamp. Before expansion, direct signature positions such as `Fn<T, Handler__Result<T>, []>` are allowed. This is a built-in signature containing a record mention, **not** permission for arbitrary nested user-generic arguments.

**Current hooks.** `FnDecl` already separates `Params`, `Ret`, `Emits`, and `Effects` (`compiler/parse.go:264–298`). `knownType` and `callee` currently recognize ordinary values and named function/extern signatures, not callable values (`compiler/types.go:162–216`). `tsTypeB` likewise has no callable mapping, and the result-union machinery already represents declared record success plus concrete errors (`compiler/emit.go:79–112, 167–221, 256–308`). These are real checker/emitter extensions, unlike G2's record-only desugaring.

**Checked falsification.** `TestGenericTypeStampsAndRouting` checks nominal stamped return names and distinct field types (`compiler/generics_types_test.go:63–99`). Conflating two same-shaped result records inside `Fn`, or allowing `Fn<int,int,[]>` by silently boxing, would invalidate that exact-identity discipline and `declaredOkShape`. Exact identities and explicit record carriers preserve it. Existing declarations without an applied `Fn` head keep their old paths.

### Q1b — Function references with explicit captured values, not source lambdas

**Verdict: one creation expression: `fnref target<TypeArgs>(named_captures)`.**

Examples of the proposed syntax:

```can
fnref std__int__divmod(divisor = divisor)
fnref handler__identity<int>()
```

Parentheses contain **captures**, not invocation arguments. Every capture is named, in the target's parameter-declaration order, and each name must identify a different declared value parameter. Exactly one target parameter must remain unbound. Its type becomes `A`; the target's declared return record becomes `R`; its complete declared error set becomes the Fn error set. Missing generic arguments are not inferred. A generic target must be specialized explicitly before it becomes a value.

The empty capture list denotes an ordinary reference to an already-unary function. Binding zero parameters of a binary function is rejected; binding every parameter is also rejected. Multiple logical invocation inputs must be placed in one explicit input record by an ordinary source wrapper. This keeps the `Fn<A,R,E>` interface unary without adding tuple types or arity overloads.

Creation evaluates the captures once, in their written order, and stores their immutable values. It does **not** execute the target. A captured divisor of zero is therefore legal to construct; invoking the resulting division callback produces the ordinary `math.zero_divisor` outcome. Runtime closure identity never depends on the address of a mutable caller environment.

For the bounded first cut, capture values are literals, variable/field references, and already-supported data constructions recursively containing those forms, including authorized brand construction under its existing rules. Calls, invocations, indexing/slicing, and computed operators inside a capture are rejected in B00. Compute such a value with an ordinary checked function first and capture its bound result field. This deliberately avoids introducing a new proof rule for a potentially trapping capture expression. No callable-bearing value may itself be captured in this first cut.

**Current hooks.** Add a distinct `Small.Kind = "fnref"`; reuse `Fname`, `TypeArgs`, and `Args` for its target and named bindings (`compiler/parse.go:18–110`). Extend the early head parsing alongside `parseCallHead`, before comparisons (`compiler/parse.go:922–964`). Keeping `target<T>(...)` adjacent reuses the existing generic-head delimiter shape (`compiler/parse.go:1465–1559`); a new `target<T> capture (...)` spelling would unnecessarily defeat that shape. Capture expressions remain AST children, so existing deep-copy machinery can own them per stamp (`compiler/expand.go:815–1020`).

**Checked falsification.** Parsing a reference as a `call` would give it the wrong execution and lint semantics: `parseCallHead` constructs `Kind: "call"`, while `lintRedundantNames` checks names precisely on call nodes (`compiler/parse.go:922–947`; `compiler/lint.go:217–375`). A distinct kind avoids both errors. `TestGenericTypeConstruction` already requires rewriting and clearing a construction's type arguments, and the same post-expansion invariant will be required of references (`compiler/generics_types_test.go:101–143`). The existing division zero-divisor row is the behavioral counterexample to eager reference creation (`std/division/division.can:29–40`): bind-zero must succeed; invoke-zero must fail with that payload.

### Q1c — Invocation is an outcome operation, not an expression call alias

**Verdict: the only invocation form is a singleton `match invoke`.**

```can
match invoke bound.operation(input = dividend)
  on Ok result => Ok(result.quotient, result.remainder)
  on math.zero_divisor failure => math.zero_divisor(failure.divisor)
```

The callee is a variable or record-field path of Fn type. There is exactly one named argument, `input`, whose value has type `A`. The callable and the input are evaluated once, in that order; only then is the target executed with the captured bindings plus the remaining argument. Invocation preserves the target's outcome and error payload identity. There is no new “callback failed” wrapper and no implicit propagation or suppression.

The success binder has record type `R`. Error binders have their declared payload shapes. All possible declared outcomes must be dispatched using the existing exhaustiveness discipline, and each committed caller row must actually witness the arm it is used to cover. A caller may explicitly translate an error into another declared error as it can in an ordinary outcome match; invocation itself never performs that translation.

B00 rejects invocation in value expressions, argument lists, contract predicates, multi-scrutinee matches, chains, `forward`, and `forward call`. It also rejects a `given` table on an invocation. Supporting a single new operation does not silently enlarge those other surfaces.

**Current hooks.** `parseMatchArms` currently chooses `MatchCall` only for a singleton `Small.Kind == "call"`, and chooses its outcome-pattern parser on that basis (`compiler/parse.go:2519–2592`). Add a distinct `MatchInvoke`, with shared outcome-pattern and outcome-domain helpers rather than a copy of the proof algorithm. Keep `MatchCall`'s named-callee invariant. Give `MatchInvoke` an explicitly resolved signature instead of inserting a fake extern or function name into `Program.Fns`. The existing `calleeSig` only stores params and return type, so the new common outcome view must also carry concrete errors (`compiler/types.go:198–216`).

**Checked falsification.** Merely parsing `invoke` as a value match selects the wrong pattern parser in `parseMatchArms`; pretending it is a named call gives graph and effect checks an empty or fabricated `Fname` (`compiler/check.go:1309–1410, 1870–1940`). A distinct node plus an explicit shared signature view prevents both shortcuts. Existing value-product, call, and chain nodes retain their meanings. Division's positive and zero rows supply a concrete two-outcome acceptance vector (`std/division/division.can:29–40`).

### Q1d — Where may callable values live?

**Verdict: admit function parameters and record fields; bound the first cut transitively.**

| Position | B00 rule |
| --- | --- |
| Ordinary function parameter | Fn itself or a record containing Fn is allowed. Such a higher-order consumer is not automatically eligible as a referenced target. |
| Record field / nested record field | Allowed, including fields introduced by G2 record specialization. |
| Function return | Only through an explicitly declared record; never bare `-> Fn<...>`. |
| Invocation input `A`, target success record `R`, or a capture | Must contain no Fn, even through nested records or sequences. This excludes higher-order returned callbacks and captured callbacks from B00. |
| Sequence element, variant payload, error payload, constant, state cell, extern parameter/return | Fn-bearing values are rejected transitively. Existing non-callable cases keep their existing rules. |
| Direct generic argument `Box<Fn<...>>` or `f<Fn<...>>` | Not admitted. B00 does not broaden G1/G2's closed-argument grammar. A plain named record may still travel through the existing generic machinery; its resulting uses remain subject to the containing-position rules above. |
| Equality, ordering, pattern inspection, serialization, capture reflection | No callable operation is admitted, including equality on an enclosing record containing a callable. Passing the value and invoking it are not reflection. |

Use a finite declaration-graph walk for “contains Fn,” with a visited set so Seq-wrapped recursive data types do not cause the checker itself to recurse forever. This is a containment check, not a new record-cycle rule.

**Checked falsification.** `knownType` recognizes sequences recursively, and the G2 Seq-self fixture is accepted (`compiler/types.go:162–197`; `compiler/generics_types_test.go:212–248`). Checking only a top-level spelling would allow a callable through `Seq<Wrapper>` or loop on the existing self fixture. Recursive containment with a visited set closes the former without changing the latter. `checkStateDecl` already restricts cells to four scalar types (`compiler/check.go:1943–1977`); this design does not create a closure-in-state exception.

## 2. Parser and two-pass expansion

### Q2a — Type syntax needs a type-aware splitter, not broader expression guesses

**Verdict: parse the Fn signature structurally at type positions, while retaining string annotations at the existing phase boundaries.**

A small shared type-shape reader must recognize a named type, Seq, a permitted user-generic mention, and the applied Fn head with two type children and one concrete error-list child. Retain canonical strings in the current declaration fields; use the structural view for validation, substitution, contained-mention collection, and rendering back to a string. This is not a second generic expander.

In particular, update `parseFields` to split field/parameter declarations with type-aware angle and square-bracket nesting after `:`. Do **not** globally make every expression `<` start a type argument. The existing `splitTopInner` only treats particular generic heads as angle groups, and `parseFields` currently delegates to it; a multi-comma Fn annotation requires deliberate handling (`compiler/parse.go:440–525, 1736–1784`). `splitTypeArgs` itself tracks angle depth, not the distinct error-list grammar (`compiler/parse.go:1800–1846`). Error-list entries must never be treated as value-type parameters.

Function-reference heads reuse the call head's explicit name/type-argument/parenthesis shape, but have their own kind and binding validation. Malformed heads must report malformed references, not fall through to comparison parsing. No lookbehind expression is needed in the parser or TextMate grammar.

**Checked falsification.** Trace `fn consumer(cb: Fn<int, Int__DivMod, [math.zero_divisor]>, n: int) ...` through the inspected `parseFields` / `splitTopInner`: treating its inner commas as parameter separators breaks it before typing. Conversely, treating all expression angle brackets as delimiters would change the comparison splitters. A type-position splitter plus the already-admitted `target<T>(...)` reference head avoids both. Preserve the existing malformed-arity / nested-instantiation rejection fixtures (`compiler/generics_types_test.go:404–514`).

### Q2b — Reuse the existing function pass, then the existing record pass

**Verdict: extend `expandGenerics`; do not lower closures through a separate monomorphizer.**

The current entry point validates templates, collects/stamps functions, rewrites their headers/calls, then collects/stamps record types and rewrites their headers/mentions (`compiler/expand.go:1–240`). Function stamps and type stamps both splice by template identity and record their `GenericBase` mapping (`compiler/expand.go:1030–1178, 1692–1750`). Preserve that sequence and ownership.

The required extensions, in order, are:

1. **Collect function references as specialization demands.** Extend the existing function-instance discovery to see `fnref` nodes in executable bodies and rows, including references nested inside record constructions. A reference specializes its target even if it is never invoked. It is an address dependency, not an execution edge. Complete explicit row binds still select target function instances; taking a reference does not invent a row for that instance.
2. **Clone and substitute as today.** Captures live in `Args`; target type arguments stay in `TypeArgs`. `cloneSmall`, `substSmall`, `substNode`, and their row/contract walks must visit those children exactly once and keep separate mutable AST trees for each instance (`compiler/expand.go:815–1020`). No capture value enters the specialization key.
3. **Rewrite reference targets with ordinary function stamps.** Generalize the appropriate kind filter in `rewriteGenericCalls` to include references for renaming only; do not classify references as invocations in call graphs. Its current renaming and complete child walk are visible at `compiler/expand.go:1186–1256`. A rewritten reference has a stamped target and no residual type arguments.
4. **Descend into Fn's input and success types in the record pass.** Extend `eachTypeSlot`'s type-shape traversal, `typeMentionsParam`, `collectTypeInstances`, `substType`, and `rewriteTypeMentions`. Retain the enclosing type-template pointer when visiting a nested signature. `Fn` is built-in, not a template to stamp; its error-list child is not a type slot (`compiler/expand.go:815–860, 1390–1610, 1795–1871`).
5. **Preserve record-instance identity and reject-unmentioned behavior.** Signature mentions contribute to the same memoized `(base,args)` collection, so `Handler__Ops<int>` can reveal `Handler__Result<int>`. No signature mention is allowed to escape rewriting merely because it sits inside a callable field. No template is silently dropped for lack of mentions.
6. **Resolve callable sites after world construction and before either proof or execution consumes them.** Supply each invocation with a checked monomorphic input/return/error signature and a completed static target-admission certificate. Error, exhaustiveness, effect, cycle, and emitter consumers must fail closed if it is absent. Q2d specifies the shared barrier; the complete internal walker/refactoring audit remains in Unresolved.

A generic record field illustrates the intended result:

```text
Source field in Handler__Ops<T>:
  transform: Fn<T, Handler__Result<T>, []>

After the existing two passes for the int instance:
  record Handler__Ops$T$int
  field transform: Fn<int, Handler__Result$T$int, []>

A reference fnref handler__identity<int>() names:
  handler__identity$T$int
```

These `$` spellings are explanatory internal names, not new legal `.can` source. Two closures using the same target/type arguments but different captures share that target stamp and have different runtime capture values. Source-template deduplication is not a claim of one emitted implementation for all element types.

**Checked falsifications.** `TestGenericTypeFnInteraction` requires function stamping to expose a record instance; `TestGenericTypePassThrough` requires a record-to-record fixpoint; `TestGenericTwoTypeTemplatesOneModule` catches stale-index replacement; `TestGenericTypeRejects` checks unmentioned templates and disallowed nested arguments (`compiler/generics_types_test.go:163–210, 342–514`). Fn traversal must enter those same collectors, not bypass them. `TestGenericTypeCrossModule` also asserts `GenericBase` and optional type-pin rewriting (`compiler/generics_types_test.go:250–314`). No new source helper declarations or alternate stamp naming scheme are introduced, so the existing `Domain__Name` exemption remains tied to `GenericBase`, not a broad allowance for `$` names.

### Q2c — Keep the precise late chain-else restriction

**Verdict: Fn syntax must not hide in late-parsed chain text.**

B00 does not add invocation or reference creation to chain-else text. Add a keyword-head-specific rejection for those new forms, and retain the existing parameter/generic-head test for old syntax. Do not restore an imprecise “contains `call` and `<`” heuristic, and do not use this feature to un-revert `select`.

**Checked falsification.** `substNode` deliberately does not transform `ChainElse`, because it remains source text (`compiler/expand.go:889–919`). The committed `TestGenericTypeChainElseGate` constructs a monomorphic function containing a generic construction in that text and expects the expansion diagnostic (`compiler/generics_types_test.go:516–540`). Rejecting new callable heads at the same boundary preserves the reason for that gate; accepting them without AST rewriting would strand a source reference or type parameter.

### Q2d — A shared static barrier must precede both proof and table execution

**Verdict: prepare and certify callable sites once per checked program, after world/authority construction and before callback-dependent proof or any table execution. Do not rely on the current per-function type-check location.**

The inspected entry points have materially different ordering. CLI `checkProgram` invokes `checkSem` for modules before its explicit `verifyExhaustiveAll` loop (`compiler/main.go:438–503`). LSP `diagnoseWith` invokes exhaustiveness before `checkSem` (`compiler/lsp.go:428–449`). Inside `checkSem`, call/effect/given checks precede `checkTypes`, which mutates constructor bindings; table execution then follows a selective execution-blocking check (`compiler/lsp.go:521–600`). Adding Fn only to `checkTypes` would leave the editor's earlier proof path without a resolved callable signature. Checking just the currently open module also cannot authorize real execution of an unprepared provider body.

Factor the necessary static preparation from execution so both entry points use one callable-preparation barrier. It must resolve monomorphic signature components, capture slots and invocation input types; complete ordinary type/authority preparation for the reachable source graph; collect address-taken targets; and enforce Q3's purity, precondition, and indirect-cycle rules. It supplies typed invocation outcomes to the existing proof machinery. It must not synthesize foreign declarations or recursively invoke `checkProgram` while evaluating a row.

Any callable-preparation error blocks callback execution regardless of whether its diagnostic appears in the old `checkSem` blocking whitelist. A missing certificate is a hard refusal in evaluator and emitter, never an invitation to continue with a guessed signature. A partial editor world may report source diagnostics but cannot execute an uncertified callback. Legacy Fn-free programs retain their established checks and output; changing unrelated diagnostic ordering is not part of this slice.

Distinguish **static eligibility to execute a row** from **final acceptance evidence**. The barrier proves that the callback operation is statically admitted; it does not require every table to have passed already, which would make callback-bearing tables circular to run. Ordinary rows then execute, their failures remain failures, and final acceptance still requires all ordinary checks, rows, coverage obligations, contract checks, and output gates. No callback run supplies missing target rows.

**Checked falsification.** The two inspected entry-point bodies above expose the ordering counterexample, while `checkStatic` explicitly diagnoses rowless function stamps and `checkSem` owns row execution (`compiler/lsp.go:484–600, 600–667`). A late type-only hook or recursive full-check-on-invoke would fail that ordering or introduce circular evaluation. A shared preflight with separate evidence completion avoids both. Exact extraction of the existing checks, full walker coverage, and diagnostic regression tests are release prerequisites, not claimed completed work.

## 3. Identity, authority, purity, and termination

### Q3a — A reference does not manufacture linkage authority

**Verdict: resolve the target as a named source function under the same source-module authority as an ordinary named dependency.**

A foreign reference requires the existing base-name revision pin in `uses`; source syntax does not add a second revision pin on `fnref`. Same-file references retain the target declaration's revision. A generic target uses the existing base-to-instance identity relation. Taking its address counts as use, even if the particular program never invokes it. Merely mentioning the Fn type does not count as taking any particular function's address.

Address dependency analysis must be distinct from execution dependency analysis. The current `calledFns` filters for call nodes, and `checkUnusedUses` consults that set (`compiler/check.go:2123–2156`). `compileAll` also fills `FnDecl.UsesHere` from named call traversal (`compiler/main.go:216–231`). Extend the appropriate use/import collection to include `fnref` targets without turning creation into an execution edge. The full downstream import planner is not audited here and remains a release prerequisite.

A consumer receiving a callable need not name its hidden concrete target in its own `uses`: the creator must have obtained a valid target, and the whole checked program must admit its execution graph. The consumer must not get authority to construct other targets or inspect the capture environment. References to extern declarations or compiler kernels are rejected; an eligible **source function** may still call an already-admitted deterministic kernel under the kernel's existing authorization rules.

**Checked falsification.** `parsePinRev` and `declaredRev` compare a pin with a declaration revision, and the type-stamp cross-module test checks both pin rewriting and `GenericBase` (`compiler/check.go:46–101`; `compiler/generics_types_test.go:250–314`). Inventing a “closure function” with a new source owner or treating all Fn consumers as users of every target would change those identities. Retaining the real target and checking authority at creation avoids both. Type pins remain optional where the current type rules make them optional; function references do not create a back door around function pins.

### Q3b — Pure means the complete reachable execution graph

**Verdict: no externs, state operations, unresolved dependencies, or nonempty effect contracts anywhere in an eligible callback graph.**

A declared error is not an effect. An error-free-looking source function that calls an extern is still ineligible. Inspect every branch, not only those taken by table rows. Closure creation cannot be used to carry a private capability across this rule.

Reuse the admissibility boundary exhibited by `checkLinkedGraph`: source functions and admitted deterministic intrinsics, excluding externs, store operations, unresolved calls, and functions with `Effects` (`compiler/linked.go:27–65`). Extend its graph input for the new invocation edges; its current `walkCalls(fn.Body)` loop is not a higher-order proof. The same completed certificate must govern ordinary table invocation, linked execution, and emitted callbacks. Do not enable execution first and add admission later.

**Checked falsification.** The source body of `checkLinkedGraph` explicitly rejects an extern even if the source wrapper's own effect list is empty. Checking only `len(target.Effects)==0` loses that rejection. B00 reuses its whole-graph criterion instead. All five inspected division functions use ordinary pure arithmetic and contain no named calls, so the selected callback target does not depend on relaxing this boundary (`std/division/division.can:29–83`).

### Q3c — Do not erase preconditions or smuggle a new recursion proof

**Verdict: reject required-precondition callback graphs and reject every cycle containing an indirect invocation edge.**

A referenced target, or a reachable source callee in its eligible graph, with a nonempty `Requires` list is refused in the bounded core. Even captures that appear to satisfy a precondition do not discharge it in B00. A future design may admit a wrapper that proves or eliminates the precondition, but merely putting a required-precondition callee behind a wrapper does not satisfy B00: the whole reachable graph is checked. B00 does not provide that discharge mechanism. No callback contract summaries are imported into callers. Contracted functions containing callable-bearing values or new callable syntax remain outside B00's proof fragment; explicit rejection must precede any claim that such a contract was verified.

For termination, use a conservative closed-world graph. For each monomorphic Fn signature, collect all admissible address-taken target stamps in the checked program, including those mentioned in fixtures. Each invocation contributes potential edges from its owner to every target of that signature. Preserve the existing direct-call edges. Reject any strongly connected component containing an indirect edge, including a possible indirect self-call. Direct self-recursion already accepted by the existing decreases checker remains eligible when it does not participate in such a component.

This intentionally rejects some terminating programs. It avoids a new higher-order termination theorem, capture-sensitive recursion inference, or specialization by runtime closure identity. Certificates are scoped to the checked program: linking additional separately checked modules requires recomputing the target sets and graph. Arbitrary JavaScript callbacks cannot be injected under a certificate for a different program.

**Checked falsifications.** `FnDecl` retains `Requires`, and `isContracted` tests both Requires and Ensures; the inspected admission classifier maps unsupported types outside its fragment (`compiler/parse.go:286–298`; `compiler/verify_admission.go:94–116, 154–176`). Signature compatibility alone does not preserve those conditions. The direct-call cycle collector keys edges by `Fname`, while `checkDecreases` recognizes specific syntactic self-calls and guards (`compiler/check.go:1309–1432, 1463–1586`). A recursive cycle hidden behind Fn would evade the former and has no proof under the latter. Explicit indirect edges and conservative rejection close that gap without weakening the accepted direct-recursion schemas. Verification is active in the inspected CLI and LSP bodies, notwithstanding older comments in the admission file (`compiler/main.go:179–196`; `compiler/lsp.go:461–467`). B00 must reject its new out-of-fragment constructs explicitly; the complete verifier walker integration remains in Unresolved.

## 4. Runtime representation and evidence

### Q4a — Store a snapshot, invoke the original checked body

**Verdict: use an opaque runtime callable; never a user-constructible environment record.**

The Go value needs a distinct callable kind with a resolved target stamp, checked signature, remaining parameter position, and a fresh binding map containing the captured values. Capture the values, not a live mutable caller-environment map. Immutable subvalues may share storage where existing ownership permits; this does not authorize a new mutable alias. Source-level field access on the callable itself must fail.

The proposed TypeScript representation is a compiler-generated unary closure thunk: evaluate each capture into a fresh local once, then invoke the existing target function with parameters reconstructed in declaration order. The thunk returns the target's existing flattened Ok/error union unchanged. It is not a new source lambda, does not add generated source decision branches, and does not become an extra independently stamped Can function. `tsTypeB` needs a callable type mapping that uses the same record/error union builders as named functions.

B00's authority claim is about checked Can programs and their trusted generated-code boundary, not hostile JavaScript using casts or mutating internal objects. Fn-bearing extern signatures and arbitrary host callback inputs are not admitted. Captured values must not be exposed by debug descriptions, serialization, or error messages; a diagnostic may name the target and mismatching capture slot, but must not dump the environment.

**Checked falsification.** `Value` has no callable carrier today, and `vField` admits only record/outcome/variant dictionaries (`compiler/eval.go:15–45, 346–365`). Reusing `Kind: "rec"` for the capture environment would expose it through existing projection. A distinct kind avoids that. The emitted outcome shape comes from declared fields and errors, not observations (`compiler/emit.go:167–221, 256–308`); returning a second nested outcome or a scalar from the thunk would contradict that ABI. The separate bind-three / bind-five parity vector in section 8 is required to catch a live-environment capture implementation.

### Q4b — Invocation executes real pure bodies in caller tests

**Verdict: invocation is explicitly linked-pure execution, with isolated callee coverage, even in an ordinary caller table row.**

The caller row supplies a checked function reference or a record containing one. Invocation executes that target's real body and its admitted pure callees. It does not consume a `given` exchange. This is an explicit new operation-specific evidence rule; ordinary named foreign `call` operations keep their existing scripted unit semantics.

Execute the callback subtree in a child context with linked-pure dispatch, no script consumption, and no target-coverage credit. Preserve the parent's resource/depth accounting rather than resetting it at every callback. On returning, mark the selected **caller invocation arm** in the parent's coverage context. Do not globally flip the caller into linked mode, and do not let a fresh child context reset a recursion/resource bound.

Every referenced generic target still needs its own routed table rows; taking a reference never manufactures them. Preserve the actual existing branch-evidence rule: passing ordinary rows share coverage across local helper flow within their module, and a checked identity relay of a local named call has a structural-certificate exception (`compiler/lsp.go:600–697, 710–760`). Do not replace this with a stricter “only the function's own rows count” rule.

The new boundary is specifically the indirect callback subtree. It receives `Cov=nil`, so executing a target through `invoke` cannot supply missing provider branch evidence. The selected invocation arm is marked in the current parent context and counts only if that enclosing ordinary row passes; the row may reach the consumer through an ordinary local helper call. Existing direct-helper coverage sharing remains intact. Linked integration runs still contribute no acceptance coverage.

Do not extend the local identity-relay exception to `MatchInvoke` in B00. If no legitimate ordinary coverage flow witnesses a required invocation error arm, compilation fails rather than synthesizing a `given` response or treating `emits` as proof that the error occurs. This limitation matters for targets whose declared error upper bound contains unrealized kinds; the follow-up cannot silently erase it.

**Checked falsifications.** `Ctx` separates Linked, Scripts, Depth, and Cov, and `markTaken` writes only when Cov is non-nil (`compiler/eval.go:261–312`). `runLinkedPure` requires a clean checked program, checks the root revision and graph, and uses discarded integration coverage (`compiler/linked.go:68–177`). Sharing caller Cov with the target or setting caller Linked globally would change those boundaries. The child-context rule makes the new behavior local to invocation and leaves the ordinary call path alone. `checkSem` merges only passing-row coverage and `relayStatus` specifically requires `MatchCall` with a local callee (`compiler/lsp.go:600–760`); retaining those branches preserves existing helper credit and the existing relay exception without granting either to an indirect target. `TestRatioPilotLinked` is an inspected example of real-body evidence through G2 stamps (`compiler/ratio_pilot_test.go:25–51`).

### Q4c — Factory rows need receipt comparison, not function equality

**Verdict: add a test-expectation comparison for callable construction receipts, separate from language equality.**

A factory row may expect `Ok(fnref target(captures))`. Compare the resolved target identity/revision, monomorphic signature, remaining argument position, and captured data values. Recurse through the declared enclosing record. This compares the construction that occurred; it is **not** a proof that two different functions are extensionally equal. Do not invoke a callback while comparing a factory expectation.

Keep language equality and ordering rejected for Fn and any record containing Fn. Use a dedicated expectation comparator rather than adding a generally usable callable equality case to `vEq`. Descriptions of a mismatch must redact captured payloads. Linked factory expectations must use the same receipt rule; invocation parity should compare observable outcomes, not serialize a JavaScript function.

**Checked falsification.** `vEq` recursively compares record fields and already distinguishes test-only sequence/byte comparison from language operators (`compiler/eval.go:367–445`). `emitEquality` separately chooses scalar or record equality by type (`compiler/emit.go:343–413`). A blind callable equality addition on either side could admit function identity as a program operation or disagree across Go and TypeScript. A type-level transitive prohibition plus a private expectation comparator prevents that drift. Factory rows remain concrete evidence instead of dropping required rows merely because their result contains a callable.

## 5. Required rejection matrix

These are acceptance conditions for the proposed feature, not reports of tests already executed. Numeric diagnostic allocations are intentionally not invented: each new rule needs its own audited registry entry, while an existing code is reused only for the same existing rule. Diagnostics must identify the source reference or invocation, not expose a synthetic stamp as the only useful location.

| Input / counterexample | Required rejection or observation | Existing boundary checked |
| --- | --- | --- |
| `Fn<int,int,[]>` | Success must be a named record; no implicit boxing. | `declaredOkShape`, `compiler/emit.go:285–308`. |
| An unknown or duplicate kind in Fn's error list | Reject the malformed error contract. A generic `E` is not an error-kind list. | Errors and types are separate registries; `compiler/types.go:25–45, 179–197`; catalogue `docs/ASTRA_STDLIB.md:190–207`. |
| A callback returning a different same-shaped record or accepting a different brand | Exact type mismatch; no nominal erasure in typing. | Record lookup / callee signatures, `compiler/types.go:179–255`; stamp-identity test, `compiler/generics_types_test.go:63–99`. |
| `fnref std__int__divmod()` | Two unbound parameters; reject. | Actual binary signature, `std/division/division.can:29`. |
| Bind both `dividend` and `divisor` | Zero unbound parameters; reject, not an eager call. | Same signature and body, `std/division/division.can:29–40`. |
| Positional, duplicated, unknown, out-of-order, or mistyped captures | Reject the particular binding rule violated. Mandatory names are not a lint suggestion. | Named-argument AST / duplicate-param handling, `compiler/parse.go:18–30, 1786–1799`; lint call handling, `compiler/lint.go:217–375`. |
| `fnref target(bound = 1 / 0)` | Computed capture expressions are outside B00, even before target execution. | Existing expression parsing distinguishes arithmetic from data, `compiler/parse.go:1120–1464`; division's explicit guard, `std/division/division.can:38–40`. |
| A capture of a record that contains Fn | Reject transitively; no hidden higher-order environment. | Record/sequence type traversal must preserve `TestGenericTypeSeqSelf`, `compiler/generics_types_test.go:212–248`. |
| A target needing `requires`, or whose graph reaches one | Reject; no capture-based discharge and no precondition erasure. | `FnDecl.Requires` / `isContracted`, `compiler/parse.go:286–298`; `compiler/verify_admission.go:94–116`. |
| A nominally pure target with an extern or state call on an untaken branch | Reject before running any caller row. | Whole-body linked graph walk, `compiler/linked.go:27–65`. |
| `given` attached to `match invoke` | Reject; no scripted callback outcome. | Separate context modes, `compiler/eval.go:261–296`; ordinary exchange shape, `compiler/check.go:1980–2020`. |
| A missing invocation outcome arm | Exhaustiveness failure. A fabricated fixture is not a fix. | Existing outcome-pattern path is explicit, `compiler/parse.go:2519–2592`; linked coverage boundary, `compiler/linked.go:68–177`. |
| An invocation cycle, including through another module or a record field | Reject a component containing an indirect edge, even if current rows terminate. | Existing direct-edge cycle graph, `compiler/check.go:1309–1432`. |
| `op == op`, or equality on `Wrapper(op)` | Reject in the language; do not reuse factory receipt equality. | Existing structural equality recursion, `compiler/eval.go:367–445`; type-directed emission, `compiler/emit.go:343–413`. |
| Change captures / target but leave factory expected receipt unchanged | Factory test fails without executing the target or disclosing captured data. | Expectation comparison is currently structural, `compiler/eval.go:367–445`; B00 adds a separately scoped receipt comparison. |
| A generic target instance with no own committed rows | Preserve the ordinary missing-evidence failure; caller invocation is not provider evidence. | Per-instance row routing, `compiler/expand.go:1030–1178`; independent integration precondition, `compiler/linked.go:68–107`. |
| New callable syntax in chain-else text | Reject at its source boundary; do not leave it for unrewritten late parsing. | `compiler/expand.go:889–919`; `compiler/generics_types_test.go:516–540`. |

## 6. Worked division pilot: separate the real collapse from the closure capability

### Q6a — Actual source-declaration collapse: two wrappers to one

**Verdict: use the existing one-field division result family as the bounded N→1 specimen. This is G2-only, not a demonstrated need for Fn.**

The exact existing declarations are (`std/division/division.can:21–27`):

```can
type Int__Value rev 1 (
  value: int
)

type Bool__Value rev 1 (
  value: bool
)
```

The proposed single family is:

```can
type Division__Value<T> rev 1 (
  value: T
)
```

**Count:** two source wrapper declarations become one source template. `Int__DivMod` stays separate because it has two fields with a different contract. The module's total record declarations become three → two source declarations; expansion still produces separate int and bool instances. No algorithm or rounding policy has been collapsed, and no function value is involved in this count.

For a complete signature-and-row specimen, the replacement module is below. The four changed public function signatures are shown at revision 2 as an explicit proposed API change; the unchanged divmod operation stays at revision 1. This is a design specimen, not an applied migration or an executed gate. The accepted revision baseline and downstream pins must be audited before choosing the actual migration revisions.

```can
mod division
  provides [std__int__divmod, std__int__mod, std__int__is_multiple, std__int__is_even, std__int__is_odd, Int__DivMod, Division__Value]
  uses []
  emits [math.zero_divisor]

error math.zero_divisor(divisor: int)

type Int__DivMod rev 1 (
  quotient: int
  remainder: int
)

type Division__Value<T> rev 1 (
  value: T
)

fn std__int__divmod(dividend: int, divisor: int) -> Int__DivMod rev 1
  emits [math.zero_divisor]
  tests
    divmod_pos_pos(7, 3) => Ok(2, 1)
    divmod_neg_pos(-7, 3) => Ok(-3, 2)
    divmod_pos_neg(7, -3) => Ok(-2, 1)
    divmod_neg_neg(-7, -3) => Ok(3, 2)
    divmod_zero_dividend(0, 5) => Ok(0, 0)
    divmod_zero_divisor(5, 0) => math.zero_divisor(0)
  match divisor == 0
    true => math.zero_divisor(divisor)
    false => Ok(dividend / divisor, dividend % divisor)

fn std__int__mod(value: int, modulus: int) -> Division__Value<int> rev 2
  emits [math.zero_divisor]
  tests
    mod_pos(7, 3) => Ok(1)
    mod_neg(-7, 3) => Ok(2)
    mod_neg_modulus(7, -3) => Ok(1)
    mod_zero(0, 5) => Ok(0)
    mod_zero_modulus(5, 0) => math.zero_divisor(0)
  match modulus == 0
    true => math.zero_divisor(modulus)
    false => Ok(value % modulus)

fn std__int__is_multiple(value: int, divisor: int) -> Division__Value<bool> rev 2
  emits [math.zero_divisor]
  tests
    multiple_yes(6, 3) => Ok(true)
    multiple_no(7, 3) => Ok(false)
    multiple_zero(0, 5) => Ok(true)
    multiple_zero_divisor(5, 0) => math.zero_divisor(0)
  match divisor == 0
    true => math.zero_divisor(divisor)
    false => Ok(value % divisor == 0)

fn std__int__is_even(value: int) -> Division__Value<bool> rev 2
  emits []
  tests
    even_pos(4) => Ok(true)
    even_neg(-4) => Ok(true)
    even_odd_pos(3) => Ok(false)
    even_odd_neg(-3) => Ok(false)
    even_zero(0) => Ok(true)
  Ok(value % 2 == 0)

fn std__int__is_odd(value: int) -> Division__Value<bool> rev 2
  emits []
  tests
    odd_pos(3) => Ok(true)
    odd_neg(-3) => Ok(true)
    odd_even_pos(4) => Ok(false)
    odd_even_neg(-4) => Ok(false)
    odd_zero(0) => Ok(false)
  Ok(value % 2 == 1)
```

**Checked falsification.** The original function bodies and all rows are visible at `std/division/division.can:29–83`; none computes with the wrapper's nominal name. The specimen changes only the wrapper declaration, provided type name, affected return annotations, and the proposed revisions. `TestGenericTypeStampsAndRouting`, `TestGenericTypeConstruction`, and the ratio pilot establish the inspected specialization pattern to preserve, not an executed result for this specimen (`compiler/generics_types_test.go:63–143`; `compiler/ratio_pilot_test.go:1–95`). Keeping `Int__DivMod` outside the family avoids lying about its quotient/remainder shape.

### Q6b — Separate Fn pilot: a divisor captured inside a result record

**Verdict: demonstrate the new capability against the real `std__int__divmod`, without claiming an existing rounding-family reduction.**

The following is exact **proposed B00 `.can` syntax**, not source accepted by the pinned compiler today. Load it together with the original division module, or the G2 specimen above whose divmod revision is unchanged.

```can
mod divcallback
  provides [divcallback__bind, divcallback__run, Divcallback__Bound]
  uses [std__int__divmod@1, Int__DivMod@1]
  emits [math.zero_divisor]

type Divcallback__Bound rev 1 (
  operation: Fn<int, Int__DivMod, [math.zero_divisor]>
)

fn divcallback__bind(divisor: int) -> Divcallback__Bound rev 1
  emits []
  tests
    positive(3) => Ok(fnref std__int__divmod(divisor = 3))
    negative(-3) => Ok(fnref std__int__divmod(divisor = -3))
    zero(0) => Ok(fnref std__int__divmod(divisor = 0))
  Ok(fnref std__int__divmod(divisor = divisor))

fn divcallback__run(bound: Divcallback__Bound, dividend: int) -> Int__DivMod rev 1
  emits [math.zero_divisor]
  tests
    positive_positive(Divcallback__Bound(fnref std__int__divmod(divisor = 3)), 7) => Ok(2, 1)
    negative_positive(Divcallback__Bound(fnref std__int__divmod(divisor = 3)), -7) => Ok(-3, 2)
    positive_negative(Divcallback__Bound(fnref std__int__divmod(divisor = -3)), 7) => Ok(-2, 1)
    negative_negative(Divcallback__Bound(fnref std__int__divmod(divisor = -3)), -7) => Ok(3, 2)
    zero_dividend(Divcallback__Bound(fnref std__int__divmod(divisor = 5)), 0) => Ok(0, 0)
    zero_divisor(Divcallback__Bound(fnref std__int__divmod(divisor = 0)), 5) => math.zero_divisor(0)
  match invoke bound.operation(input = dividend)
    on Ok result => Ok(result.quotient, result.remainder)
    on math.zero_divisor failure => math.zero_divisor(failure.divisor)
```

The six invocation expectations are copied from the semantics already exercised by division's committed rows, including all sign combinations and zero (`std/division/division.can:29–40`). Factory rows additionally distinguish creation from execution: `bind(0)` has no error, while invocation with that capture emits `math.zero_divisor(0)`. Positional record construction follows the existing lint convention; capture names and `input` are mandatory grammar rather than redundant labels (`compiler/lint.go:217–375`).

This proves the proposed surface is expressive enough to describe “a callable stored in a record with a runtime value captured.” It is not evidence that the proposed compiler changes pass, that the factory and consumer have been linked in this session, or that nonexistent rounding implementations have collapsed. Add a linked fixture that calls the factory and then the consumer, plus a TypeScript driver that keeps two separately bound closures alive, before accepting the pilot.

### Q6c — Is a closure necessarily the right representation for rounding policy?

**Verdict: necessity is not established. Do not make the language feature's acceptance depend on that assertion.**

First factor identical data carriers with G2. If the missing customer proves to be a finite, closed set of rounding modes, compare an explicit mode value and one checked dispatcher with a callback policy. If customers must supply new pure computations, Fn has a distinct purpose. No representation may collapse genuinely different rounding semantics simply because their result fields happen to match.

**Checked falsification.** The actual two-wrapper family already collapses without Fn, as section 6a shows against the inspected source. That defeats the broad premise that every division duplication is a function-value blocker. The absent rounding algorithms prevent a ruling between mode data and callbacks for those algorithms; that specific decision is retained under Unresolved.

## 7. Value-cluster sharing: options, recommendation, and missing evidence

### Q7a — Do not equate nominal API identity with accidental duplication

The following are **architecture options, not an approved migration design**. The required thirteen codec source unions are missing from the inspected `std/` tree, so there is no source-grounded verdict that their cases, precision rules, key semantics, or errors are identical. A99 calls for a follow-up decision; the catalogue also requires explicit schema mapping and exact numeric wire representation rather than treating every encoding as interchangeable (`docs/ASTRA_STDLIB.md:252–269`).

| Option | What it shares | Benefit | Falsification required before choosing it |
| --- | --- | --- | --- |
| A. One canonical semantic `Value` | A single nominal data model used by every participating codec. | Maximum sharing when models and intended public identity really coincide. | Enumerate every case and payload. Find any value a codec cannot represent, any lossy number conversion, map ordering / duplicate distinction, or changed nominal API boundary. Round-trip and rejection laws must be stated per codec, not assumed. |
| B. Keep local nominal values; share algorithms through explicit adapters | A common algorithm accepts record-shaped inputs and typed Fn operations supplied by each codec. | Separates algorithm reuse from public type unification; can be introduced one verified cluster at a time. | Produce one genuinely identical algorithm with two lossless adapters. Check that adapter errors and exact output records fit the bounded Fn signature without boxing away a distinction or adding an unproved conversion. |
| C. Shared core plus explicit format extensions | A common subset plus revisioned, explicitly represented extension cases. | Allows measured commonality without pretending every format has the same model. | Identify a stable common subset and all extension cases. Reject designs that rely on open-ended variants, generic variants, or unchecked extension payloads not supported by the chosen language slice. |

**Provisional recommendation: investigate B first, A second, C only when a concrete common-core/extension split is demonstrated.** This is a preference for the least irreversible experiment, not a claim that adapters already deduplicate the missing modules. B00 provides the ability to carry a checked operation in an adapter record; it does not supply the absent algorithms, recursive value architecture, generic variants, or conversion proofs.

Before moving this recommendation into a decided Q7 verdict, attach the exact codec files and make a matrix of nominal type owner/revision, cases and payloads, ordering and duplicate behavior, numeric preservation, accepted/rejected values, and operation-specific errors. Pick two codecs with an actual common operation; count the removed source declarations separately from the added adapters. Test both successful round trips and deliberate nonrepresentability. A single lossy “universal value” conversion falsifies an unqualified A recommendation.

Html wrappers remain outside this exercise. No type-name or structural-shape heuristic may absorb them into a generic Value family. That exclusion is a requirement of this task, not a claim that the whole html module was audited here.

## 8. Gate plan and acceptance sequence

### Q8a — Preserve existing witnesses; add feature-specific falsifications

**Verdict: no feature acceptance on parser success, one happy-path callback, or a rewritten golden alone. Each affected gate needs an independently observable failure case.**

All commands and test additions in this section are **planned**, not executed. Existing test names below refer to bodies that were read. Suggested `TestFunctionValues…` names are provisional until the repository-wide name/diagnostic audit in Unresolved is complete.

| Gate | New evidence required | Existing inspected anchor / falsification |
| --- | --- | --- |
| Unit: parsing and types | Parse a Fn parameter and generic-record field; empty and multi-kind error lists; `fnref f<T>(bound = v)`; field-path invocation. Reject malformed angle/list delimiters, primitive success, missing/duplicate/out-of-order captures, two or zero residual parameters, unknown/duplicate error kinds, wrong brands/record identities, missing generic arguments, and every prohibited invocation position. Test `<` comparisons adjacent to callable syntax so the delimiter fix cannot swallow operators. | `parseCallHead`, `parseFields`, and `splitTopInner` are separate parsing paths (`compiler/parse.go:440–525, 922–960, 1736–1784`). `TestGenericTypeRejects` pins existing invalid generic forms (`compiler/generics_types_test.go:404–514`). |
| Unit: expansion | Reference-only specialization must emit the target stamp but still reject a stamp without routed rows. Test two generic functions and two generic record templates in one module, reverse declaration order, nested Fn signature mentions, a type parameter used only inside Fn, shared target/different captures, and deep-clone independence between instances. Preserve old GenericBase routing, optional type pins, unknown/nested argument rejections, and the exact chain-else gate. | `TestGenericTypeStampsAndRouting`, `TestGenericTypeFnInteraction`, `TestGenericTwoTypeTemplatesOneModule`, `TestGenericTypeSeqSelf`, and `TestGenericTypeChainElseGate` check the mechanisms that must not fork (`compiler/generics_types_test.go:63–99, 212–248, 342–540`). |
| Unit: static admission | In both CLI and LSP entry points, reject a hidden extern/state operation on an untaken branch, a transitive required precondition, a missing target revision/provider, and an indirect cycle. A pure target with an existing proved direct self-decrease remains admissible when no indirect cycle exists. Any new static refusal must block invocation, not just add a diagnostic after executing it. | `checkLinkedGraph`, `checkGlobalCycles`, `checkDecreases`, and the two entry-point orderings are the inspected boundaries (`compiler/linked.go:27–65`; `compiler/check.go:1309–1586`; `compiler/main.go:438–503`; `compiler/lsp.go:428–449, 521–600`). |
| Unit: value and expectation semantics | Construct a zero-divisor closure without invoking it. Receipt equality must distinguish target, revision, specialization, residual slot, and capture values without executing the callback. Program equality must still reject Fn and Fn-bearing records. Changing a captured secret in a negative fixture must not print its contents. Exercise transitive forbidden positions through named wrappers, including a Seq-recursive data graph that contains no Fn and must still terminate the containment check. | Distinct `Value` kinds, field projection, recursive `vEq`, and type-directed emitted equality are separate hooks (`compiler/eval.go:15–45, 346–445`; `compiler/emit.go:343–413`). This gate catches accidental environment-record exposure and host-identity comparisons. |
| Unit: rows and coverage | A target with no rows still fails. Removing the last legitimate ordinary witness for a provider branch must fail even when a callback executes that branch. Preserve module-wide ordinary local-helper credit and the existing local identity-relay exception. An invocation branch counts only from a passing parent row; callback-subtree and integration coverage must not leak. Reject `given` on invocation; prove that a named foreign call elsewhere in the parent still consumes its own script normally. | `checkStatic`, `checkSem`, `checkCoverage`, and `relayStatus` establish the actual coverage boundaries, including their exceptions (`compiler/lsp.go:484–513, 600–760`). This is not a new own-rows-only coverage law. |
| Golden emission | Add an isolated callback fixture with its TypeScript and error catalogue. Pin imported target names, unary thunk shape, capture temporaries, original parameter order, exact success fields/errors, and Fn-bearing record types. Fn-free source must not gain runtime helpers/imports or altered outputs. Do not edit unrelated goldens to hide drift. | `TestGoldenAuthLogin`, `TestGoldenRetryLoop`, `TestGoldenCounter`, and `TestGoldenStdDivision` compile and compare committed output (`compiler/emit_golden_test.go:11–85, 190–211`). In particular, the original division golden remains unchanged while B00 is tested in an additive fixture. |
| Linked integration | Run a real source factory, carry its returned record, and invoke it through a real consumer. Cover all four dividend/divisor signs, zero dividend, and zero divisor. Check return/error payloads, imported target revision, and isolation of callback execution from parent scripts. A nested invocation must inherit resource accounting rather than reset it. Linked factory expectations must use the same private receipt comparator. | `runLinkedPure` requires an ordinary clean program and a matching root revision before real execution; it does not grant coverage credit (`compiler/linked.go:68–177`). `TestRatioPilotLinked` demonstrates a committed linked fixture over G2 output (`compiler/ratio_pilot_test.go:25–51`). |
| Go/TypeScript parity | Use the same input/expected-outcome vectors in Go and emitted TypeScript. Keep `bind(3)` and `bind(5)` alive simultaneously, invoke both, then invoke the older closure again. Include a large integer beyond host-number exactness; a separate pure record-returning decimal callback checks that canonical decimal strings survive. Test supported data captures, not serialization of JavaScript functions. Node absence is a gate failure, not a skip. | `TestRatioPilotParity` explicitly requires Node and compares bigint/exact-decimal output (`compiler/ratio_pilot_test.go:53–95`); `tsBase` preserves bigint and decimal-string representation (`compiler/emit.go:11–15`). A shared mutable environment or host-number conversion fails these vectors. |
| modcheck | Add source fixtures at a scanned root/depth, using existing base-name `provides` and `uses` pins. Check wrong/missing pins and duplicate providers through the ordinary paths. No invented closure helper names may appear in headers. Assert no accidental superseded-file skip. Compiler use analysis additionally needs a reference-only dependency test: modcheck's header scan is not a substitute for that body-level check. | The inspected scanner groups files by program and resolves std providers across roots; source pins remain base names (`tools/modcheck/main.go:1–220`). Existing generic row spelling already has a deliberate regex allowance, which must remain intact. |
| gramcheck | Add `Fn`, `fnref`, and `invoke` samples; captures and input labels; nested generic-record mentions inside Fn; and empty/multiple error lists. Compile every new regex under the existing Go-regexp gate. Add targeted grammar-order probes so Fn syntax does not steal comparison operators, generic constructors, revision pins, or strings. No lookbehind. | The gate checks representative scopes with Go regexes; it is a smoke test, not a complete TextMate precedence proof (`tools/gramcheck/main.go:1–220`). Its success alone cannot establish full editor tokenization. |
| LSP | Diagnose the positive callback fixture without synthetic-name collisions. Test invalid capture/input/error locations, UTF-16 spans after non-ASCII text, generic source-name explanations, unresolved provider edits, and repeated close/change/reopen cycles. A broken or partial world must not execute callbacks; target-set certificates cannot survive a changed world. CLI and LSP must agree on callable acceptance even though their current phase order differs. | `tokenSpan`/`spanDiag` use UTF-16 source spans; `diagnoseWith` loads providers and runs its own ordered pipeline (`compiler/lsp.go:75–142, 350–473`). `GenericBase` routing is checked separately by the G2 cross-module fixture (`compiler/generics_types_test.go:250–314`). |
| Lint | Mandatory capture names and `input` must not be reported as redundant positional names. Nested ordinary record constructions retain normal linting. `MatchInvoke` must not accidentally receive a named-call chain/forward/relay rewrite that loses its semantics. For positive fixtures require both **no findings and no skipped files**, not merely process exit zero. | `lintRedundantNames` dispatches by expression kind; `lintFiles` records parse skips, and `runLint` can return zero after printing skips (`compiler/lint.go:80–167, 217–375`). Treating a parse skip as green would conceal a broken Fn parser. |

**Checked falsification for the gate strategy.** The original division golden would detect unsolicited output changes, the ratio parity test refuses a missing runtime, and the lint runner can otherwise mask unparsed input. Those inspected bodies invalidate “one successful Go test or lint exit code proves the slice.” The gate plan therefore preserves existing fixtures, adds behavior-specific negative cases, and distinguishes skipped work from passing evidence.

### Q8b — Execution order for the eventual implementation

**Verdict: no invocation execution or std migration until its static and runtime prerequisites are present together.**

This is sequencing for a later implementation, not authorization to make these changes now:

1. Resolve the customer/source and global namespace gaps below; reserve diagnostic identities after examining the registry. Record the accepted evidence revision and baseline. Keep the actual G2-only wrapper migration independent of B00.
2. Add syntax/type-shape support and reference/type traversal to the existing two-pass expansion, with rejection behavior for not-yet-enabled invocation. Preserve all old stamps and Fn-free goldens.
3. Establish shared static preparation, binding, graph admission, and indirect-cycle refusal in CLI and LSP. No evaluator fallback may bypass this barrier.
4. Add opaque runtime construction, capture snapshots, private factory receipt comparison, and TypeScript callable types/thunks. Extend every affected walker and identity representation before accepting a new AST kind.
5. Enable invocation only with outcome proof, isolated execution/coverage, error propagation, linked vectors, and Go/TypeScript parity. Keep ordinary named calls and their evidence rules unchanged.
6. Run the complete gates on a real checkout; accept the additive closure pilot only when every gate is green. Decide any broader std migration from actual source counts and semantics, not from the existence of Fn alone.

The standard commands to retain, executed from the actual checkout root, are:

```sh
go test -count=1 ./...
go run ./tools/modcheck
go run ./tools/gramcheck
```

Add the callback unit, linked, parity, LSP, and lint probes to the Go suite so the first command actually covers them; a name in this document is not an installed gate. An optional focused `go test -count=1 ./compiler -run '^TestFunctionValues'` is useful during implementation but never replaces the full suite. Generated files are reviewed alongside their producing source after gates run; changing expected output alone is not verification.

**Checked falsification.** The factory cannot be tested by existing `vEq` when its record contains a new opaque value, and linked execution rejects unadmitted graphs before execution (`compiler/eval.go:367–445`; `compiler/linked.go:27–65, 68–177`). Enabling invocation before either prerequisite creates a value/test mismatch or an authority bypass. The order above avoids those dependencies; it does not assert that independently enabling every parser form is already safe.

## 9. Follow-up interface, open questions, and risks

### Q9a — What does the outcomes design inherit?

**Verdict: leave a callable input / success-carrier / explicit-error contract, not a claim that `and_then` is now implementable.**

The catalogue's desired shape is `Fn<T,U,E2>` inside `std__outcome__and_then`, with `E1 union E2` explicitly represented at specialization (`docs/ASTRA_STDLIB.md:190–207`). B00 provides a unary input, named-record success, a concrete declared error list, and an empty capability contract. The follow-up can compose an admitted callback at an explicit invocation point and preserve the callback's exact error identity and payload. It must not retag every error with the caller's operation merely because invocation is indirect.

The remaining interfaces must be explicit:

| Surface | B00 boundary | Separate follow-up obligation |
| --- | --- | --- |
| Successful result `U` | `R` is a named success record. Its fields are carried by the existing Ok envelope. | Define bare-success emission and any explicit adaptation between catalogue `U` and a record carrier. No implicit `Value<U>` wrapping is introduced here. |
| Errors `E1`, `E2`, union | Only a concrete literal list in a callable signature; exact set compatibility. | Define an explicit error-set parameter kind, its specialization/bind syntax, and deterministic set union. Do not pretend an ordinary G2 type parameter is already an error-set parameter. |
| Outcome as input/output value | Invocation produces the existing dispatched outcome; it does not introduce a first-class `Outcome<T,E>` value constructor. | Supply the actual outcome-value representation and transport rules needed for `and_then`, `map`, `zip`, and recovery. The missing outcomes module was not reconstructed from its alleged count. |
| Conditional composition | A callback is invoked only in the selected `match invoke` path, not when the closure is constructed. | Ensure an incoming failure does not invoke the success callback. Do not replace conditional evaluation with an eager product or use bare-return `select` as an unstated prerequisite. |
| Test evidence | Real pure callback execution, ordinary parent coverage, no indirect relay exemption. | Address unrealizable declared-error arms without weakening `emits` or inventing scripted callback evidence. Any wider relay certificate requires its own design. |
| Higher-order/effectful functions | No Fn captures, Fn-bearing input/success, host callbacks, or effectful graphs. | Design such extensions separately, including termination, capability transport, and verification. B00 is not the complete future Functions surface. |

**Checked falsification.** `declaredOkShape` still requires a record, while the catalogue's generic signatures use bare type variables and enumerated error sets (`compiler/emit.go:285–308`; `docs/ASTRA_STDLIB.md:190–207`). Claiming the 464→1 collapse follows from B00 would smuggle in both missing mechanisms and an unverified source inventory. Keeping these boundaries explicit preserves the reverted bare-return limit.

### Open questions and risks not promoted to verdicts

**Customer justification.** The real wrapper example needs G2, not Fn. The proposed divisor closure demonstrates a capability but not a measured maintenance win. Before treating B00 as the blocker for a large std deduplication, supply one real repeated algorithm whose variable computation cannot be represented more cheaply by existing data and explicit dispatch. The absent rounding families and codecs prevent that judgment in this snapshot.

**Conservative graph precision.** A signature-wide set of all address-taken targets can reject safe programs when independent callback uses share a signature. Fixture additions can also enlarge that set. The rule is deliberately conservative and its safety counterexample is examined in Q3c, but its practical rejection rate is not measured. A future flow-sensitive improvement needs its own proof; do not quietly narrow the target set to observed test calls.

**Exact error sets versus coverage.** A callable with an unrealized declared error may be well-typed yet difficult or impossible to consume under the existing evidence law. Error adapters can also change nominal error identity. B00 chooses neither automatic set narrowing nor a new invocation-relay waiver. The practical cost for real codec/outcome customers remains unmeasured.

**Ownership and resource cost.** Snapshot capture can retain large records or byte buffers. Immutable sharing is desirable only where existing runtime ownership guarantees it; hostile host mutation is outside the declared boundary, not solved by calling a closure “immutable.” Allocation, retained memory, invocation overhead, repeated editor admission, and worst-case graph size require measurement. Resource exhaustion must remain a failure, not a declared mathematical result or a termination proof.

**Revision identity and exposure.** A change to captured parameter selection, target revision, Fn result/error contract, or admitted invocation behavior must not disappear from accepted revision identity. Q3a states the authority boundary, but the fingerprint/baseline implementation was not fully audited. Capture receipts and diagnostics also need a full redaction audit; function signatures alone do not establish that no existing pretty-printer leaks a capture.

**Parser and tooling integration.** Type-aware splitting is necessary, but the complete parser/formatter/highlighter interaction is not verified by a few samples. `Small`/`Node` visitors, expression printers, lint rewrites, callee collection, proof encoders, and the emitter must all handle the new kinds or reject them explicitly. A regex smoke test is not proof of editor grammar precedence. These are release blockers where unexamined, not assumed successes.

**Value architecture.** Section 7's B-first ordering is a provisional experiment recommendation. It does not settle whether the missing codec unions are identical, whether adapters are lossless, or whether generic/recursive variants would be needed. Html nominal wrappers and the actual outcomes Compare implementation remain excluded.

## 10. Examined and omitted scope

The evidence sample is bounded. The inspected material includes the four requested design/roadmap inputs (the relevant outcome/Functions and codec tables, not all of `ASTRA_STDLIB.md`); the complete division source; G2's entry point and principal function/type stamping, substitution, and rewriting bodies; the callable-relevant parser/type/runtime/emitter areas cited above; linked execution; the G2 record tests; ratio linked/parity tests; selected golden tests; and the cited CLI/LSP, lint, modcheck, gramcheck, cycle, effect, and contract-admission bodies.

**Standard-library omissions.** No complete source audit was made of `ascii`, `html`, `quota`, `ratio`, `scalars`, `schema`, or `text`. Reading the ratio pilot tests is not reading all of the ratio module. The thirteen codec sources, outcomes implementation, and named rounding families were not available at the supplied `std/` locations in this snapshot; they were not reconstructed from a survey. Other branches, uncommitted local work, historical revisions, and every sketch directory were not searched exhaustively.

**Compiler/tooling omissions.** The complete diagnostic registry, naming checks, revision fingerprint implementation and accepted baseline, full type-checker match/proof walkers, full evaluator and test-runner equality call sites, all serializer/debug printers, complete emitter/import planner, formatter, complete TextMate grammar, all lint rewrites, provider discovery internals, build/CI configuration, and the entire existing test suite were not audited. The sections cited here do not establish exhaustive visitor coverage or a global absence of helper-name collisions. Those exact remaining integration checks belong to the next design/implementation review.

**Execution omissions.** No successful clone, compiler build, `go test`, modcheck, gramcheck, linked run, Node parity, LSP protocol session, lint run, or benchmark was performed. No remote write or source implementation was made. Inspection of a committed test establishes its asserted contract, not its passing status at this revision.

**Language exclusions.** Bare-return `select`; anonymous source lambdas; inferred captures/type arguments/errors; variadic or zero-argument callable values; function equality/reflection/serialization; closures in state, constants, sequence elements, errors, or variants; arbitrary nested generic arguments; effectful or foreign callbacks; captured callbacks; general higher-order recursion; callable contract reasoning; generic error-set algebra; generic outcome values; and html-wrapper unification are not supplied by B00. Source wrappers and existing named calls remain available only under their existing rules, not as a way around these exclusions.

## Unresolved

1. **Checkout and gate evidence.** Direct clone failed; the artifact is staged at the requested relative path, not written into a successfully cloned checkout. No execution evidence establishes that the current suite or this proposed feature is green. Apply the doc-only patch to the actual checked-out revision and run the implementation gates only when implementation exists.
2. **Requested rounding-family N→1 proof.** Supply the exact revision and files declaring `FloorQuotRem`, `TruncQuotRem`, `IntQuot`, `RatioQuot`, and `DecQuot`, their repeated algorithms, and affected consumers. They are not in the complete inspected division module. Section 6's real 2→1 wrapper collapse is G2-only and does not close the requested Fn-caused rounding-family migration.
3. **Codec and outcomes inventory.** Supply the actual thirteen codec `Value` declarations and consumers, the source behind the 29-module count, and the alleged outcomes Compare family behind 464→1. The top-level directory response cannot establish those counts or the equivalence of those APIs. Until then, section 7 remains a framed decision, not a source-verified choice.
4. **Shared preflight and visitor completeness.** The CLI/LSP ordering conflict was inspected and Q2d resolves the required architectural placement. The exact static-check extraction, all affected type/expression/proof/emitter walkers, target certificate storage, target-reference import planning, test equality call sites, and stale-editor-state invalidation have not been fully audited. No new kind may be enabled until omitted paths explicitly support or reject it.
5. **Global identities and revision baseline.** The complete name/diagnostic registry audit and accepted fingerprint implementation were not inspected. Allocate collision-free rule identities; verify references, capture-slot selection, callable signatures, and invocation semantics participate in the required baseline checks. The revision-2 choices in the G2 specimen are proposals, not verified accepted-baseline migrations.
6. **Runtime ownership and diagnostics.** Confirm capture storage ownership, resource-account inheritance, all mismatch/debug redaction paths, and Go/TypeScript behavior using actual execution. No allocation/performance/security-boundary measurement was made. The design does not establish resistance to hostile JavaScript outside checked Can linkage.
7. **Broader stdlib compatibility.** The bounded unary/pure/data-only core has not been tested against the missing codec or rounding algorithms. Whether it is sufficient for those customers, whether a finite mode value is cheaper, and how often exact-error coverage or signature-wide cycle rejection blocks useful programs remain unanswered empirical questions.
8. **Outcomes and later language surfaces.** Bare-success emission, concrete syntax/kinds for generic error sets and their union, actual first-class outcome representation, any new invocation-relay proof, effectful callbacks, and higher-order callable contracts require separate designs. B00 neither implements those interfaces nor relies on them implicitly.
