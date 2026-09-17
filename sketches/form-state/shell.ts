// GENERATED from shell.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { form__message, type Form__State, type FormResult } from "./form";
export type ShellResult = { $ail_kind: "ok"; echo: Form__State; message: string };
export type Shell__Out = { message: string; echo: Form__State };
export function shell__greet(state: Form__State): { $ail_kind: "ok"; echo: Form__State; message: string } {
  const $ail_m1: { $ail_kind: "ok"; message: string } = form__message(state);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const v = $ail_m1;
    return { $ail_kind: "ok", message: v.message, echo: state };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
