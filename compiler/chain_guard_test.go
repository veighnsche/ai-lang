package main

import (
	"strings"
	"testing"
)

// TestChainGuardsEveryStage pins the post-buildWorld invariant at
// every consumer: a MatchChain reaching emit, evaluation, or proof
// fails closed. Elaboration needs EmitsOf, so the parser cannot
// build ladders natively; these guards are the backstop.
func TestChainGuardsEveryStage(t *testing.T) {
	chain := &Node{IsMatch: true, Kind: MatchChain, Line: 7}
	if elaborated(chain) {
		t.Fatalf("elaborated(chain) = true")
	}
	if !elaborated(&Node{IsMatch: true, Kind: MatchCall}) {
		t.Fatalf("elaborated(match-call) = false")
	}

	var out []string
	if err := (&emitter{}).stmtMatch(chain, &out); err == nil ||
		!strings.Contains(err.Error(), "unelaborated") {
		t.Errorf("emit accepted an unelaborated chain: %v", err)
	}
	if _, err := evMatch(chain, map[string]*Value{}, &Ctx{}, "m__go"); err == nil ||
		!strings.Contains(err.Error(), "unelaborated") {
		t.Errorf("evaluation accepted an unelaborated chain: %v", err)
	}
	fn := &FnDecl{Name: "m__go", Body: chain}
	errs := verifyExhaustiveAll([]*Module{{Decls: []Decl{fn}}}, &Program{})
	found := false
	for _, err := range errs {
		if strings.Contains(err.Error(), "unelaborated") {
			found = true
		}
	}
	if !found {
		t.Errorf("proof accepted an unelaborated chain: %v", errs)
	}
}
