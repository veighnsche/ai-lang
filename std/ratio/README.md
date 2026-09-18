# ratio — exact fractions as records

- `ratio.can` — `mod ratio`: `std__ratio__make` (reduced
  fraction, positive denominator, zero as 0/1, rejects a zero
  denominator), `std__ratio__add`/`multiply`/`divide`, and
  `std__ratio__to_dec_exact` (exact decimal or
  `math.nonterminating_decimal`, never hidden rounding).
  Workers: `abs`, `gcd_from` (`decreases a, b by euclid`, step
  `(b, a % b)`), `seen_from` (remainder
  membership over `Seq<int>`), `digits_from` (long division with
  cycle detection; fuel of denominator + 1 always suffices, so
  exhaustion also means nonterminating), `digit_from`
  (int-to-dec digit accumulation), `fold_from` (left-to-right
  fold with the place shrunk by `* d"0.1"` — no dec division
  operator exists). Self-contained: no cross-file calls, per
  the linkage gate. Language brands wrap str scalars only, so
  the pair is a transparent `Ratio__Fraction` record held by
  construction; a forged zero denominator surfaces
  `math.zero_divisor` (hence `add`/`multiply`/`to_dec_exact`
  declare it — the one deviation from `docs/ASTRA_STDLIB.md`).
- `ratio.ts` + `errors.json` — committed golden TS prod emit
  (tests stripped). Regenerate: `go run ./compiler --out
  std/ratio std/ratio/ratio.can`; verify: `go test ./...`.

Rules: `/REQUIREMENTS.md`. Program: `docs/ASTRA_STDLIB.md`
(1.3, exact fractions as records).
