# b04 — nested generic Seq args design

## Problem

`Map<K,V>` fully generic needs `Seq<M__Pair<K,V>>` entries
(`docs/stdlib-remaining.md:32`). Surveyed 2026-09-19: single-arg
nesting (`Seq<M__Box<int>>`, including inside generic templates
and generic fns, plus `M__Box<Seq<int>>`) already works end to
end through the G2 expansion — check, eval, and emit. Only the
multi-arg literal head failed: `Seq<M__Pair<str, int>>[...]` died
with `bad sequence literal Seq<M__Pair<str: ...`.

## Root cause

`splitTopInner` protects `<...>` spans from the separator only
when `genericHeadLT` fires, which required `(` after the balanced
`>` (`compiler/parse.go:1596`). A Seq literal head ends in `[`,
so the comma inside multi-arg element types split the argument
list before `parseSeqLit` ran. (`reGenericElem` already admits
commas; the head never reached it intact.)

## Fix

`genericHeadLT` also fires on `[` after a head that itself nests
angles. The nesting requirement keeps the change a no-op for
plain `Seq<str>[` heads (no comma to protect) and `<=`
comparisons (no nested `<`), and keeps the `chainElseMentions`
gate (`compiler/expand.go:344`) tripping only on genuinely
generic heads. One spot; both callers inherit it.

## Non-goals

- `Seq<Seq<..>>` stays rejected at parse (`reGenericElem`
  excludes the `Seq` base; pinned by
  `compiler/generics_parse_test.go:276`). Plain sequence
  nesting is a separate shape with no stdlib customer today.
- `Map<K,V>` itself: a stdlib slice now unblocked, not a
  language slice.
