// GENERATED from quota.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type QuotaResult = { $can_kind: "ok" } | { $can_kind: "ok"; remaining: bigint; used: bigint } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: string } | { $can_kind: "validation.failed"; field: string; rule: string } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } | { $can_kind: "validation.not_positive"; value: bigint } | { $can_kind: "validation.negative_value"; value: bigint } | { $can_kind: "validation.dec_invalid_bounds"; lower: string; upper: string } | { $can_kind: "validation.dec_out_of_range"; value: string; lower: string; upper: string } | { $can_kind: "validation.dec_not_positive"; value: string } | { $can_kind: "validation.dec_negative_value"; value: string } | { $can_kind: "validation.invalid_length"; value: string; minimum: bigint; maximum: bigint } | { $can_kind: "validation.empty_value"; value: string } | { $can_kind: "validation.exclusive_choice" };
export type Quota__Usage = { used: bigint; remaining: bigint };
export type Validate__Pass = {};
export type Int__Value = { value: bigint };
export type Dec__Value = { value: string };
export type Str__Value = { value: string };
let Quota__used: bigint = 0n;
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
function $canDecGt(a: string, b: string): boolean {
  const A = $canDecSplit(a);
  const B = $canDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $canDecMant(A, s) > $canDecMant(B, s);
}
export function std__validate__require(condition: boolean, field: string, rule: string): { $can_kind: "ok" } | { $can_kind: "validation.failed"; field: string; rule: string } {
  if (condition) {
    return { $can_kind: "ok" };
  }
  else {
    return { $can_kind: "validation.failed", field: field, rule: rule };
  }
}
export function std__validate__int_range(value: bigint, lower: bigint, upper: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } {
  const $can_m1 = (lower <= upper);
  const $can_m2 = (value >= lower);
  const $can_m3 = (value <= upper);
  if ($can_m1 && $can_m2 && $can_m3) {
    return { $can_kind: "ok", value: value };
  }
  else if ($can_m1 && $can_m2 && !($can_m3)) {
    return { $can_kind: "validation.out_of_range", value: value, lower: lower, upper: upper };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "validation.out_of_range", value: value, lower: lower, upper: upper };
  }
  else {
    return { $can_kind: "validation.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__validate__int_nonnegative(value: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.negative_value"; value: bigint } {
  if ((value >= 0n)) {
    return { $can_kind: "ok", value: value };
  }
  else {
    return { $can_kind: "validation.negative_value", value: value };
  }
}
export function std__validate__int_positive(value: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.not_positive"; value: bigint } {
  if ((value > 0n)) {
    return { $can_kind: "ok", value: value };
  }
  else {
    return { $can_kind: "validation.not_positive", value: value };
  }
}
export function std__validate__dec_range(value: string, lower: string, upper: string): { $can_kind: "ok"; value: string } | { $can_kind: "validation.dec_invalid_bounds"; lower: string; upper: string } | { $can_kind: "validation.dec_out_of_range"; value: string; lower: string; upper: string } {
  const $can_m1 = $canDecLe(lower, upper);
  const $can_m2 = $canDecGe(value, lower);
  const $can_m3 = $canDecLe(value, upper);
  if ($can_m1 && $can_m2 && $can_m3) {
    return { $can_kind: "ok", value: value };
  }
  else if ($can_m1 && $can_m2 && !($can_m3)) {
    return { $can_kind: "validation.dec_out_of_range", value: value, lower: lower, upper: upper };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "validation.dec_out_of_range", value: value, lower: lower, upper: upper };
  }
  else {
    return { $can_kind: "validation.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__validate__dec_positive(value: string): { $can_kind: "ok"; value: string } | { $can_kind: "validation.dec_not_positive"; value: string } {
  if ($canDecGt(value, "0.0")) {
    return { $can_kind: "ok", value: value };
  }
  else {
    return { $can_kind: "validation.dec_not_positive", value: value };
  }
}
export function std__validate__dec_nonnegative(value: string): { $can_kind: "ok"; value: string } | { $can_kind: "validation.dec_negative_value"; value: string } {
  if ($canDecGe(value, "0.0")) {
    return { $can_kind: "ok", value: value };
  }
  else {
    return { $can_kind: "validation.dec_negative_value", value: value };
  }
}
export function std__validate__str_length(value: string, minimum: bigint, maximum: bigint): { $can_kind: "ok"; value: string } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.invalid_length"; value: string; minimum: bigint; maximum: bigint } {
  const $can_m1 = (minimum <= maximum);
  const $can_m2 = ((BigInt([...value].length)) >= minimum);
  const $can_m3 = ((BigInt([...value].length)) <= maximum);
  if ($can_m1 && $can_m2 && $can_m3) {
    return { $can_kind: "ok", value: value };
  }
  else if ($can_m1 && $can_m2 && !($can_m3)) {
    return { $can_kind: "validation.invalid_length", value: value, minimum: minimum, maximum: maximum };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "validation.invalid_length", value: value, minimum: minimum, maximum: maximum };
  }
  else {
    return { $can_kind: "validation.invalid_bounds", lower: minimum, upper: maximum };
  }
}
export function std__validate__str_nonempty(value: string): { $can_kind: "ok"; value: string } | { $can_kind: "validation.empty_value"; value: string } {
  if (value === "") {
    return { $can_kind: "validation.empty_value", value: value };
  }
  return { $can_kind: "ok", value: value };
}
export function std__validate__exclusive_pair(left: boolean, right: boolean): { $can_kind: "ok" } | { $can_kind: "validation.exclusive_choice" } {
  const $can_m1 = left;
  const $can_m2 = right;
  if ($can_m1 && $can_m2) {
    return { $can_kind: "validation.exclusive_choice" };
  }
  else if ($can_m1 && !($can_m2)) {
    return { $can_kind: "ok" };
  }
  else if (!($can_m1) && $can_m2) {
    return { $can_kind: "ok" };
  }
  else {
    return { $can_kind: "validation.exclusive_choice" };
  }
}
export function quota__consume(amount: bigint, quota: bigint): { $can_kind: "ok"; remaining: bigint; used: bigint } | { $can_kind: "validation.negative_value"; value: bigint } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } {
  const $can_m1: { $can_kind: "ok", value: bigint } = { $can_kind: "ok", value: Quota__used };
  switch ($can_m1.$can_kind) {
  case "ok": {
    const s = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.negative_value"; value: bigint } = std__validate__int_nonnegative(amount);
    switch ($can_m2.$can_kind) {
    case "validation.negative_value": {
      const _ = $can_m2;
      return { $can_kind: "validation.negative_value", value: amount };
    }
    case "ok": {
      const _ = $can_m2;
      const $can_m3: { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } = std__validate__int_range((s.value + amount), 0n, quota);
      switch ($can_m3.$can_kind) {
      case "validation.invalid_bounds": {
        const _ = $can_m3;
        return { $can_kind: "validation.invalid_bounds", lower: 0n, upper: quota };
      }
      case "validation.out_of_range": {
        const _ = $can_m3;
        return { $can_kind: "validation.out_of_range", value: (s.value + amount), lower: 0n, upper: quota };
      }
      case "ok": {
        const checked = $can_m3;
        Quota__used = checked.value;
        const $can_m4: { $can_kind: "ok" } = { $can_kind: "ok" };
        switch ($can_m4.$can_kind) {
        case "ok": {
          const _ = $can_m4;
          return { $can_kind: "ok", used: checked.value, remaining: (quota - checked.value) };
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
}
export function quota__usage(quota: bigint): { $can_kind: "ok"; remaining: bigint; used: bigint } | { $can_kind: "validation.negative_value"; value: bigint } {
  const $can_m1: { $can_kind: "ok", value: bigint } = { $can_kind: "ok", value: Quota__used };
  switch ($can_m1.$can_kind) {
  case "ok": {
    const s = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.negative_value"; value: bigint } = std__validate__int_nonnegative(quota);
    switch ($can_m2.$can_kind) {
    case "validation.negative_value": {
      const _ = $can_m2;
      return { $can_kind: "validation.negative_value", value: quota };
    }
    case "ok": {
      const _ = $can_m2;
      return { $can_kind: "ok", used: s.value, remaining: (quota - s.value) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  }
}
