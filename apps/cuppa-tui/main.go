package main

import "fmt"

// version is stamped at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	fmt.Println("cuppa", version)
}
