# v65: interpreted string literals `e"..."` (chatbot verdict A)

Additive only. Ordinary `"..."` keeps raw semantics
(pinned by `TestStrSemanticsEmit` and the
`tail_backslash` row in `std/html/html.ail`: `"a\b"` is
backslash+b). New `e"..."` form decodes exactly six
escapes, once, left to right; same `str` type and runtime
representation, no new primitive. Renumbered from the
verdict's v62: v62/v63 are the diagnostic fixes, v64 is
taken.

1. Scan: recognize the `e"..."` token boundary in
   `stripComment`, `braceOutsideString`,
   `splitTopInner`, operator/bracket scanners, and
   balanced-delimiter scanning (audit each; the `e`
   prefix must not leak into identifiers).
2. `parseSmall`: decode a complete `e"..."` into
   `Small{Kind: "str", Str: decoded}` with ordinary
   embedding/postfix. `parsePattern`: same decoder
   (no value/pattern split).
3. Escapes: `\"` `\\` `\n` `\r` `\t` `\0` only.
   `e"\\n"` is [92,110]; `e"\01"` is [0,49] (no
   octal); unknown escapes, dangling backslash, and
   unterminated literals are AIL1000/`CodeParse`.
   No numeric/unicode escapes, no interpolation, no
   line continuation; `d"..."` untouched.
4. `patDesc` (lsp.go): preserve the source token/span
   for interpreted patterns (decoded newlines are not
   searchable in source).
5. Emit: `normStr` already covers the six values; no
   new mapping, no runtime unescape. Update comments
   claiming all AST strings are source-raw.
6. Grammar/tests/docs: interpreted-literal samples in
   `gramcheck`, parser/emit regressions, REQUIREMENTS
   compatibility amendment (canonical-spelling
   caveat: one value, two spellings, versioned
   explicit).

Out of scope: `\b \f \v \a` and numeric controls
(deferred in writing); v66 linked-pure runner; any
4107 relaxation (kept per verdict C).

## Rollback

`git checkout -- compiler/` plus delete
`compiler/str_escapes_test.go` and `docs/v65-*`.

## Test plan

- Probe first (red): `compiler/str_escapes_test.go`
  with the verdict's five utf8 rows (raw + escaped +
  decode-once + NUL + not-octal), the two base64
  reject rows (`e"QQ==\n"`, length-eight
  `e"QUJD\r\n\r\n"`), and the `e"\x51Q=="`
  AIL1000 rejection. Pre-fix the `e` forms do not
  parse.
- Keep green: `TestStrSemanticsEmit`,
  `tail_backslash`, full `go test -count=1 ./...`.
- Post-fix: committed `.ts` unchanged by test-row
  additions (tests stripped); error-payload survival
  verified by structural equality in the rows.
- `go run ./tools/modcheck`,
  `go run ./tools/gramcheck`.
