// Keyboard shortcuts a terminal cannot tell apart, and Cmd on a Mac.
//
// A terminal sends Ctrl+Shift+S as the same byte as Ctrl+S, Ctrl+[ as Esc, and
// never sends Cmd. Bubble Tea also reads the enhanced form of a key ("CSI u":
// ESC [ code ; modifiers u), which carries every modifier. So the page catches
// the shortcuts the app uses and sends that form itself: Ctrl+Shift+S arrives as
// Ctrl+Shift+S, and on a Mac Cmd+S arrives as Ctrl+S, which is what the app
// listens for (the app shows Cmd in its menus there).
//
// This file is the same in apps/cuppa-web and apps/cuppa-desktop.

/** True on a Mac, where the command key does what Ctrl does elsewhere. */
export function isMac (): boolean {
  const platform = (navigator as Navigator & { userAgentData?: { platform?: string } }).userAgentData?.platform ?? navigator.platform ?? ''
  return /Mac/i.test(platform) || /Macintosh|Mac OS X/.test(navigator.userAgent)
}

// The letters the app binds with Ctrl: new, open, quit, preview, save, undo,
// redo, duplicate, group, ungroup, copy, paste, find.
const letters = new Set('nopqszydgucvf')
// Keys known by position, because Shift changes what they type ([ becomes {).
const byPosition: Record<string, number> = { BracketLeft: 91, BracketRight: 93 }

/** The bytes for the shortcut in a key press, or undefined when the page should leave it to the terminal. */
export function shortcutBytes (e: Pick<KeyboardEvent, 'key' | 'code' | 'ctrlKey' | 'metaKey' | 'altKey' | 'shiftKey'>, mac: boolean): string | undefined {
  const primary = mac ? e.metaKey : e.ctrlKey
  if (!primary || e.altKey || (mac && e.ctrlKey)) return undefined
  let code = byPosition[e.code]
  if (code === undefined) {
    const key = e.key.toLowerCase()
    if (key.length === 1 && letters.has(key)) code = key.charCodeAt(0)
  }
  if (code === undefined) return undefined
  const modifiers = 1 + (e.shiftKey ? 1 : 0) + 4 // 1, plus Shift 1, Alt 2, Ctrl 4
  return `\x1b[${code};${modifiers}u`
}

/** Catches shortcuts before the terminal does and passes them on as input. Returns a function that stops it. */
export function bindShortcuts (send: (bytes: string) => void, mac: boolean = isMac()): () => void {
  const onKeyDown = (e: KeyboardEvent): void => {
    const bytes = shortcutBytes(e, mac)
    if (bytes === undefined) return
    e.preventDefault()
    e.stopPropagation()
    send(bytes)
  }
  window.addEventListener('keydown', onKeyDown, true)
  return () => window.removeEventListener('keydown', onKeyDown, true)
}
