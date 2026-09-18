# reserve — S2 second consumer (reservation desk)

An independent motivating program for Schema S2 admission: a
dinner-reservation desk admitting `Reserve__Request` values
against a program-owned policy. Same frozen violation triple
and fail-fast order as the quota pilot, distinct domain and
distinct policy values (guest 1–24 scalars, party 1–8, slots
lunch/dinner). Built pre-checks as the S2 admission exhibit,
then migrated onto the S2a closed check ops per
`docs/a92-schema-s2b-call.md` with frozen expectations: the
validator binds one check per field and forwards check
outcomes unchanged.

- `reserve.can` — `mod reserve`: `reserve__request__validate`
  plus `reserve__request__admit`. Admit stays pure (validate
  then build a confirmation), so the whole admit path
  executes under linked-pure evaluation — unlike quota's
  effectful admit, no script fallback is needed for it.
- `reserve.ts` + `errors.json` — committed golden TS prod
  emit. Regenerate: `go run ./compiler --out sketches/reserve
  sketches/reserve/reserve.can std/quota/quota.can
  std/scalars/scalars.can`, then delete the co-emitted
  `quota.ts`/`scalars.ts` (same packaging gap as S1a:
  `reserve.ts` imports `./quota`, which resolves in
  whole-program compiles but dangles beside the committed
  goldens; the parity harness rewrites the specifier to
  `./quota.ts` in the temp copy only). Verify:
  `go test ./...`.
- No envelope: S1b proved nesting once; the S2 exhibit is
  the flat validator. Exhibits: `docs/a90-schema-s2-exhibits.md`;
  migration call: `docs/a92-schema-s2b-call.md`.
