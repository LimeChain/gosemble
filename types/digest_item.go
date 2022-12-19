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
