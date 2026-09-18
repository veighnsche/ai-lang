package main

import "testing"

// TestClassifyCallee pins the single callee-kind chain: every
// intrinsic name maps to its kind, intrinsics (and only
// intrinsics) report IsIntrinsic, and the local/foreign/unknown
// split follows file identity, not string matching.
func TestClassifyCallee(t *testing.T) {
	for fname, want := range map[string]CalleeKind{
		"state__get":            CalleeStoreOp,
		"state__put":            CalleeStoreOp,
		"dec__parts":            CalleeDecParts,
		"bytes__utf8__decode":   CalleeBytesDecode,
		"bytes__hex__decode":    CalleeBytesHexDecode,
		"bytes__base64__decode": CalleeBytesB64Decode,
		"bytes__hex__encode":    CalleeBytesHexEncode,
		"bytes__base64__encode": CalleeBytesB64Encode,
		"schema__asset__fields": CalleeBytesAsset,
		"bytes__utf8__export":   CalleeBytesOther,
		"bytes__utf8__encode":   CalleeBytesOther,
	} {
		if got := classifyCallee(nil, "", fname); got != want {
			t.Errorf("classifyCallee(%q) = %v, want %v", fname, got, want)
		}
		if !classifyCallee(nil, "", fname).IsIntrinsic() {
			t.Errorf("intrinsic %q reports not intrinsic", fname)
		}
	}

	prog := &Program{
		Fns:     map[string]*FnDecl{"m__help": {}, "m__far": {}},
		Externs: map[string]*ExternDecl{},
		FnFile:  map[string]string{"m__go": "m.can", "m__help": "m.can", "m__far": "o.can"},
	}
	if got := classifyCallee(prog, "m__go", "m__help"); got != CalleeLocal {
		t.Errorf("same-file helper = %v, want CalleeLocal", got)
	}
	if got := classifyCallee(prog, "m__go", "m__far"); got != CalleeForeign {
		t.Errorf("cross-file function = %v, want CalleeForeign", got)
	}
	if got := classifyCallee(prog, "m__go", "m__missing"); got != CalleeUnknown {
		t.Errorf("unresolvable name = %v, want CalleeUnknown", got)
	}
	for _, fname := range []string{"m__help", "m__far", "m__missing"} {
		if classifyCallee(prog, "m__go", fname).IsIntrinsic() {
			t.Errorf("non-intrinsic %q reports intrinsic", fname)
		}
	}
	// Order is load-bearing: an intrinsic name is never a local,
	// even when a same-file function shadows the spelling.
	shadow := &Program{
		Fns:    map[string]*FnDecl{"state__get": {}},
		FnFile: map[string]string{"m__go": "m.can", "state__get": "m.can"},
	}
	if got := classifyCallee(shadow, "m__go", "state__get"); got != CalleeStoreOp {
		t.Errorf("shadowed intrinsic = %v, want CalleeStoreOp", got)
	}
}
