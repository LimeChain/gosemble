package module

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/primitives/types"
)

type accountMutator interface {
	ensureCanWithdraw(who types.Address32, amount *big.Int, reasons types.Reasons, newBalance *big.Int) types.DispatchError
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
	issuanceBn := issuance.ToBigInt()

	sub := new(big.Int).Sub(issuanceBn, ni.ToBigInt())

	if sub.Cmp(issuanceBn) > 0 {
		sub = issuanceBn
	}

	st.TotalIssuance.Put(sc.NewU128FromBigInt(sub))
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
	issuanceBn := issuance.ToBigInt()

	add := new(big.Int).Add(issuanceBn, pi.ToBigInt())

	if add.Cmp(issuanceBn) < 0 {
		add = issuanceBn
	}

	st.TotalIssuance.Put(sc.NewU128FromBigInt(add))
}

type dustCleanerValue struct {
	AccountId         types.Address32
	NegativeImbalance negativeImbalance
	eventDepositor    types.EventDepositor
}

func newDustCleanerValue(accountId types.Address32, negativeImbalance negativeImbalance, eventDepositor types.EventDepositor) dustCleanerValue {
	return dustCleanerValue{
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
	dcv.eventDepositor.DepositEvent(events.NewEventDustLost(dcv.AccountId.FixedSequence, dcv.NegativeImbalance.Balance))
	dcv.NegativeImbalance.Drop()
}
