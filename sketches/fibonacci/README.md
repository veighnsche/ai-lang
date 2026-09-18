# fibonacci — Fibonacci without a loop

- `fib.ail` — `mod fib`: `fib__at` guards the domain (negatives
  are a declared error, not a hang) and `fib__from` walks down one
  proven self-call per step, carrying both accumulators so the
  decreases proof is a single unit step on `n` (a11). Same-file
  helper, so no call-site given (a07).
- `fib.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped).
  Regenerate: `go run ./compiler --out <dir>
  sketches/fibonacci/fib.ail`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Proof: `docs/a11-recursion.md`,
`docs/a12-contracts.md`.
