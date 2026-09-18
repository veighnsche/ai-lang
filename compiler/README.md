# canlc — the can-lang transpiler (Go, stdlib only)

Stages: `parse.go` (scanner + file parser), `check.go` (static proofs:
naming, calls, given, emits, decreases, effects, cycles), `types.go`
(exact type discipline), `eval.go` (hermetic signature-test evaluator,
per-test stores), `emit.go` (TypeScript emitter), `lsp.go` (diagnose +
prove-first gate), `main.go` (CLI + program assembly), `catalog.go`
(`errors.json`), `code.go` (`CANnnnn` registry).

From repo root:

```
go build -o /tmp/canlc ./compiler
/tmp/canlc --out /tmp/can-out sketches/auth-login/db.can sketches/auth-login/auth.can
go test ./...
```

`go test` runs the golden gates (`emit_golden_test.go` covers
`auth-login/`, `retry-loop/`, `counter/` byte-identical) plus the
diagnosis suites (`lsp_test.go`, `arith_test.go`, `helper_test.go`,
`loop_test.go`, `effects_test.go`): any parse, proof, evaluation, or
emit change that alters output or diagnostics fails the build.

Editor mode: `canlc lsp` speaks minimal LSP over stdio (initialize,
didOpen/didChange → publishDiagnostics, shutdown/exit). The
Cursor/VSCode client in `editors/vscode/` installs to the editor's
extensions dir with a signed `bin/canlc`; rebuild + reinstall + sign
after compiler changes, then verify live against `sketches/`.
