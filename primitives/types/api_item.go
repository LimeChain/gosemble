package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type ApiItem struct {
	Name    sc.FixedSequence[sc.U8] // size 8
	Version sc.U32
}

func NewApiItem(name [8]byte, version sc.U32) ApiItem {
	return ApiItem{
		Name:    sc.BytesToFixedSequenceU8(name[:]),
		Version: version,
	}
}

func (ai ApiItem) Encode(buffer *bytes.Buffer) {
	ai.Name.Encode(buffer)
	ai.Version.Encode(buffer)
}

func DecodeApiItem(buffer *bytes.Buffer) ApiItem {
	return ApiItem{
		Name:    sc.DecodeFixedSequence[sc.U8](8, buffer),
		Version: sc.DecodeU32(buffer),
	}
}

func (ai ApiItem) Bytes() []byte {
	return sc.EncodedBytes(ai)
}
