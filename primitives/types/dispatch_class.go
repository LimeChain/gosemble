package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	// A normal dispatch.
	NormalDispatch = DispatchClassValue(iota)

	// An operational dispatch.
	OperationalDispatch

	// A mandatory dispatch. These kinds of dispatch are always included regardless of their
	// weight, therefore it is critical that they are separately validated to ensure that a
	// malicious validator cannot craft a valid but impossibly heavy block. Usually this just
	// means ensuring that the extrinsic can only be included once and that it is always very
	// light.
	//
	// Do *NOT* use it for extrinsics that can be heavy.
	//
	// The only real use case for this is inherent extrinsics that are required to execute in a
	// block for the block to be valid, and it solves the issue in the case that the block
	// initialization is sufficiently heavy to mean that those inherents do not fit into the
	// block. Essentially, we assume that in these exceptional circumstances, it is better to
	// allow an overweight block to be created than to not allow any block at all to be created.
	MandatoryDispatch
)

// A generalized group of dispatch types.
type DispatchClassValue sc.U8

func (cl DispatchClassValue) Encode(buffer *bytes.Buffer) {
	sc.U8(cl).Encode(buffer)
}

func DecodeDispatchClassValue(buffer *bytes.Buffer) DispatchClassValue {
	return DispatchClassValue(sc.DecodeU8(buffer))
}

func (cl DispatchClassValue) Bytes() []byte {
	return sc.EncodedBytes(cl)
}

type DispatchClass sc.VaryingData

func NewDispatchClass(value DispatchClassValue) DispatchClass {
	return DispatchClass(sc.NewVaryingData(value))
}

func (dc DispatchClass) Encode(buffer *bytes.Buffer) {
	switch reflect.TypeOf(dc) {
	case reflect.TypeOf(NewDispatchClass(NormalDispatch)):
		sc.U8(0).Encode(buffer)
	case reflect.TypeOf(NewDispatchClass(OperationalDispatch)):
		sc.U8(1).Encode(buffer)
	case reflect.TypeOf(NewDispatchClass(MandatoryDispatch)):
		sc.U8(2).Encode(buffer)
	default:
		log.Critical("invalid DispatchClass type")
	}
}

func DecodeDispatchClass(buffer *bytes.Buffer) DispatchClass {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		return NewDispatchClass(NormalDispatch)
	case 1:
		return NewDispatchClass(OperationalDispatch)
	case 2:
		return NewDispatchClass(MandatoryDispatch)
	default:
		log.Critical("invalid DispatchClass type")
	}

	panic("unreachable")
}

func (dc DispatchClass) Bytes() []byte {
	return sc.EncodedBytes(dc)
}
