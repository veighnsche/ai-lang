# std — the blessed standard library (not sketches)

`sketches/` holds temporary use-case examples and demos. `std/`
holds blessed standard-library programs: permanent, versioned,
golden-gated modules that other programs are meant to rely on.
Nothing here is provisional. A module lands here once, with its
motivating program and checks, and never as a sketch first.

Each directory is one self-contained program (cross-file calls
stay scripted until the linkage amendment lands, so blessed
helpers ship as same-file locals; see `docs/v13-stdlib.md`):

- `quota/` — row 1: scalar validators plus the quota counter
  that consumes them.
- `scalars/` — row 2: the monomorphic scalar catalog (boolean
  logic, comparisons, predicates, int/dec arithmetic,
  conversions, bounded backoff).
- `text/` — row 3, first cut: explicit text construction
  (`std__str__concat`). Measurement and indexing wait on the
  scalar-access surface decision in `docs/v16-text.md`.
- `division/` — row 4, first cut: exact Euclidean integer
  division (`std__int__divmod` and family). Decimal division
  stays refused; `gcd` and beyond wait on fuel patterns
  (`docs/v17-division.md`).

Admission rule (from the program doc): a module is not blessed
until it has a real motivating `.ail` program, complete
producer-owned outcomes, green compile-time decision tables,
negative examples for its invariants, a checked termination
argument where it iterates, and committed golden emit
(`*.ts` + `errors.json`, frozen by `go test ./...`).

Checks: `go test ./...`, `go run ./tools/modcheck` (covers
`std/` beside `sketches/`), `go run ./tools/gramcheck`.
Program: `docs/v13-stdlib.md`.
