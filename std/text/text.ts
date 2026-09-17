// GENERATED from text.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type TextResult = { $ail_kind: "ok"; value: bigint } | { $ail_kind: "ok"; value: boolean } | { $ail_kind: "ok"; value: string } | { $ail_kind: "text.index_out_of_range"; value: string; index: bigint } | { $ail_kind: "text.invalid_slice"; value: string; start: bigint; end: bigint } | { $ail_kind: "text.not_found"; value: string; pattern: string } | { $ail_kind: "text.empty_pattern" };
export type Str__Value = { value: string };
export type Bool__Value = { value: boolean };
export type Int__Value = { value: bigint };
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
  return BigInt(cps[k].codePointAt(0));
}
function $ailStrSlice(s: string, a: bigint, b: bigint): string {
  const cps = [...s];
  const toIdx = (x: bigint): number => { if (x < 0n || x > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("str slice out of range"); return Number(x); };
  const lo = toIdx(a), hi = toIdx(b);
  if (lo > hi || hi > cps.length) throw new Error("str slice out of range");
  return cps.slice(lo, hi).join("");
}
export function std__str__concat(left: string, right: string): TextResult {
  return { $ail_kind: "ok", value: (left + right) };
}
export function std__str__is_empty(value: string): TextResult {
  if ((value === "")) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function std__str__is_whitespace(scalar: bigint): TextResult {
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
export function std__str__length_scalars(value: string): TextResult {
  return { $ail_kind: "ok", value: (BigInt([...value].length)) };
}
export function std__str__scalar_at(value: string, index: bigint): TextResult {
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
export function std__str__slice_scalars(value: string, start: bigint, end: bigint): TextResult {
  const $ail_m5 = (start >= 0n);
  const $ail_m6 = (end >= start);
  const $ail_m7 = (end <= (BigInt([...value].length)));
  if ($ail_m5 && $ail_m6 && $ail_m7) {
    return { $ail_kind: "ok", value: $ailStrSlice(value, start, end) };
  }
  else if ($ail_m5 && $ail_m6 && !($ail_m7)) {
    return { $ail_kind: "text.invalid_slice", value: value, start: start, end: end };
  }
  else if ($ail_m5 && !($ail_m6)) {
    return { $ail_kind: "text.invalid_slice", value: value, start: start, end: end };
  }
  else {
    return { $ail_kind: "text.invalid_slice", value: value, start: start, end: end };
  }
}
export function std__str__find_from(orig: string, value: string, pattern: string, pos: bigint, n: bigint): TextResult {
  if ((n <= 0n)) {
    return { $ail_kind: "text.not_found", value: orig, pattern: pattern };
  }
  else {
    if (((BigInt([...value].length)) >= (BigInt([...pattern].length)))) {
      if (($ailStrSlice(value, 0n, (BigInt([...pattern].length))) === pattern)) {
        return { $ail_kind: "ok", value: pos };
      }
      else {
        const $ail_m8: TextResult = std__str__find_from(orig, $ailStrSlice(value, 1n, (BigInt([...value].length))), pattern, (pos + 1n), (n - 1n));
        switch ($ail_m8.$ail_kind) {
        case "text.not_found": {
          const e = $ail_m8;
          return { $ail_kind: "text.not_found", value: e.value, pattern: e.pattern };
        }
        case "ok": {
          const r = $ail_m8;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
    }
    else {
      return { $ail_kind: "text.not_found", value: orig, pattern: pattern };
    }
  }
}
export function std__str__find(value: string, pattern: string): TextResult {
  if (((BigInt([...pattern].length)) === 0n)) {
    return { $ail_kind: "ok", value: 0n };
  }
  else {
    const $ail_m9: TextResult = std__str__find_from(value, value, pattern, 0n, (BigInt([...value].length)));
    switch ($ail_m9.$ail_kind) {
    case "text.not_found": {
      const e = $ail_m9;
      return { $ail_kind: "text.not_found", value: e.value, pattern: e.pattern };
    }
    case "ok": {
      const r = $ail_m9;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__str__contains(value: string, pattern: string): TextResult {
  const $ail_m10: TextResult = std__str__find(value, pattern);
  switch ($ail_m10.$ail_kind) {
  case "text.not_found": {
    const _ = $ail_m10;
    return { $ail_kind: "ok", value: false };
  }
  case "ok": {
    const _ = $ail_m10;
    return { $ail_kind: "ok", value: true };
  }
  }
}
export function std__str__starts_with(value: string, pattern: string): TextResult {
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
export function std__str__ends_with(value: string, pattern: string): TextResult {
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
export function std__str__replace_all_from(value: string, old: string, replacement: string, acc: string, n: bigint): TextResult {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: (acc + value) };
  }
  else {
    if (((BigInt([...value].length)) >= (BigInt([...old].length)))) {
      if (($ailStrSlice(value, 0n, (BigInt([...old].length))) === old)) {
        const $ail_m11: TextResult = std__str__replace_all_from($ailStrSlice(value, (BigInt([...old].length)), (BigInt([...value].length))), old, replacement, (acc + replacement), (n - 1n));
        switch ($ail_m11.$ail_kind) {
        case "ok": {
          const r = $ail_m11;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
      else {
        const $ail_m12: TextResult = std__str__replace_all_from($ailStrSlice(value, 1n, (BigInt([...value].length))), old, replacement, (acc + $ailStrSlice(value, 0n, 1n)), (n - 1n));
        switch ($ail_m12.$ail_kind) {
        case "ok": {
          const r = $ail_m12;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
    }
    else {
      return { $ail_kind: "ok", value: (acc + value) };
    }
  }
}
export function std__str__replace_all(value: string, old: string, replacement: string): TextResult {
  if (((BigInt([...old].length)) === 0n)) {
    return { $ail_kind: "text.empty_pattern" };
  }
  else {
    const $ail_m13: TextResult = std__str__replace_all_from(value, old, replacement, "", (BigInt([...value].length)));
    switch ($ail_m13.$ail_kind) {
    case "ok": {
      const r = $ail_m13;
      return { $ail_kind: "ok", value: r.value };
    }
    }
  }
}
export function std__str__trim_left(value: string, n: bigint): TextResult {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    const $ail_m14: TextResult = std__str__is_whitespace($ailStrAt(value, 0n));
    switch ($ail_m14.$ail_kind) {
    case "ok": {
      const w = $ail_m14;
      if (w.value) {
        const $ail_m15: TextResult = std__str__trim_left($ailStrSlice(value, 1n, (BigInt([...value].length))), (n - 1n));
        switch ($ail_m15.$ail_kind) {
        case "ok": {
          const r = $ail_m15;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
      else {
        const $ail_m16: TextResult = std__str__trim_right(value, (BigInt([...value].length)));
        switch ($ail_m16.$ail_kind) {
        case "ok": {
          const r = $ail_m16;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
    }
    }
  }
}
export function std__str__trim_right(value: string, n: bigint): TextResult {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: value };
  }
  else {
    const $ail_m17: TextResult = std__str__is_whitespace($ailStrAt(value, (n - 1n)));
    switch ($ail_m17.$ail_kind) {
    case "ok": {
      const w = $ail_m17;
      if (w.value) {
        const $ail_m18: TextResult = std__str__trim_right($ailStrSlice(value, 0n, (n - 1n)), (n - 1n));
        switch ($ail_m18.$ail_kind) {
        case "ok": {
          const r = $ail_m18;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
      else {
        return { $ail_kind: "ok", value: value };
      }
    }
    }
  }
}
export function std__str__trim_ascii(value: string): TextResult {
  const $ail_m19: TextResult = std__str__trim_left(value, (BigInt([...value].length)));
  switch ($ail_m19.$ail_kind) {
  case "ok": {
    const r = $ail_m19;
    return { $ail_kind: "ok", value: r.value };
  }
  }
}
export function std__str__upper_ascii_from(value: string, acc: string, n: bigint): TextResult {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    const $ail_m20: TextResult = std__str__find("abcdefghijklmnopqrstuvwxyz", $ailStrSlice(value, 0n, 1n));
    switch ($ail_m20.$ail_kind) {
    case "ok": {
      const f = $ail_m20;
      const $ail_m21: TextResult = std__str__upper_ascii_from($ailStrSlice(value, 1n, (BigInt([...value].length))), (acc + $ailStrSlice("ABCDEFGHIJKLMNOPQRSTUVWXYZ", f.value, (f.value + 1n))), (n - 1n));
      switch ($ail_m21.$ail_kind) {
      case "ok": {
        const r = $ail_m21;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
    case "text.not_found": {
      const _ = $ail_m20;
      const $ail_m22: TextResult = std__str__upper_ascii_from($ailStrSlice(value, 1n, (BigInt([...value].length))), (acc + $ailStrSlice(value, 0n, 1n)), (n - 1n));
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
export function std__str__upper_ascii(value: string): TextResult {
  const $ail_m23: TextResult = std__str__upper_ascii_from(value, "", (BigInt([...value].length)));
  switch ($ail_m23.$ail_kind) {
  case "ok": {
    const r = $ail_m23;
    return { $ail_kind: "ok", value: r.value };
  }
  }
}
export function std__str__lower_ascii_from(value: string, acc: string, n: bigint): TextResult {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    const $ail_m24: TextResult = std__str__find("ABCDEFGHIJKLMNOPQRSTUVWXYZ", $ailStrSlice(value, 0n, 1n));
    switch ($ail_m24.$ail_kind) {
    case "ok": {
      const f = $ail_m24;
      const $ail_m25: TextResult = std__str__lower_ascii_from($ailStrSlice(value, 1n, (BigInt([...value].length))), (acc + $ailStrSlice("abcdefghijklmnopqrstuvwxyz", f.value, (f.value + 1n))), (n - 1n));
      switch ($ail_m25.$ail_kind) {
      case "ok": {
        const r = $ail_m25;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
    case "text.not_found": {
      const _ = $ail_m24;
      const $ail_m26: TextResult = std__str__lower_ascii_from($ailStrSlice(value, 1n, (BigInt([...value].length))), (acc + $ailStrSlice(value, 0n, 1n)), (n - 1n));
      switch ($ail_m26.$ail_kind) {
      case "ok": {
        const r = $ail_m26;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
    }
  }
}
export function std__str__lower_ascii(value: string): TextResult {
  const $ail_m27: TextResult = std__str__lower_ascii_from(value, "", (BigInt([...value].length)));
  switch ($ail_m27.$ail_kind) {
  case "ok": {
    const r = $ail_m27;
    return { $ail_kind: "ok", value: r.value };
  }
  }
}
