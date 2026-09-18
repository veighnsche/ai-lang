// GENERATED from hello.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type HelloResult = { $ail_kind: "ok"; text: string };
export type Hello__Message = { text: string };
export function hello__world(): { $ail_kind: "ok"; text: string } {
  return { $ail_kind: "ok", text: "Hello, world!" };
}
