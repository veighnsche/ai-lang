// GENERATED from counter.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type CounterResult = { kind: "ok"; total: bigint };
export type Count__Tally = { total: bigint };
let Count__total: bigint = 0n;
export function count__bump(by: bigint): CounterResult {
  const _m1: { kind: "ok", value: bigint } = { kind: "ok", value: Count__total };
  switch (_m1.kind) {
  case "ok":
    const c = _m1;
    Count__total = (c.value + by);
    const _m2: { kind: "ok" } = { kind: "ok" };
    switch (_m2.kind) {
    case "ok":
      const _ = _m2;
      return { kind: "ok", total: (c.value + by) };
    }
  }
}
export function count__twice(by: bigint): CounterResult {
  const _m3: CounterResult = count__bump(by);
  switch (_m3.kind) {
  case "ok":
    const s = _m3;
    const _m4: CounterResult = count__bump(by);
    switch (_m4.kind) {
    case "ok":
      const t = _m4;
      return { kind: "ok", total: t.total };
    }
  }
}
