package system

import (
	"github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	BlockHashCount types.BlockHashCount
	BlockWeights   types.BlockWeights
	BlockLength    types.BlockLength
	DbWeight       types.RuntimeDbWeight
	Version        types.RuntimeVersion
}

func NewConfig(blockHashCount types.BlockHashCount, blockWeights types.BlockWeights, blockLength types.BlockLength, dbWeight types.RuntimeDbWeight, version types.RuntimeVersion) *Config {
	return &Config{
		blockHashCount,
		blockWeights,
		blockLength,
		dbWeight,
		version,
	}
}
