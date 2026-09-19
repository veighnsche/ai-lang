// GENERATED from use.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { pick__checked, type Pick__Value$T$int, type PickResult } from "./pick";
export type UseResult = { $can_kind: "ok"; value: Pick__Value$T$int } | { $can_kind: "ok"; value: string };
export type Use__Text = { value: string };
export function use__label(value: bigint): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; value: Pick__Value$T$int } | { $can_kind: "pick.denied"; value: bigint } = pick__checked(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    const $can_m2 = r.value;
    switch ($can_m2.$can_kind) {
    case "Pick__Some$T$int": {
      const s = $can_m2;
      return { $can_kind: "ok", value: "present" };
    }
    case "Pick__None$T$int": {
      const _ = $can_m2;
      return { $can_kind: "ok", value: "absent" };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "pick.denied": {
    const e = $can_m1;
    return { $can_kind: "ok", value: "denied" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__relay(value: bigint): { $can_kind: "ok"; value: Pick__Value$T$int } | { $can_kind: "pick.denied"; value: bigint } {
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
export function use__recover(value: bigint): { $can_kind: "ok"; value: Pick__Value$T$int } {
  const $can_m1: { $can_kind: "ok"; value: Pick__Value$T$int } | { $can_kind: "pick.denied"; value: bigint } = pick__checked(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "pick.denied": {
    const e = $can_m1;
    return { $can_kind: "ok", value: { $can_kind: "Pick__None$T$int" } };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
