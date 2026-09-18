// GENERATED from ascii.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type AsciiResult = { $ail_kind: "ok"; value: boolean };
// Strict boolean runtime (slice 5): eager helpers, never bare &&.
function $ailBoolAnd(a: boolean, b: boolean): boolean {
  return a && b;
}
function $ailBoolOr(a: boolean, b: boolean): boolean {
  return a || b;
}
export function std__ascii__is_digit(code: bigint): { $ail_kind: "ok"; value: boolean } {
  return { $ail_kind: "ok", value: $ailBoolAnd((code >= 48n), (code <= 57n)) };
}
export function std__ascii__is_alpha(code: bigint): { $ail_kind: "ok"; value: boolean } {
  return { $ail_kind: "ok", value: $ailBoolOr($ailBoolAnd((code >= 97n), (code <= 122n)), $ailBoolAnd((code >= 65n), (code <= 90n))) };
}
export function std__ascii__is_alnum(code: bigint): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1: { $ail_kind: "ok"; value: boolean } = std__ascii__is_alpha(code);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const a = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; value: boolean } = std__ascii__is_digit(code);
    switch ($ail_m2.$ail_kind) {
    case "ok": {
      const d = $ail_m2;
      return { $ail_kind: "ok", value: $ailBoolOr(a.value, d.value) };
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
