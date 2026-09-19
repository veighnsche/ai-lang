// GENERATED from ops.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type OpsResult = { $can_kind: "ok"; quotient: bigint; remainder: bigint } | { $can_kind: "ok"; value: string } | { $can_kind: "ops.zero_divisor"; detail: string };
export type Ops__Quot = { quotient: bigint; remainder: bigint };
export type Ops__Dec = { value: string };
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
export function ops__divmod(divisor: bigint, dividend: bigint): { $can_kind: "ok"; quotient: bigint; remainder: bigint } | { $can_kind: "ops.zero_divisor"; detail: string } {
  if (divisor === 0n) {
    return { $can_kind: "ops.zero_divisor", detail: "zero" };
  }
  return { $can_kind: "ok", quotient: $canDivMod(dividend, divisor)[0], remainder: $canDivMod(dividend, divisor)[1] };
}
export function ops__bump(step: string, x: string): { $can_kind: "ok"; value: string } {
  return { $can_kind: "ok", value: $canDecAdd(x, step) };
}
