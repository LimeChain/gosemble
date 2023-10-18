package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type InclusionFee struct {
	BaseFee           primitives.Balance
	LenFee            primitives.Balance
	AdjustedWeightFee primitives.Balance
}

func NewInclusionFee(baseFee, lenFee, adjustedWeightFee primitives.Balance) InclusionFee {
	return InclusionFee{
		baseFee,
		lenFee,
		adjustedWeightFee,
	}
}

func (i InclusionFee) Encode(buffer *bytes.Buffer) {
	i.BaseFee.Encode(buffer)
	i.LenFee.Encode(buffer)
	i.AdjustedWeightFee.Encode(buffer)
}

func (i InclusionFee) Bytes() []byte {
	return sc.EncodedBytes(i)
}

func DecodeInclusionFee(buffer *bytes.Buffer) InclusionFee {
	return InclusionFee{
		BaseFee:           sc.DecodeU128(buffer),
		LenFee:            sc.DecodeU128(buffer),
		AdjustedWeightFee: sc.DecodeU128(buffer),
	}
}

func (i InclusionFee) InclusionFee() primitives.Balance {
	return i.BaseFee.Add(i.LenFee).Add(i.AdjustedWeightFee)
}
