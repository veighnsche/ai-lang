// GENERATED from option.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type OptionResult = { $can_kind: "ok"; value: Option__Value$T$int } | { $can_kind: "ok"; value: Option__Value$T$str } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: string };
export type Option__Result$T$int = { value: bigint };
export type Option__Result$T$str = { value: string };
export type Option__Box$T$int = { value: Option__Value$T$int };
export type Option__Box$T$str = { value: Option__Value$T$str };
export type Option__Value$T$int = { $can_kind: "Option__Some$T$int"; value: bigint } | { $can_kind: "Option__None$T$int" };
export type Option__Value$T$str = { $can_kind: "Option__Some$T$str"; value: string } | { $can_kind: "Option__None$T$str" };
export function option__value_or$T$int(value: Option__Value$T$int, alternative: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Option__Some$T$int": {
    const s = $can_m1;
    return { $can_kind: "ok", value: s.value };
  }
  case "Option__None$T$int": {
    const _ = $can_m1;
    return { $can_kind: "ok", value: alternative };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function option__value_or$T$str(value: Option__Value$T$str, alternative: string): { $can_kind: "ok"; value: string } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Option__Some$T$str": {
    const s = $can_m1;
    return { $can_kind: "ok", value: s.value };
  }
  case "Option__None$T$str": {
    const _ = $can_m1;
    return { $can_kind: "ok", value: alternative };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function option__some$T$int(value: bigint): { $can_kind: "ok"; value: Option__Value$T$int } {
  return { $can_kind: "ok", value: { $can_kind: "Option__Some$T$int", value: value } };
}
export function option__some$T$str(value: string): { $can_kind: "ok"; value: Option__Value$T$str } {
  return { $can_kind: "ok", value: { $can_kind: "Option__Some$T$str", value: value } };
}
export function option__none$T$int(alternative: bigint): { $can_kind: "ok"; value: Option__Value$T$int } {
  return { $can_kind: "ok", value: { $can_kind: "Option__None$T$int" } };
}
export function option__none$T$str(alternative: string): { $can_kind: "ok"; value: Option__Value$T$str } {
  return { $can_kind: "ok", value: { $can_kind: "Option__None$T$str" } };
}
