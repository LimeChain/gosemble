package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	// Applying an extrinsic.
	PhaseApplyExtrinsic sc.U32 = iota

	// Finalizing the block.
	PhaseFinalization

	// Initializing the block.
	PhaseInitialization
)

type Phase sc.VaryingData

func NewPhase(values ...sc.Encodable) Phase {
	switch values[0] {
	case PhaseApplyExtrinsic:
		return Phase(sc.NewVaryingData(values[0:2]...))
	case PhaseFinalization, PhaseInitialization:
		return Phase(sc.NewVaryingData(values[0]))
	default:
		log.Critical("invalid phase type")
	}

	panic("unreachable")
}

func (p Phase) Encode(buffer *bytes.Buffer) {
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

func Decode(buffer *bytes.Buffer) Phase {
	b := sc.DecodeU8(buffer)

	switch b {
	case sc.U8(0):
		value := sc.DecodeU32(buffer)
		return NewPhase(PhaseApplyExtrinsic, value)
	case sc.U8(1):
		return NewPhase(PhaseFinalization)
	case sc.U8(2):
		return NewPhase(PhaseInitialization)
	default:
		log.Critical("invalid Phase type")
	}

	panic("unreachable")
}

func (p Phase) Bytes() []byte {
	return sc.EncodedBytes(p)
}
