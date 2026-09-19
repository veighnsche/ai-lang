# host-clock — shared host extern consumer demo

- `clock.can` — `mod clockdemo`: reads wall millis through pinned
  `std__clock__wall_now@1`. First cross-module host call (b02).
- `clock.ts` + `errors.json` — committed goldens. Regenerate:
  `go run ./compiler --out sketches/host-clock
  sketches/host-clock/clock.can std/host/host.can`;
  verify: `go test ./...` (`TestGoldenHostClock`).
