package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	expectedDispatchInfoBytes, _ = hex.DecodeString("04080201")

	targetDispatchInfo = DispatchInfo{
		Weight:  WeightFromParts(1, 2),
		Class:   NewDispatchClassMandatory(),
		PaysFee: NewPaysNo(),
	}
)

func Test_DispatchInfo_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	err := targetDispatchInfo.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedDispatchInfoBytes, buf.Bytes())
}

func Test_DispatchInfo_Bytes(t *testing.T) {
	assert.Equal(t, expectedDispatchInfoBytes, targetDispatchInfo.Bytes())
}

func Test_DecodeDispatchInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.Write(expectedDispatchInfoBytes)

	result, err := DecodeDispatchInfo(buf)
	assert.NoError(t, err)

	assert.Equal(t, targetDispatchInfo, result)
}

func Test_GetDispatchInfo(t *testing.T) {
	call := testCall{}

	result := GetDispatchInfo(call)

	assert.Equal(t, DispatchInfo{
		Weight:  WeightFromParts(3, 4),
		Class:   NewDispatchClassNormal(),
		PaysFee: NewPaysYes(),
	}, result)
}
