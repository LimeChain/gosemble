package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	sr25519SignatureLength = 64
)

type SignatureSr25519 struct {
	sc.FixedSequence[sc.U8] // size 64
}

func NewSignatureSr25519(values ...sc.U8) SignatureSr25519 {
	return SignatureSr25519{sc.NewFixedSequence(sr25519SignatureLength, values...)}
}

func (s SignatureSr25519) Encode(buffer *bytes.Buffer) {
	s.FixedSequence.Encode(buffer)
}

func DecodeSignatureSr25519(buffer *bytes.Buffer) SignatureSr25519 {
	s := SignatureSr25519{}
	s.FixedSequence = sc.DecodeFixedSequence[sc.U8](sr25519SignatureLength, buffer)
	return s
}

func (s SignatureSr25519) Bytes() []byte {
	return sc.EncodedBytes(s)
}
