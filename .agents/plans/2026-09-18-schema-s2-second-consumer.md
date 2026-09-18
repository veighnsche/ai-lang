## Goal

Pre-admission input for Schema S2 (closed descriptors): a second
independent concrete consumer that reuses the shared scalar
validators, so the S2 admission call can be made against concrete
duplication exhibits instead of speculation. Per `docs/a89-...` S2
gate: "a second caller of the same pilot is insufficient — the
missing source is a second independent motivating `.can` program
plus concrete duplication exhibits."

This slice builds the consumer only. The exhibits + admission
decision are a separate step after green.

## Second consumer: reservation desk

New sketch `sketches/reserve/reserve.can` (`mod reserve`), an
independent program under the same gates (modcheck walks
sketches/; emit goldens pin it):

- `Reserve__Request(guest: str, party: int, slot: str)`
- `Reserve__RequestSchema(guest_minimum, guest_maximum,
  party_lower, party_upper: int, allowed_slots: Seq<str>)`
- `reserve__request__validate` — same frozen triple
  `validation.schema_violation(path, rule, value)`, same
  fail-fast order (minimum-nonnegative → length → range →
  membership), payload reconstruction per S1a amendment 2
  (bind `_`, never read producer binders).
- `reserve__request__admit` — validate-then-confirm, kept
  pure so the whole admit path is linked-pure executable
  (stronger than quota's effectful admit, which needs
  scripts for the stateful tail).
- Baseline policy (program-owned, distinct from quota):
  guest 1–24 scalars, party 1–8, slots ["lunch", "dinner"].
- Rule tokens mirror the frozen vocabulary with own field
  paths: `str.length_scalars`, `int.closed_range`,
  `str.one_of`, `schema.guest_minimum`,
  `schema.guest_bounds`, `schema.party_bounds`.
- No envelope: S1b proved nesting once; the S2 exhibit is
  the flat validator's binding/interpretation layer.

## Sharing rules (decided from precedent)

- `uses` pins the four shared calls @1 (all rev 1 verified):
  `std__validate__str_length`, `std__validate__int_range`,
  `std__validate__str_one_of` (quota.can),
  `std__convert__int_to_str` (scalars.can). No vendoring —
  the S2 question is precisely what shared calls do NOT
  remove, so consumer 2 must be a genuine client of them.
- In-file declarations mirror the shared kinds:
  `validation.schema_violation` (raised) plus the four
  matched producer kinds (never raised, `_` binders only).
  Precedent: `math.zero_divisor` in division.can + ratio.can;
  `Str__Value` in quota.can + scalars.can.
- Cross-file calls carry per-site `given` scripts with
  reaching sets (own rows + same-file transitive callers,
  i.e. admit rows reach validate's call sites). a18-style
  verification + linked-pure vectors + node parity per §6.

## Work plan

1. `reserve.can` with 21 validate rows mirroring S1a's
   obligations (every constraint, every malformed-policy
   outcome, boundaries, priority rows, preservation, exact
   reporting incl. bigint, empty-policy, case, ws) + admit
   rows (admitted / rejected-before-confirm).
2. `errors.json` + `reserve.ts` via
   `go run ./compiler --out sketches/reserve
   sketches/reserve/reserve.can std/quota/quota.can
   std/scalars/scalars.can`, then delete co-emitted
   `quota.ts`/`scalars.ts` (same packaging gap as S1a).
3. `README.md` sketch note.
4. `TestGoldenReserve` in `emit_golden_test.go` (srcs
   reserve+quota+scalars; keep reserve.ts + errors.json).
5. Linked-pure vectors spanning reserve+quota+scalars
   (mirror `quota_request_linked_test.go`).
6. Node parity for success/kind/payloads/ordering (mirror
   `quota_request_parity_test.go`, incl. specifier
   rewrite if the same-dir emit gap bites).
7. Gates: `go test ./...`, `go run ./tools/modcheck`,
   `go run ./tools/gramcheck`.

## Non-goals

No descriptor machinery, no generic alias, no compiler
change, no second envelope, no claim the §1.5 row is
delivered. Exhibits + admission call follow in a
separate step on top of green.
