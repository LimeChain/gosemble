package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type Sr25519Signer struct {
	sc.FixedSequence[sc.U8] // size 32
}

func (s Sr25519Signer) Encode(buffer *bytes.Buffer) {
	s.FixedSequence.Encode(buffer)
}

func (s Sr25519Signer) Bytes() []byte {
	return sc.EncodedBytes(s)
}

func DecodeSr25519Signer(buffer *bytes.Buffer) (Sr25519Signer, error) {
	seq, err := sc.DecodeFixedSequence[sc.U8](32, buffer)
	if err != nil {
		return Sr25519Signer{}, err
	}
	return Sr25519Signer{seq}, nil
}

func NewSr25519Signer(values ...sc.U8) Sr25519Signer {
	if len(values) != 32 {
		log.Critical("Sr25519Signer should be of size 32")
	}
	return Sr25519Signer{sc.NewFixedSequence(32, values...)}
}
