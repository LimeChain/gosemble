package utils

import (
	"unsafe"
)

type WasmMemoryTranslator interface {
	Int64ToOffsetAndSize(offsetAndSize int64) (offset int32, size int32)
	Offset32(data []byte) int32
	BytesToOffsetAndSize(data []byte) int64
	GetWasmMemorySlice(offset int32, size int32) []byte
}

type memoryTranslator struct{}

func NewMemoryTranslator() WasmMemoryTranslator {
	return &memoryTranslator{}
}

func (m memoryTranslator) Int64ToOffsetAndSize(offsetAndSize int64) (offset int32, size int32) {
	return int32(offsetAndSize), int32(offsetAndSize >> 32)
}

func (m memoryTranslator) Offset32(data []byte) int32 {
	return int32(sliceToOffset(data))
}

func (m memoryTranslator) BytesToOffsetAndSize(data []byte) int64 {
	offset := sliceToOffset(data)
	size := len(data)
	return offsetAndSizeToInt64(int32(offset), int32(size))
}

func (m memoryTranslator) GetWasmMemorySlice(offset int32, size int32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(offset))), uintptr(size))
}

func sliceToOffset(data []byte) uintptr {
	if len(data) == 0 {
		return uintptr(unsafe.Pointer(nil))
	}

	return uintptr(unsafe.Pointer(&data[0]))
}

func offsetAndSizeToInt64(offset int32, size int32) int64 {
	return int64(offset) | (int64(size) << 32)
}
