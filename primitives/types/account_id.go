package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// AccountId It's an account ID (pubkey).
type AccountId struct {
	Ed25519Signer
	Sr25519Signer
	EcdsaSigner
	//Address32 // TODO: Varies depending on Signature (32 for ed25519 and sr25519, 33 for ecdsa)
}

func (a AccountId) Encode(buffer *bytes.Buffer) {
	a.Ed25519Signer.Encode(buffer)
	a.Sr25519Signer.Encode(buffer)
	a.EcdsaSigner.Encode(buffer)
}

func (a AccountId) Bytes() []byte {
	return sc.EncodedBytes(a)
}

func DecodeAccountId(buffer *bytes.Buffer) (AccountId, error) {
	addr32, err := DecodeEd25519(buffer)
	if err != nil {
		return AccountId{}, err
	}
	return AccountId{Ed25519Signer: addr32}, nil // TODO: length 32 or 33 depending on algorithm
}
