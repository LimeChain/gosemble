package utils

import (
	"reflect"
	"unsafe"
)

func BytesToPointerAndSize(data []byte) uint64 {
	dataPtr := uintptr(unsafe.Pointer(&data))
	dataLen := len(data)
	return PointerAndSizeToInt64(uint32(dataPtr), uint32(dataLen))
}

func PointerAndSizeToString(ptr uint32, size uint32) string {
	// We use SliceHeader, not StringHeader as it allows us to fix the capacity to what was allocated.
	// Tinygo requires these as uintptrs even if they are int fields.
	// https://github.com/tinygo-org/tinygo/issues/1284
	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  uintptr(size),
		Cap:  uintptr(size),
	}))
}

func StringToPointerAndSize(s string) (uint32, uint32) {
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return uint32(unsafePtr), uint32(len(buf))
}

func PointerAndSizeToInt64(dataPtr uint32, dataLen uint32) uint64 {
	return uint64(dataPtr) | (uint64(dataLen) << 32)
}

func Int64ToPointerAndSize(pointerAndSize uint64) (dataPtr uint32, dataLen uint32) {
	return uint32(pointerAndSize), uint32(pointerAndSize >> 32)
}
