//go:build nonwasmenv

package storage

func ChangesRoot(parent_hash int64) int64 {
	panic("not implemented")
}

func Clear(key []byte) {
	panic("not implemented")
}

func ClearPrefix(key []byte, limit []byte) {
	panic("not implemented")
}

func Exists(key []byte) int32 {
	panic("not implemented")
}

func Get(key []byte) []byte {
	panic("not implemented")
}

func NextKey(key int64) int64 {
	panic("not implemented")
}

func Read(key int64, value_out int64, offset int32) int64 {
	panic("not implemented")
}

func Root(key int32) []byte {
	panic("not implemented")
}

func Set(key []byte, value []byte) {
	panic("not implemented")
}
