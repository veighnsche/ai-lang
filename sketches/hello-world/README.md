# hello-world — the smallest complete program

- `hello.ail` — `mod hello`: one type (`Hello__Message`), one pure
  nullary function, one decision-table row. No state, no effects,
  no errors, no match. The body is a direct construction.
- `hello.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped).
  Regenerate: `go run ./compiler --out <dir>
  sketches/hello-world/hello.ail`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md` Goal + R7–R8 (logic and proof are one act).
