// GENERATED from db.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type DbResult = { $ail_kind: "ok"; failed_attempts: bigint; id: string; pw_hash: string } | { $ail_kind: "db.user_not_found"; id: string } | { $ail_kind: "db.down" };
export type Db__User = { id: string; pw_hash: string; failed_attempts: bigint };
export function db__get_user(id: string): { $ail_kind: "ok"; failed_attempts: bigint; id: string; pw_hash: string } | { $ail_kind: "db.user_not_found"; id: string } | { $ail_kind: "db.down" } {
  if (id === "u_01") {
    return { $ail_kind: "ok", id: "u_01", pw_hash: "secret", failed_attempts: 0n };
  }
  else if (id === "u_02") {
    return { $ail_kind: "ok", id: "u_02", pw_hash: "second", failed_attempts: 5n };
  }
  return { $ail_kind: "db.user_not_found", id: id };
}
