//go:build !(js && wasm)

// Command cuppa-web only exists for WebAssembly: build it with
// `node build-wasm.cjs` (GOOS=js GOARCH=wasm).
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "cuppa-web runs in a browser: build it with `node build-wasm.cjs`")
	os.Exit(2)
}
