package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type PerDispatchClassWeight struct {
	// Value for `Normal` extrinsics.
	Normal Weight
	// Value for `Operational` extrinsics.
	Operational Weight
	// Value for `Mandatory` extrinsics.
	Mandatory Weight
}

func (pdcw PerDispatchClassWeight) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		pdcw.Normal,
		pdcw.Operational,
		pdcw.Mandatory,
	)
}

func DecodePerDispatchClassWeight(buffer *bytes.Buffer, decodeWeight func(buffer *bytes.Buffer) (Weight, error)) (PerDispatchClassWeight, error) {
	normal, err := decodeWeight(buffer)
	if err != nil {
		return PerDispatchClassWeight{}, err
	}
	operational, err := decodeWeight(buffer)
	if err != nil {
		return PerDispatchClassWeight{}, err
	}
	mandatory, err := decodeWeight(buffer)
	if err != nil {
		return PerDispatchClassWeight{}, err
	}
	return PerDispatchClassWeight{
		Normal:      normal,
		Operational: operational,
		Mandatory:   mandatory,
	}, nil
}

func (pdcw PerDispatchClassWeight) Bytes() []byte {
	return sc.EncodedBytes(pdcw)
}

// Get current value for given class.
func (pdcw *PerDispatchClassWeight) Get(class DispatchClass) (*Weight, error) {
	switch class.VaryingData[0] {
	case DispatchClassNormal:
		return &pdcw.Normal, nil
	case DispatchClassOperational:
		return &pdcw.Operational, nil
	case DispatchClassMandatory:
		return &pdcw.Mandatory, nil
	default:
		return nil, newTypeError("DispatchClass")
	}
}
