package main

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ---------------------------------------------------------------- AST ------
type Arg struct {
	Name    string
	HasName bool
	V       *Small
}

type Small struct {
	Kind string // str,int,bool,dec,float,wild,binop,call,ctor,list,ref,seal,exchange,strlen,stridx,strslice
	Str  string
	// Outcome holds a scripted result for Kind exchange: the row proves
	// "this request received this permitted response" (v12).
	Outcome *Small
	// Num holds an int literal of arbitrary size (v10: ints are
	// mathematically unbounded, so literals never overflow).
	Num *big.Int
	// T holds the checker's static type for ref/binop/seal nodes
	// (v10: emit reads it to choose bigint-native vs exact-decimal
	// code, so code generation reasons over typed structure).
	T string
	B bool
	// Dec holds canonical decimal digits for Kind dec: -?\d+\.\d+ with
	// no trailing fractional zeros (d"1.50" parses to "1.5").
	Dec string
	// Seal holds the brand name for Kind seal; Str holds the literal.
	Seal  string
	Op    string
	L, R  *Small
	// Hi holds the slice end for Kind strslice (base L, start R).
	Hi *Small
	Fname string
	Args  []Arg
	Ctor  string
	Items []*Small
	Ref   []string
}

type Pattern struct {
	Kind string // wild,bool,str,variant,variantWild
	B    bool
	Str  string
	Name string
	Var  string
}

type Arm struct {
	Pat  Pattern
	Rhs  *Node
	Line int
}

type Node struct {
	IsMatch bool
	Scrut   *Small
	Arms    []Arm
	Given   map[string]*Small // nil value node = "-" unreachable
	Small   *Small
	Line    int
}

type Test struct {
	Name     string
	Args     []Arg
	Expected *Small
	Line     int
}

type Decl interface{ declKind() string }

type ErrorDecl struct {
	Name   string
	Fields [][2]string
	Line   int
}

func (d *ErrorDecl) declKind() string { return "error" }

type TypeDecl struct {
	Name   string
	Rev    int
	Fields [][2]string
	Line   int
}

func (d *TypeDecl) declKind() string { return "type" }

type FnDecl struct {
	Name     string
	Rev      int
	Params   [][2]string
	Ret      string
	Emits    []string
	Tests    []Test
	Body     *Node
	UsesHere []string
	// Effects lists the cell capabilities this function may use,
	// e.g. Count__total.read. Set from the effects metadata line.
	Effects []string
	// DecNames names the params proven to shrink on every self-call
	// (empty, none), and DecSchema names the admitted recursion shape:
	// "" is the unit loop (site passes p - 1), "euclid" is the
	// Euclidean step (site passes (b, a % b)), "narrowing" is binary
	// search (site passes (lo, mid) or (mid, hi) with mid (lo+hi)/2).
	// Set from the decreases metadata line; v19 owns the theorems.
	DecNames  []string
	DecSchema string
	Line      int
}

func (d *FnDecl) declKind() string { return "fn" }

// BrandDecl is a nominal string wrapper: brand Name is str rev N.
// An optional seals_from [B, ...] clause authorizes explicit one-way
// promotion seals from those same-module brands (v26); without it the
// brand mints from str only. Branding is proof, not runtime; the
// emitter forgets every brand.
type BrandDecl struct {
	Name      string
	Under     string
	Rev       int
	SealsFrom []string
	Line      int
}

func (d *BrandDecl) declKind() string { return "brand" }

// ExternDecl is a foreign function: declared, never defined. Calls to it
// are scripted through given tables like ail calls; the host provides
// the implementation and the TS emit imports it.
type ExternDecl struct {
	Name   string
	Rev    int
	Params [][2]string
	Ret    string
	Emits  []string
	Line   int
}

func (d *ExternDecl) declKind() string { return "extern" }

// StateDecl is module-private named storage for one base-type value:
// state Name: T = lit. Cells never appear in provides, take no pins,
// and are visible only in their own file.
type StateDecl struct {
	Name string
	Type string
	Init *Small
	Line int
}

func (d *StateDecl) declKind() string { return "state" }

type Module struct {
	File string
	// ID is the canonical source identity: the cleaned input path as
	// passed. Two inputs with different IDs are different modules even
	// when their basenames (File) match; File stays the display name.
	ID    string
	Stem  string
	Mod   string
	Hdr   map[string][]string
	Decls []Decl
}

// ---------------------------------------------------------- scanning -------
func stripComment(line string) string {
	var out strings.Builder
	inStr := false
	for i := 0; i < len(line); {
		ch := line[i]
		if inStr {
			out.WriteByte(ch)
			if ch == '\\' && i+1 < len(line) {
				out.WriteByte(line[i+1])
				i++
			} else if ch == '"' {
				inStr = false
			}
		} else {
			if ch == '"' {
				inStr = true
				out.WriteByte(ch)
			} else if ch == '/' && i+1 < len(line) && line[i+1] == '/' {
				break
			} else {
				out.WriteByte(ch)
			}
		}
		i++
	}
	return out.String()
}

func splitTop(s string, sep rune) []string {
	var parts []string
	var cur strings.Builder
	var stack []byte
	pairs := map[byte]byte{'(': ')', '[': ']'}
	inStr, esc := false, false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inStr {
			cur.WriteByte(ch)
			if esc {
				esc = false
			} else if ch == '\\' {
				esc = true
			} else if ch == '"' {
				inStr = false
			}
		} else if ch == '"' {
			inStr = true
			cur.WriteByte(ch)
		} else if closer, ok := pairs[ch]; ok {
			stack = append(stack, closer)
			cur.WriteByte(ch)
		} else if len(stack) > 0 && ch == stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
			cur.WriteByte(ch)
		} else if len(stack) == 0 && rune(ch) == sep {
			parts = append(parts, cur.String())
			cur.Reset()
		} else {
			cur.WriteByte(ch)
		}
	}
	parts = append(parts, cur.String())
	var keep []string
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			keep = append(keep, strings.TrimSpace(p))
		}
	}
	return keep
}

// findTop returns the index and matched op of the first top-level occurrence.
// topBracket finds the first [ outside strings and any paren or
// bracket depth: the start of a postfix index/slice group. Index 0
// is never reported, so list literals keep their own branch.
func topBracket(s string) (int, bool) {
	pdepth, bdepth := 0, 0
	inStr, esc := false, false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inStr {
			if esc {
				esc = false
			} else if ch == '\\' {
				esc = true
			} else if ch == '"' {
				inStr = false
			}
		} else if ch == '"' {
			inStr = true
		} else if ch == '(' {
			pdepth++
		} else if ch == ')' && pdepth > 0 {
			pdepth--
		} else if ch == '[' {
			if pdepth == 0 && bdepth == 0 {
				if i > 0 {
					return i, true
				}
				return -1, false
			}
			bdepth++
		} else if ch == ']' && bdepth > 0 {
			bdepth--
		}
	}
	return -1, false
}

// topColons lists every : outside strings and any paren or bracket
// depth. Empty sides are significant (s[1:] keeps them); splitTop
// drops empties, which would silently turn s[1:] into s[1].
func topColons(s string) []int {
	var out []int
	pdepth, bdepth := 0, 0
	inStr, esc := false, false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inStr {
			if esc {
				esc = false
			} else if ch == '\\' {
				esc = true
			} else if ch == '"' {
				inStr = false
			}
		} else if ch == '"' {
			inStr = true
		} else if ch == '(' {
			pdepth++
		} else if ch == ')' && pdepth > 0 {
			pdepth--
		} else if ch == '[' {
			bdepth++
		} else if ch == ']' && bdepth > 0 {
			bdepth--
		} else if ch == ':' && pdepth == 0 && bdepth == 0 {
			out = append(out, i)
		}
	}
	return out
}

// leadingBinOpLen reports the length of a binary operator opening s:
// two-char comparisons first, then the arithmetic ops. Zero means s
// does not start with one (= alone is binding syntax, not an op).
func leadingBinOpLen(s string) int {
	for _, op := range []string{"==", ">=", "<=", "!=", ">", "<", "+", "-", "*", "/", "%"} {
		if strings.HasPrefix(s, op) {
			return len(op)
		}
	}
	return 0
}

func findTop(s string, ops []string) (int, string) {
	var stack []byte
	pairs := map[byte]byte{'(': ')', '[': ']'}
	inStr, esc := false, false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inStr {
			if esc {
				esc = false
			} else if ch == '\\' {
				esc = true
			} else if ch == '"' {
				inStr = false
			}
		} else if ch == '"' {
			inStr = true
		} else if closer, ok := pairs[ch]; ok {
			stack = append(stack, closer)
		} else if len(stack) > 0 && ch == stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
		} else if len(stack) == 0 {
			for _, op := range ops {
				if strings.HasPrefix(s[i:], op) {
					return i, op
				}
			}
		}
	}
	return -1, ""
}

// findLastTop is findTop keeping the last top-level occurrence instead
// of the first: splitting there makes chains associate left.
func findLastTop(s string, ops []string) (int, string) {
	var stack []byte
	pairs := map[byte]byte{'(': ')', '[': ']'}
	inStr, esc := false, false
	best, bestOp := -1, ""
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inStr {
			if esc {
				esc = false
			} else if ch == '\\' {
				esc = true
			} else if ch == '"' {
				inStr = false
			}
		} else if ch == '"' {
			inStr = true
		} else if closer, ok := pairs[ch]; ok {
			stack = append(stack, closer)
		} else if len(stack) > 0 && ch == stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
		} else if len(stack) == 0 {
			for _, op := range ops {
				if strings.HasPrefix(s[i:], op) {
					best, bestOp = i, op
				}
			}
		}
	}
	return best, bestOp
}

func balanced(s string, openI int) (int, error) {
	closers := map[byte]byte{'(': ')', '[': ']'}
	want := closers[s[openI]]
	depth := 0
	inStr, esc := false, false
	for i := openI; i < len(s); i++ {
		ch := s[i]
		if inStr {
			if esc {
				esc = false
			} else if ch == '\\' {
				esc = true
			} else if ch == '"' {
				inStr = false
			}
		} else if ch == '"' {
			inStr = true
		} else if ch == s[openI] {
			depth++
		} else if ch == want {
			depth--
			if depth == 0 {
				return i, nil
			}
		}
	}
	return -1, fmt.Errorf("unbalanced %c in: %s", s[openI], s)
}

// ------------------------------------------------------- small exprs -------
var (
	reInt   = regexp.MustCompile(`^-?\d+$`)
	reWord  = regexp.MustCompile(`^[\w.]+$`)
	reName  = regexp.MustCompile(`^\w+$`)
	reDec   = regexp.MustCompile(`^d"([^"]*)"$`)
	reFloat = regexp.MustCompile(`^-?(\d+\.\d*|\.\d+|\d+[eE][+-]?\d+)$`)
	reDecNM = regexp.MustCompile(`^(-?)(\d+)\.(\d+)$`)
	reSeal  = regexp.MustCompile(`^seal\s+(\w+)\((.*)\)$`)
)

// canonDec normalizes dec digits to canonical form: no leading integer
// zeros, no trailing fractional zeros, -0 folded to 0. The shape stays
// dotted (d"1.0", never d"1") so every dec reads as fractional.
func canonDec(raw string) (string, error) {
	m := reDecNM.FindStringSubmatch(raw)
	if m == nil {
		return "", fmt.Errorf("bad dec digits %q: want d\"12.34\"", raw)
	}
	ip := strings.TrimLeft(m[2], "0")
	if ip == "" {
		ip = "0"
	}
	fp := strings.TrimRight(m[3], "0")
	if fp == "" {
		fp = "0"
	}
	if ip == "0" && fp == "0" {
		return "0.0", nil
	}
	return m[1] + ip + "." + fp, nil
}

func parseSmall(s string) (*Small, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty expression")
	}
	if strings.HasPrefix(s, `"`) {
		// A trailing index/slice group belongs to the postfix branch
		// below ("ABC"[0:1]), not to the literal: fall through when
		// the closing quote is followed by [. A literal-left binary
		// expression ("-" + tail) falls through too, unless the right
		// operand is itself quoted ("a" + "b" stays one swallowed
		// literal, as ever). Anything else keeps the old behavior.
		j := -1
		esc := false
		for k := 1; k < len(s); k++ {
			c := s[k]
			if esc {
				esc = false
			} else if c == '\\' {
				esc = true
			} else if c == '"' {
				j = k
				break
			}
		}
		if j == len(s)-1 {
			return &Small{Kind: "str", Str: s[1:j]}, nil
		}
		oldPath := true
		if j > 0 {
			rest := strings.TrimSpace(s[j+1:])
			if oplen := leadingBinOpLen(rest); oplen > 0 {
				after := strings.TrimSpace(rest[oplen:])
				if after != "" && !strings.HasPrefix(after, `"`) {
					oldPath = false
				}
			} else if strings.HasPrefix(rest, "[") {
				oldPath = false
			}
		}
		if oldPath {
			if len(s) < 2 || !strings.HasSuffix(s, `"`) {
				return nil, fmt.Errorf("bad string: %s", s)
			}
			return &Small{Kind: "str", Str: s[1 : len(s)-1]}, nil
		}
	}
	if reInt.MatchString(s) {
		n, ok := new(big.Int).SetString(s, 10)
		if !ok {
			return nil, fmt.Errorf("bad int literal %s", s)
		}
		return &Small{Kind: "int", Num: n}, nil
	}
	// Dec precedes float: d"1.5" is the one legal dotted spelling.
	if m := reDec.FindStringSubmatch(s); m != nil {
		canon, err := canonDec(m[1])
		if err != nil {
			return nil, err
		}
		return &Small{Kind: "dec", Dec: canon}, nil
	}
	// Bare dotted numbers parse (as float nodes) so checkStatic can
	// point at them with AIL6001; they never evaluate.
	if reFloat.MatchString(s) {
		return &Small{Kind: "float", Str: s}, nil
	}
	if s == "true" || s == "false" {
		return &Small{Kind: "bool", B: s == "true"}, nil
	}
	if s == "_" {
		return &Small{Kind: "wild"}, nil
	}
	// Seal precedes binops so a literal containing == stays intact.
	// The inner expression is kept raw: checkTypes enforces the
	// string-literal rule with AIL6003, and eval enforces str.
	if m := reSeal.FindStringSubmatch(s); m != nil {
		if idx := strings.Index(s, "("); idx >= 0 {
			if end, err := balanced(s, idx); err == nil && end == len(s)-1 {
				inner, err := parseSmall(m[2])
				if err != nil {
					return nil, err
				}
				return &Small{Kind: "seal", Seal: m[1], Args: []Arg{{V: inner}}}, nil
			}
		}
	}
	// Exchange precedes binops for the same reason: args bind with =
	// and outcomes may contain comparisons. One spelling per meaning:
	// every script row is `exchange args (...) outcome ...` (v12).
	if s == "exchange" || strings.HasPrefix(s, "exchange ") || strings.HasPrefix(s, "exchange\t") {
		return parseExchange(s)
	}
	if i, op := findTop(s, []string{"==", ">=", "<=", ">", "<", "!="}); i >= 0 {
		l, err := parseSmall(s[:i])
		if err != nil {
			return nil, err
		}
		r, err := parseSmall(s[i+len(op):])
		if err != nil {
			return nil, err
		}
		return &Small{Kind: "binop", Op: op, L: l, R: r}, nil
	}
	// Arithmetic binds tighter than comparisons. + and - split before
	// *, /, % (lower precedence splits first); each level splits at
	// the LAST top-level occurrence so chains associate left:
	// 10 - 3 - 2 is (10-3)-2. Unary minus exists on literals only
	// (reInt above). / and % share * precedence (v17: exact
	// Euclidean integer division; dec operands refused in checkSem).
	if i, op := findLastTop(s, []string{"+", "-"}); i > 0 {
		l, err := parseSmall(s[:i])
		if err != nil {
			return nil, err
		}
		r, err := parseSmall(s[i+len(op):])
		if err != nil {
			return nil, err
		}
		return &Small{Kind: "binop", Op: op, L: l, R: r}, nil
	}
	if i, op := findLastTop(s, []string{"*", "/", "%"}); i > 0 {
		l, err := parseSmall(s[:i])
		if err != nil {
			return nil, err
		}
		r, err := parseSmall(s[i+len(op):])
		if err != nil {
			return nil, err
		}
		return &Small{Kind: "binop", Op: op, L: l, R: r}, nil
	}
	// v20: scalar text operators. # binds tightest (this branch runs
	// only when no looser split matched, so the operand is atomic);
	// s[i] and s[a:b] are postfix at the same level, chaining left.
	if strings.HasPrefix(s, "#") {
		v, err := parseSmall(s[1:])
		if err != nil {
			return nil, err
		}
		return &Small{Kind: "strlen", L: v}, nil
	}
	if i, ok := topBracket(s); ok {
		base, err := parseSmall(s[:i])
		if err != nil {
			return nil, err
		}
		rest := s[i:]
		for len(rest) > 0 {
			if rest[0] != '[' {
				return nil, fmt.Errorf("unexpected %q after ]", rest)
			}
			end, err := balanced(rest, 0)
			if err != nil {
				return nil, err
			}
			inner := rest[1:end]
			colons := topColons(inner)
			if len(colons) > 1 {
				return nil, fmt.Errorf("bad slice: %s", rest[:end+1])
			}
			if len(colons) == 0 {
				ix, err := parseSmall(inner)
				if err != nil {
					return nil, err
				}
				base = &Small{Kind: "stridx", L: base, R: ix}
				rest = strings.TrimSpace(rest[end+1:])
				continue
			}
			lo, err := parseSmall(inner[:colons[0]])
			if err != nil {
				return nil, err
			}
			hi, err := parseSmall(inner[colons[0]+1:])
			if err != nil {
				return nil, err
			}
			base = &Small{Kind: "strslice", L: base, R: lo, Hi: hi}
			rest = strings.TrimSpace(rest[end+1:])
		}
		return base, nil
	}
	if strings.HasPrefix(s, "call ") {
		m := regexp.MustCompile(`^call\s+(\w+)\((.*)\)$`).FindStringSubmatch(s)
		if m == nil {
			return nil, fmt.Errorf("bad call syntax: %s", s)
		}
		args, err := parseArgs(m[2])
		if err != nil {
			return nil, err
		}
		return &Small{Kind: "call", Fname: m[1], Args: args}, nil
	}
	if strings.HasPrefix(s, "(") {
		if end, err := balanced(s, 0); err == nil && end == len(s)-1 {
			return parseSmall(s[1 : len(s)-1])
		}
	}
	if strings.HasPrefix(s, "[") {
		if end, err := balanced(s, 0); err == nil && end == len(s)-1 {
			var items []*Small
			for _, p := range splitTop(s[1:len(s)-1], ',') {
				it, err := parseSmall(p)
				if err != nil {
					return nil, err
				}
				items = append(items, it)
			}
			return &Small{Kind: "list", Items: items}, nil
		}
	}
	if m := regexp.MustCompile(`^([\w.]+)\((.*)\)$`).FindStringSubmatch(s); m != nil {
		if idx := strings.Index(s, "("); idx >= 0 {
			if end, err := balanced(s, idx); err == nil && end == len(s)-1 {
				args, err := parseArgs(m[2])
				if err != nil {
					return nil, err
				}
				return &Small{Kind: "ctor", Ctor: m[1], Args: args}, nil
			}
		}
	}
	if reWord.MatchString(s) {
		return &Small{Kind: "ref", Ref: strings.Split(s, ".")}, nil
	}
	return nil, fmt.Errorf("cannot parse expression: %s", s)
}

// parseExchange parses one script row: exchange args (...) outcome ....
// The args group is located by balanced parens (depth- and
// string-aware like every other grouping rule); arg values are full
// expressions. Expected args must be named: an unnamed expectation
// cannot say which parameter it pins.
func parseExchange(s string) (*Small, error) {
	rest := strings.TrimSpace(s[len("exchange"):])
	if !strings.HasPrefix(rest, "args") {
		return nil, fmt.Errorf("bad exchange row (want exchange args (...) outcome ...): %s", s)
	}
	rest = strings.TrimSpace(rest[len("args"):])
	if !strings.HasPrefix(rest, "(") {
		return nil, fmt.Errorf("bad exchange row (want exchange args (...) outcome ...): %s", s)
	}
	end, err := balanced(rest, 0)
	if err != nil {
		return nil, err
	}
	args, err := parseArgs(rest[1:end])
	if err != nil {
		return nil, err
	}
	for _, a := range args {
		if !a.HasName {
			return nil, fmt.Errorf("exchange args must be named (k = v): %s", s)
		}
	}
	rest = strings.TrimSpace(rest[end+1:])
	if !strings.HasPrefix(rest, "outcome") {
		return nil, fmt.Errorf("bad exchange row (want exchange args (...) outcome ...): %s", s)
	}
	rest = strings.TrimSpace(rest[len("outcome"):])
	if rest == "" {
		return nil, fmt.Errorf("bad exchange row (missing outcome): %s", s)
	}
	out, err := parseSmall(rest)
	if err != nil {
		return nil, err
	}
	return &Small{Kind: "exchange", Args: args, Outcome: out}, nil
}

func parseArgs(s string) ([]Arg, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	var out []Arg
	for _, part := range splitTop(s, ',') {
		i, op := findTop(part, []string{"==", ">=", "<=", "!=", "="})
		if op == "=" {
			name := strings.TrimSpace(part[:i])
			if !reName.MatchString(name) {
				return nil, fmt.Errorf("bad kwarg: %s", part)
			}
			v, err := parseSmall(part[i+1:])
			if err != nil {
				return nil, err
			}
			out = append(out, Arg{Name: name, HasName: true, V: v})
		} else {
			v, err := parseSmall(part)
			if err != nil {
				return nil, err
			}
			out = append(out, Arg{V: v})
		}
	}
	return out, nil
}

func isKwargList(args []Arg) bool {
	for _, a := range args {
		if !a.HasName {
			return false
		}
	}
	return true
}

// ------------------------------------------------------------ files --------
type row struct {
	indent int
	code   string
	line   int // 1-based source line
}

// LineError carries a 1-based source line for editor diagnostics
// (and `file:line:` CLI messages). at() attaches it once, innermost first.
type LineError struct {
	Line int
	Err  error
}

func (e *LineError) Error() string { return fmt.Sprintf("line %d: %v", e.Line, e.Err) }
func (e *LineError) Unwrap() error { return e.Err }

func at(line int, err error) error {
	if err == nil {
		return nil
	}
	var le *LineError
	if errors.As(err, &le) {
		return err
	}
	return &LineError{Line: line, Err: err}
}

var (
	reHdrLine   = regexp.MustCompile(`^(provides|uses|emits)\s*\[(.*)\]$`)
	reError     = regexp.MustCompile(`^error\s+([\w.]+)\((.*)\)$`)
	reType      = regexp.MustCompile(`^type\s+(\w+)\s+rev\s+(\d+)\s*\($`)
	reBrand     = regexp.MustCompile(`^brand\s+(\w+)\s+is\s+(\w+)\s+rev\s+(\d+)(\s+seals_from\s+\[([^\]]*)\])?$`)
	reExtern    = regexp.MustCompile(`^extern\s+(\w+)\((.*)\)\s*->\s*(\w+)\s+rev\s+(\d+)$`)
	reFn        = regexp.MustCompile(`^fn\s+(\w+)\((.*)\)\s*->\s*(\w+)\s+rev\s+(\d+)$`)
	reField     = regexp.MustCompile(`^(\w+)\s*:\s*(\w+)$`)
	reTest      = regexp.MustCompile(`^(\w+)\((.*)\)\s*=>\s*(.+)$`)
	reGiven     = regexp.MustCompile(`^(\w+)\s*=>\s*(.+)$`)
	reArm       = regexp.MustCompile(`^(?:on\s+)?(.+?)\s*=>\s*(.*)$`)
	reDecreases       = regexp.MustCompile(`^decreases\s+(\w+)$`)
	reDecreasesSchema = regexp.MustCompile(`^decreases\s+(\w+)\s*,\s*(\w+)\s+by\s+(euclid|narrowing)$`)
	reEffects   = regexp.MustCompile(`^effects\s*\[(.*)\]$`)
	reState     = regexp.MustCompile(`^state\s+(\w+)\s*:\s*(\w+)\s*=\s*(.+)$`)
	reRevWord   = regexp.MustCompile(`\brev\b`)
	rePatVar    = regexp.MustCompile(`^([\w.]+)\s+(\w+)$`)
	rePatWild   = regexp.MustCompile(`^([\w.]+)\s+_$`)
)

func parseFields(s, what string) ([][2]string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	var out [][2]string
	for _, part := range splitTop(s, ',') {
		m := reField.FindStringSubmatch(part)
		if m == nil {
			return nil, fmt.Errorf("bad %s field: %s", what, part)
		}
		out = append(out, [2]string{m[1], m[2]})
	}
	return out, nil
}

func parseModule(path string) (*Module, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	base := path
	if i := strings.LastIndex(path, "/"); i >= 0 {
		base = path[i+1:]
	}
	m, err := parseModuleText(base, string(data))
	if err != nil {
		var le *LineError
		if errors.As(err, &le) {
			return nil, fmt.Errorf("%s:%d: %v", path, le.Line, le.Err)
		}
		return nil, fmt.Errorf("%s: %v", path, err)
	}
	return m, nil
}

func parseModuleText(name, text string) (*Module, error) {
	path := name
	// Row 4: source decoding enforces valid UTF-8 before anything
	// else runs. The pipeline cannot carry malformed bytes
	// faithfully (the emitter substitutes), so fail closed at the
	// door rather than repairing silently downstream.
	if !utf8.ValidString(text) {
		return nil, at(1, fmt.Errorf("source is not valid UTF-8: decode the file as UTF-8 before compiling"))
	}
	if strings.ContainsAny(text, "{}") {
		line := 1
		for n, raw := range strings.Split(text, "\n") {
			if strings.ContainsAny(raw, "{}") {
				line = n + 1
				break
			}
		}
		return nil, at(line, fmt.Errorf("curly braces are banned, use () records"))
	}
	var rows []row
	for n, raw := range strings.Split(text, "\n") {
		code := strings.TrimRight(stripComment(raw), " \t")
		if strings.TrimSpace(code) == "" {
			continue
		}
		indent := len(code) - len(strings.TrimLeft(code, " "))
		rows = append(rows, row{indent, strings.TrimSpace(code), n + 1})
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("%s: empty file", path)
	}
	parts := strings.Fields(rows[0].code)
	if len(parts) != 2 || parts[0] != "mod" {
		return nil, at(rows[0].line, fmt.Errorf("file must open with mod <domain>"))
	}
	mod := &Module{Hdr: map[string][]string{}}
	mod.Mod = parts[1]
	base := path[strings.LastIndex(path, "/")+1:]
	mod.Stem = strings.TrimSuffix(base, ".ail")
	i := 1
	for i < len(rows) && rows[i].indent > 0 {
		m := reHdrLine.FindStringSubmatch(rows[i].code)
		if m == nil {
			return nil, at(rows[i].line, fmt.Errorf("bad mod header line: %s", rows[i].code))
		}
		mod.Hdr[m[1]] = splitTop(m[2], ',')
		i++
	}
	for _, key := range []string{"provides", "uses", "emits"} {
		if _, ok := mod.Hdr[key]; !ok {
			mod.Hdr[key] = nil
		}
	}
	for i < len(rows) {
		indent, code := rows[i].indent, rows[i].code
		if indent != 0 {
			return nil, at(rows[i].line, fmt.Errorf("top-level decl must start at column 0: %s", code))
		}
		declLine := rows[i].line
		switch {
		case strings.HasPrefix(code, "error "):
			m := reError.FindStringSubmatch(code)
			if m == nil {
				return nil, at(declLine, fmt.Errorf("bad error decl: %s", code))
			}
			fields, err := parseFields(m[2], "error")
			if err != nil {
				return nil, at(declLine, err)
			}
			mod.Decls = append(mod.Decls, &ErrorDecl{Name: m[1], Fields: fields, Line: declLine})
			i++
		case strings.HasPrefix(code, "type "):
			m := reType.FindStringSubmatch(code)
			if m == nil {
				if !reRevWord.MatchString(code) {
					return nil, at(declLine, fmt.Errorf("missing rev N: versioning is mandatory"))
				}
				return nil, at(declLine, fmt.Errorf("bad type decl: %s", code))
			}
			rev, _ := strconv.Atoi(m[2])
			decl := &TypeDecl{Name: m[1], Rev: rev, Line: declLine}
			i++
			for i < len(rows) && rows[i].indent > 0 {
				fs, err := parseFields(strings.TrimSuffix(rows[i].code, ","), "type")
				if err != nil {
					return nil, at(rows[i].line, err)
				}
				decl.Fields = append(decl.Fields, fs...)
				i++
			}
			if i >= len(rows) || rows[i].code != ")" {
				return nil, at(declLine, fmt.Errorf("type %s missing closing )", decl.Name))
			}
			i++
			mod.Decls = append(mod.Decls, decl)
		case strings.HasPrefix(code, "brand "):
			m := reBrand.FindStringSubmatch(code)
			if m == nil {
				if !reRevWord.MatchString(code) {
					return nil, at(declLine, fmt.Errorf("missing rev N: versioning is mandatory"))
				}
				return nil, at(declLine, fmt.Errorf("bad brand decl: %s", code))
			}
			rev, _ := strconv.Atoi(m[3])
			var from []string
			if m[4] != "" {
				for _, name := range strings.Split(m[5], ",") {
					if name = strings.TrimSpace(name); name != "" {
						from = append(from, name)
					}
				}
			}
			mod.Decls = append(mod.Decls, &BrandDecl{Name: m[1], Under: m[2], Rev: rev, SealsFrom: from, Line: declLine})
			i++
		case strings.HasPrefix(code, "state "):
			m := reState.FindStringSubmatch(code)
			if m == nil {
				return nil, at(declLine, fmt.Errorf("bad state decl: %s", code))
			}
			init, err := parseSmall(m[3])
			if err != nil {
				return nil, at(declLine, err)
			}
			mod.Decls = append(mod.Decls, &StateDecl{Name: m[1], Type: m[2], Init: init, Line: declLine})
			i++
		case strings.HasPrefix(code, "extern "):
			m := reExtern.FindStringSubmatch(code)
			if m == nil {
				if !reRevWord.MatchString(code) {
					return nil, at(declLine, fmt.Errorf("missing rev N: versioning is mandatory"))
				}
				return nil, at(declLine, fmt.Errorf("bad extern decl: %s", code))
			}
			rev, _ := strconv.Atoi(m[4])
			params, err := parseFields(m[2], "param")
			if err != nil {
				return nil, at(declLine, err)
			}
			ex := &ExternDecl{Name: m[1], Rev: rev, Params: params, Ret: m[3], Line: declLine}
			i++
			for i < len(rows) && rows[i].indent > 0 {
				mm := regexp.MustCompile(`^emits\s*\[(.*)\]$`).FindStringSubmatch(rows[i].code)
				if mm == nil {
					return nil, at(rows[i].line, fmt.Errorf("unexpected in extern %s: %s", ex.Name, rows[i].code))
				}
				ex.Emits = splitTop(mm[1], ',')
				i++
			}
			mod.Decls = append(mod.Decls, ex)
		case strings.HasPrefix(code, "fn "):
			m := reFn.FindStringSubmatch(code)
			if m == nil {
				if !reRevWord.MatchString(code) {
					return nil, at(declLine, fmt.Errorf("missing rev N: versioning is mandatory"))
				}
				return nil, at(declLine, fmt.Errorf("bad fn decl: %s", code))
			}
			rev, _ := strconv.Atoi(m[4])
			params, err := parseFields(m[2], "param")
			if err != nil {
				return nil, at(declLine, err)
			}
			fn := &FnDecl{Name: m[1], Rev: rev, Params: params, Ret: m[3], Line: declLine}
			i++
			for i < len(rows) && rows[i].indent > 0 {
				ind, c := rows[i].indent, rows[i].code
				metaLine := rows[i].line
				switch {
				case strings.HasPrefix(c, "emits "):
					mm := regexp.MustCompile(`^emits\s*\[(.*)\]$`).FindStringSubmatch(c)
					if mm == nil {
						return nil, at(metaLine, fmt.Errorf("bad emits: %s", c))
					}
					fn.Emits = splitTop(mm[1], ',')
					i++
				case strings.HasPrefix(c, "decreases"):
					if fn.DecNames != nil {
						return nil, at(metaLine, fmt.Errorf("duplicate decreases line"))
					}
					if mm := reDecreases.FindStringSubmatch(c); mm != nil {
						fn.DecNames = []string{mm[1]}
						i++
						break
					}
					if mm := reDecreasesSchema.FindStringSubmatch(c); mm != nil {
						fn.DecNames = []string{mm[1], mm[2]}
						fn.DecSchema = mm[3]
						i++
						break
					}
					return nil, at(metaLine, fmt.Errorf("bad decreases line: %s", c))
				case strings.HasPrefix(c, "effects "):
					mm := reEffects.FindStringSubmatch(c)
					if mm == nil {
						return nil, at(metaLine, fmt.Errorf("bad effects: %s", c))
					}
					fn.Effects = splitTop(mm[1], ',')
					i++
				case c == "tests":
					i++
					for i < len(rows) && rows[i].indent > ind {
						tline := rows[i].line
						tm := reTest.FindStringSubmatch(rows[i].code)
						if tm == nil {
							return nil, at(tline, fmt.Errorf("bad test case: %s", rows[i].code))
						}
						targs, err := parseArgs(tm[2])
						if err != nil {
							return nil, at(tline, err)
						}
						exp, err := parseSmall(tm[3])
						if err != nil {
							return nil, at(tline, err)
						}
						fn.Tests = append(fn.Tests, Test{Name: tm[1], Args: targs, Expected: exp, Line: tline})
						i++
					}
				default:
					return nil, at(metaLine, fmt.Errorf("unexpected in fn %s: %s", fn.Name, c))
				}
			}
			if i >= len(rows) || rows[i].code != "=" {
				return nil, at(declLine, fmt.Errorf("fn %s missing = body", fn.Name))
			}
			i++
			body, next, err := parseExprBlock(rows, i, -1)
			if err != nil {
				return nil, err
			}
			i = next
			for _, t := range fn.Tests {
				if !isKwargList(t.Args) {
					return nil, at(t.Line, fmt.Errorf("test %s args must be named", t.Name))
				}
			}
			fn.Body = body
			mod.Decls = append(mod.Decls, fn)
		default:
			return nil, at(declLine, fmt.Errorf("unknown top-level decl: %s", code))
		}
	}
	mod.ID = filepath.Clean(path)
	mod.File = filepath.Base(mod.ID)
	return mod, nil
}

func parsePattern(s string) (Pattern, error) {
	s = strings.TrimSpace(s)
	switch {
	case s == "_":
		return Pattern{Kind: "wild"}, nil
	case s == "true":
		return Pattern{Kind: "bool", B: true}, nil
	case s == "false":
		return Pattern{Kind: "bool"}, nil
	case strings.HasPrefix(s, `"`):
		return Pattern{Kind: "str", Str: s[1 : len(s)-1]}, nil
	}
	if m := rePatVar.FindStringSubmatch(s); m != nil && (m[1] == "Ok" || strings.Contains(m[1], ".")) {
		return Pattern{Kind: "variant", Name: m[1], Var: m[2]}, nil
	}
	if m := rePatWild.FindStringSubmatch(s); m != nil && strings.Contains(m[1], ".") {
		return Pattern{Kind: "variantWild", Name: m[1]}, nil
	}
	return Pattern{}, fmt.Errorf("bad match pattern: %s", s)
}

func parseExprBlock(rows []row, i, parentIndent int) (*Node, int, error) {
	if i >= len(rows) {
		return nil, i, fmt.Errorf("unexpected end of block")
	}
	indent, code := rows[i].indent, rows[i].code
	mline := rows[i].line
	if indent <= parentIndent {
		return nil, i, at(mline, fmt.Errorf("expected expression, found dedent: %s", code))
	}
	if strings.HasPrefix(code, "match ") {
		scrut, err := parseSmall(strings.TrimSpace(code[len("match "):]))
		if err != nil {
			return nil, i, at(mline, err)
		}
		node, next, err := parseMatchArms(rows, i+1, indent, mline, scrut)
		if err != nil {
			return nil, next, err
		}
		node.Line = mline
		return node, next, nil
	}
	sm, err := parseSmall(code)
	if err != nil {
		return nil, i, at(mline, err)
	}
	return &Node{Small: sm, Line: mline}, i + 1, nil
}

func parseMatchArms(rows []row, i, indent, mline int, scrut *Small) (*Node, int, error) {
	node := &Node{IsMatch: true, Scrut: scrut, Line: mline}
	for i < len(rows) && rows[i].indent > indent {
		ind, c := rows[i].indent, rows[i].code
		aline := rows[i].line
		if c == "given" {
			if node.Given != nil {
				return nil, i, at(aline, fmt.Errorf("duplicate given table"))
			}
			node.Given = map[string]*Small{}
			i++
			for i < len(rows) && rows[i].indent > ind {
				gline := rows[i].line
				gm := reGiven.FindStringSubmatch(rows[i].code)
				if gm == nil {
					return nil, i, at(gline, fmt.Errorf("bad given entry: %s", rows[i].code))
				}
				rhs := strings.TrimSpace(gm[2])
				if rhs == "-" {
					node.Given[gm[1]] = nil
				} else {
					sm, err := parseSmall(rhs)
					if err != nil {
						return nil, i, at(gline, err)
					}
					node.Given[gm[1]] = sm
				}
				i++
			}
			continue
		}
		m := reArm.FindStringSubmatch(c)
		if m == nil {
			return nil, i, at(aline, fmt.Errorf("bad match arm: %s", c))
		}
		pat, err := parsePattern(m[1])
		if err != nil {
			return nil, i, at(aline, err)
		}
		rest := strings.TrimSpace(m[2])
		i++
		var rhs *Node
		switch {
		case strings.HasPrefix(rest, "match "):
			sub, err := parseSmall(strings.TrimSpace(rest[len("match "):]))
			if err != nil {
				return nil, i, at(aline, err)
			}
			rhs, i, err = parseMatchArms(rows, i, ind, aline, sub)
			if err != nil {
				return nil, i, err
			}
		case rest != "":
			sm, err := parseSmall(rest)
			if err != nil {
				return nil, i, at(aline, err)
			}
			rhs = &Node{Small: sm, Line: aline}
		default:
			rhs, i, err = parseExprBlock(rows, i, ind)
			if err != nil {
				return nil, i, err
			}
		}
		node.Arms = append(node.Arms, Arm{Pat: pat, Rhs: rhs, Line: aline})
	}
	if len(node.Arms) == 0 {
		return nil, i, at(indent, fmt.Errorf("match with no arms"))
	}
	if node.Given != nil && scrut.Kind != "call" {
		return nil, i, at(node.Line, fmt.Errorf("given table on a non-call match"))
	}
	return node, i, nil
}
