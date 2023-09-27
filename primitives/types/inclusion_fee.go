package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type InclusionFee struct {
	BaseFee           Balance
	LenFee            Balance
	AdjustedWeightFee Balance
}

func NewInclusionFee(baseFee, lenFee, adjustedWeightFee Balance) InclusionFee {
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

func (i InclusionFee) InclusionFee() Balance {
	return i.BaseFee.Add(i.LenFee).Add(i.AdjustedWeightFee)
}
