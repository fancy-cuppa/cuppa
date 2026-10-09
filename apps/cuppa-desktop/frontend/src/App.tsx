import { TTY, createWailsSocket } from '@treactui/tty'

// Created once, outside the component: a new function on every render would
// reconnect the terminal each time.
const wailsSocket = createWailsSocket()

/** The whole window: one terminal showing the Cuppa terminal app, over Wails events. */
export default function App () {
  return <main className='window'><TTY createSocket={wailsSocket} label='Cuppa' className='terminal' /></main>
}
