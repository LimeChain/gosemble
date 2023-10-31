package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type FeeDetails struct {
	InclusionFee sc.Option[InclusionFee]
	Tip          primitives.Balance // not serializable
}

func (fd FeeDetails) Encode(buffer *bytes.Buffer) {
	fd.InclusionFee.Encode(buffer)
}

func (fd FeeDetails) Bytes() []byte {
	return sc.EncodedBytes(fd)
}

func DecodeFeeDetails(buffer *bytes.Buffer) (FeeDetails, error) {
	inclFee, err := sc.DecodeOptionWith(buffer, DecodeInclusionFee)
	if err != nil {
		return FeeDetails{}, err
	}
	return FeeDetails{
		InclusionFee: inclFee,
	}, nil
}

func (fd FeeDetails) FinalFee() primitives.Balance {
	sum := fd.Tip

	if fd.InclusionFee.HasValue {
		inclusionFee := fd.InclusionFee.Value.InclusionFee()
		sum = inclusionFee.Add(fd.Tip)
	}

	return sum
}
