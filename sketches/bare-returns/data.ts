// GENERATED from data.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type DataResult = { $can_kind: "ok"; value: (input: Uint8Array) => { $can_kind: "ok"; value: Uint8Array } } | { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } | { $can_kind: "ok"; value: (input: string) => { $can_kind: "ok"; value: string } } | { $can_kind: "ok"; value: Data__Row[] } | { $can_kind: "ok"; value: Uint8Array } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: bigint[] } | { $can_kind: "ok"; value: string };
export type Data__Row = { number: bigint };
export type Data__Choice$T$int = { $can_kind: "Data__Some$T$int"; value: bigint } | { $can_kind: "Data__None$T$int" };
export function data__mint(value: string): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: value };
}
export function data__bytes(value: Uint8Array): { $can_kind: "ok"; value: Uint8Array } {
  return { $can_kind: "ok", value: value };
}
export function data__items$T$Data__Row(values: Data__Row[]): { $can_kind: "ok"; value: Data__Row[] } {
  return { $can_kind: "ok", value: values };
}
export function data__items$T$int(values: bigint[]): { $can_kind: "ok"; value: bigint[] } {
  return { $can_kind: "ok", value: values };
}
export function data__add(delta: bigint, value: bigint): { $can_kind: "ok"; value: bigint } {
  return { $can_kind: "ok", value: (value + delta) };
}
export function data__id$T$Bytes(value: Uint8Array): { $can_kind: "ok"; value: Uint8Array } {
  return { $can_kind: "ok", value: value };
}
export function data__id$T$Data__Token(value: string): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: value };
}
export function data__reference$T$Bytes(): { $can_kind: "ok"; value: (input: Uint8Array) => { $can_kind: "ok"; value: Uint8Array } } {
  return { $can_kind: "ok", value: (() => { return (input) => data__id$T$Bytes(input); })() };
}
export function data__reference$T$Data__Token(): { $can_kind: "ok"; value: (input: string) => { $can_kind: "ok"; value: string } } {
  return { $can_kind: "ok", value: (() => { return (input) => data__id$T$Data__Token(input); })() };
}
export function data__factory(delta: bigint): { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } {
  return { $can_kind: "ok", value: (() => { const $can_m1 = delta; return (input) => data__add($can_m1, input); })() };
}
