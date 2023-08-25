package balances

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	MaxLocks           sc.U32
	MaxReserves        sc.U32
	ExistentialDeposit sc.U128
	StoredMap          primitives.StoredMap
}

func NewConfig(maxLocks sc.U32, maxReserves sc.U32, existentialDeposit sc.U128, storedMap primitives.StoredMap) *Config {
	return &Config{
		MaxLocks:           maxLocks,
		MaxReserves:        maxReserves,
		ExistentialDeposit: existentialDeposit,
		StoredMap:          storedMap,
	}
}
