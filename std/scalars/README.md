# std-scalars — monomorphic scalar standard library

- `scalars.ail` — `mod scalars`: boolean logic, three-way
  comparisons (`-1`/`0`/`1`, false before true), selection of
  already-computed values, int/dec predicates, int/dec
  fundamentals (abs, negate, sign, min, max, clamp, distance,
  square, power, bounded add/subtract/multiply, factorial, sum),
  the pure NOW conversions between bool, int, dec, and str, and a
  bounded backoff delay with zero, negative, saturated, and
  large-attempt cases. Every function returns a named success
  record; fallible recursion splits into a total `decreases`
  helper plus a checked entry, so no error arm is unreachable.
  String order is byte order, documented in the file.
- `scalars.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped). Regenerate: `go run ./compiler --out
  std/scalars std/scalars/scalars.ail`;
  verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/v13-stdlib.md` (row 2).
Names follow the brief verbatim (issue 5).
