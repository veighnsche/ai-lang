// GENERATED from reserve.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { std__validate__length_check, std__validate__membership_check, std__validate__range_check, type QuotaResult } from "./quota";
export type ReserveResult = { $can_kind: "ok"; guest: string; party: bigint; slot: string } | { $can_kind: "ok"; value: Reserve__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string };
export type Reserve__Request = { guest: string; party: bigint; slot: string };
export type Reserve__RequestSchema = { guest_minimum: bigint; guest_maximum: bigint; party_lower: bigint; party_upper: bigint; allowed_slots: string[] };
export type Reserve__RequestValue = { value: Reserve__Request };
export type Reserve__Confirmation = { guest: string; party: bigint; slot: string };
export function reserve__request__validate(request: Reserve__Request, schema: Reserve__RequestSchema): { $can_kind: "ok"; value: Reserve__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__length_check(request.guest, { path: "guest", tag: "guest", minimum: schema.guest_minimum, maximum: schema.guest_maximum });
  switch ($can_m1.$can_kind) {
  case "validation.schema_violation": {
    const e = $can_m1;
    return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
  }
  case "ok": {
    const _ = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__range_check(request.party, { path: "party", tag: "party", lower: schema.party_lower, upper: schema.party_upper });
    switch ($can_m2.$can_kind) {
    case "validation.schema_violation": {
      const e = $can_m2;
      return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
    }
    case "ok": {
      const _ = $can_m2;
      const $can_m3: { $can_kind: "ok"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__membership_check(request.slot, { path: "slot", allowed: schema.allowed_slots });
      switch ($can_m3.$can_kind) {
      case "validation.schema_violation": {
        const e = $can_m3;
        return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
      }
      case "ok": {
        const _ = $can_m3;
        return { $can_kind: "ok", value: request };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function reserve__request__admit(request: Reserve__Request, schema: Reserve__RequestSchema): { $can_kind: "ok"; guest: string; party: bigint; slot: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  const $can_m1: { $can_kind: "ok"; value: Reserve__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = reserve__request__validate(request, schema);
  switch ($can_m1.$can_kind) {
  case "validation.schema_violation": {
    const e = $can_m1;
    return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
  }
  case "ok": {
    const v = $can_m1;
    return { $can_kind: "ok", guest: v.value.guest, party: v.value.party, slot: v.value.slot };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
