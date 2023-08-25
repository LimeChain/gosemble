package balances

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	DbWeight           primitives.RuntimeDbWeight
	MaxLocks           sc.U32
	MaxReserves        sc.U32
	ExistentialDeposit sc.U128
	StoredMap          primitives.StoredMap
}

func NewConfig(dbWeight primitives.RuntimeDbWeight, maxLocks sc.U32, maxReserves sc.U32, existentialDeposit sc.U128, storedMap primitives.StoredMap) *Config {
	return &Config{
		DbWeight:           dbWeight,
		MaxLocks:           maxLocks,
		MaxReserves:        maxReserves,
		ExistentialDeposit: existentialDeposit,
		StoredMap:          storedMap,
	}
}
