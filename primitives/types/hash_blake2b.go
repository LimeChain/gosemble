package types

import (
	"bytes"
	"errors"

	sc "github.com/LimeChain/goscale"
)

type Blake2bHash struct {
	sc.FixedSequence[sc.U8] // size 32
}

func NewBlake2bHash(values ...sc.U8) (Blake2bHash, error) {
	if len(values) != 32 {
		return Blake2bHash{}, errors.New("Blake2bHash should be of size 32")
	}
	return Blake2bHash{sc.NewFixedSequence(32, values...)}, nil
}

func (h Blake2bHash) Encode(buffer *bytes.Buffer) error {
	return h.FixedSequence.Encode(buffer)
}

func DecodeBlake2bHash(buffer *bytes.Buffer) (Blake2bHash, error) {
	h := Blake2bHash{}
	fixedSequence, err := sc.DecodeFixedSequence[sc.U8](32, buffer)
	if err != nil {
		return Blake2bHash{}, err
	}
	h.FixedSequence = fixedSequence
	return h, nil
}

func (h Blake2bHash) Bytes() []byte {
	return sc.EncodedBytes(h)
}

func Blake2bHash69() Blake2bHash {
	hash69Bytes := []byte{69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69}
	hash69, _ := NewBlake2bHash(sc.BytesToFixedSequenceU8(hash69Bytes)...)
	return hash69
}
