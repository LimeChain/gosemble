//go:build !nonwasmenv

package hashing

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/utils"
	"golang.org/x/crypto/blake2b"
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

func Blake128(value []byte) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(value)
	r := env.ExtHashingBlake2128Version1(keyOffsetSize)
	return utils.ToWasmMemorySlice(r, 16)
}

func Blake256(value []byte) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(value)
	r := env.ExtHashingBlake2256Version1(keyOffsetSize)
	return utils.ToWasmMemorySlice(r, 32)
}

// Blake2b8 returns the first 8 bytes of the Blake2b hash of the input data
func Blake2b8(data []byte) (digest [8]byte, err error) {
	const bytes = 8
	hasher, err := blake2b.New(bytes, nil)
	if err != nil {
		return [8]byte{}, err
	}

	_, err = hasher.Write(data)
	if err != nil {
		return [8]byte{}, err
	}

	digestBytes := hasher.Sum(nil)
	copy(digest[:], digestBytes)
	return digest, nil
}

// MustBlake2b8 returns the first 8 bytes of the Blake2b hash of the input data
func MustBlake2b8(data []byte) []byte {
	digest, err := Blake2b8(data)
	if err != nil {
		log.Critical(err.Error())
	}
	return digest[:]
}
