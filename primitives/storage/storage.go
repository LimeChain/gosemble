package storage

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

func ChangesRoot(parent_hash int64) int64 { panic("Not implemented!") }

func ClearPrefix(prefix int64) { panic("Not implemented!") }

func Clear(key_data int64) { panic("Not implemented!") }

func NextKey(key int64) int64 { panic("Not implemented!") }

func Read(key int64, value_out int64, offset int32) int64 { panic("Not implemented!") }

func Root() int64 { panic("Not implemented!") }

func Exists(key int64) int32 { panic("Not implemented!") }

func Set(key []byte, value []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	valueOffsetSize := utils.BytesToOffsetAndSize(value)
	env.ExtStorageSetVersion1(keyOffsetSize, valueOffsetSize)
}

func Get(key []byte) []byte {
	psKey := utils.BytesToOffsetAndSize(key)
	psValue := env.ExtStorageGetVersion1(psKey)
	offset, size := utils.Int64ToOffsetAndSize(psValue)
	value := utils.ToWasmMemorySlice(offset, size)
	return value
}
