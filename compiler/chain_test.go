package main

import (
	"testing"
)

// a86: sequential fallible composition. A `match chain` block lists
// ordered call steps binding Ok payloads with one shared failure;
// the checker elaborates it into nested call matches, so admission,
// proofs, runs, coverage, emit, and CAN4107 see ordinary ladders.

const chainLib = `mod chain
  provides [chain__is_big, chain__is_odd, chain__fail_if_small, Chain__Bool, Chain__Out]
  uses []
  emits [chain.too_small]

error chain.too_small(value: int)

fn chain__is_big(value: int) -> Chain__Bool rev 1
  emits []
  tests
    big(10) => Ok(true)
    small(3) => Ok(false)
  Ok(value > 5)

fn chain__is_odd(value: int) -> Chain__Bool rev 1
  emits []
  tests
    odd(3) => Ok(true)
    even(4) => Ok(false)
  Ok(value % 2 == 1)

type Chain__Bool rev 1 (
  value: bool
)

type Chain__Out rev 1 (
  value: int
)

fn chain__fail_if_small(value: int) -> Chain__Out rev 1
  emits [chain.too_small]
  tests
    big(10) => Ok(10)
    small(3) => chain.too_small(3)
  match call chain__is_big(value)
    on Ok b => match b.value
      true => Ok(value)
      false => chain.too_small(value)
`

func chainDiags(t *testing.T, body string) []Diag {
	t.Helper()
	files := map[string]string{"chain.can": chainLib, "client.can": body}
	dir := writeLSPDir(t, files)
	return diagnose(dir, "client.can", body)
}

func chainErrs(t *testing.T, body string) []Diag {
	t.Helper()
	var errs []Diag
	for _, d := range chainDiags(t, body) {
		if d.Sev == "error" {
			errs = append(errs, d)
		}
	}
	return errs
}

// TestChainTwoSteps pins the whole pipeline on a threaded pair:
// parse the block, elaborate to nested matches, run every row,
// cover every arm, emit nothing surprising. Both outcomes pass.
func TestChainTwoSteps(t *testing.T) {
	body := `mod client
  provides [client__check, Client__Out]
  uses [chain__is_big@1, chain__is_odd@1]
  emits [chain.too_small]

type Client__Out rev 1 (
  value: int
)

fn client__check(value: int) -> Client__Out rev 1
  emits [chain.too_small]
  tests
    big_odd(11) => Ok(11)
    big_even(10) => chain.too_small(10)
    small(3) => chain.too_small(3)
  match chain
    call chain__is_big(value) as b when b.value
      given
        big_odd => exchange args (value = 11) outcome Ok(true)
        big_even => exchange args (value = 10) outcome Ok(true)
        small => exchange args (value = 3) outcome Ok(false)
    call chain__is_odd(value) as o when o.value
      given
        big_odd => exchange args (value = 11) outcome Ok(true)
        big_even => exchange args (value = 10) outcome Ok(false)
    then Ok(value)
    else chain.too_small(value)
`
	if errs := chainErrs(t, body); len(errs) != 0 {
		t.Fatalf("expected clean chain, got %v", errs)
	}
}

// TestChainFailureArmNeedsRow pins the step-row law through the
// existing test-per-arm rule: the guard-false arm with no selecting
// row is CAN4107, naming the arm, not a silent gap.
func TestChainFailureArmNeedsRow(t *testing.T) {
	body := `mod client
  provides [client__check, Client__Out]
  uses [chain__is_big@1]
  emits [chain.too_small]

type Client__Out rev 1 (
  value: int
)

fn client__check(value: int) -> Client__Out rev 1
  emits [chain.too_small]
  tests
    big(10) => Ok(10)
  match chain
    call chain__is_big(value) as b when b.value
      given
        big => exchange args (value = 10) outcome Ok(true)
    then Ok(value)
    else chain.too_small(value)
`
	var taken []Diag
	for _, d := range chainDiags(t, body) {
		if d.Code == CodeArmUntaken {
			taken = append(taken, d)
		}
		if d.Sev == "error" && d.Code != CodeArmUntaken {
			t.Fatalf("unexpected error diag %v", d)
		}
	}
	if len(taken) != 1 {
		t.Fatalf("expected exactly one CAN4107, got %v", taken)
	}
}

// TestChainEmittingCallee pins error-arm generation: a step whose
// callee emits takes on-arms per kind from the contract, each
// needing its selecting row like any handwritten ladder.
func TestChainEmittingCallee(t *testing.T) {
	body := `mod client
  provides [client__relay, Client__Out]
  uses [chain__fail_if_small@1]
  emits [chain.too_small]

type Client__Out rev 1 (
  value: int
)

fn client__relay(value: int) -> Client__Out rev 1
  emits [chain.too_small]
  tests
    big(10) => Ok(10)
    small(3) => chain.too_small(3)
  match chain
    call chain__fail_if_small(value) as r
      given
        big => exchange args (value = 10) outcome Ok(10)
        small => exchange args (value = 3) outcome chain.too_small(3)
    then Ok(r.value)
    else chain.too_small(value)
`
	if errs := chainErrs(t, body); len(errs) != 0 {
		t.Fatalf("expected clean emitting chain, got %v", errs)
	}
}

// TestChainGuardCallRejected pins the short-circuit boundary:
// guards stay call-free (slice-5 operands), so bind the value in
// an earlier step instead.
func TestChainGuardCallRejected(t *testing.T) {
	body := `mod client
  provides [client__check, Client__Out]
  uses [chain__is_big@1]
  emits [chain.too_small]

type Client__Out rev 1 (
  value: int
)

fn client__check(value: int) -> Client__Out rev 1
  emits [chain.too_small]
  tests
    big(10) => Ok(10)
  match chain
    call chain__is_big(value) as b when call chain__is_big(value = value)
    then Ok(value)
    else chain.too_small(value)
`
	errs := chainErrs(t, body)
	if len(errs) == 0 || !hasDiag(errs, "error", "cannot call") {
		t.Fatalf("expected guard-call rejection, got %v", errs)
	}
}

// TestChainDiscardBinderWithWhen pins the obvious nonsense: `as _`
// discards the value, so no `when` can observe it.
func TestChainDiscardBinderWithWhen(t *testing.T) {
	body := `mod client
  provides [client__check, Client__Out]
  uses [chain__is_big@1]
  emits [chain.too_small]

type Client__Out rev 1 (
  value: int
)

fn client__check(value: int) -> Client__Out rev 1
  emits [chain.too_small]
  tests
    big(10) => Ok(10)
  match chain
    call chain__is_big(value) as _ when true
    then Ok(value)
    else chain.too_small(value)
`
	errs := chainErrs(t, body)
	if len(errs) == 0 || !hasDiag(errs, "error", "named binder") {
		t.Fatalf("expected discard-binder rejection, got %v", errs)
	}
}

// TestChainMissingClauses pins the block shape: steps alone prove
// nothing without the success tail and the shared failure.
func TestChainMissingClauses(t *testing.T) {
	for _, tc := range []struct {
		name string
		tail string
	}{
		{"no then", "    else chain.too_small(value)\n"},
		{"no else", "    then Ok(value)\n"},
	} {
		body := `mod client
  provides [client__check, Client__Out]
  uses [chain__is_big@1]
  emits [chain.too_small]

type Client__Out rev 1 (
  value: int
)

fn client__check(value: int) -> Client__Out rev 1
  emits [chain.too_small]
  tests
    big(10) => Ok(10)
  match chain
    call chain__is_big(value) as b when b.value
` + tc.tail
		errs := chainErrs(t, body)
		if len(errs) == 0 || !hasDiag(errs, "error", "needs") {
			t.Fatalf("%s: expected missing-clause rejection, got %v", tc.name, errs)
		}
	}
}

// TestChainGivenScriptsSteps pins scripted evidence onto chain
// steps: a foreign step call carries the ordinary table grammar on
// the elaborated call node, so stub discipline, consistency, and
// leftover checks apply unchanged. (Approval-path chains stay
// script-free through the asset-guard test, not a syntax ban.)
func TestChainGivenScriptsSteps(t *testing.T) {
	body := `mod client
  provides [client__check, Client__Out]
  uses [chain__is_big@1]
  emits [chain.too_small]

type Client__Out rev 1 (
  value: int
)

fn client__check(value: int) -> Client__Out rev 1
  emits [chain.too_small]
  tests
    big(10) => Ok(10)
    small(3) => chain.too_small(3)
  match chain
    call chain__is_big(value) as b when b.value
      given
        big => exchange args (value = 10) outcome Ok(true)
        small => exchange args (value = 3) outcome Ok(false)
    then Ok(value)
    else chain.too_small(value)
`
	if errs := chainErrs(t, body); len(errs) != 0 {
		t.Fatalf("expected scripted chain to pass clean, got %v", errs)
	}
}

// TestChainNameBackwardCompat pins the contextual keywords: a
// variable named chain keeps its ordinary value match, since only
// `match chain` followed by a call step opens a chain block.
func TestChainNameBackwardCompat(t *testing.T) {
	client := `mod client
  provides [client__check, Client__Out]
  uses []
  emits []

type Client__Out rev 1 (
  value: int
)

fn client__check(chain: bool) -> Client__Out rev 1
  emits []
  tests
    yes(true) => Ok(1)
    no(false) => Ok(0)
  match chain
    true => Ok(1)
    false => Ok(0)
`
	files := map[string]string{"client.can": client}
	dir := writeLSPDir(t, files)
	if errs := diagnose(dir, "client.can", client); len(errs) != 0 {
		t.Fatalf("expected bool var named chain to match plainly, got %v", errs)
	}
}
