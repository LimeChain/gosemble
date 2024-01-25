package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// An object to track the currently used extrinsic weight in a block.
type ConsumedWeight PerDispatchClassWeight

func (cw ConsumedWeight) Encode(buffer *bytes.Buffer) error {
	return PerDispatchClass[Weight](cw).Encode(buffer)
}

func DecodeConsumedWeight(buffer *bytes.Buffer) (ConsumedWeight, error) {
	normal, err := DecodeWeight(buffer)
	if err != nil {
		return ConsumedWeight{}, err
	}
	operational, err := DecodeWeight(buffer)
	if err != nil {
		return ConsumedWeight{}, err
	}
	mandatory, err := DecodeWeight(buffer)
	if err != nil {
		return ConsumedWeight{}, err
	}
	return ConsumedWeight{
		Normal:      normal,
		Operational: operational,
		Mandatory:   mandatory,
	}, nil
}

func (cw ConsumedWeight) Bytes() []byte {
	return sc.EncodedBytes(cw)
}

// Get current value for given class.
func (cw *ConsumedWeight) Get(class DispatchClass) (*Weight, error) {
	switch class.VaryingData[0] {
	case DispatchClassNormal:
		return &cw.Normal, nil
	case DispatchClassOperational:
		return &cw.Operational, nil
	case DispatchClassMandatory:
		return &cw.Mandatory, nil
	default:
		return nil, newTypeError("DispatchClass")
	}
}

// Returns the total weight consumed by all extrinsics in the block.
//
// Saturates on overflow.
func (cw ConsumedWeight) Total() (Weight, error) {
	sum := WeightZero()
	for _, class := range DispatchClassAll() {
		weightForClass, err := cw.Get(class)
		if err != nil {
			return Weight{}, err
		}
		sum = sum.SaturatingAdd(*weightForClass)
	}
	return sum, nil
}

// SaturatingAdd Increase the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) SaturatingAdd(weight Weight, class DispatchClass) error {
	weightForClass, err := cw.Get(class)
	if err != nil {
		return err
	}
	weightForClass.RefTime = sc.SaturatingAddU64(weightForClass.RefTime, weight.RefTime)
	weightForClass.ProofSize = sc.SaturatingAddU64(weightForClass.ProofSize, weight.ProofSize)
	return nil
}

// Accrue Increase the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) Accrue(weight Weight, class DispatchClass) error {
	weightForClass, err := cw.Get(class)
	if err != nil {
		return err
	}
	weightForClass.SaturatingAccrue(weight)
	return nil
}

// CheckedAccrue Try to increase the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) CheckedAccrue(weight Weight, class DispatchClass) error {
	weightForClass, err := cw.Get(class)
	if err != nil {
		return err
	}
	refTime, err := sc.CheckedAddU64(weightForClass.RefTime, weight.RefTime)
	if err != nil {
		return err
	}

	proofSize, err := sc.CheckedAddU64(weightForClass.ProofSize, weight.ProofSize)
	if err != nil {
		return err
	}

	weightForClass.RefTime = refTime
	weightForClass.ProofSize = proofSize

	return nil
}

// Reduce the weight of the given class. Saturates at the numeric bounds.
func (cw *ConsumedWeight) Reduce(weight Weight, class DispatchClass) error {
	weightForClass, err := cw.Get(class)
	if err != nil {
		return err
	}
	weightForClass.SaturatingReduce(weight)
	return nil
}
