package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Call struct {
	CallIndex CallIndex
	Args      sc.Sequence[sc.U8]
}

func (c Call) Encode(buffer *bytes.Buffer) {
	c.CallIndex.Encode(buffer)
	c.Args.Encode(buffer)
}

func DecodeCall(buffer *bytes.Buffer) Call {
	c := Call{}
	c.CallIndex = DecodeCallIndex(buffer)
	c.Args = sc.DecodeSequence[sc.U8](buffer)
	return c
}

func (c Call) Bytes() []byte {
	return sc.EncodedBytes(c)
}

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
