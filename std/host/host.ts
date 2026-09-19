// GENERATED from host.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { host__hash_digest, host__mono_now, host__rand_bytes, host__wall_now } from "./host.externs";
export type HostResult = { $can_kind: "ok"; bytes: Uint8Array } | { $can_kind: "ok"; millis: bigint } | { $can_kind: "clock.unavailable" } | { $can_kind: "random.invalid_count" } | { $can_kind: "random.unavailable" } | { $can_kind: "hash.unsupported_profile" } | { $can_kind: "hash.unavailable" };
export type Clock__Instant = { millis: bigint };
export type Clock__MonotonicTick = { millis: bigint };
export type Random__Bytes = { bytes: Uint8Array };
export type Hash__Digest = { bytes: Uint8Array };
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
export function std__random__bytes(count: bigint): { $can_kind: "ok"; bytes: Uint8Array } | { $can_kind: "random.invalid_count" } | { $can_kind: "random.unavailable" } {
  if ((count < 0n)) {
    return { $can_kind: "random.invalid_count" };
  }
  else {
    const $can_m1: { $can_kind: "ok"; bytes: Uint8Array } | { $can_kind: "random.invalid_count" } | { $can_kind: "random.unavailable" } = host__rand_bytes(count);
    switch ($can_m1.$can_kind) {
    case "random.invalid_count": {
      const _ = $can_m1;
      return { $can_kind: "random.invalid_count" };
    }
    case "random.unavailable": {
      const _ = $can_m1;
      return { $can_kind: "random.unavailable" };
    }
    case "ok": {
      const r = $can_m1;
      return { $can_kind: "ok", bytes: r.bytes };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__hash__digest(data: Uint8Array, profile: string): { $can_kind: "ok"; bytes: Uint8Array } | { $can_kind: "hash.unsupported_profile" } | { $can_kind: "hash.unavailable" } {
  const $can_m1: { $can_kind: "ok"; bytes: Uint8Array } | { $can_kind: "hash.unsupported_profile" } | { $can_kind: "hash.unavailable" } = host__hash_digest(data, profile);
  switch ($can_m1.$can_kind) {
  case "hash.unsupported_profile": {
    const _ = $can_m1;
    return { $can_kind: "hash.unsupported_profile" };
  }
  case "hash.unavailable": {
    const _ = $can_m1;
    return { $can_kind: "hash.unavailable" };
  }
  case "ok": {
    const d = $can_m1;
    return { $can_kind: "ok", bytes: d.bytes };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
