package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

// Balances module events.
const (
	EventEndowed sc.U8 = iota
	EventDustLost
	EventTransfer
	EventBalanceSet
	EventReserved
	EventUnreserved
	EventReserveRepatriated
	EventDeposit
	EventWithdraw
	EventSlashed
)

func newEventEndowed(moduleIndex sc.U8, account types.PublicKey, freeBalance types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventEndowed, account, freeBalance)
}

func newEventDustLost(moduleIndex sc.U8, account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventDustLost, account, amount)
}

func newEventTransfer(moduleIndex sc.U8, from types.PublicKey, to types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventTransfer, from, to, amount)
}

func newEventBalanceSet(moduleIndex sc.U8, account types.PublicKey, free types.Balance, reserved types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventBalanceSet, account, free, reserved)
}

func newEventReserved(moduleIndex sc.U8, account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventReserved, account, amount)
}

func newEventUnreserved(moduleIndex sc.U8, account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventUnreserved, account, amount)
}

func newEventReserveRepatriated(moduleIndex sc.U8, from types.PublicKey, to types.PublicKey, amount types.Balance, destinationStatus types.BalanceStatus) types.Event {
	return types.NewEvent(moduleIndex, EventReserveRepatriated, from, to, amount, destinationStatus)
}

func newEventDeposit(moduleIndex sc.U8, account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventDeposit, account, amount)
}

func newEventWithdraw(moduleIndex sc.U8, account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventWithdraw, account, amount)
}

func newEventSlashed(moduleIndex sc.U8, account types.PublicKey, amount types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventSlashed, account, amount)
}

func DecodeEvent(moduleIndex sc.U8, buffer *bytes.Buffer) types.Event {
	decodedModuleIndex := sc.DecodeU8(buffer)
	if decodedModuleIndex != moduleIndex {
		log.Critical("invalid balances.Event module")
	}

	b := sc.DecodeU8(buffer)

	switch b {
	case EventEndowed:
		account := types.DecodePublicKey(buffer)
		freeBalance := sc.DecodeU128(buffer)
		return newEventEndowed(moduleIndex, account, freeBalance)
	case EventDustLost:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventDustLost(moduleIndex, account, amount)
	case EventTransfer:
		from := types.DecodePublicKey(buffer)
		to := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventTransfer(moduleIndex, from, to, amount)
	case EventBalanceSet:
		account := types.DecodePublicKey(buffer)
		free := sc.DecodeU128(buffer)
		reserved := sc.DecodeU128(buffer)
		return newEventBalanceSet(moduleIndex, account, free, reserved)
	case EventReserved:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventReserved(moduleIndex, account, amount)
	case EventUnreserved:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventUnreserved(moduleIndex, account, amount)
	case EventReserveRepatriated:
		from := types.DecodePublicKey(buffer)
		to := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		destinationStatus := types.DecodeBalanceStatus(buffer)
		return newEventReserveRepatriated(moduleIndex, from, to, amount, destinationStatus)
	case EventDeposit:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventDeposit(moduleIndex, account, amount)
	case EventWithdraw:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventWithdraw(moduleIndex, account, amount)
	case EventSlashed:
		account := types.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventSlashed(moduleIndex, account, amount)
	default:
		log.Critical("invalid balances.Event type")
	}

	panic("unreachable")
}
