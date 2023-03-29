package types

import sc "github.com/LimeChain/goscale"

type IdentityFee struct {
}

func (i IdentityFee) WeightToFee(weight Weight) Balance {
	return sc.NewU128FromUint64(uint64(weight.RefTime))
}
