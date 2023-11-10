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

func (s SignatureEcdsa) Encode(buffer *bytes.Buffer) error {
	return s.FixedSequence.Encode(buffer)
}

func DecodeSignatureEcdsa(buffer *bytes.Buffer) (SignatureEcdsa, error) {
	s := SignatureEcdsa{}
	seq, err := sc.DecodeFixedSequence[sc.U8](signatureEcdsaLength, buffer)
	if err != nil {
		return SignatureEcdsa{}, err
	}
	s.FixedSequence = seq
	return s, nil
}

func (s SignatureEcdsa) Bytes() []byte {
	return sc.EncodedBytes(s)
}
