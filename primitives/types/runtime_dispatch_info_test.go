package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectBytesRuntimeDispatchInfo, _ = hex.DecodeString("04080003000000000000000000000000000000")
)

var (
	runtimeDispatchInfo = RuntimeDispatchInfo{
		Weight:     WeightFromParts(1, 2),
		Class:      NewDispatchClassNormal(),
		PartialFee: sc.NewU128(3),
	}
)

func Test_RuntimeDispatchInfo_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	runtimeDispatchInfo.Encode(buffer)

	assert.Equal(t, expectBytesRuntimeDispatchInfo, buffer.Bytes())
}

func Test_DecodeRuntimeDispatchInfo(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesRuntimeDispatchInfo)

	result, err := DecodeRuntimeDispatchInfo(buffer)
	assert.Nil(t, err)

	assert.Equal(t, runtimeDispatchInfo, result)
}

func Test_RuntimeDispatchInfo_Bytes(t *testing.T) {
	result := runtimeDispatchInfo.Bytes()

	assert.Equal(t, expectBytesRuntimeDispatchInfo, result)
}
