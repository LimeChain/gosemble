package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SS58Decode(t *testing.T) {
	addr := "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY"
	wantPubKey := []byte{212, 53, 147, 199, 21, 253, 211, 28, 97, 20, 26, 189, 4, 169, 159, 214, 130, 44, 133, 88, 133, 76, 205, 227, 154, 86, 132, 231, 165, 109, 162, 125}
	wantPrefix := uint16(42)

	prefix, pubKey, err := SS58Decode(addr)
	assert.NoError(t, err)
	assert.Equal(t, wantPubKey, pubKey)
	assert.Equal(t, wantPrefix, prefix)
}
