package balances

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type consts struct {
	DbWeight           primitives.RuntimeDbWeight
	MaxLocks           sc.U32
	MaxReserves        sc.U32
	ExistentialDeposit *big.Int
}

func newConstants(dbWeight primitives.RuntimeDbWeight, maxLocks sc.U32, maxReserves sc.U32, existentialDeposit *big.Int) *consts {
	return &consts{
		DbWeight:           dbWeight,
		MaxLocks:           maxLocks,
		MaxReserves:        maxReserves,
		ExistentialDeposit: existentialDeposit,
	}
}
