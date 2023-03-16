package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

// The current block number being processed. Set by `execute_block`.
func StorageGetBlockNumber() types.BlockNumber {
	systemHash := hashing.Twox128(constants.KeySystem)
	numberHash := hashing.Twox128(constants.KeyNumber)
	return storage.GetDecode(append(systemHash, numberHash...), sc.DecodeU32)
}

// Total length (in bytes) for all extrinsics put together, for the current block.
func StorageGetAllExtrinsicsLen() sc.U32 {
	systemHash := hashing.Twox128(constants.KeySystem)
	allExtrinsicsLenHash := hashing.Twox128(constants.KeyAllExtrinsicsLen)
	return storage.GetDecode(append(systemHash, allExtrinsicsLenHash...), sc.DecodeU32)
}

func StorageSetAllExtrinsicsLen(length sc.U32) {
	systemHash := hashing.Twox128(constants.KeySystem)
	allExtrinsicsLenHash := hashing.Twox128(constants.KeyAllExtrinsicsLen)
	storage.Set(append(systemHash, allExtrinsicsLenHash...), length.Bytes())
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

// Map of block numbers to block hashes.
func StorageGetBlockHash(blockNumber sc.U32) types.Blake2bHash {
	// Module prefix
	systemHash := hashing.Twox128(constants.KeySystem)
	// Storage prefix
	blockHashHash := hashing.Twox128(constants.KeyBlockHash)
	// Block number hash
	blockNumHash := hashing.Twox64(blockNumber.Bytes())

	key := append(systemHash, blockHashHash...)
	key = append(key, blockNumHash...)
	key = append(key, blockNumber.Bytes()...)

	return storage.GetDecode(key, types.DecodeBlake2bHash)
}

// Map of block numbers to block hashes.
func StorageExistsBlockHash(blockNumber sc.U32) sc.Bool {
	// Module prefix
	systemHash := hashing.Twox128(constants.KeySystem)
	// Storage prefix
	blockHashHash := hashing.Twox128(constants.KeyBlockHash)
	// Block number hash
	blockNumHash := hashing.Twox64(blockNumber.Bytes())

	key := append(systemHash, blockHashHash...)
	key = append(key, blockNumHash...)
	key = append(key, blockNumber.Bytes()...)

	return storage.Exists(key) == 1
}

func StorageExecutionPhase() types.ExtrinsicPhase {
	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	return storage.GetDecode(append(systemHash, executionPhaseHash...), types.DecodeExtrinsicPhase)
}

func StorageSetExecutionPhase(phase types.ExtrinsicPhase) {
	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(systemHash, executionPhaseHash...), phase.Bytes())
}

func StorageClearExecutionPhase() {
	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Clear(append(systemHash, executionPhaseHash...))
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

// block weight
func StorageGetBlockWeight() types.ConsumedWeight {
	systemHash := hashing.Twox128(constants.KeySystem)
	blockWeightHash := hashing.Twox128(constants.KeyBlockWeight)
	return storage.GetDecode(append(systemHash, blockWeightHash...), types.DecodeConsumedWeight)
}

func StorageSetBlockWeight(weight types.ConsumedWeight) {
	systemHash := hashing.Twox128(constants.KeySystem)
	blockWeightHash := hashing.Twox128(constants.KeyBlockWeight)
	storage.Set(append(systemHash, blockWeightHash...), weight.Bytes())
}

func StorageClearBlockWeight() {
	systemHash := hashing.Twox128(constants.KeySystem)
	blockWeightHash := hashing.Twox128(constants.KeyBlockWeight)
	storage.Clear(append(systemHash, blockWeightHash...))
}

// Gets the index of extrinsic that is currently executing.
func StorageGetExtrinsicIndex() sc.U32 {
	return storage.GetDecode(constants.KeyExtrinsicIndex, sc.DecodeU32)
}
