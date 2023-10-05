package io

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

type Hashing interface {
	Blake128(value []byte) []byte
	Blake256(value []byte) []byte

	Twox128(value []byte) []byte
	Twox64(value []byte) []byte
}

type hashing struct{}

func NewHashing() Hashing {
	return hashing{}
}

func (h hashing) Twox64(value []byte) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(value)
	r := env.ExtHashingTwox64Version1(keyOffsetSize)
	return utils.ToWasmMemorySlice(r, 8)
}

func (h hashing) Twox128(value []byte) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(value)
	r := env.ExtHashingTwox128Version1(keyOffsetSize)
	return utils.ToWasmMemorySlice(r, 16)
}

func (h hashing) Blake128(value []byte) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(value)
	r := env.ExtHashingBlake2128Version1(keyOffsetSize)
	return utils.ToWasmMemorySlice(r, 16)
}

func (h hashing) Blake256(value []byte) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(value)
	r := env.ExtHashingBlake2256Version1(keyOffsetSize)
	return utils.ToWasmMemorySlice(r, 32)
}
