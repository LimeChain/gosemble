package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedPerDispatchClassU32Bytes, _ = hex.DecodeString("020000000300000004000000")

	targetPerDispatchClassU32 = PerDispatchClassU32{
		Normal:      sc.U32(2),
		Operational: sc.U32(3),
		Mandatory:   sc.U32(4),
	}
)

func Test_PerDispatchClassU32_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	err := targetPerDispatchClassU32.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedPerDispatchClassU32Bytes, buf.Bytes())
}

func Test_DecodePerDispatchClassU32(t *testing.T) {
	buf := bytes.NewBuffer(expectedPerDispatchClassU32Bytes)

	result, err := DecodePerDispatchClassU32(buf, sc.DecodeU32)
	assert.NoError(t, err)

	assert.Equal(t, targetPerDispatchClassU32, result)
}

func Test_PerDispatchClassU32_Bytes(t *testing.T) {
	assert.Equal(t, expectedPerDispatchClassU32Bytes, targetPerDispatchClassU32.Bytes())
}

func Test_PerDispatchClassU32_Get(t *testing.T) {
	normal, err := targetPerDispatchClassU32.Get(NewDispatchClassNormal())
	assert.NoError(t, err)
	assert.Equal(t, sc.U32(2), *normal)

	operational, err := targetPerDispatchClassU32.Get(NewDispatchClassOperational())
	assert.NoError(t, err)
	assert.Equal(t, sc.U32(3), *operational)

	mandatory, err := targetPerDispatchClassU32.Get(NewDispatchClassMandatory())
	assert.NoError(t, err)
	assert.Equal(t, sc.U32(4), *mandatory)
}

func Test_PerDispatchClassU32_Get_TypeError(t *testing.T) {
	unknownDispatchClass := DispatchClass{sc.NewVaryingData(sc.U8(3))}

	result, err := targetPerDispatchClassU32.Get(unknownDispatchClass)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())
	assert.Nil(t, result)
}
