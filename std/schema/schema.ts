// GENERATED from schema.ail by ailc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type SchemaResult = { $ail_kind: "ok"; approval: string } | { $ail_kind: "ok"; count: bigint; first: bigint } | { $ail_kind: "ok"; entry: Schema__RegistryEntry } | { $ail_kind: "ok"; entry: Schema__Revoked } | { $ail_kind: "ok"; policy: string } | { $ail_kind: "ok"; value: boolean } | { $ail_kind: "schema.asset_not_approved"; asset: Schema__AssetRequest };
export type Schema__AssetRequest = { id: string; revision: string; url: string; digest: string; role: string };
export type Schema__Site = { program: string; module: string; function: string };
export type Schema__RegistryEntry = { id: string; revision: string; url: string; digest: string; role: string; program: string; module: string; function: string; not_before: bigint; not_after: bigint };
export type Schema__Revoked = { id: string; revision: string };
export type Schema__RegistrySnapshot = { program: string; policy: string; sequence: bigint; entries: Schema__RegistryEntry[]; revoked: Schema__Revoked[] };
export type Schema__ApprovalResult = { approval: string };
export type Schema__PolicyResult = { policy: string };
export type Schema__BoolResult = { value: boolean };
export type Schema__EntryResult = { entry: Schema__RegistryEntry };
export type Schema__RevokedResult = { entry: Schema__Revoked };
export type Schema__ScanResult = { first: bigint; count: bigint };
// Byte-order string comparison: UTF-8 bytes, matching Go.
function $ailStrCmp(a: string, b: string): number {
  const A = new TextEncoder().encode(a);
  const B = new TextEncoder().encode(b);
  const n = Math.min(A.length, B.length);
  for (let i = 0; i < n; i++) {
    if (A[i] !== B[i]) {
      return A[i] < B[i] ? -1 : 1;
    }
  }
  if (A.length === B.length) {
    return 0;
  }
  return A.length < B.length ? -1 : 1;
}
function $ailStrAt(s: string, i: bigint): bigint {
  const cps = [...s];
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("str index out of range");
  const k = Number(i);
  if (k >= cps.length) throw new Error("str index out of range");
  const cp = cps[k].codePointAt(0);
  if (cp === undefined) throw new Error("str index out of range");
  return BigInt(cp);
}
function $ailStrSlice(s: string, a: bigint, b: bigint): string {
  const cps = [...s];
  const toIdx = (x: bigint): number => { if (x < 0n || x > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("str slice out of range"); return Number(x); };
  const lo = toIdx(a), hi = toIdx(b);
  if (lo > hi || hi > cps.length) throw new Error("str slice out of range");
  return cps.slice(lo, hi).join("");
}
// Sequence indexing (a38 S3): bounds throw, matching Go.
function $ailSeqAt<T>(a: T[], i: bigint): T {
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("seq index out of range");
  const k = Number(i);
  if (k >= a.length) throw new Error("seq index out of range");
  return a[k];
}
export function schema__role__check(role: string): { $ail_kind: "ok"; value: boolean } {
  if (role === "stylesheet") {
    return { $ail_kind: "ok", value: true };
  }
  else if (role === "script") {
    return { $ail_kind: "ok", value: true };
  }
  return { $ail_kind: "ok", value: false };
}
export function schema__url__prefix(url: string): { $ail_kind: "ok"; value: boolean } {
  if (((BigInt([...url].length)) < 8n)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    if (($ailStrSlice(url, 0n, 8n) === "https://")) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
}
export function schema__url__ws(code: bigint): { $ail_kind: "ok"; value: boolean } {
  if ((code === 9n)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    if ((code === 10n)) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      if ((code === 12n)) {
        return { $ail_kind: "ok", value: true };
      }
      else {
        if ((code === 13n)) {
          return { $ail_kind: "ok", value: true };
        }
        else {
          if ((code === 32n)) {
            return { $ail_kind: "ok", value: true };
          }
          else {
            return { $ail_kind: "ok", value: false };
          }
        }
      }
    }
  }
}
export function schema__url__tail_check(s: string, n: bigint, in_authority: boolean): { $ail_kind: "ok"; value: boolean } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    if (($ailStrAt(s, 0n) === 0n)) {
      return { $ail_kind: "ok", value: false };
    }
    else {
      if (($ailStrAt(s, 0n) === 35n)) {
        return { $ail_kind: "ok", value: false };
      }
      else {
        if (($ailStrAt(s, 0n) === 92n)) {
          return { $ail_kind: "ok", value: false };
        }
        else {
          const $ail_m1: { $ail_kind: "ok"; value: boolean } = schema__url__ws($ailStrAt(s, 0n));
          switch ($ail_m1.$ail_kind) {
          case "ok": {
            const w = $ail_m1;
            if (w.value) {
              return { $ail_kind: "ok", value: false };
            }
            else {
              if (($ailStrAt(s, 0n) === 64n)) {
                if (in_authority) {
                  return { $ail_kind: "ok", value: false };
                }
                else {
                  const $ail_m2: { $ail_kind: "ok"; value: boolean } = schema__url__tail_check($ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), in_authority);
                  switch ($ail_m2.$ail_kind) {
                  case "ok": {
                    const r = $ail_m2;
                    return { $ail_kind: "ok", value: r.value };
                  }
                  default: {
                    throw new Error("unreachable");
                  }
                  }
                }
              }
              else {
                if (($ailStrAt(s, 0n) === 47n)) {
                  const $ail_m3: { $ail_kind: "ok"; value: boolean } = schema__url__tail_check($ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), false);
                  switch ($ail_m3.$ail_kind) {
                  case "ok": {
                    const r = $ail_m3;
                    return { $ail_kind: "ok", value: r.value };
                  }
                  default: {
                    throw new Error("unreachable");
                  }
                  }
                }
                else {
                  const $ail_m4: { $ail_kind: "ok"; value: boolean } = schema__url__tail_check($ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n), in_authority);
                  switch ($ail_m4.$ail_kind) {
                  case "ok": {
                    const r = $ail_m4;
                    return { $ail_kind: "ok", value: r.value };
                  }
                  default: {
                    throw new Error("unreachable");
                  }
                  }
                }
              }
            }
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
      }
    }
  }
}
export function schema__digest__char(code: bigint): { $ail_kind: "ok"; value: boolean } {
  if ((code === 43n)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    if ((code === 47n)) {
      return { $ail_kind: "ok", value: true };
    }
    else {
      if ((code === 61n)) {
        return { $ail_kind: "ok", value: true };
      }
      else {
        const $ail_m1 = (code >= 48n);
        const $ail_m2 = (code <= 57n);
        if ($ail_m1 && $ail_m2) {
          return { $ail_kind: "ok", value: true };
        }
        else {
          const $ail_m3 = (code >= 65n);
          const $ail_m4 = (code <= 90n);
          if ($ail_m3 && $ail_m4) {
            return { $ail_kind: "ok", value: true };
          }
          else {
            const $ail_m5 = (code >= 97n);
            const $ail_m6 = (code <= 122n);
            if ($ail_m5 && $ail_m6) {
              return { $ail_kind: "ok", value: true };
            }
            else {
              return { $ail_kind: "ok", value: false };
            }
          }
        }
      }
    }
  }
}
export function schema__digest__tail(s: string, n: bigint): { $ail_kind: "ok"; value: boolean } {
  if ((n <= 0n)) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    const $ail_m1: { $ail_kind: "ok"; value: boolean } = schema__digest__char($ailStrAt(s, 0n));
    switch ($ail_m1.$ail_kind) {
    case "ok": {
      const c = $ail_m1;
      if (c.value) {
        const $ail_m2: { $ail_kind: "ok"; value: boolean } = schema__digest__tail($ailStrSlice(s, 1n, (BigInt([...s].length))), (n - 1n));
        switch ($ail_m2.$ail_kind) {
        case "ok": {
          const r = $ail_m2;
          return { $ail_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
      else {
        return { $ail_kind: "ok", value: false };
      }
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function schema__digest__check(digest: string): { $ail_kind: "ok"; value: boolean } {
  if (((BigInt([...digest].length)) === 71n)) {
    if (($ailStrSlice(digest, 0n, 7n) === "sha384-")) {
      const $ail_m1: { $ail_kind: "ok"; value: boolean } = schema__digest__tail($ailStrSlice(digest, 7n, (BigInt([...digest].length))), 64n);
      switch ($ail_m1.$ail_kind) {
      case "ok": {
        const r = $ail_m1;
        return { $ail_kind: "ok", value: r.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function schema__entry__matches(entry: Schema__RegistryEntry, request: Schema__AssetRequest, site: Schema__Site): { $ail_kind: "ok"; value: boolean } {
  const $ail_m1 = (entry.id === request.id);
  const $ail_m2 = (entry.revision === request.revision);
  const $ail_m3 = (entry.url === request.url);
  const $ail_m4 = (entry.digest === request.digest);
  const $ail_m5 = (entry.role === request.role);
  const $ail_m6 = (entry.program === site.program);
  const $ail_m7 = (entry.module === site.module);
  const $ail_m8 = (entry.function === site.function);
  if ($ail_m1 && $ail_m2 && $ail_m3 && $ail_m4 && $ail_m5 && $ail_m6 && $ail_m7 && $ail_m8) {
    return { $ail_kind: "ok", value: true };
  }
  else {
    return { $ail_kind: "ok", value: false };
  }
}
export function schema__snapshot__scan(entries: Schema__RegistryEntry[], request: Schema__AssetRequest, site: Schema__Site, position: bigint, fuel: bigint): { $ail_kind: "ok"; count: bigint; first: bigint } {
  if ((fuel <= 0n)) {
    return { $ail_kind: "ok", first: position, count: 0n };
  }
  else {
    if ((position < (BigInt([...entries].length)))) {
      const $ail_m1: { $ail_kind: "ok"; entry: Schema__RegistryEntry } = schema__snapshot__entry_at(entries, position);
      switch ($ail_m1.$ail_kind) {
      case "ok": {
        const at = $ail_m1;
        const $ail_m2: { $ail_kind: "ok"; value: boolean } = schema__entry__matches(at.entry, request, site);
        switch ($ail_m2.$ail_kind) {
        case "ok": {
          const m = $ail_m2;
          if (m.value) {
            const $ail_m3: { $ail_kind: "ok"; count: bigint; first: bigint } = schema__snapshot__scan(entries, request, site, (position + 1n), (fuel - 1n));
            switch ($ail_m3.$ail_kind) {
            case "ok": {
              const rest = $ail_m3;
              return { $ail_kind: "ok", first: position, count: (rest.count + 1n) };
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          else {
            const $ail_m4: { $ail_kind: "ok"; count: bigint; first: bigint } = schema__snapshot__scan(entries, request, site, (position + 1n), (fuel - 1n));
            switch ($ail_m4.$ail_kind) {
            case "ok": {
              const rest = $ail_m4;
              return { $ail_kind: "ok", first: rest.first, count: rest.count };
            }
            default: {
              throw new Error("unreachable");
            }
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
      return { $ail_kind: "ok", first: position, count: 0n };
    }
  }
}
export function schema__snapshot__entry_at(entries: Schema__RegistryEntry[], position: bigint): { $ail_kind: "ok"; entry: Schema__RegistryEntry } {
  return { $ail_kind: "ok", entry: $ailSeqAt(entries, position) };
}
export function schema__revoked__contains(revoked: Schema__Revoked[], id: string, revision: string, position: bigint, fuel: bigint): { $ail_kind: "ok"; value: boolean } {
  if ((fuel <= 0n)) {
    return { $ail_kind: "ok", value: false };
  }
  else {
    if ((position < (BigInt([...revoked].length)))) {
      const $ail_m1: { $ail_kind: "ok"; entry: Schema__Revoked } = schema__revoked__entry_at(revoked, position);
      switch ($ail_m1.$ail_kind) {
      case "ok": {
        const at = $ail_m1;
        const $ail_m2 = (at.entry.id === id);
        const $ail_m3 = (at.entry.revision === revision);
        if ($ail_m2 && $ail_m3) {
          return { $ail_kind: "ok", value: true };
        }
        else {
          const $ail_m4: { $ail_kind: "ok"; value: boolean } = schema__revoked__contains(revoked, id, revision, (position + 1n), (fuel - 1n));
          switch ($ail_m4.$ail_kind) {
          case "ok": {
            const r = $ail_m4;
            return { $ail_kind: "ok", value: r.value };
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $ail_kind: "ok", value: false };
    }
  }
}
export function schema__revoked__entry_at(revoked: Schema__Revoked[], position: bigint): { $ail_kind: "ok"; entry: Schema__Revoked } {
  return { $ail_kind: "ok", entry: $ailSeqAt(revoked, position) };
}
export function schema__policy__admit(program: string, policy: string): { $ail_kind: "ok"; policy: string } {
  return { $ail_kind: "ok", policy: ((program + "|") + policy) };
}
export function schema__asset__approve(request: Schema__AssetRequest, snapshot: Schema__RegistrySnapshot, policy: string, site: Schema__Site, head: bigint, t: bigint): { $ail_kind: "ok"; approval: string } | { $ail_kind: "schema.asset_not_approved"; asset: Schema__AssetRequest } {
  const $ail_m1: { $ail_kind: "ok"; value: boolean } = schema__role__check(request.role);
  switch ($ail_m1.$ail_kind) {
  case "ok": {
    const rc = $ail_m1;
    if (rc.value) {
      const $ail_m2: { $ail_kind: "ok"; value: boolean } = schema__url__prefix(request.url);
      switch ($ail_m2.$ail_kind) {
      case "ok": {
        const pc = $ail_m2;
        if (pc.value) {
          const $ail_m3: { $ail_kind: "ok"; value: boolean } = schema__url__tail_check(request.url, (BigInt([...request.url].length)), true);
          switch ($ail_m3.$ail_kind) {
          case "ok": {
            const tc = $ail_m3;
            if (tc.value) {
              const $ail_m4: { $ail_kind: "ok"; value: boolean } = schema__digest__check(request.digest);
              switch ($ail_m4.$ail_kind) {
              case "ok": {
                const dc = $ail_m4;
                if (dc.value) {
                  if ((snapshot.policy === policy)) {
                    if ((snapshot.program === site.program)) {
                      if ((snapshot.sequence === head)) {
                        const $ail_m5: { $ail_kind: "ok"; count: bigint; first: bigint } = schema__snapshot__scan(snapshot.entries, request, site, 0n, ((BigInt([...snapshot.entries].length)) + 1n));
                        switch ($ail_m5.$ail_kind) {
                        case "ok": {
                          const sc = $ail_m5;
                          if ((sc.count === 1n)) {
                            const $ail_m6: { $ail_kind: "ok"; entry: Schema__RegistryEntry } = schema__snapshot__entry_at(snapshot.entries, sc.first);
                            switch ($ail_m6.$ail_kind) {
                            case "ok": {
                              const found = $ail_m6;
                              const $ail_m7: { $ail_kind: "ok"; value: boolean } = schema__revoked__contains(snapshot.revoked, found.entry.id, found.entry.revision, 0n, ((BigInt([...snapshot.revoked].length)) + 1n));
                              switch ($ail_m7.$ail_kind) {
                              case "ok": {
                                const rv = $ail_m7;
                                if (rv.value) {
                                  return { $ail_kind: "schema.asset_not_approved", asset: request };
                                }
                                else {
                                  const $ail_m8 = (found.entry.not_before <= t);
                                  const $ail_m9 = (t < found.entry.not_after);
                                  if ($ail_m8 && $ail_m9) {
                                    return { $ail_kind: "ok", approval: ((((((((((((((request.id + "|") + request.revision) + "|") + request.url) + "|") + request.digest) + "|") + request.role) + "|") + site.program) + "|") + site.module) + "|") + site.function) };
                                  }
                                  else {
                                    return { $ail_kind: "schema.asset_not_approved", asset: request };
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
                            return { $ail_kind: "schema.asset_not_approved", asset: request };
                          }
                        }
                        default: {
                          throw new Error("unreachable");
                        }
                        }
                      }
                      else {
                        return { $ail_kind: "schema.asset_not_approved", asset: request };
                      }
                    }
                    else {
                      return { $ail_kind: "schema.asset_not_approved", asset: request };
                    }
                  }
                  else {
                    return { $ail_kind: "schema.asset_not_approved", asset: request };
                  }
                }
                else {
                  return { $ail_kind: "schema.asset_not_approved", asset: request };
                }
              }
              default: {
                throw new Error("unreachable");
              }
              }
            }
            else {
              return { $ail_kind: "schema.asset_not_approved", asset: request };
            }
          }
          default: {
            throw new Error("unreachable");
          }
          }
        }
        else {
          return { $ail_kind: "schema.asset_not_approved", asset: request };
        }
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
    else {
      return { $ail_kind: "schema.asset_not_approved", asset: request };
    }
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
