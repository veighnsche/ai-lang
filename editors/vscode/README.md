# can-lang editor support (Cursor / Antigravity / VSCode)

Two halves: `syntaxes/` colors (TextMate), `client/` + `bin/canlc` squiggles
(LSP client talking to `canlc lsp` over stdio).

Install (from the repo root; manual copy into the extensions dir does
NOT register — the editor only loads manifest-registered extensions):

```
go build -o editors/vscode/bin/canlc ./compiler
(cd editors/vscode && npm install --no-audit --no-fund)
(cd editors/vscode && npx -y @vscode/vsce package -o /tmp/can-lang.vsix)
cursor --install-extension /tmp/can-lang.vsix --force   # or: code --install-extension ...
antigravity-ide --install-extension /tmp/can-lang.vsix --force
codesign --force --sign - ~/.cursor/extensions/can-lang.can-lang-*/bin/canlc
codesign --force --sign - ~/.antigravity-ide/extensions/can-lang.can-lang-*/bin/canlc
```

(`codesign` re-signs the server so macOS runs it; without this the OS
kills it with "Code Signature Invalid".) Then Developer: Reload Window.
Open any `sketches/*/*.can` or `std/*/*.can` file. Broken files under
`sketches/broken-login/` should show red/yellow squiggles; clean files
show none.

See a squiggle and disagree? The diagnosis comes from `compiler/lsp.go`
(`diagnose`), proven by `compiler/lsp_test.go` — fix it there, rebuild
`bin/canlc` (`go build -o editors/vscode/bin/canlc ./compiler`), reinstall.
Strictness squiggles (redundant names, mergeable arms, chainable
ladders, tableable nests, same-outcome matches, handwritten relays,
mergeable ranges — CAN3410–3416) come from `compiler/lint.go` via the
same `diagnose` path, proven by `compiler/lint_test.go` (including the
`sketches/lint-errors/` folder golden). They publish only on files the
compiler otherwise accepts; broken files show their real errors first.

Settings: `canlc.serverPath` overrides the server binary (default: bundled
`bin/canlc`). Token colors (e.g. forcing `error` red) live in the user's
`settings.json` via `editor.tokenColorCustomizations`:

```json
"editor.tokenColorCustomizations": {
  "textMateRules": [
    {
      "scope": "keyword.declaration.error.can",
      "settings": { "foreground": "#F14C4C", "fontStyle": "bold" }
    }
  ]
}
```
