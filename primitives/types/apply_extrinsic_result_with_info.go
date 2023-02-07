package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// Same as `ApplyExtrinsicResult` but augmented with `PostDispatchInfo` on success.
type ApplyExtrinsicResultWithInfo sc.VaryingData

func NewApplyExtrinsicResultWithInfo(value sc.Encodable) ApplyExtrinsicResultWithInfo {
	// DispatchOutcome 					= 0 Outcome of dispatching the extrinsic.
	// TransactionValidityError = 1 Possible errors while checking the validity of a transaction.
	switch value.(type) {
	case DispatchOutcome, TransactionValidityError:
		return ApplyExtrinsicResultWithInfo(sc.NewVaryingData(value))
	default:
		panic("invalid ApplyExtrinsicResultWithInfo option")
	}
}

func (er ApplyExtrinsicResultWithInfo) Encode(buffer *bytes.Buffer) {
	switch er[0].(type) {
	case DispatchOutcome:
		sc.U8(0).Encode(buffer)
	case TransactionValidityError:
		sc.U8(1).Encode(buffer)
	default:
		panic("invalid ApplyExtrinsicResultWithInfo type")
	}

	er[0].Encode(buffer)
}

func DecodeApplyExtrinsicResultWithInfo(buffer *bytes.Buffer) ApplyExtrinsicResultWithInfo {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		value := DecodeDispatchOutcome(buffer)
		return NewApplyExtrinsicResultWithInfo(value)
	case 1:
		value := DecodeTransactionValidityError(buffer)
		return NewApplyExtrinsicResultWithInfo(value)
	default:
		panic("invalid ApplyExtrinsicResultWithInfo type")
	}
}

func (r ApplyExtrinsicResultWithInfo) Bytes() []byte {
	return sc.EncodedBytes(r)
}
