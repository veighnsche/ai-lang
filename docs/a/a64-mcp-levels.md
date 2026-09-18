# a64: MCP server + three-level file reading (proposal)

Status: proposal (pre-decision). No rule, no code, no golden.
Decision and sequencing are open. Amended once: landmarks,
declared agent surface, one ops table, attention memory, and
the rejected physical-sections alternative are folded in.

Goal: make "read the file" a structured, leveled operation
instead of a raw byte dump. An agent judging relevance
reads names; an agent checking behavior reads contracts;
an agent fixing a body reads source. Each level is a
computed projection of one file — never a duplicated
artifact — so views cannot drift from source. The server
itself is language-native: everything it serves traces to
a source declaration or a compiler derivation, never to
server-side heuristics.

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

## Landmarks: the 5-line start-here

Orientation dominates fresh sessions: every new agent
re-derives "what matters here" from scratch. Above L1
sits a computed landmark header per module: entry points
(provided but never called in-world), most-called fn,
primary types, most-handled errors. Derived from facts
the compiler already has — deterministic, zero
maintenance, cold-start-proof, never rots. Served first
by `can_map`, or as its own `can_landmarks` call.
Framing caution: landmarks guide but never gate. A
"start here" is a previous decision handed to a new
agent pre-loaded; the full L1 map stays one call away
so guidance cannot become a blinkered default.

## Declared agent surface (what "language-native" means)

The `mod` header already declares the code interface for
agents to read. Language-native MCP extends that: the
file declares its own agent surface in-band, and the
compiler checks it — the server serves declarations,
never guesses:

- Entry points: a checked claim (provided, tested,
  root-ish). Wrong claim = diagnostic, like any other
  declaration. This is the checkable half of "favorites";
  a bare `fav` marker no rule verifies would rot like a
  comment and is explicitly not proposed.
- Per-kind metadata: one doc line per error, severity,
  retryable or not (specified under a61 item 5; served
  here).

Objective salience lives in source if and only if the
compiler can verify it. Everything else lives in server
memory (below) or not at all.

## MCP server sketch

One operation table in the compiler (operation → handler
→ output schema), projected to both CLI flags and MCP
tools. Neither transport is first; the table is. Tool
rule: every tool maps 1:1 to a compiler phase or
artifact (parse → slices, check → diagnostics, eval →
run, catalog → trace). If canlc cannot compute it, the
server cannot serve it — no bespoke endpoints with
bespoke semantics. Conformance law: tool output
byte-equals CLI output for the same query.

Tools:

- `can_map(path)` → L1 (landmarks first).
  `can_contract(path)` → L2. `can_source(path [, fn])`
  → L3, whole file or one fn. `can_landmarks(path)` →
  the orientation header alone. The path of least
  resistance is landmarks-first, which makes token
  discipline structural instead of a system-prompt
  request agents defect from under pressure.
- `can_check(path)` → JSON diagnostics. Exists today as
  `--format json`; as a tool it becomes callable
  mid-reasoning with a schema instead of CLI arg
  guessing. Rides the diagnostic-quality work (a61 item
  1, a62, a63): every fix there improves this tool.
- `can_run(path)` → test results, structured per test.
- `can_trace(error_kind)` → the errors.json graph:
  raisers, handling arms, hitting tests. Gains mapping
  edges when a61 item 5 lands.
- `can_explain(code)` → the a61 per-code doc: rule,
  minimal violation, legal fix.

Resources for stable addressing across turns:
`can://auth.can?level=contract` — file slices citable
without re-reading. Views are read-only; edits happen on
full source, verified through `can_check`.

Stateful facilities live server-side, explicitly outside
language semantics, never presented as language truth:

- Attention memory: per-project record of what agents
  actually do (explicit stars plus implicit signal: L3
  reads, trace queries, fix locations), served as
  "previous agents found useful", with time decay so
  stale attention dies. Collaborative filtering for
  orientation; the only "favorites" that compound.
- Escalation metrics (below).

Why MCP rather than bash + `canlc`: tool schemas kill
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
contract / source reads as one decision each). Landmarks
earn their keep only if they cut L1 reads measurably.

## Non-goals

- Duplicated spec artifacts. No hand-written summaries,
  no Given-When-Then prose layer (cut from R9 in review
  as unverifiable duplication). Projections only.
- Editable views. The server never accepts an L1/L2 edit;
  there is exactly one source of truth, the `.can` file.
- New check/eval semantics. The server exposes what canlc
  computes; any behavior change is a compiler proposal
  first.
- Replacing the CLI or LSP. Transport addition, not
  migration.
- Black-box "smart" endpoints. Fixes come from compiler
  payloads (a61 item 1), never from server-side
  generation. Smart lives in the agent; the server
  serves truth.
- Physical section separators (`---` levels in source).
  Considered and rejected: the cut lands through R8's
  unit (`given` is site-attached inside bodies, so L2
  would hold expectations without scripts or reopen the
  settled call-site decision); it scatters each fn's
  signature + tests + body across three regions,
  weakening the edit loop to speed up skimming; and the
  marker is either new grammar (parser, diagnostics,
  migration) or decoration that drifts and then lies to
  head-truncation. CLI slicing gives the same token
  savings with none of the costs.

## Ordering and gates

1. Slicing + landmarks in the compiler, exposed through
   the ops table with CLI flags first — proves the
   projections without protocol work; goldens freeze all
   three levels plus landmarks for the gallery. Useful
   over bash alone, without the rest.
2. MCP transport over the same table. Conformance: tool
   outputs byte-equal the CLI flags.
3. Declared agent surface (entry checks, kind metadata
   with a61 item 5) and attention memory with decay.
4. Read-task benchmark (route / specify / fix over the
   gallery; level sufficiency + escalation rate).
5. Graduate or re-cut on the numbers from 4.

Items graduate separately; item 1 may proceed alone.
