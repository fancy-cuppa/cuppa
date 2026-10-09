import { useEffect } from 'react'
import { TTY, createWailsSocket } from '@treactui/tty'
import type { SocketLike } from '@treactui/tty'
import { bindShortcuts } from './shortcuts'

// Created once, outside the component: a new function on every render would
// reconnect the terminal each time. The page keeps the connection it made so
// that shortcuts the terminal cannot express can be sent as input.
const open = createWailsSocket()
let live: SocketLike | undefined
const wailsSocket = (url?: string): SocketLike => (live = open(url))

/** The whole window: one terminal showing the Cuppa terminal app, over Wails events. */
export default function App () {
  // Ctrl+Shift+S, Ctrl+[ and, on a Mac, Cmd+key reach the app as enhanced keys.
  useEffect(() => bindShortcuts(data => live?.send(JSON.stringify({ type: 'input', data }))), [])
  return <main className='window'><TTY createSocket={wailsSocket} label='Cuppa' className='terminal' /></main>
}
