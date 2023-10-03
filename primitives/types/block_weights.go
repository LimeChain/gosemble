package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type BlockWeights struct {
	// Base weight of block execution.
	BaseBlock Weight
	// Maximal total weight consumed by all kinds of extrinsics (without `reserved` space).
	MaxBlock Weight
	// Weight limits for extrinsics of given dispatch class.
	PerClass PerDispatchClass[WeightsPerClass]
}

func (bw BlockWeights) Encode(buffer *bytes.Buffer) {
	bw.BaseBlock.Encode(buffer)
	bw.MaxBlock.Encode(buffer)
	bw.PerClass.Encode(buffer)
}

func (bw BlockWeights) Bytes() []byte {
	return sc.EncodedBytes(bw)
}

// Get per-class weight settings.
func (bw BlockWeights) Get(class DispatchClass) *WeightsPerClass {
	if class.Is(DispatchClassNormal) {
		return &bw.PerClass.Normal
	} else if class.Is(DispatchClassOperational) {
		return &bw.PerClass.Operational
	} else if class.Is(DispatchClassMandatory) {
		return &bw.PerClass.Mandatory
	} else {
		log.Critical("Invalid dispatch class")
	}

	panic("unreachable")
}
