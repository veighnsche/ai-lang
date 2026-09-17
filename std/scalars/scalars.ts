// GENERATED from scalars.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type ScalarsResult = { $ail_kind: "ok"; remainder: bigint; root: bigint } | { $ail_kind: "ok"; value: bigint } | { $ail_kind: "ok"; value: boolean } | { $ail_kind: "ok"; value: string } | { $ail_kind: "math.invalid_bounds"; lower: bigint; upper: bigint } | { $ail_kind: "math.out_of_range"; value: bigint; lower: bigint; upper: bigint } | { $ail_kind: "math.dec_invalid_bounds"; lower: string; upper: string } | { $ail_kind: "math.dec_out_of_range"; value: string; lower: string; upper: string } | { $ail_kind: "math.negative_input"; value: bigint } | { $ail_kind: "math.nonpositive_input"; value: bigint } | { $ail_kind: "math.negative_exponent"; exponent: bigint } | { $ail_kind: "math.invalid_count"; n: bigint; k: bigint } | { $ail_kind: "math.negative_scale"; scale: bigint } | { $ail_kind: "convert.invalid_boolean"; value: string } | { $ail_kind: "convert.invalid_boolean_encoding"; value: bigint } | { $ail_kind: "convert.invalid_dec_encoding"; value: string };
export type Int__Value = { value: bigint };
export type Dec__Value = { value: string };
export type Str__Value = { value: string };
export type Bool__Value = { value: boolean };
export type Int__Root = { root: bigint; remainder: bigint };
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
  if (left) {
    if (right) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__bool__or(left: boolean, right: boolean): ScalarsResult {
  if (left) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    if (right) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
}
export function std__bool__xor(left: boolean, right: boolean): ScalarsResult {
  if (left) {
    if (right) {
      return { $ail_kind: "ok", value: false };
    }
    else {
      return { $ail_kind: "ok", value: true };
    }
  }
  else {
    if (right) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
}
export function std__bool__implies(left: boolean, right: boolean): ScalarsResult {
  if (left) {
    if (right) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
  else {
    return { $ail_kind: "ok", value: true };
  }
}
export function std__bool__equivalent(left: boolean, right: boolean): ScalarsResult {
  if (left) {
    if (right) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
  else {
    if (right) {
      return { $ail_kind: "ok", value: false };
    }
    else {
      return { $ail_kind: "ok", value: true };
    }
  }
}
export function std__compare__int(left: bigint, right: bigint): ScalarsResult {
  if ((left === right)) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    if ((left >= right)) {
      return { $ail_kind: "ok", value: 1n };
    }
    else {
      return { $ail_kind: "ok", value: -1n };
    }
  }
}
export function std__compare__dec(left: string, right: string): ScalarsResult {
  if ((left === right)) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    if ($ailDecGe(left, right)) {
      return { $ail_kind: "ok", value: 1n };
    }
    else {
      return { $ail_kind: "ok", value: -1n };
    }
  }
}
export function std__compare__str(left: string, right: string): ScalarsResult {
  if ((left === right)) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    if ($ailStrGe(left, right)) {
      return { $ail_kind: "ok", value: 1n };
    }
    else {
      return { $ail_kind: "ok", value: -1n };
    }
  }
}
export function std__compare__bool(left: boolean, right: boolean): ScalarsResult {
  if ((left === right)) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    if (left) {
      return { $ail_kind: "ok", value: 1n };
    }
    else {
      return { $ail_kind: "ok", value: -1n };
    }
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
  if ((value === 0n)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    if ((value >= 0n)) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
}
export function std__int__is_negative(value: bigint): ScalarsResult {
  if ((value === 0n)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    if ((value <= 0n)) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
}
export function std__int__in_closed_range(value: bigint, lower: bigint, upper: bigint): ScalarsResult {
  if ((lower <= upper)) {
    if ((value >= lower)) {
      if ((value <= upper)) {
        return { $ail_kind: "ok", value: true };
      }
      else {
        return { $ail_kind: "ok", value: false };
      }
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__in_open_range(value: bigint, lower: bigint, upper: bigint): ScalarsResult {
  if ((lower <= upper)) {
    if ((value <= lower)) {
      return { $ail_kind: "ok", value: false };
    }
    else {
      if ((value >= upper)) {
        return { $ail_kind: "ok", value: false };
      }
      else {
        return { $ail_kind: "ok", value: true };
      }
    }
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
  if ((value === "0.0")) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    if ($ailDecGe(value, "0.0")) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
}
export function std__dec__is_negative(value: string): ScalarsResult {
  if ((value === "0.0")) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    if ($ailDecLe(value, "0.0")) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
}
export function std__dec__in_closed_range(value: string, lower: string, upper: string): ScalarsResult {
  if ($ailDecLe(lower, upper)) {
    if ($ailDecGe(value, lower)) {
      if ($ailDecLe(value, upper)) {
        return { $ail_kind: "ok", value: true };
      }
      else {
        return { $ail_kind: "ok", value: false };
      }
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
  else {
    return { $ail_kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__in_open_range(value: string, lower: string, upper: string): ScalarsResult {
  if ($ailDecLe(lower, upper)) {
    if ($ailDecLe(value, lower)) {
      return { $ail_kind: "ok", value: false };
    }
    else {
      if ($ailDecGe(value, upper)) {
        return { $ail_kind: "ok", value: false };
      }
      else {
        return { $ail_kind: "ok", value: true };
      }
    }
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
  if ((value === 0n)) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    if ((value >= 0n)) {
      return { $ail_kind: "ok", value: 1n };
    }
    else {
      return { $ail_kind: "ok", value: -1n };
    }
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
  if ((lower <= upper)) {
    if ((value >= lower)) {
      if ((value <= upper)) {
        return { $ail_kind: "ok", value: value };
      }
      else {
        return { $ail_kind: "ok", value: upper };
      }
    }
    else {
      return { $ail_kind: "ok", value: lower };
    }
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
    const $ail_m1: ScalarsResult = std__int__pow_from(base, (exponent - 1n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const r = $ail_m1;
      return { $ail_kind: "ok", value: (base * r.value) };
    }
    }
  }
}
export function std__int__pow(base: bigint, exponent: bigint): ScalarsResult {
  if ((exponent >= 0n)) {
    const $ail_m2: ScalarsResult = std__int__pow_from(base, exponent);
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const r = $ail_m2;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_exponent", exponent: exponent };
  }
}
export function std__int__add_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): ScalarsResult {
  if ((lower <= upper)) {
    if (((left + right) >= lower)) {
      if (((left + right) <= upper)) {
        return { $ail_kind: "ok", value: (left + right) };
      }
      else {
        return { $ail_kind: "math.out_of_range", value: (left + right), lower: lower, upper: upper };
      }
    }
    else {
      return { $ail_kind: "math.out_of_range", value: (left + right), lower: lower, upper: upper };
    }
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__subtract_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): ScalarsResult {
  if ((lower <= upper)) {
    if (((left - right) >= lower)) {
      if (((left - right) <= upper)) {
        return { $ail_kind: "ok", value: (left - right) };
      }
      else {
        return { $ail_kind: "math.out_of_range", value: (left - right), lower: lower, upper: upper };
      }
    }
    else {
      return { $ail_kind: "math.out_of_range", value: (left - right), lower: lower, upper: upper };
    }
  }
  else {
    return { $ail_kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__multiply_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): ScalarsResult {
  if ((lower <= upper)) {
    if (((left * right) >= lower)) {
      if (((left * right) <= upper)) {
        return { $ail_kind: "ok", value: (left * right) };
      }
      else {
        return { $ail_kind: "math.out_of_range", value: (left * right), lower: lower, upper: upper };
      }
    }
    else {
      return { $ail_kind: "math.out_of_range", value: (left * right), lower: lower, upper: upper };
    }
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
    const $ail_m3: ScalarsResult = std__int__factorial_from((value - 1n));
    switch ($ail_m3.$ail_kind) {
    case "ok": {
      const r = $ail_m3;
      return { $ail_kind: "ok", value: (value * r.value) };
    }
    }
  }
}
export function std__int__factorial(value: bigint): ScalarsResult {
  if ((value >= 0n)) {
    const $ail_m4: ScalarsResult = std__int__factorial_from(value);
    switch ($ail_m4.$ail_kind) {
    case "ok": {
      const r = $ail_m4;
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
    const $ail_m5: ScalarsResult = std__int__sum_from((upper - 1n));
    switch ($ail_m5.$ail_kind) {
    case "ok": {
      const r = $ail_m5;
      return { $ail_kind: "ok", value: (upper + r.value) };
    }
    }
  }
}
export function std__int__sum_to(upper: bigint): ScalarsResult {
  if ((upper >= 0n)) {
    const $ail_m6: ScalarsResult = std__int__sum_from(upper);
    switch ($ail_m6.$ail_kind) {
    case "ok": {
      const r = $ail_m6;
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
    const $ail_m7: ScalarsResult = std__int__binomial_from(n, k, (i - 1n), (num * ((n - k) + i)), (den * i));
    switch ($ail_m7.$ail_kind) {
    case "ok": {
      const r = $ail_m7;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__int__binomial(n: bigint, k: bigint): ScalarsResult {
  if ((k >= 0n)) {
    if ((k <= n)) {
      const $ail_m8: ScalarsResult = std__int__binomial_from(n, k, k, 1n, 1n);
      switch ($ail_m8.$ail_kind) {
      case "ok": {
        const r = $ail_m8;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
    else {
      return { $ail_kind: "math.invalid_count", n: n, k: k };
    }
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
      const $ail_m9: ScalarsResult = std__int__is_prime_from(value, (e - 1n));
      switch ($ail_m9.$ail_kind) {
      case "ok": {
        const r = $ail_m9;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
  }
}
export function std__int__is_prime(value: bigint): ScalarsResult {
  if ((value >= 2n)) {
    const $ail_m10: ScalarsResult = std__int__is_prime_from(value, (value - 2n));
    switch ($ail_m10.$ail_kind) {
    case "ok": {
      const r = $ail_m10;
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
    const $ail_m11: ScalarsResult = std__int__next_power_of_two_pow2((e - 1n));
    switch ($ail_m11.$ail_kind) {
    case "ok": {
      const r = $ail_m11;
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
    const $ail_m12: ScalarsResult = std__int__next_power_of_two_pow2(e);
    switch ($ail_m12.$ail_kind) {
    case "ok": {
      const t = $ail_m12;
      if ((t.value >= value)) {
        const $ail_m13: ScalarsResult = std__int__next_power_of_two_from(value, (e - 1n), t.value);
        switch ($ail_m13.$ail_kind) {
        case "ok": {
          const r = $ail_m13;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
      else {
        const $ail_m14: ScalarsResult = std__int__next_power_of_two_from(value, (e - 1n), best);
        switch ($ail_m14.$ail_kind) {
        case "ok": {
          const r = $ail_m14;
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
    const $ail_m15: ScalarsResult = std__int__next_power_of_two_from(value, value, 1n);
    switch ($ail_m15.$ail_kind) {
    case "ok": {
      const r = $ail_m15;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { $ail_kind: "math.nonpositive_input", value: value };
  }
}
export function std__int__sqrt_floor_from(value: bigint, e: bigint): ScalarsResult {
  if ((e <= 0n)) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    if (((e * e) <= value)) {
      return { $ail_kind: "ok", value: e };
    }
    else {
      const $ail_m16: ScalarsResult = std__int__sqrt_floor_from(value, (e - 1n));
      switch ($ail_m16.$ail_kind) {
      case "ok": {
        const r = $ail_m16;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
  }
}
export function std__int__sqrt_floor(value: bigint): ScalarsResult {
  if ((value >= 0n)) {
    const $ail_m17: ScalarsResult = std__int__sqrt_floor_from(value, value);
    switch ($ail_m17.$ail_kind) {
    case "ok": {
      const r = $ail_m17;
      return { $ail_kind: "ok", root: r.value, remainder: (value - (r.value * r.value)) };
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_input", value: value };
  }
}
export function std__int__gcd_from(x: bigint, y: bigint, e: bigint): ScalarsResult {
  if ((e <= 0n)) {
    return { $ail_kind: "ok", value: 1n };
  }
  else {
    if (($ailDivMod(x, (e + 1n))[1] === 0n)) {
      if (($ailDivMod(y, (e + 1n))[1] === 0n)) {
        return { $ail_kind: "ok", value: (e + 1n) };
      }
      else {
        const $ail_m18: ScalarsResult = std__int__gcd_from(x, y, (e - 1n));
        switch ($ail_m18.$ail_kind) {
        case "ok": {
          const r = $ail_m18;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
    }
    else {
      const $ail_m19: ScalarsResult = std__int__gcd_from(x, y, (e - 1n));
      switch ($ail_m19.$ail_kind) {
      case "ok": {
        const r = $ail_m19;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
  }
}
export function std__int__gcd(left: bigint, right: bigint): ScalarsResult {
  const $ail_m20: ScalarsResult = std__int__abs(left);
  switch ($ail_m20.$ail_kind) {
  case "ok": {
    const a = $ail_m20;
    const $ail_m21: ScalarsResult = std__int__abs(right);
    switch ($ail_m21.$ail_kind) {
    case "ok": {
      const b = $ail_m21;
      if ((a.value === 0n)) {
        return { $ail_kind: "ok", value: b.value };
      }
      else {
        if ((b.value === 0n)) {
          return { $ail_kind: "ok", value: a.value };
        }
        else {
          const $ail_m22: ScalarsResult = std__int__gcd_from(a.value, b.value, (a.value - 1n));
          switch ($ail_m22.$ail_kind) {
          case "ok": {
            const r = $ail_m22;
            return { $ail_kind: "ok", value: r.value };
          }
          }
        }
      }
    }
    }
  }
  }
}
export function std__int__lcm(left: bigint, right: bigint): ScalarsResult {
  if ((left === 0n)) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    if ((right === 0n)) {
      return { $ail_kind: "ok", value: 0n };
    }
    else {
      const $ail_m23: ScalarsResult = std__int__gcd(left, right);
      switch ($ail_m23.$ail_kind) {
      case "ok": {
        const g = $ail_m23;
        const $ail_m24: ScalarsResult = std__int__abs(left);
        switch ($ail_m24.$ail_kind) {
        case "ok": {
          const a = $ail_m24;
          const $ail_m25: ScalarsResult = std__int__abs(right);
          switch ($ail_m25.$ail_kind) {
          case "ok": {
            const b = $ail_m25;
            return { $ail_kind: "ok", value: $ailDivMod((a.value * b.value), g.value)[0] };
          }
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
  if ((value === "0.0")) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    if ($ailDecGe(value, "0.0")) {
      return { $ail_kind: "ok", value: 1n };
    }
    else {
      return { $ail_kind: "ok", value: -1n };
    }
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
  if ($ailDecLe(lower, upper)) {
    if ($ailDecGe(value, lower)) {
      if ($ailDecLe(value, upper)) {
        return { $ail_kind: "ok", value: value };
      }
      else {
        return { $ail_kind: "ok", value: upper };
      }
    }
    else {
      return { $ail_kind: "ok", value: lower };
    }
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
    const $ail_m26: ScalarsResult = std__dec__pow_from(base, (exponent - 1n));
    switch ($ail_m26.$ail_kind) {
    case "ok": {
      const r = $ail_m26;
      return { $ail_kind: "ok", value: $ailDecMul(base, r.value) };
    }
    }
  }
}
export function std__dec__pow(base: string, exponent: bigint): ScalarsResult {
  if ((exponent >= 0n)) {
    const $ail_m27: ScalarsResult = std__dec__pow_from(base, exponent);
    switch ($ail_m27.$ail_kind) {
    case "ok": {
      const r = $ail_m27;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { $ail_kind: "math.negative_exponent", exponent: exponent };
  }
}
export function std__dec__lerp(start: string, end: string, fraction: string): ScalarsResult {
  if ($ailDecGe(fraction, "0.0")) {
    if ($ailDecLe(fraction, "1.0")) {
      return { $ail_kind: "ok", value: $ailDecAdd(start, $ailDecMul($ailDecSub(end, start), fraction)) };
    }
    else {
      return { $ail_kind: "math.dec_out_of_range", value: fraction, lower: "0.0", upper: "1.0" };
    }
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
    const $ail_m28: ScalarsResult = std__dec__scale_by_power_of_ten_from($ailDecMul(value, factor), (e - 1n), factor);
    switch ($ail_m28.$ail_kind) {
    case "ok": {
      const r = $ail_m28;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__dec__scale_by_power_of_ten(value: string, exponent: bigint): ScalarsResult {
  if ((exponent >= 0n)) {
    const $ail_m29: ScalarsResult = std__dec__scale_by_power_of_ten_from(value, exponent, "10.0");
    switch ($ail_m29.$ail_kind) {
    case "ok": {
      const r = $ail_m29;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    const $ail_m30: ScalarsResult = std__dec__scale_by_power_of_ten_from(value, (0n - exponent), "0.1");
    switch ($ail_m30.$ail_kind) {
    case "ok": {
      const r = $ail_m30;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__dec__from_parts(coefficient: bigint, scale: bigint): ScalarsResult {
  if ((scale >= 0n)) {
    const $ail_m31: ScalarsResult = std__convert__int_to_dec(coefficient);
    switch ($ail_m31.$ail_kind) {
    case "ok": {
      const whole = $ail_m31;
      const $ail_m32: ScalarsResult = std__dec__scale_by_power_of_ten(whole.value, (0n - scale));
      switch ($ail_m32.$ail_kind) {
      case "ok": {
        const moved = $ail_m32;
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
  if ((value === 0n)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    if ((value === 1n)) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "convert.invalid_boolean_encoding", value: value };
    }
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
    const $ail_m33: ScalarsResult = std__convert__int_to_dec_from((m - 1n), $ailDecAdd(acc, step), step);
    switch ($ail_m33.$ail_kind) {
    case "ok": {
      const r = $ail_m33;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__convert__int_to_dec(value: bigint): ScalarsResult {
  if ((value >= 0n)) {
    const $ail_m34: ScalarsResult = std__convert__int_to_dec_from(value, "0.0", "1.0");
    switch ($ail_m34.$ail_kind) {
    case "ok": {
      const r = $ail_m34;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
  else {
    const $ail_m35: ScalarsResult = std__convert__int_to_dec_from((0n - value), "0.0", "-1.0");
    switch ($ail_m35.$ail_kind) {
    case "ok": {
      const r = $ail_m35;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__backoff__delay(attempt: bigint, base_ms: bigint, cap_ms: bigint): ScalarsResult {
  if ((base_ms >= 0n)) {
    if ((cap_ms >= 0n)) {
      if ((attempt >= 0n)) {
        if ((attempt <= 30n)) {
          const $ail_m36: ScalarsResult = std__int__pow_from(2n, attempt);
          switch ($ail_m36.$ail_kind) {
          case "ok": {
            const p = $ail_m36;
            if (((base_ms * p.value) <= cap_ms)) {
              return { $ail_kind: "ok", value: (base_ms * p.value) };
            }
            else {
              return { $ail_kind: "ok", value: cap_ms };
            }
          }
          }
        }
        else {
          return { $ail_kind: "ok", value: cap_ms };
        }
      }
      else {
        return { $ail_kind: "math.negative_input", value: attempt };
      }
    }
    else {
      return { $ail_kind: "math.negative_input", value: cap_ms };
    }
  }
  else {
    return { $ail_kind: "math.negative_input", value: base_ms };
  }
}
