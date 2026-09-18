// Extern-stub fixture for the a70 tsc gate. Hand-written; signature
// mirrors `extern net__fetch` in sketches/retry-loop/retry.can and
// its call site in retry.ts. Never executed by the gate (`noEmit`
// check only). If the extern decl or emit shape changes, update this
// stub in the same slice.
export function net__fetch(): { $can_kind: "ok"; body: string } | { $can_kind: "net.down" } {
  throw new Error("extern stub: net__fetch is not implemented");
}
