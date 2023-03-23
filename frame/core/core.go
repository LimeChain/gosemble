/*
Core - Version 3.
*/
package core

import (
	"bytes"

	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/executive"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

type Core interface {
	Version(dataPtr int32, dataLen int32) int64
	ExecuteBlock(dataPtr int32, dataLen int32)
	InitializeBlock(dataPtr int32, dataLen int32)
}

/*
https://spec.polkadot.network/#defn-rt-core-version

SCALE encoded arguments () allocated in the Wasm VM memory, passed as:

	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
	returns a pointer-size to the SCALE-encoded (version types.VersionData) data.
*/
func Version(dataPtr int32, dataLen int32) int64 {
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
func InitializeBlock(dataPtr int32, dataLen int32) {
	data := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)

	header := primitives.DecodeHeader(buffer)
	executive.InitializeBlock(header)
}

/*
https://spec.polkadot.network/#sect-rte-core-execute-block

SCALE encoded arguments (block types.Block) allocated in the Wasm VM memory, passed as:

	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
*/
func ExecuteBlock(dataPtr int32, dataLen int32) {
	data := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)

	block := types.DecodeBlock(buffer)
	executive.ExecuteBlock(block)
}
