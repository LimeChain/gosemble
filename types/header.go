package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type BlockNumber struct {
	sc.U32
}

type Header struct {
	ParentHash     Blake2bHash
	Number         BlockNumber
	StateRoot      Hash
	ExtrinsicsRoot Hash
	Digest         Digest
}

func (h Header) Encode(buffer *bytes.Buffer) {
	h.ParentHash.Encode(buffer)
	h.Number.Encode(buffer)
	h.StateRoot.Encode(buffer)
	h.ExtrinsicsRoot.Encode(buffer)
	h.Digest.Encode(buffer)
}

func (h Header) Bytes() []byte {
	buffer := &bytes.Buffer{}
	h.Encode(buffer)

	return buffer.Bytes()
}

func DecodeHeader(buffer *bytes.Buffer) Header {
	parentHash := sc.DecodeFixedSequence[sc.U8](32, buffer)
	blockNumber := sc.DecodeCompact(buffer)
	stateRoot := sc.DecodeFixedSequence[sc.U8](32, buffer)
	extrinsicRoot := sc.DecodeFixedSequence[sc.U8](32, buffer)
	digest := DecodeDigest(buffer)

	return Header{
		ParentHash: Blake2bHash{
			parentHash,
		},
		Number: BlockNumber{
			sc.U32(blockNumber.ToBigInt().Int64()),
		},
		StateRoot:      stateRoot,
		ExtrinsicsRoot: extrinsicRoot,
		Digest:         digest,
	}
}
