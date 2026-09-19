// GENERATED from use.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { data__factory, type Data__Choice$T$int, type Data__Row, type DataResult } from "./data";
import { bridge__bool, bridge__bytes, bridge__choice, bridge__dec, bridge__int, bridge__rows, bridge__str, bridge__token } from "./bridge.externs";
export type UseResult = { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } | { $can_kind: "ok"; value: Data__Choice$T$int } | { $can_kind: "ok"; value: Data__Row[] } | { $can_kind: "ok"; value: Uint8Array } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: bigint[] } | { $can_kind: "ok"; value: boolean } | { $can_kind: "ok"; value: string } | { $can_kind: "bridge.denied" };
export function use__int(allow: boolean): { $can_kind: "ok"; value: bigint } | { $can_kind: "bridge.denied" } {
  const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "bridge.denied" } = bridge__int(allow);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "bridge.denied": {
    const e = $can_m1;
    return { $can_kind: "bridge.denied" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__str(): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } = bridge__str();
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
export function use__bool(): { $can_kind: "ok"; value: boolean } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = bridge__bool();
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
export function use__dec(): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } = bridge__dec();
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
export function use__token(): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } = bridge__token();
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
export function use__bytes(): { $can_kind: "ok"; value: Uint8Array } {
  const $can_m1: { $can_kind: "ok"; value: Uint8Array } = bridge__bytes();
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
export function use__rows(): { $can_kind: "ok"; value: Data__Row[] } {
  const $can_m1: { $can_kind: "ok"; value: Data__Row[] } = bridge__rows();
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
export function use__choice(some: boolean): { $can_kind: "ok"; value: Data__Choice$T$int } {
  const $can_m1: { $can_kind: "ok"; value: Data__Choice$T$int } = bridge__choice(some);
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
export function use__apply_bytes(cb: (input: Uint8Array) => { $can_kind: "ok"; value: Uint8Array }, value: Uint8Array): { $can_kind: "ok"; value: Uint8Array } {
  const $can_m1: { $can_kind: "ok"; value: Uint8Array } = cb(value);
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
export function use__apply_token(cb: (input: string) => { $can_kind: "ok"; value: string }, value: string): { $can_kind: "ok"; value: string } {
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
export function use__apply_items$T$Data__Row(cb: (input: Data__Row[]) => { $can_kind: "ok"; value: Data__Row[] }, values: Data__Row[]): { $can_kind: "ok"; value: Data__Row[] } {
  const $can_m1: { $can_kind: "ok"; value: Data__Row[] } = cb(values);
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
export function use__apply_items$T$int(cb: (input: bigint[]) => { $can_kind: "ok"; value: bigint[] }, values: bigint[]): { $can_kind: "ok"; value: bigint[] } {
  const $can_m1: { $can_kind: "ok"; value: bigint[] } = cb(values);
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
export function use__apply_int(cb: (input: bigint) => { $can_kind: "ok"; value: bigint }, value: bigint): { $can_kind: "ok"; value: bigint } {
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
export function use__factory_relay(delta: bigint): { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } {
  const $can_m1: { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } = data__factory(delta);
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
export function use__factory(delta: bigint, value: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: (input: bigint) => { $can_kind: "ok"; value: bigint } } = use__factory_relay(delta);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const f = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } = use__apply_int(f.value, value);
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
