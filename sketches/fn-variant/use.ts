// GENERATED from use.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { choice__choose$T$int, choice__choose$T$str, type Choice__Value$T$int, type Choice__Value$T$str, type ChoiceResult } from "./choice";
export type UseResult = { $can_kind: "ok"; cb: (input: bigint) => { $can_kind: "ok"; value: Choice__Value$T$int } } | { $can_kind: "ok"; cb: (input: string) => { $can_kind: "ok"; value: Choice__Value$T$str } } | { $can_kind: "ok"; value: Choice__Value$T$int } | { $can_kind: "ok"; value: Choice__Value$T$str } | { $can_kind: "ok"; value: bigint };
export type Use__Bound$T$int = { cb: (input: bigint) => { $can_kind: "ok"; value: Choice__Value$T$int } };
export type Use__Bound$T$str = { cb: (input: string) => { $can_kind: "ok"; value: Choice__Value$T$str } };
export type Use__Number = { value: bigint };
export function use__bind$T$int(present: boolean): { $can_kind: "ok"; cb: (input: bigint) => { $can_kind: "ok"; value: Choice__Value$T$int } } {
  return { $can_kind: "ok", cb: (() => { const $can_m1 = present; return (input) => choice__choose$T$int($can_m1, input); })() };
}
export function use__bind$T$str(present: boolean): { $can_kind: "ok"; cb: (input: string) => { $can_kind: "ok"; value: Choice__Value$T$str } } {
  return { $can_kind: "ok", cb: (() => { const $can_m1 = present; return (input) => choice__choose$T$str($can_m1, input); })() };
}
export function use__apply$T$int(cb: (input: bigint) => { $can_kind: "ok"; value: Choice__Value$T$int }, value: bigint): { $can_kind: "ok"; value: Choice__Value$T$int } {
  const $can_m1: { $can_kind: "ok"; value: Choice__Value$T$int } = cb(value);
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
export function use__apply$T$str(cb: (input: string) => { $can_kind: "ok"; value: Choice__Value$T$str }, value: string): { $can_kind: "ok"; value: Choice__Value$T$str } {
  const $can_m1: { $can_kind: "ok"; value: Choice__Value$T$str } = cb(value);
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
export function use__number(cb: (input: bigint) => { $can_kind: "ok"; value: Choice__Value$T$int } | { $can_kind: "choice.denied"; value: bigint }, value: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: Choice__Value$T$int } | { $can_kind: "choice.denied"; value: bigint } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    const $can_m2 = r.value;
    switch ($can_m2.$can_kind) {
    case "Choice__Some$T$int": {
      const s = $can_m2;
      return { $can_kind: "ok", value: s.value };
    }
    case "Choice__None$T$int": {
      const _ = $can_m2;
      return { $can_kind: "ok", value: 0n };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "choice.denied": {
    const e = $can_m1;
    return { $can_kind: "ok", value: -1n };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__checked(cb: (input: bigint) => { $can_kind: "ok"; value: Choice__Value$T$int } | { $can_kind: "choice.denied"; value: bigint }, value: bigint): { $can_kind: "ok"; value: Choice__Value$T$int } | { $can_kind: "choice.denied"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: Choice__Value$T$int } | { $can_kind: "choice.denied"; value: bigint } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "choice.denied": {
    const e = $can_m1;
    return { $can_kind: "choice.denied", value: e.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__recover(cb: (input: bigint) => { $can_kind: "ok"; value: Choice__Value$T$int } | { $can_kind: "choice.denied"; value: bigint }, value: bigint): { $can_kind: "ok"; value: Choice__Value$T$int } {
  const $can_m1: { $can_kind: "ok"; value: Choice__Value$T$int } | { $can_kind: "choice.denied"; value: bigint } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "choice.denied": {
    const e = $can_m1;
    return { $can_kind: "ok", value: { $can_kind: "Choice__None$T$int" } };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__go(present: boolean, value: bigint): { $can_kind: "ok"; value: Choice__Value$T$int } {
  const $can_m1: { $can_kind: "ok"; cb: (input: bigint) => { $can_kind: "ok"; value: Choice__Value$T$int } } = use__bind$T$int(present);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const bound = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: Choice__Value$T$int } = use__apply$T$int(bound.cb, value);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const r = $can_m2;
      return { $can_kind: "ok", value: r.value };
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
