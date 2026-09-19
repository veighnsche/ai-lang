# b02 — shared externs design

Follow-up: [B10](b10-bare-and-extern-returns.md) broadens extern successes
from records to all supported data-only return types, retaining the sharing,
pins, declaring-stem host resolution, and trust model below.

## Problem

Externs are module-local by enforcement, not by necessity: a
`std/host` module compiles but no other module can call it
(`docs/stdlib-remaining.md:34`), blocking the host shelf
(clock/random/hash/secret/log/env). The rejection is one branch in
`checkCalls` (`compiler/check.go:747-751`); everything downstream
already treats a declared foreign extern like a foreign Can
function.

## Surveyed evidence

| Mechanism | Status | Location |
|-----------|--------|----------|
| uses-pin resolution (`provides` incl. externs) | ready | `compiler/check.go:236-239, 346-366` |
| exact-rev check (`declaredRev` incl. externs) | ready | `compiler/check.go:67-94` |
| callee classification (any declared extern → `CalleeForeign`) | ready | `compiler/export.go:264-266` |
| given-table required + stub-checked via `EmitsOf` | ready | `compiler/check.go:937-1000` |
| script-consistency proof (extern params resolve) | ready | `compiler/eval.go:1806` |
| eval scripted-outcome dispatch (`prog.Uses` passes) | ready | `compiler/eval.go:1438` |
| unused-uses accounting (`calledFns` incl. extern calls) | ready | `compiler/check.go:2130-2146` |
| cycles (externs never form edges) | unaffected | `compiler/check.go:1352-1366` |
| **static admission** | **blocked** | `compiler/check.go:747-751` |
| **emit import source (caller's stem assumed)** | **blocked** | `compiler/emit.go:2236-2254` |

One test pins the locality rule: `TestDiagnoseForeignExtern`
(`compiler/lsp_test.go:1212-1224`), whose fixture calls without a
uses pin. `explain.go:167` states externs are "never shared".

## Decisions (JEV `jev-1.13.0`, Choice)

- Admission: `uses_pin` (0.99). A call to an extern declared in
  another module requires `name@rev` in `uses`, exactly like
  cross-module Can calls: exact-rev match, unpinned is
  `CodeCallNotInUses`. Same-module extern calls keep needing no
  pin (visible locality, unchanged).
- Host resolution: `declaring_stem` (0.84, runner-up `caller_stem`
  0.11). Emitted TS imports a shared extern from the declaring
  module's `.externs` stub: one host implementation owned by the
  declaring module, shared by all callers. Same-stem imports are
  byte-identical to today.

## Semantics

1. `checkCalls`: in the not-in-`prog.Fns` branch, a foreign
   extern with `prog.Uses[fname]` is admitted; without the pin it
   reports `CodeCallNotInUses` with the existing "add name@rev to
   uses" message. The "declare your own extern" error is retired.
2. No signature trust change: params/ret/emits come from the
   single `prog.Externs` declaration (double definitions already
   rejected, `compiler/check.go:193-210`).
3. Given/script/proof/eval unchanged: they already key off
   `CalleeForeign` + `EmitsOf` + `prog.Uses`.
4. Emit: `stemOf` gains extern names; the `.externs` import
   groups called externs by declaring stem — `./<decl>.externs`
   per foreign stem, `./<mod>.externs` for own. One line per
   stem, sorted, same shape as the foreign-fn import above it.
5. `explain.go` CAN3001 fix text drops "never shared".

## Non-goals

- Host implementations themselves (host shelf still needs
  writing once sharing lands).
- Transitive or version-migrating pins; rev discipline is
  unchanged.
- `std/host` module: a stdlib slice, not a language slice.

## Slices

- E1: admission + explain + tests (retire locality pin,
  pinned-pass and unpinned-`CodeCallNotInUses` fixtures).
- E2: declaring-stem emit + scripted-execution end-to-end +
  `stdlib-remaining.md` row flip to unblocked.

## Falsification

- `TestDiagnoseForeignExtern` must be updated, not deleted: its
  unpinned fixture becomes the `CodeCallNotInUses` case.
- Existing same-module extern goldens must be byte-identical
  after E2 (own-stem import line unchanged).
- A pinned foreign extern with a wrong rev must fail at
  `CodeUsesRev` before any call check runs.
