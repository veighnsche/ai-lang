// GENERATED from html.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type HtmlResult = { $ail_kind: "ok"; attribute: string } | { $ail_kind: "ok"; has: boolean } | { $ail_kind: "ok"; len: bigint; value: string } | { $ail_kind: "ok"; n: bigint; tail: string; value: string } | { $ail_kind: "ok"; name: string } | { $ail_kind: "ok"; name: string; spelling: string } | { $ail_kind: "ok"; safe: string } | { $ail_kind: "ok"; text: string } | { $ail_kind: "ok"; value: string } | { $ail_kind: "ok"; ws: boolean } | { $ail_kind: "html.invalid_attribute_name"; value: string } | { $ail_kind: "html.nul_byte"; value: string } | { $ail_kind: "html.invalid_identifier"; value: string } | { $ail_kind: "html.invalid_url"; value: string } | { $ail_kind: "html.disallowed_scheme"; value: string };
export type Html__Escaped = { value: string };
export type Html__TextResult = { text: string };
export type Html__SafeResult = { safe: string };
export type Html__NameResult = { name: string };
export type Html__ValueEscaped = { value: string };
export type Html__NamedSpelling = { name: string; spelling: string };
export type Html__AttributeResult = { attribute: string };
export type Html__BooleanNameResult = { name: string };
export type Html__BooleanSpelling = { name: string; spelling: string };
export type Html__WsVerdict = { ws: boolean };
export type Html__IdValid = { value: string };
export type Html__SchemeVerdict = { has: boolean };
export type Html__Authority = { value: string; tail: string; n: bigint };
export type Html__UrlTail = { value: string };
export type Html__CheckedUrl = { value: string; len: bigint };
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
  else if (value === "readonly") {
    return { $ail_kind: "ok", name: "readonly" };
  }
  else if (value === "required") {
    return { $ail_kind: "ok", name: "required" };
  }
  else if (value === "checked") {
    return { $ail_kind: "ok", name: "checked" };
  }
  return { $ail_kind: "html.invalid_attribute_name", value: value };
}
export function html__attribute__boolean_spelling(name: string): { $ail_kind: "ok"; name: string; spelling: string } {
  if ((name === "disabled")) {
    return { $ail_kind: "ok", name: name, spelling: "disabled" };
  }
  else {
    if ((name === "readonly")) {
      return { $ail_kind: "ok", name: name, spelling: "readonly" };
    }
    else {
      if ((name === "required")) {
        return { $ail_kind: "ok", name: name, spelling: "required" };
      }
      else {
        return { $ail_kind: "ok", name: name, spelling: "checked" };
      }
    }
  }
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
export function html__attribute__id_ws(code: bigint): { $ail_kind: "ok"; ws: boolean } {
  if ((code === 9n)) {
    return { $ail_kind: "ok", ws: true };
  }
  else {
    if ((code === 10n)) {
      return { $ail_kind: "ok", ws: true };
    }
    else {
      if ((code === 12n)) {
        return { $ail_kind: "ok", ws: true };
      }
      else {
        if ((code === 13n)) {
          return { $ail_kind: "ok", ws: true };
        }
        else {
          if ((code === 32n)) {
            return { $ail_kind: "ok", ws: true };
          }
          else {
            return { $ail_kind: "ok", ws: false };
          }
        }
      }
    }
  }
}
export function html__attribute__id_check(orig: string, s: string, n: bigint): { $ail_kind: "ok"; value: string } | { $ail_kind: "html.invalid_identifier"; value: string } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: orig };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; ws: boolean } = html__attribute__id_ws($ailStrAt(s, 0n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const c = $ail_m1;
      if (c.ws) {
        return { $ail_kind: "html.invalid_identifier", value: orig };
      }
      else {
        const $ail_m2: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.invalid_identifier"; value: string } = html__attribute__id_check(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n));
        switch ($ail_m2.$ail_kind) {
        case "ok": {
          const r = $ail_m2;
          return { $ail_kind: "ok", value: r.value };
        }
        case "html.invalid_identifier": {
          const e = $ail_m2;
          return { $ail_kind: "html.invalid_identifier", value: e.value };
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
export function html__attribute__id(value: string): { $ail_kind: "ok"; attribute: string } | { $ail_kind: "html.invalid_identifier"; value: string } | { $ail_kind: "html.nul_byte"; value: string } {
  if (((BigInt([...value].length)) === 0n)) {
    return { $ail_kind: "html.invalid_identifier", value: value };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.invalid_identifier"; value: string } = html__attribute__id_check(value, value, (BigInt([...value].length)));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const c = $ail_m1;
      const $ail_m2: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__attribute__value_from(value, value, "", (BigInt([...value].length)));
      switch ($ail_m2.$ail_kind) {
      case "ok": {
        const v = $ail_m2;
        return { $ail_kind: "ok", attribute: (("id='" + v.value) + "'") };
      }
      case "html.nul_byte": {
        const e = $ail_m2;
        return { $ail_kind: "html.nul_byte", value: value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    case "html.invalid_identifier": {
      const e = $ail_m1;
      return { $ail_kind: "html.invalid_identifier", value: e.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function html__url__scheme_token(orig: string, s: string, n: bigint, total: bigint): { $ail_kind: "ok"; has: boolean } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", has: false };
  }
  else {
    if ((n === total)) {
      if (($ailStrAt(s, 0n) <= 64n)) {
        return { $ail_kind: "ok", has: false };
      }
      else {
        if (($ailStrAt(s, 0n) <= 90n)) {
          const $ail_m1: { $ail_kind: "ok"; has: boolean } = html__url__scheme_token(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), total);
          switch ($ail_m1.$ail_kind) {
          case "ok": {
            const r = $ail_m1;
            return { $ail_kind: "ok", has: r.has };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          if (($ailStrAt(s, 0n) <= 96n)) {
            return { $ail_kind: "ok", has: false };
          }
          else {
            if (($ailStrAt(s, 0n) <= 122n)) {
              const $ail_m2: { $ail_kind: "ok"; has: boolean } = html__url__scheme_token(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), total);
              switch ($ail_m2.$ail_kind) {
              case "ok": {
                const r = $ail_m2;
                return { $ail_kind: "ok", has: r.has };
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            else {
              return { $ail_kind: "ok", has: false };
            }
          }
        }
      }
    }
    else {
      if (($ailStrAt(s, 0n) === 58n)) {
        return { $ail_kind: "ok", has: true };
      }
      else {
        if (($ailStrAt(s, 0n) === 43n)) {
          const $ail_m3: { $ail_kind: "ok"; has: boolean } = html__url__scheme_token(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), total);
          switch ($ail_m3.$ail_kind) {
          case "ok": {
            const r = $ail_m3;
            return { $ail_kind: "ok", has: r.has };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          if (($ailStrAt(s, 0n) === 45n)) {
            const $ail_m4: { $ail_kind: "ok"; has: boolean } = html__url__scheme_token(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), total);
            switch ($ail_m4.$ail_kind) {
            case "ok": {
              const r = $ail_m4;
              return { $ail_kind: "ok", has: r.has };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          else {
            if (($ailStrAt(s, 0n) === 46n)) {
              const $ail_m5: { $ail_kind: "ok"; has: boolean } = html__url__scheme_token(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), total);
              switch ($ail_m5.$ail_kind) {
              case "ok": {
                const r = $ail_m5;
                return { $ail_kind: "ok", has: r.has };
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            else {
              if (($ailStrAt(s, 0n) <= 47n)) {
                return { $ail_kind: "ok", has: false };
              }
              else {
                if (($ailStrAt(s, 0n) <= 57n)) {
                  const $ail_m6: { $ail_kind: "ok"; has: boolean } = html__url__scheme_token(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), total);
                  switch ($ail_m6.$ail_kind) {
                  case "ok": {
                    const r = $ail_m6;
                    return { $ail_kind: "ok", has: r.has };
                  }
                  default: {
                    throw new Error("unreachable");
                  }
                  }
                }
                else {
                  if (($ailStrAt(s, 0n) <= 64n)) {
                    return { $ail_kind: "ok", has: false };
                  }
                  else {
                    if (($ailStrAt(s, 0n) <= 90n)) {
                      const $ail_m7: { $ail_kind: "ok"; has: boolean } = html__url__scheme_token(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), total);
                      switch ($ail_m7.$ail_kind) {
                      case "ok": {
                        const r = $ail_m7;
                        return { $ail_kind: "ok", has: r.has };
                      }
                      default: {
                        throw new Error("unreachable");
                      }
                      }
                    }
                    else {
                      if (($ailStrAt(s, 0n) <= 96n)) {
                        return { $ail_kind: "ok", has: false };
                      }
                      else {
                        if (($ailStrAt(s, 0n) <= 122n)) {
                          const $ail_m8: { $ail_kind: "ok"; has: boolean } = html__url__scheme_token(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), total);
                          switch ($ail_m8.$ail_kind) {
                          case "ok": {
                            const r = $ail_m8;
                            return { $ail_kind: "ok", has: r.has };
                          }
                          default: {
                            throw new Error("unreachable");
                          }
                          }
                        }
                        else {
                          return { $ail_kind: "ok", has: false };
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
export function html__url__authority(orig: string, s: string, n: bigint, prev: string): { $ail_kind: "ok"; n: bigint; tail: string; value: string } | { $ail_kind: "html.invalid_url"; value: string } {
  if ((n <= 0n)) {
    if ((prev === "")) {
      return { $ail_kind: "html.invalid_url", value: orig };
    }
    else {
      if ((prev === "-")) {
        return { $ail_kind: "html.invalid_url", value: orig };
      }
      else {
        return { $ail_kind: "ok", value: orig, tail: s, n: n };
      }
    }
  }
  else {
    if (($ailStrAt(s, 0n) === 47n)) {
      if ((prev === "-")) {
        return { $ail_kind: "html.invalid_url", value: orig };
      }
      else {
        return { $ail_kind: "ok", value: orig, tail: s, n: n };
      }
    }
    else {
      if (($ailStrAt(s, 0n) === 63n)) {
        if ((prev === "-")) {
          return { $ail_kind: "html.invalid_url", value: orig };
        }
        else {
          return { $ail_kind: "ok", value: orig, tail: s, n: n };
        }
      }
      else {
        if (($ailStrAt(s, 0n) === 35n)) {
          if ((prev === "-")) {
            return { $ail_kind: "html.invalid_url", value: orig };
          }
          else {
            return { $ail_kind: "ok", value: orig, tail: s, n: n };
          }
        }
        else {
          if (($ailStrAt(s, 0n) === 0n)) {
            const $ail_m1: { $ail_kind: "ok"; n: bigint; tail: string; value: string } | { $ail_kind: "html.invalid_url"; value: string } = html__url__authority(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), prev);
            switch ($ail_m1.$ail_kind) {
            case "ok": {
              const r = $ail_m1;
              return { $ail_kind: "ok", value: r.value, tail: r.tail, n: r.n };
            }
            case "html.invalid_url": {
              const e = $ail_m1;
              return { $ail_kind: "html.invalid_url", value: e.value };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          else {
            if (($ailStrAt(s, 0n) === 46n)) {
              if ((prev === "")) {
                return { $ail_kind: "html.invalid_url", value: orig };
              }
              else {
                if ((prev === ".")) {
                  return { $ail_kind: "html.invalid_url", value: orig };
                }
                else {
                  const $ail_m2: { $ail_kind: "ok"; n: bigint; tail: string; value: string } | { $ail_kind: "html.invalid_url"; value: string } = html__url__authority(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), ".");
                  switch ($ail_m2.$ail_kind) {
                  case "ok": {
                    const r = $ail_m2;
                    return { $ail_kind: "ok", value: r.value, tail: r.tail, n: r.n };
                  }
                  case "html.invalid_url": {
                    const e = $ail_m2;
                    return { $ail_kind: "html.invalid_url", value: e.value };
                  }
                  default: {
                    throw new Error("unreachable");
                  }
                  }
                }
              }
            }
            else {
              if (($ailStrAt(s, 0n) === 45n)) {
                if ((prev === "")) {
                  return { $ail_kind: "html.invalid_url", value: orig };
                }
                else {
                  const $ail_m3: { $ail_kind: "ok"; n: bigint; tail: string; value: string } | { $ail_kind: "html.invalid_url"; value: string } = html__url__authority(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), "-");
                  switch ($ail_m3.$ail_kind) {
                  case "ok": {
                    const r = $ail_m3;
                    return { $ail_kind: "ok", value: r.value, tail: r.tail, n: r.n };
                  }
                  case "html.invalid_url": {
                    const e = $ail_m3;
                    return { $ail_kind: "html.invalid_url", value: e.value };
                  }
                  default: {
                    throw new Error("unreachable");
                  }
                  }
                }
              }
              else {
                if (($ailStrAt(s, 0n) <= 47n)) {
                  return { $ail_kind: "html.invalid_url", value: orig };
                }
                else {
                  if (($ailStrAt(s, 0n) <= 57n)) {
                    const $ail_m4: { $ail_kind: "ok"; n: bigint; tail: string; value: string } | { $ail_kind: "html.invalid_url"; value: string } = html__url__authority(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), $ailStrSlice(s, 0n, 1n));
                    switch ($ail_m4.$ail_kind) {
                    case "ok": {
                      const r = $ail_m4;
                      return { $ail_kind: "ok", value: r.value, tail: r.tail, n: r.n };
                    }
                    case "html.invalid_url": {
                      const e = $ail_m4;
                      return { $ail_kind: "html.invalid_url", value: e.value };
                    }
                    default: {
                      throw new Error("unreachable");
                    }
                    }
                  }
                  else {
                    if (($ailStrAt(s, 0n) <= 64n)) {
                      return { $ail_kind: "html.invalid_url", value: orig };
                    }
                    else {
                      if (($ailStrAt(s, 0n) <= 90n)) {
                        const $ail_m5: { $ail_kind: "ok"; n: bigint; tail: string; value: string } | { $ail_kind: "html.invalid_url"; value: string } = html__url__authority(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), $ailStrSlice(s, 0n, 1n));
                        switch ($ail_m5.$ail_kind) {
                        case "ok": {
                          const r = $ail_m5;
                          return { $ail_kind: "ok", value: r.value, tail: r.tail, n: r.n };
                        }
                        case "html.invalid_url": {
                          const e = $ail_m5;
                          return { $ail_kind: "html.invalid_url", value: e.value };
                        }
                        default: {
                          throw new Error("unreachable");
                        }
                        }
                      }
                      else {
                        if (($ailStrAt(s, 0n) <= 96n)) {
                          return { $ail_kind: "html.invalid_url", value: orig };
                        }
                        else {
                          if (($ailStrAt(s, 0n) <= 122n)) {
                            const $ail_m6: { $ail_kind: "ok"; n: bigint; tail: string; value: string } | { $ail_kind: "html.invalid_url"; value: string } = html__url__authority(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), $ailStrSlice(s, 0n, 1n));
                            switch ($ail_m6.$ail_kind) {
                            case "ok": {
                              const r = $ail_m6;
                              return { $ail_kind: "ok", value: r.value, tail: r.tail, n: r.n };
                            }
                            case "html.invalid_url": {
                              const e = $ail_m6;
                              return { $ail_kind: "html.invalid_url", value: e.value };
                            }
                            default: {
                              throw new Error("unreachable");
                            }
                            }
                          }
                          else {
                            return { $ail_kind: "html.invalid_url", value: orig };
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
export function html__url__tail(orig: string, s: string, n: bigint): { $ail_kind: "ok"; value: string } | { $ail_kind: "html.invalid_url"; value: string } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: orig };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; ws: boolean } = html__attribute__id_ws($ailStrAt(s, 0n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const c = $ail_m1;
      if (c.ws) {
        return { $ail_kind: "html.invalid_url", value: orig };
      }
      else {
        if (($ailStrAt(s, 0n) === 92n)) {
          return { $ail_kind: "html.invalid_url", value: orig };
        }
        else {
          const $ail_m2: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.invalid_url"; value: string } = html__url__tail(orig, $ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n));
          switch ($ail_m2.$ail_kind) {
          case "ok": {
            const r = $ail_m2;
            return { $ail_kind: "ok", value: r.value };
          }
          case "html.invalid_url": {
            const e = $ail_m2;
            return { $ail_kind: "html.invalid_url", value: e.value };
          }
          default: {
            throw new Error("unreachable");
          }
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
export function html__url__absolute(url: string): { $ail_kind: "ok"; len: bigint; value: string } | { $ail_kind: "html.invalid_url"; value: string } {
  if (((BigInt([...url].length)) <= 8n)) {
    return { $ail_kind: "html.invalid_url", value: url };
  }
  else {
    if (($ailStrSlice(url, 0n, 8n) === "https://")) {
      const $ail_m1: { $ail_kind: "ok"; n: bigint; tail: string; value: string } | { $ail_kind: "html.invalid_url"; value: string } = html__url__authority(url, $ailStrSlice(url, 8n, (BigInt([...url].length))), ((BigInt([...url].length)) - 8n), "");
      switch ($ail_m1.$ail_kind) {
      case "ok": {
        const a = $ail_m1;
        if ((a.tail === "")) {
          return { $ail_kind: "ok", value: a.value, len: (BigInt([...url].length)) };
        }
        else {
          const $ail_m2: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.invalid_url"; value: string } = html__url__tail(url, a.tail, a.n);
          switch ($ail_m2.$ail_kind) {
          case "ok": {
            const t = $ail_m2;
            return { $ail_kind: "ok", value: t.value, len: (BigInt([...url].length)) };
          }
          case "html.invalid_url": {
            const e = $ail_m2;
            return { $ail_kind: "html.invalid_url", value: e.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
      }
      case "html.invalid_url": {
        const e = $ail_m1;
        return { $ail_kind: "html.invalid_url", value: e.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $ail_kind: "html.invalid_url", value: url };
    }
  }
}
export function html__url__check(raw: string): { $ail_kind: "ok"; len: bigint; value: string } | { $ail_kind: "html.invalid_url"; value: string } | { $ail_kind: "html.disallowed_scheme"; value: string } {
  if (((BigInt([...raw].length)) <= 0n)) {
    return { $ail_kind: "html.invalid_url", value: raw };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; has: boolean } = html__url__scheme_token(raw, raw, (BigInt([...raw].length)), (BigInt([...raw].length)));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const t = $ail_m1;
      const $ail_m2: { $ail_kind: "ok"; len: bigint; value: string } | { $ail_kind: "html.invalid_url"; value: string } | { $ail_kind: "html.disallowed_scheme"; value: string } = html__url__admit(raw, t.has);
      switch ($ail_m2.$ail_kind) {
      case "ok": {
        const u = $ail_m2;
        return { $ail_kind: "ok", value: u.value, len: u.len };
      }
      case "html.invalid_url": {
        const e = $ail_m2;
        return { $ail_kind: "html.invalid_url", value: e.value };
      }
      case "html.disallowed_scheme": {
        const e = $ail_m2;
        return { $ail_kind: "html.disallowed_scheme", value: e.value };
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
export function html__url__admit(raw: string, has: boolean): { $ail_kind: "ok"; len: bigint; value: string } | { $ail_kind: "html.invalid_url"; value: string } | { $ail_kind: "html.disallowed_scheme"; value: string } {
  if (has) {
    if (((BigInt([...raw].length)) <= 5n)) {
      return { $ail_kind: "html.disallowed_scheme", value: raw };
    }
    else {
      if (($ailStrSlice(raw, 0n, 6n) === "https:")) {
        const $ail_m1: { $ail_kind: "ok"; len: bigint; value: string } | { $ail_kind: "html.invalid_url"; value: string } = html__url__absolute(raw);
        switch ($ail_m1.$ail_kind) {
        case "ok": {
          const u = $ail_m1;
          return { $ail_kind: "ok", value: u.value, len: u.len };
        }
        case "html.invalid_url": {
          const e = $ail_m1;
          return { $ail_kind: "html.invalid_url", value: e.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        return { $ail_kind: "html.disallowed_scheme", value: raw };
      }
    }
  }
  else {
    return { $ail_kind: "html.invalid_url", value: raw };
  }
}
export function html__attribute__href(url: string): { $ail_kind: "ok"; attribute: string } | { $ail_kind: "html.invalid_url"; value: string } | { $ail_kind: "html.disallowed_scheme"; value: string } | { $ail_kind: "html.nul_byte"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; len: bigint; value: string } | { $ail_kind: "html.invalid_url"; value: string } | { $ail_kind: "html.disallowed_scheme"; value: string } = html__url__check(url);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const u = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__attribute__value_from(u.value, u.value, "", u.len);
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const v = $ail_m2;
      return { $ail_kind: "ok", attribute: (("href='" + v.value) + "'") };
    }
    case "html.nul_byte": {
      const e = $ail_m2;
      return { $ail_kind: "html.nul_byte", value: u.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "html.invalid_url": {
    const e = $ail_m1;
    return { $ail_kind: "html.invalid_url", value: e.value };
  }
  case "html.disallowed_scheme": {
    const e = $ail_m1;
    return { $ail_kind: "html.disallowed_scheme", value: e.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function html__attribute__src(url: string): { $ail_kind: "ok"; attribute: string } | { $ail_kind: "html.invalid_url"; value: string } | { $ail_kind: "html.disallowed_scheme"; value: string } | { $ail_kind: "html.nul_byte"; value: string } {
  const $ail_m1: { $ail_kind: "ok"; len: bigint; value: string } | { $ail_kind: "html.invalid_url"; value: string } | { $ail_kind: "html.disallowed_scheme"; value: string } = html__url__check(url);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const u = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; value: string } | { $ail_kind: "html.nul_byte"; value: string } = html__attribute__value_from(u.value, u.value, "", u.len);
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const v = $ail_m2;
      return { $ail_kind: "ok", attribute: (("src='" + v.value) + "'") };
    }
    case "html.nul_byte": {
      const e = $ail_m2;
      return { $ail_kind: "html.nul_byte", value: u.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "html.invalid_url": {
    const e = $ail_m1;
    return { $ail_kind: "html.invalid_url", value: e.value };
  }
  case "html.disallowed_scheme": {
    const e = $ail_m1;
    return { $ail_kind: "html.disallowed_scheme", value: e.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function html__fragment__empty(): { $ail_kind: "ok"; safe: string } {
  return { $ail_kind: "ok", safe: "" };
}
