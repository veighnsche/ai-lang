# auth-login — password login with one retry

- `db.can` — `mod db`: user lookup + user type, own tests.
- `auth.can` — `mod auth`: login with one retry, pins `db__get_user@5`
  + `Db__User@4` + `Db__Hash@1`. Passwords are branded (`Auth__Password`,
  `Db__Hash`), the check is a foreign `extern`, secrets enter via `seal`.
  The lockout-plus-password decision lives once in `auth__verify`, a
  same-module helper with no pin and no call-site `given`.
- `db.ts` + `auth.ts` — committed golden TS prod emits (tests/given
  stripped). Regenerate: `go run ./compiler --out <dir>
  sketches/auth-login/db.can sketches/auth-login/auth.can`; verify:
  `go test ./...`.

Rules: `/REQUIREMENTS.md`.

Note (rev rule): rev hashes cover code, not comments — editing a comment
never forces a rev bump.
