package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
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

func (pdc PerDispatchClass[T]) Encode(buffer *bytes.Buffer) {
	pdc.Normal.Encode(buffer)
	pdc.Operational.Encode(buffer)
	pdc.Mandatory.Encode(buffer)
}

func DecodePerDispatchClass[T sc.Encodable](buffer *bytes.Buffer, decodeFunc func(buffer *bytes.Buffer) T) PerDispatchClass[T] {
	return PerDispatchClass[T]{
		Normal:      decodeFunc(buffer),
		Operational: decodeFunc(buffer),
		Mandatory:   decodeFunc(buffer),
	}
}

func (pdc PerDispatchClass[T]) Bytes() []byte {
	return sc.EncodedBytes(pdc)
}

// Get current value for given class.
func (pdc *PerDispatchClass[T]) Get(class DispatchClass) *T {
	switch class.VaryingData[0] {
	case DispatchClassNormal:
		return &pdc.Normal
	case DispatchClassOperational:
		return &pdc.Operational
	case DispatchClassMandatory:
		return &pdc.Mandatory
	default:
		log.Critical("invalid DispatchClass type")
	}

	panic("unreachable")
}
