package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	publicKeySerializedSize = 33
)

type EcdsaPublicKey struct {
	sc.FixedSequence[sc.U8] // size 33
}

func (s EcdsaPublicKey) SignatureType() sc.U8 {
	return PublicKeyEcdsa
}

func (s EcdsaPublicKey) Encode(buffer *bytes.Buffer) error {
	return s.FixedSequence.Encode(buffer)
}

func (s EcdsaPublicKey) Bytes() []byte {
	return sc.EncodedBytes(s)
}

func DecodeEcdsaPublicKey(buffer *bytes.Buffer) (EcdsaPublicKey, error) {
	seq, err := sc.DecodeFixedSequence[sc.U8](publicKeySerializedSize, buffer)
	if err != nil {
		return EcdsaPublicKey{}, err
	}
	return EcdsaPublicKey{seq}, nil
}

func NewEcdsaPublicKey(values ...sc.U8) (EcdsaPublicKey, error) {
	if len(values) != publicKeySerializedSize {
		return EcdsaPublicKey{}, newTypeError("EcdsaPublicKey")
	}
	return EcdsaPublicKey{sc.NewFixedSequence(publicKeySerializedSize, values...)}, nil
}
