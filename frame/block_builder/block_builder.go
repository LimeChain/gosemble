/*
BlockBuilder - Version 4.
*/
package blockbuilder

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

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
func ApplyExtrinsic(dataPtr int32, dataLen int32) int64 {
	buffer := &bytes.Buffer{}
	buffer.Write(utils.ToWasmMemorySlice(dataPtr, dataLen))
	uxt := types.DecodeUncheckedExtrinsic(buffer)

	applyExtrinsicResult := executive.ApplyExtrinsic(uxt)

	buffer.Reset()
	applyExtrinsicResult.Encode(buffer)

	return utils.BytesToOffsetAndSize(buffer.Bytes())
}

/*
https://spec.polkadot.network/#defn-rt-blockbuilder-finalize-block

SCALE encoded arguments () allocated in the Wasm VM memory, passed as:

	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
	returns a pointer-size to the SCALE-encoded (types.Header) data.
*/

// FinalizeBlock finalizes block - it is up the caller to ensure that all header fields are valid
// except state-root.
func FinalizeBlock(dataPtr int32, dataLen int32) int64 {
	system.NoteFinishedExtrinsics()

	systemHash := hashing.Twox128(constants.KeySystem)
	numberHash := hashing.Twox128(constants.KeyNumber)

	bNumber := storage.Get(append(systemHash, numberHash...))
	buf := &bytes.Buffer{}
	buf.Write(bNumber)
	blockNumber := sc.DecodeU32(buf)

	idleAndFinalizeHook(types.BlockNumber{U32: blockNumber})

	header := system.Finalize()
	encodedHeader := header.Bytes()

	return utils.BytesToOffsetAndSize(encodedHeader)
}

/*
https://spec.polkadot.network/#defn-rt-builder-inherent-extrinsics

SCALE encoded arguments (data types.InherentsData) allocated in the Wasm VM memory, passed as:

	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
	returns a pointer-size to the SCALE-encoded ([]types.Extrinsic) data.
*/
func InherentExtrinisics(dataPtr int32, dataLen int32) int64 { return 0 }

/*
https://spec.polkadot.network/#id-blockbuilder_check_inherents

SCALE encoded arguments (block types.Block, data types.InherentsData) allocated in the Wasm VM memory, passed as:

	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
	returns a pointer-size to the SCALE-encoded ([]byte) data.
*/
func CheckInherents(dataPtr int32, dataLen int32) int64 {
	return 0
}

/*
https://spec.polkadot.network/#id-blockbuilder_random_seed

SCALE encoded arguments () allocated in the Wasm VM memory, passed as:

	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
	returns a pointer-size to the SCALE-encoded ([32]byte) data.
*/
func RandomSeed(dataPtr int32, dataLen int32) int64 { return 0 }

func idleAndFinalizeHook(blockNumber types.BlockNumber) {
	systemHash := hashing.Twox128(constants.KeySystem)
	blockWeightHash := hashing.Twox128(constants.KeyBlockWeight)

	storage.Get(append(systemHash, blockWeightHash...))

	// TODO: weights
	/**
	let weight = <frame_system::Pallet<System>>::block_weight();
	let max_weight = <System::BlockWeights as frame_support::traits::Get<_>>::get().max_block;
	let remaining_weight = max_weight.saturating_sub(weight.total());

	if remaining_weight.all_gt(Weight::zero()) {
		let used_weight = <AllPalletsWithSystem as OnIdle<System::BlockNumber>>::on_idle(
			block_number,
			remaining_weight,
		);
		<frame_system::Pallet<System>>::register_extra_weight_unchecked(
			used_weight,
			DispatchClass::Mandatory,
		);
	}
	// Each pallet (babe, grandpa) has its own on_finalize that has to be implemented once it is supported
	<AllPalletsWithSystem as OnFinalize<System::BlockNumber>>::on_finalize(block_number);
	*/
}
