// Audit observation, not an assertion that trusted direct interop is a sandbox.
// Run from the repo root: node --experimental-strip-types docs/audit-probes/host-boundary.mjs
import {
  success__id$T$Seq$L$int$G$ as sequenceIdentity,
  success__id$T$Success__Pair as recordIdentity,
  success__map$T$int$T$int as mapOption,
} from '../../sketches/success-values/success.ts';

const input = [1n];
const output = sequenceIdentity(input);
input[0] = 9n;
console.log('ABI_ALIAS', output.value === input, String(output.value[0]));

let reads = 0;
const record = {
  get left() { reads++; return 7n; },
  get right() { reads++; return true; },
};
console.log('GETTER_INPUT', recordIdentity(record), reads);

let callbackEffects = 0;
mapOption({ $can_kind: 'Success__Some$T$int', value: 1n }, (value) => {
  callbackEffects++;
  return { $can_kind: 'ok', value };
});
console.log('HOST_CALLBACK_EFFECTS', callbackEffects);
