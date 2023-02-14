package system

import (
	"bytes"
	"math"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

const (
	ModuleIndex   = 0
	FunctionIndex = 0
)

func Remark(args sc.Sequence[sc.U8]) {
	// TODO:
}

func Finalize() types.Header {
	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Clear(append(systemHash, executionPhaseHash...))

	allExtrinsicsLenHash := hashing.Twox128(constants.KeyAllExtrinsicsLen)
	storage.Clear(append(systemHash, allExtrinsicsLenHash...))

	numberHash := hashing.Twox128(constants.KeyNumber)
	b := storage.Get(append(systemHash, numberHash...))
	buf := &bytes.Buffer{}
	buf.Write(b)
	blockNumber := sc.DecodeU32(buf)

	parentHashKey := hashing.Twox128(constants.KeyParentHash)
	b = storage.Get(append(systemHash, parentHashKey...))
	buf.Reset()
	buf.Write(b)

	parentHash := sc.DecodeFixedSequence[sc.U8](32, buf)

	digestHash := hashing.Twox128(constants.KeyDigest)
	b = storage.Get(append(systemHash, digestHash...))
	buf.Reset()
	buf.Write(b)

	digest := types.DecodeDigest(buf)
	buf.Reset()

	extrinsicCountHash := hashing.Twox128(constants.KeyExtrinsicCount)

	extrinsicCount := sc.U32(0)
	b = storage.Get(append(systemHash, extrinsicCountHash...))
	if len(b) > 0 {
		buf.Write(b)
		extrinsicCount = sc.DecodeU32(buf)
		buf.Reset()
	}

	extrinsicDataPrefixHash := append(systemHash, hashing.Twox128(constants.KeyExtrinsicData)...)

	extrinsics := storage.Get(append(extrinsicDataPrefixHash, hashing.Twox128(extrinsicCount.Bytes())...))

	extrinsicsRootBytes := hashing.Blake256(extrinsics)
	buf.Write(extrinsicsRootBytes)
	extrinsicsRoot := types.DecodeH256(buf)
	buf.Reset()

	blockHashCountBytes := storage.Get(append(systemHash, hashing.Twox128(constants.KeyBlockHashCount)...))

	buf.Write(blockHashCountBytes)
	blockHashCount := sc.DecodeU32(buf)
	buf.Reset()

	toRemove := blockNumber - blockHashCount
	if toRemove != 0 {
		blockNumHash := hashing.Twox64(toRemove.Bytes())
		blockNumKey := append(systemHash, hashing.Twox128(constants.KeyBlockHash)...)
		blockNumKey = append(blockNumKey, blockNumHash...)
		blockNumKey = append(blockNumKey, toRemove.Bytes()...)

		storage.Clear(blockNumKey)
	}

	storageRootBytes := storage.Root(constants.RuntimeVersion.StateVersion.Bytes())
	buf.Write(storageRootBytes)
	storageRoot := types.DecodeH256(buf)
	buf.Reset()

	return types.Header{
		ExtrinsicsRoot: extrinsicsRoot,
		StateRoot:      storageRoot,
		ParentHash:     types.Blake2bHash{FixedSequence: parentHash},
		Number:         blockNumber,
		Digest:         digest,
	}
}

func Initialize(blockNumber types.BlockNumber, parentHash types.Blake2bHash, digest types.Digest) {
	initializationPhase := sc.U32(constants.ExecutionPhaseInitialization)

	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(systemHash, executionPhaseHash...), initializationPhase.Bytes())

	storage.Set(constants.KeyExtrinsicIndex, sc.U32(0).Bytes())

	numberHash := hashing.Twox128(constants.KeyNumber)
	storage.Set(append(systemHash, numberHash...), blockNumber.Bytes())

	digestHash := hashing.Twox128(constants.KeyDigest)
	storage.Set(append(systemHash, digestHash...), digest.Bytes())

	parentHashKey := hashing.Twox128(constants.KeyParentHash)
	storage.Set(append(systemHash, parentHashKey...), parentHash.Bytes())

	blockHashKeyHash := hashing.Twox128(constants.KeyBlockHash)
	prevBlock := blockNumber - 1
	blockNumHash := hashing.Twox64(prevBlock.Bytes())
	blockNumKey := append(systemHash, blockHashKeyHash...)
	blockNumKey = append(blockNumKey, blockNumHash...)
	blockNumKey = append(blockNumKey, prevBlock.Bytes()...)

	storage.Set(blockNumKey, parentHash.Bytes())

	blockWeightHash := hashing.Twox128(constants.KeyBlockWeight)
	storage.Clear(append(systemHash, blockWeightHash...))
}

func NoteFinishedInitialize() {
	initializationPhase := sc.U32(constants.ExecutionPhaseApplyExtrinsic)

	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(systemHash, executionPhaseHash...), initializationPhase.Bytes())
}

func NoteFinishedExtrinsics() {
	value := storage.Get(constants.KeyExtrinsicIndex)
	extrinsicIndex := sc.U32(0)

	if len(value) > 1 {
		storage.Clear(constants.KeyExtrinsicIndex)
		buf := &bytes.Buffer{}
		buf.Write(value)

		extrinsicIndex = sc.DecodeU32(buf)
	}

	systemHash := hashing.Twox128(constants.KeySystem)
	extrinsicCountHash := hashing.Twox128(constants.KeyExtrinsicCount)

	storage.Set(append(systemHash, extrinsicCountHash...), extrinsicIndex.Bytes())

	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	finalizationPhase := sc.U32(constants.ExecutionPhaseInitialization)

	storage.Set(append(systemHash, executionPhaseHash...), finalizationPhase.Bytes())
}

func ResetEvents() {
	systemHash := hashing.Twox128(constants.KeySystem)
	eventsHash := hashing.Twox128(constants.KeyEvents)
	eventCountHash := hashing.Twox128(constants.KeyEventCount)
	eventTopicHash := hashing.Twox128(constants.KeyEventTopic)

	storage.Clear(append(systemHash, eventsHash...))
	storage.Clear(append(systemHash, eventCountHash...))

	limit := sc.Option[sc.U32]{
		HasValue: true,
		Value:    sc.U32(math.MaxUint32),
	}
	storage.ClearPrefix(append(systemHash, eventTopicHash...), limit.Bytes())
}

// Note what the extrinsic data of the current extrinsic index is.
//
// This is required to be called before applying an extrinsic. The data will used
// in [`finalize`] to calculate the correct extrinsics root.
func NoteExtrinsic(encodedExt []byte) {
	keySystemHash := hashing.Twox128(constants.KeySystem)
	storage.Set(append(keySystemHash, extrinsicIndexValue().Bytes()...), encodedExt)
}

// To be called immediately after an extrinsic has been applied.
//
// Emits an `ExtrinsicSuccess` or `ExtrinsicFailed` event depending on the outcome.
// The emitted event contains the post-dispatch corrected weight including
// the base-weight for its dispatch class.
func NoteAppliedExtrinsic(r *types.DispatchResultWithPostInfo[types.PostDispatchInfo], info types.DispatchInfo) {
	// TODO:
	// info.Weight = extract_actual_weight(r, &info).saturating_add(T::BlockWeights::get().get(info.class).base_extrinsic)
	// info.PaysFee = extract_actual_pays_fee(r, &info)

	// Self::deposit_event(match r {
	// 	Ok(_) => Event::ExtrinsicSuccess { dispatch_info: info },
	// 	Err(err) => {
	// 		log::trace!(
	// 			target: LOG_TARGET,
	// 			"Extrinsic failed at block({:?}): {:?}",
	// 			Self::block_number(),
	// 			err,
	// 		);
	// 		Event::ExtrinsicFailed { dispatch_error: err.error, dispatch_info: info }
	// 	},
	// });

	nextExtrinsicIndex := extrinsicIndexValue().Value + sc.U32(1)

	keySystemHash := hashing.Twox128(constants.KeySystem)

	keyExtrinsicIndex := hashing.Twox128(constants.KeyExtrinsicIndex)
	storage.Set(append(keySystemHash, keyExtrinsicIndex...), nextExtrinsicIndex.Bytes())

	keyExecutionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(keySystemHash, keyExecutionPhaseHash...), (types.NewPhase(types.PhaseApplyExtrinsic, nextExtrinsicIndex)).Bytes())
}

// Gets the index of extrinsic that is currently executing.
func extrinsicIndexValue() sc.Option[sc.U32] {
	keySystemHash := hashing.Twox128(constants.KeySystem)
	keyExtrinsicIndex := hashing.Twox128(constants.KeyExtrinsicIndex)
	value := storage.Get(append(keySystemHash, keyExtrinsicIndex...))

	if len(value) != 0 {
		buf := &bytes.Buffer{}
		buf.Write(value)
		return sc.Option[sc.U32]{HasValue: true, Value: sc.U32(sc.DecodeU8(buf))}
	} else {
		return sc.Option[sc.U32]{HasValue: false}
	}
}
