import { useEffect, useState } from 'react'
import { TTY, createWailsSocket } from '@treactui/tty'
import type { SocketLike } from '@treactui/tty'
import { startCuppa, wasmRuntime } from './wasm_runtime'
import { bindShortcuts } from './shortcuts'

// Created once, outside the component: a new function on every render would
// reconnect the terminal each time. The page keeps the connection it made so
// that shortcuts the terminal cannot express can be sent as input.
const open = createWailsSocket({ runtime: wasmRuntime })
let live: SocketLike | undefined
const socket = (url?: string): SocketLike => (live = open(url))

/** The whole page: one terminal showing the Cuppa terminal app, which runs in this tab as WebAssembly. */
export default function App () {
  const [state, setState] = useState<'loading' | 'ready' | string>('loading')

  useEffect(() => {
    startCuppa('./cuppa.wasm').then(() => setState('ready'), (error: unknown) => setState(String(error)))
  }, [])

  // Ctrl+Shift+S, Ctrl+[ and, on a Mac, Cmd+key reach the app as enhanced keys.
  useEffect(() => bindShortcuts(data => live?.send(JSON.stringify({ type: 'input', data }))), [])

  if (state === 'loading') return <main className='window'><p className='startup-error'>Loading Cuppa…</p></main>
  if (state !== 'ready') return <main className='window'><p className='startup-error'>{state}</p></main>
  return <main className='window'><TTY createSocket={socket} label='Cuppa' className='terminal' /></main>
}
