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
		checkGoldenFile(t, f, got, want)
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
		checkGoldenFile(t, f, got, want)
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
		checkGoldenFile(t, f, got, want)
	}
}

// TestGoldenReserve freezes the S2 second consumer: the reservation
// desk must transpile byte-identical. Since it pins quota validators
// and the scalars converter, those compile alongside; only reserve.ts
// and errors.json are golden-kept.
func TestGoldenReserve(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../sketches/reserve/reserve.can",
		"../std/quota/quota.can",
		"../std/scalars/scalars.can",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"reserve.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/reserve", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		checkGoldenFile(t, f, got, want)
	}
}

// TestGoldenNotify freezes the S2b third consumer: the notification
// desk must transpile byte-identical. Since it pins quota check ops,
// quota compiles alongside; only notify.ts and errors.json are
// golden-kept.
func TestGoldenNotify(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../sketches/notify/notify.can",
		"../std/quota/quota.can",
		"../std/scalars/scalars.can",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"notify.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/notify", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		checkGoldenFile(t, f, got, want)
	}
}

// TestGoldenQuotaCounter freezes row 1 of the stdlib program: the
// validation module plus its quota counter must transpile
// byte-identical, so validator payloads and the multi-shape ok union
// can never silently rot. Since S1a the quota module pins
// std__convert__int_to_str, so scalars.can compiles alongside;
// only quota.ts and errors.json are golden-kept.
func TestGoldenQuotaCounter(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../std/quota/quota.can",
		"../std/scalars/scalars.can",
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
		checkGoldenFile(t, f, got, want)
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
		checkGoldenFile(t, f, got, want)
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
		checkGoldenFile(t, f, got, want)
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
		checkGoldenFile(t, f, got, want)
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
		checkGoldenFile(t, f, got, want)
	}
}
