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
		"../sketches/auth-login/db.can",
		"../sketches/auth-login/auth.can",
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
			t.Errorf("golden mismatch: %s (re-run canlc and inspect the diff)", f)
		}
	}
}

// TestGoldenRetryLoop freezes the termination + iteration demo: the
// retry-with-fuel program must transpile byte-identical, so the loop
// shape and its recursive emit can never silently rot.
func TestGoldenRetryLoop(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../sketches/retry-loop/retry.can",
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
			t.Errorf("golden mismatch: %s (re-run canlc and inspect the diff)", f)
		}
	}
}

// TestGoldenCounter freezes the effects demo: the bounded counter
// must transpile byte-identical, so cell lets and accessor matches
// can never silently rot.
func TestGoldenCounter(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../sketches/counter/counter.can",
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
			t.Errorf("golden mismatch: %s (re-run canlc and inspect the diff)", f)
		}
	}
}

// TestGoldenQuotaCounter freezes row 1 of the stdlib program: the
// validation module plus its quota counter must transpile
// byte-identical, so validator payloads and the multi-shape ok union
// can never silently rot.
func TestGoldenQuotaCounter(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../std/quota/quota.can",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"quota.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../std/quota", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		if string(got) != string(want) {
			t.Errorf("golden mismatch: %s (re-run canlc and inspect the diff)", f)
		}
	}
}

// TestGoldenStdText freezes the text-construction cut: the concat
// module must transpile byte-identical, so its decision tables can
// never silently rot.
func TestGoldenStdText(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../std/text/text.can",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"text.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../std/text", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		if string(got) != string(want) {
			t.Errorf("golden mismatch: %s (re-run canlc and inspect the diff)", f)
		}
	}
}

// TestGoldenStdDivision freezes the division cut: the Euclidean
// module must transpile byte-identical, so its decision tables and
// the $canDivMod helper emit can never silently rot.
func TestGoldenStdDivision(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../std/division/division.can",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"division.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../std/division", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		if string(got) != string(want) {
			t.Errorf("golden mismatch: %s (re-run canlc and inspect the diff)", f)
		}
	}
}

// TestGoldenStdScalars freezes row 2 of the stdlib program: the
// monomorphic scalar catalog must transpile byte-identical, so its
// decision tables and exact-decimal emit can never silently rot.
func TestGoldenStdScalars(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../std/scalars/scalars.can",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"scalars.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../std/scalars", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		if string(got) != string(want) {
			t.Errorf("golden mismatch: %s (re-run canlc and inspect the diff)", f)
		}
	}
}

// TestGoldenStdHtml freezes the html-constructor cut: the brand and
// fragment module must transpile byte-identical, so its decision
// tables can never silently rot. (No golden covered html before
// the fragment slice; the join rows are the first to pin it.)
func TestGoldenStdHtml(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../std/html/html.can",
		"../std/ascii/ascii.can",
		"../std/scalars/scalars.can",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"html.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../std/html", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		if string(got) != string(want) {
			t.Errorf("golden mismatch: %s (re-run canlc and inspect the diff)", f)
		}
	}
}
