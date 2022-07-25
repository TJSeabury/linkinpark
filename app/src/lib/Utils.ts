export function wait(ms: number) {
  const start = Date.now();
  let now = start;
  while (now - start < ms) {
    now = Date.now();
  }
}

export const diff = (() => {
  const start = Date.now();
  return () => Date.now() - start;
})();

export const limit = (ms: number) => diff() < ms;