package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// ApplyExtrinsicResult The result of applying of an extrinsic.
//
// This type is typically used in the context of `BlockBuilder` to signal that the extrinsic
// in question cannot be included.
//
// A block containing extrinsics that have a negative inclusion outcome is invalid. A negative
// result can only occur during the block production, where such extrinsics are detected and
// removed from the block that is being created and the transaction pool.
//
// To rehash: every extrinsic in a valid block must return a positive `ApplyExtrinsicResult`.
//
// Examples of reasons preventing inclusion in a block:
//   - More block weight is required to process the extrinsic than is left in the block being built.
//     This doesn't necessarily mean that the extrinsic is invalid, since it can still be included in
//     the next block if it has enough spare weight available.
//   - The sender doesn't have enough funds to pay the transaction inclusion fee. Including such a
//     transaction in the block doesn't make sense.
//   - The extrinsic supplied a bad signature. This transaction won't become valid ever.
type ApplyExtrinsicResult sc.VaryingData // = sc.Result[DispatchOutcome, TransactionValidityError]

func NewApplyExtrinsicResult(value sc.Encodable) (ApplyExtrinsicResult, error) {
	// DispatchOutcome 					= 0 Outcome of dispatching the extrinsic.
	// TransactionValidityError = 1 Possible errors while checking the validity of a transaction.
	switch value.(type) {
	case DispatchOutcome, TransactionValidityError:
		return ApplyExtrinsicResult(sc.NewVaryingData(value)), nil
	default:
		return nil, newTypeError("ApplyExtrinsicResult")
	}
}

func (r ApplyExtrinsicResult) Encode(buffer *bytes.Buffer) error {
	switch r[0].(type) {
	case DispatchOutcome:
		err := sc.U8(0).Encode(buffer)
		if err != nil {
			return err
		}
	case TransactionValidityError:
		err := sc.U8(1).Encode(buffer)
		if err != nil {
			return err
		}
	default:
		return newTypeError("ApplyExtrinsicResult")
	}

	return r[0].Encode(buffer)
}

func DecodeApplyExtrinsicResult(buffer *bytes.Buffer) (ApplyExtrinsicResult, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return nil, err
	}

	switch b {
	case 0:
		value, err := DecodeDispatchOutcome(buffer)
		if err != nil {
			return nil, err
		}
		return NewApplyExtrinsicResult(value)
	case 1:
		err := DecodeTransactionValidityError(buffer)
		if txErr, ok := err.(TransactionValidityError); ok {
			return NewApplyExtrinsicResult(txErr)
		}
		return nil, err
	default:
		return nil, newTypeError("ApplyExtrinsicResult")
	}
}

func (r ApplyExtrinsicResult) Bytes() []byte {
	return sc.EncodedBytes(r)
}
