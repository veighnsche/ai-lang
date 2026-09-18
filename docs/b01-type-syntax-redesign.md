# B01 — type syntax redesign

Status: proposed. Full redesign of type declarations + annotations
(record/variant/brand/error decls, annotation grammar). Approved scope:
"full redesign" over the consistency packages.

Sequencing warning: the outsourced b00 function-values design grounds
in current syntax. Land this redesign first, then send the b00 prompt,
so the agent clones the new surface. Its `Fn<...>` annotation shape
must fit the grammar decided here.

## Non-negotiable semantics (surface may move, these may not)

N1. Revision identity per symbol (`docs/a67`, ASTRA_AUDIT §7):
    same-revision contract changes reject; baselines/certificates
    invalidate on body or dependency change.
N2. Flat global namespace with domain qualification: `qualifyCase`
    (`compiler/parse.go`) derives case domains from the name before
    `__`; modcheck, LSP, and lint all key on the convention.
N3. Nominal brands (proof, not runtime), closed unions with qualified
    cases, no implicit conversion (`compiler/types.go`).
N4. `provides`/`uses`/`emits` reference global names; the module
    header is out of scope.

## Current surface (grounded)

```
type Ratio__Fraction rev 1 (     type Ratio__Value<T> rev 1 (
  numerator: int                   value: T
  denominator: int               )
)
variant Login__Event rev 1 (     brand UserId is str rev 1
  case Ticked()                  error math.zero_divisor(divisor: int)
)
fn std__ratio__abs(value: int) -> Ratio__Value<int> rev 1
const std__ascii__NUL: int rev 1 = 0
```

Annotations: `str int bool dec Bytes`, `Domain__Name`, `Seq<T>`
(~200 uses, never nested), `Name<T>`. Migration scope: 442 decl
lines in std, 99 compiler test files, sketches/.

## Q1: revision spelling — `rev N` vs `@N`

Verdict: **attached `@N`** (`type Ratio__Fraction@1 (`).

`rev` stays a keyword for pins elsewhere, but decl headers move to
the pin spelling: TextMate already tokenizes `@[0-9]+` as
`constant.numeric.version`, and pins (`db__get_user@3`) prove the
shape. One revision spelling for declaring and pinning. Touches
parse, gramcheck ("rev must be a keyword" probe is re-pointed, not
deleted), TextMate, LSP, lint, and every decl line — all mechanical.

## Q2: declaration keywords — keep or rename

Verdict: **keep `type`, `variant`, `brand`, `error`**.

Renames (`record`, `union`, `newtype`) churn 442 lines for fashion.
The keywords are standard enough; the redesign budget goes to shape
and spelling, not vocabulary. Revisit only if user testing says
`variant` confuses (F# calls them unions — see ASTRA_FSHARP_BORROW).

## Q3: error shape — inline vs block

Verdict: **block form, no revision**:

```
error math.zero_divisor (
  divisor: int
)
```

Joins the decl family (block-newline like `type`); stays unrevised
because error identity is `domain.kind` and compat is tracked via
the emitter's `emits` clause (N1), not per-error revisions. Adding
error revs would be new revision-identity semantics, not syntax.

## Q4: brand shape — one-liner vs block

Verdict: **block form**:

```
brand UserId@1 (
  repr: str
  seals_from: [AdminId]
)
```

`is str` becomes `repr: str`; `seals_from` (currently a trailing
clause) becomes a row. Brands join the family; `repr` names the
backing for the day int-backed brands arrive (currently deferred).
Empty `seals_from` row omits (str-only minting, as today).

## Q5: sequence spelling — `Seq<T>` vs `[T]`

Verdict: **`[T]`**.

The one builtin collection gets builtin spelling. `Seq` remains the
checker/emitter internal name (like brands: proof vs surface); only
the annotation and literal spelling changes. `<>` stays reserved
for user generic application (`Ratio__Value<int>`), so the two
bracket kinds never collide: `[]` builtin, `<>` user. Nested `[[T]]`
stays rejected (v1 shape rule, `seqElemName`).

## Q6: `Domain__Name` — keep or module-relative short names

Verdict: **keep globals; module-relative naming is a follow-up**.

Short names in the declaring module (`type Fraction` elaborated via
the enclosing `mod`) need module-aware resolution in the checker —
a module-system redesign, not a syntax upgrade. N2 stays. The
doc proposes no new separator (`::` considered; `.` is taken by
projection and error kinds).

## Q7: field rows and generic params

Verdict: **unchanged** (`name: Type` rows, `Name<T>` params).

Standard, compact, and shared with fn params. No better shape on
offer; churning them buys nothing.

## Proposed surface (after)

```
type Ratio__Fraction@1 (         type Ratio__Value<T>@1 (
  numerator: int                   value: T
  denominator: int               )
)
variant Login__Event@1 (         brand UserId@1 (
  case Ticked()                    repr: str
)                                )
error math.zero_divisor (        fn std__ratio__abs(value: int) -> Ratio__Value<int>@1
  divisor: int
)
const std__ascii__NUL: int@1 = 0
values: [int]
```

Note: `fn`/`const` take `@N` too — one revision spelling (Q1), and
`fn` is already in the migration blast radius via annotations.

## Considered and rejected

- `type` → `record`, `variant` → `union`: fashion churn (Q2).
- Per-error revisions: new N1 semantics disguised as syntax (Q3).
- `::` or `.` domain separators: `.` taken; `::` churns N2
  machinery for no semantic gain (Q6).
- Structural/record-shorthand literals: new type system, not new
  syntax. Out of scope.
- Keeping `Seq<T>`: the single most-dated spelling in the
  language; `[T]` is unambiguous once `<>` is user-only (Q5).

## Implementation phases (on approval)

1. Grammar + scanner: TextMate `@N`-on-decl, keyword rows for
   `repr`/`seals_from`; gramcheck probes re-pointed.
2. Parse: `@N` headers, block errors/brands, `[T]` annotations;
   `parseTypeParams`/`parseTypeBinds` untouched.
3. Check/emit: internal names unchanged (`Seq`, `rev` ints) —
   elaboration absorbs the surface; `seqElemName` keeps its rule.
4. Gates: modcheck, LSP, lint probes for the new shapes.
5. Migration: 442 std decl lines, 99 test files, sketches/,
   golden files; full suite + gates green.
6. b00 handoff: send the function-values prompt against the new
   surface; its `Fn` annotation must fit Q5/Q7.

## Open questions + risks

- Q2-revisit: does `variant` confuse against F# `union`? Needs a
  reader test, not more deliberation.
- b00 sequencing: if the outsourced agent starts before this
  lands, its syntax grounding rots. Hold the prompt until green.
- `@N` on `fn` changes every fn header in std and tests; the
  mechanical migration must not smuggle semantic edits — review
  the migration diff separately from the compiler diff.
- Unresolved: exact `seals_from` row grammar when combined with
  future int-backed `repr` (deferred with int brands).
- Omitted scope: `mod`/`provides`/`uses`/`emits` headers, match/
  test/contract syntax, module-relative naming (follow-up).
