# host — the host shelf: explicit foreign observations

- `host.can` — `mod host`: `std__clock__wall_now` and
  `std__clock__monotonic_now` over `host__wall_now` /
  `host__mono_now` externs. Time points are millis ints
  (JEV instant_repr millis_int 0.97): Unix-epoch millis for
  wall, unspecified-origin millis for monotonic.
  `std__random__bytes` (CSPRNG, 1MB cap) and `std__hash__digest`
  (sealed `Hash__Profile` brand: sha256, sha512; JEV
  hash_profile brand 0.98) over `host__rand_bytes` /
  `host__hash_digest`. `std__secret__equal` (sealed
  `Secret__Value` brand behind timingSafeEqual; JEV
  secret_repr brand 1.0) and `std__env__read` (sealed
  `Env__Name`; denied is a declared upper bound, v1 hosts no
  policy) over `host__secret_equal` / `host__env_read`.
- `host.externs.ts` — real host implementations (not throwing
  stubs). Same-slice maintenance with the extern decls.
- `host.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped). Regenerate: `go run ./compiler --out std/host
  std/host/host.can`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/ASTRA_STDLIB.md` §Host.
Consumers pin the wrappers (or externs) via `uses`; see
`sketches/host-clock`.
