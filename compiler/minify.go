package main

import "strings"

// minifyTS compacts canlc's TypeScript emit: comments and unneeded
// whitespace go, everything else stays byte-identical. It runs on
// every compile automatically; each `<stem>.ts` gains a
// `<stem>.min.ts` beside it.
//
// The pass is string-aware (`'`, `"`, and backtick templates with
// `${...}` nesting): `//` or `/*` inside a literal is data, never a
// comment. It assumes the emit contains no regex literals (true of
// canlc output); a `/` outside a string only opens a comment.
//
// A gap (comment or whitespace run) survives as one space only
// where dropping it would fuse two tokens: between two word chars
// (`[A-Za-z0-9_$]`, so `return {` and `case "ok"` still minify but
// `a b` does not become `ab`), and between two `+`/`-` chars (so
// `a - -b` never becomes the decrement `a--b`).
func minifyTS(src string) string {
	const gap = -1 // whitespace/comment gap marker
	out := make([]byte, 0, len(src))
	gaps := make([]bool, 0, len(src)) // gap before out[i]
	pendingGap := false

	mark := func() {
		pendingGap = true
	}
	emit := func(b byte) {
		out = append(out, b)
		gaps = append(gaps, pendingGap)
		pendingGap = false
	}

	i, n := 0, len(src)
	// exprDepth tracks `${...}` brace nesting inside template
	// literals; >0 means the current backtick span is code again.
	var tmpl []int // stack of brace depths, one per open template
	inStr := byte(0)
	for i < n {
		c := src[i]
		if inStr != 0 {
			emit(c)
			if c == '\\' && i+1 < n {
				i++
				emit(src[i])
			} else if c == inStr {
				inStr = 0
			} else if inStr == '`' && c == '$' && i+1 < n && src[i+1] == '{' {
				// Re-enter code inside ${...}; the template
				// resumes at the matching closing brace.
				emit(src[i+1])
				i++
				tmpl = append(tmpl, 0)
				inStr = 0
			}
			i++
			continue
		}
		if len(tmpl) > 0 {
			// Inside ${...}: braces nest, strings nest.
			switch {
			case c == '"' || c == '\'' || c == '`':
				inStr = c
				emit(c)
			case c == '{':
				tmpl[len(tmpl)-1]++
				emit(c)
			case c == '}':
				if tmpl[len(tmpl)-1] == 0 {
					tmpl = tmpl[:len(tmpl)-1]
					emit(c)
					inStr = '`'
				} else {
					tmpl[len(tmpl)-1]--
					emit(c)
				}
			case c == '/' && i+1 < n && (src[i+1] == '/' || src[i+1] == '*'):
				i = skipComment(src, i)
				mark()
				continue
			case c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f':
				mark()
			default:
				emit(c)
			}
			i++
			continue
		}
		switch {
		case c == '"' || c == '\'' || c == '`':
			inStr = c
			emit(c)
		case c == '/' && i+1 < n && src[i+1] == '/':
			i = skipComment(src, i)
			mark()
			continue
		case c == '/' && i+1 < n && src[i+1] == '*':
			i = skipComment(src, i)
			mark()
			continue
		case c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f':
			mark()
		default:
			emit(c)
		}
		i++
	}

	var b strings.Builder
	b.Grow(len(out))
	for i := 0; i < len(out); i++ {
		if gaps[i] && i > 0 && needSpace(out[i-1], out[i]) {
			b.WriteByte(' ')
		}
		b.WriteByte(out[i])
	}
	return b.String()
}

// skipComment returns the index just past the `//...` or `/*...*/`
// comment opening at src[i] (src[i] == '/'). Unterminated input
// ends the file; there is nothing after it worth keeping.
func skipComment(src string, i int) int {
	if src[i+1] == '/' {
		for i < len(src) && src[i] != '\n' {
			i++
		}
		return i
	}
	for i = i + 2; i < len(src); i++ {
		if src[i] == '*' && i+1 < len(src) && src[i+1] == '/' {
			return i + 2
		}
	}
	return len(src)
}

func isWord(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' ||
		b >= '0' && b <= '9' || b == '_' || b == '$'
}

// needSpace reports whether dropping the gap between two emitted
// bytes would fuse tokens or change operators.
func needSpace(prev, next byte) bool {
	if isWord(prev) && isWord(next) {
		return true
	}
	if (prev == '+' || prev == '-') && (next == '+' || next == '-') {
		return true
	}
	return false
}
