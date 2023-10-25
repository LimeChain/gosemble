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

	targetConsumedWeight.Encode(buffer)

	assert.Equal(t, expectedConsumedWeightBytes, buffer.Bytes())
}

func Test_DecodeConsumedWeight(t *testing.T) {
	buf := bytes.NewBuffer(expectedConsumedWeightBytes)

	assert.Equal(t, targetConsumedWeight, DecodeConsumedWeight(buf))
}

func Test_ConsumedWeight_Bytes(t *testing.T) {
	assert.Equal(t, expectedConsumedWeightBytes, targetConsumedWeight.Bytes())
}

func Test_ConsumedWeight_Get_Normal(t *testing.T) {
	assert.Equal(t,
		WeightFromParts(1, 2),
		*targetConsumedWeight.Get(NewDispatchClassNormal()),
	)
}

func Test_ConsumedWeight_Get_Operational(t *testing.T) {
	assert.Equal(t,
		WeightFromParts(3, 4),
		*targetConsumedWeight.Get(NewDispatchClassOperational()),
	)
}

func Test_ConsumedWeight_Get_Mandatory(t *testing.T) {
	assert.Equal(t,
		WeightFromParts(5, 6),
		*targetConsumedWeight.Get(NewDispatchClassMandatory()),
	)
}

func Test_ConsumedWeight_Get_UnknownPanics(t *testing.T) {
	unknownDispatchClass := DispatchClass{sc.NewVaryingData(sc.U8(3))}

	assert.PanicsWithValue(t, "invalid DispatchClass type", func() {
		targetConsumedWeight.Get(unknownDispatchClass)
	})
}

func Test_ConsumedWeight_Total(t *testing.T) {
	assert.Equal(t,
		WeightFromParts(9, 12),
		targetConsumedWeight.Total(),
	)
}

func Test_ConsumedWeight_SaturatingAdd(t *testing.T) {
	targetConsumedWeight.SaturatingAdd(
		WeightFromParts(math.MaxUint64, math.MaxUint64),
		NewDispatchClassNormal(),
	)

	assert.Equal(t,
		WeightFromParts(math.MaxUint64, math.MaxUint64),
		*targetConsumedWeight.Get(NewDispatchClassNormal()),
	)

	assert.Equal(t,
		WeightFromParts(3, 4),
		*targetConsumedWeight.Get(NewDispatchClassOperational()),
	)

	assert.Equal(t,
		WeightFromParts(5, 6),
		*targetConsumedWeight.Get(NewDispatchClassMandatory()),
	)
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

	_, err := targetConsumedWeight.CheckedAccrue(
		WeightFromParts(1, 1),
		NewDispatchClassNormal(),
	)
	assert.Nil(t, err)

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

	_, err := targetConsumedWeight.CheckedAccrue(
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

	_, err := targetConsumedWeight.CheckedAccrue(
		WeightFromParts(1, math.MaxUint64),
		NewDispatchClassNormal(),
	)
	assert.Equal(t, errors.New("overflow"), err)

	_, err = targetConsumedWeight.CheckedAccrue(
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
