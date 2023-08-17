package transaction_payment

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

// TransactionPayment module events.
const (
	EventTransactionFeePaid sc.U8 = iota
)

func NewEventTransactionFeePaid(moduleIndex sc.U8, account types.PublicKey, actualFee types.Balance, tip types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventTransactionFeePaid, account, actualFee, tip)
}

func DecodeEvent(moduleIndex sc.U8, buffer *bytes.Buffer) types.Event {
	decodedModuleIndex := sc.DecodeU8(buffer)
	if decodedModuleIndex != moduleIndex {
		log.Critical("invalid transaction_payment.Event module")
	}

	b := sc.DecodeU8(buffer)

	switch b {
	case EventTransactionFeePaid:
		account := types.DecodePublicKey(buffer)
		actualFee := sc.DecodeU128(buffer)
		tip := sc.DecodeU128(buffer)
		return NewEventTransactionFeePaid(moduleIndex, account, actualFee, tip)
	default:
		log.Critical("invalid transaction_payment.Event type")
	}

	panic("unreachable")
}
