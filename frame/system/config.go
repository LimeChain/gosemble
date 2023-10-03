package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	BlockHashCount sc.U64
	BlockWeights   types.BlockWeights
	BlockLength    types.BlockLength
	DbWeight       types.RuntimeDbWeight
	Version        types.RuntimeVersion
}

func NewConfig(blockHashCount sc.U64, blockWeights types.BlockWeights, blockLength types.BlockLength, dbWeight types.RuntimeDbWeight, version types.RuntimeVersion) *Config {
	return &Config{
		blockHashCount,
		blockWeights,
		blockLength,
		dbWeight,
		version,
	}
}
