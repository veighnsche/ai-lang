// GENERATED from auth.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { db__get_user, type DbResult } from "./db";
import { auth__check_pw } from "./auth.externs";
export type AuthResult = { $ail_kind: "ok"; remaining_tries: bigint; user_id: string } | { $ail_kind: "auth.login_failed"; user_id: string } | { $ail_kind: "auth.account_locked"; user_id: string } | { $ail_kind: "auth.unavailable"; reason: string } | { $ail_kind: "auth.mismatch" };
export type Auth__Session = { user_id: string; remaining_tries: bigint };
export type Auth__Verdict = {};
export function auth__login(id: string, pw: string): { $ail_kind: "ok"; remaining_tries: bigint; user_id: string } | { $ail_kind: "auth.login_failed"; user_id: string } | { $ail_kind: "auth.account_locked"; user_id: string } | { $ail_kind: "auth.unavailable"; reason: string } {
  const $ail_m1: { $ail_kind: "ok"; failed_attempts: bigint; id: string; pw_hash: string } | { $ail_kind: "db.user_not_found"; id: string } | { $ail_kind: "db.down" } = db__get_user(id);
  switch ($ail_m1.$ail_kind) {
  case "db.user_not_found": {
    const _ = $ail_m1;
    return { $ail_kind: "auth.login_failed", user_id: id };
  }
  case "db.down": {
    const _ = $ail_m1;
    const $ail_m2: { $ail_kind: "ok"; failed_attempts: bigint; id: string; pw_hash: string } | { $ail_kind: "db.user_not_found"; id: string } | { $ail_kind: "db.down" } = db__get_user(id);
    switch ($ail_m2.$ail_kind) {
    case "db.user_not_found": {
      const _ = $ail_m2;
      return { $ail_kind: "auth.login_failed", user_id: id };
    }
    case "db.down": {
      const _ = $ail_m2;
      return { $ail_kind: "auth.unavailable", reason: "db down twice" };
    }
    case "ok": {
      const user = $ail_m2;
      const $ail_m3: { $ail_kind: "ok"; remaining_tries: bigint; user_id: string } | { $ail_kind: "auth.login_failed"; user_id: string } | { $ail_kind: "auth.account_locked"; user_id: string } = auth__verify(user.id, user.failed_attempts, pw, user.pw_hash);
      switch ($ail_m3.$ail_kind) {
      case "auth.login_failed": {
        const _ = $ail_m3;
        return { $ail_kind: "auth.login_failed", user_id: user.id };
      }
      case "auth.account_locked": {
        const _ = $ail_m3;
        return { $ail_kind: "auth.account_locked", user_id: user.id };
      }
      case "ok": {
        const s = $ail_m3;
        return { $ail_kind: "ok", user_id: s.user_id, remaining_tries: s.remaining_tries };
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
    const user = $ail_m1;
    const $ail_m4: { $ail_kind: "ok"; remaining_tries: bigint; user_id: string } | { $ail_kind: "auth.login_failed"; user_id: string } | { $ail_kind: "auth.account_locked"; user_id: string } = auth__verify(user.id, user.failed_attempts, pw, user.pw_hash);
    switch ($ail_m4.$ail_kind) {
    case "auth.login_failed": {
      const _ = $ail_m4;
      return { $ail_kind: "auth.login_failed", user_id: user.id };
    }
    case "auth.account_locked": {
      const _ = $ail_m4;
      return { $ail_kind: "auth.account_locked", user_id: user.id };
    }
    case "ok": {
      const s = $ail_m4;
      return { $ail_kind: "ok", user_id: s.user_id, remaining_tries: s.remaining_tries };
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
export function auth__verify(user_id: string, failed_attempts: bigint, pw: string, pw_hash: string): { $ail_kind: "ok"; remaining_tries: bigint; user_id: string } | { $ail_kind: "auth.login_failed"; user_id: string } | { $ail_kind: "auth.account_locked"; user_id: string } {
  if ((failed_attempts >= 3n)) {
    return { $ail_kind: "auth.account_locked", user_id: user_id };
  }
  else {
    const $ail_m1: { $ail_kind: "ok" } | { $ail_kind: "auth.mismatch" } = auth__check_pw(pw, pw_hash);
    switch ($ail_m1.$ail_kind) {
    case "auth.mismatch": {
      const _ = $ail_m1;
      return { $ail_kind: "auth.login_failed", user_id: user_id };
    }
    case "ok": {
      const ok = $ail_m1;
      return { $ail_kind: "ok", user_id: user_id, remaining_tries: (3n - failed_attempts) };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
