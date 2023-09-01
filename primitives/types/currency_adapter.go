package types

import sc "github.com/LimeChain/goscale"

// CurrencyAdapter provides an abstraction over accounts balances manipulation.
type CurrencyAdapter interface {
	// DepositIntoExisting adds free balance to `who`.
	// Returns an error if `who` is a new account.
	// Deposits an event and returns `value`.
	DepositIntoExisting(who Address32, value sc.U128) (Balance, DispatchError)
	// Withdraw removes free balance from `who` based on `reasons`.
	// If `liveness` is ExistenceRequirementKeepAlive, the remaining value must not be less than the existential deposit.
	// Checks `who` for liquidity restrictions and returns an error if they are not met.
	// Deposits a withdrawal event and returns `value`.
	Withdraw(who Address32, value sc.U128, reasons sc.U8, liveness ExistenceRequirement) (Balance, DispatchError)
}
