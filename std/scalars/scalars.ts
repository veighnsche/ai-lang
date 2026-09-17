// GENERATED from scalars.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type ScalarsResult = { kind: "ok"; value: boolean } | { kind: "math.invalid_bounds"; lower: bigint; upper: bigint } | { kind: "math.out_of_range"; value: bigint; lower: bigint; upper: bigint } | { kind: "math.dec_invalid_bounds"; lower: string; upper: string } | { kind: "math.dec_out_of_range"; value: string; lower: string; upper: string } | { kind: "math.negative_input"; value: bigint } | { kind: "math.negative_exponent"; exponent: bigint } | { kind: "convert.invalid_boolean"; value: string } | { kind: "convert.invalid_boolean_encoding"; value: bigint } | { kind: "convert.invalid_dec_encoding"; value: string };
export type Int__Value = { value: bigint };
export type Dec__Value = { value: string };
export type Str__Value = { value: string };
export type Bool__Value = { value: boolean };
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
export function std__bool__not(value: boolean): ScalarsResult {
  if (value) {
    return { kind: "ok", value: false };
  }
  else {
    return { kind: "ok", value: true };
  }
}
export function std__bool__and(left: boolean, right: boolean): ScalarsResult {
  if (left) {
    if (right) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "ok", value: false };
    }
  }
  else {
    return { kind: "ok", value: false };
  }
}
export function std__bool__or(left: boolean, right: boolean): ScalarsResult {
  if (left) {
    return { kind: "ok", value: true };
  }
  else {
    if (right) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "ok", value: false };
    }
  }
}
export function std__bool__xor(left: boolean, right: boolean): ScalarsResult {
  if (left) {
    if (right) {
      return { kind: "ok", value: false };
    }
    else {
      return { kind: "ok", value: true };
    }
  }
  else {
    if (right) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "ok", value: false };
    }
  }
}
export function std__bool__implies(left: boolean, right: boolean): ScalarsResult {
  if (left) {
    if (right) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "ok", value: false };
    }
  }
  else {
    return { kind: "ok", value: true };
  }
}
export function std__bool__equivalent(left: boolean, right: boolean): ScalarsResult {
  if (left) {
    if (right) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "ok", value: false };
    }
  }
  else {
    if (right) {
      return { kind: "ok", value: false };
    }
    else {
      return { kind: "ok", value: true };
    }
  }
}
export function std__compare__int(left: bigint, right: bigint): ScalarsResult {
  if ((left === right)) {
    return { kind: "ok", value: 0n };
  }
  else {
    if ((left >= right)) {
      return { kind: "ok", value: 1n };
    }
    else {
      return { kind: "ok", value: -1n };
    }
  }
}
export function std__compare__dec(left: string, right: string): ScalarsResult {
  if ((left === right)) {
    return { kind: "ok", value: 0n };
  }
  else {
    if ($ailDecGe(left, right)) {
      return { kind: "ok", value: 1n };
    }
    else {
      return { kind: "ok", value: -1n };
    }
  }
}
export function std__compare__str(left: string, right: string): ScalarsResult {
  if ((left === right)) {
    return { kind: "ok", value: 0n };
  }
  else {
    if ((left >= right)) {
      return { kind: "ok", value: 1n };
    }
    else {
      return { kind: "ok", value: -1n };
    }
  }
}
export function std__compare__bool(left: boolean, right: boolean): ScalarsResult {
  if ((left === right)) {
    return { kind: "ok", value: 0n };
  }
  else {
    if (left) {
      return { kind: "ok", value: 1n };
    }
    else {
      return { kind: "ok", value: -1n };
    }
  }
}
export function std__select__int(condition: boolean, when_true: bigint, when_false: bigint): ScalarsResult {
  if (condition) {
    return { kind: "ok", value: when_true };
  }
  else {
    return { kind: "ok", value: when_false };
  }
}
export function std__select__bool(condition: boolean, when_true: boolean, when_false: boolean): ScalarsResult {
  if (condition) {
    return { kind: "ok", value: when_true };
  }
  else {
    return { kind: "ok", value: when_false };
  }
}
export function std__select__str(condition: boolean, when_true: string, when_false: string): ScalarsResult {
  if (condition) {
    return { kind: "ok", value: when_true };
  }
  else {
    return { kind: "ok", value: when_false };
  }
}
export function std__select__dec(condition: boolean, when_true: string, when_false: string): ScalarsResult {
  if (condition) {
    return { kind: "ok", value: when_true };
  }
  else {
    return { kind: "ok", value: when_false };
  }
}
export function std__int__is_zero(value: bigint): ScalarsResult {
  if ((value === 0n)) {
    return { kind: "ok", value: true };
  }
  else {
    return { kind: "ok", value: false };
  }
}
export function std__int__is_positive(value: bigint): ScalarsResult {
  if ((value === 0n)) {
    return { kind: "ok", value: false };
  }
  else {
    if ((value >= 0n)) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "ok", value: false };
    }
  }
}
export function std__int__is_negative(value: bigint): ScalarsResult {
  if ((value === 0n)) {
    return { kind: "ok", value: false };
  }
  else {
    if ((value <= 0n)) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "ok", value: false };
    }
  }
}
export function std__int__in_closed_range(value: bigint, lower: bigint, upper: bigint): ScalarsResult {
  if ((lower <= upper)) {
    if ((value >= lower)) {
      if ((value <= upper)) {
        return { kind: "ok", value: true };
      }
      else {
        return { kind: "ok", value: false };
      }
    }
    else {
      return { kind: "ok", value: false };
    }
  }
  else {
    return { kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__in_open_range(value: bigint, lower: bigint, upper: bigint): ScalarsResult {
  if ((lower <= upper)) {
    if ((value <= lower)) {
      return { kind: "ok", value: false };
    }
    else {
      if ((value >= upper)) {
        return { kind: "ok", value: false };
      }
      else {
        return { kind: "ok", value: true };
      }
    }
  }
  else {
    return { kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__is_zero(value: string): ScalarsResult {
  if ((value === "0.0")) {
    return { kind: "ok", value: true };
  }
  else {
    return { kind: "ok", value: false };
  }
}
export function std__dec__is_positive(value: string): ScalarsResult {
  if ((value === "0.0")) {
    return { kind: "ok", value: false };
  }
  else {
    if ($ailDecGe(value, "0.0")) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "ok", value: false };
    }
  }
}
export function std__dec__is_negative(value: string): ScalarsResult {
  if ((value === "0.0")) {
    return { kind: "ok", value: false };
  }
  else {
    if ($ailDecLe(value, "0.0")) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "ok", value: false };
    }
  }
}
export function std__dec__in_closed_range(value: string, lower: string, upper: string): ScalarsResult {
  if ($ailDecLe(lower, upper)) {
    if ($ailDecGe(value, lower)) {
      if ($ailDecLe(value, upper)) {
        return { kind: "ok", value: true };
      }
      else {
        return { kind: "ok", value: false };
      }
    }
    else {
      return { kind: "ok", value: false };
    }
  }
  else {
    return { kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__in_open_range(value: string, lower: string, upper: string): ScalarsResult {
  if ($ailDecLe(lower, upper)) {
    if ($ailDecLe(value, lower)) {
      return { kind: "ok", value: false };
    }
    else {
      if ($ailDecGe(value, upper)) {
        return { kind: "ok", value: false };
      }
      else {
        return { kind: "ok", value: true };
      }
    }
  }
  else {
    return { kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__abs(value: bigint): ScalarsResult {
  if ((value >= 0n)) {
    return { kind: "ok", value: value };
  }
  else {
    return { kind: "ok", value: (0n - value) };
  }
}
export function std__int__negate(value: bigint): ScalarsResult {
  return { kind: "ok", value: (0n - value) };
}
export function std__int__sign(value: bigint): ScalarsResult {
  if ((value === 0n)) {
    return { kind: "ok", value: 0n };
  }
  else {
    if ((value >= 0n)) {
      return { kind: "ok", value: 1n };
    }
    else {
      return { kind: "ok", value: -1n };
    }
  }
}
export function std__int__min(left: bigint, right: bigint): ScalarsResult {
  if ((left <= right)) {
    return { kind: "ok", value: left };
  }
  else {
    return { kind: "ok", value: right };
  }
}
export function std__int__max(left: bigint, right: bigint): ScalarsResult {
  if ((left >= right)) {
    return { kind: "ok", value: left };
  }
  else {
    return { kind: "ok", value: right };
  }
}
export function std__int__clamp(value: bigint, lower: bigint, upper: bigint): ScalarsResult {
  if ((lower <= upper)) {
    if ((value >= lower)) {
      if ((value <= upper)) {
        return { kind: "ok", value: value };
      }
      else {
        return { kind: "ok", value: upper };
      }
    }
    else {
      return { kind: "ok", value: lower };
    }
  }
  else {
    return { kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__distance(left: bigint, right: bigint): ScalarsResult {
  if ((left >= right)) {
    return { kind: "ok", value: (left - right) };
  }
  else {
    return { kind: "ok", value: (right - left) };
  }
}
export function std__int__square(value: bigint): ScalarsResult {
  return { kind: "ok", value: (value * value) };
}
export function std__int__pow_from(base: bigint, exponent: bigint): ScalarsResult {
  if ((exponent <= 0n)) {
    return { kind: "ok", value: 1n };
  }
  else {
    const $ail_m1: ScalarsResult = std__int__pow_from(base, (exponent - 1n));
    switch ($ail_m1.kind) {
    case "ok": {
      const r = $ail_m1;
      return { kind: "ok", value: (base * r.value) };
    }
    }
  }
}
export function std__int__pow(base: bigint, exponent: bigint): ScalarsResult {
  if ((exponent >= 0n)) {
    const $ail_m2: ScalarsResult = std__int__pow_from(base, exponent);
    switch ($ail_m2.kind) {
    case "ok": {
      const r = $ail_m2;
      return { kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { kind: "math.negative_exponent", exponent: exponent };
  }
}
export function std__int__add_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): ScalarsResult {
  if ((lower <= upper)) {
    if (((left + right) >= lower)) {
      if (((left + right) <= upper)) {
        return { kind: "ok", value: (left + right) };
      }
      else {
        return { kind: "math.out_of_range", value: (left + right), lower: lower, upper: upper };
      }
    }
    else {
      return { kind: "math.out_of_range", value: (left + right), lower: lower, upper: upper };
    }
  }
  else {
    return { kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__subtract_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): ScalarsResult {
  if ((lower <= upper)) {
    if (((left - right) >= lower)) {
      if (((left - right) <= upper)) {
        return { kind: "ok", value: (left - right) };
      }
      else {
        return { kind: "math.out_of_range", value: (left - right), lower: lower, upper: upper };
      }
    }
    else {
      return { kind: "math.out_of_range", value: (left - right), lower: lower, upper: upper };
    }
  }
  else {
    return { kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__multiply_bounded(left: bigint, right: bigint, lower: bigint, upper: bigint): ScalarsResult {
  if ((lower <= upper)) {
    if (((left * right) >= lower)) {
      if (((left * right) <= upper)) {
        return { kind: "ok", value: (left * right) };
      }
      else {
        return { kind: "math.out_of_range", value: (left * right), lower: lower, upper: upper };
      }
    }
    else {
      return { kind: "math.out_of_range", value: (left * right), lower: lower, upper: upper };
    }
  }
  else {
    return { kind: "math.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__int__factorial_from(value: bigint): ScalarsResult {
  if ((value <= 0n)) {
    return { kind: "ok", value: 1n };
  }
  else {
    const $ail_m3: ScalarsResult = std__int__factorial_from((value - 1n));
    switch ($ail_m3.kind) {
    case "ok": {
      const r = $ail_m3;
      return { kind: "ok", value: (value * r.value) };
    }
    }
  }
}
export function std__int__factorial(value: bigint): ScalarsResult {
  if ((value >= 0n)) {
    const $ail_m4: ScalarsResult = std__int__factorial_from(value);
    switch ($ail_m4.kind) {
    case "ok": {
      const r = $ail_m4;
      return { kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { kind: "math.negative_input", value: value };
  }
}
export function std__int__sum_from(upper: bigint): ScalarsResult {
  if ((upper <= 0n)) {
    return { kind: "ok", value: 0n };
  }
  else {
    const $ail_m5: ScalarsResult = std__int__sum_from((upper - 1n));
    switch ($ail_m5.kind) {
    case "ok": {
      const r = $ail_m5;
      return { kind: "ok", value: (upper + r.value) };
    }
    }
  }
}
export function std__int__sum_to(upper: bigint): ScalarsResult {
  if ((upper >= 0n)) {
    const $ail_m6: ScalarsResult = std__int__sum_from(upper);
    switch ($ail_m6.kind) {
    case "ok": {
      const r = $ail_m6;
      return { kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { kind: "math.negative_input", value: upper };
  }
}
export function std__dec__abs(value: string): ScalarsResult {
  if ($ailDecGe(value, "0.0")) {
    return { kind: "ok", value: value };
  }
  else {
    return { kind: "ok", value: $ailDecSub("0.0", value) };
  }
}
export function std__dec__negate(value: string): ScalarsResult {
  return { kind: "ok", value: $ailDecSub("0.0", value) };
}
export function std__dec__sign(value: string): ScalarsResult {
  if ((value === "0.0")) {
    return { kind: "ok", value: 0n };
  }
  else {
    if ($ailDecGe(value, "0.0")) {
      return { kind: "ok", value: 1n };
    }
    else {
      return { kind: "ok", value: -1n };
    }
  }
}
export function std__dec__min(left: string, right: string): ScalarsResult {
  if ($ailDecLe(left, right)) {
    return { kind: "ok", value: left };
  }
  else {
    return { kind: "ok", value: right };
  }
}
export function std__dec__max(left: string, right: string): ScalarsResult {
  if ($ailDecGe(left, right)) {
    return { kind: "ok", value: left };
  }
  else {
    return { kind: "ok", value: right };
  }
}
export function std__dec__clamp(value: string, lower: string, upper: string): ScalarsResult {
  if ($ailDecLe(lower, upper)) {
    if ($ailDecGe(value, lower)) {
      if ($ailDecLe(value, upper)) {
        return { kind: "ok", value: value };
      }
      else {
        return { kind: "ok", value: upper };
      }
    }
    else {
      return { kind: "ok", value: lower };
    }
  }
  else {
    return { kind: "math.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__dec__distance(left: string, right: string): ScalarsResult {
  if ($ailDecGe(left, right)) {
    return { kind: "ok", value: $ailDecSub(left, right) };
  }
  else {
    return { kind: "ok", value: $ailDecSub(right, left) };
  }
}
export function std__dec__square(value: string): ScalarsResult {
  return { kind: "ok", value: $ailDecMul(value, value) };
}
export function std__dec__pow_from(base: string, exponent: bigint): ScalarsResult {
  if ((exponent <= 0n)) {
    return { kind: "ok", value: "1.0" };
  }
  else {
    const $ail_m7: ScalarsResult = std__dec__pow_from(base, (exponent - 1n));
    switch ($ail_m7.kind) {
    case "ok": {
      const r = $ail_m7;
      return { kind: "ok", value: $ailDecMul(base, r.value) };
    }
    }
  }
}
export function std__dec__pow(base: string, exponent: bigint): ScalarsResult {
  if ((exponent >= 0n)) {
    const $ail_m8: ScalarsResult = std__dec__pow_from(base, exponent);
    switch ($ail_m8.kind) {
    case "ok": {
      const r = $ail_m8;
      return { kind: "ok", value: r.value };
    }
    }
  }
  else {
    return { kind: "math.negative_exponent", exponent: exponent };
  }
}
export function std__dec__lerp(start: string, end: string, fraction: string): ScalarsResult {
  if ($ailDecGe(fraction, "0.0")) {
    if ($ailDecLe(fraction, "1.0")) {
      return { kind: "ok", value: $ailDecAdd(start, $ailDecMul($ailDecSub(end, start), fraction)) };
    }
    else {
      return { kind: "math.dec_out_of_range", value: fraction, lower: "0.0", upper: "1.0" };
    }
  }
  else {
    return { kind: "math.dec_out_of_range", value: fraction, lower: "0.0", upper: "1.0" };
  }
}
export function std__convert__bool_to_str(value: boolean): ScalarsResult {
  if (value) {
    return { kind: "ok", value: "true" };
  }
  else {
    return { kind: "ok", value: "false" };
  }
}
export function std__convert__str_to_bool(value: string): ScalarsResult {
  if (value === "true") {
    return { kind: "ok", value: true };
  }
  else if (value === "false") {
    return { kind: "ok", value: false };
  }
  return { kind: "convert.invalid_boolean", value: value };
}
export function std__convert__bool_to_int(value: boolean): ScalarsResult {
  if (value) {
    return { kind: "ok", value: 1n };
  }
  else {
    return { kind: "ok", value: 0n };
  }
}
export function std__convert__int_to_bool(value: bigint): ScalarsResult {
  if ((value === 0n)) {
    return { kind: "ok", value: false };
  }
  else {
    if ((value === 1n)) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "convert.invalid_boolean_encoding", value: value };
    }
  }
}
export function std__convert__bool_to_dec(value: boolean): ScalarsResult {
  if (value) {
    return { kind: "ok", value: "1.0" };
  }
  else {
    return { kind: "ok", value: "0.0" };
  }
}
export function std__convert__dec_to_bool(value: string): ScalarsResult {
  if ((value === "0.0")) {
    return { kind: "ok", value: false };
  }
  else {
    if ((value === "1.0")) {
      return { kind: "ok", value: true };
    }
    else {
      return { kind: "convert.invalid_dec_encoding", value: value };
    }
  }
}
export function std__backoff__delay(attempt: bigint, base_ms: bigint, cap_ms: bigint): ScalarsResult {
  if ((base_ms >= 0n)) {
    if ((cap_ms >= 0n)) {
      if ((attempt >= 0n)) {
        if ((attempt <= 30n)) {
          const $ail_m9: ScalarsResult = std__int__pow_from(2n, attempt);
          switch ($ail_m9.kind) {
          case "ok": {
            const p = $ail_m9;
            if (((base_ms * p.value) <= cap_ms)) {
              return { kind: "ok", value: (base_ms * p.value) };
            }
            else {
              return { kind: "ok", value: cap_ms };
            }
          }
          }
        }
        else {
          return { kind: "ok", value: cap_ms };
        }
      }
      else {
        return { kind: "math.negative_input", value: attempt };
      }
    }
    else {
      return { kind: "math.negative_input", value: cap_ms };
    }
  }
  else {
    return { kind: "math.negative_input", value: base_ms };
  }
}
