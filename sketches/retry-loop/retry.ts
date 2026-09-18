// GENERATED from retry.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { net__fetch } from "./retry.externs";
export type RetryResult = { $can_kind: "ok"; body: string } | { $can_kind: "retry.exhausted" } | { $can_kind: "net.down" };
export type Retry__Doc = { body: string };
export function retry__fetch(fuel: bigint): { $can_kind: "ok"; body: string } | { $can_kind: "retry.exhausted" } {
  if ((fuel <= 0n)) {
    return { $can_kind: "retry.exhausted" };
  }
  else {
    const $can_m1: { $can_kind: "ok"; body: string } | { $can_kind: "net.down" } = net__fetch();
    switch ($can_m1.$can_kind) {
    case "net.down": {
      const _ = $can_m1;
      const $can_m2: { $can_kind: "ok"; body: string } | { $can_kind: "retry.exhausted" } = retry__fetch((fuel - 1n));
      switch ($can_m2.$can_kind) {
      case "retry.exhausted": {
        const _ = $can_m2;
        return { $can_kind: "retry.exhausted" };
      }
      case "ok": {
        const d = $can_m2;
        return { $can_kind: "ok", body: d.body };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    case "ok": {
      const d = $can_m1;
      return { $can_kind: "ok", body: d.body };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
