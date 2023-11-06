package hashing

import (
	"golang.org/x/crypto/blake2b"
)

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
func MustBlake2b8(data []byte) [8]byte {
	digest, err := Blake2b8(data)
	if err != nil {
		panic(err)
	}
	return digest
}
