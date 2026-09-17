# v72: closed tagged unions (borrow #3) — design

Status: proposal (pre-decision). No rule, no code.
Parent: `docs/ASTRA_FSHARP_BORROW.md` §3 (finite,
monomorphic, closed unions). This doc adapts that
sketch to AIL specifics; chatbot verdict requested
before any slice.

## Goal

Make invalid state combinations unrepresentable for
the stdlib rows waiting on this (HTTP methods, form
states, validation results), reusing the existing
outcome machinery (ctors, variant patterns,
`$ail_kind` discrimination) instead of inventing a
parallel case system.

## Syntax (proposed, not parsed)

Declaration mirrors `type`, with `case` rows:

```ail
variant Login__State rev 1 (
  case Anonymous()
  case Authenticated(session: Auth__Session)
  case Locked(user_id: str, remaining_seconds: int)
)
```

Case names are qualified (`Login__Anonymous`,
`Login__Authenticated`, `Login__Locked`): the
declaration binds the short row spelling to the
qualified name, and every other position spells the
qualified name. R3 naming already admits the shape.

Construction reuses kwargs ctors:

```ail
Login__Authenticated(session = s)
Login__Anonymous()
```

Patterns reuse variant arms on value matches:

```ail
match state
  on Login__Anonymous _ => Ok(message = "Sign in")
  on Login__Authenticated a => Ok(message = a.session.user_id)
  on Login__Locked l => Ok(message = l.user_id)
```

## Key decisions

1. **Closed, checked at match time.** A value match
   over a variant-typed scrutinee must name every
   case; `_` is rejected on such matches (the
   closure is the feature — a catch-all would
   silently swallow added cases). Matches over
   non-variant scrutinees keep today's rules
   untouched, including `_`.
2. **No new value plumbing.** Constructor checking
   mirrors error ctors (exact payload: unknown and
   missing fields rejected); runtime values carry
   kind + tag exactly like errors today
   (`matchSlot`'s `err`/`ErrKind` path generalizes,
   it does not fork). Field projection (`a.session`)
   is the existing record projection.
3. **TS mapping reuses `$ail_kind`.** Each case
   emits as a union member tagged with its qualified
   name; nullary cases emit as tag-only members.
   Value matches over variants emit `switch` on the
   tag, same as call matches today. No new runtime
   helpers.
4. **Revision discipline = interface change.**
   Adding a case, removing one, or changing payload
   fields requires a rev bump; same-revision drift
   is rejected under the §2.3 identity mechanism
   once it lands (variants join the fingerprint's
   referenced shapes). Adding a case fails every
   now-incomplete match at proof time — that
   invalidation is the acceptance test.
5. **Empty variants rejected** (a case-less
   declaration admits no constructor and proves
   nothing). Nullary cases are tag-only, not
   unit-typed.
6. **Payloads are named types only** (records,
   builtins, brands, other variants). No inline
   anonymous payloads, no generics, no recursion:
   a case payload naming its own variant (directly
   or through a cycle) is rejected. Errors stay
   separate: `error` kinds remain outcome-only with
   `emits`, never cases, and never matchable outside
   call matches.
7. **Single-scrutinee only for the first cut.**
   Variant patterns in multi-scrutinee (v28) slots
   are rejected until a later slice justifies the
   cross-product proof.

## Rejected alternatives

- **Open variants / catch-all `_`:** defeats the
  add-a-case invalidation, which is the feature.
- **Unqualified case names with import resolution:**
  a second name-resolution mechanism for zero
  gain; qualification is already the R3 shape.
- **Separate tag field (not `$ail_kind`):** forks
  emit, match lowering, and the tsc gate's
  discriminated-union reading for no semantic gain.
- **Merging errors into variants:** errors carry
  `emits` upper-bound contracts and scripted
  failure injection; collapsing them removes the
  injected/actual distinction the verdicts kept.

## Slices (after verdict, each probe-first + gates)

1. v73: declaration parse + case registry
   (duplicates, empty, bad payload types rejected).
2. v74: construction + payload exactness
   (checkTypes/eval, mirroring error ctors).
3. v75: patterns + value-match proof (case
   coverage, `_` rejection, cycle rejection).
4. v76: runtime match + TS emit + tsc gate still
   green; stdlib pilot consumer (one of HTTP
   methods / form states / validation results).

## Probes (shape, not exhaustive)

- Add `Login__Expired`: every incomplete consumer
  fails proof; complete ones pass unchanged.
- `Locked` arm projecting `session`: rejected.
- Same-revision case addition: rejected (once §2.3
  lands; before that, a version-bump test).
- Emit: tag-only nullary member present;
  `switch` covers all cases; tsc strict green.

## Open questions for the chatbot

None structural — positions above are all taken.
The verdict should CONFIRM/REVISE each numbered
decision, with counterexamples where REVISE.
