// GENERATED from clock.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { std__clock__wall_now, type HostResult } from "./host";
export type ClockdemoResult = { $can_kind: "ok"; millis: bigint };
export type Clockdemo__Seen = { millis: bigint };
export function clockdemo__read_millis(): { $can_kind: "ok"; millis: bigint } | { $can_kind: "clock.unavailable" } {
  const $can_m1: { $can_kind: "ok"; millis: bigint } | { $can_kind: "clock.unavailable" } = std__clock__wall_now();
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
