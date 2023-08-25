package balances

import (
	sc "github.com/LimeChain/goscale"
)

type consts struct {
	MaxLocks           sc.U32
	MaxReserves        sc.U32
	ExistentialDeposit sc.U128
}

func newConstants(maxLocks sc.U32, maxReserves sc.U32, existentialDeposit sc.U128) *consts {
	return &consts{
		MaxLocks:           maxLocks,
		MaxReserves:        maxReserves,
		ExistentialDeposit: existentialDeposit,
	}
}
