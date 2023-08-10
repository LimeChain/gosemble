package hooks

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type OnChargeTransaction interface {
	CorrectAndDepositFee(who *primitives.Address32, correctedFee primitives.Balance, tip primitives.Balance, alreadyWithdrawn sc.Option[primitives.Balance]) primitives.TransactionValidityError
	WithdrawFee(who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, fee primitives.Balance, tip primitives.Balance) (sc.Option[primitives.Balance], primitives.TransactionValidityError)
}
