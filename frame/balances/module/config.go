package module

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	MaxLocks           sc.U32
	MaxReserves        sc.U32
	ExistentialDeposit *big.Int
	StoredMap          primitives.StoredMap
}

func NewConfig(maxLocks sc.U32, maxReserves sc.U32, existentialDeposit *big.Int, storedMap primitives.StoredMap) *Config {
	return &Config{
		MaxLocks:           maxLocks,
		MaxReserves:        maxReserves,
		ExistentialDeposit: existentialDeposit,
		StoredMap:          storedMap,
	}
}
