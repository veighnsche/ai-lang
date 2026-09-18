# sketches — can-lang use-case examples, one folder per example

- `auth-login/` — password login with one retry (`db.can` + `auth.can`).
  Clean: opening these shows zero diagnostics.
- `retry-loop/` — retry with fuel (`retry.can`): the termination +
  iteration demo. The loop is a proven self-call; minus `decreases`
  it is a cycle error. Clean: zero diagnostics.
- `counter/` — a bounded counter (`counter.can`): the effects
  demo. One private cell, declared capabilities, per-test
  isolation proved by the tables. Clean: zero diagnostics.
- `lint-errors/` — one proving fixture per `canlc lint` rule
  (`chain`, `merge`, `outcome`, `ranges`, `redundant`, `relay`,
  `table`): every file compiles and passes, and every file lints
  dirty. Open them with the editor extension installed to see each
  lint error live (red, error-severity — the linter has no warnings).
- `broken-login/` — same domain, one broken file per squiggle (`db.can`
  provider plus `missing-tests`, `bad-name`, `missing-rev`,
  `failing-test`, `missing-arm`, `dangling-test`, `dead-script`,
  `foreign-raise`, `unknown-call`, `unknown-emits`,
  `incomplete-expectation`, `missing-exchange`). Open them with the
  editor extension installed to see each diagnostic live.
  Each example folder is named by its contents.

Clean means compiler-clean AND lint-clean: gallery sketches carry
no findings of either kind (`canlc` passes, `canlc lint` exits 0).

Rules: `/REQUIREMENTS.md`. Design rationale: `CLEAN_ROOM_REVIEW.md`.

Checks: `go run ./tools/modcheck` (every `uses` resolves to
another file's `provides` + pins a rev; no `extern` re-declaration, no
`externals` sections, no curly braces, every `given` table total over
its tests).
