// GENERATED from success.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type SuccessResult = { $can_kind: "ok" } | { $can_kind: "ok"; item: Success__Pair } | { $can_kind: "ok"; item: bigint } | { $can_kind: "ok"; left: bigint; right: boolean } | { $can_kind: "ok"; value: (input: Success__Pair) => { $can_kind: "ok"; left: bigint; right: boolean } } | { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } | { $can_kind: "ok"; value: Success__Option$T$Success__Pair } | { $can_kind: "ok"; value: Success__Option$T$int } | { $can_kind: "ok"; value: Uint8Array } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: bigint[] } | { $can_kind: "ok"; value: boolean } | { $can_kind: "ok"; value: string } | { $can_kind: "option.absent" };
export type Success__Pair = { left: bigint; right: boolean };
export type Success__Value = { value: bigint };
export type Success__Empty = {};
export type Success__Box$T$Success__Pair = { item: Success__Pair };
export type Success__Box$T$int = { item: bigint };
export type Success__Option$T$Success__Pair = { $can_kind: "Success__Some$T$Success__Pair"; value: Success__Pair } | { $can_kind: "Success__None$T$Success__Pair" };
export type Success__Option$T$int = { $can_kind: "Success__Some$T$int"; value: bigint } | { $can_kind: "Success__None$T$int" };
export function success__id$T$Bytes(value: Uint8Array): { $can_kind: "ok"; value: Uint8Array } {
  return { $can_kind: "ok", value: value };
}
export function success__id$T$Seq$L$int$G$(value: bigint[]): { $can_kind: "ok"; value: bigint[] } {
  return { $can_kind: "ok", value: value };
}
export function success__id$T$Success__Empty(value: Success__Empty): { $can_kind: "ok" } {
  return (($can_m1: {  }) => ({ $can_kind: "ok" as const }))(value);
}
export function success__id$T$Success__Pair(value: Success__Pair): { $can_kind: "ok"; left: bigint; right: boolean } {
  return (($can_m1: { left: bigint; right: boolean }) => ({ $can_kind: "ok" as const, left: $can_m1.left, right: $can_m1.right }))(value);
}
export function success__id$T$Success__Token(value: string): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: value };
}
export function success__id$T$Success__Value(value: Success__Value): { $can_kind: "ok"; value: bigint } {
  return (($can_m1: { value: bigint }) => ({ $can_kind: "ok" as const, value: $can_m1.value }))(value);
}
export function success__id$T$bool(value: boolean): { $can_kind: "ok"; value: boolean } {
  return { $can_kind: "ok", value: value };
}
export function success__id$T$dec(value: string): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: value };
}
export function success__id$T$int(value: bigint): { $can_kind: "ok"; value: bigint } {
  return { $can_kind: "ok", value: value };
}
export function success__id$T$str(value: string): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: value };
}
export function success__box$T$Success__Pair(value: Success__Pair): { $can_kind: "ok"; item: Success__Pair } {
  return (($can_m1: { item: Success__Pair }) => ({ $can_kind: "ok" as const, item: $can_m1.item }))({ item: value });
}
export function success__box$T$int(value: bigint): { $can_kind: "ok"; item: bigint } {
  return (($can_m1: { item: bigint }) => ({ $can_kind: "ok" as const, item: $can_m1.item }))({ item: value });
}
export function success__reference$T$Success__Pair(): { $can_kind: "ok"; value: (input: Success__Pair) => { $can_kind: "ok"; left: bigint; right: boolean } } {
  return { $can_kind: "ok", value: (() => { return (input) => success__id$T$Success__Pair(input); })() };
}
export function success__reference$T$int(): { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } {
  return { $can_kind: "ok", value: (() => { return (input) => success__id$T$int(input); })() };
}
export function success__value_or$T$Success__Pair(value: Success__Option$T$Success__Pair, alternative: Success__Pair): { $can_kind: "ok"; left: bigint; right: boolean } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Success__Some$T$Success__Pair": {
    const s = $can_m1;
    return (($can_m2: { left: bigint; right: boolean }) => ({ $can_kind: "ok" as const, left: $can_m2.left, right: $can_m2.right }))(s.value);
  }
  case "Success__None$T$Success__Pair": {
    const _ = $can_m1;
    return (($can_m3: { left: bigint; right: boolean }) => ({ $can_kind: "ok" as const, left: $can_m3.left, right: $can_m3.right }))(alternative);
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function success__value_or$T$int(value: Success__Option$T$int, alternative: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Success__Some$T$int": {
    const s = $can_m1;
    return { $can_kind: "ok", value: s.value };
  }
  case "Success__None$T$int": {
    const _ = $can_m1;
    return { $can_kind: "ok", value: alternative };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function success__require$T$Success__Pair(value: Success__Option$T$Success__Pair): { $can_kind: "ok"; left: bigint; right: boolean } | { $can_kind: "option.absent" } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Success__Some$T$Success__Pair": {
    const s = $can_m1;
    return (($can_m2: { left: bigint; right: boolean }) => ({ $can_kind: "ok" as const, left: $can_m2.left, right: $can_m2.right }))(s.value);
  }
  case "Success__None$T$Success__Pair": {
    const _ = $can_m1;
    return { $can_kind: "option.absent" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function success__require$T$int(value: Success__Option$T$int): { $can_kind: "ok"; value: bigint } | { $can_kind: "option.absent" } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Success__Some$T$int": {
    const s = $can_m1;
    return { $can_kind: "ok", value: s.value };
  }
  case "Success__None$T$int": {
    const _ = $can_m1;
    return { $can_kind: "option.absent" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function success__pair(value: bigint): { $can_kind: "ok"; left: bigint; right: boolean } {
  return (($can_m1: { left: bigint; right: boolean }) => ({ $can_kind: "ok" as const, left: $can_m1.left, right: $can_m1.right }))({ left: value, right: (value > 0n) });
}
export function success__mapped(value: Success__Option$T$int): { $can_kind: "ok"; value: Success__Option$T$Success__Pair } {
  const $can_m1: { $can_kind: "ok"; value: Success__Option$T$Success__Pair } = success__map$T$int$T$Success__Pair(value, (() => { return (input) => success__pair(input); })());
  switch ($can_m1.$can_kind) {
  case "ok": {
    const mapped = $can_m1.value;
    return { $can_kind: "ok", value: mapped };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function success__map$T$Success__Pair$T$Success__Pair(value: Success__Option$T$Success__Pair, transform: (input: Success__Pair) => { $can_kind: "ok"; left: bigint; right: boolean }): { $can_kind: "ok"; value: Success__Option$T$Success__Pair } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Success__Some$T$Success__Pair": {
    const s = $can_m1;
    const $can_m2: { $can_kind: "ok"; left: bigint; right: boolean } = transform(s.value);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const mapped = { left: $can_m2.left, right: $can_m2.right };
      return { $can_kind: "ok", value: { $can_kind: "Success__Some$T$Success__Pair", value: mapped } };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "Success__None$T$Success__Pair": {
    const _ = $can_m1;
    return { $can_kind: "ok", value: { $can_kind: "Success__None$T$Success__Pair" } };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function success__map$T$int$T$Success__Pair(value: Success__Option$T$int, transform: (input: bigint) => { $can_kind: "ok"; left: bigint; right: boolean }): { $can_kind: "ok"; value: Success__Option$T$Success__Pair } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Success__Some$T$int": {
    const s = $can_m1;
    const $can_m2: { $can_kind: "ok"; left: bigint; right: boolean } = transform(s.value);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const mapped = { left: $can_m2.left, right: $can_m2.right };
      return { $can_kind: "ok", value: { $can_kind: "Success__Some$T$Success__Pair", value: mapped } };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "Success__None$T$int": {
    const _ = $can_m1;
    return { $can_kind: "ok", value: { $can_kind: "Success__None$T$Success__Pair" } };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function success__map$T$int$T$int(value: Success__Option$T$int, transform: (input: bigint) => { $can_kind: "ok"; value: bigint }): { $can_kind: "ok"; value: Success__Option$T$int } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Success__Some$T$int": {
    const s = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } = transform(s.value);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const mapped = $can_m2.value;
      return { $can_kind: "ok", value: { $can_kind: "Success__Some$T$int", value: mapped } };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "Success__None$T$int": {
    const _ = $can_m1;
    return { $can_kind: "ok", value: { $can_kind: "Success__None$T$int" } };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
