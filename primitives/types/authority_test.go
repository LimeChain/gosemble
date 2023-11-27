package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedAuthorityBytes, _ = hex.DecodeString(
		"01010101010101010101010101010101010101010101010101010101010101010300000000000000",
	)
)
var (
	targetAuthority = Authority{
		Id: NewAccountIdFromAddress32(
			Address32{
				FixedSequence: sc.
					NewFixedSequence[sc.U8](32,
					1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
					1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
					1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
					1, 1),
			},
		),
		Weight: 3,
	}
)

func Test_Authority_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetAuthority.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectedAuthorityBytes, buffer.Bytes())
}

func Test_Authority_Bytes(t *testing.T) {
	assert.Equal(t, expectedAuthorityBytes, targetAuthority.Bytes())
}

func Test_DecodeAuthority(t *testing.T) {
	buffer := bytes.NewBuffer(expectedAuthorityBytes)

	result, err := DecodeAuthority(buffer)
	assert.NoError(t, err)
	assert.Equal(t, targetAuthority, result)
}
