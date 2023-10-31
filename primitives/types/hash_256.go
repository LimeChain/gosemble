package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type H256 struct {
	sc.FixedSequence[sc.U8] // size 32
}

func NewH256(values ...sc.U8) H256 {
	if len(values) != 32 {
		log.Critical("H256 should be of size 32")
	}
	return H256{sc.NewFixedSequence(32, values...)}
}

func (h H256) Encode(buffer *bytes.Buffer) {
	h.FixedSequence.Encode(buffer)
}

func DecodeH256(buffer *bytes.Buffer) (H256, error) {
	h := H256{}
	fixedSequence, err := sc.DecodeFixedSequence[sc.U8](32, buffer)
	if err != nil {
		return H256{}, err
	}
	h.FixedSequence = fixedSequence
	return h, nil
}

func (h H256) Bytes() []byte {
	return sc.EncodedBytes(h)
}
