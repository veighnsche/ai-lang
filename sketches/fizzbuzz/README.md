# fizzbuzz — FizzBuzz as one flat match

- `fizzbuzz.can` — `mod fizz`: a two-slot match (a28) on
  divisibility by 3 and by 5 selects the word while the number
  rides along in the same record, so no digit renderer is needed.
  Flat by construction: every arm is witnessed by its own table
  row under the test-per-arm law (a03). Rendering an int as text
  is left to `std__convert__int_to_str`, which is recursion
  behind a helper.
- `fizzbuzz.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped).
  Regenerate: `go run ./compiler --out <dir>
  sketches/fizzbuzz/fizzbuzz.can`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Proof: `docs/a28-multi-scrutinee-match.md`,
`docs/a03-branch-coverage.md`.
