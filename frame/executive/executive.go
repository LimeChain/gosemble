package executive

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/types"
)

// InitializeBlock initialises the block execution
func InitializeBlock(header types.Header) {
	resetEvents()

	if runtimeUpgrade() {
		// TODO:
	}

	initialize(header.Number, header.ParentHash, extractPreRuntimeDigest(header.Digest))

	// TODO:

	noteFinishedInitialize()
}

func resetEvents() {
	systemHash := hashing.Twox128(constants.KeySystem)
	eventsHash := hashing.Twox128(constants.KeyEvents)
	eventCountHash := hashing.Twox128(constants.KeyEventCount)
	eventTopicHash := hashing.Twox128(constants.KeyEventTopic)

	storage.Clear(append(systemHash, eventsHash...))
	storage.Clear(append(systemHash, eventCountHash...))
	storage.ClearPrefix(append(systemHash, eventTopicHash...))
}

func runtimeUpgrade() bool {
	systemHash := hashing.Twox128(constants.KeySystem)
	lastRuntimeUpgradeHash := hashing.Twox128(constants.KeyLastRuntimeUpgrade)

	key := append(systemHash, lastRuntimeUpgradeHash...)
	last := storage.Get(key)

	buffer := &bytes.Buffer{}
	buffer.Write(last)

	rupi, err := types.DecodeLastRuntimeUpgradeInfo(buffer)
	if err != nil {
		panic(err)
	}
	buffer.Reset()

	if constants.RuntimeVersion.SpecVersion > sc.U32(rupi.SpecVersion) || rupi.SpecName != constants.RuntimeVersion.SpecName {
		sc.Compact(constants.RuntimeVersion.SpecVersion).Encode(buffer)
		constants.RuntimeVersion.SpecName.Encode(buffer)
		storage.Set(key, buffer.Bytes())

		return true
	}

	return false
}

func initialize(blockNumber types.BlockNumber, parentHash types.Blake2bHash, digest types.Digest) {
	buffer := &bytes.Buffer{}
	initializationPhase := sc.U32(constants.ExecutionPhaseInitialization)
	initializationPhase.Encode(buffer)

	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(systemHash, executionPhaseHash...), buffer.Bytes())
	buffer.Reset()

	sc.U32(0).Encode(buffer)
	storage.Set(constants.KeyExtrinsicIndex, buffer.Bytes())
	buffer.Reset()

	blockNumber.Encode(buffer)
	numberHash := hashing.Twox128(constants.KeyNumber)
	storage.Set(append(systemHash, numberHash...), buffer.Bytes())
	buffer.Reset()

	digest.Encode(buffer)
	digestHash := hashing.Twox128(constants.KeyDigest)
	storage.Set(append(systemHash, digestHash...), buffer.Bytes())
	buffer.Reset()

	parentHashKey := hashing.Twox128(constants.KeyParentHash)
	parentHash.Encode(buffer)
	storage.Set(append(systemHash, parentHashKey...), buffer.Bytes())
	buffer.Reset()

	blockHashKeyHash := hashing.Twox128(constants.KeyBlockHash)
	prevBlock := blockNumber.U32 - 1
	prevBlock.Encode(buffer)
	blockNumHash := hashing.Twox64(buffer.Bytes())
	blockNumKey := append(systemHash, blockHashKeyHash...)
	blockNumKey = append(blockNumKey, blockNumHash...)
	blockNumKey = append(blockNumKey, buffer.Bytes()...)
	buffer.Reset()
	parentHash.Encode(buffer)
	storage.Set(blockNumKey, buffer.Bytes())
	buffer.Reset()

	blockWeightHash := hashing.Twox128(constants.KeyBlockWeight)
	storage.Clear(append(systemHash, blockWeightHash...))
}

func noteFinishedInitialize() {
	buffer := &bytes.Buffer{}
	initializationPhase := sc.U32(constants.ExecutionPhaseApplyExtrinsic)
	initializationPhase.Encode(buffer)

	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(systemHash, executionPhaseHash...), buffer.Bytes())
}

func extractPreRuntimeDigest(digest types.Digest) types.Digest {
	result := types.Digest{Values: map[uint8]sc.FixedSequence[types.DigestItem]{}}
	for k, v := range digest.Values {
		if k == types.DigestTypePreRuntime {
			result.Values[k] = v
		}
	}

	return result
}
