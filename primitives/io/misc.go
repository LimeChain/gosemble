package io

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

type Misc interface {
	PrintHex([]byte)
	PrintUtf8([]byte)
	RuntimeVersion([]byte) []byte
}

type misc struct {
	memoryTranslator utils.WasmMemoryTranslator
}

func NewMisc() Misc {
	return &misc{
		memoryTranslator: utils.NewMemoryTranslator(),
	}
}

func (m misc) PrintHex(data []byte) {
	dataOffsetSize := m.memoryTranslator.BytesToOffsetAndSize(data)
	env.ExtMiscPrintHexVersion1(dataOffsetSize)
}

func (m misc) PrintUtf8(data []byte) {
	dataOffsetSize := m.memoryTranslator.BytesToOffsetAndSize(data)
	env.ExtMiscPrintUtf8Version1(dataOffsetSize)
}

func (m misc) RuntimeVersion(codeBlob []byte) []byte {
	codeBlobOffsetSize := m.memoryTranslator.BytesToOffsetAndSize(codeBlob)
	resOffsetSize := env.ExtMiscRuntimeVersionVersion1(codeBlobOffsetSize)
	offset, size := m.memoryTranslator.Int64ToOffsetAndSize(resOffsetSize)
	return m.memoryTranslator.GetWasmMemorySlice(offset, size)
}
