package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	ed25519SignatureLength = 64
)

type SignatureEd25519 struct {
	sc.FixedSequence[sc.U8] // size 64
}

func NewSignatureEd25519(values ...sc.U8) SignatureEd25519 {
	return SignatureEd25519{sc.NewFixedSequence(ed25519SignatureLength, values...)}
}

func (s SignatureEd25519) Encode(buffer *bytes.Buffer) {
	s.FixedSequence.Encode(buffer)
}

func DecodeSignatureEd25519(buffer *bytes.Buffer) SignatureEd25519 {
	s := SignatureEd25519{}
	s.FixedSequence = sc.DecodeFixedSequence[sc.U8](ed25519SignatureLength, buffer)
	return s
}

func (s SignatureEd25519) Bytes() []byte {
	return sc.EncodedBytes(s)
}
