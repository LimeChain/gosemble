package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// PostDispatchInfo Weight information that is only available post dispatch.
// NOTE: This can only be used to reduce the weight or fee, not increase it.
type PostDispatchInfo struct {
	// Actual weight consumed by a call or `None` which stands for the worst case static weight.
	ActualWeight sc.Option[Weight]

	// Whether this transaction should pay fees when all is said and done.
	PaysFee Pays
}

func (pdi PostDispatchInfo) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		pdi.ActualWeight,
		pdi.PaysFee,
	)
}

func DecodePostDispatchInfo(buffer *bytes.Buffer) (PostDispatchInfo, error) {
	actualWeight, err := sc.DecodeOptionWith(buffer, DecodeWeight)
	if err != nil {
		return PostDispatchInfo{}, err
	}
	paysFee, err := DecodePays(buffer)
	if err != nil {
		return PostDispatchInfo{}, err
	}
	return PostDispatchInfo{
		ActualWeight: actualWeight,
		PaysFee:      paysFee,
	}, nil
}

func (pdi PostDispatchInfo) Bytes() []byte {
	return sc.EncodedBytes(pdi)
}

// CalcUnspent Calculate how much (if any) weight was not used by the `Dispatchable`.
func (pdi PostDispatchInfo) CalcUnspent(info *DispatchInfo) Weight {
	return info.Weight.Sub(pdi.CalcActualWeight(info))
}

// CalcActualWeight Calculate how much weight was actually spent by the `Dispatchable`.
func (pdi PostDispatchInfo) CalcActualWeight(info *DispatchInfo) Weight {
	if pdi.ActualWeight.HasValue {
		actualWeight := pdi.ActualWeight.Value
		return actualWeight.Min(info.Weight)
	} else {
		return info.Weight
	}
}

// Pays Determine if user should actually pay fees at the end of the dispatch.
func (pdi PostDispatchInfo) Pays(info *DispatchInfo) Pays {
	// If they originally were not paying fees, or the post dispatch info
	// says they should not pay fees, then they don't pay fees.
	// This is because the pre dispatch information must contain the
	// worst case for weight and fees paid.

	if info.PaysFee == PaysNo || pdi.PaysFee == PaysNo {
		return PaysNo
	} else {
		// Otherwise they pay.
		return PaysYes
	}
}
