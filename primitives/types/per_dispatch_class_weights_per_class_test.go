package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedPerDispatchClassWeightsPerClassBytes, _ = hex.DecodeString("0408010c10011418011c200408010c10011418011c200408010c10011418011c20")

	wpc = WeightsPerClass{
		BaseExtrinsic: WeightFromParts(1, 2),
		MaxExtrinsic:  sc.NewOption[Weight](WeightFromParts(3, 4)),
		MaxTotal:      sc.NewOption[Weight](WeightFromParts(5, 6)),
		Reserved:      sc.NewOption[Weight](WeightFromParts(7, 8)),
	}

	targetPerDispatchClassWeightsPerClass = PerDispatchClassWeightsPerClass{
		Normal:      wpc,
		Operational: wpc,
		Mandatory:   wpc,
	}
)

func Test_PerDispatchClassWeightsPerClass_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	err := targetPerDispatchClassWeightsPerClass.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedPerDispatchClassWeightsPerClassBytes, buf.Bytes())
}

func Test_DecodePerDispatchClassWeightsPerClass(t *testing.T) {
	buf := bytes.NewBuffer(expectedPerDispatchClassWeightsPerClassBytes)

	result, err := DecodePerDispatchClassWeightPerClass(buf, DecodeWeightsPerClass)
	assert.NoError(t, err)

	assert.Equal(t, targetPerDispatchClassWeightsPerClass, result)
}

func Test_PerDispatchClassWeightsPerClass_Bytes(t *testing.T) {
	assert.Equal(t, expectedPerDispatchClassWeightsPerClassBytes, targetPerDispatchClassWeightsPerClass.Bytes())
}

func Test_PerDispatchClassWeightsPerClass_Get(t *testing.T) {
	normal, err := targetPerDispatchClassWeightsPerClass.Get(NewDispatchClassNormal())
	assert.NoError(t, err)
	assert.Equal(t, WeightsPerClass{
		BaseExtrinsic: WeightFromParts(1, 2),
		MaxExtrinsic:  sc.NewOption[Weight](WeightFromParts(3, 4)),
		MaxTotal:      sc.NewOption[Weight](WeightFromParts(5, 6)),
		Reserved:      sc.NewOption[Weight](WeightFromParts(7, 8)),
	}, *normal)

	operational, err := targetPerDispatchClassWeightsPerClass.Get(NewDispatchClassOperational())
	assert.NoError(t, err)
	assert.Equal(t, WeightsPerClass{
		BaseExtrinsic: WeightFromParts(1, 2),
		MaxExtrinsic:  sc.NewOption[Weight](WeightFromParts(3, 4)),
		MaxTotal:      sc.NewOption[Weight](WeightFromParts(5, 6)),
		Reserved:      sc.NewOption[Weight](WeightFromParts(7, 8)),
	}, *operational)

	mandatory, err := targetPerDispatchClassWeightsPerClass.Get(NewDispatchClassMandatory())
	assert.NoError(t, err)
	assert.Equal(t, WeightsPerClass{
		BaseExtrinsic: WeightFromParts(1, 2),
		MaxExtrinsic:  sc.NewOption[Weight](WeightFromParts(3, 4)),
		MaxTotal:      sc.NewOption[Weight](WeightFromParts(5, 6)),
		Reserved:      sc.NewOption[Weight](WeightFromParts(7, 8)),
	}, *mandatory)
}

func Test_PerDispatchClassWeightsPerClass_Get_TypeError(t *testing.T) {
	unknownDispatchClass := DispatchClass{sc.NewVaryingData(sc.U8(3))}

	result, err := targetPerDispatchClassWeightsPerClass.Get(unknownDispatchClass)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())
	assert.Nil(t, result)
}
