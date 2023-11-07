package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedPerDispatchClassBytes, _ = hex.DecodeString("010203")

	targetPerDispatchClass = PerDispatchClass[sc.U8]{
		Normal:      sc.U8(1),
		Operational: sc.U8(2),
		Mandatory:   sc.U8(3),
	}
)

func Test_PerDispatchClass_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	targetPerDispatchClass.Encode(buf)

	assert.Equal(t, expectedPerDispatchClassBytes, buf.Bytes())
}

func Test_DecodePerDispatchClass(t *testing.T) {
	buf := bytes.NewBuffer(expectedPerDispatchClassBytes)

	result, err := DecodePerDispatchClass(buf, sc.DecodeU8)
	assert.NoError(t, err)

	assert.Equal(t, targetPerDispatchClass, result)
}

func Test_PerDispatchClass_Bytes(t *testing.T) {
	assert.Equal(t, expectedPerDispatchClassBytes, targetPerDispatchClass.Bytes())
}

func Test_PerDispatchClass_Get(t *testing.T) {
	normal, err := targetPerDispatchClass.Get(NewDispatchClassNormal())
	assert.NoError(t, err)
	assert.Equal(t, sc.U8(1), *normal)

	operational, err := targetPerDispatchClass.Get(NewDispatchClassOperational())
	assert.NoError(t, err)
	assert.Equal(t, sc.U8(2), *operational)

	mandatory, err := targetPerDispatchClass.Get(NewDispatchClassMandatory())
	assert.NoError(t, err)
	assert.Equal(t, sc.U8(3), *mandatory)
}

func Test_PerDispatchClass_Get_TypeError(t *testing.T) {
	unknownDispatchClass := DispatchClass{sc.NewVaryingData(sc.U8(3))}

	result, err := targetPerDispatchClass.Get(unknownDispatchClass)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())
	assert.Nil(t, result)
}
