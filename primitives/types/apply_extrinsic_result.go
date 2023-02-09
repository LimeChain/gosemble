package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type ApplyExtrinsicResult sc.VaryingData

func NewApplyExtrinsicResult(value sc.Encodable) ApplyExtrinsicResult {
	// DispatchOutcome 					= 0 Outcome of dispatching the extrinsic.
	// TransactionValidityError = 1 Possible errors while checking the validity of a transaction.
	switch value.(type) {
	case DispatchOutcome, TransactionValidityError:
		return ApplyExtrinsicResult(sc.NewVaryingData(value))
	default:
		panic("invalid ApplyExtrinsicResult option")
	}
}

func (r ApplyExtrinsicResult) Encode(buffer *bytes.Buffer) {
	switch r[0].(type) {
	case DispatchOutcome:
		sc.U8(0).Encode(buffer)
	case TransactionValidityError:
		sc.U8(1).Encode(buffer)
	default:
		panic("invalid ApplyExtrinsicResult type")
	}

	r[0].Encode(buffer)
}

func DecodeApplyExtrinsicResult(buffer *bytes.Buffer) ApplyExtrinsicResult {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		value := DecodeDispatchOutcome(buffer)
		return NewApplyExtrinsicResult(value)
	case 1:
		value := DecodeTransactionValidityError(buffer)
		return NewApplyExtrinsicResult(value)
	default:
		panic("invalid ApplyExtrinsicResult type")
	}
}

func (r ApplyExtrinsicResult) Bytes() []byte {
	return sc.EncodedBytes(r)
}
