# host — the host shelf: explicit foreign observations

- `host.can` — `mod host`: `std__clock__wall_now` and
  `std__clock__monotonic_now` over `host__wall_now` /
  `host__mono_now` externs. Time points are millis ints
  (JEV instant_repr millis_int 0.97): Unix-epoch millis for
  wall, unspecified-origin millis for monotonic.
- `host.externs.ts` — real host implementations (not throwing
  stubs). Same-slice maintenance with the extern decls.
- `host.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped). Regenerate: `go run ./compiler --out std/host
  std/host/host.can`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/ASTRA_STDLIB.md` §Host.
Consumers pin the wrappers (or externs) via `uses`; see
`sketches/host-clock`.
