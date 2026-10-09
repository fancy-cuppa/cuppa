import { useEffect, useState } from 'react'
import { TTY, createWailsSocket } from '@treactui/tty'
import { startCuppa, wasmRuntime } from './wasm_runtime'

// Created once, outside the component: a new function on every render would
// reconnect the terminal each time.
const socket = createWailsSocket({ runtime: wasmRuntime })

/** The whole page: one terminal showing the Cuppa terminal app, which runs in this tab as WebAssembly. */
export default function App () {
  const [state, setState] = useState<'loading' | 'ready' | string>('loading')

  useEffect(() => {
    startCuppa('./cuppa.wasm').then(() => setState('ready'), (error: unknown) => setState(String(error)))
  }, [])

  if (state === 'loading') return <main className='window'><p className='startup-error'>Loading Cuppa…</p></main>
  if (state !== 'ready') return <main className='window'><p className='startup-error'>{state}</p></main>
  return <main className='window'><TTY createSocket={socket} label='Cuppa' className='terminal' /></main>
}
