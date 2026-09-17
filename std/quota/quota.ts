// GENERATED from quota.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type QuotaResult = { $ail_kind: "ok" } | { $ail_kind: "ok"; remaining: bigint; used: bigint } | { $ail_kind: "ok"; value: bigint } | { $ail_kind: "validation.failed"; field: string; rule: string } | { $ail_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $ail_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } | { $ail_kind: "validation.negative_value"; value: bigint } | { $ail_kind: "validation.empty_value"; value: string } | { $ail_kind: "validation.exclusive_choice" };
export type Quota__Usage = { used: bigint; remaining: bigint };
export type Validate__Pass = {};
export type Int__Value = { value: bigint };
export type Str__Value = { value: string };
let Quota__used: bigint = 0n;
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
