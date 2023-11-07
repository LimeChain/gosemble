package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Ed25519Signer struct {
	sc.FixedSequence[sc.U8] // size 32
}

func (s Ed25519Signer) Encode(buffer *bytes.Buffer) {
	s.FixedSequence.Encode(buffer)
}

func (s Ed25519Signer) Bytes() []byte {
	return sc.EncodedBytes(s)
}

func DecodeEd25519Signer(buffer *bytes.Buffer) (Ed25519Signer, error) {
	seq, err := sc.DecodeFixedSequence[sc.U8](32, buffer)
	if err != nil {
		return Ed25519Signer{}, err
	}
	return Ed25519Signer{seq}, nil
}

func NewEd25519Signer(values ...sc.U8) (Ed25519Signer, error) {
	if len(values) != 32 {
		return Ed25519Signer{}, newTypeError("Ed25519Signer")
	}
	return Ed25519Signer{sc.NewFixedSequence(32, values...)}, nil
}
