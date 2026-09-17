# ai-lang

A small contract-first language (`.ail`) that transpiles to TypeScript.
You write specs with named errors, exchange scripts, and proofs
(termination, effects, exact numerics); the compiler checks them and
emits typed TS plus machine-readable artifacts.

## Layout

- [`REQUIREMENTS.md`](REQUIREMENTS.md) — the language contract, versioned by amendment
- [`docs/`](docs/README.md) — version history (`a02`–`a12`), the audit trail, and reviewer notes
- [`compiler/`](compiler/README.md) — `ailc`, the Go transpiler (stdlib only)
- [`sketches/`](sketches/README.md) — example programs (auth-login, retry-loop, counter, broken-login)
- [`std/`](std/README.md) — the blessed standard library (quota, scalars)
- [`editors/vscode/`](editors/vscode/README.md) — syntax highlighting + LSP client
- [`tools/`](tools/) — grammar and module checkers

## Install

Requires Go 1.21+. Easiest, no clone needed:

```
go install github.com/veighnsche/ai-lang/compiler@latest
mv "$(go env GOPATH)/bin/compiler" ~/.local/bin/ailc
```

(`go install` names the binary after the package directory; the `mv`
gives it its real name. Make sure `~/.local/bin` is on your `PATH`.)

Or from a clone:

```
make install   # builds ./bin/ailc and copies it to ~/.local/bin/ailc
```

## Upgrade

Re-run whichever install you used — both resolve to the newest commit:

```
go install github.com/veighnsche/ai-lang/compiler@latest   # then re-mv to ailc
```

```
git pull && make install
```

`ailc --version` prints the build stamp (`make install` stamps the git
revision; plain `go install` reports `dev`). No tags or releases yet,
so "latest" means latest `main` — check the
[commits](https://github.com/veighnsche/ai-lang/commits/main) to see
what changed. Note: right after a push, `@latest` can lag the Go module
proxy by a few minutes; to upgrade immediately, pin the commit instead:
`go install github.com/veighnsche/ai-lang/compiler@<sha>`.

## Quickstart

```
ailc --out /tmp/ail-out sketches/auth-login/db.ail sketches/auth-login/auth.ail
```

From the repo root, `go test ./...` runs the golden gates
(byte-identical emit for the gallery sketches) plus the diagnosis
suites: any parse, proof, evaluation, or emit change that alters output
or diagnostics fails the build.

`go test` runs the golden gates (byte-identical emit for the gallery
sketches) plus the diagnosis suites: any parse, proof, evaluation, or
emit change that alters output or diagnostics fails the build.

## Status

Through **a12**: unbounded exact numerics (bigint/decimal agreement),
program-wide guarded termination with a returned-outcome theorem, and
producer-owned contracts with complete observations. See
[`docs/a10-numerics.md`](docs/a10-numerics.md),
[`docs/a11-recursion.md`](docs/a11-recursion.md),
[`docs/a12-contracts.md`](docs/a12-contracts.md), and the
[`docs/ASTRA_AUDIT.md`](docs/ASTRA_AUDIT.md) trail that drove them.
