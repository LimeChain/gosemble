package transaction_payment

import (
	"bytes"
	"errors"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

// TransactionPayment module events.
const (
	EventTransactionFeePaid sc.U8 = iota
)

var (
	errInvalidType   = errors.New("invalid transaction_payment.Event type")
	errInvalidModule = errors.New("invalid transaction_payment.Event module")
)

func NewEventTransactionFeePaid(moduleIndex sc.U8, account types.AccountId, actualFee types.Balance, tip types.Balance) types.Event {
	return types.NewEvent(moduleIndex, EventTransactionFeePaid, account, actualFee, tip)
}

func DecodeEvent(moduleIndex sc.U8, buffer *bytes.Buffer) (types.Event, error) {
	decodedModuleIndex, err := sc.DecodeU8(buffer)
	if err != nil {
		return types.Event{}, err
	}
	if decodedModuleIndex != moduleIndex {
		return types.Event{}, errInvalidModule
	}

	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return types.Event{}, err
	}

	switch b {
	case EventTransactionFeePaid:
		account, err := types.DecodeAccountId(buffer)
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
		return types.Event{}, errInvalidType
	}
}
