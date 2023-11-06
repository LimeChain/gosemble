package types

import (
	sc "github.com/LimeChain/goscale"
)

type DispatchResult sc.VaryingData

func NewDispatchResult(value sc.Encodable) (DispatchResult, error) {
	switch value.(type) {
	case DispatchError, DispatchErrorWithPostInfo[PostDispatchInfo]:
		return DispatchResult(sc.NewVaryingData(value)), nil
	case sc.Empty, nil:
		return DispatchResult(sc.NewVaryingData(sc.Empty{})), nil
	default:
		return DispatchResult{}, NewTypeError("DispatchResult")
	}
}
