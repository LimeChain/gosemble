package types

import (
	sc "github.com/LimeChain/goscale"
)

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
