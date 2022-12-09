/*
	Targets WebAssembly MVP
*/
package main

import (
	"bytes"

	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/types"
	"github.com/LimeChain/gosemble/utils"
)

// TODO: remove the _start export and find a way to call it from the runtime to initialize the memory.
// TinyGo requires to have a main function to compile to Wasm.
func main() {}

/*
	https://spec.polkadot.network/#defn-rt-core-version

	SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded (version types.VersionData) data.
*/
//go:export Core_version
func CoreVersion(dataPtr int32, dataLen int32) int64 {
	buffer := &bytes.Buffer{}
	constants.RuntimeVersion.Encode(buffer)
	// TODO: retain the pointer to the scaleEncVersion
	// utils.Retain(scaleEncVersion)
	return utils.BytesToOffsetAndSize(buffer.Bytes())
}

/*
https://spec.polkadot.network/#sect-rte-core-initialize-block

SCALE encoded arguments (header *types.Header) allocated in the Wasm VM memory, passed as:
	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
*/
//go:export Core_initialize_block
func CoreInitializeBlock(dataPtr int32, dataLen int32) {
	data := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := &bytes.Buffer{}
	buffer.Write(data)
	header := types.DecodeHeader(buffer)
	executive.InitializeBlock(&header)
}

/*
	https://spec.polkadot.network/#sect-rte-core-execute-block

	SCALE encoded arguments (block types.Block) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
*/
//go:export Core_execute_block
func ExecuteBlock(dataPtr int32, dataLen int32) {

}
