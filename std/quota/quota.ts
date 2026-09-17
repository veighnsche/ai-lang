// GENERATED from quota.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type QuotaResult = { $ail_kind: "ok" } | { $ail_kind: "ok"; remaining: bigint; used: bigint } | { $ail_kind: "ok"; value: bigint } | { $ail_kind: "ok"; value: string } | { $ail_kind: "validation.failed"; field: string; rule: string } | { $ail_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $ail_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } | { $ail_kind: "validation.not_positive"; value: bigint } | { $ail_kind: "validation.negative_value"; value: bigint } | { $ail_kind: "validation.dec_invalid_bounds"; lower: string; upper: string } | { $ail_kind: "validation.dec_out_of_range"; value: string; lower: string; upper: string } | { $ail_kind: "validation.dec_not_positive"; value: string } | { $ail_kind: "validation.dec_negative_value"; value: string } | { $ail_kind: "validation.empty_value"; value: string } | { $ail_kind: "validation.exclusive_choice" };
export type Quota__Usage = { used: bigint; remaining: bigint };
export type Validate__Pass = {};
export type Int__Value = { value: bigint };
export type Dec__Value = { value: string };
export type Str__Value = { value: string };
let Quota__used: bigint = 0n;
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
function $ailDecGt(a: string, b: string): boolean {
  const A = $ailDecSplit(a);
  const B = $ailDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $ailDecMant(A, s) > $ailDecMant(B, s);
}
export function std__validate__require(condition: boolean, field: string, rule: string): QuotaResult {
  if (condition) {
    return { $ail_kind: "ok" };
  }
  else {
    return { $ail_kind: "validation.failed", field: field, rule: rule };
  }
}
export function std__validate__int_range(value: bigint, lower: bigint, upper: bigint): QuotaResult {
  if ((lower <= upper)) {
    if ((value >= lower)) {
      if ((value <= upper)) {
        return { $ail_kind: "ok", value: value };
      }
      else {
        return { $ail_kind: "validation.out_of_range", value: value, lower: lower, upper: upper };
      }
    }
    else {
      return { $ail_kind: "validation.out_of_range", value: value, lower: lower, upper: upper };
    }
  }
  else {
    return { $ail_kind: "validation.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__validate__int_nonnegative(value: bigint): QuotaResult {
  if ((value >= 0n)) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    return { $ail_kind: "validation.negative_value", value: value };
  }
}
export function std__validate__int_positive(value: bigint): QuotaResult {
  if ((value > 0n)) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    return { $ail_kind: "validation.not_positive", value: value };
  }
}
export function std__validate__dec_range(value: string, lower: string, upper: string): QuotaResult {
  if ($ailDecLe(lower, upper)) {
    if ($ailDecGe(value, lower)) {
      if ($ailDecLe(value, upper)) {
        return { $ail_kind: "ok", value: value };
      }
      else {
        return { $ail_kind: "validation.dec_out_of_range", value: value, lower: lower, upper: upper };
      }
    }
    else {
      return { $ail_kind: "validation.dec_out_of_range", value: value, lower: lower, upper: upper };
    }
  }
  else {
    return { $ail_kind: "validation.dec_invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__validate__dec_positive(value: string): QuotaResult {
  if ($ailDecGt(value, "0.0")) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    return { $ail_kind: "validation.dec_not_positive", value: value };
  }
}
export function std__validate__dec_nonnegative(value: string): QuotaResult {
  if ($ailDecGe(value, "0.0")) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    return { $ail_kind: "validation.dec_negative_value", value: value };
  }
}
export function std__validate__str_nonempty(value: string): QuotaResult {
  if (value === "") {
    return { $ail_kind: "validation.empty_value", value: value };
  }
  return { $ail_kind: "ok", value: value };
}
export function std__validate__exclusive_pair(left: boolean, right: boolean): QuotaResult {
  if (left) {
    if (right) {
      return { $ail_kind: "validation.exclusive_choice" };
    }
    else {
      return { $ail_kind: "ok" };
    }
  }
  else {
    if (right) {
      return { $ail_kind: "ok" };
    }
    else {
      return { $ail_kind: "validation.exclusive_choice" };
    }
  }
}
export function quota__consume(amount: bigint, quota: bigint): QuotaResult {
  const $ail_m1: { $ail_kind: "ok", value: bigint } = { $ail_kind: "ok", value: Quota__used };
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const s = $ail_m1;
    const $ail_m2: QuotaResult = std__validate__int_nonnegative(amount);
    switch ($ail_m2.$ail_kind) {
    case "validation.negative_value": {
      const _ = $ail_m2;
      return { $ail_kind: "validation.negative_value", value: amount };
    }
    case "ok": {
      const _ = $ail_m2;
      const $ail_m3: QuotaResult = std__validate__int_range((s.value + amount), 0n, quota);
      switch ($ail_m3.$ail_kind) {
      case "validation.invalid_bounds": {
        const _ = $ail_m3;
        return { $ail_kind: "validation.invalid_bounds", lower: 0n, upper: quota };
      }
      case "validation.out_of_range": {
        const _ = $ail_m3;
        return { $ail_kind: "validation.out_of_range", value: (s.value + amount), lower: 0n, upper: quota };
      }
      case "ok": {
        const checked = $ail_m3;
        Quota__used = checked.value;
        const $ail_m4: { $ail_kind: "ok" } = { $ail_kind: "ok" };
        switch ($ail_m4.$ail_kind) {
        case "ok": {
          const _ = $ail_m4;
          return { $ail_kind: "ok", used: checked.value, remaining: (quota - checked.value) };
        }
        }
      }
      }
    }
    }
  }
  }
}
export function quota__usage(quota: bigint): QuotaResult {
  const $ail_m5: { $ail_kind: "ok", value: bigint } = { $ail_kind: "ok", value: Quota__used };
  switch ($ail_m5.$ail_kind) {
  case "ok": {
    const s = $ail_m5;
    const $ail_m6: QuotaResult = std__validate__int_nonnegative(quota);
    switch ($ail_m6.$ail_kind) {
    case "validation.negative_value": {
      const _ = $ail_m6;
      return { $ail_kind: "validation.negative_value", value: quota };
    }
    case "ok": {
      const _ = $ail_m6;
      return { $ail_kind: "ok", used: s.value, remaining: (quota - s.value) };
    }
    }
  }
  }
}
