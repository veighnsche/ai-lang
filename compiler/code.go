package main

// Diagnostic codes: stable identifiers for every squiggle, forever.
//
// Prose messages may be reworded; codes never change meaning. A code names
// the violated rule, not the phase that found it: the same rule reported by
// the editor and by `ailc --format=json` carries the same code. Families:
//
//	AIL1xxx parse (grammar, braces ban, rev presence)
//	AIL2xxx declarations (naming R3, rev pins R4, provides/uses R2)
//	AIL3xxx calls and decision tables (call resolution, given/test R8)
//	AIL4xxx errors and proof (emits R5, exhaustiveness, test runs)
//	AIL5xxx world and tooling (siblings, unused items)
//	AIL6xxx types (float ban, dec spelling, brands, seal, extern)
//
// Parse failures (AIL1000) and test-execution failures (AIL4200) are
// deliberately coarse: underneath them live dynamic runtime shapes, not
// one rule each. Everything else maps one rule to one code.
//
// When adding a code, append it to allCodes: TestCodesUnique enforces
// registry uniqueness, and TestGoldenJSONDiags freezes every code in
// golden output.
const (
	CodeParse              = "AIL1000"
	CodeFnNaming           = "AIL2001"
	CodeTypeNaming         = "AIL2002"
	CodeUsesPin            = "AIL2101"
	CodeUsesResolve        = "AIL2102"
	CodeUsesRev            = "AIL2103"
	CodeDupFn              = "AIL2201"
	CodeDupSibling         = "AIL2202"
	CodeProvidesMiss       = "AIL2301"
	CodeProvidesGhost      = "AIL2302"
	CodeUnknownCall        = "AIL3001"
	CodeCallNotInUses      = "AIL3002"
	CodeCallOutside        = "AIL3003"
	CodeCallNested         = "AIL3004"
	CodeLocalCycle         = "AIL3005"
	CodeBadDecreases       = "AIL3006"
	CodeStaleDecreases     = "AIL3007"
	CodeNoDecrease         = "AIL3008"
	CodeNoGuard            = "AIL3009"
	CodeBadBinding         = "AIL3010"
	CodeUndeclaredEffect   = "AIL3107"
	CodeStaleEffect        = "AIL3108"
	CodeNoGiven            = "AIL3101"
	CodeGivenOnLocal       = "AIL3106"
	CodeBadStub            = "AIL3102"
	CodeStubNotInEmit      = "AIL3103"
	CodeDeadScript         = "AIL3104"
	CodeDanglingTest       = "AIL3105"
	CodeNoExchange         = "AIL3109"
	CodeInconsistentScript = "AIL3110"
	CodeDupTest            = "AIL3201"
	CodeUnknownArg         = "AIL3202"
	CodeMissingArg         = "AIL3203"
	CodeBareErrorKind      = "AIL3204"
	CodeMissingTests       = "AIL3301"
	CodeUnusedUses         = "AIL3401"
	CodeUnusedParam        = "AIL3402"
	CodeForeignRaise       = "AIL4001"
	CodeUnknownKind        = "AIL4002"
	CodeProofOther         = "AIL4100"
	CodeMissingArm         = "AIL4101"
	CodeStaleArm           = "AIL4102"
	CodeBoolArms           = "AIL4103"
	CodeValueNoWild        = "AIL4104"
	CodeBadArmKind         = "AIL4105"
	CodeVariantOnVal       = "AIL4106"
	CodeArmUntaken         = "AIL4107"
	CodeInvalidRelay       = "AIL4108"
	CodeEagerScrutinee     = "AIL4109"
	CodeTestFailed         = "AIL4200"
	CodeSiblingParse       = "AIL5001"
	CodeBadCompare         = "AIL5005"
	CodeModuleCollision    = "AIL5007"
	CodeFloatLiteral       = "AIL6001"
	CodeUnknownType        = "AIL6002"
	CodeTypeMismatch       = "AIL6003"
	CodeSealForeign        = "AIL6004"
	CodeInexactDivision    = "AIL6005"
	CodeRecordCycle        = "AIL6006"
)

// allCodes mirrors every constant above; the test fails on duplicates,
// which is the only way a reused code can slip in unnoticed.
var allCodes = []string{
	CodeParse,
	CodeFnNaming, CodeTypeNaming,
	CodeUsesPin, CodeUsesResolve, CodeUsesRev,
	CodeModuleCollision,
	CodeDupFn, CodeDupSibling,
	CodeProvidesMiss, CodeProvidesGhost,
	CodeUnknownCall, CodeCallNotInUses, CodeCallOutside, CodeCallNested, CodeLocalCycle,
	CodeBadDecreases, CodeStaleDecreases, CodeNoDecrease, CodeNoGuard, CodeBadBinding,
	CodeUndeclaredEffect, CodeStaleEffect,
	CodeNoGiven, CodeGivenOnLocal, CodeBadStub, CodeStubNotInEmit, CodeDeadScript, CodeDanglingTest, CodeNoExchange, CodeInconsistentScript,
	CodeDupTest, CodeUnknownArg, CodeMissingArg, CodeBareErrorKind,
	CodeMissingTests, CodeUnusedUses, CodeUnusedParam,
	CodeForeignRaise, CodeUnknownKind,
	CodeProofOther, CodeMissingArm, CodeStaleArm, CodeBoolArms, CodeValueNoWild, CodeBadArmKind, CodeVariantOnVal, CodeArmUntaken, CodeInvalidRelay, CodeEagerScrutinee,
	CodeTestFailed,
	CodeSiblingParse, CodeBadCompare,
	CodeFloatLiteral, CodeUnknownType, CodeTypeMismatch, CodeSealForeign, CodeInexactDivision, CodeRecordCycle,
}
