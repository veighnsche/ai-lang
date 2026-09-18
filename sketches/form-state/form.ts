// GENERATED from form.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type FormResult = { $can_kind: "ok"; message: string };
export type Form__Draft = { name: string };
export type Form__Out = { message: string };
export type Form__State = { $can_kind: "Form__Empty" } | { $can_kind: "Form__Editing"; draft: Form__Draft } | { $can_kind: "Form__Submitted"; name: string };
export function form__message(state: Form__State): { $can_kind: "ok"; message: string } {
  const $can_m1 = state;
  switch ($can_m1.$can_kind) {
  case "Form__Empty": {
    const _ = $can_m1;
    return { $can_kind: "ok", message: "Start typing" };
  }
  case "Form__Editing": {
    const d = $can_m1;
    return { $can_kind: "ok", message: ("Draft: " + d.draft.name) };
  }
  case "Form__Submitted": {
    const s = $can_m1;
    return { $can_kind: "ok", message: ("Hello, " + s.name) };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
