package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type DispatchClassType sc.U8

const (
	// DispatchClassNormal A normal dispatch.
	DispatchClassNormal DispatchClassType = iota

	// DispatchClassOperational An operational dispatch.
	DispatchClassOperational

	// DispatchClassMandatory A mandatory dispatch. These kinds of dispatch are always included regardless of their
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

func (dct DispatchClassType) Encode(buffer *bytes.Buffer) error {
	return sc.U8(dct).Encode(buffer)
}

func (dct DispatchClassType) Bytes() []byte {
	return sc.EncodedBytes(dct)
}

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

func DecodeDispatchClass(buffer *bytes.Buffer) (DispatchClass, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return DispatchClass{}, err
	}

	switch DispatchClassType(b) {
	case DispatchClassNormal:
		return NewDispatchClassNormal(), nil
	case DispatchClassOperational:
		return NewDispatchClassOperational(), nil
	case DispatchClassMandatory:
		return NewDispatchClassMandatory(), nil
	default:
		return DispatchClass{}, newTypeError("DispatchClass")
	}
}

func (dc DispatchClass) Is(value DispatchClassType) (sc.Bool, error) {
	if len(dc.VaryingData) == 0 {
		return false, newTypeError("DispatchClass")
	}

	switch value {
	case DispatchClassNormal, DispatchClassOperational, DispatchClassMandatory:
		return dc.VaryingData[0] == value, nil
	default:
		return false, newTypeError("DispatchClass")
	}
}

// DispatchClassAll Returns a slice, containing all dispatch classes.
func DispatchClassAll() []DispatchClass {
	return []DispatchClass{NewDispatchClassNormal(), NewDispatchClassOperational(), NewDispatchClassMandatory()}
}
