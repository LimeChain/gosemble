package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

// An object to track the currently used extrinsic weight in a block.
type ConsumedWeight PerDispatchClass[Weight]

func (cw ConsumedWeight) Encode(buffer *bytes.Buffer) {
	PerDispatchClass[Weight](cw).Encode(buffer)
}

func DecodeConsumedWeight(buffer *bytes.Buffer) ConsumedWeight {
	return ConsumedWeight{
		Normal:      DecodeWeight(buffer),
		Operational: DecodeWeight(buffer),
		Mandatory:   DecodeWeight(buffer),
	}
}

func (cw ConsumedWeight) Bytes() []byte {
	return sc.EncodedBytes(cw)
}

// Get current value for given class.
func (cw *ConsumedWeight) Get(class DispatchClass) *Weight {
	switch class.VaryingData[0] {
	case DispatchClassNormal:
		return &cw.Normal
	case DispatchClassOperational:
		return &cw.Operational
	case DispatchClassMandatory:
		return &cw.Mandatory
	default:
		log.Critical("invalid DispatchClass type")
	}

	panic("unreachable")
}

// Returns the total weight consumed by all extrinsics in the block.
//
// Saturates on overflow.
func (cw ConsumedWeight) Total() Weight {
	sum := WeightZero()
	for _, class := range []DispatchClass{NewDispatchClassNormal(), NewDispatchClassOperational(), NewDispatchClassMandatory()} {
		sum = sum.SaturatingAdd(*cw.Get(class))
	}
	return sum
}

// SaturatingAdd Increase the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) SaturatingAdd(weight Weight, class DispatchClass) {
	weightForClass := cw.Get(class)
	weightForClass.RefTime = sc.SaturatingAddU64(weightForClass.RefTime, weight.RefTime)
	weightForClass.ProofSize = sc.SaturatingAddU64(weightForClass.ProofSize, weight.ProofSize)
}

// Accrue Increase the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) Accrue(weight Weight, class DispatchClass) {
	weightForClass := cw.Get(class)
	weightForClass.SaturatingAccrue(weight)
}

// CheckedAccrue Try to increase the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) CheckedAccrue(weight Weight, class DispatchClass) (sc.Empty, error) {
	weightForClass := cw.Get(class)
	refTime, err := sc.CheckedAddU64(weightForClass.RefTime, weight.RefTime)
	if err != nil {
		return sc.Empty{}, err
	}

	weightForClass.RefTime = refTime

	proofSize, err := sc.CheckedAddU64(weightForClass.ProofSize, weight.ProofSize)
	if err != nil {
		return sc.Empty{}, err
	}

	weightForClass.ProofSize = proofSize

	return sc.Empty{}, nil
}

// Reduce the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) Reduce(weight Weight, class DispatchClass) {
	weightForClass := cw.Get(class)
	weightForClass.SaturatingReduce(weight)
}
