// export.go: owner-authorized typed UTF-8 export (a46 S2).
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

// bytesEncodeKernel is the public UTF-8 encode intrinsic (a47 S3):
// total over str, deterministic, certificate-free.
const bytesEncodeKernel = "bytes__utf8__encode"

// bytesDecodeKernel is the public UTF-8 decode intrinsic (a50 B6):
// the first fallible kernel. Bytes in, Encoding__Text on success,
// encoding.invalid_utf8 (original payload, unchanged) on malformed
// input. Deterministic, certificate-free.
const bytesDecodeKernel = "bytes__utf8__decode"

// encodingTextRecord is the compiler-owned decode success record.
// encodingInvalidUtf8 is the compiler-owned malformed-input error,
// carrying the original Bytes payload unchanged.
const encodingTextRecord = "Encoding__Text"
const encodingInvalidUtf8 = "encoding.invalid_utf8"

// bytesHexEncodeKernel is the public hex encode intrinsic (a52 B8):
// total over Bytes, deterministic, certificate-free. Lowercase,
// byte-ordered, no text interpretation.
const bytesHexEncodeKernel = "bytes__hex__encode"

// bytesHexDecodeKernel is the public hex decode intrinsic (a55
// B10): the second fallible kernel. Even-length ASCII hex in,
// Bytes on success, encoding.invalid_hex (original string,
// unchanged) on malformed input. Deterministic,
// certificate-free.
const bytesHexDecodeKernel = "bytes__hex__decode"

// encodingInvalidHex is the compiler-owned malformed-hex error,
// carrying the original str payload unchanged.
const encodingInvalidHex = "encoding.invalid_hex"

// bytesB64EncodeKernel is the public base64 encode intrinsic (a58
// B12): total over Bytes, deterministic, certificate-free.
// Standard padded alphabet, byte-ordered, no text interpretation.
const bytesB64EncodeKernel = "bytes__base64__encode"

// bytesB64DecodeKernel is the public base64 decode intrinsic (a59
// B14): the third fallible kernel. Strict standard base64 in,
// Bytes on success, encoding.invalid_base64 (original string,
// unchanged) on malformed input. Deterministic,
// certificate-free.
const bytesB64DecodeKernel = "bytes__base64__decode"

// encodingInvalidB64 is the compiler-owned malformed-base64 error,
// carrying the original str payload unchanged.
const encodingInvalidB64 = "encoding.invalid_base64"

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
// scattered duplicate). The table is keyed by behavior, not family:
// the asset projection kernel rides the same restricted machinery
// (its name only looks out of place; renaming the table would churn
// every call site for zero behavior gain).
var bytesKernels = map[string]bytesKernel{
	bytesExportKernel: {ret: bytesValueRecord, emits: []string{}, restricted: true},
	// The asset kernel carries a static signature (unlike the export
	// kernel, whose brand varies per grant): its brands are always
	// the schema pair, so call args type-check wherever the grant
	// names resolve — a raw string smuggled where the witness belongs
	// fails statically. Execution still needs the certificate.
	assetFieldsKernel:    {params: [][2]string{{"asset", "Schema__ApprovedAsset"}, {"policy", "Schema__AssetPolicy"}}, ret: assetFieldsRecord, emits: []string{}},
	bytesEncodeKernel:    {params: [][2]string{{"value", "str"}}, ret: bytesValueRecord, emits: []string{}},
	bytesDecodeKernel:    {params: [][2]string{{"value", "Bytes"}}, ret: encodingTextRecord, emits: []string{encodingInvalidUtf8}},
	bytesHexEncodeKernel: {params: [][2]string{{"value", "Bytes"}}, ret: encodingTextRecord, emits: []string{}},
	bytesHexDecodeKernel: {params: [][2]string{{"value", "str"}}, ret: bytesValueRecord, emits: []string{encodingInvalidHex}},
	bytesB64EncodeKernel: {params: [][2]string{{"value", "Bytes"}}, ret: encodingTextRecord, emits: []string{}},
	bytesB64DecodeKernel: {params: [][2]string{{"value", "str"}}, ret: bytesValueRecord, emits: []string{encodingInvalidB64}},
}

// isBytesKernel reports any registered Bytes kernel.
func isBytesKernel(fname string) bool {
	_, ok := bytesKernels[fname]
	return ok
}

// builtinTypeDecls returns the compiler-owned record declarations.
// B2 owns Bytes__Value, B6 adds Encoding__Text, the asset bridge adds
// Asset__Fields; later slices extend this list (never a scattered
// duplicate). Callers must never insert these into source modules
// or provides.
func builtinTypeDecls() []*TypeDecl {
	return []*TypeDecl{
		{Name: bytesValueRecord, Rev: 1, Fields: [][2]string{{"value", "Bytes"}}},
		{Name: encodingTextRecord, Rev: 1, Fields: [][2]string{{"value", "str"}}},
		{Name: assetFieldsRecord, Rev: 1, Fields: [][2]string{{"url", "str"}, {"digest", "str"}, {"role", "str"}}},
	}
}

// builtinErrorDecls returns the compiler-owned error declarations
// (a50 B6). Builtin records seed every consumer first-wins; builtin
// errors follow the same rule. Source must never redeclare these,
// even identically (CodePrimitiveShadow at world build).
func builtinErrorDecls() []*ErrorDecl {
	return []*ErrorDecl{
		{Name: encodingInvalidUtf8, Fields: [][2]string{{"value", "Bytes"}}},
		{Name: encodingInvalidHex, Fields: [][2]string{{"value", "str"}}},
		{Name: encodingInvalidB64, Fields: [][2]string{{"value", "str"}}},
	}
}

// isBuiltinRecord reports a compiler-owned record name.
func isBuiltinRecord(name string) bool {
	for _, b := range builtinTypeDecls() {
		if b.Name == name {
			return true
		}
	}
	return false
}

// isBuiltinError reports a compiler-owned error kind.
func isBuiltinError(name string) bool {
	for _, b := range builtinErrorDecls() {
		if b.Name == name {
			return true
		}
	}
	return false
}

// isBytesDecode reports the fallible decode intrinsic, which needs
// its own evaluator and lowering (never the encoder path).
func isBytesDecode(fname string) bool {
	return fname == bytesDecodeKernel
}

// isBytesHexEncode reports the hex encode intrinsic, which needs
// its own evaluator and lowering (never the UTF-8 paths).
func isBytesHexEncode(fname string) bool {
	return fname == bytesHexEncodeKernel
}

// isBytesHexDecode reports the fallible hex decode intrinsic. It
// shares the generalized fallible lowering with UTF-8 decode, with
// its own codec helper, union, and error contract.
func isBytesHexDecode(fname string) bool {
	return fname == bytesHexDecodeKernel
}

// isBytesB64Encode reports the base64 encode intrinsic, which needs
// its own evaluator and lowering (never the text/hex paths).
func isBytesB64Encode(fname string) bool {
	return fname == bytesB64EncodeKernel
}

// isFallibleDecode reports a kernel lowered through the shared
// fallible path (descriptor-resolved union, per-kernel helper).
func isFallibleDecode(fname string) bool {
	return isBytesDecode(fname) || isBytesHexDecode(fname) || isBytesB64Decode(fname)
}

// isBytesB64Decode reports the fallible base64 decode intrinsic.
// It shares the generalized fallible lowering with the other
// decoders, with its own codec helper, union, and error contract.
func isBytesB64Decode(fname string) bool {
	return fname == bytesB64DecodeKernel
}

// CalleeKind is the single classification of a call scrutinee's
// callee: one ordered chain replaces the per-stage ladders in
// emit (stmtMatch), eval (evCallMatch), check (checkCalls,
// checkGiven) and linked (checkLinkedGraph), so a new callee
// kind is added once, not once per stage. The order is
// load-bearing and mirrors every ladder it replaces: store ops,
// dec__parts, bytes kernels (finest first, generic last), local
// helper, known foreign, unknown. BytesOther covers registered
// kernels with no dedicated branch (today the export kernel):
// sites keep their historical else-branch behavior for it.
// A nil prog (emit has no program context at codegen) only
// answers the intrinsic question; every non-intrinsic reports
// CalleeUnknown there, and both lower through the union path.
type CalleeKind int

const (
	CalleeStoreOp CalleeKind = iota
	CalleeDecParts
	CalleeBytesDecode
	CalleeBytesHexDecode
	CalleeBytesB64Decode
	CalleeBytesHexEncode
	CalleeBytesB64Encode
	CalleeBytesAsset
	CalleeBytesOther
	CalleeLocal
	CalleeForeign
	CalleeUnknown
)

func classifyCallee(prog *Program, owner, fname string) CalleeKind {
	switch {
	case isStoreOp(fname):
		return CalleeStoreOp
	case isDecParts(fname):
		return CalleeDecParts
	case isBytesDecode(fname):
		return CalleeBytesDecode
	case isBytesHexDecode(fname):
		return CalleeBytesHexDecode
	case isBytesB64Decode(fname):
		return CalleeBytesB64Decode
	case isBytesHexEncode(fname):
		return CalleeBytesHexEncode
	case isBytesB64Encode(fname):
		return CalleeBytesB64Encode
	case isAssetFields(fname):
		return CalleeBytesAsset
	case isBytesKernel(fname):
		return CalleeBytesOther
	}
	if localCallee(prog, owner, fname) != nil {
		return CalleeLocal
	}
	if prog != nil {
		if _, ok := prog.Fns[fname]; ok {
			return CalleeForeign
		}
		if prog.Externs[fname] != nil {
			return CalleeForeign
		}
	}
	return CalleeUnknown
}

// IsIntrinsic reports compiler-owned deterministic callees:
// store ops, dec__parts, and every bytes kernel. Intrinsics
// need no uses entry and take no given table.
func (k CalleeKind) IsIntrinsic() bool {
	switch k {
	case CalleeStoreOp, CalleeDecParts,
		CalleeBytesDecode, CalleeBytesHexDecode, CalleeBytesB64Decode,
		CalleeBytesHexEncode, CalleeBytesB64Encode,
		CalleeBytesAsset, CalleeBytesOther:
		return true
	}
	return false
}

// exportGrantSite retains a grant with its owning module: ownership is
// validated identities, not matching strings.
type exportGrantSite struct {
	g *Utf8ExportDecl
	m *Module
}

// certifyExports validates every exports_utf8 grant in the program and
// annotates each exact permitted call site. Authority failures are
// CAN6010, exporter-shape failures CAN6011. A grant that names a
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
					fmt.Sprintf("exports_utf8 grant names extern %s: exporters are can functions", g.Function), g.Function, CodeBytesExportAuthority))
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
