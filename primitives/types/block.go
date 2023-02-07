package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Block struct {
	Header     Header
	Extrinsics sc.Sequence[CheckedExtrinsic]
}

func (b Block) Encode(buffer *bytes.Buffer) {
	panic("not implemented Block Encode")
}

func DecodeBlock(buffer *bytes.Buffer) Block {
	panic("not implemented DecodeBlock")
}
