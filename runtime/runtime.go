/*
Targets WebAssembly MVP
*/
package main

import (
	"github.com/radkomih/gosemble/constants"
	"github.com/radkomih/gosemble/types"
	"github.com/radkomih/gosemble/utils"
)

/*
	SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
	returns a pointer-size to the SCALE-encoded (version types.VersionData) data.
*/
//go:export Core_version
func CoreVersion(dataPtr int32, dataLen int32) int64 {
	scaleEncVersion, err := constants.VersionDataConfig.Encode()
	utils.PanicOnError(err)
	// TODO: retain the pointer to the scaleEncVersion
	// utils.Retain(scaleEncVersion)
	return utils.BytesToOffsetAndSize(scaleEncVersion)
}

/*
SCALE encoded arguments (header *types.Header) allocated in the Wasm VM memory, passed as:
dataPtr - i32 pointer to the memory location.
dataLen - i32 length (in bytes) of the encoded arguments.
*/
//go:export Core_initialize_block
func CoreInitializeBlock(dataPtr int32, dataLen int32) {
	data := utils.ToWasmMemorySlice(dataPtr, dataLen)
	header := (&types.Header{}).Decode(data)
	_ = header
	extStorageSetVersion1(int64(123), int64(456))
	extStorageGetVersion1(int64(123))
}

/*
	SCALE encoded arguments (block types.Block) allocated in the Wasm VM memory, passed as:
	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
*/
//go:export Core_execute_block
func ExecuteBlock(dataPtr int32, dataLen int32) {

}

// TODO: remove the _start export and find a way to call it from the runtime to initialize the memory.
// TinyGo requires to have a main function to compile to Wasm.
func main() {}
