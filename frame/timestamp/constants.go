package timestamp

import sc "github.com/LimeChain/goscale"

type consts struct {
	MinimumPeriod sc.U64
}

func newConstants(minimumPeriod sc.U64) *consts {
	return &consts{
		minimumPeriod,
	}
}
