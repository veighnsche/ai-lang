package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// canlc lint goldens: exact report lines for all seven rules.
// Report text is pinned byte-for-byte (a88 golden-test capability);
// behavior first, then refusal shapes, then the proving-fixture
// folder, then idempotence.

const lintRedundant = `mod demo
  provides [demo__go, demo__pos, demo__id, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(left: int, right: int) -> Int__Value rev 1
  emits []
  tests
    named(left = 5, right = 8) => Ok(value = -3)
    mixed(5, right = 8) => Ok(value = -3)
    reordered(right = 8, left = 5) => Ok(value = -3)
  match call demo__id(value = left)
    on Ok r => Ok(value = r.value - right)

fn demo__pos(left: int, right: int) -> Int__Value rev 1
  emits []
  tests
    t(left = 5, right = 8) => Ok(value = -3)
  match call demo__id(left)
    on Ok r => Ok(value = r.value - right)

fn demo__id(value: int) -> Int__Value rev 1
  emits []
  tests
    t(value = 1) => Ok(value = 1)
  Ok(value = value)
`

func TestLintRedundantNames(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintRedundant})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:13: redundant argument name "left" (param 1 of demo__go is "left"); write positionally`,
		`demo.can:13: redundant argument name "right" (param 2 of demo__go is "right"); write positionally`,
		`demo.can:14: redundant argument name "right" (param 2 of demo__go is "right"); write positionally`,
		`demo.can:16: redundant argument name "value" (param 1 of demo__id is "value"); write positionally`,
		`demo.can:22: redundant argument name "left" (param 1 of demo__pos is "left"); write positionally`,
		`demo.can:22: redundant argument name "right" (param 2 of demo__pos is "right"); write positionally`,
		`demo.can:29: redundant argument name "value" (param 1 of demo__id is "value"); write positionally`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

const lintLadder = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: str) -> Int__Value rev 1
  emits []
  tests
    t(value = "a") => Ok(value = 1)
  match value
    "a" => Ok(value = 1)
    "b" => Ok(value = 1)
    "c" => Ok(value = 2)
    "d" => Ok(value = 1)
    _ => Ok(value = 0)
`

func TestLintMergeableLadder(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintLadder})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	var merges []string
	for _, f := range got {
		if strings.Contains(f.msg, "mergeable") {
			merges = append(merges, f.String())
		}
	}
	want := []string{
		`demo.can:15: mergeable match arms [15 16] (identical outcomes); fold into one or-pattern arm, saves 1 lines`,
	}
	if len(merges) != len(want) || merges[0] != want[0] {
		t.Fatalf("merges = %v, want %v (full report %v)", merges, want, got)
	}
}

// Diagonal bool tables stay tables: or binds tighter than the
// comma, so (true,true)|(false,false) cannot spell as one arm
// (a28 truth-table idiom survives a88).
const lintDiagonal = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(left: int, right: int) -> Int__Value rev 1
  emits []
  tests
    t(left = 1, right = 2) => Ok(value = 0)
  match left < 0, right < 0
    true, true => Ok(value = 1)
    true, false => Ok(value = -1)
    false, true => Ok(value = -1)
    false, false => Ok(value = 1)
`

func TestLintDiagonalTablesUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintDiagonal})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "mergeable") {
			t.Fatalf("diagonal table must not merge, got %v (full %v)", f, got)
		}
	}
}

// Variant arms never group: binders cannot be alpha-renamed
// into one, and case-or is not a spelled shape.
const lintVariant = `mod demo
  provides [demo__go, Form__State, Form__Msg]
  uses []
  emits []

variant Form__State rev 1 (
  case Empty()
  case Named(name: str)
)

type Form__Msg rev 1 (
  message: str
)

fn demo__go(state: Form__State) -> Form__Msg rev 1
  emits []
  tests
    t(state = Form__State.Empty()) => Ok(message = "none")
  match state
    on Form__Empty _ => Ok(message = "none")
    on Form__Named n => Ok(message = "none")
`

func TestLintVariantArmsUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintVariant})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "mergeable") {
			t.Fatalf("variant arms must not merge, got %v (full %v)", f, got)
		}
	}
}

// Non-adjacent same-outcome arms stay silent: folding across
// the interleaved "c" arm would steal its overlapping cells.
const lintGap = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: str) -> Int__Value rev 1
  emits []
  tests
    t(value = "a") => Ok(value = 1)
  match value
    "a" => Ok(value = 1)
    "z" => Ok(value = 9)
    "b" => Ok(value = 1)
`

func TestLintGappedPairUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintGap})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "mergeable") {
			t.Fatalf("gapped pair must not merge, got %v (full %v)", f, got)
		}
	}
}

const lintChain = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits [demo.failed]

error demo.failed(value: int)

type Int__Value rev 1 (
  value: int
)

fn demo__step(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(value = 1)
  Ok(value = value)

fn demo__go(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(value = 1)
  match call demo__step(value)
    on Ok a => match call demo__step(a.value)
      on Ok b => Ok(value = b.value)
      on demo.failed e => demo.failed(value = 0)
    on demo.failed e => demo.failed(value = 0)
`

func TestLintChainable(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintChain})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:22: chainable ladder: 2 sequential calls sharing one failure outcome; rewrite as match chain`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// Divergent failure payloads share no else: no chain, and the
// constant payloads are not relays either.
const lintChainDiverged = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits [demo.failed]

error demo.failed(value: int)

type Int__Value rev 1 (
  value: int
)

fn demo__step(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(value = 1)
  Ok(value = value)

fn demo__go(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(value = 1)
  match call demo__step(value)
    on Ok a => match call demo__step(a.value)
      on Ok b => Ok(value = b.value)
      on demo.failed e => demo.failed(value = 0)
    on demo.failed e => demo.failed(value = 1)
`

func TestLintChainDivergedUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintChainDiverged})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	if len(got) != 0 {
		t.Fatalf("divergent ladder must be clean, got %v", got)
	}
}

const lintTabled = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(a: int, b: int) -> Int__Value rev 1
  emits []
  tests
    t(1, 2) => Ok(value = 3)
  match a < 0
    true => match b < 0
      true => Ok(value = 0)
      false => Ok(value = 1)
    false => match b < 0
      true => Ok(value = 2)
      false => Ok(value = 3)
`

func TestLintTableable(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintTabled})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:14: nested matches share one scrutinee; fold into a multi-scrutinee table`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// Different inner scrutinees would need don't-care slots;
// impure inner scrutinees must keep their guard.
const lintTabledKept = `mod demo
  provides [demo__diff, demo__guarded, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__diff(a: int, b: int, c: int) -> Int__Value rev 1
  emits []
  tests
    t(1, 2, 3) => Ok(value = 6)
  match a < 0
    true => match b < 0
      true => Ok(value = 0)
      false => Ok(value = 1)
    false => match c < 0
      true => Ok(value = 2)
      false => Ok(value = 3)

fn demo__guarded(s: str) -> Int__Value rev 1
  emits []
  tests
    t("ab") => Ok(value = 1)
  match #s >= 1
    true => match s[0:1] == "a"
      true => Ok(value = 1)
      false => Ok(value = 0)
    false => Ok(value = 0)
`

func TestLintTabledKeptUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintTabledKept})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "multi-scrutinee") {
			t.Fatalf("kept nests must not fold, got %v (full %v)", f, got)
		}
	}
}

const lintSame = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: int) -> Int__Value rev 1
  emits []
  tests
    t(1) => Ok(value = 0)
  match value < 0
    true => Ok(value = 0)
    false => Ok(value = 0)
`

func TestLintSameOutcome(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintSame})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:14: match always yields Ok(value = 0); drop the match`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// A faulting scrutinee keeps its match even when every arm
// agrees: dropping it would drop the fault.
const lintSameGuarded = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(s: str) -> Int__Value rev 1
  emits []
  tests
    t("ab") => Ok(value = 0)
  match s[0:1] == "a"
    true => Ok(value = 0)
    false => Ok(value = 0)
`

func TestLintSameGuardedUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintSameGuarded})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	// The faulting scrutinee blocks the drop, but the adjacent
	// pair still folds: identical outcomes, scrutinee kept.
	want := []string{
		`demo.can:15: mergeable match arms [15 16] (identical outcomes); fold into one or-pattern arm, saves 1 lines`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

const lintRelay = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits [demo.failed, demo.other]

error demo.failed(value: int)
error demo.other(value: int, note: str)

fn demo__step(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(value = 1)
  Ok(value = value)

fn demo__go(value: int) -> Int__Value rev 1
  emits [demo.failed, demo.other]
  tests
    t(1) => Ok(value = 1)
  match call demo__step(value)
    on Ok a => Ok(value = a.value)
    on demo.failed e => demo.failed(value = e.value)
    on demo.other e => demo.other(value = e.value, note = "x")
`

func TestLintForwardable(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintRelay})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:21: handwritten relay of demo.failed; write forward e`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

const lintRanges = `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: int) -> Int__Value rev 1
  emits []
  tests
    t(2) => Ok(value = 1)
  match value
    1..3 => Ok(value = 1)
    4..6 => Ok(value = 1)
    8..9 => Ok(value = 2)
    _ => Ok(value = 0)
`

func TestLintRangeMerge(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintRanges})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	want := []string{
		`demo.can:15: mergeable ranges 1..3, 4..6; join into 1..6`,
	}
	if len(got) != len(want) {
		t.Fatalf("findings = %v, want %v", got, want)
	}
	for i := range got {
		if got[i].String() != want[i] {
			t.Fatalf("finding %d = %q, want %q", i, got[i].String(), want[i])
		}
	}
}

// Two sequential calls that cannot fail share no failure outcome:
// the router's divergent plain sides must not establish unity for
// the run (regression: the rejected router leaked failures=true).
const lintChainInfallible = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__step(value: int) -> Int__Value rev 1
  emits []
  tests
    t(1) => Ok(value = 1)
  Ok(value = value)

fn demo__go(value: int) -> Int__Value rev 1
  emits []
  tests
    t(1) => Ok(value = 1)
  match call demo__step(value)
    on Ok a => match call demo__step(a.value)
      on Ok b => match b.value == 0
        true => Ok(value = 0)
        false => Ok(value = b.value)
`

func TestLintChainInfallibleUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintChainInfallible})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	for _, f := range got {
		if strings.Contains(f.msg, "chainable") {
			t.Fatalf("infallible ladder must not chain, got %v (full %v)", f, got)
		}
	}
}

// A ladder of relays shares failure TEXT but no failure VALUE:
// each arm rebuilds from its own binder, and chain elaboration
// binds error arms with wildcards, so no shared else can spell
// it. Chainable stays silent; the relay rule reports each arm.
const lintChainRelay = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits [demo.failed]

error demo.failed(value: int)

type Int__Value rev 1 (
  value: int
)

fn demo__step(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(value = 1)
  Ok(value = value)

fn demo__go(value: int) -> Int__Value rev 1
  emits [demo.failed, demo.other]
  tests
    t(1) => Ok(value = 1)
  match call demo__step(value)
    on Ok a => match call demo__step(a.value)
      on Ok b => Ok(value = b.value)
      on demo.failed e => demo.failed(value = e.value)
    on demo.failed e => demo.failed(value = e.value)
`

// A ladder of forwards shares one spelling but no shared else:
// forward rebuilds from each arm's own binder, which the shared
// else scope cannot see. Chainable stays silent.
const lintChainForward = `mod demo
  provides [demo__go, demo__step, Int__Value]
  uses []
  emits [demo.failed]

error demo.failed(value: int)

type Int__Value rev 1 (
  value: int
)

fn demo__step(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(value = 1)
  Ok(value = value)

fn demo__go(value: int) -> Int__Value rev 1
  emits [demo.failed]
  tests
    t(1) => Ok(value = 1)
  match call demo__step(value)
    on Ok a => match call demo__step(a.value)
      on Ok b => Ok(value = b.value)
      on demo.failed e => forward e
    on demo.failed e => forward e
`

func TestLintChainForwardUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintChainForward})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	if len(got) != 0 {
		t.Fatalf("forward ladder must not chain, got %v", got)
	}
}

func TestLintChainRelayUntouched(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"demo.can": lintChainRelay})
	if len(skipped) != 0 {
		t.Fatalf("unexpected skips: %v", skipped)
	}
	relays := 0
	for _, f := range got {
		if strings.Contains(f.msg, "chainable") {
			t.Fatalf("relay ladder must not chain, got %v (full %v)", f, got)
		}
		if strings.Contains(f.msg, "handwritten relay") {
			relays++
			if f.code != CodeLintRelay {
				t.Fatalf("relay finding carries %q, want %q", f.code, CodeLintRelay)
			}
		}
	}
	if relays != 2 {
		t.Fatalf("expected 2 relay findings, got %v", got)
	}
}

func TestLintSkipsUnparseable(t *testing.T) {
	got, skipped := lintFiles(map[string]string{"bad.can": "mod demo\n  provides ["})
	if len(got) != 0 {
		t.Fatalf("expected no findings, got %v", got)
	}
	if len(skipped) != 1 || skipped[0] != "bad.can" {
		t.Fatalf("expected one skip, got %v", skipped)
	}
}

// TestLintErrorFixtures is the folder golden for
// sketches/lint-errors: every file must parse, every file must
// yield at least one finding, and the findings are pinned
// exactly — one rule family per file, nothing advisory-only.
// Lint findings are error-severity end to end: the folder run
// exits 1 with every finding on stdout.
func TestLintErrorFixtures(t *testing.T) {
	entries, err := os.ReadDir("../sketches/lint-errors")
	if err != nil {
		t.Fatalf("read fixture dir: %v", err)
	}
	files := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".can") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join("../sketches/lint-errors", e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		files[e.Name()] = string(raw)
	}
	if len(files) != 7 {
		t.Fatalf("want 7 proving fixtures, got %d: %v", len(files), files)
	}
	got, skipped := lintFiles(files)
	if len(skipped) != 0 {
		t.Fatalf("fixtures must parse, skipped: %v", skipped)
	}
	byFile := map[string][]string{}
	for _, f := range got {
		byFile[f.file] = append(byFile[f.file], f.String())
	}
	want := map[string][]string{
		"chain.can": {
			`chain.can:34: chainable ladder: 2 sequential calls sharing one failure outcome; rewrite as match chain`,
		},
		"merge.can": {
			`merge.can:23: mergeable match arms [23 24] (identical outcomes); fold into one or-pattern arm, saves 1 lines`,
		},
		"outcome.can": {
			`outcome.can:20: match always yields Ok(value = 0); drop the match`,
		},
		"ranges.can": {
			`ranges.can:22: mergeable ranges 1..3, 4..6; join into 1..6`,
		},
		"redundant.can": {
			`redundant.can:17: redundant argument name "left" (param 1 of redundant__sub is "left"); write positionally`,
			`redundant.can:17: redundant argument name "right" (param 2 of redundant__sub is "right"); write positionally`,
			`redundant.can:18: redundant argument name "right" (param 2 of redundant__sub is "right"); write positionally`,
			`redundant.can:20: redundant argument name "value" (param 1 of redundant__id is "value"); write positionally`,
			`redundant.can:26: redundant argument name "value" (param 1 of redundant__id is "value"); write positionally`,
		},
		"relay.can": {
			`relay.can:39: handwritten relay of relay.failed; write forward e`,
		},
		"table.can": {
			`table.can:22: nested matches share one scrutinee; fold into a multi-scrutinee table`,
		},
	}
	for name, lines := range want {
		if len(byFile[name]) != len(lines) {
			t.Fatalf("%s findings = %v, want %v", name, byFile[name], lines)
		}
		for i := range lines {
			if byFile[name][i] != lines[i] {
				t.Fatalf("%s finding %d = %q, want %q", name, i, byFile[name][i], lines[i])
			}
		}
	}
	for name := range byFile {
		if _, ok := want[name]; !ok {
			t.Fatalf("unexpected findings for %s: %v", name, byFile[name])
		}
	}
	var stdout, stderr bytes.Buffer
	if code := runLint(&stdout, &stderr, []string{"../sketches/lint-errors"}); code != 1 {
		t.Fatalf("folder lint exit = %d, want 1 (findings are errors): stdout %q stderr %q",
			code, stdout.String(), stderr.String())
	}
	if n := len(strings.Split(strings.TrimSpace(stdout.String()), "\n")); n != len(got) {
		t.Fatalf("folder lint printed %d lines, want %d (one per finding)", n, len(got))
	}
}

// TestLintSpanPositions pins the 2D editor contract: every
// lint-errors fixture finding publishes a [start, end) UTF-16
// span covering exactly the offending token — never the whole
// line. Columns are hand-checked against the fixture sources;
// the covered text is asserted too, so a locator that drifts
// onto the wrong token fails here, not in Cursor.
func TestLintSpanPositions(t *testing.T) {
	entries, err := os.ReadDir("../sketches/lint-errors")
	if err != nil {
		t.Fatalf("read fixture dir: %v", err)
	}
	texts := map[string]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".can") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join("../sketches/lint-errors", e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		texts[e.Name()] = string(raw)
	}
	type want struct {
		line, start, end int
		code, cover      string
	}
	wants := map[string][]want{
		"chain.can": {
			{34, 13, 31, CodeLintChain, "chain__step(value)"},
		},
		"merge.can": {
			{23, 4, 7, CodeLintOrFold, `"a"`},
		},
		"outcome.can": {
			{20, 2, 17, CodeLintSameOutcome, "match value < 0"},
		},
		"ranges.can": {
			{22, 4, 8, CodeLintRange, "1..3"},
		},
		"redundant.can": {
			{17, 12, 18, CodeLintRedundant, "left ="},
			{17, 22, 29, CodeLintRedundant, "right ="},
			{18, 13, 20, CodeLintRedundant, "right ="},
			{20, 27, 34, CodeLintRedundant, "value ="},
			{26, 8, 15, CodeLintRedundant, "value ="},
		},
		"relay.can": {
			{39, 25, 54, CodeLintRelay, "relay.failed(value = e.value)"},
		},
		"table.can": {
			{22, 2, 13, CodeLintTable, "match a < 0"},
		},
	}
	for name, ws := range wants {
		lines := strings.Split(texts[name], "\n")
		got := lintDiagsFor(name, texts)
		if len(got) != len(ws) {
			t.Fatalf("%s diags = %d, want %d", name, len(got), len(ws))
		}
		for i, w := range ws {
			d := got[i]
			if d.Line != w.line || d.Start != w.start || d.End != w.end || d.Code != w.code {
				t.Fatalf("%s diag %d = line %d [%d,%d) %s, want line %d [%d,%d) %s",
					name, i, d.Line, d.Start, d.End, d.Code, w.line, w.start, w.end, w.code)
			}
			if cover := lines[d.Line-1][d.Start:d.End]; cover != w.cover {
				t.Fatalf("%s diag %d covers %q, want %q", name, i, cover, w.cover)
			}
		}
	}
	// No finding in the folder may fall back to whole-line:
	// every lint squiggle must be column-precise.
	for name := range texts {
		for _, d := range lintDiagsFor(name, texts) {
			if d.End <= d.Start {
				t.Fatalf("%s:%d %s has no span (whole-line fallback)", name, d.Line, d.Code)
			}
		}
	}
}

// TestLintFindNameEq pins the redundant-name locator: the
// (skip+1)-th `name =` outside strings, with `=>`/`==` and
// dotted-prefix guards.
func TestLintFindNameEq(t *testing.T) {
	cases := []struct {
		code, name string
		skip       int
		ns, ne     int
		ok         bool
	}{
		{"ordered(left = 5, right = 8) => Ok(value = -3)", "left", 0, 8, 14, true},
		{"ordered(left = 5, right = 8) => Ok(value = -3)", "right", 0, 18, 25, true},
		{"one(value = 1) => Ok(value = 1)", "value", 0, 4, 11, true},
		{"one(value = 1) => Ok(value = 1)", "value", 1, 21, 28, true},
		{"relay.failed(value = e.value)", "value", 0, 13, 20, true},
		{"f(x == 1)", "x", 0, 0, 0, false},
		{"m(true => 1)", "true", 0, 0, 0, false},
		{`"left = 5"`, "left", 0, 0, 0, false},
		{"on relay.failed e => x", "failed", 0, 0, 0, false},
		{"f(note = \"a = b\")", "a", 0, 0, 0, false},
	}
	for _, c := range cases {
		ns, ne, ok := lintFindNameEq(c.code, c.name, c.skip)
		if ok != c.ok || ns != c.ns || ne != c.ne {
			t.Fatalf("lintFindNameEq(%q, %q, %d) = [%d,%d) %v, want [%d,%d) %v",
				c.code, c.name, c.skip, ns, ne, ok, c.ns, c.ne, c.ok)
		}
	}
}

// TestLintDiagnoseErrors pins the editor contract: lint findings
// publish through diagnose as error-severity diagnostics with
// registered CAN3410-3416 codes — this is what Cursor draws.
// Clean-world layering holds: a file with real compiler errors
// shows those first, with lint suppressed until they are fixed.
func TestLintDiagnoseErrors(t *testing.T) {
	src := `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: str) -> Int__Value rev 1
  emits []
  tests
    t("a") => Ok(value = 1)
    b("b") => Ok(value = 1)
    c("c") => Ok(value = 2)
    z("z") => Ok(value = 0)
  match value
    "a" => Ok(value = 1)
    "b" => Ok(value = 1)
    "c" => Ok(value = 2)
    _ => Ok(value = 0)
`
	dir := writeLSPDir(t, map[string]string{"demo.can": src})
	diags := diagnose(dir, "demo.can", src)
	var found *Diag
	for i, d := range diags {
		if d.Code == CodeLintOrFold {
			found = &diags[i]
		}
	}
	if found == nil {
		t.Fatalf("expected CAN3411 or-fold diag, got %v", diags)
	}
	if found.Sev != "error" {
		t.Fatalf("lint diag severity = %q, want error", found.Sev)
	}
	if found.Line != 18 {
		t.Fatalf("lint diag line = %d, want 18", found.Line)
	}
}

func TestLintDiagnoseSuppressedOnBroken(t *testing.T) {
	src := `mod demo
  provides [demo__go, Int__Value]
  uses []
  emits []

type Int__Value rev 1 (
  value: int
)

fn demo__go(value: str) -> Int__Value rev 1
  emits []
  tests
    t("a") => Ok(value = 1)
  match value
    "a" => Ok(value = 1)
    "b" => Ok(value = 1)
`
	dir := writeLSPDir(t, map[string]string{"demo.can": src})
	diags := diagnose(dir, "demo.can", src)
	for _, d := range diags {
		if strings.HasPrefix(d.Code, "CAN341") {
			t.Fatalf("broken file must show real errors first, got lint %v (full %v)", d, diags)
		}
	}
	if len(diags) == 0 {
		t.Fatalf("expected compiler errors for missing wildcard arm, got none")
	}
}

func TestLintIdempotent(t *testing.T) {
	files := map[string]string{"a.can": lintRedundant, "b.can": lintLadder}
	first, _ := lintFiles(files)
	second, _ := lintFiles(files)
	if len(first) != len(second) {
		t.Fatalf("lengths differ: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("run %d differs at %d", 2, i)
		}
	}
}

// TestLintStdClean pins the stdlib migration (can-idioms C13):
// every blessed std module lints clean — no redundant names,
// no unfolded ladders, no handwritten relays. All files lint
// together so cross-module signatures resolve as in
// `canlc lint std/`.
func TestLintStdClean(t *testing.T) {
	mods := []string{"ascii", "division", "html", "quota", "ratio", "scalars", "schema", "text"}
	files := map[string]string{}
	for _, m := range mods {
		name := m + ".can"
		raw, err := os.ReadFile(filepath.Join("../std", m, name))
		if err != nil {
			t.Fatalf("read std/%s/%s: %v", m, name, err)
		}
		files[name] = string(raw)
	}
	got, skipped := lintFiles(files)
	if len(skipped) != 0 {
		t.Fatalf("std must parse, skipped: %v", skipped)
	}
	if len(got) != 0 {
		var msgs []string
		for _, f := range got {
			msgs = append(msgs, f.String())
		}
		t.Fatalf("std must lint clean, got %d findings:\n%s", len(got), strings.Join(msgs, "\n"))
	}
}
