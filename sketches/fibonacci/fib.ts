// GENERATED from fib.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type FibResult = { $can_kind: "ok"; value: bigint } | { $can_kind: "fib.negative_input"; value: bigint };
export type Fib__Value = { value: bigint };
export function fib__at(n: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "fib.negative_input"; value: bigint } {
  if ((n < 0n)) {
    return { $can_kind: "fib.negative_input", value: n };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: bigint } = fib__from(n, 0n, 1n);
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
}
export function fib__from(n: bigint, a: bigint, b: bigint): { $can_kind: "ok"; value: bigint } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: a };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: bigint } = fib__from((n - 1n), b, (a + b));
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
}
