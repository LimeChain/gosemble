package types

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
)

type DigestItem struct {
	Engine  sc.FixedSequence[sc.U8]
	Payload sc.Sequence[sc.U8]
}

func (di DigestItem) Encode(buffer *bytes.Buffer) {
	di.Engine.Encode(buffer)
	di.Payload.Encode(buffer)
}

func (di DigestItem) Bytes() []byte {
	buffer := &bytes.Buffer{}
	di.Encode(buffer)

	return buffer.Bytes()
}
