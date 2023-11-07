package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
	"github.com/LimeChain/gosemble/primitives/types"
)

type accountMutator interface {
	ensureCanWithdraw(who types.Address32, amount sc.U128, reasons types.Reasons, newBalance sc.U128) types.DispatchError
	tryMutateAccountWithDust(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable]
	tryMutateAccount(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable]
}

type negativeImbalance struct {
	types.Balance
	totalIssuance support.StorageValue[sc.U128]
}

func newNegativeImbalance(balance types.Balance, totalIssuance support.StorageValue[sc.U128]) negativeImbalance {
	return negativeImbalance{balance, totalIssuance}
}

func (ni negativeImbalance) Drop() error {
	issuance, err := ni.totalIssuance.Get()
	if err != nil {
		return err
	}
	sub := sc.SaturatingSubU128(issuance, ni.Balance)

	ni.totalIssuance.Put(sub)
	return nil
}

type positiveImbalance struct {
	types.Balance
	totalIssuance support.StorageValue[sc.U128]
}

func newPositiveImbalance(balance types.Balance, totalIssuance support.StorageValue[sc.U128]) positiveImbalance {
	return positiveImbalance{balance, totalIssuance}
}

func (pi positiveImbalance) Drop() error {
	issuance, err := pi.totalIssuance.Get()
	if err != nil {
		return err
	}
	add := sc.SaturatingAddU128(issuance, pi.Balance)

	pi.totalIssuance.Put(add)
	return nil
}

type dustCleaner struct {
	moduleIndex       sc.U8
	accountId         types.Address32
	negativeImbalance sc.Option[negativeImbalance]
	eventDepositor    types.EventDepositor
}

func newDustCleaner(moduleId sc.U8, accountId types.Address32, negativeImbalance sc.Option[negativeImbalance], eventDepositor types.EventDepositor) dustCleaner {
	return dustCleaner{
		moduleIndex:       moduleId,
		accountId:         accountId,
		negativeImbalance: negativeImbalance,
		eventDepositor:    eventDepositor,
	}
}

func (dcv dustCleaner) Encode(buffer *bytes.Buffer) error {
	err := dcv.accountId.Encode(buffer)
	if err != nil {
		return err
	}
	return dcv.negativeImbalance.Encode(buffer)
}

func (dcv dustCleaner) Bytes() []byte {
	return sc.EncodedBytes(dcv)
}

func (dcv dustCleaner) Drop() {
	if dcv.negativeImbalance.HasValue {
		dcv.eventDepositor.DepositEvent(newEventDustLost(dcv.moduleIndex, dcv.accountId.FixedSequence, dcv.negativeImbalance.Value.Balance))
		dcv.negativeImbalance.Value.Drop()
	}
}
