# docs — ai-lang design records (reviewer index)

Start here. This folder is the language's memory: plans, per-feature
specs, and the rules each feature had to satisfy before it landed.

## Status map

| Doc | Status | One line |
|---|---|---|
| `REQUIREMENTS.md` (repo root) | Living: v0.1 freeze + ratified amendments tagged v04–v12 | The rules; see freeze note below |
| `v02-machine-artifacts.md` | Shipped | Codes, `--format=json`, `normalize`, catalog |
| `v03-branch-coverage.md` | Shipped | Test-per-arm law over green tables |
| `v04-type-discipline.md` | Shipped | Brands, `seal`, exact `dec`, no floats |
| `v05-expressiveness.md` | Landed (plan) | The gap + build order; items 1–6, 8 shipped, 7 declined |
| `v06-arithmetic.md` | Shipped | `+`, `-`, `*` exact-or-loud; division deferred with reason |
| `v07-helpers.md` | Shipped | Same-file calls, no pin, no call-site `given` |
| `v08-termination.md` | Shipped | Proven self-recursion via `decreases`; cycles refused |
| `v09-effects.md` | Shipped | Private cells, declared capabilities, per-test stores |
| `v10-numerics.md` | Shipped | One numeric semantics: unbounded ints, exact decs, exact TS emit |
| `v11-recursion.md` | Shipped | Program-wide recursion ban, guarded unit steps, returned-outcome theorem |
| `v12-contracts.md` | Shipped | Producer-owned emits, complete error expectations, exchange script rows |
| `v13-stdlib.md` | Landed (part) | Stdlib rows 0–2: linkage decision, quota-counter validation, scalar catalog |
| `v14-tsc.md` | Landed (decision) | R11 tsc clause suspended until a real gate ships; ailc is the sole verifier |
| `v15-brands.md` | Shipped | Bodies seal only their own module's brands (AIL6004); tests/scripts name any brand |
| `v16-text.md` | Landed (part) | `+` concatenates strings; Unicode-scalar indexing decided; scalar-access surface open |
| `v17-division.md` | Landed (part) | Exact Euclidean `/` and `%` on ints; dec refused (AIL6005); gcd family waits on fuel |
| `v35-typed-fragments.md` | Sketch (pre-decision) | HTMX-shaped fragment endpoints after unions, HTTP, components; revisit trigger, no rule yet |
| `v36-seq-typed-construction.md` | Shipped (S1) | Typed `Seq<T>` literals, element checking, AIL6007; length/access/append are later slices |
| `v37-seq-length.md` | Shipped (S2) | `#` counts sequence elements; `std__seq__length` realized by `#`, no second spelling |
| `v38-seq-access.md` | Shipped (S3) | `xs[i]` over `Seq<T>` yields `T`; bounded traversal under existing `decreases`; `std/seq` deferred |
| `v39-seq-append.md` | Shipped (S4) | `Seq<T> + T` copy-on-append; concat refused; compiler complete, customers next |
| `v40-str-join.md` | Shipped (customer 1) | `std__str__join` preserves order; positional first-element test; split next |
| `v41-str-split.md` | Shipped (customer 2) | `std__str__split` retains empties, leftmost policy; wrapper return; mints `text.empty_separator` |
| `v42-fragment-join.md` | Shipped (customer 3) | `html__fragment__join` over explicit children; same-brand `+`; attributes fork posed |
| `v43-named-attributes.md` | Shipped (fork verdict B) | Name-carrying wrappers; pairs minted, never parsed; `make` next |
| `v45-bytes-values.md` | Shipped (Bytes B1) | `Bytes` value admission; `Bytes(Seq<int>[...])` literals `0..255`; `Uint8Array` emit; export/codecs next |

## Reading order for a reviewer

1. `REQUIREMENTS.md` Goal + R1–R9 (the thesis and the shape).
2. `v05-expressiveness.md` (what was missing and in what order).
3. `v06` → `v09` in order (each spec pairs a power with its proof).
4. `sketches/` live shape: `auth-login/`, `retry-loop/`, `counter/`;
   `std/` blessed library: `quota/`, `scalars/`.
5. `CLEAN_ROOM_REVIEW.md` (design input; historical record, see note).

## Rules for reading (and editing)

- The v0.1 freeze means: no silent drift. Amendments land tagged with
  their version (`(v07)`), never by rewriting a ratified rule.
- Every expressive power names its proof cost; a feature whose proof
  is "future work" is a bug with a roadmap.
- One rule, one `AILnnnn` code (`compiler/code.go` is the registry).
- Verify claims mechanically: `go test ./...`,
  `go run ./tools/modcheck`, `go run ./tools/gramcheck`.
  Goldens live beside their sketches; `broken-login/` titles are
  enforced by the suite, not by inspection.
