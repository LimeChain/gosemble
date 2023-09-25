package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Block struct {
	Header     types.Header
	Extrinsics sc.Sequence[UncheckedExtrinsic]
}

func (b Block) Encode(buffer *bytes.Buffer) {
	buffer.Write(b.Header.Bytes())
	buffer.Write(b.Extrinsics.Bytes())
}

func (b Block) Bytes() []byte {
	return sc.EncodedBytes(b)
}
