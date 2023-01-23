/*
BlockBuilder - Version 4.
*/
package blockbuilder

import (
	"bytes"
	"github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
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
func ApplyExtrinsic(dataPtr int32, dataLen int32) int64 { return 0 }

// FinalizeBlock finalizes block - it is up the caller to ensure that all header fields are valid
// except state-root.
func FinalizeBlock(dataPtr int32, dataLen int32) int64 {
	noteFinishedExtrinsics()

	systemHash := hashing.Twox128(constants.KeySystem)
	numberHash := hashing.Twox128(constants.KeyNumber)

	bNumber := storage.Get(append(systemHash, numberHash...))
	buf := &bytes.Buffer{}
	buf.Write(bNumber)
	blockNumber := goscale.DecodeU32(buf)

	idleAndFinalizeHook(types.BlockNumber{U32: blockNumber})

	header := finalize()
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

func noteFinishedExtrinsics() {
	value := storage.Get(constants.KeyExtrinsicIndex)
	extrinsicIndex := goscale.U32(0)

	if len(value) > 1 {
		storage.Clear(constants.KeyExtrinsicIndex)
		buf := &bytes.Buffer{}
		buf.Write(value)

		extrinsicIndex = goscale.DecodeU32(buf)
	}

	systemHash := hashing.Twox128(constants.KeySystem)
	extrinsicCountHash := hashing.Twox128(constants.KeyExtrinsicCount)

	storage.Set(append(systemHash, extrinsicCountHash...), extrinsicIndex.Bytes())

	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	finalizationPhase := goscale.U32(constants.ExecutionPhaseInitialization)

	storage.Set(append(systemHash, executionPhaseHash...), finalizationPhase.Bytes())
}

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

func finalize() types.Header {
	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Clear(append(systemHash, executionPhaseHash...))

	allExtrinsicsLenHash := hashing.Twox128(constants.KeyAllExtrinsicsLen)
	storage.Clear(append(systemHash, allExtrinsicsLenHash...))

	numberHash := hashing.Twox128(constants.KeyNumber)
	b := storage.Get(append(systemHash, numberHash...))
	buf := &bytes.Buffer{}
	buf.Write(b)
	blockNumber := goscale.DecodeU32(buf)

	parentHashKey := hashing.Twox128(constants.KeyParentHash)
	b = storage.Get(append(systemHash, parentHashKey...))
	buf.Reset()
	buf.Write(b)

	parentHash := goscale.DecodeFixedSequence[goscale.U8](32, buf)

	digestHash := hashing.Twox128(constants.KeyDigest)
	b = storage.Get(append(systemHash, digestHash...))
	buf.Reset()
	buf.Write(b)

	digest := types.DecodeDigest(buf)
	buf.Reset()

	extrinsicCountHash := hashing.Twox128(constants.KeyExtrinsicCount)

	extrinsicCount := goscale.U32(0)
	b = storage.Get(append(systemHash, extrinsicCountHash...))
	if len(b) > 0 {
		buf.Write(b)
		extrinsicCount = goscale.DecodeU32(buf)
		buf.Reset()
	}

	extrinsicDataPrefixHash := append(systemHash, hashing.Twox128(constants.KeyExtrinsicData)...)

	extrinsics := storage.Get(append(extrinsicDataPrefixHash, hashing.Twox128(extrinsicCount.Bytes())...))

	extrinsicsRootBytes := hashing.Blake256(extrinsics)
	buf.Write(extrinsicsRootBytes)
	extrinsicsRoot := goscale.DecodeFixedSequence[goscale.U8](32, buf)
	buf.Reset()

	blockHashCountBytes := storage.Get(append(systemHash, hashing.Twox128(constants.KeyBlockHashCount)...))

	buf.Write(blockHashCountBytes)
	blockHashCount := goscale.DecodeU32(buf)
	buf.Reset()

	toRemove := blockNumber - blockHashCount
	if toRemove != 0 {
		blockNumHash := hashing.Twox64(toRemove.Bytes())
		blockNumKey := append(systemHash, hashing.Twox128(constants.KeyBlockHash)...)
		blockNumKey = append(blockNumKey, blockNumHash...)
		blockNumKey = append(blockNumKey, toRemove.Bytes()...)

		storage.Clear(blockNumKey)
	}

	storageRootBytes := storage.RootV2(constants.RuntimeVersion.StateVersion.Bytes())
	buf.Write(storageRootBytes)
	storageRoot := goscale.DecodeFixedSequence[goscale.U8](32, buf)
	buf.Reset()

	return types.Header{
		ExtrinsicsRoot: extrinsicsRoot,
		StateRoot:      storageRoot,
		ParentHash:     types.Blake2bHash{FixedSequence: parentHash},
		Number:         types.BlockNumber{U32: blockNumber},
		Digest:         digest,
	}
}
