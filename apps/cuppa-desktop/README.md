# cuppa-desktop

Cuppa in a desktop window: a [Wails](https://wails.io) v2 shell around the very
same terminal app as `cuppa-tui`, shown in a real terminal emulator by
[TReactUI](https://github.com/meta-tui/treactui). See
[ADR 0004](../../docs/adr/0004-desktop-app-wails-and-treactui.md).

```
window (web view)                       Go process
  <TTY url=…>  ──── ws://127.0.0.1:port/term/<secret> ────▶  tty-go  ──▶  shell.New(...)
  React page      ◀── App.TerminalURL() (Wails binding) ──   terminal.Start
```

## Run

You need Go, Node and the Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`;
`wails doctor` checks the rest).

```sh
npx nx run cuppa-desktop:start      # wails dev: live window, also at http://localhost:34115
npx nx run cuppa-desktop:bundle     # wails build: build/bin/cuppa-desktop(.exe)
npx nx run cuppa-desktop:test       # Go tests
```

On Linux install GTK and WebKit first (`libgtk-3-dev`, and `libwebkit2gtk-4.1-dev`
with `-tags webkit2_41` on Ubuntu 24.04). CI does this for you.

## Layout

| | |
|---|---|
| `main.go` | Starts Wails: window options and the embedded page |
| `app_binding.go` | What the page can call, and the window lifecycle (start, close, shutdown) |
| `terminal/` | The loopback WebSocket server (`Start`), and `OnQuit` which closes the window when the app quits |
| `frontend/` | The React page: one `<TTY>`. Built by `wails build`, not by Nx |

## Behaviour worth knowing

- The server only listens on `127.0.0.1`, with a random secret in the address
  and an origin check, because it runs a program.
- Closing the window asks the app to quit first, so unsaved changes prompt.
- Reloading the window (`wails dev`) keeps the design: the program is shared.
- `frontend/dist` only holds a `.gitkeep` in git so the Go code compiles; real
  assets come from `wails build`.
