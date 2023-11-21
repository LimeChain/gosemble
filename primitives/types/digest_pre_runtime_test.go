package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectBytesDigestPreRuntime, _ = hex.DecodeString("74657374200100000000000000")
)

var (
	targetDigestPreRuntime = NewDigestPreRuntime(
		consensusEngineId,
		sc.BytesToSequenceU8(sc.U64(1).Bytes()),
	)
)

func Test_DigestPreRuntime_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetDigestPreRuntime.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectBytesDigestPreRuntime, buffer.Bytes())
}

func Test_DigestPreRuntime_Bytes(t *testing.T) {
	assert.Equal(t, expectBytesDigestPreRuntime, targetDigestPreRuntime.Bytes())
}
