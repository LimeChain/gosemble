/*
Implementation of an unchecked (pre-verification) extrinsic.
*/
package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// A extrinsic right from the external world. This is unchecked and so can contain a signature.
type UncheckedExtrinsic struct {
	// TODO:
	// make it generic
	// UncheckedExtrinsic[Address, Call, Signature, Extra] where Extra: SignedExtension
	Version sc.U8

	// The signature, address, number of extrinsics have come before from
	// the same signer and an era describing the longevity of this transaction,
	// if this is a signed extrinsic.
	Signature sc.Option[ExtrinsicSignature]
	Function  Call
}

func NewUncheckedExtrinsic(function Call, signedData sc.Option[ExtrinsicSignature]) UncheckedExtrinsic {
	if signedData.HasValue {
		address, signature, extra := signedData.Value.Signer, signedData.Value.Signature, signedData.Value.Extra
		return NewSignedUncheckedExtrinsic(function, address, signature, extra)
	} else {
		return NewUnsignedUncheckedExtrinsic(function)
	}
}

// New instance of a signed extrinsic aka "transaction".
func NewSignedUncheckedExtrinsic(function Call, address MultiAddress, signature MultiSignature, extra Extra) UncheckedExtrinsic {
	return UncheckedExtrinsic{
		Version: sc.U8(ExtrinsicFormatVersion | ExtrinsicBitSigned),
		Signature: sc.NewOption[ExtrinsicSignature](
			ExtrinsicSignature{
				Signer:    address,
				Signature: signature,
				Extra:     extra,
			},
		),
		Function: function,
	}
}

// New instance of an unsigned extrinsic aka "inherent".
func NewUnsignedUncheckedExtrinsic(function Call) UncheckedExtrinsic {
	return UncheckedExtrinsic{
		Version:   sc.U8(ExtrinsicFormatVersion),
		Signature: sc.NewOption[ExtrinsicSignature](nil),
		Function:  function,
	}
}

func (uxt UncheckedExtrinsic) UnmaskedVersion() sc.U8 {
	return uxt.Version & ExtrinsicUnmaskVersion
}

func (uxt UncheckedExtrinsic) IsSigned() sc.Bool {
	return uxt.Version&ExtrinsicBitSigned == ExtrinsicBitSigned
}

func (uxt UncheckedExtrinsic) Encode(buffer *bytes.Buffer) {
	tempBuffer := &bytes.Buffer{}

	if uxt.Signature.HasValue {
		sc.U8(ExtrinsicFormatVersion | ExtrinsicBitSigned).Encode(tempBuffer)
		uxt.Signature.Encode(tempBuffer)
	} else {
		sc.U8(ExtrinsicFormatVersion & ExtrinsicUnmaskVersion).Encode(tempBuffer)
	}

	uxt.Function.Encode(tempBuffer)

	encodedLen := sc.ToCompact(uint64(tempBuffer.Len()))
	encodedLen.Encode(buffer)
	buffer.Write(tempBuffer.Bytes())
}

func DecodeUncheckedExtrinsic(buffer *bytes.Buffer) UncheckedExtrinsic {
	// This is a little more complicated than usual since the binary format must be compatible
	// with SCALE's generic `Vec<u8>` type. Basically this just means accepting that there
	// will be a prefix of vector length.
	expectedLength := sc.DecodeCompact(buffer)
	_ = expectedLength
	// beforeLength := buffer.Len()

	version, _ := buffer.ReadByte()
	isSigned := version&ExtrinsicBitSigned != 0

	if version&ExtrinsicUnmaskVersion != ExtrinsicFormatVersion {
		panic("invalid Extrinsic version")
	}

	var extSignature sc.Option[ExtrinsicSignature]
	if isSigned {
		// extSignature = sc.DecodeOption[ExtrinsicSignature](buffer)
		if hasValue := sc.DecodeU8(buffer); hasValue == 1 {
			extSignature = sc.NewOption[ExtrinsicSignature](DecodeExtrinsicSignature(buffer))
		}
	}

	function := DecodeCall(buffer)

	// if Some((beforeLength, afterLength)) = buffer.remaining_len()?.and_then(|a| beforeLength.map(|b| (b, a)))
	// {
	// 	length = beforeLength.saturating_sub(after_length)

	// 	if length != expectedLength.0 as usize {
	// 		return error("invalid length prefix".into())
	// 	}
	// }

	return UncheckedExtrinsic{
		Version:   sc.U8(version),
		Signature: extSignature,
		Function:  function,
	}
}

func (uxt UncheckedExtrinsic) Bytes() []byte {
	return sc.EncodedBytes(uxt)
}

func (uxt UncheckedExtrinsic) Check(lookup AccountIdLookup) (xt CheckedExtrinsic, err TransactionValidityError) {
	switch uxt.Signature.HasValue {
	case true:
		signed, signature, extra := uxt.Signature.Value.Signer, uxt.Signature.Value.Signature, uxt.Signature.Value.Extra

		signedAddress, err := lookup.Lookup(signed)
		if err != nil {
			return xt, err
		}

		rawPayload, err := NewSignedPayload(uxt.Function, extra)
		if err != nil {
			err = NewTransactionValidityError(NewUnknownTransaction(err))
			return xt, err
		}

		if !signature.Verify(rawPayload.UsingEncoded(), signedAddress) {
			err := NewTransactionValidityError(NewInvalidTransaction(BadProofError))
			return xt, err
		}

		function, extra, _ := rawPayload.Call, rawPayload.Extra, rawPayload.AdditionalSigned

		xt = CheckedExtrinsic{
			Signed:   sc.NewOption[AccountIdExtra](AccountIdExtra{Address32: signedAddress, Extra: extra}),
			Function: function,
		}
	case false:
		xt = CheckedExtrinsic{
			Signed:   sc.NewOption[AccountIdExtra](nil),
			Function: uxt.Function,
		}
	}

	return xt, err
}
