# form-state — first variant consumer (a76 pilot)

- `form.ail` — `mod form`: a `Form__State` union (nullary
  `Empty`, `Editing` with a record draft, `Submitted` with a
  scalar) eliminated to a message string. Every arm witnessed
  by the decision table.
- `shell.ail` — `mod shell`: pins the parent, calls
  `form__message` through `uses`, scripts variant-carrying
  exchange args in `given`, and passes the union through
  unexamined (`echo` field) beside the matched result.
- `form.ts` + `shell.ts` + `errors.json` — committed golden
  TS prod emit (tests/given stripped). The catalog carries no
  `Form` entries: data cases are never errors.
