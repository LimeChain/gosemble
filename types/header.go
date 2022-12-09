package types

import "bytes"

type BlockNumber uint32

type Header struct {
	ParentHash     Blake2bHash
	Number         BlockNumber
	StateRoot      Hash
	ExtrinsicsRoot Hash
	Digest         Digest
}

func (h Header) Encode(buffer *bytes.Buffer) {
	panic("not implemented")
}

func DecodeHeader(buffer *bytes.Buffer) Header {
	panic("not implemented")
}
