package constants

import "math/big"

const (
	MilliCents uint64 = 1_000_000_000
	Cents             = 1_000 * MilliCents
	Dollar            = 100 * Cents
)

var (
	Zero = big.NewInt(0)
)
