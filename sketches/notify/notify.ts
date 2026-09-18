// GENERATED from notify.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
import { std__validate__length_check, std__validate__membership_check, std__validate__range_check, type QuotaResult } from "./quota";
export type NotifyResult = { $can_kind: "ok"; channel: string; handle: string; retries: bigint } | { $can_kind: "ok"; value: Notify__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string };
export type Notify__Request = { handle: string; retries: bigint; channel: string };
export type Notify__RequestSchema = { handle_minimum: bigint; handle_maximum: bigint; retries_lower: bigint; retries_upper: bigint; allowed_channels: string[] };
export type Notify__RequestValue = { value: Notify__Request };
export type Notify__Confirmation = { handle: string; retries: bigint; channel: string };
export function notify__request__validate(request: Notify__Request, schema: Notify__RequestSchema): { $can_kind: "ok"; value: Notify__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  const $can_m1: { $can_kind: "ok"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__length_check(request.handle, { path: "handle", tag: "handle", minimum: schema.handle_minimum, maximum: schema.handle_maximum });
  switch ($can_m1.$can_kind) {
  case "validation.schema_violation": {
    const e = $can_m1;
    return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
  }
  case "ok": {
    const _ = $can_m1;
    const $can_m2: { $can_kind: "ok"; value: bigint } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__range_check(request.retries, { path: "retries", tag: "retries", lower: schema.retries_lower, upper: schema.retries_upper });
    switch ($can_m2.$can_kind) {
    case "validation.schema_violation": {
      const e = $can_m2;
      return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
    }
    case "ok": {
      const _ = $can_m2;
      const $can_m3: { $can_kind: "ok"; value: string } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = std__validate__membership_check(request.channel, { path: "channel", allowed: schema.allowed_channels });
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
export function notify__request__admit(request: Notify__Request, schema: Notify__RequestSchema): { $can_kind: "ok"; channel: string; handle: string; retries: bigint } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } {
  const $can_m1: { $can_kind: "ok"; value: Notify__Request } | { $can_kind: "validation.schema_violation"; path: string; rule: string; value: string } = notify__request__validate(request, schema);
  switch ($can_m1.$can_kind) {
  case "validation.schema_violation": {
    const e = $can_m1;
    return { $can_kind: "validation.schema_violation", path: e.path, rule: e.rule, value: e.value };
  }
  case "ok": {
    const v = $can_m1;
    return { $can_kind: "ok", handle: v.value.handle, retries: v.value.retries, channel: v.value.channel };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
