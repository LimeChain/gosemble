package main

import (
	_ "unsafe"

	"github.com/LimeChain/gosemble/primitives/log"
)

//go:linkname _panic runtime._panic
func _panic(msg string)

func Panic(str string) {
	log.Log(log.Critical, []byte("TODO:"), []byte(str))
	_panic(str)
}
