package types

import (
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

type DispatchResultWithPostInfo[T sc.Encodable] struct {
	HasError sc.Bool
	Ok       T
	Err      DispatchErrorWithPostInfo[T]
}

// ExtractActualWeight Extract the actual weight from a dispatch result if any or fall back to the default weight.
func ExtractActualWeight(result *DispatchResultWithPostInfo[PostDispatchInfo], info *DispatchInfo) Weight {
	var pdi PostDispatchInfo
	if result.HasError {
		err := result.Err
		pdi = err.PostInfo
	} else {
		pdi = result.Ok
	}
	return pdi.CalcActualWeight(info)
}

// ExtractActualPaysFee Extract the actual pays_fee from a dispatch result if any or fall back to the default weight.
func ExtractActualPaysFee(result *DispatchResultWithPostInfo[PostDispatchInfo], info *DispatchInfo) Pays {
	var pdi PostDispatchInfo
	if result.HasError {
		err := result.Err
		pdi = err.PostInfo
	} else {
		pdi = result.Ok
	}
	return pdi.Pays(info)
}
