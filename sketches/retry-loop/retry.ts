// GENERATED from retry.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { net__fetch } from "./retry.externs";
export type RetryResult = { kind: "ok"; body: string } | { kind: "retry.exhausted" } | { kind: "net.down" };
export type Retry__Doc = { body: string };
export function retry__fetch(fuel: bigint): RetryResult {
  if ((fuel <= 0n)) {
    return { kind: "retry.exhausted" };
  }
  else {
    const $ail_m1: { kind: "ok"; body: string } | { kind: "net.down" } = net__fetch();
    switch ($ail_m1.kind) {
    case "net.down": {
      const _ = $ail_m1;
      const $ail_m2: RetryResult = retry__fetch((fuel - 1n));
      switch ($ail_m2.kind) {
      case "retry.exhausted": {
        const _ = $ail_m2;
        return { kind: "retry.exhausted" };
      }
      case "ok": {
        const d = $ail_m2;
        return { kind: "ok", body: d.body };
      }
      }
    }
    case "ok": {
      const d = $ail_m1;
      return { kind: "ok", body: d.body };
    }
    }
  }
}
