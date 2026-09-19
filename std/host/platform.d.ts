// Platform ambient surface for the host shelf (tsc support only;
// node resolves the real builtins at runtime). The tsc gate runs
// without Node libs, so exactly the bindings the host
// implementations use are declared here — nothing more.
declare module "node:crypto" {
  export function randomBytes(n: number): Uint8Array;
  export function createHash(algorithm: string): {
    update(data: Uint8Array): { digest(): Uint8Array };
  };
  export function timingSafeEqual(a: Uint8Array, b: Uint8Array): boolean;
}
