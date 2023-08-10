package types

import sc "github.com/LimeChain/goscale"

type CurrencyAdapter interface {
	DepositIntoExisting(who Address32, value sc.U128) (Balance, DispatchError)
	Withdraw(who Address32, value sc.U128, reasons sc.U8, liveness ExistenceRequirement) (Balance, DispatchError)
}
