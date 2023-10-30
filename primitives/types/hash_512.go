package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type H512 struct {
	sc.FixedSequence[sc.U8] // size 64
}

func NewH512(values ...sc.U8) H512 {
	if len(values) != 64 {
		log.Critical("H512 should be of size 64")
	}
	return H512{sc.NewFixedSequence(64, values...)}
}

func (h H512) Encode(buffer *bytes.Buffer) {
	h.FixedSequence.Encode(buffer)
}

func DecodeH512(buffer *bytes.Buffer) H512 {
	h := H512{}
	h.FixedSequence = sc.DecodeFixedSequence[sc.U8](64, buffer)
	return h
}

func (h H512) Bytes() []byte {
	return sc.EncodedBytes(h)
}
