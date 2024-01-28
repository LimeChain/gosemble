package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type PerDispatchClassU32 struct {
	// Value for `Normal` extrinsics.
	Normal sc.U32
	// Value for `Operational` extrinsics.
	Operational sc.U32
	// Value for `Mandatory` extrinsics.
	Mandatory sc.U32
}

func (pdc PerDispatchClassU32) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		pdc.Normal,
		pdc.Operational,
		pdc.Mandatory,
	)
}

func DecodePerDispatchClassU32(buffer *bytes.Buffer, decodeU32 func(buffer *bytes.Buffer) (sc.U32, error)) (PerDispatchClassU32, error) {
	normal, err := decodeU32(buffer)
	if err != nil {
		return PerDispatchClassU32{}, err
	}
	operational, err := decodeU32(buffer)
	if err != nil {
		return PerDispatchClassU32{}, err
	}
	mandatory, err := decodeU32(buffer)
	if err != nil {
		return PerDispatchClassU32{}, err
	}
	return PerDispatchClassU32{
		Normal:      normal,
		Operational: operational,
		Mandatory:   mandatory,
	}, nil
}

func (pdc PerDispatchClassU32) Bytes() []byte {
	return sc.EncodedBytes(pdc)
}

// Get current value for given class.
func (pdc *PerDispatchClassU32) Get(class DispatchClass) (*sc.U32, error) {
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
