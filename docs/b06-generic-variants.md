# B06 — explicit generic variants

Status: implemented; commit recorded in `stdlib-remaining.md`.
B05 remains reserved for structured async. This slice extends the existing
explicit record/function generics, not outcomes or error-set algebra.

## Surface

```can
variant Option__Value<T> rev 1 (
  case Some(value: T)
  case None()
)
```

- Annotations: `Option__Value<int>`.
- Construction: `Option__Some<int>(value = 7)`, `Option__None<int>()`.
- Matching: `on Option__Some<T> s`, `on Option__None<T> _` inside a
  function parameterized by `T`; use concrete arguments outside templates.
- Every case must be covered and witnessed **per function instance**.
  A wildcard does not replace an unwitnessed case.
- Header pins name the parent (`Option__Value@1`), not the constructors.
  Case qualification still uses the parent's domain before its first `__`.
  Cases with the same source name in two parents still collide, even when
  the parents happen to instantiate at disjoint argument types.

Case arguments are explicit even for nullary cases. JEV (`jev-1.13.0`,
Choice) selected explicit patterns over contextual inference or a
parse-only slice: probability 1.00, confidence 0.99.

## Implementation

The existing monomorphizer collects records and variants together. A
common read-only field view validates parameters and discovers dependencies;
variant stamping preserves the case boundaries. Function instances substitute
arguments in patterns as well as constructors. Pattern heads participate in
the expansion site walk, including contract match trees.

`Option__Value<int>` becomes `Option__Value$T$int`; its cases become
`Option__Some$T$int` and `Option__None$T$int`. Both nominal types and runtime
tags distinguish instances. Templates never reach checking or emission.
The existing checker, evaluator, coverage proof, revision fingerprinting,
and TypeScript emitter consume ordinary monomorphic variants.

Commas within multi-argument patterns are protected only in pattern lists;
expression comparison scanning is unchanged. Source identity collisions are
checked before stamping can hide them. Headers are rewritten using the same
per-module demand tracking as generic records.

## Retained limits

- Arguments are closed monomorphic types (including the existing permitted
  `Seq<X>` shape), or bare parameters passed through a generic function.
  Nested instantiation as an argument, such as `Box<Option__Value<int>>`,
  remains rejected. A field `option: Option__Value<T>` is supported.
- No inference or implicit conversion between instances.
- Parameters must occur in payload fields; unused templates are rejected.
- Direct self-recursive payloads and variant sequences remain rejected.
- Bare variant returns remain rejected: `bare-variant returns are
  unsupported, return a record`. Use a record carrier in language examples.
  This sketch is not a claim to implement the catalogue's exact option APIs.
- No builtin `Option<T>` alias, first-class `Outcome<T,E>`, generic error
  algebra, or `Ok` record splat is introduced.

## Gates and reproducible example

`sketches/generic-option` supplies a generic optional value, record carriers,
constructors, a fallback operation, and a consumer module: 12 rows covering
both cases at both `int` and `str`. It is intentionally a language sketch,
not a new `std/option` module.

```sh
go test -count=1 ./...
go run ./tools/modcheck
go run ./tools/gramcheck
go run ./compiler --out /tmp/generic-option \
  sketches/generic-option/use.can sketches/generic-option/option.can
```

`TestGoldenGenericOption` compares both emitted `.ts` files and `errors.json`
byte-for-byte. `TestGenericOptionNodeParity` checks the generated runtime,
including a bigint above JavaScript number precision. Compiler tests also
cover multiple parameters, record/variant dependency discovery, pins,
source collisions, invalid arguments, nominal instance mismatches,
exhaustiveness, per-instance witnesses, and the retained boundaries.
