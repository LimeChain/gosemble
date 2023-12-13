package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
)

// DispatchOutcome This type specifies the outcome of dispatching a call to a module.
//
// In case of failure an error specific to the module is returned.
//
// Failure of the module call dispatching doesn't invalidate the extrinsic and it is still included
// in the block, therefore all state changes performed by the dispatched call are still persisted.
//
// For example, if the dispatching of an extrinsic involves inclusion fee payment then these
// changes are going to be preserved even if the call dispatched failed.
type DispatchOutcome sc.VaryingData //  = sc.Result[sc.Empty, DispatchError]

func NewDispatchOutcome(value sc.Encodable) (DispatchOutcome, error) {
	// None 			   = 0 - Extrinsic is valid and was submitted successfully.
	// DispatchError = 1 - Possible errors while dispatching the extrinsic.
	switch value.(type) {
	case DispatchError:
		return DispatchOutcome(sc.NewVaryingData(value)), nil
	case sc.Empty, nil:
		return DispatchOutcome(sc.NewVaryingData(sc.Empty{})), nil
	default:
		return DispatchOutcome{}, newTypeError("DispatchOutcome")
	}
}

func (o DispatchOutcome) Encode(buffer *bytes.Buffer) error {
	value := o[0]

	switch reflect.TypeOf(value) {
	case reflect.TypeOf(*new(sc.Empty)):
		return sc.U8(0).Encode(buffer)
	case reflect.TypeOf(*new(DispatchError)):
		return sc.EncodeEach(buffer, sc.U8(1), value)
	default:
		return newTypeError("DispatchOutcome")
	}
}

func DecodeDispatchOutcome(buffer *bytes.Buffer) (DispatchOutcome, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return DispatchOutcome{}, err
	}

	switch b {
	case 0:
		return NewDispatchOutcome(sc.Empty{})
	case 1:
		value, err := DecodeDispatchError(buffer)
		if err != nil {
			return DispatchOutcome{}, err
		}
		return NewDispatchOutcome(value)
	default:
		return DispatchOutcome{}, newTypeError("DispatchOutcome")
	}
}

func (o DispatchOutcome) Bytes() []byte {
	return sc.EncodedBytes(o)
}
