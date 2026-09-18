# notify — S2b third consumer (notification desk)

The first client of the S2a closed check layer: a
notification desk admitting `Notify__Request` values
against a program-owned policy (handle 3–16 scalars,
retries 0–5, channels email/sms). The validator binds
one check per field and forwards check outcomes
unchanged — no reconstruction arms, no converter
calls, no producer kinds (only
`validation.schema_violation` is declared or emitted).
Field order stays this program's chain.

- `notify.can` — `mod notify`: `notify__request__validate`
  plus the pure `notify__request__admit` (validate then
  build a dispatch ticket), so the whole admit path
  executes under linked-pure evaluation.
- `notify.ts` + `errors.json` — committed golden TS prod
  emit. Regenerate: `go run ./compiler --out sketches/notify
  sketches/notify/notify.can std/quota/quota.can
  std/scalars/scalars.can`, then delete the co-emitted
  `quota.ts`/`scalars.ts` (same packaging gap as S1a/S2:
  `notify.ts` imports `./quota` and `./scalars`, which
  resolve in whole-program compiles but dangle beside the
  committed goldens; the parity harness rewrites both
  specifiers in the temp copy only). Verify:
  `go test ./...`.
- Cost comparison against consumer 2 and the migration
  call for consumers 1–2: `docs/a92-schema-s2b-call.md`.
