//go:build !nonwasmenv

package storage

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

func Append(key []byte, value []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	valueOffsetSize := utils.BytesToOffsetAndSize(value)
	env.ExtStorageAppendVersion1(keyOffsetSize, valueOffsetSize)
}

func ChangesRoot(parent_hash int64) int64 {
	panic("not implemented")
}

func Clear(key []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	env.ExtStorageClearVersion1(keyOffsetSize)
}

func ClearPrefix(key []byte, limit []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	limitOffsetSize := utils.BytesToOffsetAndSize(limit)
	env.ExtStorageClearPrefixVersion2(keyOffsetSize, limitOffsetSize)
}

func Exists(key []byte) int32 {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	return env.ExtStorageExistsVersion1(keyOffsetSize)
}

func Get(key []byte) sc.Option[sc.Sequence[sc.U8]] {
	value := get(key)

	buffer := &bytes.Buffer{}
	buffer.Write(value)

	return sc.DecodeOption[sc.Sequence[sc.U8]](buffer)
}

// GetDecode gets the storage value and returns it decoded. The result from Get is Option<sc.Sequence[sc.U8]>.
// If the option is empty, it returns the default value T.
// If the option is not empty, it decodes it using decodeFunc and returns it.
func GetDecode[T sc.Encodable](key []byte, decodeFunc func(buffer *bytes.Buffer) T) T {
	option := Get(key)

	if !option.HasValue {
		return *new(T)
	}

	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(option.Value))

	return decodeFunc(buffer)
}

func NextKey(key int64) int64 {
	panic("not implemented")
}

func Read(key []byte, valueOut []byte, offset int32) sc.Option[sc.U32] {
	value := read(key, valueOut, offset)

	buffer := &bytes.Buffer{}
	buffer.Write(value)

	return sc.DecodeOption[sc.U32](buffer)
}

func Root(version int32) []byte {
	valueOffsetSize := env.ExtStorageRootVersion2(version)
	offset, size := utils.Int64ToOffsetAndSize(valueOffsetSize)
	value := utils.ToWasmMemorySlice(offset, size)
	return value
}

func Set(key []byte, value []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	valueOffsetSize := utils.BytesToOffsetAndSize(value)
	env.ExtStorageSetVersion1(keyOffsetSize, valueOffsetSize)
}

// TakeBytes gets the storage value. The result from Get is Option<sc.Sequence[sc.U8]>.
// If the option is empty, it returns nil.
// If the option is not empty, it clears it and returns the sequence as bytes.
func TakeBytes(key []byte) []byte {
	option := Get(key)

	if !option.HasValue {
		return nil
	}

	Clear(key)

	return sc.SequenceU8ToBytes(option.Value)
}

// TakeDecode gets the storage value and returns it decoded. The result from Get is Option<sc.Sequence[sc.U8]>.
// If the option is empty, it returns default value T.
// If the option is not empty, it clears it and returns decodeFunc(value).
func TakeDecode[T sc.Encodable](key []byte, decodeFunc func(buffer *bytes.Buffer) T) T {
	option := Get(key)

	if !option.HasValue {
		return *new(T)
	}

	Clear(key)

	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(option.Value))

	return decodeFunc(buffer)
}

// get gets the value from storage by the provided key. The wasm memory slice (value)
// represents an encoded Option<sc.Sequence[sc.U8]> (option of encoded slice).
func get(key []byte) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	valueOffsetSize := env.ExtStorageGetVersion1(keyOffsetSize)
	offset, size := utils.Int64ToOffsetAndSize(valueOffsetSize)
	value := utils.ToWasmMemorySlice(offset, size)
	return value
}

// read reads the given key value from storage, placing the value into buffer valueOut depending on offset.
// The wasm memory slice represents an encoded Option<sc.U32> representing the number of bytes left at supplied offset.
func read(key []byte, valueOut []byte, offset int32) []byte {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	valueOutOffsetSize := utils.BytesToOffsetAndSize(valueOut)

	resultOffsetSize := env.ExtStorageReadVersion1(keyOffsetSize, valueOutOffsetSize, offset)
	offset, size := utils.Int64ToOffsetAndSize(resultOffsetSize)
	value := utils.ToWasmMemorySlice(offset, size)

	return value
}
