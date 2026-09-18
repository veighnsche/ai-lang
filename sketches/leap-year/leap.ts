// GENERATED from leap.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type LeapResult = { $can_kind: "ok"; leap: boolean };
export type Leap__Answer = { leap: boolean };
// Euclidean integer division (a17): quotient and remainder with
// 0 <= r < |b| on every sign combination.
function $canDivMod(a: bigint, b: bigint): [bigint, bigint] {
  let q: bigint = a / b;
  let r: bigint = a - b * q;
  if (r < 0n) {
    q += b > 0n ? -1n : 1n;
    r = a - b * q;
  }
  return [q, r];
}
export function leap__is_leap(year: bigint): { $can_kind: "ok"; leap: boolean } {
  const $can_m1 = ($canDivMod(year, 400n)[1] === 0n);
  const $can_m2 = ($canDivMod(year, 100n)[1] === 0n);
  const $can_m3 = ($canDivMod(year, 4n)[1] === 0n);
  if ($can_m1) {
    return { $can_kind: "ok", leap: true };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", leap: false };
  }
  else if ($can_m3) {
    return { $can_kind: "ok", leap: true };
  }
  else {
    return { $can_kind: "ok", leap: false };
  }
}
