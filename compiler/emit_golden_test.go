package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGoldenAuthLogin re-transpiles the auth-login example and requires
// byte-identical output to the committed goldens. compile() already fails
// the test on any signature-test failure, so this gate covers parse,
// test evaluation, and emit in one shot.
func TestGoldenAuthLogin(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../sketches/auth-login/db.ail",
		"../sketches/auth-login/auth.ail",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"db.ts", "auth.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/auth-login", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		if string(got) != string(want) {
			t.Errorf("golden mismatch: %s (re-run ailc and inspect the diff)", f)
		}
	}
}

// TestGoldenRetryLoop freezes the termination + iteration demo: the
// retry-with-fuel program must transpile byte-identical, so the loop
// shape and its recursive emit can never silently rot.
func TestGoldenRetryLoop(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../sketches/retry-loop/retry.ail",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"retry.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/retry-loop", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		if string(got) != string(want) {
			t.Errorf("golden mismatch: %s (re-run ailc and inspect the diff)", f)
		}
	}
}

// TestGoldenCounter freezes the effects demo: the bounded counter
// must transpile byte-identical, so cell lets and accessor matches
// can never silently rot.
func TestGoldenCounter(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../sketches/counter/counter.ail",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"counter.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/counter", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		if string(got) != string(want) {
			t.Errorf("golden mismatch: %s (re-run ailc and inspect the diff)", f)
		}
	}
}
