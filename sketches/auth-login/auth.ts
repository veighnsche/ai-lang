// GENERATED from auth.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { db__get_user, type DbResult } from "./db";
import { auth__check_pw } from "./auth.externs";
export type AuthResult = { kind: "ok"; remaining_tries: bigint; user_id: string } | { kind: "auth.login_failed"; user_id: string } | { kind: "auth.account_locked"; user_id: string } | { kind: "auth.unavailable"; reason: string } | { kind: "auth.mismatch" };
export type Auth__Session = { user_id: string; remaining_tries: bigint };
export type Auth__Verdict = {};
export function auth__login(id: string, pw: string): AuthResult {
  const _m1: DbResult = db__get_user(id);
  switch (_m1.kind) {
  case "db.user_not_found":
    return { kind: "auth.login_failed", user_id: id };
  case "db.down":
    const _m2: DbResult = db__get_user(id);
    switch (_m2.kind) {
    case "db.user_not_found":
      return { kind: "auth.login_failed", user_id: id };
    case "db.down":
      return { kind: "auth.unavailable", reason: "db down twice" };
    case "ok":
      const user = _m2;
      const _m3: AuthResult = auth__verify(user.id, user.failed_attempts, pw, user.pw_hash);
      switch (_m3.kind) {
      case "auth.login_failed":
        return { kind: "auth.login_failed", user_id: user.id };
      case "auth.account_locked":
        return { kind: "auth.account_locked", user_id: user.id };
      case "ok":
        const s = _m3;
        return { kind: "ok", user_id: s.user_id, remaining_tries: s.remaining_tries };
      }
    }
  case "ok":
    const user = _m1;
    const _m4: AuthResult = auth__verify(user.id, user.failed_attempts, pw, user.pw_hash);
    switch (_m4.kind) {
    case "auth.login_failed":
      return { kind: "auth.login_failed", user_id: user.id };
    case "auth.account_locked":
      return { kind: "auth.account_locked", user_id: user.id };
    case "ok":
      const s = _m4;
      return { kind: "ok", user_id: s.user_id, remaining_tries: s.remaining_tries };
    }
  }
}
export function auth__verify(user_id: string, failed_attempts: bigint, pw: string, pw_hash: string): AuthResult {
  if ((failed_attempts >= 3n)) {
    return { kind: "auth.account_locked", user_id: user_id };
  }
  else {
    const _m5: { kind: "ok" } | { kind: "auth.mismatch" } = auth__check_pw(pw, pw_hash);
    switch (_m5.kind) {
    case "auth.mismatch":
      return { kind: "auth.login_failed", user_id: user_id };
    case "ok":
      const ok = _m5;
      return { kind: "ok", user_id: user_id, remaining_tries: (3n - failed_attempts) };
    }
  }
}
