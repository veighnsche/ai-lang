// GENERATED from html.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type HtmlResult = { $ail_kind: "ok"; safe: string } | { $ail_kind: "ok"; text: string } | { $ail_kind: "ok"; value: string };
export type Html__Escaped = { value: string };
export type Html__TextResult = { text: string };
export type Html__SafeResult = { safe: string };
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
function $ailStrSlice(s: string, a: bigint, b: bigint): string {
  const cps = [...s];
  const toIdx = (x: bigint): number => { if (x < 0n || x > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("str slice out of range"); return Number(x); };
  const lo = toIdx(a), hi = toIdx(b);
  if (lo > hi || hi > cps.length) throw new Error("str slice out of range");
  return cps.slice(lo, hi).join("");
}
export function html__text__escape_from(s: string, acc: string, n: bigint): HtmlResult {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: acc };
  }
  else {
    if (($ailStrSlice(s, 0n, 1n) === "&")) {
      const $ail_m1: HtmlResult = html__text__escape_from($ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + "&amp;"), (n - 1n));
      switch ($ail_m1.$ail_kind) {
      case "ok": {
        const r = $ail_m1;
        return { $ail_kind: "ok", value: r.value };
      }
      }
    }
    else {
      if (($ailStrSlice(s, 0n, 1n) === "<")) {
        const $ail_m2: HtmlResult = html__text__escape_from($ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + "&lt;"), (n - 1n));
        switch ($ail_m2.$ail_kind) {
        case "ok": {
          const r = $ail_m2;
          return { $ail_kind: "ok", value: r.value };
        }
        }
      }
      else {
        if (($ailStrSlice(s, 0n, 1n) === ">")) {
          const $ail_m3: HtmlResult = html__text__escape_from($ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + "&gt;"), (n - 1n));
          switch ($ail_m3.$ail_kind) {
          case "ok": {
            const r = $ail_m3;
            return { $ail_kind: "ok", value: r.value };
          }
          }
        }
        else {
          const $ail_m4: HtmlResult = html__text__escape_from($ailStrSlice(s, 1n, (BigInt([...s].length))), (acc + $ailStrSlice(s, 0n, 1n)), (n - 1n));
          switch ($ail_m4.$ail_kind) {
          case "ok": {
            const r = $ail_m4;
            return { $ail_kind: "ok", value: r.value };
          }
          }
        }
      }
    }
  }
}
export function html__text__escape(raw: string): HtmlResult {
  const $ail_m5: HtmlResult = html__text__escape_from(raw, "", (BigInt([...raw].length)));
  switch ($ail_m5.$ail_kind) {
  case "ok": {
    const r = $ail_m5;
    return { $ail_kind: "ok", text: r.value };
  }
  }
}
export function html__text__node(text: string): HtmlResult {
  return { $ail_kind: "ok", safe: text };
}
