package balances

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
)

type consts struct {
	MaxLocks           sc.U32
	MaxReserves        sc.U32
	ExistentialDeposit *big.Int
}

func newConstants(maxLocks sc.U32, maxReserves sc.U32, existentialDeposit *big.Int) *consts {
	return &consts{
		MaxLocks:           maxLocks,
		MaxReserves:        maxReserves,
		ExistentialDeposit: existentialDeposit,
	}
}
