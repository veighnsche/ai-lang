// GENERATED from json.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { std__convert__dec_to_str, std__convert__int_to_str, std__convert__str_to_dec, std__convert__str_to_int, type Bool__Value, type Dec__Value, type Int__Value, type Str__Value, type ScalarsResult } from "./scalars";
export type JsonResult = { $can_kind: "ok"; frames: Json__Frame[] } | { $can_kind: "ok"; frames: Json__PFrame[] } | { $can_kind: "ok"; next: bigint; val: Json__Value } | { $can_kind: "ok"; tag: Json__HeadTag } | { $can_kind: "ok"; value: Json__Value } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: boolean } | { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string } | { $can_kind: "json.numeric_out_of_range"; text: string } | { $can_kind: "json.invalid_syntax"; detail: string } | { $can_kind: "json.duplicate_key"; key: string };
export type Int__Value = { value: bigint };
export type Dec__Value = { value: string };
export type Str__Value = { value: string };
export type Bool__Value = { value: boolean };
export type Json__Text = { value: string };
export type Json__Field = { name: string; value: Json__Value };
export type Json__Doc = { value: Json__Value };
export type Json__Frame = { tag: Json__Tag; text: string; val: Json__Value; items: Json__Doc[]; ipos: bigint; fields: Json__Field[]; fpos: bigint };
export type Json__Frames = { frames: Json__Frame[] };
export type Json__IntSchema = { encode: (input: bigint) => { $can_kind: "ok"; value: Json__Value }; decode: (input: Json__Value) => { $can_kind: "ok"; value: bigint } | { $can_kind: "json.numeric_out_of_range"; text: string } | { $can_kind: "json.schema_mismatch"; detail: string } };
export type Json__StrSchema = { encode: (input: string) => { $can_kind: "ok"; value: Json__Value }; decode: (input: Json__Value) => { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string } };
export type Json__BoolSchema = { encode: (input: boolean) => { $can_kind: "ok"; value: Json__Value }; decode: (input: Json__Value) => { $can_kind: "ok"; value: boolean } | { $can_kind: "json.schema_mismatch"; detail: string } };
export type Json__DecSchema = { encode: (input: string) => { $can_kind: "ok"; value: Json__Value }; decode: (input: Json__Value) => { $can_kind: "ok"; value: string } | { $can_kind: "json.numeric_out_of_range"; text: string } | { $can_kind: "json.schema_mismatch"; detail: string } };
export type Json__PFrame = { tag: Json__PTag; items: Json__Doc[]; fields: Json__Field[]; keys: string[]; key: string; acc: string; val: Json__Value; ofields: Json__Field[]; okeys: string[] };
export type Json__PFrames = { frames: Json__PFrame[] };
export type Json__Head = { tag: Json__HeadTag };
export type Json__LitOut = { val: Json__Value; next: bigint };
export type Json__Value = { $can_kind: "Json__Null" } | { $can_kind: "Json__Bool"; value: boolean } | { $can_kind: "Json__Num"; text: string } | { $can_kind: "Json__Str"; value: string } | { $can_kind: "Json__Arr"; items: Json__Doc[] } | { $can_kind: "Json__Obj"; fields: Json__Field[] };
export type Json__Tag = { $can_kind: "Json__TagEmit" } | { $can_kind: "Json__TagValue" } | { $can_kind: "Json__TagItems" } | { $can_kind: "Json__TagFields" };
export type Json__PTag = { $can_kind: "Json__PArrFirst" } | { $can_kind: "Json__PArrVal" } | { $can_kind: "Json__PArrNext" } | { $can_kind: "Json__PObjFirst" } | { $can_kind: "Json__PObjKey" } | { $can_kind: "Json__PObjColon" } | { $can_kind: "Json__PObjKeyVal" } | { $can_kind: "Json__PObjNext" } | { $can_kind: "Json__PStrKey" } | { $can_kind: "Json__PStrVal" } | { $can_kind: "Json__PNumAcc" } | { $can_kind: "Json__PTail" };
export type Json__HeadTag = { $can_kind: "Json__HeadBrace" } | { $can_kind: "Json__HeadBracket" } | { $can_kind: "Json__HeadStr" } | { $can_kind: "Json__HeadLit" } | { $can_kind: "Json__HeadNum" } | { $can_kind: "Json__HeadBad" } | { $can_kind: "Json__HeadEof" };
// Strict boolean runtime (slice 5): eager helpers, never bare &&.
function $canBoolAnd(a: boolean, b: boolean): boolean {
  return a && b;
}
// Byte-order string comparison: UTF-8 bytes, matching Go.
function $canStrCmp(a: string, b: string): number {
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
function $canStrAt(s: string, i: bigint): bigint {
  const cps = [...s];
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("str index out of range");
  const k = Number(i);
  if (k >= cps.length) throw new Error("str index out of range");
  const cp = cps[k].codePointAt(0);
  if (cp === undefined) throw new Error("str index out of range");
  return BigInt(cp);
}
function $canStrSlice(s: string, a: bigint, b: bigint): string {
  const cps = [...s];
  const toIdx = (x: bigint): number => { if (x < 0n || x > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("str slice out of range"); return Number(x); };
  const lo = toIdx(a), hi = toIdx(b);
  if (lo > hi || hi > cps.length) throw new Error("str slice out of range");
  return cps.slice(lo, hi).join("");
}
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
export function std__json__pop(frames: Json__Frame[]): { $can_kind: "ok"; frames: Json__Frame[] } {
  const $can_m1: { $can_kind: "ok"; frames: Json__Frame[] } = std__json__pop_from(frames, 0n, ((BigInt([...frames].length)) - 1n), []);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", frames: r.frames };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__pop_from(frames: Json__Frame[], pos: bigint, n: bigint, acc: Json__Frame[]): { $can_kind: "ok"; frames: Json__Frame[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", frames: acc };
  }
  else {
    if ((pos < (BigInt([...frames].length)))) {
      const $can_m1: { $can_kind: "ok"; frames: Json__Frame[] } = std__json__pop_from(frames, (pos + 1n), (n - 1n), [...acc, $canSeqAt(frames, pos)]);
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        return { $can_kind: "ok", frames: r.frames };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $can_kind: "ok", frames: acc };
    }
  }
}
export function std__json__render_value(v: Json__Value): { $can_kind: "ok"; value: string } {
  const $can_m1 = v;
  switch ($can_m1.$can_kind) {
  case "Json__Null": {
    const _ = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: string } = std__json__render_step([{ tag: { $can_kind: "Json__TagValue" }, text: "", val: v, items: [], ipos: 0n, fields: [], fpos: 0n }], "", 1000000n);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const r = $can_m2;
      return { $can_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "Json__Bool": {
    const _ = $can_m1;
    const $can_m3: { $can_kind: "ok"; value: string } = std__json__render_step([{ tag: { $can_kind: "Json__TagValue" }, text: "", val: v, items: [], ipos: 0n, fields: [], fpos: 0n }], "", 1000000n);
    switch ($can_m3.$can_kind) {
    case "ok": {
      const r = $can_m3;
      return { $can_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "Json__Num": {
    const _ = $can_m1;
    const $can_m4: { $can_kind: "ok"; value: string } = std__json__render_step([{ tag: { $can_kind: "Json__TagValue" }, text: "", val: v, items: [], ipos: 0n, fields: [], fpos: 0n }], "", 1000000n);
    switch ($can_m4.$can_kind) {
    case "ok": {
      const r = $can_m4;
      return { $can_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "Json__Str": {
    const _ = $can_m1;
    const $can_m5: { $can_kind: "ok"; value: string } = std__json__render_step([{ tag: { $can_kind: "Json__TagValue" }, text: "", val: v, items: [], ipos: 0n, fields: [], fpos: 0n }], "", 1000000n);
    switch ($can_m5.$can_kind) {
    case "ok": {
      const r = $can_m5;
      return { $can_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "Json__Arr": {
    const _ = $can_m1;
    const $can_m6: { $can_kind: "ok"; value: string } = std__json__render_step([{ tag: { $can_kind: "Json__TagValue" }, text: "", val: v, items: [], ipos: 0n, fields: [], fpos: 0n }], "", 1000000n);
    switch ($can_m6.$can_kind) {
    case "ok": {
      const r = $can_m6;
      return { $can_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "Json__Obj": {
    const _ = $can_m1;
    const $can_m7: { $can_kind: "ok"; value: string } = std__json__render_step([{ tag: { $can_kind: "Json__TagValue" }, text: "", val: v, items: [], ipos: 0n, fields: [], fpos: 0n }], "", 1000000n);
    switch ($can_m7.$can_kind) {
    case "ok": {
      const r = $can_m7;
      return { $can_kind: "ok", value: r.value };
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
export function std__json__render_step(frames: Json__Frame[], acc: string, fuel: bigint): { $can_kind: "ok"; value: string } {
  if ((fuel <= 0n)) {
    return { $can_kind: "ok", value: acc };
  }
  else {
    if (((BigInt([...frames].length)) === 0n)) {
      return { $can_kind: "ok", value: acc };
    }
    else {
      const $can_m1: { $can_kind: "ok"; frames: Json__Frame[] } = std__json__pop(frames);
      switch ($can_m1.$can_kind) {
      case "ok": {
        const pf = $can_m1;
        const $can_m2 = $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).tag;
        switch ($can_m2.$can_kind) {
        case "Json__TagEmit": {
          const _ = $can_m2;
          const $can_m3: { $can_kind: "ok"; value: string } = std__json__render_step(pf.frames, (acc + $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).text), (fuel - 1n));
          switch ($can_m3.$can_kind) {
          case "ok": {
            const r = $can_m3;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        case "Json__TagValue": {
          const _ = $can_m2;
          const $can_m4 = $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).val;
          switch ($can_m4.$can_kind) {
          case "Json__Null": {
            const _ = $can_m4;
            const $can_m5: { $can_kind: "ok"; value: string } = std__json__render_step(pf.frames, (acc + "null"), (fuel - 1n));
            switch ($can_m5.$can_kind) {
            case "ok": {
              const r = $can_m5;
              return { $can_kind: "ok", value: r.value };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          case "Json__Bool": {
            const b = $can_m4;
            if (b.value) {
              const $can_m6: { $can_kind: "ok"; value: string } = std__json__render_step(pf.frames, (acc + "true"), (fuel - 1n));
              switch ($can_m6.$can_kind) {
              case "ok": {
                const r = $can_m6;
                return { $can_kind: "ok", value: r.value };
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            else {
              const $can_m7: { $can_kind: "ok"; value: string } = std__json__render_step(pf.frames, (acc + "false"), (fuel - 1n));
              switch ($can_m7.$can_kind) {
              case "ok": {
                const r = $can_m7;
                return { $can_kind: "ok", value: r.value };
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
          }
          case "Json__Num": {
            const n = $can_m4;
            const $can_m8: { $can_kind: "ok"; value: string } = std__json__render_step(pf.frames, (acc + n.text), (fuel - 1n));
            switch ($can_m8.$can_kind) {
            case "ok": {
              const r = $can_m8;
              return { $can_kind: "ok", value: r.value };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          case "Json__Str": {
            const s = $can_m4;
            const $can_m9: { $can_kind: "ok"; value: string } = std__json__escape(s.value);
            switch ($can_m9.$can_kind) {
            case "ok": {
              const e = $can_m9;
              const $can_m10: { $can_kind: "ok"; value: string } = std__json__render_step(pf.frames, (((acc + "\"") + e.value) + "\""), (fuel - 1n));
              switch ($can_m10.$can_kind) {
              case "ok": {
                const r = $can_m10;
                return { $can_kind: "ok", value: r.value };
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
          case "Json__Arr": {
            const a = $can_m4;
            const $can_m11: { $can_kind: "ok"; value: string } = std__json__render_step([...[...[...pf.frames, { tag: { $can_kind: "Json__TagEmit" }, text: "]", val: { $can_kind: "Json__Null" }, items: [], ipos: 0n, fields: [], fpos: 0n }], { tag: { $can_kind: "Json__TagItems" }, text: "", val: { $can_kind: "Json__Null" }, items: a.items, ipos: 0n, fields: [], fpos: 0n }], { tag: { $can_kind: "Json__TagEmit" }, text: "[", val: { $can_kind: "Json__Null" }, items: [], ipos: 0n, fields: [], fpos: 0n }], acc, (fuel - 1n));
            switch ($can_m11.$can_kind) {
            case "ok": {
              const r = $can_m11;
              return { $can_kind: "ok", value: r.value };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          case "Json__Obj": {
            const o = $can_m4;
            const $can_m12: { $can_kind: "ok"; value: string } = std__json__render_step([...[...[...pf.frames, { tag: { $can_kind: "Json__TagEmit" }, text: "}", val: { $can_kind: "Json__Null" }, items: [], ipos: 0n, fields: [], fpos: 0n }], { tag: { $can_kind: "Json__TagFields" }, text: "", val: { $can_kind: "Json__Null" }, items: [], ipos: 0n, fields: o.fields, fpos: 0n }], { tag: { $can_kind: "Json__TagEmit" }, text: "{", val: { $can_kind: "Json__Null" }, items: [], ipos: 0n, fields: [], fpos: 0n }], acc, (fuel - 1n));
            switch ($can_m12.$can_kind) {
            case "ok": {
              const r = $can_m12;
              return { $can_kind: "ok", value: r.value };
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
        case "Json__TagItems": {
          const _ = $can_m2;
          if (($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).ipos < (BigInt([...$canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).items].length)))) {
            if (($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).ipos === 0n)) {
              const $can_m13: { $can_kind: "ok"; value: string } = std__json__render_step([...[...pf.frames, { tag: { $can_kind: "Json__TagItems" }, text: "", val: { $can_kind: "Json__Null" }, items: $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).items, ipos: ($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).ipos + 1n), fields: [], fpos: 0n }], { tag: { $can_kind: "Json__TagValue" }, text: "", val: $canSeqAt($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).items, $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).ipos).value, items: [], ipos: 0n, fields: [], fpos: 0n }], acc, (fuel - 1n));
              switch ($can_m13.$can_kind) {
              case "ok": {
                const r = $can_m13;
                return { $can_kind: "ok", value: r.value };
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            else {
              const $can_m14: { $can_kind: "ok"; value: string } = std__json__render_step([...[...[...pf.frames, { tag: { $can_kind: "Json__TagItems" }, text: "", val: { $can_kind: "Json__Null" }, items: $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).items, ipos: ($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).ipos + 1n), fields: [], fpos: 0n }], { tag: { $can_kind: "Json__TagValue" }, text: "", val: $canSeqAt($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).items, $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).ipos).value, items: [], ipos: 0n, fields: [], fpos: 0n }], { tag: { $can_kind: "Json__TagEmit" }, text: ", ", val: { $can_kind: "Json__Null" }, items: [], ipos: 0n, fields: [], fpos: 0n }], acc, (fuel - 1n));
              switch ($can_m14.$can_kind) {
              case "ok": {
                const r = $can_m14;
                return { $can_kind: "ok", value: r.value };
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
          }
          else {
            const $can_m15: { $can_kind: "ok"; value: string } = std__json__render_step(pf.frames, acc, (fuel - 1n));
            switch ($can_m15.$can_kind) {
            case "ok": {
              const r = $can_m15;
              return { $can_kind: "ok", value: r.value };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
        }
        case "Json__TagFields": {
          const _ = $can_m2;
          if (($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fpos < (BigInt([...$canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fields].length)))) {
            const $can_m16: { $can_kind: "ok"; value: string } = std__json__escape($canSeqAt($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fields, $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fpos).name);
            switch ($can_m16.$can_kind) {
            case "ok": {
              const k = $can_m16;
              if (($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fpos === 0n)) {
                const $can_m17: { $can_kind: "ok"; value: string } = std__json__render_step([...[...[...pf.frames, { tag: { $can_kind: "Json__TagFields" }, text: "", val: { $can_kind: "Json__Null" }, items: [], ipos: 0n, fields: $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fields, fpos: ($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fpos + 1n) }], { tag: { $can_kind: "Json__TagValue" }, text: "", val: $canSeqAt($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fields, $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fpos).value, items: [], ipos: 0n, fields: [], fpos: 0n }], { tag: { $can_kind: "Json__TagEmit" }, text: (("\"" + k.value) + "\": "), val: { $can_kind: "Json__Null" }, items: [], ipos: 0n, fields: [], fpos: 0n }], acc, (fuel - 1n));
                switch ($can_m17.$can_kind) {
                case "ok": {
                  const r = $can_m17;
                  return { $can_kind: "ok", value: r.value };
                }
                default: {
                  throw new Error("unreachable");
                }
                }
              }
              else {
                const $can_m18: { $can_kind: "ok"; value: string } = std__json__render_step([...[...[...pf.frames, { tag: { $can_kind: "Json__TagFields" }, text: "", val: { $can_kind: "Json__Null" }, items: [], ipos: 0n, fields: $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fields, fpos: ($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fpos + 1n) }], { tag: { $can_kind: "Json__TagValue" }, text: "", val: $canSeqAt($canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fields, $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fpos).value, items: [], ipos: 0n, fields: [], fpos: 0n }], { tag: { $can_kind: "Json__TagEmit" }, text: (((", " + "\"") + k.value) + "\": "), val: { $can_kind: "Json__Null" }, items: [], ipos: 0n, fields: [], fpos: 0n }], acc, (fuel - 1n));
                switch ($can_m18.$can_kind) {
                case "ok": {
                  const r = $can_m18;
                  return { $can_kind: "ok", value: r.value };
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
            const $can_m19: { $can_kind: "ok"; value: string } = std__json__render_step(pf.frames, acc, (fuel - 1n));
            switch ($can_m19.$can_kind) {
            case "ok": {
              const r = $can_m19;
              return { $can_kind: "ok", value: r.value };
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
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
}
export function std__json__escape(s: string): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } = std__json__escape_from(s, 0n, (BigInt([...s].length)), "");
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
export function std__json__escape_from(value: string, pos: bigint, n: bigint, acc: string): { $can_kind: "ok"; value: string } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: acc };
  }
  else {
    if ((pos < (BigInt([...value].length)))) {
      if ($canStrAt(value, pos) === 34n) {
        const $can_m1: { $can_kind: "ok"; value: string } = std__json__escape_from(value, (pos + 1n), (n - 1n), (acc + "\\\""));
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
      else if ($canStrAt(value, pos) === 92n) {
        const $can_m2: { $can_kind: "ok"; value: string } = std__json__escape_from(value, (pos + 1n), (n - 1n), (acc + "\\\\"));
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else if ($canStrAt(value, pos) === 9n) {
        const $can_m3: { $can_kind: "ok"; value: string } = std__json__escape_from(value, (pos + 1n), (n - 1n), (acc + "\\t"));
        switch ($can_m3.$can_kind) {
        case "ok": {
          const r = $can_m3;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else if ($canStrAt(value, pos) === 10n) {
        const $can_m4: { $can_kind: "ok"; value: string } = std__json__escape_from(value, (pos + 1n), (n - 1n), (acc + "\\n"));
        switch ($can_m4.$can_kind) {
        case "ok": {
          const r = $can_m4;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else if ($canStrAt(value, pos) === 13n) {
        const $can_m5: { $can_kind: "ok"; value: string } = std__json__escape_from(value, (pos + 1n), (n - 1n), (acc + "\\r"));
        switch ($can_m5.$can_kind) {
        case "ok": {
          const r = $can_m5;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      if (($canStrAt(value, pos) < 32n)) {
        const $can_m6: { $can_kind: "ok"; value: string } = std__json__escape_hex($canStrAt(value, pos));
        switch ($can_m6.$can_kind) {
        case "ok": {
          const h = $can_m6;
          const $can_m7: { $can_kind: "ok"; value: string } = std__json__escape_from(value, (pos + 1n), (n - 1n), ((acc + "\\u00") + h.value));
          switch ($can_m7.$can_kind) {
          case "ok": {
            const r = $can_m7;
            return { $can_kind: "ok", value: r.value };
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
        const $can_m8: { $can_kind: "ok"; value: string } = std__json__escape_from(value, (pos + 1n), (n - 1n), (acc + $canStrSlice(value, pos, (pos + 1n))));
        switch ($can_m8.$can_kind) {
        case "ok": {
          const r = $can_m8;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    else {
      return { $can_kind: "ok", value: acc };
    }
  }
}
export function std__json__escape_hex(code: bigint): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } = std__json__hexdigit($canDivMod(code, 16n)[0]);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const hi = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: string } = std__json__hexdigit($canDivMod(code, 16n)[1]);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const lo = $can_m2;
      return { $can_kind: "ok", value: (hi.value + lo.value) };
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
export function std__json__hexdigit(d: bigint): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: $canSeqAt(["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"], d) };
}
export function std__json__int__encode(value: bigint): { $can_kind: "ok"; value: Json__Value } {
  const $can_m1: { $can_kind: "ok"; value: string } = std__convert__int_to_str(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: { $can_kind: "Json__Num", text: r.value } };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__int__decode(value: Json__Value): { $can_kind: "ok"; value: bigint } | { $can_kind: "json.schema_mismatch"; detail: string } | { $can_kind: "json.numeric_out_of_range"; text: string } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Json__Num": {
    const t = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } | { $can_kind: "convert.invalid_integer"; value: string } = std__convert__str_to_int(t.text);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const r = $can_m2;
      return { $can_kind: "ok", value: r.value };
    }
    case "convert.invalid_integer": {
      const _ = $can_m2;
      return { $can_kind: "json.numeric_out_of_range", text: t.text };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "Json__Null": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want int, got null" };
  }
  case "Json__Bool": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want int, got bool" };
  }
  case "Json__Str": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want int, got str" };
  }
  case "Json__Arr": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want int, got arr" };
  }
  case "Json__Obj": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want int, got obj" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__str__encode(value: string): { $can_kind: "ok"; value: Json__Value } {
  return { $can_kind: "ok", value: { $can_kind: "Json__Str", value: value } };
}
export function std__json__str__decode(value: Json__Value): { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Json__Str": {
    const s = $can_m1;
    return { $can_kind: "ok", value: s.value };
  }
  case "Json__Null": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want str, got null" };
  }
  case "Json__Bool": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want str, got bool" };
  }
  case "Json__Num": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want str, got num" };
  }
  case "Json__Arr": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want str, got arr" };
  }
  case "Json__Obj": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want str, got obj" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__bool__encode(value: boolean): { $can_kind: "ok"; value: Json__Value } {
  return { $can_kind: "ok", value: { $can_kind: "Json__Bool", value: value } };
}
export function std__json__bool__decode(value: Json__Value): { $can_kind: "ok"; value: boolean } | { $can_kind: "json.schema_mismatch"; detail: string } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Json__Bool": {
    const b = $can_m1;
    return { $can_kind: "ok", value: b.value };
  }
  case "Json__Null": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want bool, got null" };
  }
  case "Json__Num": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want bool, got num" };
  }
  case "Json__Str": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want bool, got str" };
  }
  case "Json__Arr": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want bool, got arr" };
  }
  case "Json__Obj": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want bool, got obj" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__dec__encode(value: string): { $can_kind: "ok"; value: Json__Value } {
  const $can_m1: { $can_kind: "ok"; value: string } = std__convert__dec_to_str(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: { $can_kind: "Json__Num", text: r.value } };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__dec__decode(value: Json__Value): { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string } | { $can_kind: "json.numeric_out_of_range"; text: string } {
  const $can_m1 = value;
  switch ($can_m1.$can_kind) {
  case "Json__Num": {
    const t = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: string } | { $can_kind: "convert.invalid_decimal"; value: string } = std__convert__str_to_dec(t.text);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const r = $can_m2;
      return { $can_kind: "ok", value: r.value };
    }
    case "convert.invalid_decimal": {
      const _ = $can_m2;
      return { $can_kind: "json.numeric_out_of_range", text: t.text };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  case "Json__Null": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want dec, got null" };
  }
  case "Json__Bool": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want dec, got bool" };
  }
  case "Json__Str": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want dec, got str" };
  }
  case "Json__Arr": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want dec, got arr" };
  }
  case "Json__Obj": {
    const _ = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: "want dec, got obj" };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__int__encode_apply(value: bigint, enc: (input: bigint) => { $can_kind: "ok"; value: Json__Value }): { $can_kind: "ok"; value: Json__Value } {
  const $can_m1: { $can_kind: "ok"; value: Json__Value } = enc(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const d = $can_m1;
    return { $can_kind: "ok", value: d.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__int__decode_apply(value: Json__Value, dec: (input: Json__Value) => { $can_kind: "ok"; value: bigint } | { $can_kind: "json.numeric_out_of_range"; text: string } | { $can_kind: "json.schema_mismatch"; detail: string }): { $can_kind: "ok"; value: bigint } | { $can_kind: "json.schema_mismatch"; detail: string } | { $can_kind: "json.numeric_out_of_range"; text: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "json.numeric_out_of_range"; text: string } | { $can_kind: "json.schema_mismatch"; detail: string } = dec(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "json.schema_mismatch": {
    const e = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: e.detail };
  }
  case "json.numeric_out_of_range": {
    const e = $can_m1;
    return { $can_kind: "json.numeric_out_of_range", text: e.text };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__str__encode_apply(value: string, enc: (input: string) => { $can_kind: "ok"; value: Json__Value }): { $can_kind: "ok"; value: Json__Value } {
  const $can_m1: { $can_kind: "ok"; value: Json__Value } = enc(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const d = $can_m1;
    return { $can_kind: "ok", value: d.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__str__decode_apply(value: Json__Value, dec: (input: Json__Value) => { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string }): { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string } {
  const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string } = dec(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "json.schema_mismatch": {
    const e = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: e.detail };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__bool__encode_apply(value: boolean, enc: (input: boolean) => { $can_kind: "ok"; value: Json__Value }): { $can_kind: "ok"; value: Json__Value } {
  const $can_m1: { $can_kind: "ok"; value: Json__Value } = enc(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const d = $can_m1;
    return { $can_kind: "ok", value: d.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__bool__decode_apply(value: Json__Value, dec: (input: Json__Value) => { $can_kind: "ok"; value: boolean } | { $can_kind: "json.schema_mismatch"; detail: string }): { $can_kind: "ok"; value: boolean } | { $can_kind: "json.schema_mismatch"; detail: string } {
  const $can_m1: { $can_kind: "ok"; value: boolean } | { $can_kind: "json.schema_mismatch"; detail: string } = dec(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "json.schema_mismatch": {
    const e = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: e.detail };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__dec__encode_apply(value: string, enc: (input: string) => { $can_kind: "ok"; value: Json__Value }): { $can_kind: "ok"; value: Json__Value } {
  const $can_m1: { $can_kind: "ok"; value: Json__Value } = enc(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const d = $can_m1;
    return { $can_kind: "ok", value: d.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__dec__decode_apply(value: Json__Value, dec: (input: Json__Value) => { $can_kind: "ok"; value: string } | { $can_kind: "json.numeric_out_of_range"; text: string } | { $can_kind: "json.schema_mismatch"; detail: string }): { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string } | { $can_kind: "json.numeric_out_of_range"; text: string } {
  const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "json.numeric_out_of_range"; text: string } | { $can_kind: "json.schema_mismatch"; detail: string } = dec(value);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "json.schema_mismatch": {
    const e = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: e.detail };
  }
  case "json.numeric_out_of_range": {
    const e = $can_m1;
    return { $can_kind: "json.numeric_out_of_range", text: e.text };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__int__encode_value(value: bigint, schema: Json__IntSchema): { $can_kind: "ok"; value: Json__Value } {
  const $can_m1: { $can_kind: "ok"; value: Json__Value } = std__json__int__encode_apply(value, schema.encode);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const d = $can_m1;
    return { $can_kind: "ok", value: d.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__int__decode_value(value: Json__Value, schema: Json__IntSchema): { $can_kind: "ok"; value: bigint } | { $can_kind: "json.schema_mismatch"; detail: string } | { $can_kind: "json.numeric_out_of_range"; text: string } {
  const $can_m1: { $can_kind: "ok"; value: bigint } | { $can_kind: "json.schema_mismatch"; detail: string } | { $can_kind: "json.numeric_out_of_range"; text: string } = std__json__int__decode_apply(value, schema.decode);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "json.schema_mismatch": {
    const e = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: e.detail };
  }
  case "json.numeric_out_of_range": {
    const e = $can_m1;
    return { $can_kind: "json.numeric_out_of_range", text: e.text };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__str__encode_value(value: string, schema: Json__StrSchema): { $can_kind: "ok"; value: Json__Value } {
  const $can_m1: { $can_kind: "ok"; value: Json__Value } = std__json__str__encode_apply(value, schema.encode);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const d = $can_m1;
    return { $can_kind: "ok", value: d.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__str__decode_value(value: Json__Value, schema: Json__StrSchema): { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string } {
  const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string } = std__json__str__decode_apply(value, schema.decode);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "json.schema_mismatch": {
    const e = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: e.detail };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__bool__encode_value(value: boolean, schema: Json__BoolSchema): { $can_kind: "ok"; value: Json__Value } {
  const $can_m1: { $can_kind: "ok"; value: Json__Value } = std__json__bool__encode_apply(value, schema.encode);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const d = $can_m1;
    return { $can_kind: "ok", value: d.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__bool__decode_value(value: Json__Value, schema: Json__BoolSchema): { $can_kind: "ok"; value: boolean } | { $can_kind: "json.schema_mismatch"; detail: string } {
  const $can_m1: { $can_kind: "ok"; value: boolean } | { $can_kind: "json.schema_mismatch"; detail: string } = std__json__bool__decode_apply(value, schema.decode);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "json.schema_mismatch": {
    const e = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: e.detail };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__dec__encode_value(value: string, schema: Json__DecSchema): { $can_kind: "ok"; value: Json__Value } {
  const $can_m1: { $can_kind: "ok"; value: Json__Value } = std__json__dec__encode_apply(value, schema.encode);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const d = $can_m1;
    return { $can_kind: "ok", value: d.value };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__dec__decode_value(value: Json__Value, schema: Json__DecSchema): { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string } | { $can_kind: "json.numeric_out_of_range"; text: string } {
  const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "json.schema_mismatch"; detail: string } | { $can_kind: "json.numeric_out_of_range"; text: string } = std__json__dec__decode_apply(value, schema.decode);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", value: r.value };
  }
  case "json.schema_mismatch": {
    const e = $can_m1;
    return { $can_kind: "json.schema_mismatch", detail: e.detail };
  }
  case "json.numeric_out_of_range": {
    const e = $can_m1;
    return { $can_kind: "json.numeric_out_of_range", text: e.text };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__parse_ws_skip(text: string, pos: bigint): { $can_kind: "ok"; value: bigint } {
  const $can_m1: { $can_kind: "ok"; value: bigint } = std__json__parse_ws_skip_from(text, pos, (BigInt([...text].length)), pos);
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
export function std__json__parse_ws_skip_from(text: string, pos: bigint, n: bigint, acc: bigint): { $can_kind: "ok"; value: bigint } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: acc };
  }
  else {
    if ((pos < (BigInt([...text].length)))) {
      if ($canStrAt(text, pos) === 32n) {
        const $can_m1: { $can_kind: "ok"; value: bigint } = std__json__parse_ws_skip_from(text, (pos + 1n), (n - 1n), (pos + 1n));
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
      else if ($canStrAt(text, pos) === 9n) {
        const $can_m2: { $can_kind: "ok"; value: bigint } = std__json__parse_ws_skip_from(text, (pos + 1n), (n - 1n), (pos + 1n));
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else if ($canStrAt(text, pos) === 10n) {
        const $can_m3: { $can_kind: "ok"; value: bigint } = std__json__parse_ws_skip_from(text, (pos + 1n), (n - 1n), (pos + 1n));
        switch ($can_m3.$can_kind) {
        case "ok": {
          const r = $can_m3;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else if ($canStrAt(text, pos) === 13n) {
        const $can_m4: { $can_kind: "ok"; value: bigint } = std__json__parse_ws_skip_from(text, (pos + 1n), (n - 1n), (pos + 1n));
        switch ($can_m4.$can_kind) {
        case "ok": {
          const r = $can_m4;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      return { $can_kind: "ok", value: pos };
    }
    else {
      return { $can_kind: "ok", value: acc };
    }
  }
}
export function std__json__parse_head(text: string, pos: bigint): { $can_kind: "ok"; tag: Json__HeadTag } {
  if ((pos < (BigInt([...text].length)))) {
    if ($canStrAt(text, pos) === 123n) {
      return { $can_kind: "ok", tag: { $can_kind: "Json__HeadBrace" } };
    }
    else if ($canStrAt(text, pos) === 91n) {
      return { $can_kind: "ok", tag: { $can_kind: "Json__HeadBracket" } };
    }
    else if ($canStrAt(text, pos) === 34n) {
      return { $can_kind: "ok", tag: { $can_kind: "Json__HeadStr" } };
    }
    else if ($canStrAt(text, pos) === 116n) {
      return { $can_kind: "ok", tag: { $can_kind: "Json__HeadLit" } };
    }
    else if ($canStrAt(text, pos) === 102n) {
      return { $can_kind: "ok", tag: { $can_kind: "Json__HeadLit" } };
    }
    else if ($canStrAt(text, pos) === 110n) {
      return { $can_kind: "ok", tag: { $can_kind: "Json__HeadLit" } };
    }
    else if ($canStrAt(text, pos) === 45n) {
      return { $can_kind: "ok", tag: { $can_kind: "Json__HeadNum" } };
    }
    if (($canStrAt(text, pos) >= 48n)) {
      if (($canStrAt(text, pos) <= 57n)) {
        return { $can_kind: "ok", tag: { $can_kind: "Json__HeadNum" } };
      }
      else {
        return { $can_kind: "ok", tag: { $can_kind: "Json__HeadBad" } };
      }
    }
    else {
      return { $can_kind: "ok", tag: { $can_kind: "Json__HeadBad" } };
    }
  }
  else {
    return { $can_kind: "ok", tag: { $can_kind: "Json__HeadEof" } };
  }
}
export function std__json__parse_literal(text: string, pos: bigint): { $can_kind: "ok"; next: bigint; val: Json__Value } | { $can_kind: "json.invalid_syntax"; detail: string } {
  if ($canStrAt(text, pos) === 116n) {
    if ((((BigInt([...text].length)) - pos) >= 4n)) {
      if (($canStrSlice(text, pos, (pos + 4n)) === "true")) {
        return { $can_kind: "ok", val: { $can_kind: "Json__Bool", value: true }, next: (pos + 4n) };
      }
      else {
        return { $can_kind: "json.invalid_syntax", detail: "bad literal" };
      }
    }
    else {
      return { $can_kind: "json.invalid_syntax", detail: "bad literal" };
    }
  }
  else if ($canStrAt(text, pos) === 102n) {
    if ((((BigInt([...text].length)) - pos) >= 5n)) {
      if (($canStrSlice(text, pos, (pos + 5n)) === "false")) {
        return { $can_kind: "ok", val: { $can_kind: "Json__Bool", value: false }, next: (pos + 5n) };
      }
      else {
        return { $can_kind: "json.invalid_syntax", detail: "bad literal" };
      }
    }
    else {
      return { $can_kind: "json.invalid_syntax", detail: "bad literal" };
    }
  }
  if ((((BigInt([...text].length)) - pos) >= 4n)) {
    if (($canStrSlice(text, pos, (pos + 4n)) === "null")) {
      return { $can_kind: "ok", val: { $can_kind: "Json__Null" }, next: (pos + 4n) };
    }
    else {
      return { $can_kind: "json.invalid_syntax", detail: "bad literal" };
    }
  }
  else {
    return { $can_kind: "json.invalid_syntax", detail: "bad literal" };
  }
}
export function std__json__parse_unescape(c: string): { $can_kind: "ok"; value: string } | { $can_kind: "json.invalid_syntax"; detail: string } {
  if (c === "\"") {
    return { $can_kind: "ok", value: "\"" };
  }
  else if (c === "\\") {
    return { $can_kind: "ok", value: "\\" };
  }
  else if (c === "/") {
    return { $can_kind: "ok", value: "/" };
  }
  else if (c === "n") {
    return { $can_kind: "ok", value: "\n" };
  }
  else if (c === "r") {
    return { $can_kind: "ok", value: "\r" };
  }
  else if (c === "t") {
    return { $can_kind: "ok", value: "\t" };
  }
  if ((c === "u")) {
    return { $can_kind: "json.invalid_syntax", detail: "unsupported escape" };
  }
  else {
    if ((c === "b")) {
      return { $can_kind: "json.invalid_syntax", detail: "unsupported escape" };
    }
    else {
      if ((c === "f")) {
        return { $can_kind: "json.invalid_syntax", detail: "unsupported escape" };
      }
      else {
        return { $can_kind: "json.invalid_syntax", detail: "bad escape" };
      }
    }
  }
}
export function std__json__parse_numcheck(text: string): { $can_kind: "ok"; value: boolean } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, 0n, ((BigInt([...text].length)) + 1n), 0n);
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
export function std__json__parse_numcheck_from(text: string, pos: bigint, n: bigint, state: bigint): { $can_kind: "ok"; value: boolean } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: false };
  }
  else {
    if ((pos >= (BigInt([...text].length)))) {
      if ((state === 2n)) {
        return { $can_kind: "ok", value: true };
      }
      else {
        if ((state === 3n)) {
          return { $can_kind: "ok", value: true };
        }
        else {
          if ((state === 5n)) {
            return { $can_kind: "ok", value: true };
          }
          else {
            if ((state === 8n)) {
              return { $can_kind: "ok", value: true };
            }
            else {
              return { $can_kind: "ok", value: false };
            }
          }
        }
      }
    }
    else {
      if (state === 0n) {
        if ($canStrAt(text, pos) === 45n) {
          const $can_m1: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 1n);
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
        else if ($canStrAt(text, pos) === 48n) {
          const $can_m2: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 2n);
          switch ($can_m2.$can_kind) {
          case "ok": {
            const r = $can_m2;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        if ($canBoolAnd(($canStrAt(text, pos) >= 49n), ($canStrAt(text, pos) <= 57n))) {
          const $can_m3: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 3n);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const r = $can_m3;
            return { $can_kind: "ok", value: r.value };
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
      else if (state === 1n) {
        if ($canStrAt(text, pos) === 48n) {
          const $can_m4: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 2n);
          switch ($can_m4.$can_kind) {
          case "ok": {
            const r = $can_m4;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        if ($canBoolAnd(($canStrAt(text, pos) >= 49n), ($canStrAt(text, pos) <= 57n))) {
          const $can_m5: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 3n);
          switch ($can_m5.$can_kind) {
          case "ok": {
            const r = $can_m5;
            return { $can_kind: "ok", value: r.value };
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
      else if (state === 2n) {
        if ($canStrAt(text, pos) === 46n) {
          const $can_m6: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 4n);
          switch ($can_m6.$can_kind) {
          case "ok": {
            const r = $can_m6;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else if ($canStrAt(text, pos) === 101n) {
          const $can_m7: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 6n);
          switch ($can_m7.$can_kind) {
          case "ok": {
            const r = $can_m7;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else if ($canStrAt(text, pos) === 69n) {
          const $can_m8: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 6n);
          switch ($can_m8.$can_kind) {
          case "ok": {
            const r = $can_m8;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        return { $can_kind: "ok", value: false };
      }
      else if (state === 3n) {
        if ($canStrAt(text, pos) === 46n) {
          const $can_m9: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 4n);
          switch ($can_m9.$can_kind) {
          case "ok": {
            const r = $can_m9;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else if ($canStrAt(text, pos) === 101n) {
          const $can_m10: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 6n);
          switch ($can_m10.$can_kind) {
          case "ok": {
            const r = $can_m10;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else if ($canStrAt(text, pos) === 69n) {
          const $can_m11: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 6n);
          switch ($can_m11.$can_kind) {
          case "ok": {
            const r = $can_m11;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        if ($canBoolAnd(($canStrAt(text, pos) >= 48n), ($canStrAt(text, pos) <= 57n))) {
          const $can_m12: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 3n);
          switch ($can_m12.$can_kind) {
          case "ok": {
            const r = $can_m12;
            return { $can_kind: "ok", value: r.value };
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
      else if (state === 4n) {
        if ($canBoolAnd(($canStrAt(text, pos) >= 48n), ($canStrAt(text, pos) <= 57n))) {
          const $can_m13: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 5n);
          switch ($can_m13.$can_kind) {
          case "ok": {
            const r = $can_m13;
            return { $can_kind: "ok", value: r.value };
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
      else if (state === 5n) {
        if ($canStrAt(text, pos) === 101n) {
          const $can_m14: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 6n);
          switch ($can_m14.$can_kind) {
          case "ok": {
            const r = $can_m14;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else if ($canStrAt(text, pos) === 69n) {
          const $can_m15: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 6n);
          switch ($can_m15.$can_kind) {
          case "ok": {
            const r = $can_m15;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        if ($canBoolAnd(($canStrAt(text, pos) >= 48n), ($canStrAt(text, pos) <= 57n))) {
          const $can_m16: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 5n);
          switch ($can_m16.$can_kind) {
          case "ok": {
            const r = $can_m16;
            return { $can_kind: "ok", value: r.value };
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
      else if (state === 6n) {
        if ($canStrAt(text, pos) === 43n) {
          const $can_m17: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 7n);
          switch ($can_m17.$can_kind) {
          case "ok": {
            const r = $can_m17;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else if ($canStrAt(text, pos) === 45n) {
          const $can_m18: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 7n);
          switch ($can_m18.$can_kind) {
          case "ok": {
            const r = $can_m18;
            return { $can_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        if ($canBoolAnd(($canStrAt(text, pos) >= 48n), ($canStrAt(text, pos) <= 57n))) {
          const $can_m19: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 8n);
          switch ($can_m19.$can_kind) {
          case "ok": {
            const r = $can_m19;
            return { $can_kind: "ok", value: r.value };
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
      else if (state === 7n) {
        if ($canBoolAnd(($canStrAt(text, pos) >= 48n), ($canStrAt(text, pos) <= 57n))) {
          const $can_m20: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 8n);
          switch ($can_m20.$can_kind) {
          case "ok": {
            const r = $can_m20;
            return { $can_kind: "ok", value: r.value };
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
      else if (state === 8n) {
        if ($canBoolAnd(($canStrAt(text, pos) >= 48n), ($canStrAt(text, pos) <= 57n))) {
          const $can_m21: { $can_kind: "ok"; value: boolean } = std__json__parse_numcheck_from(text, (pos + 1n), (n - 1n), 8n);
          switch ($can_m21.$can_kind) {
          case "ok": {
            const r = $can_m21;
            return { $can_kind: "ok", value: r.value };
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
      return { $can_kind: "ok", value: false };
    }
  }
}
export function std__json__parse_contains_key(keys: string[], key: string): { $can_kind: "ok"; value: boolean } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = std__json__parse_contains_key_from(keys, key, 0n, (BigInt([...keys].length)));
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
export function std__json__parse_contains_key_from(keys: string[], key: string, pos: bigint, n: bigint): { $can_kind: "ok"; value: boolean } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", value: false };
  }
  else {
    if ((pos < (BigInt([...keys].length)))) {
      if (($canSeqAt(keys, pos) === key)) {
        return { $can_kind: "ok", value: true };
      }
      else {
        const $can_m1: { $can_kind: "ok"; value: boolean } = std__json__parse_contains_key_from(keys, key, (pos + 1n), (n - 1n));
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
export function std__json__parse_pop(frames: Json__PFrame[]): { $can_kind: "ok"; frames: Json__PFrame[] } {
  const $can_m1: { $can_kind: "ok"; frames: Json__PFrame[] } = std__json__parse_pop_from(frames, 0n, ((BigInt([...frames].length)) - 1n), []);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", frames: r.frames };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__json__parse_pop_from(frames: Json__PFrame[], pos: bigint, n: bigint, acc: Json__PFrame[]): { $can_kind: "ok"; frames: Json__PFrame[] } {
  if ((n <= 0n)) {
    return { $can_kind: "ok", frames: acc };
  }
  else {
    if ((pos < (BigInt([...frames].length)))) {
      const $can_m1: { $can_kind: "ok"; frames: Json__PFrame[] } = std__json__parse_pop_from(frames, (pos + 1n), (n - 1n), [...acc, $canSeqAt(frames, pos)]);
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        return { $can_kind: "ok", frames: r.frames };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $can_kind: "ok", frames: acc };
    }
  }
}
export function std__json__parse_attach(frames: Json__PFrame[], v: Json__Value): { $can_kind: "ok"; frames: Json__PFrame[] } | { $can_kind: "json.invalid_syntax"; detail: string } {
  if (((BigInt([...frames].length)) === 0n)) {
    return { $can_kind: "ok", frames: [...frames, { tag: { $can_kind: "Json__PTail" }, items: [], fields: [], keys: [], key: "", acc: "", val: v, ofields: [], okeys: [] }] };
  }
  else {
    const $can_m1 = $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).tag;
    switch ($can_m1.$can_kind) {
    case "Json__PArrFirst": {
      const _ = $can_m1;
      const $can_m2: { $can_kind: "ok"; frames: Json__PFrame[] } = std__json__parse_pop(frames);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const pf = $can_m2;
        return { $can_kind: "ok", frames: [...pf.frames, { tag: { $can_kind: "Json__PArrNext" }, items: [...$canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).items, { value: v }], fields: [], keys: [], key: "", acc: "", val: { $can_kind: "Json__Null" }, ofields: [], okeys: [] }] };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    case "Json__PArrVal": {
      const _ = $can_m1;
      const $can_m3: { $can_kind: "ok"; frames: Json__PFrame[] } = std__json__parse_pop(frames);
      switch ($can_m3.$can_kind) {
      case "ok": {
        const pf = $can_m3;
        return { $can_kind: "ok", frames: [...pf.frames, { tag: { $can_kind: "Json__PArrNext" }, items: [...$canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).items, { value: v }], fields: [], keys: [], key: "", acc: "", val: { $can_kind: "Json__Null" }, ofields: [], okeys: [] }] };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    case "Json__PArrNext": {
      const _ = $can_m1;
      return { $can_kind: "json.invalid_syntax", detail: "trailing value" };
    }
    case "Json__PObjFirst": {
      const _ = $can_m1;
      return { $can_kind: "json.invalid_syntax", detail: "trailing value" };
    }
    case "Json__PObjKey": {
      const _ = $can_m1;
      return { $can_kind: "json.invalid_syntax", detail: "trailing value" };
    }
    case "Json__PObjColon": {
      const _ = $can_m1;
      return { $can_kind: "json.invalid_syntax", detail: "trailing value" };
    }
    case "Json__PObjKeyVal": {
      const _ = $can_m1;
      const $can_m4: { $can_kind: "ok"; frames: Json__PFrame[] } = std__json__parse_pop(frames);
      switch ($can_m4.$can_kind) {
      case "ok": {
        const pf = $can_m4;
        return { $can_kind: "ok", frames: [...pf.frames, { tag: { $can_kind: "Json__PObjNext" }, items: [], fields: [...$canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).fields, { name: $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).key, value: v }], keys: [...$canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).keys, $canSeqAt(frames, ((BigInt([...frames].length)) - 1n)).key], key: "", acc: "", val: { $can_kind: "Json__Null" }, ofields: [], okeys: [] }] };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    case "Json__PObjNext": {
      const _ = $can_m1;
      return { $can_kind: "json.invalid_syntax", detail: "trailing value" };
    }
    case "Json__PStrKey": {
      const _ = $can_m1;
      return { $can_kind: "json.invalid_syntax", detail: "trailing value" };
    }
    case "Json__PStrVal": {
      const _ = $can_m1;
      return { $can_kind: "json.invalid_syntax", detail: "trailing value" };
    }
    case "Json__PNumAcc": {
      const _ = $can_m1;
      return { $can_kind: "json.invalid_syntax", detail: "trailing value" };
    }
    case "Json__PTail": {
      const _ = $can_m1;
      return { $can_kind: "json.invalid_syntax", detail: "trailing value" };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
