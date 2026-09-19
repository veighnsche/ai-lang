// GENERATED from set.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type SetResult = { $can_kind: "ok"; members: bigint[] } | { $can_kind: "ok"; members: string[] } | { $can_kind: "ok"; value: boolean };
export type Bool__Value = { value: boolean };
export type Set__Members$T$int = { members: bigint[] };
export type Set__Members$T$str = { members: string[] };
// Sequence indexing (a38 S3): bounds throw, matching Go.
function $canSeqAt<T>(a: T[], i: bigint): T {
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("seq index out of range");
  const k = Number(i);
  if (k >= a.length) throw new Error("seq index out of range");
  return a[k];
}
export function std__set__has_from$T$int(members: bigint[], value: bigint, pos: bigint, n: bigint): { $can_kind: "ok"; value: boolean } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: false };
  }
  else {
    if ((pos < (BigInt([...members].length)))) {
      if (($canSeqAt(members, pos) === value)) {
        return { $can_kind: "ok", value: true };
      }
      else {
        const $can_m1: { $can_kind: "ok"; value: boolean } = std__set__has_from$T$int(members, value, (pos + 1n), (n - 1n));
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
      return { $can_kind: "ok", value: false };
    }
  }
}
export function std__set__has_from$T$str(members: string[], value: string, pos: bigint, n: bigint): { $can_kind: "ok"; value: boolean } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: false };
  }
  else {
    if ((pos < (BigInt([...members].length)))) {
      if (($canSeqAt(members, pos) === value)) {
        return { $can_kind: "ok", value: true };
      }
      else {
        const $can_m1: { $can_kind: "ok"; value: boolean } = std__set__has_from$T$str(members, value, (pos + 1n), (n - 1n));
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
      return { $can_kind: "ok", value: false };
    }
  }
}
export function std__set__contains$T$int(set: Set__Members$T$int, value: bigint): { $can_kind: "ok"; value: boolean } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = std__set__has_from$T$int(set.members, value, 0n, (BigInt([...set.members].length)));
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
export function std__set__contains$T$str(set: Set__Members$T$str, value: string): { $can_kind: "ok"; value: boolean } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = std__set__has_from$T$str(set.members, value, 0n, (BigInt([...set.members].length)));
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
export function std__set__union_from$T$int(right: bigint[], pos: bigint, n: bigint, acc: bigint[]): { $can_kind: "ok"; members: bigint[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", members: acc };
  }
  else {
    if ((pos < (BigInt([...right].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = std__set__has_from$T$int(acc, $canSeqAt(right, pos), 0n, (BigInt([...acc].length)));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const h = $can_m1;
        if (h.value) {
          const $can_m2: { $can_kind: "ok"; members: bigint[] } = std__set__union_from$T$int(right, (pos + 1n), (n - 1n), acc);
          switch ($can_m2.$can_kind) {
          case "ok": {
            const r = $can_m2;
            return { $can_kind: "ok", members: r.members };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          const $can_m3: { $can_kind: "ok"; members: bigint[] } = std__set__union_from$T$int(right, (pos + 1n), (n - 1n), [...acc, $canSeqAt(right, pos)]);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const r = $can_m3;
            return { $can_kind: "ok", members: r.members };
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
      return { $can_kind: "ok", members: acc };
    }
  }
}
export function std__set__union_from$T$str(right: string[], pos: bigint, n: bigint, acc: string[]): { $can_kind: "ok"; members: string[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", members: acc };
  }
  else {
    if ((pos < (BigInt([...right].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = std__set__has_from$T$str(acc, $canSeqAt(right, pos), 0n, (BigInt([...acc].length)));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const h = $can_m1;
        if (h.value) {
          const $can_m2: { $can_kind: "ok"; members: string[] } = std__set__union_from$T$str(right, (pos + 1n), (n - 1n), acc);
          switch ($can_m2.$can_kind) {
          case "ok": {
            const r = $can_m2;
            return { $can_kind: "ok", members: r.members };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          const $can_m3: { $can_kind: "ok"; members: string[] } = std__set__union_from$T$str(right, (pos + 1n), (n - 1n), [...acc, $canSeqAt(right, pos)]);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const r = $can_m3;
            return { $can_kind: "ok", members: r.members };
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
      return { $can_kind: "ok", members: acc };
    }
  }
}
export function std__set__union$T$int(left: Set__Members$T$int, right: Set__Members$T$int): { $can_kind: "ok"; members: bigint[] } {
  const $can_m1: { $can_kind: "ok"; members: bigint[] } = std__set__union_from$T$int(right.members, 0n, (BigInt([...right.members].length)), left.members);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", members: r.members };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__set__union$T$str(left: Set__Members$T$str, right: Set__Members$T$str): { $can_kind: "ok"; members: string[] } {
  const $can_m1: { $can_kind: "ok"; members: string[] } = std__set__union_from$T$str(right.members, 0n, (BigInt([...right.members].length)), left.members);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", members: r.members };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__set__intersection_from$T$int(left: bigint[], right: bigint[], pos: bigint, n: bigint, acc: bigint[]): { $can_kind: "ok"; members: bigint[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", members: acc };
  }
  else {
    if ((pos < (BigInt([...left].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = std__set__has_from$T$int(right, $canSeqAt(left, pos), 0n, (BigInt([...right].length)));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const h = $can_m1;
        if (h.value) {
          const $can_m2: { $can_kind: "ok"; members: bigint[] } = std__set__intersection_from$T$int(left, right, (pos + 1n), (n - 1n), [...acc, $canSeqAt(left, pos)]);
          switch ($can_m2.$can_kind) {
          case "ok": {
            const r = $can_m2;
            return { $can_kind: "ok", members: r.members };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          const $can_m3: { $can_kind: "ok"; members: bigint[] } = std__set__intersection_from$T$int(left, right, (pos + 1n), (n - 1n), acc);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const r = $can_m3;
            return { $can_kind: "ok", members: r.members };
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
      return { $can_kind: "ok", members: acc };
    }
  }
}
export function std__set__intersection_from$T$str(left: string[], right: string[], pos: bigint, n: bigint, acc: string[]): { $can_kind: "ok"; members: string[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", members: acc };
  }
  else {
    if ((pos < (BigInt([...left].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = std__set__has_from$T$str(right, $canSeqAt(left, pos), 0n, (BigInt([...right].length)));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const h = $can_m1;
        if (h.value) {
          const $can_m2: { $can_kind: "ok"; members: string[] } = std__set__intersection_from$T$str(left, right, (pos + 1n), (n - 1n), [...acc, $canSeqAt(left, pos)]);
          switch ($can_m2.$can_kind) {
          case "ok": {
            const r = $can_m2;
            return { $can_kind: "ok", members: r.members };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          const $can_m3: { $can_kind: "ok"; members: string[] } = std__set__intersection_from$T$str(left, right, (pos + 1n), (n - 1n), acc);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const r = $can_m3;
            return { $can_kind: "ok", members: r.members };
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
      return { $can_kind: "ok", members: acc };
    }
  }
}
export function std__set__intersection$T$int(left: Set__Members$T$int, right: Set__Members$T$int): { $can_kind: "ok"; members: bigint[] } {
  const $can_m1: { $can_kind: "ok"; members: bigint[] } = std__set__intersection_from$T$int(left.members, right.members, 0n, (BigInt([...left.members].length)), []);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", members: r.members };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__set__intersection$T$str(left: Set__Members$T$str, right: Set__Members$T$str): { $can_kind: "ok"; members: string[] } {
  const $can_m1: { $can_kind: "ok"; members: string[] } = std__set__intersection_from$T$str(left.members, right.members, 0n, (BigInt([...left.members].length)), []);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", members: r.members };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__set__difference_from$T$int(left: bigint[], right: bigint[], pos: bigint, n: bigint, acc: bigint[]): { $can_kind: "ok"; members: bigint[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", members: acc };
  }
  else {
    if ((pos < (BigInt([...left].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = std__set__has_from$T$int(right, $canSeqAt(left, pos), 0n, (BigInt([...right].length)));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const h = $can_m1;
        if (h.value) {
          const $can_m2: { $can_kind: "ok"; members: bigint[] } = std__set__difference_from$T$int(left, right, (pos + 1n), (n - 1n), acc);
          switch ($can_m2.$can_kind) {
          case "ok": {
            const r = $can_m2;
            return { $can_kind: "ok", members: r.members };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          const $can_m3: { $can_kind: "ok"; members: bigint[] } = std__set__difference_from$T$int(left, right, (pos + 1n), (n - 1n), [...acc, $canSeqAt(left, pos)]);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const r = $can_m3;
            return { $can_kind: "ok", members: r.members };
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
      return { $can_kind: "ok", members: acc };
    }
  }
}
export function std__set__difference_from$T$str(left: string[], right: string[], pos: bigint, n: bigint, acc: string[]): { $can_kind: "ok"; members: string[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", members: acc };
  }
  else {
    if ((pos < (BigInt([...left].length)))) {
      const $can_m1: { $can_kind: "ok"; value: boolean } = std__set__has_from$T$str(right, $canSeqAt(left, pos), 0n, (BigInt([...right].length)));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const h = $can_m1;
        if (h.value) {
          const $can_m2: { $can_kind: "ok"; members: string[] } = std__set__difference_from$T$str(left, right, (pos + 1n), (n - 1n), acc);
          switch ($can_m2.$can_kind) {
          case "ok": {
            const r = $can_m2;
            return { $can_kind: "ok", members: r.members };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          const $can_m3: { $can_kind: "ok"; members: string[] } = std__set__difference_from$T$str(left, right, (pos + 1n), (n - 1n), [...acc, $canSeqAt(left, pos)]);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const r = $can_m3;
            return { $can_kind: "ok", members: r.members };
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
      return { $can_kind: "ok", members: acc };
    }
  }
}
export function std__set__difference$T$int(left: Set__Members$T$int, right: Set__Members$T$int): { $can_kind: "ok"; members: bigint[] } {
  const $can_m1: { $can_kind: "ok"; members: bigint[] } = std__set__difference_from$T$int(left.members, right.members, 0n, (BigInt([...left.members].length)), []);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", members: r.members };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__set__difference$T$str(left: Set__Members$T$str, right: Set__Members$T$str): { $can_kind: "ok"; members: string[] } {
  const $can_m1: { $can_kind: "ok"; members: string[] } = std__set__difference_from$T$str(left.members, right.members, 0n, (BigInt([...left.members].length)), []);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", members: r.members };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
