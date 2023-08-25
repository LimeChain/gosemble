package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type FeeDetails struct {
	InclusionFee sc.Option[InclusionFee]

	Tip Balance // not serializable
}

func (fd FeeDetails) Encode(buffer *bytes.Buffer) {
	fd.InclusionFee.Encode(buffer)
}

func (fd FeeDetails) Bytes() []byte {
	return sc.EncodedBytes(fd)
}

func DecodeFeeDetails(buffer *bytes.Buffer) FeeDetails {
	return FeeDetails{
		InclusionFee: sc.DecodeOptionWith(buffer, DecodeInclusionFee),
	}
}

func (fd FeeDetails) FinalFee() Balance {
	sum := fd.Tip

	if fd.InclusionFee.HasValue {
		inclusionFee := fd.InclusionFee.Value.InclusionFee()
		sum = inclusionFee.Add(fd.Tip).(sc.U128)
	}

	return sum
}
