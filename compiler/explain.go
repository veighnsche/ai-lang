package main

// a71 (a61 item 1, first cut): `ailc explain CODE` prints the
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
		rule:    "Test expectations must construct the complete error value; a bare kind proves nothing about the payload (a12).",
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
	CodeRevisionIdentity: {
		rule:    "One revision owns one reviewed interface: a changed fingerprint at the same revision is drift, removals and additions are explicit, and only an accepted baseline counts.",
		violate: `adding a case to a variant without bumping its revision.`,
		fix:     "Restore the accepted interface, or publish a reviewed new revision and explicitly update affected pins. Never regenerate the baseline to silence the finding: generation is not acceptance.",
	},
	CodeContractMalformed: {
		rule:    "A contract must be well-formed before it can be proved: every predicate is Boolean over the arm binder and the function params, every declared outcome has exactly one arm, and ensures matches are Boolean case expressions.",
		violate: `requires 1, or an ensures arm missing for a declared emits kind.`,
		fix:     "Correct the declaration or clause structure. The diagnostic names the offending clause (found) and the shape the rule wanted (expected).",
	},
	CodeContractUnsupported: {
		rule:    "The first proof cut covers mathematical integers, Booleans, and finite records of those sorts, with total bodies free of externs, kernels, state, and recursion. Anything outside is rejected, never partially proved or assumed.",
		violate: `result.value == x * x, or a contracted body calling an extern.`,
		fix:     "Use supported logic and body structure, or defer that contract-bearing interface. The diagnostic names the exact operator, sort, or dependency (found).",
	},
	CodeContractUnverifiedDep: {
		rule:    "A contract is an obligation first: no caller may reason from a callee summary until that callee is verified in the proving run. Presence in the source, an accepted revision, and passing examples are not verification.",
		violate: `a contracted function calling another contracted function before any prover ran.`,
		fix:     "Verify the source callee and its closure in the proving run; do not treat its declaration as proof. The diagnostic names the caller and the unavailable callee.",
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
