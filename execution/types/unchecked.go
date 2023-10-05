package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	// ExtrinsicFormatVersion is the current version of the [`UncheckedExtrinsic`] encoded format.
	//
	// This version needs to be bumped if the encoded representation changes.
	// It ensures that if the representation is changed and the format is not known,
	// the decoding fails.
	ExtrinsicFormatVersion = 4
	ExtrinsicBitSigned     = 0b1000_0000
	ExtrinsicUnmaskVersion = 0b0111_1111
)

type UncheckedExtrinsic interface {
	sc.Encodable

	Signature() sc.Option[primitives.ExtrinsicSignature]
	Function() primitives.Call
	Extra() primitives.SignedExtra

	IsSigned() sc.Bool
	Check(lookup primitives.AccountIdLookup) (CheckedExtrinsic, primitives.TransactionValidityError)
}

type uncheckedExtrinsic struct {
	version sc.U8
	// The signature, address, number of extrinsics have come before from
	// the same signer and an era describing the longevity of this transaction,
	// if this is a signed extrinsic.
	signature sc.Option[primitives.ExtrinsicSignature]
	function  primitives.Call
	extra     primitives.SignedExtra
}

// NewSignedUncheckedExtrinsic returns a new instance of a signed extrinsic.
func NewUncheckedExtrinsic(version sc.U8, signature sc.Option[primitives.ExtrinsicSignature], function primitives.Call, extra primitives.SignedExtra) uncheckedExtrinsic {
	return uncheckedExtrinsic{
		version:   version,
		signature: signature,
		function:  function,
		extra:     extra,
	}
}

// NewUnsignedUncheckedExtrinsic returns a new instance of an unsigned extrinsic.
func NewUnsignedUncheckedExtrinsic(function primitives.Call) uncheckedExtrinsic {
	return uncheckedExtrinsic{
		version:   sc.U8(ExtrinsicFormatVersion),
		signature: sc.NewOption[primitives.ExtrinsicSignature](nil),
		function:  function,
	}
}

func (uxt uncheckedExtrinsic) IsSigned() sc.Bool {
	return uxt.signature.HasValue
}

func (uxt uncheckedExtrinsic) Encode(buffer *bytes.Buffer) {
	tempBuffer := &bytes.Buffer{}

	if uxt.Signature().HasValue {
		sc.U8(ExtrinsicFormatVersion | ExtrinsicBitSigned).Encode(tempBuffer)
		uxt.Signature().Value.Encode(tempBuffer)
	} else {
		sc.U8(ExtrinsicFormatVersion & ExtrinsicUnmaskVersion).Encode(tempBuffer)
	}

	uxt.Function().Encode(tempBuffer)

	encodedLen := sc.ToCompact(uint64(tempBuffer.Len()))
	encodedLen.Encode(buffer)
	buffer.Write(tempBuffer.Bytes())
}

func (uxt uncheckedExtrinsic) Bytes() []byte {
	return sc.EncodedBytes(uxt)
}

func (uxt uncheckedExtrinsic) Signature() sc.Option[primitives.ExtrinsicSignature] {
	return uxt.signature
}

func (uxt uncheckedExtrinsic) Function() primitives.Call {
	return uxt.function
}

func (uxt uncheckedExtrinsic) Extra() primitives.SignedExtra {
	return uxt.extra
}

func (uxt uncheckedExtrinsic) Check(lookup primitives.AccountIdLookup) (CheckedExtrinsic, primitives.TransactionValidityError) {
	if uxt.Signature().HasValue {
		signer, signature, extra := uxt.Signature().Value.Signer, uxt.Signature().Value.Signature, uxt.Signature().Value.Extra

		signedAddress, err := lookup.Lookup(signer)
		if err != nil {
			return checkedExtrinsic{}, err
		}

		rawPayload, err := NewSignedPayload(uxt.Function(), extra)
		if err != nil {
			return checkedExtrinsic{}, err
		}

		if !signature.Verify(rawPayload.UsingEncoded(), signedAddress) {
			err := primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadProof())
			return checkedExtrinsic{}, err
		}

		function, extra, _ := rawPayload.Call, rawPayload.Extra, rawPayload.AdditionalSigned

		return NewCheckedExtrinsic(sc.NewOption[primitives.Address32](signedAddress), function, extra), nil
	}

	return NewCheckedExtrinsic(sc.NewOption[primitives.Address32](nil), uxt.Function(), uxt.Extra()), nil
}
