package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	signatureEcdsaLength = 65
)

type SignatureEcdsa struct {
	sc.FixedSequence[sc.U8] // size 65
}

func NewSignatureEcdsa(values ...sc.U8) SignatureEcdsa {
	return SignatureEcdsa{sc.NewFixedSequence(signatureEcdsaLength, values...)}
}

func (s SignatureEcdsa) Encode(buffer *bytes.Buffer) {
	s.FixedSequence.Encode(buffer)
}

func DecodeSignatureEcdsa(buffer *bytes.Buffer) SignatureEcdsa {
	s := SignatureEcdsa{}
	s.FixedSequence = sc.DecodeFixedSequence[sc.U8](signatureEcdsaLength, buffer)
	return s
}

func (s SignatureEcdsa) Bytes() []byte {
	return sc.EncodedBytes(s)
}
