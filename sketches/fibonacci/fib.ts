// GENERATED from fib.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type FibResult = { $ail_kind: "ok"; value: bigint } | { $ail_kind: "fib.negative_input"; value: bigint };
export type Fib__Value = { value: bigint };
export function fib__at(n: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "fib.negative_input"; value: bigint } {
  if ((n < 0n)) {
    return { $ail_kind: "fib.negative_input", value: n };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = fib__from(n, 0n, 1n);
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const r = $ail_m1;
      return { $ail_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function fib__from(n: bigint, a: bigint, b: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: a };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = fib__from((n - 1n), b, (a + b));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const r = $ail_m1;
      return { $ail_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
