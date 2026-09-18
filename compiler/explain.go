package main

// a71 (a61 item 1, first cut) established `ailc explain CODE`:
// the rule, a minimal violation, the legal fix. Populated to
// full coverage 2026-09-18: every registered code carries a
// per-code entry grounded in its emission message. The family
// fallback stays as a safety net for codes added without docs
// (the completeness test fails first). Docs are versioned with
// the compiler; codes stay frozen.

import (
	"fmt"
	"io"
	"strings"
)

type explainEntry struct {
	rule    string
	violate string
	fix     string
}

var explainDocs = map[string]explainEntry{
	CodeForeignRaise: {
		rule:    "A function may only raise error kinds in its own emits list (R5).",
		violate: `fn m__go() emits [m.known], body raises m.stray().`,
		fix:     "Add the kind to emits, or raise a declared kind. The diagnostic names the raised kind (found) and the declared set (expected).",
	},
	CodeBareErrorKind: {
		rule:    "Test expectations must construct the complete error value; a bare kind proves nothing about the payload (a12).",
		violate: `go => m.bad  (where error m.bad(id: str))`,
		fix:     "Write m.bad(id = ...) with all fields. The diagnostic shows the complete construction (expected).",
	},
	CodeDanglingTest: {
		rule:    "Every reaching test needs a script at every foreign call (R8).",
		violate: `a test with no entry under the call's given table.`,
		fix:     "Add `test => [exchange args (...) outcome ...]` under given. The diagnostic shows the row shape (expected) and the table line (hint).",
	},
	CodeMissingArm: {
		rule:    "Call matches cover exactly the callee emits plus ok: missing or stale arms are compile errors.",
		violate: `a match on a failable call without an on <kind> arm.`,
		fix:     "Add the missing arm. The diagnostic names the kind (expected).",
	},
	CodeRevisionIdentity: {
		rule:    "One revision owns one reviewed interface: a changed fingerprint at the same revision is drift, removals and additions are explicit, and only an accepted baseline counts.",
		violate: `adding a case to a variant without bumping its revision.`,
		fix:     "Restore the accepted interface, or publish a reviewed new revision and explicitly update affected pins. Never regenerate the baseline to silence the finding: generation is not acceptance.",
	},
	CodeContractMalformed: {
		rule:    "A contract must be well-formed before it can be proved: every predicate is Boolean over the arm binder and the function params, every declared outcome has exactly one arm, and ensures matches are Boolean case expressions.",
		violate: `requires 1, or an ensures arm missing for a declared emits kind.`,
		fix:     "Correct the declaration or clause structure. The diagnostic names the offending clause (found) and the shape the rule wanted (expected).",
	},
	CodeContractUnsupported: {
		rule:    "The first proof cut covers mathematical integers, Booleans, and finite records of those sorts, with total bodies free of externs, kernels, state, and recursion. Anything outside is rejected, never partially proved or assumed.",
		violate: `result.value == x * x, or a contracted body calling an extern.`,
		fix:     "Use supported logic and body structure, or defer that contract-bearing interface. The diagnostic names the exact operator, sort, or dependency (found).",
	},
	CodeContractUnverifiedDep: {
		rule:    "A contract is an obligation first: no caller may reason from a callee summary until that callee is verified in the proving run. Presence in the source, an accepted revision, and passing examples are not verification.",
		violate: `a contracted function calling another contracted function before any prover ran.`,
		fix:     "Verify the source callee and its closure in the proving run; do not treat its declaration as proof. The diagnostic names the caller and the unavailable callee.",
	},
	CodeContractUnproven: {
		rule:    "Every exit of a contracted body proves its outcome predicate from the entry assumptions plus the path facts, and every call site proves the callee precondition. A satisfiable negation is a counterexample, not an acceptance.",
		violate: `a maximum body returning max(left,right)+1.`,
		fix:     "Correct the body, or strengthen the relevant explicit summary through reviewed revision changes. The diagnostic carries the counterexample inputs (found).",
	},
	CodeContractInconclusive: {
		rule:    "Timeout, solver failure, unknown, or a missing solver backend rejects verification as inconclusive. Inconclusive is never acceptance.",
		violate: `a proof attempt with no solver on PATH.`,
		fix:     "Rerun within the supported policy or reduce proof complexity. The diagnostic names the obligation and the cause (found).",
	},
	CodeContractInadmissibleTest: {
		rule:    "Test entry rows satisfy the tested function's requires: a row outside the admitted inputs proves nothing about the contract.",
		violate: `a row with x = -1 under requires x >= 0.`,
		fix:     "Supply an admitted row; do not skip the row or silently narrow expectations.",
	},
	CodeParse: {
		rule:    "The file must parse: grammar shapes, the braces ban, and rev presence (R1). Deliberately coarse — one code for every syntax failure.",
		violate: `any unparseable construct; the message names the site and the failing shape.`,
		fix:     "Read the message at the reported line. Frequent causes: a curly brace outside a string literal, a malformed mod header, a missing rev. Only parser failures stay line-wide; everything else narrows to a token.",
	},
	CodeFnNaming: {
		rule:    "Functions read domain__verb, lowercase domain, and never take a store-operation name (R3).",
		violate: `fn go(id: str), or fn state__get(...).`,
		fix:     "Rename to domain__verb (e.g. auth__login). Identifiers never change when revs change; the name is the stable address.",
	},
	CodeTypeNaming: {
		rule:    "Types, variants, and brands read Domain__Name (R3).",
		violate: `type user(...), brand hash is str.`,
		fix:     "Rename to Domain__Name (e.g. Db__User). Same rule covers variant and brand declarations.",
	},
	CodeConstNaming: {
		rule:    "Constants read domain__SCREAMING: a lowercase domain prefix, then an uppercase name (R3).",
		violate: `const Foo: int rev 1 = 1.`,
		fix:     "Rename to domain__NAME (e.g. std__ascii__COLON). Kinds stay visually distinct: lowercase functions, CamelCase types, SCREAMING constants.",
	},
	CodeUsesPin: {
		rule:    "Every uses entry pins a rev: name@N. Unpinned uses are errors (R4).",
		violate: `uses [db__get_user] with no @N.`,
		fix:     "Pin the rev you reviewed against (db__get_user@5). No floating versions; upgrading a pin is an explicit one-line diff.",
	},
	CodeUsesResolve: {
		rule:    "Every uses entry must resolve to exactly one provider module (R2).",
		violate: `uses [db__get_user@5] with no file providing it.`,
		fix:     "Pass the provider file to the compiler, or fix the name. Unresolved and double-provided are both compile errors.",
	},
	CodeUsesRev: {
		rule:    "A uses pin must name a rev the provider actually declares (R4).",
		violate: `uses [db__get_user@9] while the provider declares rev 5.`,
		fix:     "Pin a declared rev, or upgrade the provider (change code and rev together — one never moves without the other).",
	},
	CodeUnknownConst: {
		rule:    "A constant reference resolves to exactly one declared constant in the build.",
		violate: `Ok(value = m__NOPE) with no such const declared.`,
		fix:     "Declare the constant, or correct the name. Qualified names resolve globally; there are no locals to shadow them.",
	},
	CodeConstNotInUses: {
		rule:    "A foreign constant needs a rev pin in uses, like a foreign call (R2/R4).",
		violate: `Ok(value = lib__K) with uses [] while lib.ail provides lib__K@1.`,
		fix:     "Pin the provider rev (uses [lib__K@1]). Same-module constants need no pin.",
	},
	CodeDupFn: {
		rule:    "One name, one function per file: double definitions are rejected.",
		violate: `two fn auth__login blocks in one file.`,
		fix:     "Delete or rename one. If two modules need the name, they are separate domains — qualify accordingly.",
	},
	CodeDupSibling: {
		rule:    "One name, one provider across the build: a sibling file defining the same item is a double definition (R2).",
		violate: `db__get_user defined in both db.ail and cache.ail.`,
		fix:     "Keep one provider; the other file uses it through uses with a rev pin.",
	},
	CodeDupVariant: {
		rule:    "One name, one variant declaration: double-declared variants are rejected like double-defined functions.",
		violate: `two variant Shape blocks, in one file or across siblings.`,
		fix:     "Keep one declaration; every other module references it, never re-declares it.",
	},
	CodeCaseCollision: {
		rule:    "Case identity is unique: no case collides with a record type, no duplicate case in one variant, no case declared by two variants.",
		violate: `case Circle in both Shape and Blob, or a case named like an existing record type.`,
		fix:     "Rename the colliding case or merge the owners. Identity collisions are never resolved by qualification.",
	},
	CodeDupConst: {
		rule:    "One name, one constant per file: double definitions are rejected like double-defined functions.",
		violate: `two const m__A blocks in one file.`,
		fix:     "Delete or rename one. Cross-module sharing goes through uses with a rev pin, never re-declaration.",
	},
	CodeProvidesMiss: {
		rule:    "provides names exactly what the file defines: anything defined but unlisted is an error (R2).",
		violate: `fn auth__verify defined but absent from provides [...].`,
		fix:     "Add the name to provides, or delete the definition. (State cells are the exception: module-private, never provided.)",
	},
	CodeProvidesGhost: {
		rule:    "provides names exactly what the file defines: anything listed but never defined is an error (R2).",
		violate: `provides [auth__login] with no such fn in the file.`,
		fix:     "Define it or drop it from provides. No forward declarations, no aspirational exports.",
	},
	CodeUnknownCall: {
		rule:    "Every callee resolves: same-file helper, uses-pinned ail function, own-module extern, or state intrinsic. Unknown callees and other modules' externs are rejected.",
		violate: `call db__get_user(...) with no such function, or calling another module's extern directly.`,
		fix:     "Fix the name, add the uses pin, or declare your own extern — externs are module-local and never shared.",
	},
	CodeCallNotInUses: {
		rule:    "A function touching the outside world says so in uses: cross-file callees need a pinned entry (R6).",
		violate: `call db__get_user(...) with db__get_user@5 missing from uses.`,
		fix:     "Add name@rev to uses. Same-file helpers are exempt; everything else declares its dependency.",
	},
	CodeCallOutside: {
		rule:    "Calls live only as match scrutinees: every call's outcomes are dispatched by arms, never dropped or sequenced invisibly.",
		violate: `a bare call standing as an arm body or nested in an expression.`,
		fix:     "Wrap the call: match call f(...) with one arm per declared outcome. There is no fire-and-forget call.",
	},
	CodeCallNested: {
		rule:    "One call per scrutinee: nested calls inside a scrutinee are outside the admitted subset.",
		violate: `match call f(call g(...)).`,
		fix:     "Hoist the inner call into its own match with its own arms. Each call dispatches exactly once, at exactly one site.",
	},
	CodeLocalCycle: {
		rule:    "Helpers are acyclic; only direct self-recursion with a proven decrease is admitted — program-wide, same-file and cross-file (R7).",
		violate: `a__x calls a__y calls a__x, or an unproven self-call.`,
		fix:     "Break the cycle into a DAG, or shape the recursion as direct self-recursion with decreases, a canonical guard, and an exact step.",
	},
	CodeBadDecreases: {
		rule:    "decreases names an int param of the function — nothing else decreases (R7).",
		violate: `decreases fuel where fuel is a str, or decreases missing.`,
		fix:     "Name an int param, or add one. The measure is always an integer the sites visibly shrink.",
	},
	CodeStaleDecreases: {
		rule:    "A decreases line without a self-call is a stale proof claim: admission proofs must be load-bearing.",
		violate: `decreases fuel on a function that never calls itself.`,
		fix:     "Delete the line, or write the recursion it claims. Proof text that proves nothing is an error, not clutter.",
	},
	CodeNoDecrease: {
		rule:    "Every self-call site passes the exact admitted step: p - 1 for the unit loop (R7), (b, a % b) for euclid, (lo, mid)/(mid, hi) for narrowing (a19).",
		violate: `a self-call passing fuel unchanged, or fuel - 2.`,
		fix:     "Pass the exact step the decreases line promises. Larger or computed steps need their own blessed schema — the compiler never infers one.",
	},
	CodeNoGuard: {
		rule:    "Admitted recursion sits under the false arm of its canonical guard: p <= 0 for the unit loop, b <= 0 for euclid (R7, a19).",
		violate: `a self-call under the true arm, or with no guard match at all.`,
		fix:     "Put the self-call under the false arm of the canonical p <= 0 guard. Negative entries take the base arm and return a declared outcome.",
	},
	CodeBadBinding: {
		rule:    "Call args bind through the callee signature: arity and names resolve, positionally or named — nothing unexpected, nothing missing.",
		violate: `db__get_user(id, extra) against a one-param signature.`,
		fix:     "Match the callee's params exactly. The message names the offending arg; the signature is the contract.",
	},
	CodeUndeclaredEffect: {
		rule:    "State authority is declared beside emits: effects [C.read, C.write], transitive through local calls, no inference (R6).",
		violate: `touching Count__total with no effects line, or via a helper whose authority you did not declare.`,
		fix:     "Add the capability to effects. The message names the cell and the path (direct, helper, or uses-pin) it arrives by.",
	},
	CodeStaleEffect: {
		rule:    "Declared effects must be used: stale capabilities are errors, like stale decreases lines.",
		violate: `effects [Count__total.write] on a function that never writes it.`,
		fix:     "Delete the capability or write the access it claims. Authority text that authorizes nothing is an error.",
	},
	CodeNoGiven: {
		rule:    "Foreign calls are stubbed at the call site: every such call carries a given table, one exchange sequence per test (R8).",
		violate: `match call db__get_user(id) with no given table beneath it.`,
		fix:     "Add the given table with one row per reaching test (unreachable tests written - explicitly). Deterministic calls take no table — see AIL3106.",
	},
	CodeGivenOnLocal: {
		rule:    "Only a single foreign-call match takes given: deterministic calls (same-file helpers, state intrinsics) execute, and multi-scrutinee matches take no table.",
		violate: `a given table on match call count__bump(by), or under match x, y.`,
		fix:     "Delete the table. Helpers execute with no table; their own internal tables script every reaching test.",
	},
	CodeBadStub: {
		rule:    "Stub outcomes are Ok(..) or a declared error — nothing else scripts (R8).",
		violate: `outcome 42, or outcome m.undeclared_thing.`,
		fix:     "Script Ok(...) or a declared error kind. The message names the offending outcome and its owner.",
	},
	CodeStubNotInEmit: {
		rule:    "Stub outcomes are restricted to the callee's declared emits, resolved through uses to the provider (R8). Invented outcomes are errors.",
		violate: `outcome db.down() scripted at a call whose callee never emits db.down.`,
		fix:     "Script only what the callee declares, or fix the callee's emits. The table proves permitted responses, not imagined ones.",
	},
	CodeDeadScript: {
		rule:    "Given-table keys must equal the reaching test names: a script no test selects is dead (R8).",
		violate: `zzz => [...] under given with no test named zzz.`,
		fix:     "Delete the row or add the test. Keys and tests are a set equality, checked both directions (see AIL3105).",
	},
	CodeNoExchange: {
		rule:    "Every script row is an exchange binding expected call args to one permitted outcome: outcome-only rows prove nothing about the request (R8).",
		violate: `happy => [db.down()] with no args.`,
		fix:     "Write exchange args (...) outcome .... Expected args must arrive equal and the call must supply nothing unexpected.",
	},
	CodeInconsistentScript: {
		rule:    "Scripts must agree with what the provider computes: a script contradicting the provider is refused, not run.",
		violate: `scripting "zz" as success where the provider computes invalid_hex.`,
		fix:     "Script what the provider actually computes for those args. The message shows computed vs scripted side by side.",
	},
	CodeDupTest: {
		rule:    "Test names are unique per function: the decision table has no duplicate rows.",
		violate: `two tests named happy under one fn.`,
		fix:     "Rename one. Test names address rows from every given table — ambiguity would split the evidence.",
	},
	CodeUnknownArg: {
		rule:    "Test inputs name the function's params: unknown args are errors.",
		violate: `happy(id = "u", typo = 1) where the fn has no typo param.`,
		fix:     "Fix the name or add the param. Arg names resolve through the callee signature, tests included.",
	},
	CodeMissingArg: {
		rule:    "Test inputs supply every param: partial rows prove partial behavior and are errors.",
		violate: `happy(id = "u") where the fn also takes pw.`,
		fix:     "Supply the missing arg. Every test is a complete call specification, name(inputs) => expected.",
	},
	CodeMissingTests: {
		rule:    "Every function ships its decision table: tests and business logic are one artifact (R7).",
		violate: `a fn with no tests block.`,
		fix:     "Write the decision table: name(inputs) => expected, one row per behavior. Untested functions do not ship.",
	},
	CodeUnusedUses: {
		rule:    "uses entries must be called: unused imports are warnings, not dead weight carried silently.",
		violate: `uses [db__get_user@5] in a module that never calls it.`,
		fix:     "Delete the entry or write the call. The dependency list is exact, not aspirational.",
	},
	CodeUnusedParam: {
		rule:    "Params must be used: an ignored param is a warning that the signature promises what the body ignores.",
		violate: `fn m__go(value: int, extra: str) whose body never mentions extra.`,
		fix:     "Use it or drop it from the signature (and from every test row and call site). Signatures are exact.",
	},
	CodeUnknownKind: {
		rule:    "Every raised and declared error kind is declared somewhere: unknown kinds are rejected at both sites (R5).",
		violate: `raising m.stray() with no error m.stray decl, or listing it in emits.`,
		fix:     "Declare error m.stray(...) in its module, or fix the spelling. Kinds are values with declarations, never strings.",
	},
	CodeProofOther: {
		rule:    "Proof failures the prover cannot classify further still fail the build: unclassified is never accepted.",
		violate: `any exhaustiveness failure outside the named proof codes.`,
		fix:     "Read the message — it carries the specific complaint. If it recurs opaquely, the prover needs a named code for it.",
	},
	CodeStaleArm: {
		rule:    "Arms name outcomes the scrutinee can produce: a stale arm is a dead dispatch and a compile error.",
		violate: `on db.down ... where the callee never emits db.down.`,
		fix:     "Delete the arm or fix the callee's emits. Every arm must be reachable in principle, taken in practice (see AIL4107).",
	},
	CodeBoolArms: {
		rule:    "A bool match is exactly true plus false — no more, no fewer, and never mixed with string patterns in one slot.",
		violate: `match b with only a true arm, or true beside "a" in one slot.`,
		fix:     "Write both arms. Bool slots are finite and fully enumerated; _ has no place here.",
	},
	CodeValueNoWild: {
		rule:    "A value match with a string slot needs _ coverage of the open remainder: strings are infinite, so finite arms never suffice alone.",
		violate: `match s with "a" => ... and no _ arm.`,
		fix:     "Add the _ arm for the open remainder. The message names the uncovered slot.",
	},
	CodeBadArmKind: {
		rule:    "Arm shapes match the scrutinee: call matches take on-arms, value matches take value patterns, arity is uniform across arms.",
		violate: `a bare pattern under match call, an on-arm under match x, or arms with different pattern counts.`,
		fix:     "Match the arm kind to the match kind and keep every arm at the same arity. The message names the mismatch.",
	},
	CodeVariantOnVal: {
		rule:    "Variant patterns belong to variant matches: a case pattern on a non-variant match is a kind error.",
		violate: `on Circle ... under match x where x is not the union.`,
		fix:     "Match the union scrutinee itself — cases eliminate unions, and the arms prove membership in the parent.",
	},
	CodeArmUntaken: {
		rule:    "Test-per-arm law: every match arm must execute across the decision-table run, or carry the authorized structural certificate.",
		violate: `an arm no test reaches — dead dispatch the tables never exercise.`,
		fix:     "Add a test row that takes the arm. Certified relays are the only exemption, and certified is never reported as taken.",
	},
	CodeInvalidRelay: {
		rule:    "Identity-relay certificates are checked, not trusted: a relay-shaped arm that fails the check (wrong kind, dropped field) is invalid.",
		violate: `an on e => ... arm claiming relay shape with a mismatched kind or payload.`,
		fix:     "Address the named reason, or add a test that takes the arm. The message carries the exact failure.",
	},
	CodeEagerScrutinee: {
		rule:    "Multi-scrutinee matches evaluate every scrutinee eagerly, left to right, exactly once — unlike nested matches, a failing scrutinee faults before dispatch.",
		violate: `match s[i], b where s[i] can fault out of range.`,
		fix:     "Guard the scrutinee first (bounds check in an outer match) or nest the matches so the risky evaluation sits under its guard.",
	},
	CodeTestFailed: {
		rule:    "The compiler executes every test hermetically during the build; any failure fails the build (R7). Deliberately coarse — one code for every execution failure.",
		violate: `a test whose actual value differs from expected, or that faults, starves a script, or leaves script rows over.`,
		fix:     "Read the failure: expected vs actual first, then the given scripts (consumed head? leftovers? empty list at call time?). Evaluation is deterministic — same tables, same failure.",
	},
	CodeSiblingParse: {
		rule:    "World errors report per file without hiding the rest: a sibling that fails to parse is reported at the open file too.",
		violate: `editing auth.ail while db.ail has a syntax error.`,
		fix:     "Fix the sibling file. Execution-dependent checks stay suppressed until the world parses; nothing else hides.",
	},
	CodeBadCompare: {
		rule:    "Comparisons the emitter cannot lower are refused at compile time: byte equality is deferred, never miscompiled.",
		violate: `== over Bytes values.`,
		fix:     "Compare through an admitted route (lengths, indexed members) until byte equality lands. The refusal names the type.",
	},
	CodeModuleCollision: {
		rule:    "One canonical identity per file: the same module fed twice is a collision, not an idempotent re-add.",
		violate: `passing db.ail twice on the command line, or two spellings of one path.`,
		fix:     "Pass each file once. The message names the doubled module.",
	},
	CodeFloatLiteral: {
		rule:    "Floats are ungrammatical: exact decimals spell d\"12.34\", and bare 12.34 has no spelling at all.",
		violate: `12.5 or 1e3 in any value position.`,
		fix:     "Write the value as d\"12.50\" (no exponents — spell the digits). The message shows the exact rewrite.",
	},
	CodeUnknownType: {
		rule:    "Every named type resolves: state cells hold str, int, bool, or dec; seal names a declared brand; signatures name declared types.",
		violate: `state C: float, seal Nope("x"), or a param typed M__Missing.`,
		fix:     "Declare the type/brand or fix the spelling. Unknown types fail before any value is checked.",
	},
	CodeTypeMismatch: {
		rule:    "No implicit conversions, ever: inits match cell types literally, arithmetic takes same-type operands and yields the operand type.",
		violate: `state C: int = "0", or 1 + d"2.0".`,
		fix:     "Convert explicitly through a typed route. The message names expected vs found; the fix is always spelled out, never inferred.",
	},
	CodeSealForeign: {
		rule:    "Bodies seal only their own module's brands: minting and promotion stay inside one module, so grep seal is the whole audit.",
		violate: `seal Db__Hash("...") inside auth.ail.`,
		fix:     "Seal in the owning module and pass the brand through typed params, or declare your own brand. Cross-module minting is never admitted.",
	},
	CodeInexactDivision: {
		rule:    "Division is exact or refused: decimals divide only with exact results — divide integers, not decimals.",
		violate: `d"1.0" / d"3.0".`,
		fix:     "Restructure onto integers (quotient/remainder, scaled coefficients) or a rounded-decision fn whose contract names the rounding.",
	},
	CodeRecordCycle: {
		rule:    "Only finite monomorphic products are admitted: a record reaching itself through field types is refused before tests or output.",
		violate: `type A(f: B) beside type B(g: A).`,
		fix:     "Break the cycle. The message shows the rotation (A -> B -> A); brands and base types are leaves and never participate.",
	},
	CodeSeqLiteral: {
		rule:    "Bare [...] is script rows, never a sequence value: sequences spell Seq<T>[...] with an explicit element type.",
		violate: `["a", "b"] as a value.`,
		fix:     "Write Seq<str>[\"a\", \"b\"]. Brackets enumerate (mod lists, given rows); typed values say their type.",
	},
	CodeBytesLiteral: {
		rule:    "Bytes takes one Seq<int> literal with integer members: the shape is closed and checked literally.",
		violate: `Bytes("hi"), Bytes(Seq<str>["a"]), or a non-literal argument.`,
		fix:     "Write Bytes(Seq<int>[104, 105]) with plain integer members. Non-integers and non-literals are never admitted.",
	},
	CodeBytesElementRange: {
		rule:    "Byte members are 0..255: out-of-range members are refused element by element.",
		violate: `Bytes(Seq<int>[65, 300]).`,
		fix:     "Fix the member. The message names the index and the offending value.",
	},
	CodeBytesExportAuthority: {
		rule:    "An exports_utf8 grant needs an owning module and a known brand: authority without an owner is void.",
		violate: `a grant naming a brand no module owns, or an unknown brand.`,
		fix:     "Point the grant at the brand's owning module, or declare the brand. Grants never float.",
	},
	CodeBytesExportShape: {
		rule:    "A utf8 exporter matches the exporter predicate exactly (one granted-brand param, Bytes__Value return, empty contract, single kernel match returning the result unchanged).",
		violate: `an exporter with two params, or one that transforms before returning.`,
		fix:     "Fit the predicate; the message names the failing clause. Exporters move bytes, nothing else.",
	},
	CodeAssetBridgeAuthority: {
		rule:    "An asset_bridge grant needs two owners: the named schema module must own both brands, and the sink must live in the granting (html) module. Authority without both owners is void.",
		violate: `a grant naming a brand no loaded module owns, or a sink outside the granting module.`,
		fix:     "Point the grant at the brands' owning module and keep the sink in the granting module. Grants never float.",
	},
	CodeAssetBridgeShape: {
		rule:    "An asset sink matches the sink predicate exactly (asset then policy params, single-Html__Safe record return, one emits kind, no effects, single kernel match with one Ok arm, granted role literal inside).",
		violate: `a sink with a raw-string param, or one that never checks the granted role.`,
		fix:     "Fit the predicate; the message names the failing clause. Sinks assemble fixed elements, nothing else.",
	},
	CodeConstNonliteral: {
		rule:    "V1 constant initializers are scalar literals (int, str, dec, bool) matching the declared type. Computed, alias, and cross-constant initializers are rejected.",
		violate: `const m__B: int rev 1 = 1 + 2, or = m__A.`,
		fix:     "Write the literal value. (Allocated AIL6016: the drafts' AIL6014 collided with shipped AssetBridgeAuthority.)",
	},
	CodePrimitiveShadow: {
		rule:    "Declarations never shadow a compiler kernel or primitive: Bytes and kernel names are reserved.",
		violate: `type Bytes(...) or fn bytes__hex__decode... colliding with a kernel.`,
		fix:     "Rename the declaration. The message names the shadowed kernel.",
	},
}

// explainFamily describes each code family for codes
// without a per-code doc yet.
var explainFamily = map[string]string{
	"AIL1": "Parse and grammar rules.",
	"AIL2": "Declaration rules (naming R3, rev pins R4, provides/uses R2).",
	"AIL3": "Call and decision-table rules (resolution, given/test R8).",
	"AIL4": "Error and proof rules (emits R5, exhaustiveness, test runs).",
	"AIL5": "World and tooling rules.",
	"AIL6": "Type rules (brands, seal, extern, numerics).",
}

// explainCode returns the doc text for a registered code.
// Per-code docs win; otherwise family text that admits the
// cut. False means the code is not registered at all.
func explainCode(code string) (string, bool) {
	registered := false
	for _, c := range allCodes {
		if c == code {
			registered = true
			break
		}
	}
	if !registered {
		return "", false
	}
	if e, ok := explainDocs[code]; ok {
		return fmt.Sprintf("%s\nrule: %s\nviolation: %s\nfix: %s", code, e.rule, e.violate, e.fix), true
	}
	fam := ""
	if len(code) >= 4 {
		fam = explainFamily[code[:4]]
	}
	if fam == "" {
		fam = "No family text."
	}
	return fmt.Sprintf("%s\n%s No per-code doc yet; see the family rule and the diagnostic message.", code, fam), true
}

// runExplain implements `ailc explain CODE`.
func runExplain(out io.Writer, argv []string) int {
	if len(argv) != 1 {
		fmt.Fprintln(out, "usage: ailc explain AILnnnn")
		return 2
	}
	code := strings.ToUpper(strings.TrimSpace(argv[0]))
	text, ok := explainCode(code)
	if !ok {
		fmt.Fprintf(out, "ailc: unknown diagnostic code %q\n", argv[0])
		return 2
	}
	fmt.Fprintln(out, text)
	return 0
}
