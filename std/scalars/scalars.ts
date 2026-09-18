// GENERATED from scalars.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type ScalarsResult = { $can_kind: "ok"; coefficient: bigint; scale: bigint } | { $can_kind: "ok"; discarded: string; value: string } | { $can_kind: "ok"; exact: boolean; stripped: bigint } | { $can_kind: "ok"; fives: bigint; rest: bigint; twos: bigint } | { $can_kind: "ok"; remainder: bigint; root: bigint } | { $can_kind: "ok"; remainder: string; value: string } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: boolean } | { $can_kind: "ok"; value: string } | { $can_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "math.out_of_range"; value: bigint; lower: bigint; upper: bigint } | { $can_kind: "math.dec_invalid_bounds"; lower: string; upper: string } | { $can_kind: "math.dec_out_of_range"; value: string; lower: string; upper: string } | { $can_kind: "math.negative_input"; value: bigint } | { $can_kind: "math.nonpositive_input"; value: bigint } | { $can_kind: "math.negative_exponent"; exponent: bigint } | { $can_kind: "math.invalid_count"; n: bigint; k: bigint } | { $can_kind: "math.negative_scale"; scale: bigint } | { $can_kind: "math.dec_zero_divisor"; divisor: string } | { $can_kind: "math.nonterminating_decimal"; dividend: string; divisor: string } | { $can_kind: "math.dec_negative_input"; value: string } | { $can_kind: "math.nonrepresentable_result"; value: string } | { $can_kind: "convert.fractional_value"; value: string } | { $can_kind: "convert.invalid_boolean"; value: string } | { $can_kind: "convert.invalid_integer"; value: string } | { $can_kind: "convert.invalid_decimal"; value: string } | { $can_kind: "convert.invalid_boolean_encoding"; value: bigint } | { $can_kind: "convert.invalid_dec_encoding"; value: string };
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
function $canDecSplit(d: string): { neg: boolean; ip: string; fp: string } {
  let neg = false;
  if (d.startsWith("-")) {
    neg = true;
    d = d.slice(1);
  }
  const dot = d.indexOf(".");
  return { neg, ip: d.slice(0, dot), fp: d.slice(dot + 1) };
}
function $canDecNorm(ip: string, fp: string, neg: boolean): string {
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
function $canDecMant(p: { neg: boolean; ip: string; fp: string }, scale: number): bigint {
  let f = p.fp;
  while (f.length < scale) {
    f += "0";
  }
  const m = BigInt(p.ip + f);
  return p.neg ? -m : m;
}
function $canDecFromMant(m: bigint, scale: number): string {
  let neg = false;
  if (m < 0n) {
    neg = true;
    m = -m;
  }
  let digits = m.toString();
  while (digits.length < scale + 1) {
    digits = "0" + digits;
  }
  return $canDecNorm(digits.slice(0, digits.length - scale), digits.slice(digits.length - scale), neg);
}
function $canDecAdd(a: string, b: string): string {
  const A = $canDecSplit(a);
  const B = $canDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $canDecFromMant($canDecMant(A, s) + $canDecMant(B, s), s);
}
function $canDecSub(a: string, b: string): string {
  const A = $canDecSplit(a);
  const B = $canDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $canDecFromMant($canDecMant(A, s) - $canDecMant(B, s), s);
}
function $canDecMul(a: string, b: string): string {
  const A = $canDecSplit(a);
  const B = $canDecSplit(b);
  return $canDecFromMant($canDecMant(A, A.fp.length) * $canDecMant(B, B.fp.length), A.fp.length + B.fp.length);
}
function $canDecGe(a: string, b: string): boolean {
  const A = $canDecSplit(a);
  const B = $canDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $canDecMant(A, s) >= $canDecMant(B, s);
}
function $canDecLe(a: string, b: string): boolean {
  const A = $canDecSplit(a);
  const B = $canDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $canDecMant(A, s) <= $canDecMant(B, s);
}
function $canDecLt(a: string, b: string): boolean {
  const A = $canDecSplit(a);
  const B = $canDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $canDecMant(A, s) < $canDecMant(B, s);
}
function $canDecParts(d: string): { coefficient: bigint; scale: bigint } {
  const p = $canDecSplit(d);
  return { coefficient: $canDecMant(p, p.fp.length), scale: BigInt(p.fp.length) };
}
// Byte-order string comparison: UTF-8 bytes, matching Go.
function $canStrCmp(a: string, b: string): number {
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
function $canStrGe(a: string, b: string): boolean {
  return $canStrCmp(a, b) >= 0;
}
function $canStrAt(s: string, i: bigint): bigint {
  const cps = [...s];
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("str index out of range");
  const k = Number(i);
  if (k >= cps.length) throw new Error("str index out of range");
  const cp = cps[k].codePointAt(0);
  if (cp === undefined) throw new Error("str index out of range");
  return BigInt(cp);
}
function $canStrSlice(s: string, a: bigint, b: bigint): string {
  const cps = [...s];
  const toIdx = (x: bigint): number => { if (x < 0n || x > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("str slice out of range"); return Number(x); };
  const lo = toIdx(a), hi = toIdx(b);
  if (lo > hi || hi > cps.length) throw new Error("str slice out of range");
  return cps.slice(lo, hi).join("");
}
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
export function std__bool__not(value: boolean): { $can_kind: "ok"; value: boolean } {
  if (value) {
    return { $can_kind: "ok", value: false };
  }
  else {
    return { $can_kind: "ok", value: true };
  }
}
export function std__bool__and(left: boolean, right: boolean): { $can_kind: "ok"; value: boolean } {
  const $can_m1 = left;
  const $can_m2 = right;
  if ($can_m1 && $can_m2) {
    return { $can_kind: "ok", value: true };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "ok", value: false };
  }
  else {
    return { $can_kind: "ok", value: false };
  }
}
export function std__bool__or(left: boolean, right: boolean): { $can_kind: "ok"; value: boolean } {
  const $can_m1 = left;
  const $can_m2 = right;
  if ($can_m1) {
    return { $can_kind: "ok", value: true };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: true };
  }
  else {
    return { $can_kind: "ok", value: false };
  }
}
export function std__bool__xor(left: boolean, right: boolean): { $can_kind: "ok"; value: boolean } {
  const $can_m1 = left;
  const $can_m2 = right;
  if ($can_m1 && $can_m2) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "ok", value: true };
  }
  else if (!($can_m1) && $can_m2) {
    return { $can_kind: "ok", value: true };
  }
  else {
    return { $can_kind: "ok", value: false };
  }
}
export function std__bool__implies(left: boolean, right: boolean): { $can_kind: "ok"; value: boolean } {
  const $can_m1 = left;
  const $can_m2 = right;
  if ($can_m1 && $can_m2) {
    return { $can_kind: "ok", value: true };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "ok", value: false };
  }
  else {
    return { $can_kind: "ok", value: true };
  }
}
export function std__bool__equivalent(left: boolean, right: boolean): { $can_kind: "ok"; value: boolean } {
  const $can_m1 = left;
  const $can_m2 = right;
  if ($can_m1 && $can_m2) {
    return { $can_kind: "ok", value: true };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "ok", value: false };
  }
  else if (!($can_m1) && $can_m2) {
    return { $can_kind: "ok", value: false };
  }
  else {
    return { $can_kind: "ok", value: true };
  }
}
export function std__compare__int(left: bigint, right: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1 = (left === right);
  const $can_m2 = (left >= right);
  if ($can_m1) {
    return { $can_kind: "ok", value: 0n };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: 1n };
  }
  else {
    return { $can_kind: "ok", value: -1n };
  }
}
export function std__compare__dec(left: string, right: string): { $can_kind: "ok"; value: bigint } {
  const $can_m1 = (left === right);
  const $can_m2 = $canDecGe(left, right);
  if ($can_m1) {
    return { $can_kind: "ok", value: 0n };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: 1n };
  }
  else {
    return { $can_kind: "ok", value: -1n };
  }
}
export function std__compare__str(left: string, right: string): { $can_kind: "ok"; value: bigint } {
  const $can_m1 = (left === right);
  const $can_m2 = $canStrGe(left, right);
  if ($can_m1) {
    return { $can_kind: "ok", value: 0n };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: 1n };
  }
  else {
    return { $can_kind: "ok", value: -1n };
  }
}
export function std__compare__bool(left: boolean, right: boolean): { $can_kind: "ok"; value: bigint } {
  const $can_m1 = (left === right);
  const $can_m2 = left;
  if ($can_m1) {
    return { $can_kind: "ok", value: 0n };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: 1n };
  }
  else {
    return { $can_kind: "ok", value: -1n };
  }
}
export function std__select__int(condition: boolean, when_true: bigint, when_false: bigint): { $can_kind: "ok"; value: bigint } {
  if (condition) {
    return { $can_kind: "ok", value: when_true };
  }
  else {
    return { $can_kind: "ok", value: when_false };
  }
}
export function std__select__bool(condition: boolean, when_true: boolean, when_false: boolean): { $can_kind: "ok"; value: boolean } {
  if (condition) {
    return { $can_kind: "ok", value: when_true };
  }
  else {
    return { $can_kind: "ok", value: when_false };
  }
}
export function std__select__str(condition: boolean, when_true: string, when_false: string): { $can_kind: "ok"; value: string } {
  if (condition) {
    return { $can_kind: "ok", value: when_true };
  }
  else {
    return { $can_kind: "ok", value: when_false };
  }
}
export function std__select__dec(condition: boolean, when_true: string, when_false: string): { $can_kind: "ok"; value: string } {
  if (condition) {
    return { $can_kind: "ok", value: when_true };
  }
  else {
    return { $can_kind: "ok", value: when_false };
  }
}
export function std__int__is_zero(value: bigint): { $can_kind: "ok"; value: boolean } {
  if ((value === 0n)) {
    return { $can_kind: "ok", value: true };
  }
  else {
    return { $can_kind: "ok", value: false };
  }
}
export function std__int__is_positive(value: bigint): { $can_kind: "ok"; value: boolean } {
  const $can_m1 = (value === 0n);
  const $can_m2 = (value >= 0n);
  if ($can_m1) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: true };
  }
  else {
    return { $can_kind: "ok", value: false };
  }
}
export function std__int__is_negative(value: bigint): { $can_kind: "ok"; value: boolean } {
  const $can_m1 = (value === 0n);
  const $can_m2 = (value <= 0n);
  if ($can_m1) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: true };
  }
  else {
    return { $can_kind: "ok", value: false };
  }
}
export function std__int__in_closed_range(value: bigint, lower: bigint, upper: bigint): { $can_kind: "ok"; value: boolean } | { $can_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } {
  const $can_m1 = (lower <= upper);
  const $can_m2 = (value >= lower);
  const $can_m3 = (value <= upper);
  if ($can_m1 && $can_m2 && $can_m3) {
    return { $can_kind: "ok", value: true };
  }
  else if ($can_m1 && $can_m2 && !($can_m3)) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "ok", value: false };
  }
  else {
    return { $can_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__in_open_range(value: bigint, lower: bigint, upper: bigint): { $can_kind: "ok"; value: boolean } | { $can_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } {
  const $can_m1 = (lower <= upper);
  const $can_m2 = (value <= lower);
  const $can_m3 = (value >= upper);
  if ($can_m1 && $can_m2) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m1 && !($can_m2) && $can_m3) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m1 && !($can_m2) && !($can_m3)) {
    return { $can_kind: "ok", value: true };
  }
  else {
    return { $can_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__is_zero(value: string): { $can_kind: "ok"; value: boolean } {
  if ((value === "0.0")) {
    return { $can_kind: "ok", value: true };
  }
  else {
    return { $can_kind: "ok", value: false };
  }
}
export function std__dec__is_positive(value: string): { $can_kind: "ok"; value: boolean } {
  const $can_m1 = (value === "0.0");
  const $can_m2 = $canDecGe(value, "0.0");
  if ($can_m1) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: true };
  }
  else {
    return { $can_kind: "ok", value: false };
  }
}
export function std__dec__is_negative(value: string): { $can_kind: "ok"; value: boolean } {
  const $can_m1 = (value === "0.0");
  const $can_m2 = $canDecLe(value, "0.0");
  if ($can_m1) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: true };
  }
  else {
    return { $can_kind: "ok", value: false };
  }
}
export function std__dec__in_closed_range(value: string, lower: string, upper: string): { $can_kind: "ok"; value: boolean } | { $can_kind: "math.dec_invalid_bounds"; lower: string; upper: string } {
  const $can_m1 = $canDecLe(lower, upper);
  const $can_m2 = $canDecGe(value, lower);
  const $can_m3 = $canDecLe(value, upper);
  if ($can_m1 && $can_m2 && $can_m3) {
    return { $can_kind: "ok", value: true };
  }
  else if ($can_m1 && $can_m2 && !($can_m3)) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "ok", value: false };
  }
  else {
    return { $can_kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__in_open_range(value: string, lower: string, upper: string): { $can_kind: "ok"; value: boolean } | { $can_kind: "math.dec_invalid_bounds"; lower: string; upper: string } {
  const $can_m1 = $canDecLe(lower, upper);
  const $can_m2 = $canDecLe(value, lower);
  const $can_m3 = $canDecGe(value, upper);
  if ($can_m1 && $can_m2) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m1 && !($can_m2) && $can_m3) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m1 && !($can_m2) && !($can_m3)) {
    return { $can_kind: "ok", value: true };
  }
  else {
    return { $can_kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__abs(value: bigint): { $can_kind: "ok"; value: bigint } {
  if ((value >= 0n)) {
    return { $can_kind: "ok", value: value };
  }
  else {
    return { $can_kind: "ok", value: (-value) };
  }
}
export function std__int__negate(value: bigint): { $can_kind: "ok"; value: bigint } {
  return { $can_kind: "ok", value: (-value) };
}
export function std__int__sign(value: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1 = (value === 0n);
  const $can_m2 = (value >= 0n);
  if ($can_m1) {
    return { $can_kind: "ok", value: 0n };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: 1n };
  }
  else {
    return { $can_kind: "ok", value: -1n };
  }
}
export function std__int__min(left: bigint, right: bigint): { $can_kind: "ok"; value: bigint } {
  if ((left <= right)) {
    return { $can_kind: "ok", value: left };
  }
  else {
    return { $can_kind: "ok", value: right };
  }
}
export function std__int__max(left: bigint, right: bigint): { $can_kind: "ok"; value: bigint } {
  if ((left >= right)) {
    return { $can_kind: "ok", value: left };
  }
  else {
    return { $can_kind: "ok", value: right };
  }
}
export function std__int__clamp(value: bigint, lower: bigint, upper: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } {
  const $can_m1 = (lower <= upper);
  const $can_m2 = (value >= lower);
  const $can_m3 = (value <= upper);
  if ($can_m1 && $can_m2 && $can_m3) {
    return { $can_kind: "ok", value: value };
  }
  else if ($can_m1 && $can_m2 && !($can_m3)) {
    return { $can_kind: "ok", value: upper };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "ok", value: lower };
  }
  else {
    return { $can_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__distance(left: bigint, right: bigint): { $can_kind: "ok"; value: bigint } {
  if ((left >= right)) {
    return { $can_kind: "ok", value: (left - right) };
  }
  else {
    return { $can_kind: "ok", value: (right - left) };
  }
}
export function std__int__square(value: bigint): { $can_kind: "ok"; value: bigint } {
  return { $can_kind: "ok", value: (value * value) };
}
export function std__int__pow_from(base: bigint, exponent: bigint): { $can_kind: "ok"; value: bigint } {
  if ((exponent <= 0n)) {
    return { $can_kind: "ok", value: 1n };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__pow_from(base, (exponent - 1n));
    switch ($can_m1.$can_kind) {
    case "ok": {
      const r = $can_m1;
      return { $can_kind: "ok", value: (base * r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__int__pow(base: bigint, exponent: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "math.negative_exponent"; exponent: bigint } {
  if ((exponent >= 0n)) {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__pow_from(base, exponent);
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
  else {
    return { $can_kind: "math.negative_exponent", exponent: exponent };
  }
}
export function std__int__add_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "math.out_of_range"; value: bigint; lower: bigint; upper: bigint } {
  const $can_m1 = (lower <= upper);
  const $can_m2 = ((left + right) >= lower);
  const $can_m3 = ((left + right) <= upper);
  if ($can_m1 && $can_m2 && $can_m3) {
    return { $can_kind: "ok", value: (left + right) };
  }
  else if ($can_m1 && $can_m2 && !($can_m3)) {
    return { $can_kind: "math.out_of_range", value: (left + right), lower: lower, upper: upper };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "math.out_of_range", value: (left + right), lower: lower, upper: upper };
  }
  else {
    return { $can_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__subtract_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "math.out_of_range"; value: bigint; lower: bigint; upper: bigint } {
  const $can_m1 = (lower <= upper);
  const $can_m2 = ((left - right) >= lower);
  const $can_m3 = ((left - right) <= upper);
  if ($can_m1 && $can_m2 && $can_m3) {
    return { $can_kind: "ok", value: (left - right) };
  }
  else if ($can_m1 && $can_m2 && !($can_m3)) {
    return { $can_kind: "math.out_of_range", value: (left - right), lower: lower, upper: upper };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "math.out_of_range", value: (left - right), lower: lower, upper: upper };
  }
  else {
    return { $can_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__multiply_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "math.out_of_range"; value: bigint; lower: bigint; upper: bigint } {
  const $can_m1 = (lower <= upper);
  const $can_m2 = ((left * right) >= lower);
  const $can_m3 = ((left * right) <= upper);
  if ($can_m1 && $can_m2 && $can_m3) {
    return { $can_kind: "ok", value: (left * right) };
  }
  else if ($can_m1 && $can_m2 && !($can_m3)) {
    return { $can_kind: "math.out_of_range", value: (left * right), lower: lower, upper: upper };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "math.out_of_range", value: (left * right), lower: lower, upper: upper };
  }
  else {
    return { $can_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__factorial_from(value: bigint): { $can_kind: "ok"; value: bigint } {
  if ((value <= 0n)) {
    return { $can_kind: "ok", value: 1n };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__factorial_from((value - 1n));
    switch ($can_m1.$can_kind) {
    case "ok": {
      const r = $can_m1;
      return { $can_kind: "ok", value: (value * r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__int__factorial(value: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "math.negative_input"; value: bigint } {
  if ((value >= 0n)) {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__factorial_from(value);
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
  else {
    return { $can_kind: "math.negative_input", value: value };
  }
}
export function std__int__sum_from(upper: bigint): { $can_kind: "ok"; value: bigint } {
  if ((upper <= 0n)) {
    return { $can_kind: "ok", value: 0n };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__sum_from((upper - 1n));
    switch ($can_m1.$can_kind) {
    case "ok": {
      const r = $can_m1;
      return { $can_kind: "ok", value: (upper + r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__int__sum_to(upper: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "math.negative_input"; value: bigint } {
  if ((upper >= 0n)) {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__sum_from(upper);
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
  else {
    return { $can_kind: "math.negative_input", value: upper };
  }
}
export function std__int__binomial_from(n: bigint, k: bigint, i: bigint, num: bigint, den: bigint): { $can_kind: "ok"; value: bigint } {
  if ((i <= 0n)) {
    return { $can_kind: "ok", value: $canDivMod(num, den)[0] };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__binomial_from(n, k, (i - 1n), (num * ((n - k) + i)), (den * i));
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
export function std__int__binomial(n: bigint, k: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "math.invalid_count"; n: bigint; k: bigint } {
  const $can_m1 = (k >= 0n);
  const $can_m2 = (k <= n);
  if ($can_m1 && $can_m2) {
    const $can_m3: { $can_kind: "ok"; value: bigint } = std__int__binomial_from(n, k, k, 1n, 1n);
    switch ($can_m3.$can_kind) {
    case "ok": {
      const r = $can_m3;
      return { $can_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "math.invalid_count", n: n, k: k };
  }
  else {
    return { $can_kind: "math.invalid_count", n: n, k: k };
  }
}
export function std__int__is_prime_from(value: bigint, e: bigint): { $can_kind: "ok"; value: boolean } {
  if ((e <= 0n)) {
    return { $can_kind: "ok", value: true };
  }
  else {
    if (($canDivMod(value, (e + 1n))[1] === 0n)) {
      return { $can_kind: "ok", value: false };
    }
    else {
      const $can_m1: { $can_kind: "ok"; value: boolean } = std__int__is_prime_from(value, (e - 1n));
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
}
export function std__int__is_prime(value: bigint): { $can_kind: "ok"; value: boolean } {
  if ((value >= 2n)) {
    const $can_m1: { $can_kind: "ok"; value: boolean } = std__int__is_prime_from(value, (value - 2n));
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
  else {
    return { $can_kind: "ok", value: false };
  }
}
export function std__int__next_power_of_two_pow2(e: bigint): { $can_kind: "ok"; value: bigint } {
  if ((e <= 0n)) {
    return { $can_kind: "ok", value: 1n };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__next_power_of_two_pow2((e - 1n));
    switch ($can_m1.$can_kind) {
    case "ok": {
      const r = $can_m1;
      return { $can_kind: "ok", value: (2n * r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__int__next_power_of_two_from(value: bigint, e: bigint, best: bigint): { $can_kind: "ok"; value: bigint } {
  if ((e <= 0n)) {
    if ((1n >= value)) {
      return { $can_kind: "ok", value: 1n };
    }
    else {
      return { $can_kind: "ok", value: best };
    }
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__next_power_of_two_pow2(e);
    switch ($can_m1.$can_kind) {
    case "ok": {
      const t = $can_m1;
      if ((t.value >= value)) {
        const $can_m2: { $can_kind: "ok"; value: bigint } = std__int__next_power_of_two_from(value, (e - 1n), t.value);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $can_m3: { $can_kind: "ok"; value: bigint } = std__int__next_power_of_two_from(value, (e - 1n), best);
        switch ($can_m3.$can_kind) {
        case "ok": {
          const r = $can_m3;
          return { $can_kind: "ok", value: r.value };
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
export function std__int__next_power_of_two(value: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "math.nonpositive_input"; value: bigint } {
  if ((value >= 1n)) {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__next_power_of_two_from(value, value, 1n);
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
  else {
    return { $can_kind: "math.nonpositive_input", value: value };
  }
}
export function std__int__sqrt_floor_search(value: bigint, lo: bigint, hi: bigint): { $can_kind: "ok"; value: bigint } {
  if (((hi - lo) <= 1n)) {
    return { $can_kind: "ok", value: lo };
  }
  else {
    if ((($canDivMod((lo + hi), 2n)[0] * $canDivMod((lo + hi), 2n)[0]) <= value)) {
      const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__sqrt_floor_search(value, $canDivMod((lo + hi), 2n)[0], hi);
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
    else {
      const $can_m2: { $can_kind: "ok"; value: bigint } = std__int__sqrt_floor_search(value, lo, $canDivMod((lo + hi), 2n)[0]);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", value: r.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
}
export function std__int__sqrt_floor(value: bigint): { $can_kind: "ok"; remainder: bigint; root: bigint } | { $can_kind: "math.negative_input"; value: bigint } {
  if ((value >= 0n)) {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__sqrt_floor_search(value, 0n, (value + 1n));
    switch ($can_m1.$can_kind) {
    case "ok": {
      const r = $can_m1;
      return { $can_kind: "ok", root: r.value, remainder: (value - (r.value * r.value)) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  else {
    return { $can_kind: "math.negative_input", value: value };
  }
}
export function std__int__gcd_euclid(a: bigint, b: bigint): { $can_kind: "ok"; value: bigint } {
  if ((b <= 0n)) {
    return { $can_kind: "ok", value: a };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__gcd_euclid(b, $canDivMod(a, b)[1]);
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
export function std__int__gcd(left: bigint, right: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__int__abs(left);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const a = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } = std__int__abs(right);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const b = $can_m2;
      const $can_m3: { $can_kind: "ok"; value: bigint } = std__int__gcd_euclid(a.value, b.value);
      switch ($can_m3.$can_kind) {
      case "ok": {
        const r = $can_m3;
        return { $can_kind: "ok", value: r.value };
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
export function std__int__lcm(left: bigint, right: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1 = (left === 0n);
  const $can_m2 = (right === 0n);
  if ($can_m1) {
    return { $can_kind: "ok", value: 0n };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: 0n };
  }
  else {
    const $can_m3: { $can_kind: "ok"; value: bigint } = std__int__gcd(left, right);
    switch ($can_m3.$can_kind) {
    case "ok": {
      const g = $can_m3;
      const $can_m4: { $can_kind: "ok"; value: bigint } = std__int__abs(left);
      switch ($can_m4.$can_kind) {
      case "ok": {
        const a = $can_m4;
        const $can_m5: { $can_kind: "ok"; value: bigint } = std__int__abs(right);
        switch ($can_m5.$can_kind) {
        case "ok": {
          const b = $can_m5;
          return { $can_kind: "ok", value: $canDivMod((a.value * b.value), g.value)[0] };
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
export function std__dec__abs(value: string): { $can_kind: "ok"; value: string } {
  if ($canDecGe(value, "0.0")) {
    return { $can_kind: "ok", value: value };
  }
  else {
    return { $can_kind: "ok", value: $canDecSub("0.0", value) };
  }
}
export function std__dec__negate(value: string): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: $canDecSub("0.0", value) };
}
export function std__dec__sign(value: string): { $can_kind: "ok"; value: bigint } {
  const $can_m1 = (value === "0.0");
  const $can_m2 = $canDecGe(value, "0.0");
  if ($can_m1) {
    return { $can_kind: "ok", value: 0n };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: 1n };
  }
  else {
    return { $can_kind: "ok", value: -1n };
  }
}
export function std__dec__min(left: string, right: string): { $can_kind: "ok"; value: string } {
  if ($canDecLe(left, right)) {
    return { $can_kind: "ok", value: left };
  }
  else {
    return { $can_kind: "ok", value: right };
  }
}
export function std__dec__max(left: string, right: string): { $can_kind: "ok"; value: string } {
  if ($canDecGe(left, right)) {
    return { $can_kind: "ok", value: left };
  }
  else {
    return { $can_kind: "ok", value: right };
  }
}
export function std__dec__clamp(value: string, lower: string, upper: string): { $can_kind: "ok"; value: string } | { $can_kind: "math.dec_invalid_bounds"; lower: string; upper: string } {
  const $can_m1 = $canDecLe(lower, upper);
  const $can_m2 = $canDecGe(value, lower);
  const $can_m3 = $canDecLe(value, upper);
  if ($can_m1 && $can_m2 && $can_m3) {
    return { $can_kind: "ok", value: value };
  }
  else if ($can_m1 && $can_m2 && !($can_m3)) {
    return { $can_kind: "ok", value: upper };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "ok", value: lower };
  }
  else {
    return { $can_kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__distance(left: string, right: string): { $can_kind: "ok"; value: string } {
  if ($canDecGe(left, right)) {
    return { $can_kind: "ok", value: $canDecSub(left, right) };
  }
  else {
    return { $can_kind: "ok", value: $canDecSub(right, left) };
  }
}
export function std__dec__square(value: string): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: $canDecMul(value, value) };
}
export function std__dec__pow_from(base: string, exponent: bigint): { $can_kind: "ok"; value: string } {
  if ((exponent <= 0n)) {
    return { $can_kind: "ok", value: "1.0" };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: string } = std__dec__pow_from(base, (exponent - 1n));
    switch ($can_m1.$can_kind) {
    case "ok": {
      const r = $can_m1;
      return { $can_kind: "ok", value: $canDecMul(base, r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__dec__pow(base: string, exponent: bigint): { $can_kind: "ok"; value: string } | { $can_kind: "math.negative_exponent"; exponent: bigint } {
  if ((exponent >= 0n)) {
    const $can_m1: { $can_kind: "ok"; value: string } = std__dec__pow_from(base, exponent);
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
  else {
    return { $can_kind: "math.negative_exponent", exponent: exponent };
  }
}
export function std__dec__lerp(start: string, end: string, fraction: string): { $can_kind: "ok"; value: string } | { $can_kind: "math.dec_out_of_range"; value: string; lower: string; upper: string } {
  const $can_m1 = $canDecGe(fraction, "0.0");
  const $can_m2 = $canDecLe(fraction, "1.0");
  if ($can_m1 && $can_m2) {
    return { $can_kind: "ok", value: $canDecAdd(start, $canDecMul($canDecSub(end, start), fraction)) };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "math.dec_out_of_range", value: fraction, lower: "0.0", upper: "1.0" };
  }
  else {
    return { $can_kind: "math.dec_out_of_range", value: fraction, lower: "0.0", upper: "1.0" };
  }
}
export function std__dec__scale_by_power_of_ten_from(value: string, e: bigint, factor: string): { $can_kind: "ok"; value: string } {
  if ((e <= 0n)) {
    return { $can_kind: "ok", value: value };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten_from($canDecMul(value, factor), (e - 1n), factor);
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
export function std__dec__scale_by_power_of_ten(value: string, exponent: bigint): { $can_kind: "ok"; value: string } {
  if ((exponent >= 0n)) {
    const $can_m1: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten_from(value, exponent, "10.0");
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
  else {
    const $can_m2: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten_from(value, (-exponent), "0.1");
    switch ($can_m2.$can_kind) {
    case "ok": {
      const r = $can_m2;
      return { $can_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__dec__from_parts(coefficient: bigint, scale: bigint): { $can_kind: "ok"; value: string } | { $can_kind: "math.negative_scale"; scale: bigint } {
  if ((scale >= 0n)) {
    const $can_m1: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(coefficient);
    switch ($can_m1.$can_kind) {
    case "ok": {
      const whole = $can_m1;
      const $can_m2: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(whole.value, (-scale));
      switch ($can_m2.$can_kind) {
      case "ok": {
        const moved = $can_m2;
        return { $can_kind: "ok", value: moved.value };
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
    return { $can_kind: "math.negative_scale", scale: scale };
  }
}
export function std__dec__parts(value: string): { $can_kind: "ok"; coefficient: bigint; scale: bigint } {
  const $can_m1: { $can_kind: "ok", coefficient: bigint, scale: bigint } = { $can_kind: "ok", ...$canDecParts(value) };
  switch ($can_m1.$can_kind) {
  case "ok": {
    const p = $can_m1;
    return { $can_kind: "ok", coefficient: p.coefficient, scale: p.scale };
  }
  }
}
export function std__dec__strip_from(remaining: bigint, current: bigint, exact: boolean): { $can_kind: "ok"; exact: boolean; stripped: bigint } {
  if ((remaining <= 0n)) {
    return { $can_kind: "ok", stripped: current, exact: exact };
  }
  else {
    if (($canDivMod(current, 10n)[1] === 0n)) {
      const $can_m1: { $can_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from((remaining - 1n), $canDivMod(current, 10n)[0], exact);
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        return { $can_kind: "ok", stripped: r.stripped, exact: r.exact };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      const $can_m2: { $can_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from((remaining - 1n), $canDivMod(current, 10n)[0], false);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", stripped: r.stripped, exact: r.exact };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
}
export function std__dec__truncate(value: string): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const p = $can_m1;
    const $can_m2: { $can_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const s = $can_m2;
      if (s.exact) {
        const $can_m3: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(s.stripped);
        switch ($can_m3.$can_kind) {
        case "ok": {
          const whole = $can_m3;
          return { $can_kind: "ok", value: whole.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        if ((s.stripped >= 0n)) {
          const $can_m4: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(s.stripped);
          switch ($can_m4.$can_kind) {
          case "ok": {
            const whole = $can_m4;
            return { $can_kind: "ok", value: whole.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          const $can_m5: { $can_kind: "ok"; value: string } = std__convert__int_to_dec((s.stripped + 1n));
          switch ($can_m5.$can_kind) {
          case "ok": {
            const whole = $can_m5;
            return { $can_kind: "ok", value: whole.value };
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
export function std__dec__floor(value: string): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const p = $can_m1;
    const $can_m2: { $can_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const s = $can_m2;
      const $can_m3: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(s.stripped);
      switch ($can_m3.$can_kind) {
      case "ok": {
        const whole = $can_m3;
        return { $can_kind: "ok", value: whole.value };
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
export function std__dec__ceil(value: string): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const p = $can_m1;
    const $can_m2: { $can_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const s = $can_m2;
      if (s.exact) {
        const $can_m3: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(s.stripped);
        switch ($can_m3.$can_kind) {
        case "ok": {
          const whole = $can_m3;
          return { $can_kind: "ok", value: whole.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $can_m4: { $can_kind: "ok"; value: string } = std__convert__int_to_dec((s.stripped + 1n));
        switch ($can_m4.$can_kind) {
        case "ok": {
          const whole = $can_m4;
          return { $can_kind: "ok", value: whole.value };
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
export function std__dec__strip257_from(fuel: bigint, current: bigint, twos: bigint, fives: bigint): { $can_kind: "ok"; fives: bigint; rest: bigint; twos: bigint } {
  if ((fuel <= 0n)) {
    return { $can_kind: "ok", rest: current, twos: twos, fives: fives };
  }
  else {
    if ((current === 1n)) {
      return { $can_kind: "ok", rest: 1n, twos: twos, fives: fives };
    }
    else {
      if (($canDivMod(current, 2n)[1] === 0n)) {
        const $can_m1: { $can_kind: "ok"; fives: bigint; rest: bigint; twos: bigint } = std__dec__strip257_from((fuel - 1n), $canDivMod(current, 2n)[0], (twos + 1n), fives);
        switch ($can_m1.$can_kind) {
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", rest: r.rest, twos: r.twos, fives: r.fives };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        if (($canDivMod(current, 5n)[1] === 0n)) {
          const $can_m2: { $can_kind: "ok"; fives: bigint; rest: bigint; twos: bigint } = std__dec__strip257_from((fuel - 1n), $canDivMod(current, 5n)[0], twos, (fives + 1n));
          switch ($can_m2.$can_kind) {
          case "ok": {
            const r = $can_m2;
            return { $can_kind: "ok", rest: r.rest, twos: r.twos, fives: r.fives };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          return { $can_kind: "ok", rest: current, twos: twos, fives: fives };
        }
      }
    }
  }
}
export function std__dec__divide_exact(dividend: string, divisor: string): { $can_kind: "ok"; value: string } | { $can_kind: "math.dec_zero_divisor"; divisor: string } | { $can_kind: "math.nonterminating_decimal"; dividend: string; divisor: string } {
  const $can_m1: { $can_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(divisor);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const d = $can_m1;
    if ((d.coefficient === 0n)) {
      return { $can_kind: "math.dec_zero_divisor", divisor: divisor };
    }
    else {
      const $can_m2: { $can_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(dividend);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const n = $can_m2;
        const $can_m3: { $can_kind: "ok"; value: bigint } = std__int__pow_from(10n, d.scale);
        switch ($can_m3.$can_kind) {
        case "ok": {
          const pds = $can_m3;
          const $can_m4: { $can_kind: "ok"; value: bigint } = std__int__pow_from(10n, n.scale);
          switch ($can_m4.$can_kind) {
          case "ok": {
            const pns = $can_m4;
            const $can_m5: { $can_kind: "ok"; value: bigint } = std__int__abs((n.coefficient * pns.value));
            switch ($can_m5.$can_kind) {
            case "ok": {
              const an = $can_m5;
              const $can_m6: { $can_kind: "ok"; value: bigint } = std__int__abs((d.coefficient * pds.value));
              switch ($can_m6.$can_kind) {
              case "ok": {
                const ad = $can_m6;
                const $can_m7: { $can_kind: "ok"; value: bigint } = std__int__gcd(an.value, ad.value);
                switch ($can_m7.$can_kind) {
                case "ok": {
                  const g = $can_m7;
                  const $can_m8: { $can_kind: "ok"; fives: bigint; rest: bigint; twos: bigint } = std__dec__strip257_from($canDivMod(ad.value, g.value)[0], $canDivMod(ad.value, g.value)[0], 0n, 0n);
                  switch ($can_m8.$can_kind) {
                  case "ok": {
                    const f = $can_m8;
                    if ((f.rest === 1n)) {
                      const $can_m9: { $can_kind: "ok"; value: bigint } = std__int__max(f.twos, f.fives);
                      switch ($can_m9.$can_kind) {
                      case "ok": {
                        const m = $can_m9;
                        const $can_m10: { $can_kind: "ok"; value: bigint } = std__int__pow_from(10n, m.value);
                        switch ($can_m10.$can_kind) {
                        case "ok": {
                          const p10 = $can_m10;
                          const $can_m11: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(($canDivMod(an.value, g.value)[0] * $canDivMod(p10.value, $canDivMod(ad.value, g.value)[0])[0]));
                          switch ($can_m11.$can_kind) {
                          case "ok": {
                            const qi = $can_m11;
                            const $can_m12: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (-m.value));
                            switch ($can_m12.$can_kind) {
                            case "ok": {
                              const q = $can_m12;
                              const $can_m13 = ((n.coefficient * pns.value) < 0n);
                              const $can_m14 = ((d.coefficient * pds.value) < 0n);
                              if ($can_m13 && $can_m14) {
                                return { $can_kind: "ok", value: q.value };
                              }
                              else if ($can_m13 && !($can_m14)) {
                                return { $can_kind: "ok", value: $canDecMul(q.value, "-1.0") };
                              }
                              else if (!($can_m13) && $can_m14) {
                                return { $can_kind: "ok", value: $canDecMul(q.value, "-1.0") };
                              }
                              else {
                                return { $can_kind: "ok", value: q.value };
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
                      return { $can_kind: "math.nonterminating_decimal", dividend: dividend, divisor: divisor };
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
export function std__dec__round_half_even(value: string, scale: bigint): { $can_kind: "ok"; discarded: string; value: string } | { $can_kind: "math.negative_scale"; scale: bigint } {
  if ((scale >= 0n)) {
    const $can_m1: { $can_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
    switch ($can_m1.$can_kind) {
    case "ok": {
      const p = $can_m1;
      if ((p.scale <= scale)) {
        return { $can_kind: "ok", value: value, discarded: "0.0" };
      }
      else {
        const $can_m2: { $can_kind: "ok"; value: bigint } = std__int__pow_from(10n, (p.scale - scale));
        switch ($can_m2.$can_kind) {
        case "ok": {
          const pd = $can_m2;
          const $can_m3: { $can_kind: "ok"; value: bigint } = std__int__abs(p.coefficient);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const a = $can_m3;
            const $can_m4 = (($canDivMod(a.value, pd.value)[1] * 2n) === pd.value);
            const $can_m5 = ($canDivMod($canDivMod(a.value, pd.value)[0], 2n)[1] === 0n);
            const $can_m6 = (($canDivMod(a.value, pd.value)[1] * 2n) > pd.value);
            if ($can_m4 && $can_m5) {
              const $can_m7: { $can_kind: "ok"; value: string } = std__convert__int_to_dec($canDivMod(a.value, pd.value)[0]);
              switch ($can_m7.$can_kind) {
              case "ok": {
                const qi = $can_m7;
                const $can_m8: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($can_m8.$can_kind) {
                case "ok": {
                  const q = $can_m8;
                  if ($canDecLt(value, "0.0")) {
                    return { $can_kind: "ok", value: $canDecMul(q.value, "-1.0"), discarded: $canDecSub(value, $canDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $can_kind: "ok", value: q.value, discarded: $canDecSub(value, q.value) };
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
            else if ($can_m4 && !($can_m5)) {
              const $can_m9: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(($canDivMod(a.value, pd.value)[0] + 1n));
              switch ($can_m9.$can_kind) {
              case "ok": {
                const qi = $can_m9;
                const $can_m10: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($can_m10.$can_kind) {
                case "ok": {
                  const q = $can_m10;
                  if ($canDecLt(value, "0.0")) {
                    return { $can_kind: "ok", value: $canDecMul(q.value, "-1.0"), discarded: $canDecSub(value, $canDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $can_kind: "ok", value: q.value, discarded: $canDecSub(value, q.value) };
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
            else if (!($can_m4) && $can_m6) {
              const $can_m11: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(($canDivMod(a.value, pd.value)[0] + 1n));
              switch ($can_m11.$can_kind) {
              case "ok": {
                const qi = $can_m11;
                const $can_m12: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($can_m12.$can_kind) {
                case "ok": {
                  const q = $can_m12;
                  if ($canDecLt(value, "0.0")) {
                    return { $can_kind: "ok", value: $canDecMul(q.value, "-1.0"), discarded: $canDecSub(value, $canDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $can_kind: "ok", value: q.value, discarded: $canDecSub(value, q.value) };
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
              const $can_m13: { $can_kind: "ok"; value: string } = std__convert__int_to_dec($canDivMod(a.value, pd.value)[0]);
              switch ($can_m13.$can_kind) {
              case "ok": {
                const qi = $can_m13;
                const $can_m14: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (0n - scale));
                switch ($can_m14.$can_kind) {
                case "ok": {
                  const q = $can_m14;
                  if ($canDecLt(value, "0.0")) {
                    return { $can_kind: "ok", value: $canDecMul(q.value, "-1.0"), discarded: $canDecSub(value, $canDecMul(q.value, "-1.0")) };
                  }
                  else {
                    return { $can_kind: "ok", value: q.value, discarded: $canDecSub(value, q.value) };
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
    return { $can_kind: "math.negative_scale", scale: scale };
  }
}
export function std__dec__divide_round_half_even_result(qm: bigint, scale: bigint, num: bigint, den: bigint, dividend: string, divisor: string): { $can_kind: "ok"; remainder: string; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(qm);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const qi = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (-scale));
    switch ($can_m2.$can_kind) {
    case "ok": {
      const q = $can_m2;
      const $can_m3 = (num < 0n);
      const $can_m4 = (den < 0n);
      if ($can_m3 && $can_m4) {
        return { $can_kind: "ok", value: q.value, remainder: $canDecSub(dividend, $canDecMul(divisor, q.value)) };
      }
      else if ($can_m3 && !($can_m4)) {
        return { $can_kind: "ok", value: $canDecMul(q.value, "-1.0"), remainder: $canDecSub(dividend, $canDecMul(divisor, $canDecMul(q.value, "-1.0"))) };
      }
      else if (!($can_m3) && $can_m4) {
        return { $can_kind: "ok", value: $canDecMul(q.value, "-1.0"), remainder: $canDecSub(dividend, $canDecMul(divisor, $canDecMul(q.value, "-1.0"))) };
      }
      else {
        return { $can_kind: "ok", value: q.value, remainder: $canDecSub(dividend, $canDecMul(divisor, q.value)) };
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
export function std__dec__divide_round_half_even(dividend: string, divisor: string, scale: bigint): { $can_kind: "ok"; remainder: string; value: string } | { $can_kind: "math.dec_zero_divisor"; divisor: string } | { $can_kind: "math.negative_scale"; scale: bigint } {
  const $can_m1: { $can_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(divisor);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const d = $can_m1;
    if ((d.coefficient === 0n)) {
      return { $can_kind: "math.dec_zero_divisor", divisor: divisor };
    }
    else {
      if ((scale >= 0n)) {
        const $can_m2: { $can_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(dividend);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const n = $can_m2;
          const $can_m3: { $can_kind: "ok"; value: bigint } = std__int__pow_from(10n, d.scale);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const pds = $can_m3;
            const $can_m4: { $can_kind: "ok"; value: bigint } = std__int__pow_from(10n, n.scale);
            switch ($can_m4.$can_kind) {
            case "ok": {
              const pns = $can_m4;
              const $can_m5: { $can_kind: "ok"; value: bigint } = std__int__abs((n.coefficient * pns.value));
              switch ($can_m5.$can_kind) {
              case "ok": {
                const an = $can_m5;
                const $can_m6: { $can_kind: "ok"; value: bigint } = std__int__abs((d.coefficient * pds.value));
                switch ($can_m6.$can_kind) {
                case "ok": {
                  const ad = $can_m6;
                  const $can_m7: { $can_kind: "ok"; value: bigint } = std__int__pow_from(10n, scale);
                  switch ($can_m7.$can_kind) {
                  case "ok": {
                    const psc = $can_m7;
                    const $can_m8 = (($canDivMod((an.value * psc.value), ad.value)[1] * 2n) === ad.value);
                    const $can_m9 = ($canDivMod($canDivMod((an.value * psc.value), ad.value)[0], 2n)[1] === 0n);
                    const $can_m10 = (($canDivMod((an.value * psc.value), ad.value)[1] * 2n) > ad.value);
                    if ($can_m8 && $can_m9) {
                      const $can_m11: { $can_kind: "ok"; remainder: string; value: string } = std__dec__divide_round_half_even_result($canDivMod((an.value * psc.value), ad.value)[0], scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($can_m11.$can_kind) {
                      case "ok": {
                        const r = $can_m11;
                        return { $can_kind: "ok", value: r.value, remainder: r.remainder };
                      }
                      default: {
                        throw new Error("unreachable");
                      }
                      }
                    }
                    else if ($can_m8 && !($can_m9)) {
                      const $can_m12: { $can_kind: "ok"; remainder: string; value: string } = std__dec__divide_round_half_even_result(($canDivMod((an.value * psc.value), ad.value)[0] + 1n), scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($can_m12.$can_kind) {
                      case "ok": {
                        const r = $can_m12;
                        return { $can_kind: "ok", value: r.value, remainder: r.remainder };
                      }
                      default: {
                        throw new Error("unreachable");
                      }
                      }
                    }
                    else if (!($can_m8) && $can_m10) {
                      const $can_m13: { $can_kind: "ok"; remainder: string; value: string } = std__dec__divide_round_half_even_result(($canDivMod((an.value * psc.value), ad.value)[0] + 1n), scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($can_m13.$can_kind) {
                      case "ok": {
                        const r = $can_m13;
                        return { $can_kind: "ok", value: r.value, remainder: r.remainder };
                      }
                      default: {
                        throw new Error("unreachable");
                      }
                      }
                    }
                    else {
                      const $can_m14: { $can_kind: "ok"; remainder: string; value: string } = std__dec__divide_round_half_even_result($canDivMod((an.value * psc.value), ad.value)[0], scale, (n.coefficient * pns.value), (d.coefficient * pds.value), dividend, divisor);
                      switch ($can_m14.$can_kind) {
                      case "ok": {
                        const r = $can_m14;
                        return { $can_kind: "ok", value: r.value, remainder: r.remainder };
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
        return { $can_kind: "math.negative_scale", scale: scale };
      }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__convert__dec_to_str(value: string): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const p = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } = std__int__abs(p.coefficient);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const a = $can_m2;
      const $can_m3: { $can_kind: "ok"; value: bigint } = std__int__pow_from(10n, p.scale);
      switch ($can_m3.$can_kind) {
      case "ok": {
        const p10 = $can_m3;
        const $can_m4: { $can_kind: "ok"; value: string } = std__convert__int_to_str($canDivMod(a.value, p10.value)[0]);
        switch ($can_m4.$can_kind) {
        case "ok": {
          const ip = $can_m4;
          const $can_m5: { $can_kind: "ok"; value: string } = std__convert__int_to_str(($canDivMod(a.value, p10.value)[1] + p10.value));
          switch ($can_m5.$can_kind) {
          case "ok": {
            const fp = $can_m5;
            if ((p.coefficient < 0n)) {
              return { $can_kind: "ok", value: ((("-" + ip.value) + ".") + $canStrSlice(fp.value, 1n, (p.scale + 1n))) };
            }
            else {
              return { $can_kind: "ok", value: ((ip.value + ".") + $canStrSlice(fp.value, 1n, (p.scale + 1n))) };
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
export function std__convert__dec_to_int_exact(value: string): { $can_kind: "ok"; value: bigint } | { $can_kind: "convert.fractional_value"; value: string } {
  const $can_m1: { $can_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const p = $can_m1;
    const $can_m2: { $can_kind: "ok"; exact: boolean; stripped: bigint } = std__dec__strip_from(p.scale, p.coefficient, true);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const s = $can_m2;
      if (s.exact) {
        return { $can_kind: "ok", value: s.stripped };
      }
      else {
        return { $can_kind: "convert.fractional_value", value: value };
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
export function std__dec__sqrt_exact(value: string): { $can_kind: "ok"; value: string } | { $can_kind: "math.dec_negative_input"; value: string } | { $can_kind: "math.nonrepresentable_result"; value: string } {
  const $can_m1: { $can_kind: "ok"; coefficient: bigint; scale: bigint } = std__dec__parts(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const p = $can_m1;
    if ((p.coefficient < 0n)) {
      return { $can_kind: "math.dec_negative_input", value: value };
    }
    else {
      if (($canDivMod(p.scale, 2n)[1] === 0n)) {
        const $can_m2: { $can_kind: "ok"; value: bigint } = std__int__sqrt_floor_search(p.coefficient, 0n, (p.coefficient + 1n));
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          if (((p.coefficient - (r.value * r.value)) === 0n)) {
            const $can_m3: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(r.value);
            switch ($can_m3.$can_kind) {
            case "ok": {
              const qi = $can_m3;
              const $can_m4: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (-$canDivMod(p.scale, 2n)[0]));
              switch ($can_m4.$can_kind) {
              case "ok": {
                const q = $can_m4;
                return { $can_kind: "ok", value: q.value };
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
            return { $can_kind: "math.nonrepresentable_result", value: value };
          }
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $can_m5: { $can_kind: "ok"; value: bigint } = std__int__sqrt_floor_search((p.coefficient * 10n), 0n, ((p.coefficient * 10n) + 1n));
        switch ($can_m5.$can_kind) {
        case "ok": {
          const r = $can_m5;
          if ((((p.coefficient * 10n) - (r.value * r.value)) === 0n)) {
            const $can_m6: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(r.value);
            switch ($can_m6.$can_kind) {
            case "ok": {
              const qi = $can_m6;
              const $can_m7: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(qi.value, (-$canDivMod((p.scale + 1n), 2n)[0]));
              switch ($can_m7.$can_kind) {
              case "ok": {
                const q = $can_m7;
                return { $can_kind: "ok", value: q.value };
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
            return { $can_kind: "math.nonrepresentable_result", value: value };
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
export function std__convert__bool_to_str(value: boolean): { $can_kind: "ok"; value: string } {
  if (value) {
    return { $can_kind: "ok", value: "true" };
  }
  else {
    return { $can_kind: "ok", value: "false" };
  }
}
export function std__convert__str_to_bool(value: string): { $can_kind: "ok"; value: boolean } | { $can_kind: "convert.invalid_boolean"; value: string } {
  if (value === "true") {
    return { $can_kind: "ok", value: true };
  }
  else if (value === "false") {
    return { $can_kind: "ok", value: false };
  }
  return { $can_kind: "convert.invalid_boolean", value: value };
}
export function std__convert__bool_to_int(value: boolean): { $can_kind: "ok"; value: bigint } {
  if (value) {
    return { $can_kind: "ok", value: 1n };
  }
  else {
    return { $can_kind: "ok", value: 0n };
  }
}
export function std__convert__int_to_bool(value: bigint): { $can_kind: "ok"; value: boolean } | { $can_kind: "convert.invalid_boolean_encoding"; value: bigint } {
  const $can_m1 = (value === 0n);
  const $can_m2 = (value === 1n);
  if ($can_m1) {
    return { $can_kind: "ok", value: false };
  }
  else if ($can_m2) {
    return { $can_kind: "ok", value: true };
  }
  else {
    return { $can_kind: "convert.invalid_boolean_encoding", value: value };
  }
}
export function std__convert__bool_to_dec(value: boolean): { $can_kind: "ok"; value: string } {
  if (value) {
    return { $can_kind: "ok", value: "1.0" };
  }
  else {
    return { $can_kind: "ok", value: "0.0" };
  }
}
export function std__convert__dec_to_bool(value: string): { $can_kind: "ok"; value: boolean } | { $can_kind: "convert.invalid_dec_encoding"; value: string } {
  if ((value === "0.0")) {
    return { $can_kind: "ok", value: false };
  }
  else {
    if ((value === "1.0")) {
      return { $can_kind: "ok", value: true };
    }
    else {
      return { $can_kind: "convert.invalid_dec_encoding", value: value };
    }
  }
}
export function std__convert__int_to_dec_from(m: bigint, acc: string, step: string): { $can_kind: "ok"; value: string } {
  if ((m <= 0n)) {
    return { $can_kind: "ok", value: acc };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: string } = std__convert__int_to_dec_from((m - 1n), $canDecAdd(acc, step), step);
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
export function std__convert__int_to_str_from(n: bigint, acc: string, m: bigint): { $can_kind: "ok"; value: string } {
  if ((m <= 0n)) {
    return { $can_kind: "ok", value: acc };
  }
  else {
    if ((n <= 9n)) {
      return { $can_kind: "ok", value: ($canStrSlice("0123456789", n, (n + 1n)) + acc) };
    }
    else {
      const $can_m1: { $can_kind: "ok"; value: string } = std__convert__int_to_str_from($canDivMod(n, 10n)[0], ($canStrSlice("0123456789", $canDivMod(n, 10n)[1], ($canDivMod(n, 10n)[1] + 1n)) + acc), (m - 1n));
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
}
export function std__convert__int_to_str(value: bigint): { $can_kind: "ok"; value: string } {
  if ((value < 0n)) {
    const $can_m1: { $can_kind: "ok"; value: string } = std__convert__int_to_str_from((-value), "", ((-value) + 1n));
    switch ($can_m1.$can_kind) {
    case "ok": {
      const r = $can_m1;
      return { $can_kind: "ok", value: ("-" + r.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  else {
    const $can_m2: { $can_kind: "ok"; value: string } = std__convert__int_to_str_from(value, "", (value + 1n));
    switch ($can_m2.$can_kind) {
    case "ok": {
      const r = $can_m2;
      return { $can_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__convert__str_to_int_from(orig: string, value: string, pos: bigint, acc: bigint, n: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "convert.invalid_integer"; value: string } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: acc };
  }
  else {
    if (($canStrAt(value, pos) >= 48n)) {
      if (($canStrAt(value, pos) <= 57n)) {
        const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "convert.invalid_integer"; value: string } = std__convert__str_to_int_from(orig, value, (pos + 1n), ((acc * 10n) + ($canStrAt(value, pos) - 48n)), (n - 1n));
        switch ($can_m1.$can_kind) {
        case "convert.invalid_integer": {
          const e = $can_m1;
          return { $can_kind: "convert.invalid_integer", value: e.value };
        }
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        return { $can_kind: "convert.invalid_integer", value: orig };
      }
    }
    else {
      return { $can_kind: "convert.invalid_integer", value: orig };
    }
  }
}
export function std__convert__str_to_int(value: string): { $can_kind: "ok"; value: bigint } | { $can_kind: "convert.invalid_integer"; value: string } {
  if (((BigInt([...value].length)) === 0n)) {
    return { $can_kind: "convert.invalid_integer", value: value };
  }
  else {
    if (($canStrSlice(value, 0n, 1n) === "-")) {
      if (((BigInt([...value].length)) >= 2n)) {
        const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "convert.invalid_integer"; value: string } = std__convert__str_to_int_from(value, value, 1n, 0n, ((BigInt([...value].length)) - 1n));
        switch ($can_m1.$can_kind) {
        case "convert.invalid_integer": {
          const e = $can_m1;
          return { $can_kind: "convert.invalid_integer", value: e.value };
        }
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", value: (-r.value) };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        return { $can_kind: "convert.invalid_integer", value: value };
      }
    }
    else {
      const $can_m2: { $can_kind: "ok"; value: bigint } | { $can_kind: "convert.invalid_integer"; value: string } = std__convert__str_to_int_from(value, value, 0n, 0n, (BigInt([...value].length)));
      switch ($can_m2.$can_kind) {
      case "convert.invalid_integer": {
        const e = $can_m2;
        return { $can_kind: "convert.invalid_integer", value: e.value };
      }
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", value: r.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
}
export function std__convert__str_to_dec_from(orig: string, value: string, pos: bigint, coeff: bigint, scale: bigint, dot: bigint, n: bigint): { $can_kind: "ok"; value: string } | { $can_kind: "convert.invalid_decimal"; value: string } {
  if ((n <= 0n)) {
    if ((dot === 0n)) {
      const $can_m1: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(coeff);
      switch ($can_m1.$can_kind) {
      case "ok": {
        const whole = $can_m1;
        return { $can_kind: "ok", value: whole.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      if ((scale === 0n)) {
        return { $can_kind: "convert.invalid_decimal", value: orig };
      }
      else {
        const $can_m2: { $can_kind: "ok"; value: string } = std__convert__int_to_dec(coeff);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const whole = $can_m2;
          const $can_m3: { $can_kind: "ok"; value: string } = std__dec__scale_by_power_of_ten(whole.value, (-scale));
          switch ($can_m3.$can_kind) {
          case "ok": {
            const moved = $can_m3;
            return { $can_kind: "ok", value: moved.value };
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
    if (($canStrSlice(value, pos, (pos + 1n)) === ".")) {
      if ((dot === 0n)) {
        const $can_m4: { $can_kind: "ok"; value: string } | { $can_kind: "convert.invalid_decimal"; value: string } = std__convert__str_to_dec_from(orig, value, (pos + 1n), coeff, scale, 1n, (n - 1n));
        switch ($can_m4.$can_kind) {
        case "convert.invalid_decimal": {
          const e = $can_m4;
          return { $can_kind: "convert.invalid_decimal", value: e.value };
        }
        case "ok": {
          const r = $can_m4;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        return { $can_kind: "convert.invalid_decimal", value: orig };
      }
    }
    else {
      if (($canStrAt(value, pos) >= 48n)) {
        if (($canStrAt(value, pos) <= 57n)) {
          if ((dot === 0n)) {
            const $can_m5: { $can_kind: "ok"; value: string } | { $can_kind: "convert.invalid_decimal"; value: string } = std__convert__str_to_dec_from(orig, value, (pos + 1n), ((coeff * 10n) + ($canStrAt(value, pos) - 48n)), scale, dot, (n - 1n));
            switch ($can_m5.$can_kind) {
            case "convert.invalid_decimal": {
              const e = $can_m5;
              return { $can_kind: "convert.invalid_decimal", value: e.value };
            }
            case "ok": {
              const r = $can_m5;
              return { $can_kind: "ok", value: r.value };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          else {
            const $can_m6: { $can_kind: "ok"; value: string } | { $can_kind: "convert.invalid_decimal"; value: string } = std__convert__str_to_dec_from(orig, value, (pos + 1n), ((coeff * 10n) + ($canStrAt(value, pos) - 48n)), (scale + 1n), dot, (n - 1n));
            switch ($can_m6.$can_kind) {
            case "convert.invalid_decimal": {
              const e = $can_m6;
              return { $can_kind: "convert.invalid_decimal", value: e.value };
            }
            case "ok": {
              const r = $can_m6;
              return { $can_kind: "ok", value: r.value };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
        }
        else {
          return { $can_kind: "convert.invalid_decimal", value: orig };
        }
      }
      else {
        return { $can_kind: "convert.invalid_decimal", value: orig };
      }
    }
  }
}
export function std__convert__str_to_dec(value: string): { $can_kind: "ok"; value: string } | { $can_kind: "convert.invalid_decimal"; value: string } {
  if (((BigInt([...value].length)) === 0n)) {
    return { $can_kind: "convert.invalid_decimal", value: value };
  }
  else {
    if (($canStrSlice(value, 0n, 1n) === "-")) {
      if (((BigInt([...value].length)) >= 2n)) {
        const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "convert.invalid_decimal"; value: string } = std__convert__str_to_dec_from(value, value, 1n, 0n, 0n, 0n, ((BigInt([...value].length)) - 1n));
        switch ($can_m1.$can_kind) {
        case "convert.invalid_decimal": {
          const e = $can_m1;
          return { $can_kind: "convert.invalid_decimal", value: e.value };
        }
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", value: $canDecSub("0.0", r.value) };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        return { $can_kind: "convert.invalid_decimal", value: value };
      }
    }
    else {
      const $can_m2: { $can_kind: "ok"; value: string } | { $can_kind: "convert.invalid_decimal"; value: string } = std__convert__str_to_dec_from(value, value, 0n, 0n, 0n, 0n, (BigInt([...value].length)));
      switch ($can_m2.$can_kind) {
      case "convert.invalid_decimal": {
        const e = $can_m2;
        return { $can_kind: "convert.invalid_decimal", value: e.value };
      }
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", value: r.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
}
export function std__convert__int_to_dec(value: bigint): { $can_kind: "ok"; value: string } {
  if ((value >= 0n)) {
    const $can_m1: { $can_kind: "ok"; value: string } = std__convert__int_to_dec_from(value, "0.0", "1.0");
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
  else {
    const $can_m2: { $can_kind: "ok"; value: string } = std__convert__int_to_dec_from((-value), "0.0", "-1.0");
    switch ($can_m2.$can_kind) {
    case "ok": {
      const r = $can_m2;
      return { $can_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__backoff__delay(attempt: bigint, base_ms: bigint, cap_ms: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "math.negative_input"; value: bigint } {
  const $can_m1 = (base_ms >= 0n);
  const $can_m2 = (cap_ms >= 0n);
  const $can_m3 = (attempt >= 0n);
  const $can_m4 = (attempt <= 30n);
  if ($can_m1 && $can_m2 && $can_m3 && $can_m4) {
    const $can_m5: { $can_kind: "ok"; value: bigint } = std__int__pow_from(2n, attempt);
    switch ($can_m5.$can_kind) {
    case "ok": {
      const p = $can_m5;
      if (((base_ms * p.value) <= cap_ms)) {
        return { $can_kind: "ok", value: (base_ms * p.value) };
      }
      else {
        return { $can_kind: "ok", value: cap_ms };
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  else if ($can_m1 && $can_m2 && $can_m3 && !($can_m4)) {
    return { $can_kind: "ok", value: cap_ms };
  }
  else if ($can_m1 && $can_m2 && !($can_m3)) {
    return { $can_kind: "math.negative_input", value: attempt };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "math.negative_input", value: cap_ms };
  }
  else {
    return { $can_kind: "math.negative_input", value: base_ms };
  }
}
