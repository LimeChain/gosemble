package transaction_payment

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	OperationalFeeMultiplier sc.U8
	WeightToFee              types.WeightToFee
	LengthToFee              types.WeightToFee
	BlockWeights             types.BlockWeights
}

func NewConfig(operationalFeeMultiplier sc.U8, weightToFee, lengthToFee types.WeightToFee, blockWeights types.BlockWeights) *Config {
	return &Config{
		operationalFeeMultiplier,
		weightToFee,
		lengthToFee,
		blockWeights,
	}
}
