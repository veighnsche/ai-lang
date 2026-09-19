// GENERATED from seq.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type SeqResult = { $can_kind: "ok"; item: bigint } | { $can_kind: "ok"; item: string } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: boolean } | { $can_kind: "ok"; values: bigint[] } | { $can_kind: "ok"; values: string[] } | { $can_kind: "sequence.index_out_of_range"; index: bigint } | { $can_kind: "sequence.invalid_slice"; start: bigint; end: bigint } | { $can_kind: "sequence.not_found" };
export type Int__Value = { value: bigint };
export type Seq__Values$T$int = { values: bigint[] };
export type Seq__Values$T$str = { values: string[] };
export type Seq__Item$T$int = { item: bigint };
export type Seq__Item$T$str = { item: string };
export type Bool__Value = { value: boolean };
export type Seq__FoldStep$T$int$T$int = { acc: bigint; elem: bigint };
export type Seq__FoldStep$T$str$T$str = { acc: string; elem: string };
// Sequence indexing (a38 S3): bounds throw, matching Go.
function $canSeqAt<T>(a: T[], i: bigint): T {
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("seq index out of range");
  const k = Number(i);
  if (k >= a.length) throw new Error("seq index out of range");
  return a[k];
}
// Euclidean integer division (a17): quotient and remainder with
// 0 <= r < |b| on every sign combination.
function $canDivMod(a: bigint, b: bigint): [bigint, bigint] {
  let q: bigint = a / b;
  let r: bigint = a - b * q;
  if (r < 0n) {
    q += b > 0n ? -1n : 1n;
    r = a - b * q;
  }
  return [q, r];
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
export function std__seq__test_double(x: bigint): { $can_kind: "ok"; item: bigint } {
  return { $can_kind: "ok", item: (x + x) };
}
export function std__seq__test_length(s: string): { $can_kind: "ok"; item: bigint } {
  return { $can_kind: "ok", item: (BigInt([...s].length)) };
}
export function std__seq__test_is_even(x: bigint): { $can_kind: "ok"; value: boolean } {
  return { $can_kind: "ok", value: ($canDivMod(x, 2n)[1] === 0n) };
}
export function std__seq__test_nonempty(s: string): { $can_kind: "ok"; value: boolean } {
  return { $can_kind: "ok", value: ((BigInt([...s].length)) > 0n) };
}
export function std__seq__test_add(step: Seq__FoldStep$T$int$T$int): { $can_kind: "ok"; item: bigint } {
  return { $can_kind: "ok", item: (step.acc + step.elem) };
}
export function std__seq__test_concat(step: Seq__FoldStep$T$str$T$str): { $can_kind: "ok"; item: string } {
  return { $can_kind: "ok", item: (step.acc + step.elem) };
}
export function std__seq__map_from$T$int$T$int(values: bigint[], f: (input: bigint) => { $can_kind: "ok"; item: bigint }, pos: bigint, n: bigint, acc: bigint[]): { $can_kind: "ok"; values: bigint[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", values: acc };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; item: bigint } = f($canSeqAt(values, pos));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        const $can_m2: { $can_kind: "ok"; values: bigint[] } = std__seq__map_from$T$int$T$int(values, f, (pos + 1n), (n - 1n), [...acc, r.item]);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const s = $can_m2;
          return { $can_kind: "ok", values: s.values };
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
      return { $can_kind: "ok", values: acc };
    }
  }
}
export function std__seq__map_from$T$str$T$int(values: string[], f: (input: string) => { $can_kind: "ok"; item: bigint }, pos: bigint, n: bigint, acc: bigint[]): { $can_kind: "ok"; values: bigint[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", values: acc };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; item: bigint } = f($canSeqAt(values, pos));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        const $can_m2: { $can_kind: "ok"; values: bigint[] } = std__seq__map_from$T$str$T$int(values, f, (pos + 1n), (n - 1n), [...acc, r.item]);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const s = $can_m2;
          return { $can_kind: "ok", values: s.values };
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
      return { $can_kind: "ok", values: acc };
    }
  }
}
export function std__seq__map$T$int$T$int(values: bigint[], f: (input: bigint) => { $can_kind: "ok"; item: bigint }): { $can_kind: "ok"; values: bigint[] } {
  const $can_m1: { $can_kind: "ok"; values: bigint[] } = std__seq__map_from$T$int$T$int(values, f, 0n, (BigInt([...values].length)), []);
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
export function std__seq__map$T$str$T$int(values: string[], f: (input: string) => { $can_kind: "ok"; item: bigint }): { $can_kind: "ok"; values: bigint[] } {
  const $can_m1: { $can_kind: "ok"; values: bigint[] } = std__seq__map_from$T$str$T$int(values, f, 0n, (BigInt([...values].length)), []);
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
export function std__seq__filter_from$T$int(values: bigint[], p: (input: bigint) => { $can_kind: "ok"; value: boolean }, pos: bigint, n: bigint, acc: bigint[]): { $can_kind: "ok"; values: bigint[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", values: acc };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = p($canSeqAt(values, pos));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        if (r.value) {
          const $can_m2: { $can_kind: "ok"; values: bigint[] } = std__seq__filter_from$T$int(values, p, (pos + 1n), (n - 1n), [...acc, $canSeqAt(values, pos)]);
          switch ($can_m2.$can_kind) {
          case "ok": {
            const s = $can_m2;
            return { $can_kind: "ok", values: s.values };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          const $can_m3: { $can_kind: "ok"; values: bigint[] } = std__seq__filter_from$T$int(values, p, (pos + 1n), (n - 1n), acc);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const s = $can_m3;
            return { $can_kind: "ok", values: s.values };
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
    else {
      return { $can_kind: "ok", values: acc };
    }
  }
}
export function std__seq__filter_from$T$str(values: string[], p: (input: string) => { $can_kind: "ok"; value: boolean }, pos: bigint, n: bigint, acc: string[]): { $can_kind: "ok"; values: string[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", values: acc };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = p($canSeqAt(values, pos));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        if (r.value) {
          const $can_m2: { $can_kind: "ok"; values: string[] } = std__seq__filter_from$T$str(values, p, (pos + 1n), (n - 1n), [...acc, $canSeqAt(values, pos)]);
          switch ($can_m2.$can_kind) {
          case "ok": {
            const s = $can_m2;
            return { $can_kind: "ok", values: s.values };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          const $can_m3: { $can_kind: "ok"; values: string[] } = std__seq__filter_from$T$str(values, p, (pos + 1n), (n - 1n), acc);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const s = $can_m3;
            return { $can_kind: "ok", values: s.values };
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
    else {
      return { $can_kind: "ok", values: acc };
    }
  }
}
export function std__seq__filter$T$int(values: bigint[], p: (input: bigint) => { $can_kind: "ok"; value: boolean }): { $can_kind: "ok"; values: bigint[] } {
  const $can_m1: { $can_kind: "ok"; values: bigint[] } = std__seq__filter_from$T$int(values, p, 0n, (BigInt([...values].length)), []);
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
export function std__seq__filter$T$str(values: string[], p: (input: string) => { $can_kind: "ok"; value: boolean }): { $can_kind: "ok"; values: string[] } {
  const $can_m1: { $can_kind: "ok"; values: string[] } = std__seq__filter_from$T$str(values, p, 0n, (BigInt([...values].length)), []);
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
export function std__seq__fold_from$T$int$T$int(values: bigint[], step: (input: Seq__FoldStep$T$int$T$int) => { $can_kind: "ok"; item: bigint }, pos: bigint, n: bigint, acc: bigint): { $can_kind: "ok"; item: bigint } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", item: acc };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; item: bigint } = step({ acc: acc, elem: $canSeqAt(values, pos) });
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        const $can_m2: { $can_kind: "ok"; item: bigint } = std__seq__fold_from$T$int$T$int(values, step, (pos + 1n), (n - 1n), r.item);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const s = $can_m2;
          return { $can_kind: "ok", item: s.item };
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
      return { $can_kind: "ok", item: acc };
    }
  }
}
export function std__seq__fold_from$T$str$T$str(values: string[], step: (input: Seq__FoldStep$T$str$T$str) => { $can_kind: "ok"; item: string }, pos: bigint, n: bigint, acc: string): { $can_kind: "ok"; item: string } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", item: acc };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; item: string } = step({ acc: acc, elem: $canSeqAt(values, pos) });
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        const $can_m2: { $can_kind: "ok"; item: string } = std__seq__fold_from$T$str$T$str(values, step, (pos + 1n), (n - 1n), r.item);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const s = $can_m2;
          return { $can_kind: "ok", item: s.item };
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
      return { $can_kind: "ok", item: acc };
    }
  }
}
export function std__seq__fold$T$int$T$int(values: bigint[], initial: bigint, step: (input: Seq__FoldStep$T$int$T$int) => { $can_kind: "ok"; item: bigint }): { $can_kind: "ok"; item: bigint } {
  const $can_m1: { $can_kind: "ok"; item: bigint } = std__seq__fold_from$T$int$T$int(values, step, 0n, (BigInt([...values].length)), initial);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", item: r.item };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__seq__fold$T$str$T$str(values: string[], initial: string, step: (input: Seq__FoldStep$T$str$T$str) => { $can_kind: "ok"; item: string }): { $can_kind: "ok"; item: string } {
  const $can_m1: { $can_kind: "ok"; item: string } = std__seq__fold_from$T$str$T$str(values, step, 0n, (BigInt([...values].length)), initial);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", item: r.item };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__seq__all_from$T$int(values: bigint[], p: (input: bigint) => { $can_kind: "ok"; value: boolean }, pos: bigint, n: bigint): { $can_kind: "ok"; value: boolean } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: true };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = p($canSeqAt(values, pos));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        if (r.value) {
          const $can_m2: { $can_kind: "ok"; value: boolean } = std__seq__all_from$T$int(values, p, (pos + 1n), (n - 1n));
          switch ($can_m2.$can_kind) {
          case "ok": {
            const s = $can_m2;
            return { $can_kind: "ok", value: s.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          return { $can_kind: "ok", value: false };
        }
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $can_kind: "ok", value: true };
    }
  }
}
export function std__seq__all_from$T$str(values: string[], p: (input: string) => { $can_kind: "ok"; value: boolean }, pos: bigint, n: bigint): { $can_kind: "ok"; value: boolean } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: true };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = p($canSeqAt(values, pos));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        if (r.value) {
          const $can_m2: { $can_kind: "ok"; value: boolean } = std__seq__all_from$T$str(values, p, (pos + 1n), (n - 1n));
          switch ($can_m2.$can_kind) {
          case "ok": {
            const s = $can_m2;
            return { $can_kind: "ok", value: s.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          return { $can_kind: "ok", value: false };
        }
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $can_kind: "ok", value: true };
    }
  }
}
export function std__seq__all$T$int(values: bigint[], p: (input: bigint) => { $can_kind: "ok"; value: boolean }): { $can_kind: "ok"; value: boolean } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = std__seq__all_from$T$int(values, p, 0n, (BigInt([...values].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__seq__any_from$T$int(values: bigint[], p: (input: bigint) => { $can_kind: "ok"; value: boolean }, pos: bigint, n: bigint): { $can_kind: "ok"; value: boolean } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: false };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = p($canSeqAt(values, pos));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        if (r.value) {
          return { $can_kind: "ok", value: true };
        }
        else {
          const $can_m2: { $can_kind: "ok"; value: boolean } = std__seq__any_from$T$int(values, p, (pos + 1n), (n - 1n));
          switch ($can_m2.$can_kind) {
          case "ok": {
            const s = $can_m2;
            return { $can_kind: "ok", value: s.value };
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
    else {
      return { $can_kind: "ok", value: false };
    }
  }
}
export function std__seq__any_from$T$str(values: string[], p: (input: string) => { $can_kind: "ok"; value: boolean }, pos: bigint, n: bigint): { $can_kind: "ok"; value: boolean } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: false };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = p($canSeqAt(values, pos));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        if (r.value) {
          return { $can_kind: "ok", value: true };
        }
        else {
          const $can_m2: { $can_kind: "ok"; value: boolean } = std__seq__any_from$T$str(values, p, (pos + 1n), (n - 1n));
          switch ($can_m2.$can_kind) {
          case "ok": {
            const s = $can_m2;
            return { $can_kind: "ok", value: s.value };
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
    else {
      return { $can_kind: "ok", value: false };
    }
  }
}
export function std__seq__any$T$int(values: bigint[], p: (input: bigint) => { $can_kind: "ok"; value: boolean }): { $can_kind: "ok"; value: boolean } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = std__seq__any_from$T$int(values, p, 0n, (BigInt([...values].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__seq__find_from$T$int(values: bigint[], p: (input: bigint) => { $can_kind: "ok"; value: boolean }, pos: bigint, n: bigint): { $can_kind: "ok"; item: bigint } | { $can_kind: "sequence.not_found" } {
  if ((n <= 0n)) {
    return { $can_kind: "sequence.not_found" };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = p($canSeqAt(values, pos));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        if (r.value) {
          return { $can_kind: "ok", item: $canSeqAt(values, pos) };
        }
        else {
          const $can_m2: { $can_kind: "ok"; item: bigint } | { $can_kind: "sequence.not_found" } = std__seq__find_from$T$int(values, p, (pos + 1n), (n - 1n));
          switch ($can_m2.$can_kind) {
          case "ok": {
            const s = $can_m2;
            return { $can_kind: "ok", item: s.item };
          }
          case "sequence.not_found": {
            const _ = $can_m2;
            return { $can_kind: "sequence.not_found" };
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
    else {
      return { $can_kind: "sequence.not_found" };
    }
  }
}
export function std__seq__find_from$T$str(values: string[], p: (input: string) => { $can_kind: "ok"; value: boolean }, pos: bigint, n: bigint): { $can_kind: "ok"; item: string } | { $can_kind: "sequence.not_found" } {
  if ((n <= 0n)) {
    return { $can_kind: "sequence.not_found" };
  }
  else {
    if ((pos < (BigInt([...values].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = p($canSeqAt(values, pos));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        if (r.value) {
          return { $can_kind: "ok", item: $canSeqAt(values, pos) };
        }
        else {
          const $can_m2: { $can_kind: "ok"; item: string } | { $can_kind: "sequence.not_found" } = std__seq__find_from$T$str(values, p, (pos + 1n), (n - 1n));
          switch ($can_m2.$can_kind) {
          case "ok": {
            const s = $can_m2;
            return { $can_kind: "ok", item: s.item };
          }
          case "sequence.not_found": {
            const _ = $can_m2;
            return { $can_kind: "sequence.not_found" };
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
    else {
      return { $can_kind: "sequence.not_found" };
    }
  }
}
export function std__seq__find$T$int(values: bigint[], p: (input: bigint) => { $can_kind: "ok"; value: boolean }): { $can_kind: "ok"; item: bigint } | { $can_kind: "sequence.not_found" } {
  const $can_m1: { $can_kind: "ok"; item: bigint } | { $can_kind: "sequence.not_found" } = std__seq__find_from$T$int(values, p, 0n, (BigInt([...values].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", item: r.item };
  }
  case "sequence.not_found": {
    const _ = $can_m1;
    return { $can_kind: "sequence.not_found" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__seq__find$T$str(values: string[], p: (input: string) => { $can_kind: "ok"; value: boolean }): { $can_kind: "ok"; item: string } | { $can_kind: "sequence.not_found" } {
  const $can_m1: { $can_kind: "ok"; item: string } | { $can_kind: "sequence.not_found" } = std__seq__find_from$T$str(values, p, 0n, (BigInt([...values].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", item: r.item };
  }
  case "sequence.not_found": {
    const _ = $can_m1;
    return { $can_kind: "sequence.not_found" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
