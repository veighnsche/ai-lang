package main

import (
	"bytes"
	"math/big"
	"testing"
)

// TestGoldenNormalize freezes the canonical value rendering over the
// auth-login example. Every line below was verified by hand against the
// sketches: sorted record keys (failed_attempts first, not declaration
// order), err kinds with sorted fields, stable sort of the row labels.
func TestGoldenNormalize(t *testing.T) {
	var buf bytes.Buffer
	err := runNormalize(&buf, []string{"../sketches/auth-login/db.can", "../sketches/auth-login/auth.can"})
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	const want = `auth.auth__login/badpw => err(auth.login_failed(user_id = "u_01"))
auth.auth__login/down => err(auth.unavailable(reason = "db down twice"))
auth.auth__login/flaky => Ok(remaining_tries = 3, user_id = "u_01")
auth.auth__login/happy => Ok(remaining_tries = 3, user_id = "u_01")
auth.auth__login/locked => err(auth.account_locked(user_id = "u_02"))
auth.auth__login/missing => err(auth.login_failed(user_id = "u_99"))
auth.auth__login/retry_badpw => err(auth.login_failed(user_id = "u_01"))
auth.auth__login/retry_locked => err(auth.account_locked(user_id = "u_02"))
auth.auth__login/retry_missing => err(auth.login_failed(user_id = "u_99"))
auth.auth__verify/vh_bad => err(auth.login_failed(user_id = "u_01"))
auth.auth__verify/vh_locked => err(auth.account_locked(user_id = "u_01"))
auth.auth__verify/vh_ok => Ok(remaining_tries = 3, user_id = "u_01")
db.db__get_user/known_user => Ok(failed_attempts = 0, id = "u_01", pw_hash = "secret")
db.db__get_user/locked_user => Ok(failed_attempts = 5, id = "u_02", pw_hash = "second")
db.db__get_user/unknown_user => err(db.user_not_found(id = "u_99"))
`
	if got := buf.String(); got != want {
		t.Fatalf("normalize golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestGoldenNormalizeRetryLoop freezes the retry-loop decision table:
// every row below was verified by hand against the sketch — now hits
// the server once, later retries once, never burns fuel 2 to exhausted,
// empty stops at once with no fuel spent.
func TestGoldenNormalizeRetryLoop(t *testing.T) {
	var buf bytes.Buffer
	err := runNormalize(&buf, []string{"../sketches/retry-loop/retry.can"})
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	const want = `retry.retry__fetch/empty => err(retry.exhausted)
retry.retry__fetch/later => Ok(body = "hi")
retry.retry__fetch/never => err(retry.exhausted)
retry.retry__fetch/now => Ok(body = "hi")
`
	if got := buf.String(); got != want {
		t.Fatalf("normalize golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestGoldenNormalizeCounter freezes the counter decision table: the
// two three-rows agree by per-test isolation (a shared store would
// print 6 on the second), and double threads two bumps to 4.
func TestGoldenNormalizeCounter(t *testing.T) {
	var buf bytes.Buffer
	err := runNormalize(&buf, []string{"../sketches/counter/counter.can"})
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	const want = `counter.count__bump/three => Ok(total = 3)
counter.count__bump/three_again => Ok(total = 3)
counter.count__bump/zero => Ok(total = 0)
counter.count__twice/double => Ok(total = 4)
`
	if got := buf.String(); got != want {
		t.Fatalf("normalize golden mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// TestNormalizeValueShapes pins the canonical grammar on unit values:
// escaping, nesting, empty records, and field-less errors.
func TestNormalizeValueShapes(t *testing.T) {
	cases := []struct {
		v    *Value
		want string
	}{
		{&Value{Kind: "str", S: "a\"b\\c\nd"}, `"a\"b\\c\nd"`},
		{&Value{Kind: "int", N: big.NewInt(-3)}, "-3"},
		{&Value{Kind: "dec", D: "1.5"}, `d"1.5"`},
		{&Value{Kind: "bool", B: true}, "true"},
		{&Value{Kind: "ok", Dict: map[string]*Value{}}, "Ok()"},
		{&Value{Kind: "err", ErrKind: "db.down"}, "err(db.down)"},
		{&Value{Kind: "ok", Dict: map[string]*Value{
			"b": {Kind: "int", N: big.NewInt(1)},
			"a": {Kind: "ok", Dict: map[string]*Value{
				"z": {Kind: "bool", B: false},
			}},
		}}, `Ok(a = Ok(z = false), b = 1)`},
	}
	for _, c := range cases {
		if got := normalizeValue(c.v); got != c.want {
			t.Fatalf("normalizeValue(%v) = %q, want %q", c.v, got, c.want)
		}
	}
}
