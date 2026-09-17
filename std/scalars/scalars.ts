// GENERATED from scalars.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type ScalarsResult = { $ail_kind: "ok"; coefficient: bigint; scale: bigint } | { $ail_kind: "ok"; discarded: string; value: string } | { $ail_kind: "ok"; exact: boolean; stripped: bigint } | { $ail_kind: "ok"; fives: bigint; rest: bigint; twos: bigint } | { $ail_kind: "ok"; remainder: bigint; root: bigint } | { $ail_kind: "ok"; remainder: string; value: string } | { $ail_kind: "ok"; value: bigint } | { $ail_kind: "ok"; value: boolean } | { $ail_kind: "ok"; value: string } | { $ail_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } | { $ail_kind: "math.out_of_range"; value: bigint; lower: bigint; upper: bigint } | { $ail_kind: "math.dec_invalid_bounds"; lower: string; upper: string } | { $ail_kind: "math.dec_out_of_range"; value: string; lower: string; upper: string } | { $ail_kind: "math.negative_input"; value: bigint } | { $ail_kind: "math.nonpositive_input"; value: bigint } | { $ail_kind: "math.negative_exponent"; exponent: bigint } | { $ail_kind: "math.invalid_count"; n: bigint; k: bigint } | { $ail_kind: "math.negative_scale"; scale: bigint } | { $ail_kind: "math.dec_zero_divisor"; divisor: string } | { $ail_kind: "math.nonterminating_decimal"; dividend: string; divisor: string } | { $ail_kind: "math.dec_negative_input"; value: string } | { $ail_kind: "math.nonrepresentable_result"; value: string } | { $ail_kind: "convert.fractional_value"; value: string } | { $ail_kind: "convert.invalid_boolean"; value: string } | { $ail_kind: "convert.invalid_integer"; value: string } | { $ail_kind: "convert.invalid_decimal"; value: string } | { $ail_kind: "convert.invalid_boolean_encoding"; value: bigint } | { $ail_kind: "convert.invalid_dec_encoding"; value: string };
export type Int__Value = { value: bigint };
export type Dec__Value = { value: string };
export type Str__Value = { value: string };
export type Bool__Value = { value: boolean };
export type Int__Root = { root: bigint; remainder: bigint };
export type Dec__Parts = { coefficient: bigint; scale: bigint };
export type Dec__Strip = { stripped: bigint; exact: boolean };
export type Dec__Factors = { rest: bigint; twos: bigint; fives: bigint };
export type Dec__Rounded = { value: string; discarded: string };
export type Dec__RoundedQuotient = { value: string; remainder: string };
// Exact-decimal runtime (v10): canonical-digit strings, BigInt math.
function $ailDecSplit(d: string): { neg: boolean; ip: string; fp: string } {
  let neg = false;
  if (d.startsWith("-")) {
    neg = true;
    d = d.slice(1);
  }
  const dot = d.indexOf(".");
  return { neg, ip: d.slice(0, dot), fp: d.slice(dot + 1) };
}
function $ailDecNorm(ip: string, fp: string, neg: boolean): string {
  ip = ip.replace(/^0+(?=\d)/, "");
  fp = fp.replace(/0+$/, "");
  if (fp === "") {
    fp = "0";
  }
  if (ip === "0" && fp === "0") {
    return "0.0";
  }
  return (neg ? "-" : "") + ip + "." + fp;
}
function $ailDecMant(p: { neg: boolean; ip: string; fp: string }, scale: number): bigint {
  let f = p.fp;
  while (f.length < scale) {
    f += "0";
  }
  const m = BigInt(p.ip + f);
  return p.neg ? -m : m;
}
function $ailDecFromMant(m: bigint, scale: number): string {
  let neg = false;
  if (m < 0n) {
    neg = true;
    m = -m;
  }
  let digits = m.toString();
  while (digits.length < scale + 1) {
    digits = "0" + digits;
  }
  return $ailDecNorm(digits.slice(0, digits.length - scale), digits.slice(digits.length - scale), neg);
}
function $ailDecAdd(a: string, b: string): string {
  const A = $ailDecSplit(a);
  const B = $ailDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $ailDecFromMant($ailDecMant(A, s) + $ailDecMant(B, s), s);
}
function $ailDecSub(a: string, b: string): string {
  const A = $ailDecSplit(a);
  const B = $ailDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $ailDecFromMant($ailDecMant(A, s) - $ailDecMant(B, s), s);
}
function $ailDecMul(a: string, b: string): string {
  const A = $ailDecSplit(a);
  const B = $ailDecSplit(b);
  return $ailDecFromMant($ailDecMant(A, A.fp.length) * $ailDecMant(B, B.fp.length), A.fp.length + B.fp.length);
}
function $ailDecGe(a: string, b: string): boolean {
  const A = $ailDecSplit(a);
  const B = $ailDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $ailDecMant(A, s) >= $ailDecMant(B, s);
}
function $ailDecLe(a: string, b: string): boolean {
  const A = $ailDecSplit(a);
  const B = $ailDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $ailDecMant(A, s) <= $ailDecMant(B, s);
}
function $ailDecLt(a: string, b: string): boolean {
  const A = $ailDecSplit(a);
  const B = $ailDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $ailDecMant(A, s) < $ailDecMant(B, s);
}
function $ailDecParts(d: string): { coefficient: bigint; scale: bigint } {
  const p = $ailDecSplit(d);
  return { coefficient: $ailDecMant(p, p.fp.length), scale: BigInt(p.fp.length) };
}
// Byte-order string comparison: UTF-8 bytes, matching Go.
function $ailStrCmp(a: string, b: string): number {
  const A = new TextEncoder().encode(a);
  const B = new TextEncoder().encode(b);
  const n = Math.min(A.length, B.length);
  for (let i = 0; i < n; i++) {
    if (A[i] !== B[i]) {
      return A[i] < B[i] ? -1 : 1;
    }
  }
  if (A.length === B.length) {
    return 0;
  }
  return A.length < B.length ? -1 : 1;
}
function $ailStrGe(a: string, b: string): boolean {
  return $ailStrCmp(a, b) >= 0;
}
function $ailStrAt(s: string, i: bigint): bigint {
  const cps = [...s];
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("str index out of range");
  const k = Number(i);
  if (k >= cps.length) throw new Error("str index out of range");
  return BigInt(cps[k].codePointAt(0));
}
function $ailStrSlice(s: string, a: bigint, b: bigint): string {
  const cps = [...s];
  const toIdx = (x: bigint): number => { if (x < 0n || x > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("str slice out of range"); return Number(x); };
  const lo = toIdx(a), hi = toIdx(b);
  if (lo > hi || hi > cps.length) throw new Error("str slice out of range");
  return cps.slice(lo, hi).join("");
}
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
export function std__bool__not(value: boolean): ScalarsResult {
  if (value) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    return { $ail_kind: "ok", value: true };
  }
}
export function std__bool__and(left: boolean, right: boolean): ScalarsResult {
  const $ail_m1 = left;
  const $ail_m2 = right;
  if ($ail_m1 && $ail_m2) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__bool__or(left: boolean, right: boolean): ScalarsResult {
  const $ail_m3 = left;
  const $ail_m4 = right;
  if ($ail_m3) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m4) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__bool__xor(left: boolean, right: boolean): ScalarsResult {
  const $ail_m5 = left;
  const $ail_m6 = right;
  if ($ail_m5 && $ail_m6) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m5 && !($ail_m6)) {
    return { $ail_kind: "ok", value: true };
  }
  else if (!($ail_m5) && $ail_m6) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__bool__implies(left: boolean, right: boolean): ScalarsResult {
  const $ail_m7 = left;
  const $ail_m8 = right;
  if ($ail_m7 && $ail_m8) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m7 && !($ail_m8)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    return { $ail_kind: "ok", value: true };
  }
}
export function std__bool__equivalent(left: boolean, right: boolean): ScalarsResult {
  const $ail_m9 = left;
  const $ail_m10 = right;
  if ($ail_m9 && $ail_m10) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m9 && !($ail_m10)) {
    return { $ail_kind: "ok", value: false };
  }
  else if (!($ail_m9) && $ail_m10) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    return { $ail_kind: "ok", value: true };
  }
}
export function std__compare__int(left: bigint, right: bigint): ScalarsResult {
  const $ail_m11 = (left === right);
  const $ail_m12 = (left >= right);
  if ($ail_m11) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m12) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__compare__dec(left: string, right: string): ScalarsResult {
  const $ail_m13 = (left === right);
  const $ail_m14 = $ailDecGe(left, right);
  if ($ail_m13) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m14) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__compare__str(left: string, right: string): ScalarsResult {
  const $ail_m15 = (left === right);
  const $ail_m16 = $ailStrGe(left, right);
  if ($ail_m15) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m16) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__compare__bool(left: boolean, right: boolean): ScalarsResult {
  const $ail_m17 = (left === right);
  const $ail_m18 = left;
  if ($ail_m17) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m18) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__select__int(condition: boolean, when_true: bigint, when_false: bigint): ScalarsResult {
  if (condition) {
    return { $ail_kind: "ok", value: when_true };
  }
  else {
    return { $ail_kind: "ok", value: when_false };
  }
}
export function std__select__bool(condition: boolean, when_true: boolean, when_false: boolean): ScalarsResult {
  if (condition) {
    return { $ail_kind: "ok", value: when_true };
  }
  else {
    return { $ail_kind: "ok", value: when_false };
  }
}
export function std__select__str(condition: boolean, when_true: string, when_false: string): ScalarsResult {
  if (condition) {
    return { $ail_kind: "ok", value: when_true };
  }
  else {
    return { $ail_kind: "ok", value: when_false };
  }
}
export function std__select__dec(condition: boolean, when_true: string, when_false: string): ScalarsResult {
  if (condition) {
    return { $ail_kind: "ok", value: when_true };
  }
  else {
    return { $ail_kind: "ok", value: when_false };
  }
}
export function std__int__is_zero(value: bigint): ScalarsResult {
  if ((value === 0n)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__int__is_positive(value: bigint): ScalarsResult {
  const $ail_m19 = (value === 0n);
  const $ail_m20 = (value >= 0n);
  if ($ail_m19) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m20) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__int__is_negative(value: bigint): ScalarsResult {
  const $ail_m21 = (value === 0n);
  const $ail_m22 = (value <= 0n);
  if ($ail_m21) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m22) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__int__in_closed_range(value: bigint, lower: bigint, upper: bigint): ScalarsResult {
  const $ail_m23 = (lower <= upper);
  const $ail_m24 = (value >= lower);
  const $ail_m25 = (value <= upper);
  if ($ail_m23 && $ail_m24 && $ail_m25) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m23 && $ail_m24 && !($ail_m25)) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m23 && !($ail_m24)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__in_open_range(value: bigint, lower: bigint, upper: bigint): ScalarsResult {
  const $ail_m26 = (lower <= upper);
  const $ail_m27 = (value <= lower);
  const $ail_m28 = (value >= upper);
  if ($ail_m26 && $ail_m27) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m26 && !($ail_m27) && $ail_m28) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m26 && !($ail_m27) && !($ail_m28)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__is_zero(value: string): ScalarsResult {
  if ((value === "0.0")) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__dec__is_positive(value: string): ScalarsResult {
  const $ail_m29 = (value === "0.0");
  const $ail_m30 = $ailDecGe(value, "0.0");
  if ($ail_m29) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m30) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__dec__is_negative(value: string): ScalarsResult {
  const $ail_m31 = (value === "0.0");
  const $ail_m32 = $ailDecLe(value, "0.0");
  if ($ail_m31) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m32) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__dec__in_closed_range(value: string, lower: string, upper: string): ScalarsResult {
  const $ail_m33 = $ailDecLe(lower, upper);
  const $ail_m34 = $ailDecGe(value, lower);
  const $ail_m35 = $ailDecLe(value, upper);
  if ($ail_m33 && $ail_m34 && $ail_m35) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m33 && $ail_m34 && !($ail_m35)) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m33 && !($ail_m34)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    return { $ail_kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__in_open_range(value: string, lower: string, upper: string): ScalarsResult {
  const $ail_m36 = $ailDecLe(lower, upper);
  const $ail_m37 = $ailDecLe(value, lower);
  const $ail_m38 = $ailDecGe(value, upper);
  if ($ail_m36 && $ail_m37) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m36 && !($ail_m37) && $ail_m38) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m36 && !($ail_m37) && !($ail_m38)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__abs(value: bigint): ScalarsResult {
  if ((value >= 0n)) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    return { $ail_kind: "ok", value: (0n - value) };
  }
}
export function std__int__negate(value: bigint): ScalarsResult {
  return { $ail_kind: "ok", value: (0n - value) };
}
export function std__int__sign(value: bigint): ScalarsResult {
  const $ail_m39 = (value === 0n);
  const $ail_m40 = (value >= 0n);
  if ($ail_m39) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m40) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__int__min(left: bigint, right: bigint): ScalarsResult {
  if ((left <= right)) {
    return { $ail_kind: "ok", value: left };
  }
  else {
    return { $ail_kind: "ok", value: right };
  }
}
export function std__int__max(left: bigint, right: bigint): ScalarsResult {
  if ((left >= right)) {
    return { $ail_kind: "ok", value: left };
  }
  else {
    return { $ail_kind: "ok", value: right };
  }
}
export function std__int__clamp(value: bigint, lower: bigint, upper: bigint): ScalarsResult {
  const $ail_m41 = (lower <= upper);
  const $ail_m42 = (value >= lower);
  const $ail_m43 = (value <= upper);
  if ($ail_m41 && $ail_m42 && $ail_m43) {
    return { $ail_kind: "ok", value: value };
  }
  else if ($ail_m41 && $ail_m42 && !($ail_m43)) {
    return { $ail_kind: "ok", value: upper };
  }
  else if ($ail_m41 && !($ail_m42)) {
    return { $ail_kind: "ok", value: lower };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__distance(left: bigint, right: bigint): ScalarsResult {
  if ((left >= right)) {
    return { $ail_kind: "ok", value: (left - right) };
  }
  else {
    return { $ail_kind: "ok", value: (right - left) };
  }
}
export function std__int__square(value: bigint): ScalarsResult {
  return { $ail_kind: "ok", value: (value * value) };
}
export function std__int__pow_from(base: bigint, exponent: bigint): ScalarsResult {
  if ((exponent <= 0n)) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    const $ail_m44: ScalarsResult = std__int__pow_from(base, (exponent - 1n));
    switch ($ail_m44.$ail_kind) {
    case "ok": {
      const r = $ail_m44;
      return { $ail_kind: "ok", value: (base * r.value) };
    }
    }
  }
}
export function std__int__pow(base: bigint, exponent: bigint): ScalarsResult {
  if ((exponent >= 0n)) {
    const $ail_m45: ScalarsResult = std__int__pow_from(base, exponent);
    switch ($ail_m45.$ail_kind) {
    case "ok": {
      const r = $ail_m45;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_exponent", exponent: exponent };
  }
}
export function std__int__add_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): ScalarsResult {
  const $ail_m46 = (lower <= upper);
  const $ail_m47 = ((left + right) >= lower);
  const $ail_m48 = ((left + right) <= upper);
  if ($ail_m46 && $ail_m47 && $ail_m48) {
    return { $ail_kind: "ok", value: (left + right) };
  }
  else if ($ail_m46 && $ail_m47 && !($ail_m48)) {
    return { $ail_kind: "math.out_of_range", value: (left + right), lower: lower, upper: upper };
  }
  else if ($ail_m46 && !($ail_m47)) {
    return { $ail_kind: "math.out_of_range", value: (left + right), lower: lower, upper: upper };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__subtract_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): ScalarsResult {
  const $ail_m49 = (lower <= upper);
  const $ail_m50 = ((left - right) >= lower);
  const $ail_m51 = ((left - right) <= upper);
  if ($ail_m49 && $ail_m50 && $ail_m51) {
    return { $ail_kind: "ok", value: (left - right) };
  }
  else if ($ail_m49 && $ail_m50 && !($ail_m51)) {
    return { $ail_kind: "math.out_of_range", value: (left - right), lower: lower, upper: upper };
  }
  else if ($ail_m49 && !($ail_m50)) {
    return { $ail_kind: "math.out_of_range", value: (left - right), lower: lower, upper: upper };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__multiply_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): ScalarsResult {
  const $ail_m52 = (lower <= upper);
  const $ail_m53 = ((left * right) >= lower);
  const $ail_m54 = ((left * right) <= upper);
  if ($ail_m52 && $ail_m53 && $ail_m54) {
    return { $ail_kind: "ok", value: (left * right) };
  }
  else if ($ail_m52 && $ail_m53 && !($ail_m54)) {
    return { $ail_kind: "math.out_of_range", value: (left * right), lower: lower, upper: upper };
  }
  else if ($ail_m52 && !($ail_m53)) {
    return { $ail_kind: "math.out_of_range", value: (left * right), lower: lower, upper: upper };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__factorial_from(value: bigint): ScalarsResult {
  if ((value <= 0n)) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    const $ail_m55: ScalarsResult = std__int__factorial_from((value - 1n));
    switch ($ail_m55.$ail_kind) {
    case "ok": {
      const r = $ail_m55;
      return { $ail_kind: "ok", value: (value * r.value) };
    }
    }
  }
}
export function std__int__factorial(value: bigint): ScalarsResult {
  if ((value >= 0n)) {
    const $ail_m56: ScalarsResult = std__int__factorial_from(value);
    switch ($ail_m56.$ail_kind) {
    case "ok": {
      const r = $ail_m56;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_input", value: value };
  }
}
export function std__int__sum_from(upper: bigint): ScalarsResult {
  if ((upper <= 0n)) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    const $ail_m57: ScalarsResult = std__int__sum_from((upper - 1n));
    switch ($ail_m57.$ail_kind) {
    case "ok": {
      const r = $ail_m57;
      return { $ail_kind: "ok", value: (upper + r.value) };
    }
    }
  }
}
export function std__int__sum_to(upper: bigint): ScalarsResult {
  if ((upper >= 0n)) {
    const $ail_m58: ScalarsResult = std__int__sum_from(upper);
    switch ($ail_m58.$ail_kind) {
    case "ok": {
      const r = $ail_m58;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_input", value: upper };
  }
}
export function std__int__binomial_from(n: bigint, k: bigint, i: bigint, num: bigint, den: bigint): ScalarsResult {
  if ((i <= 0n)) {
    return { $ail_kind: "ok", value: $ailDivMod(num, den)[0] };
  }
  else {
    const $ail_m59: ScalarsResult = std__int__binomial_from(n, k, (i - 1n), (num * ((n - k) + i)), (den * i));
    switch ($ail_m59.$ail_kind) {
    case "ok": {
      const r = $ail_m59;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__int__binomial(n: bigint, k: bigint): ScalarsResult {
  const $ail_m60 = (k >= 0n);
  const $ail_m61 = (k <= n);
  if ($ail_m60 && $ail_m61) {
    const $ail_m62: ScalarsResult = std__int__binomial_from(n, k, k, 1n, 1n);
    switch ($ail_m62.$ail_kind) {
    case "ok": {
      const r = $ail_m62;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else if ($ail_m60 && !($ail_m61)) {
    return { $ail_kind: "math.invalid_count", n: n, k: k };
  }
  else {
    return { $ail_kind: "math.invalid_count", n: n, k: k };
  }
}
export function std__int__is_prime_from(value: bigint, e: bigint): ScalarsResult {
  if ((e <= 0n)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    if (($ailDivMod(value, (e + 1n))[1] === 0n)) {
      return { $ail_kind: "ok", value: false };
    }
    else {
      const $ail_m63: ScalarsResult = std__int__is_prime_from(value, (e - 1n));
      switch ($ail_m63.$ail_kind) {
      case "ok": {
        const r = $ail_m63;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
  }
}
export function std__int__is_prime(value: bigint): ScalarsResult {
  if ((value >= 2n)) {
    const $ail_m64: ScalarsResult = std__int__is_prime_from(value, (value - 2n));
    switch ($ail_m64.$ail_kind) {
    case "ok": {
      const r = $ail_m64;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__int__next_power_of_two_pow2(e: bigint): ScalarsResult {
  if ((e <= 0n)) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    const $ail_m65: ScalarsResult = std__int__next_power_of_two_pow2((e - 1n));
    switch ($ail_m65.$ail_kind) {
    case "ok": {
      const r = $ail_m65;
      return { $ail_kind: "ok", value: (2n * r.value) };
    }
    }
  }
}
export function std__int__next_power_of_two_from(value: bigint, e: bigint, best: bigint): ScalarsResult {
  if ((e <= 0n)) {
    if ((1n >= value)) {
      return { $ail_kind: "ok", value: 1n };
    }
    else {
      return { $ail_kind: "ok", value: best };
    }
  }
  else {
    const $ail_m66: ScalarsResult = std__int__next_power_of_two_pow2(e);
    switch ($ail_m66.$ail_kind) {
    case "ok": {
      const t = $ail_m66;
      if ((t.value >= value)) {
        const $ail_m67: ScalarsResult = std__int__next_power_of_two_from(value, (e - 1n), t.value);
        switch ($ail_m67.$ail_kind) {
        case "ok": {
          const r = $ail_m67;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
      else {
        const $ail_m68: ScalarsResult = std__int__next_power_of_two_from(value, (e - 1n), best);
        switch ($ail_m68.$ail_kind) {
        case "ok": {
          const r = $ail_m68;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
    }
    }
  }
}
export function std__int__next_power_of_two(value: bigint): ScalarsResult {
  if ((value >= 1n)) {
    const $ail_m69: ScalarsResult = std__int__next_power_of_two_from(value, value, 1n);
    switch ($ail_m69.$ail_kind) {
    case "ok": {
      const r = $ail_m69;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { $ail_kind: "math.nonpositive_input", value: value };
  }
}
export function std__int__sqrt_floor_search(value: bigint, lo: bigint, hi: bigint): ScalarsResult {
  if (((hi - lo) <= 1n)) {
    return { $ail_kind: "ok", value: lo };
  }
  else {
    if ((($ailDivMod((lo + hi), 2n)[0] * $ailDivMod((lo + hi), 2n)[0]) <= value)) {
      const $ail_m70: ScalarsResult = std__int__sqrt_floor_search(value, $ailDivMod((lo + hi), 2n)[0], hi);
      switch ($ail_m70.$ail_kind) {
      case "ok": {
        const r = $ail_m70;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
    else {
      const $ail_m71: ScalarsResult = std__int__sqrt_floor_search(value, lo, $ailDivMod((lo + hi), 2n)[0]);
      switch ($ail_m71.$ail_kind) {
      case "ok": {
        const r = $ail_m71;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
  }
}
export function std__int__sqrt_floor(value: bigint): ScalarsResult {
  if ((value >= 0n)) {
    const $ail_m72: ScalarsResult = std__int__sqrt_floor_search(value, 0n, (value + 1n));
    switch ($ail_m72.$ail_kind) {
    case "ok": {
      const r = $ail_m72;
      return { $ail_kind: "ok", root: r.value, remainder: (value - (r.value * r.value)) };
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_input", value: value };
  }
}
export function std__int__gcd_euclid(a: bigint, b: bigint): ScalarsResult {
  if ((b <= 0n)) {
    return { $ail_kind: "ok", value: a };
  }
  else {
    const $ail_m73: ScalarsResult = std__int__gcd_euclid(b, $ailDivMod(a, b)[1]);
    switch ($ail_m73.$ail_kind) {
    case "ok": {
      const r = $ail_m73;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__int__gcd(left: bigint, right: bigint): ScalarsResult {
  const $ail_m74: ScalarsResult = std__int__abs(left);
  switch ($ail_m74.$ail_kind) {
  case "ok": {
    const a = $ail_m74;
    const $ail_m75: ScalarsResult = std__int__abs(right);
    switch ($ail_m75.$ail_kind) {
    case "ok": {
      const b = $ail_m75;
      const $ail_m76: ScalarsResult = std__int__gcd_euclid(a.value, b.value);
      switch ($ail_m76.$ail_kind) {
      case "ok": {
        const r = $ail_m76;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
    }
  }
  }
}
export function std__int__lcm(left: bigint, right: bigint): ScalarsResult {
  const $ail_m77 = (left === 0n);
  const $ail_m78 = (right === 0n);
  if ($ail_m77) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m78) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    const $ail_m79: ScalarsResult = std__int__gcd(left, right);
    switch ($ail_m79.$ail_kind) {
    case "ok": {
      const g = $ail_m79;
      const $ail_m80: ScalarsResult = std__int__abs(left);
      switch ($ail_m80.$ail_kind) {
      case "ok": {
        const a = $ail_m80;
        const $ail_m81: ScalarsResult = std__int__abs(right);
        switch ($ail_m81.$ail_kind) {
        case "ok": {
          const b = $ail_m81;
          return { $ail_kind: "ok", value: $ailDivMod((a.value * b.value), g.value)[0] };
        }
        }
      }
      }
    }
    }
  }
}
export function std__dec__abs(value: string): ScalarsResult {
  if ($ailDecGe(value, "0.0")) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    return { $ail_kind: "ok", value: $ailDecSub("0.0", value) };
  }
}
export function std__dec__negate(value: string): ScalarsResult {
  return { $ail_kind: "ok", value: $ailDecSub("0.0", value) };
}
export function std__dec__sign(value: string): ScalarsResult {
  const $ail_m82 = (value === "0.0");
  const $ail_m83 = $ailDecGe(value, "0.0");
  if ($ail_m82) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m83) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__dec__min(left: string, right: string): ScalarsResult {
  if ($ailDecLe(left, right)) {
    return { $ail_kind: "ok", value: left };
  }
  else {
    return { $ail_kind: "ok", value: right };
  }
}
export function std__dec__max(left: string, right: string): ScalarsResult {
  if ($ailDecGe(left, right)) {
    return { $ail_kind: "ok", value: left };
  }
  else {
    return { $ail_kind: "ok", value: right };
  }
}
export function std__dec__clamp(value: string, lower: string, upper: string): ScalarsResult {
  const $ail_m84 = $ailDecLe(lower, upper);
  const $ail_m85 = $ailDecGe(value, lower);
  const $ail_m86 = $ailDecLe(value, upper);
  if ($ail_m84 && $ail_m85 && $ail_m86) {
    return { $ail_kind: "ok", value: value };
  }
  else if ($ail_m84 && $ail_m85 && !($ail_m86)) {
    return { $ail_kind: "ok", value: upper };
  }
  else if ($ail_m84 && !($ail_m85)) {
    return { $ail_kind: "ok", value: lower };
  }
  else {
    return { $ail_kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__distance(left: string, right: string): ScalarsResult {
  if ($ailDecGe(left, right)) {
    return { $ail_kind: "ok", value: $ailDecSub(left, right) };
  }
  else {
    return { $ail_kind: "ok", value: $ailDecSub(right, left) };
  }
}
export function std__dec__square(value: string): ScalarsResult {
  return { $ail_kind: "ok", value: $ailDecMul(value, value) };
}
export function std__dec__pow_from(base: string, exponent: bigint): ScalarsResult {
  if ((exponent <= 0n)) {
    return { $ail_kind: "ok", value: "1.0" };
  }
  else {
    const $ail_m87: ScalarsResult = std__dec__pow_from(base, (exponent - 1n));
    switch ($ail_m87.$ail_kind) {
    case "ok": {
      const r = $ail_m87;
      return { $ail_kind: "ok", value: $ailDecMul(base, r.value) };
    }
    }
  }
}
export function std__dec__pow(base: string, exponent: bigint): ScalarsResult {
  if ((exponent >= 0n)) {
    const $ail_m88: ScalarsResult = std__dec__pow_from(base, exponent);
    switch ($ail_m88.$ail_kind) {
    case "ok": {
      const r = $ail_m88;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_exponent", exponent: exponent };
  }
}
export function std__dec__lerp(start: string, end: string, fraction: string): ScalarsResult {
  const $ail_m89 = $ailDecGe(fraction, "0.0");
  const $ail_m90 = $ailDecLe(fraction, "1.0");
  if ($ail_m89 && $ail_m90) {
    return { $ail_kind: "ok", value: $ailDecAdd(start, $ailDecMul($ailDecSub(end, start), fraction)) };
  }
  else if ($ail_m89 && !($ail_m90)) {
    return { $ail_kind: "math.dec_out_of_range", value: fraction, lower: "0.0", upper: "1.0" };
  }
  else {
    return { $ail_kind: "math.dec_out_of_range", value: fraction, lower: "0.0", upper: "1.0" };
  }
}
export function std__dec__scale_by_power_of_ten_from(value: string, e: bigint, factor: string): ScalarsResult {
  if ((e <= 0n)) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    const $ail_m91: ScalarsResult = std__dec__scale_by_power_of_ten_from($ailDecMul(value, factor), (e - 1n), factor);
    switch ($ail_m91.$ail_kind) {
    case "ok": {
      const r = $ail_m91;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__dec__scale_by_power_of_ten(value: string, exponent: bigint): ScalarsResult {
  if ((exponent >= 0n)) {
    const $ail_m92: ScalarsResult = std__dec__scale_by_power_of_ten_from(value, exponent, "10.0");
    switch ($ail_m92.$ail_kind) {
    case "ok": {
      const r = $ail_m92;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    const $ail_m93: ScalarsResult = std__dec__scale_by_power_of_ten_from(value, (0n - exponent), "0.1");
    switch ($ail_m93.$ail_kind) {
    case "ok": {
      const r = $ail_m93;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__dec__from_parts(coefficient: bigint, scale: bigint): ScalarsResult {
  if ((scale >= 0n)) {
    const $ail_m94: ScalarsResult = std__convert__int_to_dec(coefficient);
    switch ($ail_m94.$ail_kind) {
    case "ok": {
      const whole = $ail_m94;
      const $ail_m95: ScalarsResult = std__dec__scale_by_power_of_ten(whole.value, (0n - scale));
      switch ($ail_m95.$ail_kind) {
      case "ok": {
        const moved = $ail_m95;
        return { $ail_kind: "ok", value: moved.value };
      }
      }
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_scale", scale: scale };
  }
}
export function std__dec__parts(value: string): ScalarsResult {
  const $ail_m96: { $ail_kind: "ok", coefficient: bigint, scale: bigint } = { $ail_kind: "ok", ...$ailDecParts(value) };
  switch ($ail_m96.$ail_kind) {
  case "ok": {
    const p = $ail_m96;
    return { $ail_kind: "ok", coefficient: p.coefficient, scale: p.scale };
  }
  }
}
export function std__dec__strip_from(remaining: bigint, current: bigint, exact: boolean): ScalarsResult {
  if ((remaining <= 0n)) {
    return { $ail_kind: "ok", stripped: current, exact: exact };
  }
  else {
    if (($ailDivMod(current, 10n)[1] === 0n)) {
      const $ail_m97: ScalarsResult = std__dec__strip_from((remaining - 1n), $ailDivMod(current, 10n)[0], exact);
      switch ($ail_m97.$ail_kind) {
      case "ok": {
        const r = $ail_m97;
        return { $ail_kind: "ok", stripped: r.stripped, exact: r.exact };
      }
      }
    }
    else {
      const $ail_m98: ScalarsResult = std__dec__strip_from((remaining - 1n), $ailDivMod(current, 10n)[0], false);
      switch ($ail_m98.$ail_kind) {
      case "ok": {
        const r = $ail_m98;
        return { $ail_kind: "ok", stripped: r.stripped, exact: r.exact };
      }
      }
    }
  }
}
export function std__dec__truncate(value: string): ScalarsResult {
  const $ail_m99: ScalarsResult = std__dec__parts(value);
  switch ($ail_m99.$ail_kind) {
  case "ok": {
    const p = $ail_m99;
    const $ail_m100: ScalarsResult = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($ail_m100.$ail_kind) {
    case "ok": {
      const s = $ail_m100;
      if (s.exact) {
        const $ail_m101: ScalarsResult = std__convert__int_to_dec(s.stripped);
        switch ($ail_m101.$ail_kind) {
        case "ok": {
          const whole = $ail_m101;
          return { $ail_kind: "ok", value: whole.value };
        }
        }
      }
      else {
        if ((s.stripped >= 0n)) {
          const $ail_m102: ScalarsResult = std__convert__int_to_dec(s.stripped);
          switch ($ail_m102.$ail_kind) {
          case "ok": {
            const whole = $ail_m102;
            return { $ail_kind: "ok", value: whole.value };
          }
          }
        }
        else {
          const $ail_m103: ScalarsResult = std__convert__int_to_dec((s.stripped + 1n));
          switch ($ail_m103.$ail_kind) {
          case "ok": {
            const whole = $ail_m103;
            return { $ail_kind: "ok", value: whole.value };
          }
          }
        }
      }
    }
    }
  }
  }
}
export function std__dec__floor(value: string): ScalarsResult {
  const $ail_m104: ScalarsResult = std__dec__parts(value);
  switch ($ail_m104.$ail_kind) {
  case "ok": {
    const p = $ail_m104;
    const $ail_m105: ScalarsResult = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($ail_m105.$ail_kind) {
    case "ok": {
      const s = $ail_m105;
      const $ail_m106: ScalarsResult = std__convert__int_to_dec(s.stripped);
      switch ($ail_m106.$ail_kind) {
      case "ok": {
        const whole = $ail_m106;
        return { $ail_kind: "ok", value: whole.value };
      }
      }
    }
    }
  }
  }
}
export function std__dec__ceil(value: string): ScalarsResult {
  const $ail_m107: ScalarsResult = std__dec__parts(value);
  switch ($ail_m107.$ail_kind) {
  case "ok": {
    const p = $ail_m107;
    const $ail_m108: ScalarsResult = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($ail_m108.$ail_kind) {
    case "ok": {
      const s = $ail_m108;
      if (s.exact) {
        const $ail_m109: ScalarsResult = std__convert__int_to_dec(s.stripped);
        switch ($ail_m109.$ail_kind) {
        case "ok": {
          const whole = $ail_m109;
          return { $ail_kind: "ok", value: whole.value };
        }
        }
      }
      else {
        const $ail_m110: ScalarsResult = std__convert__int_to_dec((s.stripped + 1n));
        switch ($ail_m110.$ail_kind) {
        case "ok": {
          const whole = $ail_m110;
          return { $ail_kind: "ok", value: whole.value };
        }
        }
      }
    }
    }
  }
  }
}
export function std__dec__strip257_from(fuel: bigint, current: bigint, twos: bigint, fives: bigint): ScalarsResult {
  if ((fuel <= 0n)) {
    return { $ail_kind: "ok", rest: current, twos: twos, fives: fives };
  }
  else {
    if ((current === 1n)) {
      return { $ail_kind: "ok", rest: 1n, twos: twos, fives: fives };
    }
    else {
      if (($ailDivMod(current, 2n)[1] === 0n)) {
        const $ail_m111: ScalarsResult = std__dec__strip257_from((fuel - 1n), $ailDivMod(current, 2n)[0], (twos + 1n), fives);
        switch ($ail_m111.$ail_kind) {
        case "ok": {
          const r = $ail_m111;
          return { $ail_kind: "ok", rest: r.rest, twos: r.twos, fives: r.fives };
        }
        }
      }
      else {
        if (($ailDivMod(current, 5n)[1] === 0n)) {
          const $ail_m112: ScalarsResult = std__dec__strip257_from((fuel - 1n), $ailDivMod(current, 5n)[0], twos, (fives + 1n));
          switch ($ail_m112.$ail_kind) {
          case "ok": {
            const r = $ail_m112;
            return { $ail_kind: "ok", rest: r.rest, twos: r.twos, fives: r.fives };
          }
          }
        }
        else {
          return { $ail_kind: "ok", rest: current, twos: twos, fives: fives };
        }
      }
    }
  }
}
export function std__dec__divide_exact(dividend: string, divisor: string): ScalarsResult {
  const $ail_m113: ScalarsResult = std__dec__parts(divisor);
  switch ($ail_m113.$ail_kind) {
  case "ok": {
    const d = $ail_m113;
    if ((d.coefficient === 0n)) {
      return { $ail_kind: "math.dec_zero_divisor", divisor: divisor };
    }
    else {
      const $ail_m114: ScalarsResult = std__dec__parts(dividend);
      switch ($ail_m114.$ail_kind) {
      case "ok": {
        const n = $ail_m114;
        const $ail_m115: ScalarsResult = std__int__pow_from(10n, d.scale);
        switch ($ail_m115.$ail_kind) {
        case "ok": {
          const pds = $ail_m115;
          const $ail_m116: ScalarsResult = std__int__pow_from(10n, n.scale);
          switch ($ail_m116.$ail_kind) {
          case "ok": {
            const pns = $ail_m116;
            const $ail_m117: ScalarsResult = std__int__abs((n.coefficient * pns.value));
            switch ($ail_m117.$ail_kind) {
            case "ok": {
              const an = $ail_m117;
              const $ail_m118: ScalarsResult = std__int__abs((d.coefficient * pds.value));
              switch ($ail_m118.$ail_kind) {
              case "ok": {
                const ad = $ail_m118;
                const $ail_m119: ScalarsResult = std__int__gcd(an.value, ad.value);
                switch ($ail_m119.$ail_kind) {
                case "ok": {
                  const g = $ail_m119;
                  const $ail_m120: ScalarsResult = std__dec__strip257_from($ailDivMod(ad.value, g.value)[0], $ailDivMod(ad.value, g.value)[0], 0n, 0n);
                  switch ($ail_m120.$ail_kind) {
                  case "ok": {
                    const f = $ail_m120;
                    if ((f.rest === 1n)) {
                      const $ail_m121: ScalarsResult = std__int__max(f.twos, f.fives);
                      switch ($ail_m121.$ail_kind) {
                      case "ok": {
                        const m = $ail_m121;
                        const $ail_m122: ScalarsResult = std__int__pow_from(10n, m.value);
                        switch ($ail_m122.$ail_kind) {
                        case "ok": {
                          const p10 = $ail_m122;
                          const $ail_m123: ScalarsResult = std__convert__int_to_dec(($ailDivMod(an.value, g.value)[0] * $ailDivMod(p10.value, $ailDivMod(ad.value, g.value)[0])[0]));
                          switch ($ail_m123.$ail_kind) {
                          case "ok": {
                            const qi = $ail_m123;
                            const $ail_m124: ScalarsResult = std__dec__scale_by_power_of_ten(qi.value, (0n - m.value));
                            switch ($ail_m124.$ail_kind) {
                            case "ok": {
                              const q = $ail_m124;
                              const $ail_m125 = ((n.coefficient * pns.value) < 0n);
                              const $ail_m126 = ((d.coefficient * pds.value) < 0n);
                              if ($ail_m125 && $ail_m126) {
                                return { $ail_kind: "ok", value: q.value };
                              }
                              else if ($ail_m125 && !($ail_m126)) {
                                return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0") };
                              }
                              else if (!($ail_m125) && $ail_m126) {
                                return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0") };
                              }
                              else {
                                return { $ail_kind: "ok", value: q.value };
                              }
                            }
                            }
                          }
                          }
                        }
                        }
                      }
                      }
                    }
                    else {
                      return { $ail_kind: "math.nonterminating_decimal", dividend: dividend, divisor: divisor };
                    }
                  }
                  }
                }
                }
              }
              }
            }
            }
          }
          }
        }
        }
      }
      }
    }
  }
  }
}
export function std__dec__round_half_even(value: string, scale: bigint): ScalarsResult {
  if ((scale >= 0n)) {
    const $ail_m127: ScalarsResult = std__dec__parts(value);
    switch ($ail_m127.$ail_kind) {
    case "ok": {
      const p = $ail_m127;
      if ((p.scale <= scale)) {
        return { $ail_kind: "ok", value: value, discarded: "0.0" };
      }
      else {
        const $ail_m128: ScalarsResult = std__int__pow_from(10n, (p.scale - scale));
        switch ($ail_m128.$ail_kind) {
        case "ok": {
          const pd = $ail_m128;
          const $ail_m129: ScalarsResult = std__int__abs(p.coefficient);
          switch ($ail_m129.$ail_kind) {
          case "ok": {
            const a = $ail_m129;
            const $ail_m130 = (($ailDivMod(a.value, pd.value)[1] * 2n) === pd.value);
            const $ail_m131 = ($ailDivMod($ailDivMod(a.value, pd.value)[0], 2n)[1] === 0n);
            const $ail_m132 = (($ailDivMod(a.value, pd.value)[1] * 2n) > pd.value);
            if ($ail_m130 && $ail_m131) {
              const $ail_m133: ScalarsResult = std__convert__int_to_dec($ailDivMod(a.value, pd.value)[0]);
              switch ($ail_m133.$ail_kind) {
              case "ok": {
                const qi = $ail_m133;
                const $ail_m134: ScalarsResult = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($ail_m134.$ail_kind) {
                case "ok": {
                  const q = $ail_m134;
                  if ($ailDecLt(value, "0.0")) {
                    return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), discarded: $ailDecSub(value, $ailDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $ail_kind: "ok", value: q.value, discarded: $ailDecSub(value, q.value) };
                  }
                }
                }
              }
              }
            }
            else if ($ail_m130 && !($ail_m131)) {
              const $ail_m135: ScalarsResult = std__convert__int_to_dec(($ailDivMod(a.value, pd.value)[0] + 1n));
              switch ($ail_m135.$ail_kind) {
              case "ok": {
                const qi = $ail_m135;
                const $ail_m136: ScalarsResult = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($ail_m136.$ail_kind) {
                case "ok": {
                  const q = $ail_m136;
                  if ($ailDecLt(value, "0.0")) {
                    return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), discarded: $ailDecSub(value, $ailDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $ail_kind: "ok", value: q.value, discarded: $ailDecSub(value, q.value) };
                  }
                }
                }
              }
              }
            }
            else if (!($ail_m130) && $ail_m132) {
              const $ail_m137: ScalarsResult = std__convert__int_to_dec(($ailDivMod(a.value, pd.value)[0] + 1n));
              switch ($ail_m137.$ail_kind) {
              case "ok": {
                const qi = $ail_m137;
                const $ail_m138: ScalarsResult = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($ail_m138.$ail_kind) {
                case "ok": {
                  const q = $ail_m138;
                  if ($ailDecLt(value, "0.0")) {
                    return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), discarded: $ailDecSub(value, $ailDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $ail_kind: "ok", value: q.value, discarded: $ailDecSub(value, q.value) };
                  }
                }
                }
              }
              }
            }
            else {
              const $ail_m139: ScalarsResult = std__convert__int_to_dec($ailDivMod(a.value, pd.value)[0]);
              switch ($ail_m139.$ail_kind) {
              case "ok": {
                const qi = $ail_m139;
                const $ail_m140: ScalarsResult = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($ail_m140.$ail_kind) {
                case "ok": {
                  const q = $ail_m140;
                  if ($ailDecLt(value, "0.0")) {
                    return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), discarded: $ailDecSub(value, $ailDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $ail_kind: "ok", value: q.value, discarded: $ailDecSub(value, q.value) };
                  }
                }
                }
              }
              }
            }
          }
          }
        }
        }
      }
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_scale", scale: scale };
  }
}
export function std__dec__divide_round_half_even_result(qm: bigint, scale: bigint, num: bigint, den: bigint, dividend: string, divisor: string): ScalarsResult {
  const $ail_m141: ScalarsResult = std__convert__int_to_dec(qm);
  switch ($ail_m141.$ail_kind) {
  case "ok": {
    const qi = $ail_m141;
    const $ail_m142: ScalarsResult = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
    switch ($ail_m142.$ail_kind) {
    case "ok": {
      const q = $ail_m142;
      const $ail_m143 = (num < 0n);
      const $ail_m144 = (den < 0n);
      if ($ail_m143 && $ail_m144) {
        return { $ail_kind: "ok", value: q.value, remainder: $ailDecSub(dividend, $ailDecMul(divisor, q.value)) };
      }
      else if ($ail_m143 && !($ail_m144)) {
        return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), remainder: $ailDecSub(dividend, $ailDecMul(divisor, $ailDecMul(q.value, "-1.0"))) };
      }
      else if (!($ail_m143) && $ail_m144) {
        return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), remainder: $ailDecSub(dividend, $ailDecMul(divisor, $ailDecMul(q.value, "-1.0"))) };
      }
      else {
        return { $ail_kind: "ok", value: q.value, remainder: $ailDecSub(dividend, $ailDecMul(divisor, q.value)) };
      }
    }
    }
  }
  }
}
export function std__dec__divide_round_half_even(dividend: string, divisor: string, scale: bigint): ScalarsResult {
  const $ail_m145: ScalarsResult = std__dec__parts(divisor);
  switch ($ail_m145.$ail_kind) {
  case "ok": {
    const d = $ail_m145;
    if ((d.coefficient === 0n)) {
      return { $ail_kind: "math.dec_zero_divisor", divisor: divisor };
    }
    else {
      if ((scale >= 0n)) {
        const $ail_m146: ScalarsResult = std__dec__parts(dividend);
        switch ($ail_m146.$ail_kind) {
        case "ok": {
          const n = $ail_m146;
          const $ail_m147: ScalarsResult = std__int__pow_from(10n, d.scale);
          switch ($ail_m147.$ail_kind) {
          case "ok": {
            const pds = $ail_m147;
            const $ail_m148: ScalarsResult = std__int__pow_from(10n, n.scale);
            switch ($ail_m148.$ail_kind) {
            case "ok": {
              const pns = $ail_m148;
              const $ail_m149: ScalarsResult = std__int__abs((n.coefficient * pns.value));
              switch ($ail_m149.$ail_kind) {
              case "ok": {
                const an = $ail_m149;
                const $ail_m150: ScalarsResult = std__int__abs((d.coefficient * pds.value));
                switch ($ail_m150.$ail_kind) {
                case "ok": {
                  const ad = $ail_m150;
                  const $ail_m151: ScalarsResult = std__int__pow_from(10n, scale);
                  switch ($ail_m151.$ail_kind) {
                  case "ok": {
                    const psc = $ail_m151;
                    const $ail_m152 = (($ailDivMod((an.value * psc.value), ad.value)[1] * 2n) === ad.value);
                    const $ail_m153 = ($ailDivMod($ailDivMod((an.value * psc.value), ad.value)[0], 2n)[1] === 0n);
                    const $ail_m154 = (($ailDivMod((an.value * psc.value), ad.value)[1] * 2n) > ad.value);
                    if ($ail_m152 && $ail_m153) {
                      const $ail_m155: ScalarsResult = std__dec__divide_round_half_even_result($ailDivMod((an.value * psc.value), ad.value)[0], scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($ail_m155.$ail_kind) {
                      case "ok": {
                        const r = $ail_m155;
                        return { $ail_kind: "ok", value: r.value, remainder: r.remainder };
                      }
                      }
                    }
                    else if ($ail_m152 && !($ail_m153)) {
                      const $ail_m156: ScalarsResult = std__dec__divide_round_half_even_result(($ailDivMod((an.value * psc.value), ad.value)[0] + 1n), scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($ail_m156.$ail_kind) {
                      case "ok": {
                        const r = $ail_m156;
                        return { $ail_kind: "ok", value: r.value, remainder: r.remainder };
                      }
                      }
                    }
                    else if (!($ail_m152) && $ail_m154) {
                      const $ail_m157: ScalarsResult = std__dec__divide_round_half_even_result(($ailDivMod((an.value * psc.value), ad.value)[0] + 1n), scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($ail_m157.$ail_kind) {
                      case "ok": {
                        const r = $ail_m157;
                        return { $ail_kind: "ok", value: r.value, remainder: r.remainder };
                      }
                      }
                    }
                    else {
                      const $ail_m158: ScalarsResult = std__dec__divide_round_half_even_result($ailDivMod((an.value * psc.value), ad.value)[0], scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($ail_m158.$ail_kind) {
                      case "ok": {
                        const r = $ail_m158;
                        return { $ail_kind: "ok", value: r.value, remainder: r.remainder };
                      }
                      }
                    }
                  }
                  }
                }
                }
              }
              }
            }
            }
          }
          }
        }
        }
      }
      else {
        return { $ail_kind: "math.negative_scale", scale: scale };
      }
    }
  }
  }
}
export function std__convert__dec_to_str(value: string): ScalarsResult {
  const $ail_m159: ScalarsResult = std__dec__parts(value);
  switch ($ail_m159.$ail_kind) {
  case "ok": {
    const p = $ail_m159;
    const $ail_m160: ScalarsResult = std__int__abs(p.coefficient);
    switch ($ail_m160.$ail_kind) {
    case "ok": {
      const a = $ail_m160;
      const $ail_m161: ScalarsResult = std__int__pow_from(10n, p.scale);
      switch ($ail_m161.$ail_kind) {
      case "ok": {
        const p10 = $ail_m161;
        const $ail_m162: ScalarsResult = std__convert__int_to_str($ailDivMod(a.value, p10.value)[0]);
        switch ($ail_m162.$ail_kind) {
        case "ok": {
          const ip = $ail_m162;
          const $ail_m163: ScalarsResult = std__convert__int_to_str(($ailDivMod(a.value, p10.value)[1] + p10.value));
          switch ($ail_m163.$ail_kind) {
          case "ok": {
            const fp = $ail_m163;
            if ((p.coefficient < 0n)) {
              return { $ail_kind: "ok", value: ((("-" + ip.value) + ".") + $ailStrSlice(fp.value, 1n, (p.scale + 1n))) };
            }
            else {
              return { $ail_kind: "ok", value: ((ip.value + ".") + $ailStrSlice(fp.value, 1n, (p.scale + 1n))) };
            }
          }
          }
        }
        }
      }
      }
    }
    }
  }
  }
}
export function std__convert__dec_to_int_exact(value: string): ScalarsResult {
  const $ail_m164: ScalarsResult = std__dec__parts(value);
  switch ($ail_m164.$ail_kind) {
  case "ok": {
    const p = $ail_m164;
    const $ail_m165: ScalarsResult = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($ail_m165.$ail_kind) {
    case "ok": {
      const s = $ail_m165;
      if (s.exact) {
        return { $ail_kind: "ok", value: s.stripped };
      }
      else {
        return { $ail_kind: "convert.fractional_value", value: value };
      }
    }
    }
  }
  }
}
export function std__dec__sqrt_exact(value: string): ScalarsResult {
  const $ail_m166: ScalarsResult = std__dec__parts(value);
  switch ($ail_m166.$ail_kind) {
  case "ok": {
    const p = $ail_m166;
    if ((p.coefficient < 0n)) {
      return { $ail_kind: "math.dec_negative_input", value: value };
    }
    else {
      if (($ailDivMod(p.scale, 2n)[1] === 0n)) {
        const $ail_m167: ScalarsResult = std__int__sqrt_floor_search(p.coefficient, 0n, (p.coefficient + 1n));
        switch ($ail_m167.$ail_kind) {
        case "ok": {
          const r = $ail_m167;
          if (((p.coefficient - (r.value * r.value)) === 0n)) {
            const $ail_m168: ScalarsResult = std__convert__int_to_dec(r.value);
            switch ($ail_m168.$ail_kind) {
            case "ok": {
              const qi = $ail_m168;
              const $ail_m169: ScalarsResult = std__dec__scale_by_power_of_ten(qi.value, (0n - $ailDivMod(p.scale, 2n)[0]));
              switch ($ail_m169.$ail_kind) {
              case "ok": {
                const q = $ail_m169;
                return { $ail_kind: "ok", value: q.value };
              }
              }
            }
            }
          }
          else {
            return { $ail_kind: "math.nonrepresentable_result", value: value };
          }
        }
        }
      }
      else {
        const $ail_m170: ScalarsResult = std__int__sqrt_floor_search((p.coefficient * 10n), 0n, ((p.coefficient * 10n) + 1n));
        switch ($ail_m170.$ail_kind) {
        case "ok": {
          const r = $ail_m170;
          if ((((p.coefficient * 10n) - (r.value * r.value)) === 0n)) {
            const $ail_m171: ScalarsResult = std__convert__int_to_dec(r.value);
            switch ($ail_m171.$ail_kind) {
            case "ok": {
              const qi = $ail_m171;
              const $ail_m172: ScalarsResult = std__dec__scale_by_power_of_ten(qi.value, (0n - $ailDivMod((p.scale + 1n), 2n)[0]));
              switch ($ail_m172.$ail_kind) {
              case "ok": {
                const q = $ail_m172;
                return { $ail_kind: "ok", value: q.value };
              }
              }
            }
            }
          }
          else {
            return { $ail_kind: "math.nonrepresentable_result", value: value };
          }
        }
        }
      }
    }
  }
  }
}
export function std__convert__bool_to_str(value: boolean): ScalarsResult {
  if (value) {
    return { $ail_kind: "ok", value: "true" };
  }
  else {
    return { $ail_kind: "ok", value: "false" };
  }
}
export function std__convert__str_to_bool(value: string): ScalarsResult {
  if (value === "true") {
    return { $ail_kind: "ok", value: true };
  }
  else if (value === "false") {
    return { $ail_kind: "ok", value: false };
  }
  return { $ail_kind: "convert.invalid_boolean", value: value };
}
export function std__convert__bool_to_int(value: boolean): ScalarsResult {
  if (value) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: 0n };
  }
}
export function std__convert__int_to_bool(value: bigint): ScalarsResult {
  const $ail_m173 = (value === 0n);
  const $ail_m174 = (value === 1n);
  if ($ail_m173) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m174) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "convert.invalid_boolean_encoding", value: value };
  }
}
export function std__convert__bool_to_dec(value: boolean): ScalarsResult {
  if (value) {
    return { $ail_kind: "ok", value: "1.0" };
  }
  else {
    return { $ail_kind: "ok", value: "0.0" };
  }
}
export function std__convert__dec_to_bool(value: string): ScalarsResult {
  if ((value === "0.0")) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    if ((value === "1.0")) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "convert.invalid_dec_encoding", value: value };
    }
  }
}
export function std__convert__int_to_dec_from(m: bigint, acc: string, step: string): ScalarsResult {
  if ((m <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    const $ail_m175: ScalarsResult = std__convert__int_to_dec_from((m - 1n), $ailDecAdd(acc, step), step);
    switch ($ail_m175.$ail_kind) {
    case "ok": {
      const r = $ail_m175;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__convert__int_to_str_from(n: bigint, acc: string, m: bigint): ScalarsResult {
  if ((m <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    if ((n <= 9n)) {
      return { $ail_kind: "ok", value: ($ailStrSlice("0123456789", n, (n + 1n)) + acc) };
    }
    else {
      const $ail_m176: ScalarsResult = std__convert__int_to_str_from($ailDivMod(n, 10n)[0], ($ailStrSlice("0123456789", $ailDivMod(n, 10n)[1], ($ailDivMod(n, 10n)[1] + 1n)) + acc), (m - 1n));
      switch ($ail_m176.$ail_kind) {
      case "ok": {
        const r = $ail_m176;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
  }
}
export function std__convert__int_to_str(value: bigint): ScalarsResult {
  if ((value < 0n)) {
    const $ail_m177: ScalarsResult = std__convert__int_to_str_from((0n - value), "", ((0n - value) + 1n));
    switch ($ail_m177.$ail_kind) {
    case "ok": {
      const r = $ail_m177;
      return { $ail_kind: "ok", value: ("-" + r.value) };
    }
    }
  }
  else {
    const $ail_m178: ScalarsResult = std__convert__int_to_str_from(value, "", (value + 1n));
    switch ($ail_m178.$ail_kind) {
    case "ok": {
      const r = $ail_m178;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__convert__str_to_int_from(orig: string, value: string, pos: bigint, acc: bigint, n: bigint): ScalarsResult {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    if (($ailStrAt(value, pos) >= 48n)) {
      if (($ailStrAt(value, pos) <= 57n)) {
        const $ail_m179: ScalarsResult = std__convert__str_to_int_from(orig, value, (pos + 1n), ((acc * 10n) + ($ailStrAt(value, pos) - 48n)), (n - 1n));
        switch ($ail_m179.$ail_kind) {
        case "convert.invalid_integer": {
          const e = $ail_m179;
          return { $ail_kind: "convert.invalid_integer", value: e.value };
        }
        case "ok": {
          const r = $ail_m179;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
      else {
        return { $ail_kind: "convert.invalid_integer", value: orig };
      }
    }
    else {
      return { $ail_kind: "convert.invalid_integer", value: orig };
    }
  }
}
export function std__convert__str_to_int(value: string): ScalarsResult {
  if (((BigInt([...value].length)) === 0n)) {
    return { $ail_kind: "convert.invalid_integer", value: value };
  }
  else {
    if (($ailStrSlice(value, 0n, 1n) === "-")) {
      if (((BigInt([...value].length)) >= 2n)) {
        const $ail_m180: ScalarsResult = std__convert__str_to_int_from(value, value, 1n, 0n, ((BigInt([...value].length)) - 1n));
        switch ($ail_m180.$ail_kind) {
        case "convert.invalid_integer": {
          const e = $ail_m180;
          return { $ail_kind: "convert.invalid_integer", value: e.value };
        }
        case "ok": {
          const r = $ail_m180;
          return { $ail_kind: "ok", value: (0n - r.value) };
        }
        }
      }
      else {
        return { $ail_kind: "convert.invalid_integer", value: value };
      }
    }
    else {
      const $ail_m181: ScalarsResult = std__convert__str_to_int_from(value, value, 0n, 0n, (BigInt([...value].length)));
      switch ($ail_m181.$ail_kind) {
      case "convert.invalid_integer": {
        const e = $ail_m181;
        return { $ail_kind: "convert.invalid_integer", value: e.value };
      }
      case "ok": {
        const r = $ail_m181;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
  }
}
export function std__convert__str_to_dec_from(orig: string, value: string, pos: bigint, coeff: bigint, scale: bigint, dot: bigint, n: bigint): ScalarsResult {
  if ((n <= 0n)) {
    if ((dot === 0n)) {
      const $ail_m182: ScalarsResult = std__convert__int_to_dec(coeff);
      switch ($ail_m182.$ail_kind) {
      case "ok": {
        const whole = $ail_m182;
        return { $ail_kind: "ok", value: whole.value };
      }
      }
    }
    else {
      if ((scale === 0n)) {
        return { $ail_kind: "convert.invalid_decimal", value: orig };
      }
      else {
        const $ail_m183: ScalarsResult = std__convert__int_to_dec(coeff);
        switch ($ail_m183.$ail_kind) {
        case "ok": {
          const whole = $ail_m183;
          const $ail_m184: ScalarsResult = std__dec__scale_by_power_of_ten(whole.value, (0n - scale));
          switch ($ail_m184.$ail_kind) {
          case "ok": {
            const moved = $ail_m184;
            return { $ail_kind: "ok", value: moved.value };
          }
          }
        }
        }
      }
    }
  }
  else {
    if (($ailStrSlice(value, pos, (pos + 1n)) === ".")) {
      if ((dot === 0n)) {
        const $ail_m185: ScalarsResult = std__convert__str_to_dec_from(orig, value, (pos + 1n), coeff, scale, 1n, (n - 1n));
        switch ($ail_m185.$ail_kind) {
        case "convert.invalid_decimal": {
          const e = $ail_m185;
          return { $ail_kind: "convert.invalid_decimal", value: e.value };
        }
        case "ok": {
          const r = $ail_m185;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
      else {
        return { $ail_kind: "convert.invalid_decimal", value: orig };
      }
    }
    else {
      if (($ailStrAt(value, pos) >= 48n)) {
        if (($ailStrAt(value, pos) <= 57n)) {
          if ((dot === 0n)) {
            const $ail_m186: ScalarsResult = std__convert__str_to_dec_from(orig, value, (pos + 1n), ((coeff * 10n) + ($ailStrAt(value, pos) - 48n)), scale, dot, (n - 1n));
            switch ($ail_m186.$ail_kind) {
            case "convert.invalid_decimal": {
              const e = $ail_m186;
              return { $ail_kind: "convert.invalid_decimal", value: e.value };
            }
            case "ok": {
              const r = $ail_m186;
              return { $ail_kind: "ok", value: r.value };
            }
            }
          }
          else {
            const $ail_m187: ScalarsResult = std__convert__str_to_dec_from(orig, value, (pos + 1n), ((coeff * 10n) + ($ailStrAt(value, pos) - 48n)), (scale + 1n), dot, (n - 1n));
            switch ($ail_m187.$ail_kind) {
            case "convert.invalid_decimal": {
              const e = $ail_m187;
              return { $ail_kind: "convert.invalid_decimal", value: e.value };
            }
            case "ok": {
              const r = $ail_m187;
              return { $ail_kind: "ok", value: r.value };
            }
            }
          }
        }
        else {
          return { $ail_kind: "convert.invalid_decimal", value: orig };
        }
      }
      else {
        return { $ail_kind: "convert.invalid_decimal", value: orig };
      }
    }
  }
}
export function std__convert__str_to_dec(value: string): ScalarsResult {
  if (((BigInt([...value].length)) === 0n)) {
    return { $ail_kind: "convert.invalid_decimal", value: value };
  }
  else {
    if (($ailStrSlice(value, 0n, 1n) === "-")) {
      if (((BigInt([...value].length)) >= 2n)) {
        const $ail_m188: ScalarsResult = std__convert__str_to_dec_from(value, value, 1n, 0n, 0n, 0n, ((BigInt([...value].length)) - 1n));
        switch ($ail_m188.$ail_kind) {
        case "convert.invalid_decimal": {
          const e = $ail_m188;
          return { $ail_kind: "convert.invalid_decimal", value: e.value };
        }
        case "ok": {
          const r = $ail_m188;
          return { $ail_kind: "ok", value: $ailDecSub("0.0", r.value) };
        }
        }
      }
      else {
        return { $ail_kind: "convert.invalid_decimal", value: value };
      }
    }
    else {
      const $ail_m189: ScalarsResult = std__convert__str_to_dec_from(value, value, 0n, 0n, 0n, 0n, (BigInt([...value].length)));
      switch ($ail_m189.$ail_kind) {
      case "convert.invalid_decimal": {
        const e = $ail_m189;
        return { $ail_kind: "convert.invalid_decimal", value: e.value };
      }
      case "ok": {
        const r = $ail_m189;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
  }
}
export function std__convert__int_to_dec(value: bigint): ScalarsResult {
  if ((value >= 0n)) {
    const $ail_m190: ScalarsResult = std__convert__int_to_dec_from(value, "0.0", "1.0");
    switch ($ail_m190.$ail_kind) {
    case "ok": {
      const r = $ail_m190;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    const $ail_m191: ScalarsResult = std__convert__int_to_dec_from((0n - value), "0.0", "-1.0");
    switch ($ail_m191.$ail_kind) {
    case "ok": {
      const r = $ail_m191;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__backoff__delay(attempt: bigint, base_ms: bigint, cap_ms: bigint): ScalarsResult {
  const $ail_m192 = (base_ms >= 0n);
  const $ail_m193 = (cap_ms >= 0n);
  const $ail_m194 = (attempt >= 0n);
  const $ail_m195 = (attempt <= 30n);
  if ($ail_m192 && $ail_m193 && $ail_m194 && $ail_m195) {
    const $ail_m196: ScalarsResult = std__int__pow_from(2n, attempt);
    switch ($ail_m196.$ail_kind) {
    case "ok": {
      const p = $ail_m196;
      if (((base_ms * p.value) <= cap_ms)) {
        return { $ail_kind: "ok", value: (base_ms * p.value) };
      }
      else {
        return { $ail_kind: "ok", value: cap_ms };
      }
    }
    }
  }
  else if ($ail_m192 && $ail_m193 && $ail_m194 && !($ail_m195)) {
    return { $ail_kind: "ok", value: cap_ms };
  }
  else if ($ail_m192 && $ail_m193 && !($ail_m194)) {
    return { $ail_kind: "math.negative_input", value: attempt };
  }
  else if ($ail_m192 && !($ail_m193)) {
    return { $ail_kind: "math.negative_input", value: cap_ms };
  }
  else {
    return { $ail_kind: "math.negative_input", value: base_ms };
  }
}
