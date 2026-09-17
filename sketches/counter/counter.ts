// GENERATED from counter.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type CounterResult = { $ail_kind: "ok"; total: bigint };
export type Count__Tally = { total: bigint };
let Count__total: bigint = 0n;
export function count__bump(by: bigint): CounterResult {
  const $ail_m1: { $ail_kind: "ok", value: bigint } = { $ail_kind: "ok", value: Count__total };
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const c = $ail_m1;
    Count__total = (c.value + by);
    const $ail_m2: { $ail_kind: "ok" } = { $ail_kind: "ok" };
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const _ = $ail_m2;
      return { $ail_kind: "ok", total: (c.value + by) };
    }
    }
  }
  }
}
export function count__twice(by: bigint): CounterResult {
  const $ail_m1: CounterResult = count__bump(by);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const s = $ail_m1;
    const $ail_m2: CounterResult = count__bump(by);
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const t = $ail_m2;
      return { $ail_kind: "ok", total: t.total };
    }
    }
  }
  }
}
