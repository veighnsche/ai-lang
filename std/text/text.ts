// GENERATED from text.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type TextResult = { kind: "ok"; value: string };
export type Str__Value = { value: string };
export function std__str__concat(left: string, right: string): TextResult {
  return { kind: "ok", value: (left + right) };
}
