package types

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedConsumedWeightBytes, _ = hex.DecodeString("04080c101418")

	targetConsumedWeight = ConsumedWeight{
		Normal:      WeightFromParts(1, 2),
		Operational: WeightFromParts(3, 4),
		Mandatory:   WeightFromParts(5, 6),
	}
)

func Test_ConsumedWeight_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetConsumedWeight.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectedConsumedWeightBytes, buffer.Bytes())
}

func Test_DecodeConsumedWeight(t *testing.T) {
	buf := bytes.NewBuffer(expectedConsumedWeightBytes)

	result, err := DecodeConsumedWeight(buf)
	assert.NoError(t, err)
	assert.Equal(t, targetConsumedWeight, result)
}

func Test_ConsumedWeight_Bytes(t *testing.T) {
	assert.Equal(t, expectedConsumedWeightBytes, targetConsumedWeight.Bytes())
}

func Test_ConsumedWeight_Get_Normal(t *testing.T) {
	result, err := targetConsumedWeight.Get(NewDispatchClassNormal())

	assert.NoError(t, err)
	assert.Equal(t, WeightFromParts(1, 2), *result)
}

func Test_ConsumedWeight_Get_Operational(t *testing.T) {
	result, err := targetConsumedWeight.Get(NewDispatchClassOperational())

	assert.NoError(t, err)
	assert.Equal(t, WeightFromParts(3, 4), *result)
}

func Test_ConsumedWeight_Get_Mandatory(t *testing.T) {
	result, err := targetConsumedWeight.Get(NewDispatchClassMandatory())

	assert.NoError(t, err)
	assert.Equal(t, WeightFromParts(5, 6), *result)
}

func Test_ConsumedWeight_Get_UnknownPanics(t *testing.T) {
	unknownDispatchClass := DispatchClass{sc.NewVaryingData(sc.U8(3))}

	result, err := targetConsumedWeight.Get(unknownDispatchClass)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())
	assert.Nil(t, result)
}

func Test_ConsumedWeight_Total(t *testing.T) {
	result, err := targetConsumedWeight.Total()

	assert.NoError(t, err)
	assert.Equal(t, WeightFromParts(9, 12), result)
}

func Test_ConsumedWeight_SaturatingAdd(t *testing.T) {
	targetConsumedWeight.SaturatingAdd(
		WeightFromParts(math.MaxUint64, math.MaxUint64),
		NewDispatchClassNormal(),
	)

	consumedWeightNormal, err := targetConsumedWeight.Get(NewDispatchClassNormal())
	assert.NoError(t, err)
	assert.Equal(t, WeightFromParts(math.MaxUint64, math.MaxUint64), *consumedWeightNormal)

	consumedWeightOperational, err := targetConsumedWeight.Get(NewDispatchClassOperational())
	assert.NoError(t, err)
	assert.Equal(t, WeightFromParts(3, 4), *consumedWeightOperational)

	consumedWeightMandatory, err := targetConsumedWeight.Get(NewDispatchClassMandatory())
	assert.NoError(t, err)
	assert.Equal(t, WeightFromParts(5, 6), *consumedWeightMandatory)
}

func Test_ConsumedWeight_Accrue(t *testing.T) {
	targetConsumedWeight := ConsumedWeight{
		Normal:      WeightFromParts(1, 2),
		Operational: WeightFromParts(3, 4),
		Mandatory:   WeightFromParts(5, 6),
	}

	targetConsumedWeight.Accrue(
		WeightFromParts(math.MaxUint64, math.MaxUint64),
		NewDispatchClassNormal(),
	)

	assert.Equal(t,
		WeightFromParts(math.MaxUint64, math.MaxUint64),
		targetConsumedWeight.Normal,
	)

	assert.Equal(t,
		WeightFromParts(3, 4),
		targetConsumedWeight.Operational,
	)

	assert.Equal(t,
		WeightFromParts(5, 6),
		targetConsumedWeight.Mandatory,
	)
}

func Test_ConsumedWeight_CheckedAccrue_Ok(t *testing.T) {
	targetConsumedWeight := ConsumedWeight{
		Normal:      WeightFromParts(1, 2),
		Operational: WeightFromParts(3, 4),
		Mandatory:   WeightFromParts(5, 6),
	}

	err := targetConsumedWeight.CheckedAccrue(
		WeightFromParts(1, 1),
		NewDispatchClassNormal(),
	)
	assert.NoError(t, err)

	assert.Equal(t,
		WeightFromParts(2, 3),
		targetConsumedWeight.Normal,
	)

	assert.Equal(t,
		WeightFromParts(3, 4),
		targetConsumedWeight.Operational,
	)

	assert.Equal(t,
		WeightFromParts(5, 6),
		targetConsumedWeight.Mandatory,
	)
}

func Test_ConsumedWeight_CheckedAccrue_RefTimeOverflow(t *testing.T) {
	targetConsumedWeight := ConsumedWeight{
		Normal:      WeightFromParts(1, 2),
		Operational: WeightFromParts(3, 4),
		Mandatory:   WeightFromParts(5, 6),
	}

	err := targetConsumedWeight.CheckedAccrue(
		WeightFromParts(math.MaxUint64, 1),
		NewDispatchClassNormal(),
	)
	assert.Equal(t, errors.New("overflow"), err)

	assert.Equal(t,
		WeightFromParts(1, 2),
		targetConsumedWeight.Normal,
	)

	assert.Equal(t,
		WeightFromParts(3, 4),
		targetConsumedWeight.Operational,
	)

	assert.Equal(t,
		WeightFromParts(5, 6),
		targetConsumedWeight.Mandatory,
	)
}

func Test_ConsumedWeight_CheckedAccrue_ProofSizeOverflow(t *testing.T) {
	targetConsumedWeight := ConsumedWeight{
		Normal:      WeightFromParts(1, 2),
		Operational: WeightFromParts(3, 4),
		Mandatory:   WeightFromParts(5, 6),
	}

	err := targetConsumedWeight.CheckedAccrue(
		WeightFromParts(1, math.MaxUint64),
		NewDispatchClassNormal(),
	)
	assert.Equal(t, errors.New("overflow"), err)

	err = targetConsumedWeight.CheckedAccrue(
		WeightFromParts(math.MaxUint64, 1),
		NewDispatchClassNormal(),
	)
	assert.Equal(t, errors.New("overflow"), err)

	assert.Equal(t,
		WeightFromParts(1, 2),
		targetConsumedWeight.Normal,
	)

	assert.Equal(t,
		WeightFromParts(3, 4),
		targetConsumedWeight.Operational,
	)

	assert.Equal(t,
		WeightFromParts(5, 6),
		targetConsumedWeight.Mandatory,
	)
}

func Test_ConsumedWeight_Reduce(t *testing.T) {
	targetConsumedWeight := ConsumedWeight{
		Normal:      WeightFromParts(0, 2),
		Operational: WeightFromParts(3, 4),
		Mandatory:   WeightFromParts(5, 6),
	}

	targetConsumedWeight.Reduce(
		WeightFromParts(1, 1),
		NewDispatchClassNormal(),
	)

	assert.Equal(t,
		WeightFromParts(0, 1),
		targetConsumedWeight.Normal,
	)

	assert.Equal(t,
		WeightFromParts(3, 4),
		targetConsumedWeight.Operational,
	)

	assert.Equal(t,
		WeightFromParts(5, 6),
		targetConsumedWeight.Mandatory,
	)
}
