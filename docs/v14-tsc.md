# v14 — tsc verification story (decision)

Status: landed (decision). No code change; no golden change.

## Finding

R11 carries one claim with zero mechanism: "`tsc` re-checks types and
contracts, not coverage." Nothing runs `tsc` anywhere. The repo has
no CI directory, no TypeScript toolchain, and no gate that executes
it. The claim has been aspirational since it was written.

## Decision

The clause is suspended until a gate ships. `ailc` is the sole
verifier: its check phases plus the hermetic decision-table run
are the complete proof, and emitted TypeScript is unverified
output until the reinstatement criteria below are met.

## Reinstatement criteria

A future PR reinstates the clause by shipping all three:

1. A CI job running `tsc --strict` over every committed golden
   `.ts` file plus extern-stub fixtures for each `.externs`
   import the goldens reference.
2. A green run on the current goldens, including the multi-shape
   ok unions from the v13 emitter change.
3. A documented policy for what `tsc` owns (type and contract
   shape) versus what stays exclusively in `ailc` (exhaustiveness,
   termination, decision tables).

## Priced consequence, recorded now

Current emit may not pass `--strict`: field access lands on
union-typed bindings, and ok members share one `kind`
discriminant across shapes. If the gate arrives and fails on
real output, the fix is emit narrowing or per-function result
types — not suppressions. That cost is accepted in advance so
the future gate cannot be negotiated down to a green-looking
check that proves nothing.

## Consequences (accepted by writing this down)

- REQUIREMENTS.md is untouched: the freeze rule means amendments
  land tagged, never by silent rewrite. This doc is the tag.
- No test or tooling change ships here; inventing an unrunnable
  CI job to look resolved would be exactly the aspirational
  pattern this decision retires.

## Reinstated (v70)

The clause is live again: R11's "`tsc` re-checks types and
contracts, not coverage" is now a running gate, not an
aspiration.

- Gate: `.github/workflows/tsc.yml` runs pinned
  TypeScript 5.9.2 (`tscheck/`, lockfile committed)
  `tsc --strict --noEmit` over every committed golden
  `.ts` (9 files: 5 std + 4 sketches) plus 2
  hand-written extern-stub fixtures
  (`sketches/auth-login/auth.externs.ts`,
  `sketches/retry-loop/retry.externs.ts`). Green on
  arrival: the priced emit-narrowing risk did not
  materialize, so no emit change shipped.
- Ownership: `tsc` owns type + contract shape of
  emitted output. Everything else — exhaustiveness,
  termination, decision tables, coverage — stays
  exclusively in `ailc` (`go test ./...`). A green
  `tsc` run never substitutes for the ailc proof, and
  an ailc-green program never skips the `tsc` shape
  check on changed emit.
- Follow-up (not this slice): drift-checking stubs
  against their `.ail` extern decls (a signature
  change today rots its stub silently until a human
  notices the mismatch).
