## Goal

Sequence the implementation slices for Schema's asset-approval minimum — the
machinery that unblocks the HTML Assets bundle (`html__asset__stylesheet`,
`html__asset__script`) — using the twelve trust-model rulings in
`docs/a83-astra-schema.md` as settled input. No slice may make a trust
decision; the plan turns each ruling into buildable, testable work.

## Success Criteria

- Both asset builders exist with the approved-only contract
  `(asset: ApprovedAsset, policy: AssetPolicy) → Html__Safe`: no raw-URL
  parameter, no unseal, no stringly approval, and an unapproved asset fails
  loudly at construction.
- A hostile-asset fixture set (unapproved URL, wrong hash, policy mismatch,
  registry removal, plus role mismatch and expired/revoked entries) rejects,
  each with its committed row and specified error kind.
- Rejection errors carry the offending asset only; no registry contents,
  expected digests, host lists, or alternate policies leak through CLI, LSP,
  catalog, or logs.
- Full suite green (`go test ./...`, `go run ./tools/modcheck`,
  `go run ./tools/gramcheck`); existing modules keep byte-identical behavior
  (Assets adapts to Schema, not vice versa).

## Context And Current Facts

- Gating passes: `docs/a83-astra-schema.md` covers all twelve
  `docs/schema-design-prompt.md` questions with ruling + rationale + failure
  rows and no `unresolved` markers. The required surface stays as the
  workstream scopes it: `ApprovedAsset` (opaque witness), `AssetPolicy`
  (per-program handle), registry (externally accepted snapshot), and the two
  builders.
- No `std/schema/` exists today. The std module convention (per
  `std/quota/`, `std/text/`) is `schema.ail` + `errors.json` + `schema.ts`
  (committed golden TS prod emit, regenerated via
  `go run ./compiler --out std/schema std/schema/schema.ail`) + `README.md`.
- Brand forgery is closed by `docs/a15-brands.md` (AIL6004): a function body
  may `seal` only its own file's brands; test expectations, arguments, and
  `given` rows may seal any declared brand. So only `schema.ail` bodies can
  mint `ApprovedAsset`, and only `html.ail` bodies can mint `Html__Safe`.
- Brands flow opaquely: there is no brand-to-string inspection
  (`std/html/html.ail:1-12`). A cross-module witness therefore cannot be
  consumed by ordinary calls — the Schema→HTML handoff needs an explicit
  compiler-authorized bridge, not a helper.
- Template for that bridge is Bytes B2 (`docs/bytes-plan.md` §1):
  owner-local grant (in B2, grant, brand, and function share one canonical
  `Module.ID`; the asset bridge in D2/S2 differs deliberately — dual
  authorization across two owners), exact-shape validation, checked
  authorization certificates,
  and the certificate lifecycle barrier (declare → check all modules → issue
  certificates → linkage evaluation → emit with the same certificates).
- All std modules currently declare `uses []`; cross-file reuse is scripted
  except where tests select linked-pure execution. `docs/a68-linked-pure.md`
  (`compiler/linked.go`, `Ctx.Linked`) dispatches cross-module pure calls to
  real bodies in `runLinkedPure` with an admission walk (no externs, state,
  effects, or unknown callees) — so Schema↔HTML interaction is testable
  without waiting on production linkage.
- `html__el__link` is a generic void-element wrapper over caller-supplied
  attributes; no asset builders exist, and no `Html__Url` brand exists (raw
  strings flow into href/src openly per `docs/a32-url-attributes.md`).

## Constraints And Non-goals

- The twelve rulings are input, not material: content digest AND permitted
  HTTPS origin with exact URL entry (single SHA-384 SRI digest, no
  wildcards/suffix/path-prefix); external Asset Approver principal with
  pinned key and no self-approval; one policy per program deployment; full
  binding (program + policy + registry snapshot + asset/revision + URL/digest
  + role + construction site + validity interval); construction-time
  authorization only; HTTPS-only with secure-document precondition; finite
  expiry + monotonic sequence + retained revocations + rotation by repinning;
  non-transitive approval; single exact role per approval; exact pins, no
  ranges; offending-asset-only errors; separate fixture trust root.
- Provenance machinery — principal enrollment, signing/acceptance
  configuration, trusted head/clock supply, secure document delivery,
  withdrawal of deployed artifacts — is owned by deployment operation and the
  acceptance workflow (a83 deferred table). Slices consume provenance as
  explicit inputs and never mint it.
- Non-goals (stay out of every slice): fetching, TLS/DNS enforcement,
  response/MIME/redirect handling, execution, HTTP/SQL/UI layers, inline
  script/style, closed dependency trees, runtime request mediation,
  page/DOM-bound fragments, additional roles or integrity profiles, sandbox
  or harmlessness claims.

## Key Decisions

- **D1 — Approval semantics in pure `.ail`, no compiler change.** Exact-match
  lookup over explicit snapshot values (records + `Seq`, string/int equality)
  is fully expressible. Keeps the security logic reviewable in source and the
  compiler out of policy. Rejected: native registry type or builtin
  verification (hides policy in the compiler, untestable as decision tables).
- **D2 — Two-owner bridge as a B2-shaped compiler grant.** Neither
  `html`-side raw-URL params (wall-smuggling, breaks the consumer contract)
  nor `schema`-side `Html__Safe` seals (cross-module seal, AIL6004) nor a
  generic `extern` declassifier (violates the nominal-projection caution:
  numeric projection must never become generic declassification) is
  acceptable. New declaration form (working name `asset_bridge`, spelling is
  implementer choice) authorizing exactly one pinned schema brand → one
  pinned html builder revision → one call node, with B2's identity,
  shape, and lifecycle rules.
- **D3 — Registry snapshot as an explicit value, never ambient state.**
  Lookup takes the snapshot (sequence number, entries, retained
  revocations) as a parameter. No quota-style state cell (ambient mutable
  authority would be poisonable from inside). This preserves hermetic
  evaluation and makes fixture vs production snapshots differ only in
  provenance — which is exactly the §12 boundary.
- **D4 — Policy, site, and time as compared values; trust of those values
  from outside.** In-language functions check equality/interval consistency
  (witness bindings vs supplied policy handle, site record, trusted time
  `T`); the acceptance workflow supplies the genuine articles. Until the
  provenance-certificate slice lands, builders stay test-gated and the plan
  says so explicitly.
- **D5 — One indistinguishable lookup error shape.**
  `schema.asset_not_approved` for every registry-dependent denial (unknown
  asset, wrong digest, wrong role/site grant, wrong policy, revoked,
  expired); `asset` carries the original request unchanged. Builder-level
  `html.asset_stylesheet_rejected` / `html.asset_script_rejected` carry the
  supplied witness unchanged. `AssetRegistryChangeRejected` is an
  administrative workflow result, never a callable `.ail` API.
  `SchemaAuthorityInvalid` is a compiler/acceptance diagnostic, never an
  `emits` outcome.
- **D6 — Builders emit fixed SRI output.** The granted builders carry the
  approved URL + exact integrity metadata into the element (stylesheet:
  `link` with integrity + `crossorigin`; script: `script` with integrity +
  `crossorigin`), with no caller overrides. Redirects, downgrade, and
  fallback behavior are documented retrieval obligations, not `.ail` rows.
- **D7 — Pure-first slice order.** Value model and lookup rows land before
  any compiler change, so the riskiest piece (D2 grant) is validated against
  frozen semantics rather than co-designed with them.

## Recommended Approach

Build the pure approval core first (S1), then prove the bridge with exactly
one builder (S2, the highest-risk slice), then mirror to the second builder
with role-disjointness rows (S3), then lifecycle/policy/site bindings with
the provenance-certificate design (S4), then the hostile fixture set,
fixture-root separation, and docs (S5). Each slice is independently
committable and green-gated; compiler changes are confined to S2 (grant) and
S4 (certificates/diagnostics).

## Work Plan

- **S1 — Schema value model + construction-time lookup (no compiler
  change).** New `std/schema/schema.ail`: `AssetRequest` record (identifier,
  revision, absolute URL, claimed SHA-384 SRI digest, requested role);
  `ApprovedAsset` + `AssetPolicy` opaque brands (sealed only in schema
  bodies); `RegistrySnapshot` record (sequence number, entries, retained
  revocations); `schema__asset__approve(request, snapshot, policy, site, t)`
  performing the exact-match chain — canonical-URL shape rejection
  (credentials, fragments, control characters, NUL, relatives), HTTPS-only,
  exact URL+digest+role match, expiry interval (`upper` exclusive),
  revocation/sequence check against supplied head, exact-pin match (no
  ranges/wildcards/`latest`). New kinds `schema.asset_not_approved`
  (`asset` = original request). Positive + negative `tests` rows per ruling
  (§§1, 6, 9, 10, 11); `errors.json`, `schema.ts` golden, `README.md`.
  Depends on: nothing. Proves: D1, D3, D5 lookup half.
- **S2 — Two-owner bridge grant + `stylesheet` builder (compiler change).**
  Compiler: `asset_bridge` declaration with B2-analogous validation
  (dual authorization, each side checked against its own canonical
  `Module.ID`, never by matching strings: schema-side release of the exact
  brand revision to the exact pinned builder + html-side sink authorization
  for that builder; exact builder revision; single `ApprovedAsset` +
  single `AssetPolicy` params;
  `Html__SafeResult`-shaped outcome; `emits` covers only the two builder
  rejections; body is one authorized construction, no string
  intermediates), certificate issuance after whole-program checking with the
  B2 lifecycle barrier (no trusted scripted evidence for missing authority),
  `SchemaAuthorityInvalid` diagnostics (missing trust root, illegal
  authority construction, invalid sink certification, malformed snapshot),
  grammar + LSP parity for the new form. `html`: `html__asset__stylesheet`
  as the single pinned sink, emitting the fixed `link` element with
  integrity + `crossorigin`. `runLinkedPure` vectors spanning
  schema+html files both module orders. Negative controls: ungranted sink
  body rejected at check; forged body seal still AIL6004. Depends on: S1
  types. Proves: D2, D6 (stylesheet half). Highest-risk validation in the
  plan.
- **S3 — `script` builder + role disjointness + non-transitivity rows.**
  Pin the second sink (`html__asset__script`, fixed `script` element with
  integrity + `crossorigin`); each builder supplies its own required role,
  no caller role flag. Rows: stylesheet witness → script builder fails
  `html.asset_script_rejected` and vice versa; classic-script witness spent
  as module/worker/preload unsupported; same bytes approved per-role only
  via separate entries; descendant-approval requests rejected
  (`schema.asset_not_approved`), "A and anything A loads" proposals
  rejected, no-import-certification requests rejected. Depends on: S2 grant
  pattern. Proves: §§8–9, D6 complete.
- **S4 — Lifecycle, policy/site bindings, provenance certificates.**
  In-language: validity-interval enforcement, monotonic-sequence vs
  supplied head (replayed old snapshot → `SchemaAuthorityInvalid`),
  retained-revocation records (revoked pin rejects; no resurrection by old
  snapshot; rotation = new revision + repinning, no silent retarget),
  conflicting-entries-for-one-revision rejects the snapshot, witness/policy
  snapshot-identity match (mixed contexts reject with the operation-specific
  builder kind), dev-policy-at-prod-site rejects. Compiler/acceptance
  design (paper first, then diagnostics): external approver enrollment,
  pinned trust root, proposer/approver separation records, per-program
  policy + snapshot selection immutable during construction, certified
  resolved call node (never caller-supplied site names), freshness recheck
  before deployment acceptance. Until certificates land, builders are
  test-gated — recorded in `std/schema/README.md`, never silent. Depends
  on: S1–S3. Proves: §§2–4, 7 and the D4 boundary.
- **S5 — Fixture-root separation, hostile set, docs, gates.** Harness:
  separate fixture trust root/registry; positive fixtures install explicit
  fixture snapshot values (a sealed success-shaped value alone never
  approves); `given` rows cannot override lookup; fixture-signed material
  under production keys → `SchemaAuthorityInvalid`. Committed hostile set:
  unapproved URL, wrong hash, policy mismatch, removal, role swap, expiry,
  replayed snapshot — each rejected with its row. Docs: schema README +
  decision records referenced from `docs/a83-astra-schema.md` without
  rewriting its rulings. Full gates:
  `go test ./...`, `go run ./tools/modcheck`, `go run ./tools/gramcheck`,
  TS emit regen + `tsc` per repo pins. Depends on: S1–S4. Proves: §12,
  acceptance criteria.

## Validation Plan

- S1: new `tests` rows green via `go test ./...`; hostile probes
  (host-only entry, digest-only entry, `http` URL with correct digest,
  protocol-relative/`data:`/`blob:` locators, range/`latest` selectors,
  NUL-bearing requests) each yield `schema.asset_not_approved` with the
  submitted request as payload; unknown-URL vs wrong-digest produce the
  same kind/shape (membership-oracle check by kind comparison).
- S2: `TestAssetBridge`-style falsifier — granted stylesheet construction
  succeeds; sink without grant fails the build; foreign body seal fails
  AIL6004 (negative control reusing the `sketches/broken-login` pattern);
  `runLinkedPure` schema+html vectors pass in both module orders;
  `modcheck`, `gramcheck`, golden `schema.ts`/`html.ts` byte-stable except
  intended additions.
- S3: cross-role matrices reject with the operation-specific kind carrying
  the supplied witness; diagnostic rendering of a rejected witness exposes
  no hidden authority internals.
- S4: expiry-boundary rows (valid at `t < expiry`, reject at
  `t == expiry`), revocation rows, replayed-snapshot diagnostic rows,
  mixed-snapshot (witness vs policy) rejection rows, dev-policy-at-prod
  rejection rows.
- S5: full hostile fixture set green; `given`-claims-approval control
  rejected (never authoritative success); production-acceptance with
  fixture-signed material rejected; whole suite + linters green;
  `git status` shows no changes to existing modules' behavior
  (additive-only diffs in `std/html`, none required elsewhere).

## Risks / Rollback

- **Bridge grant shape (S2) is the plan's top risk.** New declaration form
  touches grammar, check, eval, emit, and LSP. Mitigation: mirror the
  landed B2 implementation file-by-file; keep the grant asset-builders-only
  (never a general declassification facility). Rollback: revert compiler
  files to HEAD, delete grant tests — S1 stays green standalone.
- **Provenance certificates (S4) may lag the pure core.** If acceptance
  machinery is not ready, S1–S3 remain test-gated by explicit README
  notice; no slice ships a production trust claim early, and no permissive
  default is added to fill the gap.
- **Over-generalization pressure.** Each slice states its exact scope
  (two roles, one digest profile, construction-time only). Additional
  roles, profiles, or retrieval enforcement are new plans, not stretch
  goals of these slices.
- **Registry recon via diagnostics.** S2/S4 reviewers must grep new
  diagnostic strings for digests, host lists, and entry contents before
  merging; any leak fails the slice.

## Open Questions

None — every material choice above traces to `docs/a83-astra-schema.md`
rulings or to repo-verified mechanics (`docs/a15-brands.md`,
`docs/bytes-plan.md` §1, `docs/a68-linked-pure.md`, std module
convention). The `asset_bridge` keyword spelling is a non-material
implementer choice within the S2 shape constraints.
