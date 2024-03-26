export function timeout(ms: number = 1000): Promise<never> {
  return new Promise((_, reject) => setTimeout(reject, ms));
}