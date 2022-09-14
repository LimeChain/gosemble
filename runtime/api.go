/*
	TODO implement runtime APIs:
	- Core - Version 3.
	- BlockBuilder - Version 4.
*/
package main

type Core interface {
	/*
		SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded (version types.VersionData) data.
	*/
	//export Core_version
	CoreVersion(dataPtr int32, dataLen int32) int64

	/*
		SCALE encoded arguments (block types.Block) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
	*/
	//export Core_execute_block
	ExecuteBlock(dataPtr int32, dataLen int32)

	/*
		SCALE encoded arguments (header *types.Header) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
	*/
	//export Core_initialize_block
	CoreInitializeBlock(dataPtr int32, dataLen int32)
}

type BlockBuilder interface {
	/*
		SCALE encoded arguments (extrinsic types.Extrinsic) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([]byte) data.
	*/
	//export BlockBuilder_apply_extrinsic
	ApplyExtrinsic(dataPtr int32, dataLen int32) int64

	/*
		SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded (types.Header) data.
	*/
	//export BlockBuilder_finalize_block
	FinalizeBlock(dataPtr int32, dataLen int32) int64

	/*
		SCALE encoded arguments (data types.InherentsData) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([]types.Extrinsic) data.
	*/
	//export BlockBuilder_inherent_extrinisics
	InherentExtrinisics(dataPtr int32, dataLen int32) int64

	/*
		SCALE encoded arguments (block types.Block, data types.InherentsData) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([]byte) data.
	*/
	//export BlockBuilder_check_inherents
	CheckInherents(dataPtr int32, dataLen int32) int64

	/*
		SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([32]byte) data.
	*/
	//export BlockBuilder_random_seed
	RandomSeed(dataPtr int32, dataLen int32) int64
}
