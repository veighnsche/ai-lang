// GENERATED from counter.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type CounterResult = { kind: "ok"; total: bigint };
export type Count__Tally = { total: bigint };
let Count__total: bigint = 0n;
export function count__bump(by: bigint): CounterResult {
  const $ail_m1: { kind: "ok", value: bigint } = { kind: "ok", value: Count__total };
  switch ($ail_m1.kind) {
  case "ok": {
    const c = $ail_m1;
    Count__total = (c.value + by);
    const $ail_m2: { kind: "ok" } = { kind: "ok" };
    switch ($ail_m2.kind) {
    case "ok": {
      const _ = $ail_m2;
      return { kind: "ok", total: (c.value + by) };
    }
    }
  }
  }
}
export function count__twice(by: bigint): CounterResult {
  const $ail_m3: CounterResult = count__bump(by);
  switch ($ail_m3.kind) {
  case "ok": {
    const s = $ail_m3;
    const $ail_m4: CounterResult = count__bump(by);
    switch ($ail_m4.kind) {
    case "ok": {
      const t = $ail_m4;
      return { kind: "ok", total: t.total };
    }
    }
  }
  }
}
