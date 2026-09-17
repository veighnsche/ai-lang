// GENERATED from division.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type DivisionResult = { kind: "ok"; quotient: bigint; remainder: bigint } | { kind: "ok"; value: bigint } | { kind: "math.zero_divisor"; divisor: bigint };
export type Int__DivMod = { quotient: bigint; remainder: bigint };
export type Int__Value = { value: bigint };
export type Bool__Value = { value: boolean };
// Euclidean integer division (v17): quotient and remainder with
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
export function std__int__divmod(dividend: bigint, divisor: bigint): DivisionResult {
  if ((divisor === 0n)) {
    return { kind: "math.zero_divisor", divisor: divisor };
  }
  else {
    return { kind: "ok", quotient: $ailDivMod(dividend, divisor)[0], remainder: $ailDivMod(dividend, divisor)[1] };
  }
}
export function std__int__mod(value: bigint, modulus: bigint): DivisionResult {
  if ((modulus === 0n)) {
    return { kind: "math.zero_divisor", divisor: modulus };
  }
  else {
    return { kind: "ok", value: $ailDivMod(value, modulus)[1] };
  }
}
export function std__int__is_multiple(value: bigint, divisor: bigint): DivisionResult {
  if ((divisor === 0n)) {
    return { kind: "math.zero_divisor", divisor: divisor };
  }
  else {
    if (($ailDivMod(value, divisor)[1] === 0n)) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "ok", value: false };
    }
  }
}
export function std__int__is_even(value: bigint): DivisionResult {
  if (($ailDivMod(value, 2n)[1] === 0n)) {
    return { kind: "ok", value: true };
  }
  else {
    return { kind: "ok", value: false };
  }
}
export function std__int__is_odd(value: bigint): DivisionResult {
  if (($ailDivMod(value, 2n)[1] === 0n)) {
    return { kind: "ok", value: false };
  }
  else {
    return { kind: "ok", value: true };
  }
}
