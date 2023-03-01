/*
Implementation of an unchecked (pre-verification) extrinsic.
*/
package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

// A extrinsic right from the external world. This is unchecked and so can contain a signature.
//
// TODO: make it generic
// generic::UncheckedExtrinsic<Address, RuntimeCall, Signature, SignedExtra>;
type UncheckedExtrinsic struct {
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
func NewSignedUncheckedExtrinsic(function Call, address MultiAddress, signature MultiSignature, extra SignedExtra) UncheckedExtrinsic {
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
		log.Critical("invalid Extrinsic version")
	}

	var extSignature sc.Option[ExtrinsicSignature]
	if isSigned {
		extSignature = sc.DecodeOptionWith(buffer, DecodeExtrinsicSignature)
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
