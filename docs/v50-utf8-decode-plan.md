# v50: UTF-8 decode kernel (B6 plan v1 — for review)

Goal: `bytes__utf8__decode`, the first FALLIBLE codec kernel:
`Bytes` in, `Encoding__Text` on success, `encoding.invalid_utf8`
on malformed input. B7 (stdlib wrapper) is a separate slice.

## Proposed contract

- Success record: `Encoding__Text rev 1 (value: str)`.
- Error: `encoding.invalid_utf8(value: Bytes)` — carries the
  ORIGINAL payload unchanged (value-carrying fault contract, no
  index/length side-channels).
- `EmitsOf[bytes__utf8__decode] = [encoding.invalid_utf8]`,
  explicit entry with independent existence check (B2/B3 pattern).
- Deterministic: `given` rejected with `CodeGivenOnLocal`, no
  failure injection (B3 pattern).
- Exhaustiveness: exactly `ok ∪ EmitsOf[kernel]`; B6 ships
  missing-arm AND stale-arm rejection tests (no inherited
  `dec__parts` gap).

## Proposed grammar: valid iff Go `utf8.Valid` over the whole input

Invalid classes (each gets a decision-table row): overlong
encodings, encoded surrogate halves (`ED A0..BF`), codepoints
above U+10FFFF (`F4 90..` and beyond), stray continuation bytes,
5- and 6-byte leads (`F8..FF`), truncation at every length
(2-, 3-, 4-byte sequences cut after each byte). Boundary valid
rows: U+007F, U+0080, U+07FF, U+0800, U+FFFF, U+10000, U+10FFFF.

- Empty input decodes to `Ok("")` (valid, not an error).
- BOM decodes to U+FEFF preserved, never stripped.
- NUL byte decodes to U+0000 preserved, never converted
  (encoder-nul-policy applies in this direction too).

## Proposed emission

TS lowers to `new TextDecoder("utf-8", { fatal: true })`.
Parity claim: Go `utf8.Valid` accepts exactly the inputs
fatal-`TextDecoder` accepts. Node adversarial vectors (every
invalid class above + boundaries) must agree byte-for-byte
with Go outcomes — success strings identical, failures on
the same inputs.

## Test plan

- Kernel fixture module (B3 pattern): valid rows incl. all
  boundaries, BOM, NUL-first/middle/last; invalid row per
  class incl. truncated-at-each-length; empty.
- Missing-arm rejection, stale-arm rejection, no-`given`
  rejection, strict admission (non-Bytes input rejected).
- Node adversarial vectors vs Go; `errors.json` of touched
  modules byte-identical unless the new kind lands there
  (see Q1).

## Open questions for the reviewer

- Q1. Home of `Encoding__Text` + `encoding.invalid_utf8`: new
  `std/encoding/encoding.ail` module, or kernel-local
  declarations? Error kinds are module-declared; which home
  keeps the module inventory honest, and what regen churn
  does each choice cause?
- Q2. Are `Bytes`-typed error fields admissible (B1 matrix
  "all value positions")? If not, what carries the payload?
- Q3. Does the Go-`utf8.Valid` ≡ fatal-`TextDecoder` parity
  claim hold on all classes above? Counter-example or
  confirming argument required — not a citation.
- Q4. Truncated input: one `invalid_utf8` for the whole
  payload (proposed), or is per-position reporting owed?
  The proposal carries no index; object if the fault
  contract demands more.
- Q5. Empty input `Ok("")` vs error — confirm the proposal
  against the strict-grammar reading ("complete-input
  errors" cuts both ways).
