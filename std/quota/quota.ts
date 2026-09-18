// GENERATED from quota.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { std__convert__int_to_str, type ScalarsResult } from "./scalars";
export type QuotaResult = { $can_kind: "ok" } | { $can_kind: "ok"; remaining: bigint; used: bigint } | { $can_kind: "ok"; value: Quota__Envelope } | { $can_kind: "ok"; value: Quota__Request } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: string } | { $can_kind: "validation.failed"; field: string; rule: string } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } | { $can_kind: "validation.not_positive"; value: bigint } | { $can_kind: "validation.negative_value"; value: bigint } | { $can_kind: "validation.dec_invalid_bounds"; lower: string; upper: string } | { $can_kind: "validation.dec_out_of_range"; value: string; lower: string; upper: string } | { $can_kind: "validation.dec_not_positive"; value: string } | { $can_kind: "validation.dec_negative_value"; value: string } | { $can_kind: "validation.invalid_length"; value: string; minimum: bigint; maximum: bigint } | { $can_kind: "validation.empty_value"; value: string } | { $can_kind: "validation.exclusive_choice" } | { $can_kind: "validation.not_allowed"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string };
export type Quota__Usage = { used: bigint; remaining: bigint };
export type Validate__Pass = {};
export type Int__Value = { value: bigint };
export type Dec__Value = { value: string };
export type Str__Value = { value: string };
export type Quota__Request = { label: string; amount: bigint; mode: string };
export type Quota__RequestSchema = { label_minimum: bigint; label_maximum: bigint; amount_lower: bigint; amount_upper: bigint; allowed_modes: string[] };
export type Quota__RequestValue = { value: Quota__Request };
export type Quota__Envelope = { request: Quota__Request };
export type Quota__EnvelopeValue = { value: Quota__Envelope };
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
// Sequence indexing (a38 S3): bounds throw, matching Go.
function $canSeqAt<T>(a: T[], i: bigint): T {
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("seq index out of range");
  const k = Number(i);
  if (k >= a.length) throw new Error("seq index out of range");
  return a[k];
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
export function std__validate__str_one_of_from(value: string, allowed: string[], pos: bigint, n: bigint): { $can_kind: "ok"; value: string } | { $can_kind: "validation.not_allowed"; value: string } {
  if ((n <= 0n)) {
    return { $can_kind: "validation.not_allowed", value: value };
  }
  else {
    if ((pos < (BigInt([...allowed].length)))) {
      if (($canSeqAt(allowed, pos) === value)) {
        return { $can_kind: "ok", value: value };
      }
      else {
        const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "validation.not_allowed"; value: string } = std__validate__str_one_of_from(value, allowed, (pos + 1n), (n - 1n));
        switch ($can_m1.$can_kind) {
        case "validation.not_allowed": {
          const e = $can_m1;
          return { $can_kind: "validation.not_allowed", value: e.value };
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
    }
    else {
      return { $can_kind: "validation.not_allowed", value: value };
    }
  }
}
export function std__validate__str_one_of(value: string, allowed: string[]): { $can_kind: "ok"; value: string } | { $can_kind: "validation.not_allowed"; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "validation.not_allowed"; value: string } = std__validate__str_one_of_from(value, allowed, 0n, (BigInt([...allowed].length)));
  switch ($can_m1.$can_kind) {
  case "validation.not_allowed": {
    const e = $can_m1;
    return { $can_kind: "validation.not_allowed", value: e.value };
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
export function quota__request__validate(request: Quota__Request, schema: Quota__RequestSchema): { $can_kind: "ok"; value: Quota__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  if ((schema.label_minimum >= 0n)) {
    const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.invalid_length"; value: string; minimum: bigint; maximum: bigint } = std__validate__str_length(request.label, schema.label_minimum, schema.label_maximum);
    switch ($can_m1.$can_kind) {
    case "validation.invalid_bounds": {
      const _ = $can_m1;
      const $can_m2: { $can_kind: "ok"; value: string } = std__convert__int_to_str(schema.label_minimum);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const lo = $can_m2;
        const $can_m3: { $can_kind: "ok"; value: string } = std__convert__int_to_str(schema.label_maximum);
        switch ($can_m3.$can_kind) {
        case "ok": {
          const hi = $can_m3;
          return { $can_kind: "validation.schema_violation", path: "", rule: "schema.label_bounds", value: ((("minimum=" + lo.value) + ";maximum=") + hi.value) };
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
    case "validation.invalid_length": {
      const _ = $can_m1;
      return { $can_kind: "validation.schema_violation", path: "label", rule: "str.length_scalars", value: request.label };
    }
    case "ok": {
      const _ = $can_m1;
      const $can_m4: { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } = std__validate__int_range(request.amount, schema.amount_lower, schema.amount_upper);
      switch ($can_m4.$can_kind) {
      case "validation.invalid_bounds": {
        const _ = $can_m4;
        const $can_m5: { $can_kind: "ok"; value: string } = std__convert__int_to_str(schema.amount_lower);
        switch ($can_m5.$can_kind) {
        case "ok": {
          const lo = $can_m5;
          const $can_m6: { $can_kind: "ok"; value: string } = std__convert__int_to_str(schema.amount_upper);
          switch ($can_m6.$can_kind) {
          case "ok": {
            const hi = $can_m6;
            return { $can_kind: "validation.schema_violation", path: "", rule: "schema.amount_bounds", value: ((("lower=" + lo.value) + ";upper=") + hi.value) };
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
      case "validation.out_of_range": {
        const _ = $can_m4;
        const $can_m7: { $can_kind: "ok"; value: string } = std__convert__int_to_str(request.amount);
        switch ($can_m7.$can_kind) {
        case "ok": {
          const r = $can_m7;
          return { $can_kind: "validation.schema_violation", path: "amount", rule: "int.closed_range", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      case "ok": {
        const _ = $can_m4;
        const $can_m8: { $can_kind: "ok"; value: string } | { $can_kind: "validation.not_allowed"; value: string } = std__validate__str_one_of(request.mode, schema.allowed_modes);
        switch ($can_m8.$can_kind) {
        case "validation.not_allowed": {
          const _ = $can_m8;
          return { $can_kind: "validation.schema_violation", path: "mode", rule: "str.one_of", value: request.mode };
        }
        case "ok": {
          const _ = $can_m8;
          return { $can_kind: "ok", value: request };
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
    const $can_m9: { $can_kind: "ok"; value: string } = std__convert__int_to_str(schema.label_minimum);
    switch ($can_m9.$can_kind) {
    case "ok": {
      const r = $can_m9;
      return { $can_kind: "validation.schema_violation", path: "", rule: "schema.label_minimum", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function quota__request__admit(request: Quota__Request, schema: Quota__RequestSchema, quota: bigint): { $can_kind: "ok"; remaining: bigint; used: bigint } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } | { $can_kind: "validation.negative_value"; value: bigint } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } {
  const $can_m1: { $can_kind: "ok"; value: Quota__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = quota__request__validate(request, schema);
  switch ($can_m1.$can_kind) {
  case "validation.schema_violation": {
    const e = $can_m1;
    return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
  }
  case "ok": {
    const v = $can_m1;
    const $can_m2: { $can_kind: "ok"; remaining: bigint; used: bigint } | { $can_kind: "validation.negative_value"; value: bigint } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } = quota__consume(v.value.amount, quota);
    switch ($can_m2.$can_kind) {
    case "validation.negative_value": {
      const e = $can_m2;
      return { $can_kind: "validation.negative_value", value: e.value };
    }
    case "validation.invalid_bounds": {
      const e = $can_m2;
      return { $can_kind: "validation.invalid_bounds", lower: e.lower, upper: e.upper };
    }
    case "validation.out_of_range": {
      const e = $can_m2;
      return { $can_kind: "validation.out_of_range", value: e.value, lower: e.lower, upper: e.upper };
    }
    case "ok": {
      const u = $can_m2;
      return { $can_kind: "ok", used: u.used, remaining: u.remaining };
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
export function quota__envelope__validate(envelope: Quota__Envelope, schema: Quota__RequestSchema): { $can_kind: "ok"; value: Quota__Envelope } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  const $can_m1: { $can_kind: "ok"; value: Quota__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = quota__request__validate(envelope.request, schema);
  switch ($can_m1.$can_kind) {
  case "validation.schema_violation": {
    const e = $can_m1;
    if ((e.path === "")) {
      return { $can_kind: "validation.schema_violation", path: "request", rule: e.rule, value: e.value };
    }
    else {
      return { $can_kind: "validation.schema_violation", path: ("request." + e.path), rule: e.rule, value: e.value };
    }
  }
  case "ok": {
    const _ = $can_m1;
    return { $can_kind: "ok", value: envelope };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function quota__consume(amount: bigint, quota: bigint): { $can_kind: "ok"; remaining: bigint; used: bigint } | { $can_kind: "validation.negative_value"; value: bigint } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } {
  const $can_m1: { $can_kind: "ok", value: bigint } = { $can_kind: "ok", value: Quota__used };
  switch ($can_m1.$can_kind) {
  case "ok": {
    const s = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.negative_value"; value: bigint } = std__validate__int_nonnegative(amount);
    switch ($can_m2.$can_kind) {
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
    case "validation.negative_value": {
      return { $can_kind: "validation.negative_value", value: amount };
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
    case "ok": {
      const _ = $can_m2;
      return { $can_kind: "ok", used: s.value, remaining: (quota - s.value) };
    }
    case "validation.negative_value": {
      return { $can_kind: "validation.negative_value", value: quota };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  }
}
