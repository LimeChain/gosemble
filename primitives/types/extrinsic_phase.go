package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	// PhaseApplyExtrinsic Applying an extrinsic.
	PhaseApplyExtrinsic sc.U8 = iota

	// PhaseFinalization Finalizing the block.
	PhaseFinalization

	// PhaseInitialization Initializing the block.
	PhaseInitialization
)

type ExtrinsicPhase = sc.VaryingData

func NewExtrinsicPhaseApply(index sc.U32) ExtrinsicPhase {
	return sc.NewVaryingData(PhaseApplyExtrinsic, index)
}

func NewExtrinsicPhaseFinalization() ExtrinsicPhase {
	return sc.NewVaryingData(PhaseFinalization)
}

func NewExtrinsicPhaseInitialization() ExtrinsicPhase {
	return sc.NewVaryingData(PhaseInitialization)
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
		log.Critical("invalid ExtrinsicPhase type")
	}

	panic("unreachable")
}
