package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Header[N sc.Numeric] struct {
	ParentHash     Blake2bHash
	Number         N
	StateRoot      H256
	ExtrinsicsRoot H256
	Digest         Digest
}

func (h Header[N]) Encode(buffer *bytes.Buffer) {
	h.ParentHash.Encode(buffer)
	sc.ToCompact(h.Number).Encode(buffer)
	h.StateRoot.Encode(buffer)
	h.ExtrinsicsRoot.Encode(buffer)
	h.Digest.Encode(buffer)
}

func (h Header[N]) Bytes() []byte {
	buffer := &bytes.Buffer{}
	h.Encode(buffer)
	return buffer.Bytes()
}

func DecodeHeader[N sc.Numeric](buffer *bytes.Buffer) Header[N] {
	parentHash := DecodeBlake2bHash(buffer)
	blockNumber := sc.DecodeCompact(buffer)
	stateRoot := DecodeH256(buffer)
	extrinsicRoot := DecodeH256(buffer)
	digest := DecodeDigest(buffer)

	return Header[N]{
		ParentHash:     parentHash,
		Number:         sc.NewNumeric[N](sc.To[sc.U64](sc.U128(blockNumber))),
		StateRoot:      stateRoot,
		ExtrinsicsRoot: extrinsicRoot,
		Digest:         digest,
	}
}
