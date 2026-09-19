// GENERATED from use.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { success__box$T$Success__Pair, success__box$T$int, success__id$T$Success__Pair, success__id$T$int, success__reference$T$Success__Pair, success__reference$T$int, success__require$T$Success__Pair, type Success__Empty, type Success__Option$T$Success__Pair, type Success__Option$T$int, type Success__Pair, type Success__Value, type SuccessResult } from "./success";
import { use__host } from "./use.externs";
export type UseResult = { $can_kind: "ok" } | { $can_kind: "ok"; item: Success__Pair } | { $can_kind: "ok"; item: bigint } | { $can_kind: "ok"; left: bigint; right: boolean } | { $can_kind: "ok"; value: (input: Success__Pair) => { $can_kind: "ok"; left: bigint; right: boolean } } | { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } | { $can_kind: "ok"; value: Uint8Array } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: bigint[] } | { $can_kind: "option.absent" };
export function use__host_result(present: boolean): { $can_kind: "ok"; left: bigint; right: boolean } | { $can_kind: "option.absent" } {
  const $can_m1: { $can_kind: "ok"; left: bigint; right: boolean } | { $can_kind: "option.absent" } = use__host(present);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const pair = { left: $can_m1.left, right: $can_m1.right };
    return (($can_m2: { left: bigint; right: boolean }) => ({ $can_kind: "ok" as const, left: $can_m2.left, right: $can_m2.right }))(pair);
  }
  case "option.absent": {
    const e = $can_m1;
    return { $can_kind: "option.absent" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__apply$T$Bytes(cb: (input: Uint8Array) => { $can_kind: "ok"; value: Uint8Array }, value: Uint8Array): { $can_kind: "ok"; value: Uint8Array } {
  const $can_m1: { $can_kind: "ok"; value: Uint8Array } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const result = $can_m1.value;
    return { $can_kind: "ok", value: result };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__apply$T$Seq$L$int$G$(cb: (input: bigint[]) => { $can_kind: "ok"; value: bigint[] }, value: bigint[]): { $can_kind: "ok"; value: bigint[] } {
  const $can_m1: { $can_kind: "ok"; value: bigint[] } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const result = $can_m1.value;
    return { $can_kind: "ok", value: result };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__apply$T$Success__Empty(cb: (input: Success__Empty) => { $can_kind: "ok" }, value: Success__Empty): { $can_kind: "ok" } {
  const $can_m1: { $can_kind: "ok" } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const result = {  };
    return (($can_m2: {  }) => ({ $can_kind: "ok" as const }))(result);
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__apply$T$Success__Pair(cb: (input: Success__Pair) => { $can_kind: "ok"; left: bigint; right: boolean }, value: Success__Pair): { $can_kind: "ok"; left: bigint; right: boolean } {
  const $can_m1: { $can_kind: "ok"; left: bigint; right: boolean } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const result = { left: $can_m1.left, right: $can_m1.right };
    return (($can_m2: { left: bigint; right: boolean }) => ({ $can_kind: "ok" as const, left: $can_m2.left, right: $can_m2.right }))(result);
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__apply$T$Success__Value(cb: (input: Success__Value) => { $can_kind: "ok"; value: bigint }, value: Success__Value): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const result = { value: $can_m1.value };
    return (($can_m2: { value: bigint }) => ({ $can_kind: "ok" as const, value: $can_m2.value }))(result);
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__apply$T$int(cb: (input: bigint) => { $can_kind: "ok"; value: bigint }, value: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const result = $can_m1.value;
    return { $can_kind: "ok", value: result };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__relay$T$Success__Pair(value: Success__Pair): { $can_kind: "ok"; left: bigint; right: boolean } {
  const $can_m1: { $can_kind: "ok"; left: bigint; right: boolean } = use__apply$T$Success__Pair((() => { return (input) => success__id$T$Success__Pair(input); })(), value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const result = { left: $can_m1.left, right: $can_m1.right };
    return { $can_kind: "ok", left: result.left, right: result.right };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__relay$T$int(value: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = use__apply$T$int((() => { return (input) => success__id$T$int(input); })(), value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const result = $can_m1.value;
    return { $can_kind: "ok", value: result };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__box$T$Success__Pair(value: Success__Pair): { $can_kind: "ok"; item: Success__Pair } {
  const $can_m1: { $can_kind: "ok"; item: Success__Pair } = success__box$T$Success__Pair(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const box = { item: $can_m1.item };
    return (($can_m2: { item: Success__Pair }) => ({ $can_kind: "ok" as const, item: $can_m2.item }))(box);
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__box$T$int(value: bigint): { $can_kind: "ok"; item: bigint } {
  const $can_m1: { $can_kind: "ok"; item: bigint } = success__box$T$int(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const box = { item: $can_m1.item };
    return (($can_m2: { item: bigint }) => ({ $can_kind: "ok" as const, item: $can_m2.item }))(box);
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__factory$T$Success__Pair(): { $can_kind: "ok"; value: (input: Success__Pair) => { $can_kind: "ok"; left: bigint; right: boolean } } {
  const $can_m1: { $can_kind: "ok"; value: (input: Success__Pair) => { $can_kind: "ok"; left: bigint; right: boolean } } = success__reference$T$Success__Pair();
  switch ($can_m1.$can_kind) {
  case "ok": {
    const cb = $can_m1.value;
    return { $can_kind: "ok", value: cb };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__factory$T$int(): { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } {
  const $can_m1: { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } = success__reference$T$int();
  switch ($can_m1.$can_kind) {
  case "ok": {
    const cb = $can_m1.value;
    return { $can_kind: "ok", value: cb };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__factory_apply$T$Success__Pair(value: Success__Pair): { $can_kind: "ok"; left: bigint; right: boolean } {
  const $can_m1: { $can_kind: "ok"; value: (input: Success__Pair) => { $can_kind: "ok"; left: bigint; right: boolean } } = use__factory$T$Success__Pair();
  switch ($can_m1.$can_kind) {
  case "ok": {
    const cb = $can_m1.value;
    const $can_m2: { $can_kind: "ok"; left: bigint; right: boolean } = use__apply$T$Success__Pair(cb, value);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const result = { left: $can_m2.left, right: $can_m2.right };
      return (($can_m3: { left: bigint; right: boolean }) => ({ $can_kind: "ok" as const, left: $can_m3.left, right: $can_m3.right }))(result);
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__factory_apply$T$int(value: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } = use__factory$T$int();
  switch ($can_m1.$can_kind) {
  case "ok": {
    const cb = $can_m1.value;
    const $can_m2: { $can_kind: "ok"; value: bigint } = use__apply$T$int(cb, value);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const result = $can_m2.value;
      return { $can_kind: "ok", value: result };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__require$T$Success__Pair(cb: (input: Success__Option$T$Success__Pair) => { $can_kind: "ok"; left: bigint; right: boolean } | { $can_kind: "option.absent" }, value: Success__Option$T$Success__Pair): { $can_kind: "ok"; left: bigint; right: boolean } | { $can_kind: "option.absent" } {
  const $can_m1: { $can_kind: "ok"; left: bigint; right: boolean } | { $can_kind: "option.absent" } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const result = { left: $can_m1.left, right: $can_m1.right };
    return { $can_kind: "ok", left: result.left, right: result.right };
  }
  case "option.absent": {
    const e = $can_m1;
    return { $can_kind: "option.absent" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__require$T$int(cb: (input: Success__Option$T$int) => { $can_kind: "ok"; value: bigint } | { $can_kind: "option.absent" }, value: Success__Option$T$int): { $can_kind: "ok"; value: bigint } | { $can_kind: "option.absent" } {
  const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "option.absent" } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const result = $can_m1.value;
    return { $can_kind: "ok", value: result };
  }
  case "option.absent": {
    const e = $can_m1;
    return { $can_kind: "option.absent" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__recover(value: Success__Option$T$Success__Pair): { $can_kind: "ok"; left: bigint; right: boolean } {
  const $can_m1: { $can_kind: "ok"; left: bigint; right: boolean } | { $can_kind: "option.absent" } = use__require$T$Success__Pair((() => { return (input) => success__require$T$Success__Pair(input); })(), value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const result = { left: $can_m1.left, right: $can_m1.right };
    return (($can_m2: { left: bigint; right: boolean }) => ({ $can_kind: "ok" as const, left: $can_m2.left, right: $can_m2.right }))(result);
  }
  case "option.absent": {
    const e = $can_m1;
    return (($can_m3: { left: bigint; right: boolean }) => ({ $can_kind: "ok" as const, left: $can_m3.left, right: $can_m3.right }))({ left: 0n, right: false });
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
