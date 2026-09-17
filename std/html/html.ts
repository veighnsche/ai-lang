// GENERATED from html.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type HtmlResult = { $ail_kind: "ok"; attribute: string } | { $ail_kind: "ok"; name: string } | { $ail_kind: "ok"; name: string; spelling: string } | { $ail_kind: "ok"; safe: string } | { $ail_kind: "ok"; text: string } | { $ail_kind: "ok"; value: string } | { $ail_kind: "html.invalid_attribute_name"; value: string } | { $ail_kind: "html.nul_byte"; value: string };
export type Html__Escaped = { value: string };
export type Html__TextResult = { text: string };
export type Html__SafeResult = { safe: string };
export type Html__NameResult = { name: string };
export type Html__ValueEscaped = { value: string };
export type Html__NamedSpelling = { name: string; spelling: string };
export type Html__AttributeResult = { attribute: string };
export type Html__BooleanNameResult = { name: string };
export type Html__BooleanSpelling = { name: string; spelling: string };
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
export function html__text__escape_from(orig: string, s: string, acc: string, n: bigint): { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    if (($ailStrAt(s, 0n) === 0n)) {
      return { $ail_kind: "html.nul_byte", value: orig };
    }
    else {
      if (($ailStrSlice(s, 0n, 1n) === "&")) {
        const $ail_m1: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__text__escape_from(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + "&amp;"), (n - 1n));
        switch ($ail_m1.$ail_kind) {
        case "ok": {
          const r = $ail_m1;
          return { $ail_kind: "ok", value: r.value };
        }
        case "html.nul_byte": {
          const e = $ail_m1;
          return { $ail_kind: "html.nul_byte", value: e.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        if (($ailStrSlice(s, 0n, 1n) === "<")) {
          const $ail_m2: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__text__escape_from(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + "&lt;"), (n - 1n));
          switch ($ail_m2.$ail_kind) {
          case "ok": {
            const r = $ail_m2;
            return { $ail_kind: "ok", value: r.value };
          }
          case "html.nul_byte": {
            const e = $ail_m2;
            return { $ail_kind: "html.nul_byte", value: e.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          if (($ailStrSlice(s, 0n, 1n) === ">")) {
            const $ail_m3: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__text__escape_from(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + "&gt;"), (n - 1n));
            switch ($ail_m3.$ail_kind) {
            case "ok": {
              const r = $ail_m3;
              return { $ail_kind: "ok", value: r.value };
            }
            case "html.nul_byte": {
              const e = $ail_m3;
              return { $ail_kind: "html.nul_byte", value: e.value };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          else {
            const $ail_m4: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__text__escape_from(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + $ailStrSlice(s, 0n, 1n)), (n - 1n));
            switch ($ail_m4.$ail_kind) {
            case "ok": {
              const r = $ail_m4;
              return { $ail_kind: "ok", value: r.value };
            }
            case "html.nul_byte": {
              const e = $ail_m4;
              return { $ail_kind: "html.nul_byte", value: e.value };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
        }
      }
    }
  }
}
export function html__text__escape(raw: string): { $ail_kind: "ok"; text: string } | { $ail_kind: "html.nul_byte"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__text__escape_from(raw, raw, "", (BigInt([...raw].length)));
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const r = $ail_m1;
    return { $ail_kind: "ok", text: r.value };
  }
  case "html.nul_byte": {
    const e = $ail_m1;
    return { $ail_kind: "html.nul_byte", value: raw };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function html__text__node(text: string): { $ail_kind: "ok"; safe: string } {
  return { $ail_kind: "ok", safe: text };
}
export function html__attribute__name(value: string): { $ail_kind: "ok"; name: string } | { $ail_kind: "html.invalid_attribute_name"; value: string } {
  if (value === "title") {
    return { $ail_kind: "ok", name: "title" };
  }
  return { $ail_kind: "html.invalid_attribute_name", value: value };
}
export function html__attribute__value_from(orig: string, s: string, acc: string, n: bigint): { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    if (($ailStrAt(s, 0n) === 0n)) {
      return { $ail_kind: "html.nul_byte", value: orig };
    }
    else {
      if (($ailStrSlice(s, 0n, 1n) === "&")) {
        const $ail_m1: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__attribute__value_from(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + "&amp;"), (n - 1n));
        switch ($ail_m1.$ail_kind) {
        case "ok": {
          const r = $ail_m1;
          return { $ail_kind: "ok", value: r.value };
        }
        case "html.nul_byte": {
          const e = $ail_m1;
          return { $ail_kind: "html.nul_byte", value: e.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        if (($ailStrSlice(s, 0n, 1n) === "<")) {
          const $ail_m2: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__attribute__value_from(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + "&lt;"), (n - 1n));
          switch ($ail_m2.$ail_kind) {
          case "ok": {
            const r = $ail_m2;
            return { $ail_kind: "ok", value: r.value };
          }
          case "html.nul_byte": {
            const e = $ail_m2;
            return { $ail_kind: "html.nul_byte", value: e.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          if (($ailStrSlice(s, 0n, 1n) === ">")) {
            const $ail_m3: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__attribute__value_from(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + "&gt;"), (n - 1n));
            switch ($ail_m3.$ail_kind) {
            case "ok": {
              const r = $ail_m3;
              return { $ail_kind: "ok", value: r.value };
            }
            case "html.nul_byte": {
              const e = $ail_m3;
              return { $ail_kind: "html.nul_byte", value: e.value };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          else {
            if (($ailStrSlice(s, 0n, 1n) === "'")) {
              const $ail_m4: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__attribute__value_from(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + "&#39;"), (n - 1n));
              switch ($ail_m4.$ail_kind) {
              case "ok": {
                const r = $ail_m4;
                return { $ail_kind: "ok", value: r.value };
              }
              case "html.nul_byte": {
                const e = $ail_m4;
                return { $ail_kind: "html.nul_byte", value: e.value };
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            else {
              const $ail_m5: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__attribute__value_from(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + $ailStrSlice(s, 0n, 1n)), (n - 1n));
              switch ($ail_m5.$ail_kind) {
              case "ok": {
                const r = $ail_m5;
                return { $ail_kind: "ok", value: r.value };
              }
              case "html.nul_byte": {
                const e = $ail_m5;
                return { $ail_kind: "html.nul_byte", value: e.value };
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
          }
        }
      }
    }
  }
}
export function html__attribute__spelling(name: string): { $ail_kind: "ok"; name: string; spelling: string } {
  return { $ail_kind: "ok", name: name, spelling: "title" };
}
export function html__attribute__text(name: string, raw: string): { $ail_kind: "ok"; attribute: string } | { $ail_kind: "html.nul_byte"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; name: string; spelling: string } = html__attribute__spelling(name);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const w = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__attribute__value_from(raw, raw, "", (BigInt([...raw].length)));
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const v = $ail_m2;
      return { $ail_kind: "ok", attribute: (((w.spelling + "='") + v.value) + "'") };
    }
    case "html.nul_byte": {
      const e = $ail_m2;
      return { $ail_kind: "html.nul_byte", value: raw };
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
export function html__attribute__boolean_name(value: string): { $ail_kind: "ok"; name: string } | { $ail_kind: "html.invalid_attribute_name"; value: string } {
  if (value === "disabled") {
    return { $ail_kind: "ok", name: "disabled" };
  }
  return { $ail_kind: "html.invalid_attribute_name", value: value };
}
export function html__attribute__boolean_spelling(name: string): { $ail_kind: "ok"; name: string; spelling: string } {
  return { $ail_kind: "ok", name: name, spelling: "disabled" };
}
export function html__attribute__boolean(name: string, present: boolean): { $ail_kind: "ok"; attribute: string } {
  const $ail_m1: { $ail_kind: "ok"; name: string; spelling: string } = html__attribute__boolean_spelling(name);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const w = $ail_m1;
    if (present) {
      return { $ail_kind: "ok", attribute: w.spelling };
    }
    else {
      return { $ail_kind: "ok", attribute: "" };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
