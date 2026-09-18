// check.go: static cross-checks behind editor squiggles.
//
// diagnose already runs the compile-time proofs (parse, naming, rev pins,
// exhaustiveness, decision-table tests). The checks here catch what those
// proofs only trip over at test-execution time with confusing messages, or
// never catch at all: calls that resolve nowhere, given/test script
// mismatches, raises outside emits, and a mod header that disagrees with
// the file it heads. Every check returns []Diag with its own line so one
// broken file yields many precise squiggles instead of one vague one.
package main

import (
	"fmt"
	"math/big"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// qualifiedFile renders a module for diagnostics: the bare filename,
// qualified to its canonical input path when another module shares the
// basename. Single-basename worlds render exactly as before.
func qualifiedFile(mods []*Module, m *Module) string {
	for _, o := range mods {
		if o != m && o.File == m.File {
			return m.ID
		}
	}
	return m.File
}

// qualifiedID renders a module identity for diagnostics when only the
// identity string is at hand (cycle edges).
func qualifiedID(mods []*Module, id string) string {
	base := filepath.Base(id)
	for _, o := range mods {
		if o.ID != id && o.File == base {
			return id
		}
	}
	return base
}

// parsePinRev splits a uses entry into its requested revision:
// name@N yields N. ok=false when the entry carries no trailing @N,
// in which case the existing pin/resolve diagnostics own the entry.
// An all-digit pin outside int range yields -1: declared revs are
// int-parsed themselves, so no declaration can equal it and the pin
// still fails instead of resolving by bare name.
func parsePinRev(u string) (rev int, ok bool) {
	i := strings.LastIndex(u, "@")
	if i < 0 {
		return 0, false
	}
	if n, err := strconv.Atoi(u[i+1:]); err == nil {
		return n, true
	}
	if pinRe.MatchString(u) {
		return -1, true
	}
	return 0, false
}

// declaredRev reports the revision the owning module declares for a
// provided function, extern, record, or brand.
func declaredRev(owner *Module, base string) (int, bool) {
	for _, d := range owner.Decls {
		switch d := d.(type) {
		case *FnDecl:
			if d.Name == base {
				return d.Rev, true
			}
		case *ExternDecl:
			if d.Name == base {
				return d.Rev, true
			}
		case *TypeDecl:
			if d.Name == base {
				return d.Rev, true
			}
		case *VariantDecl:
			if d.Name == base {
				return d.Rev, true
			}
		case *BrandDecl:
			if d.Name == base {
				return d.Rev, true
			}
		case *ConstDecl:
			if d.Name == base {
				return d.Rev, true
			}
		}
	}
	return 0, false
}

// buildWorld is buildProgram for the editor: best-effort Program plus one
// diagnostic per broken world item, each on its own line. Double
// definitions point at the second declaration; bad uses entries point at
// the entry. Callers skip execution-dependent checks (exhaustiveness proof,
// test runs, call/given/emits checks) when the world is broken so one bad
// uses line does not cascade into false positives elsewhere.
func buildWorld(open *Module, mods []*Module, texts map[string]string) (*Program, []Diag) {
	var out []Diag
	prog := &Program{
		Fns: map[string]*FnDecl{}, Externs: map[string]*ExternDecl{},
		Brands:  map[string]string{},
		Errors:  map[string][]string{},
		EmitsOf: map[string][]string{}, Uses: map[string]bool{},
		Modules: mods, FnFile: map[string]string{}, BrandFile: map[string]string{},
		Variants: map[string]*VariantDecl{}, Cases: map[string]string{},
		Consts: map[string]*ConstDecl{}, ConstFile: map[string]string{},
	}
	provides := map[string]*Module{}
	emit := func(m *Module, d Diag) {
		d.File = qualifiedFile(mods, m)
		out = append(out, d)
	}
	// recordNames seeds the case-identity collision check (a73):
	// a qualified case name must not equal any record type name.
	// Builtins plus every declared record, collected up front so
	// declaration order never matters.
	recordNames := map[string]bool{}
	for _, b := range builtinTypeDecls() {
		recordNames[b.Name] = true
	}
	for _, m := range mods {
		for _, d := range m.Decls {
			if td, ok := d.(*TypeDecl); ok {
				recordNames[td.Name] = true
			}
		}
	}
	seenFn := map[string]*Module{}
	seenOther := map[string]bool{}
	seenVariant := map[string]*Module{}
	seenConst := map[string]*Module{}
	for _, m := range mods {
		for _, d := range m.Decls {
			name, line := declNameLine(d)
			if name == "" {
				continue
			}
			// Only functions collide loudly: the CLI errors on double
			// function definitions and silently overwrites shadowed
			// types and errors. The LSP mirrors that, so a gallery of
			// independent demos sharing one directory (and one error
			// name) still diagnoses per file instead of drowning in
			// cross-file noise. First declaration wins, deterministically.
			// A clash inside the open file points at its declaration; a
			// clash purely between siblings points at line 1 and names
			// the file, since the editor can only squiggle the open doc.
			// Externs share the function namespace: a foreign import
			// shadowing an can function (or vice versa) is the same
			// double definition, or calls could not resolve.
			_, isFn := d.(*FnDecl)
			_, isEx := d.(*ExternDecl)
			_, isVar := d.(*VariantDecl)
			_, isConst := d.(*ConstDecl)
			if isVar {
				// Variants collide loudly (decision 8): the case
				// set affects proofs, so silent first-wins would
				// prove consumers against the wrong cases.
				if _, dup := seenVariant[name]; dup {
					if m == open {
						emit(m, spanDiag(texts[m.ID], line, "error",
							fmt.Sprintf("double definition: variant %s", name), name, CodeDupVariant))
					} else {
						emit(m, spanDiag(texts[open.ID], 1, "error",
							fmt.Sprintf("sibling %s also declares variant %s (double definition)", qualifiedFile(mods, m), name), m.File, CodeDupVariant))
					}
					continue
				}
				seenVariant[name] = m
			} else if isFn || isEx {
				if isBytesKernel(name) {
					kind := "function"
					if isEx {
						kind = "extern"
					}
					emit(m, spanDiag(texts[m.ID], line, "error",
						fmt.Sprintf("%s %s shadows a compiler kernel: rename the declaration", kind, name), name, CodePrimitiveShadow))
					continue
				}
				if _, dup := seenFn[name]; dup {
					if m == open {
						emit(m, spanDiag(texts[m.ID], line, "error",
							fmt.Sprintf("double definition: %s", name), name, CodeDupFn))
					} else {
						emit(m, spanDiag(texts[open.ID], 1, "error",
							fmt.Sprintf("sibling %s also defines %s (double definition)", qualifiedFile(mods, m), name), m.File, CodeDupSibling))
					}
					continue
				}
				seenFn[name] = m
			} else if isConst {
				// Constants collide loudly like functions: a
				// shared name must resolve to exactly one
				// value, so first-wins would silently prove
				// consumers against the wrong literal.
				if _, dup := seenConst[name]; dup {
					if m == open {
						emit(m, spanDiag(texts[m.ID], line, "error",
							fmt.Sprintf("double definition: const %s", name), name, CodeDupConst))
					} else {
						emit(m, spanDiag(texts[open.ID], 1, "error",
							fmt.Sprintf("sibling %s also defines const %s (double definition)", qualifiedFile(mods, m), name), m.File, CodeDupConst))
					}
					continue
				}
				seenConst[name] = m
			} else if seenOther[name] {
				continue
			} else {
				seenOther[name] = true
			}
			switch d := d.(type) {
			case *FnDecl:
				prog.Fns[d.Name] = d
				prog.EmitsOf[d.Name] = d.Emits
				prog.FnFile[d.Name] = m.ID
				provides[d.Name] = m
			case *ConstDecl:
				prog.Consts[d.Name] = d
				prog.ConstFile[d.Name] = m.ID
				provides[d.Name] = m
			case *ExternDecl:
				prog.Externs[d.Name] = d
				prog.EmitsOf[d.Name] = d.Emits
				prog.FnFile[d.Name] = m.ID
				provides[d.Name] = m
			case *TypeDecl:
				if d.Name == "Bytes" {
					emit(m, spanDiag(texts[m.ID], line, "error",
						"type Bytes shadows the Bytes primitive: rename the declaration", d.Name, CodePrimitiveShadow))
				} else if isBuiltinRecord(d.Name) {
					emit(m, spanDiag(texts[m.ID], line, "error",
						fmt.Sprintf("type %s shadows a compiler-owned record: rename the declaration", d.Name), d.Name, CodePrimitiveShadow))
				}
				provides[d.Name] = m
			case *VariantDecl:
				if d.Name == "Bytes" {
					emit(m, spanDiag(texts[m.ID], line, "error",
						"variant Bytes shadows the Bytes primitive: rename the declaration", d.Name, CodePrimitiveShadow))
				} else if isBuiltinRecord(d.Name) {
					emit(m, spanDiag(texts[m.ID], line, "error",
						fmt.Sprintf("variant %s shadows a compiler-owned record: rename the declaration", d.Name), d.Name, CodePrimitiveShadow))
				} else if recordNames[d.Name] {
					// a74: a variant parent shares its TS type
					// name with a record of the same name, so a
					// silent collision would map one name to two
					// shapes at emit. Reject like case/record
					// collisions (decision 8); §2.3 will narrow
					// the global rule.
					emit(m, spanDiag(texts[m.ID], line, "error",
						fmt.Sprintf("variant identity collision: %s collides with record type %s", d.Name, d.Name), d.Name, CodeCaseCollision))
					continue
				}
				provides[d.Name] = m
				prog.Variants[d.Name] = d
				for _, c := range d.Cases {
					q := qualifyCase(d.Name, c.Short)
					if owner, dup := prog.Cases[q]; dup {
						if owner == d.Name {
							emit(m, spanDiag(texts[m.ID], c.Line, "error",
								fmt.Sprintf("duplicate case %s in variant %s", q, d.Name), c.Short, CodeCaseCollision))
						} else {
							emit(m, spanDiag(texts[m.ID], c.Line, "error",
								fmt.Sprintf("case identity collision: %s declared by both %s and %s", q, owner, d.Name), c.Short, CodeCaseCollision))
						}
						continue
					}
					if recordNames[q] {
						emit(m, spanDiag(texts[m.ID], c.Line, "error",
							fmt.Sprintf("case identity collision: %s collides with record type %s", q, q), c.Short, CodeCaseCollision))
						continue
					}
					prog.Cases[q] = d.Name
				}
			case *BrandDecl:
				if d.Name == "Bytes" {
					emit(m, spanDiag(texts[m.ID], line, "error",
						"brand Bytes shadows the Bytes primitive: rename the declaration", d.Name, CodePrimitiveShadow))
				} else if isBuiltinRecord(d.Name) {
					emit(m, spanDiag(texts[m.ID], line, "error",
						fmt.Sprintf("brand %s shadows a compiler-owned record: rename the declaration", d.Name), d.Name, CodePrimitiveShadow))
				}
				provides[d.Name] = m
				if _, ok := prog.Brands[d.Name]; !ok {
					prog.Brands[d.Name] = d.Under
				}
				if _, ok := prog.BrandFile[d.Name]; !ok {
					prog.BrandFile[d.Name] = m.ID
				}
			case *AssetBridgeDecl:
				// Grant-named brands erase to str like every v0
				// brand; seed the erasure map first-wins so sinks
				// emit standalone. Seal authority (BrandFile) stays
				// with real declarations only, and grants never
				// enter provides.
				for _, name := range []string{d.Asset, d.Policy} {
					if _, ok := prog.Brands[name]; !ok {
						prog.Brands[name] = "str"
					}
				}
			case *ErrorDecl:
				if isBuiltinError(d.Name) {
					emit(m, spanDiag(texts[m.ID], line, "error",
						fmt.Sprintf("error %s shadows a compiler-owned error: rename the declaration", d.Name), d.Name, CodePrimitiveShadow))
					continue
				}
				var fs []string
				for _, f := range d.Fields {
					fs = append(fs, f[0])
				}
				prog.Errors[d.Name] = fs
			}
		}
	}
	for _, m := range mods {
		for _, u := range m.Hdr["uses"] {
			line := locateLine(texts[m.ID], u, 1)
			spanText := texts[m.ID]
			if m != open {
				// Another file's problem shows on the open doc only
				// as a line-1 pointer; the message names the file.
				line = 1
				spanText = texts[open.ID]
			}
			// The pin format is validated before resolution: an
			// entry like t__f@x names no revision, so it is a
			// pin error (name@N), never a "resolves nowhere".
			if !strings.Contains(u, "@") || !pinRe.MatchString(u) {
				emit(m, spanDiag(spanText, line, "error",
					fmt.Sprintf("%s: uses %s must pin a rev (name@N)", m.File, u), u, CodeUsesPin))
				continue
			}
			base := pinRe.ReplaceAllString(u, "")
			owner, ok := provides[base]
			if !ok || owner == m {
				emit(m, spanDiag(spanText, line, "error",
					fmt.Sprintf("%s: uses %s resolves nowhere", m.File, u), base, CodeUsesResolve))
				continue
			}
			// Exact revision resolution (R4): the pin names the
			// revision the provider actually declares. A bare-name
			// hit is not enough — @999, @0, or a stale @1 against a
			// rev-2 provider all fail here, before execution or
			// emission. No historical storage: unavailable revisions
			// fail explicitly.
			if want, ok := parsePinRev(u); ok {
				if got, found := declaredRev(owner, base); found && got != want {
					emit(m, spanDiag(spanText, line, "error",
						fmt.Sprintf("%s: uses %s pins rev %d, but %s declares rev %d", m.File, u, want, base, got), u, CodeUsesRev))
					continue
				}
			}
			prog.Uses[base] = true
		}
	}
	// a46 S2: every registered kernel declares its contract
	// explicitly. Entries must exist (exhaustiveness verifies them
	// independently); an absent entry is never an empty error set.
	for name, k := range bytesKernels {
		prog.EmitsOf[name] = k.emits
	}
	// a50 B6: compiler-owned errors register globally, after source
	// declarations. Source redefinitions are shadow rejections (see
	// above), never silent overwrites, so this assignment is the
	// single authoritative field list. Catalogs enumerate prog.Errors
	// and therefore list these in every build (accepted churn: the
	// all-declarations policy is unchanged).
	for _, b := range builtinErrorDecls() {
		var fs []string
		for _, f := range b.Fields {
			fs = append(fs, f[0])
		}
		prog.Errors[b.Name] = fs
	}
	// Slice 1: resolve const-named patterns to literals once the
	// const table is complete, so the prover, test runs, and emit
	// see one shape in every pipeline (the LSP proves before
	// checkSem runs, so elaboration cannot live in checkSem
	// alone). Idempotent: checkSem re-runs it harmlessly.
	for _, m := range mods {
		for _, d := range elaborateConstPatterns(m, prog, texts[m.ID]) {
			emit(m, d)
		}
	}
	// a90: elaborate forward-call arms into full relay matches
	// first: their output arms are plain forwards the next pass
	// must expand. Same terms — once, here, before proofs, runs,
	// and emit — and refused sites stay forward-shaped, so this
	// must not re-run in checkSem either.
	for _, m := range mods {
		for _, d := range expandForwardCalls(m, prog, texts[m.ID]) {
			emit(m, d)
		}
	}
	// Slice 2: elaborate forward arms into complete constructors
	// on the same terms: once, here, before proofs, runs, and
	// emit. Invalid forwards stay forward-shaped, so this must
	// not re-run in checkSem (their diagnostics would double).
	for _, m := range mods {
		for _, d := range elaborateForwards(m, prog, texts[m.ID]) {
			emit(m, d)
		}
	}
	// a86: elaborate chain blocks into nested call matches on the
	// same terms: once, here, with callee contracts complete and
	// before any proof, run, or emit sees the function. Chains
	// vanish in place, so a second run is a no-op (and checkSem
	// must not re-run it for the same doubling reason).
	for _, m := range mods {
		for _, d := range elaborateChains(m, prog, texts[m.ID]) {
			emit(m, d)
		}
	}
	return prog, out
}

func declNameLine(d Decl) (string, int) {
	switch d := d.(type) {
	case *FnDecl:
		return d.Name, d.Line
	case *ExternDecl:
		return d.Name, d.Line
	case *TypeDecl:
		return d.Name, d.Line
	case *VariantDecl:
		return d.Name, d.Line
	case *BrandDecl:
		return d.Name, d.Line
	case *ConstDecl:
		return d.Name, d.Line
	case *AssetBridgeDecl:
		// Grants register under the sink function name (unique per
		// grant): the asset brand name would collide across the
		// stylesheet/script grants in seenOther tracking.
		return d.Function, d.Line
	case *ErrorDecl:
		return d.Name, d.Line
	}
	return "", 0
}

// headerLine finds the mod-header keyword line (provides/uses/emits/mod).
func headerLine(text, kw string, fallback int) int {
	for n, line := range strings.Split(text, "\n") {
		if t := strings.TrimSpace(line); strings.HasPrefix(t, kw+" ") || strings.HasPrefix(t, kw+"[") {
			return n + 1
		}
	}
	return fallback
}

// locateUseLine finds a code occurrence of sub, skipping header, comment,
// and declaration lines so call/error sites point at use, not declaration.
func locateUseLine(text, sub string, fallback int) int {
	for n, line := range strings.Split(text, "\n") {
		t := strings.TrimSpace(line)
		if !strings.Contains(line, sub) {
			continue
		}
		if strings.HasPrefix(t, "mod ") || strings.HasPrefix(t, "provides") ||
			strings.HasPrefix(t, "uses") || strings.HasPrefix(t, "emits") ||
			strings.HasPrefix(t, "error ") || strings.HasPrefix(t, "type ") ||
			strings.HasPrefix(t, "fn ") || strings.HasPrefix(t, "tests") ||
			strings.HasPrefix(t, "//") {
			continue
		}
		return n + 1
	}
	return locateLine(text, sub, fallback)
}

// matchNodes collects every match node in a body, innermost included.
func matchNodes(n *Node) []*Node {
	var out []*Node
	var walk func(x *Node)
	walk = func(x *Node) {
		if x == nil {
			return
		}
		if x.IsMatch {
			out = append(out, x)
			for _, a := range x.Arms {
				walk(a.Rhs)
			}
		}
	}
	walk(n)
	return out
}

// walkSmallTrees visits a Small and every nested Small (args, binop sides,
// text-operator operands, list items, exchange outcomes).
func walkSmallTrees(s *Small, f func(*Small)) {
	if s == nil {
		return
	}
	f(s)
	for _, a := range s.Args {
		walkSmallTrees(a.V, f)
	}
	if s.Kind == "binop" {
		walkSmallTrees(s.L, f)
		walkSmallTrees(s.R, f)
	}
	if s.Kind == "not" {
		walkSmallTrees(s.L, f)
	}
	if s.Kind == "neg" {
		walkSmallTrees(s.L, f)
	}
	if s.Kind == "strlen" || s.Kind == "stridx" || s.Kind == "strslice" {
		walkSmallTrees(s.L, f)
		walkSmallTrees(s.R, f)
		walkSmallTrees(s.Hi, f)
	}
	for _, it := range s.Items {
		walkSmallTrees(it, f)
	}
	if s.Kind == "exchange" {
		walkSmallTrees(s.Outcome, f)
	}
}

// bodySmalls visits every Small in executable positions exactly once:
// arm right-hand sides (at the arm's line) and plain expression bodies.
// Match scrutinees are excluded (callers handle them separately). Given
// stubs and test tables are mocks and expectations, not executed code, so
// they are excluded too.
func bodySmalls(n *Node, f func(s *Small, line int)) {
	var walk func(x *Node)
	walk = func(x *Node) {
		if x == nil {
			return
		}
		if x.IsMatch {
			for _, a := range x.Arms {
				walk(a.Rhs)
			}
			return
		}
		line := x.Line
		walkSmallTrees(x.Small, func(s *Small) { f(s, line) })
	}
	walk(n)
}

// checkModIntegrity enforces R2 from the inside: provides must name exactly
// what the file defines. Error declarations are module-scoped values, not
// provided items, so they are exempt on both sides.
func checkModIntegrity(m *Module, text string) []Diag {
	var out []Diag
	defined := map[string]bool{}
	for _, d := range m.Decls {
		switch d := d.(type) {
		case *FnDecl:
			defined[d.Name] = true
			if !hasHdr(m.Hdr["provides"], d.Name) {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("%s is defined but missing from provides", d.Name), d.Name, CodeProvidesMiss))
			}
		case *TypeDecl:
			defined[d.Name] = true
			if !hasHdr(m.Hdr["provides"], d.Name) {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("%s is defined but missing from provides", d.Name), d.Name, CodeProvidesMiss))
			}
		case *VariantDecl:
			defined[d.Name] = true
			if !hasHdr(m.Hdr["provides"], d.Name) {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("%s is defined but missing from provides", d.Name), d.Name, CodeProvidesMiss))
			}
		case *BrandDecl:
			defined[d.Name] = true
			if !hasHdr(m.Hdr["provides"], d.Name) {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("%s is defined but missing from provides", d.Name), d.Name, CodeProvidesMiss))
			}
		case *ConstDecl:
			defined[d.Name] = true
			if !hasHdr(m.Hdr["provides"], d.Name) {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("%s is defined but missing from provides", d.Name), d.Name, CodeProvidesMiss))
			}
		case *ExternDecl:
			defined[d.Name] = true
			if !hasHdr(m.Hdr["provides"], d.Name) {
				out = append(out, spanDiag(text, d.Line, "error",
					fmt.Sprintf("%s is defined but missing from provides", d.Name), d.Name, CodeProvidesMiss))
			}
		}
	}
	for _, p := range m.Hdr["provides"] {
		if !defined[p] {
			out = append(out, spanDiag(text, headerLine(text, "provides", 1), "error",
				fmt.Sprintf("provides %s but the file never defines it", p), p, CodeProvidesGhost))
		}
	}
	return out
}

func hasHdr(list []string, name string) bool {
	for _, e := range list {
		if e == name {
			return true
		}
	}
	return false
}

// resolveTestArgs gives positional test args their parameter names: a
// positional arg at list index i takes params[i], mirroring bindSlots'
// slot rule for calls. Named args keep their names. Runs in checkStatic
// before checkTestShapes, so every later phase — shapes, types, runs,
// emit, coverage, fingerprints — sees one named shape. Idempotent:
// resolved args are named, so a second run is a no-op. Three faults
// are CAN3205: a positional after a named arg, a positional past the
// arity, and a positional landing on a named-claimed slot. Unknown
// named args stay for CAN3202, missing params for CAN3203.
func resolveTestArgs(fn *FnDecl, text string) []Diag {
	var out []Diag
	for ti := range fn.Tests {
		t := &fn.Tests[ti]
		claimed := map[int]bool{}
		for _, a := range t.Args {
			if !a.HasName {
				continue
			}
			for j, p := range fn.Params {
				if p[0] == a.Name {
					claimed[j] = true
				}
			}
		}
		seenNamed := false
		for i := range t.Args {
			a := &t.Args[i]
			if a.HasName {
				seenNamed = true
				continue
			}
			if seenNamed {
				out = append(out, spanDiag(text, t.Line, "error",
					fmt.Sprintf("test %s takes positional arg after named arg", t.Name), t.Name, CodeBadPositional))
				continue
			}
			if i >= len(fn.Params) {
				out = append(out, spanDiag(text, t.Line, "error",
					fmt.Sprintf("test %s takes %d args for %d params", t.Name, len(t.Args), len(fn.Params)), t.Name, CodeBadPositional))
				continue
			}
			if claimed[i] {
				out = append(out, spanDiag(text, t.Line, "error",
					fmt.Sprintf("test %s supplies arg %s twice", t.Name, fn.Params[i][0]), t.Name, CodeBadPositional))
				continue
			}
			claimed[i] = true
			a.Name = fn.Params[i][0]
			a.HasName = true
		}
	}
	return out
}

// checkTestShapes catches decision-table mistakes without running anything:
// duplicate test names, args that match no parameter, and parameters the
// test never supplies (all three are runtime failures today). Unnamed
// args are resolution-owned faults, already reported as CAN3205.
func checkTestShapes(fn *FnDecl, text string) []Diag {
	var out []Diag
	seen := map[string]bool{}
	params := map[string]bool{}
	for _, p := range fn.Params {
		params[p[0]] = true
	}
	for _, t := range fn.Tests {
		if seen[t.Name] {
			out = append(out, spanDiag(text, t.Line, "error",
				fmt.Sprintf("duplicate test %s in %s", t.Name, fn.Name), t.Name, CodeDupTest))
		}
		seen[t.Name] = true
		got := map[string]bool{}
		for _, a := range t.Args {
			if !a.HasName {
				continue
			}
			got[a.Name] = true
			if !params[a.Name] {
				out = append(out, spanDiag(text, t.Line, "error",
					fmt.Sprintf("test %s passes unknown arg %s", t.Name, a.Name), a.Name, CodeUnknownArg))
			}
		}
		for _, p := range fn.Params {
			if !got[p[0]] {
				out = append(out, spanDiag(text, t.Line, "error",
					fmt.Sprintf("test %s is missing arg %s", t.Name, p[0]), t.Name, CodeMissingArg))
			}
		}
	}
	return out
}

// checkCalls mirrors the evaluator's call rules statically. Foreign
// scrutinees must name a known function in uses; same-file helpers
// need no pin (locality is visible, and a same-file uses entry
// resolves nowhere). Module-local externs are the exception: a
// foreign import is declared where it is called and needs no uses
// pin, but an extern from another module is unknown here (declare
// your own). Calls outside a match scrutinee are outside the v0
// subset entirely.
func checkCalls(fn *FnDecl, prog *Program, localExtern map[string]bool, text string) []Diag {
	var out []Diag
	scrut := map[*Small]bool{}
	for _, m := range matchNodes(fn.Body) {
		if m.Kind == MatchCall {
			s := m.Scruts[0]
			scrut[s] = true
			fname := s.Fname
			if classifyCallee(prog, fn.Name, fname).IsIntrinsic() {
				continue // cells resolve in checkEffects; kernels
				// need nothing; uses never applies to any of them
			}
			if _, ok := prog.Fns[fname]; !ok {
				if localExtern[fname] {
					continue
				}
				if prog.Externs[fname] != nil {
					out = append(out, spanDiag(text, m.Line, "error",
						fmt.Sprintf("%s calls extern %s from another module: declare your own extern", fn.Name, fname), fname, CodeUnknownCall))
					continue
				}
				out = append(out, spanDiag(text, m.Line, "error",
					fmt.Sprintf("%s calls unknown function %s", fn.Name, fname), fname, CodeUnknownCall))
				continue
			}
			if localCallee(prog, fn.Name, fname) != nil {
				continue
			}
			if !prog.Uses[fname] {
				out = append(out, spanDiag(text, m.Line, "error",
					fmt.Sprintf("%s calls %s which is not in uses: add name@rev to uses", fn.Name, fname), fname, CodeCallNotInUses))
			}
		}
	}
	bodySmalls(fn.Body, func(s *Small, line int) {
		_ = line
		if s.Kind == "call" && !scrut[s] {
			out = append(out, spanDiag(text, locateUseLine(text, s.Fname+"(", 1), "error",
				fmt.Sprintf("call to %s outside a match scrutinee is outside the v0 subset", s.Fname), s.Fname, CodeCallOutside))
		}
	})
	for _, m := range matchNodes(fn.Body) {
		if m.Kind == MatchCall {
			continue
		}
		for _, ms := range m.Scruts {
			if ms.Kind == "call" {
				continue // owned by the multi-match shape rule
			}
			walkSmallTrees(ms, func(s *Small) {
				if s.Kind == "call" {
					out = append(out, spanDiag(text, m.Line, "error",
						fmt.Sprintf("call to %s outside a match scrutinee is outside the v0 subset", s.Fname), s.Fname, CodeCallOutside))
				}
			})
		}
	}
	for _, m := range matchNodes(fn.Body) {
		if m.Kind != MatchCall {
			continue
		}
		ms := m.Scruts[0]
		walkSmallTrees(ms, func(s *Small) {
			if s.Kind == "call" && s != ms {
				out = append(out, spanDiag(text, m.Line, "error",
					fmt.Sprintf("nested call to %s inside a scrutinee is outside the v0 subset", s.Fname), s.Fname, CodeCallNested))
			}
		})
	}
	return out
}

// calleeUnknown reports whether fname resolves to nothing the
// checker can reason about: no function body, no extern
// declaration in any module, and no builtin op. checkCalls owns
// the CAN3001 for these; given, proof, and execution stay
// silent downstream so one root cause yields one diagnostic.
// Declared externs (module-local or otherwise) keep full
// checking: their emits make given/proof meaningful.
func calleeUnknown(prog *Program, fname string) bool {
	if _, ok := prog.Fns[fname]; ok {
		return false
	}
	if prog.Externs[fname] != nil {
		return false
	}
	if isStoreOp(fname) || isDecParts(fname) || isBytesKernel(fname) {
		return false
	}
	return true
}

// UnknownCallError marks a row failure caused solely by calling
// a function that resolves nowhere. checkSem suppresses the
// CAN4200 for it (checkCalls already reported the CAN3001) but
// still marks the function failed so coverage stays silent.
type UnknownCallError struct {
	Owner string
	Fname string
}

func (e *UnknownCallError) Error() string {
	return fmt.Sprintf("%s: %s not in uses", e.Owner, e.Fname)
}

// checkEagerScrutinee warns where a multi-scrutinee value match can
// fail before dispatch: every scrutinee evaluates eagerly, left to
// right, so a trapping operation (string index or slice, division or
// remainder, any call) in any scrutinee runs even when an outer arm
// of the equivalent nested match would have guarded it. Arity 1 is
// silent: one scrutinee evaluates exactly as the nested form does.
func checkEagerScrutinee(fn *FnDecl, text string) []Diag {
	var out []Diag
	for _, m := range matchNodes(fn.Body) {
		if m.Kind != MatchValue || len(m.Scruts) < 2 {
			continue
		}
		for i, ms := range m.Scruts {
			trap := firstTrappingOp(ms)
			if trap == "" {
				continue
			}
			out = append(out, spanDiag(text, m.Line, "warning",
				fmt.Sprintf("scrutinee %d may fail (%s) before dispatch: a multi-scrutinee match evaluates every scrutinee eagerly, unlike a nested match", i+1, trap), "match", CodeEagerScrutinee))
		}
	}
	return out
}

// isNonzeroIntLit reports a statically nonzero int divisor: arbitrarily
// large literals are fine (ints are unbounded), only zero can fail.
func isNonzeroIntLit(s *Small) bool {
	return s != nil && s.Kind == "int" && s.Num != nil && s.Num.Sign() != 0
}

// firstTrappingOp names the first operation in s that can fail at
// runtime, or "" when s is total. The set mirrors the evSmall failure
// surface: partial accessors, partial arithmetic, and calls.
func firstTrappingOp(s *Small) string {
	found := ""
	walkSmallTrees(s, func(n *Small) {
		if found != "" {
			return
		}
		switch n.Kind {
		case "stridx":
			found = "string index"
		case "strslice":
			found = "string slice"
		case "binop":
			if n.Op == "/" || n.Op == "%" {
				// A literal nonzero divisor cannot fail; parity-style
				// tables over n % 2 stay silent.
				if isNonzeroIntLit(n.R) {
					break
				}
				if n.Op == "/" {
					found = "division"
				} else {
					found = "remainder"
				}
			}
		case "call":
			found = fmt.Sprintf("call to %s", n.Fname)
		}
	})
	return found
}

// fileTests names every decision-table test in fname's file. Given
// keys are valid when they name a test in the file (a helper's inner
// table scripts caller tests); unknown names are the warning. Missing
// rows claim non-reach (a91): reaching the call without a script
// fails the test at execution instead of a static error.
func fileTests(prog *Program, fname string) map[string]bool {
	out := map[string]bool{}
	file, ok := prog.FnFile[fname]
	if !ok {
		return out
	}
	for _, m := range prog.Modules {
		if m.ID != file {
			continue
		}
		for _, d := range m.Decls {
			if fn, ok := d.(*FnDecl); ok {
				for _, t := range fn.Tests {
					out[t.Name] = true
				}
			}
		}
	}
	return out
}

// checkGiven cross-checks call-site evidence against the decision table:
// every foreign call needs a given table, every stub must be Ok or an
// error the callee can actually produce, scripts no test selects are
// dead, and the retired `-` spelling is an error. Missing rows claim
// non-reach (a91): reaching the call without a script fails the test
// at execution. Local helper calls are deterministic: a given table
// on one is an error, and none is required.
func checkGiven(fn *FnDecl, prog *Program, text string) []Diag {
	var out []Diag
	valid := fileTests(prog, fn.Name)
	for _, m := range matchNodes(fn.Body) {
		if m.Kind == MatchValue && m.Given != nil {
			out = append(out, spanDiag(text, m.Line, "error",
				fmt.Sprintf("multi-scrutinee match takes no given table: only a single call match takes given"), "given", CodeGivenOnLocal))
		}
		if m.Kind != MatchCall {
			continue
		}
		ms := m.Scruts[0]
		fname := ms.Fname
		if calleeUnknown(prog, fname) {
			continue // checkCalls owns the unknown-callee error
		}
		// Intrinsics and local helpers are deterministic: one
		// branch for the whole deterministic family, not one
		// copy per callee kind.
		if kind := classifyCallee(prog, fn.Name, fname); kind.IsIntrinsic() || kind == CalleeLocal {
			if m.Given != nil {
				what := fmt.Sprintf("call to %s takes no given table: it is deterministic", fname)
				if kind == CalleeLocal {
					what = fmt.Sprintf("call to local helper %s takes no given table: it is deterministic", fname)
				}
				out = append(out, spanDiag(text, m.Line, "error", what, fname, CodeGivenOnLocal))
			}
			continue
		}
		allowed := map[string]bool{"ok": true}
		for _, e := range prog.EmitsOf[fname] {
			allowed[e] = true
		}
		if m.Given == nil {
			out = append(out, spanDiag(text, m.Line, "error",
				fmt.Sprintf("call to %s has no given table: no test can script it", fname), fname, CodeNoGiven))
			continue
		}
		for key, sm := range m.Given {
			// The entry lives on its own `key => ...` row, below the
			// match: point there, not at the match line.
			entryLine := locateLineFrom(text, key+" =>", m.Line, m.Line)
			if sm == nil {
				// a91: the retired `-` spelling. Omission
				// claims non-reach; a dash row is noise.
				d := spanDiag(text, entryLine, "error",
					fmt.Sprintf("given row %s uses the retired `-` spelling: delete the row", key), key, CodeGivenDashRetired)
				d.Found = key + " => -"
				d.Expected = "(no row: omission claims non-reach)"
				d.Hint = "delete this row; script an exchange only if the test reaches the call"
				out = append(out, d)
				continue
			}
			checkStub(sm, fname, allowed, text, entryLine, key, fn.Name, &out)
			if !valid[key] {
				out = append(out, spanDiag(text, entryLine, "warning",
					fmt.Sprintf("script %s never runs: no test named %s in %s", key, key, fn.Name), key, CodeDeadScript))
			}
		}
	}
	return out
}

// checkScriptConsistency proves a18 linkage (issue #1): a scripted Ok
// outcome is a claim about what the provider computes, so when the
// provider body evaluates on the script args the claim must match.
// Scripted errors stay trusted: emits is an upper bound that admits
// unrealized entries (a12), and failure injection is what scripts are
// for. Anything the sandbox cannot evaluate (a foreign call needing
// its own script, unreachable state, depth exhaustion, uncomparable
// values) is trusted, never failed: a proof attempt that cannot run
// proves nothing, and a check must never fail a build it cannot model.
func checkScriptConsistency(fn *FnDecl, prog *Program, text string) []Diag {
	var out []Diag
	here := prog.FnFile[fn.Name]
	for _, m := range matchNodes(fn.Body) {
		if m.Kind != MatchCall {
			continue
		}
		ms := m.Scruts[0]
		fname := ms.Fname
		if isStoreOp(fname) || isDecParts(fname) || localCallee(prog, fn.Name, fname) != nil {
			continue
		}
		callee, ok := prog.Fns[fname]
		if !ok || prog.FnFile[fname] == here {
			continue // extern, unknown, or local: owned elsewhere
		}
		if m.Given == nil {
			continue // checkGiven owns the missing table
		}
		for key, entry := range m.Given {
			if entry == nil {
				continue
			}
			items := []*Small{entry}
			if entry.Kind == "list" {
				items = entry.Items
			}
			for _, item := range items {
				if item == nil || item.Kind != "exchange" || item.Outcome == nil {
					continue // outcome-only rows: CodeNoExchange owns
				}
				if item.Outcome.Kind != "ctor" || item.Outcome.Ctor != "Ok" {
					continue // errors are failure injection: trusted
				}
				computed, scripted, bad := contradictScriptOk(callee, item, key, prog)
				if !bad {
					continue
				}
				entryLine := locateLineFrom(text, key+" =>", m.Line, m.Line)
				out = append(out, spanDiag(text, entryLine, "error",
					fmt.Sprintf("script %s contradicts %s: provider computes %s, script says %s", key, fname, computed, scripted), key, CodeInconsistentScript))
			}
		}
	}
	return out
}

// contradictScriptOk evaluates the provider body on the script's own
// exchange args and compares the result with the scripted Ok outcome.
// It reports false (trusted) whenever the provider cannot be modeled:
// unbound args, a param the script never supplies, any evaluation
// error, or values vEq cannot compare. A panic inside the evaluator
// is a modeling failure too, never a build failure: the sandbox
// recovers and trusts.
func contradictScriptOk(callee *FnDecl, item *Small, test string, prog *Program) (computed, scripted string, bad bool) {
	defer func() {
		if recover() != nil {
			computed, scripted, bad = "", "", false
		}
	}()
	sandbox, err := freshExecCtx(prog, test, false)
	if err != nil {
		return "", "", false
	}
	// a92: the callee's file may never see the checker
	// (diagnose visits the open file only), so its body and
	// scripts bind here, bind-only — the same mutation
	// checkTypes owns. Faults stay with the callee's own
	// diagnosis; the sandbox evaluates what binding
	// resolved. Idempotent: checked trees keep their names.
	bc := newTycker(prog, "", callee.Name)
	benv := map[string]string{}
	for _, p := range callee.Params {
		benv[p[0]] = p[1]
	}
	bret := ""
	if bc.knownType(callee.Ret) {
		bret = callee.Ret
	}
	bc.exec = true
	bc.node(callee.Body, benv, bret)
	bc.checkStubs(callee, "", benv)
	argVals := map[string]*Value{}
	for _, a := range item.Args {
		v, err := evSmall(a.V, map[string]*Value{}, sandbox, callee.Name)
		if err != nil {
			return "", "", false
		}
		argVals[a.Name] = v
	}
	env := map[string]*Value{}
	for _, p := range callee.Params {
		v, ok := argVals[p[0]]
		if !ok {
			return "", "", false
		}
		env[p[0]] = v
	}
	got, err := evNode(callee.Body, env, sandbox, callee.Name)
	if err != nil {
		return "", "", false
	}
	want, err := evSmall(item.Outcome, env, sandbox, callee.Name)
	if err != nil {
		return "", "", false
	}
	eq, err := vEq(got, want)
	if err != nil || eq {
		return "", "", false
	}
	return normalizeValue(got), normalizeValue(want), true
}

// checkLocalCycles bans recursive helpers: compile-time tests run
// helper bodies inline, so a local call cycle would hang the build
// before any proof runs. Edges cover every same-file call,
// scrutinee or not; illegal positions are reported elsewhere,
// termination here.
func checkLocalCycles(m *Module, prog *Program, text string) []Diag {
	var out []Diag
	type edge struct {
		to    string
		token string
		line  int
	}
	edges := map[string][]edge{}
	add := func(from, to, token string, line int) {
		for _, e := range edges[from] {
			if e.to == to {
				return
			}
		}
		edges[from] = append(edges[from], edge{to, token, line})
	}
	for _, d := range m.Decls {
		fn, ok := d.(*FnDecl)
		if !ok {
			continue
		}
		for _, n := range matchNodes(fn.Body) {
			if n.Kind != MatchCall {
				continue
			}
			ns := n.Scruts[0]
			if localCallee(prog, fn.Name, ns.Fname) == nil {
				continue
			}
			if ns.Fname == fn.Name && len(fn.DecNames) > 0 {
				continue // proven self-recursion: checkDecreases owns it
			}
			add(fn.Name, ns.Fname, ns.Fname, n.Line)
		}
		bodySmalls(fn.Body, func(s *Small, line int) {
			if s.Kind == "call" && localCallee(prog, fn.Name, s.Fname) != nil {
				if s.Fname == fn.Name && len(fn.DecNames) > 0 {
					return // proven self-recursion: checkDecreases owns it
				}
				add(fn.Name, s.Fname, s.Fname, locateUseLine(text, s.Fname+"(", 1))
			}
		})
	}
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := map[string]int{}
	var stack []string
	var visit func(u string)
	visit = func(u string) {
		color[u] = gray
		stack = append(stack, u)
		for _, e := range edges[u] {
			switch color[e.to] {
			case gray:
				start := 0
				for i, n := range stack {
					if n == e.to {
						start = i
						break
					}
				}
				path := append(append([]string{}, stack[start:]...), e.to)
				out = append(out, spanDiag(text, e.line, "error",
					fmt.Sprintf("local call cycle %s: helpers must be acyclic", strings.Join(path, " -> ")), e.token, CodeLocalCycle))
			case white:
				visit(e.to)
			}
		}
		stack = stack[:len(stack)-1]
		color[u] = black
	}
	for _, d := range m.Decls {
		if fn, ok := d.(*FnDecl); ok && color[fn.Name] == white {
			visit(fn.Name)
		}
	}
	return out
}

// checkRecordCycles rejects cycles in the record-type graph for this
// first cut: only finite monomorphic products are admitted, so a
// record reaching itself through field types fails before tests or
// output. Brands and base types are leaves; unknown field types belong
// to the declaration check, not the cycle hunt.
func checkRecordCycles(mods []*Module, texts map[string]string) []Diag {
	fieldsOf := map[string][][2]string{}
	declLine := map[string]int{}
	declFile := map[string]string{}
	isVariant := map[string]bool{}
	for _, m := range mods {
		for _, d := range m.Decls {
			switch d := d.(type) {
			case *TypeDecl:
				if _, seen := fieldsOf[d.Name]; !seen {
					fieldsOf[d.Name] = d.Fields
					declLine[d.Name] = d.Line
					declFile[d.Name] = m.ID
				}
			case *VariantDecl:
				// Decision 6: cycle validation traverses the
				// combined record/variant graph, so a
				// variant→record→variant loop fails even
				// though the variant never names itself.
				// A variant's edges are the union of its
				// case payload types.
				if _, seen := fieldsOf[d.Name]; !seen {
					var fs [][2]string
					for _, c := range d.Cases {
						fs = append(fs, c.Fields...)
					}
					fieldsOf[d.Name] = fs
					declLine[d.Name] = d.Line
					declFile[d.Name] = m.ID
					isVariant[d.Name] = true
				}
			}
		}
	}
	isRecord := func(t string) bool {
		_, ok := fieldsOf[t]
		return ok
	}
	var out []Diag
	reported := map[string]bool{}
	var visit func(name string, stack []string)
	visit = func(name string, stack []string) {
		for _, f := range fieldsOf[name] {
			ft := f[1]
			if !isRecord(ft) {
				continue
			}
			idx := -1
			for i, s := range stack {
				if s == ft {
					idx = i
					break
				}
			}
			if idx >= 0 {
				// One diagnostic per cycle: rotate to start at the
				// smallest name so every root reports the same key.
				cyc := append(append([]string{}, stack[idx:]...), ft)
				body := cyc[:len(cyc)-1]
				mi := 0
				for i, n := range body {
					if n < body[mi] {
						mi = i
					}
				}
				rot := append(append([]string{}, body[mi:]...), body[:mi]...)
				rot = append(rot, rot[0])
				key := strings.Join(rot, "->")
				if !reported[key] {
					reported[key] = true
					text := texts[declFile[rot[0]]]
					kind := "record"
					for _, n := range rot {
						if isVariant[n] {
							kind = "variant"
							break
						}
					}
					out = append(out, spanDiag(text, declLine[rot[0]], "error",
						fmt.Sprintf("%s type cycle: %s reaches itself through field types", kind, strings.Join(rot, " -> ")),
						rot[0], CodeRecordCycle))
				}
				continue
			}
			visit(ft, append(stack, ft))
		}
	}
	var names []string
	for n := range fieldsOf {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		visit(n, []string{n})
	}
	return out
}

// checkGlobalCycles bans recursion across files (a11): the sandbox
// stubs foreign calls, so a cross-file cycle passes every per-file
// check and every test, then links into an unproved recursive cycle
// in production. Every can-to-can edge counts — same-file and
// cross-file — but cycles lying entirely in one file belong to the
// local check (one mistake, one diagnostic), so only cycles touching
// two or more files report here. Externs are host code outside the
// proof and never form edges; unknown callees belong to other codes.
// A cycle is reported once, at the call site that closes it, in the
// caller's file.
func checkGlobalCycles(mods []*Module, texts map[string]string, prog *Program) []Diag {
	var out []Diag
	type edge struct {
		to    string
		token string
		line  int
		file  string
	}
	edges := map[string][]edge{}
	add := func(from, to, token string, line int, file string) {
		for _, e := range edges[from] {
			if e.to == to {
				return
			}
		}
		edges[from] = append(edges[from], edge{to, token, line, file})
	}
	ordered := append([]*Module{}, mods...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].File < ordered[j].File })
	for _, m := range ordered {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			for _, n := range matchNodes(fn.Body) {
				if n.Kind != MatchCall {
					continue
				}
				ns := n.Scruts[0]
				if _, isEx := prog.Externs[ns.Fname]; isEx {
					continue
				}
				if _, ok := prog.Fns[ns.Fname]; !ok {
					continue
				}
				add(fn.Name, ns.Fname, ns.Fname, n.Line, m.ID)
			}
			bodySmalls(fn.Body, func(s *Small, line int) {
				if s.Kind != "call" {
					return
				}
				if _, isEx := prog.Externs[s.Fname]; isEx {
					return
				}
				if _, ok := prog.Fns[s.Fname]; !ok {
					return
				}
				add(fn.Name, s.Fname, s.Fname, locateUseLine(texts[m.ID], s.Fname+"(", 1), m.ID)
			})
		}
	}
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := map[string]int{}
	filesOf := func(path []string) map[string]bool {
		seen := map[string]bool{}
		for _, f := range path {
			seen[prog.FnFile[f]] = true
		}
		return seen
	}
	var stack []string
	var visit func(u string)
	visit = func(u string) {
		color[u] = gray
		stack = append(stack, u)
		for _, e := range edges[u] {
			switch color[e.to] {
			case gray:
				start := 0
				for i, n := range stack {
					if n == e.to {
						start = i
						break
					}
				}
				path := append(append([]string{}, stack[start:]...), e.to)
				if len(filesOf(path)) < 2 {
					continue // same-file cycle: the local check owns it
				}
				d := spanDiag(texts[e.file], e.line, "error",
					fmt.Sprintf("call cycle %s: only direct self-recursion with decreases is admitted", strings.Join(path, " -> ")), e.token, CodeLocalCycle)
				d.File = qualifiedID(mods, e.file)
				out = append(out, d)
			case white:
				visit(e.to)
			}
		}
		stack = stack[:len(stack)-1]
		color[u] = black
	}
	var seeds []string
	for _, m := range ordered {
		for _, d := range m.Decls {
			if fn, ok := d.(*FnDecl); ok {
				seeds = append(seeds, fn.Name)
			}
		}
	}
	for _, s := range seeds {
		if color[s] == white {
			visit(s)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].File != out[j].File {
			return out[i].File < out[j].File
		}
		return out[i].Line < out[j].Line
	})
	return out
}

// isDecrease reports whether a full-arity call site passes p - 1 for
// the decreases param: a named arg wins by name, else the positional
// arg at the param's index. The unit step is syntactic on purpose
// (a11: larger steps are refused even though they terminate — one
// spelling for the loop step): the compiler sees the decrease, it
// never infers one.
func isDecrease(scrut *Small, pidx int, p string) bool {
	var v *Small
	for i, a := range scrut.Args {
		if a.HasName && a.Name == p {
			v = a.V
			break
		}
		if !a.HasName && i == pidx && v == nil {
			v = a.V
		}
	}
	if v == nil {
		return false
	}
	if v.Kind != "binop" || v.Op != "-" {
		return false
	}
	if v.L.Kind != "ref" || len(v.L.Ref) != 1 || v.L.Ref[0] != p {
		return false
	}
	return v.R.Kind == "int" && v.R.Num.Cmp(big.NewInt(1)) == 0
}

// checkDecreases proves termination for self-recursion before anything
// runs: decreases must name int params, must guard a real self-call,
// and every self-call site must take the schema's canonical step
// under its canonical guard (a11 unit loop, a19 blessed schemas), so
// the proof promises a returned outcome, not a loud fault. Sites with
// the wrong arity belong to the arity rule and are skipped here, so
// one mistake yields one error family.
func checkDecreases(fn *FnDecl, prog *Program, text string) []Diag {
	if len(fn.DecNames) == 0 {
		return nil
	}
	idx := map[string]int{}
	typ := map[string]string{}
	for i, pr := range fn.Params {
		idx[pr[0]] = i
		typ[pr[0]] = pr[1]
	}
	for _, p := range fn.DecNames {
		if _, ok := idx[p]; !ok {
			return []Diag{spanDiag(text, fn.Line, "error",
				fmt.Sprintf("%s decreases %s: no such param", fn.Name, p), p, CodeBadDecreases)}
		}
		if typ[p] != "int" {
			return []Diag{spanDiag(text, fn.Line, "error",
				fmt.Sprintf("%s decreases %s: must be an int param, got %s", fn.Name, p, typ[p]), p, CodeBadDecreases)}
		}
	}
	selfAny := false
	for _, c := range walkCalls(fn.Body) {
		if c.Fname == fn.Name && localCallee(prog, fn.Name, c.Fname) != nil {
			selfAny = true
			break
		}
	}
	if !selfAny {
		return []Diag{spanDiag(text, fn.Line, "error",
			fmt.Sprintf("%s decreases %s but never calls itself", fn.Name, strings.Join(fn.DecNames, ", ")), strings.Join(fn.DecNames, ", "), CodeStaleDecreases)}
	}
	// Each schema pairs one canonical step with one canonical guard.
	// The params never rebind, so a site reached under a false guard
	// entered with the guard's negation true, and every chain lands
	// on the base arm. The guard is load-bearing, not stylistic.
	isStep := func(m *Small) bool { return false }
	isGuard := func(s *Small) bool { return false }
	stepMsg, guardMsg := "", ""
	switch fn.DecSchema {
	case "euclid":
		a, b := fn.DecNames[0], fn.DecNames[1]
		ai, bi := idx[a], idx[b]
		isStep = func(m *Small) bool { return isEuclidStep(m, ai, bi, a, b) }
		isGuard = func(s *Small) bool { return isGuardScrut(s, b) }
		stepMsg = fmt.Sprintf("%s calls itself without a euclid step: pass (%s, %s %% %s)", fn.Name, b, a, b)
		guardMsg = fmt.Sprintf("%s recurses outside the positive branch: self-calls must sit under the false arm of %s <= 0", fn.Name, b)
	case "narrowing":
		lo, hi := fn.DecNames[0], fn.DecNames[1]
		loi, hii := idx[lo], idx[hi]
		isStep = func(m *Small) bool { return isNarrowStep(m, loi, hii, lo, hi) }
		isGuard = func(s *Small) bool { return isNarrowGuard(s, lo, hi) }
		stepMsg = fmt.Sprintf("%s calls itself without a narrowing step: pass (%s, mid) or (mid, %s) with mid (%s+%s)/2", fn.Name, lo, hi, lo, hi)
		guardMsg = fmt.Sprintf("%s recurses outside the search branch: self-calls must sit under the false arm of (%s - %s) <= 1", fn.Name, hi, lo)
	default:
		p := fn.DecNames[0]
		pidx := idx[p]
		isStep = func(m *Small) bool { return isDecrease(m, pidx, p) }
		isGuard = func(s *Small) bool { return isGuardScrut(s, p) }
		stepMsg = fmt.Sprintf("%s calls itself without decreasing %s by one: pass %s - 1", fn.Name, p, p)
		guardMsg = fmt.Sprintf("%s recurses outside the positive branch: self-calls must sit under the false arm of %s <= 0", fn.Name, p)
	}
	var out []Diag
	var walk func(n *Node, guarded bool)
	walk = func(n *Node, guarded bool) {
		if n == nil || !n.IsMatch {
			return
		}
		if n.Kind == MatchCall {
			m := n.Scruts[0]
			if m.Fname == fn.Name {
				if localCallee(prog, fn.Name, m.Fname) != nil && len(m.Args) == len(fn.Params) {
					if !isStep(m) {
						out = append(out, spanDiag(text, n.Line, "error",
							stepMsg, fn.Name, CodeNoDecrease))
					} else if !guarded {
						out = append(out, spanDiag(text, n.Line, "error",
							guardMsg, fn.Name, CodeNoGuard))
					}
				}
			}
		}
		// Arity-1 matches guard on their one scrutinee, exactly as
		// before; wider tables never carried a guard, and their arms
		// stay unguarded here too.
		var gs *Small
		if len(n.Scruts) == 1 {
			gs = n.Scruts[0]
		}
		g := isGuard(gs)
		for _, a := range n.Arms {
			ag := guarded
			if g {
				if len(a.Pats) == 1 {
					ag = a.Pats[0].Kind == "bool" && !a.Pats[0].B
				} else {
					ag = false
				}
			}
			walk(a.Rhs, ag)
		}
	}
	walk(fn.Body, false)
	return out
}

// isBareRef reports a lone param reference: exactly [name].
func isBareRef(s *Small, name string) bool {
	return s != nil && s.Kind == "ref" && len(s.Ref) == 1 && s.Ref[0] == name
}

// siteArg extracts one self-call argument: a named arg wins by name,
// else the positional arg at the param's index (nil when absent).
func siteArg(m *Small, pidx int, p string) *Small {
	var v *Small
	for i, a := range m.Args {
		if a.HasName && a.Name == p {
			return a.V
		}
		if !a.HasName && i == pidx && v == nil {
			v = a.V
		}
	}
	return v
}

// isEuclidStep reports the one Euclidean shape: the site passes
// exactly (b, a % b). Under a b <= 0 false guard b is positive, and
// Euclidean % lands the new b in [0, b), so the second component
// walks a natural chain into the base arm (a19).
func isEuclidStep(m *Small, aidx, bidx int, a, b string) bool {
	if !isBareRef(siteArg(m, aidx, a), b) {
		return false
	}
	mod := siteArg(m, bidx, b)
	if mod == nil || mod.Kind != "binop" || mod.Op != "%" {
		return false
	}
	return isBareRef(mod.L, a) && isBareRef(mod.R, b)
}

// isMid reports the one canonical midpoint: (lo + hi) / 2, exact
// Euclidean integer division (ints are unbounded, so lo + hi cannot
// overflow the proof).
func isMid(v *Small, lo, hi string) bool {
	if v == nil || v.Kind != "binop" || v.Op != "/" {
		return false
	}
	if v.R == nil || v.R.Kind != "int" || v.R.Num == nil || v.R.Num.Cmp(big.NewInt(2)) != 0 {
		return false
	}
	s := v.L
	if s == nil || s.Kind != "binop" || s.Op != "+" {
		return false
	}
	return isBareRef(s.L, lo) && isBareRef(s.R, hi)
}

// isNarrowStep reports the two binary-search shapes: (lo, mid) or
// (mid, hi). With hi - lo >= 2 the midpoint sits strictly inside,
// so the bound gap strictly shrinks every site (a19).
func isNarrowStep(m *Small, loidx, hiidx int, lo, hi string) bool {
	lv := siteArg(m, loidx, lo)
	hv := siteArg(m, hiidx, hi)
	if lv == nil || hv == nil {
		return false
	}
	if isBareRef(lv, lo) && isMid(hv, lo, hi) {
		return true
	}
	return isMid(lv, lo, hi) && isBareRef(hv, hi)
}

// isNarrowGuard recognizes the one canonical search guard:
// (hi - lo) <= 1. The search is over when at most one integer
// remains between the bounds; degenerate entries (hi <= lo) take
// the base arm immediately.
func isNarrowGuard(s *Small, lo, hi string) bool {
	if s == nil || s.Kind != "binop" || s.Op != "<=" {
		return false
	}
	if s.R.Kind != "int" || s.R.Num.Cmp(big.NewInt(1)) != 0 {
		return false
	}
	d := s.L
	if d == nil || d.Kind != "binop" || d.Op != "-" {
		return false
	}
	return isBareRef(d.L, hi) && isBareRef(d.R, lo)
}

// isGuardScrut recognizes the one canonical bound guard: p <= 0. One
// spelling per meaning extends to the proof — p >= 1, p > 0, and p == 0
// spell the same bound but only this shape admits recursion.
func isGuardScrut(s *Small, p string) bool {
	if s == nil || s.Kind != "binop" || s.Op != "<=" {
		return false
	}
	if s.L.Kind != "ref" || len(s.L.Ref) != 1 || s.L.Ref[0] != p {
		return false
	}
	return s.R.Kind == "int" && s.R.Num.Sign() == 0
}

// isStoreOp reports the two blessed store builtins: neither foreign
// (stubbed) nor local (executed), but deterministic per test.
func isStoreOp(fname string) bool {
	return fname == "state__get" || fname == "state__put"
}

// isDecParts reports the blessed decimal observation kernel: neither
// foreign (stubbed) nor local (executed), but a total deterministic
// observation of its dec operand. Like the store ops it needs no
// uses entry and takes no given table; unlike them it needs no
// effects and emits nothing (a21). Totality is load-bearing below:
// verifyExhaustive wants exactly {ok} through the absent EmitsOf
// entry, so a future emits registration here must revisit that gate.
func isDecParts(fname string) bool {
	return fname == "dec__parts"
}

// storeCellName extracts the cell a store-op scrutinee names: the
// first argument must be a bare cell name, positionally. Cells are
// authority, not data, so no other shape is accepted here.
func storeCellName(scrut *Small) (string, bool) {
	if len(scrut.Args) == 0 {
		return "", false
	}
	a := scrut.Args[0]
	if a.HasName || a.V == nil || a.V.Kind != "ref" || len(a.V.Ref) != 1 {
		return "", false
	}
	return a.V.Ref[0], true
}

// cellAnywhere resolves a cell program-wide for declarations: a
// consumer declares foreign capabilities it cannot use directly,
// so existence is global while use stays file-local.
func cellAnywhere(prog *Program, name string) *StateDecl {
	for _, m := range prog.Modules {
		if s := cellInFile(prog, m.ID, name); s != nil {
			return s
		}
	}
	return nil
}

// cellInFile resolves a cell in the caller's file: cells are
// module-private, so the declaration must sit beside the use.
func cellInFile(prog *Program, file, name string) *StateDecl {
	for _, m := range prog.Modules {
		if m.ID != file {
			continue
		}
		for _, d := range m.Decls {
			if s, ok := d.(*StateDecl); ok && s.Name == name {
				return s
			}
		}
	}
	return nil
}

// callSite is one executable call position: match scrutinees only.
// Non-scrutinee calls never execute (they error), so they carry no
// effects and need no attribution.
type callSite struct {
	fname string
	line  int
}

func callSites(body *Node) []callSite {
	var out []callSite
	for _, m := range matchNodes(body) {
		if m.Kind == MatchCall {
			out = append(out, callSite{m.Scruts[0].Fname, m.Line})
		}
	}
	return out
}

// effectNeed is one required capability with its attributing line
// (always in the requiring function's own text) and the level-1
// callee it arrives through ("" when used directly). Foreign
// callees ride uses pins, not helper calls, and the message says so.
type effectNeed struct {
	line    int
	via     string
	foreign bool
}

// storeCellLenient resolves the cell for capability counting even
// when the site is otherwise malformed (e.g. named args): the
// authority is exercised, so staleness must not also fire. Shape
// errors stay with the type rule.
func storeCellLenient(scrut *Small) (string, bool) {
	if len(scrut.Args) == 0 {
		return "", false
	}
	a := scrut.Args[0]
	if a.V.Kind != "ref" || len(a.V.Ref) != 1 {
		return "", false
	}
	return a.V.Ref[0], true
}

// splitCap splits a declared capability into cell and kind.
func splitCap(e string) (cell, kind string, ok bool) {
	i := strings.LastIndex(e, ".")
	if i < 0 {
		return "", "", false
	}
	return e[:i], e[i+1:], true
}

// requiredEffects unions every capability fn needs: direct store ops
// plus transitive local callees' needs plus foreign callees' declared
// effects. First use wins, so one missing capability yields one error,
// always anchored in fn's own text.
func requiredEffects(fn *FnDecl, prog *Program) map[string]effectNeed {
	out := map[string]effectNeed{}
	add := func(cap string, need effectNeed) {
		if _, ok := out[cap]; !ok {
			out[cap] = need
		}
	}
	addDirect := func(target *FnDecl, line int, via string) {
		file := prog.FnFile[target.Name]
		for _, m := range matchNodes(target.Body) {
			if m.Kind != MatchCall {
				continue
			}
			ms := m.Scruts[0]
			if !isStoreOp(ms.Fname) {
				continue
			}
			cell, ok := storeCellLenient(ms)
			if !ok || cellInFile(prog, file, cell) == nil {
				continue
			}
			cap := cell + ".read"
			if ms.Fname == "state__put" {
				cap = cell + ".write"
			}
			ln := m.Line
			if via != "" {
				ln = line
			}
			add(cap, effectNeed{ln, via, false})
		}
	}
	addDirect(fn, 0, "")
	visited := map[string]bool{fn.Name: true}
	var merge func(name string, line int, via string)
	merge = func(name string, line int, via string) {
		if visited[name] {
			return
		}
		visited[name] = true
		cf, ok := prog.Fns[name]
		if !ok {
			return
		}
		addDirect(cf, line, via)
		for _, s := range callSites(cf.Body) {
			if localCallee(prog, name, s.fname) == nil {
				continue
			}
			merge(s.fname, line, via)
		}
	}
	for _, s := range callSites(fn.Body) {
		callee, ok := prog.Fns[s.fname]
		if !ok {
			continue
		}
		if localCallee(prog, fn.Name, s.fname) != nil {
			merge(s.fname, s.line, s.fname)
		} else if prog.Uses[s.fname] {
			for _, e := range callee.Effects {
				add(e, effectNeed{s.line, s.fname, true})
			}
		}
	}
	return out
}

// effectsEntryLine finds the function-level `effects [...]` row naming
// cap, searching down from the fn declaration like emitsEntryLine.
func effectsEntryLine(text string, fnLine int, cap string) int {
	lines := strings.Split(text, "\n")
	if fnLine < 1 {
		fnLine = 1
	}
	for n := fnLine - 1; n < len(lines); n++ {
		t := strings.TrimSpace(lines[n])
		if strings.HasPrefix(t, "effects") && strings.Contains(lines[n], cap) {
			return n + 1
		}
	}
	return fnLine
}

// checkEffects proves authority before anything runs: every capability
// fn needs (directly or through callees) must be declared, every
// declared capability must be exercised, and declarations must name
// real Cell.read|write capabilities. Undeclared cells belong to the
// unknown-callee rule and are skipped here.
func checkEffects(fn *FnDecl, prog *Program, text string) []Diag {
	var out []Diag
	file := prog.FnFile[fn.Name]
	for _, m := range matchNodes(fn.Body) {
		if m.Kind != MatchCall {
			continue
		}
		ms := m.Scruts[0]
		if !isStoreOp(ms.Fname) {
			continue
		}
		cell, ok := storeCellName(ms)
		if !ok {
			continue // cell shape belongs to the type rule
		}
		if cellInFile(prog, file, cell) == nil {
			out = append(out, spanDiag(text, m.Line, "error",
				fmt.Sprintf("%s names unknown cell %s", ms.Fname, cell), cell, CodeUnknownCall))
		}
	}
	declared := map[string]bool{}
	for _, e := range fn.Effects {
		cell, kind, ok := splitCap(e)
		if !ok || (kind != "read" && kind != "write") || cellAnywhere(prog, cell) == nil {
			out = append(out, spanDiag(text, effectsEntryLine(text, fn.Line, e), "error",
				fmt.Sprintf("%s declares unknown effect %s: want Cell.read or Cell.write", fn.Name, e), e, CodeUndeclaredEffect))
			continue
		}
		declared[e] = true
	}
	req := requiredEffects(fn, prog)
	for cap, need := range req {
		if declared[cap] {
			continue
		}
		token := "state__get"
		if strings.HasSuffix(cap, ".write") {
			token = "state__put"
		}
		msg := fmt.Sprintf("%s uses %s without declaring it: add it to effects", fn.Name, cap)
		if need.via != "" {
			via := "helper " + need.via
			if need.foreign {
				via = "uses-pin " + need.via
			}
			msg = fmt.Sprintf("%s needs %s via %s without declaring it: add it to effects", fn.Name, cap, via)
		}
		out = append(out, spanDiag(text, need.line, "error", msg, token, CodeUndeclaredEffect))
	}
	for e := range declared {
		if _, ok := req[e]; !ok {
			out = append(out, spanDiag(text, effectsEntryLine(text, fn.Line, e), "error",
				fmt.Sprintf("%s declares effects %s but nothing uses it", fn.Name, e), e, CodeStaleEffect))
		}
	}
	return out
}

// checkStateDecl validates a cell declaration: the type names a base
// type (cells hold values, not records or brands) and the init is a
// literal of exactly that type.
func checkStateDecl(s *StateDecl, text string) []Diag {
	var out []Diag
	switch s.Type {
	case "str", "int", "bool", "dec":
	default:
		out = append(out, spanDiag(text, s.Line, "error",
			fmt.Sprintf("state %s holds %s: cells hold str, int, bool, or dec", s.Name, s.Type), s.Type, CodeUnknownType))
		return out
	}
	var got string
	switch s.Init.Kind {
	case "str":
		got = "str"
	case "int":
		got = "int"
	case "bool":
		got = "bool"
	case "dec":
		got = "dec"
	default:
		out = append(out, spanDiag(text, s.Line, "error",
			fmt.Sprintf("state %s init must be a literal", s.Name), s.Name, CodeTypeMismatch))
		return out
	}
	if got != s.Type {
		out = append(out, spanDiag(text, s.Line, "error",
			fmt.Sprintf("state %s holds %s but init is %s: no implicit conversions", s.Name, s.Type, got), s.Name, CodeTypeMismatch))
	}
	return out
}

// checkStub validates one scripted row for a foreign call (a12): every
// row is an exchange binding expected call args to one permitted
// outcome, so the table proves "this request received this permitted
// response", not merely the next response. Outcome-only rows are
// refused; the outcome itself keeps the existing Ok/error rules.
func checkStub(sm *Small, fname string, allowed map[string]bool, text string, entryLine int, key, owner string, out *[]Diag) {
	if sm.Kind == "list" {
		for _, it := range sm.Items {
			checkStub(it, fname, allowed, text, entryLine, key, owner, out)
		}
		return
	}
	if sm.Kind != "exchange" {
		*out = append(*out, spanDiag(text, entryLine, "error",
			fmt.Sprintf("%s: script row must be an exchange with args and outcome", owner), key, CodeNoExchange))
		return
	}
	checkOutcome(sm.Outcome, fname, allowed, text, entryLine, key, owner, out)
}

func checkOutcome(sm *Small, fname string, allowed map[string]bool, text string, entryLine int, key, owner string, out *[]Diag) {
	switch sm.Kind {
	case "ctor":
		kind := sm.Ctor
		if kind == "Ok" {
			kind = "ok"
		}
		if kind != "ok" && !strings.Contains(kind, ".") {
			*out = append(*out, spanDiag(text, entryLine, "error",
				fmt.Sprintf("%s: stub %s is neither Ok nor an error", owner, kind), kind, CodeBadStub))
			return
		}
		if !allowed[kind] {
			*out = append(*out, spanDiag(text, entryLine, "error",
				fmt.Sprintf("%s: stub %s not in %s emits", owner, kind, fname), kind, CodeStubNotInEmit))
		}
	default:
		*out = append(*out, spanDiag(text, entryLine, "error",
			fmt.Sprintf("%s: stub must be Ok(..) or an error", owner), key, CodeBadStub))
	}
}

// checkEmits enforces R5 at the function boundary: every error value the
// body can produce must be declared in emits, and every constructed kind
// must be a declared error somewhere. Declared entries are a conservative
// upper bound (a12): unrealized entries are allowed, so no consumer stub
// can manufacture provider honesty. Every entry must still name a
// declared error. Caught values forwarded whole (on e.kind var => var)
// count as produced for their kind.
// emitsEntryLine finds the function-level `emits [...]` row naming kind,
// searching down from the fn declaration so the module header (which
// comes first and may name the same kinds) never steals the hit.
func emitsEntryLine(text string, fnLine int, kind string) int {
	lines := strings.Split(text, "\n")
	if fnLine < 1 {
		fnLine = 1
	}
	for n := fnLine - 1; n < len(lines); n++ {
		t := strings.TrimSpace(lines[n])
		if strings.HasPrefix(t, "emits") && strings.Contains(lines[n], kind) {
			return n + 1
		}
	}
	return fnLine
}

// eachRaise visits every dotted error constructor a body can produce, at
// the line that produces it: arm right-hand sides (a caught variable
// returned whole counts for its kind) plus non-call match scrutinees.
// Shared by checkEmits diagnostics and the error catalog so the two can
// never disagree on what a function raises.
func eachRaise(fn *FnDecl, f func(kind string, line int)) {
	var walk func(x *Node, line int, bound map[string]string)
	walk = func(x *Node, line int, bound map[string]string) {
		if x == nil {
			return
		}
		if x.IsMatch {
			for _, a := range x.Arms {
				inner := bound
				// Call arms carry exactly one pattern; wider value
				// tables never bind error kinds here, so only
				// arity-1 arms bind.
				if len(a.Pats) == 1 {
					if p := a.Pats[0]; p.Kind == "variant" && p.Var != "" && strings.Contains(p.Name, ".") {
						inner = map[string]string{}
						for k, v := range bound {
							inner[k] = v
						}
						inner[p.Var] = p.Name
					}
				}
				walk(a.Rhs, a.Line, inner)
			}
			return
		}
		walkSmallTrees(x.Small, func(s *Small) {
			if s.Kind == "ctor" && strings.Contains(s.Ctor, ".") {
				f(s.Ctor, line)
			}
			if s.Kind == "ref" && len(s.Ref) > 0 {
				if kind, ok := bound[s.Ref[0]]; ok {
					f(kind, line)
				}
			}
		})
	}
	walk(fn.Body, fn.Line, map[string]string{})
	for _, m := range matchNodes(fn.Body) {
		for _, ms := range m.Scruts {
			if ms.Kind == "call" {
				continue
			}
			walkSmallTrees(ms, func(s *Small) {
				if s.Kind == "ctor" && strings.Contains(s.Ctor, ".") {
					f(s.Ctor, m.Line)
				}
			})
		}
	}
}

func checkEmits(fn *FnDecl, prog *Program, text string) []Diag {
	var out []Diag
	declared := map[string]bool{}
	for _, e := range fn.Emits {
		declared[e] = true
	}
	raise := func(kind string, line int) {
		if _, ok := prog.Errors[kind]; !ok {
			out = append(out, spanDiag(text, line, "error",
				fmt.Sprintf("%s raises unknown error kind %s", fn.Name, kind), kind, CodeUnknownKind))
		} else if !declared[kind] {
			d := spanDiag(text, line, "error",
				fmt.Sprintf("%s raises %s which is not in its emits", fn.Name, kind), kind, CodeForeignRaise)
			d.Found = kind
			if len(fn.Emits) == 0 {
				d.Expected = "a declared emits kind (none declared)"
			} else {
				d.Expected = "one of: " + strings.Join(fn.Emits, ", ")
			}
			d.Hint = fmt.Sprintf("add %q to the emits [...] line, or raise a declared kind", kind)
			out = append(out, d)
		}
	}
	eachRaise(fn, raise)
	for _, e := range fn.Emits {
		if _, ok := prog.Errors[e]; !ok {
			out = append(out, spanDiag(text, emitsEntryLine(text, fn.Line, e), "error",
				fmt.Sprintf("%s declares unknown error kind %s in emits", fn.Name, e), e, CodeUnknownKind))
		}
	}
	return out
}

// calledFns is every statically visible call target, scrutinee or not.
func calledFns(fn *FnDecl) map[string]bool {
	out := map[string]bool{}
	bodySmalls(fn.Body, func(s *Small, line int) {
		if s.Kind == "call" {
			out[s.Fname] = true
		}
	})
	for _, m := range matchNodes(fn.Body) {
		if m.Kind == MatchCall {
			out[m.Scruts[0].Fname] = true
		}
	}
	return out
}

// checkUnusedUses warns on uses entries no call ever reaches. Type pins
// (Uppercase names) are exempt: types are referenced in signatures, never
// called.
func checkUnusedUses(m *Module, text string, called map[string]bool) []Diag {
	var out []Diag
	for _, u := range m.Hdr["uses"] {
		base := pinRe.ReplaceAllString(u, "")
		if base == "" || base[0] >= 'A' && base[0] <= 'Z' {
			continue
		}
		if !called[base] {
			out = append(out, spanDiag(text, locateLine(text, u, 1), "warning",
				fmt.Sprintf("uses %s but %s never calls it", u, m.Mod), u, CodeUnusedUses))
		}
	}
	return out
}

// checkUnusedParams warns on parameters the body and its scripts ignore.
func checkUnusedParams(fn *FnDecl, text string) []Diag {
	var out []Diag
	used := map[string]bool{}
	mark := func(s *Small) {
		walkSmallTrees(s, func(x *Small) {
			if x.Kind == "ref" && len(x.Ref) > 0 {
				used[x.Ref[0]] = true
			}
		})
	}
	bodySmalls(fn.Body, func(s *Small, line int) { mark(s) })
	for _, m := range matchNodes(fn.Body) {
		for _, ms := range m.Scruts {
			mark(ms)
		}
		for _, sm := range m.Given {
			mark(sm)
		}
	}
	for _, p := range fn.Params {
		if !used[p[0]] {
			out = append(out, spanDiag(text, fn.Line, "warning",
				fmt.Sprintf("%s ignores param %s", fn.Name, p[0]), p[0], CodeUnusedParam))
		}
	}
	return out
}
