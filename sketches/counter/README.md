# counter — a bounded counter

- `counter.ail` — `mod counter`: one private cell
  (`Count__total`), bumped through declared `.read`/`.write`
  capabilities. Every test starts from init, so `three` and
  `three_again` agree — the tables prove sequential tests cannot
  interfere. `count__twice` threads state through a helper, so its
  capabilities arrive transitively.
- `counter.ts` + `errors.json` — committed golden TS prod emit
  (tests/given stripped; cells are module-scope `let`s).
  Regenerate: `go run ./compiler --out <dir>
  sketches/counter/counter.ail`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Proof: `docs/a09-effects.md`.
