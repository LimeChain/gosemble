package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type DigestPreRuntime struct {
	ConsensusEngineId sc.FixedSequence[sc.U8]
	Message           sc.Sequence[sc.U8]
}

func NewDigestPreRuntime(consensusEngineId sc.FixedSequence[sc.U8], message sc.Sequence[sc.U8]) DigestPreRuntime {
	return DigestPreRuntime{consensusEngineId, message}
}

func (dpr DigestPreRuntime) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer, dpr.ConsensusEngineId, dpr.Message)
}

func (dpr DigestPreRuntime) Bytes() []byte {
	return sc.EncodedBytes(dpr)
}
