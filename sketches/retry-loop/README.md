# retry-loop — retry with fuel

- `retry.ail` — `mod retry`: fetch against a flaky server, recursing
  on `net.down` with `decreases fuel`. The loop is a proven
  self-call: base arm at `fuel <= 0`, step arm at `fuel - 1`. The
  server is a foreign `extern`, scripted per test; a site hit N
  times takes an N-element list.
- `retry.ts` + `errors.json` — committed golden TS prod emit
  (tests/given stripped). Regenerate: `go run ./compiler --out
  <dir> sketches/retry-loop/retry.ail`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Proof: `docs/a08-termination.md`.
