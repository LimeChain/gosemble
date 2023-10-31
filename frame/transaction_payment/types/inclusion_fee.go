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

func DecodeInclusionFee(buffer *bytes.Buffer) (InclusionFee, error) {
	baseFee, err := sc.DecodeU128(buffer)
	if err != nil {
		return InclusionFee{}, err
	}
	lenFee, err := sc.DecodeU128(buffer)
	if err != nil {
		return InclusionFee{}, err
	}
	adjustedWeightFee, err := sc.DecodeU128(buffer)
	if err != nil {
		return InclusionFee{}, err
	}
	return InclusionFee{
		BaseFee:           baseFee,
		LenFee:            lenFee,
		AdjustedWeightFee: adjustedWeightFee,
	}, nil
}

func (i InclusionFee) InclusionFee() primitives.Balance {
	return i.BaseFee.Add(i.LenFee).Add(i.AdjustedWeightFee)
}
