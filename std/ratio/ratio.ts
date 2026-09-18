// GENERATED from ratio.can by canlc v0.0.0. DO NOT EDIT.
// Prod emit: tests + given stripped.
export type RatioResult = { $can_kind: "ok"; denominator: bigint; numerator: bigint } | { $can_kind: "ok"; value: bigint } | { $can_kind: "ok"; value: boolean } | { $can_kind: "ok"; value: string } | { $can_kind: "ok"; values: bigint[] } | { $can_kind: "math.zero_divisor"; divisor: bigint } | { $can_kind: "math.nonterminating_decimal"; numerator: bigint; denominator: bigint };
export type Ratio__Fraction = { numerator: bigint; denominator: bigint };
export type Ratio__Int = { value: bigint };
export type Ratio__Bool = { value: boolean };
export type Ratio__Dec = { value: string };
export type Ratio__Digits = { values: bigint[] };
// Exact-decimal runtime (a10): canonical-digit strings, BigInt math.
function $canDecSplit(d: string): { neg: boolean; ip: string; fp: string } {
  let neg = false;
  if (d.startsWith("-")) {
    neg = true;
    d = d.slice(1);
  }
  const dot = d.indexOf(".");
  return { neg, ip: d.slice(0, dot), fp: d.slice(dot + 1) };
}
function $canDecNorm(ip: string, fp: string, neg: boolean): string {
  ip = ip.replace(/^0+(?=\d)/, "");
  fp = fp.replace(/0+$/, "");
  if (fp === "") {
    fp = "0";
  }
  if (ip === "0" && fp === "0") {
    return "0.0";
  }
  return (neg ? "-" : "") + ip + "." + fp;
}
function $canDecMant(p: { neg: boolean; ip: string; fp: string }, scale: number): bigint {
  let f = p.fp;
  while (f.length < scale) {
    f += "0";
  }
  const m = BigInt(p.ip + f);
  return p.neg ? -m : m;
}
function $canDecFromMant(m: bigint, scale: number): string {
  let neg = false;
  if (m < 0n) {
    neg = true;
    m = -m;
  }
  let digits = m.toString();
  while (digits.length < scale + 1) {
    digits = "0" + digits;
  }
  return $canDecNorm(digits.slice(0, digits.length - scale), digits.slice(digits.length - scale), neg);
}
function $canDecAdd(a: string, b: string): string {
  const A = $canDecSplit(a);
  const B = $canDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $canDecFromMant($canDecMant(A, s) + $canDecMant(B, s), s);
}
function $canDecSub(a: string, b: string): string {
  const A = $canDecSplit(a);
  const B = $canDecSplit(b);
  const s = Math.max(A.fp.length, B.fp.length);
  return $canDecFromMant($canDecMant(A, s) - $canDecMant(B, s), s);
}
function $canDecMul(a: string, b: string): string {
  const A = $canDecSplit(a);
  const B = $canDecSplit(b);
  return $canDecFromMant($canDecMant(A, A.fp.length) * $canDecMant(B, B.fp.length), A.fp.length + B.fp.length);
}
// Sequence indexing (a38 S3): bounds throw, matching Go.
function $canSeqAt<T>(a: T[], i: bigint): T {
  if (i < 0n || i > BigInt(Number.MAX_SAFE_INTEGER)) throw new Error("seq index out of range");
  const k = Number(i);
  if (k >= a.length) throw new Error("seq index out of range");
  return a[k];
}
// Euclidean integer division (a17): quotient and remainder with
// 0 <= r < |b| on every sign combination.
function $canDivMod(a: bigint, b: bigint): [bigint, bigint] {
  let q: bigint = a / b;
  let r: bigint = a - b * q;
  if (r < 0n) {
    q += b > 0n ? -1n : 1n;
    r = a - b * q;
  }
  return [q, r];
}
export function std__ratio__abs(value: bigint): { $can_kind: "ok"; value: bigint } {
  if ((value < 0n)) {
    return { $can_kind: "ok", value: (0n - value) };
  }
  else {
    return { $can_kind: "ok", value: value };
  }
}
export function std__ratio__gcd_from(a: bigint, b: bigint): { $can_kind: "ok"; value: bigint } {
  if ((b <= 0n)) {
    return { $can_kind: "ok", value: a };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__ratio__gcd_from(b, $canDivMod(a, b)[1]);
    switch ($can_m1.$can_kind) {
    case "ok": {
      const r = $can_m1;
      return { $can_kind: "ok", value: r.value };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__ratio__make(numerator: bigint, denominator: bigint): { $can_kind: "ok"; denominator: bigint; numerator: bigint } | { $can_kind: "math.zero_divisor"; divisor: bigint } {
  if ((denominator === 0n)) {
    return { $can_kind: "math.zero_divisor", divisor: denominator };
  }
  else {
    if ((numerator === 0n)) {
      return { $can_kind: "ok", numerator: 0n, denominator: 1n };
    }
    else {
      const $can_m1: { $can_kind: "ok"; value: bigint } = std__ratio__abs(numerator);
      switch ($can_m1.$can_kind) {
      case "ok": {
        const an = $can_m1;
        const $can_m2: { $can_kind: "ok"; value: bigint } = std__ratio__abs(denominator);
        switch ($can_m2.$can_kind) {
        case "ok": {
          const ad = $can_m2;
          const $can_m3: { $can_kind: "ok"; value: bigint } = std__ratio__gcd_from(an.value, ad.value);
          switch ($can_m3.$can_kind) {
          case "ok": {
            const g = $can_m3;
            const $can_m4 = (numerator < 0n);
            const $can_m5 = (denominator < 0n);
            if ($can_m4 && $can_m5) {
              return { $can_kind: "ok", numerator: $canDivMod(an.value, g.value)[0], denominator: $canDivMod(ad.value, g.value)[0] };
            }
            else if ($can_m4 && !($can_m5)) {
              return { $can_kind: "ok", numerator: $canDivMod((0n - an.value), g.value)[0], denominator: $canDivMod(ad.value, g.value)[0] };
            }
            else if (!($can_m4) && $can_m5) {
              return { $can_kind: "ok", numerator: $canDivMod((0n - an.value), g.value)[0], denominator: $canDivMod(ad.value, g.value)[0] };
            }
            else {
              return { $can_kind: "ok", numerator: $canDivMod(an.value, g.value)[0], denominator: $canDivMod(ad.value, g.value)[0] };
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
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
}
export function std__ratio__add(left: Ratio__Fraction, right: Ratio__Fraction): { $can_kind: "ok"; denominator: bigint; numerator: bigint } | { $can_kind: "math.zero_divisor"; divisor: bigint } {
  const $can_m1: { $can_kind: "ok"; denominator: bigint; numerator: bigint } | { $can_kind: "math.zero_divisor"; divisor: bigint } = std__ratio__make(((left.numerator * right.denominator) + (right.numerator * left.denominator)), (left.denominator * right.denominator));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", numerator: r.numerator, denominator: r.denominator };
  }
  case "math.zero_divisor": {
    const e = $can_m1;
    return { $can_kind: "math.zero_divisor", divisor: e.divisor };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__ratio__multiply(left: Ratio__Fraction, right: Ratio__Fraction): { $can_kind: "ok"; denominator: bigint; numerator: bigint } | { $can_kind: "math.zero_divisor"; divisor: bigint } {
  const $can_m1: { $can_kind: "ok"; denominator: bigint; numerator: bigint } | { $can_kind: "math.zero_divisor"; divisor: bigint } = std__ratio__make((left.numerator * right.numerator), (left.denominator * right.denominator));
  switch ($can_m1.$can_kind) {
  case "ok": {
    const r = $can_m1;
    return { $can_kind: "ok", numerator: r.numerator, denominator: r.denominator };
  }
  case "math.zero_divisor": {
    const e = $can_m1;
    return { $can_kind: "math.zero_divisor", divisor: e.divisor };
  }
  default: {
    throw new Error("unreachable");
  }
  }
}
export function std__ratio__divide(left: Ratio__Fraction, right: Ratio__Fraction): { $can_kind: "ok"; denominator: bigint; numerator: bigint } | { $can_kind: "math.zero_divisor"; divisor: bigint } {
  if ((right.numerator === 0n)) {
    return { $can_kind: "math.zero_divisor", divisor: right.numerator };
  }
  else {
    const $can_m1: { $can_kind: "ok"; denominator: bigint; numerator: bigint } | { $can_kind: "math.zero_divisor"; divisor: bigint } = std__ratio__make((left.numerator * right.denominator), (left.denominator * right.numerator));
    switch ($can_m1.$can_kind) {
    case "ok": {
      const r = $can_m1;
      return { $can_kind: "ok", numerator: r.numerator, denominator: r.denominator };
    }
    case "math.zero_divisor": {
      const e = $can_m1;
      return { $can_kind: "math.zero_divisor", divisor: e.divisor };
    }
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
export function std__ratio__seen_from(seen: bigint[], target: bigint, pos: bigint, fuel: bigint): { $can_kind: "ok"; value: boolean } {
  if ((fuel <= 0n)) {
    return { $can_kind: "ok", value: false };
  }
  else {
    if ((pos < (BigInt([...seen].length)))) {
      if (($canSeqAt(seen, pos) === target)) {
        return { $can_kind: "ok", value: true };
      }
      else {
        const $can_m1: { $can_kind: "ok"; value: boolean } = std__ratio__seen_from(seen, target, (pos + 1n), (fuel - 1n));
        switch ($can_m1.$can_kind) {
        case "ok": {
          const r = $can_m1;
          return { $can_kind: "ok", value: r.value };
        }
        default: {
          throw new Error("unreachable");
        }
        }
      }
    }
    else {
      return { $can_kind: "ok", value: false };
    }
  }
}
export function std__ratio__digits_from(orig_n: bigint, orig_d: bigint, rem: bigint, denom: bigint, digits: bigint[], seen: bigint[], fuel: bigint): { $can_kind: "ok"; values: bigint[] } | { $can_kind: "math.nonterminating_decimal"; numerator: bigint; denominator: bigint } {
  if ((fuel <= 0n)) {
    return { $can_kind: "math.nonterminating_decimal", numerator: orig_n, denominator: orig_d };
  }
  else {
    if ((rem === 0n)) {
      return { $can_kind: "ok", values: digits };
    }
    else {
      const $can_m1: { $can_kind: "ok"; value: boolean } = std__ratio__seen_from(seen, rem, 0n, ((BigInt([...seen].length)) + 1n));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const hit = $can_m1;
        if (hit.value) {
          return { $can_kind: "math.nonterminating_decimal", numerator: orig_n, denominator: orig_d };
        }
        else {
          const $can_m2: { $can_kind: "ok"; values: bigint[] } | { $can_kind: "math.nonterminating_decimal"; numerator: bigint; denominator: bigint } = std__ratio__digits_from(orig_n, orig_d, ((rem * 10n) - ($canDivMod((rem * 10n), denom)[0] * denom)), denom, [...digits, $canDivMod((rem * 10n), denom)[0]], [...seen, rem], (fuel - 1n));
          switch ($can_m2.$can_kind) {
          case "ok": {
            const r = $can_m2;
            return { $can_kind: "ok", values: r.values };
          }
          case "math.nonterminating_decimal": {
            const e = $can_m2;
            return { $can_kind: "math.nonterminating_decimal", numerator: e.numerator, denominator: e.denominator };
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
  }
}
export function std__ratio__digit_from(remaining: bigint, acc: string, fuel: bigint): { $can_kind: "ok"; value: string } {
  if ((fuel <= 0n)) {
    return { $can_kind: "ok", value: acc };
  }
  else {
    if ((remaining <= 0n)) {
      return { $can_kind: "ok", value: acc };
    }
    else {
      const $can_m1: { $can_kind: "ok"; value: string } = std__ratio__digit_from((remaining - 1n), $canDecAdd(acc, "1.0"), (fuel - 1n));
      switch ($can_m1.$can_kind) {
      case "ok": {
        const r = $can_m1;
        return { $can_kind: "ok", value: r.value };
      }
      default: {
        throw new Error("unreachable");
      }
      }
    }
  }
}
export function std__ratio__fold_from(digits: bigint[], pos: bigint, place: string, acc: string, fuel: bigint): { $can_kind: "ok"; value: string } {
  if ((fuel <= 0n)) {
    return { $can_kind: "ok", value: acc };
  }
  else {
    if ((pos < (BigInt([...digits].length)))) {
      const $can_m1: { $can_kind: "ok"; value: string } = std__ratio__digit_from($canSeqAt(digits, pos), "0.0", 10n);
      switch ($can_m1.$can_kind) {
      case "ok": {
        const d = $can_m1;
        const $can_m2: { $can_kind: "ok"; value: string } = std__ratio__fold_from(digits, (pos + 1n), $canDecMul(place, "0.1"), $canDecAdd(acc, $canDecMul(d.value, place)), (fuel - 1n));
        switch ($can_m2.$can_kind) {
        case "ok": {
          const r = $can_m2;
          return { $can_kind: "ok", value: r.value };
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
      return { $can_kind: "ok", value: acc };
    }
  }
}
export function std__ratio__to_dec_exact(value: Ratio__Fraction): { $can_kind: "ok"; value: string } | { $can_kind: "math.zero_divisor"; divisor: bigint } | { $can_kind: "math.nonterminating_decimal"; numerator: bigint; denominator: bigint } {
  if ((value.denominator === 0n)) {
    return { $can_kind: "math.zero_divisor", divisor: value.denominator };
  }
  else {
    const $can_m1: { $can_kind: "ok"; value: bigint } = std__ratio__abs(value.numerator);
    switch ($can_m1.$can_kind) {
    case "ok": {
      const an = $can_m1;
      const $can_m2: { $can_kind: "ok"; value: bigint } = std__ratio__abs(value.denominator);
      switch ($can_m2.$can_kind) {
      case "ok": {
        const ad = $can_m2;
        const $can_m3: { $can_kind: "ok"; value: string } = std__ratio__digit_from($canDivMod(an.value, ad.value)[0], "0.0", ($canDivMod(an.value, ad.value)[0] + 1n));
        switch ($can_m3.$can_kind) {
        case "ok": {
          const ip = $can_m3;
          const $can_m4: { $can_kind: "ok"; values: bigint[] } | { $can_kind: "math.nonterminating_decimal"; numerator: bigint; denominator: bigint } = std__ratio__digits_from(value.numerator, value.denominator, (an.value - ($canDivMod(an.value, ad.value)[0] * ad.value)), ad.value, [], [], (ad.value + 1n));
          switch ($can_m4.$can_kind) {
          case "ok": {
            const dg = $can_m4;
            const $can_m5: { $can_kind: "ok"; value: string } = std__ratio__fold_from(dg.values, 0n, "0.1", "0.0", ((BigInt([...dg.values].length)) + 1n));
            switch ($can_m5.$can_kind) {
            case "ok": {
              const fr = $can_m5;
              const $can_m6 = (value.numerator < 0n);
              const $can_m7 = (value.denominator < 0n);
              if ($can_m6 && $can_m7) {
                return { $can_kind: "ok", value: $canDecAdd(ip.value, fr.value) };
              }
              else if ($can_m6 && !($can_m7)) {
                return { $can_kind: "ok", value: $canDecSub("0.0", $canDecAdd(ip.value, fr.value)) };
              }
              else if (!($can_m6) && $can_m7) {
                return { $can_kind: "ok", value: $canDecSub("0.0", $canDecAdd(ip.value, fr.value)) };
              }
              else {
                return { $can_kind: "ok", value: $canDecAdd(ip.value, fr.value) };
              }
            }
            default: {
              throw new Error("unreachable");
            }
            }
          }
          case "math.nonterminating_decimal": {
            const e = $can_m4;
            return { $can_kind: "math.nonterminating_decimal", numerator: e.numerator, denominator: e.denominator };
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
    default: {
      throw new Error("unreachable");
    }
    }
  }
}
