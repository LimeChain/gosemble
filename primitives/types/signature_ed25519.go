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

func (s SignatureEd25519) Encode(buffer *bytes.Buffer) {
	s.FixedSequence.Encode(buffer)
}

func DecodeSignatureEd25519(buffer *bytes.Buffer) SignatureEd25519 {
	s := SignatureEd25519{}
	s.FixedSequence = sc.DecodeFixedSequence[sc.U8](signatureEd25519Length, buffer)
	return s
}

func (s SignatureEd25519) Bytes() []byte {
	return sc.EncodedBytes(s)
}
