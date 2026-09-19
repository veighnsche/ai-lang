// GENERATED from use.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { ops__bump, ops__divmod, type Ops__Dec, type Ops__Quot, type OpsResult } from "./ops";
export type UseResult = { $can_kind: "ok"; cb: (input: bigint) => { $can_kind: "ok"; quotient: bigint; remainder: bigint } | { $can_kind: "ops.zero_divisor"; detail: string } } | { $can_kind: "ok"; cb: (input: string) => { $can_kind: "ok"; value: string } } | { $can_kind: "ok"; quotient: bigint; remainder: bigint } | { $can_kind: "ok"; value: string };
export type Use__Cb = { cb: (input: bigint) => { $can_kind: "ok"; quotient: bigint; remainder: bigint } | { $can_kind: "ops.zero_divisor"; detail: string } };
export type Use__CbDec = { cb: (input: string) => { $can_kind: "ok"; value: string } };
export function use__bind(d: bigint): { $can_kind: "ok"; cb: (input: bigint) => { $can_kind: "ok"; quotient: bigint; remainder: bigint } | { $can_kind: "ops.zero_divisor"; detail: string } } {
  return { $can_kind: "ok", cb: (() => { const $can_m1 = d; return (input) => ops__divmod($can_m1, input); })() };
}
export function use__run(cb: (input: bigint) => { $can_kind: "ok"; quotient: bigint; remainder: bigint } | { $can_kind: "ops.zero_divisor"; detail: string }, n: bigint): { $can_kind: "ok"; quotient: bigint; remainder: bigint } {
  const $can_m1: { $can_kind: "ok"; quotient: bigint; remainder: bigint } | { $can_kind: "ops.zero_divisor"; detail: string } = cb(n);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", quotient: r.quotient, remainder: r.remainder };
  }
  case "ops.zero_divisor": {
    const _ = $can_m1;
    return { $can_kind: "ok", quotient: 0n, remainder: 0n };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function use__go(d: bigint, n: bigint): { $can_kind: "ok"; quotient: bigint; remainder: bigint } {
  const $can_m1: { $can_kind: "ok"; cb: (input: bigint) => { $can_kind: "ok"; quotient: bigint; remainder: bigint } | { $can_kind: "ops.zero_divisor"; detail: string } } = use__bind(d);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const h = $can_m1;
    const $can_m2: { $can_kind: "ok"; quotient: bigint; remainder: bigint } = use__run(h.cb, n);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const r = $can_m2;
      return { $can_kind: "ok", quotient: r.quotient, remainder: r.remainder };
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
export function use__binddec(s: string): { $can_kind: "ok"; cb: (input: string) => { $can_kind: "ok"; value: string } } {
  return { $can_kind: "ok", cb: (() => { const $can_m1 = s; return (input) => ops__bump($can_m1, input); })() };
}
export function use__rundec(cb: (input: string) => { $can_kind: "ok"; value: string }, x: string): { $can_kind: "ok"; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } = cb(x);
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
