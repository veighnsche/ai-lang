# json — JSON value AST, render, scalar codecs, schemas

- `json.can` — `mod json`: `Json__Value` is the parsed tree
  (`Null`/`Bool`/`Num`/`Str`/`Arr`/`Obj`; `Num` carries the wire
  text verbatim, never a lossy float). `std__json__render_value`
  prints any tree through a single self-recursive frame machine:
  only direct self-recursion is admitted, so the natural
  value/fields/array mutual recursion is inexpressible and an
  explicit `Seq<Json__Frame>` stack carries the pending work
  (tag-dispatched records — variant sequences are not admitted;
  stack top is the back since sequence concatenation is not in
  v1). Each step burns one fuel from `std__json__RENDER_FUEL`
  (1M); exhaustion raises payloadless
  `json.render_budget_exhausted`. `std__json__escape` handles
  `"`/`\`/TAB/LF/CR plus `\u00XX` for other C0 controls via an
  armless hex table. Typed codecs bridge scalars to `Json__Doc`
  (`int`/`str`/`bool`/`dec` × encode/decode); decodes reject
  mistyped trees with `json.schema_mismatch` and bad numeric
  text with `json.numeric_out_of_range` (translated from
  `convert.*`). Schemas are a monomorphic family (JEV 0.74 over
  generic-uniform): `Json__IntSchema` & co. carry exact
  encode/decode callbacks invoked through bare-name apply
  wrappers (`invoke` heads must be bare names).
- `json.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped). Regenerate: `go run ./compiler --out
  /tmp/jg std/json/json.can std/scalars/scalars.can`, copy
  `json.ts` + `errors.json`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/ASTRA_STDLIB.md` §1.8.
Byte-level parse (`invalid_syntax`, `duplicate_key`) and the
`std__json__encode`/`decode` drivers land in a later slice.
