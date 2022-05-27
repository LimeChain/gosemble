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
	CoreVersion(dataPtr uint32, dataLen uint32) uint64

	/*
		SCALE encoded arguments (block types.Block) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
	*/
	//export Core_execute_block
	ExecuteBlock(dataPtr uint32, dataLen uint32)

	/*
		SCALE encoded arguments (header *types.Header) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
	*/
	//export Core_initialize_block
	CoreInitializeBlock(dataPtr uint32, dataLen uint32)
}

type BlockBuilder interface {
	/*
		SCALE encoded arguments (extrinsic types.Extrinsic) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([]byte) data.
	*/
	//export BlockBuilder_apply_extrinsic
	ApplyExtrinsic(dataPtr uint32, dataLen uint32) uint64

	/*
		SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded (types.Header) data.
	*/
	//export BlockBuilder_finalize_block
	FinalizeBlock(dataPtr uint32, dataLen uint32) uint64

	/*
		SCALE encoded arguments (data types.InherentsData) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([]types.Extrinsic) data.
	*/
	//export BlockBuilder_inherent_extrinisics
	InherentExtrinisics(dataPtr uint32, dataLen uint32) uint64

	/*
		SCALE encoded arguments (block types.Block, data types.InherentsData) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([]byte) data.
	*/
	//export BlockBuilder_check_inherents
	CheckInherents(dataPtr uint32, dataLen uint32) uint64

	/*
		SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([32]byte) data.
	*/
	//export BlockBuilder_random_seed
	RandomSeed(dataPtr uint32, dataLen uint32) uint64
}
