package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	// A normal dispatch.
	NormalDispatch = sc.U8(iota)

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

const (
	// Transactor will pay related fees.
	PaysYes = iota

	// Transactor will NOT pay related fees.
	PaysNo
)

type DispatchResult sc.VaryingData

func NewDispatchResult(value sc.Encodable) DispatchResult {
	switch value.(type) {
	case DispatchError:
		return DispatchResult(sc.NewVaryingData(value))
	case sc.Empty, nil:
		return DispatchResult(sc.NewVaryingData(sc.Empty{}))
	default:
		panic("invalid DispatchResult option")
	}
}

func (r DispatchResult) Encode(buffer *bytes.Buffer) {
	value := r[0]

	switch value {
	case NormalDispatch:
		sc.U8(0).Encode(buffer)
	case OperationalDispatch:
		sc.U8(1).Encode(buffer)
	case MandatoryDispatch:
		sc.U8(2).Encode(buffer)
	default:
		panic("invalid DispatchResult type")
	}
}

func DecodeDispatchResult(buffer *bytes.Buffer) DispatchResult {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		return NewDispatchResult(NormalDispatch)
	case 1:
		return NewDispatchResult(OperationalDispatch)
	case 2:
		return NewDispatchResult(MandatoryDispatch)
	default:
		panic("invalid DispatchResult type")
	}
}

func (r DispatchResult) Bytes() []byte {
	return sc.EncodedBytes(r)
}

func (r DispatchResult) PostDispatch() (ok Pre, err TransactionValidityError) {
	// TODO:
	ok = Pre{}
	return ok, err
}

// Result of a `Dispatchable` which contains the `DispatchResult` and additional information about
// the `Dispatchable` that is only known post dispatch.
// type DispatchErrorWithPostInfo = TransactionValidityError
type DispatchResultWithInfo sc.VaryingData

func NewDispatchResultWithInfo(value sc.Encodable) DispatchError {
	switch value.(type) {
	case UnknownError, DataLookupError, BadOriginError, CustomModuleError:
		return DispatchError(sc.NewVaryingData(value))
	default:
		panic("invalid DispatchError type")
	}
}

type PostInfo struct{}
