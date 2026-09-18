# schema — construction-time asset approval (S1: pure core)

- `schema.ail` — `mod schema`: the asset-approval value model and lookup.
  `Schema__AssetRequest` is untrusted identifying data (plain record, no
  authority). `Schema__ApprovedAsset` is the opaque witness, sealed only in
  this file's bodies and only on the exact-match path of
  `schema__asset__approve`. `Schema__AssetPolicy` is the opaque policy
  handle; the snapshot carries the expected handle and approval checks
  agreement by brand equality, never by inspection. The snapshot, site,
  head sequence, and time arrive as explicit values — their provenance is
  the deployment acceptance workflow's job (S4), not this file's.
- Checks, cheapest first (all failures collapse to the one kind below):
  closed role gate (`stylesheet` | `script`), absolute-`https` prefix,
  tail scan (no NUL, fragment, backslash, whitespace, or userinfo),
  single-SHA-384 SRI form (`sha384-` + 64 base64 scalars), policy-handle
  agreement, program agreement, head-sequence currency, exactly-one entry
  match, revocation walk, validity interval (upper bound exclusive).
- `error schema.asset_not_approved(asset: Schema__AssetRequest)` is the
  only lookup failure: unknown asset, wrong digest, wrong role/site,
  wrong policy, revoked, expired, stale, and conflicting snapshots all
  produce this same shape carrying the original request unchanged. No
  registered digest, host list, or alternate entry ever leaves the module.
- NUL is rejected, never stripped; the `tail_nul` row carries a literal
  NUL byte (invisible — verify with `tr -d -c '\000' | wc -c`).
- `schema.ts` + `errors.json` — committed golden TS prod emit. Regenerate:
  `go run ./compiler --out std/schema std/schema/schema.ail`; verify:
  `go test ./...`.
- Separator safety (S2): the witness joins its eight fields with `|`
  for the projection kernel, so `schema__tokens__check` refuses `|`
  (and NUL) in ids, revisions, and site fields; the URL tail and the
  digest alphabet exclude it too. A field holding the separator would
  make the kernel split ambiguous, so approval refuses it.

Rules: `docs/a83-astra-schema.md` (trust rulings, §§1, 6, 9–11 in this
slice). Plan: `.agents/plans/2026-09-18-schema-asset-slices.md` (S1).
Builders (`html__asset__stylesheet`, `html__asset__script`), the two-owner
bridge grant, and provenance certificates arrive in S2–S4; until then this
core is test-gated and makes no production trust claim.
