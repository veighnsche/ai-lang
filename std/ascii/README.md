# ascii — named ASCII scalar bounds

- `ascii.ail` — `mod ascii`: `std__ascii__HASH`, `PLUS`, `MINUS`,
  `DOT`, `SLASH`, `ZERO`, `NINE`, `COLON`, `QUESTION`, `AT`, `A`,
  `Z`, `BACKTICK`, `SMALL_A`, `SMALL_Z` (lowercase needs marking;
  bare `A`/`Z` are uppercase). Consumers pin
  them like any `uses` entry (`std__ascii__COLON@1`).
  `html__url__scheme_token` compares against these names; the
  emitted program is byte-identical to the magic literals.
  Predicates (`std__ascii__is_digit`, `is_alpha`, `is_alnum`)
  sit beside the constants, returning the scalars
  `Bool__Value` wrapper by pin. Digit and alpha are leaves
  over `and`/`or` of comparisons; alnum composes explicitly
  bound call results (no operand calls), alpha outer.
