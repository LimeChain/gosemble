package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/utils"
)

type BlockWeights struct {
	// Base weight of block execution.
	BaseBlock Weight
	// Maximal total weight consumed by all kinds of extrinsics (without `reserved` space).
	MaxBlock Weight
	// Weight limits for extrinsics of given dispatch class.
	PerClass PerDispatchClass[WeightsPerClass]
}

func (bw BlockWeights) Encode(buffer *bytes.Buffer) error {
	return utils.EncodeEach(buffer,
		bw.BaseBlock,
		bw.MaxBlock,
		bw.PerClass,
	)
}

func (bw BlockWeights) Bytes() []byte {
	return sc.EncodedBytes(bw)
}

// Get per-class weight settings.
func (bw BlockWeights) Get(class DispatchClass) (*WeightsPerClass, error) {
	isNormalDispatch, err := class.Is(DispatchClassNormal)
	if err != nil {
		return nil, err
	}
	if isNormalDispatch {
		return &bw.PerClass.Normal, nil
	}

	isOperationalDispatch, err := class.Is(DispatchClassOperational)
	if err != nil {
		return nil, err
	}
	if isOperationalDispatch {
		return &bw.PerClass.Operational, nil
	}

	isMandatoryDispatch, err := class.Is(DispatchClassMandatory)
	if err != nil {
		return nil, err
	}
	if isMandatoryDispatch {
		return &bw.PerClass.Mandatory, nil
	}

	return nil, newTypeError("DispatchClass")
}
