package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type CallIndex struct {
	ModuleIndex   sc.U8
	FunctionIndex sc.U8
}

func (ci CallIndex) Encode(buffer *bytes.Buffer) {
	ci.ModuleIndex.Encode(buffer)
	ci.FunctionIndex.Encode(buffer)
}

func DecodeCallIndex(buffer *bytes.Buffer) CallIndex {
	ci := CallIndex{}
	ci.ModuleIndex = sc.DecodeU8(buffer)
	ci.FunctionIndex = sc.DecodeU8(buffer)
	return ci
}

func (ci CallIndex) Bytes() []byte {
	return sc.EncodedBytes(ci)
}
