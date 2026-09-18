// GENERATED from fizzbuzz.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type FizzResult = { $ail_kind: "ok"; n: bigint; word: string };
export type Fizz__Label = { word: string; n: bigint };
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
export function fizz__label(n: bigint): { $ail_kind: "ok"; n: bigint; word: string } {
  const $ail_m1 = ($ailDivMod(n, 3n)[1] === 0n);
  const $ail_m2 = ($ailDivMod(n, 5n)[1] === 0n);
  if ($ail_m1 && $ail_m2) {
    return { $ail_kind: "ok", word: "FizzBuzz", n: n };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "ok", word: "Fizz", n: n };
  }
  else if (!($ail_m1) && $ail_m2) {
    return { $ail_kind: "ok", word: "Buzz", n: n };
  }
  else {
    return { $ail_kind: "ok", word: "", n: n };
  }
}
