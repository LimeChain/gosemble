package balances

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type consts struct {
	DbWeight           primitives.RuntimeDbWeight
	MaxLocks           sc.U32
	MaxReserves        sc.U32
	ExistentialDeposit sc.U128
}

type metadataConstants struct {
	ExistentialDeposit primitives.ExistentialDeposit
	MaxLocks           primitives.MaxLocks
	MaxReserves        primitives.MaxReserves
}

func newConstants(dbWeight primitives.RuntimeDbWeight, maxLocks sc.U32, maxReserves sc.U32, existentialDeposit sc.U128) *consts {
	return &consts{
		DbWeight:           dbWeight,
		MaxLocks:           maxLocks,
		MaxReserves:        maxReserves,
		ExistentialDeposit: existentialDeposit,
	}
}
