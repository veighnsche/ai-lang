// GENERATED from ascii.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type AsciiResult = { $can_kind: "ok"; value: boolean };
// Strict boolean runtime (slice 5): eager helpers, never bare &&.
function $canBoolAnd(a: boolean, b: boolean): boolean {
  return a && b;
}
function $canBoolOr(a: boolean, b: boolean): boolean {
  return a || b;
}
export function std__ascii__is_digit(code: bigint): { $can_kind: "ok"; value: boolean } {
  return { $can_kind: "ok", value: $canBoolAnd((code >= 48n), (code <= 57n)) };
}
export function std__ascii__is_alpha(code: bigint): { $can_kind: "ok"; value: boolean } {
  return { $can_kind: "ok", value: $canBoolOr($canBoolAnd((code >= 97n), (code <= 122n)), $canBoolAnd((code >= 65n), (code <= 90n))) };
}
export function std__ascii__is_alnum(code: bigint): { $can_kind: "ok"; value: boolean } {
  const $can_m1: { $can_kind: "ok"; value: boolean } = std__ascii__is_alpha(code);
  switch ($can_m1.$can_kind) {
  case "ok": {
    const a = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: boolean } = std__ascii__is_digit(code);
    switch ($can_m2.$can_kind) {
    case "ok": {
      const d = $can_m2;
      return { $can_kind: "ok", value: $canBoolOr(a.value, d.value) };
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
