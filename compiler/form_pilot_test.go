package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// a76: pilot consumer. form-state-to-message proves a73–a75 compose:
// construction, elimination, cross-module pins and values, linked
// execution, and emit agree on one program. No new language surface:
// if these probes need any, that is an a73–a75 bug, not scope.

const pilotFormSrc = `mod form
  provides [form__message, Form__Out, Form__State, Form__Draft]
  uses []
  emits []

variant Form__State rev 1 (
  case Empty()
  case Editing(draft: Form__Draft)
  case Submitted(name: str)
)

type Form__Draft rev 1 (
  name: str
)

type Form__Out rev 1 (
  message: str
)

fn form__message(state: Form__State) -> Form__Out rev 1
  emits []
  tests
    empty(Form__Empty()) => Ok(message = "Start typing")
    editing(Form__Editing(draft = Form__Draft(name = "Al"))) => Ok(message = "Draft: Al")
    submitted(Form__Submitted(name = "Bo")) => Ok(message = "Hello, Bo")
  match state
    on Form__Empty _ => Ok(message = "Start typing")
    on Form__Editing d => Ok(message = "Draft: " + d.draft.name)
    on Form__Submitted s => Ok(message = "Hello, " + s.name)
`

const pilotShellSrc = `mod shell
  provides [shell__greet, Shell__Out]
  uses [form__message@1, Form__Out@1, Form__State@1, Form__Draft@1]
  emits []

type Shell__Out rev 1 (
  message: str
  echo: Form__State
)

fn shell__greet(state: Form__State) -> Shell__Out rev 1
  emits []
  tests
    e(Form__Empty()) => Ok(message = "Start typing", echo = Form__Empty())
    s(Form__Submitted(name = "Bo")) => Ok(message = "Hello, Bo", echo = Form__Submitted(name = "Bo"))
  match call form__message(state)
    given
      e => [exchange args (state = Form__Empty()) outcome Ok(message = "Start typing")]
      s => [exchange args (state = Form__Submitted(name = "Bo")) outcome Ok(message = "Hello, Bo")]
    on Ok v => Ok(message = v.message, echo = state)
`

func pilotFiles() map[string]string {
	return map[string]string{"form.can": pilotFormSrc, "shell.can": pilotShellSrc}
}

// TestPilotDiagnoseClean pins the composition: provider and
// consumer diagnose clean together, including variant-carrying
// exchange args and a passthrough union field.
func TestPilotDiagnoseClean(t *testing.T) {
	dir := writeLSPDir(t, pilotFiles())
	if diags := diagnose(dir, "form.can", pilotFormSrc); len(diags) != 0 {
		t.Fatalf("form: expected no diagnostics, got %v", diags)
	}
	if diags := diagnose(dir, "shell.can", pilotShellSrc); len(diags) != 0 {
		t.Fatalf("shell: expected no diagnostics, got %v", diags)
	}
}

// TestPilotLinked pins linked execution: the consumer runs the real
// provider over a nullary and a payload case, no scripts.
func TestPilotLinked(t *testing.T) {
	if err := runLinkedPure(t, pilotFiles(),
		[]string{"form.can", "shell.can"},
		"shell__greet", 1,
		map[string]string{"state": `Form__Empty()`},
		`Ok(message = "Start typing", echo = Form__Empty())`); err != nil {
		t.Fatalf("linked empty: %v", err)
	}
	if err := runLinkedPure(t, pilotFiles(),
		[]string{"form.can", "shell.can"},
		"shell__greet", 1,
		map[string]string{"state": `Form__Submitted(name = "Bo")`},
		`Ok(message = "Hello, Bo", echo = Form__Submitted(name = "Bo"))`); err != nil {
		t.Fatalf("linked submitted: %v", err)
	}
}

// TestPilotLinkedContradiction pins the CAN3110 shape in linkage: a
// wrong message fails, and a wrong-tag echo fails on tag identity.
func TestPilotLinkedContradiction(t *testing.T) {
	err := runLinkedPure(t, pilotFiles(),
		[]string{"form.can", "shell.can"},
		"shell__greet", 1,
		map[string]string{"state": `Form__Empty()`},
		`Ok(message = "WRONG", echo = Form__Empty())`)
	if err == nil || !strings.Contains(err.Error(), "payload mismatch") {
		t.Fatalf("expected payload mismatch, got %v", err)
	}
	err = runLinkedPure(t, pilotFiles(),
		[]string{"form.can", "shell.can"},
		"shell__greet", 1,
		map[string]string{"state": `Form__Empty()`},
		`Ok(message = "Start typing", echo = Form__Submitted(name = "Bo"))`)
	if err == nil || !strings.Contains(err.Error(), "payload mismatch") {
		t.Fatalf("expected tag-identity mismatch, got %v", err)
	}
}

// TestPilotCatalogAbsent pins catalog silence for the composed
// program: eliminating and passing cases creates no entries.
func TestPilotCatalogAbsent(t *testing.T) {
	dir := t.TempDir()
	files := pilotFiles()
	var paths []string
	for name, src := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(src), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		paths = append(paths, filepath.Join(dir, name))
	}
	out := t.TempDir()
	if err := compile(out, paths); err != nil {
		t.Fatalf("compile: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(out, "errors.json"))
	if err != nil {
		t.Fatalf("read errors.json: %v", err)
	}
	if strings.Contains(string(raw), "Form") {
		t.Fatalf("pilot leaked into errors catalog: %s", raw)
	}
}

// TestGoldenFormState freezes the pilot artifacts: the committed
// sketches must transpile byte-identical, so the first real
// variant consumer can never silently rot.
func TestGoldenFormState(t *testing.T) {
	dir := t.TempDir()
	srcs := []string{
		"../sketches/form-state/form.can",
		"../sketches/form-state/shell.can",
	}
	if err := compile(dir, srcs); err != nil {
		t.Fatalf("compile: %v", err)
	}
	for _, f := range []string{"form.ts", "shell.ts", "errors.json"} {
		got, err := os.ReadFile(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("read fresh %s: %v", f, err)
		}
		want, err := os.ReadFile(filepath.Join("../sketches/form-state", f))
		if err != nil {
			t.Fatalf("read golden %s: %v", f, err)
		}
		checkGoldenFile(t, f, got, want)
	}
}
