# map — insertion-ordered maps with generic keys and values

- `map.can` — `mod map`: `std__map__get`, `std__map__insert`,
  `std__map__replace`, `std__map__remove`, and `std__map__entries`
  over `Map<K,V>` (one `Seq<Map__Pair<K,V>>` entries spine).
  A map is its ordered entries; replace keeps position, remove
  closes ranks, insert never overwrites. Errors are payloadless:
  error payloads cannot be generic.
- `map.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped). Regenerate: `go run ./compiler --out std/map
  std/map/map.can`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/ASTRA_STDLIB.md` §1.8.
