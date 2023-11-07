package types

import (
	"bytes"
	"errors"

	sc "github.com/LimeChain/goscale"
)

type H512 struct {
	sc.FixedSequence[sc.U8] // size 64
}

func NewH512(values ...sc.U8) (H512, error) {
	if len(values) != 64 {
		return H512{}, errors.New("H512 should be of size 64")
	}
	return H512{sc.NewFixedSequence(64, values...)}, nil
}

func (h H512) Encode(buffer *bytes.Buffer) {
	h.FixedSequence.Encode(buffer)
}

func DecodeH512(buffer *bytes.Buffer) (H512, error) {
	h := H512{}
	fixedSequence, err := sc.DecodeFixedSequence[sc.U8](64, buffer)
	if err != nil {
		return H512{}, err
	}
	h.FixedSequence = fixedSequence
	return h, nil
}

func (h H512) Bytes() []byte {
	return sc.EncodedBytes(h)
}
