# text — explicit text construction and scalar access

- `text.ail` — `mod text`: `std__str__concat`, `is_empty`,
  `is_whitespace` (int-denoted scalars: literals are raw, so
  tab/LF/CR never appear in a table), `length_scalars`,
  `scalar_at`, `slice_scalars`, `find_from`/`find`, `contains`,
  `starts_with`, `ends_with`, `replace_all_from`/`replace_all`
  (non-overlapping, left to right), `trim_left`/`trim_right`/
  `trim_ascii` (space, tab, LF, CR), and `upper/lower_ascii`
  (`_from` workers over alphabet literals via same-file `find`).
  Concat preserves every scalar; safety lives in context
  encoders, never here. `find` locates the empty pattern at 0;
  `replace` rejects it. Split/join wait on Seq, graphemes/
  casefold/normalize on pinned data, hex/base64/utf8 on Bytes.
- `text.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped). Regenerate: `go run ./compiler --out
  std/text std/text/text.ail`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/v16-text.md`,
surface: `docs/v20-text-operators.md` (`#`, `s[i]`, `s[a:b]`).
