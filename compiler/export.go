// export.go: owner-authorized typed UTF-8 export (v46 S2).
//
// A grant (exports_utf8 Brand via fn@rev) authorizes one exact-shape
// function to disclose one brand as Bytes through the
// bytes__utf8__export kernel. certifyExports is a whole-program pass
// that runs after world build and before any linkage evaluation, in
// both the CLI and the editor pipelines: certificates must exist
// before contradictScriptOk can evaluate anything, or an uncertified
// refusal would launder into trusted script evidence.
//
// A certificate is the ExportBrand annotation on the exact permitted
// call node, set only here after full validation. Nothing is cached
// across runs: every checkProgram/diagnose builds a fresh world, so
// grant removal, signature, or body changes invalidate by
// reconstruction even when fn@rev is unchanged.
package main

import (
	"fmt"
)

// bytesExportKernel is the restricted export intrinsic. It is not a
// user-visible generic: the permitted nominal input comes from the
// certificate for the current call site, never a global signature.
const bytesExportKernel = "bytes__utf8__export"

// bytesValueRecord is the compiler-owned success record shared by the
// export kernel and (from B3) the UTF-8 codec family.
const bytesValueRecord = "Bytes__Value"

// isBytesExport reports the restricted export intrinsic.
func isBytesExport(fname string) bool {
	return fname == bytesExportKernel
}

// bytesEncodeKernel is the public UTF-8 encode intrinsic (v47 S3):
// total over str, deterministic, certificate-free.
const bytesEncodeKernel = "bytes__utf8__encode"

// bytesKernel describes one compiler kernel: its static signature,
// result record, declared emits, and whether calls need a grant
// certificate. Only the export kernel is restricted; public kernels
// check through the ordinary signature machinery.
type bytesKernel struct {
	params     [][2]string
	ret        string
	emits      []string
	restricted bool
}

// bytesKernels is the authority for kernel registration: call rules,
// given rules, EmitsOf entries, exhaustiveness, dispatch, and
// lowering all key off this table. Later slices add rows (never a
// scattered duplicate).
var bytesKernels = map[string]bytesKernel{
	bytesExportKernel: {ret: bytesValueRecord, emits: []string{}, restricted: true},
	bytesEncodeKernel: {params: [][2]string{{"value", "str"}}, ret: bytesValueRecord, emits: []string{}},
}

// isBytesKernel reports any registered Bytes kernel.
func isBytesKernel(fname string) bool {
	_, ok := bytesKernels[fname]
	return ok
}

// builtinTypeDecls returns the compiler-owned record declarations.
// B2 owns Bytes__Value only; later slices extend this list (never a
// scattered duplicate). Callers must never insert these into source
// modules or provides.
func builtinTypeDecls() []*TypeDecl {
	return []*TypeDecl{
		{Name: bytesValueRecord, Rev: 1, Fields: [][2]string{{"value", "Bytes"}}},
	}
}

// exportGrantSite retains a grant with its owning module: ownership is
// validated identities, not matching strings.
type exportGrantSite struct {
	g *Utf8ExportDecl
	m *Module
}

// certifyExports validates every exports_utf8 grant in the program and
// annotates each exact permitted call site. Authority failures are
// AIL6010, exporter-shape failures AIL6011. A grant that names a
// function (even invalidly) owns that function's diagnostics: the
// uncertified-call rule fires only where no grant names the function,
// so one problem reports once.
func certifyExports(mods []*Module, prog *Program, texts map[string]string) []Diag {
	var out []Diag
	emit := func(m *Module, d Diag) {
		d.File = qualifiedFile(mods, m)
		out = append(out, d)
	}
	var grants []exportGrantSite
	for _, m := range mods {
		for _, d := range m.Decls {
			if g, ok := d.(*Utf8ExportDecl); ok {
				grants = append(grants, exportGrantSite{g, m})
			}
		}
	}
	// Candidates are scanned directly from declarations: duplicate
	// brands are first-wins silent elsewhere, so ambiguity must be
	// detected here before that information is lost.
	brandMods := map[string][]*Module{}
	fnDecls := map[string][]*FnDecl{}
	fnMods := map[string][]*Module{}
	for _, m := range mods {
		for _, d := range m.Decls {
			switch d := d.(type) {
			case *BrandDecl:
				brandMods[d.Name] = append(brandMods[d.Name], m)
			case *FnDecl:
				fnDecls[d.Name] = append(fnDecls[d.Name], d)
				fnMods[d.Name] = append(fnMods[d.Name], m)
			}
		}
	}
	named := map[string]bool{}
	for _, gs := range grants {
		g := gs.g
		named[g.Function] = true
		text := texts[gs.m.ID]
		if gs.m.ID == "" {
			emit(gs.m, spanDiag(text, g.Line, "error",
				"exports_utf8 grant has no owning module", g.Brand, CodeBytesExportAuthority))
			continue
		}
		owners, ok := brandMods[g.Brand]
		if !ok {
			emit(gs.m, spanDiag(text, g.Line, "error",
				fmt.Sprintf("exports_utf8 grant names unknown brand %s", g.Brand), g.Brand, CodeBytesExportAuthority))
			continue
		}
		if len(owners) != 1 {
			emit(gs.m, spanDiag(text, g.Line, "error",
				fmt.Sprintf("exports_utf8 grant names ambiguous brand %s", g.Brand), g.Brand, CodeBytesExportAuthority))
			continue
		}
		fns, ok := fnDecls[g.Function]
		if !ok {
			if prog.Externs[g.Function] != nil {
				emit(gs.m, spanDiag(text, g.Line, "error",
					fmt.Sprintf("exports_utf8 grant names extern %s: exporters are ail functions", g.Function), g.Function, CodeBytesExportAuthority))
				continue
			}
			emit(gs.m, spanDiag(text, g.Line, "error",
				fmt.Sprintf("exports_utf8 grant names unknown function %s", g.Function), g.Function, CodeBytesExportAuthority))
			continue
		}
		if len(fns) != 1 {
			emit(gs.m, spanDiag(text, g.Line, "error",
				fmt.Sprintf("exports_utf8 grant names ambiguous function %s", g.Function), g.Function, CodeBytesExportAuthority))
			continue
		}
		fn := fns[0]
		fnOwner := fnMods[g.Function][0]
		if g.Revision != fn.Rev {
			emit(gs.m, spanDiag(text, g.Line, "error",
				fmt.Sprintf("exports_utf8 grant names %s@%d but %s declares rev %d", g.Function, g.Revision, g.Function, fn.Rev), g.Function, CodeBytesExportAuthority))
			continue
		}
		if gs.m.ID != owners[0].ID || gs.m.ID != fnOwner.ID || owners[0].ID == "" || fnOwner.ID == "" {
			emit(gs.m, spanDiag(text, g.Line, "error",
				fmt.Sprintf("exports_utf8 %s via %s@%d is not owner-local: grant, brand, and function must share one module", g.Brand, g.Function, g.Revision), g.Brand, CodeBytesExportAuthority))
			continue
		}
		scrut, reason := exportShape(fn, g.Brand)
		if reason != "" {
			emit(fnOwner, spanDiag(texts[fnOwner.ID], fn.Line, "error",
				fmt.Sprintf("%s is not a valid utf8 exporter: %s", fn.Name, reason), fn.Name, CodeBytesExportShape))
			continue
		}
		scrut.ExportBrand = g.Brand
	}
	// Uncertified calls: every bytes__utf8__export scrutinee must carry
	// a certificate, unless a grant already owns this function's
	// diagnostics. Calls outside scrutinees belong to the position
	// rule, not this pass.
	for _, m := range mods {
		text := texts[m.ID]
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok || named[fn.Name] {
				continue
			}
			for _, n := range matchNodes(fn.Body) {
				if n.Kind != MatchCall || len(n.Scruts) == 0 {
					continue
				}
				s := n.Scruts[0]
				if s.Kind != "call" || !isBytesExport(s.Fname) || s.ExportBrand != "" {
					continue
				}
				emit(m, spanDiag(text, n.Line, "error",
					fmt.Sprintf("call to %s in %s is not authorized by a valid exports_utf8 grant", s.Fname, fn.Name), s.Fname, CodeBytesExportAuthority))
			}
		}
	}
	return out
}

// exportShape checks the exact exporter predicate: one parameter of
// exactly the granted brand, the compiler-owned result record, an
// empty contract, and one call match over the kernel returning the
// result unchanged. It returns the permitted call node, or a reason.
func exportShape(fn *FnDecl, brand string) (*Small, string) {
	if len(fn.Params) != 1 {
		return nil, "want exactly one parameter"
	}
	param := fn.Params[0][0]
	if fn.Params[0][1] != brand {
		return nil, fmt.Sprintf("parameter %s must have type %s", param, brand)
	}
	if fn.Ret != bytesValueRecord {
		return nil, fmt.Sprintf("must return %s", bytesValueRecord)
	}
	if len(fn.Emits) != 0 {
		return nil, "must declare emits []"
	}
	if len(fn.Effects) != 0 {
		return nil, "must declare no effects"
	}
	if len(fn.DecNames) != 0 || fn.DecSchema != "" {
		return nil, "must declare no decreases clause"
	}
	n := fn.Body
	if n == nil || !n.IsMatch || n.Kind != MatchCall || len(n.Scruts) != 1 {
		return nil, fmt.Sprintf("body must be a single call match over %s", bytesExportKernel)
	}
	scrut := n.Scruts[0]
	if scrut.Kind != "call" || scrut.Fname != bytesExportKernel {
		return nil, fmt.Sprintf("body must be a single call match over %s", bytesExportKernel)
	}
	if len(scrut.Args) != 1 || scrut.Args[0].HasName || !isBareRef(scrut.Args[0].V, param) {
		return nil, fmt.Sprintf("export call takes the bare parameter %s", param)
	}
	if n.Given != nil {
		return nil, "export call takes no given table"
	}
	if len(n.Arms) != 1 {
		return nil, "want exactly one Ok arm"
	}
	arm := n.Arms[0]
	if len(arm.Pats) != 1 {
		return nil, "want exactly one Ok arm"
	}
	p := arm.Pats[0]
	if p.Kind != "variant" || p.Name != "Ok" || p.Var == "" {
		return nil, "want exactly one Ok arm"
	}
	rhs := arm.Rhs
	if rhs == nil || rhs.IsMatch {
		return nil, "must return the export result unchanged"
	}
	v := rhs.Small
	if v == nil || v.Kind != "ctor" || v.Ctor != "Ok" || len(v.Args) != 1 {
		return nil, "must return the export result unchanged"
	}
	a := v.Args[0]
	if a.Name != "value" || a.V == nil || a.V.Kind != "ref" || len(a.V.Ref) != 2 || a.V.Ref[0] != p.Var || a.V.Ref[1] != "value" {
		return nil, "must return the export result unchanged"
	}
	return scrut, ""
}
