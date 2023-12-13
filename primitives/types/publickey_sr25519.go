package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Sr25519PublicKey struct {
	sc.FixedSequence[sc.U8] // size 32
}

func (s Sr25519PublicKey) SignatureType() sc.U8 {
	return PublicKeySr25519
}

func (s Sr25519PublicKey) Encode(buffer *bytes.Buffer) error {
	return s.FixedSequence.Encode(buffer)
}

func (s Sr25519PublicKey) Bytes() []byte {
	return sc.EncodedBytes(s)
}

func DecodeSr25519PublicKey(buffer *bytes.Buffer) (Sr25519PublicKey, error) {
	seq, err := sc.DecodeFixedSequence[sc.U8](32, buffer)
	if err != nil {
		return Sr25519PublicKey{}, err
	}
	return Sr25519PublicKey{seq}, nil
}

func NewSr25519PublicKey(values ...sc.U8) (Sr25519PublicKey, error) {
	if len(values) != 32 {
		return Sr25519PublicKey{}, newTypeError("Sr25519PublicKey")
	}
	return Sr25519PublicKey{sc.NewFixedSequence(32, values...)}, nil
}
