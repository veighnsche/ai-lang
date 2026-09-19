# seq — immutable ordered sequences

- `seq.can` — `mod seq`: `std__seq__empty`, `std__seq__singleton`,
  `std__seq__length`, `std__seq__get`, `std__seq__append`,
  `std__seq__concat`, and `std__seq__slice`, each generic over the
  element type with int + str decision tables. Bare-Seq returns are
  rejected by the language, so sequences travel in `Seq__Values<T>`
  and elements in `Seq__Item<T>`. Slice is half-open, same rule as
  text slicing; invalid bounds never clamp.
- `seq.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped). Regenerate: `go run ./compiler --out std/seq
  std/seq/seq.can`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/ASTRA_STDLIB.md` §1.8.
Callback transforms (map/filter/fold/find/all/any/sort/unique) wait
on function values; maps and sets wait on collection types.
