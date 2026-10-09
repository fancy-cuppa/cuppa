#!/usr/bin/env node
// Starts dist/cuppa.wasm in Node, connects over the same bridge the page uses,
// and checks that the first screen has the editor on it. It proves the Go side
// runs as WebAssembly; the page itself (xterm.js) is not involved.
'use strict'
const { readFileSync } = require('node:fs')
const { join } = require('node:path')
const { spawnSync } = require('node:child_process')

const dist = join(__dirname, 'dist')
const goroot = spawnSync('go', ['env', 'GOROOT'], { encoding: 'utf8' }).stdout.trim()
require(join(goroot, 'lib', 'wasm', 'wasm_exec.js'))

const frames = []
globalThis.cuppaBridge = { down: (_name, data) => frames.push(JSON.parse(JSON.parse(data).f ?? 'null')) }
const sleep = ms => new Promise(resolve => setTimeout(resolve, ms))
const fail = message => { console.error('smoke-wasm:', message); process.exit(1) }

const go = new Go()
WebAssembly.instantiate(readFileSync(join(dist, 'cuppa.wasm')), go.importObject).then(async ({ instance }) => {
  void go.run(instance)
  for (let waited = 0; !globalThis.cuppaBridge.up; waited += 10) {
    if (waited > 10000) fail('the program did not start')
    await sleep(10)
  }
  const up = (n, t, f) => globalThis.cuppaBridge.up('treactui:up', JSON.stringify({ c: 'smoke', n, t, f }))
  up(0, 'open')
  up(1, 'frame', JSON.stringify({ type: 'resize', cols: 120, rows: 40 }))
  await sleep(2500)

  const output = frames.filter(f => f?.type === 'output').map(f => f.data).join('')
  for (const want of ['COMPONENTS', 'DETAILS', 'Untitled']) {
    if (!output.includes(want)) fail(`the first screen lacks "${want}"`)
  }
  console.log(`ok: ${output.length} bytes of terminal output, the editor is on screen`)
  process.exit(0)
})
