# v64: MCP server + three-level file reading (proposal)

Status: proposal (pre-decision). No rule, no code, no golden.
Decision and sequencing are open. This proposes two ideas
that multiply: progressive-disclosure file views (L1/L2/L3)
served through an MCP server, so agents pay context only
for the depth the task needs.

Goal: make "read the file" a structured, leveled operation
instead of a raw byte dump. An agent judging relevance
reads names; an agent checking behavior reads contracts;
an agent fixing a body reads source. Each level is a
computed projection of one file — never a duplicated
artifact — so views cannot drift from source.

## The levels, against current syntax

The grammar is already nearly level-sliced: header, then
declarations (signature + `tests`), then `=` bodies. The
levels fall out of it:

- L1 — map. `mod` header, type names + field types, fn
  names + param/return types + `emits` + `effects`,
  error names, test names only. Enough to route ("where
  is `auth.mismatch` handled?"), plan, and judge
  relevance. Estimated 10–20% of file tokens.
- L2 — contract. L1 plus full `tests` rows and `given`
  tables: the complete behavioral spec without
  implementation. Enough to answer "what happens on
  `db.down` twice?" without reading one match arm. R7+R8
  make this level possible; no mainstream language can
  offer it, because nowhere else are tests and stubs
  part of the signature block.
- L3 — source. Everything, bodies included. Needed only
  when the contract is wrong or the implementation is
  suspect.

R9 already demands L1 ("judge a file from formal lines
only") plus a layout contract to enforce. This proposal
supersedes the enforcement half: with AST-aware slicing,
the server serves any level of any file regardless of
ordering, plus finer cuts for free (one fn body, one fn's
tests, all `emits` in a module). Layout stays canonical
by formatter; levels never depend on it.

## MCP server sketch

Thin transport over existing `ailc` internals: same Go
package, new wire protocol. The day it duplicates
check/eval logic instead of calling it, it is a bug.

Tools:

- `ail_map(path)` → L1. `ail_contract(path)` → L2.
  `ail_source(path [, fn])` → L3, whole file or one fn.
  The path of least resistance is L1-first, which makes
  token discipline structural instead of a system-prompt
  request agents defect from under pressure.
- `ail_check(path)` → JSON diagnostics. Exists today as
  `--format json`; as a tool it becomes callable
  mid-reasoning with a schema instead of CLI arg
  guessing. Rides the diagnostic-quality work (v61 item
  1, v62, v63): every fix there improves this tool.
- `ail_run(path)` → test results, structured per test.
- `ail_trace(error_kind)` → the errors.json graph:
  raisers, handling arms, hitting tests. Gains mapping
  edges when v61 item 5 lands.
- `ail_explain(code)` → the v61 per-code doc: rule,
  minimal violation, legal fix.

Resources for stable addressing across turns:
`ail://auth.ail?level=contract` — file slices citable
without re-reading. Views are read-only; edits happen on
full source, verified through `ail_check`.

Why MCP rather than bash + `ailc`: tool schemas kill
CLI-syntax guessing and grep-parse brittleness; L1-first
becomes the default path rather than a convention; one
integration serves every harness (Claude Code, Cursor,
own agent) instead of per-editor LSP work.

## The bet that must be measured

Levels save tokens only if L1/L2 suffice for most tasks.
If agents constantly escalate to L3, levels add
roundtrips instead of removing context. Instrument from
day one: escalation rate per level per task class. If L1
answers fewer than ~70% of reads, the level boundaries
are wrong, not the idea — re-cut before adding levels.
Never add L1.5; three levels is the whole point (map /
contract / source reads as one decision each).

## Non-goals

- Duplicated spec artifacts. No hand-written summaries,
  no Given-When-Then prose layer (cut from R9 in review
  as unverifiable duplication). Projections only.
- Editable views. The server never accepts an L1/L2 edit;
  there is exactly one source of truth, the `.ail` file.
- New check/eval semantics. The server exposes what ailc
  computes; any behavior change is a compiler proposal
  first.
- Replacing the CLI or LSP. Transport addition, not
  migration.

## Ordering and gates

1. AST slicing in the compiler (`ailc map/contract`
   as CLI flags first — proves the projections without
   any protocol work; goldens freeze all three levels
   for the gallery sketches).
2. MCP transport over the same entry points. Conformance:
   tool outputs byte-equal the CLI flags.
3. Escalation instrumentation + a read-task benchmark
   (route / specify / fix tasks over the gallery;
   record level sufficiency per task).
4. Graduate or re-cut on the numbers from 3.

Item 1 may proceed alone and is useful without the rest:
leveled CLI output already serves agents over bash.
