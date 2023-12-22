package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	publicKeySr25519Length = 32
)

type Sr25519PublicKey struct {
	sc.FixedSequence[sc.U8] // size 32
}

func NewSr25519PublicKey(values ...sc.U8) (Sr25519PublicKey, error) {
	if len(values) != publicKeySr25519Length {
		return Sr25519PublicKey{}, newTypeError("Sr25519PublicKey")
	}
	return Sr25519PublicKey{sc.NewFixedSequence(publicKeySr25519Length, values...)}, nil
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
	seq, err := sc.DecodeFixedSequence[sc.U8](publicKeySr25519Length, buffer)
	if err != nil {
		return Sr25519PublicKey{}, err
	}
	return Sr25519PublicKey{seq}, nil
}

func DecodeSequenceSr25519PublicKey(buffer *bytes.Buffer) (sc.Sequence[Sr25519PublicKey], error) {
	return sc.DecodeSequenceWith[Sr25519PublicKey](buffer, DecodeSr25519PublicKey)
}
