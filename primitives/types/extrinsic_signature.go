package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// The signature is a varying data type indicating the used signature type,
// followed by the signature created by the extrinsic author (the sender).
type ExtrinsicSignature struct {
	// sc.U64

	// is the 32-byte address of the sender of the extrinsic
	// as described in https://docs.substrate.io/reference/address-formats/
	// AccountId AccountId // size 32
	Signer    MultiAddress
	Signature MultiSignature
	Extra     Extra
}

func (s ExtrinsicSignature) Encode(buffer *bytes.Buffer) {
	s.Signer.Encode(buffer)
	s.Signature.Encode(buffer) // panic(len(s.Signature))
	s.Extra.Encode(buffer)
}

func DecodeExtrinsicSignature(buffer *bytes.Buffer) ExtrinsicSignature {
	s := ExtrinsicSignature{}
	s.Signer = DecodeMultiAddress(buffer)
	s.Signature = DecodeMultiSignature(buffer)
	s.Extra = DecodeExtra(buffer)
	return s
}

func (s ExtrinsicSignature) Bytes() []byte {
	return sc.EncodedBytes(s)
}
