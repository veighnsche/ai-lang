package main

import "testing"

// The sibling-parse suppression keys on the typed diagnostic code, never
// on message wording: rewording the human message must not detach it,
// and no other message may trigger it.
func TestNeedsSiblingWarnKeysOnCode(t *testing.T) {
	resolve := Diag{Sev: "error", Msg: "m: uses v resolves nowhere", Code: CodeUsesResolve}
	if !needsSiblingWarn([]Diag{resolve}) {
		t.Fatal("CodeUsesResolve error should arm the sibling warning")
	}
	reworded := Diag{Sev: "error", Msg: "totally different wording", Code: CodeUsesResolve}
	if !needsSiblingWarn([]Diag{reworded}) {
		t.Fatal("reworded CodeUsesResolve error should still arm the sibling warning")
	}
	impostor := Diag{Sev: "error", Msg: "x resolves nowhere", Code: CodeUnknownCall}
	if needsSiblingWarn([]Diag{impostor}) {
		t.Fatal("message text alone must not arm the sibling warning")
	}
	soft := Diag{Sev: "warning", Msg: "m: uses v resolves nowhere", Code: CodeUsesResolve}
	if needsSiblingWarn([]Diag{soft}) {
		t.Fatal("non-error CodeUsesResolve must not arm the sibling warning")
	}
	if needsSiblingWarn(nil) {
		t.Fatal("empty world must not arm the sibling warning")
	}
}
