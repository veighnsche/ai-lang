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

// Platform crypto primitives (ambiently declared in platform.d.ts
// because the tsc gate runs without Node libs; node resolves the
// real builtin at runtime). Randomness is always node:crypto —
// never Math.random — and hashing is always the platform digest,
// never a vendored reimplementation. Both are synchronous,
// matching extern calls.
import { randomBytes, createHash } from "node:crypto";

// Upper bound for one random_bytes call: past it the host reports
// invalid_count (the negative half rejects deterministically in
// Can). Mirrors the capped decision-table row.
const RANDOM_MAX_BYTES = 1048576;

// Revisioned algorithm contracts: exact names only. Anything else
// reports unsupported_profile; createHash never sees it.
const HASH_ALGORITHMS = ["sha256", "sha512"] as const;

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

export function host__rand_bytes(count: bigint):
  | { $can_kind: "ok"; bytes: Uint8Array }
  | { $can_kind: "random.invalid_count" }
  | { $can_kind: "random.unavailable" } {
  try {
    if (count < 0n || count > BigInt(RANDOM_MAX_BYTES)) {
      return { $can_kind: "random.invalid_count" };
    }
    return {
      $can_kind: "ok",
      bytes: new Uint8Array(randomBytes(Number(count))),
    };
  } catch {
    return { $can_kind: "random.unavailable" };
  }
}

export function host__hash_digest(
  data: Uint8Array,
  profile: string,
):
  | { $can_kind: "ok"; bytes: Uint8Array }
  | { $can_kind: "hash.unsupported_profile" }
  | { $can_kind: "hash.unavailable" } {
  try {
    if (!(HASH_ALGORITHMS as readonly string[]).includes(profile)) {
      return { $can_kind: "hash.unsupported_profile" };
    }
    return {
      $can_kind: "ok",
      bytes: new Uint8Array(createHash(profile).update(data).digest()),
    };
  } catch {
    return { $can_kind: "hash.unavailable" };
  }
}
