package main

// v71 (v61 item 1, first cut): `ailc explain CODE` prints the
// per-code doc — the rule, a minimal violation, the legal
// fix. Four payload codes carry full entries; every other
// registered code falls back to family text that admits
// the cut. Docs are versioned with the compiler; codes
// stay frozen.

import (
	"fmt"
	"io"
	"strings"
)

type explainEntry struct {
	rule    string
	violate string
	fix     string
}

var explainDocs = map[string]explainEntry{
	CodeForeignRaise: {
		rule:    "A function may only raise error kinds in its own emits list (R5).",
		violate: `fn m__go() emits [m.known], body raises m.stray().`,
		fix:     "Add the kind to emits, or raise a declared kind. The diagnostic names the raised kind (found) and the declared set (expected).",
	},
	CodeBareErrorKind: {
		rule:    "Test expectations must construct the complete error value; a bare kind proves nothing about the payload (v12).",
		violate: `go => m.bad  (where error m.bad(id: str))`,
		fix:     "Write m.bad(id = ...) with all fields. The diagnostic shows the complete construction (expected).",
	},
	CodeDanglingTest: {
		rule:    "Every reaching test needs a script at every foreign call (R8).",
		violate: `a test with no entry under the call's given table.`,
		fix:     "Add `test => [exchange args (...) outcome ...]` under given. The diagnostic shows the row shape (expected) and the table line (hint).",
	},
	CodeMissingArm: {
		rule:    "Call matches cover exactly the callee emits plus ok: missing or stale arms are compile errors.",
		violate: `a match on a failable call without an on <kind> arm.`,
		fix:     "Add the missing arm. The diagnostic names the kind (expected).",
	},
}

// explainFamily describes each code family for codes
// without a per-code doc yet.
var explainFamily = map[string]string{
	"AIL1": "Parse and grammar rules.",
	"AIL2": "Declaration rules (naming R3, rev pins R4, provides/uses R2).",
	"AIL3": "Call and decision-table rules (resolution, given/test R8).",
	"AIL4": "Error and proof rules (emits R5, exhaustiveness, test runs).",
	"AIL5": "World and tooling rules.",
	"AIL6": "Type rules (brands, seal, extern, numerics).",
}

// explainCode returns the doc text for a registered code.
// Per-code docs win; otherwise family text that admits the
// cut. False means the code is not registered at all.
func explainCode(code string) (string, bool) {
	registered := false
	for _, c := range allCodes {
		if c == code {
			registered = true
			break
		}
	}
	if !registered {
		return "", false
	}
	if e, ok := explainDocs[code]; ok {
		return fmt.Sprintf("%s\nrule: %s\nviolation: %s\nfix: %s", code, e.rule, e.violate, e.fix), true
	}
	fam := ""
	if len(code) >= 4 {
		fam = explainFamily[code[:4]]
	}
	if fam == "" {
		fam = "No family text."
	}
	return fmt.Sprintf("%s\n%s No per-code doc yet; see the family rule and the diagnostic message.", code, fam), true
}

// runExplain implements `ailc explain CODE`.
func runExplain(out io.Writer, argv []string) int {
	if len(argv) != 1 {
		fmt.Fprintln(out, "usage: ailc explain AILnnnn")
		return 2
	}
	code := strings.ToUpper(strings.TrimSpace(argv[0]))
	text, ok := explainCode(code)
	if !ok {
		fmt.Fprintf(out, "ailc: unknown diagnostic code %q\n", argv[0])
		return 2
	}
	fmt.Fprintln(out, text)
	return 0
}
