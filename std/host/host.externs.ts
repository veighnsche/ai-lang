// Host implementations for `extern host__wall_now` and
// `extern host__mono_now` in std/host/host.can. Real observations,
// not stubs: signatures mirror the extern decls (int erases to
// bigint) and their call sites in host.ts. `performance` is
// ambiently declared because the tsc gate runs without DOM/Node
// libs; Node >= 16 provides it globally. Monotonicity is the
// platform's: trussed by TestStdHostNodeSmoke (non-decreasing
// across calls), never by wall time. If an extern decl or emit
// shape changes, update this file in the same slice.

// Minimal ambient for the monotonic clock (Node global).
declare const performance: { now(): number };

export function host__wall_now():
  | { $can_kind: "ok"; millis: bigint }
  | { $can_kind: "clock.unavailable" } {
  return { $can_kind: "ok", millis: BigInt(Date.now()) };
}

export function host__mono_now():
  | { $can_kind: "ok"; millis: bigint }
  | { $can_kind: "clock.unavailable" } {
  return { $can_kind: "ok", millis: BigInt(Math.floor(performance.now())) };
}
