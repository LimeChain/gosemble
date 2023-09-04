package aura

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type consts struct {
	DbWeight      primitives.RuntimeDbWeight
	MinimumPeriod sc.U64
}

func newConstants(dbWeight primitives.RuntimeDbWeight, minimumPeriod sc.U64) *consts {
	return &consts{
		DbWeight:      dbWeight,
		MinimumPeriod: minimumPeriod,
	}
}
