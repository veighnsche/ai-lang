// GENERATED from quota.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { std__convert__int_to_str, type ScalarsResult } from "./scalars";
export type QuotaResult = { $can_kind: "ok" } | { $can_kind: "ok"; failures: Validate__Failure[]; passed: boolean } | { $can_kind: "ok"; remaining: bigint; used: bigint } | { $can_kind: "ok"; value: Quota__Envelope } | { $can_kind: "ok"; value: Quota__Request } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: string } | { $can_kind: "ok"; values: Validate__Failure[] } | { $can_kind: "validation.failed"; field: string; rule: string } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } | { $can_kind: "validation.not_positive"; value: bigint } | { $can_kind: "validation.negative_value"; value: bigint } | { $can_kind: "validation.dec_invalid_bounds"; lower: string; upper: string } | { $can_kind: "validation.dec_out_of_range"; value: string; lower: string; upper: string } | { $can_kind: "validation.dec_not_positive"; value: string } | { $can_kind: "validation.dec_negative_value"; value: string } | { $can_kind: "validation.invalid_length"; value: string; minimum: bigint; maximum: bigint } | { $can_kind: "validation.empty_value"; value: string } | { $can_kind: "validation.exclusive_choice" } | { $can_kind: "validation.not_allowed"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string };
export type Quota__Usage = { used: bigint; remaining: bigint };
export type Validate__Pass = {};
export type Int__Value = { value: bigint };
export type Dec__Value = { value: string };
export type Str__Value = { value: string };
export type Validate__Length = { path: string; tag: string; minimum: bigint; maximum: bigint };
export type Validate__Range = { path: string; tag: string; lower: bigint; upper: bigint };
export type Validate__Membership = { path: string; allowed: string[] };
export type Validate__Outcome = { passed: boolean; path: string; rule: string; detail: string };
export type Validate__Failure = { path: string; rule: string; detail: string };
export type Validate__Report = { passed: boolean; failures: Validate__Failure[] };
export type Validate__Failures = { values: Validate__Failure[] };
export type Quota__Request = { label: string; amount: bigint; mode: string };
export type Validate__StrSchema = { path: string; tag: string; minimum: bigint; maximum: bigint; allowed: string[]; restrict: boolean };
export type Validate__IntSchema = { path: string; tag: string; lower: bigint; upper: bigint };
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
export function std__validate__length_check(value: string, constraint: Validate__Length): { $can_kind: "ok"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  if ((constraint.minimum >= 0n)) {
    const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.invalid_length"; value: string; minimum: bigint; maximum: bigint } = std__validate__str_length(value, constraint.minimum, constraint.maximum);
    switch ($can_m1.$can_kind) {
    case "validation.invalid_bounds": {
      const _ = $can_m1;
      const $can_m2: { $can_kind: "ok"; value: string } = std__convert__int_to_str(constraint.minimum);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const lo = $can_m2;
        const $can_m3: { $can_kind: "ok"; value: string } = std__convert__int_to_str(constraint.maximum);
        switch ($can_m3.$can_kind) {
        case "ok": {
          const hi = $can_m3;
          return { $can_kind: "validation.schema_violation", path: "", rule: (("schema." + constraint.tag) + "_bounds"), value: ((("minimum=" + lo.value) + ";maximum=") + hi.value) };
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
      return { $can_kind: "validation.schema_violation", path: constraint.path, rule: "str.length_scalars", value: value };
    }
    case "ok": {
      const v = $can_m1;
      return { $can_kind: "ok", value: v.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  else {
    const $can_m4: { $can_kind: "ok"; value: string } = std__convert__int_to_str(constraint.minimum);
    switch ($can_m4.$can_kind) {
    case "ok": {
      const r = $can_m4;
      return { $can_kind: "validation.schema_violation", path: "", rule: (("schema." + constraint.tag) + "_minimum"), value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__validate__range_check(value: bigint, constraint: Validate__Range): { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } = std__validate__int_range(value, constraint.lower, constraint.upper);
  switch ($can_m1.$can_kind) {
  case "validation.invalid_bounds": {
    const _ = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: string } = std__convert__int_to_str(constraint.lower);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const lo = $can_m2;
      const $can_m3: { $can_kind: "ok"; value: string } = std__convert__int_to_str(constraint.upper);
      switch ($can_m3.$can_kind) {
      case "ok": {
        const hi = $can_m3;
        return { $can_kind: "validation.schema_violation", path: "", rule: (("schema." + constraint.tag) + "_bounds"), value: ((("lower=" + lo.value) + ";upper=") + hi.value) };
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
    const _ = $can_m1;
    const $can_m4: { $can_kind: "ok"; value: string } = std__convert__int_to_str(value);
    switch ($can_m4.$can_kind) {
    case "ok": {
      const r = $can_m4;
      return { $can_kind: "validation.schema_violation", path: constraint.path, rule: "int.closed_range", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "ok": {
    const v = $can_m1;
    return { $can_kind: "ok", value: v.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__validate__membership_check(value: string, constraint: Validate__Membership): { $can_kind: "ok"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "validation.not_allowed"; value: string } = std__validate__str_one_of(value, constraint.allowed);
  switch ($can_m1.$can_kind) {
  case "validation.not_allowed": {
    const _ = $can_m1;
    return { $can_kind: "validation.schema_violation", path: constraint.path, rule: "str.one_of", value: value };
  }
  case "ok": {
    const v = $can_m1;
    return { $can_kind: "ok", value: v.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__validate__all_from(checks: Validate__Outcome[], pos: bigint, n: bigint, acc: Validate__Failure[]): { $can_kind: "ok"; failures: Validate__Failure[]; passed: boolean } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", passed: ((BigInt([...acc].length)) === 0n), failures: acc };
  }
  else {
    if ((pos < (BigInt([...checks].length)))) {
      const $can_m1: { $can_kind: "ok"; values: Validate__Failure[] } = std__validate__all_push($canSeqAt(checks, pos), acc);
      switch ($can_m1.$can_kind) {
      case "ok": {
        const p = $can_m1;
        const $can_m2: { $can_kind: "ok"; failures: Validate__Failure[]; passed: boolean } = std__validate__all_from(checks, (pos + 1n), (n - 1n), p.values);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", passed: r.passed, failures: r.failures };
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
      return { $can_kind: "ok", passed: ((BigInt([...acc].length)) === 0n), failures: acc };
    }
  }
}
export function std__validate__all_push(current: Validate__Outcome, acc: Validate__Failure[]): { $can_kind: "ok"; values: Validate__Failure[] } {
  if (current.passed) {
    return { $can_kind: "ok", values: acc };
  }
  else {
    return { $can_kind: "ok", values: [...acc, { path: current.path, rule: current.rule, detail: current.detail }] };
  }
}
export function std__validate__all(checks: Validate__Outcome[]): { $can_kind: "ok"; failures: Validate__Failure[]; passed: boolean } {
  const $can_m1: { $can_kind: "ok"; failures: Validate__Failure[]; passed: boolean } = std__validate__all_from(checks, 0n, (BigInt([...checks].length)), []);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", passed: r.passed, failures: r.failures };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__validate__schema__str(value: string, schema: Validate__StrSchema): { $can_kind: "ok"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__length_check(value, { path: schema.path, tag: schema.tag, minimum: schema.minimum, maximum: schema.maximum });
  switch ($can_m1.$can_kind) {
  case "validation.schema_violation": {
    const e = $can_m1;
    return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
  }
  case "ok": {
    const v = $can_m1;
    if (schema.restrict) {
      const $can_m2: { $can_kind: "ok"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__membership_check(v.value, { path: schema.path, allowed: schema.allowed });
      switch ($can_m2.$can_kind) {
      case "validation.schema_violation": {
        const e = $can_m2;
        return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
      }
      case "ok": {
        const m = $can_m2;
        return { $can_kind: "ok", value: m.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $can_kind: "ok", value: v.value };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__validate__schema__int(value: bigint, schema: Validate__IntSchema): { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__range_check(value, { path: schema.path, tag: schema.tag, lower: schema.lower, upper: schema.upper });
  switch ($can_m1.$can_kind) {
  case "validation.schema_violation": {
    const e = $can_m1;
    return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
  }
  case "ok": {
    const v = $can_m1;
    return { $can_kind: "ok", value: v.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function quota__request__validate(request: Quota__Request, schema: Quota__RequestSchema): { $can_kind: "ok"; value: Quota__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__length_check(request.label, { path: "label", tag: "label", minimum: schema.label_minimum, maximum: schema.label_maximum });
  switch ($can_m1.$can_kind) {
  case "validation.schema_violation": {
    const e = $can_m1;
    return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
  }
  case "ok": {
    const _ = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__range_check(request.amount, { path: "amount", tag: "amount", lower: schema.amount_lower, upper: schema.amount_upper });
    switch ($can_m2.$can_kind) {
    case "validation.schema_violation": {
      const e = $can_m2;
      return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
    }
    case "ok": {
      const _ = $can_m2;
      const $can_m3: { $can_kind: "ok"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__membership_check(request.mode, { path: "mode", allowed: schema.allowed_modes });
      switch ($can_m3.$can_kind) {
      case "validation.schema_violation": {
        const e = $can_m3;
        return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
      }
      case "ok": {
        const _ = $can_m3;
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
