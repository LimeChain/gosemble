package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectBytesRuntimeVersion, _ = hex.DecodeString("44746573742d737065632d76657273696f6e44746573742d696d706c2d76657273696f6e0100000002000000030000000474657374696e6730060000000400000005")
)

var (
	runtimeVersion = RuntimeVersion{
		SpecName:           "test-spec-version",
		ImplName:           "test-impl-version",
		AuthoringVersion:   1,
		SpecVersion:        2,
		ImplVersion:        3,
		TransactionVersion: 4,
		StateVersion:       5,
		Apis: sc.Sequence[ApiItem]{
			{
				Name:    sc.BytesToFixedSequenceU8(apiName[:]),
				Version: 6,
			},
		},
	}
)

func Test_RuntimeVersion_SetApis(t *testing.T) {
	target := RuntimeVersion{
		SpecName:           "test-spec-version",
		ImplName:           "test-impl-version",
		AuthoringVersion:   1,
		SpecVersion:        2,
		ImplVersion:        3,
		TransactionVersion: 4,
		StateVersion:       5,
	}

	apiItems := sc.Sequence[ApiItem]{
		{
			Name:    sc.BytesToFixedSequenceU8(apiName[:]),
			Version: 6,
		},
	}

	target.SetApis(apiItems)

	assert.Equal(t, runtimeVersion, target)
}

func Test_RuntimeVersion_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	runtimeVersion.Encode(buffer)

	assert.Equal(t, expectBytesRuntimeVersion, buffer.Bytes())
}

func Test_DecodeRuntimeVersion(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesRuntimeVersion)

	result, err := DecodeRuntimeVersion(buffer)
	assert.Nil(t, err)

	assert.Equal(t, runtimeVersion, result)
}

func Test_Runtime_Bytes(t *testing.T) {
	result := runtimeVersion.Bytes()

	assert.Equal(t, expectBytesRuntimeVersion, result)
}
