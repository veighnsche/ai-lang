// Command gramcheck verifies the ai-lang TextMate grammar: every syntactic
// class needs its own scope, and each scope's regex must fire on a
// representative sample. Smoke test only — it does NOT verify rule
// precedence (first-match-wins in TextMate), so keep specific rules before
// general ones by hand when editing the grammar.
//
// Run from anywhere inside the repo: go run ./tools/gramcheck
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var jsonFiles = []string{
	"package.json",
	"language-configuration.json",
	"syntaxes/ail.tmGrammar.json",
}

type sampleCase struct {
	scope   string
	samples []string
}

var cases = []sampleCase{
	{"keyword.control.ail", []string{"on", "rev", "match", "given", "brand", "const", "forward", "extern", "seal", "seals_from", "decreases", "state", "effects", "exchange", "args", "outcome", "exports_utf8", "asset_bridge", "and", "or", "not"}},
	{"storage.type.primitive.ail", []string{"str", "int", "dec"}},
	{"constant.numeric.integer.ail", []string{"42"}},
	{"constant.numeric.decimal.ail", []string{`d"1.5"`}},
	{"constant.language.ail", []string{"true", "Ok", "_"}},
	{"keyword.operator.ail", []string{"=>", "->", "==", "+", "-", "*", "/", "%"}},
	{"entity.name.function.ail", []string{"db__get_user", "auth__login"}},
	{"entity.name.type.ail", []string{"Db__User", "Auth__Session"}},
	{"entity.name.tag.ail", []string{"user_id:", "happy:"}},
	{"constant.numeric.version.ail", []string{"@3"}},
	{"variable.other.readwrite.ail", []string{"id", "pw", "user"}},
	{"string.quoted.double.ail", []string{`"secret"`, `e"a\nb"`}},
	{"comment.line.double-slash.ail", []string{"// note"}},
}

func loadJSON(dir, name string) (map[string]any, error) {
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return nil, err
	}
	var v map[string]any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("%s: %v", name, err)
	}
	return v, nil
}

func patternsOf(grammar map[string]any) []map[string]any {
	var out []map[string]any
	if list, ok := grammar["patterns"].([]any); ok {
		for _, p := range list {
			if m, ok := p.(map[string]any); ok {
				out = append(out, m)
			}
		}
	}
	return out
}

func strOf(v any) string {
	s, _ := v.(string)
	return s
}

// check verifies the grammar in dir and returns every violation found.
func check(dir string) []string {
	var errs []string
	for _, name := range jsonFiles {
		if _, err := loadJSON(dir, name); err != nil {
			errs = append(errs, fmt.Sprintf("%v", err))
		}
	}
	grammar, err := loadJSON(dir, "syntaxes/ail.tmGrammar.json")
	if err != nil {
		return errs
	}
	patterns := patternsOf(grammar)

	var errorRules []map[string]any
	for _, p := range patterns {
		if strOf(p["name"]) == "keyword.declaration.error.ail" {
			errorRules = append(errorRules, p)
		}
	}
	if len(errorRules) != 1 {
		errs = append(errs, "expected exactly one error rule")
	} else if !strings.Contains(strOf(errorRules[0]["match"]), "error") {
		errs = append(errs, "error rule must match the error keyword")
	}
	for _, p := range patterns {
		if strOf(p["name"]) == "keyword.control.ail" &&
			strings.Contains(strOf(p["match"]), "error") {
			errs = append(errs, "generic keyword rule must not also match error")
		}
	}
	var keywordMatches []string
	var pinRules []map[string]any
	for _, p := range patterns {
		switch strOf(p["name"]) {
		case "keyword.control.ail":
			keywordMatches = append(keywordMatches, strOf(p["match"]))
		case "constant.numeric.version.ail":
			pinRules = append(pinRules, p)
		}
	}
	foundRev := false
	for _, m := range keywordMatches {
		if strings.Contains(m, "rev") {
			foundRev = true
		}
	}
	if !foundRev {
		errs = append(errs, "rev must be a keyword")
	}
	if len(pinRules) != 1 {
		errs = append(errs, "expected one version-pin rule")
	} else if ok, _ := regexp.MatchString(strOf(pinRules[0]["match"]), "db__get_user@3"); !ok {
		errs = append(errs, "pin rule must match @N")
	}

	scopes := map[string][]string{}
	for _, p := range patterns {
		match := strOf(p["match"])
		if match == "" {
			match = strOf(p["begin"])
		}
		if match == "" {
			continue
		}
		var names []string
		if n := strOf(p["name"]); n != "" {
			names = append(names, n)
		}
		if caps, ok := p["captures"].(map[string]any); ok {
			for _, c := range caps {
				if cm, ok := c.(map[string]any); ok {
					if n := strOf(cm["name"]); n != "" {
						names = append(names, n)
					}
				}
			}
		}
		for _, n := range names {
			scopes[n] = append(scopes[n], match)
		}
	}
	for _, c := range cases {
		rules := scopes[c.scope]
		if len(rules) == 0 {
			errs = append(errs, fmt.Sprintf("no rule for %s", c.scope))
			continue
		}
		for _, sample := range c.samples {
			hit := false
			for _, m := range rules {
				if ok, _ := regexp.MatchString(m, sample); ok {
					hit = true
					break
				}
			}
			if !hit {
				errs = append(errs, fmt.Sprintf("%q matches nothing under %s", sample, c.scope))
			}
		}
	}

	var blocks []map[string]any
	for _, p := range patterns {
		if _, ok := p["begin"]; ok && strings.Contains(strOf(p["begin"]), "given") {
			blocks = append(blocks, p)
		}
	}
	if len(blocks) != 1 {
		errs = append(errs, "expected one given-block rule")
	} else {
		var inner []string
		if list, ok := blocks[0]["patterns"].([]any); ok {
			for _, q := range list {
				if qm, ok := q.(map[string]any); ok {
					if caps, ok := qm["captures"].(map[string]any); ok {
						for _, c := range caps {
							if cm, ok := c.(map[string]any); ok {
								inner = append(inner, strOf(cm["name"]))
							}
						}
					}
				}
			}
		}
		has := func(n string) bool {
			for _, v := range inner {
				if v == n {
					return true
				}
			}
			return false
		}
		if !has("entity.name.tag.ail") {
			errs = append(errs, "given-block must scope test-case keys")
		}
		if !has("constant.language.unreachable.ail") {
			errs = append(errs, "given-block must scope the - marker")
		}
	}
	sort.Strings(errs)
	return errs
}

// repoRoot walks up from the working directory to the repo root (go.mod).
func repoRoot() (string, error) {
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
			return "", fmt.Errorf("repo root (go.mod) not found")
		}
		dir = parent
	}
}

func main() {
	root, err := repoRoot()
	if err != nil {
		fmt.Println("GRAMMAR CHECK FAILED")
		fmt.Println(" -", err)
		os.Exit(1)
	}
	if errs := check(filepath.Join(root, "editors", "vscode")); len(errs) > 0 {
		fmt.Println("GRAMMAR CHECK FAILED")
		for _, e := range errs {
			fmt.Println(" -", e)
		}
		os.Exit(1)
	}
	fmt.Println("grammar OK: error, rev, pins, and core syntax highlighted")
}
