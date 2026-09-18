// GENERATED from auth.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { db__get_user, type DbResult } from "./db";
import { auth__check_pw } from "./auth.externs";
export type AuthResult = { $can_kind: "ok"; remaining_tries: bigint; user_id: string } | { $can_kind: "auth.login_failed"; user_id: string } | { $can_kind: "auth.account_locked"; user_id: string } | { $can_kind: "auth.unavailable"; reason: string } | { $can_kind: "auth.mismatch" };
export type Auth__Session = { user_id: string; remaining_tries: bigint };
export type Auth__Verdict = {};
export function auth__login(id: string, pw: string): { $can_kind: "ok"; remaining_tries: bigint; user_id: string } | { $can_kind: "auth.login_failed"; user_id: string } | { $can_kind: "auth.account_locked"; user_id: string } | { $can_kind: "auth.unavailable"; reason: string } {
  const $can_m1: { $can_kind: "ok"; failed_attempts: bigint; id: string; pw_hash: string } | { $can_kind: "db.user_not_found"; id: string } | { $can_kind: "db.down" } = db__get_user(id);
  switch ($can_m1.$can_kind) {
  case "db.user_not_found": {
    const _ = $can_m1;
    return { $can_kind: "auth.login_failed", user_id: id };
  }
  case "db.down": {
    const _ = $can_m1;
    const $can_m2: { $can_kind: "ok"; failed_attempts: bigint; id: string; pw_hash: string } | { $can_kind: "db.user_not_found"; id: string } | { $can_kind: "db.down" } = db__get_user(id);
    switch ($can_m2.$can_kind) {
    case "db.user_not_found": {
      const _ = $can_m2;
      return { $can_kind: "auth.login_failed", user_id: id };
    }
    case "db.down": {
      const _ = $can_m2;
      return { $can_kind: "auth.unavailable", reason: "db down twice" };
    }
    case "ok": {
      const user = $can_m2;
      const $can_m3: { $can_kind: "ok"; remaining_tries: bigint; user_id: string } | { $can_kind: "auth.login_failed"; user_id: string } | { $can_kind: "auth.account_locked"; user_id: string } = auth__verify(user.id, user.failed_attempts, pw, user.pw_hash);
      switch ($can_m3.$can_kind) {
      case "auth.login_failed": {
        const _ = $can_m3;
        return { $can_kind: "auth.login_failed", user_id: user.id };
      }
      case "auth.account_locked": {
        const _ = $can_m3;
        return { $can_kind: "auth.account_locked", user_id: user.id };
      }
      case "ok": {
        const s = $can_m3;
        return { $can_kind: "ok", user_id: s.user_id, remaining_tries: s.remaining_tries };
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
  case "ok": {
    const user = $can_m1;
    const $can_m4: { $can_kind: "ok"; remaining_tries: bigint; user_id: string } | { $can_kind: "auth.login_failed"; user_id: string } | { $can_kind: "auth.account_locked"; user_id: string } = auth__verify(user.id, user.failed_attempts, pw, user.pw_hash);
    switch ($can_m4.$can_kind) {
    case "auth.login_failed": {
      const _ = $can_m4;
      return { $can_kind: "auth.login_failed", user_id: user.id };
    }
    case "auth.account_locked": {
      const _ = $can_m4;
      return { $can_kind: "auth.account_locked", user_id: user.id };
    }
    case "ok": {
      const s = $can_m4;
      return { $can_kind: "ok", user_id: s.user_id, remaining_tries: s.remaining_tries };
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
export function auth__verify(user_id: string, failed_attempts: bigint, pw: string, pw_hash: string): { $can_kind: "ok"; remaining_tries: bigint; user_id: string } | { $can_kind: "auth.login_failed"; user_id: string } | { $can_kind: "auth.account_locked"; user_id: string } {
  if ((failed_attempts >= 3n)) {
    return { $can_kind: "auth.account_locked", user_id: user_id };
  }
  else {
    const $can_m1: { $can_kind: "ok" } | { $can_kind: "auth.mismatch" } = auth__check_pw(pw, pw_hash);
    switch ($can_m1.$can_kind) {
    case "auth.mismatch": {
      const _ = $can_m1;
      return { $can_kind: "auth.login_failed", user_id: user_id };
    }
    case "ok": {
      const ok = $can_m1;
      return { $can_kind: "ok", user_id: user_id, remaining_tries: (3n - failed_attempts) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
