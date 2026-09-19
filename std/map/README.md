# map — insertion-ordered string-keyed maps

- `map.can` — `mod map`: `std__map__get`, `std__map__insert`,
  `std__map__replace`, `std__map__remove`, and `std__map__entries`
  over `Map__Entries<V>` (parallel keys/values, generic values).
  Fully generic `Map<K,V>` is unbuildable today: the Seq-elem rule
  rejects `Seq<Map__Entry<K,V>>`. A map is its ordered entries;
  replace keeps position, remove closes ranks, insert never
  overwrites.
- `map.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped). Regenerate: `go run ./compiler --out std/map
  std/map/map.can`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/ASTRA_STDLIB.md` §1.8.
Int-keyed maps wait on a second family; generic keys wait on
nested-argument Seq elements.
