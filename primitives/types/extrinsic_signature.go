package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// ExtrinsicSignature The signature is a varying data type indicating the used signature type,
// followed by the signature created by the extrinsic author (the sender).
type ExtrinsicSignature struct {
	// is the 32-byte address of the sender of the extrinsic
	// as described in https://docs.substrate.io/reference/address-formats/
	Signer    MultiAddress
	Signature MultiSignature
	Extra     SignedExtra
}

func (s ExtrinsicSignature) Encode(buffer *bytes.Buffer) {
	s.Signer.Encode(buffer)
	s.Signature.Encode(buffer)
	s.Extra.Encode(buffer)
}

func DecodeExtrinsicSignature(extra SignedExtra, buffer *bytes.Buffer) (ExtrinsicSignature, error) {
	s := ExtrinsicSignature{}
	signer, err := DecodeMultiAddress(buffer)
	if err != nil {
		return ExtrinsicSignature{}, err
	}
	s.Signer = signer
	signature, err := DecodeMultiSignature(buffer)
	if err != nil {
		return ExtrinsicSignature{}, err
	}
	s.Signature = signature

	s.Extra = extra
	s.Extra.Decode(buffer)

	return s, nil
}

func (s ExtrinsicSignature) Bytes() []byte {
	return sc.EncodedBytes(s)
}
