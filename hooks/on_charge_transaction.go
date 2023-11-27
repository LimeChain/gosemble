package hooks

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type OnChargeTransaction interface {
	CorrectAndDepositFee(who primitives.AccountId, correctedFee primitives.Balance, tip primitives.Balance, alreadyWithdrawn sc.Option[primitives.Balance]) error
	WithdrawFee(who primitives.AccountId, call primitives.Call, info *primitives.DispatchInfo, fee primitives.Balance, tip primitives.Balance) (sc.Option[primitives.Balance], error)
}
