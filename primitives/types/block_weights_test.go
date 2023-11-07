package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedBlockWeightsBytes, _ = hex.DecodeString(
		"04080c101418011c20012428012c303438013c40014448014c505458015c60016468016c70",
	)
)

var (
	targetBlockWeights = BlockWeights{
		BaseBlock: WeightFromParts(1, 2),
		MaxBlock:  WeightFromParts(3, 4),
		PerClass: PerDispatchClass[WeightsPerClass]{
			Normal: WeightsPerClass{
				BaseExtrinsic: WeightFromParts(5, 6),
				MaxExtrinsic:  sc.NewOption[Weight](WeightFromParts(7, 8)),
				MaxTotal:      sc.NewOption[Weight](WeightFromParts(9, 10)),
				Reserved:      sc.NewOption[Weight](WeightFromParts(11, 12)),
			},
			Operational: WeightsPerClass{
				BaseExtrinsic: WeightFromParts(13, 14),
				MaxExtrinsic:  sc.NewOption[Weight](WeightFromParts(15, 16)),
				MaxTotal:      sc.NewOption[Weight](WeightFromParts(17, 18)),
				Reserved:      sc.NewOption[Weight](WeightFromParts(19, 20)),
			},
			Mandatory: WeightsPerClass{
				BaseExtrinsic: WeightFromParts(21, 22),
				MaxExtrinsic:  sc.NewOption[Weight](WeightFromParts(23, 24)),
				MaxTotal:      sc.NewOption[Weight](WeightFromParts(25, 26)),
				Reserved:      sc.NewOption[Weight](WeightFromParts(27, 28)),
			},
		},
	}
)

func Test_BlockWeights_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetBlockWeights.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectedBlockWeightsBytes, buffer.Bytes())
}

func Test_BlockWeights_Bytes(t *testing.T) {
	assert.Equal(t, expectedBlockWeightsBytes, targetBlockWeights.Bytes())
}

func Test_BlockWeights_Get(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       DispatchClass
		expectation Weight
	}{
		{
			label:       "Normal",
			input:       NewDispatchClassNormal(),
			expectation: targetBlockWeights.PerClass.Normal.BaseExtrinsic,
		},
		{
			label:       "Operational",
			input:       NewDispatchClassOperational(),
			expectation: targetBlockWeights.PerClass.Operational.BaseExtrinsic,
		},
		{
			label:       "Mandatory",
			input:       NewDispatchClassMandatory(),
			expectation: targetBlockWeights.PerClass.Mandatory.BaseExtrinsic,
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			result, err := targetBlockWeights.Get(testExample.input)

			assert.NoError(t, err)
			assert.Equal(t, testExample.expectation, result.BaseExtrinsic)
		})
	}
}

func Test_BlockWeights_Get_TypeError(t *testing.T) {
	unknownDispatchClass := DispatchClass{sc.NewVaryingData(sc.U8(3))}

	res, err := targetBlockWeights.Get(unknownDispatchClass)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())
	assert.Nil(t, res)
}
