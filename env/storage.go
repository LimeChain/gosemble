package env

/*
	Storage: Interface for manipulating the storage from within the runtime.
*/
//go:wasm-module env
//go:export ext_storage_changes_root_version_1
func extStorageChangesRootVersion1(parent_hash int64) int64

func ExtStorageChangesRootVersion1(parent_hash int64) int64 {
	return extStorageChangesRootVersion1(parent_hash)
}

//go:wasm-module env
//go:export ext_storage_clear_prefix_version_1
func extStorageClearPrefixVersion1(prefix int64)

func ExtStorageClearPrefixVersion1(prefix int64) {
	extStorageClearPrefixVersion1(prefix)
}

//go:wasm-module env
//go:export ext_storage_clear_version_1
func extStorageClearVersion1(key_data int64)

func ExtStorageClearVersion1(key_data int64) {
	extStorageClearVersion1(key_data)
}

//go:wasm-module env
//go:export ext_storage_get_version_1
func extStorageGetVersion1(key int64) int64

func ExtStorageGetVersion1(key int64) int64 {
	return extStorageGetVersion1(key)
}

//go:wasm-module env
//go:export ext_storage_next_key_version_1
func extStorageNextKeyVersion1(key int64) int64

func ExtStorageNextKeyVersion1(key int64) int64 {
	return extStorageNextKeyVersion1(key)
}

//go:wasm-module env
//go:export ext_storage_read_version_1
func extStorageReadVersion1(key int64, value_out int64, offset int32) int64

func ExtStorageReadVersion1(key int64, value_out int64, offset int32) int64 {
	return extStorageReadVersion1(key, value_out, offset)
}

//go:wasm-module env
//go:export ext_storage_root_version_1
func extStorageRootVersion1() int64

func ExtStorageRootVersion1() int64 {
	return extStorageRootVersion1()
}

//go:wasm-module env
//go:export ext_storage_set_version_1
func extStorageSetVersion1(key int64, value int64)

func ExtStorageSetVersion1(key int64, value int64) {
	extStorageSetVersion1(key, value)
}

//go:wasm-module env
//go:export ext_storage_exists_version_1
func extStorageExistsVersion1(key int64) int32

func ExtStorageExistsVersion1(key int64) int32 {
	return extStorageExistsVersion1(key)
}
