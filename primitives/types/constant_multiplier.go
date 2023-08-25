package types

import sc "github.com/LimeChain/goscale"

type ConstantMultiplier struct {
	Multiplier Balance
}

func NewConstantMultiplier(multiplier Balance) ConstantMultiplier {
	return ConstantMultiplier{
		Multiplier: multiplier,
	}
}

func (cm ConstantMultiplier) WeightToFee(weight Weight) Balance {
	return sc.NewU128FromUint64(uint64(weight.RefTime)).Mul(cm.Multiplier).(Balance)
}
