package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Block struct {
	Header     Header
	Extrinsics sc.Sequence[UncheckedExtrinsic]
}

func (b Block) Encode(buffer *bytes.Buffer) {
	buffer.Write(b.Header.Bytes())
	buffer.Write(b.Extrinsics.Bytes())
}

func (b Block) Bytes() []byte {
	return sc.EncodedBytes(b)
}

func DecodeBlock(buffer *bytes.Buffer) Block {
	return Block{
		Header:     DecodeHeader(buffer),
		Extrinsics: sc.DecodeSequence[UncheckedExtrinsic](buffer),
	}
}
