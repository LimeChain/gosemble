//go:build nonwasmenv

package hashing

import (
	"github.com/ChainSafe/gossamer/lib/common"
)

func Twox128(value []byte) []byte {
	h, _ := common.Twox128Hash(value)
	return h[:]
}

func Twox64(value []byte) []byte {
	h, _ := common.Twox64(value)
	return h[:]
}

func Blake128(value []byte) []byte {
	h, _ := common.Blake2b128(value)
	return h[:]
}

func Blake256(value []byte) []byte {
	h, _ := common.Blake2bHash(value)
	return h[:]
}

// Blake2b8 returns the first 8 bytes of the Blake2b hash of the input data
func Blake2b8(data []byte) (digest [8]byte, err error) {
	return common.Blake2b8(data)
}

// MustBlake2b8 returns the first 8 bytes of the Blake2b hash of the input data
func MustBlake2b8(data []byte) []byte {
	digest := common.MustBlake2b8(data)
	return digest[:]
}
