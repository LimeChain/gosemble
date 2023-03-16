package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

func DepositEvent(event types.Event) {
	depositEventIndexed([]types.H256{}, event)
}

func depositEventIndexed(topics []types.H256, event types.Event) {
	blockNumber := StorageGetBlockNumber()
	if blockNumber == 0 {
		return
	}

	eventRecord := types.EventRecord{
		Phase:  StorageExecutionPhase(),
		Event:  event,
		Topics: topics,
	}

	oldEventCount := storageEventCount()
	newEventCount := oldEventCount + 1 // checked_add
	if newEventCount < oldEventCount {
		return
	}

	storageSetEventCount(newEventCount)

	storageAppendEvent(eventRecord)

	topicValue := sc.NewVaryingData(blockNumber, oldEventCount)
	for _, topic := range topics {
		storageAppendTopic(topic, topicValue)
	}
}
