package system

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/system"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

const (
	EventExtrinsicSuccess sc.U8 = iota
	EventExtrinsicFailed
	EventCodeUpdated
	EventNewAccount
	EventKilledAccount
	EventRemarked
)

type event sc.VaryingData

func (e event) Encode(buffer *bytes.Buffer) {
	if len(e) < 2 {
		log.Critical("cannot system.Event encode")
	}

	switch e[1] {
	case EventExtrinsicSuccess:
		e[0].Encode(buffer)
		e[1].Encode(buffer)
		e[2].Encode(buffer)
	case EventExtrinsicFailed:
		e[0].Encode(buffer)
		e[1].Encode(buffer)
		e[2].Encode(buffer)
		e[3].Encode(buffer)
	case EventCodeUpdated:
		e[0].Encode(buffer)
		e[1].Encode(buffer)
	case EventNewAccount:
		e[0].Encode(buffer)
		e[1].Encode(buffer)
		e[2].Encode(buffer)
	case EventKilledAccount:
		e[0].Encode(buffer)
		e[1].Encode(buffer)
		e[2].Encode(buffer)
	case EventRemarked:
		e[0].Encode(buffer)
		e[1].Encode(buffer)
		e[2].Encode(buffer)
		e[3].Encode(buffer)
	default:
		log.Critical("invalid system.Event type")
	}
}

func (e event) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func NewEventExtrinsicSuccess(dispatchInfo types.DispatchInfo) event {
	return event{system.ModuleIndex, EventCodeUpdated, dispatchInfo}
}

func NewEventExtrinsicFailed(dispatchError types.DispatchError, dispatchInfo types.DispatchInfo) types.Event {
	return event{system.ModuleIndex, EventCodeUpdated, dispatchError, dispatchInfo}
}

func NewEventCodeUpdated() types.Event {
	return event{system.ModuleIndex, EventCodeUpdated}
}

func NewEventNewAccount(account types.PublicKey) types.Event {
	return event{system.ModuleIndex, EventNewAccount, account}
}

func NewEventKilledAccount(account types.PublicKey) types.Event {
	return event{system.ModuleIndex, EventKilledAccount, account}
}

func NewEventRemarked(sender types.PublicKey, hash types.H256) types.Event {
	return event{system.ModuleIndex, EventRemarked, sender, hash}
}

func DecodeEvent(buffer *bytes.Buffer) types.Event {
	module := sc.DecodeU8(buffer)
	if module != system.ModuleIndex {
		log.Critical("invalid system.Event")
	}

	b := sc.DecodeU8(buffer)

	switch b {
	case EventExtrinsicSuccess:
		dispatchInfo := types.DecodeDispatchInfo(buffer)
		return NewEventExtrinsicSuccess(dispatchInfo)
	case EventExtrinsicFailed:
		dispatchErr := types.DecodeDispatchError(buffer)
		dispatchInfo := types.DecodeDispatchInfo(buffer)
		return NewEventExtrinsicFailed(dispatchErr, dispatchInfo)
	case EventCodeUpdated:
		return NewEventCodeUpdated()
	case EventNewAccount:
		account := types.DecodePublicKey(buffer)
		return NewEventNewAccount(account)
	case EventKilledAccount:
		account := types.DecodePublicKey(buffer)
		return NewEventKilledAccount(account)
	case EventRemarked:
		account := types.DecodePublicKey(buffer)
		hash := types.DecodeH256(buffer)
		return NewEventRemarked(account, hash)
	default:
		log.Critical("invalid system.Event type")
	}

	panic("unreachable")
}
