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
	panic("not implemented Header Encode")
}

func DecodeHeader(buffer *bytes.Buffer) Header {
	parentHash := sc.DecodeFixedSequence[sc.U8](32, buffer)
	number := sc.DecodeCompact(buffer)
	stateRoot := sc.DecodeFixedSequence[sc.U8](32, buffer)
	extrinsicRoot := sc.DecodeFixedSequence[sc.U8](32, buffer)
	digest := DecodeDigest(buffer)

	return Header{
		ParentHash: Blake2bHash{
			parentHash,
		},
		Number: BlockNumber{
			sc.U32(number),
		},
		StateRoot:      stateRoot,
		ExtrinsicsRoot: extrinsicRoot,
		Digest:         digest,
	}
}
