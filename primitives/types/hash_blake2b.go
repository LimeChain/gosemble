package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

type Blake2bHash struct {
	sc.FixedSequence[sc.U8] // size 32
}

func NewBlake2bHash(values ...sc.U8) Blake2bHash {
	if len(values) != 32 {
		log.Critical("Blake2bHash should be of size 32")
	}
	return Blake2bHash{sc.NewFixedSequence(32, values...)}
}

func (h Blake2bHash) Encode(buffer *bytes.Buffer) {
	h.FixedSequence.Encode(buffer)
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
