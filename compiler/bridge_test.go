package main

// Asset bridge falsifier (S2 slice plan): the grant, the certificate,
// and the kernel each fail loud on exactly their own violation. The
// .ail rows in std/html prove the builder behavior; these tests prove
// the compiler machinery around it.

import (
	"os"
	"strings"
	"testing"
)

// bridgeSink is a minimal granted sink: fixed assembly from projected
// fields with the granted role literal gating. Both params are used
// (the kernel takes asset and policy positionally).
const bridgeSink = `mod sink
  provides [sink__css, Sink__Res, Html__Safe]
  uses []
  emits [sink.rejected]

error sink.rejected(asset: Schema__ApprovedAsset)

type Sink__Res rev 1 (
  safe: Html__Safe
)

brand Html__Safe is str rev 1

asset_bridge Schema__ApprovedAsset, Schema__AssetPolicy from schema via sink__css@1 for stylesheet

fn sink__css(asset: Schema__ApprovedAsset, policy: Schema__AssetPolicy) -> Sink__Res rev 1
  emits [sink.rejected]
  tests
    css(asset = seal Schema__ApprovedAsset("a|b|https://h/x|d|stylesheet|p|m|f"), policy = seal Schema__AssetPolicy("p|q")) => Ok(safe = seal Html__Safe("x"))
    css_role_mismatch(asset = seal Schema__ApprovedAsset("a|b|https://h/x|d|script|p|m|f"), policy = seal Schema__AssetPolicy("p|q")) => sink.rejected(asset = seal Schema__ApprovedAsset("a|b|https://h/x|d|script|p|m|f"))
=
  match call schema__asset__fields(asset, policy)
    on Ok f => match f.role
      "stylesheet" => Ok(safe = seal Html__Safe("x"))
      _ => sink.rejected(asset = asset)
`

func bridgeProgram(t *testing.T, files map[string]string) ([]*Module, []Diag) {
	t.Helper()
	var mods []*Module
	texts := map[string]string{}
	for name, src := range files {
		m, err := parseModuleText(name, src)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		mods = append(mods, m)
		texts[m.ID] = src
	}
	_, collected := checkProgram(mods, texts, nil, nil)
	return mods, collected
}

func bridgeErrs(collected []Diag) []string {
	var out []string
	for _, d := range collected {
		if d.Sev == "error" {
			out = append(out, d.Code+" "+d.Msg)
		}
	}
	return out
}

func wantBridgeDiag(t *testing.T, collected []Diag, code, sub string) {
	t.Helper()
	for _, d := range collected {
		if d.Sev == "error" && d.Code == code && strings.Contains(d.Msg, sub) {
			return
		}
	}
	t.Fatalf("expected %s containing %q, got %v", code, sub, bridgeErrs(collected))
}

func wantBridgeClean(t *testing.T, collected []Diag) {
	t.Helper()
	if errs := bridgeErrs(collected); len(errs) > 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

// TestAssetBridgeStandalone pins the golden path: a granted sink
// compiles alone on the grant's claim (schema need not load).
func TestAssetBridgeStandalone(t *testing.T) {
	_, collected := bridgeProgram(t, map[string]string{"sink.ail": bridgeSink})
	wantBridgeClean(t, collected)
}

// TestAssetBridgeUngranted pins the certificate refusal: a kernel
// call in a sink no grant names fails at check time.
func TestAssetBridgeUngranted(t *testing.T) {
	src := `mod lone
  provides [lone__css]
  uses []
  emits []

type Lone__Res rev 1 (
  safe: Html__Safe
)

brand Html__Safe is str rev 1

fn lone__css(asset: str) -> Lone__Res rev 1
  emits []
  tests
    css(asset = "a|b|c|d|e|f|g|h") => Ok(safe = seal Html__Safe("x"))
=
  match call schema__asset__fields(asset, asset)
    on Ok f => Ok(safe = seal Html__Safe("x"))
`
	_, collected := bridgeProgram(t, map[string]string{"lone.ail": src})
	wantBridgeDiag(t, collected, CodeAssetBridgeAuthority, "not authorized by a valid asset_bridge grant")
}

// TestAssetBridgeShape pins the sink predicate: a sink with a
// raw-string parameter where the witness belongs fails shape.
func TestAssetBridgeShape(t *testing.T) {
	src := `mod wide
  provides [wide__css]
  uses []
  emits [wide.rejected]

error wide.rejected(asset: str)

type Wide__Res rev 1 (
  safe: Html__Safe
)

brand Html__Safe is str rev 1

asset_bridge Schema__ApprovedAsset, Schema__AssetPolicy from schema via wide__css@1 for stylesheet

fn wide__css(asset: str, policy: str) -> Wide__Res rev 1
  emits [wide.rejected]
  tests
    css(asset = "u", policy = "p|q") => Ok(safe = seal Html__Safe("x"))
=
  match call schema__asset__fields(asset, policy)
    on Ok f => match f.role
      "stylesheet" => Ok(safe = seal Html__Safe("x"))
      _ => wide.rejected(asset = asset)
`
	_, collected := bridgeProgram(t, map[string]string{"wide.ail": src})
	wantBridgeDiag(t, collected, CodeAssetBridgeShape, "must have type Schema__ApprovedAsset")
}

// TestAssetBridgeRole pins the role half of shape: a sink that never
// gates the projected role on the granted literal fails shape.
func TestAssetBridgeRole(t *testing.T) {
	src := `mod norole
  provides [norole__css]
  uses []
  emits [norole.rejected]

error norole.rejected(asset: Schema__ApprovedAsset)

type Norole__Res rev 1 (
  safe: Html__Safe
)

brand Html__Safe is str rev 1

asset_bridge Schema__ApprovedAsset, Schema__AssetPolicy from schema via norole__css@1 for stylesheet

fn norole__css(asset: Schema__ApprovedAsset, policy: Schema__AssetPolicy) -> Norole__Res rev 1
  emits [norole.rejected]
  tests
    css(asset = seal Schema__ApprovedAsset("a|b|https://h/x|d|script|p|m|f"), policy = seal Schema__AssetPolicy("p|q")) => norole.rejected(asset = seal Schema__ApprovedAsset("a|b|https://h/x|d|script|p|m|f"))
=
  match call schema__asset__fields(asset, policy)
    on Ok f => norole.rejected(asset = asset)
`
	_, collected := bridgeProgram(t, map[string]string{"norole.ail": src})
	wantBridgeDiag(t, collected, CodeAssetBridgeShape, "must gate the projected role")
}

// TestAssetBridgeSelfGrant pins two-owner authority: a grant whose
// owner is the granting module itself is void.
func TestAssetBridgeSelfGrant(t *testing.T) {
	src := `mod self
  provides [self__css]
  uses []
  emits [self.rejected]

error self.rejected(asset: Schema__ApprovedAsset)

type Self__Res rev 1 (
  safe: Html__Safe
)

brand Html__Safe is str rev 1
brand Schema__ApprovedAsset is str rev 1
brand Schema__AssetPolicy is str rev 1

asset_bridge Schema__ApprovedAsset, Schema__AssetPolicy from self via self__css@1 for stylesheet

fn self__css(asset: Schema__ApprovedAsset, policy: Schema__AssetPolicy) -> Self__Res rev 1
  emits [self.rejected]
  tests
    css(asset = seal Schema__ApprovedAsset("a|b|https://h/x|d|stylesheet|p|m|f"), policy = seal Schema__AssetPolicy("p|q")) => Ok(safe = seal Html__Safe("x"))
=
  match call schema__asset__fields(asset, policy)
    on Ok f => match f.role
      "stylesheet" => Ok(safe = seal Html__Safe("x"))
      _ => self.rejected(asset = asset)
`
	_, collected := bridgeProgram(t, map[string]string{"self.ail": src})
	wantBridgeDiag(t, collected, CodeAssetBridgeAuthority, "not two-owner")
}

// TestAssetBridgeWrongOwner pins owner verification where loaded: a
// grant claiming an owner that declares neither brand fails.
func TestAssetBridgeWrongOwner(t *testing.T) {
	owner := `mod schema
  provides [Schema__ApprovedAsset, Schema__AssetPolicy]
  uses []
  emits []

brand Schema__ApprovedAsset is str rev 1
brand Schema__AssetPolicy is str rev 1
`
	// The fixture owner uses mod name "schema", so the correct
	// grant stays clean; point the grant elsewhere to fail.
	bad := strings.Replace(bridgeSink, "from schema via", "from impostor via", 1)
	_, collected := bridgeProgram(t, map[string]string{
		"schema.ail": owner,
		"sink.ail":   bad,
	})
	wantBridgeDiag(t, collected, CodeAssetBridgeAuthority, "outside owner module impostor")
}

// TestAssetBridgeUnknownSink pins grant ownership of diagnostics: a
// grant naming a function that resolves nowhere fails authority.
func TestAssetBridgeUnknownSink(t *testing.T) {
	bad := strings.Replace(bridgeSink, "via sink__css@1", "via sink__ghost@1", 1)
	_, collected := bridgeProgram(t, map[string]string{"sink.ail": bad})
	wantBridgeDiag(t, collected, CodeAssetBridgeAuthority, "unknown function sink__ghost")
}

// TestAssetBridgeMalformedWitness pins kernel refusal: a sealed but
// ill-formed witness (wrong field count) fails loud at runtime, never
// as a language outcome. Malformed is failed tooling, not a rejected
// asset: well-formed rejections stay declared builder errors.
func TestAssetBridgeMalformedWitness(t *testing.T) {
	owner := `mod schema
  provides [Schema__ApprovedAsset, Schema__AssetPolicy]
  uses []
  emits []

brand Schema__ApprovedAsset is str rev 1
brand Schema__AssetPolicy is str rev 1
`
	files := map[string]string{"schema.ail": owner, "sink.ail": bridgeSink}
	err := runLinkedPure(t, files, []string{"schema.ail", "sink.ail"}, "sink__css", 1,
		map[string]string{
			"asset":  `seal Schema__ApprovedAsset("only|three|parts")`,
			"policy": `seal Schema__AssetPolicy("p|q")`,
		}, `Ok(safe = "x")`)
	if err == nil {
		t.Fatal("expected malformed-witness refusal, got pass")
	}
	if !strings.Contains(err.Error(), "well-formed approval witness") {
		t.Fatalf("refusal must name the malformed witness, got %v", err)
	}
}

// TestAssetBridgeMalformedPolicy pins the handle shape check: a
// policy handle without admit form fails loud at runtime.
func TestAssetBridgeMalformedPolicy(t *testing.T) {
	owner := `mod schema
  provides [Schema__ApprovedAsset, Schema__AssetPolicy]
  uses []
  emits []

brand Schema__ApprovedAsset is str rev 1
brand Schema__AssetPolicy is str rev 1
`
	files := map[string]string{"schema.ail": owner, "sink.ail": bridgeSink}
	err := runLinkedPure(t, files, []string{"schema.ail", "sink.ail"}, "sink__css", 1,
		map[string]string{
			"asset":  `seal Schema__ApprovedAsset("a|b|https://h/x|d|stylesheet|p|m|f")`,
			"policy": `seal Schema__AssetPolicy("bare")`,
		}, `Ok(safe = "x")`)
	if err == nil {
		t.Fatal("expected malformed-policy refusal, got pass")
	}
	if !strings.Contains(err.Error(), "admitted policy handle") {
		t.Fatalf("refusal must name the malformed handle, got %v", err)
	}
}

// TestAssetBridgeRealModules pins no drift between the fixture shape
// and the shipped modules: the real stylesheet sink projects across
// the real owner module in both file orders.
func TestAssetBridgeRealModules(t *testing.T) {
	htmlSrc, err := os.ReadFile("../std/html/html.ail")
	if err != nil {
		t.Fatal(err)
	}
	schemaSrc, err := os.ReadFile("../std/schema/schema.ail")
	if err != nil {
		t.Fatal(err)
	}
	// Slice 1: html.ail pins ascii consts; the provider loads
	// beside the real modules.
	asciiSrc, err := os.ReadFile("../std/ascii/ascii.ail")
	if err != nil {
		t.Fatal(err)
	}
	digest := "sha384-" + strings.Repeat("A", 64)
	witness := "app-css|1.0.0|https://cdn.example/app.css|" + digest + "|stylesheet|shop|pages|home"
	element := "<link rel='stylesheet' href='https://cdn.example/app.css' integrity='" + digest + "' crossorigin='anonymous'>"
	files := map[string]string{"html.ail": string(htmlSrc), "schema.ail": string(schemaSrc), "ascii.ail": string(asciiSrc)}
	for _, order := range [][]string{
		{"schema.ail", "html.ail", "ascii.ail"},
		{"html.ail", "schema.ail", "ascii.ail"},
	} {
		err := runLinkedPure(t, files, order, "html__asset__stylesheet", 1,
			map[string]string{
				"asset":  `seal Schema__ApprovedAsset("` + witness + `")`,
				"policy": `seal Schema__AssetPolicy("shop|prod")`,
			}, `Ok(safe = "`+element+`")`)
		if err != nil {
			t.Fatalf("order %v: %v", order, err)
		}
	}
}

// TestAssetBridgeLinked pins end-to-end projection across modules in
// both orders: the kernel discloses approved fields to the sink.
func TestAssetBridgeLinked(t *testing.T) {
	owner := `mod schema
  provides [Schema__ApprovedAsset, Schema__AssetPolicy]
  uses []
  emits []

brand Schema__ApprovedAsset is str rev 1
brand Schema__AssetPolicy is str rev 1
`
	// The bridgeSink program doubles as the linked root: run it
	// against the loaded owner module in both file orders.
	for _, order := range [][]string{
		{"schema.ail", "sink.ail"},
		{"sink.ail", "schema.ail"},
	} {
		files := map[string]string{"schema.ail": owner, "sink.ail": bridgeSink}
		err := runLinkedPure(t, files, order, "sink__css", 1,
			map[string]string{
				"asset":  `seal Schema__ApprovedAsset("a|b|https://h/x|d|stylesheet|p|m|f")`,
				"policy": `seal Schema__AssetPolicy("p|q")`,
			}, `Ok(safe = "x")`)
		if err != nil {
			t.Fatalf("order %v: %v", order, err)
		}
	}
}
