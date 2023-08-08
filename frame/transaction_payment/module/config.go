package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	OperationalFeeMultiplier sc.U8
	WeightToFee              types.WeightToFee
	LengthToFee              types.WeightToFee
	BlockWeights             system.BlockWeights
}

func NewConfig(operationalFeeMultiplier sc.U8, weightToFee, lengthToFee types.WeightToFee, blockWeights system.BlockWeights) *Config {
	return &Config{
		operationalFeeMultiplier,
		weightToFee,
		lengthToFee,
		blockWeights,
	}
}
