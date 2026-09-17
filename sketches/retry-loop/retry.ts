// GENERATED from retry.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { net__fetch } from "./retry.externs";
export type RetryResult = { $ail_kind: "ok"; body: string } | { $ail_kind: "retry.exhausted" } | { $ail_kind: "net.down" };
export type Retry__Doc = { body: string };
export function retry__fetch(fuel: bigint): RetryResult {
  if ((fuel <= 0n)) {
    return { $ail_kind: "retry.exhausted" };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; body: string } | { $ail_kind: "net.down" } = net__fetch();
    switch ($ail_m1.$ail_kind) {
    case "net.down": {
      const _ = $ail_m1;
      const $ail_m2: RetryResult = retry__fetch((fuel - 1n));
      switch ($ail_m2.$ail_kind) {
      case "retry.exhausted": {
        const _ = $ail_m2;
        return { $ail_kind: "retry.exhausted" };
      }
      case "ok": {
        const d = $ail_m2;
        return { $ail_kind: "ok", body: d.body };
      }
      }
    }
    case "ok": {
      const d = $ail_m1;
      return { $ail_kind: "ok", body: d.body };
    }
    }
  }
}
