# quota-counter — validation plus a bounded counter

- `quota.ail` — `mod quota`: monomorphic scalar validators
  (`validate__require`, `validate__int_range`,
  `validate__int_nonnegative`, `validate__str_nonempty`,
  `validate__exclusive_pair`) plus a quota counter
  (`quota__consume`, `quota__usage`) that reuses them through
  same-file local calls, so every call executes its body.
  Validators return the accepted value or a producer-owned typed
  error; bounds are inclusive and reversed bounds fail instead of
  being silently swapped. Every test starts from init.
- `quota.ts` + `errors.json` — committed golden TS prod emit
  (tests/given stripped; the cell is a module-scope `let`).
  Regenerate: `go run ./compiler --out std/quota
  std/quota/quota.ail`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/v13-stdlib.md` (row 1).
