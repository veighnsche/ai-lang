// GENERATED from form.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type FormResult = { $ail_kind: "ok"; message: string };
export type Form__Draft = { name: string };
export type Form__Out = { message: string };
export type Form__State = { $ail_kind: "Form__Empty" } | { $ail_kind: "Form__Editing"; draft: Form__Draft } | { $ail_kind: "Form__Submitted"; name: string };
export function form__message(state: Form__State): { $ail_kind: "ok"; message: string } {
  const $ail_m1 = state;
  switch ($ail_m1.$ail_kind) {
  case "Form__Empty": {
    const _ = $ail_m1;
    return { $ail_kind: "ok", message: "Start typing" };
  }
  case "Form__Editing": {
    const d = $ail_m1;
    return { $ail_kind: "ok", message: ("Draft: " + d.draft.name) };
  }
  case "Form__Submitted": {
    const s = $ail_m1;
    return { $ail_kind: "ok", message: ("Hello, " + s.name) };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
