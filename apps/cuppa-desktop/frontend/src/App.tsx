import { useEffect, useState } from 'react'
import { TTY } from '@treactui/tty'
import { TerminalURL } from '../wailsjs/go/main/App'

/** The whole window: one terminal showing the Cuppa terminal app. */
export default function App () {
  const [url, setUrl] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    TerminalURL().then(setUrl).catch((e: unknown) => setError(String(e)))
  }, [])

  if (error) {
    return <p role='alert' className='startup-error'>Cuppa could not start its terminal: {error}</p>
  }
  return <main className='window'>{url && <TTY url={url} label='Cuppa' className='terminal' />}</main>
}
