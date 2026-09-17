# text — explicit text construction

- `text.ail` — `mod text`: `std__str__concat`, the basic explicit
  construction operation, with decision tables covering empties,
  multi-scalar text, and markup-passthrough. Concat preserves
  every scalar; safety lives in context encoders, never here.
- `text.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped). Regenerate: `go run ./compiler --out
  std/text std/text/text.ail`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/v16-text.md`.
Measurement and indexing wait on the scalar-access surface
decision recorded there.
