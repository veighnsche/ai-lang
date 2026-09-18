// Package scan holds the source-text scanners shared by the compiler
// and the repo tools: one implementation of the R1 brace scan feeds
// canlc's parser and modcheck's repo gate, so the ban cannot drift
// between the two. RepoRoot serves the tools' repo-root walk-up.
package scan

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// BraceOutsideString reports whether line carries { or } outside a
// string literal. Braces inside "..." are data (JSON, CSS, templates),
// never delimiters; braces in code or comments stay banned (R1).
// String tracking matches the comment rule: " opens, \ escapes the
// next byte, " closes, and // outside a string starts a comment whose
// quotes never toggle string state.
func BraceOutsideString(line string) bool {
	inStr := false
	for i := 0; i < len(line); {
		ch := line[i]
		if inStr {
			if ch == '\\' && i+1 < len(line) {
				i++
			} else if ch == '"' {
				inStr = false
			}
		} else {
			if ch == '"' {
				inStr = true
			} else if ch == '/' && i+1 < len(line) && line[i+1] == '/' {
				rest := line[i:]
				return strings.ContainsAny(rest, "{}")
			} else if ch == '{' || ch == '}' {
				return true
			}
		}
		i++
	}
	return false
}

// HasBraceOutsideString reports whether any line of body carries a
// brace outside a string literal.
func HasBraceOutsideString(body string) bool {
	for _, line := range strings.Split(body, "\n") {
		if BraceOutsideString(line) {
			return true
		}
	}
	return false
}

// RepoRoot walks up from the working directory to the repo root
// (go.mod), for tools that run from anywhere inside the repo.
func RepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("repo root (go.mod) not found")
		}
		dir = parent
	}
}
