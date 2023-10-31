package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedDigestItemBytes = []byte{
		0x74, 0x65, 0x73, 0x74, // Engine
		0x20, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // Payload
	}
)

var (
	targetDigestItem = DigestItem{
		Engine:  sc.BytesToFixedSequenceU8([]byte{'t', 'e', 's', 't'}),
		Payload: sc.BytesToSequenceU8(sc.U64(1).Bytes()),
	}
)

func Test_DigestItem_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetDigestItem.Encode(buffer)

	assert.Equal(t, expectedDigestItemBytes, buffer.Bytes())
}

func Test_DigestItem_Bytes(t *testing.T) {
	assert.Equal(t, expectedDigestItemBytes, targetDigestItem.Bytes())
}

func Test_DecodeDigestItem(t *testing.T) {
	buffer := bytes.NewBuffer(expectedDigestItemBytes)

	result, err := DecodeDigestItem(buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetDigestItem, result)
}
