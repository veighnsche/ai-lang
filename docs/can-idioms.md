# can-idioms — how `.can` code is written here

Status: living. Three tiers, and the tier matters more than the
rule. **Law** fails the build with a named `CANnnnn` (look any
of them up with `canlc explain CANnnnn`). **Convention** is
unanimous blessed-`std/` practice held by review, with file
pointers. **Judgment** is a labeled decision made where the
language forced one — follow the documentation pattern, not the
decision. Anything here without a citation is suspect; challenge
it against the cited source, never against taste. An AI coder
that "cleans up" against general instincts instead of this file
goes reverse — it happened on 2026-09-18 (tier 2, rule C1),
which is why this file exists.

## Tier 1 — Law

### Calls dispatch, never sequence

Every call lives as a `match` scrutinee with dispatched arms;
bare calls in bodies or nested `f(call g(...))` are refused
(CAN3003, CAN4109; fix: wrap each call in its own match).
Callees resolve to a same-file helper, a `uses`-pinned
function, an own-module extern, or a state intrinsic — other
modules' externs are uncallable (CAN3001, CAN3002). Positional
or named args both bind through the signature (CAN3010).

### Match arms are exact

A call match covers exactly the callee's emits plus `Ok` — no
missing arms, no stale ones (CAN4101, CAN4102). A bool slot is
exactly `true` plus `false`, never `_` (CAN4103). A string slot
needs `_` for the open remainder (CAN4104). Arm kinds fit the
scrutinee (on-arms for calls, value patterns for values) with
uniform arity (CAN4105). First match wins; an arm earlier arms
fully cover is rejected before tests run, as is an
or-alternative that adds no space; wildcards stay their own arm,
never an alternative (CAN4111, CAN4112). Integer ranges are
closed `LOW..HIGH` with a strict gap — `5..5` must be the
singleton `5` — over literals or visible integer constants
(CAN4110); blessed code ranges over `std__ascii__` consts
(`std/html/html.can`, `html__url__scheme_token`) and
or-alternatives over consts (`html__attribute__id_ws`).

### Errors are declared, raised in-list

A body raises only kinds in its own `emits` (CAN4001); every
raised and listed kind is declared somewhere (CAN4002:
"Declare error m.stray(...) in its module"). Dotted namespaces
are shared practice, and the declaration may sit in the raising
module: `error net.down()` lives in `mod retry`
(`sketches/retry-loop/retry.can`), `error
encoding.invalid_base64url(value: str)` in `mod text`
(`std/text/text.can`). Module-header `emits` is parsed but
never checked against anything — `std/text/text.can` omits the
kernel `encoding.*` kinds its functions emit — so it is a
document, not a contract.

### Every arm is executed or certified

Each match arm must run at least once across the function's
decision tables (CAN4107), or carry the relay certificate: a
bound error arm of a *local* call whose body rebuilds the same
kind with every field unchanged (CAN4108; `compiler/lsp.go`,
`relayStatus`). The short spelling is `forward`:

```can
on math.zero_divisor e => forward e
```

`forward` is the whole RHS with exactly the arm's binder —
projections, outer names, and value matches are refused
(CAN3011; `compiler/forward.go` elaborates it into the complete
constructor, so admission and coverage see the handwritten
shape). Handwritten same-kind reconstruction
(`std/text/text.can`, `std__base64__decode`) elaborates to the
same thing; write `forward`, and write the reconstruction out
only when remapping kinds or payloads.

### Recursion has three blessed schemas

Helpers are acyclic; only direct self-recursion is admitted
(CAN3005). `decreases` names int param(s) and must be
load-bearing (CAN3006, CAN3007). The site passes the exact
step for its schema (CAN3008): unit loop `p - 1`
(`decreases fuel`); Euclid `(b, a % b)` under `decreases a, b
by euclid` (`std/scalars/scalars.can`, `std__int__gcd_euclid`);
narrowing `(lo, mid)`/`(mid, hi)` (`decreases lo, hi by
narrowing`). The self-call sits under the false arm of the
canonical guard — `p <= 0`, `b <= 0` (CAN3009) — and guard
tracking follows arity-1 matches only: wider tables never
carry a guard (`compiler/check.go`, recursion walk). So a
guard that protects recursion is always a single-scrutinee
`match`, even in a file that tables everything else.

### Multi scrutinees evaluate eagerly

Every scrutinee of a multi-match evaluates left to right,
exactly once, before dispatch — unlike nesting, a faulting
scrutinee faults first (CAN4109). Guard risky evaluation in an
outer match or nest it under its guard. For pure value
scrutinees there is no observable difference; reach for the
table (tier 2, C1).

### Tests are complete specifications

Every function ships a `tests` block (CAN3301). Rows are
uniquely named with complete args (CAN3201–CAN3203) and
complete error construction — `math.zero_divisor(divisor =
0)`, never a bare kind (CAN3204). Tables execute hermetically
during the build; any failure fails it (CAN4200). Foreign
calls (uses-pinned or extern) carry one `given` table per call
with an `exchange args (...) outcome ...` row per reaching
test, keys equal to test names in both directions, outcomes
restricted to the callee's emits (CAN3101–CAN3105, CAN3109,
CAN3110; unreachable tests script `-`, per `docs/a07-helpers.md`).
Local calls execute — no `given` on them, none on
multi-scrutinee matches (CAN3106). Hence pure `std/` code has
no `given` at all.

### Names, surface, revisions

Functions read `domain__verb`, types/variants/brands
`Domain__Name`, constants `domain__SCREAMING` (CAN2001–2003);
nothing shadows a kernel or primitive (CAN6012). `provides`
names exactly what the file defines, no more, no less
(CAN2301, CAN2302) — helpers included. Cross-file names pin
revs, functions and types alike (`uses [db__get_user@5,
Db__User@4, Db__Hash@1]`, `sketches/auth-login/`), one
provider per name (CAN2101–2103, CAN2201/2202). Every
declaration carries `rev`; a changed fingerprint at the same
rev is drift (CAN6013). State authority rides `effects
[Cell.read, Cell.write]` beside `emits`, transitive through
locals, never stale (CAN3107, CAN3108;
`sketches/counter/counter.can` for the shape). Constants
initialize from scalar literals only (CAN6016;
`std/ascii/ascii.can`). Bodies seal only their own module's
brands; test and script rows may name any brand as checked
data (CAN6004, `docs/a15-brands.md`).

### Values spell their types

No implicit conversions, ever — arithmetic takes same-type
operands and yields the operand type (CAN6003). Decimals
spell `d"12.50"`; bare `12.5` is ungrammatical (CAN6001).
Decimal `/` and `%` have no exact result: divide integers,
which divide Euclidean with `%` as the remainder (CAN6005,
`docs/a17-division.md`; `std__int__divmod` pins `dividend %
divisor`). `==` over `Bytes` is refused — compare lengths and
indexed members (CAN5005). Records are finite monomorphic
products, never cyclic (CAN6006). Sequences spell
`Seq<T>[...]` with explicit element type (CAN6007); `Bytes`
takes one `Bytes(Seq<int>[...])` literal with members `0..255`
(CAN6008, CAN6009). `+` concatenates strings, appends
`Seq<T> + T`, and joins same-brand str-backed values
(`docs/a16-text.md`, `docs/a39-seq-append.md`,
`docs/a42-fragment-join.md`); `-` and `*` stay numeric-only;
unary `-` takes int or dec; non-arithmetic operators yield
bool. `#` counts string scalars and sequence elements;
`[a:b]` slices, `[i]` reads. Pinned rows (`=> Ok(...) pinned`)
are trusted acceptance under `docs/a87-acceptance-authority.md`
(draft) and report weakening loudly (CAN6017).

## Tier 2 — Convention

### C1. Pure bool tables stay tables

A decision with no intermediate structure is one
multi-scrutinee `match` in counting order (`true,true /
true,false / false,true / false,false`), per the stated goal
of `docs/a28-multi-scrutinee-match.md`. Blessed tables:
`std__ratio__make` sign handling, `std__dec__sign`,
`std__int__lcm`. Nesting is for guards (tier 1: recursion
guards are arity-1; CAN4109 fault avoidance) and for arms that
run calls first (`std__str__starts_with`) — never for a pure
table. Partial wildcards follow `std__int__lcm` (`true, _`
then exhaustive rest); a full `_, _` catch-all has no
precedent — enumerate explicitly.

### C2. Bool fields match directly

`on Ok w => match w.value` with `true`/`false` arms
(`std__str__trim_left`, `html__attributes__make_from`). No `==
true`, no int flags.

### C3. Guards mirror blessed shapes

Slice under a length comparison (`match #pattern <= #value`
then `value[#value - #pattern:#value]`,
`std__str__ends_with`); upward seq walks under `match position
< #values` with `values[position]`, position flowing from zero
(`std__str__join_from`, `html__fragment__join_from`) — one
upper guard, no lower-bound ritual. Mirror, don't invent: a
guard shape the checker cannot see through fails the build,
and the fix is always a dumber guard.

### C4. Fuel workers pin exhaustion

`(work, acc, fuel)` with `decreases`, base arm returning the
accumulator, and a named row for the spent budget
(`zero_budget`, `short`, `over`, `empty`). Fuel is a
bound, not a cost model: generous fuel over fast-converging
work is normal; linear-or-worse fuel is honest only with the
cost said aloud (NOW† in `docs/ASTRA_STDLIB.md`; `std/ratio`
documents it per function).

### C5. Test rows are behavior-named

`snake_case` names after the edge: `zero_budget`,
`fidelity`, `overlap`, `trunc_a`, `padbits2`, `miss`.
Uniqueness is law (CAN3201); the naming is convention.

### C6. Not idioms (do not "fix")

Call-match arm order is mixed in blessed code (`Ok` first in
`std__base64__decode`, errors first in
`std__convert__str_to_dec`) — order carries no meaning.
`forward` versus handwritten same-kind reconstruction is one
spelling elaborated (`std/html` writes `forward`, `std/text`
writes it out). Positional versus named call args are both
bound through the signature; test rows always name theirs.

## Tier 3 — Judgment pattern

When the language cannot spell the specified shape, implement
the closest shape and document the gap in-file with its
trigger. The pattern, from `std/ratio/ratio.can`: brands wrap
str scalars only, so a numerator/denominator pair is a
transparent `Ratio__Fraction` record held by construction
through `make`; the gap (hand-forged zero denominators) is
closed by surfacing `math.zero_divisor` at every entry point,
and the header comment says exactly that, including where it
deviates from `docs/ASTRA_STDLIB.md`. A gap note names the
trigger values, the chosen behavior, and the wishlist delta —
never just the behavior.
