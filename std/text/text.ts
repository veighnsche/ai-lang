// GENERATED from text.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type TextResult = { $ail_kind: "ok"; value: Uint8Array } | { $ail_kind: "ok"; value: bigint } | { $ail_kind: "ok"; value: boolean } | { $ail_kind: "ok"; value: string } | { $ail_kind: "ok"; values: string[] } | { $ail_kind: "text.index_out_of_range"; value: string; index: bigint } | { $ail_kind: "text.invalid_slice"; value: string; start: bigint; end: bigint } | { $ail_kind: "text.not_found"; value: string; pattern: string } | { $ail_kind: "text.empty_pattern" } | { $ail_kind: "text.empty_separator" } | { $ail_kind: "encoding.invalid_utf8"; value: Uint8Array };
export type Str__Value = { value: string };
export type Bool__Value = { value: boolean };
export type Int__Value = { value: bigint };
export type Split__Result = { values: string[] };
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
// Sequence indexing (v38 S3): bounds throw, matching Go.
function $ailSeqAt<T>(a: T[], i: bigint): T {
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("seq index out of range");
  const k = Number(i);
  if (k >= a.length) throw new Error("seq index out of range");
  return a[k];
}
function $ailUtf8Decode(value: Uint8Array): { $ail_kind: "ok"; value: string } | { $ail_kind: "encoding.invalid_utf8"; value: Uint8Array } {
  let i = 0;
  const n = value.length;
  let valid = true;
  while (i < n && valid) {
    const b0 = value[i];
    if (b0 < 0x80) { i += 1; continue; }
    let need = 0; let lo = 0x80; let hi = 0xBF;
    if (b0 >= 0xC2 && b0 <= 0xDF) { need = 1; }
    else if (b0 === 0xE0) { need = 2; lo = 0xA0; }
    else if (b0 >= 0xE1 && b0 <= 0xEC) { need = 2; }
    else if (b0 === 0xED) { need = 2; hi = 0x9F; }
    else if (b0 >= 0xEE && b0 <= 0xEF) { need = 2; }
    else if (b0 === 0xF0) { need = 3; lo = 0x90; }
    else if (b0 >= 0xF1 && b0 <= 0xF3) { need = 3; }
    else if (b0 === 0xF4) { need = 3; hi = 0x8F; }
    else { valid = false; break; }
    if (i + need >= n) { valid = false; break; }
    const b1 = value[i + 1];
    if (b1 < lo || b1 > hi) { valid = false; break; }
    for (let k = 2; k <= need; k++) {
      const b = value[i + k];
      if (b < 0x80 || b > 0xBF) { valid = false; break; }
    }
    i += 1 + need;
  }
  if (!valid) return { $ail_kind: "encoding.invalid_utf8", value: value };
  return { $ail_kind: "ok", value: new TextDecoder("utf-8", { fatal: true, ignoreBOM: true }).decode(value) };
}
function $ailHexEncode(value: Uint8Array): string {
  const digits = "0123456789abcdef";
  let out = "";
  for (let i = 0; i < value.length; i++) {
    out += digits[value[i] >> 4] + digits[value[i] & 15];
  }
  return out;
}
export function std__str__concat(left: string, right: string): { $ail_kind: "ok"; value: string } {
  return { $ail_kind: "ok", value: (left + right) };
}
export function std__str__is_empty(value: string): { $ail_kind: "ok"; value: boolean } {
  if ((value === "")) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__str__is_whitespace(scalar: bigint): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1 = (scalar === 32n);
  const $ail_m2 = (scalar === 9n);
  const $ail_m3 = (scalar === 10n);
  const $ail_m4 = (scalar === 13n);
  if ($ail_m1) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m2) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m3) {
    return { $ail_kind: "ok", value: true };
  }
  else if ($ail_m4) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__str__length_scalars(value: string): { $ail_kind: "ok"; value: bigint } {
  return { $ail_kind: "ok", value: (BigInt([...value].length)) };
}
export function std__str__scalar_at(value: string, index: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "text.index_out_of_range"; value: string; index: bigint } {
  if ((index >= 0n)) {
    if ((index < (BigInt([...value].length)))) {
      return { $ail_kind: "ok", value: $ailStrAt(value, index) };
    }
    else {
      return { $ail_kind: "text.index_out_of_range", value: value, index: index };
    }
  }
  else {
    return { $ail_kind: "text.index_out_of_range", value: value, index: index };
  }
}
export function std__str__slice_scalars(value: string, start: bigint, end: bigint): { $ail_kind: "ok"; value: string } | { $ail_kind: "text.invalid_slice"; value: string; start: bigint; end: bigint } {
  const $ail_m1 = (start >= 0n);
  const $ail_m2 = (end >= start);
  const $ail_m3 = (end <= (BigInt([...value].length)));
  if ($ail_m1 && $ail_m2 && $ail_m3) {
    return { $ail_kind: "ok", value: $ailStrSlice(value, start, end) };
  }
  else if ($ail_m1 && $ail_m2 && !($ail_m3)) {
    return { $ail_kind: "text.invalid_slice", value: value, start: start, end: end };
  }
  else if ($ail_m1 && !($ail_m2)) {
    return { $ail_kind: "text.invalid_slice", value: value, start: start, end: end };
  }
  else {
    return { $ail_kind: "text.invalid_slice", value: value, start: start, end: end };
  }
}
export function std__str__find_from(orig: string, value: string, pattern: string, pos: bigint, n: bigint): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "text.not_found"; value: string; pattern: string } {
  if ((n <= 0n)) {
    return { $ail_kind: "text.not_found", value: orig, pattern: pattern };
  }
  else {
    if (((BigInt([...value].length)) >= (BigInt([...pattern].length)))) {
      if (($ailStrSlice(value, 0n, (BigInt([...pattern].length))) === pattern)) {
        return { $ail_kind: "ok", value: pos };
      }
      else {
        const $ail_m1: { $ail_kind: "ok"; value: bigint } | { $ail_kind: "text.not_found"; value: string; pattern: string } = std__str__find_from(orig, $ailStrSlice(value, 1n, (BigInt([...value].length))), pattern, (pos + 1n), (n - 1n));
        switch ($ail_m1.$ail_kind) {
        case "text.not_found": {
          const e = $ail_m1;
          return { $ail_kind: "text.not_found", value: e.value, pattern: e.pattern };
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
    }
    else {
      return { $ail_kind: "text.not_found", value: orig, pattern: pattern };
    }
  }
}
export function std__str__find(value: string, pattern: string): { $ail_kind: "ok"; value: bigint } | { $ail_kind: "text.not_found"; value: string; pattern: string } {
  if (((BigInt([...pattern].length)) === 0n)) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } | { $ail_kind: "text.not_found"; value: string; pattern: string } = std__str__find_from(value, value, pattern, 0n, (BigInt([...value].length)));
    switch ($ail_m1.$ail_kind) {
    case "text.not_found": {
      const e = $ail_m1;
      return { $ail_kind: "text.not_found", value: e.value, pattern: e.pattern };
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
}
export function std__str__contains(value: string, pattern: string): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1: { $ail_kind: "ok"; value: bigint } | { $ail_kind: "text.not_found"; value: string; pattern: string } = std__str__find(value, pattern);
  switch ($ail_m1.$ail_kind) {
  case "text.not_found": {
    const _ = $ail_m1;
    return { $ail_kind: "ok", value: false };
  }
  case "ok": {
    const _ = $ail_m1;
    return { $ail_kind: "ok", value: true };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__str__starts_with(value: string, pattern: string): { $ail_kind: "ok"; value: boolean } {
  if (((BigInt([...pattern].length)) <= (BigInt([...value].length)))) {
    if (($ailStrSlice(value, 0n, (BigInt([...pattern].length))) === pattern)) {
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
export function std__str__ends_with(value: string, pattern: string): { $ail_kind: "ok"; value: boolean } {
  if (((BigInt([...pattern].length)) <= (BigInt([...value].length)))) {
    if (($ailStrSlice(value, ((BigInt([...value].length)) - (BigInt([...pattern].length))), (BigInt([...value].length))) === pattern)) {
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
export function std__str__replace_all_from(value: string, old: string, replacement: string, acc: string, n: bigint): { $ail_kind: "ok"; value: string } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: (acc + value) };
  }
  else {
    if (((BigInt([...value].length)) >= (BigInt([...old].length)))) {
      if (($ailStrSlice(value, 0n, (BigInt([...old].length))) === old)) {
        const $ail_m1: { $ail_kind: "ok"; value: string } = std__str__replace_all_from($ailStrSlice(value, (BigInt([...old].length)), (BigInt([...value].length))), old, replacement, (acc + replacement), (n - 1n));
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
        const $ail_m2: { $ail_kind: "ok"; value: string } = std__str__replace_all_from($ailStrSlice(value, 1n, (BigInt([...value].length))), old, replacement, (acc + $ailStrSlice(value, 0n, 1n)), (n - 1n));
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
    else {
      return { $ail_kind: "ok", value: (acc + value) };
    }
  }
}
export function std__str__replace_all(value: string, old: string, replacement: string): { $ail_kind: "ok"; value: string } | { $ail_kind: "text.empty_pattern" } {
  if (((BigInt([...old].length)) === 0n)) {
    return { $ail_kind: "text.empty_pattern" };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: string } = std__str__replace_all_from(value, old, replacement, "", (BigInt([...value].length)));
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
export function std__str__trim_left(value: string, n: bigint): { $ail_kind: "ok"; value: string } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: boolean } = std__str__is_whitespace($ailStrAt(value, 0n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const w = $ail_m1;
      if (w.value) {
        const $ail_m2: { $ail_kind: "ok"; value: string } = std__str__trim_left($ailStrSlice(value, 1n, (BigInt([...value].length))), (n - 1n));
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
        const $ail_m3: { $ail_kind: "ok"; value: string } = std__str__trim_right(value, (BigInt([...value].length)));
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
export function std__str__trim_right(value: string, n: bigint): { $ail_kind: "ok"; value: string } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: boolean } = std__str__is_whitespace($ailStrAt(value, (n - 1n)));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const w = $ail_m1;
      if (w.value) {
        const $ail_m2: { $ail_kind: "ok"; value: string } = std__str__trim_right($ailStrSlice(value, 0n, (n - 1n)), (n - 1n));
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
        return { $ail_kind: "ok", value: value };
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__str__trim_ascii(value: string): { $ail_kind: "ok"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; value: string } = std__str__trim_left(value, (BigInt([...value].length)));
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
export function std__str__upper_ascii_from(value: string, acc: string, n: bigint): { $ail_kind: "ok"; value: string } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } | { $ail_kind: "text.not_found"; value: string; pattern: string } = std__str__find("abcdefghijklmnopqrstuvwxyz", $ailStrSlice(value, 0n, 1n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const f = $ail_m1;
      const $ail_m2: { $ail_kind: "ok"; value: string } = std__str__upper_ascii_from($ailStrSlice(value, 1n, (BigInt([...value].length))), (acc + $ailStrSlice("ABCDEFGHIJKLMNOPQRSTUVWXYZ", f.value, (f.value + 1n))), (n - 1n));
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
    case "text.not_found": {
      const _ = $ail_m1;
      const $ail_m3: { $ail_kind: "ok"; value: string } = std__str__upper_ascii_from($ailStrSlice(value, 1n, (BigInt([...value].length))), (acc + $ailStrSlice(value, 0n, 1n)), (n - 1n));
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
}
export function std__str__upper_ascii(value: string): { $ail_kind: "ok"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; value: string } = std__str__upper_ascii_from(value, "", (BigInt([...value].length)));
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
export function std__str__lower_ascii_from(value: string, acc: string, n: bigint): { $ail_kind: "ok"; value: string } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: bigint } | { $ail_kind: "text.not_found"; value: string; pattern: string } = std__str__find("ABCDEFGHIJKLMNOPQRSTUVWXYZ", $ailStrSlice(value, 0n, 1n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const f = $ail_m1;
      const $ail_m2: { $ail_kind: "ok"; value: string } = std__str__lower_ascii_from($ailStrSlice(value, 1n, (BigInt([...value].length))), (acc + $ailStrSlice("abcdefghijklmnopqrstuvwxyz", f.value, (f.value + 1n))), (n - 1n));
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
    case "text.not_found": {
      const _ = $ail_m1;
      const $ail_m3: { $ail_kind: "ok"; value: string } = std__str__lower_ascii_from($ailStrSlice(value, 1n, (BigInt([...value].length))), (acc + $ailStrSlice(value, 0n, 1n)), (n - 1n));
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
}
export function std__str__lower_ascii(value: string): { $ail_kind: "ok"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; value: string } = std__str__lower_ascii_from(value, "", (BigInt([...value].length)));
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
export function std__str__join_from(values: string[], separator: string, position: bigint, fuel: bigint, acc: string): { $ail_kind: "ok"; value: string } {
  if ((fuel <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    if ((position < (BigInt([...values].length)))) {
      if ((position === 0n)) {
        const $ail_m1: { $ail_kind: "ok"; value: string } = std__str__join_from(values, separator, (position + 1n), (fuel - 1n), (acc + $ailSeqAt(values, position)));
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
        const $ail_m2: { $ail_kind: "ok"; value: string } = std__str__join_from(values, separator, (position + 1n), (fuel - 1n), ((acc + separator) + $ailSeqAt(values, position)));
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
    else {
      return { $ail_kind: "ok", value: acc };
    }
  }
}
export function std__str__join(values: string[], separator: string): { $ail_kind: "ok"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; value: string } = std__str__join_from(values, separator, 0n, ((BigInt([...values].length)) + 1n), "");
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
export function std__str__split_from(value: string, separator: string, current: string, acc: string[], fuel: bigint): { $ail_kind: "ok"; values: string[] } {
  if ((fuel <= 0n)) {
    return { $ail_kind: "ok", values: [...acc, current] };
  }
  else {
    if ((value === "")) {
      return { $ail_kind: "ok", values: [...acc, current] };
    }
    else {
      if (((BigInt([...value].length)) >= (BigInt([...separator].length)))) {
        if (($ailStrSlice(value, 0n, (BigInt([...separator].length))) === separator)) {
          const $ail_m1: { $ail_kind: "ok"; values: string[] } = std__str__split_from($ailStrSlice(value, (BigInt([...separator].length)), (BigInt([...value].length))), separator, "", [...acc, current], (fuel - 1n));
          switch ($ail_m1.$ail_kind) {
          case "ok": {
            const r = $ail_m1;
            return { $ail_kind: "ok", values: r.values };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          const $ail_m2: { $ail_kind: "ok"; values: string[] } = std__str__split_from($ailStrSlice(value, 1n, (BigInt([...value].length))), separator, (current + $ailStrSlice(value, 0n, 1n)), acc, (fuel - 1n));
          switch ($ail_m2.$ail_kind) {
          case "ok": {
            const r = $ail_m2;
            return { $ail_kind: "ok", values: r.values };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
      }
      else {
        const $ail_m3: { $ail_kind: "ok"; values: string[] } = std__str__split_from($ailStrSlice(value, 1n, (BigInt([...value].length))), separator, (current + $ailStrSlice(value, 0n, 1n)), acc, (fuel - 1n));
        switch ($ail_m3.$ail_kind) {
        case "ok": {
          const r = $ail_m3;
          return { $ail_kind: "ok", values: r.values };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
  }
}
export function std__str__split(value: string, separator: string): { $ail_kind: "ok"; values: string[] } | { $ail_kind: "text.empty_separator" } {
  if ((separator === "")) {
    return { $ail_kind: "text.empty_separator" };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; values: string[] } = std__str__split_from(value, separator, "", [], ((BigInt([...value].length)) + 1n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const r = $ail_m1;
      return { $ail_kind: "ok", values: r.values };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__utf8__encode(value: string): { $ail_kind: "ok"; value: Uint8Array } {
  const $ail_m1: { $ail_kind: "ok", value: Uint8Array } = { $ail_kind: "ok", value: new TextEncoder().encode(value) };
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const r = $ail_m1;
    return { $ail_kind: "ok", value: r.value };
  }
  }
}
export function std__utf8__decode(value: Uint8Array): { $ail_kind: "ok"; value: string } | { $ail_kind: "encoding.invalid_utf8"; value: Uint8Array } {
  const $ail_m1: { $ail_kind: "ok"; value: string } | { $ail_kind: "encoding.invalid_utf8"; value: Uint8Array } = $ailUtf8Decode(value);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const r = $ail_m1;
    return { $ail_kind: "ok", value: r.value };
  }
  case "encoding.invalid_utf8": {
    const e = $ail_m1;
    return { $ail_kind: "encoding.invalid_utf8", value: e.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__hex__encode(value: Uint8Array): { $ail_kind: "ok"; value: string } {
  const $ail_m1: { $ail_kind: "ok", value: string } = { $ail_kind: "ok", value: $ailHexEncode(value) };
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const r = $ail_m1;
    return { $ail_kind: "ok", value: r.value };
  }
  }
}
