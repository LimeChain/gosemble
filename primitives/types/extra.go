package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// Extra data, E, is a tuple containing additional meta data about the extrinsic and the system it is meant to be executed in.
// E := (Tmor, N, Pt)
type Extra struct {
	// Tmor: contains the SCALE encoded mortality of the extrinsic
	// Mortality sc.Sequence[sc.U8]
	Era ExtrinsicEra

	// N: a compact integer containing the nonce of the sender.
	// The nonce must be incremented by one for each extrinsic created,
	// otherwise the Polkadot network will reject the extrinsic.
	Nonce sc.Compact // sc.U64

	// Pt: a compact integer containing the transactor pay including tip.
	Fee sc.Compact // sc.U64
}

// type SignedExtra struct {
// 	NonZeroSender MultiAddress
// 	SpecVersion   sc.U32
// 	TxVersion     sc.U32
// 	Genesis       H256
// 	Era           ExtrinsicEra
// 	Nonce         sc.Compact
// 	// Weight
// 	TransactionPayment sc.Compact
// }

func (e Extra) Encode(buffer *bytes.Buffer) {
	e.Era.Encode(buffer)
	e.Nonce.Encode(buffer)
	e.Fee.Encode(buffer)
}

func DecodeExtra(buffer *bytes.Buffer) Extra {
	e := Extra{}
	e.Era = DecodeExtrinsicEra(buffer)
	e.Nonce = sc.DecodeCompact(buffer)
	e.Fee = sc.DecodeCompact(buffer)
	return e
}

func (e Extra) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e Extra) AdditionalSigned() (AdditionalSigned, TransactionValidityError) {
	return AdditionalSigned{
		SpecVersion:   sc.U32(RuntimeVersion{}.SpecVersion),
		FormatVersion: ExtrinsicFormatVersion,
		// GenesisHash:   H256(),
		// BlockHash: H256(),
		// TransactionVersion sc.U32
		// BlockNumber
	}, nil
}

// Information on a transaction's validity and, if valid, on how it relates to other transactions.
func (_ Extra) Validate(who *Address32, call *Call, info *DispatchInfo, length sc.Compact) (ok ValidTransaction, err TransactionValidityError) {
	valid := DefaultValidTransaction()

	ok, err = who.Validate()
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	ok, err = call.Validate()
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	ok, err = info.Validate()
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	ok, err = Length(length).Validate()
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	return valid, err
}

func (_ Extra) ValidateUnsigned(call *Call, info *DispatchInfo, length sc.Compact) (ok ValidTransaction, err TransactionValidityError) {
	valid := DefaultValidTransaction()

	ok, err = call.ValidateUnsigned()
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	ok, err = info.ValidateUnsigned()
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	ok, err = Length(length).ValidateUnsigned()
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	return valid, err
}

// Do any pre-flight stuff for a signed transaction.
//
// Make sure to perform the same checks as in [`Validate`].
func (_ Extra) PreDispatch(e Extra, who *Address32, call *Call, info *DispatchInfo, length sc.Compact) (ok Pre, err TransactionValidityError) {
	ok, err = who.PreDispatch()
	if err != nil {
		return ok, err
	}

	ok, err = call.PreDispatch()
	if err != nil {
		return ok, err
	}

	ok, err = info.PreDispatch()
	if err != nil {
		return ok, err
	}

	ok, err = Length(length).PreDispatch()
	if err != nil {
		return ok, err
	}

	return ok, err
}

func (_ Extra) PreDispatchUnsigned(call *Call, info *DispatchInfo, length sc.Compact) (ok Pre, err TransactionValidityError) {
	// Extra{}.ValidateUnsigned(call, info, length)

	ok, err = call.PreDispatchUnsigned()
	if err != nil {
		return ok, err
	}

	ok, err = info.PreDispatchUnsigned()
	if err != nil {
		return ok, err
	}

	ok, err = Length(length).PreDispatchUnsigned()
	if err != nil {
		return ok, err
	}

	return ok, err
}

func (_ Extra) PostDispatch(pre sc.Option[Pre], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) (ok Pre, err TransactionValidityError) {
	switch pre.HasValue {
	case true:
		// ok, err = pre.Value.PostDispatch()
		// if err != nil {
		// 	return ok, err
		// }

		ok, err = info.PostDispatch()
		if err != nil {
			return ok, err
		}

		ok, err = postInfo.PostDispatch()
		if err != nil {
			return ok, err
		}

		ok, err = Length(length).PostDispatch()
		if err != nil {
			return ok, err
		}

		ok, err = result.PostDispatch()
		if err != nil {
			return ok, err
		}

	case false:
		// sc.Empty
		// info
		// postInfo
		// length
		// result
	}

	return ok, err
}
