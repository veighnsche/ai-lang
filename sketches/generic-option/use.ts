// GENERATED from use.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { option__value_or$T$int, type Option__Value$T$int, type Option__Value$T$str, type OptionResult } from "./option";
export type UseResult = { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: string };
export type Use__Number = { value: bigint };
export type Use__Text = { value: string };
export function use__number(value: Option__Value$T$int): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = option__value_or$T$int(value, 9n);
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
export function use__text(value: Option__Value$T$str): { $can_kind: "ok"; value: string } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Option__Some$T$str": {
    const s = $can_m1;
    return { $can_kind: "ok", value: s.value };
  }
  case "Option__None$T$str": {
    const _ = $can_m1;
    return { $can_kind: "ok", value: "fallback" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
