// Trusted synchronous host implementation; Can rows script, Node tests execute it.
export function use__host(present: boolean):
  | { $can_kind: "ok"; left: bigint; right: boolean }
  | { $can_kind: "option.absent" } {
  return present
    ? { $can_kind: "ok", left: 7n, right: true }
    : { $can_kind: "option.absent" };
}
