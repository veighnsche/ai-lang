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
Byte-level parse landed in S11b: `std__json__parse_value`
guards empty input, then `std__json__parse_step` runs a
single self-recursive 12-state machine over `Seq<Json__PFrame>`
(empty dispatch + `ArrFirst`/`ArrVal`/`ArrNext` + `ObjFirst`/
`ObjKey`/`ObjColon`/`ObjKeyVal`/`ObjNext` + `StrKey`/`StrVal`/
`NumAcc` + `Tail`). Value states push nested frames; key and
punctuation continuations replace the top frame (a lingering
`ObjFirst` under `StrKey` breaks `parse_attach`, found
in-slice). `StrKey` inherits the parent's fields/keys for
duplicate detection; numbers accumulate raw and validate
through the 9-state numcheck; `Tail` rejects trailing values.
Only the `std__json__encode`/`decode` text drivers remain
(slice 12).
