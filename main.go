package main

import "github.com/radkomih/gosemble/dev"

func main() {
	dev.RunInWazmer("build/dev_runtime.wasm")
}
