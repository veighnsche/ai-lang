// GENERATED from reserve.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { std__validate__int_range, std__validate__str_length, std__validate__str_one_of, type QuotaResult } from "./quota";
import { std__convert__int_to_str, type ScalarsResult } from "./scalars";
export type ReserveResult = { $can_kind: "ok"; guest: string; party: bigint; slot: string } | { $can_kind: "ok"; value: Reserve__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.invalid_length"; value: string; minimum: bigint; maximum: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } | { $can_kind: "validation.not_allowed"; value: string };
export type Reserve__Request = { guest: string; party: bigint; slot: string };
export type Reserve__RequestSchema = { guest_minimum: bigint; guest_maximum: bigint; party_lower: bigint; party_upper: bigint; allowed_slots: string[] };
export type Reserve__RequestValue = { value: Reserve__Request };
export type Reserve__Confirmation = { guest: string; party: bigint; slot: string };
export function reserve__request__validate(request: Reserve__Request, schema: Reserve__RequestSchema): { $can_kind: "ok"; value: Reserve__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  if ((schema.guest_minimum >= 0n)) {
    const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.invalid_length"; value: string; minimum: bigint; maximum: bigint } = std__validate__str_length(request.guest, schema.guest_minimum, schema.guest_maximum);
    switch ($can_m1.$can_kind) {
    case "validation.invalid_bounds": {
      const _ = $can_m1;
      const $can_m2: { $can_kind: "ok"; value: string } = std__convert__int_to_str(schema.guest_minimum);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const lo = $can_m2;
        const $can_m3: { $can_kind: "ok"; value: string } = std__convert__int_to_str(schema.guest_maximum);
        switch ($can_m3.$can_kind) {
        case "ok": {
          const hi = $can_m3;
          return { $can_kind: "validation.schema_violation", path: "", rule: "schema.guest_bounds", value: ((("minimum=" + lo.value) + ";maximum=") + hi.value) };
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
    case "validation.invalid_length": {
      const _ = $can_m1;
      return { $can_kind: "validation.schema_violation", path: "guest", rule: "str.length_scalars", value: request.guest };
    }
    case "ok": {
      const _ = $can_m1;
      const $can_m4: { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.invalid_bounds"; lower: bigint; upper: bigint } | { $can_kind: "validation.out_of_range"; value: bigint; lower: bigint; upper: bigint } = std__validate__int_range(request.party, schema.party_lower, schema.party_upper);
      switch ($can_m4.$can_kind) {
      case "validation.invalid_bounds": {
        const _ = $can_m4;
        const $can_m5: { $can_kind: "ok"; value: string } = std__convert__int_to_str(schema.party_lower);
        switch ($can_m5.$can_kind) {
        case "ok": {
          const lo = $can_m5;
          const $can_m6: { $can_kind: "ok"; value: string } = std__convert__int_to_str(schema.party_upper);
          switch ($can_m6.$can_kind) {
          case "ok": {
            const hi = $can_m6;
            return { $can_kind: "validation.schema_violation", path: "", rule: "schema.party_bounds", value: ((("lower=" + lo.value) + ";upper=") + hi.value) };
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
      case "validation.out_of_range": {
        const _ = $can_m4;
        const $can_m7: { $can_kind: "ok"; value: string } = std__convert__int_to_str(request.party);
        switch ($can_m7.$can_kind) {
        case "ok": {
          const r = $can_m7;
          return { $can_kind: "validation.schema_violation", path: "party", rule: "int.closed_range", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      case "ok": {
        const _ = $can_m4;
        const $can_m8: { $can_kind: "ok"; value: string } | { $can_kind: "validation.not_allowed"; value: string } = std__validate__str_one_of(request.slot, schema.allowed_slots);
        switch ($can_m8.$can_kind) {
        case "validation.not_allowed": {
          const _ = $can_m8;
          return { $can_kind: "validation.schema_violation", path: "slot", rule: "str.one_of", value: request.slot };
        }
        case "ok": {
          const _ = $can_m8;
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
  else {
    const $can_m9: { $can_kind: "ok"; value: string } = std__convert__int_to_str(schema.guest_minimum);
    switch ($can_m9.$can_kind) {
    case "ok": {
      const r = $can_m9;
      return { $can_kind: "validation.schema_violation", path: "", rule: "schema.guest_minimum", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
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
