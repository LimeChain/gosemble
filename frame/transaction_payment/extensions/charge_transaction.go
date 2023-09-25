package extensions

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type chargeTransaction struct {
	currencyAdapter primitives.CurrencyAdapter
}

func newChargeTransaction(currencyAdapter primitives.CurrencyAdapter) chargeTransaction {
	return chargeTransaction{currencyAdapter: currencyAdapter}
}

func (ct chargeTransaction) WithdrawFee(who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, fee primitives.Balance, tip primitives.Balance) (sc.Option[primitives.Balance], primitives.TransactionValidityError) {
	if fee.Eq(constants.Zero) {
		return sc.NewOption[primitives.Balance](nil), nil
	}

	withdrawReasons := primitives.WithdrawReasonsTransactionPayment
	if tip.Eq(constants.Zero) {
		withdrawReasons = primitives.WithdrawReasonsTransactionPayment
	} else {
		withdrawReasons = primitives.WithdrawReasonsTransactionPayment | primitives.WithdrawReasonsTip
	}

	imbalance, err := ct.currencyAdapter.Withdraw(*who, fee, sc.U8(withdrawReasons), primitives.ExistenceRequirementKeepAlive)
	if err != nil {
		return sc.NewOption[primitives.Balance](nil), primitives.NewTransactionValidityError(primitives.NewInvalidTransactionPayment())
	}

	return sc.NewOption[primitives.Balance](imbalance), nil
}

func (ct chargeTransaction) CorrectAndDepositFee(who *primitives.Address32, correctedFee primitives.Balance, tip primitives.Balance, alreadyWithdrawn sc.Option[primitives.Balance]) primitives.TransactionValidityError {
	if alreadyWithdrawn.HasValue {
		alreadyPaidNegativeImbalance := alreadyWithdrawn.Value
		refundAmount := alreadyPaidNegativeImbalance.Sub(correctedFee)

		refundPositiveImbalance, err := ct.currencyAdapter.DepositIntoExisting(*who, refundAmount)
		if err != nil {
			return primitives.NewTransactionValidityError(primitives.NewInvalidTransactionPayment())
		}

		if alreadyPaidNegativeImbalance.Lt(refundPositiveImbalance) {
			return primitives.NewTransactionValidityError(primitives.NewInvalidTransactionPayment())
		}
	}
	return nil
}
