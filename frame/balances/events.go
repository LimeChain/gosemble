package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/balances/types"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
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

const (
	errInvalidEventModule = "invalid balances.Event module"
	errInvalidEventType   = "invalid balances.Event type"
)

func newEventEndowed(moduleIndex sc.U8, account primitives.PublicKey, freeBalance primitives.Balance) primitives.Event {
	return primitives.NewEvent(moduleIndex, EventEndowed, account, freeBalance)
}

func newEventDustLost(moduleIndex sc.U8, account primitives.PublicKey, amount primitives.Balance) primitives.Event {
	return primitives.NewEvent(moduleIndex, EventDustLost, account, amount)
}

func newEventTransfer(moduleIndex sc.U8, from primitives.PublicKey, to primitives.PublicKey, amount primitives.Balance) primitives.Event {
	return primitives.NewEvent(moduleIndex, EventTransfer, from, to, amount)
}

func newEventBalanceSet(moduleIndex sc.U8, account primitives.PublicKey, free primitives.Balance, reserved primitives.Balance) primitives.Event {
	return primitives.NewEvent(moduleIndex, EventBalanceSet, account, free, reserved)
}

func newEventReserved(moduleIndex sc.U8, account primitives.PublicKey, amount primitives.Balance) primitives.Event {
	return primitives.NewEvent(moduleIndex, EventReserved, account, amount)
}

func newEventUnreserved(moduleIndex sc.U8, account primitives.PublicKey, amount primitives.Balance) primitives.Event {
	return primitives.NewEvent(moduleIndex, EventUnreserved, account, amount)
}

func newEventReserveRepatriated(moduleIndex sc.U8, from primitives.PublicKey, to primitives.PublicKey, amount primitives.Balance, destinationStatus types.BalanceStatus) primitives.Event {
	return primitives.NewEvent(moduleIndex, EventReserveRepatriated, from, to, amount, destinationStatus)
}

func newEventDeposit(moduleIndex sc.U8, account primitives.PublicKey, amount primitives.Balance) primitives.Event {
	return primitives.NewEvent(moduleIndex, EventDeposit, account, amount)
}

func newEventWithdraw(moduleIndex sc.U8, account primitives.PublicKey, amount primitives.Balance) primitives.Event {
	return primitives.NewEvent(moduleIndex, EventWithdraw, account, amount)
}

func newEventSlashed(moduleIndex sc.U8, account primitives.PublicKey, amount primitives.Balance) primitives.Event {
	return primitives.NewEvent(moduleIndex, EventSlashed, account, amount)
}

func DecodeEvent(moduleIndex sc.U8, buffer *bytes.Buffer) primitives.Event {
	decodedModuleIndex := sc.DecodeU8(buffer)
	if decodedModuleIndex != moduleIndex {
		log.Critical(errInvalidEventModule)
	}

	b := sc.DecodeU8(buffer)

	switch b {
	case EventEndowed:
		account := primitives.DecodePublicKey(buffer)
		freeBalance := sc.DecodeU128(buffer)
		return newEventEndowed(moduleIndex, account, freeBalance)
	case EventDustLost:
		account := primitives.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventDustLost(moduleIndex, account, amount)
	case EventTransfer:
		from := primitives.DecodePublicKey(buffer)
		to := primitives.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventTransfer(moduleIndex, from, to, amount)
	case EventBalanceSet:
		account := primitives.DecodePublicKey(buffer)
		free := sc.DecodeU128(buffer)
		reserved := sc.DecodeU128(buffer)
		return newEventBalanceSet(moduleIndex, account, free, reserved)
	case EventReserved:
		account := primitives.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventReserved(moduleIndex, account, amount)
	case EventUnreserved:
		account := primitives.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventUnreserved(moduleIndex, account, amount)
	case EventReserveRepatriated:
		from := primitives.DecodePublicKey(buffer)
		to := primitives.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		destinationStatus := types.DecodeBalanceStatus(buffer)
		return newEventReserveRepatriated(moduleIndex, from, to, amount, destinationStatus)
	case EventDeposit:
		account := primitives.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventDeposit(moduleIndex, account, amount)
	case EventWithdraw:
		account := primitives.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventWithdraw(moduleIndex, account, amount)
	case EventSlashed:
		account := primitives.DecodePublicKey(buffer)
		amount := sc.DecodeU128(buffer)
		return newEventSlashed(moduleIndex, account, amount)
	default:
		log.Critical(errInvalidEventType)
	}

	panic("unreachable")
}
