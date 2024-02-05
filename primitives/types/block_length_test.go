package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedBlockLengthBytes, _ = hex.DecodeString(
		"000000000100000002000000",
	)
)

var (
	targetBlockLength = BlockLength{
		Max: PerDispatchClassU32{
			Normal:      sc.U32(0),
			Operational: sc.U32(1),
			Mandatory:   sc.U32(2),
		},
	}
)

func Test_BlockLength_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetBlockLength.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectedBlockLengthBytes, buffer.Bytes())
}

func Test_BlockLength_Bytes(t *testing.T) {
	assert.Equal(t, expectedBlockLengthBytes, targetBlockLength.Bytes())
}
