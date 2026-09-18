package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepoGrammar(t *testing.T) {
	if _, err := os.Stat("../../editors/vscode/syntaxes/can.tmGrammar.json"); err != nil {
		t.Skip("not in repo checkout")
	}
	if errs := check("../../editors/vscode"); len(errs) > 0 {
		t.Fatalf("repo grammar failed: %v", errs)
	}
}

func writeGrammar(t *testing.T, patterns string) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		"package.json":                `{"name": "x"}`,
		"language-configuration.json": `{}`,
		"syntaxes/can.tmGrammar.json": `{"scopeName": "source.can", "patterns": [` + patterns + `]}`,
	}
	for name, body := range files {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

const errorRule = `{"match": "\\berror\\b", "name": "keyword.declaration.error.can"}`
const keywordRule = `{"match": "\\b(fn|on|rev)\\b", "name": "keyword.control.can"}`

func TestMissingErrorRule(t *testing.T) {
	dir := writeGrammar(t, keywordRule)
	if errs := check(dir); !contains(errs, "exactly one error rule") {
		t.Fatalf("expected error-rule violation, got %v", errs)
	}
}

func TestMissingRevKeyword(t *testing.T) {
	noRev := `{"match": "\\b(fn|on)\\b", "name": "keyword.control.can"}`
	pin := `{"match": "@[0-9]+", "name": "constant.numeric.version.can"}`
	dir := writeGrammar(t, errorRule+","+noRev+","+pin)
	if errs := check(dir); !contains(errs, "rev must be a keyword") {
		t.Fatalf("expected rev violation, got %v", errs)
	}
}

func TestBadJSON(t *testing.T) {
	dir := writeGrammar(t, keywordRule)
	bad := filepath.Join(dir, "syntaxes", "can.tmGrammar.json")
	if err := os.WriteFile(bad, []byte(`{not json`), 0o644); err != nil {
		t.Fatal(err)
	}
	if errs := check(dir); len(errs) == 0 {
		t.Fatal("expected JSON violation, got none")
	}
}

func contains(errs []string, sub string) bool {
	for _, e := range errs {
		if strings.Contains(e, sub) {
			return true
		}
	}
	return false
}
