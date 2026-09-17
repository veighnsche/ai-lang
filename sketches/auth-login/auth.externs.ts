// Extern-stub fixture for the a70 tsc gate. Hand-written; signature
// mirrors `extern auth__check_pw` in sketches/auth-login/auth.ail
// (brands erase to string) and its call site in auth.ts. Never
// executed by the gate (`noEmit` check only). If the extern decl or
// emit shape changes, update this stub in the same slice.
export function auth__check_pw(
  pw: string,
  hash: string,
): { $ail_kind: "ok" } | { $ail_kind: "auth.mismatch" } {
  throw new Error("extern stub: auth__check_pw is not implemented");
}
