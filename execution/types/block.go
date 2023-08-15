package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Block[N sc.Numeric] struct {
	Header     types.Header[N]
	Extrinsics sc.Sequence[UncheckedExtrinsic]
}

func (b Block[N]) Encode(buffer *bytes.Buffer) {
	buffer.Write(b.Header.Bytes())
	buffer.Write(b.Extrinsics.Bytes())
}

func (b Block[N]) Bytes() []byte {
	return sc.EncodedBytes(b)
}
