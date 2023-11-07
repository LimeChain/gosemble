package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func NewBlock(header types.Header, extrinsics sc.Sequence[types.UncheckedExtrinsic]) primitives.Block {
	return block{
		header:     header,
		extrinsics: extrinsics,
	}
}

type block struct {
	header     types.Header
	extrinsics sc.Sequence[types.UncheckedExtrinsic]
}

func (b block) Encode(buffer *bytes.Buffer) error {
	_, err := buffer.Write(b.header.Bytes())
	if err != nil {
		return err
	}
	_, err = buffer.Write(b.extrinsics.Bytes())
	return err
}

func (b block) Bytes() []byte {
	return sc.EncodedBytes(b)
}

func (b block) Header() types.Header {
	return b.header
}

func (b block) Extrinsics() sc.Sequence[types.UncheckedExtrinsic] {
	return b.extrinsics
}
