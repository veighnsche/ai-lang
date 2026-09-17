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
// Exact-decimal runtime (a10): canonical-digit strings, BigInt math.
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
  const cp = cps[k].codePointAt(0);
  if (cp === undefined) throw new Error("str index out of range");
  return BigInt(cp);
}
function $ailStrSlice(s: string, a: bigint, b: bigint): string {
  const cps = [...s];
  const toIdx = (x: bigint): number => { if (x < 0n || x > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("str slice out of range"); return Number(x); };
  const lo = toIdx(a), hi = toIdx(b);
  if (lo > hi || hi > cps.length) throw new Error("str slice out of range");
  return cps.slice(lo, hi).join("");
}
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
export function std__bool__not(value: boolean): { $ail_kind: "ok"; value: boolean } {
  if (value) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    return { $ail_kind: "ok", value: true };
  }
}
export function std__bool__and(left: boolean, right: boolean): { $ail_kind: "ok"; value: boolean } {
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
export function std__bool__or(left: boolean, right: boolean): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1 = left;
  const $ail_m2 = right;
  if ($ail_m1) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__bool__xor(left: boolean, right: boolean): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1 = left;
  const $ail_m2 = right;
  if ($ail_m1 && $ail_m2) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "ok", value: true };
  }
  else if (!($ail_m1) && $ail_m2) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__bool__implies(left: boolean, right: boolean): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1 = left;
  const $ail_m2 = right;
  if ($ail_m1 && $ail_m2) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    return { $ail_kind: "ok", value: true };
  }
}
export function std__bool__equivalent(left: boolean, right: boolean): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1 = left;
  const $ail_m2 = right;
  if ($ail_m1 && $ail_m2) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "ok", value: false };
  }
  else if (!($ail_m1) && $ail_m2) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    return { $ail_kind: "ok", value: true };
  }
}
export function std__compare__int(left: bigint, right: bigint): { $ail_kind: "ok"; value: bigint } {
  const $ail_m1 = (left === right);
  const $ail_m2 = (left >= right);
  if ($ail_m1) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__compare__dec(left: string, right: string): { $ail_kind: "ok"; value: bigint } {
  const $ail_m1 = (left === right);
  const $ail_m2 = $ailDecGe(left, right);
  if ($ail_m1) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__compare__str(left: string, right: string): { $ail_kind: "ok"; value: bigint } {
  const $ail_m1 = (left === right);
  const $ail_m2 = $ailStrGe(left, right);
  if ($ail_m1) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__compare__bool(left: boolean, right: boolean): { $ail_kind: "ok"; value: bigint } {
  const $ail_m1 = (left === right);
  const $ail_m2 = left;
  if ($ail_m1) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__select__int(condition: boolean, when_true: bigint, when_false: bigint): { $ail_kind: "ok"; value: bigint } {
  if (condition) {
    return { $ail_kind: "ok", value: when_true };
  }
  else {
    return { $ail_kind: "ok", value: when_false };
  }
}
export function std__select__bool(condition: boolean, when_true: boolean, when_false: boolean): { $ail_kind: "ok"; value: boolean } {
  if (condition) {
    return { $ail_kind: "ok", value: when_true };
  }
  else {
    return { $ail_kind: "ok", value: when_false };
  }
}
export function std__select__str(condition: boolean, when_true: string, when_false: string): { $ail_kind: "ok"; value: string } {
  if (condition) {
    return { $ail_kind: "ok", value: when_true };
  }
  else {
    return { $ail_kind: "ok", value: when_false };
  }
}
export function std__select__dec(condition: boolean, when_true: string, when_false: string): { $ail_kind: "ok"; value: string } {
  if (condition) {
    return { $ail_kind: "ok", value: when_true };
  }
  else {
    return { $ail_kind: "ok", value: when_false };
  }
}
export function std__int__is_zero(value: bigint): { $ail_kind: "ok"; value: boolean } {
  if ((value === 0n)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__int__is_positive(value: bigint): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1 = (value === 0n);
  const $ail_m2 = (value >= 0n);
  if ($ail_m1) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__int__is_negative(value: bigint): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1 = (value === 0n);
  const $ail_m2 = (value <= 0n);
  if ($ail_m1) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__int__in_closed_range(value: bigint, lower: bigint, upper: bigint): { $ail_kind: "ok"; value: boolean } | { $ail_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } {
  const $ail_m1 = (lower <= upper);
  const $ail_m2 = (value >= lower);
  const $ail_m3 = (value <= upper);
  if ($ail_m1 && $ail_m2 && $ail_m3) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m1 && $ail_m2 && !($ail_m3)) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__in_open_range(value: bigint, lower: bigint, upper: bigint): { $ail_kind: "ok"; value: boolean } | { $ail_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } {
  const $ail_m1 = (lower <= upper);
  const $ail_m2 = (value <= lower);
  const $ail_m3 = (value >= upper);
  if ($ail_m1 && $ail_m2) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m1 && !($ail_m2) && $ail_m3) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m1 && !($ail_m2) && !($ail_m3)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__is_zero(value: string): { $ail_kind: "ok"; value: boolean } {
  if ((value === "0.0")) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__dec__is_positive(value: string): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1 = (value === "0.0");
  const $ail_m2 = $ailDecGe(value, "0.0");
  if ($ail_m1) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__dec__is_negative(value: string): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1 = (value === "0.0");
  const $ail_m2 = $ailDecLe(value, "0.0");
  if ($ail_m1) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__dec__in_closed_range(value: string, lower: string, upper: string): { $ail_kind: "ok"; value: boolean } | { $ail_kind: "math.dec_invalid_bounds"; lower: string; upper: string } {
  const $ail_m1 = $ailDecLe(lower, upper);
  const $ail_m2 = $ailDecGe(value, lower);
  const $ail_m3 = $ailDecLe(value, upper);
  if ($ail_m1 && $ail_m2 && $ail_m3) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m1 && $ail_m2 && !($ail_m3)) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    return { $ail_kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__in_open_range(value: string, lower: string, upper: string): { $ail_kind: "ok"; value: boolean } | { $ail_kind: "math.dec_invalid_bounds"; lower: string; upper: string } {
  const $ail_m1 = $ailDecLe(lower, upper);
  const $ail_m2 = $ailDecLe(value, lower);
  const $ail_m3 = $ailDecGe(value, upper);
  if ($ail_m1 && $ail_m2) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m1 && !($ail_m2) && $ail_m3) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m1 && !($ail_m2) && !($ail_m3)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__abs(value: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((value >= 0n)) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    return { $ail_kind: "ok", value: (0n - value) };
  }
}
export function std__int__negate(value: bigint): { $ail_kind: "ok"; value: bigint } {
  return { $ail_kind: "ok", value: (0n - value) };
}
export function std__int__sign(value: bigint): { $ail_kind: "ok"; value: bigint } {
  const $ail_m1 = (value === 0n);
  const $ail_m2 = (value >= 0n);
  if ($ail_m1) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__int__min(left: bigint, right: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((left <= right)) {
    return { $ail_kind: "ok", value: left };
  }
  else {
    return { $ail_kind: "ok", value: right };
  }
}
export function std__int__max(left: bigint, right: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((left >= right)) {
    return { $ail_kind: "ok", value: left };
  }
  else {
    return { $ail_kind: "ok", value: right };
  }
}
export function std__int__clamp(value: bigint, lower: bigint, upper: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } {
  const $ail_m1 = (lower <= upper);
  const $ail_m2 = (value >= lower);
  const $ail_m3 = (value <= upper);
  if ($ail_m1 && $ail_m2 && $ail_m3) {
    return { $ail_kind: "ok", value: value };
  }
  else if ($ail_m1 && $ail_m2 && !($ail_m3)) {
    return { $ail_kind: "ok", value: upper };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "ok", value: lower };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__distance(left: bigint, right: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((left >= right)) {
    return { $ail_kind: "ok", value: (left - right) };
  }
  else {
    return { $ail_kind: "ok", value: (right - left) };
  }
}
export function std__int__square(value: bigint): { $ail_kind: "ok"; value: bigint } {
  return { $ail_kind: "ok", value: (value * value) };
}
export function std__int__pow_from(base: bigint, exponent: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((exponent <= 0n)) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__pow_from(base, (exponent - 1n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const r = $ail_m1;
      return { $ail_kind: "ok", value: (base * r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__int__pow(base: bigint, exponent: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "math.negative_exponent"; exponent: bigint } {
  if ((exponent >= 0n)) {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__pow_from(base, exponent);
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
  else {
    return { $ail_kind: "math.negative_exponent", exponent: exponent };
  }
}
export function std__int__add_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } | { $ail_kind: "math.out_of_range"; value: bigint; lower: bigint; upper: bigint } {
  const $ail_m1 = (lower <= upper);
  const $ail_m2 = ((left + right) >= lower);
  const $ail_m3 = ((left + right) <= upper);
  if ($ail_m1 && $ail_m2 && $ail_m3) {
    return { $ail_kind: "ok", value: (left + right) };
  }
  else if ($ail_m1 && $ail_m2 && !($ail_m3)) {
    return { $ail_kind: "math.out_of_range", value: (left + right), lower: lower, upper: upper };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "math.out_of_range", value: (left + right), lower: lower, upper: upper };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__subtract_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } | { $ail_kind: "math.out_of_range"; value: bigint; lower: bigint; upper: bigint } {
  const $ail_m1 = (lower <= upper);
  const $ail_m2 = ((left - right) >= lower);
  const $ail_m3 = ((left - right) <= upper);
  if ($ail_m1 && $ail_m2 && $ail_m3) {
    return { $ail_kind: "ok", value: (left - right) };
  }
  else if ($ail_m1 && $ail_m2 && !($ail_m3)) {
    return { $ail_kind: "math.out_of_range", value: (left - right), lower: lower, upper: upper };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "math.out_of_range", value: (left - right), lower: lower, upper: upper };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__multiply_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } | { $ail_kind: "math.out_of_range"; value: bigint; lower: bigint; upper: bigint } {
  const $ail_m1 = (lower <= upper);
  const $ail_m2 = ((left * right) >= lower);
  const $ail_m3 = ((left * right) <= upper);
  if ($ail_m1 && $ail_m2 && $ail_m3) {
    return { $ail_kind: "ok", value: (left * right) };
  }
  else if ($ail_m1 && $ail_m2 && !($ail_m3)) {
    return { $ail_kind: "math.out_of_range", value: (left * right), lower: lower, upper: upper };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "math.out_of_range", value: (left * right), lower: lower, upper: upper };
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__factorial_from(value: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((value <= 0n)) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__factorial_from((value - 1n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const r = $ail_m1;
      return { $ail_kind: "ok", value: (value * r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__int__factorial(value: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "math.negative_input"; value: bigint } {
  if ((value >= 0n)) {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__factorial_from(value);
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
  else {
    return { $ail_kind: "math.negative_input", value: value };
  }
}
export function std__int__sum_from(upper: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((upper <= 0n)) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__sum_from((upper - 1n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const r = $ail_m1;
      return { $ail_kind: "ok", value: (upper + r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__int__sum_to(upper: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "math.negative_input"; value: bigint } {
  if ((upper >= 0n)) {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__sum_from(upper);
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
  else {
    return { $ail_kind: "math.negative_input", value: upper };
  }
}
export function std__int__binomial_from(n: bigint, k: bigint, i: bigint, num: bigint, den: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((i <= 0n)) {
    return { $ail_kind: "ok", value: $ailDivMod(num, den)[0] };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__binomial_from(n, k, (i - 1n), (num * ((n - k) + i)), (den * i));
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
export function std__int__binomial(n: bigint, k: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "math.invalid_count"; n: bigint; k: bigint } {
  const $ail_m1 = (k >= 0n);
  const $ail_m2 = (k <= n);
  if ($ail_m1 && $ail_m2) {
    const $ail_m3: { $ail_kind: "ok"; value: bigint } = std__int__binomial_from(n, k, k, 1n, 1n);
    switch ($ail_m3.$ail_kind) {
    case "ok": {
      const r = $ail_m3;
      return { $ail_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "math.invalid_count", n: n, k: k };
  }
  else {
    return { $ail_kind: "math.invalid_count", n: n, k: k };
  }
}
export function std__int__is_prime_from(value: bigint, e: bigint): { $ail_kind: "ok"; value: boolean } {
  if ((e <= 0n)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    if (($ailDivMod(value, (e + 1n))[1] === 0n)) {
      return { $ail_kind: "ok", value: false };
    }
    else {
      const $ail_m1: { $ail_kind: "ok"; value: boolean } = std__int__is_prime_from(value, (e - 1n));
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
}
export function std__int__is_prime(value: bigint): { $ail_kind: "ok"; value: boolean } {
  if ((value >= 2n)) {
    const $ail_m1: { $ail_kind: "ok"; value: boolean } = std__int__is_prime_from(value, (value - 2n));
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
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__int__next_power_of_two_pow2(e: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((e <= 0n)) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__next_power_of_two_pow2((e - 1n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const r = $ail_m1;
      return { $ail_kind: "ok", value: (2n * r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__int__next_power_of_two_from(value: bigint, e: bigint, best: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((e <= 0n)) {
    if ((1n >= value)) {
      return { $ail_kind: "ok", value: 1n };
    }
    else {
      return { $ail_kind: "ok", value: best };
    }
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__next_power_of_two_pow2(e);
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const t = $ail_m1;
      if ((t.value >= value)) {
        const $ail_m2: { $ail_kind: "ok"; value: bigint } = std__int__next_power_of_two_from(value, (e - 1n), t.value);
        switch ($ail_m2.$ail_kind) {
        case "ok": {
          const r = $ail_m2;
          return { $ail_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $ail_m3: { $ail_kind: "ok"; value: bigint } = std__int__next_power_of_two_from(value, (e - 1n), best);
        switch ($ail_m3.$ail_kind) {
        case "ok": {
          const r = $ail_m3;
          return { $ail_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__int__next_power_of_two(value: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "math.nonpositive_input"; value: bigint } {
  if ((value >= 1n)) {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__next_power_of_two_from(value, value, 1n);
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
  else {
    return { $ail_kind: "math.nonpositive_input", value: value };
  }
}
export function std__int__sqrt_floor_search(value: bigint, lo: bigint, hi: bigint): { $ail_kind: "ok"; value: bigint } {
  if (((hi - lo) <= 1n)) {
    return { $ail_kind: "ok", value: lo };
  }
  else {
    if ((($ailDivMod((lo + hi), 2n)[0] * $ailDivMod((lo + hi), 2n)[0]) <= value)) {
      const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__sqrt_floor_search(value, $ailDivMod((lo + hi), 2n)[0], hi);
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
    else {
      const $ail_m2: { $ail_kind: "ok"; value: bigint } = std__int__sqrt_floor_search(value, lo, $ailDivMod((lo + hi), 2n)[0]);
      switch ($ail_m2.$ail_kind) {
      case "ok": {
        const r = $ail_m2;
        return { $ail_kind: "ok", value: r.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
}
export function std__int__sqrt_floor(value: bigint): { $ail_kind: "ok"; remainder: bigint; root: bigint } | { $ail_kind: "math.negative_input"; value: bigint } {
  if ((value >= 0n)) {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__sqrt_floor_search(value, 0n, (value + 1n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const r = $ail_m1;
      return { $ail_kind: "ok", root: r.value, remainder: (value - (r.value * r.value)) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_input", value: value };
  }
}
export function std__int__gcd_euclid(a: bigint, b: bigint): { $ail_kind: "ok"; value: bigint } {
  if ((b <= 0n)) {
    return { $ail_kind: "ok", value: a };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__gcd_euclid(b, $ailDivMod(a, b)[1]);
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
export function std__int__gcd(left: bigint, right: bigint): { $ail_kind: "ok"; value: bigint } {
  const $ail_m1: { $ail_kind: "ok"; value: bigint } = std__int__abs(left);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const a = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; value: bigint } = std__int__abs(right);
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const b = $ail_m2;
      const $ail_m3: { $ail_kind: "ok"; value: bigint } = std__int__gcd_euclid(a.value, b.value);
      switch ($ail_m3.$ail_kind) {
      case "ok": {
        const r = $ail_m3;
        return { $ail_kind: "ok", value: r.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__int__lcm(left: bigint, right: bigint): { $ail_kind: "ok"; value: bigint } {
  const $ail_m1 = (left === 0n);
  const $ail_m2 = (right === 0n);
  if ($ail_m1) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    const $ail_m3: { $ail_kind: "ok"; value: bigint } = std__int__gcd(left, right);
    switch ($ail_m3.$ail_kind) {
    case "ok": {
      const g = $ail_m3;
      const $ail_m4: { $ail_kind: "ok"; value: bigint } = std__int__abs(left);
      switch ($ail_m4.$ail_kind) {
      case "ok": {
        const a = $ail_m4;
        const $ail_m5: { $ail_kind: "ok"; value: bigint } = std__int__abs(right);
        switch ($ail_m5.$ail_kind) {
        case "ok": {
          const b = $ail_m5;
          return { $ail_kind: "ok", value: $ailDivMod((a.value * b.value), g.value)[0] };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__dec__abs(value: string): { $ail_kind: "ok"; value: string } {
  if ($ailDecGe(value, "0.0")) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    return { $ail_kind: "ok", value: $ailDecSub("0.0", value) };
  }
}
export function std__dec__negate(value: string): { $ail_kind: "ok"; value: string } {
  return { $ail_kind: "ok", value: $ailDecSub("0.0", value) };
}
export function std__dec__sign(value: string): { $ail_kind: "ok"; value: bigint } {
  const $ail_m1 = (value === "0.0");
  const $ail_m2 = $ailDecGe(value, "0.0");
  if ($ail_m1) {
    return { $ail_kind: "ok", value: 0n };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: -1n };
  }
}
export function std__dec__min(left: string, right: string): { $ail_kind: "ok"; value: string } {
  if ($ailDecLe(left, right)) {
    return { $ail_kind: "ok", value: left };
  }
  else {
    return { $ail_kind: "ok", value: right };
  }
}
export function std__dec__max(left: string, right: string): { $ail_kind: "ok"; value: string } {
  if ($ailDecGe(left, right)) {
    return { $ail_kind: "ok", value: left };
  }
  else {
    return { $ail_kind: "ok", value: right };
  }
}
export function std__dec__clamp(value: string, lower: string, upper: string): { $ail_kind: "ok"; value: string } | { $ail_kind: "math.dec_invalid_bounds"; lower: string; upper: string } {
  const $ail_m1 = $ailDecLe(lower, upper);
  const $ail_m2 = $ailDecGe(value, lower);
  const $ail_m3 = $ailDecLe(value, upper);
  if ($ail_m1 && $ail_m2 && $ail_m3) {
    return { $ail_kind: "ok", value: value };
  }
  else if ($ail_m1 && $ail_m2 && !($ail_m3)) {
    return { $ail_kind: "ok", value: upper };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "ok", value: lower };
  }
  else {
    return { $ail_kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__distance(left: string, right: string): { $ail_kind: "ok"; value: string } {
  if ($ailDecGe(left, right)) {
    return { $ail_kind: "ok", value: $ailDecSub(left, right) };
  }
  else {
    return { $ail_kind: "ok", value: $ailDecSub(right, left) };
  }
}
export function std__dec__square(value: string): { $ail_kind: "ok"; value: string } {
  return { $ail_kind: "ok", value: $ailDecMul(value, value) };
}
export function std__dec__pow_from(base: string, exponent: bigint): { $ail_kind: "ok"; value: string } {
  if ((exponent <= 0n)) {
    return { $ail_kind: "ok", value: "1.0" };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: string } = std__dec__pow_from(base, (exponent - 1n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const r = $ail_m1;
      return { $ail_kind: "ok", value: $ailDecMul(base, r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__dec__pow(base: string, exponent: bigint): { $ail_kind: "ok"; value: string } | { $ail_kind: "math.negative_exponent"; exponent: bigint } {
  if ((exponent >= 0n)) {
    const $ail_m1: { $ail_kind: "ok"; value: string } = std__dec__pow_from(base, exponent);
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
  else {
    return { $ail_kind: "math.negative_exponent", exponent: exponent };
  }
}
export function std__dec__lerp(start: string, end: string, fraction: string): { $ail_kind: "ok"; value: string } | { $ail_kind: "math.dec_out_of_range"; value: string; lower: string; upper: string } {
  const $ail_m1 = $ailDecGe(fraction, "0.0");
  const $ail_m2 = $ailDecLe(fraction, "1.0");
  if ($ail_m1 && $ail_m2) {
    return { $ail_kind: "ok", value: $ailDecAdd(start, $ailDecMul($ailDecSub(end, start), fraction)) };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "math.dec_out_of_range", value: fraction, lower: "0.0", upper: "1.0" };
  }
  else {
    return { $ail_kind: "math.dec_out_of_range", value: fraction, lower: "0.0", upper: "1.0" };
  }
}
export function std__dec__scale_by_power_of_ten_from(value: string, e: bigint, factor: string): { $ail_kind: "ok"; value: string } {
  if ((e <= 0n)) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten_from($ailDecMul(value, factor), (e - 1n), factor);
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
export function std__dec__scale_by_power_of_ten(value: string, exponent: bigint): { $ail_kind: "ok"; value: string } {
  if ((exponent >= 0n)) {
    const $ail_m1: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten_from(value, exponent, "10.0");
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
  else {
    const $ail_m2: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten_from(value, (0n - exponent), "0.1");
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const r = $ail_m2;
      return { $ail_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__dec__from_parts(coefficient: bigint, scale: bigint): { $ail_kind: "ok"; value: string } | { $ail_kind: "math.negative_scale"; scale: bigint } {
  if ((scale >= 0n)) {
    const $ail_m1: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(coefficient);
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const whole = $ail_m1;
      const $ail_m2: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(whole.value, (0n - scale));
      switch ($ail_m2.$ail_kind) {
      case "ok": {
        const moved = $ail_m2;
        return { $ail_kind: "ok", value: moved.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_scale", scale: scale };
  }
}
export function std__dec__parts(value: string): { $ail_kind: "ok"; coefficient: bigint; scale: bigint } {
  const $ail_m1: { $ail_kind: "ok", coefficient: bigint, scale: bigint } = { $ail_kind: "ok", ...$ailDecParts(value) };
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const p = $ail_m1;
    return { $ail_kind: "ok", coefficient: p.coefficient, scale: p.scale };
  }
  }
}
export function std__dec__strip_from(remaining: bigint, current: bigint, exact: boolean): { $ail_kind: "ok"; exact: boolean; stripped: bigint } {
  if ((remaining <= 0n)) {
    return { $ail_kind: "ok", stripped: current, exact: exact };
  }
  else {
    if (($ailDivMod(current, 10n)[1] === 0n)) {
      const $ail_m1: { $ail_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from((remaining - 1n), $ailDivMod(current, 10n)[0], exact);
      switch ($ail_m1.$ail_kind) {
      case "ok": {
        const r = $ail_m1;
        return { $ail_kind: "ok", stripped: r.stripped, exact: r.exact };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      const $ail_m2: { $ail_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from((remaining - 1n), $ailDivMod(current, 10n)[0], false);
      switch ($ail_m2.$ail_kind) {
      case "ok": {
        const r = $ail_m2;
        return { $ail_kind: "ok", stripped: r.stripped, exact: r.exact };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
}
export function std__dec__truncate(value: string): { $ail_kind: "ok"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const p = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const s = $ail_m2;
      if (s.exact) {
        const $ail_m3: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(s.stripped);
        switch ($ail_m3.$ail_kind) {
        case "ok": {
          const whole = $ail_m3;
          return { $ail_kind: "ok", value: whole.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        if ((s.stripped >= 0n)) {
          const $ail_m4: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(s.stripped);
          switch ($ail_m4.$ail_kind) {
          case "ok": {
            const whole = $ail_m4;
            return { $ail_kind: "ok", value: whole.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          const $ail_m5: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec((s.stripped + 1n));
          switch ($ail_m5.$ail_kind) {
          case "ok": {
            const whole = $ail_m5;
            return { $ail_kind: "ok", value: whole.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__dec__floor(value: string): { $ail_kind: "ok"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const p = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const s = $ail_m2;
      const $ail_m3: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(s.stripped);
      switch ($ail_m3.$ail_kind) {
      case "ok": {
        const whole = $ail_m3;
        return { $ail_kind: "ok", value: whole.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__dec__ceil(value: string): { $ail_kind: "ok"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const p = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const s = $ail_m2;
      if (s.exact) {
        const $ail_m3: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(s.stripped);
        switch ($ail_m3.$ail_kind) {
        case "ok": {
          const whole = $ail_m3;
          return { $ail_kind: "ok", value: whole.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $ail_m4: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec((s.stripped + 1n));
        switch ($ail_m4.$ail_kind) {
        case "ok": {
          const whole = $ail_m4;
          return { $ail_kind: "ok", value: whole.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__dec__strip257_from(fuel: bigint, current: bigint, twos: bigint, fives: bigint): { $ail_kind: "ok"; fives: bigint; rest: bigint; twos: bigint } {
  if ((fuel <= 0n)) {
    return { $ail_kind: "ok", rest: current, twos: twos, fives: fives };
  }
  else {
    if ((current === 1n)) {
      return { $ail_kind: "ok", rest: 1n, twos: twos, fives: fives };
    }
    else {
      if (($ailDivMod(current, 2n)[1] === 0n)) {
        const $ail_m1: { $ail_kind: "ok"; fives: bigint; rest: bigint; twos: bigint } = std__dec__strip257_from((fuel - 1n), $ailDivMod(current, 2n)[0], (twos + 1n), fives);
        switch ($ail_m1.$ail_kind) {
        case "ok": {
          const r = $ail_m1;
          return { $ail_kind: "ok", rest: r.rest, twos: r.twos, fives: r.fives };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        if (($ailDivMod(current, 5n)[1] === 0n)) {
          const $ail_m2: { $ail_kind: "ok"; fives: bigint; rest: bigint; twos: bigint } = std__dec__strip257_from((fuel - 1n), $ailDivMod(current, 5n)[0], twos, (fives + 1n));
          switch ($ail_m2.$ail_kind) {
          case "ok": {
            const r = $ail_m2;
            return { $ail_kind: "ok", rest: r.rest, twos: r.twos, fives: r.fives };
          }
          default: {
            throw new Error("unreachable");
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
export function std__dec__divide_exact(dividend: string, divisor: string): { $ail_kind: "ok"; value: string } | { $ail_kind: "math.dec_zero_divisor"; divisor: string } | { $ail_kind: "math.nonterminating_decimal"; dividend: string; divisor: string } {
  const $ail_m1: { $ail_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(divisor);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const d = $ail_m1;
    if ((d.coefficient === 0n)) {
      return { $ail_kind: "math.dec_zero_divisor", divisor: divisor };
    }
    else {
      const $ail_m2: { $ail_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(dividend);
      switch ($ail_m2.$ail_kind) {
      case "ok": {
        const n = $ail_m2;
        const $ail_m3: { $ail_kind: "ok"; value: bigint } = std__int__pow_from(10n, d.scale);
        switch ($ail_m3.$ail_kind) {
        case "ok": {
          const pds = $ail_m3;
          const $ail_m4: { $ail_kind: "ok"; value: bigint } = std__int__pow_from(10n, n.scale);
          switch ($ail_m4.$ail_kind) {
          case "ok": {
            const pns = $ail_m4;
            const $ail_m5: { $ail_kind: "ok"; value: bigint } = std__int__abs((n.coefficient * pns.value));
            switch ($ail_m5.$ail_kind) {
            case "ok": {
              const an = $ail_m5;
              const $ail_m6: { $ail_kind: "ok"; value: bigint } = std__int__abs((d.coefficient * pds.value));
              switch ($ail_m6.$ail_kind) {
              case "ok": {
                const ad = $ail_m6;
                const $ail_m7: { $ail_kind: "ok"; value: bigint } = std__int__gcd(an.value, ad.value);
                switch ($ail_m7.$ail_kind) {
                case "ok": {
                  const g = $ail_m7;
                  const $ail_m8: { $ail_kind: "ok"; fives: bigint; rest: bigint; twos: bigint } = std__dec__strip257_from($ailDivMod(ad.value, g.value)[0], $ailDivMod(ad.value, g.value)[0], 0n, 0n);
                  switch ($ail_m8.$ail_kind) {
                  case "ok": {
                    const f = $ail_m8;
                    if ((f.rest === 1n)) {
                      const $ail_m9: { $ail_kind: "ok"; value: bigint } = std__int__max(f.twos, f.fives);
                      switch ($ail_m9.$ail_kind) {
                      case "ok": {
                        const m = $ail_m9;
                        const $ail_m10: { $ail_kind: "ok"; value: bigint } = std__int__pow_from(10n, m.value);
                        switch ($ail_m10.$ail_kind) {
                        case "ok": {
                          const p10 = $ail_m10;
                          const $ail_m11: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(($ailDivMod(an.value, g.value)[0] * $ailDivMod(p10.value, $ailDivMod(ad.value, g.value)[0])[0]));
                          switch ($ail_m11.$ail_kind) {
                          case "ok": {
                            const qi = $ail_m11;
                            const $ail_m12: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - m.value));
                            switch ($ail_m12.$ail_kind) {
                            case "ok": {
                              const q = $ail_m12;
                              const $ail_m13 = ((n.coefficient * pns.value) < 0n);
                              const $ail_m14 = ((d.coefficient * pds.value) < 0n);
                              if ($ail_m13 && $ail_m14) {
                                return { $ail_kind: "ok", value: q.value };
                              }
                              else if ($ail_m13 && !($ail_m14)) {
                                return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0") };
                              }
                              else if (!($ail_m13) && $ail_m14) {
                                return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0") };
                              }
                              else {
                                return { $ail_kind: "ok", value: q.value };
                              }
                            }
                            default: {
                              throw new Error("unreachable");
                            }
                            }
                          }
                          default: {
                            throw new Error("unreachable");
                          }
                          }
                        }
                        default: {
                          throw new Error("unreachable");
                        }
                        }
                      }
                      default: {
                        throw new Error("unreachable");
                      }
                      }
                    }
                    else {
                      return { $ail_kind: "math.nonterminating_decimal", dividend: dividend, divisor: divisor };
                    }
                  }
                  default: {
                    throw new Error("unreachable");
                  }
                  }
                }
                default: {
                  throw new Error("unreachable");
                }
                }
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__dec__round_half_even(value: string, scale: bigint): { $ail_kind: "ok"; discarded: string; value: string } | { $ail_kind: "math.negative_scale"; scale: bigint } {
  if ((scale >= 0n)) {
    const $ail_m1: { $ail_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const p = $ail_m1;
      if ((p.scale <= scale)) {
        return { $ail_kind: "ok", value: value, discarded: "0.0" };
      }
      else {
        const $ail_m2: { $ail_kind: "ok"; value: bigint } = std__int__pow_from(10n, (p.scale - scale));
        switch ($ail_m2.$ail_kind) {
        case "ok": {
          const pd = $ail_m2;
          const $ail_m3: { $ail_kind: "ok"; value: bigint } = std__int__abs(p.coefficient);
          switch ($ail_m3.$ail_kind) {
          case "ok": {
            const a = $ail_m3;
            const $ail_m4 = (($ailDivMod(a.value, pd.value)[1] * 2n) === pd.value);
            const $ail_m5 = ($ailDivMod($ailDivMod(a.value, pd.value)[0], 2n)[1] === 0n);
            const $ail_m6 = (($ailDivMod(a.value, pd.value)[1] * 2n) > pd.value);
            if ($ail_m4 && $ail_m5) {
              const $ail_m7: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec($ailDivMod(a.value, pd.value)[0]);
              switch ($ail_m7.$ail_kind) {
              case "ok": {
                const qi = $ail_m7;
                const $ail_m8: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($ail_m8.$ail_kind) {
                case "ok": {
                  const q = $ail_m8;
                  if ($ailDecLt(value, "0.0")) {
                    return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), discarded: $ailDecSub(value, $ailDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $ail_kind: "ok", value: q.value, discarded: $ailDecSub(value, q.value) };
                  }
                }
                default: {
                  throw new Error("unreachable");
                }
                }
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            else if ($ail_m4 && !($ail_m5)) {
              const $ail_m9: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(($ailDivMod(a.value, pd.value)[0] + 1n));
              switch ($ail_m9.$ail_kind) {
              case "ok": {
                const qi = $ail_m9;
                const $ail_m10: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($ail_m10.$ail_kind) {
                case "ok": {
                  const q = $ail_m10;
                  if ($ailDecLt(value, "0.0")) {
                    return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), discarded: $ailDecSub(value, $ailDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $ail_kind: "ok", value: q.value, discarded: $ailDecSub(value, q.value) };
                  }
                }
                default: {
                  throw new Error("unreachable");
                }
                }
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            else if (!($ail_m4) && $ail_m6) {
              const $ail_m11: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(($ailDivMod(a.value, pd.value)[0] + 1n));
              switch ($ail_m11.$ail_kind) {
              case "ok": {
                const qi = $ail_m11;
                const $ail_m12: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($ail_m12.$ail_kind) {
                case "ok": {
                  const q = $ail_m12;
                  if ($ailDecLt(value, "0.0")) {
                    return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), discarded: $ailDecSub(value, $ailDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $ail_kind: "ok", value: q.value, discarded: $ailDecSub(value, q.value) };
                  }
                }
                default: {
                  throw new Error("unreachable");
                }
                }
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            else {
              const $ail_m13: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec($ailDivMod(a.value, pd.value)[0]);
              switch ($ail_m13.$ail_kind) {
              case "ok": {
                const qi = $ail_m13;
                const $ail_m14: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($ail_m14.$ail_kind) {
                case "ok": {
                  const q = $ail_m14;
                  if ($ailDecLt(value, "0.0")) {
                    return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), discarded: $ailDecSub(value, $ailDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $ail_kind: "ok", value: q.value, discarded: $ailDecSub(value, q.value) };
                  }
                }
                default: {
                  throw new Error("unreachable");
                }
                }
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_scale", scale: scale };
  }
}
export function std__dec__divide_round_half_even_result(qm: bigint, scale: bigint, num: bigint, den: bigint, dividend: string, divisor: string): { $ail_kind: "ok"; remainder: string; value: string } {
  const $ail_m1: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(qm);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const qi = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const q = $ail_m2;
      const $ail_m3 = (num < 0n);
      const $ail_m4 = (den < 0n);
      if ($ail_m3 && $ail_m4) {
        return { $ail_kind: "ok", value: q.value, remainder: $ailDecSub(dividend, $ailDecMul(divisor, q.value)) };
      }
      else if ($ail_m3 && !($ail_m4)) {
        return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), remainder: $ailDecSub(dividend, $ailDecMul(divisor, $ailDecMul(q.value, "-1.0"))) };
      }
      else if (!($ail_m3) && $ail_m4) {
        return { $ail_kind: "ok", value: $ailDecMul(q.value, "-1.0"), remainder: $ailDecSub(dividend, $ailDecMul(divisor, $ailDecMul(q.value, "-1.0"))) };
      }
      else {
        return { $ail_kind: "ok", value: q.value, remainder: $ailDecSub(dividend, $ailDecMul(divisor, q.value)) };
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__dec__divide_round_half_even(dividend: string, divisor: string, scale: bigint): { $ail_kind: "ok"; remainder: string; value: string } | { $ail_kind: "math.dec_zero_divisor"; divisor: string } | { $ail_kind: "math.negative_scale"; scale: bigint } {
  const $ail_m1: { $ail_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(divisor);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const d = $ail_m1;
    if ((d.coefficient === 0n)) {
      return { $ail_kind: "math.dec_zero_divisor", divisor: divisor };
    }
    else {
      if ((scale >= 0n)) {
        const $ail_m2: { $ail_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(dividend);
        switch ($ail_m2.$ail_kind) {
        case "ok": {
          const n = $ail_m2;
          const $ail_m3: { $ail_kind: "ok"; value: bigint } = std__int__pow_from(10n, d.scale);
          switch ($ail_m3.$ail_kind) {
          case "ok": {
            const pds = $ail_m3;
            const $ail_m4: { $ail_kind: "ok"; value: bigint } = std__int__pow_from(10n, n.scale);
            switch ($ail_m4.$ail_kind) {
            case "ok": {
              const pns = $ail_m4;
              const $ail_m5: { $ail_kind: "ok"; value: bigint } = std__int__abs((n.coefficient * pns.value));
              switch ($ail_m5.$ail_kind) {
              case "ok": {
                const an = $ail_m5;
                const $ail_m6: { $ail_kind: "ok"; value: bigint } = std__int__abs((d.coefficient * pds.value));
                switch ($ail_m6.$ail_kind) {
                case "ok": {
                  const ad = $ail_m6;
                  const $ail_m7: { $ail_kind: "ok"; value: bigint } = std__int__pow_from(10n, scale);
                  switch ($ail_m7.$ail_kind) {
                  case "ok": {
                    const psc = $ail_m7;
                    const $ail_m8 = (($ailDivMod((an.value * psc.value), ad.value)[1] * 2n) === ad.value);
                    const $ail_m9 = ($ailDivMod($ailDivMod((an.value * psc.value), ad.value)[0], 2n)[1] === 0n);
                    const $ail_m10 = (($ailDivMod((an.value * psc.value), ad.value)[1] * 2n) > ad.value);
                    if ($ail_m8 && $ail_m9) {
                      const $ail_m11: { $ail_kind: "ok"; remainder: string; value: string } = std__dec__divide_round_half_even_result($ailDivMod((an.value * psc.value), ad.value)[0], scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($ail_m11.$ail_kind) {
                      case "ok": {
                        const r = $ail_m11;
                        return { $ail_kind: "ok", value: r.value, remainder: r.remainder };
                      }
                      default: {
                        throw new Error("unreachable");
                      }
                      }
                    }
                    else if ($ail_m8 && !($ail_m9)) {
                      const $ail_m12: { $ail_kind: "ok"; remainder: string; value: string } = std__dec__divide_round_half_even_result(($ailDivMod((an.value * psc.value), ad.value)[0] + 1n), scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($ail_m12.$ail_kind) {
                      case "ok": {
                        const r = $ail_m12;
                        return { $ail_kind: "ok", value: r.value, remainder: r.remainder };
                      }
                      default: {
                        throw new Error("unreachable");
                      }
                      }
                    }
                    else if (!($ail_m8) && $ail_m10) {
                      const $ail_m13: { $ail_kind: "ok"; remainder: string; value: string } = std__dec__divide_round_half_even_result(($ailDivMod((an.value * psc.value), ad.value)[0] + 1n), scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($ail_m13.$ail_kind) {
                      case "ok": {
                        const r = $ail_m13;
                        return { $ail_kind: "ok", value: r.value, remainder: r.remainder };
                      }
                      default: {
                        throw new Error("unreachable");
                      }
                      }
                    }
                    else {
                      const $ail_m14: { $ail_kind: "ok"; remainder: string; value: string } = std__dec__divide_round_half_even_result($ailDivMod((an.value * psc.value), ad.value)[0], scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($ail_m14.$ail_kind) {
                      case "ok": {
                        const r = $ail_m14;
                        return { $ail_kind: "ok", value: r.value, remainder: r.remainder };
                      }
                      default: {
                        throw new Error("unreachable");
                      }
                      }
                    }
                  }
                  default: {
                    throw new Error("unreachable");
                  }
                  }
                }
                default: {
                  throw new Error("unreachable");
                }
                }
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        return { $ail_kind: "math.negative_scale", scale: scale };
      }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__convert__dec_to_str(value: string): { $ail_kind: "ok"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const p = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; value: bigint } = std__int__abs(p.coefficient);
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const a = $ail_m2;
      const $ail_m3: { $ail_kind: "ok"; value: bigint } = std__int__pow_from(10n, p.scale);
      switch ($ail_m3.$ail_kind) {
      case "ok": {
        const p10 = $ail_m3;
        const $ail_m4: { $ail_kind: "ok"; value: string } = std__convert__int_to_str($ailDivMod(a.value, p10.value)[0]);
        switch ($ail_m4.$ail_kind) {
        case "ok": {
          const ip = $ail_m4;
          const $ail_m5: { $ail_kind: "ok"; value: string } = std__convert__int_to_str(($ailDivMod(a.value, p10.value)[1] + p10.value));
          switch ($ail_m5.$ail_kind) {
          case "ok": {
            const fp = $ail_m5;
            if ((p.coefficient < 0n)) {
              return { $ail_kind: "ok", value: ((("-" + ip.value) + ".") + $ailStrSlice(fp.value, 1n, (p.scale + 1n))) };
            }
            else {
              return { $ail_kind: "ok", value: ((ip.value + ".") + $ailStrSlice(fp.value, 1n, (p.scale + 1n))) };
            }
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__convert__dec_to_int_exact(value: string): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "convert.fractional_value"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const p = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const s = $ail_m2;
      if (s.exact) {
        return { $ail_kind: "ok", value: s.stripped };
      }
      else {
        return { $ail_kind: "convert.fractional_value", value: value };
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__dec__sqrt_exact(value: string): { $ail_kind: "ok"; value: string } | { $ail_kind: "math.dec_negative_input"; value: string } | { $ail_kind: "math.nonrepresentable_result"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const p = $ail_m1;
    if ((p.coefficient < 0n)) {
      return { $ail_kind: "math.dec_negative_input", value: value };
    }
    else {
      if (($ailDivMod(p.scale, 2n)[1] === 0n)) {
        const $ail_m2: { $ail_kind: "ok"; value: bigint } = std__int__sqrt_floor_search(p.coefficient, 0n, (p.coefficient + 1n));
        switch ($ail_m2.$ail_kind) {
        case "ok": {
          const r = $ail_m2;
          if (((p.coefficient - (r.value * r.value)) === 0n)) {
            const $ail_m3: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(r.value);
            switch ($ail_m3.$ail_kind) {
            case "ok": {
              const qi = $ail_m3;
              const $ail_m4: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - $ailDivMod(p.scale, 2n)[0]));
              switch ($ail_m4.$ail_kind) {
              case "ok": {
                const q = $ail_m4;
                return { $ail_kind: "ok", value: q.value };
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          else {
            return { $ail_kind: "math.nonrepresentable_result", value: value };
          }
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $ail_m5: { $ail_kind: "ok"; value: bigint } = std__int__sqrt_floor_search((p.coefficient * 10n), 0n, ((p.coefficient * 10n) + 1n));
        switch ($ail_m5.$ail_kind) {
        case "ok": {
          const r = $ail_m5;
          if ((((p.coefficient * 10n) - (r.value * r.value)) === 0n)) {
            const $ail_m6: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(r.value);
            switch ($ail_m6.$ail_kind) {
            case "ok": {
              const qi = $ail_m6;
              const $ail_m7: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - $ailDivMod((p.scale + 1n), 2n)[0]));
              switch ($ail_m7.$ail_kind) {
              case "ok": {
                const q = $ail_m7;
                return { $ail_kind: "ok", value: q.value };
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          else {
            return { $ail_kind: "math.nonrepresentable_result", value: value };
          }
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__convert__bool_to_str(value: boolean): { $ail_kind: "ok"; value: string } {
  if (value) {
    return { $ail_kind: "ok", value: "true" };
  }
  else {
    return { $ail_kind: "ok", value: "false" };
  }
}
export function std__convert__str_to_bool(value: string): { $ail_kind: "ok"; value: boolean } | { $ail_kind: "convert.invalid_boolean"; value: string } {
  if (value === "true") {
    return { $ail_kind: "ok", value: true };
  }
  else if (value === "false") {
    return { $ail_kind: "ok", value: false };
  }
  return { $ail_kind: "convert.invalid_boolean", value: value };
}
export function std__convert__bool_to_int(value: boolean): { $ail_kind: "ok"; value: bigint } {
  if (value) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    return { $ail_kind: "ok", value: 0n };
  }
}
export function std__convert__int_to_bool(value: bigint): { $ail_kind: "ok"; value: boolean } | { $ail_kind: "convert.invalid_boolean_encoding"; value: bigint } {
  const $ail_m1 = (value === 0n);
  const $ail_m2 = (value === 1n);
  if ($ail_m1) {
    return { $ail_kind: "ok", value: false };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "convert.invalid_boolean_encoding", value: value };
  }
}
export function std__convert__bool_to_dec(value: boolean): { $ail_kind: "ok"; value: string } {
  if (value) {
    return { $ail_kind: "ok", value: "1.0" };
  }
  else {
    return { $ail_kind: "ok", value: "0.0" };
  }
}
export function std__convert__dec_to_bool(value: string): { $ail_kind: "ok"; value: boolean } | { $ail_kind: "convert.invalid_dec_encoding"; value: string } {
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
export function std__convert__int_to_dec_from(m: bigint, acc: string, step: string): { $ail_kind: "ok"; value: string } {
  if ((m <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec_from((m - 1n), $ailDecAdd(acc, step), step);
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
export function std__convert__int_to_str_from(n: bigint, acc: string, m: bigint): { $ail_kind: "ok"; value: string } {
  if ((m <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    if ((n <= 9n)) {
      return { $ail_kind: "ok", value: ($ailStrSlice("0123456789", n, (n + 1n)) + acc) };
    }
    else {
      const $ail_m1: { $ail_kind: "ok"; value: string } = std__convert__int_to_str_from($ailDivMod(n, 10n)[0], ($ailStrSlice("0123456789", $ailDivMod(n, 10n)[1], ($ailDivMod(n, 10n)[1] + 1n)) + acc), (m - 1n));
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
}
export function std__convert__int_to_str(value: bigint): { $ail_kind: "ok"; value: string } {
  if ((value < 0n)) {
    const $ail_m1: { $ail_kind: "ok"; value: string } = std__convert__int_to_str_from((0n - value), "", ((0n - value) + 1n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const r = $ail_m1;
      return { $ail_kind: "ok", value: ("-" + r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  else {
    const $ail_m2: { $ail_kind: "ok"; value: string } = std__convert__int_to_str_from(value, "", (value + 1n));
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const r = $ail_m2;
      return { $ail_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__convert__str_to_int_from(orig: string, value: string, pos: bigint, acc: bigint, n: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "convert.invalid_integer"; value: string } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    if (($ailStrAt(value, pos) >= 48n)) {
      if (($ailStrAt(value, pos) <= 57n)) {
        const $ail_m1: { $ail_kind: "ok"; value: bigint } | { $ail_kind: "convert.invalid_integer"; value: string } = std__convert__str_to_int_from(orig, value, (pos + 1n), ((acc * 10n) + ($ailStrAt(value, pos) - 48n)), (n - 1n));
        switch ($ail_m1.$ail_kind) {
        case "convert.invalid_integer": {
          const e = $ail_m1;
          return { $ail_kind: "convert.invalid_integer", value: e.value };
        }
        case "ok": {
          const r = $ail_m1;
          return { $ail_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
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
export function std__convert__str_to_int(value: string): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "convert.invalid_integer"; value: string } {
  if (((BigInt([...value].length)) === 0n)) {
    return { $ail_kind: "convert.invalid_integer", value: value };
  }
  else {
    if (($ailStrSlice(value, 0n, 1n) === "-")) {
      if (((BigInt([...value].length)) >= 2n)) {
        const $ail_m1: { $ail_kind: "ok"; value: bigint } | { $ail_kind: "convert.invalid_integer"; value: string } = std__convert__str_to_int_from(value, value, 1n, 0n, ((BigInt([...value].length)) - 1n));
        switch ($ail_m1.$ail_kind) {
        case "convert.invalid_integer": {
          const e = $ail_m1;
          return { $ail_kind: "convert.invalid_integer", value: e.value };
        }
        case "ok": {
          const r = $ail_m1;
          return { $ail_kind: "ok", value: (0n - r.value) };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        return { $ail_kind: "convert.invalid_integer", value: value };
      }
    }
    else {
      const $ail_m2: { $ail_kind: "ok"; value: bigint } | { $ail_kind: "convert.invalid_integer"; value: string } = std__convert__str_to_int_from(value, value, 0n, 0n, (BigInt([...value].length)));
      switch ($ail_m2.$ail_kind) {
      case "convert.invalid_integer": {
        const e = $ail_m2;
        return { $ail_kind: "convert.invalid_integer", value: e.value };
      }
      case "ok": {
        const r = $ail_m2;
        return { $ail_kind: "ok", value: r.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
}
export function std__convert__str_to_dec_from(orig: string, value: string, pos: bigint, coeff: bigint, scale: bigint, dot: bigint, n: bigint): { $ail_kind: "ok"; value: string } | { $ail_kind: "convert.invalid_decimal"; value: string } {
  if ((n <= 0n)) {
    if ((dot === 0n)) {
      const $ail_m1: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(coeff);
      switch ($ail_m1.$ail_kind) {
      case "ok": {
        const whole = $ail_m1;
        return { $ail_kind: "ok", value: whole.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      if ((scale === 0n)) {
        return { $ail_kind: "convert.invalid_decimal", value: orig };
      }
      else {
        const $ail_m2: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec(coeff);
        switch ($ail_m2.$ail_kind) {
        case "ok": {
          const whole = $ail_m2;
          const $ail_m3: { $ail_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(whole.value, (0n - scale));
          switch ($ail_m3.$ail_kind) {
          case "ok": {
            const moved = $ail_m3;
            return { $ail_kind: "ok", value: moved.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
  }
  else {
    if (($ailStrSlice(value, pos, (pos + 1n)) === ".")) {
      if ((dot === 0n)) {
        const $ail_m4: { $ail_kind: "ok"; value: string } | { $ail_kind: "convert.invalid_decimal"; value: string } = std__convert__str_to_dec_from(orig, value, (pos + 1n), coeff, scale, 1n, (n - 1n));
        switch ($ail_m4.$ail_kind) {
        case "convert.invalid_decimal": {
          const e = $ail_m4;
          return { $ail_kind: "convert.invalid_decimal", value: e.value };
        }
        case "ok": {
          const r = $ail_m4;
          return { $ail_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
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
            const $ail_m5: { $ail_kind: "ok"; value: string } | { $ail_kind: "convert.invalid_decimal"; value: string } = std__convert__str_to_dec_from(orig, value, (pos + 1n), ((coeff * 10n) + ($ailStrAt(value, pos) - 48n)), scale, dot, (n - 1n));
            switch ($ail_m5.$ail_kind) {
            case "convert.invalid_decimal": {
              const e = $ail_m5;
              return { $ail_kind: "convert.invalid_decimal", value: e.value };
            }
            case "ok": {
              const r = $ail_m5;
              return { $ail_kind: "ok", value: r.value };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          else {
            const $ail_m6: { $ail_kind: "ok"; value: string } | { $ail_kind: "convert.invalid_decimal"; value: string } = std__convert__str_to_dec_from(orig, value, (pos + 1n), ((coeff * 10n) + ($ailStrAt(value, pos) - 48n)), (scale + 1n), dot, (n - 1n));
            switch ($ail_m6.$ail_kind) {
            case "convert.invalid_decimal": {
              const e = $ail_m6;
              return { $ail_kind: "convert.invalid_decimal", value: e.value };
            }
            case "ok": {
              const r = $ail_m6;
              return { $ail_kind: "ok", value: r.value };
            }
            default: {
              throw new Error("unreachable");
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
export function std__convert__str_to_dec(value: string): { $ail_kind: "ok"; value: string } | { $ail_kind: "convert.invalid_decimal"; value: string } {
  if (((BigInt([...value].length)) === 0n)) {
    return { $ail_kind: "convert.invalid_decimal", value: value };
  }
  else {
    if (($ailStrSlice(value, 0n, 1n) === "-")) {
      if (((BigInt([...value].length)) >= 2n)) {
        const $ail_m1: { $ail_kind: "ok"; value: string } | { $ail_kind: "convert.invalid_decimal"; value: string } = std__convert__str_to_dec_from(value, value, 1n, 0n, 0n, 0n, ((BigInt([...value].length)) - 1n));
        switch ($ail_m1.$ail_kind) {
        case "convert.invalid_decimal": {
          const e = $ail_m1;
          return { $ail_kind: "convert.invalid_decimal", value: e.value };
        }
        case "ok": {
          const r = $ail_m1;
          return { $ail_kind: "ok", value: $ailDecSub("0.0", r.value) };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        return { $ail_kind: "convert.invalid_decimal", value: value };
      }
    }
    else {
      const $ail_m2: { $ail_kind: "ok"; value: string } | { $ail_kind: "convert.invalid_decimal"; value: string } = std__convert__str_to_dec_from(value, value, 0n, 0n, 0n, 0n, (BigInt([...value].length)));
      switch ($ail_m2.$ail_kind) {
      case "convert.invalid_decimal": {
        const e = $ail_m2;
        return { $ail_kind: "convert.invalid_decimal", value: e.value };
      }
      case "ok": {
        const r = $ail_m2;
        return { $ail_kind: "ok", value: r.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
}
export function std__convert__int_to_dec(value: bigint): { $ail_kind: "ok"; value: string } {
  if ((value >= 0n)) {
    const $ail_m1: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec_from(value, "0.0", "1.0");
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
  else {
    const $ail_m2: { $ail_kind: "ok"; value: string } = std__convert__int_to_dec_from((0n - value), "0.0", "-1.0");
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const r = $ail_m2;
      return { $ail_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__backoff__delay(attempt: bigint, base_ms: bigint, cap_ms: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "math.negative_input"; value: bigint } {
  const $ail_m1 = (base_ms >= 0n);
  const $ail_m2 = (cap_ms >= 0n);
  const $ail_m3 = (attempt >= 0n);
  const $ail_m4 = (attempt <= 30n);
  if ($ail_m1 && $ail_m2 && $ail_m3 && $ail_m4) {
    const $ail_m5: { $ail_kind: "ok"; value: bigint } = std__int__pow_from(2n, attempt);
    switch ($ail_m5.$ail_kind) {
    case "ok": {
      const p = $ail_m5;
      if (((base_ms * p.value) <= cap_ms)) {
        return { $ail_kind: "ok", value: (base_ms * p.value) };
      }
      else {
        return { $ail_kind: "ok", value: cap_ms };
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  else if ($ail_m1 && $ail_m2 && $ail_m3 && !($ail_m4)) {
    return { $ail_kind: "ok", value: cap_ms };
  }
  else if ($ail_m1 && $ail_m2 && !($ail_m3)) {
    return { $ail_kind: "math.negative_input", value: attempt };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "math.negative_input", value: cap_ms };
  }
  else {
    return { $ail_kind: "math.negative_input", value: base_ms };
  }
}
