// GENERATED from choice.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type ChoiceResult = { $can_kind: "ok"; value: Choice__Value$T$int } | { $can_kind: "ok"; value: Choice__Value$T$str } | { $can_kind: "choice.denied"; value: bigint };
export type Choice__Value$T$int = { $can_kind: "Choice__Some$T$int"; value: bigint } | { $can_kind: "Choice__None$T$int" };
export type Choice__Value$T$str = { $can_kind: "Choice__Some$T$str"; value: string } | { $can_kind: "Choice__None$T$str" };
export function choice__choose$T$int(present: boolean, value: bigint): { $can_kind: "ok"; value: Choice__Value$T$int } {
  if (present) {
    return { $can_kind: "ok", value: { $can_kind: "Choice__Some$T$int", value: value } };
  }
  else {
    return { $can_kind: "ok", value: { $can_kind: "Choice__None$T$int" } };
  }
}
export function choice__choose$T$str(present: boolean, value: string): { $can_kind: "ok"; value: Choice__Value$T$str } {
  if (present) {
    return { $can_kind: "ok", value: { $can_kind: "Choice__Some$T$str", value: value } };
  }
  else {
    return { $can_kind: "ok", value: { $can_kind: "Choice__None$T$str" } };
  }
}
export function choice__checked(value: bigint): { $can_kind: "ok"; value: Choice__Value$T$int } | { $can_kind: "choice.denied"; value: bigint } {
  const $can_m1 = (value < 0n);
  const $can_m2 = (value === 0n);
  if ($can_m1) {
    return { $can_kind: "choice.denied", value: value };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: { $can_kind: "Choice__None$T$int" } };
  }
  else {
    return { $can_kind: "ok", value: { $can_kind: "Choice__Some$T$int", value: value } };
  }
}
