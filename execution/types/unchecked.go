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

type UncheckedExtrinsic struct {
	Version sc.U8

	// The signature, address, number of extrinsics have come before from
	// the same signer and an era describing the longevity of this transaction,
	// if this is a signed extrinsic.
	Signature sc.Option[primitives.ExtrinsicSignature]
	Function  primitives.Call
	Extra     primitives.SignedExtra
}

// NewSignedUncheckedExtrinsic returns a new instance of a signed extrinsic.
func NewUncheckedExtrinsic(version sc.U8, signature sc.Option[primitives.ExtrinsicSignature], function primitives.Call, extra primitives.SignedExtra) UncheckedExtrinsic {
	return UncheckedExtrinsic{
		Version:   version,
		Signature: signature,
		Function:  function,
		Extra:     extra,
	}
}

// NewUnsignedUncheckedExtrinsic returns a new instance of an unsigned extrinsic.
func NewUnsignedUncheckedExtrinsic(function primitives.Call) UncheckedExtrinsic {
	return UncheckedExtrinsic{
		Version:   sc.U8(ExtrinsicFormatVersion),
		Signature: sc.NewOption[primitives.ExtrinsicSignature](nil),
		Function:  function,
	}
}

func (uxt UncheckedExtrinsic) UnmaskedVersion() sc.U8 {
	return uxt.Version & ExtrinsicUnmaskVersion
}

func (uxt UncheckedExtrinsic) IsSigned() sc.Bool {
	return uxt.Signature.HasValue
}

func (uxt UncheckedExtrinsic) Encode(buffer *bytes.Buffer) {
	tempBuffer := &bytes.Buffer{}

	if uxt.Signature.HasValue {
		sc.U8(ExtrinsicFormatVersion | ExtrinsicBitSigned).Encode(tempBuffer)
		uxt.Signature.Value.Encode(tempBuffer)
	} else {
		sc.U8(ExtrinsicFormatVersion & ExtrinsicUnmaskVersion).Encode(tempBuffer)
	}

	uxt.Function.Encode(tempBuffer)

	encodedLen := sc.ToCompact(uint64(tempBuffer.Len()))
	encodedLen.Encode(buffer)
	buffer.Write(tempBuffer.Bytes())
}

func (uxt UncheckedExtrinsic) Bytes() []byte {
	return sc.EncodedBytes(uxt)
}
