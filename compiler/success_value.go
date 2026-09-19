package main

import (
	"fmt"
	"strings"
)

// B11's typed Ok is an outcome boundary, never a data constructor. The
// annotation is exact; it does not change the existing declared success ABI.
func (c *tycker) checkWholeOk(s *Small, want string, line int, env map[string]string) {
	fail := func(msg string) {
		c.out = append(c.out, spanDiag(c.text, line, "error", msg, "Ok", CodeTypeMismatch))
	}
	if len(s.TypeArgs) != 1 || !c.knownType(s.TypeArgs[0]) {
		fail("typed Ok needs exactly one known success type")
		return
	}
	typ := s.TypeArgs[0]
	if want == "" || !sameType(typ, want) {
		fail(fmt.Sprintf("typed Ok<%s> requires a matching success boundary, got %s", typ, want))
	}
	if len(s.Args) != 1 || (s.Args[0].HasName && s.Args[0].Name != "value") {
		fail("typed Ok takes exactly one whole value (positional or value =)")
		for _, a := range s.Args {
			c.value(a.V, "", line, env, "typed Ok value")
		}
		return
	}
	a := &s.Args[0]
	c.value(a.V, "", line, env, "typed Ok value")
	if got, ok := c.typeOf(a.V, env); !ok {
		fail(fmt.Sprintf("typed Ok needs a data value of type %s, not an outcome", typ))
	} else if !sameType(got, typ) {
		c.mismatch(line, "typed Ok value", got, typ, tokenOf(a.V))
	}
	a.Name, a.HasName = "value", true
}

func (c *tycker) patternSuccessType(p Pattern, ret string, line int) string {
	if len(p.TypeArgs) == 0 {
		return c.successBinderType(ret)
	}
	if len(p.TypeArgs) != 1 || !c.knownType(p.TypeArgs[0]) || !sameType(p.TypeArgs[0], ret) {
		c.out = append(c.out, spanDiag(c.text, line, "error",
			fmt.Sprintf("typed Ok pattern must name the exact success type %s", ret), "Ok", CodeTypeMismatch))
	}
	return ret
}

// Kernel-only synthetic payloads have no source-namable success type.
func (c *tycker) rejectKernelWholeOk(p Pattern, line int) {
	if len(p.TypeArgs) != 0 {
		c.out = append(c.out, spanDiag(c.text, line, "error",
			"typed Ok pattern requires a declared success type, not a synthetic kernel payload", "Ok", CodeTypeMismatch))
	}
}

// Proof lowering alone expands a whole record into projections. Runtime
// lowering evaluates the operand once; proof expressions are pure terms.
func wholeOkArgs(s *Small, recs map[string][][2]string) []Arg {
	if s.Ctor != "Ok" || len(s.TypeArgs) != 1 || len(s.Args) != 1 {
		return s.Args
	}
	v := s.Args[0].V
	if fields, ok := recs[s.TypeArgs[0]]; ok {
		args := make([]Arg, 0, len(fields))
		for _, f := range fields {
			args = append(args, Arg{Name: f[0], HasName: true, V: &Small{Kind: "proj", L: v, Field: f[0]}})
		}
		return args
	}
	return []Arg{{Name: "value", HasName: true, V: v}}
}

func (e *emitter) wholeSuccessBinding(p Pattern, tmp string) (string, error) {
	if p.Name != "Ok" || len(p.TypeArgs) == 0 {
		return tmp, nil
	}
	if len(p.TypeArgs) != 1 {
		return "", fmt.Errorf("typed Ok pattern needs one success type")
	}
	if fields, ok := e.recs[p.TypeArgs[0]]; ok {
		return tsRecordProjection(tmp, fields, false), nil
	}
	if !valueSuccess(p.TypeArgs[0], e.variants[p.TypeArgs[0]] != nil || e.brands[p.TypeArgs[0]] != "") {
		return "", fmt.Errorf("unknown typed Ok success %s", p.TypeArgs[0])
	}
	return tmp + ".value", nil
}

func (e *emitter) emitWholeOk(s *Small) (string, error) {
	if len(s.TypeArgs) != 1 || len(s.Args) != 1 || (s.Args[0].HasName && s.Args[0].Name != "value") {
		return "", fmt.Errorf("typed Ok takes one type and one whole value")
	}
	v, err := e.emitValue(s.Args[0].V)
	if err != nil {
		return "", err
	}
	if fields, ok := e.recs[s.TypeArgs[0]]; ok {
		// Structural parameter annotation needs only the field dependencies
		// already imported for the declared flattened result. No extra nominal
		// record import, spread, or repeated evaluation of the operand.
		var types []string
		for _, f := range fields {
			t, err := tsTypeB(f[1], e.brands, e.recs, e.variants, e.errTypes)
			if err != nil {
				return "", err
			}
			types = append(types, tsField(f[0], t))
		}
		tmp := e.fresh()
		return "((" + tmp + ": { " + strings.Join(types, "; ") + " }) => (" + tsRecordProjection(tmp, fields, true) + "))(" + v + ")", nil
	}
	if !valueSuccess(s.TypeArgs[0], e.variants[s.TypeArgs[0]] != nil || e.brands[s.TypeArgs[0]] != "") {
		return "", fmt.Errorf("unknown typed Ok success %s", s.TypeArgs[0])
	}
	return "{ " + tsTag + `: "ok", value: ` + v + " }", nil
}

func tsRecordProjection(v string, fields [][2]string, outcome bool) string {
	parts := []string{}
	if outcome {
		parts = append(parts, tsTag+`: "ok" as const`)
	}
	for _, f := range fields {
		parts = append(parts, tsField(f[0], v+"."+f[0]))
	}
	return "{ " + strings.Join(parts, ", ") + " }"
}
