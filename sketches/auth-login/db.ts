// GENERATED from db.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type DbResult = { kind: "ok"; failed_attempts: bigint; id: string; pw_hash: string } | { kind: "db.user_not_found"; id: string } | { kind: "db.down" };
export type Db__User = { id: string; pw_hash: string; failed_attempts: bigint };
export function db__get_user(id: string): DbResult {
  if (id === "u_01") {
    return { kind: "ok", id: "u_01", pw_hash: "secret", failed_attempts: 0n };
  }
  else if (id === "u_02") {
    return { kind: "ok", id: "u_02", pw_hash: "second", failed_attempts: 5n };
  }
  return { kind: "db.user_not_found", id: id };
}
