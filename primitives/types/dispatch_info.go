package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// DispatchInfo A bundle of static information collected from the `#[pallet::weight]` attributes.
type DispatchInfo struct {
	// Weight of this transaction.
	Weight Weight

	// Class of this transaction.
	Class DispatchClass

	// Does this transaction pay fees.
	PaysFee Pays
}

func (di DispatchInfo) Encode(buffer *bytes.Buffer) {
	di.Weight.Encode(buffer)
	di.Class.Encode(buffer)
	di.PaysFee.Encode(buffer)
}

func DecodeDispatchInfo(buffer *bytes.Buffer) DispatchInfo {
	di := DispatchInfo{}
	di.Weight = DecodeWeight(buffer)
	di.Class = DecodeDispatchClass(buffer)
	di.PaysFee = DecodePays(buffer)
	return di
}

func (di DispatchInfo) Bytes() []byte {
	return sc.EncodedBytes(di)
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

// GetDispatchInfo returns the DispatchInfo of the given call.
// Uses call's BaseWeight to calculate all the information in DispatchInfo
func GetDispatchInfo(call Call) DispatchInfo {
	baseWeight := call.BaseWeight()

	return DispatchInfo{
		Weight:  call.WeighData(baseWeight),
		Class:   call.ClassifyDispatch(baseWeight),
		PaysFee: call.PaysFee(baseWeight),
	}
}
