// GENERATED from map.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type MapResult = { $can_kind: "ok"; item: bigint } | { $can_kind: "ok"; item: string } | { $can_kind: "ok"; keys: string[]; values: bigint[] } | { $can_kind: "ok"; keys: string[]; values: string[] } | { $can_kind: "ok"; value: bigint } | { $can_kind: "map.key_absent"; key: string } | { $can_kind: "map.key_exists"; key: string };
export type Int__Value = { value: bigint };
export type Map__Entries$T$int = { keys: string[]; values: bigint[] };
export type Map__Entries$T$str = { keys: string[]; values: string[] };
export type Map__Item$T$int = { item: bigint };
export type Map__Item$T$str = { item: string };
// Sequence indexing (a38 S3): bounds throw, matching Go.
function $canSeqAt<T>(a: T[], i: bigint): T {
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("seq index out of range");
  const k = Number(i);
  if (k >= a.length) throw new Error("seq index out of range");
  return a[k];
}
export function std__map__index_from$T$int(map: Map__Entries$T$int, key: string, pos: bigint, n: bigint): { $can_kind: "ok"; value: bigint } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: -1n };
  }
  else {
    if ((pos < (BigInt([...map.keys].length)))) {
      if (($canSeqAt(map.keys, pos) === key)) {
        return { $can_kind: "ok", value: pos };
      }
      else {
        const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$int(map, key, (pos + 1n), (n - 1n));
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
export function std__map__index_from$T$str(map: Map__Entries$T$str, key: string, pos: bigint, n: bigint): { $can_kind: "ok"; value: bigint } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: -1n };
  }
  else {
    if ((pos < (BigInt([...map.keys].length)))) {
      if (($canSeqAt(map.keys, pos) === key)) {
        return { $can_kind: "ok", value: pos };
      }
      else {
        const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$str(map, key, (pos + 1n), (n - 1n));
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
export function std__map__get$T$int(map: Map__Entries$T$int, key: string): { $can_kind: "ok"; item: bigint } | { $can_kind: "map.key_absent"; key: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$int(map, key, 0n, (BigInt([...map.keys].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent", key: key };
    }
    else {
      return { $can_kind: "ok", item: $canSeqAt(map.values, i.value) };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__map__get$T$str(map: Map__Entries$T$str, key: string): { $can_kind: "ok"; item: string } | { $can_kind: "map.key_absent"; key: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$str(map, key, 0n, (BigInt([...map.keys].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent", key: key };
    }
    else {
      return { $can_kind: "ok", item: $canSeqAt(map.values, i.value) };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__map__insert$T$int(map: Map__Entries$T$int, key: string, value: bigint): { $can_kind: "ok"; keys: string[]; values: bigint[] } | { $can_kind: "map.key_exists"; key: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$int(map, key, 0n, (BigInt([...map.keys].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "ok", keys: [...map.keys, key], values: [...map.values, value] };
    }
    else {
      return { $can_kind: "map.key_exists", key: key };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__map__insert$T$str(map: Map__Entries$T$str, key: string, value: string): { $can_kind: "ok"; keys: string[]; values: string[] } | { $can_kind: "map.key_exists"; key: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$str(map, key, 0n, (BigInt([...map.keys].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "ok", keys: [...map.keys, key], values: [...map.values, value] };
    }
    else {
      return { $can_kind: "map.key_exists", key: key };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__map__replace_from$T$int(map: Map__Entries$T$int, key: string, value: bigint, pos: bigint, n: bigint, acc_keys: string[], acc_values: bigint[]): { $can_kind: "ok"; keys: string[]; values: bigint[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", keys: acc_keys, values: acc_values };
  }
  else {
    if ((pos < (BigInt([...map.keys].length)))) {
      if (($canSeqAt(map.keys, pos) === key)) {
        const $can_m1: { $can_kind: "ok"; keys: string[]; values: bigint[] } = std__map__replace_from$T$int(map, key, value, (pos + 1n), (n - 1n), [...acc_keys, $canSeqAt(map.keys, pos)], [...acc_values, value]);
        switch ($can_m1.$can_kind) {
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", keys: r.keys, values: r.values };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $can_m2: { $can_kind: "ok"; keys: string[]; values: bigint[] } = std__map__replace_from$T$int(map, key, value, (pos + 1n), (n - 1n), [...acc_keys, $canSeqAt(map.keys, pos)], [...acc_values, $canSeqAt(map.values, pos)]);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", keys: r.keys, values: r.values };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    else {
      return { $can_kind: "ok", keys: acc_keys, values: acc_values };
    }
  }
}
export function std__map__replace_from$T$str(map: Map__Entries$T$str, key: string, value: string, pos: bigint, n: bigint, acc_keys: string[], acc_values: string[]): { $can_kind: "ok"; keys: string[]; values: string[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", keys: acc_keys, values: acc_values };
  }
  else {
    if ((pos < (BigInt([...map.keys].length)))) {
      if (($canSeqAt(map.keys, pos) === key)) {
        const $can_m1: { $can_kind: "ok"; keys: string[]; values: string[] } = std__map__replace_from$T$str(map, key, value, (pos + 1n), (n - 1n), [...acc_keys, $canSeqAt(map.keys, pos)], [...acc_values, value]);
        switch ($can_m1.$can_kind) {
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", keys: r.keys, values: r.values };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $can_m2: { $can_kind: "ok"; keys: string[]; values: string[] } = std__map__replace_from$T$str(map, key, value, (pos + 1n), (n - 1n), [...acc_keys, $canSeqAt(map.keys, pos)], [...acc_values, $canSeqAt(map.values, pos)]);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", keys: r.keys, values: r.values };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    else {
      return { $can_kind: "ok", keys: acc_keys, values: acc_values };
    }
  }
}
export function std__map__replace$T$int(map: Map__Entries$T$int, key: string, value: bigint): { $can_kind: "ok"; keys: string[]; values: bigint[] } | { $can_kind: "map.key_absent"; key: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$int(map, key, 0n, (BigInt([...map.keys].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent", key: key };
    }
    else {
      const $can_m2: { $can_kind: "ok"; keys: string[]; values: bigint[] } = std__map__replace_from$T$int(map, key, value, 0n, (BigInt([...map.keys].length)), [], []);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", keys: r.keys, values: r.values };
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
export function std__map__replace$T$str(map: Map__Entries$T$str, key: string, value: string): { $can_kind: "ok"; keys: string[]; values: string[] } | { $can_kind: "map.key_absent"; key: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$str(map, key, 0n, (BigInt([...map.keys].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent", key: key };
    }
    else {
      const $can_m2: { $can_kind: "ok"; keys: string[]; values: string[] } = std__map__replace_from$T$str(map, key, value, 0n, (BigInt([...map.keys].length)), [], []);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", keys: r.keys, values: r.values };
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
export function std__map__remove_from$T$int(map: Map__Entries$T$int, key: string, pos: bigint, n: bigint, acc_keys: string[], acc_values: bigint[]): { $can_kind: "ok"; keys: string[]; values: bigint[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", keys: acc_keys, values: acc_values };
  }
  else {
    if ((pos < (BigInt([...map.keys].length)))) {
      if (($canSeqAt(map.keys, pos) === key)) {
        const $can_m1: { $can_kind: "ok"; keys: string[]; values: bigint[] } = std__map__remove_from$T$int(map, key, (pos + 1n), (n - 1n), acc_keys, acc_values);
        switch ($can_m1.$can_kind) {
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", keys: r.keys, values: r.values };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $can_m2: { $can_kind: "ok"; keys: string[]; values: bigint[] } = std__map__remove_from$T$int(map, key, (pos + 1n), (n - 1n), [...acc_keys, $canSeqAt(map.keys, pos)], [...acc_values, $canSeqAt(map.values, pos)]);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", keys: r.keys, values: r.values };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    else {
      return { $can_kind: "ok", keys: acc_keys, values: acc_values };
    }
  }
}
export function std__map__remove_from$T$str(map: Map__Entries$T$str, key: string, pos: bigint, n: bigint, acc_keys: string[], acc_values: string[]): { $can_kind: "ok"; keys: string[]; values: string[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", keys: acc_keys, values: acc_values };
  }
  else {
    if ((pos < (BigInt([...map.keys].length)))) {
      if (($canSeqAt(map.keys, pos) === key)) {
        const $can_m1: { $can_kind: "ok"; keys: string[]; values: string[] } = std__map__remove_from$T$str(map, key, (pos + 1n), (n - 1n), acc_keys, acc_values);
        switch ($can_m1.$can_kind) {
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", keys: r.keys, values: r.values };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        const $can_m2: { $can_kind: "ok"; keys: string[]; values: string[] } = std__map__remove_from$T$str(map, key, (pos + 1n), (n - 1n), [...acc_keys, $canSeqAt(map.keys, pos)], [...acc_values, $canSeqAt(map.values, pos)]);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", keys: r.keys, values: r.values };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    else {
      return { $can_kind: "ok", keys: acc_keys, values: acc_values };
    }
  }
}
export function std__map__remove$T$int(map: Map__Entries$T$int, key: string): { $can_kind: "ok"; keys: string[]; values: bigint[] } | { $can_kind: "map.key_absent"; key: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$int(map, key, 0n, (BigInt([...map.keys].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent", key: key };
    }
    else {
      const $can_m2: { $can_kind: "ok"; keys: string[]; values: bigint[] } = std__map__remove_from$T$int(map, key, 0n, (BigInt([...map.keys].length)), [], []);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", keys: r.keys, values: r.values };
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
export function std__map__remove$T$str(map: Map__Entries$T$str, key: string): { $can_kind: "ok"; keys: string[]; values: string[] } | { $can_kind: "map.key_absent"; key: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__map__index_from$T$str(map, key, 0n, (BigInt([...map.keys].length)));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const i = $can_m1;
    if ((i.value < 0n)) {
      return { $can_kind: "map.key_absent", key: key };
    }
    else {
      const $can_m2: { $can_kind: "ok"; keys: string[]; values: string[] } = std__map__remove_from$T$str(map, key, 0n, (BigInt([...map.keys].length)), [], []);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const r = $can_m2;
        return { $can_kind: "ok", keys: r.keys, values: r.values };
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
export function std__map__entries$T$int(map: Map__Entries$T$int): { $can_kind: "ok"; keys: string[]; values: bigint[] } {
  return { $can_kind: "ok", keys: map.keys, values: map.values };
}
export function std__map__entries$T$str(map: Map__Entries$T$str): { $can_kind: "ok"; keys: string[]; values: string[] } {
  return { $can_kind: "ok", keys: map.keys, values: map.values };
}
