package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type BlockNumber sc.U32

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
	panic("not implemented DecodeHeader")
}
