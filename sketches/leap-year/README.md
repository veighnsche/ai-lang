# leap-year — Gregorian rule as one flat match

- `leap.can` — `mod leap`: a three-slot match (a28) on divisibility
  by 400, 100, and 4. The century exception reads top to bottom;
  wildcards keep four arms exhaustive, each witnessed by its own
  table row (a03). Naming follows the std `is_*` predicates
  (`std__int__is_zero`, `std__int__is_prime`).
- `leap.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped).
  Regenerate: `go run ./compiler --out <dir>
  sketches/leap-year/leap.can`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Proof: `docs/a28-multi-scrutinee-match.md`,
`docs/a03-branch-coverage.md`.
