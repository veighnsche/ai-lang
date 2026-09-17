// GENERATED from quota.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type QuotaResult = { kind: "ok" } | { kind: "ok"; remaining: bigint; used: bigint } | { kind: "ok"; value: bigint } | { kind: "validation.failed"; field: string; rule: string } | { kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } | { kind: "validation.negative_value"; value: bigint } | { kind: "validation.empty_value"; value: string } | { kind: "validation.exclusive_choice" };
export type Quota__Usage = { used: bigint; remaining: bigint };
export type Validate__Pass = {};
export type Int__Value = { value: bigint };
export type Str__Value = { value: string };
let Quota__used: bigint = 0n;
export function std__validate__require(condition: boolean, field: string, rule: string): QuotaResult {
  if (condition) {
    return { kind: "ok" };
  }
  else {
    return { kind: "validation.failed", field: field, rule: rule };
  }
}
export function std__validate__int_range(value: bigint, lower: bigint, upper: bigint): QuotaResult {
  if ((lower <= upper)) {
    if ((value >= lower)) {
      if ((value <= upper)) {
        return { kind: "ok", value: value };
      }
      else {
        return { kind: "validation.out_of_range", value: value, lower: lower, upper: upper };
      }
    }
    else {
      return { kind: "validation.out_of_range", value: value, lower: lower, upper: upper };
    }
  }
  else {
    return { kind: "validation.invalid_bounds", lower: lower, upper: upper };
  }
}
export function std__validate__int_nonnegative(value: bigint): QuotaResult {
  if ((value >= 0n)) {
    return { kind: "ok", value: value };
  }
  else {
    return { kind: "validation.negative_value", value: value };
  }
}
export function std__validate__str_nonempty(value: string): QuotaResult {
  if (value === "") {
    return { kind: "validation.empty_value", value: value };
  }
  return { kind: "ok", value: value };
}
export function std__validate__exclusive_pair(left: boolean, right: boolean): QuotaResult {
  if (left) {
    if (right) {
      return { kind: "validation.exclusive_choice" };
    }
    else {
      return { kind: "ok" };
    }
  }
  else {
    if (right) {
      return { kind: "ok" };
    }
    else {
      return { kind: "validation.exclusive_choice" };
    }
  }
}
export function quota__consume(amount: bigint, quota: bigint): QuotaResult {
  const $ail_m1: { kind: "ok", value: bigint } = { kind: "ok", value: Quota__used };
  switch ($ail_m1.kind) {
  case "ok": {
    const s = $ail_m1;
    const $ail_m2: QuotaResult = std__validate__int_nonnegative(amount);
    switch ($ail_m2.kind) {
    case "validation.negative_value": {
      const _ = $ail_m2;
      return { kind: "validation.negative_value", value: amount };
    }
    case "ok": {
      const _ = $ail_m2;
      const $ail_m3: QuotaResult = std__validate__int_range((s.value + amount), 0n, quota);
      switch ($ail_m3.kind) {
      case "validation.invalid_bounds": {
        const _ = $ail_m3;
        return { kind: "validation.invalid_bounds", lower: 0n, upper: quota };
      }
      case "validation.out_of_range": {
        const _ = $ail_m3;
        return { kind: "validation.out_of_range", value: (s.value + amount), lower: 0n, upper: quota };
      }
      case "ok": {
        const checked = $ail_m3;
        Quota__used = checked.value;
        const $ail_m4: { kind: "ok" } = { kind: "ok" };
        switch ($ail_m4.kind) {
        case "ok": {
          const _ = $ail_m4;
          return { kind: "ok", used: checked.value, remaining: (quota - checked.value) };
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
  const $ail_m5: { kind: "ok", value: bigint } = { kind: "ok", value: Quota__used };
  switch ($ail_m5.kind) {
  case "ok": {
    const s = $ail_m5;
    const $ail_m6: QuotaResult = std__validate__int_nonnegative(quota);
    switch ($ail_m6.kind) {
    case "validation.negative_value": {
      const _ = $ail_m6;
      return { kind: "validation.negative_value", value: quota };
    }
    case "ok": {
      const _ = $ail_m6;
      return { kind: "ok", used: s.value, remaining: (quota - s.value) };
    }
    }
  }
  }
}
