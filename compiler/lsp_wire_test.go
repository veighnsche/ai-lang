package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// Slice 0: diagnostic transport. Published frames must carry the
// stable code and the a71 payloads; payload-free lines stay
// byte-identical. Probes first.

func publishFrame(t *testing.T, text string, diags []Diag) []map[string]any {
	t.Helper()
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	if err := publishDiagnostics(w, "file:///m.can", text, diags); err != nil {
		t.Fatalf("publish: %v", err)
	}
	raw := buf.String()
	head, body, ok := strings.Cut(raw, "\r\n\r\n")
	if !ok {
		t.Fatalf("no frame header in %q", raw)
	}
	if !strings.HasPrefix(head, "Content-Length: ") {
		t.Fatalf("bad frame header %q", head)
	}
	var frame struct {
		Params struct {
			Diagnostics []map[string]any `json:"diagnostics"`
		} `json:"params"`
	}
	if err := json.Unmarshal([]byte(body), &frame); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return frame.Params.Diagnostics
}

// TestPublishCarriesCode pins the stable code on the wire.
func TestPublishCarriesCode(t *testing.T) {
	got := publishFrame(t, "ok", []Diag{{Line: 1, Sev: "error", Msg: "bad arm", Code: "CAN4107"}})
	if len(got) != 1 || got[0]["code"] != "CAN4107" {
		t.Fatalf("wire diagnostic carries no code: %v", got)
	}
}

// TestPublishCarriesPayloads pins Expected/Found/Hint in data.
func TestPublishCarriesPayloads(t *testing.T) {
	got := publishFrame(t, "ok", []Diag{{
		Line: 1, Sev: "error", Msg: "bad arm", Code: "CAN4107",
		Expected: "one arm per outcome", Found: "two on Ok", Hint: "merge",
	}})
	data, ok := got[0]["data"].(map[string]any)
	if !ok {
		t.Fatalf("wire diagnostic carries no data: %v", got[0])
	}
	for k, want := range map[string]string{
		"expected": "one arm per outcome", "found": "two on Ok", "hint": "merge",
	} {
		if data[k] != want {
			t.Fatalf("data[%q] = %v, want %q", k, data[k], want)
		}
	}
}

// TestPublishOmitsEmpty pins backward compat: payload-free lines
// carry neither code nor data keys.
func TestPublishOmitsEmpty(t *testing.T) {
	got := publishFrame(t, "ok", []Diag{{Line: 1, Sev: "warning", Msg: "unused"}})
	if len(got) != 1 {
		t.Fatalf("expected one diagnostic, got %v", got)
	}
	for _, k := range []string{"code", "data"} {
		if _, present := got[0][k]; present {
			t.Fatalf("payload-free line carries %q: %v", k, got[0])
		}
	}
}
