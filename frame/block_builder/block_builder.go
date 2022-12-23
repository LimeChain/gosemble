/*
BlockBuilder - Version 4.
*/
package blockbuilder

type BlockBuilder interface {
	ApplyExtrinsic(dataPtr int32, dataLen int32) int64
	FinalizeBlock(dataPtr int32, dataLen int32) int64
	InherentExtrinisics(dataPtr int32, dataLen int32) int64
	CheckInherents(dataPtr int32, dataLen int32) int64
	RandomSeed(dataPtr int32, dataLen int32) int64
}

/*
	https://spec.polkadot.network/#sect-rte-apply-extrinsic

	SCALE encoded arguments (extrinsic types.Extrinsic) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([]byte) data.
*/
func ApplyExtrinsic(dataPtr int32, dataLen int32) int64

/*
	https://spec.polkadot.network/#defn-rt-blockbuilder-finalize-block

	SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded (types.Header) data.
*/
func FinalizeBlock(dataPtr int32, dataLen int32) int64

/*
	https://spec.polkadot.network/#defn-rt-builder-inherent-extrinsics

	SCALE encoded arguments (data types.InherentsData) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([]types.Extrinsic) data.
*/
func InherentExtrinisics(dataPtr int32, dataLen int32) int64

/*
	https://spec.polkadot.network/#id-blockbuilder_check_inherents

	SCALE encoded arguments (block types.Block, data types.InherentsData) allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([]byte) data.
*/
func CheckInherents(dataPtr int32, dataLen int32) int64

/*
	https://spec.polkadot.network/#id-blockbuilder_random_seed

	SCALE encoded arguments () allocated in the Wasm VM memory, passed as:
		dataPtr - i32 pointer to the memory location.
		dataLen - i32 length (in bytes) of the encoded arguments.
		returns a pointer-size to the SCALE-encoded ([32]byte) data.
*/
func RandomSeed(dataPtr int32, dataLen int32) int64
