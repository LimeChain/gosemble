package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type IncRefStatus sc.U8

const (
	IncRefStatusCreated IncRefStatus = iota
	IncRefStatusExisted
)

func (irs IncRefStatus) Encode(buffer *bytes.Buffer) error {
	return sc.U8(irs).Encode(buffer)
}

func (irs IncRefStatus) Bytes() []byte {
	return sc.EncodedBytes(irs)
}
