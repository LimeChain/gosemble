package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	// A normal dispatch.
	DispatchClassNormal sc.U8 = iota

	// An operational dispatch.
	DispatchClassOperational

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
	DispatchClassMandatory
)

// A generalized group of dispatch types.
type DispatchClass struct {
	sc.VaryingData
}

func NewDispatchClassNormal() DispatchClass {
	return DispatchClass{sc.NewVaryingData(DispatchClassNormal)}
}

func NewDispatchClassOperational() DispatchClass {
	return DispatchClass{sc.NewVaryingData(DispatchClassOperational)}
}

func NewDispatchClassMandatory() DispatchClass {
	return DispatchClass{sc.NewVaryingData(DispatchClassMandatory)}
}

func DecodeDispatchClass(buffer *bytes.Buffer) DispatchClass {
	b := sc.DecodeU8(buffer)

	switch b {
	case DispatchClassNormal:
		return NewDispatchClassNormal()
	case DispatchClassOperational:
		return NewDispatchClassOperational()
	case DispatchClassMandatory:
		return NewDispatchClassMandatory()
	default:
		log.Critical("invalid DispatchClass type")
	}

	panic("unreachable")
}

func (dc DispatchClass) Is(value sc.U8) sc.Bool {
	// TODO: type safety
	switch value {
	case DispatchClassNormal, DispatchClassOperational, DispatchClassMandatory:
		return dc.VaryingData[0] == value
	default:
		log.Critical("invalid DispatchClass value")
	}

	panic("unreachable")
}

// A struct holding value for each `DispatchClass`.
type PerDispatchClass[T sc.Encodable] struct {
	// Value for `Normal` extrinsics.
	Normal T

	// Value for `Operational` extrinsics.
	Operational T

	// Value for `Mandatory` extrinsics.
	Mandatory T
}
