package frame

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

func DepositEvent(event types.Event, blockNumber types.BlockNumber) {
	depositEventIndexed([]types.H256{}, event, blockNumber)
}

func depositEventIndexed(topics []types.H256, event types.Event, blockNumber types.BlockNumber) {
	if blockNumber == 0 {
		return
	}

	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)

	phase := storage.GetDecode(append(systemHash, executionPhaseHash...), types.DecodeExtrinsicPhase)

	eventRecord := types.EventRecord{
		Phase:  phase,
		Event:  event,
		Topics: topics,
	}

	eventCountHash := hashing.Twox128(constants.KeyEventCount)

	oldEventCount := storage.GetDecode(append(systemHash, eventCountHash...), sc.DecodeU32)
	newEventCount := oldEventCount + 1 // checked_add
	if newEventCount < oldEventCount {
		return
	}

	storage.Set(append(systemHash, eventCountHash...), newEventCount.Bytes())

	storage.Append(append(systemHash, hashing.Twox128(constants.KeyEvents)...), eventRecord.Bytes())

	eventTopicsHash := hashing.Twox128(constants.KeyEventTopics)
	eventTopicsPrefix := append(systemHash, eventTopicsHash...)

	for _, topic := range topics {
		storage.Append(append(eventTopicsPrefix, topic.Bytes()...), sc.NewVaryingData(blockNumber, oldEventCount).Bytes())
	}
}
