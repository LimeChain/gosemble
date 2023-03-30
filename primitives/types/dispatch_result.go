package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type DispatchResult sc.VaryingData

func NewDispatchResult(value sc.Encodable) DispatchResult {
	switch value.(type) {
	case DispatchError, DispatchErrorWithPostInfo[PostDispatchInfo]:
		return DispatchResult(sc.NewVaryingData(value))
	case sc.Empty, nil:
		return DispatchResult(sc.NewVaryingData(sc.Empty{}))
	default:
		log.Critical("invalid DispatchResult type")
	}

	panic("unreachable")
}

func (r DispatchResult) Encode(buffer *bytes.Buffer) {
	// TODO:
}

func DecodeDispatchResult(buffer *bytes.Buffer) DispatchResult {
	// TODO:
	return DispatchResult{}
}

func (r DispatchResult) Bytes() []byte {
	return sc.EncodedBytes(r)
}

// Result of a `Dispatchable` which contains the `DispatchResult` and additional information about
// the `Dispatchable` that is only known post dispatch.
//
// type DispatchResultWithInfo sc.VaryingData

// func NewDispatchResultWithInfo(value sc.Encodable) DispatchResultWithInfo {
// 	switch value.(type) {
// 	case UnknownError, DataLookupError, BadOriginError, CustomModuleError:
// 		return DispatchResultWithInfo(sc.NewVaryingData(value))
// 	default:
// 		log.Critical("invalid DispatchResultWithInfo type")
// 	}

// }
type DispatchResultWithPostInfo[T sc.Encodable] struct {
	HasError sc.Bool
	Ok       T
	Err      DispatchErrorWithPostInfo[T]
}

func (r DispatchResultWithPostInfo[T]) Encode(buffer *bytes.Buffer) {
	r.HasError.Encode(buffer)

	if r.HasError {
		r.Err.Encode(buffer)
	} else {
		r.Ok.Encode(buffer)
	}
}

func DecodeDispatchResultWithPostInfo[T sc.Encodable](buffer *bytes.Buffer) DispatchResultWithPostInfo[T] {
	hasError := sc.DecodeBool(buffer)

	if hasError {
		// TODO: finish this
		return DispatchResultWithPostInfo[T]{}
	} else {
		return DispatchResultWithPostInfo[T]{}
	}
}

func (r DispatchResultWithPostInfo[T]) Bytes() []byte {
	return sc.EncodedBytes(r)
}
