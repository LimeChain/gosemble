package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
)

const (
	// PhaseApplyExtrinsic Applying an extrinsic.
	PhaseApplyExtrinsic sc.U8 = iota

	// PhaseFinalization Finalizing the block.
	PhaseFinalization

	// PhaseInitialization Initializing the block.
	PhaseInitialization
)

type ExtrinsicPhase struct {
	sc.VaryingData
}

func NewExtrinsicPhaseApply(index sc.U32) ExtrinsicPhase {
	return ExtrinsicPhase{sc.NewVaryingData(PhaseApplyExtrinsic, index)}
}

func NewExtrinsicPhaseFinalization() ExtrinsicPhase {
	return ExtrinsicPhase{sc.NewVaryingData(PhaseFinalization)}
}

func NewExtrinsicPhaseInitialization() ExtrinsicPhase {
	return ExtrinsicPhase{sc.NewVaryingData(PhaseInitialization)}
}

func DecodeExtrinsicPhase(buffer *bytes.Buffer) (ExtrinsicPhase, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return ExtrinsicPhase{}, err
	}

	switch b {
	case PhaseApplyExtrinsic:
		index, err := sc.DecodeU32(buffer)
		if err != nil {
			return ExtrinsicPhase{}, err
		}
		return NewExtrinsicPhaseApply(index), nil
	case PhaseFinalization:
		return NewExtrinsicPhaseFinalization(), nil
	case PhaseInitialization:
		return NewExtrinsicPhaseInitialization(), nil
	default:
		return ExtrinsicPhase{}, newTypeError("ExtrinsicPhase")
	}
}

func (p *ExtrinsicPhase) Metadata() *MetadataTypeDefinition {
	def := NewMetadataTypeDefinitionVariant(
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
	return &def
}
