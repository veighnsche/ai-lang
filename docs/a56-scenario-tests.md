# a56 — Scenario tests (Gherkin-shaped journeys, sketch)

Status: sketch (pre-decision). No rule, no code, no golden. This
records what first-class given/when/then could look like, why it
is NOT a replacement for `tests` + `given`, and the conditions
under which the real parts graduate to proposals.

## Observation

The language is already two-thirds Gherkin. `tests` rows are
`Then` (call specifications: `name(inputs) => expected`); `given`
tables are `Given` (scripted exchanges per test). What is missing
is a named `When` and multi-step sequencing: every test today is
one call, so sequencing hides inside wrapper fns (`count__twice`
exists only to bump twice in one test).

## Sketch (not specified)

A `scenario` groups one journey: labeled steps whose payloads are
the typed constructs the compiler already checks.

```
scenario happy
  given user u_01 exists
    exchange db__get_user args (id = "u_01")
      outcome Ok(...)
  when login with correct password
    call (id = "u_01", pw = seal Auth__Password("secret"))
  then session with 3 tries
    expect Ok(user_id = "u_01", remaining_tries = 3)
```

Step sentences are labels, not code; the indented block under
each is the checkable payload. Multi-step scenarios chain
when/then pairs over real cell state instead of fresh-init per
test. `Background` is a shared `given` prefix; `Scenario
Outline` + `Examples` desugars to N test rows.

## Verdict: not better as a replacement

Contract tables win for unit-level testing and stay the kernel:

- Density: 9 login specs are 9 lines plus factored tables; the
  scenario form is ~10x the lines for the same coverage.
- Exhaustiveness is checkable: "given keys equal test names" is
  a set comparison. A missing scenario is invisible; a missing
  table row is a diagnostic.
- No prose drift: a test line cannot lie about itself, but an
  unchecked label can ("locked out user" above `expect Ok`).
  Unchecked prose in a contract-first language is a liability.

What Gherkin adds decomposes into sequencing (real expressive
power) and prose labels (cosmetic; the reader of `.can` is an
engineer checking contracts, not a stakeholder reading
features). Take the power without the prose:

1. `seq` tests: multi-call tests (call, assert, call, assert)
   over real cell state. Closes the expressive gap with no
   natural language involved.
2. Optional test descriptions plus an emitted `.feature`
   projection: readability as compiler output, never as source
   of truth.

## Why journeys will matter anyway

Composition bugs live at the seams, and real projects are
mostly seams. Contract tests prove each fn in isolation; they
say nothing about whether `db → auth → session → HTTP
response` agree end to end. Every new row (HTTP, components)
multiplies the seams, so the demand for multistep acceptance
journeys grows with the language. But journeys rot fastest (any
hop change breaks them), so they stay a thin top layer: one per
critical path, combinatorial bulk underneath in contract
tables. Determinism is non-negotiable: every boundary scripted
per journey, no shared world, no unscripted IO.

## Non-goals

- Regex-to-step-definition glue, shared mutable `World`, or
  any natural-language semantics. Steps are typed calls, not
  sentences the runner interprets.
- Replacing `tests` + `given`. Scenarios layer above them;
  the kernel does not shrink.
- Cross-module journeys before the module story supports
  them. A journey calling several fns across modules is a new
  module-scope block with pins on every callee, not a fatter
  `tests` block — and it needs an effects/proof story first.

## Prerequisites, in order

1. `seq` tests (single-fn, multi-call) — the smallest real step;
   admissible on its own once designed.
2. HTTP row + component row — the seams journeys exist to cover.
3. Closed tagged unions (`variant`) — journey outcomes and
   mapped errors want them.
4. Cross-module `journey` block — only after 1–3 land.

## Feasibility trigger

Revisit when the HTTP and component rows have landed and a real
sketch has a journey that single-call tests cannot express
(DB → auth → endpoint with intermediate assertions). First
check whether `seq` tests plus existing tables already cover
it. Graduate to a proposal only with that consumer in hand; a
journey feature without a journey is a bug with a roadmap.
