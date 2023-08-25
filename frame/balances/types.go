package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type accountMutator interface {
	ensureCanWithdraw(who types.Address32, amount sc.U128, reasons types.Reasons, newBalance sc.U128) types.DispatchError
	tryMutateAccountWithDust(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable]
	tryMutateAccount(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable]
}

type negativeImbalance struct {
	types.Balance
}

func newNegativeImbalance(balance types.Balance) negativeImbalance {
	return negativeImbalance{balance}
}

func (ni negativeImbalance) Drop() {
	st := newStorage() // TODO: revise
	issuance := st.TotalIssuance.Get()

	sub := issuance.Sub(ni)
	if sub.Gt(issuance) {
		sub = issuance
	}

	st.TotalIssuance.Put(sub.(sc.U128))
}

type positiveImbalance struct {
	types.Balance
}

func newPositiveImbalance(balance types.Balance) positiveImbalance {
	return positiveImbalance{balance}
}

func (pi positiveImbalance) Drop() {
	st := newStorage() // TODO: revise
	issuance := st.TotalIssuance.Get()

	add := issuance.Add(pi)
	if add.Lt(issuance) {
		add = issuance
	}

	st.TotalIssuance.Put(add.(sc.U128))
}

type dustCleanerValue struct {
	moduleIndex       sc.U8
	AccountId         types.Address32
	NegativeImbalance negativeImbalance
	eventDepositor    types.EventDepositor
}

func newDustCleanerValue(moduleId sc.U8, accountId types.Address32, negativeImbalance negativeImbalance, eventDepositor types.EventDepositor) dustCleanerValue {
	return dustCleanerValue{
		moduleIndex:       moduleId,
		AccountId:         accountId,
		NegativeImbalance: negativeImbalance,
		eventDepositor:    eventDepositor,
	}
}

func (dcv dustCleanerValue) Encode(buffer *bytes.Buffer) {
	dcv.AccountId.Encode(buffer)
	dcv.NegativeImbalance.Encode(buffer)
}

func (dcv dustCleanerValue) Bytes() []byte {
	return sc.EncodedBytes(dcv)
}

func (dcv dustCleanerValue) Drop() {
	dcv.eventDepositor.DepositEvent(newEventDustLost(dcv.moduleIndex, dcv.AccountId.FixedSequence, dcv.NegativeImbalance.Balance))
	dcv.NegativeImbalance.Drop()
}
