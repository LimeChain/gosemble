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

func NewEventTransactionFeePaid(moduleIndex sc.U8, account types.AccountId[types.SignerAddress], actualFee types.Balance, tip types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventTransactionFeePaid, account, actualFee, tip)
}

func DecodeEvent[S types.SignerAddress](moduleIndex sc.U8, buffer *bytes.Buffer) (types.Event, error) {
	decodedModuleIndex, err := sc.DecodeU8(buffer)
	if err != nil {
		return types.Event{}, err
	}
	if decodedModuleIndex != moduleIndex {
		log.Critical("invalid transaction_payment.Event module")
	}

	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return types.Event{}, err
	}

	switch b {
	case EventTransactionFeePaid:
		account, err := types.DecodeAccountId[S](buffer)
		if err != nil {
			return types.Event{}, err
		}
		actualFee, err := sc.DecodeU128(buffer)
		if err != nil {
			return types.Event{}, err
		}
		tip, err := sc.DecodeU128(buffer)
		if err != nil {
			return types.Event{}, err
		}
		return NewEventTransactionFeePaid(moduleIndex, account, actualFee, tip), nil
	default:
		log.Critical("invalid transaction_payment.Event type")
	}

	panic("unreachable")
}
