# ai-lang editor support (Cursor / VSCode)

Two halves: `syntaxes/` colors (TextMate), `client/` + `bin/ailc` squiggles
(LSP client talking to `ailc lsp` over stdio).

Install (from the repo root; manual copy into the extensions dir does
NOT register — the editor only loads manifest-registered extensions):

```
go build -o editors/vscode/bin/ailc ./compiler
(cd editors/vscode && npm install --no-audit --no-fund)
(cd editors/vscode && npx -y @vscode/vsce package -o /tmp/ai-lang.vsix)
cursor --install-extension /tmp/ai-lang.vsix --force   # or: code --install-extension ...
codesign --force --sign - ~/.cursor/extensions/ai-lang.ai-lang-*/bin/ailc
```

(`codesign` re-signs the server so macOS runs it; without this the OS
kills it with "Code Signature Invalid".) Then Developer: Reload Window.
Open any `sketches/*/*.ail` or `std/*/*.ail` file. Broken files under
`sketches/broken-login/` should show red/yellow squiggles; clean files
show none.

See a squiggle and disagree? The diagnosis comes from `compiler/lsp.go`
(`diagnose`), proven by `compiler/lsp_test.go` — fix it there, rebuild
`bin/ailc` (`go build -o editors/vscode/bin/ailc ./compiler`), reinstall.

Settings: `ailc.serverPath` overrides the server binary (default: bundled
`bin/ailc`). Token colors (e.g. forcing `error` red) live in the user's
`settings.json` via `editor.tokenColorCustomizations`:

```json
"editor.tokenColorCustomizations": {
  "textMateRules": [
    {
      "scope": "keyword.declaration.error.ail",
      "settings": { "foreground": "#F14C4C", "fontStyle": "bold" }
    }
  ]
}
```
