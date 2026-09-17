# html — constructor-controlled brands, starting with text

- `html.ail` — `mod html`: `brand Html__Text`, `html__text__escape`
  plus its `_from` worker. The worker scans scalar by scalar
  (front-consumption shape, `upper_ascii` precedent), emitting
  `&amp;`, `&lt;`, `&gt;` and passing everything else through,
  astral plane included. Raw strings enter `Html__Text` only
  through the entry; the brand erases to string and flows
  opaquely (never inspected, only passed to brand-typed params).
  NUL is rejected with `html.nul_byte` — the one scalar this
  serialization path cannot preserve as that scalar. All other
  controls pass through per HTML text semantics (each would need
  its own kind and justification to reject); quotes stay
  unescaped here by the text-context contract. See
  `docs/encoder-nul-policy.md`.
- `html__text__node` promotes `Html__Text` to the second brand,
  `Html__Safe is str rev 1 seals_from [Html__Text]` — a serialized
  fragment for ordinary child-fragment boundaries, with no authority
  for script, style, attribute, or URL contexts. Promotion is
  relabeling (`erase(node(t)) = erase(t)`): no double escape, no
  normalization. One-way, exact, same-module, non-transitive; see
  `docs/v26-html-node.md`.
- `html__attribute__name` gates `Html__TextAttributeName`: a small
  exact lowercase allowlist (currently `{title}`), each member with
  its own justification; everything else is `invalid_attribute_name`.
  See `docs/v27-attribute-name.md`.
- `html.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped; errors empty — the module is total).
  Regenerate: `go run ./compiler --out std/html
  std/html/html.ail`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/v25-html-text.md`,
`docs/v26-html-node.md`, `docs/v27-attribute-name.md`,
`docs/encoder-nul-policy.md`, brand scope: `docs/v15-brands.md`.
