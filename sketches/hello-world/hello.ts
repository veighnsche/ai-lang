// GENERATED from hello.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type HelloResult = { $can_kind: "ok"; text: string };
export type Hello__Message = { text: string };
export function hello__world(): { $can_kind: "ok"; text: string } {
  return { $can_kind: "ok", text: "Hello, world!" };
}
