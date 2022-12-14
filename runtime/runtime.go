/*
Targets WebAssembly MVP
*/
package main

import "github.com/LimeChain/gosemble/frame/core"

// TODO:
// remove the _start export and find a way to call it from the runtime to initialize the memory.
// TinyGo requires to have a main function to compile to Wasm.
func main() {}

//go:export Core_version
func CoreVersion(dataPtr int32, dataLen int32) int64 {
	return core.Version(dataPtr, dataLen)
}

//go:export Core_initialize_block
func CoreInitializeBlock(dataPtr int32, dataLen int32) {
	core.InitializeBlock(dataPtr, dataLen)
}

//go:export Core_execute_block
func CoreExecuteBlock(dataPtr int32, dataLen int32)

//go:export BlockBuilder_apply_extrinsic
func BlockBuilderApplyExtrinsic(dataPtr int32, dataLen int32) int64

//go:export BlockBuilder_finalize_block
func BlockBuilderFinalizeBlock(dataPtr int32, dataLen int32) int64

//go:export BlockBuilder_inherent_extrinisics
func BlockBuilderInherentExtrinisics(dataPtr int32, dataLen int32) int64

//go:export BlockBuilder_check_inherents
func BlockBuilderCheckInherents(dataPtr int32, dataLen int32) int64

//go:export BlockBuilder_random_seed
func BlockBuilderRandomSeed(dataPtr int32, dataLen int32) int64
