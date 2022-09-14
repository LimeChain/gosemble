package utils

import (
	"unsafe"
)

var alivePointers = map[uintptr]interface{}{}

func Retain(data []byte) {
	ptr := &data[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	alivePointers[unsafePtr] = data
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func Int64ToOffsetAndSize(offsetAndSize int64) (offset int32, size int32) {
	return int32(offsetAndSize), int32(offsetAndSize >> 32)
}

func OffsetAndSizeToInt64(offset int32, size int32) int64 {
	return int64(offset) | (int64(size) << 32)
}

func SliceToOffset(data []byte) uintptr {
	return uintptr(unsafe.Pointer(&data[0]))
}

func BytesToOffsetAndSize(data []byte) int64 {
	offset := SliceToOffset(data)
	size := len(data)
	return OffsetAndSizeToInt64(int32(offset), int32(size))
}

func StringToOffsetAndSize(str string) int64 {
	data := []byte(str)
	offset := SliceToOffset(data)
	size := len(data)
	return OffsetAndSizeToInt64(int32(offset), int32(size))
}

// func PointerAndSizeToString(ptr int32, size int32) string {
// 	// We use SliceHeader, not StringHeader as it allows us to fix the capacity to what was allocated.
// 	// Tinygo requires these as uintptrs even if they are int fields.
// 	// https://github.com/tinygo-org/tinygo/issues/1284
// 	return *(*string)(unsafe.Pointer(&reflect.SliceHeader{
// 		Data: uintptr(ptr),
// 		Len:  uintptr(size),
// 		Cap:  uintptr(size),
// 	}))
// }

func ToWasmMemorySlice(offset int32, size int32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(offset))), uintptr(size))
}

func WriteToMemory(offset int32, size int32, data [256]byte) {
	memory := ToWasmMemorySlice(offset, size)

	for i := int32(0); i < size; i++ {
		memory[i] = byte(data[i])

		// ptr := (*byte)(unsafe.Pointer(uintptr(offset) + uintptr(i)))
		// *ptr = byte(data[i])

		// ptr := unsafe.Pointer(uintptr(offset))
		// bs := (*[MAX_ARRAY_SIZE]byte)(ptr)[:]
	}
}
