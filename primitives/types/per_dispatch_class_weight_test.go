package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedPerDispatchClassWeightBytes, _ = hex.DecodeString("04080c101418")

	targetPerDispatchClassWeight = PerDispatchClassWeight{
		Normal:      Weight{sc.NewU64(1), sc.NewU64(2)},
		Operational: Weight{sc.NewU64(3), sc.NewU64(4)},
		Mandatory:   Weight{sc.NewU64(5), sc.NewU64(6)},
	}
)

func Test_PerDispatchClassWeight_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	err := targetPerDispatchClassWeight.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedPerDispatchClassWeightBytes, buf.Bytes())
}

func Test_DecodePerDispatchClassWeight(t *testing.T) {
	buf := bytes.NewBuffer(expectedPerDispatchClassWeightBytes)

	result, err := DecodePerDispatchClassWeight(buf, DecodeWeight)
	assert.NoError(t, err)

	assert.Equal(t, targetPerDispatchClassWeight, result)
}

func Test_PerDispatchClassWeight_Bytes(t *testing.T) {
	assert.Equal(t, expectedPerDispatchClassWeightBytes, targetPerDispatchClassWeight.Bytes())
}

func Test_PerDispatchClassWeight_Get(t *testing.T) {
	normal, err := targetPerDispatchClassWeight.Get(NewDispatchClassNormal())
	assert.NoError(t, err)
	assert.Equal(t, Weight{sc.NewU64(1), sc.NewU64(2)}, *normal)

	operational, err := targetPerDispatchClassWeight.Get(NewDispatchClassOperational())
	assert.NoError(t, err)
	assert.Equal(t, Weight{sc.NewU64(3), sc.NewU64(4)}, *operational)

	mandatory, err := targetPerDispatchClassWeight.Get(NewDispatchClassMandatory())
	assert.NoError(t, err)
	assert.Equal(t, Weight{sc.NewU64(5), sc.NewU64(6)}, *mandatory)
}

func Test_PerDispatchClassWeight_Get_TypeError(t *testing.T) {
	unknownDispatchClass := DispatchClass{sc.NewVaryingData(sc.U8(3))}

	result, err := targetPerDispatchClassWeight.Get(unknownDispatchClass)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())
	assert.Nil(t, result)
}
