// Extern-stub fixture for the v70 tsc gate. Hand-written; signature
// mirrors `extern net__fetch` in sketches/retry-loop/retry.ail and
// its call site in retry.ts. Never executed by the gate (`noEmit`
// check only). If the extern decl or emit shape changes, update this
// stub in the same slice.
export function net__fetch(): { $ail_kind: "ok"; body: string } | { $ail_kind: "net.down" } {
  throw new Error("extern stub: net__fetch is not implemented");
}
