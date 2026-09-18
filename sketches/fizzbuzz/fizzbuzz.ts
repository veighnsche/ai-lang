// GENERATED from fizzbuzz.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type FizzResult = { $can_kind: "ok"; n: bigint; word: string };
export type Fizz__Label = { word: string; n: bigint };
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
export function fizz__label(n: bigint): { $can_kind: "ok"; n: bigint; word: string } {
  const $can_m1 = ($canDivMod(n, 3n)[1] === 0n);
  const $can_m2 = ($canDivMod(n, 5n)[1] === 0n);
  if ($can_m1 && $can_m2) {
    return { $can_kind: "ok", word: "FizzBuzz", n: n };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "ok", word: "Fizz", n: n };
  }
  else if (!($can_m1) && $can_m2) {
    return { $can_kind: "ok", word: "Buzz", n: n };
  }
  else {
    return { $can_kind: "ok", word: "", n: n };
  }
}
