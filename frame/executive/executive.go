package executive

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/types"
)

// InitializeBlock initialises a block with the given header,
// starting the execution of a particular block.
func InitializeBlock(header types.Header) {
	resetEvents()

	if runtimeUpgrade() {
		// TODO: weight
	}

	initialize(header.Number, header.ParentHash, extractPreRuntimeDigest(header.Digest))

	// TODO: weight

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
		value := append(
			sc.Compact(constants.RuntimeVersion.SpecVersion).Bytes(),
			constants.RuntimeVersion.SpecName.Bytes()...)
		storage.Set(key, value)

		return true
	}

	return false
}

func initialize(blockNumber types.BlockNumber, parentHash types.Blake2bHash, digest types.Digest) {
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
	prevBlock := blockNumber.U32 - 1
	blockNumHash := hashing.Twox64(prevBlock.Bytes())
	blockNumKey := append(systemHash, blockHashKeyHash...)
	blockNumKey = append(blockNumKey, blockNumHash...)
	blockNumKey = append(blockNumKey, prevBlock.Bytes()...)

	storage.Set(blockNumKey, parentHash.Bytes())

	blockWeightHash := hashing.Twox128(constants.KeyBlockWeight)
	storage.Clear(append(systemHash, blockWeightHash...))
}

func noteFinishedInitialize() {
	initializationPhase := sc.U32(constants.ExecutionPhaseApplyExtrinsic)

	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(systemHash, executionPhaseHash...), initializationPhase.Bytes())
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
