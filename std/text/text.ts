// GENERATED from text.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type TextResult = { $ail_kind: "ok"; value: boolean } | { $ail_kind: "ok"; value: string };
export type Str__Value = { value: string };
export type Bool__Value = { value: boolean };
export function std__str__concat(left: string, right: string): TextResult {
  return { $ail_kind: "ok", value: (left + right) };
}
export function std__str__is_empty(value: string): TextResult {
  if ((value === "")) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
