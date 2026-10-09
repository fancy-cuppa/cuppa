// The part of the Wails runtime that createWailsSocket uses (EventsOn,
// EventsEmit), backed by the Go program running in this page as WebAssembly.
// The Go side (events_js.go) exposes globalThis.cuppaBridge.up(name, data) to
// receive events and calls cuppaBridge.down(name, data) to send them.

type Bridge = { up?: (name: string, data: string) => void, down?: (name: string, data: string) => void }
type Callback = (...data: unknown[]) => void

declare global {
  // wasm_exec.js, shipped with Go, defines this.
  class Go {
    importObject: WebAssembly.Imports
    run (instance: WebAssembly.Instance): Promise<void>
  }
  // eslint-disable-next-line no-var
  var cuppaBridge: Bridge | undefined
}

const listeners = new Map<string, Set<Callback>>()

export const wasmRuntime = {
  EventsOn (name: string, callback: Callback): () => void {
    const set = listeners.get(name) ?? new Set<Callback>()
    set.add(callback)
    listeners.set(name, set)
    return () => { set.delete(callback) }
  },
  EventsEmit (name: string, ...data: unknown[]): void {
    globalThis.cuppaBridge?.up?.(name, String(data[0]))
  }
}

/** Downloads and starts the Go program; resolves once it can receive events. */
export async function startCuppa (url: string): Promise<void> {
  const bridge: Bridge = globalThis.cuppaBridge ?? {}
  globalThis.cuppaBridge = bridge
  bridge.down = (name, data) => { listeners.get(name)?.forEach(callback => callback(data)) }

  const go = new Go()
  const { instance } = await WebAssembly.instantiateStreaming(fetch(url), go.importObject)
  void go.run(instance) // never resolves while the program runs
  for (let waited = 0; bridge.up === undefined; waited += 10) {
    if (waited > 10_000) throw new Error('Cuppa did not start')
    await new Promise(resolve => setTimeout(resolve, 10))
  }
}
