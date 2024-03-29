package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/stretchr/testify/assert"
)

func Test_NewExtrinsicPhaseApply(t *testing.T) {
	index := sc.U32(2)
	expect := ExtrinsicPhase{sc.NewVaryingData(PhaseApplyExtrinsic, index)}

	assert.Equal(t, expect, NewExtrinsicPhaseApply(index))
}

func Test_NewExtrinsicPhaseFinalization(t *testing.T) {
	expect := ExtrinsicPhase{sc.NewVaryingData(PhaseFinalization)}

	assert.Equal(t, expect, NewExtrinsicPhaseFinalization())
}

func Test_NewExtrinsicPhaseInitialization(t *testing.T) {
	expect := ExtrinsicPhase{sc.NewVaryingData(PhaseInitialization)}

	assert.Equal(t, expect, NewExtrinsicPhaseInitialization())
}

func Test_DecodeExtrinsicPhase_ApplyExtrinsic(t *testing.T) {
	index := sc.U32(2)

	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)
	buffer.Write(index.Bytes())

	result, err := DecodeExtrinsicPhase(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewExtrinsicPhaseApply(index), result)
}

func Test_DecodeExtrinsicPhase_Finalization(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(1)

	result, err := DecodeExtrinsicPhase(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewExtrinsicPhaseFinalization(), result)
}

func Test_DecodeExtrinsicPhase_Initialization(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(2)

	result, err := DecodeExtrinsicPhase(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewExtrinsicPhaseInitialization(), result)
}

func Test_DecodeExtrinsicPhase_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(3)

	res, err := DecodeExtrinsicPhase(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'ExtrinsicPhase' type", err.Error())
	assert.Equal(t, ExtrinsicPhase{}, res)
}

func Test_DecodeExtrinsicPhase_Empty(t *testing.T) {
	buffer := &bytes.Buffer{}

	res, err := DecodeExtrinsicPhase(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, ExtrinsicPhase{}, res)
}

func Test_ExtrinsicPhase_MetadataDefinition(t *testing.T) {
	expect := NewMetadataTypeDefinitionVariant(
		sc.Sequence[MetadataDefinitionVariant]{
			NewMetadataDefinitionVariant(
				"ApplyExtrinsic",
				sc.Sequence[MetadataTypeDefinitionField]{
					NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),
				},
				PhaseApplyExtrinsic,
				"Phase.ApplyExtrinsic"),
			NewMetadataDefinitionVariant(
				"Finalization",
				sc.Sequence[MetadataTypeDefinitionField]{},
				PhaseFinalization,
				"Phase.Finalization"),
			NewMetadataDefinitionVariant(
				"Initialization",
				sc.Sequence[MetadataTypeDefinitionField]{},
				PhaseInitialization,
				"Phase.Initialization"),
		})

	assert.Equal(t, &expect, new(ExtrinsicPhase).MetadataDefinition())
}
