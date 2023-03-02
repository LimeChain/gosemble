package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	// Applying an extrinsic.
	PhaseApplyExtrinsic sc.U8 = iota

	// Finalizing the block.
	PhaseFinalization

	// Initializing the block.
	PhaseInitialization
)

type ExtrinsicPhase sc.VaryingData

func NewExtrinsicPhase(values ...sc.Encodable) ExtrinsicPhase {
	switch values[0] {
	case PhaseApplyExtrinsic:
		return ExtrinsicPhase(sc.NewVaryingData(values[0:2]...))
	case PhaseFinalization, PhaseInitialization:
		return ExtrinsicPhase(sc.NewVaryingData(values[0]))
	default:
		log.Critical("invalid phase type")
	}

	panic("unreachable")
}

func (p ExtrinsicPhase) Encode(buffer *bytes.Buffer) {
	switch p[0] {
	case PhaseApplyExtrinsic:
		sc.U8(0).Encode(buffer)
		p[1].Encode(buffer)
	case PhaseFinalization:
		sc.U8(1).Encode(buffer)
	case PhaseInitialization:
		sc.U8(2).Encode(buffer)
	}
}

func DecodeExtrinsicPhase(buffer *bytes.Buffer) ExtrinsicPhase {
	b := sc.DecodeU8(buffer)

	switch b {
	case PhaseApplyExtrinsic:
		value := sc.DecodeU32(buffer)
		return NewExtrinsicPhase(PhaseApplyExtrinsic, value)
	case PhaseFinalization:
		return NewExtrinsicPhase(PhaseFinalization)
	case PhaseInitialization:
		return NewExtrinsicPhase(PhaseInitialization)
	default:
		log.Critical("invalid Phase type")
	}

	panic("unreachable")
}

func (p ExtrinsicPhase) Bytes() []byte {
	return sc.EncodedBytes(p)
}
