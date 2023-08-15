package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

// StorageGetBlockNumber returns the current block number being processed. Set by `execute_block`.
func StorageGetBlockNumber() types.BlockNumber {
	systemHash := hashing.Twox128(constants.KeySystem)
	numberHash := hashing.Twox128(constants.KeyNumber)
	return storage.GetDecode(append(systemHash, numberHash...), sc.DecodeU32)
}

func StorageGetAccount(who types.PublicKey) types.AccountInfo {
	systemHash := hashing.Twox128(constants.KeySystem)
	accountHash := hashing.Twox128(constants.KeyAccount)

	whoBytes := sc.FixedSequenceU8ToBytes(who)

	key := append(systemHash, accountHash...)
	key = append(key, hashing.Blake128(whoBytes)...)
	key = append(key, whoBytes...)

	return storage.GetDecode(key, types.DecodeAccountInfo)
}

func StorageSetAccount(who types.PublicKey, account types.AccountInfo) {
	systemHash := hashing.Twox128(constants.KeySystem)
	accountHash := hashing.Twox128(constants.KeyAccount)

	whoBytes := sc.FixedSequenceU8ToBytes(who)

	key := append(systemHash, accountHash...)
	key = append(key, hashing.Blake128(whoBytes)...)
	key = append(key, whoBytes...)

	storage.Set(key, account.Bytes())
}

func StorageExecutionPhase() types.ExtrinsicPhase {
	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	return storage.GetDecode(append(systemHash, executionPhaseHash...), types.DecodeExtrinsicPhase)
}

func storageEventCount() sc.U32 {
	systemHash := hashing.Twox128(constants.KeySystem)
	eventCountHash := hashing.Twox128(constants.KeyEventCount)

	key := append(systemHash, eventCountHash...)
	return storage.GetDecode(key, sc.DecodeU32)
}

func storageSetEventCount(eventCount sc.U32) {
	systemHash := hashing.Twox128(constants.KeySystem)
	eventCountHash := hashing.Twox128(constants.KeyEventCount)
	key := append(systemHash, eventCountHash...)
	storage.Set(key, eventCount.Bytes())
}

func storageAppendEvent(eventRecord types.EventRecord) {
	systemHash := hashing.Twox128(constants.KeySystem)
	key := append(systemHash, hashing.Twox128(constants.KeyEvents)...)
	storage.Append(key, eventRecord.Bytes())
}

func storageAppendTopic(topic types.H256, value sc.VaryingData) {
	systemHash := hashing.Twox128(constants.KeySystem)
	eventTopicsHash := hashing.Twox128(constants.KeyEventTopics)

	eventTopicsPrefix := append(systemHash, eventTopicsHash...)

	key := append(eventTopicsPrefix, topic.Bytes()...)
	storage.Append(key, value.Bytes())
}
