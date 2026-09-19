// GENERATED from use.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { scalar__add, scalar__checked, scalar__positive, type ScalarResult } from "./scalar";
export type UseResult = { $can_kind: "ok"; cb: (input: bigint) => { $can_kind: "ok"; value: bigint } } | { $can_kind: "ok"; positive: boolean } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: boolean } | { $can_kind: "ok"; value: string } | { $can_kind: "scalar.negative"; value: bigint };
export type Use__Bound = { cb: (input: bigint) => { $can_kind: "ok"; value: bigint } };
export type Use__Report = { positive: boolean };
export function use__apply$T$bool$T$bool(cb: (input: boolean) => { $can_kind: "ok"; value: boolean }, value: boolean): { $can_kind: "ok"; value: boolean } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = cb(value);
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
export function use__apply$T$dec$T$dec(cb: (input: string) => { $can_kind: "ok"; value: string }, value: string): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } = cb(value);
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
export function use__apply$T$int$T$bool(cb: (input: bigint) => { $can_kind: "ok"; value: boolean }, value: bigint): { $can_kind: "ok"; value: boolean } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = cb(value);
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
export function use__apply$T$int$T$int(cb: (input: bigint) => { $can_kind: "ok"; value: bigint }, value: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = cb(value);
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
export function use__apply$T$str$T$str(cb: (input: string) => { $can_kind: "ok"; value: string }, value: string): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } = cb(value);
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
export function use__bind(delta: bigint): { $can_kind: "ok"; cb: (input: bigint) => { $can_kind: "ok"; value: bigint } } {
  return { $can_kind: "ok", cb: (() => { const $can_m1 = delta; return (input) => scalar__add($can_m1, input); })() };
}
export function use__go(delta: bigint, value: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; cb: (input: bigint) => { $can_kind: "ok"; value: bigint } } = use__bind(delta);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const bound = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } = use__apply$T$int$T$int(bound.cb, value);
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
export function use__relay(value: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "scalar.negative"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "scalar.negative"; value: bigint } = scalar__checked(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "scalar.negative": {
    const e = $can_m1;
    return { $can_kind: "scalar.negative", value: e.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__checked(cb: (input: bigint) => { $can_kind: "ok"; value: bigint } | { $can_kind: "scalar.negative"; value: bigint }, value: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "scalar.negative"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "scalar.negative"; value: bigint } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "scalar.negative": {
    const e = $can_m1;
    return { $can_kind: "scalar.negative", value: e.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__recover(cb: (input: bigint) => { $can_kind: "ok"; value: bigint } | { $can_kind: "scalar.negative"; value: bigint }, value: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "scalar.negative"; value: bigint } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "scalar.negative": {
    const e = $can_m1;
    return { $can_kind: "ok", value: (e.value + 3n) };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__report(value: bigint): { $can_kind: "ok"; positive: boolean } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = scalar__positive(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", positive: r.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__classify(cb: (input: bigint) => { $can_kind: "ok"; value: boolean }, value: bigint): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = cb(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    if (r.value) {
      return { $can_kind: "ok", value: "positive" };
    }
    else {
      return { $can_kind: "ok", value: "not positive" };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__sugar(value: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "scalar.negative"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "scalar.negative"; value: bigint } = use__relay(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "scalar.negative": {
    const e = $can_m1;
    return { $can_kind: "scalar.negative", value: e.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
