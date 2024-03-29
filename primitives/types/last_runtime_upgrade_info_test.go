package types

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectBytesLastRuntimeUpgradeInfo, _ = hex.DecodeString("044c746573742d6c7275692d737065632d6e616d65")
)

var (
	lrui = LastRuntimeUpgradeInfo{
		SpecVersion: sc.Compact{Number: sc.U32(1)},
		SpecName:    "test-lrui-spec-name",
	}

	lruiInvalidNumber = LastRuntimeUpgradeInfo{
		SpecVersion: sc.Compact{Number: sc.U8(1)},
		SpecName:    "test-lrui-spec-name",
	}
)

func Test_LastRuntimeUpgradeInfo_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := lrui.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectBytesLastRuntimeUpgradeInfo, buffer.Bytes())
}

func Test_LastRuntimeUpgradeInfo_Encode_InvalidSpecVersionNumber(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := lruiInvalidNumber.Encode(buffer)
	assert.Equal(t, err, errors.New("invalid SpecVersion of LastRuntimeUpgradeInfo"))
}

func Test_DecodeLastRuntimeUpgradeInfo(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesLastRuntimeUpgradeInfo)

	result, err := DecodeLastRuntimeUpgradeInfo(buffer)
	assert.NoError(t, err)

	assert.Equal(t, lrui, result)
}

func Test_LastRuntimeUpgradeInfo_Bytes(t *testing.T) {
	assert.Equal(t, expectBytesLastRuntimeUpgradeInfo, lrui.Bytes())
}
