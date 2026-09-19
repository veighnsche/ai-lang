# B05 — Structured async joins, explicit authority, and terminal outcomes

Status: **proposed bounded design, for approval**. No implementation is authorized by this document.
Date: 2026-09-19.
Evidence revision: `veighnsche/can-lang@8fccbb79aee193ef6f7c64ea09e441acb87e16fd`.
Deliverable path: `docs/b05-async-design.md`; `b05` was free in the complete `docs/` listing at that revision.
Semantic verification: **source inspection and checked design falsifications only**. No compiler, test, lint, editor, linkage, TypeScript, or grammar gate was run.

**Decision in one sentence:** admit an all-terminal `match join` over named operations, with a complete outcome product and checked noninterfering authority; keep incomplete work entirely inside the runtime, and admit host work only under a revision-bound, trusted settlement contract.

This is the async design, not its implementation plan. Sections D1–D5 settle the requested surface, static discipline, evidence model, host boundary, and acceptance gates. The gates describe what later implementation must demonstrate, not a sequence of implementation slices.

## 0. Evidence boundary and scope rulings

### Q0a — Which input statements are current facts?

**Verdict: retain the requested constraints, but distinguish current source from historical design text.**

The evidence revision was selected once. References in this document are repository paths, named sections, or source symbols at that revision. A checked falsification means that the proposed rule was compared with those inspected sources and with the stated counterexample; it does not mean that proposed syntax was compiled. The recommendations below are not maintainer ratification or an external probabilistic verdict.

| Input or tempting inference | Inspected evidence | Ruling for B05 |
| --- | --- | --- |
| The audit already specifies async. | `docs/a/ASTRA_AUDIT.md`, **Task 3 §4**, calls its code an illustrative join fragment. Task 4's Async row adds an explicit liveness assumption or bound. | Adopt the obligations; rederive the syntax and execution rules. Do not treat the unrelated Task 1 §4 as the async section. |
| Host implementations are still unwritten. | `std/host/host.can`, `std/host/host.externs.ts`, and the host-shelf row of `docs/stdlib-remaining.md` contain synchronous implementations. | The missing boundary is bounded host completion and concurrent authority, not all host code. B02's old non-goal does not describe current absence. |
| Every arm must literally execute. | `docs/a/a03-branch-coverage.md`, retrospective amendment, distinguishes execution from the specifically authorized identity-relay certificate. | Preserve executed / certified / uncovered categories. B05 grants no new join certificate and fabricates no hits. |
| B00's printed invocation form is the current parser form. | `compiler/parse.go`, `parseInvokeHead`, explicitly records `invoke <bare-name> with <arg>` as the settled deviation from B00's sketch. | Do not reopen or copy the obsolete invocation spelling. Join branches are named calls, not invocations. |
| B01's proposed declaration syntax has landed. | `docs/b01-type-syntax-redesign.md` is proposed; the inspected host source uses `Seq<T>`, inline errors, and `brand ... is str`. | Samples use the inspected declaration family. B05 is not a type-syntax migration. |
| Function bodies need the audit's `=` separator. | `compiler/parse.go` explicitly rejects a lone `=` after function metadata. | No such separator appears in the proposed samples. |
| A11's unit-step rule is the entire current recursion policy. | `docs/a/a19-recursion-schemas.md` adds Euclid and narrowing; `FnDecl.DecNames/DecSchema` represent them. | Preserve all currently admitted schemas. A join is no escape from their guards, steps, or whole-program cycle checks. |
| An extern without declared cell effects is pure. | `ExternDecl` has no effects or completion fields; `requiredEffects` reasons about state and Can callees, not a complete host footprint. | Absence is not a concurrency certificate. Existing externs are ineligible inside task execution until the new boundary explicitly admits them. |

**Checked falsification.** The actual host implementations falsify the brief's absence claim. The parser's explicit separator rejection falsifies a verbatim adoption of the audit sample. The a03 amendment falsifies equating structural certification with execution. These corrections change neither the user's exclusions nor the audit's async obligations. [S1–S5, S9–S11]

### Q0b — Does this include Resources or reorder the standard library?

**Verdict: keep Resources separate, but settle the cleanup obligation that async cannot safely omit.**

There are no source resource handles, cancellation tokens, lifetime binders, user cancellation operations, deadline values, races, streams, channels, or detached tasks here. Time limits are literal metadata on a host contract, not values that a Can program stores or passes. A host adapter's internal cleanup mechanism is part of its trusted contract, not a new Can capability value.

The catalogue's **Async** dependency label includes explicit cancellation/deadlines. B05 intentionally supplies only the structured-join and bounded-host part of that wider label. It does **not** mark that catalogue dependency, HTTP, SQL, or UI as fully unblocked. The catalogue and a13's adoption order remain intact: HTML and string/brand work do not wait for this design. General partitioned host resources, cancellation, and acquisition/use/release still require the separate Resources decision.

**Checked falsification.** `ASTRA_STDLIB.md`, “Dependency labels,” names both cancellation/deadlines and opaque resource protocols. Calling this document the complete catalogue Async surface would erase those requirements. Conversely, postponing the meaning of timeout completion would already leave Task 3 §4's join-once obligation unsatisfied. D4 specifies settlement without adding the postponed value surface. [S1, S12–S14]

## D1. Surface and look-and-feel

### Q1a — What is the one concurrency construct?

**Verdict: a distinct `match join <site>` with explicitly typed, named-call branches and positional outcome-product arms.**

```can
match join reads
  left: Left__Value
    effects [Left__total.read]
    call read__left(left_down)
  right: Right__Value
    effects [Right__total.read]
    call read__right(right_down)
  on Ok l, Ok r => Ok(l.value, r.value)
  on left.down _, Ok _ => pair.left_failed()
  on Ok _, right.down _ => pair.right_failed()
  on left.down _, right.down _ => pair.both_failed()
```

`reads` is a **local site identifier**, not a type, expression, result object, or global declaration. Site names are unique within a function; branch labels are unique within a site. Both use lower-case identifiers with optional digits/underscores after the initial letter. Generic stamps retain their source labels under distinct function identities. Labels cannot be read as expressions and never denote running work.

Each branch has exactly one existing named `call` head. Its return annotation must equal the resolved, nominal success-record type. Its mandatory `effects [...]`, including explicit `[]`, lists the checked authority delegated to that branch. A direct generic call must retain the existing explicit specialization rules; the join does not infer types or errors.

The arm has one leading `on`, then one ordinary outcome pattern per branch, separated by commas in declaration order. `left.down` above is the **real declared error kind**, not a branch-qualified pseudo-outcome. `Ok l` binds the left record in its slot. Identical error kinds in two branches remain distinguishable by slot position. There is no tuple or outcome-product type available in value positions.

A branch has no anonymous source body, completion function, or independently executable identity. Sequential logic, guards, error translation, and inner joins belong in ordinary named Can functions. This is lexical ownership of work, not a new function-value mechanism.

The audit sketch is adopted or replaced element by element:

| Sketch element | B05 treatment | Grounding |
| --- | --- | --- |
| `match join` | Adopt as a new match family. | Outcomes remain exhaustively handled control flow, not a new value API. |
| Nominal-looking `Pair__Read` after the head | Replace with a local site identifier. | Stable ownership/trace identity is needed; an exported task/result type is not. |
| `order [left, right]` | Remove. Declaration order fixes inputs and slots; Q3c separately defines causal execution order. | One duplicated order list cannot establish real completion order. |
| Named, typed branches | Adopt, with exactly one resolved named call each. | Existing nominal success records and function declarations supply the operation contract. |
| Per-branch effects | Adopt only with actual resolved capabilities and transitive checking. | a09's cell identities and the audit's write/read rule, not annotation text alone. |
| Branch-qualified pseudo-outcomes | Replace with positional slots containing real `Ok` or declared error patterns. | Reuse existing error identity; do not invent `left.Ok` as an outcome name. |
| Explicit combined outcomes | Adopt and generalize to the complete finite product. | No implicit failure-selection policy; Q2b prices the case count. |

**Checked falsification.** `parseMatchArmsKind` already splits comma-delimited pattern slots, but chooses outcome-pattern parsing only for `MatchCall`/`MatchInvoke`; ordinary call matches have singleton semantics. Thus spelling this as an existing value product or multiple ordinary call scrutinees would select the wrong checks. A distinct `MatchJoin` is a necessary future addition. The audit's `left.Ok` convention would also conflate a local label with the language's existing `Ok`/declared-error identity; positional slots avoid that extra namespace. [S1, S4, S5]

### Q1b — What does a complete two-branch example look like?

**Verdict: show ordinary deterministic Can tasks first, with real rows for every product arm.**

The following is a proposed, self-contained illustrative module. It is not a new repository sketch and has not been compiled. Its new syntax is confined to the join. Both readers have an actual failure path; no injected completion is needed to make the four parent arms witnessable.

```can
mod joinread
  provides [read__left, read__right, pair__read, Left__Value, Right__Value, Pair__Value]
  uses []
  emits [left.down, right.down, pair.left_failed, pair.right_failed, pair.both_failed]

error left.down()
error right.down()
error pair.left_failed()
error pair.right_failed()
error pair.both_failed()

type Left__Value rev 1 (
  value: int
)

type Right__Value rev 1 (
  value: int
)

type Pair__Value rev 1 (
  left: int
  right: int
)

state Left__total: int = 10
state Right__total: int = 20

fn read__left(down: bool) -> Left__Value rev 1
  effects [Left__total.read]
  emits [left.down]
  tests
    available(false) => Ok(10)
    unavailable(true) => left.down()
  match down
    true => left.down()
    false => match call state__get(Left__total)
      on Ok value => Ok(value.value)

fn read__right(down: bool) -> Right__Value rev 1
  effects [Right__total.read]
  emits [right.down]
  tests
    available(false) => Ok(20)
    unavailable(true) => right.down()
  match down
    true => right.down()
    false => match call state__get(Right__total)
      on Ok value => Ok(value.value)

fn pair__read(left_down: bool, right_down: bool) -> Pair__Value rev 1
  effects [Left__total.read, Right__total.read]
  emits [pair.left_failed, pair.right_failed, pair.both_failed]
  tests
    both_available(false, false) => Ok(10, 20)
    left_unavailable(true, false) => pair.left_failed()
    right_unavailable(false, true) => pair.right_failed()
    both_unavailable(true, true) => pair.both_failed()
  match join reads
    left: Left__Value
      effects [Left__total.read]
      call read__left(left_down)
    right: Right__Value
      effects [Right__total.read]
      call read__right(right_down)
    on Ok l, Ok r => Ok(l.value, r.value)
    on left.down _, Ok _ => pair.left_failed()
    on Ok _, right.down _ => pair.right_failed()
    on left.down _, right.down _ => pair.both_failed()
```

The module exposes its functions and records, not its cells. A result binder contains a completed immutable value. Returning `pair.left_failed()` does not imply that the right branch was skipped; that arm is reached only after the right branch has also terminated. Neither branch carries `given`, because both are same-file Can calls. The state builtins inside them also take no scripts.

**Checked falsification.** a09's local-call/store distinction makes a scripted `read__left` an invalid example of the existing evidence discipline. Its private scalar cells and deterministic `state__get` support this example's chosen effects and values. The four boolean rows reach four real combinations; collapsing them into a left-error fallback would lose the deliberately distinct both-failed result. Current declaration/body conventions are checked against `std/host/host.can` and the parser, not the old audit separator. [S4, S9, S10]

### Q1c — How do singleton wrappers, nesting, and conditional use work?

**Verdict: reject an empty join and a redundant native singleton; permit a singleton only when its branch directly names a bounded host extern. Nest through ordinary functions.**

| Shape | Ruling | Reason |
| --- | --- | --- |
| Zero branches | Reject. | No task and no outcome product to dispatch; do not invent a unit-effect construct. |
| One branch naming an ordinary Can function | Reject. | Use its ordinary `match call`. A wrapper whose body contains a join is still an ordinary function to its callers. |
| One branch naming a completion-contracted extern | Admit. | This is the one explicit source boundary that acquires a terminal host outcome under D4's contract. |
| Two or more branches | Admit subject to D2. | Cardinality does not waive authority, totality, coverage, or host-contract checks. |

**Singleton host sample.** `net__fetch@2` and its `Fetch__Value`, `net.down`, and `net.timeout` declarations are the proposed contract in D4. In a separate consumer module, `uses` includes the extern and success-type pins. The function below is a declaration fragment; the enclosing module header and the D4 provider are its stated context.

```can
fn fetch__bounded(key: str) -> Fetch__Value rev 1
  effects [Host__World.use]
  emits [net.down, net.timeout]
  tests
    found("alpha") => Ok("payload")
    unavailable("alpha") => net.down()
    expired("alpha") => net.timeout()
  match join fetch
    request: Fetch__Value
      effects [Host__World.use]
      call net__fetch(key)
      given
        found => [exchange args (key = "alpha") outcome Ok("payload")]
        unavailable => [exchange args (key = "alpha") outcome net.down()]
        expired => [exchange args (key = "alpha") outcome net.timeout()]
    on Ok result => Ok(result.body)
    on net.down _ => net.down()
    on net.timeout _ => net.timeout()
```

`given` is an optional **branch-owned block after its call, at the same indentation as `effects` and `call`**. It scripts that one named foreign operation. There is no join-wide table that fabricates the aggregate result. The `expired` row observes a terminal timeout outcome; it does not run a timer or demonstrate the host's settlement promise.

**Nested-join sample.** Given the complete `pair__read` above and an ordinary pure `stamp__value() -> Stamp__Value` returning `Ok(7)`, the following body fragment has an outer join whose `pair` child performs the inner `reads` join:

```can
match join report
  pair: Pair__Value
    effects [Left__total.read, Right__total.read]
    call pair__read(left_down, right_down)
  stamp: Stamp__Value
    effects []
    call stamp__value()
  on Ok values, Ok mark => Ok(values.left, values.right, mark.value)
  on pair.left_failed _, Ok _ => pair.left_failed()
  on pair.right_failed _, Ok _ => pair.right_failed()
  on pair.both_failed _, Ok _ => pair.both_failed()
```

The containing function has a three-field named success record, those three declared errors, the union of the shown effects, and four corresponding boolean rows. The `stamp` helper has an ordinary row of its own. These obligations are not satisfied merely by printing the body fragment. There is deliberately **no inline branch-body alias** for this nesting: the existing named-function boundary supplies return shape, effects, errors, recursion admission, and evidence ownership.

**Join in an ordinary outcome arm.** Given an ordinary `gate__check(allowed)` returning either `Gate__Pass` or `gate.denied`, this body fragment starts neither reader on a denied path:

```can
match call gate__check(allowed)
  on gate.denied _ => gate.denied()
  on Ok _ => match join reads
    left: Left__Value
      effects [Left__total.read]
      call read__left(left_down)
    right: Right__Value
      effects [Right__total.read]
      call read__right(right_down)
    on Ok l, Ok r => Ok(l.value, r.value)
    on left.down _, Ok _ => pair.left_failed()
    on Ok _, right.down _ => pair.right_failed()
    on left.down _, right.down _ => pair.both_failed()
```

The containing function needs one denied row plus the four admitted combinations. Its effects remain conservative over the whole body. A join inside the selected arm does not make the gate eager. A foreign gate retains its own ordinary `given`; the join changes no upstream call semantics.

**Checked falsification.** B00's conditional-composition boundary and a11's call-graph rule rule out extracting task completion into callback values or pretending an eager product preserves a conditional call chain. The gate example preserves conditional acquisition; the nested example creates ordinary named graph edges and cannot hide a cycle. B02 supports the provider pin in the singleton sample, but its existing extern signature alone does not supply D4's new contract. [S2–S5, S7, S16]

### Q1d — What is explicitly rejected, and why?

**Verdict: reject `async`/`await`, `async fn`, colored function types/call sites, and any source-level suspend/await operation; reject pending values and bare task creation.**

The user's constraint is decisive. The proof-discipline reason is narrower than a claim that other language designs are impossible: B05's checked control nodes must account for every initiated operation, authority transfer, terminal outcome, and path to the enclosing result. Introducing an `await` expression, `async fn` modifier, colored return type, or an implicit rejection channel is not an authorized extension of those nodes. A different language could verify explicit await states; that possibility does not authorize this forbidden surface or make it necessary here.

A `Pending<T>`, promise, future, task handle, or deferred value would require additional transport, lifetime, and join-once rules for arguments, records, collections, errors, and exits. B05 instead makes the owner and all terminal consumers lexical. These forms are rejected, not given placeholder semantics:

```text
async fn ...
await operation(...)
spawn operation(...)
detach operation(...)
let work = pending operation(...)
return task_handle
```

`join` is also rejected in expressions, constructor arguments, type annotations, contract predicates, `forward`, `select`, and a branch completion callback. Task inputs, task success records, task error payloads, and host payloads must be transitively Fn-free. A label cannot become a callable or a value. Pure function-value use **inside** an ordinary task body remains subject to all existing B00 rules and to D2's complete target-graph admission; it is not a task-completion mechanism.

**Checked falsification.** B00 §9 and Unresolved explicitly exclude effectful/host callbacks and Fn-bearing invocation input/success. A callback disguised as a join branch would cross that boundary. `Value` has completed data/outcome and opaque Fn representations but no pending representation; adding one is not a free spelling change. This ruling also rejects the brief's overbroad inference that await syntax is intrinsically unverifiable: the user constraint and this workspace's specified proof costs suffice without that assertion. [S2, S3, S15]

### Q1e — Which grammar elements are genuinely new?

**Verdict: reserve a new join family; do not broaden existing outcome or value matches accidentally.**

The following is grammar notation for review, not an edited grammar file:

```text
join           := 'match join' local_site NEWLINE
                  INDENT branch+ product_arm+ DEDENT
branch         := local_label ':' success_record NEWLINE
                  INDENT 'effects' capability_list NEWLINE
                  named_call NEWLINE
                  [ 'given' NEWLINE INDENT exchange_rows DEDENT ]
                  DEDENT
product_arm    := 'on' outcome_slot (',' outcome_slot)* '=>' rhs
outcome_slot   := 'Ok' binder_or_discard
                | declared_error_kind binder_or_discard
```

`named_call` reuses the existing `call name<TypeArgs>(args)` head and argument conventions. The only branch expression is that call; `rhs` is the existing result-or-nested-match block. Branches precede arms, annotations precede the call, and `given` follows it. Empty joins, duplicate labels, wrong arity, and invalid singleton kind are specific admission failures. Ordinary `match invoke`, `match chain`, value-product `|` patterns, and ordinary `given` placement do not change.

Eventual grammar touches are: the `match join` head; local site/branch declarations; branch-owned metadata and scripts; product-arm parsing with exactly one `on`; the completion metadata in D4; and the special host authority spelling. Parser, formatter, expression visitors, expansion/cloning, use collection, diagnostics, fingerprinting, highlighter, and lint must recognize or explicitly reject these additions. Naming the touches does not implement them.

**Checked falsification.** `parseValuePattern` deliberately admits `|` only for value patterns, while `parseMatchArmsKind` routes outcome patterns separately. Reusing the value matcher without preserving that boundary would admit a second unpriced proof abstraction. `isMetaHead` and `parseExprBlock` also show that metadata and body parsing are distinct; a branch declaration cannot be made real by changing only a highlighting regex. [S4, S5]

## D2. Static discipline and semantic guarantees

### Q2a — How is every initiated task joined exactly once?

**Verdict: enforce a lexical, all-terminal barrier, not a linear value that users must remember to consume.**

For each dynamic join activation, the semantics has the following internal states:

```text
prepare inputs -> register all children -> run children
               -> all children terminal -> dispatch one product arm -> retire
```

The following are mandatory rules, not scheduler preferences.

1. Resolve all branch signatures, types, effects, contracts, and recursion/definedness admission before executing the program. At runtime, prepare argument values exactly once, branch by branch and argument by argument in declaration order, before any child starts. Preparation accepts only the total, pure, Fn-free expression subset described in Q2d. It cannot consume an external observation or use a sibling result.
2. Allocate one internal outcome slot per branch. Register every child before running any child body or accepting a completion that could dispatch the parent. A declared error from the first child must not prevent another declared child from starting.
3. Each child has one terminal transition. A nested child is terminal only after all of its descendants are terminal. The child keeps its delegated authority until then.
4. The parent evaluates no arm body until every slot is terminal. No early return, error propagation, first-success selection, or failure-triggered sibling cancellation is implicit in the join.
5. Dispatch exactly one product arm and consume every slot once. Retire that activation before executing the arm body. A join in that body is a new activation, not reuse of old slots.

The internal identity is the checked world/root invocation, owning function specialization, local site, parent-owned activation number, and branch ordinal, recursively extended for nesting. It is not a source value. Activation counters belong to their logical parent, so racing siblings cannot change each other's identities.

A duplicate completion, attempted slot overwrite, unknown activation, or already-retired receipt is a runtime/host-contract violation. It must not silently become a second outcome or a successful join. D4 specifies the failure boundary. “Exactly once” here means language task initiation and terminal consumption; it does **not** promise exactly-once network delivery or exactly-once remote business effects.

**Checker formulation.** There is no independent spawn or join-use node to count. Each `MatchJoin` owns a fixed nonempty branch vector and a checked product dispatcher. The typed node receives a world-specific admission certificate only when every branch has one resolved target, one success shape, one outcome domain, one verified authority summary, and the required liveness/definedness admission. Every evaluator/emitter entry rejects an absent or stale certificate. Ordinary named calls in the branch vector are real graph edges, not inert metadata.

**Checked falsification.** Lowering to a chain that returns on the first error violates rule 2 or 4. Returning a timeout while the worker remains able to write violates rule 3. Permitting a branch to hand its slot to another function recreates the forbidden pending-value surface. `MatchChain` is explicitly sequential in `parse.go`, so elaborating a join into that node would misrepresent its all-child semantics. [S1, S3–S5]

### Q2b — What precisely must the outcome product cover?

**Verdict: enumerate the full declared tag product, with one explicit arm for each tuple and no new compression rule.**

For branch `i`, let `O_i = {Ok} ∪ emits(target_i)`. Its payload shape is attached to each tag by the target declaration. The join domain is the Cartesian product `O_1 × ... × O_n`. This is a proof domain, not a passable value type.

A join arm must name one valid tag in every slot. Duplicate tuples, missing slots, extra slots, wrong branch error kinds, wrong payload access, and duplicate payload binders are rejected. Empty payloads bind `_`; nonempty payloads may also be discarded with `_`. Discarding a payload does not discard its tag. Bare `_` as an outcome slot, whole-product catch-alls, `else`, outcome `|` alternatives, and “any error” arms are not admitted in B05. Therefore accepted arms do not overlap, and arm source order does not select a failure policy.

For two operations with `Ok`, `down`, and `timeout`, all nine tuples remain:

| Left slot | Right slot | Required explicit dispatch |
| --- | --- | --- |
| `Ok` | `Ok` | Both succeeded. |
| `Ok` | `right.down` | Left success and right service failure. |
| `Ok` | `right.timeout` | Left success and right terminal timeout. |
| `left.down` | `Ok` | Left service failure and right success. |
| `left.down` | `right.down` | Both service failures. |
| `left.down` | `right.timeout` | Left service failure and right timeout. |
| `left.timeout` | `Ok` | Left timeout and right success. |
| `left.timeout` | `right.down` | Left timeout and right service failure. |
| `left.timeout` | `right.timeout` | Both terminal timeouts. |

This table describes outcome checking; two host branches still separately fail B05's conservative host-authority rule. The same three-outcome product can be exercised by two terminating native classifiers that return those declared kinds.

The compiler may compare canonical tuple sets and report bounded, deterministically ordered missing witnesses without enumerating them all in a diagnostic. It must not truncate the domain and claim exhaustiveness. A proof-work limit is a visible failure/inconclusive result, not permission to ignore remaining tuples. No arbitrary maximum branch count is presented as a solution to the combinatorial problem.

**No unjustified abstraction is adopted.** An ordinary named adapter may intentionally translate several errors into one declared error, with its own complete handling and rows. That changes the API's observable distinctions and must be explicit and revision-relevant. Wrapping outcomes in existing ordinary data merely moves the classification obligation; it does not establish a new join certificate. B05 neither invents first-class generic outcomes nor silently applies a first-error adapter.

There are two real costs:

- The required tuple count is `∏ |O_i|`. Three-outcome branches give `3^n` tuples, and every explicit product arm needs legitimate evidence. A small diagnostic witness budget does not reduce this cost.
- Shared immutable inputs can correlate outcomes. Two identical local classifiers called with the same boolean can each have two declared outcomes while their mixed product tuples are unreachable. Their join can be exhaustive yet fail coverage. An unrealized declared error produces the same problem. Tests must not invent native outcomes, and B05 grants no relational or join-relay certificate to escape it. Refactoring, a deliberately coarser API, or a separately justified future proof abstraction is necessary.

**Checked falsification.** The audit explicitly prices the nine-case example. a03's amendment separates evidence from proof certificates; a81 prohibits using scripts as facts. Treating a foreign script, a probability, or an observed correlation as grounds to delete a product cell would cross both boundaries. The current outcome parser also excludes the value matcher's `|` abstraction. This conservative design contains no general cure for product explosion or correlated unreachable cells, and does not claim one. [S1, S5, S8]

### Q2c — What counts as disjoint branch authority?

**Verdict: partition the resolved resource identities, check writes against sibling reads and writes, and reject unknown host footprints.**

For native cells, define `R_i` and `W_i` from the complete checked effect summary of branch `i`. For every distinct `i,j`, require:

```text
W_i ∩ (R_j ∪ W_j) = ∅
W_j ∩ (R_i ∪ W_i) = ∅
```

Read/read sharing is allowed. Read/write, write/read, and write/write sharing are rejected. A declaration spelling is not the resource identity: the checker resolves the actual cell and owning module's canonical identity. Different wrappers, imports, aliases in diagnostics, or same-basename modules cannot manufacture independent cells.

The branch annotation must equal the callee's declared effect set after canonical resolution; an explicit empty set is required for an empty declaration. The whole-program effect check must additionally establish that the declaration covers the actual transitive body. Neither a smaller branch annotation nor a stale declaration hides a use. The enclosing function declares the union of its branch and arm effects under the existing visibility and staleness discipline. The same complete traversal includes nested joins and any already-admitted pure Fn target sets.

The parent performs no effects while children run. A nested join may delegate only authority already held by its containing task and must satisfy its own pairwise rule. Completed descendants do not release a task's authority to an ancestor that has not passed its barrier. Outcome-arm effects occur after restoration of the join's authority and retain ordinary sequential semantics.

**Private-cell access remains private.** A caller can account for an imported accessor's declared cell effect without gaining permission to call `state__get` or `state__put` on that private cell directly. B05 does not add exported cells, pointers to cells, or resource values.

**Host authority, bounded first cut.** Introduce one built-in static authority identity, `Host__World.use`, meaning exclusive access to the host boundary. It is neither a state cell nor a value and needs no `provides`/`uses` declaration of its own. All completion-contracted externs in D4 use it, and Can callers declare it transitively. It conflicts with itself. Consequently two concurrent branches that may reach host code are rejected even when their extern names, URLs, or declared business errors differ. A single host-using branch may overlap with native branches holding disjoint Can-cell authority; host adapters may not access Can cells or reenter the world.

A missing host footprint is not `[]`. An old synchronous extern anywhere in a task's possible call closure is rejected until that operation has an explicit conforming bounded contract. Pure-looking clock, random, environment, logging, and hashing calls receive no implicit exception. More precise host domains, read-sharing permissions, connection-level partition, locking, transaction protocols, and commutativity proofs are **deferred to the separate resource/protocol design**. B05's conservative rule is useful but does not enable arbitrary parallel HTTP/SQL fan-out.

**Checked falsification.** a09 and `checkEffects` validate real `Cell.read|write` capabilities, not arbitrary annotations. `ExternDecl` has no effect field, and `requiredEffects` is not a host-footprint analysis. Copying the audit's illustrative effect names into annotations would therefore not make them meaningful. Two aliases for one database connection falsify name-based host separation; two wrappers around the same private cell falsify function-name-based separation. The formulas above are the audit's actual write/read rule, rather than the stronger and unnecessary rejection of read/read sharing. [S1, S4, S9, S10]

### Q2d — What native-task termination claim is admitted?

**Verdict: preserve the whole-program termination policy, additionally reject unproved partial computation inside tasks, and separate progress from elapsed time.**

A native task's complete possible execution closure must satisfy all ordinary static checks and the following task-admission obligations. These are a proposed new admission rule, not a claim that the current compiler already certifies total task bodies.

**Finite control.** Include every branch call in the whole-program call graph. Preserve the existing direct self-recursion schemas: unit, Euclid, and narrowing. Their exact guard and argument shapes remain necessary at a recursive branch call just as at an ordinary call. Reject mutual/indirect cycles not admitted by the existing policy. Inspect all already-certified B00 targets, never only the targets used by rows. A nested join does not reset the caller's ranking proof or resource accounting.

**Defined primitive operations.** A finite computation that faults is not a declared terminal outcome. The initial task-safe primitive set is explicit: well-typed immutable data construction/projection, authorized seals, exact addition/subtraction/multiplication, supported comparisons and boolean operations, lengths, exhaustive matches, and well-typed scalar store access. Each is admitted only in its already-valid value/type positions. Other operations require an independently established total-kernel or local definedness rule; no unknown kernel gets a default approval.

For the first cut, division/remainder with a statically nonzero integer-literal denominator is permitted. The exact guarded remainder site certified by the existing Euclid schema is also permitted, and narrowing's literal divisor is already nonzero. Other dynamic division/remainder and index/slice operations are rejected in task execution unless an existing, inspected checker certificate establishes their definedness at that exact site. A successful row, a callee name, or an unverified `requires` annotation is not such a certificate. No new general partial-operation solver is hidden in B05. The same total-expression restriction applies to join argument preparation.

**Admitted host leaves.** The native-only admission scope has none. Under D4, every possible host leaf must instead have its bounded settlement contract and the enclosing effect summary must include `Host__World.use`. Direct ordinary calls to old externs cannot bypass the singleton boundary. Host-assisted task completion is conditional on every such contract, not a native termination proof about the external service.

**Progress assumptions, explicitly named.** The native theorem assumes a fair runtime, finite execution of each admitted primitive on its finite input, and adequate runtime resources. With finite branching and the existing well-founded recursion proofs, a native task has a finite computation tree. Joining finitely many such trees adds finite work and yields a terminal declared outcome. This is an eventual-completion argument under stated assumptions, not a measured or statically enforced millisecond limit.

Let `W_i` denote a branch's finite native work under those assumptions. A join's native work includes the sum `Σ W_i`, preparation, and its selected arm; it is not automatically `max W_i` merely because concurrency is allowed. A host leaf has D4's explicit elapsed-time assumption. Sequential host calls accumulate those assumptions along the admitted path. Neither the scheduler nor this document silently promises parallel CPU speedup or a constant end-to-end deadline for an input-dependent number of operations.

**Explicit deferral:** a mechanically checked numerical CPU/allocation/frame budget or end-to-end real-time bound for arbitrary native tasks is outside B05 and belongs with resource budgets. Existing `decreases` metadata, the structural proof, and the named progress assumptions are the current liveness basis. Host leaves, unlike native steps, are not admitted with an unbounded “eventually” promise. Sandbox limits remain visible resource failures, never mathematical outcomes or proof certificates.

**Checked falsification.** `evArith` has an explicit division-by-zero evaluation error, so a11's returned-outcome wording cannot be used as evidence that every finite arithmetic expression is defined. A test using divisor `1` cannot prove a task safe for divisor `0`. Conversely, rejecting every remainder operation would wrongly discard a19's guarded Euclid theorem. The bounded task-admission rule preserves that theorem while refusing unsupported sites. a08/a13/a19 separately identify termination and resource cost; B05 does not promote a depth backstop or fair scheduling assumption into a wall-clock proof. [S6, S7, S13, S15]

### Q2e — Can another root invocation or a host callback violate the partition?

**Verdict: serialize root admission within a checked runtime world, preserve persistent state across roots, and prohibit host reentry.**

Sibling separation alone is insufficient: two simultaneous roots could each acquire the same cells. The initial join-capable runtime therefore admits **one active root per world**. A root retains that admission until all its work is terminal. A host embedding may queue a later request outside the world or reject it at the embedding boundary; it cannot execute that request against the same store concurrently. This is a runtime-world admission rule, not a Can handle, extra result variant, or function annotation.

The world owns its cells. Production cells persist across sequential roots as they do across existing calls; they are not reset at each join or each production root. Decision-table rows still start from a fresh store. A host adapter cannot call back into Can while holding a root's host authority. Other processes and independently instantiated worlds are outside the cell-partition theorem; sharing their external resources needs an explicit host protocol, not the assumption that two runtime IDs imply independent databases.

Native branch effects are not transactional. If one branch writes its private partition and another returns a declared error, the completed write remains. Each branch observes its own prior writes. For pairwise noninterfering native branches, final cell state is independent of their interleaving; this does not establish an arbitrary invariant relating two cells unless that invariant is separately proved. No snapshot rollback, speculative commit, or hidden transactional store is introduced.

**Checked falsification.** a09's production module-scope cells and fresh test stores demonstrate why test isolation is not production-root exclusion. Running two separately valid joins against one cell store defeats sibling-only checks; the world admission rule closes that specific hole. Reusing a timed-out worker's host authority before it is quiescent would reopen it, which is why D4 makes cleanup part of terminal completion. [S1, S9, S10]

### Q2f — Which proofs may claim to understand a join?

**Verdict: structural join admission and observed evidence are supplied here; unsupported universal contract reasoning fails closed.**

The new static checker proves well-formed ownership, tag-product completeness, authority separation, task-graph admission, and the presence of the required host contracts. It does not prove a host implementation honest, eliminate unrealizable error combinations from a provider's `emits`, or synthesize function `ensures` from test scripts.

A universal contract verifier that has no semantics for `MatchJoin` must mark such a function inadmissible/unsupported for that verifier, rather than skip the node and report success. The function may still be admitted by the ordinary language and structural-join checks when it does not request an unsupported proof surface. Existing independently proved pure helpers keep their summaries under their existing rules; no joint postcondition about concurrent state is inferred here. A future join-summary rule needs its own design, especially for shared invariants and host assumptions.

All structural failures block execution in both CLI and LSP. Every AST walker must either include join branches and arms or explicitly refuse the new kind. Certificates are bound to the resolved world, specialization, callee dependencies, effects, outcome sets, and host contract metadata; editing any of those invalidates them.

**Checked falsification.** a81 allows a summary only after its own obligations discharge and expressly excludes `given` from proof facts. A visitor that recognizes only `MatchCall` while silently walking past `MatchJoin` can miss a host call, an effect, or a cycle. The current `requiredEffects` and parser match-kind switches make that integration risk concrete. An unsupported proof is therefore not reported as a failed program behavior or as a proved contract. [S4, S8, S10]

## D3. Test scripting, traces, and coverage

### Q3a — What is scripted, and where does the script live?

**Verdict: reuse request-aware `given` exchanges at foreign operation sites; do not introduce scripted join results or scripted native completion.**

| Operation under test | Evidence behavior |
| --- | --- |
| Same-file named Can task | Execute its real body, including nested joins and its reached foreign sites. No branch `given`. |
| Deterministic store or total kernel | Execute the operation. No fabricated result or timing script. |
| Cross-file Can task in ordinary unit evidence | Retain the existing foreign-call observation/honesty rules and branch-owned `given`; report it as scripted observation, not execution of the provider. |
| Completion-contracted host extern | Use its branch-owned `given` exchanges in decision tables. No live host, real sleeping, or real timer. |
| Real linked acceptance run | Execute every Can body. Only the explicitly selected host test boundary supplies host receipts; no fallback to a Can provider's `given`. |

Each selected item has the existing shape `exchange args (...) outcome ...`. Arguments are checked against the **already prepared invocation values**. Success and error payloads must have the declared shapes; errors must belong to the callee's `emits`. The join receives the completed outcome produced by this existing call-kind discipline and fills the corresponding slot. It never accepts a row saying “the whole pair is successful” without the branch calls occurring.

A cross-file script is counterfactual evidence at a declared boundary. It does not demonstrate the provider's termination, run its state mutations, or make its unreachable error real. This existing limitation is especially important after a join whose caller subsequently reads state. Such unit rows do not establish linked state behavior; D5 requires separate real-body linked/state/parity vectors. Static partition checking inspects real transitive authority even when a unit row substitutes an outcome.

**Checked falsification.** `evScriptOutcome` reads `node.Given[ctx.Test]`, consumes a per-site stream, validates the exchange request, and accepts declared error injection. `evLinkedOutcome` instead executes a real Can body without script fallback. Mixing those paths would turn injected outcomes into provider execution evidence. a09's store rule likewise rejects scripts where the store itself determines the result. [S8, S10, S15]

### Q3b — What does omission mean for an eager join?

**Verdict: omission still claims non-reach; it never means non-start, infinite work, success, or an implicit timeout.**

A join on an unselected outer path is not activated, so its sites need no entry for that root row. Once a join is activated, **every direct branch call is reached** under normal admitted execution, including siblings of a failed branch. Every reached foreign site therefore needs its applicable exchange. Missing a selected entry, exhausting its list, using the retired dash spelling, supplying a malformed request/outcome, or naming a nonexistent row fails through the relevant rule; missing an observation is not recovered by inventing `timeout`.

Partial tables remain partial. Unknown-key validation and the existing advisory treatment of statically stale valid keys remain distinct; B05 does not silently rewrite a91's diagnostics for join-free code. For rows that exercise a join, add an exact-consumption acceptance check for **activated** script streams: after the whole canonical root run, a reached stream may contain neither an unconsumed listed exchange nor a doubly consumed exchange. Do not check a shared stream for leftovers at the end of its first child; another logical child or later loop iteration can still consume it. A never-activated stream is not converted into a fabricated task.

This is an explicit new evidence obligation on join-exercising rows, not a claim that all legacy streams already enforce end-of-run balance. A row with missing/excess completion evidence is red and contributes no new join-arm coverage.

**Checked falsification.** a91 makes missing entries legal until reached and hard-bans `-`. The actual evaluator reports “no script for call” or “script exhausted”; neither denotes a task outcome. Skipping the right child because only a left-failure exchange was supplied would violate Q2a even though it could make the row look green. [S8, S15]

### Q3c — How are repeated sites deterministic without `order [...]`?

**Verdict: use a declaration-order reference interpreter and a branch-addressed observation ledger; do not let a racing task consume the next shared script item.**

The normal decision-table interpreter is a deterministic **reference execution**, not a real-time concurrency test. It prepares and registers all children, then evaluates their logical task trees in declaration order, depth first; each native child runs to completion before the next reference child. A scripted foreign receipt is available immediately to that reference run. All children still run before product dispatch. This is a legal linearization of the admitted noninterfering task semantics, not a promise that production effects occur in that total order.

The reference run keeps existing root-row resolution and per-call-site script-list order. It allocates each selected exchange to an internal event identity consisting of the logical task ancestry, static call site, and occurrence within that task. Thus a helper call site reached by two children consumes two entries in a deterministic **reference** order. It does not share a mutable cursor among concurrently executing production-test children.

For scheduler/parity tests, record those allocations as an immutable observation ledger. Then run a fresh execution/store with the same inputs, delivering the ledger's receipts by their event identities under another fair completion order. Request arguments must still match at the addressed site; a missing, extra, or mismatched event is a test failure. Only enabled events can be delivered. A changed schedule must not be allowed to allocate the script again in arrival order. The replay starts from fresh state, not the store already mutated by the reference run.

Example falsifier: the same foreign helper site is reached with `key = "left"` and `key = "right"`. Its row lists a left request/response followed by a right request/response. If the right task reaches it first in a shuffled runtime, a shared live cursor would either fail request matching or swap observations. Reference allocation plus event-addressed delivery prevents both without adding a new source row-key syntax. Request mismatch stays a failure; the harness must not search ahead for an item that happens to fit.

This is a required semantic model for the future test harness, not a prototype or a claim that current `Ctx.Scripts` supports it. Join-free source retains its current single execution order. Existing local-helper and loop stream behavior is retained in reference execution.

**Production guarantee.** Branch registration and argument preparation have declaration order. Within a branch, ordinary execution has causal order. Every terminal event precedes product dispatch. Sibling terminal events are otherwise unordered. A report may sort independent events by branch identity for stable display, but that display order must be labelled canonical rather than recorded chronological order. If timing or ordering is observed through host calls, that information is an external observation; B05 does not promise the same observation under every real-world schedule.

For fixed, correctly addressed observations and admitted noninterfering native effects, different fair interleavings preserve the product tags/payloads, selected arm, and final native store. The argument is that steps on disjoint writable authority commute, shared native reads have no concurrent writer, and completed input values are immutable. It does not prove host liveness or correctness of an omitted resource protocol. Several permutation tests are useful counterexamples, not an exhaustive proof of all schedules.

**Checked falsification.** The current evaluator keys scripts by node and test, which is correct for sequential evaluation but cannot simply be shared across racing children. The audit explicitly identifies new trace semantics as a cost and never defines what its illustrative `order` orders. Reinterpreting that word as completion order would either serialize work or lie about timing. A canonical reference order plus an explicit causal production trace preserves request-aware scripting without promising a global effect order. [S1, S15]

### Q3d — What counts as join coverage or evidence of liveness?

**Verdict: actual completed dispatch in a green row covers a join arm; scripts and schedule permutations prove no liveness.**

Credit a join arm only when every child has supplied a valid terminal outcome, the tuple selects that arm, and the enclosing row completes successfully with balanced activated streams. A registered child, a scripted exchange merely present in a table, an arm that was never selected, or an invalid host receipt earns no join credit. Retain the existing policy for suppressing coverage cascades on red rows and for distinguishing executed, authorized-certified, and uncovered arms elsewhere.

Native local execution retains the repository's existing local-helper evidence rules. Imported scripted providers and linked/parity runs acquire no new provider-arm credit merely because they sit beneath a join. Pure Fn execution retains B00's evidence boundaries. B05 introduces no aggregate relay certificate and no indirect way to create one.

A terminal timeout row means: **given a conforming host receipt with this declared outcome, this program took this path**. It does not prove that a timer fired, a worker stopped, a remote request was undone, or any host implementation meets its bound. Those are separate host-conformance obligations. Likewise a native task that returns promptly in a row still needs Q2d's static admission, and a resource limit exceeded during a row is not a native timeout outcome.

**Checked falsification.** a03's current amendment prohibits fabricated execution hits and limits the structural exception; a81 excludes scripts from proof facts. Removing the last actual witness for a product arm must therefore fail even if its tuple remains listed in scripts or a solver has a theorem about a different function. A row that returns `net.timeout()` immediately cannot discharge a 2-second production contract. [S3, S8]

## D4. Host boundary and TypeScript lowering

### Q4a — What is the completion contract on a shared extern?

**Verdict: a revision-bound trusted contract with an operation cutoff, a final settlement bound, and one explicit declared timeout kind.**

This is proposed provider syntax, not an addition to `std/host` and not a claim that a network adapter exists:

```can
mod fetchhost
  provides [net__fetch, Fetch__Value]
  uses []
  emits [net.down, net.timeout]

error net.down()
error net.timeout()

type Fetch__Value rev 1 (
  body: str
)

extern net__fetch(key: str) -> Fetch__Value rev 2
  effects [Host__World.use]
  emits [net.down, net.timeout]
  completion trusted
    timeout_ms 2000
    settle_ms 2100
    timeout net.timeout
```

The consumer module for Q1c has `uses [net__fetch@2, Fetch__Value@1]`. The numbers are illustrative contract choices, not measurements or accepted revision history. No separate completion symbol is added to `provides`: the contract identity is the declaring module identity plus **the extern name and exact revision**, including its metadata in revision identity. A changed bound or settlement promise is not a same-revision implementation detail.

Canonical well-formedness rules:

- `timeout_ms T` and `settle_ms S` are integer literals in milliseconds, with `0 < T <= S`. No floats, inferred units, clamping, dynamic deadline expressions, or resource tokens. A host binding must reject an unsupported exact range rather than silently round or narrow it.
- `timeout e.kind` names a declared, zero-payload member of this extern's `emits`. It is not an extra hidden outcome. B05 defines this one bounded terminal contract form; unbounded eventual-completion declarations are not admitted.
- `effects [Host__World.use]` is mandatory for the extern. The host has no Can-cell authority. Its parameters and success/error payloads are transitively Fn-free and remain subject to existing exact type/brand rules.
- A completion-contracted extern is initiated only as a direct join branch. Its callers can wrap it in ordinary Can functions, which remain ordinary functions at every further call site. Propagated host **authority** follows the existing effect discipline; there is no propagated completion return type or function modifier.

The trusted host contract starts when the host service accepts the immutable request. Before `T`, it may produce a normal declared outcome **only when the operation is already settled**. At `T`, if it has not done so, the timeout disposition wins; a normal result at exactly the cutoff does not change that disposition. After internal stopping/cleanup, the service must produce the declared timeout receipt no later than `S`. Every outcome, not only success, is terminal and settled. `S - T` is the explicit cleanup allowance.

**Settled** means that the operation and its descendants retain no authority to perform further effects within the contract's modeled boundary, mutate supplied/returned buffers, deliver a second completion, or reenter Can. A timeout is not a success-shaped indication that an abort request was merely sent. Earlier completed side effects are not rolled back. A mutating remote request does not gain a rollback, idempotency, or exactly-once guarantee from this spelling.

The host contract must be honest about remote activity: if late remote actions can still exercise authority included in the promised boundary, local socket closure alone cannot certify settlement. Such an operation is not admitted under this contract until a stronger protocol or isolation mechanism establishes the obligation. The illustrative fetch is therefore an interface proposal, not evidence that ordinary HTTP cancellation suffices.

**Checked falsification.** B02 provides exact pins and declaring-stem implementation ownership, not completion semantics. The current `ExternDecl` cannot represent these clauses, so they require an explicit future parser/checker extension. A one-number timeout that reports immediately while cleanup continues fails join ownership; separating `T` and `S` exposes rather than conceals that cost. Injecting `net.timeout()` in a row proves neither the bound nor settlement. [S1, S4, S8, S9, S16]

### Q4b — What is trusted, what is checked, and what happens on a broken contract?

**Verdict: check the boundary's declared shape and receipts; name the remaining trust and never convert violations into ordinary Can success.**

| Obligation | Owner |
| --- | --- |
| Correct exact revision, effect authority, timeout kind, bound shape, payload types, join placement | Can static checker. |
| One request identity and one terminal receipt; no slot overwrite; owned data transfer | Runtime and host adapter, with executable conformance tests. |
| No ordinary failure/rejection leaking outside declared operation outcomes | Adapter maps anticipated host failures to the declared set before submitting a receipt. |
| Completion, cleanup, absence of later effects, and the real bound under the selected platform conditions | Explicit trusted host contract; not a type theorem or decision-table proof. |
| No same-world reentry or concurrent root bypass | Checked runtime/embedding boundary. |
| No callable or mutable-reference alias smuggled through values | Static Fn containment checks plus host validation/ownership enforcement. |

The host binding's identity and ABI version belong in the build's evidence. The emitted import still comes from the **declaring** module's `.externs` stem, including when a different module calls it through a `uses` pin. Existing synchronous host exports are not silently changed into a different ABI or grandfathered as bounded tasks.

Incoming and outgoing values are validated/owned snapshots. In particular, the host must not retain a mutable `Uint8Array` alias and later alter a Can input or completed result. Brand authority is not weakened by transport; no new unseal or generalized byte-export right is created. Error payloads keep their declared identity and shape; no wrapper silently discards them.

A duplicate/malformed receipt, unexpected host rejection, reentry, or reported failure to settle is a **runtime contract failure**, not a new undeclared Can error tag. No normal root outcome is issued on that path. The world is failed and cannot admit another root while unretired authority might remain. The runtime may report the fault to its embedding and contain/terminate the world, but must not claim healthy join completion simply because it stopped looking at the worker. A failed environment falls outside the conditional completion theorem; B05 does not pretend it can force a dishonest host to finish.

Not every violation is dynamically detectable. A host falsely claiming that remote effects ceased can evade a local receipt validator. The contract is named **trusted** for that reason. Conformance tests, isolation, and receipt validation make the boundary accountable; they do not turn an external service's liveness into a Can proof.

**Checked falsification.** The current host source contains real synchronous observations and mutable host byte carriers, but no bounded worker-settlement protocol. Reusing a declared `Ok` shape alone cannot establish that an adapter is quiescent or that bytes remain immutable. B02's declaration-owner rule also defeats resolving a shared adapter against the caller's stem, which could select a different implementation of the promised contract. [S9, S15, S16]

### Q4c — Can a timeout be implemented as a timer racing an unretired worker?

**Verdict: no. A terminal receipt requires settlement; timer expiry alone is insufficient.**

The timer is an observation owned by the host contract. It fixes the timeout disposition; it does not independently close the task's outcome slot. The worker must stop or lose the relevant authority under a mechanism that actually establishes settlement, and the timeout receipt is accepted only as that settled outcome. Other sibling tasks continue to their own terminal outcomes. There is no first-timeout-wins join or automatic sibling cancellation.

A second pitfall is the scheduler itself. Running an arbitrarily expensive exact-integer operation on the same nonpreemptive event loop as a deadline timer can delay the timer and cleanup beyond the bound. Finite native computation does not imply bounded timer latency. Therefore a conforming bounded-host backend needs an independently progressing host/timer service or an explicitly adequate bounded-progress execution arrangement. A naive timer added to the current synchronous TypeScript emitter does not meet this contract.

The contract's `S` bounds production of the settled host receipt under the selected host platform's **explicit bounded-progress and clock assumptions**. It does not bound when a busy Can task consumes that receipt, how long unrelated native branches take, or when the entire root returns. Fairness alone is not a millisecond scheduling bound. Platform outage, resource exhaustion, and violation of the promised service progress are failures outside the trusted guarantee, not silently widened deadlines.

**Checked falsification.** A worker that writes after a timeout result falsifies the audit's joined-authority obligation. A long native operation that prevents the timer from running falsifies deriving a real-time bound from a08's finite-control proof. Neither problem is repaired by a happy-path timeout row. These are host/backend acceptance blockers, not deferred questions about source cancellation-token syntax. [S1, S6–S9, S15]

### Q4d — What does TypeScript lowering conceptually produce?

**Verdict: a uniform private machine ABI for a join-capable linked world, not two classes of Can functions or promise-returning Can APIs.**

Every Can function in such a world uses the same internal execution convention. Its source signature still denotes a named completed success record and declared errors. Ordinary calls remain ordinary outcome dispatch; a caller is not rewritten into a different source function category because a callee contains a join.

| Can construct | Conceptual emitted TypeScript machinery |
| --- | --- |
| Ordinary named function | A private entry/control frame with typed complete parameters, locals, and terminal result routing. |
| Ordinary `match call` | Push/enter the callee through that same convention and dispatch its terminal outcome at the recorded return site. |
| `match join` | Prepare values; allocate a typed child vector and fixed outcome slots; register all children; retain the parent control position until all slots are terminal. |
| Native child | Run checked Can control steps using its delegated store authority. No user callback is passed to the scheduler. |
| Bounded extern child | Submit a request with exact declaring-contract identity to the host service; receive a validated, owned terminal receipt at the private slot. |
| Product arm | Inspect the completed slot tags in declaration order and enter the statically checked arm with its typed payload bindings. |
| Root entry/exit | The embedding starts one world-owned invocation and receives one terminal result or a distinct runtime-contract fault report. No guest-visible unfinished value is returned. |

Frame pointers, return positions, event identities, and scheduler queues are backend control state, like an interpreter's stack, not Can data, Fn values, outcome payloads, or values available to reflection/serialization. Private runtime delivery machinery does not authorize source completion callbacks. A host implementation may use platform completion facilities behind the adapter, but their promise/handle types never appear in Can signatures or records.

This is a real backend/linkage cost. A join-capable world must not mix direct-return and machine-call modules by guessing which functions might perform host work. All Can modules in that linked world use its uniform ABI. A completely join-free build may keep the existing emitter and byte-identical goldens; that is a whole-artifact compatibility choice, not per-function source coloring. Cross-ABI interoperability is not silently promised; re-emit the checked world or use a separately specified embedding boundary.

TypeScript emission keeps exact bigint and decimal representations, declared success/error payloads, and `$can_kind`-compatible terminal discrimination. No first-class generic outcome feature, public task class, or generated source lambda API is added. Host imports retain B02's declaring-stem ownership. A scheduler is not an excuse to introduce implicit error propagation or catch-and-discard behavior.

`tsc` checks the resulting static shapes. Can remains responsible for exhaustiveness, evidence, task ownership, authority, and its native termination rules. Shape checking cannot prove exactly-once receipt handling, absence of late effects, scheduler fairness, or the semantic equivalence of the machine lowering; those need the distinct structural and runtime/parity gates in D5.

**Checked falsification.** The inspected current host exports return completed `$can_kind` unions synchronously. Simply substituting a host promise at one call site would break the consumer's current outcome inspection and propagate a different callable contract. A uniform internal machine convention avoids that per-function split but has the stated whole-world cost. a14 explicitly reserves behavioral proof to Can rather than TypeScript. No runnable lowering or runtime implementation has been produced here. [S3–S5, S9, S16, S17]

### Q4e — Which revision and compatibility boundaries change?

**Verdict: make concurrency-relevant metadata and executable dependencies revision-relevant; keep all changes prospective.**

A changed branch target, branch order, prepared arguments, effect delegation, outcome handling, specialization, or completion contract participates in executable identity. Host-boundary changes include timeout/settlement bounds, allowed authority, timeout error, value-ownership policy, and adapter ABI. A changed dependency must invalidate the join certificate even when the caller's textual body did not change.

No completion contract is inferred for an existing extern. Existing join-free programs retain their semantics and source forms. A program that introduces a join must satisfy the new checks and use the appropriate checked runtime boundary; adding the syntax alone does not upgrade a legacy host stub, certify cross-root safety, or make an existing pure linked runner support state/host work.

**Checked falsification.** B02's exact uses-pin discipline and B00's revision/certificate boundary would be undermined by caching admission solely on a caller's name or by changing a timeout promise at the same declared revision. The full current fingerprint implementation was not audited in this review; D5 makes inclusion and invalidation an explicit acceptance requirement rather than claiming it already works. [S3, S16]

## D5. Gate plan and acceptance

### Q5a — What evidence must a later implementation supply?

**Verdict: require independent structural, behavioral, host-boundary, and tooling falsifications; a parsed sample or a changed golden is not acceptance.**

Everything in this table is **planned evidence**, not an executed result. Existing mechanics are anchors to preserve; proposed tests are requirements, not assertions that test functions already exist. New diagnostic identities must be allocated after checking the live registry, with one independent rule per code and corresponding explanation/span tests. This document reserves no unverified `CAN` numbers.

| Gate | Positive evidence required | Mutation or counterexample that must fail |
| --- | --- | --- |
| **Parsing and nominal types** | Parse all four D1 sample shapes; local labels; generic named targets with explicit stamps; branch-owned exchanges; exact success records; zero-field and discarded payloads. Preserve the current ordinary function body and invocation forms. | Empty join; redundant native singleton; duplicate site/branch/binder; misplaced metadata or script; a branch with two calls or an inline body; wrong nominal result; Fn hidden through a nested input/result/error record; a label used as a value; every Q1d exclusion. |
| **Expansion and complete traversal** | Clone nested joins and scripts independently under generic expansion. Resolve branch dependencies and nominal return stamps. Every checker, printer, use collector, proof encoder, and emitter sees the new node or refuses it deliberately. | One generic stamp mutates another; a call referenced only from a join disappears from imports, effects, or cycle checking; stale admission survives a provider edit; an unhandled node is accepted by a default switch arm. |
| **Ownership and initiation** | Register all children; preserve once-only prepared values; complete each slot once; dispatch only after the last terminal receipt; nest and retire activation identities correctly. | Immediate left error suppresses right initiation; duplicate/late receipt overwrites a slot; a parent returns with an active descendant; a branch consumes a sibling result during preparation; recursion resets task accounting. |
| **Outcome product and payload checking** | Four real native combinations from Q1b and nine combinations from two independent three-outcome classifiers. Diagnose deterministic missing tuples and preserve full payloads, including equal error names in different slots. | Delete each tuple in turn; duplicate one; use a wrong error kind or slot arity; accept a whole-slot wildcard/OR/catch-all; narrow from test observations; silently truncate a large product when a proof limit is reached. |
| **Authority and root isolation** | Read/read sharing; disjoint write/read partitions; transitive local/imported summaries; nested subset delegation; persistent state between serialized roots; a host branch alongside disjoint native work. | Same-cell read/write or write/write via different wrappers; hidden transitive host use; two distinct extern names treated as independent host authority; host access to Can cells; same-world reentry; simultaneous roots against one store; automatic rollback on sibling error. |
| **Native termination and definedness** | Existing guarded unit, Euclid, and narrowing recursion remain admissible at branch calls. Total primitive subset and established exact-site certificates admit. Structural refusal blocks both CLI and LSP execution. | An indirect/cross-module cycle through a join; wrong guard or step; finite division-by-zero/index fault hidden behind passing rows; an unclassified kernel accepted as total; a script treated as proof of a real provider's termination. |
| **Value and store semantics** | Preparation and capture of complete values exactly once; independent branch effects; own-write visibility; post-barrier parent reads; no writable alias in complete record/byte results. | Re-evaluate an input at completion; reset cells at every join; lose a completed write on another branch's error; read a sibling's concurrently written cell; retain a mutable host alias that changes a completed value. |
| **Scripts and deterministic replay** | Preserve local/store execution and foreign exchange checks. Partial keys work on unselected paths. Repeated shared helper sites allocate observations in reference order and replay by logical event identity from fresh state. | Omitted sibling exchange suppresses a task; missing script manufactures timeout; activated stream is exhausted or left over; a racing child consumes another request's exchange; replay searches ahead for a matching request; replay reuses already-mutated reference state. |
| **Coverage and evidence categories** | Actual green-row dispatch covers each product arm. Existing local-helper, B00, and authorized relay boundaries remain intact. Reports distinguish executed, certified, uncovered, scripted, and linked evidence. | A listed but unconsumed exchange creates a hit; a red row credits a join arm; provider or linked runs launder coverage; correlated unreachable product arms are silently certified; a scripted timeout claims liveness. |
| **Host contract shape and conformance** | Exact `uses` pins, declaration ownership, full typed receipts, normal settlement before `T`, timeout at the cutoff, cleanup within `S`, and strict ownership of inputs/results. Bind a real adapter only with its separate conformance evidence. | Unknown timeout kind; invalid bounds or range narrowing; no contract; unexpected host rejection; duplicate result; result at `T` overriding timeout disposition; receipt before cleanup; late effect/buffer mutation/reentry; service cannot progress while native work runs. |
| **Golden TypeScript and shape gate** | New isolated join goldens show the uniform machine convention, declaration-order slots, exact terminal unions, private control state, and declaring-stem adapters. Join-free goldens stay unchanged. The repository-pinned TypeScript gate checks all affected outputs and adapter fixtures. | An old direct-return module is called through a guessed ABI; a Can-facing signature exposes unfinished work; implicit failure propagation; incorrect discriminant/record narrowing; an `any`/assertion suppresses the mismatch; unrelated golden edits hide drift. |
| **Real linked execution** | A separate admitted structured runner executes native providers, nested joins, stateful helpers, and ordinary calls to join-containing functions. Host behavior is isolated behind a controlled host boundary. No script fallback for Can providers. | Treat the existing linked-pure runner as already state/host-capable; return a provider outcome without entering the caller's arm; substitute a fake provider result to satisfy a state assertion; leak linked coverage into ordinary proof credit. |
| **Go / TypeScript parity and scheduler variation** | Same input vectors, addressed observations, results, payloads, final stores, and causal traces across the reference model and generated runtime. Include large bigint, exact decimal, owned bytes, nested activations, timeout settlement, and both terminal orders of independent branches. | Arrival-order failure selection; swapped payloads; host-number conversion; shared mutable buffers; wrong post-join state; a missing runtime silently skips the gate. Test runtime contract failures separately from ordinary Can outcomes. |
| **modcheck** | Branch-only imported calls/types are counted under ordinary provider/revision rules; the new fixture is in a scanned location. No local task/site label appears in a module header. | Wrong/missing provider pin; wrong revision; a join-only dependency is treated as unused; a misplaced fixture is skipped and reported green. |
| **gramcheck and formatting** | New join/branch/product/host-contract tokens receive the intended scopes without stealing strings, dotted errors, generic heads, ordinary value products, or current invocation syntax. Normalize preserves ownership and slot order. | Treat a regex smoke test as a proof of full parsing; accidentally admit value-pattern alternatives in outcome slots; lose a branch `given` or reorder branches during formatting. |
| **LSP and diagnostics** | Same admission as CLI; useful source spans for wrong slots, effects, targets, and contracts; UTF-16 correctness; edited providers invalidate certificates; incomplete worlds never run work. | A partial edit causes host execution; old-world certificates persist; a generic synthetic name leaks as the only explanation; one missing effect produces cascades instead of its real source diagnostic. |
| **Lint and proposed rewrites** | Positive fixtures have no unwanted findings **and no skipped files**. Necessary branch/type/effect annotations are understood. Ordinary constructions retain their normal lint rules. | Lint exits successfully after skipping unparsed joins; a chain/forward rewrite changes an all-child join into early exit; a rewrite invents callbacks or drops a product case. |
| **Revision identity and diagnostics exposure** | Branch order/targets/effects and all completion metadata affect the proper executable identity; wrong pins fail before execution; reports identify the trust boundary without exposing secret data. | A bound changes at the same accepted revision; host ABI changes retain a cached certificate; a completion-mismatch report prints secrets or private value representations; same-basename modules share an accidental identity. |

Retain the existing root-level gates when implementation is eventually authorized:

```sh
go test -count=1 ./...
go run ./tools/modcheck
go run ./tools/gramcheck
```

Use the repository's pinned `tscheck` / TypeScript workflow for the shape gate; an unrelated globally installed checker is not a substitute. The new unit, linked, host-conformance, parity, LSP, and lint checks must actually be attached to an executed gate before their names count as evidence. Do not run a hypothetical probe here, alter an expected output to make it pass, or report the repository's historical green statements as this review's own execution.

**Checked falsification.** B00 §8 identifies lint parse skips, linked evidence isolation, whole-world admission, and missing-runtime parity checks as real acceptance hazards. The inspected current evaluator explicitly excludes state/externs from its linked-pure admission and shares site-script cursors in ordinary execution. Thus reusing the existing gate names alone would not test the new semantics. a14 independently limits TypeScript's role to shape. [S3, S8, S15, S17]

### Q5b — How does every audit obligation map to a verdict?

**Verdict: none is silently dropped; partial scope is named at the affected obligation.**

The source is `docs/a/ASTRA_AUDIT.md`, **Task 3 §4**, plus the **Task 4 Async** row. The motivation, proof obligations, claimed benefits, complications, and liveness wording are all accounted for below.

| Audit statement or consequence | B05 ruling / evidence obligation | Status within this design |
| --- | --- | --- |
| Async is deferred; current state semantics are single-threaded. | Q0b keeps implementation priority; Q2c/Q2e supply the new interference and world-admission rules rather than assuming existing store safety. | Addressed prospectively; nothing implemented. |
| Every spawned task is joined exactly once. | Q1d prohibits independent task creation; Q2a requires registration, one terminal slot per child, all-terminal dispatch, and recursive retirement. | Addressed structurally, conditional on runtime/host conformance. |
| Outcome product is exhaustive. | Q2b uses the complete declared tag product and rejects missing/duplicate/abstract catch-all tuples. | Addressed; no new proof abstraction. |
| Branch authority is explicitly partitioned. | Mandatory branch effects; Q2c resolves canonical identities and checks complete transitive summaries. | Addressed for native cells and one conservative host domain. |
| Each branch's writes are disjoint from other reads and writes. | Q2c states both pairwise intersections; Q2e excludes hidden same-world concurrent roots and reentry. | Addressed; read/read allowed. |
| A separately specified protocol could allow sharing. | No such exemption exists in B05. More precise host partition and safe-sharing protocols belong to Resources. | Explicitly deferred; potentially overlapping work is rejected, not assumed safe. |
| Two three-outcome operations imply nine cases. | Q2b enumerates them and records the exponential and correlated-unreachable costs. | Addressed honestly; not solved by a hidden fallback. |
| No implicit first error wins, forgotten/detached work, or scheduler-dependent failure selection. | Q2a collects all outcomes; Q2b requires the caller's policy for each tuple; Q3c addresses observation allocation and completion-order independence under stated conditions. | Addressed. Real host observations are not claimed schedule-invariant. |
| New trace semantics are required. | Q3c defines causal order, canonical reference order, stable logical identities, and fresh-state event-addressed replay. | Addressed as a specification; future trace/parity gates required. |
| Join goldens and task result types are required. | D5 requires isolated goldens and exact internal terminal-slot types. Q1a/Q4d explicitly keep task-control/result-product carriers private rather than inventing source pending or generic outcome values. | Addressed by adopting the need for typed internal results, not a first-class task type. |
| Concurrent state access needs interference rules. | Q2c/Q2e define delegation, root exclusion, immutable values, persistent stores, and absence of automatic rollback. | Addressed. Arbitrary cross-cell invariants are not newly proved. |
| A join type does not prove an external service terminates. | Q4a–Q4c mark completion and bounded service progress as trusted, with checked receipt shape and settlement obligations. | Addressed; external liveness is not advertised as a Can proof. |
| Start with terminating Can tasks. | Q2d requires the existing whole-program ranking/graph rules and explicit definedness of task primitives. | Addressed; neither finite-faulting tasks nor test-only termination claims are admitted. |
| Host work needs trusted completion or timeout. | Q4a defines one bounded terminal contract with declared timeout and separate settlement bound; unknown/unbounded host leaves reject. | Addressed at the contract level; no adapter is certified by this review. |
| Never smuggle liveness assumptions into proof-carrying claims. | Q2d/Q3d/Q4b separately report native progress assumptions, trusted host progress, runtime failures, and observed evidence. | Addressed. |
| Every task has an explicit liveness assumption or bound. | Native tasks: checked finite control/definedness plus explicitly named progress/resource assumptions. Host leaves: the named `T/S` contract and platform bounded-progress assumptions. | Addressed at that strength. A numerical native CPU/memory budget or whole-root real-time bound is explicitly deferred to resource budgets, not claimed from this design. |

**Checked falsification.** A checklist restricted to the audit's four headline proof obligations would omit trace semantics, the nine-case cost, the internal result representation, and the explicit host caveat. Comparing the complete Task 3 §4 passage and Task 4 Async row against this matrix exposes those otherwise easy omissions. The deferrals above have rejection consequences; none authorizes unchecked work. [S1]

### Q5c — What remains risky or deliberately outside this approval?

**Verdict: approve or reject the bounded semantics as a whole; do not mistake the following exclusions for already available capabilities.**

| Risk / excluded capability | Consequence and reopening condition |
| --- | --- |
| Product explosion and correlated unreachable tuples | Useful safe computations can be unconsumable under the evidence law. No measured reduction is claimed. A wider relational or coverage certificate needs an independent design and falsification, not invented hits. |
| One host authority and serialized roots | Host fan-out and same-world throughput are deliberately restricted. HTTP/SQL/UI planning can use this seam but cannot declare its resource/concurrency requirements finished. Reopen with concrete resource identities, authority transport, aliasing, and safe-sharing protocols. |
| Bounded host progress and cleanup | The existing synchronous host shelf is not a bounded adapter. A concrete binding and independently progressing runtime must establish the named assumptions before host admission is enabled. Timer-only completion is refused. |
| Numerical native budgets and end-to-end deadlines | Finite control can still consume excessive CPU, memory, or task frames. Runtime limits remain failures; a mathematical resource-budget surface is separate work. No timing guarantee is inferred from a passing row. |
| Universal postconditions over joins/state | Structural safety is not a universal relational contract. The current contract verifier must refuse unsupported joins until a specific summary/composition design exists. |
| Full compiler/runtime integration | The complete fingerprint, parser visitors, formatter, target/import planner, TypeScript runtime, and diagnostic/redaction paths were not exhaustively audited. Their named gate obligations are release blockers, not assumed mechanisms. |
| API/ABI migration | A join-capable world needs a uniform runtime convention and a controlled embedding boundary. Compatibility with separately emitted legacy modules is not promised without re-emission or a separate checked bridge. |
| Resource handles, user cancellation, deadlines as values, races, streams, and channels | Excluded. B05's internal cleanup obligations do not create source operations for any of them. Reopen through Resources or a dedicated concurrency extension. |

**Approval boundary.** The reviewable artifact is D1–D5 with the obligations and exclusions above. Approval would settle the bounded design; it would not authorize code changes, migrations, gate execution on imaginary features, or an implementation sequence in this document. A separate implementation-planning task follows design approval.

**Checked falsification.** a05 requires a design before a build and refuses feature claims whose proof is merely future work. Here the bounded admitted subset has stated rules; omitted subsets remain rejected or outside its theorem. `ASTRA_STDLIB.md` expressly treats catalogue entries as proposed contracts, and B00 separately records unmeasured customer and integration risks. The approval boundary therefore does not mark a roadmap item shipped or erase the Resources dependency. [S3, S6, S12–S14]

## Source ledger and inspection limits

All repository references below refer to `8fccbb79aee193ef6f7c64ea09e441acb87e16fd`. The GitHub connector was used for source reads after a direct clone failed to resolve GitHub. The document and its additive patch are local deliverables; no successful checkout, remote commit, or repository write is claimed.

| ID | Inspected source and relevant location |
| --- | --- |
| S1 | `docs/a/ASTRA_AUDIT.md`, Task 3 §4 “Structured async with exhaustive joins and authority partitioning,” and Task 4's Async question/obligation. Focused source reads covered the full join passage and its immediately following implications. |
| S2 | The user-supplied async-design brief accompanying this review: Goal, Success Criteria, Constraints and Non-goals, D1–D5, Validation Plan, and Resources split question. This is a task constraint, not an implementation fact. |
| S3 | `docs/b00-function-values-design.md`, §0–§1, §8 gate/acceptance table, §9 boundaries, examined/omitted scope and Unresolved. Its older evidence revision and parser sketch are not treated as current code. |
| S4 | `compiler/parse.go`, AST and declaration definitions including `MatchKind`, `Node`, `FnDecl`, `ExternDecl`, `StateDecl`, and canonical module identity. |
| S5 | `compiler/parse.go`, function-body separator rejection, `isMetaHead`, `parseValuePattern`, `parseExprBlock`, `parseInvokeHead`, `parseScrutList`, `parseMatchArms`, and `parseMatchArmsKind`; generic/comma splitting was also inspected. |
| S6 | `docs/a/a05-expressiveness.md`, “Rule of the road,” “Explicitly later,” and Consequences; `docs/a/a08-termination.md`, proof-gated admission, amended halting argument, script-list behavior, and resource/host limits. |
| S7 | `docs/a/a11-recursion.md`, whole-program graph rule and guarded unit recursion; `docs/a/a19-recursion-schemas.md`, Rule, Theorems, and Cost note for unit/Euclid/narrowing. |
| S8 | `docs/a/a03-branch-coverage.md`, observed hit semantics and retrospective amendment; `docs/a/a81-proof-obligations.md`, Obligations and script/proof separation; `docs/a/a91-given-omission-amendment.md`, S0 rulings, runtime omission behavior, and stale-script follow-up. |
| S9 | `std/host/host.can`, declarations and clock/random/hash/secret/environment call evidence; `std/host/host.externs.ts`, synchronous implementations and actual bigint/byte/terminal-union boundary. The remaining host module tail was not needed for the absence-claim correction. |
| S10 | `docs/a/a09-effects.md`, Rule, Shape, Isolation argument, Static rules, and production boundary; `compiler/check.go`, `requiredEffects`, `checkEffects`, state declaration validation, and request-aware stub checking comments. |
| S11 | `docs/b01-type-syntax-redesign.md`, proposed status, existing/proposed surfaces, and sequencing warning. No claim that its contradictory draft details override the current parser. |
| S12 | `docs/ASTRA_STDLIB.md`, opening program statement, “Reading the catalogue,” signature notation, and Dependency labels. The whole catalogue was not re-audited. |
| S13 | `docs/a/a13-stdlib.md`, linkage evidence boundary, named-success-record convention, termination/cost distinction, and HTTP/SQL/UI deferrals/adoption order. Historical shipment statements are not new gate results. |
| S14 | `docs/stdlib-remaining.md`, host-shelf implementation record and HTTP/SQL/UI blockers. Its historical green assertions were read, not independently reproduced. |
| S15 | `compiler/eval.go`, `Value`, exact arithmetic including the division-by-zero failure, `evLinkedOutcome`, `evScriptOutcome`, and the corresponding call-kind dispatch excerpt. This was not an exhaustive evaluator/runtime audit. |
| S16 | `docs/b02-shared-externs-design.md`, admission, declaring-stem ownership, Semantics, Non-goals, and Falsification. |
| S17 | `docs/a/a14-tsc.md`, reinstatement, ownership of shape versus behavior, and checker-upgrade note. No TypeScript workflow was executed. |

**Unperformed work:** no compiler build, Go tests, modcheck, gramcheck, lint run, LSP session, solver run, Node/TypeScript execution, linked test, emitted golden, prototype, benchmark, real host timeout experiment, or complete namespace/fingerprint/redaction audit. No file outside this design document is proposed for modification by the accompanying patch. The illustrative `.can` blocks are design specimens only and live solely in this Markdown document.
