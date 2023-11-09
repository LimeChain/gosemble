package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Ed25519PublicKey struct {
	sc.FixedSequence[sc.U8] // size 32
}

func (s Ed25519PublicKey) SignatureType() sc.U8 {
	return PublicKeyEd25519
}

func (s Ed25519PublicKey) Encode(buffer *bytes.Buffer) {
	s.FixedSequence.Encode(buffer)
}

func (s Ed25519PublicKey) Bytes() []byte {
	return sc.EncodedBytes(s)
}

func DecodeEd25519PublicKey(buffer *bytes.Buffer) (Ed25519PublicKey, error) {
	seq, err := sc.DecodeFixedSequence[sc.U8](32, buffer)
	if err != nil {
		return Ed25519PublicKey{}, err
	}
	return Ed25519PublicKey{seq}, nil
}

func NewEd25519PublicKey(values ...sc.U8) (Ed25519PublicKey, error) {
	if len(values) != 32 {
		return Ed25519PublicKey{}, newTypeError("Ed25519PublicKey")
	}
	return Ed25519PublicKey{sc.NewFixedSequence(32, values...)}, nil
}
