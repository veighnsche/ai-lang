// Trusted synchronous host implementations for the executable sketch.
// Can tests script this boundary; Node parity exercises these implementations.
import type { Data__Choice$T$int, Data__Row } from "./data";

type Ok<T> = { $can_kind: "ok"; value: T };
export function bridge__int(allow: boolean): Ok<bigint> | { $can_kind: "bridge.denied" } {
  return allow ? { $can_kind: "ok", value: 9007199254740993n } : { $can_kind: "bridge.denied" };
}
export function bridge__str(): Ok<string> {
  return { $can_kind: "ok", value: "hello" };
}
export function bridge__bool(): Ok<boolean> {
  return { $can_kind: "ok", value: false };
}
export function bridge__dec(): Ok<string> {
  return { $can_kind: "ok", value: "0.125" };
}
export function bridge__token(): Ok<string> {
  return { $can_kind: "ok", value: "a" };
}
export function bridge__bytes(): Ok<Uint8Array> {
  return { $can_kind: "ok", value: new Uint8Array([0, 255]) };
}
export function bridge__rows(): Ok<Data__Row[]> {
  return { $can_kind: "ok", value: [{ number: 7n }] };
}
export function bridge__choice(some: boolean): Ok<Data__Choice$T$int> {
  return { $can_kind: "ok", value: some ? { $can_kind: "Data__Some$T$int", value: 7n } : { $can_kind: "Data__None$T$int" } };
}
