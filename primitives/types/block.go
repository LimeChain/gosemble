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
	header := DecodeHeader(buffer)

	size := sc.DecodeCompact(buffer)
	length := size.ToBigInt()
	extrinsics := make([]UncheckedExtrinsic, length.Int64())

	for i := 0; i < len(extrinsics); i++ {
		extrinsics[i] = DecodeUncheckedExtrinsic(buffer)
	}

	return Block{
		Header:     header,
		Extrinsics: extrinsics,
	}
}
