package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type DecRefStatus sc.U8

const (
	DecRefStatusReaped DecRefStatus = iota
	DecRefStatusExists
)

func (drs DecRefStatus) Encode(buffer *bytes.Buffer) error {
	return sc.U8(drs).Encode(buffer)
}

func (drs DecRefStatus) Bytes() []byte {
	return sc.EncodedBytes(drs)
}
