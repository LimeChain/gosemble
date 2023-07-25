package module

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
)

type Config struct {
	MaxLocks           sc.U32
	MaxReserves        sc.U32
	ExistentialDeposit *big.Int
}

func NewConfig(maxLocks sc.U32, maxReserves sc.U32, existentialDeposit *big.Int) *Config {
	return &Config{
		MaxLocks:           maxLocks,
		MaxReserves:        maxReserves,
		ExistentialDeposit: existentialDeposit,
	}
}
