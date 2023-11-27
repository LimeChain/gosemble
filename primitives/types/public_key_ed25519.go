package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	publicKeyEd25519Length = 32
)

type Ed25519PublicKey struct {
	sc.FixedSequence[sc.U8] // size 32
}

func NewEd25519PublicKey(values ...sc.U8) (Ed25519PublicKey, error) {
	if len(values) != publicKeyEd25519Length {
		return Ed25519PublicKey{}, newTypeError("Ed25519PublicKey")
	}
	return Ed25519PublicKey{sc.NewFixedSequence(publicKeyEd25519Length, values...)}, nil
}

func (s Ed25519PublicKey) SignatureType() sc.U8 {
	return PublicKeyEd25519
}

func (s Ed25519PublicKey) Encode(buffer *bytes.Buffer) error {
	return s.FixedSequence.Encode(buffer)
}

func (s Ed25519PublicKey) Bytes() []byte {
	return sc.EncodedBytes(s)
}

func DecodeEd25519PublicKey(buffer *bytes.Buffer) (Ed25519PublicKey, error) {
	seq, err := sc.DecodeFixedSequence[sc.U8](publicKeyEd25519Length, buffer)
	if err != nil {
		return Ed25519PublicKey{}, err
	}
	return Ed25519PublicKey{seq}, nil
}
