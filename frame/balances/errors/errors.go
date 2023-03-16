package errors

import sc "github.com/LimeChain/goscale"

const (
	ErrorVestingBalance sc.U8 = iota
	ErrorLiquidityRestrictions
	ErrorInsufficientBalance
	ErrorExistentialDeposit
	ErrorKeepAlive
	ErrorExistingVestingSchedule
	ErrorDeadAccount
	ErrorTooManyReserves
)
