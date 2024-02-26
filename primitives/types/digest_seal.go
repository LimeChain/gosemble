package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type DigestSeal struct {
	ConsensusEngineId sc.FixedSequence[sc.U8]
	Message           sc.Sequence[sc.U8]
}

func NewDigestSeal(consensusEngineId sc.FixedSequence[sc.U8], message sc.Sequence[sc.U8]) DigestSeal {
	return DigestSeal{consensusEngineId, message}
}

func (ds DigestSeal) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer, ds.ConsensusEngineId, ds.Message)
}

func (ds DigestSeal) Bytes() []byte {
	return sc.EncodedBytes(ds)
}
