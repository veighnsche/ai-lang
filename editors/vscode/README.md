# ai-lang editor support (Cursor / VSCode)

Two halves: `syntaxes/` colors (TextMate), `client/` + `bin/ailc` squiggles
(LSP client talking to `ailc lsp` over stdio).

Install: copy this folder's contents (including `node_modules/` and `bin/`)
into the editor's extensions dir as `ai-lang.ai-lang-<version>`, then
re-sign the server so macOS runs it (`codesign --force --sign - <extdir>/bin/ailc`;
without this the OS kills it with "Code Signature Invalid"), then
Developer: Reload Window. Open any `sketches/*/*.ail` file. Broken files
under `sketches/broken-login/` should show red/yellow squiggles; clean
files show none.

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
