package main

// a86: sequential fallible composition. A `match chain` block lists
// ordered call steps binding Ok payloads with one shared failure:
//
//	match chain
//	  call schema__role__check(request.role) as rc when rc.value
//	  call schema__snapshot__scan(...) as sc when sc.count == 1
//	then <tail using the binders>
//	else <one shared failure expression>
//
// `chain`, `as`, `when`, `then`, `else` are contextual only: `match
// chain` followed by a deeper `call ...` line parses as a chain,
// anything else parses exactly as before (a bool variable named
// chain keeps its value match), and the words stay ordinary
// identifiers everywhere else. Steps are single lines, so `given`
// tables cannot occur in chain position by construction.
//
// The parser builds a MatchChain node; elaborateChains (in
// buildWorld beside elaborateForwards) rewrites it into the nested
// `match call` ladder the source abbreviates, with error arms per
// emitted kind from the callee contract. Admission, the prover,
// test runs, coverage, emit, and AIL4107 therefore see ordinary
// handwritten shape — including the test-per-arm law, which is how
// every chain step's failure earns its selecting row with no new
// enforcement mechanism.

import (
	"fmt"
	"strings"
)

// parseRhs parses one arm-result position: an inline nested match, a
// single small, or a deeper block. Extracted verbatim from
// parseMatchArms so chain tails share the arm-result grammar.
func parseRhs(rows []row, i, ind, aline int, rest string) (*Node, int, error) {
	var rhs *Node
	switch {
	case strings.HasPrefix(rest, "match "):
		if rest == "match chain" && isChainBlock(rows, i, ind) {
			return parseChainBlock(rows, i, ind, aline)
		}
		sub, err := parseScrutList(strings.TrimSpace(rest[len("match "):]))
		if err != nil {
			return nil, i, at(aline, err)
		}
		rhs, i, err = parseMatchArms(rows, i, ind, aline, sub)
		if err != nil {
			return nil, i, err
		}
	case rest != "":
		var sm *Small
		if rest == "forward" || strings.HasPrefix(rest, "forward ") {
			// Slice 2: `forward v` is an arm-RHS shape only (see
			// parseMatchArms); chain tails share the restriction.
			sm = &Small{Kind: "forward", Str: strings.TrimSpace(strings.TrimPrefix(rest, "forward"))}
		} else {
			var err error
			sm, err = parseSmall(rest)
			if err != nil {
				return nil, i, at(aline, err)
			}
		}
		rhs = &Node{Small: sm, Line: aline}
	default:
		var err error
		rhs, i, err = parseExprBlock(rows, i, ind)
		if err != nil {
			return nil, i, err
		}
	}
	return rhs, i, nil
}

// parseGivenBlock parses a `given` line plus its deeper exchange
// entries. Extracted verbatim from parseMatchArms so chain steps
// script foreign calls with the identical table grammar.
func parseGivenBlock(rows []row, i, ind int) (map[string]*Small, int, error) {
	given := map[string]*Small{}
	i++
	for i < len(rows) && rows[i].indent > ind {
		gline := rows[i].line
		gm := reGiven.FindStringSubmatch(rows[i].code)
		if gm == nil {
			return nil, i, at(gline, fmt.Errorf("bad given entry: %s", rows[i].code))
		}
		rhs := strings.TrimSpace(gm[2])
		if rhs == "-" {
			given[gm[1]] = nil
		} else {
			sm, err := parseSmall(rhs)
			if err != nil {
				return nil, i, at(gline, err)
			}
			given[gm[1]] = sm
		}
		i++
	}
	return given, i, nil
}

// isChainBlock reports whether a `match chain` line opens a chain:
// the next row is deeper and starts a call step. Anything else is
// an ordinary value match (a variable named chain keeps working).
func isChainBlock(rows []row, i, indent int) bool {
	if i >= len(rows) || rows[i].indent <= indent {
		return false
	}
	return strings.HasPrefix(rows[i].code, "call ")
}

// parseChainBlock parses steps plus then/else into a MatchChain
// node. Steps bind Ok payloads; the tail continues on full success
// and the shared else (a single expression) aborts every failure.
// Steps, then, and else all sit deeper than the `match chain`
// line, arm-like, so enclosing arms can never collide with them.
func parseChainBlock(rows []row, i, indent, mline int) (*Node, int, error) {
	node := &Node{IsMatch: true, Kind: MatchChain, Line: mline}
	var tail *Node
	var elseText string
	var elseLine int
	for i < len(rows) && rows[i].indent > indent {
		c := rows[i].code
		sline := rows[i].line
		stepIndent := rows[i].indent
		switch {
		case c == "then" || strings.HasPrefix(c, "then "):
			if tail != nil {
				return nil, i, at(sline, fmt.Errorf("match chain takes one then"))
			}
			if elseText != "" {
				return nil, i, at(sline, fmt.Errorf("match chain needs then before else"))
			}
			rest := strings.TrimSpace(strings.TrimPrefix(c, "then"))
			if rest == "" {
				var err error
				tail, i, err = parseExprBlock(rows, i+1, rows[i].indent)
				if err != nil {
					return nil, i, err
				}
				continue
			}
			var err error
			tail, i, err = parseRhs(rows, i+1, rows[i].indent, sline, rest)
			if err != nil {
				return nil, i, err
			}
			if tail.Small != nil && tail.Small.Kind == "forward" {
				return nil, i, at(sline, fmt.Errorf("chain tail cannot be a bare forward: forward needs a call-outcome arm"))
			}
		case c == "else" || strings.HasPrefix(c, "else "):
			if elseText != "" {
				return nil, i, at(sline, fmt.Errorf("match chain takes one else"))
			}
			if tail == nil {
				return nil, i, at(sline, fmt.Errorf("match chain needs then before else"))
			}
			rest := strings.TrimSpace(strings.TrimPrefix(c, "else"))
			if rest == "" {
				return nil, i, at(sline, fmt.Errorf("chain else must be a single expression"))
			}
			sm, err := parseSmall(rest)
			if err != nil {
				return nil, i, at(sline, err)
			}
			if sm.Kind == "forward" {
				return nil, i, at(sline, fmt.Errorf("chain else cannot forward: no arm binder is in scope"))
			}
			elseText, elseLine = rest, sline
			i++
		case !strings.HasPrefix(c, "call "):
			return nil, i, at(sline, fmt.Errorf("chain step must be a call: call F(args) as binder [when guard], got %s", c))
		default:
			step, err := parseChainStep(c, sline)
			if err != nil {
				return nil, i, err
			}
			if tail != nil {
				return nil, i, at(sline, fmt.Errorf("match chain takes no steps after then"))
			}
			node.ChainSteps = append(node.ChainSteps, step)
			i++
			// A deeper `given` table scripts this step's call
			// with the ordinary table grammar; the elaborated
			// call node carries it, so stub discipline,
			// consistency, and leftover checks apply unchanged.
			// (Approval-path chains stay script-free through
			// the existing asset-guard test, not a ban here.)
			if i < len(rows) && rows[i].indent > stepIndent && rows[i].code == "given" {
				var err error
				step.Given, i, err = parseGivenBlock(rows, i, rows[i].indent)
				if err != nil {
					return nil, i, err
				}
				node.ChainSteps[len(node.ChainSteps)-1] = step
			}
		}
	}
	if len(node.ChainSteps) == 0 {
		return nil, i, at(mline, fmt.Errorf("match chain needs at least one call step"))
	}
	if tail == nil {
		return nil, i, at(mline, fmt.Errorf("match chain needs then"))
	}
	if elseText == "" {
		return nil, i, at(mline, fmt.Errorf("match chain needs else"))
	}
	node.ChainTail = tail
	node.ChainElse = elseText
	node.ChainElseLine = elseLine
	return node, i, nil
}

// parseChainStep parses `call F(args) as binder [when guard]`:
// top-level `as`/`when` splits stay outside strings and parens, so
// quoted text and nested calls never confuse the step boundary.
func parseChainStep(c string, sline int) (ChainStep, error) {
	var step ChainStep
	step.Line = sline
	ai, _ := findTopWord(c, []string{"as"})
	if ai < 0 {
		return step, at(sline, fmt.Errorf("chain step needs as: call F(args) as binder [when guard], got %s", c))
	}
	// The call text keeps its `call ` prefix: parseSmall only
	// reads Kind call with it, exactly as match scrutinees do.
	callText := strings.TrimSpace(c[:ai])
	rest := strings.TrimSpace(c[ai+len("as"):])
	wi, _ := findTopWord(rest, []string{"when"})
	guardText := ""
	binder := rest
	if wi >= 0 {
		binder = strings.TrimSpace(rest[:wi])
		guardText = strings.TrimSpace(rest[wi+len("when"):])
	}
	if !isChainBinder(binder) {
		return step, at(sline, fmt.Errorf("chain step needs a binder: call F(args) as binder [when guard], got %s", c))
	}
	if binder == "_" && guardText != "" {
		return step, at(sline, fmt.Errorf("chain when needs a named binder: as _ discards the value"))
	}
	call, err := parseSmall(callText)
	if err != nil {
		return step, at(sline, err)
	}
	if call.Kind != "call" {
		return step, at(sline, fmt.Errorf("chain step must be a call, got %s", callText))
	}
	step.Call = call
	step.Binder = binder
	if guardText == "" {
		return step, nil
	}
	guard, err := parseSmall(guardText)
	if err != nil {
		return step, at(sline, err)
	}
	if bad := guardImpurity(guard); bad != "" {
		return step, at(sline, fmt.Errorf("chain guard cannot %s: bind the value in an earlier step", bad))
	}
	step.Guard = guard
	return step, nil
}

// isChainBinder accepts arm-binder shapes: identifiers and `_`.
// Keyword collisions resolve downstream exactly as handwritten
// `on Ok` binders, so the chain introduces no parallel rule.
func isChainBinder(s string) bool {
	if s == "_" {
		return true
	}
	for i := 0; i < len(s); i++ {
		ch := s[i]
		ok := ch == '_' || ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' && i > 0
		if !ok {
			return false
		}
	}
	return len(s) > 0
}

// guardImpurity rejects what a short-circuit guard must not hold:
// calls (slice-5 operands stay call-free) and forward (an arm-RHS
// shape only). Sealed and pure constructors stay welcome.
func guardImpurity(sm *Small) string {
	bad := ""
	walkSmallTrees(sm, func(n *Small) {
		if bad != "" {
			return
		}
		switch n.Kind {
		case "call":
			bad = fmt.Sprintf("call %s", n.Fname)
		case "forward":
			bad = "forward a value"
		}
	})
	return bad
}

// elaborateChains rewrites chain blocks into nested call matches in
// place, deepest chains first so nested tails elaborate before the
// steps enclosing them. It runs once in buildWorld beside
// elaborateForwards, after every declaration (and callee contract)
// it reads is complete and before any proof, test run, or emit
// sees the function: downstream meets ordinary ladders only.
// Elaborated nodes carry step lines, so diagnostics and coverage
// point at the step a row must witness, never at generated code.
func elaborateChains(open *Module, prog *Program, text string) []Diag {
	var out []Diag
	fail := func(line int, format string, args ...any) {
		out = append(out, spanDiag(text, line, "error", fmt.Sprintf(format, args...), "chain", CodeChainElab))
	}
	for _, d := range open.Decls {
		fn, ok := d.(*FnDecl)
		if !ok || fn.Body == nil {
			continue
		}
		elabChainNode(fn.Body, prog, text, fail)
	}
	return out
}

// elabChainNode elaborates post-order: nested tails first, then the
// chain itself, so every MatchChain vanishes in one pass and a
// second run is a no-op. Ordinary matches recurse through arm
// results; chains recurse through the tail, whose nested chains go
// first because the tail is consumed, not shared, by levels.
func elabChainNode(n *Node, prog *Program, text string, fail func(int, string, ...any)) {
	if n == nil {
		return
	}
	if n.Kind == MatchChain {
		elabChainNode(n.ChainTail, prog, text, fail)
		elabChain(n, prog, fail)
		return
	}
	if n.IsMatch {
		for _, a := range n.Arms {
			elabChainNode(a.Rhs, prog, text, fail)
		}
	}
}

// elabChain builds the ladder inside-out: the tail is the deepest
// continuation, each step wraps it in its call match, and every
// failure level re-parses the shared else text for an unaliased
// Small. Error arms cover each emitted kind with a wildcard binder:
// payloads stay reachable only through handwritten ladders.
func elabChain(n *Node, prog *Program, fail func(int, string, ...any)) {
	cont := n.ChainTail
	for i := len(n.ChainSteps) - 1; i >= 0; i-- {
		step := n.ChainSteps[i]
		inner := cont
		if step.Guard != nil {
			inner = &Node{IsMatch: true, Kind: MatchValue,
				Scruts: []*Small{step.Guard}, Line: step.Line,
				Arms: []Arm{
					{Pats: []Pattern{{Kind: "bool", B: true}}, Rhs: cont, Line: step.Line},
					{Pats: []Pattern{{Kind: "bool", B: false}}, Rhs: chainElse(step, n, fail), Line: step.Line},
				}}
		}
		arms := []Arm{{Pats: []Pattern{{Kind: "variant", Name: "Ok", Var: step.Binder}}, Rhs: inner, Line: step.Line}}
		for _, e := range prog.EmitsOf[step.Call.Fname] {
			arms = append(arms, Arm{Pats: []Pattern{{Kind: "variantWild", Name: e}}, Rhs: chainElse(step, n, fail), Line: step.Line})
		}
		cont = &Node{IsMatch: true, Kind: MatchCall,
			Scruts: []*Small{step.Call}, Arms: arms, Given: step.Given, Line: step.Line}
	}
	// In-place rewrite keeps parent pointers stable; every field
	// the ladder carries moves across, including the step's given
	// table (dropped here it would read as script-free downstream).
	n.Kind = MatchCall
	n.Scruts = cont.Scruts
	n.Arms = cont.Arms
	n.Given = cont.Given
	n.ChainSteps = nil
	n.ChainTail = nil
	n.ChainElse = ""
}

// chainElse parses one copy of the shared failure for a level. The
// text validated at parse, so failure here is unreachable; it fails
// closed with the else line rather than aliasing one Small across
// levels (in-place const elaboration must never see shared nodes).
func chainElse(step ChainStep, n *Node, fail func(int, string, ...any)) *Node {
	sm, err := parseSmall(n.ChainElse)
	if err != nil {
		fail(n.ChainElseLine, "%s", err.Error())
		sm = &Small{Kind: "bool", B: false}
	}
	return &Node{Small: sm, Line: step.Line}
}
