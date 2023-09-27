package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	BlockHashCount sc.U64
	BlockWeights   BlockWeights
	BlockLength    BlockLength
	DbWeight       types.RuntimeDbWeight
	Version        types.RuntimeVersion
}

func NewConfig(blockHashCount sc.U64, blockWeights BlockWeights, blockLength BlockLength, dbWeight types.RuntimeDbWeight, version types.RuntimeVersion) *Config {
	return &Config{
		blockHashCount,
		blockWeights,
		blockLength,
		dbWeight,
		version,
	}
}
