package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// A struct holding value for each `DispatchClass`.
type PerDispatchClass[T sc.Encodable] struct {
	// Value for `Normal` extrinsics.
	Normal T
	// Value for `Operational` extrinsics.
	Operational T
	// Value for `Mandatory` extrinsics.
	Mandatory T
}

func (pdc PerDispatchClass[T]) Encode(buffer *bytes.Buffer) error {
	err := pdc.Normal.Encode(buffer)
	if err != nil {
		return err
	}
	err = pdc.Operational.Encode(buffer)
	if err != nil {
		return err
	}
	return pdc.Mandatory.Encode(buffer)
}

func DecodePerDispatchClass[T sc.Encodable](buffer *bytes.Buffer, decodeFunc func(buffer *bytes.Buffer) (T, error)) (PerDispatchClass[T], error) {
	normal, err := decodeFunc(buffer)
	if err != nil {
		return PerDispatchClass[T]{}, err
	}
	operational, err := decodeFunc(buffer)
	if err != nil {
		return PerDispatchClass[T]{}, err
	}
	mandatory, err := decodeFunc(buffer)
	if err != nil {
		return PerDispatchClass[T]{}, err
	}
	return PerDispatchClass[T]{
		Normal:      normal,
		Operational: operational,
		Mandatory:   mandatory,
	}, nil
}

func (pdc PerDispatchClass[T]) Bytes() []byte {
	return sc.EncodedBytes(pdc)
}

// Get current value for given class.
func (pdc *PerDispatchClass[T]) Get(class DispatchClass) (*T, error) {
	switch class.VaryingData[0] {
	case DispatchClassNormal:
		return &pdc.Normal, nil
	case DispatchClassOperational:
		return &pdc.Operational, nil
	case DispatchClassMandatory:
		return &pdc.Mandatory, nil
	default:
		return nil, newTypeError("DispatchClass")
	}
}
