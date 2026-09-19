# b03 — `seq[i].field` projection design

## Problem

Indexing yields an element but nothing projects off it:
`m[i].field` dies in the postfix loop with `unexpected ".field"
after ]` (`compiler/parse.go:1443`), because field projection
exists only as `ref` paths resolved from environment names
(`compiler/types.go:330`). Slice 2 worked around it via `_push`
(`docs/stdlib-remaining.md:41`).

## Representation

A new `Small` kind `proj` with `L` = base expression and
`Field` = one segment. Chains nest left (`m[i].a.b` is
`proj(proj(stridx, a), b)`); the postfix loop alternates brackets
and `.word`, so `m[i].f[j]` parses too. No desugar into `ref` is
possible: `ref` paths resolve from names, never computed bases.

## Site table (one case per existing `stridx` site)

| Site | `proj` behavior |
|------|-----------------|
| `parse.go` postfix loop | accept `.word` after `]`; strict ident, else `bad projection` |
| `types.go` `typeOf` | field type from `typeOf(L)` via `recs`; silent on unknown |
| `types.go` `value` | visit base; non-record → `cannot project`; missing → `no field f on T` (mirrors ref, `types.go:711`); set `s.T`; the existing `want` tail compares |
| `types.go` `tokenOf` | `Field` |
| `eval.go` `evSmall` | `vField(evSmall(L), Field)` (fail-closed, mirrors ref) |
| `emit.go` `emitValue` | `<base>.<field>` (bases are calls/refs/postfix; unions are structural) |
| `check.go` `walkSmallTrees`, `catalog.go` `stubKindsOf`, `eval.go` `walkCalls` | descend `L` |
| `revision.go` `canonSmall` | `proj(<base>,<field>)` |
| `verify_admission.go` | `outside`, like `stridx` (its base contains an index) |
| `check.go` `firstTrappingOp`, `lint.go` `lintPureSmall` | no entry: the inner index is found by the walk; `proj` itself is total |

Codes mirror existing rules: `CodeTypeMismatch` for bad base
(like `stridx`) and unknown field (like `ref`). No new code.

## Non-goals

- Projection off non-bracket bases (`(f()).x`, `call f(x).y`):
  the parser accepts `.word` only after `]` in the postfix loop.
- `err:`/`case` bases: index results are plain annotations, so
  `recs` lookup is complete.
- Slices of records (`m[a:b].f`): slices are str-only, so the
  slice itself rejects first and `proj` stays silent — same
  rule, no special case.
