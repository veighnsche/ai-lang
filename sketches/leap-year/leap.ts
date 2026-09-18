// GENERATED from leap.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type LeapResult = { $ail_kind: "ok"; leap: boolean };
export type Leap__Answer = { leap: boolean };
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
export function leap__is_leap(year: bigint): { $ail_kind: "ok"; leap: boolean } {
  const $ail_m1 = ($ailDivMod(year, 400n)[1] === 0n);
  const $ail_m2 = ($ailDivMod(year, 100n)[1] === 0n);
  const $ail_m3 = ($ailDivMod(year, 4n)[1] === 0n);
  if ($ail_m1) {
    return { $ail_kind: "ok", leap: true };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", leap: false };
  }
  else if ($ail_m3) {
    return { $ail_kind: "ok", leap: true };
  }
  else {
    return { $ail_kind: "ok", leap: false };
  }
}
