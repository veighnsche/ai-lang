# Seam-fix loop — one bolted-on finding at a time

Run this loop for exactly ONE finding per session, then STOP. The session
goal owns the queue; after you stop, the next session picks up the next point.

## The 12-point queue (in order)

1. Kernel/call dispatch ladder — one parameter codec path (`emit`, `eval`, `export`, `check`, `linked`).
2. Linked/scripted call duality — one call semantics with explicit mode.
3. `$can` runtime + decimal duplication — one decimal semantics module.
4. Result unions — single Result-type rule (per-function vs module-wide).
5. Chain as extra pass — parse ladders natively.
6. Admit-gate-cover per feature — admit natively, one coverage rule.
7. Equality/ctor/literal triple rules — one rule each, shared by all stages.
8. Verifier gaps — admit only what the prover discharges.
9. Baseline/pins — one baseline story.
10. Diag/LSP strings — typed codes, one suppression path, CLI/LSP parity.
11. Payload/registry — uniform payload + sequenced registry.
12. Tooling/CLI/stubs — one shared impl, generated stubs.

Full evidence refs live in the session goal objective.

## Loop steps (one finding)

1. **Research.** Open every implementation AND test body cited for the
   finding. Search output only locates candidates. Write down the seam
   symptoms (special-case branches, duplicated paths, naming warts) with
   file:line refs. Never let a test you write define correctness; derive
   the contract from existing callers and sibling code.
2. **Brainstorm.** Propose 2–3 ways to make it seamless. Score each on
   blast radius (files touched, golden churn), risk to gates
   (`go test ./...`, modcheck, gramcheck, tsc), and reversibility.
   Pick the smallest correct fix at the root cause.
3. **Plan.** Write the implementation plan as a short checklist: files to
   change, contract to preserve, goldens that will move, gates to run.
4. **Implement.** Make the change with the editing tools. Reproduce the
   seam first where feasible; keep the diff minimal; no unrelated edits.
   Regen goldens only where the new behavior requires it (`errors.json`
   must stay byte-identical unless the finding owns it).
5. **Clean-room review.** Spawn ONE subagent with a falsification brief:
   "break this change or find contract violations; read the bodies, do not
   trust my summary." Give it the diff and the finding. Wait for its result.
6. **Fix issues.** Address every review hit or record why it is not a hit.
   Re-run the full gates (`go test ./...`, modcheck, gramcheck, tsc where
   touched) after each fix round.
7. **Commit.** `git status` to confirm only intended files changed, then
   commit with a message naming the finding and the seamless direction.
   Never amend, push, rebase, or touch history.
8. **STOP.** Report progress on the session goal (percent + next point),
   end the turn. Do NOT start the next point.
