package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	publicKeySerializedSize = 33
)

type EcdsaSigner struct {
	sc.FixedSequence[sc.U8] // size 33
}

func (s EcdsaSigner) SignatureType() sc.U8 {
	return PublicKeyEcdsa
}

func (s EcdsaSigner) Encode(buffer *bytes.Buffer) {
	s.FixedSequence.Encode(buffer)
}

func (s EcdsaSigner) Bytes() []byte {
	return sc.EncodedBytes(s)
}

func DecodeEcdsaSigner(buffer *bytes.Buffer) (EcdsaSigner, error) {
	seq, err := sc.DecodeFixedSequence[sc.U8](publicKeySerializedSize, buffer)
	if err != nil {
		return EcdsaSigner{}, err
	}
	return EcdsaSigner{seq}, nil
}

func NewEcdsaSigner(values ...sc.U8) (EcdsaSigner, error) {
	if len(values) != publicKeySerializedSize {
		return EcdsaSigner{}, newTypeError("EcdsaSigner")
	}
	return EcdsaSigner{sc.NewFixedSequence(publicKeySerializedSize, values...)}, nil
}
