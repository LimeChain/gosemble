package balances

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	DbWeight           primitives.RuntimeDbWeight
	MaxLocks           sc.U32
	MaxReserves        sc.U32
	ExistentialDeposit *big.Int
	StoredMap          primitives.StoredMap
}

func NewConfig(dbWeight primitives.RuntimeDbWeight, maxLocks sc.U32, maxReserves sc.U32, existentialDeposit *big.Int, storedMap primitives.StoredMap) *Config {
	return &Config{
		DbWeight:           dbWeight,
		MaxLocks:           maxLocks,
		MaxReserves:        maxReserves,
		ExistentialDeposit: existentialDeposit,
		StoredMap:          storedMap,
	}
}
