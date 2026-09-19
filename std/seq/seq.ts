// GENERATED from seq.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type SeqResult = { $can_kind: "ok"; item: bigint } | { $can_kind: "ok"; item: string } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; values: bigint[] } | { $can_kind: "ok"; values: string[] } | { $can_kind: "sequence.index_out_of_range"; index: bigint } | { $can_kind: "sequence.invalid_slice"; start: bigint; end: bigint };
export type Int__Value = { value: bigint };
export type Seq__Values$T$int = { values: bigint[] };
export type Seq__Values$T$str = { values: string[] };
export type Seq__Item$T$int = { item: bigint };
export type Seq__Item$T$str = { item: string };
// Sequence indexing (a38 S3): bounds throw, matching Go.
function $canSeqAt<T>(a: T[], i: bigint): T {
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("seq index out of range");
  const k = Number(i);
  if (k >= a.length) throw new Error("seq index out of range");
  return a[k];
}
export function std__seq__empty$T$int(): { $can_kind: "ok"; values: bigint[] } {
  return { $can_kind: "ok", values: [] };
}
export function std__seq__empty$T$str(): { $can_kind: "ok"; values: string[] } {
  return { $can_kind: "ok", values: [] };
}
export function std__seq__singleton$T$int(value: bigint): { $can_kind: "ok"; values: bigint[] } {
  return { $can_kind: "ok", values: [value] };
}
export function std__seq__singleton$T$str(value: string): { $can_kind: "ok"; values: string[] } {
  return { $can_kind: "ok", values: [value] };
}
export function std__seq__length$T$int(values: bigint[]): { $can_kind: "ok"; value: bigint } {
  return { $can_kind: "ok", value: (BigInt([...values].length)) };
}
export function std__seq__length$T$str(values: string[]): { $can_kind: "ok"; value: bigint } {
  return { $can_kind: "ok", value: (BigInt([...values].length)) };
}
export function std__seq__get$T$int(values: bigint[], index: bigint): { $can_kind: "ok"; item: bigint } | { $can_kind: "sequence.index_out_of_range"; index: bigint } {
  const $can_m1 = (index < 0n);
  const $can_m2 = (index >= (BigInt([...values].length)));
  if ($can_m1) {
    return { $can_kind: "sequence.index_out_of_range", index: index };
  }
  else if ($can_m2) {
    return { $can_kind: "sequence.index_out_of_range", index: index };
  }
  else {
    return { $can_kind: "ok", item: $canSeqAt(values, index) };
  }
}
export function std__seq__get$T$str(values: string[], index: bigint): { $can_kind: "ok"; item: string } | { $can_kind: "sequence.index_out_of_range"; index: bigint } {
  const $can_m1 = (index < 0n);
  const $can_m2 = (index >= (BigInt([...values].length)));
  if ($can_m1) {
    return { $can_kind: "sequence.index_out_of_range", index: index };
  }
  else if ($can_m2) {
    return { $can_kind: "sequence.index_out_of_range", index: index };
  }
  else {
    return { $can_kind: "ok", item: $canSeqAt(values, index) };
  }
}
export function std__seq__append$T$int(values: bigint[], item: bigint): { $can_kind: "ok"; values: bigint[] } {
  return { $can_kind: "ok", values: [...values, item] };
}
export function std__seq__append$T$str(values: string[], item: string): { $can_kind: "ok"; values: string[] } {
  return { $can_kind: "ok", values: [...values, item] };
}
export function std__seq__concat_from$T$int(left: bigint[], right: bigint[], pos: bigint, n: bigint, acc: bigint[]): { $can_kind: "ok"; values: bigint[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", values: acc };
  }
  else {
    if ((pos < (BigInt([...right].length)))) {
      const $can_m1: { $can_kind: "ok"; values: bigint[] } = std__seq__concat_from$T$int(left, right, (pos + 1n), (n - 1n), [...acc, $canSeqAt(right, pos)]);
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        return { $can_kind: "ok", values: r.values };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $can_kind: "ok", values: acc };
    }
  }
}
export function std__seq__concat_from$T$str(left: string[], right: string[], pos: bigint, n: bigint, acc: string[]): { $can_kind: "ok"; values: string[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", values: acc };
  }
  else {
    if ((pos < (BigInt([...right].length)))) {
      const $can_m1: { $can_kind: "ok"; values: string[] } = std__seq__concat_from$T$str(left, right, (pos + 1n), (n - 1n), [...acc, $canSeqAt(right, pos)]);
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        return { $can_kind: "ok", values: r.values };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $can_kind: "ok", values: acc };
    }
  }
}
export function std__seq__concat$T$int(left: bigint[], right: bigint[]): { $can_kind: "ok"; values: bigint[] } {
  const $can_m1: { $can_kind: "ok"; values: bigint[] } = std__seq__concat_from$T$int(left, right, 0n, (BigInt([...right].length)), left);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", values: r.values };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__seq__concat$T$str(left: string[], right: string[]): { $can_kind: "ok"; values: string[] } {
  const $can_m1: { $can_kind: "ok"; values: string[] } = std__seq__concat_from$T$str(left, right, 0n, (BigInt([...right].length)), left);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", values: r.values };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__seq__slice_from$T$int(values: bigint[], end: bigint, pos: bigint, n: bigint, acc: bigint[]): { $can_kind: "ok"; values: bigint[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", values: acc };
  }
  else {
    const $can_m1 = (pos < end);
    const $can_m2 = (pos < (BigInt([...values].length)));
    if ($can_m1 && $can_m2) {
      const $can_m3: { $can_kind: "ok"; values: bigint[] } = std__seq__slice_from$T$int(values, end, (pos + 1n), (n - 1n), [...acc, $canSeqAt(values, pos)]);
      switch ($can_m3.$can_kind) {
      case "ok": {
        const r = $can_m3;
        return { $can_kind: "ok", values: r.values };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $can_kind: "ok", values: acc };
    }
  }
}
export function std__seq__slice_from$T$str(values: string[], end: bigint, pos: bigint, n: bigint, acc: string[]): { $can_kind: "ok"; values: string[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", values: acc };
  }
  else {
    const $can_m1 = (pos < end);
    const $can_m2 = (pos < (BigInt([...values].length)));
    if ($can_m1 && $can_m2) {
      const $can_m3: { $can_kind: "ok"; values: string[] } = std__seq__slice_from$T$str(values, end, (pos + 1n), (n - 1n), [...acc, $canSeqAt(values, pos)]);
      switch ($can_m3.$can_kind) {
      case "ok": {
        const r = $can_m3;
        return { $can_kind: "ok", values: r.values };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $can_kind: "ok", values: acc };
    }
  }
}
export function std__seq__slice$T$int(values: bigint[], start: bigint, end: bigint): { $can_kind: "ok"; values: bigint[] } | { $can_kind: "sequence.invalid_slice"; start: bigint; end: bigint } {
  const $can_m1 = (start < 0n);
  const $can_m2 = (start > end);
  const $can_m3 = (end > (BigInt([...values].length)));
  if ($can_m1) {
    return { $can_kind: "sequence.invalid_slice", start: start, end: end };
  }
  else if ($can_m2) {
    return { $can_kind: "sequence.invalid_slice", start: start, end: end };
  }
  else if ($can_m3) {
    return { $can_kind: "sequence.invalid_slice", start: start, end: end };
  }
  else {
    const $can_m4: { $can_kind: "ok"; values: bigint[] } = std__seq__slice_from$T$int(values, end, start, (end - start), []);
    switch ($can_m4.$can_kind) {
    case "ok": {
      const r = $can_m4;
      return { $can_kind: "ok", values: r.values };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__seq__slice$T$str(values: string[], start: bigint, end: bigint): { $can_kind: "ok"; values: string[] } | { $can_kind: "sequence.invalid_slice"; start: bigint; end: bigint } {
  const $can_m1 = (start < 0n);
  const $can_m2 = (start > end);
  const $can_m3 = (end > (BigInt([...values].length)));
  if ($can_m1) {
    return { $can_kind: "sequence.invalid_slice", start: start, end: end };
  }
  else if ($can_m2) {
    return { $can_kind: "sequence.invalid_slice", start: start, end: end };
  }
  else if ($can_m3) {
    return { $can_kind: "sequence.invalid_slice", start: start, end: end };
  }
  else {
    const $can_m4: { $can_kind: "ok"; values: string[] } = std__seq__slice_from$T$str(values, end, start, (end - start), []);
    switch ($can_m4.$can_kind) {
    case "ok": {
      const r = $can_m4;
      return { $can_kind: "ok", values: r.values };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
