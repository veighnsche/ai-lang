// GENERATED from bridge.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { type Data__Choice$T$int, type Data__Row, type DataResult } from "./data";
import { bridge__bytes } from "./bridge.externs";
export type BridgeResult = { $can_kind: "ok"; value: Uint8Array } | { $can_kind: "bridge.denied" };
export function bridge__local(): { $can_kind: "ok"; value: Uint8Array } {
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
