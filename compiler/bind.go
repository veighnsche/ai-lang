package main

import "fmt"

// bindSlots resolves a call's argument list against the callee's
// parameter list, returning one slot index per source argument: slot[i]
// is the parameter position bound to the i-th argument in source order.
// Positional arguments bind by source index, named arguments bind by
// parameter name. Every parameter receives exactly one argument:
// unknown names, duplicate assignments (including mixed-form double
// assignment of one slot by position and by name), out-of-range
// positionals, and omissions are all errors.
//
// This is the single binding rule shared by checking, evaluation, and
// emission. Argument expressions still evaluate in source order; only
// the binding vector is canonicalized, so a reordered named call means
// the same thing in every phase. Script (exchange) rows are outside
// this rule: they are all-named by grammar and proven by name at each
// hit in checkExchangeArgs.
func bindSlots(fname string, args []Arg, params [][2]string) ([]int, error) {
	index := make(map[string]int, len(params))
	for j, p := range params {
		index[p[0]] = j
	}
	slots := make([]int, len(args))
	assigned := make([]bool, len(params))
	for i, a := range args {
		slot := -1
		if a.HasName {
			j, ok := index[a.Name]
			if !ok {
				return nil, fmt.Errorf("call %s has no param %s", fname, a.Name)
			}
			slot = j
		} else {
			if i >= len(params) {
				return nil, fmt.Errorf("call %s takes %d args for %d params", fname, len(args), len(params))
			}
			slot = i
		}
		if assigned[slot] {
			return nil, fmt.Errorf("call %s supplies arg %s twice", fname, params[slot][0])
		}
		assigned[slot] = true
		slots[i] = slot
	}
	for j, p := range params {
		if !assigned[j] {
			return nil, fmt.Errorf("call %s is missing arg %s", fname, p[0])
		}
	}
	return slots, nil
}

// bindCtorArgs names a constructor's positional arguments from the
// declaration's field order (a92): positional i binds fields[i],
// named arguments keep their names. It reports only positional
// faults — overflow and double supply against a named claim —
// with the same wording as the static constructor check; named
// verification (unknown or missing fields) stays with the caller,
// exactly like the static split between binding and checking.
// Untyped contexts (given tables, test expectations) evaluate raw,
// so evaluation binds here what checking binds nowhere. Ok takes
// the conventional single field [value]; completeness is never
// required here (Ok() stays the empty ok dict, as before).
func bindCtorArgs(ctor string, args []Arg, fields []string) ([]string, error) {
	names := make([]string, len(args))
	claimed := map[string]bool{}
	for _, a := range args {
		if a.HasName {
			claimed[a.Name] = true
		}
	}
	for i, a := range args {
		if a.HasName {
			names[i] = a.Name
			continue
		}
		if i >= len(fields) {
			return nil, fmt.Errorf("%s takes %d args for %d fields", ctor, len(args), len(fields))
		}
		if claimed[fields[i]] {
			return nil, fmt.Errorf("%s supplies field %s twice", ctor, fields[i])
		}
		claimed[fields[i]] = true
		names[i] = fields[i]
	}
	return names, nil
}
