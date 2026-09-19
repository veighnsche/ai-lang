package main

// Diagnostic codes: stable identifiers for every squiggle, forever.
//
// Prose messages may be reworded; codes never change meaning. A code names
// the violated rule, not the phase that found it: the same rule reported by
// the editor and by `canlc --format=json` carries the same code. Families:
//
//	CAN1xxx parse (grammar, braces ban, rev presence)
//	CAN2xxx declarations (naming R3, rev pins R4, provides/uses R2)
//	CAN3xxx calls and decision tables (call resolution, given/test R8)
//	CAN4xxx errors and proof (emits R5, exhaustiveness, test runs)
//	CAN5xxx world and tooling (siblings, unused items)
//	CAN6xxx types (float ban, dec spelling, brands, seal, extern)
//
// Parse failures (CAN1000) and test-execution failures (CAN4200) are
// deliberately coarse: underneath them live dynamic runtime shapes, not
// one rule each. Everything else maps one rule to one code.
//
// When adding a code, append it to allCodes: TestCodesUnique enforces
// registry uniqueness, and TestGoldenJSONDiags freezes every code in
// golden output.
const (
	CodeParse                    = "CAN1000"
	CodeFnNaming                 = "CAN2001"
	CodeTypeNaming               = "CAN2002"
	CodeConstNaming              = "CAN2003"
	CodeUsesPin                  = "CAN2101"
	CodeUsesResolve              = "CAN2102"
	CodeUsesRev                  = "CAN2103"
	CodeUnknownConst             = "CAN2104"
	CodeConstNotInUses           = "CAN2105"
	CodeDupFn                    = "CAN2201"
	CodeDupSibling               = "CAN2202"
	CodeDupVariant               = "CAN2203"
	CodeCaseCollision            = "CAN2204"
	CodeDupConst                 = "CAN2205"
	CodeProvidesMiss             = "CAN2301"
	CodeProvidesGhost            = "CAN2302"
	CodeUnknownCall              = "CAN3001"
	CodeCallNotInUses            = "CAN3002"
	CodeCallOutside              = "CAN3003"
	CodeCallNested               = "CAN3004"
	CodeLocalCycle               = "CAN3005"
	CodeBadDecreases             = "CAN3006"
	CodeStaleDecreases           = "CAN3007"
	CodeNoDecrease               = "CAN3008"
	CodeNoGuard                  = "CAN3009"
	CodeBadBinding               = "CAN3010"
	CodeBadForward               = "CAN3011"
	CodeChainElab                = "CAN3012"
	CodeBadForwardCall           = "CAN3013"
	CodeGenericExpand            = "CAN3014"
	CodeUndeclaredEffect         = "CAN3107"
	CodeStaleEffect              = "CAN3108"
	CodeNoGiven                  = "CAN3101"
	CodeGivenOnLocal             = "CAN3106"
	CodeBadStub                  = "CAN3102"
	CodeStubNotInEmit            = "CAN3103"
	CodeDeadScript               = "CAN3104"
	// CAN3105 (CodeDanglingTest) retired by a91: omission is legal,
	// so "test with no script" is no longer a diagnostic. Numbers
	// are never reused.
	CodeNoExchange               = "CAN3109"
	CodeInconsistentScript       = "CAN3110"
	CodeGivenDashRetired         = "CAN3111"
	CodeGivenOnInvoke            = "CAN3112"
	CodeDupTest                  = "CAN3201"
	CodeUnknownArg               = "CAN3202"
	CodeMissingArg               = "CAN3203"
	CodeBareErrorKind            = "CAN3204"
	CodeBadPositional            = "CAN3205"
	CodeMissingTests             = "CAN3301"
	CodeUnusedUses               = "CAN3401"
	CodeUnusedParam              = "CAN3402"
	CodeLintRedundant            = "CAN3410"
	CodeLintOrFold               = "CAN3411"
	CodeLintChain                = "CAN3412"
	CodeLintTable                = "CAN3413"
	CodeLintSameOutcome          = "CAN3414"
	CodeLintRelay                = "CAN3415"
	CodeLintRange                = "CAN3416"
	CodeLintRestate              = "CAN3417"
	CodeLintLadder               = "CAN3418"
	CodeLintRelayCall            = "CAN3419"
	CodeLintUnreachedKey         = "CAN3420"
	CodeForeignRaise             = "CAN4001"
	CodeUnknownKind              = "CAN4002"
	CodeProofOther               = "CAN4100"
	CodeMissingArm               = "CAN4101"
	CodeStaleArm                 = "CAN4102"
	CodeBoolArms                 = "CAN4103"
	CodeValueNoWild              = "CAN4104"
	CodeBadArmKind               = "CAN4105"
	CodeVariantOnVal             = "CAN4106"
	CodeArmUntaken               = "CAN4107"
	CodeInvalidRelay             = "CAN4108"
	CodeEagerScrutinee           = "CAN4109"
	CodeBadRange                 = "CAN4110"
	CodeUselessArm               = "CAN4111"
	CodeUselessAlt               = "CAN4112"
	CodeTestFailed               = "CAN4200"
	CodeSiblingParse             = "CAN5001"
	CodeBadCompare               = "CAN5005"
	CodeModuleCollision          = "CAN5007"
	CodeFloatLiteral             = "CAN6001"
	CodeUnknownType              = "CAN6002"
	CodeTypeMismatch             = "CAN6003"
	CodeSealForeign              = "CAN6004"
	CodeInexactDivision          = "CAN6005"
	CodeRecordCycle              = "CAN6006"
	CodeSeqLiteral               = "CAN6007"
	CodeBytesLiteral             = "CAN6008"
	CodeBytesElementRange        = "CAN6009"
	CodeBytesExportAuthority     = "CAN6010"
	CodeBytesExportShape         = "CAN6011"
	CodeAssetBridgeAuthority     = "CAN6014"
	CodeAssetBridgeShape         = "CAN6015"
	CodeConstNonliteral          = "CAN6016"
	CodePinnedWeakened           = "CAN6017"
	// CAN6018 (CodeFnValueDeferred) retired by b00 I2: invocation
	// proves statically now, so "deferred" is no longer a diagnostic.
	// Numbers are never reused.
	CodeFnHeadInvalid            = "CAN6019"
	CodeFnResidualArity          = "CAN6020"
	CodeFnComputedCapture        = "CAN6021"
	CodeFnContainment            = "CAN6022"
	CodeFnTargetRefused          = "CAN6023"
	CodePrimitiveShadow          = "CAN6012"
	CodeRevisionIdentity         = "CAN6013"
	CodeContractMalformed        = "CAN4301"
	CodeContractUnsupported      = "CAN4302"
	CodeContractUnverifiedDep    = "CAN4303"
	CodeContractUnproven         = "CAN4304"
	CodeContractInconclusive     = "CAN4305"
	CodeContractInadmissibleTest = "CAN4306"
)

// allCodes registers every constant above in numeric sequence;
// TestCodesUnique fails on duplicates (the only way a reused code
// slips in unnoticed) and TestAllCodesSequenced fails on disorder,
// so gaps and collisions surface at a glance.
var allCodes = []string{
	CodeParse,
	CodeFnNaming, CodeTypeNaming, CodeConstNaming,
	CodeUsesPin, CodeUsesResolve, CodeUsesRev, CodeUnknownConst, CodeConstNotInUses,
	CodeDupFn, CodeDupSibling, CodeDupVariant, CodeCaseCollision, CodeDupConst,
	CodeProvidesMiss, CodeProvidesGhost,
	CodeUnknownCall, CodeCallNotInUses, CodeCallOutside, CodeCallNested, CodeLocalCycle,
	CodeBadDecreases, CodeStaleDecreases, CodeNoDecrease, CodeNoGuard, CodeBadBinding, CodeBadForward, CodeChainElab, CodeBadForwardCall, CodeGenericExpand,
	CodeNoGiven, CodeBadStub, CodeStubNotInEmit, CodeDeadScript, CodeGivenOnLocal,
	CodeUndeclaredEffect, CodeStaleEffect, CodeNoExchange, CodeInconsistentScript, CodeGivenDashRetired, CodeGivenOnInvoke,
	CodeDupTest, CodeUnknownArg, CodeMissingArg, CodeBareErrorKind, CodeBadPositional,
	CodeMissingTests, CodeUnusedUses, CodeUnusedParam,
	CodeLintRedundant, CodeLintOrFold, CodeLintChain, CodeLintTable,
	CodeLintSameOutcome, CodeLintRelay, CodeLintRange, CodeLintRestate, CodeLintLadder, CodeLintRelayCall, CodeLintUnreachedKey,
	CodeForeignRaise, CodeUnknownKind,
	CodeProofOther, CodeMissingArm, CodeStaleArm, CodeBoolArms, CodeValueNoWild, CodeBadArmKind, CodeVariantOnVal, CodeArmUntaken, CodeInvalidRelay, CodeEagerScrutinee, CodeBadRange, CodeUselessArm, CodeUselessAlt,
	CodeTestFailed,
	CodeContractMalformed, CodeContractUnsupported, CodeContractUnverifiedDep,
	CodeContractUnproven, CodeContractInconclusive, CodeContractInadmissibleTest,
	CodeSiblingParse, CodeBadCompare, CodeModuleCollision,
	CodeFloatLiteral, CodeUnknownType, CodeTypeMismatch, CodeSealForeign, CodeInexactDivision, CodeRecordCycle, CodeSeqLiteral,
	CodeBytesLiteral, CodeBytesElementRange, CodeBytesExportAuthority, CodeBytesExportShape, CodePrimitiveShadow, CodeRevisionIdentity,
	CodeAssetBridgeAuthority, CodeAssetBridgeShape, CodeConstNonliteral, CodePinnedWeakened,
	CodeFnHeadInvalid,
	CodeFnResidualArity,
	CodeFnComputedCapture,
	CodeFnContainment,
	CodeFnTargetRefused,
}
