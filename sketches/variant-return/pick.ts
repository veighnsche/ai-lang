// GENERATED from pick.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type PickResult = { $can_kind: "ok"; value: Pick__Value$T$int } | { $can_kind: "ok"; value: Pick__Value$T$str } | { $can_kind: "pick.denied"; value: bigint };
export type Pick__Value$T$int = { $can_kind: "Pick__Some$T$int"; value: bigint } | { $can_kind: "Pick__None$T$int" };
export type Pick__Value$T$str = { $can_kind: "Pick__Some$T$str"; value: string } | { $can_kind: "Pick__None$T$str" };
export function pick__choose$T$int(present: boolean, value: bigint): { $can_kind: "ok"; value: Pick__Value$T$int } {
  if (present) {
    return { $can_kind: "ok", value: { $can_kind: "Pick__Some$T$int", value: value } };
  }
  else {
    return { $can_kind: "ok", value: { $can_kind: "Pick__None$T$int" } };
  }
}
export function pick__choose$T$str(present: boolean, value: string): { $can_kind: "ok"; value: Pick__Value$T$str } {
  if (present) {
    return { $can_kind: "ok", value: { $can_kind: "Pick__Some$T$str", value: value } };
  }
  else {
    return { $can_kind: "ok", value: { $can_kind: "Pick__None$T$str" } };
  }
}
export function pick__identity$T$int(value: Pick__Value$T$int): { $can_kind: "ok"; value: Pick__Value$T$int } {
  return { $can_kind: "ok", value: value };
}
export function pick__identity$T$str(value: Pick__Value$T$str): { $can_kind: "ok"; value: Pick__Value$T$str } {
  return { $can_kind: "ok", value: value };
}
export function pick__checked(value: bigint): { $can_kind: "ok"; value: Pick__Value$T$int } | { $can_kind: "pick.denied"; value: bigint } {
  const $can_m1 = (value < 0n);
  const $can_m2 = (value === 0n);
  if ($can_m1) {
    return { $can_kind: "pick.denied", value: value };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: { $can_kind: "Pick__None$T$int" } };
  }
  else {
    return { $can_kind: "ok", value: { $can_kind: "Pick__Some$T$int", value: value } };
  }
}
export function pick__relay(value: bigint): { $can_kind: "ok"; value: Pick__Value$T$int } | { $can_kind: "pick.denied"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: Pick__Value$T$int } | { $can_kind: "pick.denied"; value: bigint } = pick__checked(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "pick.denied": {
    const e = $can_m1;
    return { $can_kind: "pick.denied", value: e.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function pick__relay_generic$T$int(present: boolean, value: bigint): { $can_kind: "ok"; value: Pick__Value$T$int } {
  const $can_m1: { $can_kind: "ok"; value: Pick__Value$T$int } = pick__choose$T$int(present, value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function pick__relay_generic$T$str(present: boolean, value: string): { $can_kind: "ok"; value: Pick__Value$T$str } {
  const $can_m1: { $can_kind: "ok"; value: Pick__Value$T$str } = pick__choose$T$str(present, value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
