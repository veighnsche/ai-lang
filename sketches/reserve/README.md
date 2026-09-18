# reserve — S2 second consumer (reservation desk)

An independent motivating program for Schema S2 admission: a
dinner-reservation desk admitting `Reserve__Request` values
against a program-owned policy. Same frozen violation triple
and fail-fast order as the quota pilot, distinct domain and
distinct policy values (guest 1–24 scalars, party 1–8, slots
lunch/dinner) — the S2 question is what binding and
interpretation work repeats per record that the shared scalar
calls do not remove, so this sketch imports all four shared
calls (`uses` pins into `std/quota` and `std/scalars`) and
vendors nothing.

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
  `reserve.ts` imports `./quota` and `./scalars`, which
  resolve in whole-program compiles but dangle beside the
  committed goldens; the parity harness rewrites both
  specifiers to `./quota.ts`/`./scalars.ts` in the temp copy
  only). Verify: `go test ./...`.
- No envelope: S1b proved nesting once; the S2 exhibit is
  the flat validator. No descriptor machinery here — the
  exhibits and the admission call follow on top of green
  (see `.agents/plans/2026-09-18-schema-s2-second-consumer.md`).
