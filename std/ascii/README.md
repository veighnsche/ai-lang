# ascii — named ASCII scalar bounds

- `ascii.ail` — `mod ascii`: `std__ascii__PLUS`, `MINUS`, `DOT`,
  `SLASH`, `NINE`, `COLON`, `AT`, `Z`, `BACKTICK`, `SMALL_Z`
  (lowercase z needs marking; `Z` is uppercase). Consumers pin
  them like any `uses` entry (`std__ascii__COLON@1`).
  `html__url__scheme_token` compares against these names; the
  emitted program is byte-identical to the magic literals.
  Predicates (`is_digit`, `is_alpha`, `is_alnum`) join this
  module with the boolean-operator slice, which needs them.
