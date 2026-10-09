# cuppa-web (spike, #98)

The Cuppa terminal app compiled to WebAssembly and shown in a web page: no
server, only static files. The page is the same TReactUI terminal the desktop
app uses (`@treactui/tty`); only the transport differs.

```
node build-wasm.cjs          # dist/cuppa.wasm, copied with wasm_exec.js to frontend/public
node smoke-wasm.cjs          # runs the wasm in Node and checks the first screen
cd frontend && npm install && npm run build && npm run preview
```

## How it fits together

- `main_js.go` runs `shell.New(...)` behind tty-go's `BindShared`, exactly like
  `cuppa-desktop` does behind Wails.
- `events_js.go` is tty-go's two-method `Events` bus over `syscall/js`:
  the page calls `cuppaBridge.up(name, data)`, Go calls `cuppaBridge.down(name, data)`.
- `frontend/src/wasm_runtime.ts` gives `createWailsSocket` the `EventsOn` /
  `EventsEmit` it expects (its `runtime` option), so the page needs no new socket.

## Two patches, applied at build time

Bubble Tea v2.1.0 and `atotto/clipboard` do not build for `js/wasm`.
`build-wasm.cjs` copies them next to the build, adds the files in `wasmpatch/`
and builds against the copies through a temporary `go.work` `replace`.
Nothing is forked or committed. (`go build -overlay` cannot be used: Go refuses
overlays inside the module cache.)

1. `tty_js.go`: no terminal to put in raw mode, no process to suspend, no resize signal.
2. `mapNl := false`: Bubble Tea assumes a line feed returns the cursor to column
   0 everywhere except Windows. tty-go rewrites each line feed as `ESC D`
   (same column), which is only right under the Windows assumption. Without this
   the screen is scrambled. **tty-go has the same mismatch on Linux and macOS
   (see the finding below).**

## Result

Works: the editor draws, the mouse drags a component from the palette onto the
canvas, menus open, the details panel edits. Checked in headless Chrome.

Not done, in rough order of effort:

- **Size:** `cuppa.wasm` is 12.5 MB (3.3 MB gzipped). Fine for a demo, slow to start on a phone.
- **Files:** Open and Save need browser pickers and downloads; there is no filesystem.
- **Clipboard:** `clipboard_js.go` returns an error. Needs the async browser API behind a user gesture.
- **Threads:** Go wasm is single-threaded and runs on the page's main thread. A Web Worker would keep the page responsive.
- **Settings and packs:** layout and user packs are read from disk; they are skipped here.
- **Accessibility:** the desktop app wraps the model with `terminal.Describe` (screen-reader snapshots). It lives in `cuppa-desktop`, so it is not used here yet.

## Finding for tty-go (not verified on Linux or macOS)

`session.rawLineFeeds` rewrites `\n` as `ESC D` for every OS, but Bubble Tea
v2.1.0 only emits "line feed, same column" on Windows
(`mapNl := runtime.GOOS != "windows" && p.ttyInput == nil`). On Linux and macOS
it assumes the line feed also returns to column 0, so the rewrite should apply
only when `runtime.GOOS == "windows"`. The browser build hit exactly this
(`GOOS=js`); `cuppa-desktop` on Linux and macOS very likely does too. Read from
the code only: the desktop app has not been run on those systems.
