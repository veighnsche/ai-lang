# set — sets with explicit ordering

- `set.can` — `mod set`: `std__set__contains`,
  `std__set__union`, `std__set__intersection`, and
  `std__set__difference` over `Set__Members<T>`. Members keep
  first-occurrence order; union appends only absent members;
  intersection and difference preserve the left order. Equality
  is per-instance `==`, shared with maps.
- `set.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped). Regenerate: `go run ./compiler --out std/set
  std/set/set.can`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/ASTRA_STDLIB.md` §1.8.
