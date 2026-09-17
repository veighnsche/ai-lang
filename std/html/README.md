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
  exact lowercase allowlist (currently `title`), each member with
  its own justification; everything else is `invalid_attribute_name`.
  See `docs/v27-attribute-name.md`.
- `html__attribute__text` composes a brand-typed name with an
  encoded value into `Html__Attribute`, canonical single-quoted
  form (`title='...'`). The spelling is reconstructed from the
  closed admitted domain, not extracted from the brand. See
  `docs/v29-attribute-value.md`.
- `html__attribute__boolean_name` gates
  `Html__BooleanAttributeName` (currently `disabled`,
  `readonly`, `required`, `checked`);
  `html__attribute__boolean` serializes presence as the spelling
  and absence as the empty contribution. Syntactic guarantee
  only — never inertness, never element applicability. See
  `docs/v30-boolean-attribute.md`.
- `html__attribute__id` validates an identifier (nonempty, no
  ASCII whitespace) and serializes it as `id='...'`, reusing the
  shared value worker. Uniqueness needs a tree and stays out.
  New error `html.invalid_identifier`. See
  `docs/v31-identifier-attribute.md`.
- `html__attribute__href` / `html__attribute__src` validate
  absolute-`https` URLs under a restricted ASCII authority
  profile and serialize them reusing the shared worker. Fused
  raw input: no `Html__Url` brand exists yet because no
  consumer could serialize one. New errors `html.invalid_url`,
  `html.disallowed_scheme`. See `docs/v32-url-attributes.md`.
- `html__fragment__empty` seals `""` as the zero `Html__Safe`
  fragment. See `docs/v33-empty-fragment.md`.
- `html.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped; errors empty — the module is total).
  Regenerate: `go run ./compiler --out std/html
  std/html/html.ail`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/v25-html-text.md`,
`docs/v26-html-node.md`, `docs/v27-attribute-name.md`,
`docs/v29-attribute-value.md`, `docs/v30-boolean-attribute.md`,
`docs/v31-identifier-attribute.md`, `docs/v32-url-attributes.md`,
`docs/v33-empty-fragment.md`, `docs/v34-boolean-names.md`,
`docs/encoder-nul-policy.md`, brand scope: `docs/v15-brands.md`.
