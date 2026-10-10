#!/usr/bin/env node
// Builds cuppa.wasm for the browser. Bubble Tea and atotto/clipboard have no
// js/wasm support, so this copies them, adds the files in wasmpatch/, and builds
// against the copies through a temporary go.work (a `replace` there).
// Nothing is forked or committed, and the module cache is not touched.
'use strict'
const { spawnSync } = require('node:child_process')
const { cpSync, mkdirSync, readFileSync, rmSync, writeFileSync, chmodSync, readdirSync, statSync } = require('node:fs')
const { join, resolve } = require('node:path')

const here = __dirname
const root = resolve(here, '..', '..')
const out = resolve(process.argv[2] || join(here, 'dist'))
const slash = p => p.split('\\').join('/')

function go (args, env) {
  const r = spawnSync('go', args, { cwd: here, encoding: 'utf8', env: { ...process.env, ...env }, stdio: ['ignore', 'pipe', 'inherit'] })
  if (r.status !== 0) process.exit(r.status || 1)
  return r.stdout.trim()
}

const dir = module => {
  go(['mod', 'download', module], { GOOS: '', GOARCH: '' }) // a fresh checkout has not downloaded it yet
  return go(['list', '-m', '-f', '{{.Dir}}', module], { GOOS: '', GOARCH: '' })
}

function writable (path) {
  chmodSync(path, 0o777)
  if (statSync(path).isDirectory()) for (const name of readdirSync(path)) writable(join(path, name))
}

function patched (module, name, patch, file) {
  const target = join(out, 'patched', name)
  rmSync(target, { recursive: true, force: true })
  cpSync(dir(module), target, { recursive: true })
  writable(target)
  cpSync(join(here, 'wasmpatch', patch), join(target, file))
  return target
}

mkdirSync(out, { recursive: true })
const tea = patched('charm.land/bubbletea/v2', 'bubbletea', 'bubbletea_tty_js.go.txt', 'tty_js.go')
const clip = patched('github.com/atotto/clipboard', 'clipboard', 'clipboard_js.go.txt', 'clipboard_js.go')

// The repository's go.work, with absolute paths and the two replacements.
const source = readFileSync(join(root, 'go.work'), 'utf8')
const version = /^go (\S+)/m.exec(source)[1]
const uses = [...source.matchAll(/^[ \t]*\.\/(\S+)[ \t]*$/gm)].map(m => join(root, m[1]))
const work = [
  `go ${version}`,
  '',
  'use (',
  ...uses.map(u => `\t${JSON.stringify(u)}`),
  ')',
  '',
  `replace charm.land/bubbletea/v2 => ${slash(tea)}`,
  `replace github.com/atotto/clipboard => ${slash(clip)}`,
  ''
].join('\n')
const workFile = join(out, 'go.work')
writeFileSync(workFile, work)

go(['build', '-trimpath', '-ldflags', '-s -w', '-o', join(out, 'cuppa.wasm'), '.'], { GOOS: 'js', GOARCH: 'wasm', GOWORK: workFile })
console.log('built', join(out, 'cuppa.wasm'))

// The page serves both from its public folder; wasm_exec.js must match the Go that built the wasm.
const pub = join(here, 'frontend', 'public')
mkdirSync(pub, { recursive: true })
cpSync(join(out, 'cuppa.wasm'), join(pub, 'cuppa.wasm'))
cpSync(join(go(['env', 'GOROOT']), 'lib', 'wasm', 'wasm_exec.js'), join(pub, 'wasm_exec.js'))
console.log('copied cuppa.wasm and wasm_exec.js to', pub)
