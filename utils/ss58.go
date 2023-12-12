package utils

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/blake2b"
)

// SS58Decode decodes an SS58 checksum value into its data and format.
func SS58Decode(address string) (uint16, []byte, error) {
	data := base58.Decode(address)
	if len(data) < 2 {
		return 0, nil, fmt.Errorf("expected at least 2 bytes in base58 decoded address")
	}

	prefixLen := int8(0)
	ident := uint16(0)
	if data[0] <= 63 {
		prefixLen = 1
		ident = uint16(data[0])
	} else if data[0] < 127 {
		lower := (data[0] << 2) | (data[1] >> 6)
		upper := data[1] & 0b00111111
		prefixLen = 2
		ident = uint16(lower) | (uint16(upper) << 8)
	} else {
		return 0, nil, fmt.Errorf("invalid address")
	}

	checkSumLength := 2
	hash := ss58hash(data[:len(data)-checkSumLength])
	checksum := hash[:checkSumLength]

	givenChecksum := data[len(data)-checkSumLength:]
	if !bytes.Equal(givenChecksum, checksum) {
		return 0, nil, fmt.Errorf("checksum mismatch: expected %v but got %v", checksum, givenChecksum)
	}

	return ident, data[prefixLen : len(data)-checkSumLength], nil
}

func ss58hash(data []byte) [64]byte {
	// Adapted from https://github.com/paritytech/substrate/blob/e6def65920d30029e42d498cb07cec5dd433b927/primitives/core/src/crypto.rs#L369
	prefix := []byte("SS58PRE")
	return blake2b.Sum512(append(prefix, data...))
}
