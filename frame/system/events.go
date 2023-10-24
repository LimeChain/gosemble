package system

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

// System module events.
const (
	EventExtrinsicSuccess sc.U8 = iota
	EventExtrinsicFailed
	EventCodeUpdated
	EventNewAccount
	EventKilledAccount
	EventRemarked
)

const (
	errInvalidEventModule = "invalid system.Event module"
	errInvalidEventType   = "invalid system.Event type"
)

func newEventExtrinsicSuccess(moduleIndex sc.U8, dispatchInfo types.DispatchInfo) types.Event {
	return types.NewEvent(moduleIndex, EventExtrinsicSuccess, dispatchInfo)
}

func newEventExtrinsicFailed(moduleIndex sc.U8, dispatchError types.DispatchError, dispatchInfo types.DispatchInfo) types.Event {
	return types.NewEvent(moduleIndex, EventExtrinsicFailed, dispatchError, dispatchInfo)
}

func newEventCodeUpdated(moduleIndex sc.U8) types.Event {
	return types.NewEvent(moduleIndex, EventCodeUpdated)
}

func newEventNewAccount(moduleIndex sc.U8, account types.PublicKey) types.Event {
	return types.NewEvent(moduleIndex, EventNewAccount, account)
}

func newEventKilledAccount(moduleIndex sc.U8, account types.PublicKey) types.Event {
	return types.NewEvent(moduleIndex, EventKilledAccount, account)
}

func newEventRemarked(moduleIndex sc.U8, sender types.PublicKey, hash types.H256) types.Event {
	return types.NewEvent(moduleIndex, EventRemarked, sender, hash)
}

func DecodeEvent(moduleIndex sc.U8, buffer *bytes.Buffer) (types.Event, error) {
	decodedModuleIndex, err := sc.DecodeU8(buffer)
	if err != nil {
		return types.Event{}, err
	}
	if decodedModuleIndex != moduleIndex {
		log.Critical(errInvalidEventModule)
	}

	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return types.Event{}, err
	}

	switch b {
	case EventExtrinsicSuccess:
		dispatchInfo := types.DecodeDispatchInfo(buffer)
		return newEventExtrinsicSuccess(moduleIndex, dispatchInfo), nil
	case EventExtrinsicFailed:
		dispatchErr := types.DecodeDispatchError(buffer)
		dispatchInfo := types.DecodeDispatchInfo(buffer)
		return newEventExtrinsicFailed(moduleIndex, dispatchErr, dispatchInfo), nil
	case EventCodeUpdated:
		return newEventCodeUpdated(moduleIndex), nil
	case EventNewAccount:
		account := types.DecodePublicKey(buffer)
		return newEventNewAccount(moduleIndex, account), nil
	case EventKilledAccount:
		account := types.DecodePublicKey(buffer)
		return newEventKilledAccount(moduleIndex, account), nil
	case EventRemarked:
		account := types.DecodePublicKey(buffer)
		hash := types.DecodeH256(buffer)
		return newEventRemarked(moduleIndex, account, hash), nil
	default:
		log.Critical(errInvalidEventType)
	}

	panic("unreachable")
}
