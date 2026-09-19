# b05 — structured joins and host completion contracts (design)

Status: design only. No compiler, grammar, stdlib, sketch, tool, or generated
file is changed or authorized by this document.

Repository snapshot inspected: `main@8fccbb79aee193ef6f7c64ea09e441acb87e16fd`
(2026-09-19). `docs/b00` through `docs/b04` exist at this snapshot, so `b05`
is the next free b-series number.

No gates were run for this design. File/line references below name the inspected
snapshot. “Checked falsification” means the proposed rule was traced against the
cited source and existing language rule; it does not mean the proposed feature
was compiled.

---

## Q0 — What is the feature, and what is not being reopened?

### Q0a — The concurrency surface is one lexical join form

**Verdict: add one structured `match join` form. A task can exist only as a
lexically declared branch of that form; no source value represents work that
has not completed.**

A join starts all of its branches, waits for every branch to produce one ordinary
Can outcome, then dispatches once over the product of those outcomes. The branch
outcomes never become a record, sequence, function value, task handle, or other
passable value. The only values that leave the join are ordinary payload binders
inside the selected arm.

This adopts the direction of `docs/a/ASTRA_AUDIT.md` §4 but not its sample
syntax wholesale. In particular, B05 keeps the useful idea — a match variant
whose arms cover a product of call outcomes — and re-derives the rest from the
current parser, effect system, evidence rules, and b00 boundary.

The design is deliberately ahead of implementation. `a05` says no expressive
build starts without its own spec document and still lists this area as later
work (`docs/a/a05-expressiveness.md:62-76`). The stdlib plan now has concrete
HTTP/SQL/UI demand, but those layers still depend on Resources as recorded in
`docs/stdlib-remaining.md:35-40`; B05 does not reorder those layers.

**Checked falsification.** `MatchKind` already names distinct match families and
`parseExprBlock` already recognizes special `match chain` / `match invoke`
heads before ordinary scrutinee parsing (`compiler/parse.go:166-194,
2621-2661`). A new `MatchJoin` family fits the existing control-flow split
without pretending that multiple call scrutinees are an ordinary value match.
`a28` explicitly excluded multi-call products because scripting and outcome
products needed their own design (`docs/a/a28-multi-scrutinee-match.md:83-86`).
B05 is that separate design rather than a silent relaxation of a28.

### Q0b — The rejected surface stays rejected

**Verdict: reject `async fn`, `await`, source-level function coloring, explicit
suspend points, `Task<T>`/future/promise source types, bare `spawn`, `detach`,
and callback-shaped completion.**

The rejection is semantic, not aesthetic:

- there is no passable representation of incomplete work, so nothing can be
  stored, returned, captured, compared, serialized, or forgotten;
- there is no second function type or modifier whose presence must propagate
  through signatures;
- there is no scheduler-visible source point inside an expression;
- there is no completion callback parameter, preserving b00's ban on host and
  effectful callbacks.

B00 deliberately limits function values to pure Can graphs and names host
callbacks / higher-order effectful graphs as separate work
(`docs/b00-function-values-design.md:34-37,529-532`). B05 does not widen `Fn`
inputs, returns, captures, extern positions, or invocation rules. A join branch
is a `call`, never an `invoke`.

**Checked falsification.** Treating a task as another `Fn` or record carrier
would immediately cross b00's transitive Fn-bearing position bans
(`docs/b00-function-values-design.md:111-114`) and would create a value that
could outlive the lexical join, defeating the audit's exactly-once obligation.
Keeping branches inside one AST node makes that illegal state unrepresentable.

### Q0c — Resources remain a separate design

**Verdict: keep Resources separate. B05 has no cancellation token, deadline
value, resource handle, acquisition scope, or protocol exception.**

The only time bound introduced here is a host-extern contract that references an
ordinary `int` parameter. It is metadata on the host boundary, not a new value
kind. Resource-scoped cancellation, shared connection ownership, and protocols
that permit otherwise-conflicting concurrent effects remain the Resources
workstream.

This split is intentional: B05 can prove lexical joining, outcome coverage, Can
state partitioning, and finite completion assumptions without inventing dynamic
ownership before the HTTP/SQL/resource model exists.

**Checked falsification.** The remaining-stdlib ledger says HTTP and SQL need
both this feature and Resources (`docs/stdlib-remaining.md:35-39`). Claiming
that B05 alone solves connection, transaction, stream, or cancellation lifetime
would contradict that dependency rather than remove it.

---

## Q1 — What is the canonical source syntax?

### Q1a — `match join` has labeled branches and labeled product arms

**Verdict: the canonical head is exactly `match join`; no descriptor follows
it. Each branch declares a unique label, its expected success record, its exact
Can-state authority, and one `call`. Arms qualify every outcome by branch
label.**

Canonical two-branch shape:

```can
match join
  left: Left__Value
    effects [Left__cache.read]
    call left__read(key)
  right: Right__Value
    effects [Right__cache.read]
    call right__read(key)
  on left.Ok l, right.Ok r =>
    Ok(left = l.value, right = r.value)
  on left.down _, right.Ok _ =>
    pair.left_failed()
  on left.Ok _, right.down _ =>
    pair.right_failed()
  on left.down _, right.down _ =>
    pair.both_failed()
```

The grammar contract is:

```text
join-head       := "match join"
join-branch     := LABEL ":" SUCCESS_TYPE
                   "effects [" EFFECTS? "]"
                   CALL_LINE
                   GIVEN_BLOCK?
join-arm        := "on" join-slot ("," join-slot)* "=>" rhs
join-slot       := LABEL ".Ok" binder
                 | LABEL "." ERROR_KIND binder
                 | LABEL "._"
binder          := IDENT | "_"
```

Rules:

1. Branch labels are unique identifiers local to that join.
2. Branch order is source order and is semantically significant only for
   argument evaluation and launch order; arm meaning is label-qualified.
3. `SUCCESS_TYPE` must exactly equal the called function/extern return record.
   It is an explicit cross-check, not an inferred alias.
4. `effects [...]` is mandatory, including `effects []` for a branch with no
   Can-state authority.
5. The branch body is exactly one `call` plus an optional ordinary `given`
   block when that callee kind is scripted today. No arbitrary expression,
   nested control block, or function-value invocation lives in the branch
   declaration.
6. Every arm has exactly one slot per branch, in branch declaration order.
   Labels make accidental reorder visible; permutations are rejected rather
   than normalized.
7. Join arms require the `on` prefix even though some existing match forms
   tolerate it being omitted. B05 gets one spelling.

The audit sketch's `Pair__Read` descriptor is rejected: no current checker,
coverage rule, test script, or emitter needs a nominal join identity. The join
site already has source identity, and branch labels provide the names required
for product diagnostics and scripts.

**Checked falsification.** Current `parseMatchArmsKind` already separates arm
parsing after a match family has been decided (`compiler/parse.go:2726-2788`).
Using label-qualified join slots avoids overloading ordinary `parsePattern`,
whose call patterns have no branch namespace. Reusing bare positional `Ok, Ok`
patterns would make branch reorder a meaning-changing edit with no textual
signal; the labels prevent that.

### Q1b — The audit's `order [...]` line is not admitted

**Verdict: branch declaration order fixes argument evaluation and launch order;
completion order is not observable in Can, so an `order [...]` clause would
create a second spelling for information the language does not expose.**

When a join is entered:

1. evaluate every branch call argument expression left-to-right by branch
   declaration and parameter-expression order;
2. only after every argument value exists successfully, launch each branch in
   branch declaration order;
3. collect every branch outcome;
4. perform one product-arm dispatch.

No branch is launched while later branch arguments are still being evaluated.
That avoids a partial-start state if a later argument expression fails loudly.
Once launch begins, language behavior cannot inspect which branch completes
first. The selected arm depends only on the final labeled outcome tuple.

There is therefore no scheduler race deciding “which error wins.” If policy
cares about the combination, the product table says so. If policy intentionally
prioritizes one outcome, that priority is written as overlapping source patterns
and ordinary first-match semantics, not inferred from completion timing.

**Checked falsification.** The audit's sample includes `order [left, right]`,
but §4's actual obligations are exactly-once joining, exhaustive outcome
products, and authority partitioning (`docs/a/ASTRA_AUDIT.md:519-528`). None
requires completion order to be observable. Existing value matches already use
source-ordered first-match dispatch; introducing a second order list would add a
staleness problem without strengthening those obligations.

### Q1c — Zero branches reject; one branch is legal

**Verdict: a join with zero branches is malformed. A one-branch join is legal
and is the canonical structured wrapper when a caller wants join semantics
around one completion-bounded host operation.**

Example host declaration and one-branch wrapper:

```can
error net.timeout()
error net.unavailable()

type Net__Request rev 1 (
  path: str
)

type Net__Response rev 1 (
  body: str
)

extern host__fetch(req: Net__Request, timeout_ms: int) -> Net__Response rev 1
  completes timeout_ms else net.timeout
  rejects net.unavailable
  emits [net.timeout, net.unavailable]

fn fetch__bounded(req: Net__Request, timeout_ms: int) -> Net__Response rev 1
  emits [net.timeout, net.unavailable]
  tests
    ok(Net__Request("/"), 1000) => Ok(body = "hi")
    timed(Net__Request("/"), 5) => net.timeout()
  match join
    request: Net__Response
      effects []
      call host__fetch(req, timeout_ms)
      given
        ok => [exchange args (req = Net__Request("/"), timeout_ms = 1000) outcome Ok(body = "hi")]
        timed => [exchange args (req = Net__Request("/"), timeout_ms = 5) outcome net.timeout()]
    on request.net.timeout _ => net.timeout()
    on request.net.unavailable _ => net.unavailable()
    on request.Ok r => Ok(body = r.body)
```

The one-branch form does not invent a pending value and does not grant early
return. It still waits for its declared branch and dispatches an ordinary
outcome. Its value is structural uniformity: the same form used for parallel
work can wrap the host completion contract before a later Resources design adds
cancellation or acquisition scopes.

**Checked falsification.** Rejecting one branch would force a second
host-completion construct even though the join product of one domain is just
that domain. Admitting zero branches would produce no call, no meaningful
outcome domain, and a match whose only purpose could be a hidden constant
expression; ordinary Can already expresses that without a concurrency form.

### Q1d — Nested joins compose only through ordinary RHS nesting

**Verdict: joins nest exactly where other expression blocks nest: in an arm RHS
or inside a called Can function. A branch declaration itself remains one call.**

Nested join after a successful outer product:

```can
match join
  account: Account__Value
    effects []
    call account__read(id)
  profile: Profile__Value
    effects []
    call profile__read(id)
  on account.missing _, profile._ => user.missing()
  on account.Ok a, profile.missing _ => user.profile_missing()
  on account.Ok a, profile.Ok p => match join
    avatar: Avatar__Value
      effects []
      call avatar__read(p.avatar_id)
    badges: Badge__List
      effects []
      call badges__read(a.id)
    on avatar.Ok av, badges.Ok bs => Ok(account = a, profile = p, avatar = av, badges = bs)
    on avatar.missing _, badges._ => user.avatar_missing()
    on avatar.Ok _, badges.unavailable _ => user.badges_unavailable()
```

Join inside an ordinary outcome arm:

```can
match call auth__prepare(input)
  on auth.invalid _ => auth.invalid()
  on Ok prepared => match join
    policy: Policy__Value
      effects []
      call policy__read(prepared.policy_id)
    subject: Subject__Value
      effects []
      call subject__read(prepared.subject_id)
    on policy.Ok p, subject.Ok s => Ok(policy = p, subject = s)
    on policy.missing _, subject._ => auth.policy_missing()
    on policy.Ok _, subject.missing _ => auth.subject_missing()
```

This keeps one expression grammar: arms may contain another expression block;
branches are not mini-functions.

**Checked falsification.** `parseExprBlock` and `parseRhs` already make nested
matches ordinary RHS structure. Allowing an arbitrary branch body would create
a second function-body grammar, obscure the branch's one task root, and make
its static authority set harder to compare with one declared callee.

---

## Q2 — How are outcome products checked?

### Q2a — Each branch contributes one finite outcome domain

**Verdict: a branch domain is exactly `{Ok} ∪ callee.emits`; the join domain is
the labeled Cartesian product of those branch domains.**

`Ok` carries the callee's declared success record. Each declared error kind
carries its existing error payload shape. The branch's explicit success type is
checked against the callee declaration; it does not change the domain.

An arm slot may be:

- `<label>.Ok x` — exact success, binding the success record;
- `<label>.<error-kind> e` — exact declared error, binding that error payload;
- the same exact pattern with `_` instead of a binder;
- `<label>._` — any outcome from that one branch, binding nothing.

No outcome OR-pattern is added in B05. `<label>._` is the one grouping device.
It is enough to express “the other branch does not matter” without introducing
binder-unification rules across different outcome payload types.

**Checked falsification.** Ordinary multi-value matching is already specified
as a product-pattern table (`docs/a/a28-multi-scrutinee-match.md:10-13`). Its
checker uses symbolic product-space subtraction instead of enumerating every
cell (`docs/a/a28-multi-scrutinee-match.md:127-130`). Extending that algebra
from bool/string slot domains to finite call-outcome slot domains is smaller
and more consistent than inventing a second exhaustiveness engine.

### Q2b — Exhaustiveness is symbolic; source policy is still explicit

**Verdict: prove join exhaustiveness with the a28 space-subtraction model,
generalized to finite labeled outcome sets. Do not construct the whole product
just to check it. Diagnose up to three concrete missing tuples.**

For each source arm in order:

1. translate every slot to a subset of its branch outcome domain;
2. take the product of those subsets;
3. subtract product space already covered by earlier arms;
4. add the useful remainder to `covered`;
5. after the last arm, require `total - covered` to be empty.

A witness renders with labels, for example:

```text
missing join outcome (left=net.down, right=Ok)
```

A wildcard arm may cover many product cells. That is allowed only because its
RHS explicitly discards that branch's outcome. Ordinary first-match behavior
settles intentional overlaps; an arm completely shadowed by earlier arms stays
dead evidence and is caught by the existing coverage law.

This contains checker complexity, not semantic complexity. If two branches each
have three outcomes and the program genuinely distinguishes all nine
combinations, the source still needs nine behaviors. B05 does not invent an
automatic “first failure,” “collect errors,” or “prefer left” abstraction. If
multiple cells intentionally share one behavior, `<label>._` can combine them.
A later abstraction must justify its own proof and payload rules.

**Checked falsification.** The audit explicitly calls out the 3×3 nine-case
cost and warns against hand-waving it (`docs/a/ASTRA_AUDIT.md:519-528`). Reusing
an already-shipped symbolic product checker prevents internal Cartesian blowup,
but it cannot honestly remove source cases whose outputs differ. The wildcard
only compresses cases with the same written behavior.

### Q2c — No early failure path exists

**Verdict: every launched branch reaches one terminal outcome before arm
dispatch, even when another branch has already produced an error.**

There is no implicit fail-fast behavior. A branch error is data for the product
match, not permission to abandon siblings. This is what makes the exactly-once
join structural and deterministic without cancellation machinery.

A future Resources design may add a checked cancellation protocol, but it must
still define what terminal outcome a cancelled branch contributes before the
join dispatches. B05 does not reserve an implicit cancellation outcome.

**Checked falsification.** An early return after the first error would make
some launched branch unjoined and would make selected failure depend on runtime
completion timing. Both contradict the audit's §4 obligations. There is no
current task value on which a safe “later join” could be recorded, and B05 does
not create one.

---

## Q3 — What enforces exactly-once joining and authority partitioning?

### Q3a — Exactly once is a syntax invariant, not a linear task-handle analysis

**Verdict: the language never exposes a task handle. Entering `match join`
creates all branch tasks and that same node consumes all of them once. There is
nothing for a separate linear checker to track.**

The checker therefore proves exactly once by validating the node shape:

- at least one branch;
- one root call per branch;
- unique labels;
- one product dispatch owned by that node;
- no branch reference outside the node;
- no syntax for returning, storing, forwarding, or manually joining a branch.

The implementation may internally assign branch IDs, but those IDs are not Can
values and never appear in types or emitted public contracts.

**Checked falsification.** A handle-based design would need new value positions,
lifetime rules, alias rules, and forgotten-handle diagnostics. B00 already has a
transitive “contains Fn” walk because first-class callables can travel through
records and sequences; repeating that complexity for task handles would be the
opposite of the audit's structured-only requirement. Lexical ownership removes
the alias problem instead of checking it after creation.

### Q3b — Branch `effects` is an exact local authority declaration

**Verdict: every branch must spell the exact transitive Can-state capability set
of its root callee. Missing capabilities reuse the existing undeclared-effect
rule; extras reuse the stale-effect rule.**

For a branch root `call f(...)`:

```text
branch.effects == effects(f)
```

where `effects(f)` is the same already-declared transitive capability set a09
requires callers to expose. The branch line is intentionally redundant with the
callee signature: it makes the concurrency partition visible at the join site
and makes a callee effect expansion break the join until reviewed.

The enclosing function still obeys the existing rule: its function-level
`effects [...]` must cover the union of all sequential and joined work. B05 does
not infer missing function authority from branch declarations.

**Checked falsification.** A09 defines `.read` / `.write` capabilities and says
callee effects are transitive rather than inferred
(`docs/a/a09-effects.md:56-68`). Reusing that set gives B05 a real authority
vocabulary. Omitting the branch line would force reviewers to chase every
callee just to see whether two branches conflict; accepting a superset would
let stale “authority” text pretend a partition was reviewed when it was not.

### Q3c — Concurrent Can-state conflicts are rejected pairwise

**Verdict: reads may overlap reads. Any write must be disjoint from every
sibling read or write of the same cell. B05 has no protocol exception.**

For every pair of branches `A`, `B`:

```text
writes(A) ∩ (reads(B) ∪ writes(B)) = ∅
writes(B) ∩ (reads(A) ∪ writes(A)) = ∅
```

Thus:

- `X.read` + `X.read` — legal;
- `X.write` + `X.read` — rejected;
- `X.read` + `X.write` — rejected;
- `X.write` + `X.write` — rejected;
- accesses to different cells — legal.

Each branch sees immutable argument values. Module cells remain shared runtime
state; the rule is what makes their concurrent interaction scheduler-independent.
After the join, writes to distinct cells are all visible to the selected arm.

A separately proved synchronization/resource protocol could relax this later.
B05 deliberately has no “unsafe,” lock, atomic, or protocol escape hatch.

**Checked falsification.** A09's test isolation argument is explicitly
single-threaded and calls out that concurrency ordering questions were outside
that design (`docs/a/a09-effects.md:69-90`). The audit supplies exactly this
write-vs-read/write partition rule for the new surface
(`docs/a/ASTRA_AUDIT.md:519-522`). Allowing same-cell write/read here would make
existing deterministic store reasoning false without any replacement proof.

### Q3d — Host-side sharing is trusted only through a narrow isolation contract

**Verdict: a multi-branch join may reach a host extern only when every such
extern declares `concurrency isolated`. A one-branch join does not need that
marker. Dynamic shared-resource concurrency remains blocked for Resources.**

`concurrency isolated` means: overlapping invocations of that extern do not
share mutable host authority whose interaction can change a Can-visible outcome
or externally ordered effect. It is a host conformance promise, not something
Can proves from TypeScript.

Examples that can plausibly carry it: independent stateless request adapters or
independent timers. An ordered logger, one mutable cursor, a transaction, a
socket stream, or a connection-scoped SQL operation must not claim it merely to
pass the checker; those need a resource/ownership model.

The rule is transitive at the join site. If a branch calls Can function `f` and
`f` eventually reaches host extern `h`, the multi-branch admission check sees
`h` in the reachable graph and requires the marker.

**Checked falsification.** A09 explicitly says extern bodies are host-side and
outside its state proof (`docs/a/a09-effects.md:56-68`). Pretending Can cell
capabilities also describe a database connection or logging sink would be a
false proof. Conversely, banning every host operation from a multi-branch join
would make the feature unusable for its eventual network customers. A narrow,
explicit trusted-isolation claim states the missing fact without inventing
resource handles early.

---

## Q4 — What is the liveness rule?

### Q4a — Existing Can termination is the branch proof for Can code

**Verdict: a branch rooted in Can code is join-admissible only if its reachable
Can graph satisfies the existing whole-program termination policy; B05 adds no
second recursion theorem.**

Can already rejects call-graph cycles except admitted direct recursion with a
proved decrease shape (`docs/a/a11-recursion.md`). For join purposes, finite Can
computation is therefore inherited. Nested joins do not weaken that proof: each
nested join must independently satisfy B05, and recursion around it still obeys
the existing decrease rule.

The join checker additionally walks reachable host externs because Can's
termination theorem does not prove host completion.

**Checked falsification.** A08/a11 distinguish source termination from runtime
resource bounds. Re-proving recursion inside B05 would fork the termination
policy. Treating “the test completed once” as liveness would also violate the
verification boundary: execution evidence is not a universal termination
argument.

### Q4b — Host participation requires an explicit completion bound

**Verdict: any host extern reachable from a join branch must declare a bounded
completion contract. The canonical extern metadata is:**

```can
extern host__fetch(req: Net__Request, timeout_ms: int) -> Net__Response rev 1
  completes timeout_ms else net.timeout
  rejects net.unavailable
  concurrency isolated
  emits [net.timeout, net.unavailable]
```

`completes P else E` means:

1. `P` names one declared `int` parameter of this extern;
2. a non-positive value means immediate `E` without trusting provider progress;
3. for a positive value, a conforming host adapter resolves to an ordinary
   declared outcome no later than that host-provided bound; if provider work has
   not produced one, the adapter resolves `E`;
4. `E` is a declared zero-field error kind present in this extern's `emits`;
5. the bound is a trusted host contract. Decision tables can exercise its
   result, but do not prove clock or scheduler behavior.

`rejects R` means a host throw/rejected completion is normalized to the declared
zero-field error `R`, also present in `emits`. A conforming host adapter exposes
no second language-level failure channel.

Both metadata lines are revision-significant extern contract material. A rev
bump is required when the bound parameter, timeout mapping, rejection mapping,
or isolation promise changes in a way that changes the contract.

An extern without `completes` remains legal for ordinary sequential Can use under
existing rules, but it is not join-admissible. B05 does not retroactively claim
liveness for old host code.

**Checked falsification.** The audit says host completion is not proved by the
join type and requires an explicit trusted completion or timeout contract
(`docs/a/ASTRA_AUDIT.md:525-528`). B02 already centralizes a shared extern's
contract and implementation at its declaring stem while preserving exact `uses`
pins (`docs/b02-shared-externs-design.md:39-60`). Putting completion metadata on
the extern extends the existing trust boundary instead of inventing per-call
host promises in Can source.

### Q4c — Host bounds do not become proof facts

**Verdict: the verifier may use the existence of an admitted completion
contract to establish the qualitative claim “this host edge is bounded,” but it
must not derive functional postconditions from scripted completions or claim
that tests proved the host clock bound.**

B05 does not attempt arithmetic composition of wall-clock bounds through call
graphs. A join can therefore be proved finite without advertising “returns in N
milliseconds” as a Can theorem. Numeric resource budgeting is a separate proof
feature.

A blocked host event loop or dishonest adapter is a host conformance failure.
The contract makes that trust explicit; it does not turn a JavaScript timer into
a theorem.

**Checked falsification.** A81 states that `given` scripts are never proof facts
(`docs/a/a81-proof-obligations.md:19-22`). Treating a scripted timeout row as
proof that a real provider honors the bound would cross exactly that boundary.

---

## Q5 — How do decision tables script joins?

### Q5a — Reuse branch-local `given`; do not add a second completion table

**Verdict: each join branch uses the ordinary `given` / `exchange` grammar when
its callee kind is scripted today. That exchange is the observed completion
outcome for that branch. There is no join-wide `complete`, schedule, or timing
block.**

Example:

```can
match join
  left: Left__Value
    effects []
    call host__left(id, timeout_ms)
    given
      both_ok => [exchange args (id = 7, timeout_ms = 100) outcome Ok(value = "L")]
      left_down => [exchange args (id = 7, timeout_ms = 100) outcome left.down()]
  right: Right__Value
    effects []
    call host__right(id, timeout_ms)
    given
      both_ok => [exchange args (id = 7, timeout_ms = 100) outcome Ok(value = "R")]
      left_down => [exchange args (id = 7, timeout_ms = 100) outcome Ok(value = "R")]
  on left.Ok l, right.Ok r => Ok(left = l.value, right = r.value)
  on left.down _, right.Ok _ => pair.left_failed()
```

The script says “for this test, this branch completion is observed.” It does not
say when the completion happens and does not certify that production will
complete. A timeout test scripts the declared timeout outcome like any other
foreign outcome.

Local deterministic Can callees keep executing their bodies under the current
unit semantics. Foreign Can/extern branches keep the current scripted boundary.
B05 changes where multiple such completions are collected, not the trust level
of a callee kind.

**Checked falsification.** A67 keeps hermetic `given` rows specifically so a
caller can exercise any permitted outcome independent of provider behavior
(`docs/a/a67-testing-contracts-design.md:17-23`). Adding a second completion
syntax for the same observed outcome would create two spellings and risk a
false distinction between “call result” and “task result.”

### Q5b — Omission keeps the a91 meaning

**Verdict: a missing branch `given` entry claims non-reach. If the join and that
scripted branch are reached, the test fails for missing script. Unknown keys
remain errors/lints exactly as today.**

A join reaches every declared branch once it reaches launch, so a reaching test
must provide an exchange for each branch whose callee requires scripting. No
static total table is reintroduced.

**Checked falsification.** A91 deliberately removed total `given` tables and
made omission the non-reach claim (`docs/a/a91-given-omission-amendment.md:14-17`).
Making joins total again would revive the boilerplate rule one day after it was
retired and would create inconsistent omission semantics between `match call`
and `match join`.

### Q5c — Completion order is not scripted because it has no Can meaning

**Verdict: the hermetic evaluator may obtain branch outcomes in declaration
order and then dispatch the product. It does not simulate elapsed time or
completion races.**

This is semantically sufficient because:

- all call arguments were fixed before launch;
- Can-state conflicts are statically rejected;
- multi-branch host externs carry the trusted isolation promise;
- no source expression can observe a branch before the full join;
- no arm is selected until every branch has an outcome.

A future feature that exposes races, first-completion selection, streaming, or
cancellation would need a new design and new evidence semantics. None is hidden
inside B05.

**Checked falsification.** If tests scripted completion order while source could
not observe it, the table would contain evidence that cannot affect behavior.
If source *could* observe it, the selected outcome could become scheduler
sensitive and the state-disjointness argument would no longer be sufficient.
B05 chooses the smaller semantic surface.

### Q5d — Join-arm coverage remains observed execution

**Verdict: CAN4107-style coverage applies to join arms unchanged: an arm counts
only when a passing committed row collects the full branch product, selects that
specific arm, and enters its RHS. Script presence is not coverage.**

A wildcard product arm is one arm for coverage purposes. The checker proves the
space it denotes; the evidence law proves that the written behavior is executed
at least once. B05 does not silently strengthen coverage to every product cell.
If a project wants cell-by-cell tests, that is a separate policy/tooling layer.

**Checked falsification.** A03 records coverage at actual arm entry, and a67
explicitly says script presence is not coverage (`docs/a/a03-branch-coverage.md:6-9,52-55`;
`docs/a/a67-testing-contracts-design.md:25-31`). Counting every branch exchange
as taking a join arm would let two unrelated scripts “cover” product behavior
that never ran.

---

## Q6 — What is the TypeScript boundary and lowering shape?

### Q6a — Can keeps one call contract; the target uses one uniform completion ABI

**Verdict: when B05 is eventually implemented, generated Can callables use one
uniform completion-capable TypeScript ABI. The source type system still has
only ordinary Can returns and declared error outcomes.**

The reason is structural. A JS target cannot run genuine host work concurrently
while preserving a synchronous result ABI for some transitive callers and a
different result ABI for others without creating two hidden call classes. B05
rejects that distinction at the Can level. The emitter therefore pays the cost
once: every generated Can callable participates in the same target completion
protocol, including code that resolves immediately.

A concrete lowering may use `Promise<Result>` for that internal/public generated
TS ABI. `Promise` never becomes a Can type and cannot appear in `.can` records,
parameters, `Fn`, errors, state, tests, or pattern binders.

This is a future compatibility cost, not a design-time behavior change. The B05
implementation plan must treat the ABI migration as an explicit emitter slice,
regenerate all affected goldens, update extern adapters, and run the strict TS
gate. It must not quietly infer two generated function classes.

**Checked falsification.** The current TS gate assigns type/contract shape to
TypeScript and keeps exhaustiveness, termination, decision tables, and coverage
in canlc (`docs/a/a70-tsc-gate.md:15-18`). That is the right ownership split:
TS can verify the uniform target ABI, but it cannot substitute for B05's source
proofs. The repository's current `tscheck/tsconfig.json` targets ES2020, where
`Promise` is available.

### Q6b — Join lowering evaluates arguments, launches all branches, then dispatches

**Verdict: lower a join in four phases matching Q1b. A representative target
shape is:**

```ts
const leftArg = /* evaluate Can argument */;
const rightArg = /* evaluate Can argument */;

const leftRun = left__read(leftArg);
const rightRun = right__read(rightArg);

return Promise.all([leftRun, rightRun]).then(([left, right]) => {
  // compiler-generated dispatch over left.$can_kind and right.$can_kind
  // with the source arm order preserved
});
```

The exact helper names are an implementation detail; the semantic requirements
are not:

- all argument values exist before the first branch launch;
- launch order is declaration order;
- every branch is collected;
- result dispatch is label/product based;
- no host promise is constructible or bindable in Can source;
- branch outcomes remain the existing Ok/error unions.

Compiler-generated target closures used by `.then` are emitter machinery, not
B00 `Fn` values. They cannot be supplied by source or host code through a Can
function parameter.

**Checked falsification.** Lowering branches one after another all the way to
completion would preserve results for pure work but would not implement the
concurrency feature. Returning after the first error would violate Q2c.
Launching a branch before later argument evaluation would introduce a partial
start state that Can has no recovery syntax for.

### Q6c — Completion-contracted extern adapters resolve ordinary outcome unions

**Verdict: a completion-contracted extern's generated/hand-written TypeScript
adapter participates in the same completion ABI and resolves only the declared
Ok/error union. The adapter owns timeout and rejection normalization required by
its metadata.**

B02's declaring-stem rule remains: a shared extern is implemented by the module
that declares it, and callers pin the revision normally
(`docs/b02-shared-externs-design.md:39-60`). B05 adds contract fields; it does
not move host implementation ownership to callers.

The strict TS gate should check the adapter's result shape, but host liveness and
isolation remain trusted runtime contracts. Dedicated host smoke tests can
falsify obvious nonconformance (for example, a deliberately never-resolving
provider must yield the declared timeout through the adapter), but those tests
do not become Can proof facts.

**Checked falsification.** Letting host code return a raw target promise into a
Can binder would introduce exactly the first-class incomplete value B05 bans.
Letting a rejected target promise bypass the declared error union would create a
new failure channel. The adapter boundary must normalize before Can sees a
result.

### Q6d — State after a join is defined only for admitted disjoint effects

**Verdict: the selected arm observes the store after all branch writes have
completed. Because conflicting same-cell access is rejected, that final store
is independent of branch interleaving.**

The unit evaluator may execute branch bodies in declaration order while still
modeling the language semantics correctly: under the admitted effect rule, any
other interleaving produces the same Can-visible store and outcome product.
This is a simulation convenience, not a source ordering guarantee.

**Checked falsification.** Without the pairwise conflict rule, a sequential unit
evaluator could accidentally pick one deterministic order while emitted host
execution exhibits another. A09's current fresh-store determinism would then be
insufficient evidence. Q3c removes that mismatch before execution.

---

## Q7 — How does B05 interact with existing language surfaces?

### Q7a — Ordinary `match call` remains the sequential one-call form

**Verdict: `match call` is unchanged. It consumes one call outcome before its
arms run. `match join` is used when multiple calls must be launched as one
structured group, or when a caller intentionally uses the one-branch structured
form.**

A completion-contracted extern may still be consumed by ordinary `match call`;
its host adapter produces the same declared outcome union. B05 does not create a
new call keyword.

**Checked falsification.** Requiring a new call spelling for every host
completion contract would spread the feature outside the one intended control
surface and force callers to know a target/backend property. The source already
has an explicit effectful call boundary.

### Q7b — Function values remain outside join branches

**Verdict: a B05 branch root is `call`, not `invoke`. B00 function values remain
pure, unary, and governed by their existing invocation match.**

A later design could decide whether pure `Fn` invocations are useful join roots.
B05 does not need them to satisfy HTTP/SQL/UI planning, and allowing them now
would reopen callback graph/lifetime questions the prompt explicitly excludes.

**Checked falsification.** B00 names higher-order/effectful functions as a
follow-up and rejects Fn-bearing values across numerous effectful positions
(`docs/b00-function-values-design.md:529-532`). Using `invoke` as a task branch
would blur that boundary for no required B05 customer.

### Q7c — Linked-pure integration may admit only pure joins

**Verdict: the future linked runner may execute a join only when every reachable
branch remains inside its existing pure/deterministic graph. A join reaching an
extern or state operation remains refused by `runLinkedPure`.**

For a pure join, linked execution can evaluate all real branch bodies, collect
their outcomes, and dispatch; it still must not let linked traces satisfy local
CAN4107 coverage.

**Checked falsification.** `checkLinkedGraph` currently rejects externs, state
ops, and effectful functions and runs real pure bodies only
(`compiler/linked.go`). Weakening that just because the call is under a join
would turn an integration runner into an effects runner without a trust model.

### Q7d — Revision identity must include the new host contract fields

**Verdict: extern `completes`, `rejects`, and `concurrency` metadata are part of
revision identity and exact-rev compatibility. Join source structure itself is
body behavior and follows the existing function revision/body rules.**

A caller pinned to a host extern revision must not silently receive a different
timeout mapping or isolation promise under the same rev.

**Checked falsification.** B02 relies on exact revision pins for shared externs.
Leaving the new metadata out of canonical identity would let the very contracts
that justify join admission drift without invalidating a pin.

---

## Q8 — What diagnostics are required?

### Q8a — Reuse existing rule families where the semantic rule is already the same

**Verdict: do not allocate a new code merely because a rule occurs under a
join. Reuse existing families for known rules; add new codes only for genuinely
new join obligations. Exact numeric assignments wait for the implementation
registry audit.**

Planned mapping:

| Condition | Diagnostic family |
| --- | --- |
| branch call arg/type mismatch | existing call/type diagnostics |
| branch success annotation != callee return | existing type mismatch family |
| missing branch effect | existing undeclared-effect rule |
| extra branch effect | existing stale-effect rule |
| bad arm arity/label/slot shape | existing match-arm shape family |
| unknown branch outcome | existing call-outcome/arm family |
| uncovered product | existing missing-arm family, with labeled tuple witness |
| untaken join arm | existing CAN4107 coverage rule |
| zero branches | parse error |
| conflicting sibling effects | **new join-conflict code** |
| join reaches extern without completion contract | **new join-liveness code** |
| multi-branch join reaches extern without isolation contract | **new join-host-isolation code** |

One rule gets one code. The three new rules are independent and should not be
collapsed into one “bad join” diagnostic.

**Checked falsification.** A28 reused existing match codes for product value
matches rather than creating tuple-specific duplicates. A09 likewise separates
missing effect authority from stale authority. B05 follows those precedents and
adds codes only where no existing rule states the failure.

---

## Q9 — What is the implementation gate plan?

### Q9a — Acceptance is layered; parser success is never enough

**Verdict: implementation, when separately approved, ships in slices whose
failures are independently observable. This document authorizes none of them.**

| Gate / slice | Required positive evidence | Required falsification |
| --- | --- | --- |
| Parse / AST | `MatchJoin` head; one and multiple branches; mandatory label/type/effects/call; branch-local `given`; labeled arms; nested RHS forms; extern completion metadata | reject descriptor after `match join`; reject `order [...]`; zero branches; duplicate labels; missing/duplicate effects or call; branch `invoke`; malformed labeled slots; completion metadata naming a non-`int`/unknown parameter |
| Type checking | branch args type-check; success annotation equals callee return; Ok/error binders get existing payload types | wrong success record; unknown emitted error; binder misuse; Fn-bearing/task-like invented positions remain rejected |
| Static admission | exact branch effects; read/read sharing; disjoint writes; transitive join-reachable graph accepted when finite | write/read and write/write same-cell conflicts; join-reachable extern without `completes`; multi-branch reachable extern without `concurrency isolated`; stale/missing branch effects |
| Exhaustiveness | 2×2 and 3×3 outcome domains; wildcard compression; up to three labeled witnesses; symbolic checker does not materialize full product | missing one exact tuple; shadowed wildcard/exact arm later fails evidence; unknown label/permuted arm labels reject |
| Unit execution / scripts | local Can branches execute; scripted branches consume ordinary exchanges; nested joins; one-branch timeout outcome | omitted reaching branch script fails at execution; unknown/stale keys keep a91 behavior; script never satisfies liveness proof |
| Coverage | each committed join arm taken by at least one passing row | one never-selected join arm reports existing coverage failure even if all branch `given` entries exist |
| Termination / liveness | proved recursive Can branch accepted; completion-contracted host branch admitted | unproved cycle still fails before tests; old extern without completion metadata refused only when join-reachable |
| Emit golden | all args evaluated before launches; launches in branch order; collection of every branch; product dispatch; nested join | mutant sequentializes full branches; mutant returns on first error; mutant starts first branch before later args; all must fail golden/behavior probes |
| Host adapter | timeout and rejected host completion normalize to declared zero-field errors; shared extern still resolves from declaring stem | wrong-stem import; raw host completion leaks into Can result; undeclared timeout/rejection kind; dishonest isolated fixture caught by targeted smoke where mechanically testable |
| Linked | pure join executes real branch bodies under explicit linked test | extern/state/effectful join remains refused; linked trace does not count for local arm coverage |
| Revision identity / parity | extern contract metadata changes canonical revision material consistently across CLI/editor | changing timeout/rejection/isolation metadata without rev produces revision-identity failure rather than silent drift |
| LSP | new join-specific codes block consistently with CLI; spans point at branch/effect/extern contract | no duplicate cascade for one root cause; malformed branch reports at branch site |
| Lint | stale-script reachability walks branch calls; existing lint behavior preserved elsewhere | unreachable join-branch script key flagged; valid partial tables stay clean |
| modcheck | module/uses/provides semantics unchanged; new source construct accepted in modules | no implicit host ownership/pin bypass through a join |
| gramcheck | `join`, `completes`, `rejects`, `concurrency` scopes/keywords pinned | grammar does not accidentally bless rejected descriptor/order/task spellings |
| tsc | strict check passes regenerated goldens/adapters under the uniform target completion ABI | no `any`/suppression escape to force green; result unions and adapter signatures stay strict |

Planned full gates after each independently shippable implementation slice:

```text
go test -count=1 ./...
go run ./tools/modcheck
go run ./tools/gramcheck
strict TypeScript gate from tscheck/
LSP and lint regression suites
```

No gate is claimed green by B05.

**Checked falsification.** B00's acceptance plan requires independently
observable failures per affected phase rather than parser/happy-path proof.
A70 assigns target shape to the strict TS gate while keeping Can semantics in
canlc. The table above preserves those ownership lines.

### Q9b — Suggested implementation sequence

**Verdict: if this design is approved, plan the build in the following order;
do not combine them into one giant change.**

1. **J1 — syntax/AST only, fail closed downstream.** Parse/print/canonicalize
   the join and extern metadata, but static admission rejects execution until
   later slices are present.
2. **J2 — product typing + symbolic exhaustiveness.** Branch domains, labeled
   patterns, witnesses, no execution.
3. **J3 — authority + liveness admission.** Exact branch effects, pairwise
   conflict checks, reachable-extern completion/isolation checks, revision
   identity.
4. **J4 — hermetic evaluator + evidence.** Branch-local scripts, product
   execution, coverage, omission/lint behavior; still no target emit acceptance.
5. **J5 — uniform target completion ABI.** Migrate generated Can calls and host
   adapters as one explicit compatibility slice; strict TS green before join
   target lowering is enabled.
6. **J6 — join target lowering + goldens.** Concurrent launch/collect/product
   dispatch, state behavior, nested joins.
7. **J7 — linked-pure, LSP, lint, modcheck, gramcheck parity and final docs.**

J5 is intentionally before J6: otherwise a join emitter would depend on an
unstated per-function target distinction, which B05 rejects.

**Checked falsification.** The parser already fails closed on not-yet-enabled
`MatchInvoke` stages in b00-era implementation comments, establishing a local
precedent for syntax landing before execution only when every downstream phase
refuses the new kind. The sequence above uses the same discipline and makes the
large ABI migration explicit rather than coupling it to product semantics.

---

## Q10 — What remains deliberately unresolved?

### Q10a — No automatic failure-combination abstraction

**Verdict: closed for B05: none.**

The product matcher plus `<label>._` is the only compression. If a 3×3 join has
nine semantically different results, nine source behaviors are the cost. A
future `collect`, failure-set, priority, or applicative abstraction needs its own
payload and evidence design.

### Q10b — No dynamic shared-resource concurrency

**Verdict: deferred to Resources, with a hard B05 boundary.**

`concurrency isolated` covers only host operations that can honestly promise
invocation independence. Shared transactions, streams, handles, connection
pools with semantic ordering, lock ownership, and cancellation scopes wait for
Resources. B05 supplies no escape hatch.

### Q10c — No quantitative whole-program time theorem

**Verdict: deferred.**

Host externs name explicit bounds; Can code has termination proofs. B05 does not
sum those bounds through sequential paths, take maxima through joins, account
for scheduler overhead, or prove a function-level millisecond budget. That is a
resource-budget proof problem, not required to keep a join finite.

### Q10d — No race/select/streaming form

**Verdict: deferred with no reserved syntax.**

B05 waits for all branches. It does not expose first completion, partial
results, progress streams, or cancellation. Those would make completion order
observable and invalidate several simplifications above.

---

## Audit-obligation trace

| `ASTRA_AUDIT.md` §4 obligation | B05 ruling | Status |
| --- | --- | --- |
| every spawned task joined exactly once | Q3a: no task handle exists; the lexical `MatchJoin` node creates and consumes every branch | **addressed structurally** |
| exhaustive combined outcome product | Q2a/Q2b: finite labeled branch domains + symbolic product subtraction + concrete missing witnesses | **addressed statically** |
| branch authority explicitly partitioned | Q3b/Q3c: mandatory exact branch effect sets + pairwise write/read-write disjointness; Q3d makes host independence an explicit trusted contract | **addressed for Can state; host trust named explicitly** |
| explicit liveness basis / no smuggled host termination | Q4a: existing Can termination; Q4b: join-reachable host externs require `completes … else …` and rejection normalization; Q4c keeps scripts out of proof | **addressed** |

The audit's proposed protocol exception to state disjointness is **deferred**:
B05 has no synchronization/resource protocol to prove it. The safe first cut is
rejection, not an unproved exception.

---

## Final design record

1. One lexical `match join` surface; no task values.
2. Branch labels, explicit success record, mandatory exact `effects`, one root
   `call`, optional ordinary `given`.
3. No join descriptor and no separate order clause; declaration order fixes
   argument evaluation and launch, completion order is unobservable.
4. Zero branches reject; one branch is valid.
5. Arms use label-qualified outcome slots; `<label>._` is the only B05 product
   compression.
6. Exhaustiveness reuses a28-style symbolic product subtraction; no automatic
   failure policy.
7. All launched branches complete before dispatch; no early failure exit.
8. Exactly-once joining is structural, not a handle-analysis problem.
9. Same-cell write/read or write/write across siblings rejects; read/read is
   legal.
10. Multi-branch host participation requires `concurrency isolated`; dynamic
    shared resources wait for Resources.
11. Join-reachable host externs require `completes <int-param> else <error>` and
    `rejects <error>`; both mapped errors are declared zero-field outcomes.
12. Branch `given` rows are completion observations, never liveness proof;
    omission keeps a91 semantics; join-arm coverage remains observed execution.
13. The eventual TS backend uses one uniform completion-capable callable ABI;
    target promises never become Can values.
14. B02 declaring-stem ownership and exact revision pins remain; new host
    contract metadata is revision-significant.
15. Cancellation, resource handles, protocol exceptions, race/select,
    streaming, and quantitative time budgets are outside B05.

Approval of B05 authorizes a separate implementation-plan pass only. It does
not authorize compiler or library changes by itself.
