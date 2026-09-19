package main

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
)

// G1 explicit generics (a93): monomorphic expansion before
// checking. Every generic fn stamps one ordinary copy per
// distinct instantiation; templates never reach checking,
// and every downstream phase consumes plain monomorphic AST.
// Runs at the top of checkProgram and diagnoseWith, so CLI,
// tests, linked evaluation, and the editor share it.

// mangleInstance names one stamped copy. $ never occurs in
// source identifiers (\w+), so stamps cannot collide with
// declared names; the escape keeps the mapping injective.
func mangleInstance(base string, args []string) string {
	esc := func(a string) string {
		r := strings.NewReplacer("<", "$L$", ">", "$G$", ",", "$C$")
		return r.Replace(a)
	}
	var b strings.Builder
	b.WriteString(base)
	for _, a := range args {
		b.WriteString("$T$")
		b.WriteString(esc(a))
	}
	return b.String()
}

type genericInfo struct {
	decl *FnDecl
	mod  *Module
}

// describeStamp renders a stamped name back into source terms
// for diagnostics: m__sel$T$str reads as generic m__sel at
// <str>. Falls back to the raw name when the shape surprises.
func describeStamp(stamp, base string) string {
	rest, ok := strings.CutPrefix(stamp, base+"$T$")
	if !ok {
		return stamp
	}
	r := strings.NewReplacer("$L$", "<", "$G$", ">", "$C$", ",")
	parts := strings.Split(rest, "$T$")
	for i, p := range parts {
		parts[i] = r.Replace(p)
	}
	return "generic " + base + " at <" + strings.Join(parts, ",") + ">"
}

// expandGenerics stamps every generic fn in mods. On success it
// rewrites modules in place (stamps replace templates, call
// sites and headers name stamps) and returns nil; on failure it
// returns diagnostics and leaves mods unusable (callers abort).
func expandGenerics(mods []*Module, texts map[string]string) []Diag {
	var out []Diag
	emit := func(m *Module, line int, token, format string, args ...any) {
		d := spanDiag(texts[m.ID], line, "error", fmt.Sprintf(format, args...), token, CodeGenericExpand)
		d.File = qualifiedFile(mods, m)
		out = append(out, d)
	}
	gens := map[string]*genericInfo{}
	fnBase := map[string]*Module{}
	for _, m := range mods {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			if _, dup := fnBase[fn.Name]; dup {
				_, wasGeneric := gens[fn.Name]
				if wasGeneric || len(fn.TypeParams) > 0 {
					emit(m, fn.Line, fn.Name, "generic %s is already declared: base names are unique", fn.Name)
				}
				continue
			}
			fnBase[fn.Name] = m
			if len(fn.TypeParams) > 0 {
				gens[fn.Name] = &genericInfo{decl: fn, mod: m}
			}
		}
	}
	tgens := map[string]*genericTypeInfo{}
	typeBase := map[string]*Module{}
	for _, m := range mods {
		for _, d := range m.Decls {
			td := typeTemplateView(d)
			if td == nil {
				continue
			}
			if _, dup := typeBase[td.Name]; dup {
				_, wasGeneric := tgens[td.Name]
				if wasGeneric || len(td.TypeParams) > 0 {
					emit(m, td.Line, td.Name, "generic %s is already declared: base names are unique", td.Name)
				}
				continue
			}
			typeBase[td.Name] = m
			if len(td.TypeParams) > 0 {
				tgens[td.Name] = &genericTypeInfo{decl: td, source: d, mod: m}
			}
		}
	}
	if len(gens) == 0 && len(tgens) == 0 && !hasGenericSyntax(mods) {
		return nil
	}
	checkGenericCaseIdentities(mods, emit)
	types := genericTypeRegistry(mods)
	for base, g := range gens {
		checkGenericTemplate(base, g, emit)
	}
	for base, g := range tgens {
		checkGenericTypeTemplate(base, g, emit)
	}
	// Chain-else gate for monomorphic functions: generic ones are
	// gated in checkGenericTemplate. Late-parsed else text never
	// rewrites, so any generic head there is rejected up front.
	for _, m := range mods {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			if _, ok := gens[fn.Name]; ok {
				continue
			}
			for _, n := range matchNodes(fn.Body) {
				if n.Kind == MatchChain && chainElseMentions(n.ChainElse, nil) {
					emit(m, n.ChainElseLine, fn.Name, "chain else in %s must not mention generic calls or constructions", fn.Name)
				}
			}
		}
	}
	if len(out) > 0 {
		return out
	}
	known := collectInstances(mods, gens, tgens, types, emit)
	if len(out) > 0 {
		return out
	}
	for base := range gens {
		if len(known[base]) == 0 {
			g := gens[base]
			emit(g.mod, g.decl.Line, base, "generic %s is never instantiated: give it a row", base)
		}
	}
	if len(out) > 0 {
		return out
	}
	if diags := stampGenerics(gens, known, types, mods, texts); len(diags) > 0 {
		out = append(out, diags...)
		return out
	}
	rewriteGenericHeaders(mods, gens, known, texts, &out)
	rewriteGenericCalls(mods, gens)
	// G2 type phase: fn stamps are concrete, so every remaining
	// mention is closed except inside type templates, which
	// resolve per enclosing instance to a fixpoint.
	tknown, tused := collectTypeInstances(mods, tgens, types, emit)
	if len(out) > 0 {
		return out
	}
	for base := range tgens {
		if len(tknown[base]) == 0 {
			g := tgens[base]
			emit(g.mod, g.decl.Line, base, "generic %s is never instantiated: annotate or construct an instance", base)
		}
	}
	if len(out) > 0 {
		return out
	}
	stampGenericTypes(tgens, tknown, mods)
	rewriteGenericTypeHeaders(mods, tgens, tknown, tused, texts, &out)
	rewriteTypeMentions(mods, tgens, tknown)
	return out
}

// hasGenericSyntax reports whether any generic surface exists:
// a parameterized decl, a row with binds, or a call with type
// args. Downstream phases ignore TypeArgs entirely, so a stray
// `call f<str>(x)` with no generic in the program must still
// reach expansion (which rejects it) rather than silently
// degrading to a monomorphic call.
func hasGenericSyntax(mods []*Module) bool {
	for _, m := range mods {
		for _, d := range m.Decls {
			if td, ok := d.(*TypeDecl); ok && len(td.TypeParams) > 0 {
				return true
			}
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			if len(fn.TypeParams) > 0 {
				return true
			}
			for i := range fn.Tests {
				if len(fn.Tests[i].TypeBinds) > 0 {
					return true
				}
			}
			hit := false
			everySmall(fn, func(st smallSite) {
				if len(st.s.TypeArgs) > 0 || hasMentionSyntax(st.s.Elem) {
					hit = true
				}
			})
			if hit {
				return true
			}
		}
		hit := false
		everyModuleSmall(m, func(st smallSite) {
			if len(st.s.TypeArgs) > 0 || hasMentionSyntax(st.s.Elem) {
				hit = true
			}
		})
		if hit {
			return true
		}
	}
	// G2: a stray `Box<str>` annotation with no template in the
	// program must still reach expansion (which rejects the
	// unknown base) rather than failing downstream as a mystery.
	hit := false
	eachTypeSlot(mods, func(slot typeSlot) {
		if hasMentionSyntax(*slot.ref) {
			hit = true
		}
	})
	return hit
}

// hasMentionSyntax reports whether a type string holds a generic
// mention attempt: a `<` with a word before it other than Seq.
// Seq shapes are self-identical across expansion, never mentions.
func hasMentionSyntax(t string) bool {
	for i := 0; i < len(t); i++ {
		if t[i] != '<' {
			continue
		}
		j := i - 1
		for j >= 0 && (isWordChar(t[j]) || t[j] == '.') {
			j--
		}
		j++
		if j >= i {
			continue
		}
		if w := t[j:i]; w != "Seq" {
			return true
		}
	}
	return false
}

// genericTypeRegistry collects every name a G1 type argument may
// name: base types, Bytes, declared records, brands, variants,
// and builtin record shapes. Error kinds are not values.
func genericTypeRegistry(mods []*Module) map[string]bool {
	types := map[string]bool{"str": true, "int": true, "bool": true, "dec": true, "Bytes": true}
	for _, b := range builtinTypeDecls() {
		types[b.Name] = true
	}
	for _, m := range mods {
		for _, d := range m.Decls {
			switch d := d.(type) {
			case *TypeDecl:
				// G2: templates are not admittable arguments. A
				// bare base fills no parameter (closedArg fails
				// it as unknown), and instantiations nest.
				if len(d.TypeParams) == 0 {
					types[d.Name] = true
				}
			case *BrandDecl:
				types[d.Name] = true
			case *VariantDecl:
				if len(d.TypeParams) == 0 {
					types[d.Name] = true
				}
			}
		}
	}
	return types
}

// checkGenericTemplate enforces the template rules: every param
// is used in the signature, names collide with nothing usable
// in annotation position, emits stay concrete, the base name is
// unique, and chain-else text (parsed later, as source) mentions
// neither params nor generic calls.
func checkGenericTemplate(base string, g *genericInfo, emit func(*Module, int, string, string, ...any)) {
	fn := g.decl
	// No type-namespace collision check: params forbid `_`
	// while every type name (record, variant, brand) requires
	// `__`, so the charsets are disjoint by construction.
	seenParam := map[string]bool{}
	for _, p := range fn.TypeParams {
		seenParam[p] = true
	}
	for _, p := range fn.Params {
		if seenParam[p[0]] {
			emit(g.mod, fn.Line, p[0], "value parameter %s of %s shadows a type parameter", p[0], base)
		}
	}
	sigUses := map[string]bool{}
	markSig := func(t string) {
		for _, p := range fn.TypeParams {
			if typeMentionsParam(t, p) {
				sigUses[p] = true
			}
		}
	}
	for _, p := range fn.Params {
		markSig(p[1])
	}
	markSig(fn.Ret)
	for _, p := range fn.TypeParams {
		if !sigUses[p] {
			emit(g.mod, fn.Line, p, "type parameter %s of %s is never used in its signature", p, base)
		}
	}
	for _, e := range fn.Emits {
		if seenParam[e] {
			emit(g.mod, fn.Line, e, "emits %s of %s names a type parameter: error-set parameters are not admitted", e, base)
		}
	}
	for _, n := range matchNodes(fn.Body) {
		if n.Kind == MatchChain && chainElseMentions(n.ChainElse, fn.TypeParams) {
			emit(g.mod, n.ChainElseLine, base, "chain else in generic %s must not mention type parameters, generic calls, or generic constructions", base)
		}
	}
}

// chainElseMentions reports whether else text (reparsed later by
// chain elaboration) names a param or a generic call. Word-ish
// matching is enough: this is a rejection gate, not a parser.
func chainElseMentions(elseText string, params []string) bool {
	// A generic head (call or construction) never elaborates:
	// precise shape, so a call beside a comparison no longer
	// trips the gate.
	for i := 0; i < len(elseText); i++ {
		if elseText[i] == '<' && genericHeadLT(elseText, i) {
			return true
		}
	}
	for _, p := range params {
		for _, tok := range strings.FieldsFunc(elseText, func(r rune) bool {
			return !(r == '_' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z')
		}) {
			if tok == p {
				return true
			}
		}
	}
	return false
}

// smallSite is one expression position that may hold a generic
// call or a parameter mention, with its diagnostic line.
type smallSite struct {
	s    *Small
	line int
}

// everySmall yields every Small rooted at fn (body, tests, givens,
// contracts), plus construction-shaped proxies for case patterns.
// Updates to proxy heads/args write back to their patterns. Module-level
// const and state initializers are walked by everyModuleSmall.
func everySmall(fn *FnDecl, yield func(smallSite)) {
	var walk func(s *Small, line int)
	walk = func(s *Small, line int) {
		if s == nil {
			return
		}
		yield(smallSite{s, line})
		for _, a := range s.Args {
			walk(a.V, line)
		}
		walk(s.L, line)
		walk(s.R, line)
		walk(s.Hi, line)
		for _, it := range s.Items {
			walk(it, line)
		}
		walk(s.Outcome, line)
	}
	var walkNode func(n *Node)
	walkNode = func(n *Node) {
		if n == nil {
			return
		}
		for _, sc := range n.Scruts {
			walk(sc, n.Line)
		}
		walk(n.InvokeArg, n.Line)
		for i := range n.Arms {
			a := &n.Arms[i]
			// Case heads share explicit-argument validation and stamping
			// with constructors, but stay patterns in the downstream AST.
			var walkPattern func(*Pattern)
			walkPattern = func(p *Pattern) {
				if p.isCase() {
					s := &Small{Kind: "ctor", Ctor: p.Name, TypeArgs: p.TypeArgs}
					yield(smallSite{s, a.Line})
					p.Name, p.TypeArgs = s.Ctor, s.TypeArgs
				}
				for j := range p.Alts {
					walkPattern(&p.Alts[j])
				}
			}
			for j := range a.Pats {
				walkPattern(&a.Pats[j])
			}
			walkNode(a.Rhs)
		}
		for _, g := range n.Given {
			walk(g, n.Line)
		}
		walk(n.Small, n.Line)
		for _, st := range n.ChainSteps {
			walk(st.Call, n.Line)
			walk(st.Guard, n.Line)
			for _, g := range st.Given {
				walk(g, n.Line)
			}
		}
		walkNode(n.ChainTail)
	}
	walkNode(fn.Body)
	for i := range fn.Tests {
		t := &fn.Tests[i]
		for _, a := range t.Args {
			walk(a.V, t.Line)
		}
		walk(t.Expected, t.Line)
	}
	for _, r := range fn.Requires {
		walk(r, fn.Line)
	}
	for _, a := range fn.Ensures {
		for _, p := range a.Preds {
			walk(p, a.Line)
		}
		for _, m := range a.Matches {
			walkNode(m)
		}
	}
}

// everyModuleSmall yields const values and state initializers,
// the only Small roots outside function bodies.
func everyModuleSmall(m *Module, yield func(smallSite)) {
	var walk func(s *Small, line int)
	walk = func(s *Small, line int) {
		if s == nil {
			return
		}
		yield(smallSite{s, line})
		for _, a := range s.Args {
			walk(a.V, line)
		}
		walk(s.L, line)
		walk(s.R, line)
		walk(s.Hi, line)
		for _, it := range s.Items {
			walk(it, line)
		}
		walk(s.Outcome, line)
	}
	for _, d := range m.Decls {
		switch d := d.(type) {
		case *ConstDecl:
			walk(d.Value, d.Line)
		case *StateDecl:
			walk(d.Init, d.Line)
		}
	}
}

// closedArg validates one type argument that must name a known
// monomorphic type: no parameters, no nested instantiation, no
// variant sequences (the checkDeclFields Decision 6 deferral,
// extended to instantiation). Seq<X> admits one plain level.
func closedArg(arg string, types map[string]bool, variants map[string]bool) error {
	if elem, ok := seqElemName(arg); ok {
		if _, nested := seqElemName(elem); nested {
			return fmt.Errorf("nested instantiation %s is not admitted", arg)
		}
		if !types[elem] {
			return fmt.Errorf("unknown type %s", elem)
		}
		if variants[elem] {
			return fmt.Errorf("Seq<%s> is deferred: variant sequences are not admitted", elem)
		}
		return nil
	}
	if strings.Contains(arg, "<") {
		return fmt.Errorf("nested instantiation %s is not admitted", arg)
	}
	if !types[arg] {
		return fmt.Errorf("unknown type %s", arg)
	}
	return nil
}

// variantNames collects declared variant parents for the
// Seq<Variant> instantiation rule.
func variantNames(mods []*Module) map[string]bool {
	out := map[string]bool{}
	for _, m := range mods {
		for _, d := range m.Decls {
			if vd, ok := d.(*VariantDecl); ok {
				out[vd.Name] = true
			}
		}
	}
	return out
}

// collectInstances gathers every concrete instantiation: rows
// pin their own, monomorphic call sites name theirs, and
// pass-through calls inside generics resolve through a
// fixpoint over enclosing stamps. Only closed and
// bare-parameter forms are admitted, so the pair space is
// finite and the loop terminates when no new pair appears.
func collectInstances(mods []*Module, gens map[string]*genericInfo, tgens map[string]*genericTypeInfo, types map[string]bool, emit func(*Module, int, string, string, ...any)) map[string][][]string {
	variants := variantNames(mods)
	known := map[string][][]string{}
	seen := map[string]bool{}
	add := func(base string, args []string) bool {
		key := base + "\x00" + strings.Join(args, "\x00")
		if seen[key] {
			return false
		}
		seen[key] = true
		known[base] = append(known[base], append([]string{}, args...))
		return true
	}
	genericOf := func(name string) *genericInfo { return gens[name] }
	// Rows seed their own generic; binds must be complete and
	// concrete, and monomorphic rows take no binds.
	for _, m := range mods {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			g := genericOf(fn.Name)
			for i := range fn.Tests {
				t := &fn.Tests[i]
				if g == nil {
					if len(t.TypeBinds) > 0 {
						emit(m, t.Line, t.Name, "row %s of monomorphic %s takes no type binds", t.Name, fn.Name)
					}
					continue
				}
				args, ok := resolveBinds(fn, t, emit, m)
				if !ok {
					continue
				}
				fine := true
				for _, a := range args {
					if err := closedArg(a, types, variants); err != nil {
						emit(m, t.Line, t.Name, "row %s of %s: %v", t.Name, fn.Name, err)
						fine = false
						break
					}
				}
				if fine {
					add(fn.Name, args)
				}
			}
		}
	}
	// Call sites: arity, param-or-closed form, self-call identity.
	// Unknown callees are left for CodeUnknownCall downstream.
	type pending struct {
		mod    *Module
		caller *FnDecl
		site   smallSite
		callee string
	}
	var calls []pending
	// considerCtor validates generic construction forms while the
	// caller is still known: bare caller parameters or closed
	// types, mirroring calls. Seeding is the type phase's job;
	// this pass only rejects malformed shapes with caller scope.
	considerCtor := func(m *Module, caller *FnDecl, site smallSite) {
		s := site.s
		g, _ := genericConstructor(tgens, s.Ctor)
		if g == nil {
			if len(s.TypeArgs) > 0 {
				// Unlike calls (whose unknown callees fail
				// loudly downstream), an unknown constructed
				// base would degrade: error ctors and Ok
				// resolve with TypeArgs ignored. Fail here.
				if types[s.Ctor] {
					emit(m, site.line, s.Ctor, "type arguments on monomorphic type %s", s.Ctor)
				} else {
					emit(m, site.line, s.Ctor, "unknown type %s", s.Ctor)
				}
			}
			return
		}
		if len(s.TypeArgs) == 0 {
			emit(m, site.line, s.Ctor, "generic %s needs explicit type arguments", s.Ctor)
			return
		}
		if len(s.TypeArgs) != len(g.decl.TypeParams) {
			emit(m, site.line, s.Ctor, "generic %s takes %d type arguments, got %d", s.Ctor, len(g.decl.TypeParams), len(s.TypeArgs))
			return
		}
		params := map[string]bool{}
		if caller != nil {
			if _, ok := gens[caller.Name]; ok {
				for _, p := range caller.TypeParams {
					params[p] = true
				}
			}
		}
		for _, a := range s.TypeArgs {
			if params[a] {
				continue
			}
			if elem, ok := seqElemName(a); ok {
				if _, nested := seqElemName(elem); nested {
					emit(m, site.line, s.Ctor, "construction of %s: nested instantiation %s is not admitted", s.Ctor, a)
					return
				}
				if params[elem] {
					emit(m, site.line, s.Ctor, "construction of %s: type arguments must be bare parameters or closed types", s.Ctor)
					return
				}
			} else if strings.Contains(a, "<") {
				emit(m, site.line, s.Ctor, "construction of %s: nested instantiation %s is not admitted", s.Ctor, a)
				return
			}
			if err := closedArg(a, types, variants); err != nil {
				if len(params) > 0 && !strings.Contains(a, "<") {
					emit(m, site.line, s.Ctor, "construction of %s: %s is not a type parameter of %s", s.Ctor, a, caller.Name)
				} else {
					emit(m, site.line, s.Ctor, "construction of %s: %v", s.Ctor, err)
				}
				return
			}
		}
	}
	consider := func(m *Module, caller *FnDecl, site smallSite) {
		s := site.s
		if s.Kind == "ctor" {
			considerCtor(m, caller, site)
			return
		}
		// A reference demands its target's stamp exactly like a
		// call: same arity rules, same monomorphic rejection, same
		// self-reference repetition. The reference never executes
		// in stage 1, but its target must still exist.
		if s.Kind != "call" && s.Kind != "fnref" {
			return
		}
		g := genericOf(s.Fname)
		if g == nil {
			if len(s.TypeArgs) > 0 {
				if _, knownFn := fnBaseOf(mods, s.Fname); knownFn {
					emit(m, site.line, s.Fname, "type arguments on monomorphic function %s", s.Fname)
				}
			}
			return
		}
		if len(s.TypeArgs) == 0 {
			emit(m, site.line, s.Fname, "generic %s needs explicit type arguments", s.Fname)
			return
		}
		if len(s.TypeArgs) != len(g.decl.TypeParams) {
			emit(m, site.line, s.Fname, "generic %s takes %d type arguments, got %d", s.Fname, len(g.decl.TypeParams), len(s.TypeArgs))
			return
		}
		if caller != nil && caller.Name == s.Fname {
			for i, a := range s.TypeArgs {
				if a != g.decl.TypeParams[i] {
					what := "call"
					if s.Kind == "fnref" {
						what = "reference"
					}
					emit(m, site.line, s.Fname, "recursive %s of %s must repeat its own type arguments", what, s.Fname)
					return
				}
			}
		}
		calls = append(calls, pending{m, caller, site, s.Fname})
	}
	for _, m := range mods {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			everySmall(fn, func(st smallSite) { consider(m, fn, st) })
		}
		everyModuleSmall(m, func(st smallSite) { consider(m, nil, st) })
	}
	// Validate every call's forms once (arity already checked):
	// each arg is a bare caller parameter or a closed type.
	// Mixed, nested, and unknown forms fail here; only valid
	// calls enter the fixpoint, so errors never duplicate.
	type valid struct {
		call      pending
		usesParam bool
	}
	var valids []valid
	failed := false
	for _, c := range calls {
		params := map[string]bool{}
		if c.caller != nil && genericOf(c.caller.Name) != nil {
			for _, p := range c.caller.TypeParams {
				params[p] = true
			}
		}
		usesParam := false
		formsOK := true
		for _, a := range c.site.s.TypeArgs {
			if params[a] {
				usesParam = true
				continue
			}
			if elem, ok := seqElemName(a); ok {
				if _, nested := seqElemName(elem); nested {
					emit(c.mod, c.site.line, c.callee, "call of %s: nested instantiation %s is not admitted", c.callee, a)
					formsOK = false
					break
				}
				if params[elem] {
					emit(c.mod, c.site.line, c.callee, "call of %s: type arguments must be bare parameters or closed types", c.callee)
					formsOK = false
					break
				}
			} else if strings.Contains(a, "<") {
				emit(c.mod, c.site.line, c.callee, "call of %s: nested instantiation %s is not admitted", c.callee, a)
				formsOK = false
				break
			}
			if err := closedArg(a, types, variants); err != nil {
				// A bare unknown name inside a generic caller
				// is an unknown parameter, not an unknown
				// type: name the scope it missed.
				if len(params) > 0 && !strings.Contains(a, "<") {
					emit(c.mod, c.site.line, c.callee, "call of %s: %s is not a type parameter of %s", c.callee, a, c.caller.Name)
				} else {
					emit(c.mod, c.site.line, c.callee, "call of %s: %v", c.callee, err)
				}
				formsOK = false
				break
			}
		}
		if !formsOK {
			failed = true
			continue
		}
		valids = append(valids, valid{c, usesParam})
	}
	if failed {
		return known
	}
	// Concrete calls seed directly; pass-through resolves per
	// enclosing instance to a fixpoint. Only closed forms and
	// bare parameters survived validation, so resolution
	// synthesizes no new type expressions and the memoized pair
	// space is finite: the loop ends when nothing new appears.
	for changed := true; changed; {
		changed = false
		for _, v := range valids {
			if !v.usesParam {
				if add(v.call.callee, v.call.site.s.TypeArgs) {
					changed = true
				}
				continue
			}
			for _, inst := range known[v.call.caller.Name] {
				sub := map[string]string{}
				for i, p := range v.call.caller.TypeParams {
					sub[p] = inst[i]
				}
				resolved := make([]string, len(v.call.site.s.TypeArgs))
				for i, a := range v.call.site.s.TypeArgs {
					if r, ok := sub[a]; ok {
						resolved[i] = r
					} else {
						resolved[i] = a
					}
				}
				if add(v.call.callee, resolved) {
					changed = true
				}
			}
		}
	}
	return known
}

// resolveBinds orders a row's P=E pairs by the decl's parameter
// list, enforcing known names, completeness, and no dupes.
func resolveBinds(fn *FnDecl, t *Test, emit func(*Module, int, string, string, ...any), m *Module) ([]string, bool) {
	idx := map[string]int{}
	for i, p := range fn.TypeParams {
		idx[p] = i
	}
	args := make([]string, len(fn.TypeParams))
	filled := map[string]bool{}
	for _, b := range t.TypeBinds {
		at, ok := idx[b[0]]
		if !ok {
			emit(m, t.Line, t.Name, "row %s of %s pins unknown type parameter %s", t.Name, fn.Name, b[0])
			return nil, false
		}
		if filled[b[0]] {
			emit(m, t.Line, t.Name, "row %s of %s pins %s twice", t.Name, fn.Name, b[0])
			return nil, false
		}
		filled[b[0]] = true
		args[at] = b[1]
	}
	for _, p := range fn.TypeParams {
		if !filled[p] {
			emit(m, t.Line, t.Name, "row %s of %s must pin every type parameter: missing %s", t.Name, fn.Name, p)
			return nil, false
		}
	}
	return args, true
}

// fnBaseOf reports whether any module declares fn (any shape).
func fnBaseOf(mods []*Module, name string) (*Module, bool) {
	for _, m := range mods {
		for _, d := range m.Decls {
			if fn, ok := d.(*FnDecl); ok && fn.Name == name {
				return m, true
			}
		}
	}
	return nil, false
}

// substType replaces whole parameters: T becomes its argument,
// Seq<T> becomes Seq<arg>. Annotation grammar admits only these
// two positions, so anything else passes through untouched.
func substType(t string, sub map[string]string) string {
	if r, ok := sub[t]; ok {
		return r
	}
	if elem, ok := seqElemName(t); ok {
		if r, ok := sub[elem]; ok {
			return "Seq<" + r + ">"
		}
		// G2: Seq-wrapped mentions substitute their element.
		if inner, args, ok := splitMention(elem); ok {
			out := make([]string, len(args))
			for i, a := range args {
				out[i] = substType(strings.TrimSpace(a), sub)
			}
			return "Seq<" + inner + "<" + strings.Join(out, ",") + ">>"
		}
		return t
	}
	// G2: generic mentions substitute per argument.
	if base, args, ok := splitMention(t); ok {
		out := make([]string, len(args))
		for i, a := range args {
			out[i] = substType(strings.TrimSpace(a), sub)
		}
		return base + "<" + strings.Join(out, ",") + ">"
	}
	return t
}

// substSmall rewrites one expression in place: call type args,
// sequence element names, and constructor heads that name a
// parameter. Refs, seals, literals, and operator structure are
// type-free positions and pass through.
func substSmall(s *Small, sub map[string]string) {
	if s == nil {
		return
	}
	for i, a := range s.TypeArgs {
		if r, ok := sub[a]; ok {
			s.TypeArgs[i] = r
		}
	}
	if r, ok := sub[s.Ctor]; ok && s.Ctor != "" {
		s.Ctor = r
	}
	if s.Elem != "" {
		s.Elem = substType(s.Elem, sub)
	}
	for _, a := range s.Args {
		substSmall(a.V, sub)
	}
	substSmall(s.L, sub)
	substSmall(s.R, sub)
	substSmall(s.Hi, sub)
	for _, it := range s.Items {
		substSmall(it, sub)
	}
	substSmall(s.Outcome, sub)
}

// substNode rewrites a match tree in place. ChainElse is source
// text, not AST: templates mentioning params or generic calls
// there are rejected at template validation, so nothing here
// can hide a parameter.
func substNode(n *Node, sub map[string]string) {
	if n == nil {
		return
	}
	for _, sc := range n.Scruts {
		substSmall(sc, sub)
	}
	// Invoke arguments substitute like scrutinees: a generic
	// consumer invoking with a constructed generic argument
	// stamps closed args, never the template's parameters.
	substSmall(n.InvokeArg, sub)
	for i := range n.Arms {
		a := &n.Arms[i]
		var substPattern func(*Pattern)
		substPattern = func(p *Pattern) {
			for j, arg := range p.TypeArgs {
				p.TypeArgs[j] = substType(arg, sub)
			}
			for j := range p.Alts {
				substPattern(&p.Alts[j])
			}
		}
		for j := range a.Pats {
			substPattern(&a.Pats[j])
		}
		substNode(a.Rhs, sub)
	}
	for _, g := range n.Given {
		substSmall(g, sub)
	}
	substSmall(n.Small, sub)
	for i := range n.ChainSteps {
		substSmall(n.ChainSteps[i].Call, sub)
		substSmall(n.ChainSteps[i].Guard, sub)
		for _, g := range n.ChainSteps[i].Given {
			substSmall(g, sub)
		}
	}
	substNode(n.ChainTail, sub)
}

func cloneBig(n *big.Int) *big.Int {
	if n == nil {
		return nil
	}
	return new(big.Int).Set(n)
}

func cloneArgs(args []Arg) []Arg {
	if args == nil {
		return nil
	}
	out := make([]Arg, len(args))
	for i, a := range args {
		out[i] = Arg{Name: a.Name, HasName: a.HasName, V: cloneSmall(a.V)}
	}
	return out
}

// cloneSmall deep-copies an expression: every stamp owns its
// whole tree, since checking annotates Small.T in place and
// shared subtrees would cross-contaminate instances.
func cloneSmall(s *Small) *Small {
	if s == nil {
		return nil
	}
	c := *s
	c.Num = cloneBig(s.Num)
	c.TypeArgs = append([]string{}, s.TypeArgs...)
	c.Args = cloneArgs(s.Args)
	c.L = cloneSmall(s.L)
	c.R = cloneSmall(s.R)
	c.Hi = cloneSmall(s.Hi)
	if s.Items != nil {
		c.Items = make([]*Small, len(s.Items))
		for i, it := range s.Items {
			c.Items[i] = cloneSmall(it)
		}
	}
	c.Outcome = cloneSmall(s.Outcome)
	return &c
}

func clonePattern(p Pattern) Pattern {
	c := p
	c.TypeArgs = append([]string{}, p.TypeArgs...)
	c.Num = cloneBig(p.Num)
	c.Hi = cloneBig(p.Hi)
	if p.Alts != nil {
		c.Alts = make([]Pattern, len(p.Alts))
		for i, a := range p.Alts {
			c.Alts[i] = clonePattern(a)
		}
	}
	return c
}

func cloneGiven(g map[string]*Small) map[string]*Small {
	if g == nil {
		return nil
	}
	out := make(map[string]*Small, len(g))
	for k, v := range g {
		out[k] = cloneSmall(v)
	}
	return out
}

func cloneNode(n *Node) *Node {
	if n == nil {
		return nil
	}
	c := *n
	if n.Scruts != nil {
		c.Scruts = make([]*Small, len(n.Scruts))
		for i, s := range n.Scruts {
			c.Scruts[i] = cloneSmall(s)
		}
	}
	if n.Arms != nil {
		c.Arms = make([]Arm, len(n.Arms))
		for i, a := range n.Arms {
			c.Arms[i] = Arm{Line: a.Line, Rhs: cloneNode(a.Rhs)}
			if a.Pats != nil {
				c.Arms[i].Pats = make([]Pattern, len(a.Pats))
				for j, p := range a.Pats {
					c.Arms[i].Pats[j] = clonePattern(p)
				}
			}
		}
	}
	c.Given = cloneGiven(n.Given)
	c.Small = cloneSmall(n.Small)
	c.InvokeArg = cloneSmall(n.InvokeArg)
	if n.ChainSteps != nil {
		c.ChainSteps = make([]ChainStep, len(n.ChainSteps))
		for i, st := range n.ChainSteps {
			c.ChainSteps[i] = ChainStep{
				Call:   cloneSmall(st.Call),
				Binder: st.Binder,
				Guard:  cloneSmall(st.Guard),
				Given:  cloneGiven(st.Given),
			}
		}
	}
	c.ChainTail = cloneNode(n.ChainTail)
	c.analysis = nil
	return &c
}

func cloneTest(t Test) Test {
	c := t
	c.TypeBinds = append([][2]string{}, t.TypeBinds...)
	c.Args = cloneArgs(t.Args)
	c.Expected = cloneSmall(t.Expected)
	return c
}

// stampGenerics replaces every template with its sorted stamps.
// Each stamp substitutes its instance through signature and
// body, keeps only its routed rows (binds stripped: stamps are
// monomorphic), and inherits rev and line for error locality.
// Substituted signature positions re-validate as closed types,
// so Seq<T> at T=Seq<str> fails here with its instantiation
// named instead of confusing downstream checks.
func stampGenerics(gens map[string]*genericInfo, known map[string][][]string, types map[string]bool, mods []*Module, texts map[string]string) []Diag {
	var out []Diag
	variants := variantNames(mods)
	for base, g := range gens {
		for _, args := range known[base] {
			sub := map[string]string{}
			for i, p := range g.decl.TypeParams {
				sub[p] = args[i]
			}
			check := func(t string) {
				// Only param-mentioning positions validate
				// here; closed positions keep their downstream
				// codes (CAN6002 etc.) untouched.
				mentions := false
				if _, ok := sub[t]; ok {
					mentions = true
				} else if elem, ok := seqElemName(t); ok {
					_, mentions = sub[elem]
				}
				if !mentions {
					return
				}
				if err := closedArg(substType(t, sub), types, variants); err != nil {
					d := spanDiag(texts[g.mod.ID], g.decl.Line, "error", fmt.Sprintf("generic %s at <%s>: %v", base, strings.Join(args, ","), err), base, CodeGenericExpand)
					d.File = qualifiedFile(mods, g.mod)
					out = append(out, d)
				}
			}
			for _, p := range g.decl.Params {
				check(p[1])
			}
			check(g.decl.Ret)
		}
	}
	if len(out) > 0 {
		return out
	}
	for base, g := range gens {
		insts := append([][]string{}, known[base]...)
		sort.Slice(insts, func(i, j int) bool {
			return mangleInstance(base, insts[i]) < mangleInstance(base, insts[j])
		})
		// Row routing: each row's resolved args select its stamp.
		byKey := map[string][]Test{}
		for _, t := range g.decl.Tests {
			idx := map[string]int{}
			for i, p := range g.decl.TypeParams {
				idx[p] = i
			}
			args := make([]string, len(g.decl.TypeParams))
			for _, b := range t.TypeBinds {
				args[idx[b[0]]] = b[1]
			}
			byKey[strings.Join(args, "\x00")] = append(byKey[strings.Join(args, "\x00")], t)
		}
		var stamps []Decl
		for _, args := range insts {
			sub := map[string]string{}
			for i, p := range g.decl.TypeParams {
				sub[p] = args[i]
			}
			st := &FnDecl{
				Name:       mangleInstance(base, args),
				Rev:        g.decl.Rev,
				TypeParams: nil,
				Ret:        substType(g.decl.Ret, sub),
				Emits:      append([]string{}, g.decl.Emits...),
				Effects:    append([]string{}, g.decl.Effects...),
				DecNames:   append([]string{}, g.decl.DecNames...),
				DecSchema:  g.decl.DecSchema,
				Line:       g.decl.Line,
			}
			for _, p := range g.decl.Params {
				st.Params = append(st.Params, [2]string{p[0], substType(p[1], sub)})
			}
			st.Body = cloneNode(g.decl.Body)
			substNode(st.Body, sub)
			for _, t := range byKey[strings.Join(args, "\x00")] {
				rt := cloneTest(t)
				rt.TypeBinds = nil
				substTest(&rt, sub)
				st.Tests = append(st.Tests, rt)
			}
			for _, r := range g.decl.Requires {
				c := cloneSmall(r)
				substSmall(c, sub)
				st.Requires = append(st.Requires, c)
			}
			for _, a := range g.decl.Ensures {
				ca := ContractArm{Outcome: a.Outcome, Bind: a.Bind, Line: a.Line}
				for _, p := range a.Preds {
					c := cloneSmall(p)
					substSmall(c, sub)
					ca.Preds = append(ca.Preds, c)
				}
				for _, m := range a.Matches {
					c := cloneNode(m)
					substNode(c, sub)
					ca.Matches = append(ca.Matches, c)
				}
				st.Ensures = append(st.Ensures, ca)
			}
			stamps = append(stamps, st)
		}
		m := g.mod
		// Identity splice: a recorded index goes stale the
		// moment a same-module sibling stamps first (map
		// order), so the template is located by pointer.
		// Each iteration removes exactly its own template,
		// so the end state is order-independent.
		at := -1
		for i, d := range m.Decls {
			if d == Decl(g.decl) {
				at = i
				break
			}
		}
		if at < 0 {
			continue
		}
		decls := make([]Decl, 0, len(m.Decls)-1+len(stamps))
		decls = append(decls, m.Decls[:at]...)
		decls = append(decls, stamps...)
		decls = append(decls, m.Decls[at+1:]...)
		m.Decls = decls
		if m.GenericBase == nil {
			m.GenericBase = map[string]string{}
		}
		for _, args := range insts {
			m.GenericBase[mangleInstance(base, args)] = base
		}
	}
	return out
}

// substTest rewrites a routed row's args and expectation.
func substTest(t *Test, sub map[string]string) {
	for _, a := range t.Args {
		substSmall(a.V, sub)
	}
	substSmall(t.Expected, sub)
}

// rewriteGenericCalls renames every generic call site to its
// stamped copy. Monomorphic callers name concrete args
// directly; stamps were substituted at stamping, so their
// args are concrete here too.
func rewriteGenericCalls(mods []*Module, gens map[string]*genericInfo) {
	rewrite := func(s *Small) {
		if s == nil || (s.Kind != "call" && s.Kind != "fnref") || len(s.TypeArgs) == 0 {
			return
		}
		if _, ok := gens[s.Fname]; !ok {
			return
		}
		s.Fname = mangleInstance(s.Fname, s.TypeArgs)
		s.TypeArgs = nil
	}
	var walk func(s *Small)
	walk = func(s *Small) {
		if s == nil {
			return
		}
		rewrite(s)
		for _, a := range s.Args {
			walk(a.V)
		}
		walk(s.L)
		walk(s.R)
		walk(s.Hi)
		for _, it := range s.Items {
			walk(it)
		}
		walk(s.Outcome)
	}
	var walkNode func(n *Node)
	walkNode = func(n *Node) {
		if n == nil {
			return
		}
		for _, sc := range n.Scruts {
			walk(sc)
		}
		walk(n.InvokeArg)
		for _, a := range n.Arms {
			walkNode(a.Rhs)
		}
		for _, g := range n.Given {
			walk(g)
		}
		walk(n.Small)
		for i := range n.ChainSteps {
			walk(n.ChainSteps[i].Call)
			walk(n.ChainSteps[i].Guard)
			for _, g := range n.ChainSteps[i].Given {
				walk(g)
			}
		}
		walkNode(n.ChainTail)
	}
	for _, m := range mods {
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			walkNode(fn.Body)
			for i := range fn.Tests {
				for _, a := range fn.Tests[i].Args {
					walk(a.V)
				}
				walk(fn.Tests[i].Expected)
			}
			for _, r := range fn.Requires {
				walk(r)
			}
			for _, a := range fn.Ensures {
				for _, p := range a.Preds {
					walk(p)
				}
				for _, mt := range a.Matches {
					walkNode(mt)
				}
			}
		}
		everyModuleSmall(m, func(st smallSite) { walk(st.s) })
	}
}

// rewriteGenericHeaders swaps base names for stamps in provides
// and per-module used stamps in uses. Runs before call
// rewriting, so used sets read TypeArgs directly. A base pin with no
// pinned call/reference use drops with the same warning a dead
// monomorphic pin gets (CodeUnusedUses, same severity).
func rewriteGenericHeaders(mods []*Module, gens map[string]*genericInfo, known map[string][][]string, texts map[string]string, out *[]Diag) {
	for base, g := range gens {
		var stamps []string
		for _, args := range known[base] {
			stamps = append(stamps, mangleInstance(base, args))
		}
		sort.Strings(stamps)
		var provides []string
		for _, p := range g.mod.Hdr["provides"] {
			if p == base {
				provides = append(provides, stamps...)
			} else {
				provides = append(provides, p)
			}
		}
		g.mod.Hdr["provides"] = provides
	}
	for _, m := range mods {
		used := map[string]map[string]bool{}
		note := func(s *Small) {
			if s == nil || (s.Kind != "call" && s.Kind != "fnref") || len(s.TypeArgs) == 0 {
				return
			}
			if _, ok := gens[s.Fname]; !ok {
				return
			}
			if used[s.Fname] == nil {
				used[s.Fname] = map[string]bool{}
			}
			used[s.Fname][mangleInstance(s.Fname, s.TypeArgs)] = true
		}
		// Same positions calledFns polices: leaf Smalls plus
		// MatchCall scrutinees. Test-arg and script calls need
		// no pins, exactly like monomorphic calls. Function references
		// do need pins in every position, including test data (B00).
		for _, d := range m.Decls {
			fn, ok := d.(*FnDecl)
			if !ok {
				continue
			}
			everySmall(fn, func(st smallSite) {
				if st.s.Kind == "fnref" {
					note(st.s)
				}
			})
			bodySmalls(fn.Body, func(s *Small, line int) { note(s) })
			for _, n := range matchNodes(fn.Body) {
				if n.Kind == MatchCall && len(n.Scruts) > 0 {
					note(n.Scruts[0])
				}
			}
		}
		var uses []string
		for _, u := range m.Hdr["uses"] {
			base := pinRe.ReplaceAllString(u, "")
			if _, ok := gens[base]; !ok {
				uses = append(uses, u)
				continue
			}
			set := used[base]
			if len(set) == 0 {
				d := spanDiag(texts[m.ID], locateLine(texts[m.ID], u, 1), "warning",
					fmt.Sprintf("uses %s but %s never instantiates it", u, m.Mod), u, CodeUnusedUses)
				d.File = qualifiedFile(mods, m)
				*out = append(*out, d)
				continue
			}
			suffix := ""
			if i := strings.LastIndex(u, "@"); i >= 0 {
				suffix = u[i:]
			}
			var mangled []string
			for stamp := range set {
				mangled = append(mangled, stamp+suffix)
			}
			sort.Strings(mangled)
			uses = append(uses, mangled...)
		}
		m.Hdr["uses"] = uses
	}
}

// Explicit generic records (a97) and variants (b06): the type phase runs after
// the fn phase, so every mention it sees is closed except inside
// type templates, which resolve per enclosing instance to a
// fixpoint. Stamps are plain records/variants; every downstream phase
// consumes them untouched.

// splitMention splits a whole type string of the shape Base<args>
// into its base and top-level arguments. ok=false for plain names
// and unbalanced shapes.
func splitMention(t string) (string, []string, bool) {
	lt := strings.IndexByte(t, '<')
	if lt <= 0 || !strings.HasSuffix(t, ">") {
		return "", nil, false
	}
	base := t[:lt]
	if !reWord.MatchString(base) {
		return "", nil, false
	}
	args, err := splitTypeArgs(t[lt+1 : len(t)-1])
	if err != nil {
		return "", nil, false
	}
	return base, args, true
}

// mentionOf extracts the generic mention in a type string, if it
// is one: Base<args> with a non-Seq base, unwrapping one Seq
// layer first. Plain names and Seq shapes are not mentions.
func mentionOf(t string) (string, []string, bool) {
	if strings.HasPrefix(t, "Seq<") && strings.HasSuffix(t, ">") {
		if elem, args, ok := splitMention(t); ok && elem == "Seq" && len(args) == 1 {
			t = strings.TrimSpace(args[0])
		} else {
			return "", nil, false
		}
	}
	// B00 stage 1a: the Fn<head> is builtin, never a template to
	// stamp, so it is not a mention. Descending into its input
	// and success types for instance collection (and substituting
	// through it) lands with expansion; until then a generic
	// record field holding Fn keeps its literal parameters.
	if _, _, _, ok := fnTypeShape(t); ok {
		return "", nil, false
	}
	base, args, ok := splitMention(t)
	if !ok || base == "Seq" {
		return "", nil, false
	}
	return base, args, true
}

// typeMentionsParam reports whether the type string names the
// parameter: bare, Seq-wrapped, inside a generic mention's
// arguments, or inside a function head's input or success,
// recursing so Box<Seq<T>> and Fn<T, M__O, []> count.
func typeMentionsParam(t, p string) bool {
	if t == p {
		return true
	}
	if elem, ok := seqElemName(t); ok {
		return typeMentionsParam(elem, p)
	}
	if a, r, _, ok := fnTypeShape(t); ok {
		return typeMentionsParam(strings.TrimSpace(a), p) ||
			typeMentionsParam(strings.TrimSpace(r), p)
	}
	if _, args, ok := mentionOf(t); ok {
		for _, a := range args {
			if typeMentionsParam(strings.TrimSpace(a), p) {
				return true
			}
		}
	}
	return false
}

// typeSlot is one annotation position holding a type string: a
// place a generic mention may appear. ref points at the live
// string so collection and rewriting share one walk.
type typeSlot struct {
	mod   *Module
	line  int
	token string
	tpl   *TypeDecl // non-nil when inside a type template
	state bool
	ref   *string
}

func eachTypeSlot(mods []*Module, yield func(typeSlot)) {
	for _, m := range mods {
		for _, d := range m.Decls {
			switch d := d.(type) {
			case *FnDecl:
				for i := range d.Params {
					yield(typeSlot{m, d.Line, d.Params[i][0], nil, false, &d.Params[i][1]})
				}
				yield(typeSlot{m, d.Line, d.Name, nil, false, &d.Ret})
			case *TypeDecl:
				var tpl *TypeDecl
				if len(d.TypeParams) > 0 {
					tpl = d
				}
				for i := range d.Fields {
					yield(typeSlot{m, d.Line, d.Fields[i][0], tpl, false, &d.Fields[i][1]})
				}
			case *VariantDecl:
				var tpl *TypeDecl
				if len(d.TypeParams) > 0 {
					tpl = typeTemplateView(d)
				}
				for ci := range d.Cases {
					for i := range d.Cases[ci].Fields {
						yield(typeSlot{m, d.Cases[ci].Line, d.Cases[ci].Fields[i][0], tpl, false, &d.Cases[ci].Fields[i][1]})
					}
				}
			case *ErrorDecl:
				for i := range d.Fields {
					yield(typeSlot{m, d.Line, d.Fields[i][0], nil, false, &d.Fields[i][1]})
				}
			case *ConstDecl:
				yield(typeSlot{m, d.Line, d.Name, nil, false, &d.Type})
			case *ExternDecl:
				for i := range d.Params {
					yield(typeSlot{m, d.Line, d.Params[i][0], nil, false, &d.Params[i][1]})
				}
				yield(typeSlot{m, d.Line, d.Name, nil, false, &d.Ret})
			case *StateDecl:
				yield(typeSlot{m, d.Line, d.Name, nil, true, &d.Type})
			}
		}
	}
}

// typeTemplateView shares field/parameter validation and collection.
// A variant's flattened fields are read-only; stamps retain its cases.
func typeTemplateView(d Decl) *TypeDecl {
	switch d := d.(type) {
	case *TypeDecl:
		return d
	case *VariantDecl:
		t := &TypeDecl{Name: d.Name, Rev: d.Rev, TypeParams: d.TypeParams, Line: d.Line}
		for _, c := range d.Cases {
			t.Fields = append(t.Fields, c.Fields...)
		}
		return t
	}
	return nil
}

type genericTypeInfo struct {
	decl   *TypeDecl // common view, never emitted for a variant
	source Decl
	mod    *Module
}

// A variant parent is a type, not a constructor; its qualified cases
// resolve to the owning template for validation and instance demand.
func genericConstructor(gens map[string]*genericTypeInfo, name string) (*genericTypeInfo, string) {
	if g := gens[name]; g != nil {
		if _, variant := g.source.(*VariantDecl); !variant {
			return g, name
		}
	}
	for base, g := range gens {
		if v, ok := g.source.(*VariantDecl); ok {
			for _, c := range v.Cases {
				if qualifyCase(base, c.Short) == name {
					return g, base
				}
			}
		}
	}
	return nil, ""
}

// Source collisions must not disappear behind disjoint instance args.
func checkGenericCaseIdentities(mods []*Module, emit func(*Module, int, string, string, ...any)) {
	type owner struct {
		name    string
		generic bool
	}
	owners := map[string]owner{}
	records := map[string]bool{}
	for _, m := range mods {
		for _, d := range m.Decls {
			if d, ok := d.(*TypeDecl); ok {
				records[d.Name] = true
			}
		}
	}
	for _, m := range mods {
		for _, d := range m.Decls {
			v, ok := d.(*VariantDecl)
			if !ok {
				continue
			}
			generic := len(v.TypeParams) > 0
			for _, c := range v.Cases {
				q := qualifyCase(v.Name, c.Short)
				if old, exists := owners[q]; exists && (generic || old.generic) {
					emit(m, c.Line, c.Short, "case identity collision: %s declared by both %s and %s", q, old.name, v.Name)
				}
				if generic && records[q] {
					emit(m, c.Line, c.Short, "case identity collision: %s collides with record type %s", q, q)
				}
				owners[q] = owner{v.Name, generic}
			}
		}
	}
}

// checkGenericTypeTemplate enforces the template rules for a generic
// record or variant: every parameter is used in the fields, and the template
// never fields its own base directly. Types have no termination
// argument, so direct recursion is rejected; Seq-wrapped
// self-reference stays legal, exactly like monomorphic records.
func checkGenericTypeTemplate(base string, g *genericTypeInfo, emit func(*Module, int, string, string, ...any)) {
	if _, variant := g.source.(*VariantDecl); variant && !typeNameRe.MatchString(base) {
		emit(g.mod, g.decl.Line, base, "variant name %q must match Domain__Name", base)
	}
	uses := map[string]bool{}
	for _, f := range g.decl.Fields {
		for _, p := range g.decl.TypeParams {
			if typeMentionsParam(f[1], p) {
				uses[p] = true
			}
		}
		if mbase, _, ok := splitMention(f[1]); ok && mbase == base {
			emit(g.mod, g.decl.Line, f[0], "generic %s fields its own base: recursive types are not admitted", base)
		}
	}
	for _, p := range g.decl.TypeParams {
		if !uses[p] {
			emit(g.mod, g.decl.Line, p, "type parameter %s of %s is never used in its fields", p, base)
		}
	}
}

// pairKey memoizes (base, args) pairs across collection and
// rewriting, mirroring the fn instance memo.
func pairKey(base string, args []string) string {
	return base + "\x00" + strings.Join(args, "\x00")
}

// collectTypeInstances gathers every (base, args) instance a generic
// record needs, plus the per-module used sets for pin rewriting.
// Closed mentions seed directly; mentions inside type templates
// resolve per enclosing instance to a fixpoint. Only closed forms
// and bare own-parameters survive validation, so the memoized pair
// space is finite: the loop ends when nothing new appears.
func collectTypeInstances(mods []*Module, tgens map[string]*genericTypeInfo, types map[string]bool, emit func(*Module, int, string, string, ...any)) (map[string][][]string, map[*Module]map[string]map[string]bool) {
	variants := variantNames(mods)
	known := map[string][][]string{}
	seen := map[string]bool{}
	add := func(base string, args []string) bool {
		key := pairKey(base, args)
		if seen[key] {
			return false
		}
		seen[key] = true
		known[base] = append(known[base], append([]string{}, args...))
		return true
	}
	used := map[*Module]map[string]map[string]bool{}
	note := func(m *Module, base string, args []string) {
		if used[m] == nil {
			used[m] = map[string]map[string]bool{}
		}
		if used[m][base] == nil {
			used[m][base] = map[string]bool{}
		}
		used[m][base][mangleInstance(base, args)] = true
	}
	failed := false
	// seedClosed validates one closed mention and records it.
	seedClosed := func(m *Module, line int, token, base string, args []string) {
		g := tgens[base]
		if g == nil {
			if types[base] {
				emit(m, line, token, "type arguments on monomorphic type %s", base)
			} else {
				emit(m, line, token, "unknown type %s", base)
			}
			failed = true
			return
		}
		if len(args) != len(g.decl.TypeParams) {
			emit(m, line, token, "generic %s takes %d type arguments, got %d", base, len(g.decl.TypeParams), len(args))
			failed = true
			return
		}
		for _, a := range args {
			if err := closedArg(strings.TrimSpace(a), types, variants); err != nil {
				emit(m, line, token, "%v", err)
				failed = true
				return
			}
		}
		add(base, args)
		note(m, base, args)
	}
	// Template-relative mentions resolve per enclosing instance
	// below; the defining module owns the resolved reference.
	type relative struct {
		mod   *Module
		tpl   *TypeDecl
		line  int
		token string
		base  string
		args  []string
	}
	var relatives []relative
	// handleTypeMention records one generic mention found in a slot:
	// closed mentions demand stamps now, template mentions carrying
	// parameters resolve per enclosing instance below.
	handleTypeMention := func(slot typeSlot, base string, args []string) {
		if slot.tpl == nil {
			seedClosed(slot.mod, slot.line, slot.token, base, args)
			return
		}
		params := map[string]bool{}
		for _, p := range slot.tpl.TypeParams {
			params[p] = true
		}
		usesParam := false
		for _, a := range args {
			a = strings.TrimSpace(a)
			if params[a] {
				usesParam = true
				continue
			}
			// Annotations admit Seq-wrapped parameters (unlike
			// call and construction argument lists, which take
			// bare parameters or closed types only).
			if elem, isSeq := seqElemName(a); isSeq && params[elem] {
				usesParam = true
				continue
			}
			if err := closedArg(a, types, variants); err != nil {
				emit(slot.mod, slot.line, slot.token, "%v", err)
				failed = true
				return
			}
		}
		g := tgens[base]
		if g == nil {
			if types[base] {
				emit(slot.mod, slot.line, slot.token, "type arguments on monomorphic type %s", base)
			} else {
				emit(slot.mod, slot.line, slot.token, "unknown type %s", base)
			}
			failed = true
			return
		}
		if len(args) != len(g.decl.TypeParams) {
			emit(slot.mod, slot.line, slot.token, "generic %s takes %d type arguments, got %d", base, len(g.decl.TypeParams), len(args))
			failed = true
			return
		}
		if !usesParam {
			add(base, args)
			note(slot.mod, base, args)
			return
		}
		relatives = append(relatives, relative{slot.mod, slot.tpl, slot.line, slot.token, base, args})
	}
	eachTypeSlot(mods, func(slot typeSlot) {
		t := strings.TrimSpace(*slot.ref)
		// A function head is not itself a mention, but its input
		// and success positions hold ordinary mentions that
		// demand stamps exactly as top-level slots do.
		if a, r, _, ok := fnTypeShape(t); ok {
			for _, sub := range []string{a, r} {
				if base, args, ok := mentionOf(strings.TrimSpace(sub)); ok {
					handleTypeMention(slot, base, args)
				}
			}
			return
		}
		base, args, ok := mentionOf(t)
		if !ok {
			return
		}
		if slot.state {
			// Same shape as checkStateDecl, pre-rewrite so the
			// message names the source mention, never a stamp.
			emit(slot.mod, slot.line, t, "state %s holds %s: cells hold str, int, bool, or dec", slot.token, t)
			failed = true
			return
		}
		handleTypeMention(slot, base, args)
	})
	considerSmall := func(m *Module, s *Small, line int) {
		if s == nil {
			return
		}
		// Generic constructions seed their instance. Forms were
		// validated in the fn pass; re-check closed args here
		// so walker drift fails closed, never silently.
		if s.Kind == "ctor" && len(s.TypeArgs) > 0 {
			_, base := genericConstructor(tgens, s.Ctor)
			if base == "" {
				return
			}
			for _, a := range s.TypeArgs {
				if err := closedArg(a, types, variants); err != nil {
					emit(m, line, s.Ctor, "construction of %s: %v", s.Ctor, err)
					failed = true
					return
				}
			}
			add(base, s.TypeArgs)
			note(m, base, s.TypeArgs)
		}
		if s.Elem != "" {
			if base, args, ok := mentionOf(s.Elem); ok {
				seedClosed(m, line, "Seq", base, args)
			}
		}
	}
	for _, m := range mods {
		for _, d := range m.Decls {
			if fn, ok := d.(*FnDecl); ok {
				everySmall(fn, func(st smallSite) { considerSmall(m, st.s, st.line) })
			}
		}
		everyModuleSmall(m, func(st smallSite) { considerSmall(m, st.s, st.line) })
	}
	if failed {
		return known, used
	}
	for changed := true; changed; {
		changed = false
		for _, r := range relatives {
			for _, inst := range known[r.tpl.Name] {
				sub := map[string]string{}
				for i, p := range r.tpl.TypeParams {
					sub[p] = inst[i]
				}
				resolved := make([]string, len(r.args))
				ok := true
				for i, a := range r.args {
					resolved[i] = substType(strings.TrimSpace(a), sub)
					if err := closedArg(resolved[i], types, variants); err != nil {
						emit(r.mod, r.line, r.token, "%v", err)
						failed = true
						ok = false
						break
					}
				}
				if !ok {
					continue
				}
				if add(r.base, resolved) {
					changed = true
				}
				note(r.mod, r.base, resolved)
			}
		}
	}
	return known, used
}

// stampGenericTypes replaces each type template with one plain
// record or variant per instance, substituting fields per instance. Stamps
// take the template rev and splice by identity, so sibling order
// never strands a template (a93 amendment 6 applies here too).
func stampGenericTypes(tgens map[string]*genericTypeInfo, known map[string][][]string, mods []*Module) {
	for base, g := range tgens {
		insts := append([][]string{}, known[base]...)
		sort.Slice(insts, func(i, j int) bool {
			return mangleInstance(base, insts[i]) < mangleInstance(base, insts[j])
		})
		var stamps []Decl
		for _, args := range insts {
			sub := map[string]string{}
			for i, p := range g.decl.TypeParams {
				sub[p] = args[i]
			}
			st := &TypeDecl{
				Name:       mangleInstance(base, args),
				Rev:        g.decl.Rev,
				TypeParams: nil,
				Line:       g.decl.Line,
			}
			if v, ok := g.source.(*VariantDecl); ok {
				vs := &VariantDecl{Name: st.Name, Rev: st.Rev, Line: st.Line}
				for _, c := range v.Cases {
					cs := VariantCase{Short: c.Short, Line: c.Line}
					for _, f := range c.Fields {
						cs.Fields = append(cs.Fields, [2]string{f[0], substType(f[1], sub)})
					}
					vs.Cases = append(vs.Cases, cs)
				}
				stamps = append(stamps, vs)
			} else {
				for _, f := range g.decl.Fields {
					st.Fields = append(st.Fields, [2]string{f[0], substType(f[1], sub)})
				}
				stamps = append(stamps, st)
			}
		}
		m := g.mod
		at := -1
		for i, d := range m.Decls {
			if d == g.source {
				at = i
				break
			}
		}
		if at < 0 {
			continue
		}
		decls := make([]Decl, 0, len(m.Decls)-1+len(stamps))
		decls = append(decls, m.Decls[:at]...)
		decls = append(decls, stamps...)
		decls = append(decls, m.Decls[at+1:]...)
		m.Decls = decls
		if m.GenericBase == nil {
			m.GenericBase = map[string]string{}
		}
		for _, args := range insts {
			m.GenericBase[mangleInstance(base, args)] = base
		}
	}
}

// rewriteGenericTypeHeaders swaps base names for stamps in provides
// and per-module used stamps in uses, mirroring the fn pass.
// Type pins are optional, but when present they resolve against
// stamps exactly like fn pins, so text-level tools (which see
// base names) and the compiler agree.
func rewriteGenericTypeHeaders(mods []*Module, tgens map[string]*genericTypeInfo, known map[string][][]string, used map[*Module]map[string]map[string]bool, texts map[string]string, out *[]Diag) {
	for base, g := range tgens {
		var stamps []string
		for _, args := range known[base] {
			stamps = append(stamps, mangleInstance(base, args))
		}
		sort.Strings(stamps)
		var provides []string
		for _, p := range g.mod.Hdr["provides"] {
			if p == base {
				provides = append(provides, stamps...)
			} else {
				provides = append(provides, p)
			}
		}
		g.mod.Hdr["provides"] = provides
	}
	for _, m := range mods {
		var uses []string
		for _, u := range m.Hdr["uses"] {
			base := pinRe.ReplaceAllString(u, "")
			if _, ok := tgens[base]; !ok {
				uses = append(uses, u)
				continue
			}
			set := used[m][base]
			if len(set) == 0 {
				d := spanDiag(texts[m.ID], locateLine(texts[m.ID], u, 1), "warning",
					fmt.Sprintf("uses %s but %s never instantiates it", u, m.Mod), u, CodeUnusedUses)
				d.File = qualifiedFile(mods, m)
				*out = append(*out, d)
				continue
			}
			suffix := ""
			if i := strings.LastIndex(u, "@"); i >= 0 {
				suffix = u[i:]
			}
			var mangled []string
			for stamp := range set {
				mangled = append(mangled, stamp+suffix)
			}
			sort.Strings(mangled)
			uses = append(uses, mangled...)
		}
		m.Hdr["uses"] = uses
	}
}

// rewriteTypeMentions renames every generic mention to its stamped
// copy: annotation strings, generic constructions, and sequence
// element types. Runs after stamping, so every valid mention
// resolves; invalid ones aborted the phase before.
func rewriteTypeMentions(mods []*Module, tgens map[string]*genericTypeInfo, known map[string][][]string) {
	stampOf := map[string]string{}
	for base, insts := range known {
		if _, ok := tgens[base]; !ok {
			continue
		}
		for _, args := range insts {
			stampOf[pairKey(base, args)] = mangleInstance(base, args)
		}
	}
	var rewriteStr func(t string) string
	rewriteStr = func(t string) string {
		// Function heads rewrite through their input and success
		// positions; untouched heads return byte-identical so no
		// spacing churn reaches string-compared annotations.
		if a, r, e, ok := fnTypeShape(t); ok {
			ra, rb := rewriteStr(strings.TrimSpace(a)), rewriteStr(strings.TrimSpace(r))
			if ra == strings.TrimSpace(a) && rb == strings.TrimSpace(r) {
				return t
			}
			return "Fn<" + ra + "," + rb + "," + strings.TrimSpace(e) + ">"
		}
		base, args, ok := mentionOf(t)
		if !ok {
			return t
		}
		stamp, ok := stampOf[pairKey(base, args)]
		if !ok {
			return t
		}
		if strings.HasPrefix(t, "Seq<") {
			return "Seq<" + stamp + ">"
		}
		return stamp
	}
	eachTypeSlot(mods, func(slot typeSlot) {
		*slot.ref = rewriteStr(*slot.ref)
	})
	rewriteOne := func(s *Small) {
		if s == nil {
			return
		}
		if s.Kind == "ctor" && len(s.TypeArgs) > 0 {
			if _, base := genericConstructor(tgens, s.Ctor); base != "" {
				if _, ok := stampOf[pairKey(base, s.TypeArgs)]; ok {
					s.Ctor = mangleInstance(s.Ctor, s.TypeArgs)
					s.TypeArgs = nil
				}
			}
		}
		if s.Elem != "" {
			s.Elem = rewriteStr(s.Elem)
		}
	}
	for _, m := range mods {
		for _, d := range m.Decls {
			if fn, ok := d.(*FnDecl); ok {
				everySmall(fn, func(st smallSite) { rewriteOne(st.s) })
			}
		}
		everyModuleSmall(m, func(st smallSite) { rewriteOne(st.s) })
	}
}
