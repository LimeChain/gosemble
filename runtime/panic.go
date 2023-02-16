package main

import (
	_ "unsafe"

	"github.com/LimeChain/gosemble/primitives/log"
)

//go:linkname _panic runtime._panic
func _panic(msg string)

func Panic(str string) {
	log.Critical("TODO:", str)
	_panic(str)
}
