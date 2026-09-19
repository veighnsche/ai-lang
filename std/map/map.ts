// GENERATED from map.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type MapResult = { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } | { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } | { $can_kind: "ok"; item: bigint } | { $can_kind: "ok"; item: string } | { $can_kind: "ok"; value: bigint } | { $can_kind: "map.key_absent" } | { $can_kind: "map.key_exists" };
export type Int__Value = { value: bigint };
export type Map__Pair$T$int$T$str = { key: bigint; val: string };
export type Map__Pair$T$str$T$int = { key: string; val: bigint };
export type Map$T$int$T$str = { entries: Map__Pair$T$int$T$str[] };
export type Map$T$str$T$int = { entries: Map__Pair$T$str$T$int[] };
export type Map__Item$T$int = { item: bigint };
export type Map__Item$T$str = { item: string };
// Sequence indexing (a38 S3): bounds throw, matching Go.
function $canSeqAt<T>(a: T[], i: bigint): T {
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("seq index out of range");
  const k = Number(i);
  if (k >= a.length) throw new Error("seq index out of range");
  return a[k];
}
export function std__map__index_from$T$int$T$str(map: Map$T$int$T$str, key: bigint, pos: bigint, n: bigint): { $can_kind: "ok"; value: bigint } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: -1n };
  }
  else {
    if ((pos < (BigInt([...map.entries].length)))) {
      if (($canSeqAt(map.entries, pos).key === key)) {
        return { $can_kind: "ok", value: pos };
      }
      else {
        const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$int$T$str(map, key, (pos + 1n), (n - 1n));
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
    }
    else {
      return { $can_kind: "ok", value: -1n };
    }
  }
}
export function std__map__index_from$T$str$T$int(map: Map$T$str$T$int, key: string, pos: bigint, n: bigint): { $can_kind: "ok"; value: bigint } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: -1n };
  }
  else {
    if ((pos < (BigInt([...map.entries].length)))) {
      if (($canSeqAt(map.entries, pos).key === key)) {
        return { $can_kind: "ok", value: pos };
      }
      else {
        const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$str$T$int(map, key, (pos + 1n), (n - 1n));
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
    }
    else {
      return { $can_kind: "ok", value: -1n };
    }
  }
}
export function std__map__get$T$int$T$str(map: Map$T$int$T$str, key: bigint): { $can_kind: "ok"; item: string } | { $can_kind: "map.key_absent" } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$int$T$str(map, key, 0n, (BigInt([...map.entries].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent" };
    }
    else {
      return { $can_kind: "ok", item: $canSeqAt(map.entries, i.value).val };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__map__get$T$str$T$int(map: Map$T$str$T$int, key: string): { $can_kind: "ok"; item: bigint } | { $can_kind: "map.key_absent" } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$str$T$int(map, key, 0n, (BigInt([...map.entries].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent" };
    }
    else {
      return { $can_kind: "ok", item: $canSeqAt(map.entries, i.value).val };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__map__insert$T$int$T$str(map: Map$T$int$T$str, key: bigint, value: string): { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } | { $can_kind: "map.key_exists" } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$int$T$str(map, key, 0n, (BigInt([...map.entries].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "ok", entries: [...map.entries, { key: key, val: value }] };
    }
    else {
      return { $can_kind: "map.key_exists" };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__map__insert$T$str$T$int(map: Map$T$str$T$int, key: string, value: bigint): { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } | { $can_kind: "map.key_exists" } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$str$T$int(map, key, 0n, (BigInt([...map.entries].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "ok", entries: [...map.entries, { key: key, val: value }] };
    }
    else {
      return { $can_kind: "map.key_exists" };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__map__replace_from$T$int$T$str(map: Map$T$int$T$str, key: bigint, value: string, pos: bigint, n: bigint, acc: Map__Pair$T$int$T$str[]): { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", entries: acc };
  }
  else {
    if ((pos < (BigInt([...map.entries].length)))) {
      if (($canSeqAt(map.entries, pos).key === key)) {
        const $can_m1: { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } = std__map__replace_from$T$int$T$str(map, key, value, (pos + 1n), (n - 1n), [...acc, { key: key, val: value }]);
        switch ($can_m1.$can_kind) {
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", entries: r.entries };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $can_m2: { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } = std__map__replace_from$T$int$T$str(map, key, value, (pos + 1n), (n - 1n), [...acc, $canSeqAt(map.entries, pos)]);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", entries: r.entries };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    else {
      return { $can_kind: "ok", entries: acc };
    }
  }
}
export function std__map__replace_from$T$str$T$int(map: Map$T$str$T$int, key: string, value: bigint, pos: bigint, n: bigint, acc: Map__Pair$T$str$T$int[]): { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", entries: acc };
  }
  else {
    if ((pos < (BigInt([...map.entries].length)))) {
      if (($canSeqAt(map.entries, pos).key === key)) {
        const $can_m1: { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } = std__map__replace_from$T$str$T$int(map, key, value, (pos + 1n), (n - 1n), [...acc, { key: key, val: value }]);
        switch ($can_m1.$can_kind) {
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", entries: r.entries };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $can_m2: { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } = std__map__replace_from$T$str$T$int(map, key, value, (pos + 1n), (n - 1n), [...acc, $canSeqAt(map.entries, pos)]);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", entries: r.entries };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    else {
      return { $can_kind: "ok", entries: acc };
    }
  }
}
export function std__map__replace$T$int$T$str(map: Map$T$int$T$str, key: bigint, value: string): { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } | { $can_kind: "map.key_absent" } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$int$T$str(map, key, 0n, (BigInt([...map.entries].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent" };
    }
    else {
      const $can_m2: { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } = std__map__replace_from$T$int$T$str(map, key, value, 0n, (BigInt([...map.entries].length)), []);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", entries: r.entries };
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
export function std__map__replace$T$str$T$int(map: Map$T$str$T$int, key: string, value: bigint): { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } | { $can_kind: "map.key_absent" } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$str$T$int(map, key, 0n, (BigInt([...map.entries].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent" };
    }
    else {
      const $can_m2: { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } = std__map__replace_from$T$str$T$int(map, key, value, 0n, (BigInt([...map.entries].length)), []);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", entries: r.entries };
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
export function std__map__remove_from$T$int$T$str(map: Map$T$int$T$str, key: bigint, pos: bigint, n: bigint, acc: Map__Pair$T$int$T$str[]): { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", entries: acc };
  }
  else {
    if ((pos < (BigInt([...map.entries].length)))) {
      if (($canSeqAt(map.entries, pos).key === key)) {
        const $can_m1: { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } = std__map__remove_from$T$int$T$str(map, key, (pos + 1n), (n - 1n), acc);
        switch ($can_m1.$can_kind) {
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", entries: r.entries };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $can_m2: { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } = std__map__remove_from$T$int$T$str(map, key, (pos + 1n), (n - 1n), [...acc, $canSeqAt(map.entries, pos)]);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", entries: r.entries };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    else {
      return { $can_kind: "ok", entries: acc };
    }
  }
}
export function std__map__remove_from$T$str$T$int(map: Map$T$str$T$int, key: string, pos: bigint, n: bigint, acc: Map__Pair$T$str$T$int[]): { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", entries: acc };
  }
  else {
    if ((pos < (BigInt([...map.entries].length)))) {
      if (($canSeqAt(map.entries, pos).key === key)) {
        const $can_m1: { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } = std__map__remove_from$T$str$T$int(map, key, (pos + 1n), (n - 1n), acc);
        switch ($can_m1.$can_kind) {
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", entries: r.entries };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $can_m2: { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } = std__map__remove_from$T$str$T$int(map, key, (pos + 1n), (n - 1n), [...acc, $canSeqAt(map.entries, pos)]);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", entries: r.entries };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    else {
      return { $can_kind: "ok", entries: acc };
    }
  }
}
export function std__map__remove$T$int$T$str(map: Map$T$int$T$str, key: bigint): { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } | { $can_kind: "map.key_absent" } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$int$T$str(map, key, 0n, (BigInt([...map.entries].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent" };
    }
    else {
      const $can_m2: { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } = std__map__remove_from$T$int$T$str(map, key, 0n, (BigInt([...map.entries].length)), []);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", entries: r.entries };
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
export function std__map__remove$T$str$T$int(map: Map$T$str$T$int, key: string): { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } | { $can_kind: "map.key_absent" } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$str$T$int(map, key, 0n, (BigInt([...map.entries].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent" };
    }
    else {
      const $can_m2: { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } = std__map__remove_from$T$str$T$int(map, key, 0n, (BigInt([...map.entries].length)), []);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", entries: r.entries };
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
export function std__map__entries$T$int$T$str(map: Map$T$int$T$str): { $can_kind: "ok"; entries: Map__Pair$T$int$T$str[] } {
  return { $can_kind: "ok", entries: map.entries };
}
export function std__map__entries$T$str$T$int(map: Map$T$str$T$int): { $can_kind: "ok"; entries: Map__Pair$T$str$T$int[] } {
  return { $can_kind: "ok", entries: map.entries };
}
