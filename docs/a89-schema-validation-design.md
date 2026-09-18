# a89: std__validate__schema — design (proposal)

Status: proposal (pre-decision). No rule, no code.
Parent: `docs/ASTRA_STDLIB.md` §1.5 (`validate__schema`),
§6 order 9 (schemas and canonical codecs).

## Goal

Define the reusable schema-validation surface so §1.5 can
close: one surface with structured paths and complete
payloads, per the catalogue. This doc settles the design;
the S1 pilot below is the first and only authorized code
slice, and it starts after this design is reviewed.

## Why the catalogue signature cannot land as written

`(value: T, schema: Schema<T>) → T !
[validation.schema_violation]` needs three things that do
not exist:

1. **Generics (order 8).** `Schema<T>` cannot be named and
   no `fn` abstracts over `T`. `knownType`
   (`compiler/types.go`) admits base types, records,
   brands, `Bytes`, and `Seq` over plain elements only.
   Unlike `one_of`/`all`, monomorphizing does not rescue
   the generic surface: a hand validator for one record is
   application code, not a reusable surface. The
   shape-polymorphism is the feature.
2. **Descriptors and traversal (order 9).** There are no
   schema descriptor values (field specs, nested schemas,
   constraints, revisions) and no way to walk a record
   generically: field access is static (`r.field`), with
   no field-by-name projection or reflection. The
   existing `std/schema` is asset approval only, and its
   charter calls that "the asset minimum of Schema, not
   the general layer" (`docs/schema-design-prompt.md`).
3. **A trust boundary.** Inside `.can`, the checker
   already proves record shapes at compile time, so
   runtime structural validation of an already-typed
   value is vacuous. Validation earns its keep at a
   decode boundary (wire bytes → `T`), and no codec
   (`json__decode`, form decode, `row__decode`) exists
   yet.

## Key clarification: what schema validation checks

Shape is static; a schema is a reusable **constraint
bundle plus structured reporting**. Two framings:

- (a) Value is already `T`; the schema holds per-field
  constraints; the validator returns `T` or a violation
  with a path. Matches the catalogue signature.
- (b) Untyped wire data → `T` construction. That is
  decode (`json__decode`, `row__decode`), a separate
  operation.

Recommend (a): `validate__schema` checks constraints and
never constructs `T` from untyped data. Decode operations
own construction; descriptors are shared if and when both
exist. Nothing in S1 pre-builds for (b).

## Proposed pilot (S1): monomorphic, one record, no language changes

One concrete record with str, int, and enum-like fields
(candidate: an account-profile or config record, fixed in
S0 review), one plain-record "schema" value holding its
constraints (bounds, length limits, `allowed: Seq<str>`),
one validator that walks the known fields with static
access and ordinary calls to the existing validators
(`str_length`, `int_range`, plus `str_one_of` once that
slice lands — note the ordering dependency in the open
questions).

Why this shape: no new surface; schema-as-plain-record
proves descriptor-as-value; static traversal; per-field
calls reuse blessed validators instead of re-implementing
bounds. Candidate-A spirit throughout.

## Design decisions

1. **Schema-as-value: plain records (CONFIRMED for S1).**
   The pilot schema is an ordinary rev'd record. Closed
   descriptor records (`FieldSpec` rows in a `Seq`) wait
   for a second consumer; generic `Schema<T>` waits for
   order 8. Rejected: a new `schema` declaration form now
   (premature — one consumer); an opaque brand
   (approval-style opacity fits authority, not
   constraints: tests and decision tables must see the
   bounds).
2. **Traversal without reflection: static per-field calls
   (CONFIRMED for S1).** Generic traversal is explicitly
   unsolved here and waits on order 8+. Rejected: runtime
   field projection (new surface, unmotivated by one
   record); codegen (none exists).
3. **Violation payload: single kind, rendered value
   (RECOMMENDED).** One kind has one field list, but the
   offending value's type varies per field. Minting
   per-type kinds (`schema_violation_str/int/...`)
   explodes; carrying the raw value cannot fit one kind.
   Recommend `validation.schema_violation(path: str,
   rule: str, value: str)` with the offender rendered
   through the canonical converters (`int_to_str`,
   `dec_to_str`, `bool_to_str` all exist and are exact or
   canonical). Complete enough for every scalar, one
   kind, printable.
4. **Path encoding: dot-joined str (RECOMMENDED for S1).**
   `"address.zip"`, empty for the root. A `Seq<str>` path
   is the nicer value, but no existing error carries a
   `Seq` field and that path is untested — do not make
   the pilot the probe. Revisit if `Seq` error fields get
   proven elsewhere.
5. **Nesting: composition by ordinary calls
   (CONFIRMED).** A nested record's validator runs and
   the caller prefixes its path on relay. No descriptor
   nesting in S1; the pilot record stays flat, with at
   most one nested call proving the prefix rule.
6. **Revisioning: type rev only (CONFIRMED for S1).**
   Descriptor revision is the constraint record's own
   `rev`. No separate `SchemaRevision` until order 9.
7. **Relation to codecs: strict split (CONFIRMED).**
   `validate__schema` never touches untyped data; decode
   ops never share code with it until S3. Shared
   descriptors are deferred, not implied.

## Non-goals (all slices here)

Business policies (password strength, usernames, email
domains — the catalogue forbids them in generic
validation); generic `Schema<T>`; migration; JSON/form/row
codecs; `Seq<Variant>` outcomes; any compiler change in
S1.

## Slices

- **S0 (this doc): design review.** Acceptance: decisions
  1–7 ruled, pilot record chosen, payload kind fixed,
  open questions closed. No code.
- **S1: pilot validator.** Monomorphic schema record +
  validator + motivating program + decision tables with
  negative rows for every constraint and the path-prefix
  rule + committed goldens. No compiler changes, no new
  diagnostic codes. Acceptance: the §6 admission items
  for the pilot; `go test ./...`, `modcheck`,
  `gramcheck` green.
- **S2: descriptor records.** Only if a second consumer
  appears. Closed `FieldSpec` rows, still monomorphic.
  No speculative blessing.
- **S3: generic surface.** Once order 8 lands, migrate
  the pilot to `Schema<T>`; keep per-type payload kinds
  only if the one-kind rule forces it.
- **S4: codec integration.** Decode operations share
  descriptors with validation. Owned by order 9, not by
  this doc.

## Open questions (must close before S1)

1. Pilot record choice: account-profile vs config. Needs
   str + int + enum-like fields and a real motivating
   program. Propose both candidates with their programs
   at S0 review.
2. Payload: confirm `(path, rule, value: str)` single
   kind, including the empty-path-for-root convention.
3. Ordering vs the `one_of`/`all` slices: the pilot
   reuses `str_one_of` for membership. Either land
   `one_of` first or inline membership in S1 with a
   named follow-up — decide at S0, do not silently
   duplicate the membership contract.

## Rollback

This doc authorizes no code, so there is nothing to roll
back. If S0 review rejects the pilot-first shape, the
fallback is the catalogue order: `validate__schema`
waits for orders 8+9 and §1.5 stays open except it.
