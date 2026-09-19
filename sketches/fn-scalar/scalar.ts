// GENERATED from scalar.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type ScalarResult = { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: boolean } | { $can_kind: "ok"; value: string } | { $can_kind: "scalar.negative"; value: bigint };
// Exact-decimal runtime (a10): canonical-digit strings, BigInt math.
function $canDecSplit(d: string): { neg: boolean; ip: string; fp: string } {
  let neg = false;
  if (d.startsWith("-")) {
    neg = true;
    d = d.slice(1);
  }
  const dot = d.indexOf(".");
  return { neg, ip: d.slice(0, dot), fp: d.slice(dot + 1) };
}
function $canDecNorm(ip: string, fp: string, neg: boolean): string {
  ip = ip.replace(/^0+(?=\d)/, "");
  fp = fp.replace(/0+$/, "");
  if (fp === "") {
    fp = "0";
  }
  if (ip === "0" && fp === "0") {
    return "0.0";
  }
  return (neg ? "-" : "") + ip + "." + fp;
}
function $canDecMant(p: { neg: boolean; ip: string; fp: string }, scale: number): bigint {
  let f = p.fp;
  while (f.length < scale) {
    f += "0";
  }
  const m = BigInt(p.ip + f);
  return p.neg ? -m : m;
}
function $canDecFromMant(m: bigint, scale: number): string {
  let neg = false;
  if (m < 0n) {
    neg = true;
    m = -m;
  }
  let digits = m.toString();
  while (digits.length < scale + 1) {
    digits = "0" + digits;
  }
  return $canDecNorm(digits.slice(0, digits.length - scale), digits.slice(digits.length - scale), neg);
}
function $canDecAdd(a: string, b: string): string {
  const A = $canDecSplit(a);
  const B = $canDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $canDecFromMant($canDecMant(A, s) + $canDecMant(B, s), s);
}
export function scalar__identity$T$bool(value: boolean): { $can_kind: "ok"; value: boolean } {
  return { $can_kind: "ok", value: value };
}
export function scalar__identity$T$dec(value: string): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: value };
}
export function scalar__identity$T$int(value: bigint): { $can_kind: "ok"; value: bigint } {
  return { $can_kind: "ok", value: value };
}
export function scalar__identity$T$str(value: string): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: value };
}
export function scalar__add(delta: bigint, value: bigint): { $can_kind: "ok"; value: bigint } {
  return { $can_kind: "ok", value: (value + delta) };
}
export function scalar__positive(value: bigint): { $can_kind: "ok"; value: boolean } {
  return { $can_kind: "ok", value: (value > 0n) };
}
export function scalar__twice(value: string): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: $canDecAdd(value, value) };
}
export function scalar__checked(value: bigint): { $can_kind: "ok"; value: bigint } | { $can_kind: "scalar.negative"; value: bigint } {
  if ((value < 0n)) {
    return { $can_kind: "scalar.negative", value: value };
  }
  else {
    return { $can_kind: "ok", value: value };
  }
}
