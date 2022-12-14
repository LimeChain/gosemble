package hashing

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

func Twox128(value []byte) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(value)
	r := env.ExtHashingTwox128Version1(keyOffsetSize)

	return utils.ToWasmMemorySlice(r, 16)
}

func Twox64(value []byte) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(value)

	r := env.ExtHashingTwox64Version1(keyOffsetSize)

	return utils.ToWasmMemorySlice(r, 8)
}
