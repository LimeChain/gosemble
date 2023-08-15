package module

import sc "github.com/LimeChain/goscale"

type consts struct {
	OperationalFeeMultiplier sc.U8
}

func newConstants(operationalFeeMultiplier sc.U8) *consts {
	return &consts{
		operationalFeeMultiplier,
	}
}
