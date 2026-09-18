// GENERATED from division.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type DivisionResult = { $ail_kind: "ok"; quotient: bigint; remainder: bigint } | { $ail_kind: "ok"; value: bigint } | { $ail_kind: "ok"; value: boolean } | { $ail_kind: "math.zero_divisor"; divisor: bigint };
export type Int__DivMod = { quotient: bigint; remainder: bigint };
export type Int__Value = { value: bigint };
export type Bool__Value = { value: boolean };
// Euclidean integer division (a17): quotient and remainder with
// 0 <= r < |b| on every sign combination.
function $ailDivMod(a: bigint, b: bigint): [bigint, bigint] {
  let q: bigint = a / b;
  let r: bigint = a - b * q;
  if (r < 0n) {
    q += b > 0n ? -1n : 1n;
    r = a - b * q;
  }
  return [q, r];
}
export function std__int__divmod(dividend: bigint, divisor: bigint): { $ail_kind: "ok"; quotient: bigint; remainder: bigint } | { $ail_kind: "math.zero_divisor"; divisor: bigint } {
  if ((divisor === 0n)) {
    return { $ail_kind: "math.zero_divisor", divisor: divisor };
  }
  else {
    return { $ail_kind: "ok", quotient: $ailDivMod(dividend, divisor)[0], remainder: $ailDivMod(dividend, divisor)[1] };
  }
}
export function std__int__mod(value: bigint, modulus: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "math.zero_divisor"; divisor: bigint } {
  if ((modulus === 0n)) {
    return { $ail_kind: "math.zero_divisor", divisor: modulus };
  }
  else {
    return { $ail_kind: "ok", value: $ailDivMod(value, modulus)[1] };
  }
}
export function std__int__is_multiple(value: bigint, divisor: bigint): { $ail_kind: "ok"; value: boolean } | { $ail_kind: "math.zero_divisor"; divisor: bigint } {
  if ((divisor === 0n)) {
    return { $ail_kind: "math.zero_divisor", divisor: divisor };
  }
  else {
    return { $ail_kind: "ok", value: ($ailDivMod(value, divisor)[1] === 0n) };
  }
}
export function std__int__is_even(value: bigint): { $ail_kind: "ok"; value: boolean } {
  return { $ail_kind: "ok", value: ($ailDivMod(value, 2n)[1] === 0n) };
}
export function std__int__is_odd(value: bigint): { $ail_kind: "ok"; value: boolean } {
  return { $ail_kind: "ok", value: ($ailDivMod(value, 2n)[1] === 1n) };
}
