export async function register() {
  console.log({ env: process.env.NEXT_RUNTIME })
  await import('./instrumentation.node')
}
