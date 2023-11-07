package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	signatureEd25519Length = 64
)

type SignatureEd25519 struct {
	sc.FixedSequence[sc.U8] // size 64
}

func NewSignatureEd25519(values ...sc.U8) SignatureEd25519 {
	return SignatureEd25519{sc.NewFixedSequence(signatureEd25519Length, values...)}
}

func (s SignatureEd25519) Encode(buffer *bytes.Buffer) error {
	return s.FixedSequence.Encode(buffer)
}

func DecodeSignatureEd25519(buffer *bytes.Buffer) (SignatureEd25519, error) {
	s := SignatureEd25519{}
	seq, err := sc.DecodeFixedSequence[sc.U8](signatureEd25519Length, buffer)
	if err != nil {
		return SignatureEd25519{}, err
	}
	s.FixedSequence = seq

	return s, nil
}

func (s SignatureEd25519) Bytes() []byte {
	return sc.EncodedBytes(s)
}
