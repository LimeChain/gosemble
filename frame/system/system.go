package system

import (
	"bytes"
	"github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
	"math"
)

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

	storageRootBytes := storage.Root(constants.RuntimeVersion.StateVersion.Bytes())
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

func Initialize(blockNumber types.BlockNumber, parentHash types.Blake2bHash, digest types.Digest) {
	initializationPhase := goscale.U32(constants.ExecutionPhaseInitialization)

	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(systemHash, executionPhaseHash...), initializationPhase.Bytes())

	storage.Set(constants.KeyExtrinsicIndex, goscale.U32(0).Bytes())

	numberHash := hashing.Twox128(constants.KeyNumber)
	storage.Set(append(systemHash, numberHash...), blockNumber.Bytes())

	digestHash := hashing.Twox128(constants.KeyDigest)
	storage.Set(append(systemHash, digestHash...), digest.Bytes())

	parentHashKey := hashing.Twox128(constants.KeyParentHash)
	storage.Set(append(systemHash, parentHashKey...), parentHash.Bytes())

	blockHashKeyHash := hashing.Twox128(constants.KeyBlockHash)
	prevBlock := blockNumber.U32 - 1
	blockNumHash := hashing.Twox64(prevBlock.Bytes())
	blockNumKey := append(systemHash, blockHashKeyHash...)
	blockNumKey = append(blockNumKey, blockNumHash...)
	blockNumKey = append(blockNumKey, prevBlock.Bytes()...)

	storage.Set(blockNumKey, parentHash.Bytes())

	blockWeightHash := hashing.Twox128(constants.KeyBlockWeight)
	storage.Clear(append(systemHash, blockWeightHash...))
}

func NoteFinishedInitialize() {
	initializationPhase := goscale.U32(constants.ExecutionPhaseApplyExtrinsic)

	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(systemHash, executionPhaseHash...), initializationPhase.Bytes())
}

func NoteFinishedExtrinsics() {
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

func ResetEvents() {
	systemHash := hashing.Twox128(constants.KeySystem)
	eventsHash := hashing.Twox128(constants.KeyEvents)
	eventCountHash := hashing.Twox128(constants.KeyEventCount)
	eventTopicHash := hashing.Twox128(constants.KeyEventTopic)

	storage.Clear(append(systemHash, eventsHash...))
	storage.Clear(append(systemHash, eventCountHash...))

	limit := goscale.Option[goscale.U32]{
		HasValue: true,
		Value:    goscale.U32(math.MaxUint32),
	}
	storage.ClearPrefix(append(systemHash, eventTopicHash...), limit.Bytes())
}
