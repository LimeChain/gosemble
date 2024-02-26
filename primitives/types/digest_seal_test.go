package types

import (
	"bytes"
	"encoding/hex"
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	expectBytesDigestSeal, _ = hex.DecodeString("74657374200100000000000000")
)

var (
	targetDigestSeal = NewDigestSeal(
		consensusEngineId,
		sc.BytesToSequenceU8(sc.U64(1).Bytes()),
	)
)

func Test_DigestSeal_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetDigestSeal.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectBytesDigestSeal, buffer.Bytes())
}

func Test_DigestSeal_Bytes(t *testing.T) {
	assert.Equal(t, expectBytesDigestSeal, targetDigestSeal.Bytes())
}
