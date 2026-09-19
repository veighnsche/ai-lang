// GENERATED from host.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { host__mono_now, host__wall_now } from "./host.externs";
export type HostResult = { $can_kind: "ok"; millis: bigint } | { $can_kind: "clock.unavailable" };
export type Clock__Instant = { millis: bigint };
export type Clock__MonotonicTick = { millis: bigint };
export function std__clock__wall_now(): { $can_kind: "ok"; millis: bigint } | { $can_kind: "clock.unavailable" } {
  const $can_m1: { $can_kind: "ok"; millis: bigint } | { $can_kind: "clock.unavailable" } = host__wall_now();
  switch ($can_m1.$can_kind) {
  case "clock.unavailable": {
    const _ = $can_m1;
    return { $can_kind: "clock.unavailable" };
  }
  case "ok": {
    const t = $can_m1;
    return { $can_kind: "ok", millis: t.millis };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__clock__monotonic_now(): { $can_kind: "ok"; millis: bigint } | { $can_kind: "clock.unavailable" } {
  const $can_m1: { $can_kind: "ok"; millis: bigint } | { $can_kind: "clock.unavailable" } = host__mono_now();
  switch ($can_m1.$can_kind) {
  case "clock.unavailable": {
    const _ = $can_m1;
    return { $can_kind: "clock.unavailable" };
  }
  case "ok": {
    const t = $can_m1;
    return { $can_kind: "ok", millis: t.millis };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
